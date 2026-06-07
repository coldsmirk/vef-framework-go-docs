package main

import (
	"encoding/json"
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

const (
	i18nPackage          = "github.com/coldsmirk/vef-framework-go/i18n"
	localesPackage       = "github.com/coldsmirk/vef-framework-go/i18n/locales"
	i18nFingerprint      = "67076d160fa9ebb52302f2873faebeb52fd2be12e1386f8f2825b01130da74ae"
	localesFingerprint   = "71bc852d645cd378f46663e5998b7e22b9c7372906fdaf6a04ded7acc122ade0"
	i18nTopLevelCount    = 12
	i18nFieldCount       = 1
	i18nMethodCount      = 2
	i18nEntryCount       = 15
	localesTopLevelCount = 1
	localesFieldCount    = 0
	localesMethodCount   = 0
	localesEntryCount    = 1
	localeMessageCount   = 160
)

type corpus struct {
	label   string
	content string
}

type packageSurface struct {
	consts  []string
	funcs   []string
	types   []string
	vars    []string
	fields  []string
	methods []string
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishDocs := readCorpus("English i18n docs", filepath.Join(docsRoot, "docs/features/i18n.md"))
	chineseDocs := readCorpus("Chinese i18n docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/i18n.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	expectedI18n := packageSurface{
		consts: []string{
			"DefaultLanguage",
		},
		funcs: []string{
			"CurrentLanguage",
			"GetSupportedLanguages",
			"IsLanguageSupported",
			"New",
			"SetLanguage",
			"T",
			"Te",
		},
		types: []string{
			"Config",
			"Translator",
		},
		vars: []string{
			"ErrMessageIDEmpty",
			"ErrUnsupportedLanguage",
		},
		fields: []string{
			"Config.Locales",
		},
		methods: []string{
			"Translator.T",
			"Translator.Te",
		},
	}
	expectedLocales := packageSurface{
		vars: []string{
			"EmbedLocales",
		},
	}

	var failures []string
	i18nSurface := exportedPackageSurface(filepath.Join(sourceRoot, "i18n"))
	localesSurface := exportedPackageSurface(filepath.Join(sourceRoot, "i18n/locales"))

	failures = append(failures, compareNames("i18n const", i18nSurface.consts, expectedI18n.consts)...)
	failures = append(failures, compareNames("i18n func", i18nSurface.funcs, expectedI18n.funcs)...)
	failures = append(failures, compareNames("i18n type", i18nSurface.types, expectedI18n.types)...)
	failures = append(failures, compareNames("i18n var", i18nSurface.vars, expectedI18n.vars)...)
	failures = append(failures, compareNames("i18n exported field", i18nSurface.fields, expectedI18n.fields)...)
	failures = append(failures, compareNames("i18n exported method", i18nSurface.methods, expectedI18n.methods)...)
	failures = append(failures, compareSurfaceCounts("i18n", i18nSurface, i18nTopLevelCount, i18nFieldCount, i18nMethodCount, i18nEntryCount)...)

	failures = append(failures, compareNames("i18n/locales const", localesSurface.consts, expectedLocales.consts)...)
	failures = append(failures, compareNames("i18n/locales func", localesSurface.funcs, expectedLocales.funcs)...)
	failures = append(failures, compareNames("i18n/locales type", localesSurface.types, expectedLocales.types)...)
	failures = append(failures, compareNames("i18n/locales var", localesSurface.vars, expectedLocales.vars)...)
	failures = append(failures, compareNames("i18n/locales exported field", localesSurface.fields, expectedLocales.fields)...)
	failures = append(failures, compareNames("i18n/locales exported method", localesSurface.methods, expectedLocales.methods)...)
	failures = append(failures, compareSurfaceCounts("i18n/locales", localesSurface, localesTopLevelCount, localesFieldCount, localesMethodCount, localesEntryCount)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, verifyI18nReferences(doc, expectedI18n)...)
		failures = append(failures, verifyLocalesReferences(doc, expectedLocales)...)
		failures = append(failures, verifySurfaceMentioned(doc, "i18n", i18nSurface)...)
		failures = append(failures, verifySurfaceMentioned(doc, "locales", localesSurface)...)
	}

	failures = append(failures, missingTerms(englishDocs, englishDocTerms())...)
	failures = append(failures, missingTerms(chineseDocs, chineseDocTerms())...)
	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, verifyLocaleCatalogs(sourceRoot)...)
	failures = append(failures, verifyAuditArtifacts(sourceRoot, docsRoot)...)
	failures = append(failures, runPackageTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("i18n contracts verified")
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/i18n",
		"TYPE Config : github.com/coldsmirk/vef-framework-go/i18n.Config",
		"FIELD Locales : embed.FS [field_order=1 tag=\"\"]",
		"FUNC CurrentLanguage : func() string",
		"CONST DefaultLanguage : untyped string = \"zh-CN\"",
		"VAR ErrMessageIDEmpty : error",
		"VAR ErrUnsupportedLanguage : error",
		"FUNC GetSupportedLanguages : func() []string",
		"FUNC IsLanguageSupported : func(languageCode string) bool",
		"FUNC New : func(config github.com/coldsmirk/vef-framework-go/i18n.Config) (github.com/coldsmirk/vef-framework-go/i18n.Translator, error)",
		"FUNC SetLanguage : func(languageCode string) error",
		"FUNC T : func(messageID string, templateData ...map[string]any) string",
		"FUNC Te : func(messageID string, templateData ...map[string]any) (string, error)",
		"TYPE Translator : github.com/coldsmirk/vef-framework-go/i18n.Translator",
		"METHOD T : func(messageID string, templateData ...map[string]any) string",
		"METHOD Te : func(messageID string, templateData ...map[string]any) (string, error)",
		"## github.com/coldsmirk/vef-framework-go/i18n/locales",
		"VAR EmbedLocales : embed.FS",
	}
}

