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
	"regexp"
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

	englishDocs := readCorpus("English contextx docs", filepath.Join(docsRoot, "docs/advanced/extending-parameters.md"))
	chineseDocs := readCorpus("Chinese contextx docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/advanced/extending-parameters.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	contextxDir := filepath.Join(sourceRoot, "contextx")

	expectedFuncs := []string{
		"DB", "DataPermApplier", "Logger", "Principal",
		"RequestID", "RequestIP", "RequestMethod", "RequestPath",
		"SetDB", "SetDataPermApplier", "SetLogger", "SetPrincipal",
		"SetRequestID", "SetRequestIP", "SetRequestMethod", "SetRequestPath",
	}
	expectedConsts := []string{
		"KeyRequest", "KeyRequestID", "KeyRequestIP", "KeyPrincipal", "KeyLogger",
		"KeyDB", "KeyDataPermApplier", "KeyRequestMethod", "KeyRequestPath",
	}
	expectedConstValues := map[string]int64{
		"KeyRequest":         0,
		"KeyRequestID":       1,
		"KeyRequestIP":       2,
		"KeyPrincipal":       3,
		"KeyLogger":          4,
		"KeyDB":              5,
		"KeyDataPermApplier": 6,
		"KeyRequestMethod":   7,
		"KeyRequestPath":     8,
	}

	var failures []string

	exported := exportedPackageSurface(contextxDir)
	failures = append(failures, compareNames("contextx func", exported.funcs, expectedFuncs)...)
	failures = append(failures, compareNames("contextx const", exported.consts, expectedConsts)...)
	if len(exported.types) > 0 {
		failures = append(failures, "contextx package should not expose types, found: "+strings.Join(exported.types, ", "))
	}
	if len(exported.vars) > 0 {
		failures = append(failures, "contextx package should not expose vars, found: "+strings.Join(exported.vars, ", "))
	}
	if len(exported.methods) > 0 {
		failures = append(failures, "contextx package should not expose methods, found: "+strings.Join(exported.methods, ", "))
	}
	if len(exported.fields) > 0 {
		failures = append(failures, "contextx package should not expose fields, found: "+strings.Join(exported.fields, ", "))
	}

	constValues := exportedConstValues(contextxDir)
	for _, name := range expectedConsts {
		got, ok := constValues[name]
		if !ok {
			failures = append(failures, "missing typed constant value for "+name)
			continue
		}
		if got != expectedConstValues[name] {
			failures = append(failures, fmt.Sprintf("constant %s value drifted: got %d, want %d", name, got, expectedConstValues[name]))
		}
	}

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}

	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
		failures = append(failures, forbidPhantomContextxSymbols(doc, append(expectedFuncs, expectedConsts...))...)
		failures = append(failures, forbids(doc, []string{
			"contextx.Request(", "contextx.SetRequest(",
			"`Request`", "`SetRequest`",
			"contextKey", "setValue",
		})...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"unexported key type",
		"cannot construct additional keys of the same type",
		"No Request or SetRequest accessor exists in `contextx`",
		"the package itself does not read or write this key",
		"`KeyRequest = 0`", "`KeyRequestID = 1`", "`KeyRequestIP = 2`",
		"`KeyPrincipal = 3`", "`KeyLogger = 4`", "`KeyDB = 5`",
		"`KeyDataPermApplier = 6`", "`KeyRequestMethod = 7`", "`KeyRequestPath = 8`",
		"cannot distinguish \"unset\" from \"explicitly set to an empty string\"",
		"`Logger` and `DB` first return a correctly typed context value",
		"Fallbacks are scanned left to right",
		"nil and typed nil fallbacks are skipped",
		"a typed nil value already stored in the context still wins",
		"The original context is unchanged",
		"Always keep the returned context when using a standard context",
		"same Fiber context",
		"Signature auth binds both values into signature verification",
		"Signature auth also uses it for IP whitelist checks",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"未导出的 key 类型",
		"不能构造同类型的新 key",
		"`contextx` 中不存在 Request 或 SetRequest accessor",
		"这个 package 自身不读写该 key",
		"`KeyRequest = 0`", "`KeyRequestID = 1`", "`KeyRequestIP = 2`",
		"`KeyPrincipal = 3`", "`KeyLogger = 4`", "`KeyDB = 5`",
		"`KeyDataPermApplier = 6`", "`KeyRequestMethod = 7`", "`KeyRequestPath = 8`",
		"无法区分“未设置”和“显式设置为空字符串”",
		"`Logger` 和 `DB` 会先返回 context 中类型正确的值",
		"fallbacks 按从左到右扫描",
		"nil 和 typed nil fallbacks 会被跳过",
		"如果 typed nil 已经存进 context，类型断言成功后仍会直接返回",
		"原 context 不会被修改",
		"使用标准 context 时，必须保留返回值",
		"同一个 Fiber context",
		"Signature auth 会把两者绑定进签名校验",
		"Signature auth 也会用它做 IP whitelist 检查",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "contextx/contextx.go",
			terms: []string{
				"type contextKey int",
				"KeyRequest contextKey = iota",
				"KeyRequestID", "KeyRequestIP", "KeyPrincipal", "KeyLogger", "KeyDB",
				"KeyDataPermApplier", "KeyRequestMethod", "KeyRequestPath",
				"if c, ok := ctx.(fiber.Ctx); ok",
				"c.Locals(key, value)",
				"return context.WithValue(ctx, key, value)",
				"ctx.Value(KeyRequestID).(string)",
				"ctx.Value(KeyRequestIP).(string)",
				"ctx.Value(KeyRequestMethod).(string)",
				"ctx.Value(KeyRequestPath).(string)",
				"ctx.Value(KeyPrincipal).(*security.Principal)",
				"ctx.Value(KeyDataPermApplier).(security.DataPermissionApplier)",
				"ctx.Value(KeyLogger).(logx.Logger)",
				"ctx.Value(KeyDB).(orm.DB)",
				"reflectx.IsNotEmpty(fallback)",
			},
		},
		{
			path: "contextx/contextx_test.go",
			terms: []string{
				"ReturnsEmptyStringFromEmptyContext",
				"ReturnsEmptyStringForWrongType",
				"SkipsTypedNilFallback",
				"ContextValueTakesPrecedenceOverFallback",
				"WorksWithFiberContext",
			},
		},
		{
			path: "internal/middleware/logger.go",
			terms: []string{
				"contextx.SetLogger(ctx, logger)",
				"contextx.SetRequestID(ctx, requestID)",
				"contextx.SetLogger(",
				"contextx.SetRequestID(ctx.Context(), requestID)",
			},
		},
		{
			path: "internal/api/middleware/auth.go",
			terms: []string{
				"contextx.SetRequestIP(ctx.Context(), httpx.GetIP(ctx))",
				"contextx.SetRequestMethod(reqCtx, ctx.Method())",
				"contextx.SetRequestPath(reqCtx, ctx.Path())",
				"contextx.SetPrincipal(ctx, principal)",
				"contextx.SetPrincipal(ctx.Context(), principal)",
			},
		},
		{
			path: "internal/api/middleware/contextual.go",
			terms: []string{
				"contextx.SetDB(ctx, db)",
				"contextx.SetDB(ctx.Context(), db)",
				"contextx.SetLogger(ctx, scopedLogger)",
				"contextx.SetLogger(ctx.Context(), scopedLogger)",
			},
		},
		{
			path: "internal/api/middleware/data_permission.go",
			terms: []string{
				"contextx.SetDataPermApplier(ctx, applier)",
				"contextx.SetDataPermApplier(ctx.Context(), applier)",
			},
		},
		{
			path: "internal/security/signature_authenticator.go",
			terms: []string{
				"method := contextx.RequestMethod(ctx)",
				"path := contextx.RequestPath(ctx)",
				"requestIP := contextx.RequestIP(ctx)",
				"return security.ErrIPNotAllowed",
			},
		},
	}
	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	failures = append(failures, forbidSourceAccessors(filepath.Join(sourceRoot, "contextx", "contextx.go"))...)
	failures = append(failures, runRuntimeChecks(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("contextx contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("Contextx contract docs verified: %d public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
		len(expectedFuncs)+len(expectedConsts), len(exported.methods), len(exported.fields), len(sourceChecks))
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
		panic(fmt.Errorf("failed to parse contextx package: %w", err))
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

func exportedConstValues(dir string) map[string]int64 {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse contextx package for const values: %w", err))
	}

	values := make(map[string]int64)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.CONST {
					continue
				}

				lastWasIota := false
				for iotaIndex, spec := range genDecl.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					if len(valueSpec.Values) > 0 {
						lastWasIota = isIotaExpr(valueSpec.Values[0])
					}
					if !lastWasIota {
						continue
					}
					for _, name := range valueSpec.Names {
						if name.IsExported() {
							values[name.Name] = int64(iotaIndex)
						}
					}
				}
			}
		}
	}

	return values
}

