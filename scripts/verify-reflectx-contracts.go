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
	reflectxPackage     = "github.com/coldsmirk/vef-framework-go/reflectx"
	reflectxFingerprint = "bb62b3bd50f5b54c5af99deb16b7cfb61fa52e69f92e3ab789dc81c744f6d3de"
	reflectxTopLevel    = 78
	reflectxFields      = 11
	reflectxMethods     = 0
	reflectxEntries     = 89

	englishSmallUtilitiesPath = "docs/utilities/small-utilities.md"
	chineseSmallUtilitiesPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/small-utilities.md"
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
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadReflectxAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadReflectxManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contracts := loadReflectxContracts(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveReflectxEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contracts)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
	}
	for _, doc := range []corpus{englishSmallUtilities, chineseSmallUtilities} {
		failures = append(failures, verifyDocumentedSurface(doc, entries)...)
		failures = append(failures, verifyNoPhantomReflectxRefs(doc, entries)...)
		failures = append(failures, verifyNoPhantomReflectxFields(doc, entries)...)
		failures = append(failures, missingTerms(doc, smallUtilitiesTerms())...)
	}

	failures = append(failures, verifyContractLedger(review, contracts, sourceRoot)...)
	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("reflectx contracts verified")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != reflectxPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != reflectxTopLevel || entry.Fields != reflectxFields ||
		entry.Methods != reflectxMethods || entry.Fingerprint != reflectxFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			reflectxTopLevel, reflectxFields, reflectxMethods, reflectxFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != reflectxPackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != reflectxTopLevel ||
		review.ReviewedSurface.Fields != reflectxFields ||
		review.ReviewedSurface.Methods != reflectxMethods ||
		review.ReviewedSurface.EntryCount != reflectxEntries ||
		review.ReviewedSurface.Fingerprint != reflectxFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !sameSet(review.ContractIDs, reflectxContractIDs()) {
		failures = append(failures, fmt.Sprintf("contract package review contract ids mismatch: got %v want %v", review.ContractIDs, reflectxContractIDs()))
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != reflectxEntries {
		failures = append(failures, fmt.Sprintf("reflectx audit entry count mismatch: got %d want %d", len(entries), reflectxEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != reflectxPackage {
			failures = append(failures, "non-reflectx audit entry passed into reflectx verifier: "+entry.ID)
		}
		if strings.Contains(entry.Package, "/internal/") {
			failures = append(failures, "internal package included in reflectx verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate reflectx audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "reflectx audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != reflectxTopLevel || counts["field"] != reflectxFields || counts["method"] != reflectxMethods {
		failures = append(failures, fmt.Sprintf("reflectx audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	failures = append(failures, verifyExpectedFields(entries)...)

	return failures
}

func verifyExpectedFields(entries []auditEntry) []string {
	expected := []string{
		"TagConfig.Name",
		"TagConfig.Value",
		"VisitorConfig.Recursive",
		"VisitorConfig.DiveTag",
		"VisitorConfig.MaxDepth",
		"Visitor.VisitStruct",
		"Visitor.VisitField",
		"Visitor.VisitMethod",
		"TypeVisitor.VisitStructType",
		"TypeVisitor.VisitFieldType",
		"TypeVisitor.VisitMethodType",
	}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Kind == "field" {
			seen[entry.Symbol] = true
		}
	}

	var failures []string
	for _, symbol := range expected {
		if !seen[symbol] {
			failures = append(failures, "reflectx audit entries missing field "+symbol)
		}
	}

	return failures
}

func verifyCoverage(
	entries []auditEntry,
	manifestEntry manifestEntry,
	review contractPackageReview,
	contracts map[string]contractEntry,
) []string {
	var failures []string
	expected := []string{englishSmallUtilitiesPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest reflectx coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review reflectx coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	for _, id := range reflectxContractIDs() {
		contract, ok := contracts[id]
		if !ok {
			failures = append(failures, "missing reflectx contract "+id)
			continue
		}
		if !sameSet(contract.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("contract entry %s coverage mismatch: got %v want %v", id, contract.Coverage, expected))
		}
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, reflectxPackage)
	if section == "" {
		return []string{index.label + " missing reflectx package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s reflectx index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}
	if strings.Contains(section, "/internal/") {
		failures = append(failures, index.label+" reflectx index includes internal package")
	}

	return failures
}

func verifyDocumentedSurface(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		ref := entry.Symbol
		if entry.Kind == "top" {
			ref = "reflectx." + entry.Symbol
		}
		if !hasCodeReference(doc.content, ref) {
			failures = append(failures, fmt.Sprintf("%s missing audited reflectx entry `%s`", doc.label, ref))
		}
	}

	return failures
}

func verifyNoPhantomReflectxRefs(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		if entry.Kind == "top" {
			valid["reflectx."+entry.Symbol] = true
		}
	}

	var failures []string
	for _, ref := range reflectxReferencesFromMarkdownCode(doc.content) {
		if valid[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown reflectx public API: %s", doc.label, ref))
	}

	return failures
}

func verifyNoPhantomReflectxFields(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		if entry.Kind == "field" {
			valid[entry.Symbol] = true
		}
	}

	var failures []string
	for _, ref := range reflectxFieldReferencesFromMarkdownCode(doc.content) {
		if valid[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown reflectx exported field: %s", doc.label, ref))
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
		reflectxFingerprint,
		"78 exported top-level",
		"11 exported fields",
		"no exported\nmethods",
		"reflectx.ToString",
		"reflectx.ToBoolE",
		"reflectx.ToDecimalE",
		"reflectx.ErrCannotConvertType",
		"reflectx.SetStringValue",
		"reflectx.Equal",
		"reflectx.Contains",
		"reflectx.VisitFor",
		"TagConfig.Name",
		"VisitorConfig.Recursive",
		"Visitor.VisitMethod",
		"TypeVisitor.VisitMethodType",
		"github.com/spf13/cast",
		"To*E",
		"non-E variants",
		"decimal.Zero",
		"nil pointer/interface",
		"decimal.NewFromAny",
		"reflectx.IsSimilarType",
		"reflectx.IsTypeCompatible",
		"addressable pointer copy",
		"pointer method set",
		"fresh",
		"nil string pointers/slices/maps",
		"cross-category",
		"convertible map keys",
		"reflectx.Continue = 0",
		"reflectx.Stop = 1",
		"reflectx.SkipChildren = 2",
		"TagConfig{Name: \"visit\", Value: \"dive\"}",
		"VisitorConfig.MaxDepth == 0",
		"depth >= MaxDepth",
		"StructField.Index",
		"absolute index path",
	}
}

func verifyContractLedger(review contractPackageReview, contracts map[string]contractEntry, sourceRoot string) []string {
	var failures []string
	expectedCoverage := []string{englishSmallUtilitiesPath}
	allEvidence := append([]string{}, review.SourceEvidence...)

	for _, id := range reflectxContractIDs() {
		contract, ok := contracts[id]
		if !ok {
			failures = append(failures, "missing reflectx contract "+id)
			continue
		}
		if contract.ID != id {
			failures = append(failures, fmt.Sprintf("reflectx contract id mismatch: got %q want %q", contract.ID, id))
		}
		if contract.Package != reflectxPackage {
			failures = append(failures, fmt.Sprintf("reflectx contract package mismatch for %s: got %q", id, contract.Package))
		}
		if contract.Kind != "runtime-contract" {
			failures = append(failures, fmt.Sprintf("reflectx contract kind mismatch for %s: got %q", id, contract.Kind))
		}
		if contract.Disposition != "documented:semantic-contract" {
			failures = append(failures, fmt.Sprintf("reflectx contract disposition mismatch for %s: got %q", id, contract.Disposition))
		}
		if !sameSet(contract.Coverage, expectedCoverage) {
			failures = append(failures, fmt.Sprintf("reflectx contract coverage mismatch for %s: got %v want %v", id, contract.Coverage, expectedCoverage))
		}
		for _, term := range expectedContractTerms(id) {
			if !contains(contract.Terms, term) {
				failures = append(failures, fmt.Sprintf("reflectx contract %s missing term %s", id, term))
			}
		}
		allEvidence = append(allEvidence, contract.SourceEvidence...)
		allEvidence = append(allEvidence, contract.TestEvidence...)
	}

	for _, item := range allEvidence {
		path, lineText, ok := strings.Cut(item, ":")
		if !ok || lineText == "" {
			failures = append(failures, "reflectx contract evidence missing line number: "+item)
			continue
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, "reflectx contract evidence missing file: "+item)
		}
	}

	return failures
}

func expectedContractTerms(id string) []string {
	switch id {
	case reflectxPackage + "#runtime-contract:cast-alias-zero-vs-error":
		return []string{"github.com/spf13/cast", "To*E", "non-E variants", "zero value"}
	case reflectxPackage + "#runtime-contract:decimal-conversion-nil-and-error":
		return []string{"reflectx.ToDecimalE", "decimal.Zero", "nil pointer/interface", "decimal.NewFromAny", "reflectx.ToDecimal"}
	case reflectxPackage + "#runtime-contract:type-compatibility-and-conversion":
		return []string{"reflectx.Indirect", "reflectx.IsPointerToStruct", "reflectx.IsSimilarType", "reflectx.IsTypeCompatible", "reflectx.ConvertValue", "reflectx.ErrCannotConvertType"}
	case reflectxPackage + "#runtime-contract:method-discovery":
		return []string{"reflectx.FindMethod", "addressable pointer copy", "reflectx.CollectMethods", "pointer method set"}
	case reflectxPackage + "#runtime-contract:string-field-accessors":
		return []string{"reflectx.IsStringType", "reflectx.SetStringValue", "fresh", "nil string pointers/slices/maps", "no-op"}
	case reflectxPackage + "#runtime-contract:empty-equal-contains":
		return []string{"reflectx.IsEmpty", "*string", "reflectx.Equal", "cross-category", "reflectx.Contains", "convertible map keys"}
	case reflectxPackage + "#runtime-contract:struct-visitor-dive-depth-contract":
		return []string{"reflectx.Continue = 0", "reflectx.Stop = 1", "reflectx.SkipChildren = 2", "TagConfig{Name: \"visit\", Value: \"dive\"}", "VisitorConfig.MaxDepth == 0", "depth >= MaxDepth", "StructField.Index", "absolute index path", "pointer method set"}
	default:
		panic("unknown reflectx contract id " + id)
	}
}

func verifySourceContracts(sourceRoot string) []string {
	files := map[string]string{
		"reflectx/convert.go":           readSourceFile(sourceRoot, "reflectx/convert.go").content,
		"reflectx/reflectx.go":          readSourceFile(sourceRoot, "reflectx/reflectx.go").content,
		"reflectx/string_field.go":      readSourceFile(sourceRoot, "reflectx/string_field.go").content,
		"reflectx/value.go":             readSourceFile(sourceRoot, "reflectx/value.go").content,
		"reflectx/visitor.go":           readSourceFile(sourceRoot, "reflectx/visitor.go").content,
		"reflectx/convert_test.go":      readSourceFile(sourceRoot, "reflectx/convert_test.go").content,
		"reflectx/reflectx_test.go":     readSourceFile(sourceRoot, "reflectx/reflectx_test.go").content,
		"reflectx/string_field_test.go": readSourceFile(sourceRoot, "reflectx/string_field_test.go").content,
		"reflectx/value_test.go":        readSourceFile(sourceRoot, "reflectx/value_test.go").content,
		"reflectx/visitor_test.go":      readSourceFile(sourceRoot, "reflectx/visitor_test.go").content,
	}

	checks := []struct {
		file string
		term string
	}{
		{"reflectx/convert.go", `ToString  = cast.ToString`},
		{"reflectx/convert.go", `ToBoolE = cast.ToBoolE`},
		{"reflectx/convert.go", `return decimal.Zero, nil`},
		{"reflectx/convert.go", `return decimal.NewFromAny(value)`},
		{"reflectx/convert.go", `v, _ := ToDecimalE(value)`},
		{"reflectx/reflectx.go", `func Indirect(t reflect.Type) reflect.Type`},
		{"reflectx/reflectx.go", `return t != nil && t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct`},
		{"reflectx/reflectx.go", `strings.IndexByte(name1, '[')`},
		{"reflectx/reflectx.go", `addressablePointer(target)`},
		{"reflectx/reflectx.go", `target = target.Elem()`},
		{"reflectx/reflectx.go", `sourceType == targetType || sourceType.AssignableTo(targetType)`},
		{"reflectx/reflectx.go", `targetType.Kind() == reflect.Interface`},
		{"reflectx/reflectx.go", `return reflect.Zero(targetType), nil`},
		{"reflectx/reflectx.go", `reflect.New(targetType.Elem())`},
		{"reflectx/reflectx.go", `fmt.Errorf("%w: %s -> %s", ErrCannotConvertType, sourceType, targetType)`},
		{"reflectx/string_field.go", `return t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.String`},
		{"reflectx/string_field.go", `return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.String`},
		{"reflectx/string_field.go", `t.Key().Kind() == reflect.String`},
		{"reflectx/string_field.go", `return "", false`},
		{"reflectx/string_field.go", `strValue := s`},
		{"reflectx/string_field.go", `v.Set(reflect.ValueOf(&strValue))`},
		{"reflectx/value.go", `return isEmpty(rv)`},
		{"reflectx/value.go", `elem.Kind() == reflect.String`},
		{"reflectx/value.go", `isSignedInt(ka) && isSignedInt(kb)`},
		{"reflectx/value.go", `isUnsignedInt(ka) && isUnsignedInt(kb)`},
		{"reflectx/value.go", `isFloat(ka) && isFloat(kb)`},
		{"reflectx/value.go", `return va.Interface() == vb.Interface()`},
		{"reflectx/value.go", `strings.Contains(rv.String(), ev.String())`},
		{"reflectx/value.go", `if !ev.Type().ConvertibleTo(mapKeyType)`},
		{"reflectx/visitor.go", `Continue VisitAction = iota`},
		{"reflectx/visitor.go", `SkipChildren`},
		{"reflectx/visitor.go", `Recursive: true`},
		{"reflectx/visitor.go", `DiveTag:   TagConfig{Name: "visit", Value: "dive"}`},
		{"reflectx/visitor.go", `if config.MaxDepth > 0 && depth >= config.MaxDepth`},
		{"reflectx/visitor.go", `if ancestors.Contains(targetType)`},
		{"reflectx/visitor.go", `if !field.CanInterface()`},
		{"reflectx/visitor.go", `case SkipChildren:`},
		{"reflectx/visitor.go", `if field.Anonymous`},
		{"reflectx/visitor.go", `field.Tag.Get(diveTag.Name) == diveTag.Value`},
		{"reflectx/visitor.go", `field.Index = buildAbsoluteIndexPath(parentIndexPath, field)`},
		{"reflectx/visitor.go", `ptrTarget := addressablePointer(target)`},
		{"reflectx/visitor.go", `ptrType := reflect.PointerTo(targetType)`},
		{"reflectx/convert_test.go", `func TestToDecimalE`},
		{"reflectx/convert_test.go", `func TestToDecimal`},
		{"reflectx/reflectx_test.go", `func TestIsSimilarType`},
		{"reflectx/reflectx_test.go", `func TestFindMethod`},
		{"reflectx/reflectx_test.go", `func TestCollectMethods`},
		{"reflectx/reflectx_test.go", `func TestConvertValue`},
		{"reflectx/string_field_test.go", `func TestSetStringValue`},
		{"reflectx/value_test.go", `func TestEqual`},
		{"reflectx/value_test.go", `func TestContains`},
		{"reflectx/visitor_test.go", `func TestVisitFieldIndexPathAnonymousEmbedded`},
		{"reflectx/visitor_test.go", `func TestVisitCyclicReference`},
		{"reflectx/visitor_test.go", `func TestVisitTypeMethodStopAction`},
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
	return runCommand(sourceRoot, "go", "test", "./reflectx")
}

func loadReflectxAuditEntries(path string) []auditEntry {
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
		if entry.Package == reflectxPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})
	if len(entries) == 0 {
		panic("API audit ledger missing reflectx entries")
	}

	return entries
}

func loadReflectxManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == reflectxPackage {
			return entry
		}
	}

	panic("API audit manifest missing reflectx package")
}

func loadReflectxContracts(path string) (contractPackageReview, map[string]contractEntry) {
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
		if item.Package == reflectxPackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing reflectx package review")
	}

	contracts := make(map[string]contractEntry)
	for _, item := range ledger.Entries {
		if item.Package == reflectxPackage {
			contracts[item.ID] = item
		}
	}
	if len(contracts) == 0 {
		panic("API contract ledger missing reflectx contract entries")
	}

	return review, contracts
}

func loadLiveReflectxEntry(sourceRoot, docsRoot string) manifestEntry {
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
		if entry.Package == reflectxPackage {
			return entry
		}
	}

	panic("live API inventory missing reflectx package")
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

func reflectxReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`reflectx\.[A-Z][A-Za-z0-9_]*`)
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

func reflectxFieldReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`(?:TagConfig|VisitorConfig|Visitor|TypeVisitor)\.[A-Z][A-Za-z0-9_]*`)
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

func reflectxContractIDs() []string {
	return []string{
		reflectxPackage + "#runtime-contract:cast-alias-zero-vs-error",
		reflectxPackage + "#runtime-contract:decimal-conversion-nil-and-error",
		reflectxPackage + "#runtime-contract:type-compatibility-and-conversion",
		reflectxPackage + "#runtime-contract:method-discovery",
		reflectxPackage + "#runtime-contract:string-field-accessors",
		reflectxPackage + "#runtime-contract:empty-equal-contains",
		reflectxPackage + "#runtime-contract:struct-visitor-dive-depth-contract",
	}
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
