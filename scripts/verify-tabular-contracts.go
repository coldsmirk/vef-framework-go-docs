package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	tabularPackage = "github.com/coldsmirk/vef-framework-go/tabular"

	tabularFingerprint = "4bd760b4d107c89e997d9354f50b7a39c3fe441019e4b0e676bb6aeaf2ed5b94"
	tabularTopLevel    = 68
	tabularFields      = 37
	tabularMethods     = 38
	tabularEntries     = 143

	tabularGroupedEntries              = 75
	tabularGroupedFields               = 37
	tabularGroupedMethods              = 38
	tabularGroupedReceivers            = 20
	tabularGroupedSignatureFingerprint = "ad9d25935821a4995a105325e8ee31250becdab85063bc9b6b8c37c77a09c89a"
	tabularGroupedReceiverFingerprint  = "1d5eac3f748671a7761a9fbf07c39fbf4d86e4aa35b8f79a46e600a1fd6fa67d"

	englishTabularPath = "docs/features/tabular.md"
	chineseTabularPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/tabular.md"
	englishIndexPath   = "docs/reference/public-api-index.md"
	chineseIndexPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
)

type corpus struct {
	label   string
	content string
}

type auditLedger struct {
	Entries []auditEntry `json:"entries"`
}

type auditEntry struct {
	ID          string   `json:"id"`
	Package     string   `json:"package"`
	Kind        string   `json:"kind"`
	Symbol      string   `json:"symbol"`
	Signature   string   `json:"signature"`
	Disposition string   `json:"disposition"`
	Coverage    []string `json:"coverage"`
}

type manifest struct {
	Packages []manifestEntry `json:"packages"`
}

type manifestEntry struct {
	Package     string   `json:"package"`
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
	Disposition    string   `json:"disposition"`
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	TestEvidence   []string `json:"test_evidence"`
	Terms          []string `json:"terms"`
}