func englishDocTerms() []string {
	return []string{
		"12 top-level exported symbols, 1 exported field, and 2 exported methods",
		"public surface fingerprint is `" + i18nFingerprint + "`",
		"1 top-level exported symbol, 0 exported fields, and 0 exported methods",
		"public surface fingerprint is `" + localesFingerprint + "`",
		"`i18n.GetSupportedLanguages()` returns a copy",
		"`en.json` and `zh-CN.json` catalogs contain the same 160 message IDs",
		"`i18n.DefaultLanguage` | default language constant (`zh-CN`)",
		"`i18n.GetSupportedLanguages()` | return a copy of the supported language codes (`zh-CN`, `en`)",
		"`i18n.IsLanguageSupported(code)` | report whether a language code is supported",
		"`i18n.CurrentLanguage()` | read the current global language code",
		"`i18n.SetLanguage(code)` | atomically switch the process-level global translator",
		"`i18n.T(messageID, data...)` | translate with graceful fallback to the message ID",
		"`i18n.Te(messageID, data...)` | translate with explicit error reporting",
		"`i18n.New(config)` | create a dedicated translator instance from `i18n.Config`",
		"`i18n.Config` | constructor config type for dedicated translators",
		"`i18n.Config.Locales` | embedded locale file set used by `i18n.New(config)`",
		"`i18n.Translator` | interface implemented by dedicated translator instances",
		"`i18n.Translator.T(messageID, data...)` | translate with graceful fallback on a dedicated translator",
		"`i18n.Translator.Te(messageID, data...)` | translate with explicit errors on a dedicated translator",
		"`i18n.ErrUnsupportedLanguage` | sentinel wrapped when a requested language is unsupported",
		"`i18n.ErrMessageIDEmpty` | sentinel returned when `Te` receives an empty message ID",
		"`locales.EmbedLocales` | framework-shipped embedded locale file set",
		"`i18n.T(\"missing_key\")` | logs the translation error and returns `\"missing_key\"` as the fallback",
		"`i18n.Te(\"missing_key\")` | returns `\"\"` and an error wrapping the localization failure",
		"`i18n.Te(\"\")` | returns `\"\"` and `i18n.ErrMessageIDEmpty`",
		"uses only the first `map[string]any` value",
		"`i18n.New(config)` creates an independent `i18n.Translator`",
		"it does not change the process-level global translator",
		"`i18n.Config.Locales` must provide every supported locale file",
		"an empty `embed.FS` or a missing supported file makes `i18n.New(config)` return an error",
		"`i18n.SetLanguage(...)` wraps `i18n.ErrUnsupportedLanguage`",
		"`i18n.Te(\"\")` returns `i18n.ErrMessageIDEmpty` directly",
		"`VEF_I18N_LANGUAGE` is exposed by the `config` package as `config.EnvI18NLanguage`; it is not an `i18n` package API",
		"If the embedded catalogs cannot be loaded, initialization panics",
		"Passing an empty string to `i18n.SetLanguage(\"\")` re-reads `VEF_I18N_LANGUAGE`",
		"When the language changes, the active locale file set is preserved",
		"load the active translator through an atomic pointer",
		"concurrent translation calls and language switches race-free",
	}
}

