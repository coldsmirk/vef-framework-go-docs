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

	englishDocs := readCorpus("English copier docs", filepath.Join(docsRoot, "docs/utilities/copier.md"))
	chineseDocs := readCorpus("Chinese copier docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/copier.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	copierDir := filepath.Join(sourceRoot, "copier")
	expectedSurface := packageSurface{
		funcs: []string{
			"Copy",
			"WithCaseInsensitive",
			"WithDeepCopy",
			"WithFieldNameMapping",
			"WithIgnoreEmpty",
			"WithTypeConverters",
		},
		types: []string{"CopyOption", "FieldNameMapping", "TypeConverter"},
	}

	var failures []string
	exported := exportedPackageSurface(copierDir)
	failures = append(failures, compareNames("copier const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("copier func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("copier type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("copier var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("copier method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("copier exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"9 top-level exported symbols",
		"no exported fields",
		"no exported\nmethods",
		"fingerprint is\n`44b6cf428fb9c642afca0cd25257c8ade57c9ac855b3ecc67cf575c1323fdf58`",
		"`copier.Copy(src, dst, options...)`",
		"`copier.CopyWithOption(dst, src, opt)`",
		"`dst` must be a pointer destination",
		"`copier.TypeConverter`",
		"`copier.FieldNameMapping`",
		"`WithIgnoreEmpty()`",
		"`IgnoreEmpty = true`",
		"`copier.WithDeepCopy()`",
		"`DeepCopy = true`",
		"`copier.WithCaseInsensitive()`",
		"`CaseSensitive = false`",
		"default copying remains case-sensitive",
		"`WithFieldNameMapping(...)`",
		"Appends mappings",
		"`WithTypeConverters(...)`",
		"Appends custom converters after the built-in converters",
		"Options are applied in the order they are passed",
		"`WithFieldNameMapping(...)` also appends mappings",
		"value ↔ pointer conversions",
		"`decimal.Decimal` ↔ `*decimal.Decimal`",
		"`time.Time` ↔ `*time.Time`",
		"`timex.DateTime` ↔ `*timex.DateTime`",
		"`timex.Date` ↔ `*timex.Date`",
		"`timex.Time` ↔ `*timex.Time`",
		"Value-to-pointer converters allocate a new local value",
		"Pointer-to-value converters dereference non-nil pointers",
		"If the source pointer is nil, the converter returns the zero value",
		"`decimal.Zero`",
		"zero `time.Time` /\n`timex.DateTime`",
		"`Create`, `CreateMany`,\n`Update`, and `UpdateMany`",
		"The update builders use `WithIgnoreEmpty()`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"9 个\ntop-level exported symbols",
		"没有 exported fields",
		"没有 exported methods",
		"fingerprint 是\n`44b6cf428fb9c642afca0cd25257c8ade57c9ac855b3ecc67cf575c1323fdf58`",
		"`copier.Copy(src, dst, options...)`",
		"`copier.CopyWithOption(dst, src, opt)`",
		"`dst` 必须是 pointer destination",
		"`copier.TypeConverter`",
		"`copier.FieldNameMapping`",
		"`copier.WithIgnoreEmpty()`",
		"`IgnoreEmpty = true`",
		"`copier.WithDeepCopy()`",
		"`DeepCopy = true`",
		"`copier.WithCaseInsensitive()`",
		"`CaseSensitive = false`",
		"默认复制仍然是 case-sensitive",
		"`copier.WithFieldNameMapping(...)`",
		"追加 mappings",
		"`copier.WithTypeConverters(...)`",
		"追加自定义 converters，而不是替换内置 converters",
		"options 会按传入顺序应用",
		"`WithFieldNameMapping(...)` 也会追加 mappings",
		"值 ↔ 指针自动转换器",
		"`decimal.Decimal` ↔ `*decimal.Decimal`",
		"`time.Time` ↔ `*time.Time`",
		"`timex.DateTime` ↔ `*timex.DateTime`",
		"`timex.Date` ↔ `*timex.Date`",
		"`timex.Time` ↔ `*timex.Time`",
		"value-to-pointer converter 会分配一个新的局部值",
		"pointer-to-value converter 会解引用非 nil 指针",
		"source pointer 为 nil",
		"返回目标类型的零值",
		"`decimal.Zero`",
		"zero `time.Time` / `timex.DateTime`",
		"`Create`、`CreateMany`、`Update` 和 `UpdateMany`",
		"`WithIgnoreEmpty()`",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "copier/copier.go",
			terms: []string{
				"CopyOption func(option *copier.Option)",
				"TypeConverter = copier.TypeConverter",
				"FieldNameMapping = copier.FieldNameMapping",
				"func WithIgnoreEmpty() CopyOption",
				"option.IgnoreEmpty = true",
				"func WithDeepCopy() CopyOption",
				"option.DeepCopy = true",
				"func WithCaseInsensitive() CopyOption",
				"option.CaseSensitive = false",
				"func WithFieldNameMapping(mappings ...FieldNameMapping) CopyOption",
				"option.FieldNameMapping = append(option.FieldNameMapping, mappings...)",
				"func WithTypeConverters(converters ...TypeConverter) CopyOption",
				"option.Converters = append(option.Converters, converters...)",
				"func Copy(src, dst any, options ...CopyOption) error",
				"CaseSensitive: true",
				"Converters:    defaultConverters",
				"for _, apply := range options",
				"apply(&opt)",
				"return copier.CopyWithOption(dst, src, opt)",
			},
		},
		{
			path: "copier/converters.go",
			terms: []string{
				"func makeValueToPtrConverter[T any]() TypeConverter",
				"return &v, nil",
				"func makePtrToValueConverter[T any]() TypeConverter",
				"if p := src.(*T); p != nil",
				"return *p, nil",
				"return lo.Empty[T](), nil",
				"stringToStringPtrConverter",
				"boolToBoolPtrConverter",
				"intToIntPtrConverter",
				"uintToUintPtrConverter",
				"float64ToFloat64PtrConverter",
				"decimalToDecimalPtrConverter",
				"timeToTimePtrConverter",
				"dateTimeToDateTimePtrConverter",
				"dateToDatePtrConverter",
				"timexTimeToTimexTimePtrConverter",
			},
		},
		{
			path: "copier/copier_test.go",
			terms: []string{
				"TestCopyBasic",
				"TestCopyValueToPtr",
				"StringToPtr",
				"BoolToPtr",
				"IntToPtr",
				"DecimalToPtr",
				"TimeToPtr",
				"DateTimeToPtr",
				"DateToPtr",
				"TimexTimeToPtr",
				"TestCopyPtrToValue",
				"NilStringPtrToValue",
				"NilBoolPtrToValue",
				"NilInt64PtrToValue",
				"NilDecimalPtrToValue",
				"NilTimePtrToValue",
				"NilDateTimePtrToValue",
				"TestCopyOptions",
				"IgnoreEmpty",
				"CaseInsensitive",
				"TestCopyDeepCopy",
				"DeepCopySlice",
				"DeepCopyNestedStruct",
				"TestCopyFieldNameMapping",
				"TestCopyError",
				"NonPointerDestination",
			},
		},
		{
			path:  "crud/create.go",
			terms: []string{"copier.Copy(&params, &model)"},
		},
		{
			path:  "crud/create_many.go",
			terms: []string{"copier.Copy(&params.List[i], &models[i])"},
		},
		{
			path: "crud/update.go",
			terms: []string{
				"copier.Copy(&params, &model)",
				"copier.Copy(&model, &oldModel, copier.WithIgnoreEmpty())",
			},
		},
		{
			path: "crud/update_many.go",
			terms: []string{
				"copier.Copy(&params.List[i], &models[i])",
				"copier.Copy(&models[i], &oldModels[i], copier.WithIgnoreEmpty())",
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
		panic(fmt.Errorf("copier contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Copier contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse copier package: %w", err))
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
		"`copier.Copy(src, dst, options...)`",
		"`copier.CopyOption`",
		"`copier.TypeConverter`",
		"`copier.FieldNameMapping`",
		"`copier.WithIgnoreEmpty()`",
		"`copier.WithDeepCopy()`",
		"`copier.WithCaseInsensitive()`",
		"`copier.WithFieldNameMapping(...)`",
		"`copier.WithTypeConverters(...)`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/copier",
		"FUNC Copy : func(src any, dst any, options ...github.com/coldsmirk/vef-framework-go/copier.CopyOption) error",
		"TYPE CopyOption : github.com/coldsmirk/vef-framework-go/copier.CopyOption",
		"TYPE FieldNameMapping : github.com/coldsmirk/vef-framework-go/copier.FieldNameMapping",
		"TYPE TypeConverter : github.com/coldsmirk/vef-framework-go/copier.TypeConverter",
		"FUNC WithCaseInsensitive : func() github.com/coldsmirk/vef-framework-go/copier.CopyOption",
		"FUNC WithDeepCopy : func() github.com/coldsmirk/vef-framework-go/copier.CopyOption",
		"FUNC WithFieldNameMapping : func(mappings ...github.com/coldsmirk/vef-framework-go/copier.FieldNameMapping) github.com/coldsmirk/vef-framework-go/copier.CopyOption",
		"FUNC WithIgnoreEmpty : func() github.com/coldsmirk/vef-framework-go/copier.CopyOption",
		"FUNC WithTypeConverters : func(converters ...github.com/coldsmirk/vef-framework-go/copier.TypeConverter) github.com/coldsmirk/vef-framework-go/copier.CopyOption",
	}
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./copier")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./copier failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func compareNames(label string, actual, expected []string) []string {
	actualSet := toSet(actual)
	expectedSet := toSet(expected)
	var failures []string

	for _, name := range expected {
		if !actualSet[name] {
			failures = append(failures, fmt.Sprintf("missing %s %s", label, name))
		}
	}
	for _, name := range actual {
		if !expectedSet[name] {
			failures = append(failures, fmt.Sprintf("unexpected %s %s", label, name))
		}
	}

	return failures
}

func toSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}

	return result
}

func missingTerms(doc corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !containsTerm(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term %q", doc.label, term))
		}
	}

	return failures
}

func containsTerm(content string, term string) bool {
	if strings.Contains(content, term) {
		return true
	}

	return strings.Contains(normalizeSpace(content), normalizeSpace(term))
}

func normalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func receiverBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverBaseName(t.X)
	case *ast.IndexExpr:
		return receiverBaseName(t.X)
	case *ast.IndexListExpr:
		return receiverBaseName(t.X)
	default:
		return ""
	}
}

func readCorpus(label string, path string) corpus {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(data)}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("failed to resolve %s: %w", path, err))
	}

	return filepath.Clean(abs)
}
