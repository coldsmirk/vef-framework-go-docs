package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	logxPackage     = "github.com/coldsmirk/vef-framework-go/logx"
	logxFingerprint = "4ff9c19b53d9985911e2985c2763802337d6a69f783962bb73b9ee7424481eaf"
	logxTopLevel    = 7
	logxFields      = 0
	logxMethods     = 15
	logxEntries     = 22

	englishSmallUtilitiesPath = "docs/utilities/small-utilities.md"
	chineseSmallUtilitiesPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/small-utilities.md"
	englishExtensionPath      = "docs/reference/extension-points.md"
	chineseExtensionPath      = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/extension-points.md"
	englishIndexPath          = "docs/reference/public-api-index.md"
	chineseIndexPath          = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
)

type corpus struct {
	label   string
	path    string
	content string
}

type auditLedger struct {
	Entries []auditEntry `json:"entries"`
}

type auditEntry struct {
	ID        string   `json:"id"`
	Package   string   `json:"package"`
	Kind      string   `json:"kind"`
	Symbol    string   `json:"symbol"`
	Signature string   `json:"signature"`
	Coverage  []string `json:"coverage"`
}

type manifest struct {
	Packages []manifestEntry `json:"packages"`
}

type manifestEntry struct {
	Package     string   `json:"package"`
	Tier        string   `json:"tier"`
	Coverage    []string `json:"coverage"`
	TopLevel    int      `json:"top_level"`
	Fields      int      `json:"fields"`
	Methods     int      `json:"methods"`
	Fingerprint string   `json:"fingerprint"`
}

type contractLedger struct {
	PackageReviews []contractPackageReview `json:"package_reviews"`
	Entries        []contractEntry         `json:"entries"`
}

type contractPackageReview struct {
	Package         string        `json:"package"`
	Disposition     string        `json:"disposition"`
	ReviewedSurface reviewSurface `json:"reviewed_surface"`
	Coverage        []string      `json:"coverage"`
	SourceEvidence  []string      `json:"source_evidence"`
	ContractIDs     []string      `json:"contract_ids"`
}

type reviewSurface struct {
	TopLevel    int    `json:"top_level"`
	Fields      int    `json:"fields"`
	Methods     int    `json:"methods"`
	EntryCount  int    `json:"entry_count"`
	Fingerprint string `json:"fingerprint"`
}

