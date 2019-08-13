// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-options Generates Go code using a package as a graceful options.
// Given the name of a type T
// go-options will create a new self-contained Go source file implementing
//	func apply(*Pill)
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// For example, given this snippet,
//
//	package painkiller
//
//
//	type Pill struct{}
//
//
// running this command
//
//	go-options -type=Pill
//
// in the same directory will create the file pill_options.go, in package painkiller,
// containing a definition of
//
//	var _default_Pill_value = func() (val Pill) { return }()
//
//	// A PillOptions sets options.
//	type PillOptions interface {
//		apply(*Pill)
//	}
//
//	// EmptyPillOptions does not alter the configuration. It can be embedded
//	// in another structure to build custom options.
//	//
//	// This API is EXPERIMENTAL.
//	type EmptyPillOptions struct{}
//
//	func (EmptyPillOptions) apply(*Pill) {}
//
//	// PillOptionFunc wraps a function that modifies PillOptionFunc into an
//	// implementation of the PillOptions interface.
//	type PillOptionFunc func(*Number)
//
//	func (f PillOptionFunc) apply(do *Pill) {
//		f(do)
//	}
//
//	func (o *Pill) ApplyOptions(options ...PillOption) *Pill {
//		for _, opt := range options {
//			if opt == nil {
//				continue
//			}
//			opt.apply(o)
//		}
//		return o
//	}

//
// Typically this process would be run using go generate, like this:
//
//	//go:generate go-options -type=Pill
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default output file is t_string.go,
// where t is the lower-cased name of the first type listed. It can be overridden
// with the -output flag.
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	strings_ "github.com/searKing/golang/go/strings"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeInfos   = flag.String("type", "", "comma-separated list of type names; must be set")
	output      = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	trimprefix  = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	linecomment = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of go-options:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-options [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttp://godoc.org/github.com/searKing/go-options\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	goOptionsToolName = "go-options"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-options: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeInfos) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// type <key, value> type <key, value>
	types := newTypeInfo(*typeInfos)
	if len(types) == 0 {
		flag.Usage()
		os.Exit(3)
	}

	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	g := Generator{
		trimPrefix:  *trimprefix,
		lineComment: *linecomment,
	}
	// TODO(suzmue): accept other patterns for packages (directories, list of files, import paths, etc).
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	g.parsePackage(args, tags)

	// Print the header and package clause.
	g.Printf("// Code generated by \"%s %s\"; DO NOT EDIT.\n", goOptionsToolName, strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")

	// Run generate for each type.
	for _, typeInfo := range types {
		g.generate(typeInfo)
	}

	// Format the output.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_options.go", types[0].eleName)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, target, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.

	trimPrefix  string
	lineComment bool
}

// Printf format & write to the buf in this generator
func (g *Generator) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeInfo typeInfo
	values   []Value // Accumulator for constant values of that type.

	trimPrefix  string
	lineComment bool
}

// Package holds a single parsed package and associated files and ast files.
type Package struct {
	// Name is the package name as it appears in the package source code.
	name string

	// Defs maps identifiers to the objects they define (including
	// package names, dots "." of dot-imports, and blank "_" identifiers).
	// For identifiers that do not denote objects (e.g., the package name
	// in package clauses, or symbolic variables t in t := x.(type) of
	// type switch headers), the corresponding objects are nil.
	//
	// For an embedded field, Defs returns the field *Var it defines.
	//
	// Invariant: Defs[id] == nil || Defs[id].Pos() == id.Pos()
	defs map[*ast.Ident]types.Object

	// Ast files to which this package contains.
	files []*File
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:        file,
			pkg:         g.pkg,
			trimPrefix:  g.trimPrefix,
			lineComment: g.lineComment,
		}
	}
}

// generate produces the String method for the named type.
func (g *Generator) generate(typeInfo typeInfo) {
	// <key, value>
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeInfo = typeInfo
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %+v", typeInfo)
	}
	g.buildOneRun(values[0])
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

