package main

import (
	"bytes"
	"flag"
	"fmt"
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

	jsSource := readCorpus("js package source", filepath.Join(sourceRoot, "js/js.go"))
	libsSource := readCorpus("js library source", filepath.Join(sourceRoot, "js/libs.go"))
	sourceSetup := corpus{
		label:   "js setup source",
		content: jsSource.content + "\n" + libsSource.content,
	}
	englishDocs := readCorpus("English JS docs", filepath.Join(docsRoot, "docs/features/js-engine.md"))
	chineseDocs := readCorpus("Chinese JS docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/js-engine.md"))
	docs := []corpus{englishDocs, chineseDocs}

	checks := []struct {
		corpus corpus
		terms  []string
	}{
		{
			corpus: jsSource,
			terms: []string{
				"Runtime    = goja.Runtime",
				"Value      = goja.Value",
				"Object     = goja.Object",
				"Program    = goja.Program",
				"AstProgram = ast.Program",
				"Compile     = goja.Compile",
				"MustCompile = goja.MustCompile",
				"IsNaN       = goja.IsNaN",
				"IsString    = goja.IsString",
				"IsBigInt    = goja.IsBigInt",
				"IsNumber    = goja.IsNumber",
				"IsInfinity  = goja.IsInfinity",
				"IsUndefined = goja.IsUndefined",
				"IsNull      = goja.IsNull",
				"vm := goja.New()",
				"vm.SetParserOptions(parser.WithDisableSourceMaps)",
				"vm.SetFieldNameMapper(goja.TagFieldNameMapper(\"json\", true))",
				"if _, err := vm.RunProgram(lib); err != nil",
				"return nil, err",
				"return goja.Parse(name, src, parser.WithDisableSourceMaps)",
			},
		},
		{
			corpus: libsSource,
			terms: []string{
				"//go:embed libs/day.v1_11_19.js",
				"//go:embed libs/big.v7_0_1.js",
				"//go:embed libs/utils.v12_7_0.js",
				"//go:embed libs/validator.v13_15_20.js",
				"compiledDayJs       = MustCompile(\"day.js\", string(dayJs), true)",
				"compiledBigJs       = MustCompile(\"big.js\", string(bigJs), true)",
				"compiledUtilsJs     = MustCompile(\"utils.js\", string(utilsJs), true)",
				"compiledValidatorJs = MustCompile(\"validator.js\", string(validatorJs), true)",
			},
		},
	}

	docTerms := []string{
		"goja.New()",
		"parser.WithDisableSourceMaps",
		"goja.TagFieldNameMapper(\"json\", true)",
		"dayjs",
		"Big",
		"utils",
		"validator",
		"libs/day.v1_11_19.js",
		"libs/big.v7_0_1.js",
		"libs/utils.v12_7_0.js",
		"libs/validator.v13_15_20.js",
		"Day.js 1.11.19",
		"big.js 7.0.1",
		"utils 12.7.0",
		"validator.js 13.15.20",
		"Node-style module loader",
		"require.NewRegistry",
		"console",
		"fs",
		"net",
		"timers",
		"vm.Set(...)",
		"Runtime.Interrupt(...)",
		"Runtime.ClearInterrupt()",
		"v0.0.0-20260311135729-065cd970411c",
		"js.Compile",
		"js.MustCompile",
		"js.Is*",
		"public API index",
	}
	englishDocTerms := []string{
		"pre-compiled browser bundles executed in this order",
		"returns `nil` and the first load error",
		"does not register native modules",
		"does not enable a `console` shim",
		"does not install Node APIs",
		"not a sandbox policy",
		"do not inherit `js.New()` runtime parser options",
		"VEF does not wrap individual library functions",
	}
	chineseDocTerms := []string{
		"按顺序执行预编译的 browser bundles",
		"返回 `nil` 和第一",
		"不会注册 native modules",
		"不会启用 `console` shim",
		"不会安装",
		"不是 sandbox policy",
		"runtime parser options",
		"不会逐个包装 library functions",
	}

	var failures []string
	for _, check := range checks {
		failures = append(failures, missingTerms(check.corpus, check.terms)...)
	}
	failures = append(failures, missingOrderedTerms(jsSource, []string{
		"compiledDayJs",
		"compiledBigJs",
		"compiledUtilsJs",
		"compiledValidatorJs",
	})...)
	failures = append(failures, forbiddenTerms(sourceSetup, []string{
		"require.NewRegistry",
		"RegisterNativeModule",
		"console.Enable",
		"vm.Set(\"",
	})...)

	for _, doc := range docs {
		failures = append(failures, missingTerms(doc, docTerms)...)
	}
	failures = append(failures, missingTerms(englishDocs, englishDocTerms)...)
	failures = append(failures, missingTerms(chineseDocs, chineseDocTerms)...)
	failures = append(failures, runGoTest(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("JS contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Println("JS contract docs verified: 2 source files, 2 doc mirrors, go test ./js")
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
		if !containsTerm(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term: %s", c.label, term))
		}
	}

	return failures
}

func missingOrderedTerms(c corpus, terms []string) []string {
	var failures []string
	offset := 0
	for _, term := range terms {
		idx := strings.Index(c.content[offset:], term)
		if idx < 0 {
			failures = append(failures, fmt.Sprintf("%s missing ordered term after offset %d: %s", c.label, offset, term))
			continue
		}
		offset += idx + len(term)
	}

	return failures
}

func forbiddenTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if strings.Contains(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s contains forbidden setup term: %s", c.label, term))
		}
	}

	return failures
}

func containsTerm(content, term string) bool {
	if strings.Contains(content, term) {
		return true
	}

	return strings.Contains(normalizeWhitespace(content), normalizeWhitespace(term))
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func runGoTest(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./js")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./js failed: %v\n%s", err, strings.TrimSpace(output.String()))}
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
