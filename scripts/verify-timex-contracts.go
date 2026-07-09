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
	timexPackage = "github.com/coldsmirk/vef-framework-go/timex"

	timexFingerprint = "760ac8fb520f880a74c370205f924a149e2450a97cbb2c894425f2497d6a96e4"
	timexTopLevel    = 20
	timexFields      = 0
	timexMethods     = 136
	timexEntries     = 156

	timexGroupedEntries              = 136
	timexGroupedMethods              = 136
	timexGroupedReceivers            = 3
	timexGroupedSignatureFingerprint = "019c36b4c473fee9e3141aac255efd34b802467315ce3f88636bcaf3f2ff0c7f"
	timexGroupedReceiverFingerprint  = "5912e53f6b7f0fccca6f9b02332f57c3063ebe523d3f5a3f33717a734a9d2d2d"

	englishTimexPath = "docs/utilities/timex.md"
	chineseTimexPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/timex.md"
	englishIndexPath = "docs/reference/public-api-index.md"
	chineseIndexPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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

	englishDocs := readCorpus("English timex docs", filepath.Join(docsRoot, englishTimexPath))
	chineseDocs := readCorpus("Chinese timex docs", filepath.Join(docsRoot, chineseTimexPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[timexPackage]
	liveAuditEntries := timexEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	timexEntries := timexEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(timexEntries)...)
	failures = append(failures, verifyLiveAuditEntries(timexEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedTimexSurface(timexEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, timexEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, timexEntries)...)
	failures = append(failures, verifyTimexDocs(timexEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Timex contract docs verified: 156 public entries, 136 grouped methods, Date/Time/DateTime wire contracts")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != timexPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != timexTopLevel ||
		entry.Fields != timexFields ||
		entry.Methods != timexMethods ||
		entry.Fingerprint != timexFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s timex surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			timexTopLevel, timexFields, timexMethods, timexFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != timexPackage {
			continue
		}

		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, timexCoverage()) {
			failures = append(failures, fmt.Sprintf("timex manifest coverage mismatch: got %v want %v", entry.Coverage, timexCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing timex package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != timexPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "timex contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != timexTopLevel ||
			review.ReviewedSurface.Fields != timexFields ||
			review.ReviewedSurface.Methods != timexMethods ||
			review.ReviewedSurface.EntryCount != timexEntries ||
			review.ReviewedSurface.Fingerprint != timexFingerprint {
			failures = append(failures, fmt.Sprintf(
				"timex contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, timexCoverage()) {
			failures = append(failures, fmt.Sprintf("timex contract review coverage mismatch: got %v want %v", review.Coverage, timexCoverage()))
		}
		if !contains(review.ContractIDs, timexContractID()) {
			failures = append(failures, "timex contract review missing contract id "+timexContractID())
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing timex package review")
	}

	var foundContract bool
	for _, entry := range contracts.Entries {
		if entry.ID != timexContractID() {
			continue
		}
		foundContract = true
		if entry.Package != timexPackage || entry.Kind != "wire-contract" {
			failures = append(failures, fmt.Sprintf("timex contract entry shape mismatch: package=%s kind=%s", entry.Package, entry.Kind))
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "timex contract entry disposition mismatch: "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, timexCoverage()) {
			failures = append(failures, fmt.Sprintf("timex contract coverage mismatch: got %v want %v", entry.Coverage, timexCoverage()))
		}
		for _, term := range []string{"time.DateTime", "time.DateOnly", "time.TimeOnly", "Between", "json.Marshaler"} {
			if !contains(entry.Terms, term) {
				failures = append(failures, "timex contract missing term "+term)
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
	}
	if !foundContract {
		failures = append(failures, "contract ledger missing timex contract entry")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != timexEntries {
		failures = append(failures, fmt.Sprintf("timex audit entry count mismatch: got %d want %d", len(entries), timexEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != timexPackage {
			failures = append(failures, "non-timex audit entry passed into timex verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate timex audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "timex audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, timexCoverage()) {
			failures = append(failures, fmt.Sprintf("timex audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, timexCoverage()))
		}
	}
	if counts["top"] != timexTopLevel || counts["field"] != timexFields || counts["method"] != timexMethods {
		failures = append(failures, fmt.Sprintf("timex audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != timexTopLevel ||
		dispositionCounts["grouped:type-member-family"] != timexGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"timex audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			timexTopLevel,
			timexGroupedEntries,
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
			failures = append(failures, fmt.Sprintf("timex missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"timex live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
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
			failures = append(failures, fmt.Sprintf("timex extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedTimexSurface(entries []auditEntry) []string {
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
			failures = append(failures, fmt.Sprintf("timex grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("timex grouped type-member surface", rows, timexGroupedEntries, timexGroupedSignatureFingerprint)...)
	if kindCounts["method"] != timexGroupedMethods || kindCounts["field"] != 0 {
		failures = append(failures, fmt.Sprintf("timex grouped kind counts mismatch: got fields/methods=%d/%d want=0/%d", kindCounts["field"], kindCounts["method"], timexGroupedMethods))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("timex grouped receiver/type families", receiverRows, timexGroupedReceivers, timexGroupedReceiverFingerprint)...)

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, timexPackage)
	if section == "" {
		return []string{index.label + " missing timex package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s timex index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyTimexDocs(entries []auditEntry, docs []corpus) []string {
	var topSymbols []string
	for _, entry := range entries {
		if entry.Kind == "top" {
			topSymbols = append(topSymbols, entry.Symbol)
		}
	}
	sort.Strings(topSymbols)

	var failures []string
	for _, doc := range docs {
		for _, symbol := range topSymbols {
			if !strings.Contains(doc.content, "`"+symbol+"`") {
				failures = append(failures, doc.label+" missing top-level timex symbol `"+symbol+"`")
			}
		}
		for _, term := range []string{
			"`time.DateTime`",
			"`time.DateOnly`",
			"`time.TimeOnly`",
			"`Parse`",
			"`ParseDate`",
			"`ParseTime`",
			"`MarshalJSON`",
			"`UnmarshalJSON`",
			"`MarshalText`",
			"`UnmarshalText`",
			"`Scan`",
			"`Value`",
			"`Between`",
			"open interval",
			"Sunday",
			"Saturday",
			"`Unwrap`",
			"`Format`",
			"`String`",
			"`UnixMilli`",
			"`UnixMicro`",
			"`UnixNano`",
			"`ToDuration`",
			"no `T` separator",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing timex semantic term "+term)
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
			path: "timex/common.go",
			terms: []string{
				"dateTimeLayout = time.DateTime",
				"dateLayout     = time.DateOnly",
				"timeLayout     = time.TimeOnly",
				"parseTimeWithFallback",
				"cast.ToTimeInDefaultLocationE",
				"unquoteJSON",
			},
		},
		{
			path: "timex/date.go",
			terms: []string{
				"type Date time.Time",
				"Date is a civil date",
				"return ErrInvalidDateFormat",
				"func (d Date) Between(start, end Date) bool",
				"return d.After(start) && d.Before(end)",
				"BeginOfWeek returns the beginning of the week (Sunday)",
				"EndOfWeek returns the end of the week (Saturday)",
				"func ParseDate(value string, pattern ...string) (Date, error)",
			},
		},
		{
			path: "timex/time.go",
			terms: []string{
				"type Time time.Time",
				"Time is a civil clock value",
				"return ErrInvalidTimeFormat",
				"func (t Time) ToDuration() time.Duration",
				"return t.After(start) && t.Before(end)",
				"func ParseTime(value string, pattern ...string) (Time, error)",
			},
		},
		{
			path: "timex/datetime.go",
			terms: []string{
				"type DateTime time.Time",
				"return ErrInvalidDateTimeFormat",
				"func (dt DateTime) Between(start, end DateTime) bool",
				"return dt.After(start) && dt.Before(end)",
				"func FromUnixMilli(msec int64) DateTime",
				"func FromUnixMicro(usec int64) DateTime",
				"func Parse(value string, pattern ...string) (DateTime, error)",
			},
		},
		{
			path: "timex/errors.go",
			terms: []string{
				"ErrInvalidDateFormat = errors.New(\"invalid date format\")",
				"ErrInvalidDateTimeFormat = errors.New(\"invalid datetime format\")",
				"ErrInvalidTimeFormat = errors.New(\"invalid time format\")",
				"ErrFailedScan = errors.New(\"failed to scan value\")",
				"ErrUnsupportedDestType = errors.New(\"unsupported destination type\")",
			},
		},
		{
			path: "timex/date_test.go",
			terms: []string{
				"Between is exclusive of the start endpoint",
				"Between is exclusive of the end endpoint",
				"UnmarshalText must share UnmarshalJSON's strictness",
			},
		},
		{
			path: "timex/datetime_test.go",
			terms: []string{
				"Between is exclusive of the start endpoint",
				"Between is exclusive of the end endpoint",
			},
		},
		{
			path: "timex/time_test.go",
			terms: []string{
				"Between is exclusive of the start endpoint",
				"Between is exclusive of the end endpoint",
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
	cmd := exec.Command("go", "test", "./timex")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./timex failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func timexEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == timexPackage {
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

func timexContractID() string {
	return timexPackage + "#wire-contract:date-time-json-text-format"
}

func timexCoverage() []string {
	return []string{englishTimexPath}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