func chineseDocTerms() []string {
	return []string{
		"12 个 top-level exported symbols、1 个 exported field、2 个 exported methods",
		"public surface fingerprint 是 `" + i18nFingerprint + "`",
		"1 个 top-level exported symbol、0 个 exported fields、0 个 exported methods",
		"public surface fingerprint 是 `" + localesFingerprint + "`",
		"`i18n.GetSupportedLanguages()` 返回 supported-language slice 的副本",
		"`en.json` 和 `zh-CN.json` catalogs 各包含同一组 160 个 message IDs",
		"`i18n.DefaultLanguage` | 默认语言常量（`zh-CN`）",
		"`i18n.GetSupportedLanguages()` | 返回支持的语言代码副本（`zh-CN`、`en`）",
		"`i18n.IsLanguageSupported(code)` | 判断语言代码是否受支持",
		"`i18n.CurrentLanguage()` | 读取当前全局语言代码",
		"`i18n.SetLanguage(code)` | 原子切换进程级全局 translator",
		"`i18n.T(messageID, data...)` | 翻译并在失败时回退到 message ID",
		"`i18n.Te(messageID, data...)` | 翻译并显式返回错误",
		"`i18n.New(config)` | 根据 `i18n.Config` 创建独立 translator 实例",
		"`i18n.Config` | 独立 translator 构造配置类型",
		"`i18n.Config.Locales` | `i18n.New(config)` 使用的嵌入式 locale file set",
		"`i18n.Translator` | 独立 translator 实例实现的接口",
		"`i18n.Translator.T(messageID, data...)` | 在独立 translator 上翻译，并在失败时回退",
		"`i18n.Translator.Te(messageID, data...)` | 在独立 translator 上翻译，并显式返回错误",
		"`i18n.ErrUnsupportedLanguage` | 请求语言不受支持时被包装的 sentinel",
		"`i18n.ErrMessageIDEmpty` | `Te` 收到空 message ID 时返回的 sentinel",
		"`locales.EmbedLocales` | 框架随包提供的嵌入式 locale file set",
		"`i18n.T(\"missing_key\")` | 记录翻译错误，并把 `\"missing_key\"` 作为 fallback 返回",
		"`i18n.Te(\"missing_key\")` | 返回 `\"\"` 和包装后的 localization failure",
		"`i18n.Te(\"\")` | 返回 `\"\"` 和 `i18n.ErrMessageIDEmpty`",
		"实现只使用第一个 `map[string]any`",
		"`i18n.New(config)` 会创建独立的 `i18n.Translator`",
		"并且不会改变包级 `i18n.T(...)` 与 `i18n.Te(...)` 使用的进程级全局 translator",
		"`i18n.Config.Locales` 必须提供所有受支持的 locale 文件",
		"空 `embed.FS` 或缺少某个受支持文件都会让 `i18n.New(config)` 返回 error",
		"`i18n.SetLanguage(...)` 会包装 `i18n.ErrUnsupportedLanguage`",
		"`i18n.Te(\"\")` 会直接返回 `i18n.ErrMessageIDEmpty`",
		"`VEF_I18N_LANGUAGE` 由 `config` 包以 `config.EnvI18NLanguage` 暴露；它不是 `i18n` 包 API",
		"如果内置 catalogs 无法加载，初始化会 panic",
		"传入空字符串 `i18n.SetLanguage(\"\")` 会重新读取 `VEF_I18N_LANGUAGE`",
		"语言切换时会保留当前 active locale file set",
		"通过 atomic pointer 读取当前 translator",
		"翻译调用和语言切换并发时不会产生 data race",
	}
}

