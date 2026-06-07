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

	pageDir := filepath.Join(sourceRoot, "page")
	expectedSurface := packageSurface{
		consts: []string{"DefaultPageNumber", "DefaultPageSize", "MaxPageSize"},
		funcs:  []string{"New"},
		types:  []string{"Page", "Pageable"},
		fields: []string{"Page.Items", "Page.Page", "Page.Size", "Page.Total", "Pageable.Page", "Pageable.Size"},
		methods: []string{
			"Page.HasNext",
			"Page.HasPrevious",
			"Page.TotalPages",
			"Pageable.Normalize",
			"Pageable.Offset",
		},
	}

	var failures []string
	exported := exportedPackageSurface(pageDir)
	failures = append(failures, compareNames("page const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("page func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("page type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("page var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("page method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("page exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms(exported))...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms(exported))...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"6 top-level exported symbols",
		"6 exported fields",
		"5 exported\nmethods",
		"no exported variables",
		"`Pageable.Page` is 1-based",
		"`Normalize(size...)` mutates the receiver in place",
		"`Normalize` resets `Page < 1` to `DefaultPageNumber`",
		"uses only the first optional fallback size",
		"`DefaultPageSize` when no fallback is provided",
		"above `MaxPageSize` are clamped to `MaxPageSize`",
		"does not re-validate a negative custom fallback",
		"`Offset()` is a plain `(Page - 1) * Size` calculation",
		"`New` copies `pageable.Page`, `pageable.Size`, and `total`",
		"converts nil `items` to a non-nil empty slice",
		"non-nil `items` are used\n  as provided and are not cloned",
		"`TotalPages()` returns `0` when `Size == 0`",
		"ceiling\n  division of `Total / Size`",
		"`HasNext()` compares `Page < TotalPages()`",
		"`HasPrevious()` returns true when `Page > 1`",
		"does not validate negative totals",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"6 个 top-level exported symbols",
		"6 个 exported fields",
		"5 个\nexported methods",
		"没有 exported variables",
		"`Pageable.Page` 是 1-based 页码",
		"`Normalize(size...)` 会原地修改 receiver",
		"`Normalize` 会把 `Page < 1` 重置为 `DefaultPageNumber`",
		"只使用第一个可选 fallback size",
		"fallback 时使用 `DefaultPageSize`",
		"大于 `MaxPageSize` 的值会被截到 `MaxPageSize`",
		"不会对负数自定义 fallback 再做一次校正",
		"`Offset()` 只是 `(Page - 1) * Size` 计算",
		"`New` 会把 `pageable.Page`、`pageable.Size` 和 `total` 复制进响应",
		"nil `items` 转成非 nil 空 slice",
		"非 nil `items` 按原样使用",
		"`TotalPages()` 在 `Size == 0` 时返回 `0`",
		"按 `Total / Size` 向上取整",
		"`HasNext()` 判断 `Page < TotalPages()`",
		"`HasPrevious()` 在 `Page > 1` 时返回 true",
		"不会校验负数 total",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "page/page.go",
			terms: []string{
				"DefaultPageNumber int = 1",
				"DefaultPageSize int = 15",
				"MaxPageSize int = 1000",
				"type Pageable struct",
				"Page int `json:\"page\"`",
				"Size int `json:\"size\"`",
				"func (p *Pageable) Normalize(size ...int)",
				"if p.Page < 1",
				"p.Page = DefaultPageNumber",
				"if p.Size < 1",
				"p.Size = size[0]",
				"p.Size = DefaultPageSize",
				"if p.Size > MaxPageSize",
				"p.Size = MaxPageSize",
				"func (p Pageable) Offset() int",
				"return (p.Page - 1) * p.Size",
				"type Page[T any] struct",
				"Total int64 `json:\"total\"`",
				"Items []T   `json:\"items\"`",
				"func (page Page[T]) TotalPages() int",
				"if page.Size == 0",
				"func (page Page[T]) HasNext() bool",
				"return page.Page < page.TotalPages()",
				"func (page Page[T]) HasPrevious() bool",
				"return page.Page > 1",
				"func New[T any](pageable Pageable, total int64, items []T) Page[T]",
				"if items == nil",
				"items = []T{}",
			},
		},
		{
			path: "page/page_test.go",
			terms: []string{
				"PageLessThanOne",
				"SizeLessThanOne",
				"SizeExceedsMaximum",
				"TestPageableOffset",
				"TestNewPageWithNilItems",
				"ZeroSize",
				"TestPageHasNext",
				"TestPageHasPrevious",
				"TestPageableJSONMarshaling",
				"TestPageJSONMarshaling",
				"WithStructType",
				"ApiPaginationWorkflow",
				"EmptyResultSet",
			},
		},
		{
			path: "crud/find_page.go",
			terms: []string{
				"func (a *findPageOperation[TModel, TSearch]) findPage(db orm.DB) (func(ctx fiber.Ctx, db orm.DB, transformer mold.Transformer, pageable page.Pageable, search TSearch, meta api.Meta) error, error)",
				"pageable.Normalize(a.defaultPageSize)",
				"query  = db.NewSelect().Model(&models).SelectModelColumns().Paginate(pageable)",
				"return result.Ok(page.New(pageable, total, []any{})).Response(ctx)",
				"return result.Ok(page.New(pageable, total, typedModels)).Response(ctx)",
				"return result.Ok(page.New(pageable, total, items)).Response(ctx)",
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
		panic(fmt.Errorf("page contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Page contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse page package: %w", err))
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
	for _, name := range surface.consts {
		switch name {
		case "DefaultPageNumber":
			terms = append(terms, "`DefaultPageNumber`", "`int = 1`")
		case "DefaultPageSize":
			terms = append(terms, "`DefaultPageSize`", "`int = 15`")
		case "MaxPageSize":
			terms = append(terms, "`MaxPageSize`", "`int = 1000`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.types {
		switch name {
		case "Page":
			terms = append(terms, "`Page[T]`", "`type Page[T any] struct { Page int json:\"page\"; Size int json:\"size\"; Total int64 json:\"total\"; Items []T json:\"items\" }`")
		case "Pageable":
			terms = append(terms, "`Pageable`", "`type Pageable struct { Page int json:\"page\"; Size int json:\"size\" }`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.funcs {
		switch name {
		case "New":
			terms = append(terms, "`New`", "`page.New[T any](pageable page.Pageable, total int64, items []T) page.Page[T]`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.fields {
		terms = append(terms, "`"+name+"`")
	}
	for _, name := range surface.methods {
		switch name {
		case "Page.HasNext":
			terms = append(terms, "`Page.HasNext`", "`func (page Page[T]) HasNext() bool`")
		case "Page.HasPrevious":
			terms = append(terms, "`Page.HasPrevious`", "`func (page Page[T]) HasPrevious() bool`")
		case "Page.TotalPages":
			terms = append(terms, "`Page.TotalPages`", "`func (page Page[T]) TotalPages() int`")
		case "Pageable.Normalize":
			terms = append(terms, "`Pageable.Normalize`", "`func (p *Pageable) Normalize(size ...int)`")
		case "Pageable.Offset":
			terms = append(terms, "`Pageable.Offset`", "`func (p Pageable) Offset() int`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}

	return terms
}

func publicIndexTerms(surface packageSurface) []string {
	terms := []string{"## github.com/coldsmirk/vef-framework-go/page"}
	for _, name := range surface.consts {
		switch name {
		case "DefaultPageNumber":
			terms = append(terms, "CONST DefaultPageNumber : int = 1")
		case "DefaultPageSize":
			terms = append(terms, "CONST DefaultPageSize : int = 15")
		case "MaxPageSize":
			terms = append(terms, "CONST MaxPageSize : int = 1000")
		default:
			terms = append(terms, "CONST "+name+" :")
		}
	}
	for _, name := range surface.funcs {
		switch name {
		case "New":
			terms = append(terms, "FUNC New : func[T any](pageable github.com/coldsmirk/vef-framework-go/page.Pageable, total int64, items []T) github.com/coldsmirk/vef-framework-go/page.Page[T]")
		default:
			terms = append(terms, "FUNC "+name+" :")
		}
	}
	for _, name := range surface.types {
		switch name {
		case "Page":
			terms = append(terms, "TYPE Page : github.com/coldsmirk/vef-framework-go/page.Page[T any]")
		case "Pageable":
			terms = append(terms, "TYPE Pageable : github.com/coldsmirk/vef-framework-go/page.Pageable")
		default:
			terms = append(terms, "TYPE "+name+" :")
		}
	}
	for _, name := range surface.fields {
		switch name {
		case "Page.Page":
			terms = append(terms, "FIELD Page : int")
		case "Page.Size", "Pageable.Size":
			terms = append(terms, "FIELD Size : int")
		case "Page.Total":
			terms = append(terms, "FIELD Total : int64")
		case "Page.Items":
			terms = append(terms, "FIELD Items : []T")
		case "Pageable.Page":
			terms = append(terms, "FIELD Page : int")
		default:
			terms = append(terms, "FIELD "+memberName(name)+" :")
		}
	}
	for _, name := range surface.methods {
		switch name {
		case "Page.HasNext":
			terms = append(terms, "METHOD HasNext : func() bool")
		case "Page.HasPrevious":
			terms = append(terms, "METHOD HasPrevious : func() bool")
		case "Page.TotalPages":
			terms = append(terms, "METHOD TotalPages : func() int")
		case "Pageable.Normalize":
			terms = append(terms, "METHOD Normalize : func(size ...int)")
		case "Pageable.Offset":
			terms = append(terms, "METHOD Offset : func() int")
		default:
			terms = append(terms, "METHOD "+memberName(name)+" :")
		}
	}

	return terms
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./page")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./page failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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
