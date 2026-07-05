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
	"regexp"
	"sort"
	"strings"
)

const (
	configPackage                     = "github.com/coldsmirk/vef-framework-go/config"
	configFingerprint                 = "f0c4b5df8283faa4a53bbeb3c0a86f03df34c384b81253792d827db1fdd61a65"
	configTopLevel                    = 53
	configFields                      = 112
	configMethods                     = 23
	configEntries                     = 188
	configGroupedEntries              = 135
	configGroupedFields               = 112
	configGroupedMethods              = 23
	configGroupedReceivers            = 21
	configGroupedSignatureFingerprint = "f3a1fa4f2281d8c0d89c1a3163e03bda423541b95368d9d15eaefb866d0d2dac"
	configGroupedReceiverFingerprint  = "7272919e35fb155fd48e83b8285f875daf39b904d7f48fbc236fa358d77e74c5"

	englishReferencePath = "docs/reference/configuration-reference.md"
	chineseReferencePath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/configuration-reference.md"
	englishGuidePath     = "docs/getting-started/configuration.md"
	chineseGuidePath     = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/getting-started/configuration.md"
	englishIndexPath     = "docs/reference/public-api-index.md"
	chineseIndexPath     = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	Terms          []string `json:"terms"`
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishReference := readCorpus("English configuration reference", docsRoot, englishReferencePath)
	chineseReference := readCorpus("Chinese configuration reference", docsRoot, chineseReferencePath)
	englishGuide := readCorpus("English getting-started configuration", docsRoot, englishGuidePath)
	chineseGuide := readCorpus("Chinese getting-started configuration", docsRoot, chineseGuidePath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadConfigAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadConfigManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadConfigContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveConfigEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyGroupedConfigSurface(entries, []corpus{englishReference, chineseReference})...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
	}
	for _, reference := range []corpus{englishReference, chineseReference} {
		failures = append(failures, verifyReferenceSurface(reference, entries)...)
		failures = append(failures, verifyNoPhantomConfigRefs(reference, entries)...)
		failures = append(failures, missingTerms(reference, referenceTerms(reference.label))...)
	}
	for _, guide := range []corpus{englishGuide, chineseGuide} {
		failures = append(failures, verifyNoPhantomConfigRefs(guide, entries)...)
		failures = append(failures, missingTerms(guide, guideTerms(guide.label))...)
	}

	failures = append(failures, verifyContractLedger(review, contract, sourceRoot)...)
	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, runExecutableConfigChecks(sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("config contracts verified: 135 grouped configuration entries")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != configPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != configTopLevel || entry.Fields != configFields ||
		entry.Methods != configMethods || entry.Fingerprint != configFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			configTopLevel, configFields, configMethods, configFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != configPackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != configTopLevel ||
		review.ReviewedSurface.Fields != configFields ||
		review.ReviewedSurface.Methods != configMethods ||
		review.ReviewedSurface.EntryCount != configEntries ||
		review.ReviewedSurface.Fingerprint != configFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, configContractID()) {
		failures = append(failures, "contract package review missing config contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != configEntries {
		failures = append(failures, fmt.Sprintf("config audit entry count mismatch: got %d want %d", len(entries), configEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != configPackage {
			failures = append(failures, "non-config audit entry passed into config verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate config audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "config audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != configTopLevel || counts["field"] != configFields || counts["method"] != configMethods {
		failures = append(failures, fmt.Sprintf("config audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}

	return failures
}

func verifyGroupedConfigSurface(entries []auditEntry, docs []corpus) []string {
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
			failures = append(failures, fmt.Sprintf("config grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}

		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint(
		"config grouped configuration surface",
		rows,
		configGroupedEntries,
		configGroupedSignatureFingerprint,
	)...)
	if kindCounts["field"] != configGroupedFields || kindCounts["method"] != configGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"config grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			configGroupedFields,
			configGroupedMethods,
		))
	}

	receiverRows := make([]string, 0, len(receiverCounts))
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint(
		"config grouped receiver families",
		receiverRows,
		configGroupedReceivers,
		configGroupedReceiverFingerprint,
	)...)

	commonTerms := []string{
		"135 grouped configuration entries",
		"112 exported configuration fields",
		"23 exported configuration methods",
	}
	for _, doc := range docs {
		for _, term := range commonTerms {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped config audit term "+term)
			}
		}
		if strings.HasPrefix(doc.label, "Chinese") {
			if !containsNormalized(doc.content, "21 个 config struct/interface families") {
				failures = append(failures, doc.label+" missing grouped config audit term 21 个 config struct/interface families")
			}
			continue
		}
		for _, term := range []string{
			"135 grouped configuration entries",
			"21 config struct/interface families",
			"112 exported configuration fields",
			"23 exported configuration methods",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped config audit term "+term)
			}
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
	expected := []string{englishReferencePath, englishGuidePath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest config coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review config coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry config coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, configPackage)
	if section == "" {
		return []string{index.label + " missing config package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s config index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyReferenceSurface(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		term := "`config." + entry.Symbol + "`"
		if !strings.Contains(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing audited config entry %s", doc.label, term))
		}
	}

	return failures
}

func verifyNoPhantomConfigRefs(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		valid["config."+entry.Symbol] = true
	}

	var failures []string
	refs := configReferencesFromMarkdownCode(doc.content)
	for _, ref := range refs {
		if valid[ref] {
			continue
		}
		if ref == "config.StorageConfig.Effective" {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown config public API: %s", doc.label, ref))
	}

	return failures
}

func referenceTerms(label string) []string {
	common := []string{
		configFingerprint,
		"`DataSourcesConfig.Map`",
		"`config.PrimaryDataSourceName`",
		"`config.Config.Unmarshal(key, target)`",
		"`config.DataSourcesConfig.Primary()`",
		"`config.ApprovalConfig.ApplyDefaults()`",
		"`config.StorageConfig.Effective...`",
		"`config.EventConfig.Validate()`",
		"ErrInboxRetentionTooShort",
		"oracle",
		"sqlserver",
		"enabled = true",
		"nil `*redis.Client`",
		"32mib",
		"access token",
		"refresh token",
		"DefaultMaxUploadSize",
		"1073741824",
		"1 GiB",
		"90d",
		"168h",
		"path.Match",
	}
	if strings.HasPrefix(label, "Chinese") {
		return append(common,
			"53 个 top-level exported symbols",
			"112 个 exported fields",
			"23 个 exported methods",
			"public surface fingerprint 是 `"+configFingerprint+"`",
			"方法会返回零值 `config.DataSourceConfig`",
			"它不会启用 `AutoMigrate`",
			"框架注入的是 nil `*redis.Client`",
			"内置 JWT token generator 签发的 access token 固定 `30m` 过期",
			"Bearer token",
		)
	}

	return append(common,
		"53 top-level exported symbols",
		"112 exported fields",
		"23 exported methods",
		"public surface fingerprint is `"+configFingerprint+"`",
		"the method returns the zero `config.DataSourceConfig`",
		"It does not enable `AutoMigrate`",
		"the framework provides a nil `*redis.Client`",
		"access tokens issued by the built-in JWT token generator expire after `30m`",
		"Bearer auth",
	)
}

func guideTerms(label string) []string {
	common := []string{
		"Configuration Reference",
		"vef.data_sources.primary",
		"config.DBKind",
		"oracle",
		"sqlserver",
		"database.ErrUnsupportedDBKind",
		"token_expires",
		"refresh-token",
		"`30m`",
		"vef.redis.enabled",
		"nil `*redis.Client`",
		"skips startup `PING`",
		"enabled = true",
		"require_auth",
		"Bearer auth",
		"AutoMigrate",
		"32mib",
	}
	if strings.HasPrefix(label, "Chinese") {
		return []string{
			"配置参考",
			"vef.data_sources.primary",
			"config.DBKind",
			"oracle",
			"sqlserver",
			"database.ErrUnsupportedDBKind",
			"token_expires",
			"refresh token",
			"`30m`",
			"vef.redis.enabled",
			"nil `*redis.Client`",
			"跳过启动 `PING`",
			"enabled = true",
			"require_auth",
			"Bearer auth",
			"AutoMigrate",
			"32mib",
		}
	}

	return common
}

func verifyContractLedger(review contractPackageReview, contract contractEntry, sourceRoot string) []string {
	var failures []string
	if contract.ID != configContractID() {
		failures = append(failures, fmt.Sprintf("config contract entry id mismatch: got %q", contract.ID))
	}
	if contract.Package != configPackage {
		failures = append(failures, fmt.Sprintf("config contract entry package mismatch: got %q", contract.Package))
	}
	if contract.Kind != "configuration-contract" {
		failures = append(failures, fmt.Sprintf("config contract entry kind mismatch: got %q", contract.Kind))
	}
	if contract.Contract != "configuration keys, environment override, and effective defaults" {
		failures = append(failures, fmt.Sprintf("config contract label mismatch: got %q", contract.Contract))
	}

	expectedEvidence := configSourceEvidence(sourceRoot)
	if !sameSet(review.SourceEvidence, expectedEvidence) {
		failures = append(failures, fmt.Sprintf("config package review source evidence mismatch: got %v want %v", review.SourceEvidence, expectedEvidence))
	}
	if !sameSet(contract.SourceEvidence, expectedEvidence) {
		failures = append(failures, fmt.Sprintf("config contract source evidence mismatch: got %v want %v", contract.SourceEvidence, expectedEvidence))
	}
	for _, term := range []string{
		"vef.app",
		"trusted_proxies",
		"vef.data_sources.primary",
		"PrimaryDataSourceName",
		"VEF_CONFIG_PATH",
		"VEF_I18N_LANGUAGE",
		"token_expires",
		"vef.redis",
		"enabled",
		"vef.storage",
		"DefaultMaxUploadSize",
		"EffectiveMaxUploadSize",
		"require_auth",
		"ApplyDefaults()",
		"ErrInboxRetentionTooShort",
		"sample_interval",
	} {
		if !contains(contract.Terms, term) {
			failures = append(failures, "config contract missing term "+term)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	checks := []struct {
		path  string
		terms []string
	}{
		{
			path: "internal/config/config.go",
			terms: []string{
				"c.TagName = \"config\"",
				"c.IgnoreUntaggedFields = true",
				"v.AddConfigPath(\"./configs\")",
				"v.AddConfigPath(\"$\" + config.EnvConfigPath)",
				"v.AddConfigPath(\".\")",
				"v.AddConfigPath(\"../configs\")",
			},
		},
		{
				path: "config/env.go",
				terms: []string{
					"EnvPrefix       = \"VEF\"",
					"EnvNodeID       = EnvPrefix + \"_NODE_ID\"",
					"EnvLogLevel     = EnvPrefix + \"_LOG_LEVEL\"",
					"EnvConfigPath   = EnvPrefix + \"_CONFIG_PATH\"",
					"EnvI18NLanguage = EnvPrefix + \"_I18N_LANGUAGE\"",
				},
			},
		{
			path: "config/data_sources.go",
			terms: []string{
				"const PrimaryDataSourceName = \"primary\"",
				"Oracle    DBKind = \"oracle\"",
				"SQLServer DBKind = \"sqlserver\"",
				"Postgres  DBKind = \"postgres\"",
				"MySQL     DBKind = \"mysql\"",
				"SQLite    DBKind = \"sqlite\"",
				"Map map[string]DataSourceConfig",
				"return c.Map[PrimaryDataSourceName]",
			},
		},
		{
			path: "internal/config/data_sources.go",
			terms: []string{
				"ErrPrimaryDataSourceMissing",
				"cfg.Unmarshal(\"vef.data_sources\", &sources)",
				"sources[config.PrimaryDataSourceName]",
				"return &config.DataSourcesConfig{Map: sources}, nil",
			},
		},
		{
			path: "internal/database/provider.go",
			terms: []string{
				"registry.register(sqlite.NewProvider())",
				"registry.register(postgres.NewProvider())",
				"registry.register(mysql.NewProvider())",
			},
		},
		{
			path: "internal/database/factory.go",
			terms: []string{
				"provider, exists := registry.lookup(cfg.Kind)",
				"return nil, newUnsupportedDBKindError(cfg.Kind)",
			},
		},
		{
			path: "internal/app/fiber.go",
			terms: []string{
				"bodyLimitStr := lo.CoalesceOrEmpty(strings.TrimSpace(cfg.BodyLimit), \"32mib\")",
				"trustProxy := len(cfg.TrustedProxies) > 0",
				"proxyHeader = fiber.HeaderXForwardedFor",
				"TrustProxy:       trustProxy",
				"TrustProxyConfig: fiber.TrustProxyConfig{Proxies: cfg.TrustedProxies}",
			},
		},
		{
			path: "internal/redis/redis.go",
			terms: []string{
				"if !cfg.Enabled",
				"return nil",
				"Network:               lo.CoalesceOrEmpty(cfg.Network, \"tcp\")",
				"lo.CoalesceOrEmpty(cfg.Host, \"/run/redis/redis.sock\")",
				"lo.CoalesceOrEmpty(cfg.Host, \"127.0.0.1\")",
				"lo.CoalesceOrEmpty(cfg.Port, 6379)",
			},
		},
		{
			path: "internal/redis/module.go",
			terms: []string{
				"if client == nil",
				"return nil",
				"client.Ping(ctx).Err()",
			},
		},
		{
			path: "config/mcp.go",
			terms: []string{
				"RequireAuth *bool `config:\"require_auth\"`",
			},
		},
		{
			path: "internal/mcp/handler.go",
			terms: []string{
				"if params.MCPConfig.RequireAuth == nil || *params.MCPConfig.RequireAuth",
				"httpHandler = applyAuthMiddleware(httpHandler, params.AuthManager)",
			},
		},
		{
			path: "internal/security/module.go",
			terms: []string{
				"cfg.TokenExpires = RefreshTokenExpires",
				"cfg.RefreshNotBefore = AccessTokenExpires / 2",
				"cfg.LoginRateLimit = 6",
				"cfg.RefreshRateLimit = 1",
				"case \"\":",
				"generated, err := security.GenerateSecret()",
				"case security.DefaultJWTSecret:",
			},
		},
		{
			path: "internal/security/jwt_token_generator.go",
			terms: []string{
				"AccessTokenExpires  = time.Minute * 30",
				"RefreshTokenExpires = time.Hour * 24 * 7",
				"return g.jwt.Generate(claimsBuilder, AccessTokenExpires, 0)",
				"return g.jwt.Generate(claimsBuilder, g.refreshExpires, g.refreshNotBefore)",
			},
		},
		{
			path: "config/storage.go",
			terms: []string{
				"StorageMinIO      StorageProvider = \"minio\"",
				"StorageMemory     StorageProvider = \"memory\"",
				"StorageFilesystem StorageProvider = \"filesystem\"",
				"DefaultMaxUploadSize int64 = 1024 * 1024 * 1024",
				"DefaultClaimTTL         time.Duration = 24 * time.Hour",
				"DefaultMaxPendingClaims int           = 100",
				"DefaultSweepInterval  time.Duration = 5 * time.Minute",
				"DefaultSweepBatchSize int           = 200",
				"DefaultDeleteWorkerInterval time.Duration = 5 * time.Minute",
				"DefaultDeleteBatchSize      int           = 100",
				"DefaultDeleteConcurrency    int           = 8",
				"DefaultDeleteMaxAttempts    int           = 12",
				"DefaultDeleteLeaseWindow    time.Duration = 5 * time.Minute",
			},
		},
		{
			path: "internal/storage/storage.go",
			terms: []string{
				"if provider == \"\"",
				"provider = config.StorageMemory",
				"storage provider not configured; defaulting to in-memory storage",
				"case config.StorageMinIO:",
				"case config.StorageMemory:",
				"case config.StorageFilesystem:",
			},
		},
		{
			path: "internal/storage/filesystem/service.go",
			terms: []string{
				"if root == \"\"",
				"root = \"./storage\"",
			},
		},
		{
			path: "internal/storage/minio/service.go",
			terms: []string{
				"bucket: lo.CoalesceOrEmpty(cfg.Bucket, appCfg.Name, \"vef-app\")",
			},
		},
		{
			path: "config/approval.go",
			terms: []string{
				"func (c *ApprovalConfig) ApplyDefaults()",
				"c.TimeoutScanInterval = time.Minute",
				"c.PreWarningScanInterval = 5 * time.Minute",
				"c.CleanupScanInterval = 24 * time.Hour",
				"c.DelegationMaxDepth = 10",
				"c.FormSnapshotRetention = 90 * 24 * time.Hour",
				"c.UrgeRecordRetention = 30 * 24 * time.Hour",
				"c.CCRecordRetention = 90 * 24 * time.Hour",
			},
		},
		{
			path: "internal/config/constructors.go",
			terms: []string{
				"approvalConfig.ApplyDefaults()",
				"return unmarshalConfig(cfg, \"vef.event\", new(config.EventConfig))",
			},
		},
		{
			path: "config/event.go",
			terms: []string{
				"return cmp.Or(c.DefaultTransport, \"memory\")",
				"return coalescePositive(c.AsyncQueueSize, 4096)",
				"return coalescePositive(c.AsyncWorkers, 4)",
				"return coalescePositive(c.PublishTimeout, 5*time.Second)",
				"return coalescePositive(c.CleanupInterval, time.Hour)",
				"return coalescePositive(c.CompletedTTL, 7*24*time.Hour)",
				"return coalescePositive(c.Retention, 7*24*time.Hour)",
				"return coalescePositive(c.ProcessingLease, 10*time.Minute)",
				"ErrInboxRetentionTooShort",
				"if !c.Middleware.Inbox || !c.Transports.Outbox.Enabled",
				"maxRetries = 10",
				"backoffSecs := math.Pow(2, float64(maxRetries+1)) - 2",
				"return time.Duration(math.MaxInt64)",
			},
		},
		{
			path: "event/transport/memory/memory.go",
			terms: []string{
				"return 1024",
				"return FullPolicyError",
			},
		},
		{
			path: "event/transport/outbox/outbox.go",
			terms: []string{
				"return 10 * time.Second",
				"return 10",
				"return 100",
				"return 4",
				"return 15 * time.Second",
				"return \"memory\"",
			},
		},
		{
			path: "event/transport/redisstream/redis_stream.go",
			terms: []string{
				"return \"vef:events:\"",
				"return 5 * time.Second",
				"return 60 * time.Second",
				"return 30 * time.Second",
				"return 64",
				"return \"0\"",
			},
		},
		{
			path: "internal/monitor/config.go",
			terms: []string{
				"DefaultSampleInterval = 10 * time.Second",
				"DefaultSampleDuration = 2 * time.Second",
			},
		},
		{
			path: "internal/monitor/service.go",
			terms: []string{
				"resolved := DefaultConfig()",
				"if cfg.SampleInterval > 0",
				"if cfg.SampleDuration > 0",
				"resolved.ExcludedMounts = cfg.ExcludedMounts",
			},
		},
	}

	var failures []string
	for _, check := range checks {
		doc := readSourceFile(sourceRoot, check.path)
		failures = append(failures, missingTerms(doc, check.terms)...)
	}

	return failures
}

func runExecutableConfigChecks(sourceRoot string) []string {
	code := `package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/coldsmirk/vef-framework-go/config"
)

func check(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}

func main() {
	ds := (&config.DataSourcesConfig{Map: map[string]config.DataSourceConfig{
		config.PrimaryDataSourceName: {Kind: config.SQLite},
	}}).Primary()
	check(ds.Kind == config.SQLite, "Primary() did not return the primary data source")
	check((&config.DataSourcesConfig{}).Primary() == (config.DataSourceConfig{}), "missing primary should return zero data source config")

	approval := &config.ApprovalConfig{AutoMigrate: true}
	approval.ApplyDefaults()
	check(approval.AutoMigrate, "ApplyDefaults should not disable AutoMigrate")
	check(approval.TimeoutScanInterval == time.Minute, "approval timeout scan default mismatch")
	check(approval.PreWarningScanInterval == 5*time.Minute, "approval pre-warning default mismatch")
	check(approval.CleanupScanInterval == 24*time.Hour, "approval cleanup default mismatch")
	check(approval.DelegationMaxDepth == 10, "approval delegation default mismatch")
	check(approval.FormSnapshotRetention == 90*24*time.Hour, "approval form snapshot retention mismatch")
	check(approval.UrgeRecordRetention == 30*24*time.Hour, "approval urge retention mismatch")
	check(approval.CCRecordRetention == 90*24*time.Hour, "approval CC retention mismatch")

	storage := config.StorageConfig{
		MaxUploadSize:        -1,
		ClaimTTL:             -time.Second,
		MaxPendingClaims:     -1,
		SweepInterval:        -time.Second,
		SweepBatchSize:       -1,
		DeleteWorkerInterval: -time.Second,
		DeleteBatchSize:      -1,
		DeleteConcurrency:    -1,
		DeleteMaxAttempts:    -1,
		DeleteLeaseWindow:    -time.Second,
	}
	check(storage.EffectiveMaxUploadSize() == config.DefaultMaxUploadSize, "storage max upload default mismatch")
	check(storage.EffectiveClaimTTL() == config.DefaultClaimTTL, "storage claim TTL default mismatch")
	check(storage.EffectiveMaxPendingClaims() == config.DefaultMaxPendingClaims, "storage pending claims default mismatch")
	check(storage.EffectiveSweepInterval() == config.DefaultSweepInterval, "storage sweep interval default mismatch")
	check(storage.EffectiveSweepBatchSize() == config.DefaultSweepBatchSize, "storage sweep batch size default mismatch")
	check(storage.EffectiveDeleteWorkerInterval() == config.DefaultDeleteWorkerInterval, "storage delete worker default mismatch")
	check(storage.EffectiveDeleteBatchSize() == config.DefaultDeleteBatchSize, "storage delete batch size default mismatch")
	check(storage.EffectiveDeleteConcurrency() == config.DefaultDeleteConcurrency, "storage delete concurrency default mismatch")
	check(storage.EffectiveDeleteMaxAttempts() == config.DefaultDeleteMaxAttempts, "storage delete attempts default mismatch")
	check(storage.EffectiveDeleteLeaseWindow() == config.DefaultDeleteLeaseWindow, "storage delete lease default mismatch")

	event := config.EventConfig{}
	check(event.EffectiveDefaultTransport() == "memory", "event default transport mismatch")
	check(event.EffectiveAsyncQueueSize() == 4096, "event async queue size mismatch")
	check(event.EffectiveAsyncWorkers() == 4, "event async workers mismatch")
	check(event.EffectivePublishTimeout() == 5*time.Second, "event publish timeout mismatch")
	outbox := config.EventOutboxTransportConfig{}
	check(outbox.EffectiveCleanupInterval() == time.Hour, "outbox cleanup default mismatch")
	check(outbox.EffectiveCompletedTTL() == 7*24*time.Hour, "outbox completed TTL mismatch")
	inbox := config.EventInboxConfig{}
	check(inbox.EffectiveRetention() == 7*24*time.Hour, "inbox retention default mismatch")
	check(inbox.EffectiveProcessingLease() == 10*time.Minute, "inbox processing lease mismatch")
	check(inbox.EffectiveCleanupInterval() == time.Hour, "inbox cleanup default mismatch")

	event.Middleware.Inbox = true
	event.Transports.Outbox.Enabled = true
	event.Inbox.Retention = time.Minute
	err := event.Validate()
	check(errors.Is(err, config.ErrInboxRetentionTooShort), fmt.Sprintf("Validate should reject short retention, got %v", err))
}
`

	path, err := writeTempGo(code)
	if err != nil {
		return []string{"failed to write executable config check: " + err.Error()}
	}
	defer os.Remove(path)

	return runCommand(sourceRoot, "go", "run", path)
}

func runSourceTests(sourceRoot string) []string {
	commands := [][]string{
		{"go", "test", "./config"},
		{"go", "test", "./internal/app", "-run", "TestCreateFiberAppDefaultBodyLimit"},
		{"go", "test", "./internal/redis", "-run", "Test(NewClientDisabled|BuildRedisAddr)$"},
		{"go", "test", "./internal/mcp", "-run", "TestMCPAuthMode"},
		{"go", "test", "./internal/monitor", "-run", "TestResolveConfig"},
	}

	var failures []string
	for _, command := range commands {
		failures = append(failures, runCommand(sourceRoot, command[0], command[1:]...)...)
	}

	return failures
}

func loadConfigAuditEntries(path string) []auditEntry {
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
		if entry.Package == configPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind != entries[j].Kind {
			return entries[i].Kind < entries[j].Kind
		}

		return entries[i].Symbol < entries[j].Symbol
	})

	return entries
}

func loadConfigManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == configPackage {
			return entry
		}
	}

	panic("API audit manifest missing config package")
}

func loadConfigContract(path string) (contractPackageReview, contractEntry) {
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
		if item.Package == configPackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing config package review")
	}

	var contract contractEntry
	contractFound := false
	for _, item := range ledger.Entries {
		if item.ID == configContractID() {
			contract = item
			contractFound = true

			break
		}
	}
	if !contractFound {
		panic("API contract ledger missing config contract entry")
	}

	return review, contract
}

func loadLiveConfigEntry(sourceRoot, docsRoot string) manifestEntry {
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
		if entry.Package == configPackage {
			return entry
		}
	}

	panic("live API inventory missing config package")
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

func configReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`config\.[A-Z][A-Za-z0-9_]*(?:\.[A-Z][A-Za-z0-9_]*)?`)
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

func configSourceEvidence(sourceRoot string) []string {
	evidence := []string{
		"config/app.go:4",
		"config/approval.go:9",
		"config/cors.go:4",
		"config/data_sources.go:5",
		"config/env.go:5",
		"config/event.go:6",
		"config/mcp.go:4",
		"config/monitor.go:5",
		"config/redis.go:4",
		"config/security.go:5",
		"config/storage.go:6",
	}
	for _, item := range evidence {
		path, _, _ := strings.Cut(item, ":")
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			panic(err)
		}
	}

	return evidence
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

func configContractID() string {
	return configPackage + "#configuration-contract:config-keys-and-effective-defaults"
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}

	return false
}

func containsNormalized(content, term string) bool {
	return strings.Contains(content, term) ||
		strings.Contains(strings.Join(strings.Fields(content), " "), strings.Join(strings.Fields(term), " "))
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

func writeTempGo(code string) (string, error) {
	file, err := os.CreateTemp("", "verify-config-contracts-*.go")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(code); err != nil {
		return "", err
	}

	return file.Name(), nil
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