func verifySourceContracts(sourceRoot string) []string {
	checks := []struct {
		path  string
		terms []string
	}{
		{
			path: "config/env.go",
			terms: []string{
				"EnvKeyPrefix    = \"VEF\"",
				"EnvI18NLanguage = EnvKeyPrefix + \"_I18N_LANGUAGE\"",
			},
		},
		{
			path: "i18n/i18n.go",
			terms: []string{
				"const DefaultLanguage = \"zh-CN\"",
				"supportedLanguages = []string{\"zh-CN\", \"en\"}",
				"current atomic.Pointer[state]",
				"preferredLanguage := lo.CoalesceOrEmpty(os.Getenv(vefconfig.EnvI18NLanguage), DefaultLanguage)",
				"st, err := newState(locales.EmbedLocales, preferredLanguage)",
				"panic(err)",
				"type Config struct",
				"Locales embed.FS",
				"bundle := i18n.NewBundle(language.SimplifiedChinese)",
				"bundle.RegisterUnmarshalFunc(\"json\", json.Unmarshal)",
				"filename := fmt.Sprintf(\"%s.json\", lang)",
				"bundle.LoadMessageFileFS(localesFS, filename)",
				"return current.Load().translator.T(messageID, templateData...)",
				"return current.Load().translator.Te(messageID, templateData...)",
				"result := make([]string, len(supportedLanguages))",
				"copy(result, supportedLanguages)",
				"return slices.Contains(supportedLanguages, languageCode)",
				"return current.Load().language",
				"if languageCode == \"\"",
				"languageCode = lo.CoalesceOrEmpty(os.Getenv(vefconfig.EnvI18NLanguage), DefaultLanguage)",
				"return fmt.Errorf(\"%w: %s (supported: %v)\", ErrUnsupportedLanguage, languageCode, supportedLanguages)",
				"st, err := newState(current.Load().locales, languageCode)",
				"current.Store(st)",
			},
		},
		{
			path: "i18n/translator.go",
			terms: []string{
				"type Translator interface",
				"T(messageID string, templateData ...map[string]any) string",
				"Te(messageID string, templateData ...map[string]any) (string, error)",
				"message, err := t.Te(messageID, templateData...)",
				"return messageID",
				"if messageID == \"\"",
				"return \"\", ErrMessageIDEmpty",
				"if len(templateData) > 0",
				"data = templateData[0]",
				"MessageID:    messageID",
				"TemplateData: data",
				"return \"\", fmt.Errorf(\"translation failed for messageID %q: %w\", messageID, err)",
				"func New(config Config) (Translator, error)",
				"preferredLanguage := lo.CoalesceOrEmpty(os.Getenv(vefconfig.EnvI18NLanguage), DefaultLanguage)",
				"st, err := newState(config.Locales, preferredLanguage)",
				"return st.translator, nil",
			},
		},
		{
			path: "i18n/errors.go",
			terms: []string{
				"ErrUnsupportedLanguage = errors.New(\"unsupported language code\")",
				"ErrMessageIDEmpty = errors.New(\"messageID cannot be empty\")",
			},
		},
		{
			path: "i18n/locales/locales.go",
			terms: []string{
				"//go:embed *.json",
				"var EmbedLocales embed.FS",
			},
		},
		{
			path: "i18n/i18n_test.go",
			terms: []string{
				"ConfigFieldAccess",
				"ConfigWithEmptyLocales",
				"SetToEmptyStringUsesDefault",
				"SetToUnsupportedLanguage",
				"TestCurrentLanguage",
				"TestSetLanguagePreservesLocales",
				"Should return a copy, not the original slice",
				"TFunctionWithInvalidMessageID",
				"TEFunctionWithInvalidMessageID",
				"TEFunctionWithEmptyMessageID",
				"Interface implementation consistency",
				"Direct translator uses environment/default language, not the globally set language",
			},
		},
	}

	var failures []string
	for _, check := range checks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	return failures
}

