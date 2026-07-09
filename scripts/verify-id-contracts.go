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

	englishDocs := readCorpus("English ID docs", filepath.Join(docsRoot, "docs/utilities/id-generation.md"))
	chineseDocs := readCorpus("Chinese ID docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/id-generation.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	idDir := filepath.Join(sourceRoot, "id")
	expectedSurface := packageSurface{
		consts: []string{
			"DefaultRandomIDGeneratorAlphabet",
			"DefaultRandomIDGeneratorLength",
		},
		funcs: []string{
			"Generate",
			"GenerateUUID",
			"NewRandomIDGenerator",
			"NewUUIDGenerator",
			"NewXIDGenerator",
			"WithAlphabet",
			"WithLength",
		},
		types: []string{
			"IDGenerator",
			"RandomIDGeneratorOption",
		},
		vars: []string{
			"DefaultUUIDGenerator",
			"DefaultXIDGenerator",
		},
		methods: []string{"IDGenerator.Generate"},
	}

	var failures []string
	exported := exportedPackageSurface(idDir)
	failures = append(failures, compareNames("id const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("id func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("id type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("id var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("id method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("id exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"Delegates to `DefaultXIDGenerator.Generate()`",
		"20-character XID",
		"Delegates to `DefaultUUIDGenerator.Generate()`",
		"UUID v7 string",
		"Returns an `IDGenerator` that wraps `xid.New().String()`",
		"uses `uuid.NewV7()` and panics if UUID creation fails",
		"applies options in order",
		"`id.DefaultRandomIDGeneratorAlphabet`",
		"`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`",
		"`id.DefaultRandomIDGeneratorLength`",
		"empty alphabet or zero length panics when\n`Generate()` is called",
		"`go-nanoid/v2` `MustGenerate`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"委托 `DefaultXIDGenerator.Generate()`",
		"20 字符 XID",
		"委托 `DefaultUUIDGenerator.Generate()`",
		"UUID v7 string",
		"包装 `xid.New().String()`",
		"使用 `uuid.NewV7()` 的 `IDGenerator`；UUID 创建失败时 panic",
		"按顺序应用 options",
		"`id.DefaultRandomIDGeneratorAlphabet`",
		"`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`",
		"`id.DefaultRandomIDGeneratorLength`",
		"empty alphabet 或 zero length 会在调用 `Generate()` 时 panic",
		"`go-nanoid/v2` 的 `MustGenerate`",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "id/id.go",
			terms: []string{
				"type IDGenerator interface",
				"Generate() string",
				"func Generate() string",
				"return DefaultXIDGenerator.Generate()",
				"func GenerateUUID() string",
				"return DefaultUUIDGenerator.Generate()",
			},
		},
		{
			path: "id/uuid.go",
			terms: []string{
				"var DefaultUUIDGenerator = NewUUIDGenerator()",
				"type uuidGenerator struct{}",
				"func (*uuidGenerator) Generate() string",
				"id, err := uuid.NewV7()",
				"panic(fmt.Errorf(\"failed to generate UUID: %w\", err))",
				"return id.String()",
				"func NewUUIDGenerator() IDGenerator",
				"return &uuidGenerator{}",
			},
		},
		{
			path: "id/xid.go",
			terms: []string{
				"var DefaultXIDGenerator = NewXIDGenerator()",
				"type xidGenerator struct{}",
				"func (*xidGenerator) Generate() string",
				"return xid.New().String()",
				"func NewXIDGenerator() IDGenerator",
				"return &xidGenerator{}",
			},
		},
		{
			path: "id/random.go",
			terms: []string{
				"DefaultRandomIDGeneratorAlphabet = \"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\"",
				"DefaultRandomIDGeneratorLength   = 32",
				"type randomIDGenerator struct",
				"alphabet string",
				"length   int",
				"func (g *randomIDGenerator) Generate() string",
				"return nid.MustGenerate(g.alphabet, g.length)",
				"type RandomIDGeneratorOption func(*randomIDGenerator)",
				"func WithAlphabet(alphabet string) RandomIDGeneratorOption",
				"g.alphabet = alphabet",
				"func WithLength(length int) RandomIDGeneratorOption",
				"g.length = length",
				"func NewRandomIDGenerator(opts ...RandomIDGeneratorOption) IDGenerator",
				"alphabet: DefaultRandomIDGeneratorAlphabet",
				"length:   DefaultRandomIDGeneratorLength",
				"for _, opt := range opts",
				"opt(g)",
			},
		},
		{
			path: "id/id_test.go",
			terms: []string{
				"TestGenerate",
				"UseXIDGeneratorByDefault",
				"TestGenerateUUID",
				"TestDefaultGenerators",
				"TestConcurrentGeneration",
			},
		},
		{
			path: "id/random_test.go",
			terms: []string{
				"TestRandomIDGenerator",
				"CreateWithCustomAlphabetAndLength",
				"DefaultValues",
				"SingleCharacterAlphabet",
				"ThreadSafe",
			},
		},
		{
			path: "id/edge_cases_test.go",
			terms: []string{
				"TestRandomIdGeneratorEdgeCases",
				"EmptyAlphabet",
				"ZeroLength",
				"UnicodeCharacters",
				"TestInterfaceCompliance",
				"AllGeneratorsImplementInterface",
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
		panic(fmt.Errorf("id contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("ID contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse id package: %w", err))
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
		"`id.IDGenerator`",
		"`IDGenerator.Generate()`",
		"`id.Generate()`",
		"`id.GenerateUUID()`",
		"`id.DefaultXIDGenerator`",
		"`id.DefaultUUIDGenerator`",
		"`id.NewXIDGenerator()`",
		"`id.NewUUIDGenerator()`",
		"`id.NewRandomIDGenerator(opts...)`",
		"`id.RandomIDGeneratorOption`",
		"`id.WithAlphabet(alphabet)`",
		"`id.WithLength(length)`",
		"`id.DefaultRandomIDGeneratorAlphabet`",
		"`id.DefaultRandomIDGeneratorLength`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/id",
		"CONST DefaultRandomIDGeneratorAlphabet : untyped string = \"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\"",
		"CONST DefaultRandomIDGeneratorLength : untyped int = 32",
		"VAR DefaultUUIDGenerator : github.com/coldsmirk/vef-framework-go/id.IDGenerator",
		"VAR DefaultXIDGenerator : github.com/coldsmirk/vef-framework-go/id.IDGenerator",
		"FUNC Generate : func() string",
		"FUNC GenerateUUID : func() string",
		"TYPE IDGenerator : github.com/coldsmirk/vef-framework-go/id.IDGenerator",
		"METHOD Generate : func() string",
		"FUNC NewRandomIDGenerator : func(opts ...github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption) github.com/coldsmirk/vef-framework-go/id.IDGenerator",
		"FUNC NewUUIDGenerator : func() github.com/coldsmirk/vef-framework-go/id.IDGenerator",
		"FUNC NewXIDGenerator : func() github.com/coldsmirk/vef-framework-go/id.IDGenerator",
		"TYPE RandomIDGeneratorOption : github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption",
		"FUNC WithAlphabet : func(alphabet string) github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption",
		"FUNC WithLength : func(length int) github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption",
	}
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./id")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./id failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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
