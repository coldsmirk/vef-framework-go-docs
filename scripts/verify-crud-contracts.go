package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	crudPackage                     = "github.com/coldsmirk/vef-framework-go/crud"
	crudFingerprint                 = "2f13a6adc30cc0203584e004a78b2090a461e74bf21291533cd955a85ea3105d"
	crudTopLevel                    = 129
	crudFields                      = 36
	crudMethods                     = 263
	crudEntries                     = 428
	crudGroupedEntries              = 299
	crudGroupedFields               = 36
	crudGroupedMethods              = 263
	crudGroupedReceivers            = 27
	crudGroupedSignatureFingerprint = "ad294250b01cbbd2f67019f48524ca9771c867a2b68d99b1017b977326a2006d"
	crudGroupedReceiverFingerprint  = "f20b3dc2a881bf0b206babbe89a687cde293e74a641c95dece429dd8ecd692e7"

	englishCrudPath         = "docs/data-access/crud.md"
	chineseCrudPath         = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/crud.md"
	englishHooksPath        = "docs/data-access/hooks.md"
	chineseHooksPath        = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/hooks.md"
	englishTransactionsPath = "docs/data-access/transactions.md"
	chineseTransactionsPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/transactions.md"
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
	ID       string   `json:"id"`
	Package  string   `json:"package"`
	Kind     string   `json:"kind"`
	Coverage []string `json:"coverage"`
	Terms    []string `json:"terms"`
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

	englishCrud := readCorpus("English CRUD docs", filepath.Join(docsRoot, englishCrudPath))
	chineseCrud := readCorpus("Chinese CRUD docs", filepath.Join(docsRoot, chineseCrudPath))
	englishHooks := readCorpus("English hooks docs", filepath.Join(docsRoot, englishHooksPath))
	chineseHooks := readCorpus("Chinese hooks docs", filepath.Join(docsRoot, chineseHooksPath))
	englishTransactions := readCorpus("English transactions docs", filepath.Join(docsRoot, englishTransactionsPath))
	chineseTransactions := readCorpus("Chinese transactions docs", filepath.Join(docsRoot, chineseTransactionsPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	englishTopicDocs := combineCorpora("English CRUD topic docs", englishCrud, englishHooks, englishTransactions)
	chineseTopicDocs := combineCorpora("Chinese CRUD topic docs", chineseCrud, chineseHooks, chineseTransactions)

	auditEntries := loadCrudAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadCrudManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadCrudContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveInventory(sourceRoot, docsRoot)[crudPackage]

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(auditEntries)...)
	failures = append(failures, verifyGroupedCRUDSurface(auditEntries)...)
	failures = append(failures, verifyCoverage(auditEntries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, auditEntries)...)
	}
	for _, doc := range []corpus{englishTopicDocs, chineseTopicDocs} {
		failures = append(failures, verifyDocumentedTopLevel(doc, auditEntries)...)
		failures = append(failures, verifySemanticDocTerms(doc, auditEntries)...)
	}
	for _, doc := range []corpus{englishCrud, chineseCrud} {
		failures = append(failures, verifyActionRows(doc, auditEntries)...)
		failures = append(failures, verifyProcessorSignatures(doc, auditEntries)...)
		failures = append(failures, verifyErrorCodeRows(doc, auditEntries)...)
	}

	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, runRuntimeChecks(sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("CRUD contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("CRUD contract docs verified: %d public entries, 299 grouped builder entries, source-derived action/error/DTO/hook contracts, 6 doc corpora\n", len(auditEntries))
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != crudPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != crudTopLevel || entry.Fields != crudFields ||
		entry.Methods != crudMethods || entry.Fingerprint != crudFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			crudTopLevel, crudFields, crudMethods, crudFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != crudPackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != crudTopLevel ||
		review.ReviewedSurface.Fields != crudFields ||
		review.ReviewedSurface.Methods != crudMethods ||
		review.ReviewedSurface.EntryCount != crudEntries ||
		review.ReviewedSurface.Fingerprint != crudFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, crudContractID()) {
		failures = append(failures, "contract package review missing CRUD contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != crudEntries {
		failures = append(failures, fmt.Sprintf("CRUD audit entry count mismatch: got %d want %d", len(entries), crudEntries))
	}
	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != crudPackage {
			failures = append(failures, "non-CRUD audit entry passed into CRUD verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate CRUD audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "CRUD audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != crudTopLevel || counts["field"] != crudFields || counts["method"] != crudMethods {
		failures = append(failures, fmt.Sprintf("CRUD audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}

	return failures
}

func verifyGroupedCRUDSurface(entries []auditEntry) []string {
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
			failures = append(failures, fmt.Sprintf("CRUD grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}

		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint(
		"CRUD grouped builder surface",
		rows,
		crudGroupedEntries,
		crudGroupedSignatureFingerprint,
	)...)
	if kindCounts["field"] != crudGroupedFields || kindCounts["method"] != crudGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"CRUD grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			crudGroupedFields,
			crudGroupedMethods,
		))
	}

	receiverRows := make([]string, 0, len(receiverCounts))
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint(
		"CRUD grouped receiver families",
		receiverRows,
		crudGroupedReceivers,
		crudGroupedReceiverFingerprint,
	)...)

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

func receiverForSymbol(symbol string) (string, bool) {
	receiver, _, ok := strings.Cut(symbol, ".")
	if !ok || receiver == "" {
		return "", false
	}

	return receiver, true
}

func verifyCoverage(
	entries []auditEntry,
	manifestEntry manifestEntry,
	review contractPackageReview,
	contract contractEntry,
) []string {
	var failures []string
	expected := []string{englishCrudPath, "docs/data-access/query-builder.md", englishHooksPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest CRUD coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review CRUD coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry CRUD coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, crudPackage)
	if section == "" {
		return []string{index.label + " missing CRUD package section"}
	}
	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s CRUD index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
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
		if !documentMentionsSymbol(doc.content, entry.Symbol) {
			failures = append(failures, doc.label+" missing top-level public CRUD symbol `"+entry.Symbol+"`")
		}
	}

	return failures
}

func verifySemanticDocTerms(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		if entry.Kind != "field" {
			continue
		}
		field := strings.TrimPrefix(entry.Symbol, entryTypePrefix(entry.Symbol))
		if field == "" {
			continue
		}
		if !strings.Contains(doc.content, "`"+entry.Symbol+"`") &&
			!strings.Contains(doc.content, "`"+field+"`") &&
			!strings.Contains(doc.content, "`"+jsonFieldFromSignature(entry.Signature)+"`") {
			failures = append(failures, doc.label+" missing field contract for "+entry.Symbol)
		}
	}

	commonTerms := []string{
		"`CreateManyParams[TParams]`",
		"`UpdateManyParams[TParams]`",
		"`DeleteManyParams`",
		"`DataOption`",
		"`TreeDataOption`",
		"`DataOptionConfig`",
		"`DataOptionColumnMapping`",
		"`FindOperationConfig`",
		"`FindOperationOption`",
		"`QueryPartsConfig`",
		"`QueryPart`",
		"`QueryRoot`",
		"`QueryBase`",
		"`QueryRecursive`",
		"`QueryAll`",
		"`ApplyDataPermission`",
		"`GetAuditUserNameRelations`",
		"`WithDefaultPageSize(size)`",
		"`WithDefaultColumnMapping(mapping)`",
		"`WithIDColumn(name)`",
		"`WithParentIDColumn(name)`",
		"`WithDefaultFormat(...)`",
		"`WithExcelOptions(...)`",
		"`WithCsvOptions(...)`",
		"`WithFilenameBuilder(...)`",
		"`FormatExcel`",
		"`FormatCsv`",
		"`excel`",
		"`csv`",
		"`multipart`",
		"`params.file`",
		"`format`",
		"`errors`",
		"`RunInTx`",
		"`DisableDataPerm()`",
		"`WithAuditUserNames`",
		"`WithQueryApplier`",
		"`WithRelation`",
		"`WithSelect(column)`",
		"`WithSelectAs(column, alias)`",
	}
	failures = append(failures, missingTerms(doc, commonTerms)...)

	return failures
}

func verifyActionRows(doc corpus, entries []auditEntry) []string {
	var failures []string
	actionSymbols := actionConstantSymbols(entries)
	for _, symbol := range actionSymbols {
		entry := entryBySymbol(entries, symbol)
		value := constantValue(entry.Signature)
		if value == "" {
			failures = append(failures, "could not parse action constant value for "+symbol)
			continue
		}
		if !strings.Contains(doc.content, "`"+symbol+"`") {
			failures = append(failures, doc.label+" missing action constant `"+symbol+"`")
		}
		if !strings.Contains(doc.content, "`"+value+"`") {
			failures = append(failures, doc.label+" missing action constant value `"+value+"` for "+symbol)
		}
	}

	return failures
}

func verifyProcessorSignatures(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		if entry.Kind != "top" || !isProcessorSignatureSymbol(entry.Symbol) {
			continue
		}
		shortSignature := shortCrudSignature(entry.Signature)
		if !strings.Contains(doc.content, shortSignature) {
			failures = append(failures, doc.label+" missing processor signature "+shortSignature)
		}
	}

	return failures
}

func verifyErrorCodeRows(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		if entry.Kind != "top" || !strings.HasPrefix(entry.Symbol, "ErrCode") {
			continue
		}
		value := constantValue(entry.Signature)
		if value == "" {
			failures = append(failures, "could not parse CRUD error code value for "+entry.Symbol)
			continue
		}
		if !strings.Contains(doc.content, "`"+entry.Symbol+"`") {
			failures = append(failures, doc.label+" missing CRUD error code `"+entry.Symbol+"`")
		}
		if !strings.Contains(doc.content, "`"+value+"`") {
			failures = append(failures, doc.label+" missing CRUD error code value `"+value+"` for "+entry.Symbol)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "crud/constants.go",
			terms: []string{
				"FormatExcel TabularFormat = \"excel\"",
				"FormatCsv   TabularFormat = \"csv\"",
				"ErrCodeProcessorInvalidReturn = 2400",
				"RPCActionCreate          = \"create\"",
				"RESTActionUpdate          = \"put /:\" + IDColumn",
				"RESTActionFindTreeOptions = \"get /tree/options\"",
				"IDColumn          = orm.ColumnID",
				"ParentIDColumn    = \"parent_id\"",
				"LabelColumn       = \"label\"",
				"ValueColumn       = \"value\"",
				"DescriptionColumn = \"description\"",
			},
		},
		{
			path: "crud/builder.go",
			terms: []string{
				"type Builder[T any] interface",
				"Action(action string) T",
				"EnableAudit() T",
				"Timeout(timeout time.Duration) T",
				"Public() T",
				"RequiredPermission(token string) T",
				"RateLimit(maxRequests int, period time.Duration) T",
				"if err := api.ValidateActionName(action, b.kind); err != nil",
				"panic(err)",
				"RateLimit:          b.rateLimit",
			},
		},
		{
			path: "crud/crud.go",
			terms: []string{
				"func NewCreate[TModel, TParams any](kind ...api.Kind) Create[TModel, TParams]",
				"func NewFindTree[TModel, TSearch any](",
				"treeBuilder func(flatModels []TModel) []TModel",
				"idColumn:       IDColumn",
				"parentIDColumn: ParentIDColumn",
				"return api.Action(getAction(RPCActionImport, RESTActionImport, kind...))",
			},
		},
		{
			path: "crud/find.go",
			terms: []string{
				"Setup(db orm.DB, config *FindOperationConfig, opts ...*FindOperationOption) error",
				"ConfigureQuery(query orm.SelectQuery, search TSearch, meta api.Meta, ctx fiber.Ctx, part QueryPart) error",
				"WithProcessor(processor Processor[TProcessorIn, TSearch]) TOperation",
				"WithOptions(opts ...*FindOperationOption) TOperation",
				"WithSelect(column string, parts ...QueryPart) TOperation",
				"WithSelectAs(column, alias string, parts ...QueryPart) TOperation",
				"WithDefaultSort(sort ...*sortx.OrderSpec) TOperation",
				"DisableDataPerm() TOperation",
				"WithAuditUserNames(userModel any, nameColumn ...string) TOperation",
				"withDataPerm(condParts...)",
				"return ErrModelNoPrimaryKey",
				"return ErrAuditUserCompositePK",
			},
		},
		{
			path: "crud/find_option.go",
			terms: []string{
				"QueryRoot QueryPart = iota",
				"QueryBase",
				"QueryRecursive",
				"QueryAll",
				"func resolveQueryParts(parts ...QueryPart) []QueryPart",
				"return []QueryPart{QueryRoot}",
				"func withSearchApplier[TSearch any](parts ...QueryPart) *FindOperationOption",
				"errSearchTypeMismatch",
				"func withDataPerm(parts ...QueryPart) *FindOperationOption",
				"func withAuditUserNames(userModel any, nameColumn string, parts ...QueryPart) *FindOperationOption",
			},
		},
		{
			path: "crud/find_page.go",
			terms: []string{
				"WithDefaultPageSize(size int) FindPage[TModel, TSearch]",
				"pageable.Normalize(a.defaultPageSize)",
				"ErrCodeProcessorInvalidReturn",
				"fiber.StatusInternalServerError",
			},
		},
		{
			path: "crud/find_options.go",
			terms: []string{
				"WithDefaultColumnMapping(mapping *DataOptionColumnMapping) FindOptions[TModel, TSearch]",
				"mergeOptionColumnMapping(&config.DataOptionColumnMapping, a.defaultColumnMapping)",
				"query.Limit(maxOptionsLimit)",
				"options = []DataOption{}",
			},
		},
		{
			path: "crud/find_tree.go",
			terms: []string{
				"WithIDColumn(name string) FindTree[TModel, TSearch]",
				"WithParentIDColumn(name string) FindTree[TModel, TSearch]",
				"[]QueryPart{QueryBase, QueryRecursive}",
				"[]QueryPart{QueryBase}",
				"errColumnNotFound",
				"models := a.treeBuilder(flatModels)",
			},
		},
		{
			path: "crud/find_tree_options.go",
			terms: []string{
				"WithDefaultColumnMapping(mapping *DataOptionColumnMapping) FindTreeOptions[TModel, TSearch]",
				"WithIDColumn(name string) FindTreeOptions[TModel, TSearch]",
				"WithParentIDColumn(name string) FindTreeOptions[TModel, TSearch]",
				"[]QueryPart{QueryBase}",
				"treeOptions := tree.Build(flatOptions, treeAdapter)",
			},
		},
		{
			path: "crud/export.go",
			terms: []string{
				"WithDefaultFormat(format TabularFormat) Export[TModel, TSearch]",
				"WithExcelOptions(opts ...excel.ExportOption) Export[TModel, TSearch]",
				"WithCsvOptions(opts ...csv.ExportOption) Export[TModel, TSearch]",
				"WithFilenameBuilder(builder FilenameBuilder[TSearch]) Export[TModel, TSearch]",
				"lo.CoalesceOrEmpty(config.Format, a.defaultFormat, FormatExcel)",
				"return ErrUnsupportedExportFormat",
				"query.Limit(maxQueryLimit)",
				"fiber.HeaderContentDisposition",
			},
		},
		{
			path: "crud/import.go",
			terms: []string{
				"File *multipart.FileHeader `json:\"file\"`",
				"Format TabularFormat `json:\"format\"`",
				"if httpx.IsJSON(ctx)",
				"return ErrImportRequiresMultipart",
				"return ErrImportRequiresFile",
				"return ErrUnsupportedImportFormat",
				"return ErrFileOpenFailed",
				"return ErrImportTypeAssertionFailed",
				"ErrCodeImportValidationFailed",
				"i18n.T(\"crud_import_validation_failed\")",
				"fiber.StatusUnprocessableEntity",
				"return db.RunInTx(ctx.Context(), func(txCtx context.Context, tx orm.DB) error",
				"return result.Ok(",
				"\"total\": len(models)",
			},
		},
		{
			path: "crud/params.go",
			terms: []string{
				"List []TParams `json:\"list\" validate:\"required,min=1,dive\" label_i18n:\"crud_batch_create_list\"`",
				"PKs []any `json:\"pks\" validate:\"required,min=1\" label_i18n:\"crud_batch_delete_pks\"`",
				"Sort []sortx.OrderSpec `json:\"sort\"`",
			},
		},
		{
			path: "crud/option.go",
			terms: []string{
				"Label string `json:\"label\" bun:\"label\"`",
				"Value string `json:\"value\" bun:\"value\"`",
				"Description string `json:\"description,omitempty\" bun:\"description\"`",
				"Meta map[string]any `json:\"meta,omitempty\" bun:\"meta\"`",
				"LabelColumn string `json:\"labelColumn\"`",
				"ValueColumn string `json:\"valueColumn\"`",
				"DescriptionColumn string `json:\"descriptionColumn\"`",
				"MetaColumns []string `json:\"metaColumns\"`",
				"ID string `json:\"-\" bun:\"id\"`",
				"ParentID *string `json:\"-\" bun:\"parent_id\"`",
				"Children []TreeDataOption `json:\"children,omitempty\" bun:\"-\"`",
			},
		},
		{
			path: "crud/api_errors.go",
			terms: []string{
				"ErrCodeFieldNotExistInModel           = 2401",
				"ErrCodePrimaryKeyRequired             = 2402",
				"ErrCodeCompositePrimaryKeyRequiresMap = 2403",
				"ErrCodeUnsupportedExportFormat        = 2404",
				"ErrCodeImportRequiresMultipart        = 2405",
				"ErrCodeImportRequiresFile             = 2406",
				"ErrCodeUnsupportedImportFormat        = 2407",
				"ErrCodeFileOpenFailed                 = 2408",
				"ErrCodeImportTypeAssertionFailed      = 2409",
				"ErrCodeImportValidationFailed         = 2410",
			},
		},
	}

	var failures []string
	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	return failures
}