func verifyLocaleCatalogs(sourceRoot string) []string {
	en, failures := readLocale(filepath.Join(sourceRoot, "i18n/locales/en.json"))
	zh, zhFailures := readLocale(filepath.Join(sourceRoot, "i18n/locales/zh-CN.json"))
	failures = append(failures, zhFailures...)
	if len(failures) > 0 {
		return failures
	}

	if len(en) != localeMessageCount {
		failures = append(failures, fmt.Sprintf("en locale message count mismatch: got %d want %d", len(en), localeMessageCount))
	}
	if len(zh) != localeMessageCount {
		failures = append(failures, fmt.Sprintf("zh-CN locale message count mismatch: got %d want %d", len(zh), localeMessageCount))
	}

	for key := range en {
		if _, ok := zh[key]; !ok {
			failures = append(failures, "zh-CN locale missing key "+key)
		}
	}
	for key := range zh {
		if _, ok := en[key]; !ok {
			failures = append(failures, "en locale missing key "+key)
		}
	}

	for _, key := range []string{
		"ok",
		"access_denied",
		"validator_phone_number",
		"api_request_resource",
		"api_request_action",
		"api_request_version",
		"unsupported_media_type",
	} {
		if _, ok := en[key]; !ok {
			failures = append(failures, "en locale missing required key "+key)
		}
		if _, ok := zh[key]; !ok {
			failures = append(failures, "zh-CN locale missing required key "+key)
		}
	}

	return failures
}

func readLocale(path string) (map[string]any, []string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, []string{fmt.Sprintf("failed to read locale %s: %v", path, err)}
	}

	var messages map[string]any
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, []string{fmt.Sprintf("failed to parse locale %s: %v", path, err)}
	}

	return messages, nil
}

func verifyAuditArtifacts(sourceRoot, docsRoot string) []string {
	var failures []string
	failures = append(failures, verifyManifest(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))...)
	failures = append(failures, verifyContractLedger(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))...)

	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", filepath.Join(docsRoot, "scripts/api-audit-manifest.json"),
		"-ledger", filepath.Join(docsRoot, "scripts/api-audit-ledger.json"),
		"-contract-ledger", filepath.Join(docsRoot, "scripts/api-contract-ledger.json"),
		"-print-current",
	)
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return append(failures, fmt.Sprintf("verify-api-audit -print-current failed: %v\n%s", err, strings.TrimSpace(string(output))))
	}

	current := string(output)
	for _, term := range []string{
		"\"package\": \"" + i18nPackage + "\"",
		"\"top_level\": 12",
		"\"fields\": 1",
		"\"methods\": 2",
		"\"fingerprint\": \"" + i18nFingerprint + "\"",
		"\"package\": \"" + localesPackage + "\"",
		"\"top_level\": 1",
		"\"fields\": 0",
		"\"methods\": 0",
		"\"fingerprint\": \"" + localesFingerprint + "\"",
	} {
		if !strings.Contains(current, term) {
			failures = append(failures, "current API inventory missing term "+term)
		}
	}

	return failures
}

