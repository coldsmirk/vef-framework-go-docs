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
	moldPackage = "github.com/coldsmirk/vef-framework-go/mold"

	moldFingerprint = "42c01cc63b24bafe689f8da49350ea05675c9089cf88d6f02f5ac9ea478d0859"
	moldTopLevel    = 17
	moldFields      = 1
	moldMethods     = 23
	moldEntries     = 41

	moldGroupedEntries              = 24
	moldGroupedFields               = 1
	moldGroupedMethods              = 23
	moldGroupedReceivers            = 12
	moldGroupedSignatureFingerprint = "c13ba6272799d3e2c0caecc57285bc41e21286a5d29975421c0910c29950fe5d"
	moldGroupedReceiverFingerprint  = "894ac12de35dbfce1502fb244901535292680957f9b5650e02c2b351af0a100f"

	englishMoldPath  = "docs/features/mold.md"
	chineseMoldPath  = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/mold.md"
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

	englishDocs := readCorpus("English mold docs", filepath.Join(docsRoot, englishMoldPath))
	chineseDocs := readCorpus("Chinese mold docs", filepath.Join(docsRoot, chineseMoldPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[moldPackage]
	liveAuditEntries := moldEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	moldEntries := moldEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(moldEntries)...)
	failures = append(failures, verifyLiveAuditEntries(moldEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedMoldSurface(moldEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, moldEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, moldEntries)...)
	failures = append(failures, verifyMoldDocs(moldEntries, englishDocs, chineseDocs)...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Mold contract docs verified: 41 public entries, 24 grouped transformer entries, tag grammar and default transformer boundary")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != moldPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != moldTopLevel ||
		entry.Fields != moldFields ||
		entry.Methods != moldMethods ||
		entry.Fingerprint != moldFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s mold surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			moldTopLevel, moldFields, moldMethods, moldFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != moldPackage {
			continue
		}

		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, moldCoverage()) {
			failures = append(failures, fmt.Sprintf("mold manifest coverage mismatch: got %v want %v", entry.Coverage, moldCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing mold package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	expectedContracts := map[string][]string{
		moldPackage + "#event-contract:dictionary-cache-invalidation": {
			"DictionaryChangedEvent",
			"PublishDictionaryChangedEvent",
			"Keys",
			"Resolve",
		},
		moldPackage + "#field-contract:transformer-inputs-and-traversal": {
			"Transformer.Struct",
			"Transformer.Field",
			"time.Time",
			"nil",
		},
		moldPackage + "#string-contract:mold-tag-grammar": {
			"dive",
			"keys",
			"endkeys",
			"0x2C",
		},
		moldPackage + "#string-contract:translate-transformer": {
			"translate=dict:",
			"<Field>Name",
			"user?",
			"dict:status?",
		},
	}

	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != moldPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "mold contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != moldTopLevel ||
			review.ReviewedSurface.Fields != moldFields ||
			review.ReviewedSurface.Methods != moldMethods ||
			review.ReviewedSurface.EntryCount != moldEntries ||
			review.ReviewedSurface.Fingerprint != moldFingerprint {
			failures = append(failures, fmt.Sprintf(
				"mold contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, moldCoverage()) {
			failures = append(failures, fmt.Sprintf("mold contract review coverage mismatch: got %v want %v", review.Coverage, moldCoverage()))
		}
		if !sameSet(review.ContractIDs, sortedKeys(expectedContracts)) {
			failures = append(failures, fmt.Sprintf("mold contract ids mismatch: got %v want %v", review.ContractIDs, sortedKeys(expectedContracts)))
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing mold package review")
	}

	foundContracts := map[string]bool{}
	for _, entry := range contracts.Entries {
		terms, ok := expectedContracts[entry.ID]
		if !ok {
			continue
		}
		foundContracts[entry.ID] = true
		if entry.Package != moldPackage {
			failures = append(failures, fmt.Sprintf("mold contract entry package mismatch for %s: %s", entry.ID, entry.Package))
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "mold contract entry disposition mismatch for "+entry.ID+": "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, moldCoverage()) {
			failures = append(failures, fmt.Sprintf("mold contract coverage mismatch for %s: got %v want %v", entry.ID, entry.Coverage, moldCoverage()))
		}
		for _, term := range terms {
			if !contains(entry.Terms, term) {
				failures = append(failures, fmt.Sprintf("mold contract %s missing term %s", entry.ID, term))
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.TestEvidence)...)
	}
	for id := range expectedContracts {
		if !foundContracts[id] {
			failures = append(failures, "contract ledger missing mold contract entry "+id)
		}
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != moldEntries {
		failures = append(failures, fmt.Sprintf("mold audit entry count mismatch: got %d want %d", len(entries), moldEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != moldPackage {
			failures = append(failures, "non-mold audit entry passed into mold verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate mold audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "mold audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, moldCoverage()) {
			failures = append(failures, fmt.Sprintf("mold audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, moldCoverage()))
		}
	}
	if counts["top"] != moldTopLevel || counts["field"] != moldFields || counts["method"] != moldMethods {
		failures = append(failures, fmt.Sprintf("mold audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != moldTopLevel ||
		dispositionCounts["grouped:type-member-family"] != moldGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"mold audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			moldTopLevel,
			moldGroupedEntries,
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
			failures = append(failures, fmt.Sprintf("mold missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"mold live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
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
			failures = append(failures, fmt.Sprintf("mold extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedMoldSurface(entries []auditEntry, docs []corpus) []string {
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
			failures = append(failures, fmt.Sprintf("mold grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("mold grouped type-member surface", rows, moldGroupedEntries, moldGroupedSignatureFingerprint)...)
	if kindCounts["field"] != moldGroupedFields || kindCounts["method"] != moldGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"mold grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			moldGroupedFields,
			moldGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("mold grouped receiver/type families", receiverRows, moldGroupedReceivers, moldGroupedReceiverFingerprint)...)

	for _, doc := range docs {
		for _, term := range []string{
			"41 public mold entries",
			"24 grouped mold field/method entries",
			"12 mold receiver/type families",
			"1 exported mold field entry",
			"23 exported mold method entries",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped mold audit term "+term)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, moldPackage)
	if section == "" {
		return []string{index.label + " missing mold package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s mold index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyMoldDocs(entries []auditEntry, englishDocs, chineseDocs corpus) []string {
	var topSymbols []string
	for _, entry := range entries {
		if entry.Kind == "top" {
			topSymbols = append(topSymbols, entry.Symbol)
		}
	}
	sort.Strings(topSymbols)

	commonTerms := []string{
		"`mold`",
		"`Transformer.Struct`",
		"`Transformer.Field`",
		"`time.Time`",
		"`mold:\"-\"`",
		"`dive`",
		"`keys`",
		"`endkeys`",
		"`0x2C`",
		"`translate`",
		"translate=dict:",
		"`<Field>Name`",
		"`user?`",
		"`dict:status?`",
		"`DictionaryChangedEvent`",
		"`PublishDictionaryChangedEvent`",
		"Keys",
		"`Resolve`",
		"`DictionaryLoader`",
		"`event.Bus`",
		"`NewCachedDictionaryResolver`",
		"`DictionaryLoaderFunc`",
		"`DictionaryTranslator`",
		"`DictionaryResolver`",
		"mold:\"expr=price * qty\"",
		"`expression.Engine`",
		"`vef:mold:field_transformers`",
		"`expr`",
	}
	englishTerms := []string{
		"declaration order",
		"not provided by the `mold` module",
		"other field transformers must be registered",
	}
	chineseTerms := []string{
		"声明顺序",
		"不是 `mold` module 单独提供",
		"其他 field transformer",
	}

	var failures []string
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		for _, symbol := range topSymbols {
			if !strings.Contains(doc.content, "`"+symbol+"`") &&
				!strings.Contains(doc.content, "`mold."+symbol+"`") {
				failures = append(failures, doc.label+" missing top-level mold symbol `"+symbol+"`")
			}
		}
		for _, term := range commonTerms {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing mold semantic term "+term)
			}
		}
	}
	for _, term := range englishTerms {
		if !containsNormalized(englishDocs.content, term) {
			failures = append(failures, englishDocs.label+" missing mold semantic term "+term)
		}
	}
	for _, term := range chineseTerms {
		if !containsNormalized(chineseDocs.content, term) {
			failures = append(failures, chineseDocs.label+" missing mold semantic term "+term)
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
			path: "mold/transformer.go",
			terms: []string{
				"type Transformer interface",
				"Struct(ctx context.Context, value any) error",
				"Field(ctx context.Context, value any, tags string) error",
				"type FieldTransformer interface",
				"Tag() string",
				"Transform(ctx context.Context, fl FieldLevel) error",
				"type StructTransformer interface",
				"type Interceptor interface",
				"type FieldLevel interface",
				"SiblingField(name string) (reflect.Value, bool)",
				"Struct() reflect.Value",
				"type StructLevel interface",
				"type Func func(ctx context.Context, fl FieldLevel) error",
				"type StructLevelFunc func(ctx context.Context, sl StructLevel) error",
				"type InterceptorFunc func(current reflect.Value) (inner reflect.Value)",
			},
		},
		{
			path: "mold/translator.go",
			terms: []string{
				"type Translator interface",
				"Supports(kind string) bool",
				"Translate(ctx context.Context, kind, value string) (string, error)",
				"type DictionaryResolver interface",
				"Resolve(ctx context.Context, key, code string) (string, error)",
				"type DictionaryLoader interface",
				"Load(ctx context.Context, key string) (map[string]string, error)",
			},
		},
		{
			path: "mold/cached_dictionary_resolver.go",
			terms: []string{
				"eventTypeDictionaryChanged = \"vef.translate.dictionary.changed\"",
				"type DictionaryLoaderFunc func(ctx context.Context, key string) (map[string]string, error)",
				"func (f DictionaryLoaderFunc) Load(ctx context.Context, key string) (map[string]string, error)",
				"type DictionaryChangedEvent struct",
				"Keys []string `json:\"keys\"`",
				"func (*DictionaryChangedEvent) EventType() string { return eventTypeDictionaryChanged }",
				"func PublishDictionaryChangedEvent(ctx context.Context, bus event.Bus, keys ...string) error",
				"func NewCachedDictionaryResolver(",
				"panic(\"NewCachedDictionaryResolver requires a non-nil DictionaryLoader, but got nil\")",
				"panic(\"NewCachedDictionaryResolver requires a non-nil event.Bus, but got nil\")",
				"event.SubscribeTyped[*DictionaryChangedEvent](bus, resolver.handleInvalidation)",
				"func (r *CachedDictionaryResolver) Resolve(ctx context.Context, key, code string) (string, error)",
				"if key == \"\" || code == \"\"",
				"return r.cache.Invalidate(ctx, evt.Keys...)",
			},
		},
		{
			path: "internal/mold/module.go",
			terms: []string{
				"fx.Decorate(",
				"func(loader mold.DictionaryLoader, bus event.Bus) mold.DictionaryResolver",
				"fx.ParamTags(`optional:\"true\"`)",
				"NewTransformer",
				"`group:\"vef:mold:field_transformers\"`",
				"`group:\"vef:mold:struct_transformers\"`",
				"`group:\"vef:mold:interceptors\"`",
				"NewTranslateTransformer",
				"`group:\"vef:mold:translators\"`",
				"NewDictionaryTranslator",
			},
		},
		{
			path: "internal/mold/mold.go",
			terms: []string{
				"tagName:         \"mold\"",
				"func New() *MoldTransformer",
				"func (t *MoldTransformer) Register(tag string, fn mold.Func)",
				"panic(\"mold: transformation tag cannot be empty\")",
				"panic(\"mold: transformation function cannot be nil\")",
				"func (t *MoldTransformer) RegisterAlias(alias, tags string)",
				"func (t *MoldTransformer) RegisterStructLevel(fn mold.StructLevelFunc, types ...any)",
				"func (t *MoldTransformer) RegisterInterceptor(fn mold.InterceptorFunc, types ...any)",
				"func (t *MoldTransformer) Struct(ctx context.Context, v any) error",
				"orig.Kind() != reflect.Pointer || orig.IsNil()",
				"val.Kind() != reflect.Struct || typ == timeType",
				"func (t *MoldTransformer) Field(ctx context.Context, v any, tags string) error",
				"if tags == \"\" || tags == ignoreTag",
				"return t.handleDive(ctx, current, kind, ct)",
				"return t.traverseStruct(ctx, current, original)",
			},
		},
		{
			path: "internal/mold/restricted.go",
			terms: []string{
				"diveTag            = \"dive\"",
				"tagSeparator       = \",\"",
				"ignoreTag          = \"-\"",
				"tagKeySeparator    = \"=\"",
				"utf8HexComma       = \"0x2C\"",
				"keysTag            = \"keys\"",
				"endKeysTag         = \"endkeys\"",
			},
		},
		{
			path: "internal/mold/cache.go",
			terms: []string{
				"strings.Split(tagString, tagSeparator)",
				"case diveTag:",
				"case keysTag:",
				"case endKeysTag:",
				"strings.SplitN(tag, tagKeySeparator, 2)",
				"strings.ReplaceAll(vals[1], utf8HexComma, \",\")",
			},
		},
		{
			path: "internal/mold/translate.go",
			terms: []string{
				"translatedFieldNameSuffix = \"Name\"",
				"func (*TranslateTransformer) Tag() string",
				"return \"translate\"",
				"func (t *TranslateTransformer) Transform(ctx context.Context, fl mold.FieldLevel) error",
				"translatedFieldName := name + translatedFieldNameSuffix",
				"if strings.HasSuffix(kind, \"?\")",
				"func (t *TranslateTransformer) transformStringSlice",
				"NewTranslateTransformer(translators []mold.Translator) mold.FieldTransformer",
			},
		},
		{
			path: "internal/mold/dictionary_translator.go",
			terms: []string{
				"dictKeyPrefix = \"dict:\"",
				"type DictionaryTranslator struct",
				"func (*DictionaryTranslator) Supports(kind string) bool",
				"return strings.HasPrefix(kind, dictKeyPrefix)",
				"func (t *DictionaryTranslator) Translate(ctx context.Context, kind, value string) (string, error)",
				"dictKey := kind[len(dictKeyPrefix):]",
				"func NewDictionaryTranslator(resolver mold.DictionaryResolver) mold.Translator",
			},
		},
		{
			path: "internal/expression/transformer.go",
			terms: []string{
				"fieldTransformerTag = \"expr\"",
				"func NewFieldTransformer(engine expression.Engine) mold.FieldTransformer",
				"func (*fieldTransformer) Tag() string",
				"return fieldTransformerTag",
				"source := fl.Param()",
				"env = s.Interface()",
				"t.engine.Evaluate(ctx, source, env)",
				"value.Decode(field.Addr().Interface())",
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
	cmd := exec.Command("go", "test", "./mold", "./internal/mold")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./mold ./internal/mold failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func moldEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == moldPackage {
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

func moldCoverage() []string {
	return []string{englishMoldPath}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
