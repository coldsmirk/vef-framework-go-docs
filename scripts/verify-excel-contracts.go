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
	excelPackage     = "github.com/coldsmirk/vef-framework-go/excel"
	excelFingerprint = "a449ebeda509ae9b0a2c7bfa083c70b45bc4635bdb49ff1d674400aced129324"
	excelTopLevel    = 17
	excelFields      = 0
	excelMethods     = 0
	excelEntries     = 17

	englishExcelPath    = "docs/features/excel.md"
	chineseExcelPath    = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/excel.md"
	englishCsvExcelPath = "docs/features/csv-excel.md"
	chineseCsvExcelPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/csv-excel.md"
	englishIndexPath    = "docs/reference/public-api-index.md"
	chineseIndexPath    = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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

	englishExcel := readCorpus("English Excel docs", docsRoot, englishExcelPath)
	chineseExcel := readCorpus("Chinese Excel docs", docsRoot, chineseExcelPath)
	englishCsvExcel := readCorpus("English CSV/Excel docs", docsRoot, englishCsvExcelPath)
	chineseCsvExcel := readCorpus("Chinese CSV/Excel docs", docsRoot, chineseCsvExcelPath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadExcelAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadExcelManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadExcelContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveExcelEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
	}
	for _, doc := range []corpus{englishExcel, chineseExcel, englishCsvExcel, chineseCsvExcel} {
		failures = append(failures, verifyDocumentedSurface(doc, entries)...)
		failures = append(failures, verifyNoPhantomExcelRefs(doc, entries)...)
		failures = append(failures, missingTerms(doc, excelDocTerms(doc.label))...)
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

	fmt.Println("excel contracts verified")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != excelPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != excelTopLevel || entry.Fields != excelFields ||
		entry.Methods != excelMethods || entry.Fingerprint != excelFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			excelTopLevel, excelFields, excelMethods, excelFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != excelPackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != excelTopLevel ||
		review.ReviewedSurface.Fields != excelFields ||
		review.ReviewedSurface.Methods != excelMethods ||
		review.ReviewedSurface.EntryCount != excelEntries ||
		review.ReviewedSurface.Fingerprint != excelFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, excelContractID()) {
		failures = append(failures, "contract package review missing excel contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != excelEntries {
		failures = append(failures, fmt.Sprintf("excel audit entry count mismatch: got %d want %d", len(entries), excelEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != excelPackage {
			failures = append(failures, "non-excel audit entry passed into excel verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate excel audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "excel audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != excelTopLevel || counts["field"] != excelFields || counts["method"] != excelMethods {
		failures = append(failures, fmt.Sprintf("excel audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
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
	expected := []string{englishExcelPath, englishCsvExcelPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest excel coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review excel coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry excel coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, excelPackage)
	if section == "" {
		return []string{index.label + " missing excel package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s excel index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyDocumentedSurface(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		ref := "excel." + entry.Symbol
		if !hasCodeReference(doc.content, ref) {
			failures = append(failures, fmt.Sprintf("%s missing audited excel entry `%s`", doc.label, ref))
		}
	}

	return failures
}

func verifyNoPhantomExcelRefs(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		valid["excel."+entry.Symbol] = true
	}

	var failures []string
	refs := excelReferencesFromMarkdownCode(doc.content)
	for _, ref := range refs {
		if valid[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown excel public API: %s", doc.label, ref))
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

func excelDocTerms(label string) []string {
	common := []string{
		excelFingerprint,
		"tabular.Importer",
		"tabular.Exporter",
		"[]tabular.ColumnSpec",
		"[]tabular.MapOption",
		"tabular.NewMapAdapterFromSpecs",
		"excel.NewMapImporter(specs, mapOpts, opts...)",
		"excel.WithImportSheetName",
		"excel.WithImportSheetIndex",
		"excel.WithSkipRows",
		"excel.WithoutHeader",
		"excel.WithoutTrimSpace",
		"excel.ErrSheetIndexOutOfRange",
		"tabular.ErrNoDataRowsFound",
		"tabular.ErrDuplicateHeaderName",
		"native typed cell",
		"timex.Time",
		"FormatterFn",
	}

	if strings.Contains(label, "CSV/Excel") {
		return []string{
			"excel.NewMapImporter(specs, mapOpts, opts...)",
			"excel.ExportOption",
			"excel.ImportOption",
			"excel.WithImportSheetName",
			"excel.WithImportSheetIndex",
			"excel.WithSkipRows",
			"excel.WithoutHeader",
			"excel.WithoutTrimSpace",
			"excel.ErrSheetIndexOutOfRange",
			"native typed cell",
			"timex.Time",
			"FormatterFn",
		}
	}

	return common
}

func verifyContractLedger(review contractPackageReview, contract contractEntry, sourceRoot string) []string {
	var failures []string
	expectedCoverage := []string{englishExcelPath, englishCsvExcelPath}
	if contract.ID != excelContractID() {
		failures = append(failures, fmt.Sprintf("excel contract id mismatch: got %q", contract.ID))
	}
	if contract.Package != excelPackage {
		failures = append(failures, fmt.Sprintf("excel contract package mismatch: got %q", contract.Package))
	}
	if contract.Kind != "runtime-contract" {
		failures = append(failures, fmt.Sprintf("excel contract kind mismatch: got %q", contract.Kind))
	}
	if contract.Disposition != "documented:semantic-contract" {
		failures = append(failures, fmt.Sprintf("excel contract disposition mismatch: got %q", contract.Disposition))
	}
	if !sameSet(contract.Coverage, expectedCoverage) {
		failures = append(failures, fmt.Sprintf("excel contract coverage mismatch: got %v want %v", contract.Coverage, expectedCoverage))
	}
	for _, term := range []string{
		"WithImportSheetName",
		"WithImportSheetIndex",
		"ErrSheetIndexOutOfRange",
		"WithoutHeader",
		"WithoutTrimSpace",
		"NewMapImporter",
		"[]tabular.MapOption",
		"native typed cell",
		"timex.Time",
		"FormatterFn",
	} {
		if !contains(contract.Terms, term) {
			failures = append(failures, "excel contract terms missing "+term)
		}
	}

	allEvidence := append([]string{}, review.SourceEvidence...)
	allEvidence = append(allEvidence, contract.SourceEvidence...)
	allEvidence = append(allEvidence, contract.TestEvidence...)
	for _, item := range allEvidence {
		path, lineText, ok := strings.Cut(item, ":")
		if !ok || lineText == "" {
			failures = append(failures, "excel contract evidence missing line number: "+item)
			continue
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, "excel contract evidence missing file: "+item)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	files := map[string]string{
		"excel.go":         readSourceFile(sourceRoot, "excel/excel.go").content,
		"options.go":       readSourceFile(sourceRoot, "excel/options.go").content,
		"importer.go":      readSourceFile(sourceRoot, "excel/importer.go").content,
		"exporter.go":      readSourceFile(sourceRoot, "excel/exporter.go").content,
		"errors.go":        readSourceFile(sourceRoot, "excel/errors.go").content,
		"excel_test.go":    readSourceFile(sourceRoot, "excel/excel_test.go").content,
		"options_test.go":  readSourceFile(sourceRoot, "excel/options_test.go").content,
		"importer_test.go": readSourceFile(sourceRoot, "excel/importer_test.go").content,
		"exporter_test.go": readSourceFile(sourceRoot, "excel/exporter_test.go").content,
	}

	checks := []struct {
		file string
		term string
	}{
		{"errors.go", `ErrSheetIndexOutOfRange = errors.New("sheet index out of range")`},
		{"excel.go", `func NewMapImporter(`},
		{"excel.go", `specs []tabular.ColumnSpec, mapOpts []tabular.MapOption, opts ...ImportOption`},
		{"excel.go", `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)`},
		{"excel.go", `tabular.NewTypedImporter[T](NewImporterFor[T](opts...))`},
		{"excel.go", `tabular.NewTypedExporter[T](NewExporterFor[T](opts...))`},
		{"options.go", `type ExportOption func(*exportConfig)`},
		{"options.go", `type ImportOption func(*importConfig)`},
		{"options.go", `o.sheetName = name`},
		{"options.go", `o.sheetIndex = index`},
		{"options.go", `o.skipRows = max(rows, 0)`},
		{"options.go", `o.hasHeader = false`},
		{"options.go", `o.trimSpace = false`},
		{"importer.go", `sheetIndex: 0`},
		{"importer.go", `hasHeader:  true`},
		{"importer.go", `trimSpace:  true`},
		{"importer.go", `if sheetName == ""`},
		{"importer.go", `ErrSheetIndexOutOfRange`},
		{"importer.go", `tabular.ErrNoDataRowsFound`},
		{"importer.go", `tabular.BuildHeaderMapping(`},
		{"importer.go", `tabular.DefaultPositionalMapping(schema)`},
		{"importer.go", `tabular.IsEmptyRow(row, i.options.trimSpace)`},
		{"importer.go", `tabular.ParseRowOptions{TrimSpace: i.options.trimSpace}`},
		{"exporter.go", `sheetName: "Sheet1"`},
		{"exporter.go", `f.SetSheetName("Sheet1", e.options.sheetName)`},
		{"exporter.go", `column.Format == "" && tabular.IsDefaultFormatter(column, e.formatters)`},
		{"exporter.go", `nativeCellValue(raw)`},
		{"exporter.go", `case timex.DateTime:`},
		{"exporter.go", `case timex.Date:`},
		{"exporter.go", `timex.Time (time-of-day) is deliberately left unwrapped`},
		{"excel_test.go", `NewMapImporter(`},
		{"excel_test.go", `[]tabular.MapOption{tabular.WithRowValidator`},
		{"options_test.go", `WithSkipRowsClampsNegative`},
		{"options_test.go", `WithoutHeader`},
		{"options_test.go", `WithoutTrimSpace`},
		{"importer_test.go", `WithImportSheetName`},
		{"importer_test.go", `EmptyRowsSkipped`},
		{"exporter_test.go", `TestExporterNativeCellTypes`},
		{"exporter_test.go", `ExplicitFormatColumnStaysString`},
		{"exporter_test.go", `CustomFormatterColumnStaysString`},
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
	return runCommand(sourceRoot, "go", "test", "./excel")
}

func loadExcelAuditEntries(path string) []auditEntry {
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
		if entry.Package == excelPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})
	if len(entries) == 0 {
		panic("API audit ledger missing excel entries")
	}

	return entries
}

func loadExcelManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == excelPackage {
			return entry
		}
	}

	panic("API audit manifest missing excel package")
}

func loadExcelContract(path string) (contractPackageReview, contractEntry) {
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
		if item.Package == excelPackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing excel package review")
	}

	var contract contractEntry
	contractFound := false
	for _, item := range ledger.Entries {
		if item.ID == excelContractID() {
			contract = item
			contractFound = true

			break
		}
	}
	if !contractFound {
		panic("API contract ledger missing excel contract entry")
	}

	return review, contract
}

func loadLiveExcelEntry(sourceRoot, docsRoot string) manifestEntry {
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
		if entry.Package == excelPackage {
			return entry
		}
	}

	panic("live API inventory missing excel package")
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

func excelReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`excel\.[A-Z][A-Za-z0-9_]*`)
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

func excelContractID() string {
	return excelPackage + "#runtime-contract:excel-options-and-native-cells"
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