func verifyManifest(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("failed to read API audit manifest: %v", err)}
	}
	var manifest struct {
		Packages []struct {
			Package     string   `json:"package"`
			Coverage    []string `json:"coverage"`
			TopLevel    int      `json:"top_level"`
			Fields      int      `json:"fields"`
			Methods     int      `json:"methods"`
			Fingerprint string   `json:"fingerprint"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return []string{fmt.Sprintf("failed to parse API audit manifest: %v", err)}
	}

	expected := map[string]struct {
		topLevel    int
		fields      int
		methods     int
		fingerprint string
	}{
		i18nPackage:    {i18nTopLevelCount, i18nFieldCount, i18nMethodCount, i18nFingerprint},
		localesPackage: {localesTopLevelCount, localesFieldCount, localesMethodCount, localesFingerprint},
	}

	seen := map[string]bool{}
	var failures []string
	for _, entry := range manifest.Packages {
		want, ok := expected[entry.Package]
		if !ok {
			continue
		}
		seen[entry.Package] = true
		if entry.TopLevel != want.topLevel || entry.Fields != want.fields ||
			entry.Methods != want.methods || entry.Fingerprint != want.fingerprint {
			failures = append(failures, fmt.Sprintf("manifest surface mismatch for %s", entry.Package))
		}
		if !contains(entry.Coverage, "docs/features/i18n.md") {
			failures = append(failures, fmt.Sprintf("manifest coverage for %s does not include docs/features/i18n.md", entry.Package))
		}
	}
	for pkg := range expected {
		if !seen[pkg] {
			failures = append(failures, "manifest missing package "+pkg)
		}
	}

	return failures
}

func verifyContractLedger(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("failed to read API contract ledger: %v", err)}
	}
	var ledger struct {
		PackageReviews []struct {
			Package         string `json:"package"`
			ReviewedSurface struct {
				TopLevel    int    `json:"top_level"`
				Fields      int    `json:"fields"`
				Methods     int    `json:"methods"`
				EntryCount  int    `json:"entry_count"`
				Fingerprint string `json:"fingerprint"`
			} `json:"reviewed_surface"`
			Coverage    []string `json:"coverage"`
			ContractIDs []string `json:"contract_ids"`
		} `json:"package_reviews"`
		Entries []struct {
			ID       string   `json:"id"`
			Package  string   `json:"package"`
			Coverage []string `json:"coverage"`
			Terms    []string `json:"terms"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(data, &ledger); err != nil {
		return []string{fmt.Sprintf("failed to parse API contract ledger: %v", err)}
	}

	var failures []string
	reviews := map[string]bool{}
	for _, review := range ledger.PackageReviews {
		switch review.Package {
		case i18nPackage:
			reviews[review.Package] = true
			if review.ReviewedSurface.TopLevel != i18nTopLevelCount ||
				review.ReviewedSurface.Fields != i18nFieldCount ||
				review.ReviewedSurface.Methods != i18nMethodCount ||
				review.ReviewedSurface.EntryCount != i18nEntryCount ||
				review.ReviewedSurface.Fingerprint != i18nFingerprint {
				failures = append(failures, "contract ledger reviewed surface mismatch for "+i18nPackage)
			}
			if !contains(review.Coverage, "docs/features/i18n.md") {
				failures = append(failures, "contract ledger coverage missing docs/features/i18n.md for "+i18nPackage)
			}
			if !contains(review.ContractIDs, i18nPackage+"#runtime-contract:language-selection-and-fallback") {
				failures = append(failures, "contract ledger missing i18n language-selection contract id")
			}
		case localesPackage:
			reviews[review.Package] = true
			if review.ReviewedSurface.TopLevel != localesTopLevelCount ||
				review.ReviewedSurface.Fields != localesFieldCount ||
				review.ReviewedSurface.Methods != localesMethodCount ||
				review.ReviewedSurface.EntryCount != localesEntryCount ||
				review.ReviewedSurface.Fingerprint != localesFingerprint {
				failures = append(failures, "contract ledger reviewed surface mismatch for "+localesPackage)
			}
			if !contains(review.Coverage, "docs/features/i18n.md") {
				failures = append(failures, "contract ledger coverage missing docs/features/i18n.md for "+localesPackage)
			}
			if !contains(review.ContractIDs, localesPackage+"#runtime-contract:embedded-locale-files") {
				failures = append(failures, "contract ledger missing locales embedded-file contract id")
			}
		}
	}
	for _, pkg := range []string{i18nPackage, localesPackage} {
		if !reviews[pkg] {
			failures = append(failures, "contract ledger missing package review for "+pkg)
		}
	}

	requiredEntries := map[string][]string{
		i18nPackage + "#runtime-contract:language-selection-and-fallback": {
			"zh-CN", "en", "VEF_I18N_LANGUAGE", "SetLanguage", "message ID",
		},
		localesPackage + "#runtime-contract:embedded-locale-files": {
			"locales.EmbedLocales", "*.json", "zh-CN", "en",
		},
	}
	seenEntries := map[string]bool{}
	for _, entry := range ledger.Entries {
		terms, ok := requiredEntries[entry.ID]
		if !ok {
			continue
		}
		seenEntries[entry.ID] = true
		if !contains(entry.Coverage, "docs/features/i18n.md") {
			failures = append(failures, "contract ledger entry coverage missing docs/features/i18n.md for "+entry.ID)
		}
		for _, term := range terms {
			if !contains(entry.Terms, term) {
				failures = append(failures, fmt.Sprintf("contract ledger entry %s missing term %q", entry.ID, term))
			}
		}
	}
	for id := range requiredEntries {
		if !seenEntries[id] {
			failures = append(failures, "contract ledger missing entry "+id)
		}
	}

	return failures
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./i18n", "./i18n/locales")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./i18n ./i18n/locales failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func exportedPackageSurface(dir string) packageSurface {
	fset := token.NewFileSet()
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(fmt.Errorf("failed to read source dir %s: %w", dir, err))
	}

	surface := packageSurface{}
	fieldSet := map[string]bool{}
	methodSet := map[string]bool{}

	for _, entry := range files {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			panic(fmt.Errorf("failed to parse %s: %w", filepath.Join(dir, name), err))
		}

		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				switch d.Tok {
				case token.CONST:
					for _, spec := range d.Specs {
						valueSpec := spec.(*ast.ValueSpec)
						for _, name := range valueSpec.Names {
							if name.IsExported() {
								surface.consts = append(surface.consts, name.Name)
							}
						}
					}
				case token.VAR:
					for _, spec := range d.Specs {
						valueSpec := spec.(*ast.ValueSpec)
						for _, name := range valueSpec.Names {
							if name.IsExported() {
								surface.vars = append(surface.vars, name.Name)
							}
						}
					}
				case token.TYPE:
					for _, spec := range d.Specs {
						typeSpec := spec.(*ast.TypeSpec)
						if !typeSpec.Name.IsExported() {
							continue
						}
						typeName := typeSpec.Name.Name
						surface.types = append(surface.types, typeName)

						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							for _, field := range structType.Fields.List {
								for _, fieldName := range exportedFieldNames(field) {
									fieldSet[typeName+"."+fieldName] = true
								}
							}
						}
						if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							for _, method := range interfaceType.Methods.List {
								for _, methodName := range exportedFieldNames(method) {
									methodSet[typeName+"."+methodName] = true
								}
							}
						}
					}
				}
			case *ast.FuncDecl:
				if !d.Name.IsExported() {
					continue
				}
				if d.Recv == nil {
					surface.funcs = append(surface.funcs, d.Name.Name)
					continue
				}
				receiver := receiverName(d.Recv.List[0].Type)
				if receiver != "" && ast.IsExported(receiver) {
					methodSet[receiver+"."+d.Name.Name] = true
				}
			}
		}
	}

	surface.fields = sortedKeys(fieldSet)
	surface.methods = sortedKeys(methodSet)
	sort.Strings(surface.consts)
	sort.Strings(surface.funcs)
	sort.Strings(surface.types)
	sort.Strings(surface.vars)

	return surface
}

