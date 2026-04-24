package gen

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

// CodeWriter represents a Go source file being generated.
type CodeWriter struct {
	filename          string                          // 文件名
	goImportPath      GoImportPath                    // 当前文件的导入路径
	buf               bytes.Buffer                    // 存储生成的代码内容
	packageNames      map[GoImportPath]GoPackageName  // 记录导入路径到包名的映射
	usedPackageNames  map[GoPackageName]bool          // 记录已使用的包名，避免冲突
	manualImports     map[GoImportPath]bool           // 手动导入的包路径
	ImportRewriteFunc func(GoImportPath) GoImportPath // 导入路径重写函数
}

// NewCodeWriter creates a new CodeWriter instance for generating Go source files.
//
// Parameters:
//   - filename: the name of the file being generated
//   - goImportPath: the import path of the current file
//
// Returns:
//   - *CodeWriter: a new CodeWriter instance
func NewCodeWriter(filename string, goImportPath GoImportPath) *CodeWriter {
	g := &CodeWriter{
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

// P prints arguments to the buffer, appending a newline after each call.
// Special handling for GoIdent types, converting them to qualified names.
//
// Parameters:
//   - v: values to print
func (g *CodeWriter) P(v ...any) {
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

// Import adds a manual import package path.
//
// Parameters:
//   - importPath: the import path to add
func (g *CodeWriter) Import(importPath GoImportPath) {
	g.manualImports[importPath] = true
}

// Write implements the io.Writer interface, writing bytes to the internal buffer.
//
// Parameters:
//   - p: bytes to write
//
// Returns:
//   - int: number of bytes written
//   - error: any write error
func (g *CodeWriter) Write(p []byte) (n int, err error) {
	return g.buf.Write(p)
}

// QualifiedGoIdent returns the qualified name (package.name) for the given identifier.
// If the identifier is in the same package, only the identifier name is returned.
//
// Parameters:
//   - ident: the Go identifier to qualify
//
// Returns:
//   - string: the qualified name
func (g *CodeWriter) QualifiedGoIdent(ident GoIdent) string {
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

// Content returns the generated file content, automatically processing import declarations if the file is Go source.
//
// Returns:
//   - []byte: the generated file content
//   - error: any parsing or formatting error
func (g *CodeWriter) Content() ([]byte, error) {
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

// Comments represents comments in Go code.
type Comments string

// String converts comments to Go comment format.
//
// Returns:
//   - string: the formatted comments
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

// GoIdent represents a Go identifier, including its name and import path.
type GoIdent struct {
	GoName       string       // 标识符名称
	GoImportPath GoImportPath // 所属包的导入路径
}

// String returns the string representation of the identifier.
//
// Returns:
//   - string: formatted as "importPath".name
func (id GoIdent) String() string { return fmt.Sprintf("%q.%v", id.GoImportPath, id.GoName) }

// GoImportPath represents a Go import path.
type GoImportPath string

// String returns the string representation of the import path.
//
// Returns:
//   - string: the quoted import path
func (p GoImportPath) String() string { return strconv.Quote(string(p)) }

// Ident creates an identifier under this import path with the given name.
//
// Parameters:
//   - s: the identifier name
//
// Returns:
//   - GoIdent: the created identifier
func (p GoImportPath) Ident(s string) GoIdent {
	return GoIdent{GoName: s, GoImportPath: p}
}

// GoPackageName represents a Go package name.
type GoPackageName string

// cleanPackageName 清理包名，确保其符合Go命名规则
func cleanPackageName(name string) GoPackageName {
	return GoPackageName(stringx.GoSanitized(name))
}
