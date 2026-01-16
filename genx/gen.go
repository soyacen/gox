package genx

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/soyacen/gox/stringx"
)

// GeneratedFile 表示一个正在生成的Go源文件
type GeneratedFile struct {
	filename          string                           // 文件名
	goImportPath      GoImportPath                   // 当前文件的导入路径
	buf               bytes.Buffer                   // 存储生成的代码内容
	packageNames      map[GoImportPath]GoPackageName // 记录导入路径到包名的映射
	usedPackageNames  map[GoPackageName]bool         // 记录已使用的包名，避免冲突
	manualImports     map[GoImportPath]bool          // 手动导入的包路径
	ImportRewriteFunc func(GoImportPath) GoImportPath // 导入路径重写函数
}

// NewGeneratedFile 创建一个新的生成文件实例
func NewGeneratedFile(filename string, goImportPath GoImportPath) *GeneratedFile {
	g := &GeneratedFile{
		filename:         filename,
		goImportPath:     goImportPath,
		packageNames:     make(map[GoImportPath]GoPackageName),
		usedPackageNames: make(map[GoPackageName]bool),
		manualImports:    make(map[GoImportPath]bool),
	}

	// 将Go预声明标识符标记为已使用
	for _, s := range types.Universe.Names() {
		g.usedPackageNames[GoPackageName(s)] = true
	}
	return g
}

// P 将参数打印到缓冲区，每行结束添加换行符
// 特殊处理 GoIdent 类型，将其转换为适当的限定名称
func (g *GeneratedFile) P(v ...any) {
	for _, x := range v {
		switch x := x.(type) {
		case GoIdent:
			fmt.Fprint(&g.buf, g.QualifiedGoIdent(x))
		default:
			fmt.Fprint(&g.buf, x)
		}
	}
	fmt.Fprintln(&g.buf)
}

// Import 添加一个手动导入的包路径
func (g *GeneratedFile) Import(importPath GoImportPath) {
	g.manualImports[importPath] = true
}

// Write 实现 io.Writer 接口，将字节写入内部缓冲区
func (g *GeneratedFile) Write(p []byte) (n int, err error) {
	return g.buf.Write(p)
}

// QualifiedGoIdent 返回给定标识符的限定名称（包名.标识符）
// 如果标识符在同一包中，则只返回标识符名称
func (g *GeneratedFile) QualifiedGoIdent(ident GoIdent) string {
	// 如果标识符属于当前包，则直接返回名称
	if ident.GoImportPath == g.goImportPath {
		return ident.GoName
	}
	
	// 检查是否已有对应的包名
	if packageName, ok := g.packageNames[ident.GoImportPath]; ok {
		return string(packageName) + "." + ident.GoName
	}
	
	// 根据导入路径的基名生成包名
	packageName := cleanPackageName(path.Base(string(ident.GoImportPath)))
	
	// 如果包名已被使用，则添加数字后缀直到找到唯一名称
	originalName := packageName
	for i := 1; g.usedPackageNames[packageName]; i++ {
		packageName = originalName + GoPackageName(strconv.Itoa(i))
	}
	
	g.packageNames[ident.GoImportPath] = packageName
	g.usedPackageNames[packageName] = true
	return string(packageName) + "." + ident.GoName
}

// Content 获取生成的文件内容，如果文件是Go源码则自动处理导入声明
func (g *GeneratedFile) Content() ([]byte, error) {
	if !strings.HasSuffix(g.filename, ".go") {
		return g.buf.Bytes(), nil
	}

	original := g.buf.Bytes()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", original, parser.ParseComments)
	if err != nil {
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		return nil, fmt.Errorf("%v: unparsable Go source: %v\n%v", g.filename, err, src.String())
	}

	// 准备导入路径列表
	var importPaths [][2]string
	rewriteImport := func(importPath string) string {
		if f := g.ImportRewriteFunc; f != nil {
			return string(f(GoImportPath(importPath)))
		}
		return importPath
	}
	
	// 添加由标识符引用产生的导入
	for importPath, packageName := range g.packageNames {
		pkgName := string(packageName)
		pkgPath := rewriteImport(string(importPath))
		importPaths = append(importPaths, [2]string{pkgName, pkgPath})
	}
	
	// 添加手动导入的包（如果没有被引用过，使用 _ 作为空白导入）
	for importPath := range g.manualImports {
		if _, ok := g.packageNames[importPath]; !ok {
			pkgPath := rewriteImport(string(importPath))
			importPaths = append(importPaths, [2]string{"_", pkgPath})
		}
	}
	
	// 按路径排序导入
	sort.Slice(importPaths, func(i, j int) bool {
		return importPaths[i][1] < importPaths[j][1]
	})

	if len(importPaths) > 0 {
		pos := file.Package
		tokFile := fset.File(file.Package)
		pkgLine := tokFile.Line(file.Package)
		for _, c := range file.Comments {
			if tokFile.Line(c.Pos()) > pkgLine {
				break
			}
			pos = c.End()
		}

		impDecl := &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: pos,
			Lparen: pos,
			Rparen: pos,
		}
		for _, importPath := range importPaths {
			impDecl.Specs = append(impDecl.Specs, &ast.ImportSpec{
				Name: &ast.Ident{
					Name:    importPath[0],
					NamePos: pos,
				},
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    strconv.Quote(importPath[1]),
					ValuePos: pos,
				},
				EndPos: pos,
			})
		}
		file.Decls = append([]ast.Decl{impDecl}, file.Decls...)
	}

	var out bytes.Buffer
	if err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(&out, fset, file); err != nil {
		return nil, fmt.Errorf("%v: can not reformat Go source: %v", g.filename, err)
	}
	return out.Bytes(), nil
}

// Comments 表示Go代码中的注释
type Comments string

// String 将注释转换为Go代码中的注释格式
func (c Comments) String() string {
	if c == "" {
		return ""
	}
	var b []byte
	for _, line := range strings.Split(strings.TrimSuffix(string(c), "\n"), "\n") {
		b = append(b, "//"...)
		b = append(b, line...)
		b = append(b, "\n"...)
	}
	return string(b)
}

// GoIdent 表示Go代码中的标识符，包括名称和所属的导入路径
type GoIdent struct {
	GoName       string       // 标识符名称
	GoImportPath GoImportPath // 所属包的导入路径
}

// String 返回标识符的字符串表示
func (id GoIdent) String() string { return fmt.Sprintf("%q.%v", id.GoImportPath, id.GoName) }

// GoImportPath 表示Go代码中的导入路径
type GoImportPath string

// String 返回导入路径的字符串表示
func (p GoImportPath) String() string { return strconv.Quote(string(p)) }

// Ident 根据名称创建一个在该导入路径下的标识符
func (p GoImportPath) Ident(s string) GoIdent {
	return GoIdent{GoName: s, GoImportPath: p}
}

// GoPackageName 表示Go代码中的包名
type GoPackageName string

// cleanPackageName 清理包名，确保其符合Go命名规则
func cleanPackageName(name string) GoPackageName {
	return GoPackageName(stringx.GoSanitized(name))
}