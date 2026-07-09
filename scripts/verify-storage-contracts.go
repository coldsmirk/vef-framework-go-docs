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
	storagePackage = "github.com/coldsmirk/vef-framework-go/storage"

	storageFingerprint = "38d5a444291850b9a3b564e6e885c9e3c9bb516237d850a448d1b8d67ad53faf"
	storageTopLevel    = 101
	storageFields      = 50
	storageMethods     = 31
	storageEntries     = 182

	storageGroupedEntries              = 81
	storageGroupedFields               = 50
	storageGroupedMethods              = 31
	storageGroupedReceivers            = 29
	storageGroupedSignatureFingerprint = "d1a8f855b01c098ff7f221f5548abbe6df79a7b53aaf7536d4845467dad2316b"
	storageGroupedReceiverFingerprint  = "1a12dd3b3db4ae47817d0c6435678a7b6060d44f6cd9756adbfe3b5543d13c8d"

	englishStoragePath    = "docs/infrastructure/storage.md"
	chineseStoragePath    = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/infrastructure/storage.md"
	englishConfigPath     = "docs/reference/configuration-reference.md"
	chineseConfigPath     = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/configuration-reference.md"
	englishBuiltInsPath   = "docs/reference/built-in-resources.md"
	chineseBuiltInsPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/built-in-resources.md"
	englishExtensionsPath = "docs/reference/extension-points.md"
	chineseExtensionsPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/extension-points.md"
	englishIndexPath      = "docs/reference/public-api-index.md"
	chineseIndexPath      = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
)

type corpus struct {
	label   string
	content string
}

