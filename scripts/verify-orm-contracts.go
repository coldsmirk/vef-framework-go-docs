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
	"strings"
)

const (
	ormPackage = "github.com/coldsmirk/vef-framework-go/orm"

	ormGroupedEntryCount           = 1352
	ormGroupedSignatureFingerprint = "d6c50843e9466dda73a696b1c15df3c0bceec85e126ab62f48a05963d21b20a3"
	ormGroupedReceiverFingerprint  = "874cde2d4e0726680812365c2050f8258681d2cac60ba9fcbdf974cbb74b3f58"

	ormBunGroupedEntryCount           = 237
	ormBunGroupedSignatureFingerprint = "69a4d7e969fec812a3eb5504fe01b697d44ed517ee60428ef88193b086b36a5c"

	ormVEFOwnedGroupedEntryCount           = 1115
	ormVEFOwnedGroupedSignatureFingerprint = "a20609b911785f56fa4eae16d869b9bf7f04ed03a51b27977b27b603ea1bd0ed"
)

type corpus struct {
	label   string
	content string
}

type auditLedger struct {
	Entries []auditLedgerEntry `json:"entries"`
}

type auditLedgerEntry struct {
	Package     string `json:"package"`
	Kind        string `json:"kind"`
	Symbol      string `json:"symbol"`
	Signature   string `json:"signature"`
	Disposition string `json:"disposition"`
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	publicBun := readCorpus("public ORM Bun aliases", filepath.Join(sourceRoot, "orm/bun.go"))
	publicFacade := readCorpus("public ORM facade", filepath.Join(sourceRoot, "orm/orm.go"))
	queryTrait := readCorpus("internal ORM query traits", filepath.Join(sourceRoot, "internal/orm/query_trait.go"))
	goMod := readCorpus("source go.mod", filepath.Join(sourceRoot, "go.mod"))
	englishDocs := readCorpora("English ORM docs", docsRoot, []string{
		"docs/data-access/orm-querying.md",
		"docs/data-access/orm-mutations.md",
		"docs/data-access/orm-expressions.md",
		"docs/data-access/orm-ddl.md",
		"docs/data-access/query-builder.md",
	})
	chineseDocs := readCorpora("Chinese ORM docs", docsRoot, []string{
		"i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/orm-querying.md",
		"i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/orm-mutations.md",
		"i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/orm-expressions.md",
		"i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/orm-ddl.md",
		"i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/query-builder.md",
	})
	docs := []corpus{englishDocs, chineseDocs}

	var failures []string
	failures = append(failures, missingTerms(publicBun, []string{
		"BaseModel         = bun.BaseModel",
		"BunSelectQuery    = bun.SelectQuery",
		"BunInsertQuery    = bun.InsertQuery",
		"BunUpdateQuery    = bun.UpdateQuery",
		"BunDeleteQuery    = bun.DeleteQuery",
		"Table             = schema.Table",
		"Field             = schema.Field",
		"Relation          = schema.Relation",
		"Dialect           = schema.Dialect",
		"type BeforeSelectHook interface",
		"BeforeSelect(ctx context.Context, query *BunSelectQuery) error",
		"AfterSelect(ctx context.Context, query *BunSelectQuery) error",
		"BeforeInsert(ctx context.Context, query *BunInsertQuery) error",
		"AfterInsert(ctx context.Context, query *BunInsertQuery) error",
		"BeforeUpdate(ctx context.Context, query *BunUpdateQuery) error",
		"AfterUpdate(ctx context.Context, query *BunUpdateQuery) error",
		"BeforeDelete(ctx context.Context, query *BunDeleteQuery) error",
		"AfterDelete(ctx context.Context, query *BunDeleteQuery) error",
	})...)
	failures = append(failures, missingTerms(publicFacade, []string{
		"DB                          = orm.DB",
		"SelectQuery                 = orm.SelectQuery",
		"InsertQuery                 = orm.InsertQuery",
		"UpdateQuery                 = orm.UpdateQuery",
		"DeleteQuery                 = orm.DeleteQuery",
		"MergeQuery                  = orm.MergeQuery",
		"ConditionBuilder            = orm.ConditionBuilder",
		"ExprBuilder                 = orm.ExprBuilder",
		"WindowCountBuilder          = orm.WindowCountBuilder",
		"DataTypeDef       = orm.DataTypeDef",
		"ApplySort = orm.ApplySort",
		"DataType = orm.DataType",
	})...)
	failures = append(failures, missingTerms(queryTrait, []string{
		"ScanAndCount(ctx context.Context, dest ...any) (int64, error)",
		"Count(ctx context.Context) (int64, error)",
		"Exists(ctx context.Context) (bool, error)",
	})...)
	failures = append(failures, missingTerms(goMod, []string{
		"github.com/uptrace/bun v1.3.0",
	})...)
	failures = append(failures, verifyGroupedORMMethodSurface(docsRoot)...)

	commonDocTerms := []string{
		"VEF-owned ORM method families",
		"receiver/category",
		"orm.SelectQuery",
		"orm.InsertQuery",
		"orm.UpdateQuery",
		"orm.DeleteQuery",
		"orm.MergeQuery",
		"condition",
		"expression",
		"aggregate",
		"window-builder",
		"orm.BunSelectQuery",
		"orm.BunInsertQuery",
		"orm.BunUpdateQuery",
		"orm.BunDeleteQuery",
		"github.com/uptrace/bun",
		"v1.3.0",
		"Bun/schema",
		"Table",
		"Field",
		"Relation",
		"Dialect",
		"orm.SelectQuery.Count",
		"int64",
		"public API index",
		"ScanAndCount(ctx)",
		"Count(ctx)",
		"WithValues(name, model)",
		"WithOrderedValues(name, model)",
		"QueryBuilder.DB()",
	}
	englishDocTerms := []string{
		"Do not read VEF query-interface behavior from the Bun aliases",
		"These rows describe VEF `orm.SelectQuery` execution methods",
		"not the upstream `orm.BunSelectQuery` pass-through alias",
	}
	chineseDocTerms := []string{
		"不要用 Bun aliases 推断 VEF",
		"以下行描述的是 VEF `orm.SelectQuery` 执行方法",
		"不是上游 `orm.BunSelectQuery` pass-through alias",
	}

	for _, doc := range docs {
		failures = append(failures, missingTerms(doc, commonDocTerms)...)
	}
	failures = append(failures, missingTerms(englishDocs, englishDocTerms)...)
	failures = append(failures, missingTerms(chineseDocs, chineseDocTerms)...)
	failures = append(failures, runGoTest(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("ORM contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Println("ORM contract docs verified: public facade aliases, VEF/Bun execution boundary, 1,352 grouped method entries, go test ./orm")
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func readCorpora(label, root string, paths []string) corpus {
	parts := make([]string, 0, len(paths))
	for _, path := range paths {
		part := readCorpus(label, filepath.Join(root, path))
		parts = append(parts, part.content)
	}

	return corpus{label: label, content: strings.Join(parts, "\n")}
}

func loadAuditLedger(path string) auditLedger {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read API audit ledger at %s: %w", path, err))
	}

	var ledger auditLedger
	if err := json.Unmarshal(content, &ledger); err != nil {
		panic(fmt.Errorf("failed to parse API audit ledger at %s: %w", path, err))
	}

	return ledger
}

func verifyGroupedORMMethodSurface(docsRoot string) []string {
	ledger := loadAuditLedger(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))

	var groupedRows []string
	var bunRows []string
	var vefOwnedRows []string
	receiverCounts := map[string]int{}

	var failures []string
	for _, entry := range ledger.Entries {
		if entry.Package != ormPackage || entry.Disposition != "grouped:builder-member-family" {
			continue
		}

		if entry.Kind != "method" {
			failures = append(failures, fmt.Sprintf("ORM grouped entry %s has non-method kind %q", entry.Symbol, entry.Kind))
		}

		receiver, ok := splitReceiver(entry.Symbol)
		if !ok {
			failures = append(failures, fmt.Sprintf("ORM grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}

		row := strings.Join([]string{entry.Symbol, entry.Kind, entry.Signature}, "\t")
		groupedRows = append(groupedRows, row)
		receiverCounts[receiver]++

		if strings.HasPrefix(receiver, "Bun") {
			if !isBunPassThroughReceiver(receiver) {
				failures = append(failures, fmt.Sprintf("ORM grouped entry uses unexpected Bun receiver %q", receiver))
			}
			bunRows = append(bunRows, row)
			continue
		}

		vefOwnedRows = append(vefOwnedRows, row)
	}

	failures = append(failures, verifyCountAndFingerprint(
		"ORM grouped method surface",
		groupedRows,
		ormGroupedEntryCount,
		ormGroupedSignatureFingerprint,
	)...)
	failures = append(failures, verifyCountAndFingerprint(
		"ORM Bun pass-through grouped method surface",
		bunRows,
		ormBunGroupedEntryCount,
		ormBunGroupedSignatureFingerprint,
	)...)
	failures = append(failures, verifyCountAndFingerprint(
		"ORM VEF-owned grouped method surface",
		vefOwnedRows,
		ormVEFOwnedGroupedEntryCount,
		ormVEFOwnedGroupedSignatureFingerprint,
	)...)

	receiverRows := make([]string, 0, len(receiverCounts))
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyCountAndFingerprint(
		"ORM grouped receiver families",
		receiverRows,
		105,
		ormGroupedReceiverFingerprint,
	)...)

	return failures
}

func verifyCountAndFingerprint(label string, rows []string, wantCount int, wantFingerprint string) []string {
	gotFingerprint := fingerprint(rows)
	var failures []string
	if len(rows) != wantCount {
		failures = append(failures, fmt.Sprintf("%s count mismatch: got %d want %d", label, len(rows), wantCount))
	}
	if gotFingerprint != wantFingerprint {
		failures = append(failures, fmt.Sprintf("%s fingerprint mismatch: got %s want %s", label, gotFingerprint, wantFingerprint))
	}

	return failures
}

func fingerprint(rows []string) string {
	sorted := append([]string(nil), rows...)
	sort.Strings(sorted)

	hash := sha256.New()
	for _, row := range sorted {
		hash.Write([]byte(row))
		hash.Write([]byte("\n"))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func splitReceiver(symbol string) (string, bool) {
	parts := strings.SplitN(symbol, ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}

	return parts[0], true
}

func isBunPassThroughReceiver(receiver string) bool {
	switch receiver {
	case "BunSelectQuery", "BunInsertQuery", "BunUpdateQuery", "BunDeleteQuery":
		return true
	default:
		return false
	}
}

func missingTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !containsTerm(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term: %s", c.label, term))
		}
	}

	return failures
}

func containsTerm(content, term string) bool {
	if strings.Contains(content, term) {
		return true
	}

	return strings.Contains(normalizeWhitespace(content), normalizeWhitespace(term))
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func runGoTest(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./orm")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./orm failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