type liveInventoryEntry struct {
	Package     string   `json:"package"`
	Coverage    []string `json:"coverage"`
	TopLevel    int      `json:"top_level"`
	Fields      int      `json:"fields"`
	Methods     int      `json:"methods"`
	Fingerprint string   `json:"fingerprint"`
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)
	manifestPath := filepath.Join(docsRoot, "scripts/api-audit-manifest.json")
	auditLedgerPath := filepath.Join(docsRoot, "scripts/api-audit-ledger.json")
	contractLedgerPath := filepath.Join(docsRoot, "scripts/api-contract-ledger.json")

	englishDocs := readCorpus("English tabular docs", filepath.Join(docsRoot, englishTabularPath))
	chineseDocs := readCorpus("Chinese tabular docs", filepath.Join(docsRoot, chineseTabularPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[tabularPackage]
	liveAuditEntries := tabularEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	tabularEntries := tabularEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(tabularEntries)...)
	failures = append(failures, verifyLiveAuditEntries(tabularEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedTabularSurface(tabularEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, tabularEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, tabularEntries)...)
	failures = append(failures, verifyTabularDocs([]corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Tabular contract docs verified: 143 public entries, 75 grouped entries, schema/mapping/tag contracts")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != tabularPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != tabularTopLevel ||
		entry.Fields != tabularFields ||
		entry.Methods != tabularMethods ||
		entry.Fingerprint != tabularFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s tabular surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			tabularTopLevel, tabularFields, tabularMethods, tabularFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != tabularPackage {
			continue
		}

		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, tabularCoverage()) {
			failures = append(failures, fmt.Sprintf("tabular manifest coverage mismatch: got %v want %v", entry.Coverage, tabularCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing tabular package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	expectedContracts := map[string][]string{
		tabularPackage + "#field-contract:dynamic-schema-and-map-validation": {
			"ColumnSpec",
			"Required",
			"Validators",
			"RowValidator",
			"errors.Join",
		},
		tabularPackage + "#runtime-contract:header-mapping-and-row-import": {
			"BuildHeaderMapping",
			"DefaultPositionalMapping",
			"WithoutHeader",
			"ErrDuplicateHeaderName",
			"ImportError",
		},
		tabularPackage + "#string-contract:tabular-tag-grammar": {
			"key=value",
			"dive",
			"ID;order=1",
			"ErrDuplicateColumnKey",
		},
	}

	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != tabularPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "tabular contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != tabularTopLevel ||
			review.ReviewedSurface.Fields != tabularFields ||
			review.ReviewedSurface.Methods != tabularMethods ||
			review.ReviewedSurface.EntryCount != tabularEntries ||
			review.ReviewedSurface.Fingerprint != tabularFingerprint {
			failures = append(failures, fmt.Sprintf(
				"tabular contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, tabularCoverage()) {
			failures = append(failures, fmt.Sprintf("tabular contract review coverage mismatch: got %v want %v", review.Coverage, tabularCoverage()))
		}
		if !sameSet(review.ContractIDs, sortedKeys(expectedContracts)) {
			failures = append(failures, fmt.Sprintf("tabular contract ids mismatch: got %v want %v", review.ContractIDs, sortedKeys(expectedContracts)))
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing tabular package review")
	}

	foundContracts := map[string]bool{}
	for _, entry := range contracts.Entries {
		terms, ok := expectedContracts[entry.ID]
		if !ok {
			continue
		}
		foundContracts[entry.ID] = true
		if entry.Package != tabularPackage {
			failures = append(failures, fmt.Sprintf("tabular contract entry package mismatch for %s: %s", entry.ID, entry.Package))
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "tabular contract entry disposition mismatch for "+entry.ID+": "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, tabularCoverage()) {
			failures = append(failures, fmt.Sprintf("tabular contract coverage mismatch for %s: got %v want %v", entry.ID, entry.Coverage, tabularCoverage()))
		}
		for _, term := range terms {
			if !contains(entry.Terms, term) {
				failures = append(failures, fmt.Sprintf("tabular contract %s missing term %s", entry.ID, term))
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.TestEvidence)...)
	}
	for id := range expectedContracts {
		if !foundContracts[id] {
			failures = append(failures, "contract ledger missing tabular contract entry "+id)
		}
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != tabularEntries {
		failures = append(failures, fmt.Sprintf("tabular audit entry count mismatch: got %d want %d", len(entries), tabularEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != tabularPackage {
			failures = append(failures, "non-tabular audit entry passed into tabular verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate tabular audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "tabular audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, tabularCoverage()) {
			failures = append(failures, fmt.Sprintf("tabular audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, tabularCoverage()))
		}
	}
	if counts["top"] != tabularTopLevel || counts["field"] != tabularFields || counts["method"] != tabularMethods {
		failures = append(failures, fmt.Sprintf("tabular audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != tabularTopLevel ||
		dispositionCounts["grouped:type-member-family"] != tabularGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"tabular audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			tabularTopLevel,
			tabularGroupedEntries,
		))
	}

	return failures
}

func verifyLiveAuditEntries(ledgerEntries, liveEntries []auditEntry) []string {
	ledgerByID := entriesByID(ledgerEntries)
	liveByID := entriesByID(liveEntries)
	var failures []string

	for id, live := range liveByID {
		ledger, ok := ledgerByID[id]
		if !ok {
			failures = append(failures, fmt.Sprintf("tabular missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"tabular live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
				id,
				ledger.Kind,
				ledger.Symbol,
				ledger.Signature,
				live.Kind,
				live.Symbol,
				live.Signature,
			))
		}
	}
	for id, ledger := range ledgerByID {
		if _, ok := liveByID[id]; !ok {
			failures = append(failures, fmt.Sprintf("tabular extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedTabularSurface(entries []auditEntry, docs []corpus) []string {
	var rows []string
	receiverCounts := map[string]int{}
	kindCounts := map[string]int{}
	var failures []string
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Disposition, "grouped:") {
			continue
		}
		receiver, ok := receiverForSymbol(entry.Symbol)
		if !ok {
			failures = append(failures, fmt.Sprintf("tabular grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("tabular grouped type-member surface", rows, tabularGroupedEntries, tabularGroupedSignatureFingerprint)...)
	if kindCounts["field"] != tabularGroupedFields || kindCounts["method"] != tabularGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"tabular grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			tabularGroupedFields,
			tabularGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("tabular grouped receiver/type families", receiverRows, tabularGroupedReceivers, tabularGroupedReceiverFingerprint)...)

	for _, doc := range docs {
		for _, term := range []string{
			"143 public tabular entries",
			"75 grouped tabular field/method entries",
			"20 tabular receiver/type families",
			"37 exported tabular field entries",
			"38 exported tabular method entries",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped tabular audit term "+term)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, tabularPackage)
	if section == "" {
		return []string{index.label + " missing tabular package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s tabular index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyTabularDocs(docs []corpus) []string {
	var failures []string
	for _, doc := range docs {
		for _, term := range []string{
			"`tabular`",
			"`tabular:\"-\"`",
			"`tabular:\"dive\"`",
			"`key=value`",
			"`tabular:\"name=ID;order=1\"`",
			"`NewSchemaFromSpecs`",
			"`ColumnSpec.Required`",
			"`Validators`",
			"`RowValidator`",
			"`errors.Join`",
			"`BuildHeaderMapping`",
			"`DefaultPositionalMapping`",
			"`WithoutHeader()`",
			"`ErrDuplicateHeaderName`",
			"`ParseRow`",
			"`ImportError`",
			"`Column.Default`",
			"`ColumnSpec.Default`",
			"`MappingOptions.TrimSpace`",
			"`ParseRowOptions.TrimSpace`",
			"`TypedExporter.Inner`",
			"`TypedImporter.Inner`",
			"`ImportError.Unwrap`",
			"`ExportError.Unwrap`",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing tabular semantic term "+term)
			}
		}
	}

	return failures
}

func verifySourceTerms(sourceRoot string) []string {
	checks := []struct {
		path  string
		terms []string
	}{
		{
			path: "tabular/constants.go",
			terms: []string{
				"TagTabular = \"tabular\"",
				"AttrDive      = \"dive\"",
				"IgnoreField = \"-\"",
			},
		},
		{
			path: "tabular/spec.go",
			terms: []string{
				"type ColumnSpec struct",
				"Key string",
				"Type reflect.Type",
				"Required bool",
				"Validators []CellValidator",
				"return nil, fmt.Errorf(\"%w: spec #%d\", ErrMissingColumnKey, i)",
				"return nil, fmt.Errorf(\"%w: %s\", ErrMissingColumnType, spec.Key)",
				"return nil, fmt.Errorf(\"%w: %s\", ErrDuplicateColumnKey, spec.Key)",
				"return nil, fmt.Errorf(\"%w: %s\", ErrDuplicateHeaderName, name)",
			},
		},
		{
			path: "tabular/mapping.go",
			terms: []string{
				"type MappingOptions struct",
				"TrimSpace bool",
				"headerName = strings.TrimSpace(headerName)",
				"if headerName == \"\"",
				"return nil, fmt.Errorf(\"%w: %s\", ErrDuplicateHeaderName, headerName)",
				"func DefaultPositionalMapping(schema *Schema) map[int]int",
			},
		},
		{
			path: "tabular/parse_row.go",
			terms: []string{
				"type ParseRowOptions struct",
				"cellValue = strings.TrimSpace(cellValue)",
				"if cellValue == \"\" && column.Default != \"\"",
				"if cellValue == \"\"",
				"Err:    fmt.Errorf(\"parse value: %w\", err)",
				"When the returned slice is non-empty",
			},
		},
		{
			path: "tabular/errors.go",
			terms: []string{
				"ErrDataMustBeSlice = errors.New(\"data must be a slice\")",
				"ErrDuplicateColumnKey = errors.New(\"duplicate column key\")",
				"ErrDuplicateHeaderName = errors.New(\"duplicate header name\")",
				"ErrRequiredMissing = errors.New(\"required value is missing\")",
				"ErrTypedRowMismatch = errors.New(\"importer returned unexpected row type\")",
				"func (e ImportError) Unwrap() error",
				"func (e ExportError) Unwrap() error",
			},
		},
		{
			path: "tabular/map_adapter.go",
			terms: []string{
				"errors.Join(errs...)",
				"WithRowValidator",
				"ErrRequiredMissing",
			},
		},
		{
			path: "tabular/parser.go",
			terms: []string{
				"strx.ParseTag(tag)",
				"TagTabular",
				"IgnoreField",
				"AttrDive",
			},
		},
	}

	var failures []string
	for _, check := range checks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		for _, term := range check.terms {
			if !strings.Contains(source.content, term) {
				failures = append(failures, source.label+" missing source term "+term)
			}
		}
	}

	return failures
}

func runGoTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./tabular", "./csv", "./excel")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./tabular ./csv ./excel failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func tabularEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == tabularPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

	return entries
}

func loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath string) map[string]manifestEntry {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", manifestPath,
		"-ledger", auditLedgerPath,
		"-contract-ledger", contractLedgerPath,
		"-print-current",
	)
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("verify-api-audit -print-current failed: %w\n%s", err, strings.TrimSpace(string(output))))
	}

	payload := "[" + strings.TrimSpace(string(output)) + "]"
	var entries []liveInventoryEntry
	if err := json.Unmarshal([]byte(payload), &entries); err != nil {
		panic(fmt.Errorf("parse live inventory: %w", err))
	}

	result := map[string]manifestEntry{}
	for _, entry := range entries {
		result[entry.Package] = manifestEntry{
			Package:     entry.Package,
			Coverage:    entry.Coverage,
			TopLevel:    entry.TopLevel,
			Fields:      entry.Fields,
			Methods:     entry.Methods,
			Fingerprint: entry.Fingerprint,
		}
	}

	return result
}

func loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath string) auditLedger {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", manifestPath,
		"-print-ledger",
	)
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("verify-api-audit -print-ledger failed: %w\n%s", err, strings.TrimSpace(string(output))))
	}

	var ledger auditLedger
	if err := json.Unmarshal(output, &ledger); err != nil {
		panic(fmt.Errorf("parse live audit ledger: %w", err))
	}

	return ledger
}

func verifySourceEvidence(sourceRoot string, evidence []string) []string {
	var failures []string
	for _, item := range evidence {
		path, lineText, ok := strings.Cut(item, ":")
		if !ok {
			failures = append(failures, "bad source evidence format "+item)
			continue
		}
		if _, err := strconv.Atoi(lineText); err != nil {
			failures = append(failures, "bad source evidence line "+item)
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, "source evidence missing "+item)
		}
	}

	return failures
}

func verifyGroupedFingerprint(label string, rows []string, wantCount int, wantFingerprint string) []string {
	gotFingerprint := fingerprintRows(rows)
	var failures []string
	if len(rows) != wantCount {
		failures = append(failures, fmt.Sprintf("%s count mismatch: got %d want %d", label, len(rows), wantCount))
	}
	if gotFingerprint != wantFingerprint {
		failures = append(failures, fmt.Sprintf("%s fingerprint mismatch: got %s want %s", label, gotFingerprint, wantFingerprint))
	}

	return failures
}

func fingerprintRows(rows []string) string {
	sorted := append([]string(nil), rows...)
	sort.Strings(sorted)

	hash := sha256.New()
	for _, row := range sorted {
		hash.Write([]byte(row))
		hash.Write([]byte("\n"))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func entriesByID(entries []auditEntry) map[string]auditEntry {
	result := map[string]auditEntry{}
	for _, entry := range entries {
		result[entry.ID] = entry
	}

	return result
}

func receiverForSymbol(symbol string) (string, bool) {
	receiver, _, ok := strings.Cut(symbol, ".")
	if !ok || receiver == "" {
		return "", false
	}

	return receiver, true
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

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func loadJSON[T any](path string) T {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var result T
	if err := json.Unmarshal(content, &result); err != nil {
		panic(err)
	}

	return result
}

func sameSet(got, want []string) bool {
	got = sortedUnique(got)
	want = sortedUnique(want)
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}

	return true
}

func sortedKeys[V any](m map[string]V) []string {
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	sort.Strings(result)

	return result
}

func sortedUnique(values []string) []string {
	set := sliceSet(values)
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)

	return result
}

func sliceSet(values []string) map[string]bool {
	result := map[string]bool{}
	for _, value := range values {
		result[value] = true
	}

	return result
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}

func containsNormalized(content, term string) bool {
	return strings.Contains(content, term) ||
		strings.Contains(strings.Join(strings.Fields(content), " "), strings.Join(strings.Fields(term), " "))
}

func tabularCoverage() []string {
	return []string{englishTabularPath}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
