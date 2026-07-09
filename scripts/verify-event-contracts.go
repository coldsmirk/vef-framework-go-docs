package main

import (
	"bytes"
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
	eventDocsPath         = "docs/infrastructure/event-bus.md"
	chineseEventDocsPath  = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/infrastructure/event-bus.md"
	englishIndexPath      = "docs/reference/public-api-index.md"
	chineseIndexPath      = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
	eventContractCoverage = eventDocsPath
)

type corpus struct {
	label   string
	content string
}

type eventPackage struct {
	pkg         string
	topLevel    int
	fields      int
	methods     int
	entries     int
	fingerprint string
	contractID  string
}

var eventPackages = []eventPackage{
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/event",
		topLevel:    44,
		fields:      28,
		methods:     10,
		entries:     82,
		fingerprint: "490bc23b8a2ca2918b6114a3ea4ba3263c33e3ae1258fb4e784474853ceb0a58",
		contractID:  "github.com/coldsmirk/vef-framework-go/event#runtime-contract:event-routing-publish-subscribe",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/event/inbox",
		topLevel:    13,
		fields:      9,
		methods:     4,
		entries:     26,
		fingerprint: "ecaf2505c06ba2fd7420ff3912d385a68a1d6acd2a96983add3d9885e3426835",
		contractID:  "github.com/coldsmirk/vef-framework-go/event/inbox#runtime-contract:inbox-dedupe-record-status",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/event/middleware",
		topLevel:    15,
		fields:      0,
		methods:     7,
		entries:     22,
		fingerprint: "30e45d70410a8b34acd03d590c3e2e6b9ee7dc6521bd713faed6b7859bf28946",
		contractID:  "github.com/coldsmirk/vef-framework-go/event/middleware#runtime-contract:event-middleware-order-and-trace-context",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/event/transport",
		topLevel:    10,
		fields:      18,
		methods:     17,
		entries:     45,
		fingerprint: "b7d275442c0cf34f0210bf9d1a9b3b7229c10cf9ebbf446fac80afad725b53af",
		contractID:  "github.com/coldsmirk/vef-framework-go/event/transport#runtime-contract:transport-frame-capabilities",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/event/transport/memory",
		topLevel:    6,
		fields:      3,
		methods:     2,
		entries:     11,
		fingerprint: "81cbba959853c34db2310fa2a9933d0e45eee7a0dc8cc0b7e458f71e97a1e238",
		contractID:  "github.com/coldsmirk/vef-framework-go/event/transport/memory#runtime-contract:memory-transport-name-and-full-policy",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/event/transport/outbox",
		topLevel:    10,
		fields:      23,
		methods:     12,
		entries:     45,
		fingerprint: "76ed794b65fe67f487f225c129d2802be55fbb448212a1a457b9740298ade2cc",
		contractID:  "github.com/coldsmirk/vef-framework-go/event/transport/outbox#runtime-contract:outbox-transport-status-and-sink",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/event/transport/redisstream",
		topLevel:    2,
		fields:      13,
		methods:     10,
		entries:     25,
		fingerprint: "af99f02ce92f12b046e9fe3438f3f87878aebf2fda33d5a37043089761f1a20f",
		contractID:  "github.com/coldsmirk/vef-framework-go/event/transport/redisstream#runtime-contract:redis-stream-transport-config",
	},
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

	checks := []struct {
		sourcePath      string
		sourceTerms     []string
		docTerms        []string
		englishDocTerms []string
		chineseDocTerms []string
	}{
		{
			sourcePath: "event/event.go",
			sourceTerms: []string{
				"Event interface", "EventType() string", "Envelope struct",
				"ID string", "Type string", "Source string", "OccurredAt time.Time",
				"PublishedAt time.Time", "TraceID string", "SpanID  string",
				"CorrelationID string", "Headers map[string]string", "Payload Event",
				"RawPayload struct", "Body []byte", "Handler func",
				"ErrorSink func", "Unsubscribe func", "Bus interface",
				"Publish(ctx context.Context, evt Event, opts ...PublishOption)",
				"PublishBatch(ctx context.Context, evts []Event, opts ...PublishOption)",
				"Subscribe(eventType string, h Handler, opts ...SubscribeOption)",
			},
			docTerms: []string{
				"type Bus interface", "Publish(ctx context.Context, evt Event, opts ...PublishOption)",
				"PublishBatch(ctx context.Context, evts []Event, opts ...PublishOption)",
				"Subscribe(eventType string, h Handler, opts ...SubscribeOption)",
				"EventType()", "Envelope", "RawPayload", "Handler",
				"ErrorSink", "Unsubscribe", "ID", "Type", "Source",
				"OccurredAt", "PublishedAt", "TraceID", "SpanID",
				"CorrelationID", "Headers", "Payload", "Body",
			},
		},
		{
			sourcePath: "event/publish.go",
			sourceTerms: []string{
				"PublishOption func(*PublishConfig)", "PublishConfig struct",
				"Tx orm.DB", "Async bool", "Source string", "OccurredAt time.Time",
				"CorrelationID string", "Headers map[string]string",
				"ApplyPublishOptions", "WithTx", "WithAsync", "WithSource",
				"WithOccurredAt", "WithCorrelationID", "WithHeaders",
				"maps.Copy(c.Headers, h)",
			},
			docTerms: []string{
				"PublishOption", "PublishConfig", "event.WithTx(tx orm.DB)",
				"event.WithAsync()", "event.WithSource(name)",
				"event.WithOccurredAt(t)", "event.WithCorrelationID(id)",
				"event.WithHeaders(map)", "ApplyPublishOptions",
			},
			englishDocTerms: []string{"left-to-right", "later options win"},
			chineseDocTerms: []string{"从左到右", "后面的值覆盖"},
		},
		{
			sourcePath: "event/subscribe.go",
			sourceTerms: []string{
				"SubscribeOption func(*SubscribeConfig)", "SubscribeConfig struct",
				"Group string", "Concurrency int", "ApplySubscribeOptions",
				"WithGroup", "WithConcurrency", "if n > 0",
			},
			docTerms: []string{
				"SubscribeOption", "SubscribeConfig", "event.WithGroup(name)",
				"event.WithConcurrency(n)", "ApplySubscribeOptions",
			},
			englishDocTerms: []string{"worker count", "Defaults to 1"},
			chineseDocTerms: []string{"worker 数量", "默认 1"},
		},
		{
			sourcePath: "event/typed.go",
			sourceTerms: []string{
				"TypedHandler[T Event]", "SubscribeTyped[T Event]", "ErrNilTypeParameter",
				"decodePayload[T]", "RawPayload", "json.Unmarshal(raw.Body",
				"AsEvents[T Event]", "ErrUnknownPayload",
			},
			docTerms: []string{
				"SubscribeTyped[T]", "TypedHandler", "RawPayload",
				"canonical JSON body", "ErrUnknownPayload", "ErrNilTypeParameter",
				"AsEvents",
			},
			englishDocTerms: []string{"pointer type", "value type"},
			chineseDocTerms: []string{"指针类型", "值类型"},
		},
		{
			sourcePath: "event/errors.go",
			sourceTerms: []string{
				"ErrBusNotStarted", "ErrBusAlreadyStarted", "ErrTxRequired",
				"ErrTransportNotFound", "ErrAsyncQueueFull", "ErrQueueFull",
				"ErrHandlerPanic", "ErrShutdownTimeout", "ErrNoRouteMatched",
				"ErrUnknownPayload", "ErrPayloadTooLarge", "ErrInvalidEventType",
				"ErrNilTypeParameter", "ErrGroupRequired", "ErrTxAsyncMutex",
			},
			docTerms: []string{
				"event.ErrBusNotStarted", "event.ErrBusAlreadyStarted",
				"event.ErrTxRequired", "event.ErrTransportNotFound",
				"event.ErrAsyncQueueFull", "event.ErrQueueFull",
				"event.ErrHandlerPanic", "event.ErrShutdownTimeout",
				"event.ErrNoRouteMatched", "event.ErrUnknownPayload",
				"event.ErrPayloadTooLarge", "event.ErrInvalidEventType",
				"event.ErrNilTypeParameter", "event.ErrGroupRequired",
				"event.ErrTxAsyncMutex",
			},
		},
		{
			sourcePath: "event/metrics.go",
			sourceTerms: []string{
				"MetricsRecorder interface", "PublishObserved(eventType string, err error)",
				"ConsumeObserved(eventType string, elapsed time.Duration, err error)",
			},
			docTerms: []string{
				"MetricsRecorder", "PublishObserved", "ConsumeObserved",
			},
			englishDocTerms: []string{"counters and latency histograms"},
			chineseDocTerms: []string{"计数和延迟直方图"},
		},
		{
			sourcePath: "event/route_inspector.go",
			sourceTerms: []string{
				"RouteInspector interface", "HasTransactionalRoute(eventType string) bool",
				"HasSubscribableTransport(eventType string) bool",
			},
			docTerms: []string{
				"type RouteInspector interface", "HasTransactionalRoute(eventType string) bool",
				"HasSubscribableTransport(eventType string) bool", "ErrTxRequired",
				"ErrNoRouteMatched",
			},
			englishDocTerms: []string{"fail fast"},
			chineseDocTerms: []string{"快速失败"},
		},
		{
			sourcePath: "event/middleware/middleware.go",
			sourceTerms: []string{
				"OrderRecover = -100", "OrderTracing = -50", "OrderLogging = -25",
				"OrderMetrics = 0", "OrderInbox = 100", "PublishHandler",
				"PublishMiddleware interface", "ConsumeHandler", "ConsumeMiddleware interface",
				"Applies(caps transport.Capabilities) bool", "ChainPublish",
				"ChainConsume", "slices.SortStableFunc",
			},
			docTerms: []string{
				"PublishHandler", "ConsumeHandler", "PublishMiddleware",
				"ConsumeMiddleware", "ChainPublish", "ChainConsume",
				"OrderLogging", "OrderTracing", "OrderMetrics", "OrderRecover",
				"OrderInbox",
			},
			englishDocTerms: []string{"registration order"},
			chineseDocTerms: []string{"注册顺序"},
		},
		{
			sourcePath: "event/middleware/tracing_context.go",
			sourceTerms: []string{
				"WithTraceID", "WithIncomingTraceID", "TraceIDFromContext",
				"IncomingTraceIDFromContext",
			},
			docTerms: []string{
				"TraceIDFromContext", "IncomingTraceIDFromContext",
				"WithTraceID", "WithIncomingTraceID", "tracing_strict",
			},
		},
		{
			sourcePath: "event/transport/transport.go",
			sourceTerms: []string{
				"ErrSubscribeUnsupported", "EventTypePattern", "Frame struct",
				"Delivery interface", "Frame() Frame", "Attempt() int",
				"Ack(ctx context.Context)", "Nack(ctx context.Context, retryAfter time.Duration, err error)",
				"ConsumeFunc", "SubscribeConfig", "Capabilities struct",
				"Durable bool", "Transactional bool", "Ordered bool",
				"AtLeastOnce bool", "SupportsGroups bool", "PublishOnly bool",
				"Transport interface", "Name() string", "Capabilities() Capabilities",
				"Start(ctx context.Context)", "Stop(ctx context.Context)",
				"Publish(ctx context.Context, frames []Frame)",
				"Subscribe(eventType, group string", "TxTransport interface",
				"PublishTx(ctx context.Context, tx orm.DB, frames []Frame)",
			},
			docTerms: []string{
				"ErrSubscribeUnsupported", "EventTypePattern", "Frame",
				"Delivery", "Capabilities", "Durable", "Transactional",
				"Ordered", "AtLeastOnce", "SupportsGroups", "PublishOnly",
				"Transport", "TxTransport", "ConsumeFunc", "Unsubscribe",
			},
		},
		{
			sourcePath: "event/inbox/inbox.go",
			sourceTerms: []string{
				"StatusProcessing Status = \"processing\"", "StatusCompleted Status = \"completed\"",
				"AcquireResultAcquired AcquireResult = \"acquired\"",
				"AcquireResultCompleted AcquireResult = \"completed\"",
				"AcquireResultInProgress AcquireResult = \"in_progress\"",
				"Record struct", "EventID string", "ConsumerGroup string",
				"LockID string", "LockedUntil", "CompletedAt",
				"Repository interface", "Acquire(ctx context.Context",
				"MarkCompleted", "Release", "DeleteOlderThan",
			},
			docTerms: []string{
				"StatusProcessing", "StatusCompleted", "AcquireResultAcquired",
				"AcquireResultCompleted", "AcquireResultInProgress",
				"processing", "completed", "acquired", "in_progress",
				"eventId", "consumerGroup", "status", "lockId",
				"lockedUntil", "completedAt", "Repository",
			},
		},
		{
			sourcePath: "event/inbox/errors.go",
			sourceTerms: []string{
				"ErrInProgress", "ErrLockLost", "ErrUnknownAcquireResult",
				"ErrMissingLockID",
			},
			docTerms: []string{
				"ErrInProgress", "ErrLockLost", "ErrMissingLockID",
				"ErrUnknownAcquireResult",
			},
		},
		{
			sourcePath: "event/transport/memory/memory.go",
			sourceTerms: []string{
				"Name = \"memory\"", "FullPolicyError FullPolicy = \"error\"",
				"FullPolicyBlock FullPolicy = \"block\"",
				"FullPolicyDropOldest FullPolicy = \"drop_oldest\"",
				"QueueSize int", "FullPolicy FullPolicy", "PublishTimeout time.Duration",
				"EffectiveQueueSize", "return 1024", "EffectiveFullPolicy",
			},
			docTerms: []string{
				"`memory`", "FullPolicyError", "FullPolicyBlock",
				"FullPolicyDropOldest", "`error`", "`block`", "`drop_oldest`",
				"QueueSize", "PublishTimeout", "1024",
			},
		},
		{
			sourcePath: "event/transport/outbox/outbox.go",
			sourceTerms: []string{
				"Name = \"outbox\"", "StatusPending Status = \"pending\"",
				"StatusProcessing Status = \"processing\"", "StatusCompleted Status = \"completed\"",
				"StatusFailed Status = \"failed\"", "StatusDead Status = \"dead\"",
				"Record struct", "EventID", "EventType", "CorrelationID",
				"Payload", "json.RawMessage", "RetryCount", "LastError",
				"ProcessedAt", "RetryAfter", "OccurredAt", "Repository interface",
				"InsertBatch", "InsertBatchTx", "ClaimBatch", "MarkCompleted",
				"MarkFailed", "DeleteCompletedOlderThan", "RelayInterval",
				"MaxRetries", "BatchSize", "LeaseMultiplier", "MinLease",
				"SinkName", "EffectiveSinkName", "return \"memory\"",
			},
			docTerms: []string{
				"`outbox`", "StatusPending", "StatusProcessing", "StatusCompleted",
				"StatusFailed", "StatusDead", "`pending`", "`processing`",
				"`completed`", "`failed`", "`dead`", "eventId",
				"eventType", "correlationId", "payload", "retryCount",
				"lastError", "processedAt", "retryAfter", "occurredAt",
				"RelayInterval", "MaxRetries", "BatchSize", "LeaseMultiplier",
				"MinLease", "SinkName",
			},
			englishDocTerms: []string{"defaults to `memory`"},
			chineseDocTerms: []string{"默认使用 `memory`"},
		},
		{
			sourcePath: "event/transport/redisstream/redis_stream.go",
			sourceTerms: []string{
				"Name = \"redis_stream\"", "StreamPrefix string", "MaxLenApprox int64",
				"BlockTimeout time.Duration", "ClaimIdle time.Duration",
				"ClaimInterval time.Duration", "ClaimBatchSize int64",
				"ReaperConcurrency int", "HandlerTimeout time.Duration",
				"SetupTimeout time.Duration", "ConsumerID string", "StartID string", "EffectiveStreamPrefix",
				"return \"vef:events:\"", "EffectiveBlockTimeout", "5 * time.Second",
				"EffectiveClaimIdle", "60 * time.Second", "EffectiveClaimInterval",
				"30 * time.Second", "EffectiveClaimBatchSize", "return 64",
				"EffectiveReaperConcurrency", "return 4", "EffectiveSetupTimeout",
				"EffectiveStartID", "return \"0\"", "StreamKey",
			},
			docTerms: []string{
				"`redis_stream`", "StreamPrefix", "MaxLenApprox", "BlockTimeout",
				"ClaimIdle", "ClaimInterval", "ClaimBatchSize", "ReaperConcurrency",
				"HandlerTimeout", "SetupTimeout", "ConsumerID", "StartID",
				"stream_prefix", "max_len_approx", "block_timeout", "claim_idle",
				"claim_interval", "claim_batch_size", "reaper_concurrency",
				"handler_timeout", "setup_timeout", "consumer_id", "start_id",
				"vef:events:", "`0`",
			},
		},
		{
			sourcePath: "internal/event/bus.go",
			sourceTerms: []string{
				"maxFrameBodyBytes = 1 << 20", "maxHeaderEntries    = 32",
				"maxHeaderKeyBytes   = 128", "maxHeaderValueBytes = 1024",
				"ErrBusAlreadyStarted", "ErrBusNotStarted", "ErrTxAsyncMutex",
				"ErrTxRequired", "ErrGroupRequired", "ErrNoRouteMatched",
				"ErrPayloadTooLarge", "ErrInvalidEventType", "ErrQueueFull",
				"context.WithoutCancel(ctx)", "WithAsync",
			},
			docTerms: []string{
				"1 MiB", "32", "128 bytes", "1024 bytes",
				"ErrBusAlreadyStarted", "ErrBusNotStarted", "ErrTxAsyncMutex",
				"ErrTxRequired", "ErrGroupRequired", "ErrNoRouteMatched",
				"ErrPayloadTooLarge", "ErrInvalidEventType", "ErrQueueFull",
			},
			englishDocTerms: []string{"cancellation-detached context"},
			chineseDocTerms: []string{"脱离请求取消的 context"},
		},
		{
			sourcePath: "internal/event/router.go",
			sourceTerms: []string{
				"path.Match", "first matching rule wins", "EffectiveDefaultTransport",
				"event.ErrTransportNotFound",
			},
			docTerms: []string{
				"path.Match", "default_transport",
			},
			englishDocTerms: []string{"first matching rule wins", "unknown transport"},
			chineseDocTerms: []string{"第一条匹配命中即停止", "未知 transport"},
		},
		{
			sourcePath: "internal/event/outbox_module.go",
			sourceTerms: []string{
				"ErrOutboxSinkRouteMismatch", "sink missing from outbox-bearing route",
				"validateOutboxSinkRoute", "outbox.sink", "PublishOnly",
				"subscribable", "vef:event:outbox:relay", "vef:event:outbox:cleanup",
			},
			docTerms: []string{
				"outbox `sink`", "vef:event:outbox:relay", "vef:event:outbox:cleanup",
			},
			englishDocTerms: []string{"same transport list", "validates this at\nstartup"},
			chineseDocTerms: []string{"同一个 transports 列表", "启动时校验"},
		},
		{
			sourcePath: "internal/event/transport/outbox/transport.go",
			sourceTerms: []string{
				"ErrSinkNotConfigured", "ErrDLQReentry", "ErrInvalidFrameBody",
				"PublishOnly:    true", "SetSink", "Subscribe is unsupported",
				"transport.ErrSubscribeUnsupported", "json.Valid(body)",
			},
			docTerms: []string{
				"publish-only", "sink", "ErrSubscribeUnsupported",
				"vef.dlq", "loop guard",
			},
			englishDocTerms: []string{"JSON-shaped frame bodies"},
			chineseDocTerms: []string{"JSON 形状的 frame body"},
		},
		{
			sourcePath: "internal/event/transport/outbox/relay.go",
			sourceTerms: []string{
				"maxLastErrorBytes = 256", "dlqHeader = \"vef.dlq\"",
				"defaultDLQTopic", "\"vef-dlq.\" + eventType", "backoffFor",
				"maxBackoff = time.Hour", "redactError", "StatusDead",
				"DLQ forward", "MarkFailed",
			},
			docTerms: []string{
				"vef-dlq.<eventType>", "vef.dlq=1", "dead", "lastError",
				"256 bytes",
			},
			englishDocTerms: []string{"exponential backoff", "capped at 1h", "DLQ forwarding"},
			chineseDocTerms: []string{"指数退避", "最高 1h", "DLQ 转发"},
		},
		{
			sourcePath: "internal/event/transport/redisstream/transport.go",
			sourceTerms: []string{
				"maxFrameBytes = 1 << 20", "Ping(ctx)", "XAddArgs",
				"MaxLenApprox", "XGroupCreateMkStream", "EffectiveStartID",
				"consumer := prefix + \"-\" + id.GenerateUUID()",
				"Poison-message policy", "XAck", "invalid event type",
			},
			docTerms: []string{
				"PING", "XADD MAXLEN ~", "XGROUP", "start_id",
				"consumer_id", "poison message", "XACK", "1 MiB",
			},
			englishDocTerms: []string{"UUID suffix"},
			chineseDocTerms: []string{"UUID 后缀"},
		},
		{
			sourcePath: "internal/event/transport/redisstream/reaper.go",
			sourceTerms: []string{
				"XPendingExt", "EffectiveClaimIdle", "EffectiveClaimBatchSize",
				"XClaim", "RetryCount", "attempt := max",
			},
			docTerms: []string{
				"reaper", "XCLAIM", "claim_idle", "claim_interval",
				"claim_batch_size", "pending",
			},
		},
		{
			sourcePath: "internal/event/middleware/inbox.go",
			sourceTerms: []string{
				"Capabilities", "caps.AtLeastOnce", "consumerGroupFromContext",
				"AcquireResultAcquired", "AcquireResultCompleted",
				"AcquireResultInProgress", "ErrInProgress", "ErrMissingLockID",
				"ErrUnknownAcquireResult", "ErrLockLost", "Release",
			},
			docTerms: []string{
				"AtLeastOnce", "Capabilities.AtLeastOnce",
				"AcquireResultAcquired", "AcquireResultCompleted",
				"AcquireResultInProgress", "ErrInProgress", "ErrMissingLockID",
				"ErrUnknownAcquireResult", "ErrLockLost", "processing lease",
			},
			englishDocTerms: []string{"Inbox middleware"},
			chineseDocTerms: []string{"Inbox 中间件"},
		},
	}

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)
	englishDocs := readCorpus("English event docs", filepath.Join(docsRoot, eventDocsPath))
	chineseDocs := readCorpus("Chinese event docs", filepath.Join(docsRoot, chineseEventDocsPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))
	docs := []corpus{englishDocs, chineseDocs}

	auditEntries := loadEventAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntries := loadEventManifestEntries(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	reviews, contracts := loadEventContracts(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntries := loadLiveEventEntries(sourceRoot, docsRoot)

	var failures []string
	for _, pkg := range eventPackages {
		failures = append(failures, verifySurfaceEntry("live public API inventory", pkg, liveEntries[pkg.pkg])...)
		failures = append(failures, verifySurfaceEntry("API audit manifest", pkg, manifestEntries[pkg.pkg])...)
		failures = append(failures, verifyReviewSurface(pkg, reviews[pkg.pkg])...)
		failures = append(failures, verifyAuditEntries(pkg, auditEntries[pkg.pkg])...)
		failures = append(failures, verifyCoverage(pkg, auditEntries[pkg.pkg], manifestEntries[pkg.pkg], reviews[pkg.pkg], contracts[pkg.contractID])...)
		for _, index := range []corpus{englishIndex, chineseIndex} {
			failures = append(failures, verifyGeneratedIndexSection(index, pkg, auditEntries[pkg.pkg])...)
		}
	}
	for _, doc := range docs {
		failures = append(failures, verifyEventContractTerms(doc, contracts)...)
	}

	for _, check := range checks {
		source := readCorpus(check.sourcePath, filepath.Join(sourceRoot, check.sourcePath))
		failures = append(failures, missingTerms(source, check.sourceTerms)...)
		for _, doc := range docs {
			failures = append(failures, missingTerms(doc, check.docTerms)...)
		}
		failures = append(failures, missingTerms(englishDocs, check.englishDocTerms)...)
		failures = append(failures, missingTerms(chineseDocs, check.chineseDocTerms)...)
	}
	failures = append(failures, verifyContractLedger(reviews, contracts, sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("event contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("Event contract docs verified: %d packages, %d source files, %d doc mirrors\n", len(eventPackages), len(checks), len(docs))
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

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}

func verifySurfaceEntry(label string, pkg eventPackage, entry manifestEntry) []string {
	var failures []string
	if entry.Package != pkg.pkg {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q want %q", label, entry.Package, pkg.pkg))
	}
	if entry.TopLevel != pkg.topLevel || entry.Fields != pkg.fields ||
		entry.Methods != pkg.methods || entry.Fingerprint != pkg.fingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch for %s: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			pkg.pkg,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			pkg.topLevel, pkg.fields, pkg.methods, pkg.fingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(pkg eventPackage, review contractPackageReview) []string {
	var failures []string
	if review.Package != pkg.pkg {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q want %q", review.Package, pkg.pkg))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch for %s: got %q", pkg.pkg, review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != pkg.topLevel ||
		review.ReviewedSurface.Fields != pkg.fields ||
		review.ReviewedSurface.Methods != pkg.methods ||
		review.ReviewedSurface.EntryCount != pkg.entries ||
		review.ReviewedSurface.Fingerprint != pkg.fingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch for %s: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			pkg.pkg,
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !sameSet(review.ContractIDs, []string{pkg.contractID}) {
		failures = append(failures, fmt.Sprintf("contract package review contract ids mismatch for %s: got %v want %v", pkg.pkg, review.ContractIDs, []string{pkg.contractID}))
	}

	return failures
}

func verifyAuditEntries(pkg eventPackage, entries []auditEntry) []string {
	var failures []string
	if len(entries) != pkg.entries {
		failures = append(failures, fmt.Sprintf("event audit entry count mismatch for %s: got %d want %d", pkg.pkg, len(entries), pkg.entries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != pkg.pkg {
			failures = append(failures, "non-event package entry passed into event verifier: "+entry.ID)
		}
		if strings.Contains(entry.Package, "/internal/") {
			failures = append(failures, "internal package included in event verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate event audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "event audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != pkg.topLevel || counts["field"] != pkg.fields || counts["method"] != pkg.methods {
		failures = append(failures, fmt.Sprintf("event audit kind counts mismatch for %s: top/field/method=%d/%d/%d", pkg.pkg, counts["top"], counts["field"], counts["method"]))
	}

	return failures
}

func verifyCoverage(
	pkg eventPackage,
	entries []auditEntry,
	manifestEntry manifestEntry,
	review contractPackageReview,
	contract contractEntry,
) []string {
	var failures []string
	expected := []string{eventContractCoverage}
	if manifestEntry.Tier != "feature" {
		failures = append(failures, fmt.Sprintf("manifest tier mismatch for %s: got %q want feature", pkg.pkg, manifestEntry.Tier))
	}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest event coverage mismatch for %s: got %v want %v", pkg.pkg, manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review event coverage mismatch for %s: got %v want %v", pkg.pkg, review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry event coverage mismatch for %s: got %v want %v", pkg.pkg, contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, pkg eventPackage, entries []auditEntry) []string {
	section := packageSection(index.content, pkg.pkg)
	if section == "" {
		return []string{index.label + " missing event package section for " + pkg.pkg}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s event index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyEventContractTerms(doc corpus, contracts map[string]contractEntry) []string {
	var failures []string
	for _, pkg := range eventPackages {
		contract := contracts[pkg.contractID]
		for _, term := range contract.Terms {
			if !strings.Contains(doc.content, term) {
				failures = append(failures, fmt.Sprintf("%s missing event contract term for %s: %s", doc.label, contract.ID, term))
			}
		}
	}

	return failures
}

func verifyContractLedger(reviews map[string]contractPackageReview, contracts map[string]contractEntry, sourceRoot string) []string {
	var failures []string
	for _, pkg := range eventPackages {
		review := reviews[pkg.pkg]
		contract := contracts[pkg.contractID]
		if contract.ID != pkg.contractID {
			failures = append(failures, fmt.Sprintf("event contract id mismatch for %s: got %q", pkg.pkg, contract.ID))
		}
		if contract.Package != pkg.pkg {
			failures = append(failures, fmt.Sprintf("event contract package mismatch for %s: got %q", pkg.pkg, contract.Package))
		}
		if contract.Kind != "runtime-contract" {
			failures = append(failures, fmt.Sprintf("event contract kind mismatch for %s: got %q", pkg.pkg, contract.Kind))
		}
		if contract.Disposition != "documented:semantic-contract" {
			failures = append(failures, fmt.Sprintf("event contract disposition mismatch for %s: got %q", pkg.pkg, contract.Disposition))
		}
		if len(contract.Terms) == 0 {
			failures = append(failures, "event contract terms empty for "+contract.ID)
		}

		allEvidence := append([]string{}, review.SourceEvidence...)
		allEvidence = append(allEvidence, contract.SourceEvidence...)
		allEvidence = append(allEvidence, contract.TestEvidence...)
		for _, item := range allEvidence {
			path, lineText, ok := strings.Cut(item, ":")
			if !ok || lineText == "" {
				failures = append(failures, "event contract evidence missing line number: "+item)
				continue
			}
			if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
				failures = append(failures, "event contract evidence missing file: "+item)
			}
		}
	}

	return failures
}

func runSourceTests(sourceRoot string) []string {
	return runCommand(sourceRoot, "go", "test", "./event/...")
}

func loadEventAuditEntries(path string) map[string][]auditEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger auditLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	result := map[string][]auditEntry{}
	for _, entry := range ledger.Entries {
		if _, ok := eventPackageByName(entry.Package); ok {
			result[entry.Package] = append(result[entry.Package], entry)
		}
	}
	for pkg := range result {
		sort.Slice(result[pkg], func(i, j int) bool {
			return result[pkg][i].ID < result[pkg][j].ID
		})
	}

	return result
}

func loadEventManifestEntries(path string) map[string]manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}

	result := map[string]manifestEntry{}
	for _, entry := range m.Packages {
		if _, ok := eventPackageByName(entry.Package); ok {
			result[entry.Package] = entry
		}
	}

	return result
}

func loadEventContracts(path string) (map[string]contractPackageReview, map[string]contractEntry) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger contractLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	reviews := map[string]contractPackageReview{}
	for _, review := range ledger.PackageReviews {
		if _, ok := eventPackageByName(review.Package); ok {
			reviews[review.Package] = review
		}
	}

	contracts := map[string]contractEntry{}
	for _, item := range ledger.Entries {
		for _, pkg := range eventPackages {
			if item.ID == pkg.contractID {
				contracts[item.ID] = item
			}
		}
	}

	return reviews, contracts
}

func loadLiveEventEntries(sourceRoot, docsRoot string) map[string]manifestEntry {
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

	result := map[string]manifestEntry{}
	for _, entry := range entries {
		if _, ok := eventPackageByName(entry.Package); ok {
			result[entry.Package] = entry
		}
	}

	return result
}

func eventPackageByName(name string) (eventPackage, bool) {
	for _, pkg := range eventPackages {
		if pkg.pkg == name {
			return pkg, true
		}
	}

	return eventPackage{}, false
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