type docSet struct {
	featureDocs   []corpus
	configDocs    []corpus
	builtInDocs   []corpus
	extensionDocs []corpus
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
	Disposition    string   `json:"disposition"`
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
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

	checks := []struct {
		sourcePath      string
		docGroup        string
		sourceTerms     []string
		docTerms        []string
		englishDocTerms []string
		chineseDocTerms []string
	}{
		{
			sourcePath: "config/storage.go",
			docGroup:   "config",
			sourceTerms: []string{
				"StorageMinIO      StorageProvider = \"minio\"",
				"StorageMemory     StorageProvider = \"memory\"",
				"StorageFilesystem StorageProvider = \"filesystem\"",
				"AutoMigrate bool", "DefaultMaxUploadSize int64 = 1024 * 1024 * 1024",
				"DefaultClaimTTL         time.Duration = 24 * time.Hour",
				"DefaultMaxPendingClaims int           = 100",
				"DefaultSweepInterval  time.Duration = 5 * time.Minute",
				"DefaultSweepBatchSize int           = 200",
				"DefaultDeleteWorkerInterval time.Duration = 5 * time.Minute",
				"DefaultDeleteBatchSize      int           = 100",
				"DefaultDeleteConcurrency    int           = 8",
				"DefaultDeleteMaxAttempts    int           = 12",
				"DefaultDeleteLeaseWindow    time.Duration = 5 * time.Minute",
				"coalescePositive", "EffectiveMaxUploadSize", "EffectiveClaimTTL",
			},
			docTerms: []string{
				"memory", "filesystem", "minio", "vef.storage.auto_migrate",
				"sys_storage_upload_claim", "sys_storage_upload_part",
				"sys_storage_pending_delete", "1 GiB", "24h", "100",
				"5m", "200", "8", "12", "Effective...",
			},
			englishDocTerms: []string{"logs a warning", "objects are lost on restart"},
			chineseDocTerms: []string{"输出 warning", "对象会在进程重启后丢失"},
		},
		{
			sourcePath: "internal/storage/storage.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"provider == \"\"", "provider = config.StorageMemory",
				"storage provider not configured; defaulting to in-memory storage",
				"ErrUnsupportedStorageProvider", "config.StorageMinIO",
				"config.StorageMemory", "config.StorageFilesystem",
			},
			docTerms: []string{
				"`storage.provider`", "`memory`",
				"`filesystem`", "`minio`",
			},
			englishDocTerms: []string{"`storage.provider` selects the backend", "defaults to `memory` and logs a warning"},
			chineseDocTerms: []string{"`storage.provider` 选择后端", "默认 `memory` 并输出 warning"},
		},
		{
			sourcePath: "storage/service.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"Service interface", "PutObject", "GetObject", "DeleteObject",
				"DeleteObjects", "CopyObject", "StatObject", "Multipart interface",
				"GetObject(ctx context.Context, opts GetObjectOptions) (io.ReadCloser, *ObjectInfo, error)",
				"Callers MUST close the reader and MUST nil-check the ObjectInfo",
				"PartSize() int64", "MaxPartCount() int", "InitMultipart",
				"PutPart", "CompleteMultipart", "AbortMultipart",
				"MultipartFor", "ErrPartTooSmall", "ErrPartETagMismatch",
				"ErrPartNumberOutOfRange", "ErrUploadSessionNotFound",
				"AbortMultipart is idempotent",
			},
			docTerms: []string{
				"type Service interface", "PutObject", "GetObject", "DeleteObject",
				"DeleteObjects", "CopyObject", "StatObject", "type Multipart interface",
				"GetObject(ctx, opts GetObjectOptions) (io.ReadCloser, *ObjectInfo, error)",
				"PartSize() int64", "MaxPartCount() int", "InitMultipart",
				"PutPart", "CompleteMultipart", "AbortMultipart", "MultipartFor",
				"ErrPartTooSmall", "ErrPartETagMismatch", "ErrPartNumberOutOfRange",
				"ErrUploadSessionNotFound", "AbortMultipart",
			},
			englishDocTerms: []string{"last-writer-wins", "idempotent"},
			chineseDocTerms: []string{"last-writer-wins", "幂等"},
		},
		{
			sourcePath: "internal/storage/resource.go",
			docGroup:   "built-in",
			sourceTerms: []string{
				"api.NewRPCResource", "\"sys/storage\"", "Action: \"init_upload\"",
				"Action: \"upload_part\"", "Action: \"list_parts\"",
				"Action: \"complete_upload\"", "Action: \"abort_upload\"",
				"InitUploadParams", "Filename    string", "Size        int64",
				"ContentType string", "Public      bool", "InitUploadResult",
				"ClaimID", "PartSize", "PartCount", "ExpiresAt",
				"Status:           store.ClaimStatusPending",
				"SetUploadID", "params.Public && !r.cfg.AllowPublicUploads",
				"sanitizeContentType", "storage.ErrPublicUploadsNotAllowed",
				"storage.ErrTooManyPendingUploads", "storage.ErrUploadRequiresMultipart",
				"storage.ErrUploadRequiresFile", "storage.ErrUploadPartTooLarge",
				"storage.ErrUploadPartTooSmall", "storage.ErrUploadPartsIncomplete",
				"CompleteUploadResult", "OriginalFilename", "ErrUploadSessionNotFound",
				"storage.ErrUploadSizeMismatch",
			},
			docTerms: []string{
				"`sys/storage`", "`init_upload`", "`upload_part`", "`list_parts`",
				"`complete_upload`", "`abort_upload`", "`claimId`", "`partSize`",
				"`partCount`", "`expiresAt`", "`originalFilename`",
				"`public`", "`vef.storage.allow_public_uploads`",
				"`contentType`", "`application/octet-stream`", "`partNumber`",
				"`file`", "`metadata`", "ErrCodeUploadSizeMismatch",
			},
			englishDocTerms: []string{
				"init_upload -> upload_part ->\ncomplete_upload",
				"create a pending claim", "Small files still return `partCount = 1`",
				"no\nsingle-PUT HTTP action", "unsafe same-origin types",
			},
			chineseDocTerms: []string{
				"init_upload -> upload_part ->\ncomplete_upload",
				"创建 pending claim", "小文件也会返回 `partCount = 1`",
				"不存在单次 PUT HTTP\naction", "同源不安全类型",
			},
		},
		{
			sourcePath: "internal/storage/proxy_middleware.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"return \"storage_proxy\"", "return 900",
				"router.Get(\"/storage/files/+\"", "url.PathUnescape",
				"isValidObjectKey", "storage.PublicPrefix", "FileACL.CanRead",
				"X-Content-Type-Options", "nosniff", "public, max-age=3600, immutable",
				"private, no-store", "Do NOT send ETag for private files",
			},
			docTerms: []string{
				"GET /storage/files/<key>", "storage_proxy", "900",
				"`pub/*`", "`FileACL.CanRead`",
				"X-Content-Type-Options: nosniff",
				"Cache-Control: public, max-age=3600, immutable",
				"Cache-Control: private, no-store",
			},
			englishDocTerms: []string{"URL-decodes", "no `ETag`"},
			chineseDocTerms: []string{"URL 解码", "不返回 `ETag`"},
		},
		{
			sourcePath: "storage/file_acl.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"PublicPrefix = \"pub/\"", "PrivatePrefix = \"priv/\"",
				"FileACL interface", "CanRead", "DefaultFileACL",
				"strings.HasPrefix(key, PublicPrefix)",
				"vef.SupplyFileACL",
			},
			docTerms: []string{
				"storage.PublicPrefix", "`pub/`", "storage.PrivatePrefix",
				"`priv/`", "storage.FileACL", "CanRead", "storage.DefaultFileACL",
				"vef.SupplyFileACL(...)",
			},
			englishDocTerms: []string{"default denies every `priv/*` read"},
			chineseDocTerms: []string{"默认实现拒绝所有 `priv/*` 读"},
		},
		{
			sourcePath: "storage/file_refs.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"MetaTypeUploadedFile MetaType = \"uploaded_file\"",
				"MetaTypeRichText MetaType = \"rich_text\"",
				"MetaTypeMarkdown MetaType = \"markdown\"",
				"const tagMeta = \"meta\"", "strx.WithBareValueMode(strx.BareAsKey)",
				"WithDiveTag(tagMeta, \"dive\")", "reflectx.IsStringSliceType",
				"reflectx.IsStringMapType", "map[string]string",
			},
			docTerms: []string{
				"meta:\"uploaded_file\"", "meta:\"rich_text\"", "meta:\"markdown\"",
				"meta:\"dive\"", "string` / `*string` / `[]string` / `map[string]string`",
				"map", "value",
			},
			englishDocTerms: []string{"Unsupported field shapes are ignored"},
			chineseDocTerms: []string{"不支持的字段形态会被忽略"},
		},
		{
			sourcePath: "storage/url_key_mapper.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"DefaultProxyPrefix = \"/storage/files/\"",
				"ProxyURLKeyMapper", "IdentityURLKeyMapper",
				"URLToKey", "KeyToURL", "parsed.Scheme != \"\"",
				"return m.prefix() + key",
			},
			docTerms: []string{
				"URLKeyMapper", "storage.ProxyURLKeyMapper",
				"storage.IdentityURLKeyMapper", "storage.DefaultProxyPrefix",
				"URLToKey", "KeyToURL", "/storage/files/<key>",
				"ReplaceHtmlURLs", "ReplaceMarkdownURLs",
			},
		},
		{
			sourcePath: "storage/metadata.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"func CanonicalizeMetadataKeys(m map[string]string) map[string]string",
				"textproto.CanonicalMIMEHeaderKey",
				"if len(m) == 0",
				"out[textproto.CanonicalMIMEHeaderKey(k)] = v",
			},
			docTerms: []string{
				"CanonicalizeMetadataKeys",
				"S3/HTTP-header canonical form",
				"provider-neutral",
				"author",
			},
		},
		{
			sourcePath: "storage/events.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"EventTypeFileClaimed = \"vef.storage.file.claimed\"",
				"EventTypeFileDeleted = \"vef.storage.file.deleted\"",
				"EventTypeDeleteDeadLetter = \"vef.storage.delete.dead_letter\"",
				"FileClaimedEvent", "FileDeletedEvent", "DeleteDeadLetterEvent",
				"PendingDeleteID", "LastError", "EventType() string",
			},
			docTerms: []string{
				"vef.storage.file.claimed", "vef.storage.file.deleted",
				"vef.storage.delete.dead_letter", "FileClaimedEvent",
				"FileDeletedEvent", "DeleteDeadLetterEvent", "pendingDeleteId",
				"lastError", "event.WithTx", "event.WithGroup",
			},
		},
		{
			sourcePath: "internal/storage/module.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"ErrEventRouteNotTransactional", "verifyEventRouting",
				"storage.EventTypeFileClaimed", "storage.EventTypeFileDeleted",
				"storage.EventTypeDeleteDeadLetter", "inspector.HasTransactionalRoute",
				"vef.event.transports.outbox.enabled=true", "vef.storage.*",
				"vef.event.default_transport=\\\"outbox\\\"",
			},
			docTerms: []string{
				"vef.storage.file.claimed", "vef.storage.file.deleted",
				"vef.storage.delete.dead_letter", "vef.storage.*", "outbox",
			},
			englishDocTerms: []string{"transactional event transport", "fails fast at startup"},
			chineseDocTerms: []string{"事务性 event transport", "启动失败"},
		},
		{
			sourcePath: "internal/storage/module.go",
			docGroup:   "extension",
			sourceTerms: []string{
				"newDefaultFileACL", "newDefaultURLKeyMapper",
				"vef:api:resources", "vef:app:middlewares",
			},
			docTerms: []string{
				"vef.SupplyFileACL", "vef.SupplyURLKeyMapper",
				"vef:api:resources", "vef:app:middlewares",
			},
		},
		{
			sourcePath: "internal/storage/worker/claim_sweeper.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"ClaimSweeper", "ListExpired", "MarkUploadedIfPendingExpired",
				"DeleteIfPendingExpired", "DeleteReasonClaimExpired",
				"canRecoverCompletedClaim", "StatObject", "info.Size != claim.Size",
			},
			docTerms: []string{
				"claim sweeper", "DeleteReasonClaimExpired",
			},
			englishDocTerms: []string{"expired-but-completed multipart object", "marking it uploaded"},
			chineseDocTerms: []string{"已过期但对象已经完成的 multipart\nclaim 恢复为 uploaded"},
		},
		{
			sourcePath: "internal/storage/worker/delete_worker.go",
			docGroup:   "feature",
			sourceTerms: []string{
				"DeleteWorker", "Lease", "EffectiveDeleteConcurrency",
				"DeleteObject", "item.Reason", "computeBackoff",
				"30 * time.Second", "1 * time.Hour", "DeleteMaxAttempts",
				"NewFileDeletedEvent", "NewDeleteDeadLetterEvent",
				"classifyDeleteError", "access_denied", "bucket_not_found",
				"session_not_found", "transient",
			},
			docTerms: []string{
				"DeleteWorker", "vef.storage.file.deleted",
				"vef.storage.delete.dead_letter", "DeleteReason", "access_denied",
				"bucket_not_found", "session_not_found", "transient",
			},
			englishDocTerms: []string{"retry/backoff", "queue row is already retired"},
			chineseDocTerms: []string{"队列行已经退役"},
		},
		{
			sourcePath: "internal/storage/migration/migration.go",
			docGroup:   "config",
			sourceTerms: []string{
				"expectedTables", "sys_storage_upload_claim",
				"sys_storage_upload_part", "sys_storage_pending_delete",
				"Label:          \"storage\"", "Idempotent",
			},
			docTerms: []string{
				"sys_storage_upload_claim", "sys_storage_upload_part",
				"sys_storage_pending_delete",
			},
			englishDocTerms: []string{"idempotent"},
			chineseDocTerms: []string{"幂等"},
		},
	}

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)
	manifestPath := filepath.Join(docsRoot, "scripts/api-audit-manifest.json")
	auditLedgerPath := filepath.Join(docsRoot, "scripts/api-audit-ledger.json")
	contractLedgerPath := filepath.Join(docsRoot, "scripts/api-contract-ledger.json")

	englishDocs := readCorpus("English storage docs", filepath.Join(docsRoot, englishStoragePath))
	chineseDocs := readCorpus("Chinese storage docs", filepath.Join(docsRoot, chineseStoragePath))
	englishConfig := readCorpus("English configuration docs", filepath.Join(docsRoot, englishConfigPath))
	chineseConfig := readCorpus("Chinese configuration docs", filepath.Join(docsRoot, chineseConfigPath))
	englishBuiltIns := readCorpus("English built-in resources docs", filepath.Join(docsRoot, englishBuiltInsPath))
	chineseBuiltIns := readCorpus("Chinese built-in resources docs", filepath.Join(docsRoot, chineseBuiltInsPath))
	englishExtensions := readCorpus("English extension points docs", filepath.Join(docsRoot, englishExtensionsPath))
	chineseExtensions := readCorpus("Chinese extension points docs", filepath.Join(docsRoot, chineseExtensionsPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	runtime := loadJSON[runtimeLedger](filepath.Join(docsRoot, "scripts/runtime-api-ledger.json"))
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[storagePackage]
	liveAuditEntries := storageEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	storageEntries := storageEntriesFromAudit(audit)

	docs := docSet{
		featureDocs:   []corpus{englishDocs, chineseDocs},
		configDocs:    []corpus{englishConfig, chineseConfig},
		builtInDocs:   []corpus{englishBuiltIns, chineseBuiltIns},
		extensionDocs: []corpus{englishExtensions, chineseExtensions},
	}

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(storageEntries)...)
	failures = append(failures, verifyLiveAuditEntries(storageEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedStorageSurface(storageEntries)...)
	failures = append(failures, verifyRuntimeActions(runtime, []corpus{englishDocs, chineseDocs, englishBuiltIns, chineseBuiltIns})...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, storageEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, storageEntries)...)
	for _, check := range checks {
		source := readCorpus(check.sourcePath, filepath.Join(sourceRoot, check.sourcePath))
		failures = append(failures, missingTerms(source, check.sourceTerms)...)
		for _, doc := range docsForCheck(check.sourcePath, check.docGroup, docs) {
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

	fmt.Printf(
		"Storage contract docs verified: %d public entries, %d grouped field/method entries, sys/storage runtime actions, %d source files, %d doc corpora\n",
		len(storageEntries),
		storageGroupedEntries,
		len(checks),
		8,
	)
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != storagePackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != storageTopLevel ||
		entry.Fields != storageFields ||
		entry.Methods != storageMethods ||
		entry.Fingerprint != storageFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s storage surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			storageTopLevel, storageFields, storageMethods, storageFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != storagePackage {
			continue
		}

		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, storageCoverage()) {
			failures = append(failures, fmt.Sprintf("storage manifest coverage mismatch: got %v want %v", entry.Coverage, storageCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing storage package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != storagePackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "storage contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != storageTopLevel ||
			review.ReviewedSurface.Fields != storageFields ||
			review.ReviewedSurface.Methods != storageMethods ||
			review.ReviewedSurface.EntryCount != storageEntries ||
			review.ReviewedSurface.Fingerprint != storageFingerprint {
			failures = append(failures, fmt.Sprintf(
				"storage contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, storageCoverage()) {
			failures = append(failures, fmt.Sprintf("storage contract review coverage mismatch: got %v want %v", review.Coverage, storageCoverage()))
		}
		if !contains(review.ContractIDs, storageContractID()) {
			failures = append(failures, "storage contract review missing contract id "+storageContractID())
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing storage package review")
	}

	var foundContract bool
	for _, entry := range contracts.Entries {
		if entry.ID != storageContractID() {
			continue
		}
		foundContract = true
		if entry.Package != storagePackage || entry.Kind != "dynamic-resource" {
			failures = append(failures, fmt.Sprintf("storage contract entry shape mismatch: package=%s kind=%s", entry.Package, entry.Kind))
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "storage contract entry disposition mismatch: "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, storageCoverage()) {
			failures = append(failures, fmt.Sprintf("storage contract coverage mismatch: got %v want %v", entry.Coverage, storageCoverage()))
		}
		for _, term := range []string{
			"sys/storage",
			"/storage/files/",
			"pub/",
			"priv/",
			`meta:"uploaded_file"`,
			`meta:"dive"`,
			"vef.storage.delete.dead_letter",
		} {
			if !contains(entry.Terms, term) {
				failures = append(failures, "storage contract missing term "+term)
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
	}
	if !foundContract {
		failures = append(failures, "contract ledger missing storage contract entry")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != storageEntries {
		failures = append(failures, fmt.Sprintf("storage audit entry count mismatch: got %d want %d", len(entries), storageEntries))
	}

	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != storagePackage {
			failures = append(failures, "non-storage audit entry passed into storage verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate storage audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "storage audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, storageCoverage()) {
			failures = append(failures, fmt.Sprintf("storage audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, storageCoverage()))
		}
	}
	if counts["top"] != storageTopLevel || counts["field"] != storageFields || counts["method"] != storageMethods {
		failures = append(failures, fmt.Sprintf("storage audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != storageTopLevel ||
		dispositionCounts["grouped:type-member-family"] != storageGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"storage audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			storageTopLevel,
			storageGroupedEntries,
		))
	}
	if dispositionCounts["documented:top-level"]+dispositionCounts["grouped:type-member-family"] != storageEntries {
		failures = append(failures, "storage audit has entries outside the reviewed dispositions")
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
			failures = append(failures, fmt.Sprintf("storage missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"storage live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
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
			failures = append(failures, fmt.Sprintf("storage extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedStorageSurface(entries []auditEntry) []string {
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
			failures = append(failures, fmt.Sprintf("storage grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("storage grouped type-member surface", rows, storageGroupedEntries, storageGroupedSignatureFingerprint)...)
	if kindCounts["field"] != storageGroupedFields || kindCounts["method"] != storageGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"storage grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			storageGroupedFields,
			storageGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("storage grouped receiver/type families", receiverRows, storageGroupedReceivers, storageGroupedReceiverFingerprint)...)

	return failures
}

func verifyRuntimeActions(runtime runtimeLedger, docs []corpus) []string {
	actions := map[string]bool{}
	var resourceFound bool
	for _, entry := range runtime.Entries {
		if entry.Category == "built-in resource" && entry.Name == "sys/storage" && entry.Value == "rpc" {
			resourceFound = true
		}
		if entry.Category == "built-in resource action" && strings.HasPrefix(entry.Name, "sys/storage/") {
			actions[strings.TrimPrefix(entry.Name, "sys/storage/")] = true
		}
	}

	wantActions := []string{
		"abort_upload",
		"complete_upload",
		"init_upload",
		"list_parts",
		"upload_part",
	}
	var failures []string
	if !resourceFound {
		failures = append(failures, "runtime ledger missing sys/storage built-in resource")
	}
	failures = append(failures, compareSets("runtime sys/storage actions", sortedKeys(actions), wantActions)...)
	for _, doc := range docs {
		if !strings.Contains(doc.content, "`sys/storage`") {
			failures = append(failures, doc.label+" missing sys/storage")
		}
		for _, action := range wantActions {
			if !strings.Contains(doc.content, "`"+action+"`") {
				failures = append(failures, doc.label+" missing storage action "+action)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, storagePackage)
	if section == "" {
		return []string{index.label + " missing storage package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s storage index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func runGoTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./storage", "./internal/storage/...")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./storage ./internal/storage/... failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func storageEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == storagePackage {
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

func storageContractID() string {
	return storagePackage + "#dynamic-resource:storage-resource-upload-proxy-lifecycle"
}

func storageCoverage() []string {
	return []string{englishStoragePath, englishBuiltInsPath, englishExtensionsPath}
}

func docsForCheck(sourcePath, docGroup string, docs docSet) []corpus {
	switch docGroup {
	case "feature":
		return docs.featureDocs
	case "config":
		return docs.configDocs
	case "built-in":
		return docs.builtInDocs
	case "extension":
		return docs.extensionDocs
	default:
		panic(fmt.Sprintf("storage check for %s has unknown docGroup %q", sourcePath, docGroup))
	}
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
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

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