func isIotaExpr(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "iota"
}

func publicDocSurfaceTerms() []string {
	return []string{
		"`KeyRequest`", "`KeyRequestID`", "`KeyRequestIP`",
		"`KeyPrincipal`", "`KeyLogger`", "`KeyDB`",
		"`KeyDataPermApplier`", "`KeyRequestMethod`", "`KeyRequestPath`",
		"`RequestID`", "`contextx.RequestID(ctx context.Context) string`",
		"`SetRequestID`", "`contextx.SetRequestID(ctx context.Context, requestID string) context.Context`",
		"`RequestIP`", "`contextx.RequestIP(ctx context.Context) string`",
		"`SetRequestIP`", "`contextx.SetRequestIP(ctx context.Context, ip string) context.Context`",
		"`RequestMethod`", "`contextx.RequestMethod(ctx context.Context) string`",
		"`SetRequestMethod`", "`contextx.SetRequestMethod(ctx context.Context, method string) context.Context`",
		"`RequestPath`", "`contextx.RequestPath(ctx context.Context) string`",
		"`SetRequestPath`", "`contextx.SetRequestPath(ctx context.Context, path string) context.Context`",
		"`Principal`", "`contextx.Principal(ctx context.Context) *security.Principal`",
		"`SetPrincipal`", "`contextx.SetPrincipal(ctx context.Context, principal *security.Principal) context.Context`",
		"`Logger`", "`contextx.Logger(ctx context.Context, fallbacks ...logx.Logger) logx.Logger`",
		"`SetLogger`", "`contextx.SetLogger(ctx context.Context, logger logx.Logger) context.Context`",
		"`DB`", "`contextx.DB(ctx context.Context, fallbacks ...orm.DB) orm.DB`",
		"`SetDB`", "`contextx.SetDB(ctx context.Context, db orm.DB) context.Context`",
		"`DataPermApplier`", "`contextx.DataPermApplier(ctx context.Context) security.DataPermissionApplier`",
		"`SetDataPermApplier`", "`contextx.SetDataPermApplier(ctx context.Context, applier security.DataPermissionApplier) context.Context`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/contextx",
		"FUNC DB : func(ctx context.Context, fallbacks ...github.com/coldsmirk/vef-framework-go/orm.DB) github.com/coldsmirk/vef-framework-go/orm.DB",
		"FUNC DataPermApplier : func(ctx context.Context) github.com/coldsmirk/vef-framework-go/security.DataPermissionApplier",
		"CONST KeyDB : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 5",
		"CONST KeyDataPermApplier : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 6",
		"CONST KeyLogger : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 4",
		"CONST KeyPrincipal : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 3",
		"CONST KeyRequest : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 0",
		"CONST KeyRequestID : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 1",
		"CONST KeyRequestIP : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 2",
		"CONST KeyRequestMethod : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 7",
		"CONST KeyRequestPath : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 8",
		"FUNC Logger : func(ctx context.Context, fallbacks ...github.com/coldsmirk/vef-framework-go/logx.Logger) github.com/coldsmirk/vef-framework-go/logx.Logger",
		"FUNC Principal : func(ctx context.Context) *github.com/coldsmirk/vef-framework-go/security.Principal",
		"FUNC RequestID : func(ctx context.Context) string",
		"FUNC RequestIP : func(ctx context.Context) string",
		"FUNC RequestMethod : func(ctx context.Context) string",
		"FUNC RequestPath : func(ctx context.Context) string",
		"FUNC SetDB : func(ctx context.Context, db github.com/coldsmirk/vef-framework-go/orm.DB) context.Context",
		"FUNC SetDataPermApplier : func(ctx context.Context, applier github.com/coldsmirk/vef-framework-go/security.DataPermissionApplier) context.Context",
		"FUNC SetLogger : func(ctx context.Context, logger github.com/coldsmirk/vef-framework-go/logx.Logger) context.Context",
		"FUNC SetPrincipal : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) context.Context",
		"FUNC SetRequestID : func(ctx context.Context, requestID string) context.Context",
		"FUNC SetRequestIP : func(ctx context.Context, ip string) context.Context",
		"FUNC SetRequestMethod : func(ctx context.Context, method string) context.Context",
		"FUNC SetRequestPath : func(ctx context.Context, path string) context.Context",
	}
}

