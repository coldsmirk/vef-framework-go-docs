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

	englishDocs := readCorpus("English small utilities docs", filepath.Join(docsRoot, "docs/utilities/small-helpers.md"))
	chineseDocs := readCorpus("Chinese small utilities docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/small-helpers.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	monadDir := filepath.Join(sourceRoot, "monad")
	expectedSurface := packageSurface{
		funcs:  []string{"NewRange"},
		types:  []string{"Range"},
		fields: []string{"Range.End", "Range.Start"},
		methods: []string{
			"Range.Contains",
			"Range.Intersection",
			"Range.IsEmpty",
			"Range.IsNotEmpty",
			"Range.IsValid",
			"Range.Overlaps",
		},
	}

	var failures []string
	exported := exportedPackageSurface(monadDir)
	failures = append(failures, compareNames("monad const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("monad func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("monad type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("monad var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("monad method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("monad exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms(exported))...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms(exported))...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"`Range[T]` is constrained to `cmp.Ordered`",
		"ranges are inclusive at both ends: `[Start, End]`",
		"exported fields have no JSON tags",
		"default Go JSON encoding uses `Start`\n  and `End`",
		"`NewRange(start, end)` stores the two bounds exactly as provided and does not\n  reorder them",
		"`IsValid()` and `IsNotEmpty()` return `Start <= End`",
		"`IsEmpty()` returns `Start > End`",
		"`Contains(value)` checks `Start <= value && value <= End`",
		"`Overlaps(other)` checks `Start <= other.End && other.Start <= End`",
		"adjacent\n  ranges that share one endpoint overlap",
		"`Intersection(other)` first calls `Overlaps(other)`",
		"`Intersection` returns `Range[T]{}, false`",
		"returned zero range carries no business meaning",
		"`max(Start, other.Start)` to `min(End, other.End)`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"`Range[T]` 约束为 `cmp.Ordered`",
		"区间两端都是闭区间：`[Start, End]`",
		"exported fields 没有 JSON tags",
		"默认 JSON 编码会使用 `Start` 和 `End`",
		"`NewRange(start, end)` 会按原样保存两个边界，不会重排",
		"`IsValid()` 和 `IsNotEmpty()` 返回 `Start <= End`",
		"`IsEmpty()` 返回 `Start > End`",
		"`Contains(value)` 检查 `Start <= value && value <= End`",
		"`Overlaps(other)` 检查 `Start <= other.End && other.Start <= End`",
		"共享一个\n  端点的相邻区间也算重叠",
		"`Intersection(other)` 会先调用 `Overlaps(other)`",
		"`Intersection` 返回 `Range[T]{}, false`",
		"返回的 zero range\n  没有业务意义",
		"`max(Start, other.Start)` 到\n  `min(End, other.End)`",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "monad/range.go",
			terms: []string{
				"import \"cmp\"",
				"type Range[T cmp.Ordered] struct",
				"Start T",
				"End   T",
				"func NewRange[T cmp.Ordered](start, end T) Range[T]",
				"Start: start",
				"End:   end",
				"func (r Range[T]) Contains(value T) bool",
				"return r.Start <= value && value <= r.End",
				"func (r Range[T]) IsValid() bool",
				"return r.Start <= r.End",
				"func (r Range[T]) IsEmpty() bool",
				"return r.Start > r.End",
				"func (r Range[T]) IsNotEmpty() bool",
				"func (r Range[T]) Overlaps(other Range[T]) bool",
				"return r.Start <= other.End && other.Start <= r.End",
				"func (r Range[T]) Intersection(other Range[T]) (Range[T], bool)",
				"if !r.Overlaps(other)",
				"return Range[T]{}, false",
				"Start: max(r.Start, other.Start)",
				"End:   min(r.End, other.End)",
			},
		},
		{
			path: "monad/range_test.go",
			terms: []string{
				"TestNewRange",
				"StringRange",
				"FloatRange",
				"ValueAtStart",
				"ValueAtEnd",
				"SingleValueRangeContains",
				"InvalidRange",
				"TestRangeIsEmpty",
				"AdjacentRanges",
				"SinglePointOverlap",
				"TestRangeIntersection",
				"Range[int]{}",
				"TestRangeIntersectionString",
				"TestRangeJSONMarshaling",
				"TestRangeWithDifferentTypes",
				"TestRangeEdgeCases",
				"NegativeRangeIntersection",
				"TestRangeStringOperations",
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
		panic(fmt.Errorf("monad contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Monad contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
		topLevelPublic, len(exported.methods), len(exported.fields), len(sourceChecks))
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
		panic(fmt.Errorf("failed to parse monad package: %w", err))
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

func publicDocSurfaceTerms(surface packageSurface) []string {
	var terms []string
	for _, name := range surface.types {
		switch name {
		case "Range":
			terms = append(terms, "`Range[T]`", "`type Range[T cmp.Ordered] struct`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.funcs {
		switch name {
		case "NewRange":
			terms = append(terms, "`NewRange`", "`monad.NewRange[T cmp.Ordered](start T, end T) monad.Range[T]`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.fields {
		switch name {
		case "Range.Start":
			terms = append(terms, "`Range.Start`", "`T`")
		case "Range.End":
			terms = append(terms, "`Range.End`", "`T`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.methods {
		switch name {
		case "Range.Contains":
			terms = append(terms, "`Range.Contains`", "`func (r Range[T]) Contains(value T) bool`")
		case "Range.IsValid":
			terms = append(terms, "`Range.IsValid`", "`func (r Range[T]) IsValid() bool`")
		case "Range.IsEmpty":
			terms = append(terms, "`Range.IsEmpty`", "`func (r Range[T]) IsEmpty() bool`")
		case "Range.IsNotEmpty":
			terms = append(terms, "`Range.IsNotEmpty`", "`func (r Range[T]) IsNotEmpty() bool`")
		case "Range.Overlaps":
			terms = append(terms, "`Range.Overlaps`", "`func (r Range[T]) Overlaps(other monad.Range[T]) bool`")
		case "Range.Intersection":
			terms = append(terms, "`Range.Intersection`", "`func (r Range[T]) Intersection(other monad.Range[T]) (monad.Range[T], bool)`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}

	return terms
}

func publicIndexTerms(surface packageSurface) []string {
	terms := []string{"## github.com/coldsmirk/vef-framework-go/monad"}
	for _, name := range surface.funcs {
		switch name {
		case "NewRange":
			terms = append(terms, "FUNC NewRange : func[T cmp.Ordered](start T, end T) github.com/coldsmirk/vef-framework-go/monad.Range[T]")
		default:
			terms = append(terms, "FUNC "+name+" :")
		}
	}
	for _, name := range surface.types {
		switch name {
		case "Range":
			terms = append(terms, "TYPE Range : github.com/coldsmirk/vef-framework-go/monad.Range[T cmp.Ordered]")
		default:
			terms = append(terms, "TYPE "+name+" :")
		}
	}
	for _, name := range surface.fields {
		switch name {
		case "Range.Start":
			terms = append(terms, "FIELD Start : T")
		case "Range.End":
			terms = append(terms, "FIELD End : T")
		default:
			terms = append(terms, "FIELD "+memberName(name)+" :")
		}
	}
	for _, name := range surface.methods {
		switch name {
		case "Range.Contains":
			terms = append(terms, "METHOD Contains : func(value T) bool")
		case "Range.Intersection":
			terms = append(terms, "METHOD Intersection : func(other github.com/coldsmirk/vef-framework-go/monad.Range[T]) (github.com/coldsmirk/vef-framework-go/monad.Range[T], bool)")
		case "Range.IsEmpty":
			terms = append(terms, "METHOD IsEmpty : func() bool")
		case "Range.IsNotEmpty":
			terms = append(terms, "METHOD IsNotEmpty : func() bool")
		case "Range.IsValid":
			terms = append(terms, "METHOD IsValid : func() bool")
		case "Range.Overlaps":
			terms = append(terms, "METHOD Overlaps : func(other github.com/coldsmirk/vef-framework-go/monad.Range[T]) bool")
		default:
			terms = append(terms, "METHOD "+memberName(name)+" :")
		}
	}

	return terms
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./monad")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./monad failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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

func memberName(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[i+1:]
	}

	return name
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
