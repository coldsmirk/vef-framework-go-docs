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
	csvPackage     = "github.com/coldsmirk/vef-framework-go/csv"
	csvFingerprint = "625d27224a8fbc9542243e3ffabba202710b5feba0b34d2d4e1ca0c43630f978"
	csvTopLevel    = 18
	csvFields      = 0
	csvMethods     = 0
	csvEntries     = 18

	englishCsvPath   = "docs/data-tools/tabular.md"
	chineseCsvPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-tools/tabular.md"
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

	englishCsv := readCorpus("English CSV docs", docsRoot, englishCsvPath)
	chineseCsv := readCorpus("Chinese CSV docs", docsRoot, chineseCsvPath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadCsvAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadCsvManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadCsvContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveCsvEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
	}
	for _, doc := range []corpus{englishCsv, chineseCsv} {
		failures = append(failures, verifyDocumentedSurface(doc, entries)...)
		failures = append(failures, verifyNoPhantomCsvRefs(doc, entries)...)
		failures = append(failures, verifyNoCsvOwnedSharedErrors(doc)...)
		failures = append(failures, verifyNoMapOptsOnExporter(doc)...)
		failures = append(failures, missingTerms(doc, csvDocTerms())...)
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

	fmt.Println("csv contracts verified")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != csvPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != csvTopLevel || entry.Fields != csvFields ||
		entry.Methods != csvMethods || entry.Fingerprint != csvFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			csvTopLevel, csvFields, csvMethods, csvFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != csvPackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != csvTopLevel ||
		review.ReviewedSurface.Fields != csvFields ||
		review.ReviewedSurface.Methods != csvMethods ||
		review.ReviewedSurface.EntryCount != csvEntries ||
		review.ReviewedSurface.Fingerprint != csvFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, csvContractID()) {
		failures = append(failures, "contract package review missing csv contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != csvEntries {
		failures = append(failures, fmt.Sprintf("csv audit entry count mismatch: got %d want %d", len(entries), csvEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != csvPackage {
			failures = append(failures, "non-csv audit entry passed into csv verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate csv audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "csv audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != csvTopLevel || counts["field"] != csvFields || counts["method"] != csvMethods {
		failures = append(failures, fmt.Sprintf("csv audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
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
	expected := []string{englishCsvPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest csv coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review csv coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry csv coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, csvPackage)
	if section == "" {
		return []string{index.label + " missing csv package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s csv index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyDocumentedSurface(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		ref := "csv." + entry.Symbol
		if !hasCodeReference(doc.content, ref) {
			failures = append(failures, fmt.Sprintf("%s missing audited csv entry `%s`", doc.label, ref))
		}
	}

	return failures
}

func verifyNoPhantomCsvRefs(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		valid["csv."+entry.Symbol] = true
	}

	var failures []string
	refs := csvReferencesFromMarkdownCode(doc.content)
	for _, ref := range refs {
		if valid[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown csv public API: %s", doc.label, ref))
	}

	return failures
}

func verifyNoCsvOwnedSharedErrors(doc corpus) []string {
	re := regexp.MustCompile(`csv\.Err[A-Za-z0-9_]*`)
	if refs := re.FindAllString(doc.content, -1); len(refs) > 0 {
		sort.Strings(refs)
		return []string{fmt.Sprintf("%s uses csv-owned shared error names: %v", doc.label, refs)}
	}

	return nil
}

func verifyNoMapOptsOnExporter(doc corpus) []string {
	for _, part := range markdownCodeParts(doc.content) {
		for _, line := range strings.Split(part, "\n") {
			if strings.Contains(line, "csv.NewMapExporter") && strings.Contains(line, "mapOpts") {
				return []string{doc.label + " documents mapOpts on csv.NewMapExporter"}
			}
		}
	}

	return nil
}

func hasCodeReference(content, ref string) bool {
	return strings.Contains(content, "`"+ref+"`") ||
		strings.Contains(content, "`"+ref+"(") ||
		strings.Contains(content, "`"+ref+"()") ||
		strings.Contains(content, "`"+ref+"[") ||
		codePartsContainReference(content, ref)
}

func csvDocTerms() []string {
	return []string{
		"csv.NewMapImporter(specs, mapOpts, opts...)",
		"csv.NewMapExporter(specs, opts...)",
		"csv.NewImporter",
		"csv.NewExporter",
		"csv.NewImporterFor",
		"csv.NewExporterFor",
		"csv.NewTypedImporterFor",
		"csv.NewTypedExporterFor",
		"csv.ExportOption",
		"csv.ImportOption",
		"csv.WithImportDelimiter",
		"csv.WithExportDelimiter",
		"csv.WithSkipRows",
		"csv.WithComment",
		"csv.WithCRLF",
		"csv.WithoutHeader",
		"csv.WithoutTrimSpace",
		"csv.WithoutWriteHeader",
		"tabular.Importer",
		"tabular.Exporter",
		"[]tabular.ColumnSpec",
		"[]tabular.MapOption",
		"tabular.NewMapAdapterFromSpecs",
		"ReadAll",
		"FieldsPerRecord = -1",
		"TrimLeadingSpace",
		"tabular.ErrDataMustBeSlice",
		"tabular.ErrNoDataRowsFound",
		"tabular.ErrDuplicateHeaderName",
		"tabular.ErrUnsetField",
		"flush CSV writer",
	}
}

func verifyContractLedger(review contractPackageReview, contract contractEntry, sourceRoot string) []string {
	var failures []string
	expectedCoverage := []string{englishCsvPath}
	if contract.ID != csvContractID() {
		failures = append(failures, fmt.Sprintf("csv contract id mismatch: got %q", contract.ID))
	}
	if contract.Package != csvPackage {
		failures = append(failures, fmt.Sprintf("csv contract package mismatch: got %q", contract.Package))
	}
	if contract.Kind != "runtime-contract" {
		failures = append(failures, fmt.Sprintf("csv contract kind mismatch: got %q", contract.Kind))
	}
	if contract.Disposition != "documented:semantic-contract" {
		failures = append(failures, fmt.Sprintf("csv contract disposition mismatch: got %q", contract.Disposition))
	}
	if !sameSet(contract.Coverage, expectedCoverage) {
		failures = append(failures, fmt.Sprintf("csv contract coverage mismatch: got %v want %v", contract.Coverage, expectedCoverage))
	}
	for _, term := range []string{
		"NewMapImporter",
		"NewMapExporter",
		"[]tabular.MapOption",
		"WithImportDelimiter",
		"WithSkipRows",
		"WithoutHeader",
		"WithoutTrimSpace",
		"WithoutWriteHeader",
		"WithCRLF",
		"ReadAll",
		"FieldsPerRecord",
		"WithComment",
		"flush CSV writer",
	} {
		if !contains(contract.Terms, term) {
			failures = append(failures, "csv contract terms missing "+term)
		}
	}

	allEvidence := append([]string{}, review.SourceEvidence...)
	allEvidence = append(allEvidence, contract.SourceEvidence...)
	allEvidence = append(allEvidence, contract.TestEvidence...)
	for _, item := range allEvidence {
		path, lineText, ok := strings.Cut(item, ":")
		if !ok || lineText == "" {
			failures = append(failures, "csv contract evidence missing line number: "+item)
			continue
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, "csv contract evidence missing file: "+item)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	files := map[string]string{
		"csv.go":           readSourceFile(sourceRoot, "csv/csv.go").content,
		"options.go":       readSourceFile(sourceRoot, "csv/options.go").content,
		"importer.go":      readSourceFile(sourceRoot, "csv/importer.go").content,
		"exporter.go":      readSourceFile(sourceRoot, "csv/exporter.go").content,
		"csv_test.go":      readSourceFile(sourceRoot, "csv/csv_test.go").content,
		"options_test.go":  readSourceFile(sourceRoot, "csv/options_test.go").content,
		"importer_test.go": readSourceFile(sourceRoot, "csv/importer_test.go").content,
		"exporter_test.go": readSourceFile(sourceRoot, "csv/exporter_test.go").content,
	}

	checks := []struct {
		file string
		term string
	}{
		{"csv.go", `func NewMapImporter(`},
		{"csv.go", `specs []tabular.ColumnSpec, mapOpts []tabular.MapOption, opts ...ImportOption`},
		{"csv.go", `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)`},
		{"csv.go", `func NewMapExporter(`},
		{"csv.go", `specs []tabular.ColumnSpec, opts ...ExportOption`},
		{"csv.go", `tabular.NewMapAdapterFromSpecs(specs)`},
		{"csv.go", `tabular.NewTypedImporter[T](NewImporterFor[T](opts...))`},
		{"csv.go", `tabular.NewTypedExporter[T](NewExporterFor[T](opts...))`},
		{"options.go", `type ImportOption func(*importConfig)`},
		{"options.go", `type ExportOption func(*exportConfig)`},
		{"options.go", `o.delimiter = delimiter`},
		{"options.go", `o.hasHeader = false`},
		{"options.go", `o.skipRows = max(rows, 0)`},
		{"options.go", `o.trimSpace = false`},
		{"options.go", `o.comment = comment`},
		{"options.go", `o.writeHeader = false`},
		{"options.go", `o.useCRLF = true`},
		{"importer.go", `delimiter: ','`},
		{"importer.go", `hasHeader: true`},
		{"importer.go", `trimSpace: true`},
		{"importer.go", `comment:   0`},
		{"importer.go", `csvReader.Comma = i.options.delimiter`},
		{"importer.go", `csvReader.Comment = i.options.comment`},
		{"importer.go", `csvReader.FieldsPerRecord = -1`},
		{"importer.go", `csvReader.ReadAll()`},
		{"importer.go", `tabular.ImportRows(rows, i.adapter, i.parsers, tabular.ImportRowsOptions{`},
		{"importer.go", `SkipRows:  i.options.skipRows`},
		{"importer.go", `HasHeader: i.options.hasHeader`},
		{"importer.go", `TrimSpace: i.options.trimSpace`},
		{"exporter.go", `delimiter:   ','`},
		{"exporter.go", `writeHeader: true`},
		{"exporter.go", `useCRLF:     false`},
		{"exporter.go", `csvWriter.Comma = e.options.delimiter`},
		{"exporter.go", `csvWriter.UseCRLF = e.options.useCRLF`},
		{"exporter.go", `csvWriter.Flush()`},
		{"exporter.go", `fmt.Errorf("flush CSV writer: %w", err)`},
		{"csv_test.go", `NewMapImporter(`},
		{"csv_test.go", `[]tabular.MapOption{tabular.WithRowValidator`},
		{"options_test.go", `WithSkipRowsClampsNegative`},
		{"options_test.go", `WithoutHeader`},
		{"options_test.go", `WithoutTrimSpace`},
		{"options_test.go", `WithCRLF`},
		{"importer_test.go", `CustomDelimiter`},
		{"importer_test.go", `WithoutHeader`},
		{"importer_test.go", `WithSkipRows`},
		{"importer_test.go", `EmptyRowsSkipped`},
		{"importer_test.go", `IgnoresUnknownAndMissingColumns`},
		{"exporter_test.go", `WithoutWriteHeader`},
		{"exporter_test.go", `NullPointerValuesEmitEmptyCells`},
	}

	var failures []string
	for _, check := range checks {
		if !strings.Contains(files[check.file], check.term) {
			failures = append(failures, fmt.Sprintf("%s missing source contract term: %s", check.file, check.term))
		}
	}

	if strings.Contains(files["csv.go"], `NewMapExporter(
	specs []tabular.ColumnSpec, mapOpts []tabular.MapOption`) {
		failures = append(failures, "csv source unexpectedly accepts mapOpts on NewMapExporter")
	}

	return failures
}

func runSourceTests(sourceRoot string) []string {
	return runCommand(sourceRoot, "go", "test", "./csv")
}

func loadCsvAuditEntries(path string) []auditEntry {
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
		if entry.Package == csvPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})
	if len(entries) == 0 {
		panic("API audit ledger missing csv entries")
	}

	return entries
}

func loadCsvManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == csvPackage {
			return entry
		}
	}

	panic("API audit manifest missing csv package")
}

func loadCsvContract(path string) (contractPackageReview, contractEntry) {
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
		if item.Package == csvPackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing csv package review")
	}

	var contract contractEntry
	contractFound := false
	for _, item := range ledger.Entries {
		if item.ID == csvContractID() {
			contract = item
			contractFound = true

			break
		}
	}
	if !contractFound {
		panic("API contract ledger missing csv contract entry")
	}

	return review, contract
}

func loadLiveCsvEntry(sourceRoot, docsRoot string) manifestEntry {
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
		if entry.Package == csvPackage {
			return entry
		}
	}

	panic("live API inventory missing csv package")
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

func csvReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`csv\.[A-Z][A-Za-z0-9_]*`)
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

func csvContractID() string {
	return csvPackage + "#runtime-contract:csv-options-and-memory"
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