func exportedFieldNames(field *ast.Field) []string {
	var result []string
	if len(field.Names) == 0 {
		name := receiverName(field.Type)
		if name != "" && ast.IsExported(name) {
			result = append(result, name)
		}
		return result
	}

	for _, name := range field.Names {
		if name.IsExported() {
			result = append(result, name.Name)
		}
	}

	return result
}

func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverName(t.X)
	case *ast.IndexExpr:
		return receiverName(t.X)
	case *ast.IndexListExpr:
		return receiverName(t.X)
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func compareSurfaceCounts(label string, surface packageSurface, topLevel, fields, methods, entries int) []string {
	actualTopLevel := len(surface.consts) + len(surface.funcs) + len(surface.types) + len(surface.vars)
	actualEntries := actualTopLevel + len(surface.fields) + len(surface.methods)
	var failures []string
	if actualTopLevel != topLevel {
		failures = append(failures, fmt.Sprintf("%s top-level count mismatch: got %d want %d", label, actualTopLevel, topLevel))
	}
	if len(surface.fields) != fields {
		failures = append(failures, fmt.Sprintf("%s exported field count mismatch: got %d want %d", label, len(surface.fields), fields))
	}
	if len(surface.methods) != methods {
		failures = append(failures, fmt.Sprintf("%s exported method count mismatch: got %d want %d", label, len(surface.methods), methods))
	}
	if actualEntries != entries {
		failures = append(failures, fmt.Sprintf("%s entry count mismatch: got %d want %d", label, actualEntries, entries))
	}

	return failures
}