func (g *Generator) goimport(src []byte) []byte {
	var opt = &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}

	res, err := imports.Process("", src, opt)
	if err != nil {
		log.Fatalf("process import: %s", err)
	}

	return res
}

// Value represents a declared constant.
type Value struct {
	originalName string // The name of the constant.
	name         string // The name with trimmed prefix.
	str          string // The string representation given by the "go/constant" package.

	eleImport string // import path of the atomic.Value type.
	eleName   string // Name of the atomic.Value type.
}

func (v *Value) String() string {
	return v.str
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	// Token must be in IMPORT, CONST, TYPE, VAR
	if !ok || decl.Tok != token.TYPE {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		tspec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is TYPE.
		typ = tspec.Name.Name

		if typ != f.typeInfo.eleName {
			// This is not the type we're looking for.
			continue
		}
		v := Value{
			originalName: typ,
			str:          typ,
		}
		if c := tspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
			v.name = strings.TrimSpace(c.Text())
		} else {
			v.name = strings.TrimPrefix(v.originalName, f.trimPrefix)
		}
		v.eleName = v.name
		f.values = append(f.values, v)
	}
	return false
}

// Helpers

// declareNameVar declares the concatenated names
// strings representing the runs of values
func (g *Generator) declareNameVar(run Value) string {
	defaultValName, defaultValDecl := g.createValAndNameDecl(run)
	g.Printf("var %s\n", defaultValDecl)
	return defaultValName
}

// createValAndNameDecl returns the pair of declarations for the run. The caller will add "var".
func (g *Generator) createValAndNameDecl(val Value) (string, string) {
	defaultValName := fmt.Sprintf("_default_%s_value", val.eleName)
	defaultValDecl := fmt.Sprintf("%s = func() (val %s) { return }()", defaultValName, val.eleName)

	return defaultValName, defaultValDecl
}

// buildOneRun generates the variables and String method for a single run of contiguous values.
func (g *Generator) buildOneRun(value Value) {
	//values := run
	g.Printf("\n")
	if strings.TrimSpace(value.eleImport) != "" {
		g.Printf(stringImport, value.eleImport)
	}

	g.declareNameVar(value)

	//The generated code is simple enough to write as a Printf format.
	optionInterfaceName := strings_.CamelCaseSlice(value.eleName, "option")
	g.Printf(stringOneRun, value.eleName, optionInterfaceName)
	if strings.TrimSpace(value.eleImport) != "" {
		g.Printf("\n")
		g.Printf(stringApplyOptionsAsCFunction, value.eleName, optionInterfaceName)
	} else {
		g.Printf("\n")
		g.Printf(stringApplyOptionsAsMemberFunction, value.eleName, optionInterfaceName)
	}
}

// Arguments to format are:
//	[1]: import path
const stringImport = `import "%s"
`

// Arguments to format are:
//	[1]: option type name
//	[2]: optionInterface type name
const stringOneRun = `// A %[2]s sets options.
type %[2]s interface {
	apply(*%[1]s)
}

// Empty%[2]s does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type Empty%[2]s struct{}

func (Empty%[2]s) apply(*%[1]s) {}

// %[2]sFunc wraps a function that modifies %[1]s into an
// implementation of the %[2]s interface.
type %[2]sFunc func(*%[1]s)

func (f %[2]sFunc) apply(do *%[1]s) {
	f(do)
}
`

// Arguments to format are:
//	[1]: option type name
//	[2]: optionInterface type name
const stringApplyOptionsAsMemberFunction = `func (o *%[1]s) ApplyOptions(options ...%[2]s) *%[1]s {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
`

// Arguments to format are:
//	[1]: option type name
//	[2]: optionInterface type name
const stringApplyOptionsAsCFunction = `func ApplyOptions(o *%[1]s, options ...%[2]s) *%[1]s {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
`