func runRuntimeChecks(sourceRoot string) []string {
	tmpDir, err := os.MkdirTemp("", "verify-crud-contracts-*")
	if err != nil {
		return []string{fmt.Sprintf("failed to create temp module: %v", err)}
	}
	defer os.RemoveAll(tmpDir)

	goMod := fmt.Sprintf(`module verifycrudcontracts

go 1.25.0

require github.com/coldsmirk/vef-framework-go v0.0.0

replace github.com/coldsmirk/vef-framework-go => %s
`, sourceRoot)
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o600); err != nil {
		return []string{fmt.Sprintf("failed to write temp go.mod: %v", err)}
	}

	mainGo := `package main

import (
	"fmt"
	"time"

	"github.com/coldsmirk/vef-framework-go/api"
	"github.com/coldsmirk/vef-framework-go/crud"
	"github.com/coldsmirk/vef-framework-go/orm"
)

type Model struct {
	ID        string ` + "`" + `bun:"id,pk"` + "`" + `
	CreatedBy string ` + "`" + `bun:"created_by"` + "`" + `
	UpdatedBy string ` + "`" + `bun:"updated_by"` + "`" + `
}

type Params struct{}
type Search struct{}

func mustAction(label string, spec api.OperationSpec, want string) {
	if spec.Action != want {
		panic(fmt.Sprintf("%s action got %q want %q", label, spec.Action, want))
	}
}

func main() {
	mustAction("create", crud.NewCreate[Model, Params]().Build(func() {}), crud.RPCActionCreate)
	mustAction("update", crud.NewUpdate[Model, Params]().Build(func() {}), crud.RPCActionUpdate)
	mustAction("delete", crud.NewDelete[Model]().Build(func() {}), crud.RPCActionDelete)
	mustAction("create_many", crud.NewCreateMany[Model, Params]().Build(func() {}), crud.RPCActionCreateMany)
	mustAction("update_many", crud.NewUpdateMany[Model, Params]().Build(func() {}), crud.RPCActionUpdateMany)
	mustAction("delete_many", crud.NewDeleteMany[Model]().Build(func() {}), crud.RPCActionDeleteMany)
	mustAction("find_one", crud.NewFindOne[Model, Search]().Build(func() {}), crud.RPCActionFindOne)
	mustAction("find_all", crud.NewFindAll[Model, Search]().Build(func() {}), crud.RPCActionFindAll)
	mustAction("find_page", crud.NewFindPage[Model, Search]().Build(func() {}), crud.RPCActionFindPage)
	mustAction("find_options", crud.NewFindOptions[Model, Search]().Build(func() {}), crud.RPCActionFindOptions)
	mustAction("find_tree", crud.NewFindTree[Model, Search](func(models []Model) []Model { return models }).Build(func() {}), crud.RPCActionFindTree)
	mustAction("find_tree_options", crud.NewFindTreeOptions[Model, Search]().Build(func() {}), crud.RPCActionFindTreeOptions)
	mustAction("export", crud.NewExport[Model, Search]().Build(func() {}), crud.RPCActionExport)
	mustAction("import", crud.NewImport[Model]().Build(func() {}), crud.RPCActionImport)

	spec := crud.NewCreate[Model, Params]().
		EnableAudit().
		Timeout(time.Second).
		Public().
		RequiredPermission("sys:test").
		RateLimit(7, time.Minute).
		Build(func() {})
	if !spec.EnableAudit || spec.Timeout != time.Second || !spec.Public || spec.RequiredPermission != "sys:test" {
		panic(fmt.Sprintf("base builder flags not preserved: %#v", spec))
	}
	if spec.RateLimit == nil || spec.RateLimit.Max != 7 || spec.RateLimit.Period != time.Minute {
		panic(fmt.Sprintf("base builder rate limit not preserved: %#v", spec.RateLimit))
	}

	relations := crud.GetAuditUserNameRelations((*Model)(nil), "display_name")
	if len(relations) != 2 {
		panic(fmt.Sprintf("audit relation count got %d", len(relations)))
	}
	if relations[0].Alias != "creator" || relations[0].JoinType != orm.JoinLeft || relations[0].ForeignColumn != orm.ColumnCreatedBy ||
		len(relations[0].SelectedColumns) != 1 || relations[0].SelectedColumns[0].Alias != orm.ColumnCreatedByName ||
		relations[1].Alias != "updater" || relations[1].ForeignColumn != orm.ColumnUpdatedBy ||
		len(relations[1].SelectedColumns) != 1 || relations[1].SelectedColumns[0].Alias != orm.ColumnUpdatedByName {
		panic(fmt.Sprintf("audit relations changed: %#v", relations))
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainGo), 0o600); err != nil {
		return []string{fmt.Sprintf("failed to write temp main.go: %v", err)}
	}

	cmd := exec.Command("go", "run", "-mod=mod", ".")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("CRUD runtime contract check failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func runSourceTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./crud/...", "-count=1")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./crud/... failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func loadCrudAuditEntries(path string) []auditEntry {
	ledger := loadJSON[auditLedger](path)
	var entries []auditEntry
	for _, entry := range ledger.Entries {
		if entry.Package == crudPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

	return entries
}

func loadCrudManifestEntry(path string) manifestEntry {
	m := loadJSON[manifest](path)
	for _, entry := range m.Packages {
		if entry.Package == crudPackage {
			return entry
		}
	}
	panic("API audit manifest missing CRUD package")
}

func loadCrudContract(path string) (contractPackageReview, contractEntry) {
	ledger := loadJSON[contractLedger](path)
	var review contractPackageReview
	for _, item := range ledger.PackageReviews {
		if item.Package == crudPackage {
			review = item
			break
		}
	}
	if review.Package == "" {
		panic("API contract ledger missing CRUD package review")
	}
	var contract contractEntry
	for _, item := range ledger.Entries {
		if item.ID == crudContractID() {
			contract = item
			break
		}
	}
	if contract.ID == "" {
		panic("API contract ledger missing CRUD contract entry")
	}

	return review, contract
}

func crudContractID() string {
	return crudPackage + "#dynamic-resource:crud-action-format-import-export"
}

func loadLiveInventory(sourceRoot, docsRoot string) map[string]manifestEntry {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
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

func actionConstantSymbols(entries []auditEntry) []string {
	var symbols []string
	for _, entry := range entries {
		if entry.Kind == "top" && (strings.HasPrefix(entry.Symbol, "RPCAction") || strings.HasPrefix(entry.Symbol, "RESTAction")) {
			symbols = append(symbols, entry.Symbol)
		}
	}
	sort.Strings(symbols)

	return symbols
}

func entryBySymbol(entries []auditEntry, symbol string) auditEntry {
	for _, entry := range entries {
		if entry.Symbol == symbol {
			return entry
		}
	}

	return auditEntry{}
}

func constantValue(signature string) string {
	_, value, ok := strings.Cut(signature, " = ")
	if !ok {
		return ""
	}
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		unquoted, err := strconv.Unquote(value)
		if err == nil {
			return unquoted
		}
	}

	return value
}

func shortCrudSignature(signature string) string {
	signature = strings.ReplaceAll(signature, crudPackage+".", "")
	signature = strings.ReplaceAll(signature, "github.com/coldsmirk/vef-framework-go/orm.", "orm.")
	signature = strings.ReplaceAll(signature, "github.com/gofiber/fiber/v3.", "fiber.")
	prefix := regexp.MustCompile(`^[A-Za-z0-9_]+ : `)

	return prefix.ReplaceAllString(signature, "")
}

func entryTypePrefix(symbol string) string {
	prefix, _, ok := strings.Cut(symbol, ".")
	if !ok {
		return ""
	}

	return prefix + "."
}

func jsonFieldFromSignature(signature string) string {
	match := regexp.MustCompile(`json:\\\"([^\\\",]+)`).FindStringSubmatch(signature)
	if len(match) == 2 {
		return match[1]
	}

	return ""
}

