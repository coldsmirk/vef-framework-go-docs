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
	decimalPackage     = "github.com/coldsmirk/vef-framework-go/decimal"
	decimalFingerprint = "ea79b685aa80a0df3929fb69b8a3e0941805ce057e5c7ebf81a853da467d8401"
	decimalTopLevel    = 22
	decimalFields      = 0
	decimalMethods     = 70
	decimalEntries     = 92

	englishDecimalPath = "docs/utilities/decimal.md"
	chineseDecimalPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/decimal.md"
	englishIndexPath   = "docs/reference/public-api-index.md"
	chineseIndexPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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

	englishDocs := readCorpus("English decimal docs", docsRoot, englishDecimalPath)
	chineseDocs := readCorpus("Chinese decimal docs", docsRoot, chineseDecimalPath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadDecimalAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadDecimalManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadDecimalContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveDecimalEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, verifyDocumentedTopLevel(doc, entries)...)
		failures = append(failures, verifyDocumentedMethods(doc, entries)...)
		failures = append(failures, verifyNoPhantomDecimalRefs(doc, entries)...)
		failures = append(failures, missingTerms(doc, decimalDocTerms(doc.label))...)
	}

	failures = append(failures, verifyContractLedger(review, contract, sourceRoot)...)
	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, verifyUpstreamContracts(sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("decimal contracts verified")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != decimalPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != decimalTopLevel || entry.Fields != decimalFields ||
		entry.Methods != decimalMethods || entry.Fingerprint != decimalFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			decimalTopLevel, decimalFields, decimalMethods, decimalFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != decimalPackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != decimalTopLevel ||
		review.ReviewedSurface.Fields != decimalFields ||
		review.ReviewedSurface.Methods != decimalMethods ||
		review.ReviewedSurface.EntryCount != decimalEntries ||
		review.ReviewedSurface.Fingerprint != decimalFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, decimalContractID()) {
		failures = append(failures, "contract package review missing decimal contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != decimalEntries {
		failures = append(failures, fmt.Sprintf("decimal audit entry count mismatch: got %d want %d", len(entries), decimalEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != decimalPackage {
			failures = append(failures, "non-decimal audit entry passed into decimal verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate decimal audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "decimal audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != decimalTopLevel || counts["field"] != decimalFields || counts["method"] != decimalMethods {
		failures = append(failures, fmt.Sprintf("decimal audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
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
	expected := []string{englishDecimalPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest decimal coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review decimal coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry decimal coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, decimalPackage)
	if section == "" {
		return []string{index.label + " missing decimal package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s decimal index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyDocumentedTopLevel(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		if entry.Kind != "top" {
			continue
		}
		ref := "decimal." + entry.Symbol
		if !hasCodeReference(doc.content, ref) {
			failures = append(failures, fmt.Sprintf("%s missing audited decimal top-level entry `%s`", doc.label, ref))
		}
	}

	return failures
}

func verifyDocumentedMethods(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		if entry.Kind != "method" {
			continue
		}
		if !hasCodeReference(doc.content, entry.Symbol) {
			failures = append(failures, fmt.Sprintf("%s missing audited decimal method `%s`", doc.label, entry.Symbol))
		}
	}

	return failures
}

func verifyNoPhantomDecimalRefs(doc corpus, entries []auditEntry) []string {
	validTop := map[string]bool{}
	validMethods := map[string]bool{}
	for _, entry := range entries {
		switch entry.Kind {
		case "top":
			validTop["decimal."+entry.Symbol] = true
		case "method":
			validMethods[entry.Symbol] = true
		}
	}

	var failures []string
	for _, ref := range decimalReferencesFromMarkdownCode(doc.content) {
		if validTop[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown decimal public API: %s", doc.label, ref))
	}
	for _, ref := range decimalMethodReferencesFromMarkdownCode(doc.content) {
		if validMethods[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown Decimal public method: %s", doc.label, ref))
	}

	return failures
}

func hasCodeReference(content, ref string) bool {
	return strings.Contains(content, "`"+ref+"`") ||
		strings.Contains(content, "`"+ref+"(") ||
		strings.Contains(content, "`"+ref+"()")
}

func decimalDocTerms(label string) []string {
	common := []string{
		decimalFingerprint,
		"shopspring/decimal v1.4.0",
		"type alias",
		"DivisionPrecision",
		"MarshalJSONWithoutQuotes",
		"PowPrecisionNegativeExponent",
		"ExpMaxIterations",
		"*regexp.Regexp",
		"regexp.MustCompile(\"[$,]\")",
		"nil `*Decimal`",
		"decimal: unsupported type",
		"division by zero",
		"quoted JSON string",
		"Decimal.Float64",
		"Decimal.InexactFloat64",
	}
	if strings.HasPrefix(label, "Chinese") {
		return append(common,
			"22 个 top-level symbols",
			"0 个 exported fields",
			"70 个 exported methods",
			"没有\nexported fields",
			"NaN 或 infinity",
		)
	}

	return append(common,
		"22 top-level symbols",
		"0 exported fields",
		"70 exported methods",
		"no exported fields",
		"NaN or infinity",
	)
}

func verifyContractLedger(review contractPackageReview, contract contractEntry, sourceRoot string) []string {
	var failures []string
	if contract.ID != decimalContractID() {
		failures = append(failures, fmt.Sprintf("decimal contract entry id mismatch: got %q", contract.ID))
	}
	if contract.Package != decimalPackage {
		failures = append(failures, fmt.Sprintf("decimal contract entry package mismatch: got %q", contract.Package))
	}
	if contract.Kind != "runtime-contract" {
		failures = append(failures, fmt.Sprintf("decimal contract entry kind mismatch: got %q", contract.Kind))
	}
	if contract.Contract != "NewFromAny supported input families and unsupported-type behavior" {
		failures = append(failures, fmt.Sprintf("decimal contract label mismatch: got %q", contract.Contract))
	}

	expectedEvidence := decimalSourceEvidence(sourceRoot)
	if !sameSet(review.SourceEvidence, expectedEvidence) {
		failures = append(failures, fmt.Sprintf("decimal package review source evidence mismatch: got %v want %v", review.SourceEvidence, expectedEvidence))
	}
	if !sameSet(contract.SourceEvidence, expectedEvidence) {
		failures = append(failures, fmt.Sprintf("decimal contract source evidence mismatch: got %v want %v", contract.SourceEvidence, expectedEvidence))
	}
	expectedTestEvidence := []string{"decimal/convert_test.go:16"}
	if !sameSet(contract.TestEvidence, expectedTestEvidence) {
		failures = append(failures, fmt.Sprintf("decimal contract test evidence mismatch: got %v want %v", contract.TestEvidence, expectedTestEvidence))
	}
	for _, term := range []string{
		"shopspring/decimal v1.4.0",
		"type alias",
		"top-level symbols",
		"exported methods",
		"decimal.NewFromFormattedString",
		"*regexp.Regexp",
		"decimal.NewFromAny",
		"fmt.Stringer",
		"nil `*Decimal`",
		"decimal: unsupported type",
		"decimal.MustFromAny",
		"DivisionPrecision",
		"MarshalJSONWithoutQuotes",
		"division by zero",
		"quoted JSON string",
		"Decimal.Float64",
		"Decimal.InexactFloat64",
	} {
		if !contains(contract.Terms, term) {
			failures = append(failures, "decimal contract missing term "+term)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	checks := []struct {
		path  string
		terms []string
	}{
		{
			path: "go.mod",
			terms: []string{
				"github.com/shopspring/decimal v1.4.0",
			},
		},
		{
			path: "decimal/decimal.go",
			terms: []string{
				"type Decimal = decimal.Decimal",
				"Zero = decimal.Zero",
				"One  = decimal.NewFromInt(1)",
				"New                      = decimal.New",
				"NewFromFloat             = decimal.NewFromFloat",
				"NewFromFloat32           = decimal.NewFromFloat32",
				"NewFromFloatWithExponent = decimal.NewFromFloatWithExponent",
				"NewFromInt               = decimal.NewFromInt",
				"NewFromInt32             = decimal.NewFromInt32",
				"NewFromUint64            = decimal.NewFromUint64",
				"NewFromBigInt            = decimal.NewFromBigInt",
				"NewFromBigRat            = decimal.NewFromBigRat",
				"NewFromString            = decimal.NewFromString",
				"NewFromFormattedString   = decimal.NewFromFormattedString",
				"RequireFromString        = decimal.RequireFromString",
				"Max         = decimal.Max",
				"Min         = decimal.Min",
				"Sum         = decimal.Sum",
				"Avg         = decimal.Avg",
				"RescalePair = decimal.RescalePair",
			},
		},
		{
			path: "decimal/convert.go",
			terms: []string{
				"func NewFromAny(v any) (Decimal, error)",
				"case Decimal:",
				"case *Decimal:",
				"return Zero, nil",
				"case int:",
				"case uint64:",
				"case float32:",
				"case float64:",
				"case string:",
				"case []byte:",
				"case bool:",
				"return One, nil",
				"if s, ok := v.(fmt.Stringer); ok",
				"return NewFromString(s.String())",
				"return Zero, fmt.Errorf(\"%w: %T\", errUnsupportedType, v)",
				"func MustFromAny(v any) Decimal",
				"panic(err)",
			},
		},
		{
			path: "decimal/errors.go",
			terms: []string{
				"var errUnsupportedType = errors.New(\"decimal: unsupported type\")",
			},
		},
		{
			path: "decimal/convert_test.go",
			terms: []string{
				"DecimalPointerNil",
				"Uint64ExceedsInt64Max",
				"StringInvalid",
				"ByteSliceInvalid",
				"BoolTrue",
				"BoolFalse",
				"StringerInvalid",
				"UnsupportedStruct",
				"NilInput",
				"PanicsOnUnsupportedType",
			},
		},
	}

	var failures []string
	for _, check := range checks {
		source := readSourceFile(sourceRoot, check.path)
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	return failures
}

func verifyUpstreamContracts(sourceRoot string) []string {
	modPath, err := runCommandOutput(sourceRoot, "go", "env", "GOMODCACHE")
	if err != nil {
		return []string{err.Error()}
	}
	upstreamRoot := filepath.Join(strings.TrimSpace(modPath), "github.com/shopspring/decimal@v1.4.0")
	upstream := readSourceFile(upstreamRoot, "decimal.go")
	return missingTerms(upstream, []string{
		"var DivisionPrecision = 16",
		"var MarshalJSONWithoutQuotes = false",
		"var ExpMaxIterations = 1000",
		"func NewFromFormattedString(value string, replRegexp *regexp.Regexp) (Decimal, error)",
		"func (d Decimal) Div(d2 Decimal) Decimal",
		"return d.DivRound(d2, int32(DivisionPrecision))",
		"panic(\"decimal division by 0\")",
		"func (d Decimal) Float64() (f float64, exact bool)",
		"func (d Decimal) InexactFloat64() float64",
		"func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error",
		"func (d Decimal) MarshalJSON() ([]byte, error)",
		"if MarshalJSONWithoutQuotes",
	})
}

func runSourceTests(sourceRoot string) []string {
	return runCommand(sourceRoot, "go", "test", "./decimal")
}

func loadDecimalAuditEntries(path string) []auditEntry {
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
		if entry.Package == decimalPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind != entries[j].Kind {
			return entries[i].Kind < entries[j].Kind
		}

		return entries[i].Symbol < entries[j].Symbol
	})

	return entries
}

func loadDecimalManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == decimalPackage {
			return entry
		}
	}

	panic("API audit manifest missing decimal package")
}

func loadDecimalContract(path string) (contractPackageReview, contractEntry) {
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
		if item.Package == decimalPackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing decimal package review")
	}

	var contract contractEntry
	contractFound := false
	for _, item := range ledger.Entries {
		if item.ID == decimalContractID() {
			contract = item
			contractFound = true

			break
		}
	}
	if !contractFound {
		panic("API contract ledger missing decimal contract entry")
	}

	return review, contract
}

func loadLiveDecimalEntry(sourceRoot, docsRoot string) manifestEntry {
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
		if entry.Package == decimalPackage {
			return entry
		}
	}

	panic("live API inventory missing decimal package")
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

func decimalReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`decimal\.[A-Z][A-Za-z0-9_]*`)
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

func decimalMethodReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`Decimal\.[A-Z][A-Za-z0-9_]*`)
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

func decimalSourceEvidence(sourceRoot string) []string {
	evidence := []string{
		"decimal/decimal.go:5",
		"decimal/convert.go:8",
		"decimal/errors.go:5",
		"decimal/convert_test.go:16",
		"go.mod:41",
	}
	for _, item := range evidence {
		path, _, _ := strings.Cut(item, ":")
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			panic(err)
		}
	}

	return evidence
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

func readSourceFile(root, relPath string) corpus {
	content, err := os.ReadFile(filepath.Join(root, relPath))
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

func decimalContractID() string {
	return decimalPackage + "#runtime-contract:any-conversion-and-error-behavior"
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

func runCommandOutput(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s failed: %v\n%s", strings.Join(append([]string{name}, args...), " "), err, strings.TrimSpace(output.String()))
	}

	return output.String(), nil
}