func verifyI18nReferences(doc corpus, surface packageSurface) []string {
	allowed := map[string]bool{}
	for _, group := range [][]string{surface.consts, surface.funcs, surface.types, surface.vars} {
		for _, name := range group {
			allowed[name] = true
		}
	}

	return verifyPackageReferences(doc, "i18n", allowed)
}

func verifyLocalesReferences(doc corpus, surface packageSurface) []string {
	allowed := map[string]bool{}
	for _, group := range [][]string{surface.consts, surface.funcs, surface.types, surface.vars} {
		for _, name := range group {
			allowed[name] = true
		}
	}

	return verifyPackageReferences(doc, "locales", allowed)
}

func verifyPackageReferences(doc corpus, pkg string, allowed map[string]bool) []string {
	seen := make(map[string]bool)
	re := regexp.MustCompile(regexp.QuoteMeta(pkg) + `\.([A-Z][A-Za-z0-9_]*)`)
	for _, match := range re.FindAllStringSubmatch(doc.content, -1) {
		seen[match[1]] = true
	}

	var failures []string
	for name := range seen {
		if !allowed[name] {
			failures = append(failures, fmt.Sprintf("%s references unknown %s public symbol %s.%s", doc.label, pkg, pkg, name))
		}
	}

	return failures
}

func verifySurfaceMentioned(doc corpus, pkg string, surface packageSurface) []string {
	var failures []string
	for _, name := range topLevelNames(surface) {
		if !containsTerm(doc.content, pkg+"."+name) {
			failures = append(failures, fmt.Sprintf("%s missing exported %s symbol %s.%s", doc.label, pkg, pkg, name))
		}
	}

	for _, field := range surface.fields {
		if !containsTerm(doc.content, pkg+"."+field) {
			failures = append(failures, fmt.Sprintf("%s missing exported %s field %s.%s", doc.label, pkg, pkg, field))
		}
	}

	for _, method := range surface.methods {
		if !containsTerm(doc.content, pkg+"."+method) {
			failures = append(failures, fmt.Sprintf("%s missing exported %s method %s.%s", doc.label, pkg, pkg, method))
		}
	}

	return failures
}

func topLevelNames(surface packageSurface) []string {
	names := append([]string{}, surface.consts...)
	names = append(names, surface.funcs...)
	names = append(names, surface.types...)
	names = append(names, surface.vars...)
	sort.Strings(names)

	return names
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

func readCorpus(label, path string) corpus {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s (%s): %w", label, path, err))
	}

	return corpus{label: label, content: string(data)}
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
