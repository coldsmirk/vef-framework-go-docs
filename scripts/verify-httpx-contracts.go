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

	httpxDir := filepath.Join(sourceRoot, "httpx")
	expectedFuncs := []string{"GetIP", "IsJSON", "IsMultipart"}

	var failures []string
	exported := exportedPackageSurface(httpxDir)
	failures = append(failures, compareNames("httpx func", exported.funcs, expectedFuncs)...)
	if len(exported.consts) > 0 {
		failures = append(failures, "httpx package should not expose consts, found: "+strings.Join(exported.consts, ", "))
	}
	if len(exported.types) > 0 {
		failures = append(failures, "httpx package should not expose types, found: "+strings.Join(exported.types, ", "))
	}
	if len(exported.vars) > 0 {
		failures = append(failures, "httpx package should not expose vars, found: "+strings.Join(exported.vars, ", "))
	}
	if len(exported.methods) > 0 {
		failures = append(failures, "httpx package should not expose methods, found: "+strings.Join(exported.methods, ", "))
	}
	if len(exported.fields) > 0 {
		failures = append(failures, "httpx package should not expose fields, found: "+strings.Join(exported.fields, ", "))
	}

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"delegates to Fiber's `ctx.Is(\"json\")`",
		"accepts standard JSON",
		"content types including charset variants",
		"`strings.HasPrefix(...)`",
		"boundary parameters such as `multipart/form-data; boundary=...`",
		"`GetIP` delegates to `ctx.IP()`",
		"`vef.app.trusted_proxies` controls whether proxy headers are trusted",
		"raw client-supplied `X-Forwarded-For` is ignored",
		"Fiber may honor `X-Forwarded-For` according",
		"to its proxy settings",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"委托 Fiber 的 `ctx.Is(\"json\")`",
		"接受标准 JSON content type",
		"`strings.HasPrefix(...)`",
		"`multipart/form-data; boundary=...`",
		"`GetIP` 委托 `ctx.IP()`",
		"`vef.app.trusted_proxies` 控制 proxy headers 是否可信",
		"客户端直接伪造的 `X-Forwarded-For` 会被忽略",
		"Fiber 会按自己的 proxy settings 处理 `X-Forwarded-For`",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "httpx/content_type.go",
			terms: []string{
				"func IsJSON(ctx fiber.Ctx) bool",
				"return ctx.Is(\"json\")",
				"func IsMultipart(ctx fiber.Ctx) bool",
				"return strings.HasPrefix(ctx.Get(fiber.HeaderContentType), fiber.MIMEMultipartForm)",
			},
		},
		{
			path: "httpx/content_type_test.go",
			terms: []string{
				"ApplicationJson",
				"ApplicationJsonWithCharset",
				"MissingContentType",
				"NonJsonContentType",
				"MultipartFormData",
				"fiber.MIMEMultipartForm+\"; boundary=MyBoundary\"",
			},
		},
		{
			path: "httpx/ip.go",
			terms: []string{
				"func GetIP(ctx fiber.Ctx) string",
				"return ctx.IP()",
				"trusted proxies (vef.app.trusted_proxies)",
				"X-Forwarded-For",
			},
		},
		{
			path: "httpx/ip_test.go",
			terms: []string{
				"IgnoresUntrustedXForwardedFor",
				"DirectIPWhenNoHeader",
				"HonorsXForwardedForFromTrustedProxy",
				"TrustProxy: true",
				"ProxyHeader:      fiber.HeaderXForwardedFor",
			},
		},
		{
			path: "internal/api/middleware/auth.go",
			terms: []string{
				"contextx.SetRequestIP(ctx.Context(), httpx.GetIP(ctx))",
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
		panic(fmt.Errorf("httpx contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("HTTPX contract docs verified: %d public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse httpx package: %w", err))
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
		"`IsJSON`", "`httpx.IsJSON(ctx fiber.Ctx) bool`",
		"`IsMultipart`", "`httpx.IsMultipart(ctx fiber.Ctx) bool`",
		"`GetIP`", "`httpx.GetIP(ctx fiber.Ctx) string`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/httpx",
		"FUNC GetIP : func(ctx github.com/gofiber/fiber/v3.Ctx) string",
		"FUNC IsJSON : func(ctx github.com/gofiber/fiber/v3.Ctx) bool",
		"FUNC IsMultipart : func(ctx github.com/gofiber/fiber/v3.Ctx) bool",
	}
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./httpx")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./httpx failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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
