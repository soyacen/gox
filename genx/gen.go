// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package genx

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/soyacen/gox/stringx"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"path"
	"sort"
	"strconv"
	"strings"
)

// A GeneratedFile represents a Go source file to be generated.
type GeneratedFile struct {
	filename         string
	goImportPath     GoImportPath
	buf              bytes.Buffer
	packageNames     map[GoImportPath]GoPackageName
	usedPackageNames map[GoPackageName]bool
	manualImports    map[GoImportPath]bool
	// ImportRewriteFunc is called with the import path of each package
	// imported by a generated file. It returns the import path to use
	// for this package.
	ImportRewriteFunc func(GoImportPath) GoImportPath
}

// NewGeneratedFile creates a new generated file with the given filename
// and import path.
func NewGeneratedFile(filename string, goImportPath GoImportPath) *GeneratedFile {
	g := &GeneratedFile{
		filename:         filename,
		goImportPath:     goImportPath,
		packageNames:     make(map[GoImportPath]GoPackageName),
		usedPackageNames: make(map[GoPackageName]bool),
		manualImports:    make(map[GoImportPath]bool),
	}

	// All predeclared identifiers in Go are already used.
	for _, s := range types.Universe.Names() {
		g.usedPackageNames[GoPackageName(s)] = true
	}
	return g
}

// P prints a line to the generated output. It converts each parameter to a
// string following the same rules as [fmt.Print]. It never inserts spaces
// between parameters.
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

// Import ensures a package is imported by the generated file.
//
// Packages referenced by [GeneratedFile.QualifiedGoIdent] are automatically imported.
// Explicitly importing a package with Import is generally only necessary
// when the import will be blank (import _ "package").
func (g *GeneratedFile) Import(importPath GoImportPath) {
	g.manualImports[importPath] = true
}

// Write implements [io.Writer].
func (g *GeneratedFile) Write(p []byte) (n int, err error) {
	return g.buf.Write(p)
}

// QualifiedGoIdent returns the string to use for a Go identifier.
//
// If the identifier is from a different Go package than the generated file,
// the returned name will be qualified (package.name) and an import statement
// for the identifier's package will be included in the file.
func (g *GeneratedFile) QualifiedGoIdent(ident GoIdent) string {
	if ident.GoImportPath == g.goImportPath {
		return ident.GoName
	}
	if packageName, ok := g.packageNames[ident.GoImportPath]; ok {
		return string(packageName) + "." + ident.GoName
	}
	packageName := cleanPackageName(path.Base(string(ident.GoImportPath)))
	for i, orig := 1, packageName; g.usedPackageNames[packageName]; i++ {
		packageName = orig + GoPackageName(strconv.Itoa(i))
	}
	g.packageNames[ident.GoImportPath] = packageName
	g.usedPackageNames[packageName] = true
	return string(packageName) + "." + ident.GoName
}

// Content returns the contents of the generated file.
func (g *GeneratedFile) Content() ([]byte, error) {
	if !strings.HasSuffix(g.filename, ".go") {
		return g.buf.Bytes(), nil
	}

	// Reformat generated code.
	original := g.buf.Bytes()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", original, parser.ParseComments)
	if err != nil {
		// Print out the bad code with line numbers.
		// This should never happen in practice, but it can while changing generated code
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		return nil, fmt.Errorf("%v: unparsable Go source: %v\n%v", g.filename, err, src.String())
	}

	// Collect a sorted list of all imports.
	var importPaths [][2]string
	rewriteImport := func(importPath string) string {
		if f := g.ImportRewriteFunc; f != nil {
			return string(f(GoImportPath(importPath)))
		}
		return importPath
	}
	for importPath := range g.packageNames {
		pkgName := string(g.packageNames[GoImportPath(importPath)])
		pkgPath := rewriteImport(string(importPath))
		importPaths = append(importPaths, [2]string{pkgName, pkgPath})
	}
	for importPath := range g.manualImports {
		if _, ok := g.packageNames[importPath]; !ok {
			pkgPath := rewriteImport(string(importPath))
			importPaths = append(importPaths, [2]string{"_", pkgPath})
		}
	}
	sort.Slice(importPaths, func(i, j int) bool {
		return importPaths[i][1] < importPaths[j][1]
	})

	// Modify the AST to include a new import block.
	if len(importPaths) > 0 {
		// Insert block after package statement or
		// possible comment attached to the end of the package statement.
		pos := file.Package
		tokFile := fset.File(file.Package)
		pkgLine := tokFile.Line(file.Package)
		for _, c := range file.Comments {
			if tokFile.Line(c.Pos()) > pkgLine {
				break
			}
			pos = c.End()
		}

		// Construct the import block.
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

// Comments is a comments string as provided by protoc.
type Comments string

// String formats the comments by inserting // to the start of each line,
// ensuring that there is a trailing newline.
// An empty comment is formatted as an empty string.
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

// A GoIdent is a Go identifier, consisting of a name and import path.
// The name is a single identifier and may not be a dot-qualified selector.
type GoIdent struct {
	GoName       string
	GoImportPath GoImportPath
}

func (id GoIdent) String() string { return fmt.Sprintf("%q.%v", id.GoImportPath, id.GoName) }

// A GoImportPath is the import path of a Go package.
// For example: "google.golang.org/protobuf/compiler/protogen"
type GoImportPath string

func (p GoImportPath) String() string { return strconv.Quote(string(p)) }

// Ident returns a GoIdent with s as the GoName and p as the GoImportPath.
func (p GoImportPath) Ident(s string) GoIdent {
	return GoIdent{GoName: s, GoImportPath: p}
}

// A GoPackageName is the name of a Go package. e.g., "protobuf".
type GoPackageName string

// cleanPackageName converts a string to a valid Go package name.
func cleanPackageName(name string) GoPackageName {
	return GoPackageName(stringx.GoSanitized(name))
}