type contractEntry struct {
	ID             string   `json:"id"`
	Package        string   `json:"package"`
	Kind           string   `json:"kind"`
	Contract       string   `json:"contract"`
	Disposition    string   `json:"disposition"`
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	TestEvidence   []string `json:"test_evidence"`
	Terms          []string `json:"terms"`
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishSmallUtilities := readCorpus("English small utilities docs", docsRoot, englishSmallUtilitiesPath)
	chineseSmallUtilities := readCorpus("Chinese small utilities docs", docsRoot, chineseSmallUtilitiesPath)
	englishExtension := readCorpus("English extension points docs", docsRoot, englishExtensionPath)
	chineseExtension := readCorpus("Chinese extension points docs", docsRoot, chineseExtensionPath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadLogxAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadLogxManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadLogxContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveLogxEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
		failures = append(failures, verifyRootNamedLoggerIndex(index)...)
	}
	for _, doc := range []corpus{englishSmallUtilities, chineseSmallUtilities} {
		failures = append(failures, verifyDocumentedSurface(doc, entries)...)
		failures = append(failures, verifyNoPhantomLogxRefs(doc, entries)...)
		failures = append(failures, verifyNoPhantomLoggerMembers(doc, entries)...)
		failures = append(failures, missingTerms(doc, smallUtilitiesTerms())...)
	}
	for _, doc := range []corpus{englishExtension, chineseExtension} {
		failures = append(failures, verifyNoPhantomLogxRefs(doc, entries)...)
		failures = append(failures, verifyExtensionTerms(doc)...)
	}
	failures = append(failures, verifyNoPhantomLogxRefsInAllDocs(docsRoot, entries)...)

	failures = append(failures, verifyContractLedger(review, contract, sourceRoot)...)
	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("logx contracts verified")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != logxPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != logxTopLevel || entry.Fields != logxFields ||
		entry.Methods != logxMethods || entry.Fingerprint != logxFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			logxTopLevel, logxFields, logxMethods, logxFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != logxPackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != logxTopLevel ||
		review.ReviewedSurface.Fields != logxFields ||
		review.ReviewedSurface.Methods != logxMethods ||
		review.ReviewedSurface.EntryCount != logxEntries ||
		review.ReviewedSurface.Fingerprint != logxFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, logxContractID()) {
		failures = append(failures, "contract package review missing logx contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != logxEntries {
		failures = append(failures, fmt.Sprintf("logx audit entry count mismatch: got %d want %d", len(entries), logxEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != logxPackage {
			failures = append(failures, "non-logx audit entry passed into logx verifier: "+entry.ID)
		}
		if strings.Contains(entry.Package, "/internal/") {
			failures = append(failures, "internal package included in logx verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate logx audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "logx audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != logxTopLevel || counts["field"] != logxFields || counts["method"] != logxMethods {
		failures = append(failures, fmt.Sprintf("logx audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	failures = append(failures, verifyExpectedEntries(entries)...)

	return failures
}

func verifyExpectedEntries(entries []auditEntry) []string {
	expected := map[string]string{
		"Level":                 "Level : github.com/coldsmirk/vef-framework-go/logx.Level",
		"LevelDebug":            "LevelDebug : github.com/coldsmirk/vef-framework-go/logx.Level = 1",
		"LevelInfo":             "LevelInfo : github.com/coldsmirk/vef-framework-go/logx.Level = 2",
		"LevelWarn":             "LevelWarn : github.com/coldsmirk/vef-framework-go/logx.Level = 3",
		"LevelError":            "LevelError : github.com/coldsmirk/vef-framework-go/logx.Level = 4",
		"LevelPanic":            "LevelPanic : github.com/coldsmirk/vef-framework-go/logx.Level = 5",
		"Logger":                "Logger : github.com/coldsmirk/vef-framework-go/logx.Logger",
		"Level.String":          "String : func() string",
		"Logger.Named":          "Named : func(name string) github.com/coldsmirk/vef-framework-go/logx.Logger",
		"Logger.WithCallerSkip": "WithCallerSkip : func(skip int) github.com/coldsmirk/vef-framework-go/logx.Logger",
		"Logger.Enabled":        "Enabled : func(level github.com/coldsmirk/vef-framework-go/logx.Level) bool",
		"Logger.Sync":           "Sync : func()",
		"Logger.Debug":          "Debug : func(message string)",
		"Logger.Debugf":         "Debugf : func(template string, args ...any)",
		"Logger.Info":           "Info : func(message string)",
		"Logger.Infof":          "Infof : func(template string, args ...any)",
		"Logger.Warn":           "Warn : func(message string)",
		"Logger.Warnf":          "Warnf : func(template string, args ...any)",
		"Logger.Error":          "Error : func(message string)",
		"Logger.Errorf":         "Errorf : func(template string, args ...any)",
		"Logger.Panic":          "Panic : func(message string)",
		"Logger.Panicf":         "Panicf : func(template string, args ...any)",
	}

	seen := map[string]string{}
	for _, entry := range entries {
		seen[entry.Symbol] = entry.Signature
	}

	var failures []string
	for symbol, signature := range expected {
		got, ok := seen[symbol]
		if !ok {
			failures = append(failures, "logx audit entries missing "+symbol)
			continue
		}
		if got != signature {
			failures = append(failures, fmt.Sprintf("logx audit signature mismatch for %s: got %q want %q", symbol, got, signature))
		}
	}

	return failures
}

func verifyCoverage(
	entries []auditEntry,
	manifestEntry manifestEntry,
	review contractPackageReview,
	contract contractEntry,
) []string {
	var failures []string
	expected := []string{englishSmallUtilitiesPath, englishExtensionPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest logx coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review logx coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry logx coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, logxPackage)
	if section == "" {
		return []string{index.label + " missing logx package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s logx index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}
	if strings.Contains(section, "internal/logx") {
		failures = append(failures, index.label+" logx index includes internal/logx")
	}

	return failures
}

func verifyRootNamedLoggerIndex(index corpus) []string {
	rootSection := packageSection(index.content, "github.com/coldsmirk/vef-framework-go")
	if rootSection == "" {
		return []string{index.label + " missing root package section"}
	}
	want := "FUNC NamedLogger : func(name string) github.com/coldsmirk/vef-framework-go/logx.Logger"
	if !strings.Contains(rootSection, want) {
		return []string{index.label + " root package index missing NamedLogger signature"}
	}

	return nil
}

func verifyDocumentedSurface(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		ref := entry.Symbol
		if entry.Kind == "top" {
			ref = "logx." + entry.Symbol
		}
		if !hasCodeReference(doc.content, ref) {
			failures = append(failures, fmt.Sprintf("%s missing audited logx entry `%s`", doc.label, ref))
		}
	}

	return failures
}

func verifyNoPhantomLogxRefs(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		if entry.Kind == "top" {
			valid["logx."+entry.Symbol] = true
		}
	}

	var failures []string
	for _, ref := range logxReferencesFromMarkdownCode(doc.content) {
		if valid[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown logx public API: %s", doc.label, ref))
	}
	if strings.Contains(doc.content, "internal/logx") {
		failures = append(failures, doc.label+" documents internal/logx as user-facing API")
	}

	return failures
}

func verifyNoPhantomLoggerMembers(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		if entry.Kind == "method" {
			valid[entry.Symbol] = true
		}
	}

	var failures []string
	for _, ref := range loggerMemberReferencesFromMarkdownCode(doc.content) {
		if valid[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown Logger/Level member: %s", doc.label, ref))
	}

	return failures
}

func verifyNoPhantomLogxRefsInAllDocs(docsRoot string, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		if entry.Kind == "top" {
			valid["logx."+entry.Symbol] = true
		}
	}

	var failures []string
	for _, root := range []string{
		filepath.Join(docsRoot, "docs"),
		filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current"),
	} {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(docsRoot, path)
			if err != nil {
				rel = path
			}
			doc := corpus{
				label:   filepath.ToSlash(rel),
				path:    filepath.ToSlash(rel),
				content: string(data),
			}
			for _, ref := range logxReferencesFromMarkdownCode(doc.content) {
				if !valid[ref] {
					failures = append(failures, fmt.Sprintf("%s references unknown logx public API: %s", doc.label, ref))
				}
			}
			if strings.Contains(doc.content, "internal/logx") {
				failures = append(failures, doc.label+" documents internal/logx as user-facing API")
			}

			return nil
		})
		if err != nil {
			failures = append(failures, "failed to scan docs for logx refs: "+err.Error())
		}
	}

	return failures
}

func hasCodeReference(content, ref string) bool {
	return strings.Contains(content, "`"+ref+"`") ||
		strings.Contains(content, "`"+ref+"(") ||
		strings.Contains(content, "`"+ref+"()") ||
		strings.Contains(content, "`"+ref+"[") ||
		codePartsContainReference(content, ref)
}

func smallUtilitiesTerms() []string {
	return []string{
		logxFingerprint,
		"7 exported top-level",
		"15 exported",
		"no exported fields",
		"`logx.Level`",
		"`logx.LevelDebug = 1`",
		"`logx.LevelInfo = 2`",
		"`logx.LevelWarn = 3`",
		"`logx.LevelError = 4`",
		"`logx.LevelPanic = 5`",
		"`logx.Logger`",
		"`Level.String() string`",
		"`unknown`",
		"zero value",
		"`Logger.Named(name string) logx.Logger`",
		"`Logger.WithCallerSkip(skip int) logx.Logger`",
		"`Logger.Enabled(level logx.Level) bool`",
		"`Logger.Sync()`",
		"`Logger.Debug(message string)`",
		"`Logger.Debugf(template string, args ...any)`",
		"`Logger.Info(message string)`",
		"`Logger.Infof(template string, args ...any)`",
		"`Logger.Warn(message string)`",
		"`Logger.Warnf(template string, args ...any)`",
		"`Logger.Error(message string)`",
		"`Logger.Errorf(template string, args ...any)`",
		"`Logger.Panic(message string)`",
		"`Logger.Panicf(template string, args ...any)`",
		"`vef.NamedLogger(name string) logx.Logger`",
	}
}

func verifyExtensionTerms(doc corpus) []string {
	return missingTerms(doc, []string{
		"`vef.NamedLogger(name)`",
		"`logx.Logger`",
		"`Level`",
		"`Level.String()`",
		"`Logger`",
	})
}

func verifyContractLedger(review contractPackageReview, contract contractEntry, sourceRoot string) []string {
	var failures []string
	expectedCoverage := []string{englishSmallUtilitiesPath, englishExtensionPath}
	if contract.ID != logxContractID() {
		failures = append(failures, fmt.Sprintf("logx contract id mismatch: got %q", contract.ID))
	}
	if contract.Package != logxPackage {
		failures = append(failures, fmt.Sprintf("logx contract package mismatch: got %q", contract.Package))
	}
	if contract.Kind != "runtime-contract" {
		failures = append(failures, fmt.Sprintf("logx contract kind mismatch: got %q", contract.Kind))
	}
	if contract.Disposition != "documented:semantic-contract" {
		failures = append(failures, fmt.Sprintf("logx contract disposition mismatch: got %q", contract.Disposition))
	}
	if !sameSet(contract.Coverage, expectedCoverage) {
		failures = append(failures, fmt.Sprintf("logx contract coverage mismatch: got %v want %v", contract.Coverage, expectedCoverage))
	}
	for _, term := range []string{
		"logx.Level",
		"logx.LevelDebug = 1",
		"logx.LevelInfo = 2",
		"logx.LevelWarn = 3",
		"logx.LevelError = 4",
		"logx.LevelPanic = 5",
		"logx.Logger",
		"Level.String()",
		"unknown",
		"Logger.Named(name string) logx.Logger",
		"Logger.Sync()",
		"Logger.Debugf(template string, args ...any)",
		"Logger.Panicf(template string, args ...any)",
		"vef.NamedLogger(name string) logx.Logger",
	} {
		if !contains(contract.Terms, term) {
			failures = append(failures, "logx contract terms missing "+term)
		}
	}

	allEvidence := append([]string{}, review.SourceEvidence...)
	allEvidence = append(allEvidence, contract.SourceEvidence...)
	allEvidence = append(allEvidence, contract.TestEvidence...)
	for _, item := range allEvidence {
		path, lineText, ok := strings.Cut(item, ":")
		if !ok || lineText == "" {
			failures = append(failures, "logx contract evidence missing line number: "+item)
			continue
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, "logx contract evidence missing file: "+item)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	files := map[string]string{
		"logx/logger.go":                readSourceFile(sourceRoot, "logx/logger.go").content,
		"logx/logger_test.go":           readSourceFile(sourceRoot, "logx/logger_test.go").content,
		"log.go":                        readSourceFile(sourceRoot, "log.go").content,
		"internal/logx/logger.go":       readSourceFile(sourceRoot, "internal/logx/logger.go").content,
		"internal/logx/wrapper.go":      readSourceFile(sourceRoot, "internal/logx/wrapper.go").content,
		"internal/logx/discard_test.go": readSourceFile(sourceRoot, "internal/logx/discard_test.go").content,
	}

	checks := []struct {
		file string
		term string
	}{
		{"logx/logger.go", "type Level int8"},
		{"logx/logger.go", "LevelDebug Level = iota + 1"},
		{"logx/logger.go", "LevelInfo"},
		{"logx/logger.go", "LevelWarn"},
		{"logx/logger.go", "LevelError"},
		{"logx/logger.go", "LevelPanic"},
		{"logx/logger.go", `func (l Level) String() string`},
		{"logx/logger.go", `case LevelDebug:`},
		{"logx/logger.go", `return "debug"`},
		{"logx/logger.go", `case LevelInfo:`},
		{"logx/logger.go", `return "info"`},
		{"logx/logger.go", `case LevelWarn:`},
		{"logx/logger.go", `return "warn"`},
		{"logx/logger.go", `case LevelError:`},
		{"logx/logger.go", `return "error"`},
		{"logx/logger.go", `case LevelPanic:`},
		{"logx/logger.go", `return "panic"`},
		{"logx/logger.go", `return "unknown"`},
		{"logx/logger.go", `type Logger interface {`},
		{"logx/logger.go", `Named(name string) Logger`},
		{"logx/logger.go", `WithCallerSkip(skip int) Logger`},
		{"logx/logger.go", `Enabled(level Level) bool`},
		{"logx/logger.go", `Sync()`},
		{"logx/logger.go", `Debug(message string)`},
		{"logx/logger.go", `Debugf(template string, args ...any)`},
		{"logx/logger.go", `Info(message string)`},
		{"logx/logger.go", `Infof(template string, args ...any)`},
		{"logx/logger.go", `Warn(message string)`},
		{"logx/logger.go", `Warnf(template string, args ...any)`},
		{"logx/logger.go", `Error(message string)`},
		{"logx/logger.go", `Errorf(template string, args ...any)`},
		{"logx/logger.go", `Panic(message string)`},
		{"logx/logger.go", `Panicf(template string, args ...any)`},
		{"logx/logger_test.go", `{"Debug", LevelDebug, "debug"}`},
		{"logx/logger_test.go", `{"Info", LevelInfo, "info"}`},
		{"logx/logger_test.go", `{"Warn", LevelWarn, "warn"}`},
		{"logx/logger_test.go", `{"Error", LevelError, "error"}`},
		{"logx/logger_test.go", `{"Panic", LevelPanic, "panic"}`},
		{"logx/logger_test.go", `{"Unknown", Level(0), "unknown"}`},
		{"logx/logger_test.go", `{"OutOfRange", Level(127), "unknown"}`},
		{"log.go", `func NamedLogger(name string) logx.Logger`},
		{"log.go", `return ilogx.Named(name)`},
		{"internal/logx/logger.go", `func Named(name string) logx.Logger`},
		{"internal/logx/wrapper.go", `func (l *zapLogger) Panic(message string)`},
		{"internal/logx/wrapper.go", `l.logger.Panic(message)`},
		{"internal/logx/wrapper.go", `func (l *zapLogger) Panicf(template string, args ...any)`},
		{"internal/logx/wrapper.go", `l.logger.Panicf(template, args...)`},
		{"internal/logx/discard_test.go", `assert.PanicsWithValue(t, "boom", func() {`},
		{"internal/logx/discard_test.go", `logger.Panic("boom")`},
		{"internal/logx/discard_test.go", `assert.PanicsWithValue(t, "boom 42", func() {`},
		{"internal/logx/discard_test.go", `logger.Panicf("boom %d", 42)`},
	}

	var failures []string
	for _, check := range checks {
		if !strings.Contains(files[check.file], check.term) {
			failures = append(failures, fmt.Sprintf("%s missing source contract term: %s", check.file, check.term))
		}
	}

	return failures
}

func runSourceTests(sourceRoot string) []string {
	return runCommand(sourceRoot, "go", "test", "./logx", "./internal/logx")
}

func loadLogxAuditEntries(path string) []auditEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger auditLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	var entries []auditEntry
	for _, entry := range ledger.Entries {
		if entry.Package == logxPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})
	if len(entries) == 0 {
		panic("API audit ledger missing logx entries")
	}

	return entries
}

func loadLogxManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == logxPackage {
			return entry
		}
	}

	panic("API audit manifest missing logx package")
}

func loadLogxContract(path string) (contractPackageReview, contractEntry) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger contractLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	var review contractPackageReview
	reviewFound := false
	for _, item := range ledger.PackageReviews {
		if item.Package == logxPackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing logx package review")
	}

	var contract contractEntry
	contractFound := false
	for _, item := range ledger.Entries {
		if item.ID == logxContractID() {
			contract = item
			contractFound = true

			break
		}
	}
	if !contractFound {
		panic("API contract ledger missing logx contract entry")
	}

	return review, contract
}

func loadLiveLogxEntry(sourceRoot, docsRoot string) manifestEntry {
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
		panic(fmt.Errorf("verify-api-audit -print-current failed: %w\n%s", err, strings.TrimSpace(string(output))))
	}

	var entries []manifestEntry
	payload := "[" + strings.TrimSpace(string(output)) + "]"
	if err := json.Unmarshal([]byte(payload), &entries); err != nil {
		panic(fmt.Errorf("failed to parse verify-api-audit -print-current output: %w", err))
	}
	for _, entry := range entries {
		if entry.Package == logxPackage {
			return entry
		}
	}

	panic("live API inventory missing logx package")
}

func packageSection(content, pkg string) string {
	marker := "## " + pkg
	start := strings.Index(content, marker)
	if start < 0 {
		return ""
	}
	rest := content[start:]
	next := strings.Index(rest[len(marker):], "\n## ")
	if next < 0 {
		return rest
	}

	return rest[:len(marker)+next]
}

func logxReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`logx\.[A-Z][A-Za-z0-9_]*`)
	seen := map[string]bool{}
	for _, part := range codeParts {
		for _, ref := range re.FindAllString(part, -1) {
			seen[ref] = true
		}
	}

	refs := make([]string, 0, len(seen))
	for ref := range seen {
		refs = append(refs, ref)
	}
	sort.Strings(refs)

	return refs
}

func loggerMemberReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`(?:Logger|Level)\.[A-Z][A-Za-z0-9_]*`)
	seen := map[string]bool{}
	for _, part := range codeParts {
		for _, ref := range re.FindAllString(part, -1) {
			seen[ref] = true
		}
	}

	refs := make([]string, 0, len(seen))
	for ref := range seen {
		refs = append(refs, ref)
	}
	sort.Strings(refs)

	return refs
}

func codePartsContainReference(content, ref string) bool {
	for _, part := range markdownCodeParts(content) {
		if !strings.Contains(part, ref) {
			continue
		}
		for _, suffix := range []string{"(", "[", " ", "\n", "\r", "\t"} {
			if strings.Contains(part, ref+suffix) {
				return true
			}
		}
		if strings.TrimSpace(part) == ref {
			return true
		}
	}

	return false
}

func markdownCodeParts(content string) []string {
	var parts []string
	lines := strings.Split(content, "\n")
	inFence := false
	var fence strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inFence {
				parts = append(parts, fence.String())
				fence.Reset()
				inFence = false
			} else {
				inFence = true
			}

			continue
		}

		if inFence {
			fence.WriteString(line)
			fence.WriteByte('\n')

			continue
		}

		parts = append(parts, inlineCodeParts(line)...)
	}

	return parts
}

