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
	searchPackage = "github.com/coldsmirk/vef-framework-go/search"

	searchFingerprint = "e5daf3aeeade69b48eb2c0419003c930df8746c739eab98b1e131c22f9424cab"
	searchTopLevel    = 43
	searchFields      = 0
	searchMethods     = 1
	searchEntries     = 44

	searchGroupedEntries              = 1
	searchGroupedFields               = 0
	searchGroupedMethods              = 1
	searchGroupedReceivers            = 1
	searchGroupedSignatureFingerprint = "8394796f33a5ff2bd6063f3809eba0cc0215d9665bf4394f53901a6f64aecb6f"
	searchGroupedReceiverFingerprint  = "ed72604c91cf6fc5e4aa124631f1f7f3430b41a0994e75a313ed615157ddf6e2"

	englishQueryBuilderPath = "docs/guide/query-builder.md"
	chineseQueryBuilderPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/query-builder.md"
	englishIndexPath        = "docs/reference/public-api-index.md"
	chineseIndexPath        = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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

	englishDocs := readCorpus("English query builder docs", filepath.Join(docsRoot, englishQueryBuilderPath))
	chineseDocs := readCorpus("Chinese query builder docs", filepath.Join(docsRoot, chineseQueryBuilderPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[searchPackage]
	liveAuditEntries := searchEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	searchEntries := searchEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(searchEntries)...)
	failures = append(failures, verifyLiveAuditEntries(searchEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedSearchSurface(searchEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, searchEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, searchEntries)...)
	failures = append(failures, verifySearchDocs(searchEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Search contract docs verified: 44 public entries, 1 grouped Search.Apply method, tag grammar and apply semantics")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != searchPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != searchTopLevel ||
		entry.Fields != searchFields ||
		entry.Methods != searchMethods ||
		entry.Fingerprint != searchFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s search surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			searchTopLevel, searchFields, searchMethods, searchFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != searchPackage {
			continue
		}
		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, searchCoverage()) {
			failures = append(failures, fmt.Sprintf("search manifest coverage mismatch: got %v want %v", entry.Coverage, searchCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing search package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	contractID := searchPackage + "#string-contract:search-tag-grammar-and-operators"
	expectedTerms := []string{
		"search:\"",
		"operator=",
		"column=",
		"params=",
		"Search.Apply",
		"api.P",
		"WithSpacePairDelimiter",
		"monad.Range",
		"notBetween",
		"iContains",
		"isNull",
	}

	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != searchPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "search contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != searchTopLevel ||
			review.ReviewedSurface.Fields != searchFields ||
			review.ReviewedSurface.Methods != searchMethods ||
			review.ReviewedSurface.EntryCount != searchEntries ||
			review.ReviewedSurface.Fingerprint != searchFingerprint {
			failures = append(failures, fmt.Sprintf(
				"search contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, searchCoverage()) {
			failures = append(failures, fmt.Sprintf("search contract review coverage mismatch: got %v want %v", review.Coverage, searchCoverage()))
		}
		if !sameSet(review.ContractIDs, []string{contractID}) {
			failures = append(failures, fmt.Sprintf("search contract ids mismatch: got %v want %v", review.ContractIDs, []string{contractID}))
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing search package review")
	}

	var foundEntry bool
	for _, entry := range contracts.Entries {
		if entry.ID != contractID {
			continue
		}
		foundEntry = true
		if entry.Package != searchPackage {
			failures = append(failures, "search contract entry package mismatch: "+entry.Package)
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "search contract entry disposition mismatch: "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, searchCoverage()) {
			failures = append(failures, fmt.Sprintf("search contract coverage mismatch: got %v want %v", entry.Coverage, searchCoverage()))
		}
		for _, term := range expectedTerms {
			if !contains(entry.Terms, term) {
				failures = append(failures, "search contract entry missing term "+term)
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
	}
	if !foundEntry {
		failures = append(failures, "contract ledger missing search contract entry "+contractID)
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != searchEntries {
		failures = append(failures, fmt.Sprintf("search audit entry count mismatch: got %d want %d", len(entries), searchEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != searchPackage {
			failures = append(failures, "non-search audit entry passed into search verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate search audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "search audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, searchCoverage()) {
			failures = append(failures, fmt.Sprintf("search audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, searchCoverage()))
		}
	}
	if counts["top"] != searchTopLevel || counts["field"] != searchFields || counts["method"] != searchMethods {
		failures = append(failures, fmt.Sprintf("search audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != searchTopLevel ||
		dispositionCounts["grouped:builder-member-family"] != searchGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"search audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:builder-member-family"],
			searchTopLevel,
			searchGroupedEntries,
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
			failures = append(failures, fmt.Sprintf("search missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"search live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
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
			failures = append(failures, fmt.Sprintf("search extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedSearchSurface(entries []auditEntry, docs []corpus) []string {
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
			failures = append(failures, fmt.Sprintf("search grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("search grouped type-member surface", rows, searchGroupedEntries, searchGroupedSignatureFingerprint)...)
	if kindCounts["field"] != searchGroupedFields || kindCounts["method"] != searchGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"search grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			searchGroupedFields,
			searchGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("search grouped receiver/type families", receiverRows, searchGroupedReceivers, searchGroupedReceiverFingerprint)...)

	for _, doc := range docs {
		for _, term := range []string{
			"44 public search entries",
			"1 grouped search method entry",
			"1 search receiver/type family",
			"0 exported search field entries",
			"1 exported search method entry",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped search audit term "+term)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, searchPackage)
	if section == "" {
		return []string{index.label + " missing search package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s search index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifySearchDocs(entries []auditEntry, docs []corpus) []string {
	var topSymbols []string
	for _, entry := range entries {
		if entry.Kind == "top" {
			topSymbols = append(topSymbols, entry.Symbol)
		}
	}
	sort.Strings(topSymbols)

	terms := []string{
		"`search.Operator`",
		"`Search.Apply(...)`",
		"`search.New`",
		"`search.NewFor[T]`",
		"`search.Applier`",
		"`TagSearch`",
		"`IgnoreField`",
		"`AttrOperator`",
		"`AttrColumn`",
		"`AttrAlias`",
		"`AttrParams`",
		"`AttrDive`",
		"`ParamDelimiter`",
		"`ParamType`",
		"`TypeInt`",
		"`TypeDecimal`",
		"`TypeDate`",
		"`TypeDateTime`",
		"`TypeTime`",
		"`eq`",
		"`notBetween`",
		"`iContains`",
		"`isNull`",
		"`isNotNull`",
		"`params=delimiter:| type:int`",
		"`api.P`",
		"`monad.Range[T]`",
		"`type:int`",
		"`type:date`",
		"`dive`",
		"`-`",
		"nil pointer",
		"Unknown operator",
	}

	var failures []string
	for _, doc := range docs {
		for _, symbol := range topSymbols {
			if !strings.Contains(doc.content, "`"+symbol+"`") &&
				!strings.Contains(doc.content, "`search."+symbol+"`") &&
				!strings.Contains(doc.content, "search."+symbol) &&
				!strings.Contains(doc.content, symbol+"(") {
				failures = append(failures, doc.label+" missing top-level search symbol `"+symbol+"`")
			}
		}
		for _, term := range terms {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing search semantic term "+term)
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
			path: "search/constants.go",
			terms: []string{
				"type Operator string",
				"Equals             Operator = \"eq\"",
				"NotEquals          Operator = \"neq\"",
				"GreaterThan        Operator = \"gt\"",
				"GreaterThanOrEqual Operator = \"gte\"",
				"LessThan           Operator = \"lt\"",
				"LessThanOrEqual    Operator = \"lte\"",
				"Between    Operator = \"between\"",
				"NotBetween Operator = \"notBetween\"",
				"In    Operator = \"in\"",
				"NotIn Operator = \"notIn\"",
				"IsNull    Operator = \"isNull\"",
				"IsNotNull Operator = \"isNotNull\"",
				"Contains      Operator = \"contains\"",
				"NotContains   Operator = \"notContains\"",
				"StartsWith    Operator = \"startsWith\"",
				"NotStartsWith Operator = \"notStartsWith\"",
				"EndsWith      Operator = \"endsWith\"",
				"NotEndsWith   Operator = \"notEndsWith\"",
				"ContainsIgnoreCase      Operator = \"iContains\"",
				"NotContainsIgnoreCase   Operator = \"iNotContains\"",
				"StartsWithIgnoreCase    Operator = \"iStartsWith\"",
				"NotStartsWithIgnoreCase Operator = \"iNotStartsWith\"",
				"EndsWithIgnoreCase      Operator = \"iEndsWith\"",
				"NotEndsWithIgnoreCase   Operator = \"iNotEndsWith\"",
				"TagSearch = \"search\"",
				"AttrDive     = \"dive\"",
				"AttrAlias    = \"alias\"",
				"AttrColumn   = \"column\"",
				"AttrOperator = \"operator\"",
				"AttrParams   = \"params\"",
				"ParamDelimiter = \"delimiter\"",
				"ParamType      = \"type\"",
				"IgnoreField = \"-\"",
				"TypeInt      = \"int\"",
				"TypeDecimal  = \"dec\"",
				"TypeDate     = \"date\"",
				"TypeDateTime = \"datetime\"",
				"TypeTime     = \"time\"",
			},
		},
		{
			path: "search/parser.go",
			terms: []string{
				"var apiInType = reflect.TypeFor[api.P]()",
				"func New(typ reflect.Type) Search",
				"typ = reflectx.Indirect(typ)",
				"return Search{}",
				"func NewFor[T any]() Search",
				"return New(reflect.TypeFor[T]())",
				"if field.Anonymous && field.Type == apiInType",
				"if tag == IgnoreField",
				"if tag == AttrDive",
				"attrs := strx.ParseTag(tag)",
				"reflectx.WithDiveTag(TagSearch, AttrDive)",
				"[]string{lo.SnakeCase(field.Name)}",
				"strings.Split(column, \"|\")",
				"operator := lo.CoalesceOrEmpty(attrs[AttrOperator], attrs[strx.DefaultKey], string(Equals))",
				"params = strx.ParseTag(attrs[AttrParams]",
				"strx.WithSpacePairDelimiter()",
				"strx.WithValueDelimiter(':')",
			},
		},
		{
			path: "search/search.go",
			terms: []string{
				"func (f Search) Apply(cb orm.ConditionBuilder, target any, defaultAlias ...string)",
				"if value.Kind() != reflect.Struct",
				"return",
				"if field.Kind() == reflect.Pointer && field.IsNil()",
				"alias := getColumnAlias(c.alias, defaultAlias...)",
				"dbx.ColumnWithAlias(column, alias)",
				"case Equals, NotEquals, GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual:",
				"case Between, NotBetween:",
				"case In, NotIn:",
				"case IsNull, IsNotNull:",
				"Unknown operator %q for columns %v, condition ignored",
				"delimiter := lo.CoalesceOrEmpty(conditionParams[ParamDelimiter], \",\")",
				"for value := range strings.SplitSeq(slice, delimiter)",
				"case TypeInt:",
				"cb.ApplyIf(shouldApply",
				"if content == \"\"",
				"cb.Group(func(cb orm.ConditionBuilder)",
				"applyLikeMethod(useOr, cb.OrContains, cb.Contains, column, content)",
			},
		},
		{
			path: "search/range.go",
			terms: []string{
				"rangeType            = reflect.TypeFor[monad.Range[int]]()",
				"reflectx.IsSimilarType(valueType, rangeType)",
				"return parseStringRange(value.String(), conditionParams)",
				"return parseSliceRange(value)",
				"delimiter := lo.CoalesceOrEmpty(conditionParams[ParamDelimiter], \",\")",
				"strings.SplitN(value, delimiter, 2)",
				"TypeInt:      parseIntRange",
				"TypeDecimal:  parseDecimalRange",
				"TypeDate:     parseDateRange",
				"TypeTime:     parseTimeRange",
				"TypeDateTime: parseDateTimeRange",
				"if value.Len() != 2",
				"time.DateOnly",
				"time.TimeOnly",
				"time.DateTime",
			},
		},
		{
			path: "search/applier.go",
			terms: []string{
				"func Applier[T any]() func(T) orm.ApplyFunc[orm.ConditionBuilder]",
				"f := NewFor[T]()",
				"f.Apply(cb, value)",
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
	cmd := exec.Command("go", "test", "./search")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./search failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func searchEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == searchPackage {
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
	set := map[string]bool{}
	for _, value := range values {
		set[value] = true
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)

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

func searchCoverage() []string {
	return []string{englishQueryBuilderPath}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