func forbidPhantomContextxSymbols(c corpus, allowed []string) []string {
	allowedSet := set(allowed)
	allowedSet["contextx"] = true

	re := regexp.MustCompile(`\bcontextx\.([A-Za-z_][A-Za-z0-9_]*)\b`)
	matches := re.FindAllStringSubmatch(c.content, -1)
	var failures []string
	for _, match := range matches {
		if !allowedSet[match[1]] {
			failures = append(failures, fmt.Sprintf("%s references non-public or nonexistent contextx symbol: %s", c.label, match[0]))
		}
	}

	return failures
}

func forbidSourceAccessors(path string) []string {
	content := readCorpus("contextx source accessor scan", path)
	forbidden := []string{
		"func Request(ctx context.Context)",
		"func SetRequest(ctx context.Context",
	}

	return forbids(content, forbidden)
}

func runRuntimeChecks(sourceRoot string) []string {
	tmpDir, err := os.MkdirTemp("", "verify-contextx-contracts-*")
	if err != nil {
		return []string{fmt.Sprintf("failed to create temp module: %v", err)}
	}
	defer os.RemoveAll(tmpDir)

	goMod := fmt.Sprintf(`module verifycontextxcontracts

go 1.26.1

require github.com/coldsmirk/vef-framework-go v0.0.0

replace github.com/coldsmirk/vef-framework-go => %s
`, sourceRoot)
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o600); err != nil {
		return []string{fmt.Sprintf("failed to write temp go.mod: %v", err)}
	}

	mainGo := `package main

import (
	"context"
	"fmt"

	"github.com/coldsmirk/vef-framework-go/contextx"
	"github.com/coldsmirk/vef-framework-go/logx"
)

type mockLogger struct{}

func (*mockLogger) Named(name string) logx.Logger { return &mockLogger{} }
func (*mockLogger) WithCallerSkip(skip int) logx.Logger { return &mockLogger{} }
func (*mockLogger) Enabled(level logx.Level) bool { return true }
func (*mockLogger) Sync() {}
func (*mockLogger) Debug(message string) {}
func (*mockLogger) Debugf(template string, args ...any) {}
func (*mockLogger) Info(message string) {}
func (*mockLogger) Infof(template string, args ...any) {}
func (*mockLogger) Warn(message string) {}
func (*mockLogger) Warnf(template string, args ...any) {}
func (*mockLogger) Error(message string) {}
func (*mockLogger) Errorf(template string, args ...any) {}
func (*mockLogger) Panic(message string) { panic(message) }
func (*mockLogger) Panicf(template string, args ...any) { panic(fmt.Sprintf(template, args...)) }

func main() {
	ctx := context.Background()
	_ = contextx.SetRequestID(ctx, "not-captured")
	if got := contextx.RequestID(ctx); got != "" {
		panic(fmt.Sprintf("stdlib SetRequestID without capture should not mutate original context, got %q", got))
	}

	ctx = contextx.SetRequestID(ctx, "captured")
	if got := contextx.RequestID(ctx); got != "captured" {
		panic(fmt.Sprintf("stdlib SetRequestID with capture should store value, got %q", got))
	}

	if got := contextx.RequestID(context.WithValue(context.Background(), contextx.KeyRequestID, 123)); got != "" {
		panic(fmt.Sprintf("wrong type RequestID should return empty string, got %q", got))
	}

	if got := contextx.RequestMethod(context.WithValue(context.Background(), contextx.KeyRequestMethod, 123)); got != "" {
		panic(fmt.Sprintf("wrong type RequestMethod should return empty string, got %q", got))
	}

	if got := contextx.RequestPath(context.WithValue(context.Background(), contextx.KeyRequestPath, 123)); got != "" {
		panic(fmt.Sprintf("wrong type RequestPath should return empty string, got %q", got))
	}

	var typedNil *mockLogger
	fallback := &mockLogger{}
	if got := contextx.Logger(context.Background(), typedNil, fallback); got != fallback {
		panic("Logger should skip typed nil fallbacks and return the first non-empty fallback")
	}

	ctx = contextx.SetLogger(context.Background(), typedNil)
	if got := contextx.Logger(ctx, fallback); got != logx.Logger(typedNil) {
		panic("Logger should return typed nil stored in context before checking fallbacks")
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainGo), 0o600); err != nil {
		return []string{fmt.Sprintf("failed to write temp main.go: %v", err)}
	}

	cmd := exec.Command("go", "run", "-mod=mod", ".")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("contextx runtime contract check failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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

func forbids(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if strings.Contains(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s must not include forbidden term: %s", c.label, term))
		}
	}

	return failures
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
