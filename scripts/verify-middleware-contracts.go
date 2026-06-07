package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	middlewarePackage     = "github.com/coldsmirk/vef-framework-go/middleware"
	middlewareFingerprint = "2b84c28d70bf6ca996dc263173670f3e39ddf36493968f0d68c951af3d6bdb8e"
	middlewareTopLevel    = 1
	middlewareFields      = 3
	middlewareMethods     = 0
	middlewareEntries     = 4

	englishSPAPath   = "docs/advanced/spa-integration.md"
	chineseSPAPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/advanced/spa-integration.md"
	englishIndexPath = "docs/reference/public-api-index.md"
	chineseIndexPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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

	englishSPA := readCorpus("English SPA docs", docsRoot, englishSPAPath)
	chineseSPA := readCorpus("Chinese SPA docs", docsRoot, chineseSPAPath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadMiddlewareAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadMiddlewareManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadMiddlewareContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveMiddlewareEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
	}
	for _, doc := range []corpus{englishSPA, chineseSPA} {
		failures = append(failures, verifyDocumentedSurface(doc, entries)...)
		failures = append(failures, verifySPATerms(doc)...)
		failures = append(failures, verifyNoStaleExcludePathsLanguage(doc)...)
	}

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

	fmt.Println("middleware contracts verified")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != middlewarePackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != middlewareTopLevel || entry.Fields != middlewareFields ||
		entry.Methods != middlewareMethods || entry.Fingerprint != middlewareFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			middlewareTopLevel, middlewareFields, middlewareMethods, middlewareFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != middlewarePackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != middlewareTopLevel ||
		review.ReviewedSurface.Fields != middlewareFields ||
		review.ReviewedSurface.Methods != middlewareMethods ||
		review.ReviewedSurface.EntryCount != middlewareEntries ||
		review.ReviewedSurface.Fingerprint != middlewareFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, middlewareContractID()) {
		failures = append(failures, "contract package review missing middleware contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != middlewareEntries {
		failures = append(failures, fmt.Sprintf("middleware audit entry count mismatch: got %d want %d", len(entries), middlewareEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != middlewarePackage {
			failures = append(failures, "non-middleware audit entry passed into middleware verifier: "+entry.ID)
		}
		if strings.Contains(entry.Package, "/internal/") {
			failures = append(failures, "internal package included in middleware verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate middleware audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "middleware audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != middlewareTopLevel || counts["field"] != middlewareFields || counts["method"] != middlewareMethods {
		failures = append(failures, fmt.Sprintf("middleware audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	failures = append(failures, verifyExpectedEntries(entries)...)

	return failures
}

func verifyExpectedEntries(entries []auditEntry) []string {
	expected := map[string]string{
		"SPAConfig":              "SPAConfig : github.com/coldsmirk/vef-framework-go/middleware.SPAConfig",
		"SPAConfig.Path":         `Path : string [field_order=1 tag=""]`,
		"SPAConfig.Fs":           `Fs : io/fs.FS [field_order=2 tag=""]`,
		"SPAConfig.ExcludePaths": `ExcludePaths : []string [field_order=3 tag=""]`,
	}

	seen := map[string]string{}
	for _, entry := range entries {
		seen[entry.Symbol] = entry.Signature
	}

	var failures []string
	for symbol, signature := range expected {
		got, ok := seen[symbol]
		if !ok {
			failures = append(failures, "middleware audit entries missing "+symbol)
			continue
		}
		if got != signature {
			failures = append(failures, fmt.Sprintf("middleware audit signature mismatch for %s: got %q want %q", symbol, got, signature))
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
	expected := []string{englishSPAPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest middleware coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review middleware coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry middleware coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, middlewarePackage)
	if section == "" {
		return []string{index.label + " missing middleware package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s middleware index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}
	if strings.Contains(section, "/internal/") {
		failures = append(failures, index.label+" middleware index includes internal package path")
	}

	return failures
}

func verifyDocumentedSurface(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		ref := entry.Symbol
		if entry.Kind == "top" {
			ref = "middleware." + entry.Symbol
		}
		if !hasCodeReference(doc.content, ref) {
			failures = append(failures, fmt.Sprintf("%s missing audited middleware entry `%s`", doc.label, ref))
		}
	}

	return failures
}

func verifySPATerms(doc corpus) []string {
	var failures []string
	failures = append(failures, missingTerms(doc, []string{
		middlewareFingerprint,
		"`middleware.SPAConfig`",
		"`SPAConfig.Path`",
		"`Path string`",
		"`SPAConfig.Fs`",
		"`Fs fs.FS`",
		"`SPAConfig.ExcludePaths`",
		"`ExcludePaths []string`",
		"`vef.ProvideSPAConfig(...)`",
		"`vef.SupplySPAConfigs(...)`",
		"`vef:spa`",
		"`index.html`",
		"`/static/*`",
		"`/api`",
		"`/ws`",
		"`/apidocs`",
	})...)

	alternatives := map[string][]string{
		"exported top-level count": {"1 exported top-level", "1 个 exported top-level"},
		"exported fields count":    {"3 exported field", "3 个 exported fields"},
		"exported methods count":   {"0 exported method", "0 个 exported methods"},
	}
	for label, terms := range alternatives {
		if !containsAnyTerm(doc.content, terms) {
			failures = append(failures, fmt.Sprintf("%s missing %s term: one of %v", doc.label, label, terms))
		}
	}

	return failures
}

func verifyNoStaleExcludePathsLanguage(doc corpus) []string {
	staleTerms := []string{
		"does not actively use",
		"not actively use",
		"not currently enforced",
		"already-enforced",
		"Current limitation",
		"当前限制",
		"并没有真正使用",
		"没有真正使用",
		"未生效",
	}

	var failures []string
	for _, term := range staleTerms {
		if strings.Contains(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s contains stale ExcludePaths limitation wording: %s", doc.label, term))
		}
	}

	return failures
}

func verifyContractLedger(review contractPackageReview, contract contractEntry, sourceRoot string) []string {
	var failures []string
	expectedCoverage := []string{englishSPAPath}
	if contract.ID != middlewareContractID() {
		failures = append(failures, fmt.Sprintf("middleware contract id mismatch: got %q", contract.ID))
	}
	if contract.Package != middlewarePackage {
		failures = append(failures, fmt.Sprintf("middleware contract package mismatch: got %q", contract.Package))
	}
	if contract.Kind != "runtime-contract" {
		failures = append(failures, fmt.Sprintf("middleware contract kind mismatch: got %q", contract.Kind))
	}
	if contract.Disposition != "documented:semantic-contract" {
		failures = append(failures, fmt.Sprintf("middleware contract disposition mismatch: got %q", contract.Disposition))
	}
	if !sameSet(contract.Coverage, expectedCoverage) {
		failures = append(failures, fmt.Sprintf("middleware contract coverage mismatch: got %v want %v", contract.Coverage, expectedCoverage))
	}
	for _, term := range []string{
		"SPAConfig",
		"SPAConfig.Path",
		"Path string",
		"SPAConfig.Fs",
		"Fs fs.FS",
		"SPAConfig.ExcludePaths",
		"ExcludePaths []string",
		"vef:spa",
		"ProvideSPAConfig",
		"SupplySPAConfigs",
		"index.html",
		"static",
		"fallback",
		"path-segment boundaries",
		"empty exclusion prefixes are ignored",
		"trailing slash",
	} {
		if !contains(contract.Terms, term) {
			failures = append(failures, "middleware contract terms missing "+term)
		}
	}

	allEvidence := append([]string{}, review.SourceEvidence...)
	allEvidence = append(allEvidence, contract.SourceEvidence...)
	allEvidence = append(allEvidence, contract.TestEvidence...)
	for _, item := range allEvidence {
		path, lineText, ok := strings.Cut(item, ":")
		if !ok || lineText == "" {
			failures = append(failures, "middleware contract evidence missing line number: "+item)
			continue
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, "middleware contract evidence missing file: "+item)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	files := map[string]string{
		"middleware/spa.go":               readSourceFile(sourceRoot, "middleware/spa.go").content,
		"internal/middleware/spa.go":      readSourceFile(sourceRoot, "internal/middleware/spa.go").content,
		"internal/middleware/spa_test.go": readSourceFile(sourceRoot, "internal/middleware/spa_test.go").content,
		"di.go":                           readSourceFile(sourceRoot, "di.go").content,
	}

	checks := []struct {
		file string
		term string
	}{
		{"middleware/spa.go", "type SPAConfig struct {"},
		{"middleware/spa.go", "Path string"},
		{"middleware/spa.go", "Fs fs.FS"},
		{"middleware/spa.go", "ExcludePaths []string"},
		{"internal/middleware/spa.go", "func spaEntryPath(config *middleware.SPAConfig) string"},
		{"internal/middleware/spa.go", `if config.Path == ""`},
		{"internal/middleware/spa.go", `return "/"`},
		{"internal/middleware/spa.go", `path.Join(spaEntryPath(config), "static") + "/"`},
		{"internal/middleware/spa.go", "func (s *spaMiddleware) Apply(router fiber.Router)"},
		{"internal/middleware/spa.go", "hasAnyPrefix(reqPath, config.ExcludePaths)"},
		{"internal/middleware/spa.go", "ctx.Path(entry)"},
		{"internal/middleware/spa.go", "return ctx.RestartRouting()"},
		{"internal/middleware/spa.go", `group.Get("/", static.New("index.html"`},
		{"internal/middleware/spa.go", "FS:            config.Fs"},
		{"internal/middleware/spa.go", `group.Get("/static/*", static.New("", static.Config{`},
		{"internal/middleware/spa.go", "func hasAnyPrefix(reqPath string, prefixes []string) bool"},
		{"internal/middleware/spa.go", `if prefix == ""`},
		{"internal/middleware/spa.go", `strings.TrimSuffix(prefix, "/")`},
		{"internal/middleware/spa.go", `reqPath == p || strings.HasPrefix(reqPath, p+"/")`},
		{"internal/middleware/spa_test.go", "TestSpaEntryPath"},
		{"internal/middleware/spa_test.go", `path: "", want: "/"`},
		{"internal/middleware/spa_test.go", "TestSpaStaticPrefix"},
		{"internal/middleware/spa_test.go", "TestHasAnyPrefix"},
		{"internal/middleware/spa_test.go", "SharedStringPrefixNotMatched"},
		{"internal/middleware/spa_test.go", "TrailingSlashPrefixNormalized"},
		{"internal/middleware/spa_test.go", "EmptyPrefixIgnored"},
		{"internal/middleware/spa_test.go", "TestSPAMiddlewareApply"},
		{"internal/middleware/spa_test.go", "ExcludedAPIPathIsNotRewritten"},
		{"internal/middleware/spa_test.go", `ExcludePaths: []string{"/api"}`},
		{"internal/middleware/spa_test.go", "ExistingAPIRouteStillReachable"},
		{"di.go", "func ProvideSPAConfig(constructor any, paramTags ...string) fx.Option"},
		{"di.go", "func SupplySPAConfigs(config *middleware.SPAConfig, configs ...*middleware.SPAConfig) fx.Option"},
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
	return runCommand(sourceRoot, "go", "test", "./middleware", "./internal/middleware")
}

func loadMiddlewareAuditEntries(path string) []auditEntry {
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
		if entry.Package == middlewarePackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})
	if len(entries) == 0 {
		panic("API audit ledger missing middleware entries")
	}

	return entries
}

func loadMiddlewareManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == middlewarePackage {
			return entry
		}
	}

	panic("API audit manifest missing middleware package")
}

func loadMiddlewareContract(path string) (contractPackageReview, contractEntry) {
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
		if item.Package == middlewarePackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing middleware package review")
	}

	var contract contractEntry
	contractFound := false
	for _, item := range ledger.Entries {
		if item.ID == middlewareContractID() {
			contract = item
			contractFound = true

			break
		}
	}
	if !contractFound {
		panic("API contract ledger missing middleware contract entry")
	}

	return review, contract
}

func loadLiveMiddlewareEntry(sourceRoot, docsRoot string) manifestEntry {
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
		if entry.Package == middlewarePackage {
			return entry
		}
	}

	panic("live API inventory missing middleware package")
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

func hasCodeReference(content, ref string) bool {
	return strings.Contains(content, "`"+ref+"`") ||
		strings.Contains(content, "`"+ref+"(") ||
		strings.Contains(content, "`"+ref+"()") ||
		strings.Contains(content, "`"+ref+"[") ||
		codePartsContainReference(content, ref)
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

func containsAnyTerm(content string, terms []string) bool {
	for _, term := range terms {
		if strings.Contains(content, term) {
			return true
		}
	}

	return false
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

func middlewareContractID() string {
	return middlewarePackage + "#runtime-contract:spa-middleware-config"
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
