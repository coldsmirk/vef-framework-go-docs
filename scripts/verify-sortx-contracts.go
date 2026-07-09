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

	sortxDir := filepath.Join(sourceRoot, "sortx")
	expectedSurface := packageSurface{
		consts: []string{"NullsDefault", "NullsFirst", "NullsLast", "OrderAsc", "OrderDesc"},
		types:  []string{"NullsOrder", "OrderDirection", "OrderSpec"},
		vars:   []string{"ErrInvalidOrderDirection"},
		fields: []string{"OrderSpec.Column", "OrderSpec.Direction", "OrderSpec.NullsOrder"},
		methods: []string{
			"NullsOrder.String",
			"OrderDirection.MarshalJSON",
			"OrderDirection.MarshalText",
			"OrderDirection.String",
			"OrderDirection.UnmarshalJSON",
			"OrderDirection.UnmarshalText",
			"OrderSpec.IsValid",
		},
	}

	var failures []string
	exported := exportedPackageSurface(sortxDir)
	failures = append(failures, compareNames("sortx const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("sortx func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("sortx type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("sortx var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("sortx method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("sortx exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms(exported))...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms(exported))...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"`OrderDirection.String()` returns `DESC` only for `OrderDesc`",
		"every other\n  value renders as `ASC`",
		"`MarshalText()` lowercases the result of `String()`",
		"`UnmarshalText(text)` trims surrounding whitespace",
		"accepts `asc` /\n  `desc` case-insensitively",
		"invalid text returns an error wrapping `ErrInvalidOrderDirection`",
		"`MarshalJSON()` delegates to `MarshalText()` and emits a JSON string",
		"`UnmarshalJSON(data)` requires a JSON string",
		"`OrderDirection must be a JSON string` error",
		"`NullsDefault.String()` returns an empty string",
		"`NullsFirst.String()` returns `NULLS FIRST`",
		"`NullsLast.String()` returns `NULLS LAST`",
		"any other `NullsOrder` value also returns an empty string",
		"`OrderSpec.IsValid()` checks only `Column != \"\"`",
		"does not validate the\n  column as a SQL identifier",
		"validate `Direction` / `NullsOrder`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"`OrderDirection.String()` 只有在值为 `OrderDesc` 时返回 `DESC`",
		"其他值都会\n  渲染为 `ASC`",
		"`MarshalText()` 会把 `String()` 的结果转成小写",
		"`UnmarshalText(text)` 会 trim 首尾空白",
		"大小写不敏感方式接受 `asc` /\n  `desc`",
		"非法 text 会返回包装了 `ErrInvalidOrderDirection` 的错误",
		"`MarshalJSON()` 委托 `MarshalText()`",
		"`UnmarshalJSON(data)` 要求输入是 JSON string",
		"`OrderDirection must be a JSON string` 错误",
		"`NullsDefault.String()` 返回空字符串",
		"`NullsFirst.String()` 返回 `NULLS FIRST`",
		"`NullsLast.String()` 返回 `NULLS LAST`",
		"其他 `NullsOrder` 值也返回空字符串",
		"`OrderSpec.IsValid()` 只检查 `Column != \"\"`",
		"不会把 column 当 SQL\n  identifier 校验",
		"不会校验 `Direction` / `NullsOrder`",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "sortx/order.go",
			terms: []string{
				"type OrderDirection int",
				"OrderAsc OrderDirection = iota",
				"OrderDesc",
				"func (od OrderDirection) String() string",
				"case OrderDesc:",
				"return \"DESC\"",
				"return \"ASC\"",
				"func (od OrderDirection) MarshalText() ([]byte, error)",
				"return []byte(strings.ToLower(od.String())), nil",
				"func (od *OrderDirection) UnmarshalText(text []byte) error",
				"strings.ToUpper(strings.TrimSpace(string(text)))",
				"return fmt.Errorf(\"%w: %q (expected \\\"asc\\\" or \\\"desc\\\")\", ErrInvalidOrderDirection, string(text))",
				"func (od OrderDirection) MarshalJSON() ([]byte, error)",
				"return json.Marshal(string(text))",
				"func (od *OrderDirection) UnmarshalJSON(data []byte) error",
				"return fmt.Errorf(\"OrderDirection must be a JSON string: %w\", err)",
				"type NullsOrder int",
				"NullsDefault NullsOrder = iota",
				"NullsFirst",
				"NullsLast",
				"func (no NullsOrder) String() string",
				"return \"NULLS FIRST\"",
				"return \"NULLS LAST\"",
				"type OrderSpec struct",
				"Column string",
				"Direction OrderDirection",
				"NullsOrder NullsOrder",
				"func (spec OrderSpec) IsValid() bool",
				"return spec.Column != \"\"",
			},
		},
		{
			path: "sortx/errors.go",
			terms: []string{
				"var ErrInvalidOrderDirection = errors.New(\"invalid OrderDirection value\")",
			},
		},
		{
			path: "sortx/order_test.go",
			terms: []string{
				"TestOrderDirectionString",
				"TestOrderDirectionMarshalText",
				"TestOrderDirectionUnmarshalText",
				"WithLeadingSpace",
				"WithTrailingSpace",
				"InvalidValue",
				"EmptyString",
				"TestOrderDirectionMarshalJSON",
				"TestOrderDirectionUnmarshalJSON",
				"NotAString",
				"BooleanValue",
				"NullValue",
				"TestOrderDirectionJSONRoundTrip",
				"TestNullsOrderString",
				"DefaultNullsOrder",
				"TestOrderSpecIsValid",
				"InvalidWithoutColumn",
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
		panic(fmt.Errorf("sortx contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Sortx contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse sortx package: %w", err))
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
		case "OrderAsc":
			terms = append(terms, "`OrderAsc`", "`sortx.OrderDirection = 0`")
		case "OrderDesc":
			terms = append(terms, "`OrderDesc`", "`sortx.OrderDirection = 1`")
		case "NullsDefault":
			terms = append(terms, "`NullsDefault`", "`sortx.NullsOrder = 0`")
		case "NullsFirst":
			terms = append(terms, "`NullsFirst`", "`sortx.NullsOrder = 1`")
		case "NullsLast":
			terms = append(terms, "`NullsLast`", "`sortx.NullsOrder = 2`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.vars {
		switch name {
		case "ErrInvalidOrderDirection":
			terms = append(terms, "`ErrInvalidOrderDirection`", "exported `error` var")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.types {
		switch name {
		case "OrderDirection":
			terms = append(terms, "`OrderDirection`", "`type OrderDirection int`")
		case "NullsOrder":
			terms = append(terms, "`NullsOrder`", "`type NullsOrder int`")
		case "OrderSpec":
			terms = append(terms, "`OrderSpec`", "`type OrderSpec struct`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.fields {
		switch name {
		case "OrderSpec.Column":
			terms = append(terms, "`OrderSpec.Column`", "`string`")
		case "OrderSpec.Direction":
			terms = append(terms, "`OrderSpec.Direction`", "`sortx.OrderDirection`")
		case "OrderSpec.NullsOrder":
			terms = append(terms, "`OrderSpec.NullsOrder`", "`sortx.NullsOrder`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.methods {
		switch name {
		case "OrderDirection.String":
			terms = append(terms, "`OrderDirection.String`", "`func (od OrderDirection) String() string`")
		case "OrderDirection.MarshalText":
			terms = append(terms, "`OrderDirection.MarshalText`", "`func (od OrderDirection) MarshalText() ([]byte, error)`")
		case "OrderDirection.UnmarshalText":
			terms = append(terms, "`OrderDirection.UnmarshalText`", "`func (od *OrderDirection) UnmarshalText(text []byte) error`")
		case "OrderDirection.MarshalJSON":
			terms = append(terms, "`OrderDirection.MarshalJSON`", "`func (od OrderDirection) MarshalJSON() ([]byte, error)`")
		case "OrderDirection.UnmarshalJSON":
			terms = append(terms, "`OrderDirection.UnmarshalJSON`", "`func (od *OrderDirection) UnmarshalJSON(data []byte) error`")
		case "NullsOrder.String":
			terms = append(terms, "`NullsOrder.String`", "`func (no NullsOrder) String() string`")
		case "OrderSpec.IsValid":
			terms = append(terms, "`OrderSpec.IsValid`", "`func (spec OrderSpec) IsValid() bool`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}

	return terms
}

func publicIndexTerms(surface packageSurface) []string {
	terms := []string{"## github.com/coldsmirk/vef-framework-go/sortx"}
	for _, name := range surface.vars {
		switch name {
		case "ErrInvalidOrderDirection":
			terms = append(terms, "VAR ErrInvalidOrderDirection : error")
		default:
			terms = append(terms, "VAR "+name+" :")
		}
	}
	for _, name := range surface.consts {
		switch name {
		case "OrderAsc":
			terms = append(terms, "CONST OrderAsc : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection = 0")
		case "OrderDesc":
			terms = append(terms, "CONST OrderDesc : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection = 1")
		case "NullsDefault":
			terms = append(terms, "CONST NullsDefault : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder = 0")
		case "NullsFirst":
			terms = append(terms, "CONST NullsFirst : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder = 1")
		case "NullsLast":
			terms = append(terms, "CONST NullsLast : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder = 2")
		default:
			terms = append(terms, "CONST "+name+" :")
		}
	}
	for _, name := range surface.types {
		switch name {
		case "OrderDirection":
			terms = append(terms, "TYPE OrderDirection : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection")
		case "NullsOrder":
			terms = append(terms, "TYPE NullsOrder : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder")
		case "OrderSpec":
			terms = append(terms, "TYPE OrderSpec : github.com/coldsmirk/vef-framework-go/sortx.OrderSpec")
		default:
			terms = append(terms, "TYPE "+name+" :")
		}
	}
	for _, name := range surface.fields {
		switch name {
		case "OrderSpec.Column":
			terms = append(terms, "FIELD Column : string")
		case "OrderSpec.Direction":
			terms = append(terms, "FIELD Direction : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection")
		case "OrderSpec.NullsOrder":
			terms = append(terms, "FIELD NullsOrder : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder")
		default:
			terms = append(terms, "FIELD "+memberName(name)+" :")
		}
	}
	for _, name := range surface.methods {
		switch name {
		case "OrderDirection.String", "NullsOrder.String":
			terms = append(terms, "METHOD String : func() string")
		case "OrderDirection.MarshalText":
			terms = append(terms, "METHOD MarshalText : func() ([]byte, error)")
		case "OrderDirection.UnmarshalText":
			terms = append(terms, "METHOD UnmarshalText : func(text []byte) error")
		case "OrderDirection.MarshalJSON":
			terms = append(terms, "METHOD MarshalJSON : func() ([]byte, error)")
		case "OrderDirection.UnmarshalJSON":
			terms = append(terms, "METHOD UnmarshalJSON : func(data []byte) error")
		case "OrderSpec.IsValid":
			terms = append(terms, "METHOD IsValid : func() bool")
		default:
			terms = append(terms, "METHOD "+memberName(name)+" :")
		}
	}

	return terms
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./sortx")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./sortx failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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