func packageSection(content, pkg string) string {
	marker := "## " + pkg
	start := strings.Index(content, marker)
	if start < 0 {
		marker = "### `" + pkg + "`"
		start = strings.Index(content, marker)
		if start < 0 {
			return ""
		}
	}
	rest := content[start:]
	next := regexp.MustCompile(`(?m)^##+ `).FindAllStringIndex(rest[len(marker):], -1)
	if len(next) == 0 {
		return rest
	}

	return rest[:len(marker)+next[0][0]]
}

func documentMentionsSymbol(content, symbol string) bool {
	needles := []string{
		"`" + symbol + "`",
		"`" + symbol,
		"`" + symbol + "[",
		"type " + symbol,
		"func " + symbol,
		symbol + " :",
	}
	for _, needle := range needles {
		if strings.Contains(content, needle) {
			return true
		}
	}

	return false
}

func isProcessorSignatureSymbol(symbol string) bool {
	return symbol == "Processor" || symbol == "FilenameBuilder" || strings.HasSuffix(symbol, "Processor")
}

func parseGoFile(path string) *ast.File {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	return file
}

func combineCorpora(label string, corpora ...corpus) corpus {
	var b strings.Builder
	for _, c := range corpora {
		b.WriteString(c.content)
		b.WriteString("\n")
	}

	return corpus{label: label, content: b.String()}
}

func missingTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !strings.Contains(c.content, term) {
			failures = append(failures, c.label+" missing term "+term)
		}
	}

	return failures
}

func sameSet(left, right []string) bool {
	left = sortedUnique(left)
	right = sortedUnique(right)
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
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


func readCorpus(label, path string) corpus {
	return corpus{label: label, content: readFile(path)}
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func loadJSON[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}

	return result
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()

	return output.String(), err
}
