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
	mcpPackage = "github.com/coldsmirk/vef-framework-go/mcp"

	mcpFingerprint = "40d5ab2a18c99a3c369d907963fdde2f413699f63bef1a2cca9ae9174f0a2067"
	mcpTopLevel    = 44
	mcpFields      = 11
	mcpMethods     = 64
	mcpEntries     = 119

	mcpGroupedEntries              = 75
	mcpGroupedFields               = 11
	mcpGroupedMethods              = 64
	mcpGroupedReceivers            = 27
	mcpGroupedSignatureFingerprint = "e74c2f6abe36c57b9d13f45b4fd99fb934bcbbb549536a6048bb5c99edfdd803"
	mcpGroupedReceiverFingerprint  = "40cd9bdc2846dd7956f2512d95c1ad927fb400623a402dcbbf12698b99b87410"

	englishMCPPath   = "docs/features/mcp.md"
	chineseMCPPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/mcp.md"
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

	checks := []struct {
		sourcePath      string
		sourceTerms     []string
		docTerms        []string
		englishDocTerms []string
		chineseDocTerms []string
	}{
		{
			sourcePath: "mcp/mcp.go",
			sourceTerms: []string{
				"ToolProvider interface", "Tools() []ToolDefinition",
				"ResourceProvider interface", "Resources() []ResourceDefinition",
				"ResourceTemplateProvider interface", "ResourceTemplates() []ResourceTemplateDefinition",
				"PromptProvider interface", "Prompts() []PromptDefinition",
				"ServerInfo struct", "Name         string", "Version      string",
				"Instructions string", "type (", "Server         = mcp.Server",
				"ServerOptions  = mcp.ServerOptions", "ServerSession  = mcp.ServerSession",
				"Implementation = mcp.Implementation", "Tool            = mcp.Tool",
				"ToolHandler     = mcp.ToolHandler", "CallToolRequest = mcp.CallToolRequest",
				"CallToolResult  = mcp.CallToolResult", "Resource            = mcp.Resource",
				"ResourceTemplate    = mcp.ResourceTemplate", "Prompt           = mcp.Prompt",
				"PromptHandler    = mcp.PromptHandler", "Content      = mcp.Content",
				"TextContent  = mcp.TextContent", "ImageContent = mcp.ImageContent",
				"AudioContent = mcp.AudioContent", "ResourceNotFoundError = mcp.ResourceNotFoundError",
				"NewToolResultText", "IsError: true",
			},
			docTerms: []string{
				"mcp.ToolProvider", "mcp.ResourceProvider", "mcp.ResourceTemplateProvider",
				"mcp.PromptProvider", "mcp.ToolDefinition", "mcp.ResourceDefinition",
				"mcp.ResourceTemplateDefinition", "mcp.PromptDefinition", "mcp.ServerInfo",
				"Server", "ServerOptions", "ServerSession", "Implementation",
				"Tool", "CallToolRequest", "CallToolResult", "ToolHandler",
				"Resource", "ReadResourceRequest", "ReadResourceResult", "ResourceHandler",
				"ResourceTemplate", "Prompt", "PromptArgument", "PromptMessage",
				"Role", "GetPromptParams", "GetPromptRequest", "GetPromptResult",
				"PromptHandler", "Content", "TextContent", "ImageContent",
				"AudioContent", "Annotations", "mcp.ResourceNotFoundError",
				"mcp.NewToolResultText(text)", "mcp.NewToolResultError(message)",
			},
		},
		{
			sourcePath: "mcp/auth.go",
			sourceTerms: []string{
				"GetPrincipalFromContext", "TokenInfoFromContext(ctx)",
				"tokenInfo.Extra[\"principal\"].(*security.Principal)",
				"return security.PrincipalAnonymous", "DBWithOperator",
				"orm.PlaceholderKeyOperator", "principal.ID",
			},
			docTerms: []string{
				"mcp.GetPrincipalFromContext(ctx)", "security.PrincipalAnonymous",
				"mcp.DBWithOperator(ctx, db)", "orm.PlaceholderKeyOperator",
			},
		},
		{
			sourcePath: "mcp/schema.go",
			sourceTerms: []string{
				"Anonymous: true", "DoNotReference: true",
				"RequiredFromJSONSchemaTags: true", "SchemaFor[T any]()",
				"SchemaOf(v any)", "if v == nil", "return nil",
				"MustSchemaFor[T any]()", "MustSchemaOf(v any)",
				"panic(\"mcp: failed to generate schema\")", "delete(resultVal, \"$schema\")",
				"jsonschema_description", "jsonschema_extras",
			},
			docTerms: []string{
				"mcp.SchemaFor[T]()", "mcp.SchemaOf(v)", "mcp.MustSchemaFor[T]()",
				"mcp.MustSchemaOf(v)", "jsonschema:\"required\"", "$ref", "$id",
				"$schema", "mcp.SchemaOf(nil)", "nil", "mcp.MustSchemaOf(nil)",
				"mcp: failed to generate schema", "required", "nullable",
				"title=...", "description=...", "type=...", "anchor=...",
				"default=...", "example=...", "enum=...", "oneof_required=...",
				"anyof_required=...", "oneof_ref=...", "oneof_type=...",
				"anyof_ref=...", "anyof_type=...", "minLength=...",
				"maxLength=...", "pattern=...", "format=...", "readOnly=true",
				"writeOnly=true", "minimum=...", "maximum=...", "exclusiveMinimum=...",
				"exclusiveMaximum=...", "multipleOf=...", "minItems=...",
				"maxItems=...", "uniqueItems=true", "jsonschema_description",
				"jsonschema_extras",
			},
		},
		{
			sourcePath: "config/mcp.go",
			sourceTerms: []string{
				"MCPConfig struct", "Enabled bool `config:\"enabled\"`",
				"RequireAuth *bool `config:\"require_auth\"`",
				"unset value defaults to secure", "allow anonymous access",
			},
			docTerms: []string{
				"vef.mcp.enabled = true", "vef.mcp.require_auth = false",
			},
			englishDocTerms: []string{"secure by default"},
			chineseDocTerms: []string{"默认是安全的"},
		},
		{
			sourcePath: "internal/mcp/server.go",
			sourceTerms: []string{
				"if !params.MCPConfig.Enabled", "smcp.NewServer",
				"Name:    getServerName(params)", "Version: getServerVersion(params)",
				"Instructions: getInstructions(params)", "registerTools",
				"registerResources", "registerResourceTemplates", "registerPrompts",
				"return \"vef-mcp-server\"", "return \"v1.0.0\"",
			},
			docTerms: []string{
				"mcp.ServerInfo", "vef.app.name", "vef-mcp-server", "v1.0.0",
			},
			englishDocTerms: []string{"MCP server only activates", "default instructions are empty"},
			chineseDocTerms: []string{"MCP server 只有在", "默认 instructions 为空"},
		},
		{
			sourcePath: "internal/mcp/handler.go",
			sourceTerms: []string{
				"mcp.NewStreamableHTTPHandler", "new(mcp.StreamableHTTPOptions)",
				"params.MCPConfig.RequireAuth == nil || *params.MCPConfig.RequireAuth",
				"applyAuthMiddleware", "auth.RequireBearerToken",
			},
			docTerms: []string{
				"Streamable HTTP MCP endpoint", "require_auth",
			},
			englishDocTerms: []string{"Bearer-token verification", "framework token authentication path"},
			chineseDocTerms: []string{"Bearer token 校验", "框架的 token authentication path"},
		},
		{
			sourcePath: "internal/mcp/middleware.go",
			sourceTerms: []string{
				"const mcpPath = \"/mcp\"", "return \"mcp\"",
				"return 500", "router.All(mcpPath, m.handler.FiberHandler())",
				"all methods",
			},
			docTerms: []string{
				"/mcp", "order `500`",
			},
			englishDocTerms: []string{"all HTTP methods"},
			chineseDocTerms: []string{"所有 HTTP method"},
		},
		{
			sourcePath: "internal/mcp/tools/query.go",
			sourceTerms: []string{
				"Name:        \"database_query\"", "read-only (SELECT)",
				"InputSchema: mcp.MustSchemaFor[QueryArgs]()", "SQL    string",
				"Params []any", "if args.SQL == \"\"", "Sql parameter is required",
				"sqlguard.EnsureReadOnly", "Only read-only (SELECT) queries are permitted",
				"db := mcp.DBWithOperator(ctx, t.db)", "convertByteSlices",
				"mcp.NewToolResultText(string(jsonBytes))", "mcp.NewToolResultError",
			},
			docTerms: []string{
				"database_query", "read-only `SELECT`", "params",
				"MCP tool error result", "mcp.DBWithOperator(...)",
			},
			englishDocTerms: []string{
				"`sql` is required", "SQL query string using `?` placeholders",
				"JSON text content", "UTF-8 byte slices", "Base64 strings",
			},
			chineseDocTerms: []string{
				"`sql` 必填", "使用 `?` 占位符", "JSON 文本内容", "UTF-8 的 `[]byte`", "Base64 string",
			},
		},
		{
			sourcePath: "internal/orm/sqlguard/guard.go",
			sourceTerms: []string{
				"func EnsureReadOnly", "ErrSQLParseFailed", "ErrNotReadOnly",
				"data-modifying CTEs", "firstDangerousFunction",
				"pg_read_file", "pg_sleep", "nextval", "setval",
			},
			docTerms: []string{
				"fail-closed", "data-modifying CTE", "pg_read_file",
				"pg_sleep", "nextval", "setval",
			},
		},
		{
			sourcePath: "internal/mcp/prompts/naming_master.go",
			sourceTerms: []string{
				"Name:        \"naming-master\"", "Senior IT naming expert",
				"database objects", "audit fields", "indexes", "constraints",
				"foreign key strategies", "Role:    mcp.Role(\"user\")",
				"Text: namingMasterPromptContent",
			},
			docTerms: []string{
				"naming-master",
			},
			englishDocTerms: []string{
				"code identifiers", "database objects", "audit fields",
				"indexes", "constraints", "foreign key strategy",
			},
			chineseDocTerms: []string{
				"代码 identifier", "数据库对象", "审计字段",
				"索引", "约束", "外键策略",
			},
		},
		{
			sourcePath: "internal/mcp/handler_test.go",
			sourceTerms: []string{
				"TestUnsetRequireAuthDefaultsToSecure", "http.StatusUnauthorized",
				"TestExplicitFalseAllowsAnonymous", "require_auth=false must allow anonymous access",
			},
			docTerms: []string{
				"vef.mcp.require_auth", "vef.mcp.require_auth = false",
				"Bearer <token>", "bearer <token>",
			},
			englishDocTerms: []string{
				"If `vef.mcp.require_auth` is omitted or\nset to `true`",
				"Set `vef.mcp.require_auth = false` only",
				"secure by default",
			},
			chineseDocTerms: []string{
				"`vef.mcp.require_auth` 未配置或设置为 `true`",
				"`vef.mcp.require_auth = false`",
				"默认是安全的",
			},
		},
		{
			sourcePath: "internal/mcp/mcp_test.go",
			sourceTerms: []string{
				"LowerCaseBearer", "NoPrefix", "Should accept lowercase bearer prefix",
				"Should reject token without Bearer prefix",
			},
			docTerms: []string{
				"Bearer <token>", "bearer <token>",
			},
			englishDocTerms: []string{
				"raw token without a Bearer prefix",
			},
			chineseDocTerms: []string{
				"没有 Bearer\nprefix 的裸 token",
			},
		},
	}

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)
	manifestPath := filepath.Join(docsRoot, "scripts/api-audit-manifest.json")
	auditLedgerPath := filepath.Join(docsRoot, "scripts/api-audit-ledger.json")
	contractLedgerPath := filepath.Join(docsRoot, "scripts/api-contract-ledger.json")

	englishDocs := readCorpus("English MCP docs", filepath.Join(docsRoot, englishMCPPath))
	chineseDocs := readCorpus("Chinese MCP docs", filepath.Join(docsRoot, chineseMCPPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))
	docs := []corpus{englishDocs, chineseDocs}

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[mcpPackage]
	liveAuditEntries := mcpEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	mcpEntries := mcpEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(mcpEntries)...)
	failures = append(failures, verifyLiveAuditEntries(mcpEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedMCPSurface(mcpEntries, docs)...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, mcpEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, mcpEntries)...)
	failures = append(failures, verifyMCPDocs(docs)...)
	for _, check := range checks {
		source := readCorpus(check.sourcePath, filepath.Join(sourceRoot, check.sourcePath))
		failures = append(failures, missingTerms(source, check.sourceTerms)...)
		for _, doc := range docs {
			failures = append(failures, missingTerms(doc, check.docTerms)...)
		}
		failures = append(failures, missingTerms(englishDocs, check.englishDocTerms)...)
		failures = append(failures, missingTerms(chineseDocs, check.chineseDocTerms)...)
	}
	failures = append(failures, runGoTests(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Printf("MCP contract docs verified: %d public entries, %d grouped entries, %d source files, %d doc mirrors\n", len(mcpEntries), mcpGroupedEntries, len(checks), len(docs))
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != mcpPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != mcpTopLevel ||
		entry.Fields != mcpFields ||
		entry.Methods != mcpMethods ||
		entry.Fingerprint != mcpFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s MCP surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			mcpTopLevel, mcpFields, mcpMethods, mcpFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != mcpPackage {
			continue
		}

		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, mcpCoverage()) {
			failures = append(failures, fmt.Sprintf("MCP manifest coverage mismatch: got %v want %v", entry.Coverage, mcpCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing MCP package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != mcpPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "MCP contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != mcpTopLevel ||
			review.ReviewedSurface.Fields != mcpFields ||
			review.ReviewedSurface.Methods != mcpMethods ||
			review.ReviewedSurface.EntryCount != mcpEntries ||
			review.ReviewedSurface.Fingerprint != mcpFingerprint {
			failures = append(failures, fmt.Sprintf(
				"MCP contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, mcpCoverage()) {
			failures = append(failures, fmt.Sprintf("MCP contract review coverage mismatch: got %v want %v", review.Coverage, mcpCoverage()))
		}
		if !contains(review.ContractIDs, mcpContractID()) {
			failures = append(failures, "MCP contract review missing contract id "+mcpContractID())
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing MCP package review")
	}

	var foundContract bool
	for _, entry := range contracts.Entries {
		if entry.ID != mcpContractID() {
			continue
		}
		foundContract = true
		if entry.Package != mcpPackage || entry.Kind != "dynamic-resource" {
			failures = append(failures, fmt.Sprintf("MCP contract entry shape mismatch: package=%s kind=%s", entry.Package, entry.Kind))
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "MCP contract entry disposition mismatch: "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, mcpCoverage()) {
			failures = append(failures, fmt.Sprintf("MCP contract coverage mismatch: got %v want %v", entry.Coverage, mcpCoverage()))
		}
		for _, term := range []string{
			"/mcp",
			"Streamable HTTP",
			"order `500`",
			"HTTP method",
			"vef.mcp.enabled",
			"require_auth",
			"Bearer",
			"vef-mcp-server",
			"v1.0.0",
			"database_query",
			"naming-master",
			"SchemaFor",
			"SchemaOf(nil)",
			"MustSchemaOf(nil)",
			`jsonschema:"required"`,
			"$schema",
			"PrincipalAnonymous",
			"DBWithOperator",
			"orm.PlaceholderKeyOperator",
		} {
			if !contains(entry.Terms, term) {
				failures = append(failures, "MCP contract missing term "+term)
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
	}
	if !foundContract {
		failures = append(failures, "contract ledger missing MCP contract entry")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != mcpEntries {
		failures = append(failures, fmt.Sprintf("MCP audit entry count mismatch: got %d want %d", len(entries), mcpEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != mcpPackage {
			failures = append(failures, "non-MCP audit entry passed into MCP verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate MCP audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "MCP audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, mcpCoverage()) {
			failures = append(failures, fmt.Sprintf("MCP audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, mcpCoverage()))
		}
	}
	if counts["top"] != mcpTopLevel || counts["field"] != mcpFields || counts["method"] != mcpMethods {
		failures = append(failures, fmt.Sprintf("MCP audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != mcpTopLevel ||
		dispositionCounts["grouped:type-member-family"] != mcpGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"MCP audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			mcpTopLevel,
			mcpGroupedEntries,
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
			failures = append(failures, fmt.Sprintf("MCP missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"MCP live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
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
			failures = append(failures, fmt.Sprintf("MCP extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedMCPSurface(entries []auditEntry, docs []corpus) []string {
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
			failures = append(failures, fmt.Sprintf("MCP grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("MCP grouped type-member surface", rows, mcpGroupedEntries, mcpGroupedSignatureFingerprint)...)
	if kindCounts["field"] != mcpGroupedFields || kindCounts["method"] != mcpGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"MCP grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			mcpGroupedFields,
			mcpGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("MCP grouped receiver/type families", receiverRows, mcpGroupedReceivers, mcpGroupedReceiverFingerprint)...)

	for _, doc := range docs {
		for _, term := range []string{
			"119 public MCP entries",
			"75 grouped MCP field/method entries",
			"27 MCP receiver/type families",
			"11 exported MCP field entries",
			"64 exported MCP method entries",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped MCP audit term "+term)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, mcpPackage)
	if section == "" {
		return []string{index.label + " missing MCP package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s MCP index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyMCPDocs(docs []corpus) []string {
	var failures []string
	for _, doc := range docs {
		for _, term := range []string{
			"MCP SDK pass-through surface",
			"promoted SDK methods",
			"public API index",
			"`ToolDefinition.Tool`",
			"`ToolDefinition.Handler`",
			"`ResourceDefinition.Resource`",
			"`ResourceDefinition.Handler`",
			"`ResourceTemplateDefinition.Template`",
			"`PromptDefinition.Prompt`",
			"`ServerInfo.Name`",
			"`ServerInfo.Version`",
			"`ServerInfo.Instructions`",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing MCP semantic term "+term)
			}
		}
	}

	return failures
}

func runGoTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./mcp", "./internal/mcp/...")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./mcp ./internal/mcp/... failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func mcpEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == mcpPackage {
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

func missingTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !strings.Contains(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term: %s", c.label, term))
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

func mcpContractID() string {
	return mcpPackage + "#dynamic-resource:mcp-endpoint-tool-resource-contract"
}

func mcpCoverage() []string {
	return []string{englishMCPPath}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
