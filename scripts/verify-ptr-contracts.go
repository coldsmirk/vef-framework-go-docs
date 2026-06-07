package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type corpus struct {
	label   string
	content string
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishDocs := readCorpus("English small utilities docs", filepath.Join(docsRoot, "docs/utilities/small-utilities.md"))
	chineseDocs := readCorpus("Chinese small utilities docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/small-utilities.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	ptrDir := filepath.Join(sourceRoot, "ptr")
	expectedFuncs := []string{"Coalesce", "Equal", "Of", "Value", "ValueOrElse", "Zero"}

	var failures []string
	exported := exportedPackageSurface(ptrDir)
	failures = append(failures, compareNames("ptr func", exported.funcs, expectedFuncs)...)
	if len(exported.consts) > 0 {
		failures = append(failures, "ptr package should not expose consts, found: "+strings.Join(exported.consts, ", "))
	}
	if len(exported.types) > 0 {
		failures = append(failures, "ptr package should not expose types, found: "+strings.Join(exported.types, ", "))
	}
	if len(exported.vars) > 0 {
		failures = append(failures, "ptr package should not expose vars, found: "+strings.Join(exported.vars, ", "))
	}
	if len(exported.methods) > 0 {
		failures = append(failures, "ptr package should not expose methods, found: "+strings.Join(exported.methods, ", "))
	}
	if len(exported.fields) > 0 {
		failures = append(failures, "ptr package should not expose fields, found: "+strings.Join(exported.fields, ", "))
	}

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"6 exported functions",
		"no exported types",
		"no exported fields",
		"no exported methods",
		"`ptr.Of` requires `T comparable`",
		"returns `nil` for zero values",
		"pointer to a copy of `v`",
		"`ptr.Value` checks fallbacks from left to right",
		"A non-nil primary pointer wins",
		"over all fallbacks",
		"A non-nil fallback pointing at a zero value is still used",
		"When every pointer is nil, it returns `ptr.Zero[T]()`",
		"`ptr.ValueOrElse` is lazy",
		"`fn` must be non-nil",
		"`ptr.Equal` returns true when both pointers are nil",
		"false when exactly one",
		"is nil",
		"compares `*a == *b`",
		"`ptr.Coalesce` returns the exact first non-nil pointer",
		"With no",
		"arguments, or with only nil arguments, it returns nil",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"6 个 exported functions",
		"没有 exported types",
		"没有 exported fields",
		"没有 exported methods",
		"`ptr.Of` 要求 `T comparable`",
		"zero values，它返回 `nil`",
		"指向 `v`",
		"`ptr.Value` 会从左到右检查 fallbacks",
		"非 nil 的 primary pointer 优先于所有",
		"指向 zero value 的非 nil fallback 仍然会被使用",
		"返回 `ptr.Zero[T]()`",
		"`ptr.ValueOrElse` 是 lazy 的",
		"`fn` 必须非 nil",
		"`ptr.Equal` 在两个指针都为 nil 时返回 true",
		"只有一个为 nil 时返回 false",
		"比较 `*a == *b`",
		"`ptr.Coalesce` 返回准确的第一个非 nil 指针",
		"没有参数或参数全是\nnil 时返回 nil",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "ptr/ptr.go",
			terms: []string{
				"func Of[T comparable](v T) *T",
				"if v == zero",
				"return nil",
				"return new(v)",
				"func Zero[T any]() T",
				"func Value[T any](p *T, fallbacks ...*T) T",
				"if p != nil",
				"for _, fb := range fallbacks",
				"func ValueOrElse[T any](p *T, fn func() T) T",
				"return fn()",
				"func Equal[T comparable](a, b *T) bool",
				"if a == nil && b == nil",
				"if a == nil || b == nil",
				"return *a == *b",
				"func Coalesce[T any](ptrs ...*T) *T",
			},
		},
		{
			path: "ptr/ptr_test.go",
			terms: []string{
				"EmptyString",
				"NonEmptyString",
				"TestOfInt",
				"TestOfBool",
				"TestZero",
				"NilWithMultipleFallbacks",
				"EmptyStringFallback",
				"Should not call fallback function when pointer is non-nil",
				"BothNil",
				"FirstNil",
				"SamePointer",
				"ReturnsSamePointer",
			},
		},
	}
	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	failures = append(failures, runPackageTests(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("ptr contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("PTR contract docs verified: %d public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
		len(expectedFuncs), len(exported.methods), len(exported.fields), len(sourceChecks))
}

type packageSurface struct {
	consts  []string
	funcs   []string
	types   []string
	vars    []string
	methods []string
	fields  []string
}

func exportedPackageSurface(dir string) packageSurface {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse ptr package: %w", err))
	}

	consts := make(map[string]bool)
	funcs := make(map[string]bool)
	typesMap := make(map[string]bool)
	vars := make(map[string]bool)
	methods := make(map[string]bool)
	fields := make(map[string]bool)

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Recv == nil {
						if d.Name.IsExported() {
							funcs[d.Name.Name] = true
						}
						continue
					}
					if d.Name.IsExported() {
						recv := receiverBaseName(d.Recv.List[0].Type)
						if ast.IsExported(recv) {
							methods[recv+"."+d.Name.Name] = true
						}
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if s.Name.IsExported() {
								typesMap[s.Name.Name] = true
								if structType, ok := s.Type.(*ast.StructType); ok {
									for _, field := range structType.Fields.List {
										for _, name := range field.Names {
											if name.IsExported() {
												fields[s.Name.Name+"."+name.Name] = true
											}
										}
									}
								}
								if iface, ok := s.Type.(*ast.InterfaceType); ok {
									for _, field := range iface.Methods.List {
										for _, name := range field.Names {
											if name.IsExported() {
												methods[s.Name.Name+"."+name.Name] = true
											}
										}
									}
								}
							}
						case *ast.ValueSpec:
							for _, name := range s.Names {
								if !name.IsExported() {
									continue
								}
								if d.Tok == token.CONST {
									consts[name.Name] = true
								} else {
									vars[name.Name] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return packageSurface{
		consts:  sortedKeys(consts),
		funcs:   sortedKeys(funcs),
		types:   sortedKeys(typesMap),
		vars:    sortedKeys(vars),
		methods: sortedKeys(methods),
		fields:  sortedKeys(fields),
	}
}

func publicDocSurfaceTerms() []string {
	return []string{
		"`Of`", "`ptr.Of[T comparable](v T) *T`",
		"`Zero`", "`ptr.Zero[T any]() T`",
		"`Value`", "`ptr.Value[T any](p *T, fallbacks ...*T) T`",
		"`ValueOrElse`", "`ptr.ValueOrElse[T any](p *T, fn func() T) T`",
		"`Equal`", "`ptr.Equal[T comparable](a *T, b *T) bool`",
		"`Coalesce`", "`ptr.Coalesce[T any](ptrs ...*T) *T`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/ptr",
		"FUNC Coalesce : func[T any](ptrs ...*T) *T",
		"FUNC Equal : func[T comparable](a *T, b *T) bool",
		"FUNC Of : func[T comparable](v T) *T",
		"FUNC Value : func[T any](p *T, fallbacks ...*T) T",
		"FUNC ValueOrElse : func[T any](p *T, fn func() T) T",
		"FUNC Zero : func[T any]() T",
	}
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./ptr")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./ptr failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func receiverBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverBaseName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return receiverBaseName(t.X)
	case *ast.IndexListExpr:
		return receiverBaseName(t.X)
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func compareNames(label string, current, expected []string) []string {
	currentSet := set(current)
	expectedSet := set(expected)
	var failures []string
	for _, name := range current {
		if !expectedSet[name] {
			failures = append(failures, "unexpected public "+label+": "+name)
		}
	}
	for _, name := range expected {
		if !currentSet[name] {
			failures = append(failures, "missing expected public "+label+": "+name)
		}
	}

	return failures
}

func sortedKeys(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)

	return result
}

func set(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}

	return result
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func missingTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		failures = append(failures, missingTerm(c, term)...)
	}

	return failures
}

func missingTerm(c corpus, term string) []string {
	if !strings.Contains(c.content, term) {
		return []string{fmt.Sprintf("%s missing term: %s", c.label, term)}
	}

	return nil
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