func inlineCodeParts(line string) []string {
	var parts []string
	for {
		start := strings.IndexByte(line, '`')
		if start < 0 {
			return parts
		}
		rest := line[start+1:]
		end := strings.IndexByte(rest, '`')
		if end < 0 {
			return parts
		}
		parts = append(parts, rest[:end])
		line = rest[end+1:]
	}
}

func missingTerms(doc corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !strings.Contains(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term: %s", doc.label, term))
		}
	}

	return failures
}

func readCorpus(label, root, relPath string) corpus {
	content, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		panic(err)
	}

	return corpus{label: label, path: relPath, content: string(content)}
}

func readSourceFile(sourceRoot, relPath string) corpus {
	content, err := os.ReadFile(filepath.Join(sourceRoot, relPath))
	if err != nil {
		panic(err)
	}

	return corpus{label: relPath, path: relPath, content: string(content)}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}

func logxContractID() string {
	return logxPackage + "#runtime-contract:logger-interface-and-level-strings"
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}

	return false
}

func sameSet(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	gotCopy := append([]string(nil), got...)
	wantCopy := append([]string(nil), want...)
	sort.Strings(gotCopy)
	sort.Strings(wantCopy)
	for i := range gotCopy {
		if gotCopy[i] != wantCopy[i] {
			return false
		}
	}

	return true
}

func runCommand(dir, name string, args ...string) []string {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("%s failed: %v\n%s", strings.Join(append([]string{name}, args...), " "), err, strings.TrimSpace(output.String()))}
	}

	return nil
}
