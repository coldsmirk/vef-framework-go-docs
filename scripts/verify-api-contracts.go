package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const apiPackage = "github.com/coldsmirk/vef-framework-go/api"

type corpus struct {
	label   string
	content string
}

type auditLedger struct {
	Entries []auditLedgerEntry `json:"entries"`
}

type auditLedgerEntry struct {
	ID        string `json:"id"`
	Package   string `json:"package"`
	Kind      string `json:"kind"`
	Symbol    string `json:"symbol"`
	Signature string `json:"signature"`
}

type contractLedger struct {
	PackageReviews []contractPackageReview `json:"package_reviews"`
}

type contractPackageReview struct {
	Package         string              `json:"package"`
	ReviewedSurface publicSurfaceReview `json:"reviewed_surface"`
}

type publicSurfaceReview struct {
	TopLevel    int    `json:"top_level"`
	Fields      int    `json:"fields"`
	Methods     int    `json:"methods"`
	EntryCount  int    `json:"entry_count"`
	Fingerprint string `json:"fingerprint"`
}

type indexEntry struct {
	kind      string
	symbol    string
	signature string
	line      string
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishAPIDocs := readCorpus("English API docs", filepath.Join(docsRoot, "docs/guide/api.md"))
	chineseAPIDocs := readCorpus("Chinese API docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/api.md"))
	englishHandlerDocs := readCorpus("English custom handlers docs", filepath.Join(docsRoot, "docs/guide/custom-handlers.md"))
	chineseHandlerDocs := readCorpus("Chinese custom handlers docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/custom-handlers.md"))
	englishRoutingDocs := readCorpus("English routing docs", filepath.Join(docsRoot, "docs/guide/routing.md"))
	chineseRoutingDocs := readCorpus("Chinese routing docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/routing.md"))
	englishParamsDocs := readCorpus("English params and meta docs", filepath.Join(docsRoot, "docs/guide/params-and-meta.md"))
	chineseParamsDocs := readCorpus("Chinese params and meta docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/params-and-meta.md"))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	ledger := readAuditLedger(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	contract := readContractLedger(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	expectedEntries := apiLedgerEntries(ledger)

	var failures []string
	failures = append(failures, verifyContractReview(contract)...)
	failures = append(failures, verifyIndexSection(englishIndex, expectedEntries)...)
	failures = append(failures, verifyIndexSection(chineseIndex, expectedEntries)...)
	failures = append(failures, runAPIAuditGate(sourceRoot, docsRoot)...)

	englishTopicDocs := combineCorpora(
		"English API topic docs",
		englishAPIDocs,
		englishHandlerDocs,
		englishRoutingDocs,
		englishParamsDocs,
	)
	chineseTopicDocs := combineCorpora(
		"Chinese API topic docs",
		chineseAPIDocs,
		chineseHandlerDocs,
		chineseRoutingDocs,
		chineseParamsDocs,
	)
	failures = append(failures, missingTerms(englishTopicDocs, englishPublicDocSurfaceTerms())...)
	failures = append(failures, missingTerms(chineseTopicDocs, chinesePublicDocSurfaceTerms())...)

	failures = append(failures, missingTerms(englishAPIDocs, []string{
		"66 top-level exported symbols",
		"45 exported fields",
		"40 exported methods",
		"fingerprint is `a8ef51431b3e8661cfd8687b36c7b9a5458651788c2454fbd165c700f19f5b3e`",
		"`api.NewRPCResource` and `api.NewRESTResource` validate the resource name",
		"construction time",
		"They panic when validation fails",
		"`api.ValidateActionName(action, kind)`",
		"RPC/REST action validation",
		"`Kind.String()` returns `rpc`",
		"`unknown`",
		"`get`, `post`,",
		"`trace`, `connect`, and `all`",
		"Sub-resource paths may contain `/`",
		"each segment must use kebab-case",
		"dynamic Fiber params such as `/:id` are not accepted",
		"custom `OperationsProvider` implementations should still produce action strings that already satisfy",
		"`Identifier.String()` is promoted to `Operation.String()` and `Request.String()`",
		"`Action` | required; direct `WithOperations(...)` specs",
		"`30s`",
		"`Max=100`, `Period=5m`",
		"`HasRateLimit()` returns true only when `RateLimit != nil` and `RateLimit.Max > 0`",
		"`RequiresAuth()` assumes `Auth` is non-nil",
		"`Auth.Strategy != api.AuthStrategyNone`",
		"`Params.Decode` and `Meta.Decode` require a pointer to a struct",
	})...)
	failures = append(failures, missingTerms(chineseAPIDocs, []string{
		"66 个 top-level exported symbols",
		"45 个 exported fields",
		"40 个 exported methods",
		"fingerprint 是 `a8ef51431b3e8661cfd8687b36c7b9a5458651788c2454fbd165c700f19f5b3e`",
		"`api.NewRPCResource` 和 `api.NewRESTResource` 会在构造期校验 resource",
		"校验失败会 panic",
		"`api.ValidateActionName(action, kind)`",
		"RPC/REST action",
		"`Kind.String()` 对 `KindRPC` 返回 `rpc`",
		"`unknown`",
		"`trace`、`connect` 和 `all`",
		"sub-resource 路径可以包含 `/`",
		"每一段都必须是 kebab-case",
		"动态 Fiber 参数如 `/:id` 不会被公开 validator 接受",
		"自定义 `OperationsProvider` 仍应产出已经满足",
		"`Identifier.String()` 会提升为 `Operation.String()` 和 `Request.String()`",
		"`Action` | 必填；直接 `WithOperations(...)`",
		"`30s`",
		"`Max=100`、`Period=5m`",
		"`HasRateLimit()` 只有在 `RateLimit != nil` 且 `RateLimit.Max > 0` 时才返回 true",
		"`RequiresAuth()` 假定 `Auth` 非 nil",
		"`Auth.Strategy != api.AuthStrategyNone`",
		"`Params.Decode` 和 `Meta.Decode` 都要求传入 struct 指针",
	})...)
	failures = append(failures, missingTerms(englishHandlerDocs, []string{
		"REST does not infer handler methods",
		"REST action format",
		"`get admin/users`",
		"sub-resource paths may contain `/`, but each segment must use kebab-case",
		"Any other return shape is invalid for a direct handler",
		"`fiber.Ctx`",
		"`orm.DB`",
		"`logx.Logger`",
		"`*security.Principal`",
		"`api.Params`",
		"`api.Meta`",
		"typed struct embedding `api.P`",
		"typed struct embedding `api.M`",
		"`page.Pageable`",
		"`cron.Scheduler`",
		"`event.Bus`",
		"`mold.Transformer`",
		"`storage.Service`",
		"`datasource.Registry`",
		"direct field match",
		"tagged dive field match",
		"embedded field match",
		"`storage.Files`",
	})...)
	failures = append(failures, missingTerms(chineseHandlerDocs, []string{
		"REST 不会从 action 字符串推导方法名",
		"REST action 格式",
		"`get admin/users`",
		"sub-resource path 可以包含 `/`，但每一段都必须使用 kebab-case",
		"其他返回形态都不是合法的直接 handler",
		"`fiber.Ctx`",
		"`orm.DB`",
		"`logx.Logger`",
		"`*security.Principal`",
		"`api.Params`",
		"`api.Meta`",
		"嵌入 `api.P` 的 typed struct",
		"嵌入 `api.M` 的 typed struct",
		"`page.Pageable`",
		"`cron.Scheduler`",
		"`event.Bus`",
		"`mold.Transformer`",
		"`storage.Service`",
		"`datasource.Registry`",
		"direct field match",
		"tagged dive field match",
		"embedded field match",
		"`storage.Files`",
	})...)
	failures = append(failures, missingTerms(englishRoutingDocs, []string{
		"`DefaultRPCEndpoint`",
		"`/api`",
		"`FormKeyParams`",
		"`FormKeyMeta`",
		"`X-Meta-*`",
		"lowercased",
		"`/api/<resource>/<subpath>`",
	})...)
	failures = append(failures, missingTerms(chineseRoutingDocs, []string{
		"`DefaultRPCEndpoint`",
		"`/api`",
		"`FormKeyParams`",
		"`FormKeyMeta`",
		"`X-Meta-*`",
		"小写",
		"`/api/<resource>/<subpath>`",
	})...)
	failures = append(failures, missingTerms(englishParamsDocs, []string{
		"`api.P`",
		"`api.M`",
		"`Params.Decode`",
		"`Meta.Decode`",
		"pointer to a struct",
	})...)
	failures = append(failures, missingTerms(chineseParamsDocs, []string{
		"`api.P`",
		"`api.M`",
		"`Params.Decode`",
		"`Meta.Decode`",
		"struct 指针",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "api/resource.go",
			terms: []string{
				"type Resource interface",
				"type RouterStrategy interface",
				"type Engine interface",
				"type Kind uint8",
				"func (k Kind) String() string",
				"validHTTPVerbs = collections.NewHashSetFrom(",
				"\"trace\"",
				"\"connect\"",
				"\"all\"",
				"func ValidateActionName(action string, kind Kind) error",
				"snakeCasePattern.MatchString(action)",
				"strings.SplitN(action, \" \", 3)",
				"restResourceNamePattern.MatchString(subRes)",
				"ValidateActionName(op.Action, r.kind)",
				"func NewRESTResource(name string, opts ...ResourceOption) Resource",
				"func NewRPCResource(name string, opts ...ResourceOption) Resource",
				"panic(err)",
				"func WithVersion(v string) ResourceOption",
				"func WithOperations(ops ...OperationSpec) ResourceOption",
				"func WithAuth(auth *AuthConfig) ResourceOption",
			},
		},
		{
			path: "api/operation.go",
			terms: []string{
				"type OperationSpec struct",
				"Action string",
				"RequiredPermission string",
				"RateLimit *RateLimitConfig",
				"type Operation struct",
				"Identifier",
				"func (o *Operation) HasRateLimit() bool",
				"return o.RateLimit != nil && o.RateLimit.Max > 0",
				"func (o *Operation) RequiresAuth() bool",
				"return o.Auth.Strategy != AuthStrategyNone",
				"type RateLimitConfig struct",
				"type OperationsProvider interface",
				"type OperationsCollector interface",
			},
		},
		{
			path: "api/request.go",
			terms: []string{
				"type Identifier struct",
				"Resource string `json:\"resource\"",
				"Action   string `json:\"action\"",
				"Version  string `json:\"version\"",
				"func (id Identifier) String() string",
				"return id.Resource + \":\" + id.Action + \":\" + id.Version",
				"type Params map[string]any",
				"func (p Params) Decode(out any) error",
				"type Meta map[string]any",
				"func (m Meta) Decode(out any) error",
				"reflectx.IsPointerToStruct",
				"type Request struct",
				"Identifier",
				"Params Params `json:\"params\"`",
				"Meta   Meta   `json:\"meta\"`",
				"func (r *Request) GetParam(key string) (any, bool)",
				"func (r *Request) GetMeta(key string) (any, bool)",
			},
		},
		{
			path: "api/auth.go",
			terms: []string{
				"AuthStrategyNone      = \"none\"",
				"AuthStrategyBearer    = \"bearer\"",
				"AuthStrategySignature = \"signature\"",
				"type AuthConfig struct",
				"func Public() *AuthConfig",
				"func BearerAuth() *AuthConfig",
				"func SignatureAuth() *AuthConfig",
				"func (c *AuthConfig) Clone() *AuthConfig",
				"type AuthStrategy interface",
				"type AuthStrategyRegistry interface",
			},
		},
		{
			path: "api/audit.go",
			terms: []string{
				"type AuditEvent struct",
				"func (*AuditEvent) EventType() string",
				"func SubscribeAuditEvent(",
				"bus event.Bus,",
			},
		},
		{
			path: "api/api_errors.go",
			terms: []string{
				"var (",
				"ErrInvalidRequestParams",
				"ErrInvalidRequestMeta",
			},
		},
		{
			path: "api/errors.go",
			terms: []string{
				"ErrEmptyResourceName",
				"ErrInvalidResourceName",
				"ErrInvalidActionName",
				"ErrInvalidParamsType",
				"ErrInvalidMetaType",
			},
		},
		{
			path: "api/header.go",
			terms: []string{
				"HeaderXMetaPrefix = \"X-Meta-\"",
				"HeaderXTimestamp  = \"X-Timestamp\"",
				"HeaderXNonce      = \"X-Nonce\"",
				"HeaderXSignature  = \"X-Signature\"",
				"HeaderXAppID      = \"X-App-ID\"",
			},
		},
		{
			path: "api/version.go",
			terms: []string{
				"VersionV1 = \"v1\"",
				"VersionV9 = \"v9\"",
			},
		},
		{
			path: "api/sentinel.go",
			terms: []string{
				"type P struct{}",
				"type M struct{}",
			},
		},
		{
			path: "api/handler.go",
			terms: []string{
				"type Middleware interface",
				"type HandlerResolver interface",
				"type HandlerAdapter interface",
				"type HandlerParamResolver interface",
				"type FactoryParamResolver interface",
			},
		},
		{
			path: "internal/api/engine.go",
			terms: []string{
				"type engine struct",
				"defaultVersion: api.VersionV1",
				"defaultTimeout: 30 * time.Second",
				"defaultAuth:    api.BearerAuth()",
				"Max:    100",
				"Period: 5 * time.Minute",
				"func (e *engine) Register(resources ...api.Resource) error",
				"buildOperation",
				"if spec.Action == \"\"",
				"api.Public()",
				"ac.Options[shared.AuthOptionRequiredPermission] = spec.RequiredPermission",
			},
		},
		{
			path: "internal/api/router/rest.go",
			terms: []string{
				"func (*REST) parseAction(action string) (method, subPath string)",
				"method = strings.ToUpper(method)",
				"func (*REST) buildPath(resource, subPath string) string",
				"return \"/\" + resource + subPath",
				"strings.CutPrefix(key, api.HeaderXMetaPrefix)",
				"req.Meta[strings.ToLower(metaKey)] = values[0]",
			},
		},
		{
			path: "internal/api/router/rpc.go",
			terms: []string{
				"DefaultRPCEndpoint = \"/api\"",
				"FormKeyParams = \"params\"",
				"FormKeyMeta   = \"meta\"",
			},
		},
		{
			path: "internal/api/param/module.go",
			terms: []string{
				"NewCtxResolver()",
				"NewDBResolver()",
				"NewLoggerResolver()",
				"NewPrincipalResolver()",
				"NewParamsResolver()",
				"NewMetaResolver()",
				"NewSchedulerResolver",
				"NewBusResolver",
				"NewTransformerResolver",
				"NewStorageResolver",
				"NewDataSourcesResolver",
				"NewDBFactoryResolver",
				"NewSchedulerFactoryResolver",
				"NewBusFactoryResolver",
				"NewTransformerFactoryResolver",
				"NewStorageFactoryResolver",
				"NewFilesFactoryResolver",
				"NewDataSourcesFactoryResolver",
			},
		},
		{
			path: "internal/api/param/resolver_manager.go",
			terms: []string{
				"embedsAPIParams(paramType)",
				"buildParamsResolver(paramType)",
				"embedsAPIMeta(paramType) || isBuiltinMetaType(paramType)",
				"buildMetaResolver(paramType)",
			},
		},
		{
			path: "internal/api/param/helpers.go",
			terms: []string{
				"searchDirectFields",
				"searchTaggedFields",
				"searchEmbeddedFields",
				"reflectx.WithDiveTag(\"api\", \"dive\")",
				"apiParamsType = reflect.TypeFor[api.P]()",
				"apiMetaType   = reflect.TypeFor[api.M]()",
				"reflect.TypeFor[page.Pageable]()",
			},
		},
		{
			path: "internal/api/collector/resource_provider.go",
			terms: []string{
				"type ResourceProviderCollector struct",
				"func (*ResourceProviderCollector) Collect(resource api.Resource) []api.OperationSpec",
				"specs := resource.Operations()",
				"return specs",
			},
		},
		{
			path: "internal/api/collector/embedded_provider.go",
			terms: []string{
				"type EmbeddedProviderCollector struct",
				"api.OperationsProvider",
				"provider.Provide()",
			},
		},
	}
	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	failures = append(failures, runPackageTests(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("api contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("API contract docs verified: 66 top-level public symbols, 40 public methods, 45 public fields, %d source/runtime files, 8 doc mirrors\n", len(sourceChecks))
}

func readAuditLedger(path string) auditLedger {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read audit ledger %s: %w", path, err))
	}

	var ledger auditLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(fmt.Errorf("failed to parse audit ledger %s: %w", path, err))
	}

	return ledger
}

func readContractLedger(path string) contractLedger {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read contract ledger %s: %w", path, err))
	}

	var ledger contractLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(fmt.Errorf("failed to parse contract ledger %s: %w", path, err))
	}

	return ledger
}

func apiLedgerEntries(ledger auditLedger) map[string]auditLedgerEntry {
	entries := make(map[string]auditLedgerEntry)
	for _, entry := range ledger.Entries {
		if entry.Package != apiPackage {
			continue
		}
		key := entry.Kind + ":" + entry.Symbol
		if _, exists := entries[key]; exists {
			panic("duplicate API audit ledger entry " + key)
		}
		entries[key] = entry
	}

	if len(entries) != 151 {
		panic(fmt.Sprintf("expected 151 API audit ledger entries, got %d", len(entries)))
	}

	return entries
}

func verifyContractReview(ledger contractLedger) []string {
	for _, review := range ledger.PackageReviews {
		if review.Package != apiPackage {
			continue
		}

		surface := review.ReviewedSurface
		var failures []string
		if surface.TopLevel != 66 {
			failures = append(failures, fmt.Sprintf("contract review top_level mismatch: got %d", surface.TopLevel))
		}
		if surface.Fields != 45 {
			failures = append(failures, fmt.Sprintf("contract review fields mismatch: got %d", surface.Fields))
		}
		if surface.Methods != 40 {
			failures = append(failures, fmt.Sprintf("contract review methods mismatch: got %d", surface.Methods))
		}
		if surface.EntryCount != 151 {
			failures = append(failures, fmt.Sprintf("contract review entry_count mismatch: got %d", surface.EntryCount))
		}
		if surface.Fingerprint != "a8ef51431b3e8661cfd8687b36c7b9a5458651788c2454fbd165c700f19f5b3e" {
			failures = append(failures, "contract review fingerprint mismatch: got "+surface.Fingerprint)
		}

		return failures
	}

	return []string{"missing API package review in api-contract-ledger.json"}
}

func verifyIndexSection(doc corpus, expected map[string]auditLedgerEntry) []string {
	section := extractPackageSection(doc, "## "+apiPackage)
	actual := parseIndexSection(doc.label, section)
	var failures []string

	for key, entry := range expected {
		found, ok := actual[key]
		if !ok {
			failures = append(failures, fmt.Sprintf("%s missing API index entry %s", doc.label, key))

			continue
		}
		if found.signature != entry.Signature {
			failures = append(failures, fmt.Sprintf("%s stale API index signature for %s: got %q want %q", doc.label, key, found.signature, entry.Signature))
		}
	}
	for key := range actual {
		if _, ok := expected[key]; !ok {
			failures = append(failures, fmt.Sprintf("%s has stale API index entry %s", doc.label, key))
		}
	}

	requiredPromoted := []string{
		"field:Operation.Action",
		"field:Operation.Resource",
		"field:Operation.Version",
		"field:Request.Action",
		"field:Request.Resource",
		"field:Request.Version",
		"method:Operation.String",
		"method:Request.String",
	}
	for _, key := range requiredPromoted {
		if _, ok := actual[key]; !ok {
			failures = append(failures, fmt.Sprintf("%s missing promoted API index entry %s", doc.label, key))
		}
	}

	return failures
}

func extractPackageSection(doc corpus, heading string) string {
	start := strings.Index(doc.content, heading+"\n")
	if start < 0 {
		panic(fmt.Sprintf("%s missing section heading %q", doc.label, heading))
	}
	section := doc.content[start:]
	next := strings.Index(section[len(heading)+1:], "\n## ")
	if next >= 0 {
		section = section[:len(heading)+1+next]
	}

	return section
}

func parseIndexSection(label string, section string) map[string]indexEntry {
	entries := make(map[string]indexEntry)
	currentType := ""
	for _, line := range strings.Split(section, "\n") {
		switch {
		case strings.HasPrefix(line, "TYPE "):
			signature := strings.TrimPrefix(line, "TYPE ")
			symbol := symbolFromSignature(signature)
			currentType = symbol
			addIndexEntry(label, entries, indexEntry{kind: "top", symbol: symbol, signature: signature, line: line})
		case strings.HasPrefix(line, "FUNC "):
			signature := strings.TrimPrefix(line, "FUNC ")
			symbol := symbolFromSignature(signature)
			currentType = ""
			addIndexEntry(label, entries, indexEntry{kind: "top", symbol: symbol, signature: signature, line: line})
		case strings.HasPrefix(line, "CONST "):
			signature := strings.TrimPrefix(line, "CONST ")
			symbol := symbolFromSignature(signature)
			currentType = ""
			addIndexEntry(label, entries, indexEntry{kind: "top", symbol: symbol, signature: signature, line: line})
		case strings.HasPrefix(line, "VAR "):
			signature := strings.TrimPrefix(line, "VAR ")
			symbol := symbolFromSignature(signature)
			currentType = ""
			addIndexEntry(label, entries, indexEntry{kind: "top", symbol: symbol, signature: signature, line: line})
		case strings.HasPrefix(line, "  FIELD "):
			if currentType == "" {
				panic(fmt.Sprintf("%s contains FIELD without current TYPE: %q", label, line))
			}
			signature := strings.TrimPrefix(line, "  FIELD ")
			symbol := currentType + "." + symbolFromSignature(signature)
			addIndexEntry(label, entries, indexEntry{kind: "field", symbol: symbol, signature: signature, line: line})
		case strings.HasPrefix(line, "  METHOD "):
			if currentType == "" {
				panic(fmt.Sprintf("%s contains METHOD without current TYPE: %q", label, line))
			}
			signature := strings.TrimPrefix(line, "  METHOD ")
			symbol := currentType + "." + symbolFromSignature(signature)
			addIndexEntry(label, entries, indexEntry{kind: "method", symbol: symbol, signature: signature, line: line})
		}
	}

	return entries
}

func addIndexEntry(label string, entries map[string]indexEntry, entry indexEntry) {
	key := entry.kind + ":" + entry.symbol
	if _, exists := entries[key]; exists {
		panic(fmt.Sprintf("%s contains duplicate API index entry %s", label, key))
	}
	entries[key] = entry
}

func symbolFromSignature(signature string) string {
	if i := strings.Index(signature, " : "); i >= 0 {
		return signature[:i]
	}
	if i := strings.Index(signature, " = "); i >= 0 {
		return signature[:i]
	}

	return signature
}

func englishPublicDocSurfaceTerms() []string {
	return []string{
		"`api.Resource`",
		"`api.RouterStrategy`",
		"`api.Engine`",
		"`api.Kind`",
		"`api.KindRPC`",
		"`api.KindREST`",
		"`api.ValidateActionName(action, kind) error`",
		"`api.NewRPCResource(name, opts...)`",
		"`api.NewRESTResource(name, opts...)`",
		"`api.WithVersion(v)`",
		"`api.WithAuth(config)`",
		"`api.WithOperations(specs...)`",
		"`api.OperationSpec`",
		"`api.Operation`",
		"`api.RateLimitConfig`",
		"`api.OperationsProvider`",
		"`api.OperationsCollector`",
		"`api.Identifier`",
		"`api.Request`",
		"`api.Params`",
		"`api.Meta`",
		"`api.P`",
		"`api.M`",
		"`api.AuthConfig`",
		"`api.Public()`",
		"`api.BearerAuth()`",
		"`api.SignatureAuth()`",
		"`api.AuthStrategy`",
		"`api.AuthStrategyRegistry`",
		"`api.Middleware`",
		"`api.HandlerResolver`",
		"`api.HandlerAdapter`",
		"`api.HandlerParamResolver`",
		"`api.FactoryParamResolver`",
		"`api.AuditEvent`",
		"`api.SubscribeAuditEvent`",
		"`api.HeaderXMetaPrefix`",
		"`api.HeaderXTimestamp`",
		"`api.HeaderXNonce`",
		"`api.HeaderXSignature`",
		"`api.HeaderXAppID`",
		"`api.VersionV1`",
		"`api.VersionV9`",
		"`api.ErrInvalidRequestParams`",
		"`api.ErrInvalidRequestMeta`",
		"`api.ErrInvalidParamsType`",
		"`api.ErrInvalidMetaType`",
	}
}

func chinesePublicDocSurfaceTerms() []string {
	return []string{
		"`api.Resource`",
		"`api.RouterStrategy`",
		"`api.Engine`",
		"`api.Kind`",
		"`api.KindRPC`",
		"`api.KindREST`",
		"`api.ValidateActionName(action, kind) error`",
		"`api.NewRPCResource(name, opts...)`",
		"`api.NewRESTResource(name, opts...)`",
		"`api.WithVersion(v)`",
		"`api.WithAuth(config)`",
		"`api.WithOperations(specs...)`",
		"`api.OperationSpec`",
		"`api.Operation`",
		"`api.RateLimitConfig`",
		"`api.OperationsProvider`",
		"`api.OperationsCollector`",
		"`api.Identifier`",
		"`api.Request`",
		"`api.Params`",
		"`api.Meta`",
		"`api.P`",
		"`api.M`",
		"`api.AuthConfig`",
		"`api.Public()`",
		"`api.BearerAuth()`",
		"`api.SignatureAuth()`",
		"`api.AuthStrategy`",
		"`api.AuthStrategyRegistry`",
		"`api.Middleware`",
		"`api.HandlerResolver`",
		"`api.HandlerAdapter`",
		"`api.HandlerParamResolver`",
		"`api.FactoryParamResolver`",
		"`api.AuditEvent`",
		"`api.SubscribeAuditEvent`",
		"`api.HeaderXMetaPrefix`",
		"`api.HeaderXTimestamp`",
		"`api.HeaderXNonce`",
		"`api.HeaderXSignature`",
		"`api.HeaderXAppID`",
		"`api.VersionV1`",
		"`api.VersionV9`",
		"`api.ErrInvalidRequestParams`",
		"`api.ErrInvalidRequestMeta`",
		"`api.ErrInvalidParamsType`",
		"`api.ErrInvalidMetaType`",
	}
}

func combineCorpora(label string, docs ...corpus) corpus {
	var content strings.Builder
	for _, doc := range docs {
		content.WriteString("\n--- ")
		content.WriteString(doc.label)
		content.WriteString(" ---\n")
		content.WriteString(doc.content)
	}

	return corpus{label: label, content: content.String()}
}

func missingTerms(doc corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !containsTerm(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term %q", doc.label, term))
		}
	}

	return failures
}

func containsTerm(content string, term string) bool {
	if strings.Contains(content, term) {
		return true
	}

	return strings.Contains(normalizeSpace(content), normalizeSpace(term))
}

func normalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func runAPIAuditGate(sourceRoot string, docsRoot string) []string {
	cmd := exec.Command(
		"go",
		"run",
		filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", filepath.Join(docsRoot, "scripts/api-audit-manifest.json"),
		"-ledger", filepath.Join(docsRoot, "scripts/api-audit-ledger.json"),
		"-contract-ledger", filepath.Join(docsRoot, "scripts/api-contract-ledger.json"),
	)
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("global API audit gate failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./api", "./internal/api/...")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./api ./internal/api/... failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func readCorpus(label string, path string) corpus {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(data)}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("failed to resolve %s: %w", path, err))
	}

	return filepath.Clean(abs)
}
