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
	schemaPackage = "github.com/coldsmirk/vef-framework-go/schema"

	schemaFingerprint = "7e874740c91b2aa73d968ffe782ae27256691b6611fa235973111940f27f188e"
	schemaTopLevel    = 12
	schemaFields      = 38
	schemaMethods     = 3
	schemaEntries     = 53

	schemaGroupedEntries              = 41
	schemaGroupedFields               = 38
	schemaGroupedMethods              = 3
	schemaGroupedReceivers            = 10
	schemaGroupedSignatureFingerprint = "fbb57af25a5c2f3364726e554b1c7cd5b6adf69c15e8e2f9877efe6a7a491fda"
	schemaGroupedReceiverFingerprint  = "06a6a4cd745e4651d3f9962cd32f1f1c07233320ee735e08e3d8e05c9adda521"

	englishSchemaPath   = "docs/features/schema.md"
	chineseSchemaPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/schema.md"
	englishBuiltInsPath = "docs/reference/built-in-resources.md"
	chineseBuiltInsPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/built-in-resources.md"
	englishIndexPath    = "docs/reference/public-api-index.md"
	chineseIndexPath    = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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

type runtimeLedger struct {
	Entries []runtimeEntry `json:"entries"`
}

type runtimeEntry struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	Value    string `json:"value"`
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

	englishDocs := readCorpus("English schema docs", filepath.Join(docsRoot, englishSchemaPath))
	chineseDocs := readCorpus("Chinese schema docs", filepath.Join(docsRoot, chineseSchemaPath))
	englishBuiltIns := readCorpus("English built-in resources docs", filepath.Join(docsRoot, englishBuiltInsPath))
	chineseBuiltIns := readCorpus("Chinese built-in resources docs", filepath.Join(docsRoot, chineseBuiltInsPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	runtime := loadJSON[runtimeLedger](filepath.Join(docsRoot, "scripts/runtime-api-ledger.json"))
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[schemaPackage]
	liveAuditEntries := schemaEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	schemaEntries := schemaEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(schemaEntries)...)
	failures = append(failures, verifyLiveAuditEntries(schemaEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedSchemaSurface(schemaEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifyRuntimeActions(runtime, []corpus{englishDocs, chineseDocs, englishBuiltIns, chineseBuiltIns})...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, schemaEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, schemaEntries)...)
	failures = append(failures, verifySchemaDocs(schemaEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Schema contract docs verified: 53 public entries, 41 grouped DTO/service entries, sys/schema runtime actions")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != schemaPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != schemaTopLevel ||
		entry.Fields != schemaFields ||
		entry.Methods != schemaMethods ||
		entry.Fingerprint != schemaFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s schema surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			schemaTopLevel, schemaFields, schemaMethods, schemaFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != schemaPackage {
			continue
		}

		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, schemaCoverage()) {
			failures = append(failures, fmt.Sprintf("schema manifest coverage mismatch: got %v want %v", entry.Coverage, schemaCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing schema package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	expectedContracts := map[string][]string{
		schemaPackage + "#dynamic-resource:sys-schema": {
			"sys/schema",
			"list_tables",
			"get_table_schema",
			"list_views",
			"60",
		},
		schemaPackage + "#field-contract:schema-dto-wire-shape": {
			"omitempty",
			"primaryKey",
			"isAutoIncrement",
			"onDelete",
			"columns",
		},
		schemaPackage + "#runtime-contract:primary-datasource-inspection": {
			"primary data source",
			"PostgreSQL",
			"MySQL",
			"SQLite",
			"information_schema.views",
		},
	}

	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != schemaPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "schema contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != schemaTopLevel ||
			review.ReviewedSurface.Fields != schemaFields ||
			review.ReviewedSurface.Methods != schemaMethods ||
			review.ReviewedSurface.EntryCount != schemaEntries ||
			review.ReviewedSurface.Fingerprint != schemaFingerprint {
			failures = append(failures, fmt.Sprintf(
				"schema contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, schemaCoverage()) {
			failures = append(failures, fmt.Sprintf("schema contract review coverage mismatch: got %v want %v", review.Coverage, schemaCoverage()))
		}
		if !sameSet(review.ContractIDs, sortedKeys(expectedContracts)) {
			failures = append(failures, fmt.Sprintf("schema contract ids mismatch: got %v want %v", review.ContractIDs, sortedKeys(expectedContracts)))
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing schema package review")
	}

	foundContracts := map[string]bool{}
	for _, entry := range contracts.Entries {
		terms, ok := expectedContracts[entry.ID]
		if !ok {
			continue
		}
		foundContracts[entry.ID] = true
		if entry.Package != schemaPackage {
			failures = append(failures, fmt.Sprintf("schema contract entry package mismatch for %s: %s", entry.ID, entry.Package))
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "schema contract entry disposition mismatch for "+entry.ID+": "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, schemaCoverageForContract(entry.ID)) {
			failures = append(failures, fmt.Sprintf("schema contract coverage mismatch for %s: got %v want %v", entry.ID, entry.Coverage, schemaCoverageForContract(entry.ID)))
		}
		for _, term := range terms {
			if !contains(entry.Terms, term) {
				failures = append(failures, fmt.Sprintf("schema contract %s missing term %s", entry.ID, term))
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.TestEvidence)...)
	}
	for id := range expectedContracts {
		if !foundContracts[id] {
			failures = append(failures, "contract ledger missing schema contract entry "+id)
		}
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != schemaEntries {
		failures = append(failures, fmt.Sprintf("schema audit entry count mismatch: got %d want %d", len(entries), schemaEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != schemaPackage {
			failures = append(failures, "non-schema audit entry passed into schema verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate schema audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "schema audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, schemaCoverage()) {
			failures = append(failures, fmt.Sprintf("schema audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, schemaCoverage()))
		}
	}
	if counts["top"] != schemaTopLevel || counts["field"] != schemaFields || counts["method"] != schemaMethods {
		failures = append(failures, fmt.Sprintf("schema audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != schemaTopLevel ||
		dispositionCounts["grouped:type-member-family"] != schemaGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"schema audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			schemaTopLevel,
			schemaGroupedEntries,
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
			failures = append(failures, fmt.Sprintf("schema missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"schema live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
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
			failures = append(failures, fmt.Sprintf("schema extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedSchemaSurface(entries []auditEntry, docs []corpus) []string {
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
			failures = append(failures, fmt.Sprintf("schema grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("schema grouped type-member surface", rows, schemaGroupedEntries, schemaGroupedSignatureFingerprint)...)
	if kindCounts["field"] != schemaGroupedFields || kindCounts["method"] != schemaGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"schema grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			schemaGroupedFields,
			schemaGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("schema grouped receiver/type families", receiverRows, schemaGroupedReceivers, schemaGroupedReceiverFingerprint)...)

	for _, doc := range docs {
		for _, term := range []string{
			"53 public schema entries",
			"41 grouped schema field/method entries",
			"10 schema receiver/type families",
			"38 exported schema field entries",
			"3 exported schema method entries",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped schema audit term "+term)
			}
		}
	}

	return failures
}

func verifyRuntimeActions(runtime runtimeLedger, docs []corpus) []string {
	actions := map[string]bool{}
	var resourceFound bool
	for _, entry := range runtime.Entries {
		if entry.Category == "built-in resource" && entry.Name == "sys/schema" && entry.Value == "rpc" {
			resourceFound = true
		}
		if entry.Category == "built-in resource action" && strings.HasPrefix(entry.Name, "sys/schema/") {
			actions[strings.TrimPrefix(entry.Name, "sys/schema/")] = true
		}
	}

	wantActions := []string{"get_table_schema", "list_tables", "list_views"}
	var failures []string
	if !resourceFound {
		failures = append(failures, "runtime ledger missing sys/schema built-in resource")
	}
	failures = append(failures, compareSets("runtime sys/schema actions", sortedKeys(actions), wantActions)...)
	for _, doc := range docs {
		if !strings.Contains(doc.content, "`sys/schema`") {
			failures = append(failures, doc.label+" missing sys/schema")
		}
		for _, action := range wantActions {
			if !strings.Contains(doc.content, "`"+action+"`") {
				failures = append(failures, doc.label+" missing schema action "+action)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, schemaPackage)
	if section == "" {
		return []string{index.label + " missing schema package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s schema index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifySchemaDocs(entries []auditEntry, docs []corpus) []string {
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
			if !strings.Contains(doc.content, "`"+symbol+"`") &&
				!strings.Contains(doc.content, "`schema."+symbol+"`") {
				failures = append(failures, doc.label+" missing top-level schema symbol `"+symbol+"`")
			}
		}
		for _, term := range []string{
			"primary data source",
			"PostgreSQL",
			"MySQL",
			"SQLite",
			"`information_schema.views`",
			"`DATABASE()`",
			"`sqlite_schema`",
			"`sqlite_%`",
			"`omitempty`",
			"`primaryKey`",
			"`isAutoIncrement`",
			"`onDelete`",
			"`ErrCodeTableNotFound`",
			"`ErrTableNotFound`",
			"`Max = 60`",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing schema semantic term "+term)
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
			path: "schema/service.go",
			terms: []string{
				"type Table struct",
				"type TableSchema struct",
				"type Column struct",
				"`json:\"primaryKey,omitempty\"`",
				"`json:\"isAutoIncrement,omitempty\"`",
				"`json:\"onDelete,omitempty\"`",
				"type Service interface",
				"ListTables(ctx context.Context)",
				"GetTableSchema(ctx context.Context, name string)",
				"ListViews(ctx context.Context)",
			},
		},
		{
			path: "schema/api_errors.go",
			terms: []string{
				"ErrCodeTableNotFound = 2300",
				"ErrTableNotFound = result.Err(",
				"i18n.T(\"schema_table_not_found\")",
			},
		},
		{
			path: "internal/schema/resource.go",
			terms: []string{
				"api.NewRPCResource(",
				"\"sys/schema\"",
				"Action: \"list_tables\", RateLimit: &api.RateLimitConfig{Max: 60}",
				"Action: \"get_table_schema\", RateLimit: &api.RateLimitConfig{Max: 60}",
				"Action: \"list_views\", RateLimit: &api.RateLimitConfig{Max: 60}",
				"Name string `json:\"name\" validate:\"required\"`",
				"return schema.ErrTableNotFound",
			},
		},
		{
			path: "internal/schema/service.go",
			terms: []string{
				"primary := dataSources.Primary()",
				"NewInspector(db, primary.Kind, primary.Schema)",
				"convertTable(table)",
				"IsAutoIncrement: hasAutoIncrement(col)",
				"referentialActionToString",
				"case \"serial\", \"bigserial\", \"smallserial\"",
			},
		},
		{
			path: "internal/schema/inspector.go",
			terms: []string{
				"schema = lo.CoalesceOrEmpty(schemaName, \"public\")",
				"schema = \"main\"",
				"information_schema.views AS v",
				"v.table_schema = DATABASE()",
				"sqlite_schema",
				"name NOT LIKE 'sqlite_%'",
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
	cmd := exec.Command("go", "test", "./schema", "./internal/schema")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./schema ./internal/schema failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func schemaEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == schemaPackage {
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

func compareSets(label string, got, want []string) []string {
	got = sortedUnique(got)
	want = sortedUnique(want)
	gotSet := sliceSet(got)
	wantSet := sliceSet(want)

	var failures []string
	for _, item := range want {
		if !gotSet[item] {
			failures = append(failures, label+" missing "+item)
		}
	}
	for _, item := range got {
		if !wantSet[item] {
			failures = append(failures, label+" has phantom "+item)
		}
	}

	return failures
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

func schemaCoverage() []string {
	return []string{englishSchemaPath, englishBuiltInsPath}
}

func schemaCoverageForContract(id string) []string {
	if strings.HasSuffix(id, "#runtime-contract:primary-datasource-inspection") {
		return []string{englishSchemaPath}
	}

	return schemaCoverage()
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
