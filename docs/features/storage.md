---
sidebar_position: 4
---

# File Storage

VEF ships a provider-neutral storage abstraction, three built-in providers, a multipart upload protocol, a typed CRUD lifecycle facade for keeping model file references in sync with the backend, and a transactional outbox for downstream cleanup.

> The storage module went through a heavy overhaul after v0.21: the legacy `Promoter[T]` was replaced by `Files` / `FilesFor[T]`, the upload protocol unified on chunked multipart with an explicit claim/queue lifecycle, principal authorization was threaded through the lifecycle, and the on-the-wire `Consume` / `Enqueue` surface was renamed and pruned. This page describes the current public surface; older snapshots are not API-compatible.

## Supported Providers

| Provider value | Backend |
| --- | --- |
| `memory` | in-process map; tests and ephemeral demos |
| `filesystem` | local filesystem |
| `minio` | MinIO / S3-compatible object storage |

`storage.provider` selects the backend. Without configuration the module
defaults to `memory` and logs a warning; objects are lost on restart.

Set `vef.storage.auto_migrate = true` when the storage tables should be created
by the module at startup. The migration is idempotent and checks
`sys_storage_upload_claim`, `sys_storage_upload_part`, and
`sys_storage_pending_delete`.

## `storage.Service` Interface

Application code depends on `storage.Service`, never on a provider-specific type:

```go
type Service interface {
    PutObject(ctx, opts PutObjectOptions) (*ObjectInfo, error)
    GetObject(ctx, opts GetObjectOptions) (io.ReadCloser, error)
    DeleteObject(ctx, opts DeleteObjectOptions) error
    DeleteObjects(ctx, opts DeleteObjectsOptions) error
    CopyObject(ctx, opts CopyObjectOptions) (*ObjectInfo, error)
    StatObject(ctx, opts StatObjectOptions) (*ObjectInfo, error)
}
```

Option types: `PutObjectOptions`, `GetObjectOptions`, `DeleteObjectOptions`, `DeleteObjectsOptions`, `CopyObjectOptions`, `StatObjectOptions`. Use the option struct for every call — direct positional arguments are not supported on purpose so that adding fields stays additive.

## Multipart Upload

The framework's upload protocol is **chunked multipart only** — the original single-PUT upload was removed in v0.21 (`refactor(storage): unify upload protocol on multipart`). Every backend implements `storage.Multipart`:

```go
type Multipart interface {
    PartSize() int64
    MaxPartCount() int
    InitMultipart(ctx, opts InitMultipartOptions) (*MultipartSession, error)
    PutPart(ctx, opts PutPartOptions) (*PartInfo, error)
    CompleteMultipart(ctx, opts CompleteMultipartOptions) (*ObjectInfo, error)
    AbortMultipart(ctx, opts AbortMultipartOptions) error
}
```

Obtain the typed handle with `storage.MultipartFor(svc)` (returns `nil` when the backend does not implement chunked uploads). The contract guarantees:

- Distinct part numbers may upload concurrently; same-part calls are last-writer-wins.
- Every non-final part must be at least `PartSize()` bytes.
- `CompleteMultipart` verifies every recorded `(PartNumber, ETag)` and that parts cover `1..N` contiguously.
- Sessions close after `CompleteMultipart` or `AbortMultipart`; further calls return `ErrUploadSessionNotFound`. `AbortMultipart` is idempotent.

> The `sys/storage.list_parts` RPC action exists to let clients resume an in-flight upload, but it is served from the framework's part-store table, not from a `ListParts` method on `storage.Multipart` — the backend interface itself only exposes the six methods above.

## Built-In Resource: `sys/storage`

The storage module registers an RPC resource with the multipart upload actions:

| Action | Purpose |
| --- | --- |
| `init_upload` | create a pending claim, open a multipart session, and return opaque `claimId` plus the negotiated `partSize` |
| `upload_part` | upload one part of an open session |
| `list_parts` | inspect parts already uploaded for a session |
| `complete_upload` | seal a session; the server assembles the final manifest from recorded parts |
| `abort_upload` | abort and release a session |

Download is served via the proxy middleware described below.

All HTTP uploads use this same protocol: `init_upload -> upload_part ->
complete_upload`. Small files still return `partCount = 1`; there is no
single-PUT HTTP action. The `public` flag defaults to private behavior, and
`vef.storage.allow_public_uploads` must be true before clients can request
`pub/` keys.

Client `contentType` values are sanitized. Safe binary, image, audio, video,
font, archive, and PDF types are accepted; unsafe same-origin types such as
`text/html` and `application/javascript` are replaced by extension detection or
`application/octet-stream`.

## Visibility Prefixes

Object keys carry their intended visibility as a prefix:

| Constant | Value | Meaning |
| --- | --- | --- |
| `storage.PublicPrefix` | `pub/` | world-readable; default ACL grants read |
| `storage.PrivatePrefix` | `priv/` | controlled by business state via `FileACL` |

The storage resource emits keys under `pub/` or `priv/` depending on the upload's `public` flag. Proxy downloads serve `pub/*` anonymously and call `FileACL` for non-public keys; the storage backend itself does not enforce visibility.

## FileACL

`storage.FileACL` decides whether a principal may read a private key.

```go
type FileACL interface {
    CanRead(ctx context.Context, principal *security.Principal, key string) (bool, error)
}
```

Default behavior (`storage.DefaultFileACL`): grant read access only to keys under `pub/`. The proxy short-circuits `pub/*` before calling the ACL so public files work without an auth token; business code overrides `FileACL` via `vef.SupplyFileACL(...)` for private keys and ownership-aware reads.

## Storage Proxy Middleware

The module mounts an app-level download route:

```
GET /storage/files/<key>
```

Behavior:

| Surface | Behavior |
| --- | --- |
| Routing | app middleware named `storage_proxy` at order `900`; not an RPC action and not dispatched by the API engine |
| Key validation | URL-decodes `<key>` once; rejects empty keys, absolute paths, `..` segments, backslashes, NUL bytes, redundant slashes, and trailing slashes |
| Access | serves `pub/*` anonymously; every other key calls `FileACL.CanRead` with the request principal |
| Content type | uses backend metadata or extension detection, then sanitizes unsafe types to `application/octet-stream`; always sends `X-Content-Type-Options: nosniff` |
| Cache headers | `pub/*` gets `Cache-Control: public, max-age=3600, immutable` and an `ETag` when stat data has one; non-public keys get `Cache-Control: private, no-store` and no `ETag` |

## Upload Claims and Pending Delete (Lifecycle)

`init_upload` persists an `upload_claim` row owned by the calling principal with
status `pending`. `complete_upload` marks that same claim as `uploaded`. Until
the business model adopts the key (via `Files.OnCreate` / `OnUpdate`), the
object lives in a quarantined state — a periodic sweeper either recovers an
expired-but-completed multipart object by marking it uploaded, or enqueues the
abandoned object for asynchronous deletion (`DeleteReasonClaimExpired`).

Business writes therefore split into two transactional surfaces:

- **Claim consumer**: deletes the `upload_claim` row in the same transaction as the business insert.
- **Delete enqueuer**: inserts a `pending_delete` row for objects that should be reclaimed asynchronously (replaced field values, deleted business rows).

A background `DeleteWorker` then drains `pending_delete` rows against the backend and applies retry/backoff. Successfully drained rows emit `vef.storage.file.deleted`; rows that exhaust the retry budget are removed from the queue after emitting `vef.storage.delete.dead_letter`, which is the durable signal for manual investigation.

Storage fails fast at startup unless `vef.storage.file.claimed`,
`vef.storage.file.deleted`, and `vef.storage.delete.dead_letter` route through a
transactional event transport. In practice, enable the outbox transport and add
a route for `vef.storage.*` to `outbox`, or set the default event transport to
`outbox`.

## `Files` and `FilesFor[T]`

The high-level CRUD lifecycle facade — this replaced the older `Promoter[T]`:

```go
type Files interface {
    OnCreate(ctx, tx orm.DB, principal *security.Principal, model any) error
    OnUpdate(ctx, tx orm.DB, principal *security.Principal, oldModel, newModel any) error
    OnDelete(ctx, tx orm.DB, model any) error
}
```

Key semantics:

- All three methods **must run inside a business transaction** (`orm.DB.RunInTx`). The supplied `tx` is the business-DB instance, so claim consumption and pending-delete bookkeeping commit or roll back atomically with the business write.
- `OnCreate` / `OnUpdate` take a `*security.Principal` — only claims owned by that principal can be adopted. Nil / anonymous principals fail with `ErrAccessDenied`. Background jobs that legitimately operate on behalf of the system pass a synthetic system principal explicitly.
- `OnDelete` does not consume claims and therefore takes no principal; row ownership must be verified at the CRUD layer first.
- `FileClaimedEvent` is published through the **outbox transport inside the caller's transaction** (`event.WithTx`) — subscribers see the event only if the business transaction commits.

### Typed counterpart

`storage.FilesFor[T]` resolves the meta spec once at construction so the per-call reflect lookup disappears:

```go
files := storage.NewFilesFor[User](filesFacade)
err := files.OnCreate(ctx, tx, principal, &user)
```

CRUD lifecycle hooks were migrated to `FilesFor[T]` in v0.22 (`refactor(crud): use FilesFor[T] for typed file lifecycle hooks`); custom hooks should follow the same pattern.

## Two Ways To Claim A File

When the user uploads a file, the framework keeps it in a "pending" state until your business code claims it. There are two ways to do that — pick the one that matches what your code already has:

- **Have a model struct?** Pass it to `Files` / `FilesFor[T]` and the framework will figure out the file fields by itself.
- **Just have a file key (or a list of keys)?** Call `ClaimConsumer.Consume(...)` directly.

Both end up doing the same thing — the second one is just the manual version of the first. Use whichever fits the call site better.

### Way 1: Pass in the struct (the easy way)

Tag the file fields with `meta:"uploaded_file"`, then hand the struct to `FilesFor[T]`. That's it.

```go
type Article struct {
    orm.FullAuditedModel
    CoverImage string   `json:"coverImage" bun:"cover_image" meta:"uploaded_file"`
    Gallery    []string `json:"gallery"    bun:"gallery,array" meta:"uploaded_file"`
    Body       string   `json:"body"       bun:"body"        meta:"rich_text"`
}

files := storage.NewFilesFor[Article](filesFacade)

err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
    if _, err := tx.NewInsert().Model(article).Exec(ctx); err != nil {
        return err
    }
    // Claim every file referenced by `article` in one call.
    return files.OnCreate(ctx, tx, principal, article)
})
```

On update, pass both the old and the new model — the framework claims the new files and queues the replaced ones for deletion:

```go
err := files.OnUpdate(ctx, tx, principal, oldArticle, newArticle)
```

On delete, pass the model — every referenced file gets queued for deletion:

```go
err := files.OnDelete(ctx, tx, article)
```

This is what regular CRUD already uses under the hood. If a struct fits, this is what you want.

### Way 2: Pass in the file key (when there's no struct)

Sometimes you don't have a model — maybe it's a background job, a custom upload flow, or you just want to claim one specific key. Inject `storage.ClaimConsumer` and call `Consume` with a slice of keys:

```go
err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
    if _, err := tx.NewInsert().Model(report).Exec(ctx); err != nil {
        return err
    }
    // Claim the file directly by its key.
    return claims.Consume(ctx, tx, principal, []string{report.FileKey})
})
```

If you also need to delete a file (e.g. the previous version), use `storage.DeleteEnqueuer`:

```go
err := deletes.Enqueue(ctx, tx,
    []string{oldKey},
    storage.DeleteReasonReplaced, // or DeleteReasonDeleted
)
```

A few things to keep in mind:

- Always call these inside `RunInTx` and pass the same `tx` — that's how the claim and your business write commit together.
- `Consume` only succeeds for files uploaded by **the same principal**. Trying to claim someone else's file returns `storage.ErrClaimNotFound`.
- A nil or anonymous principal returns `storage.ErrAccessDenied`. Background jobs need to construct a real system principal first.
- Empty / nil key slices are fine — they do nothing.
- Use `DeleteReasonReplaced` when overwriting a field, `DeleteReasonDeleted` when removing the owning record. `DeleteReasonClaimExpired` is for the framework only, don't pass it.

> If you ever catch yourself writing reflection to scan a struct's file fields, stop — that's exactly what `FilesFor[T]` does. Switch back to Way 1.

## Meta-Tagged Model Fields

Fields participate in the lifecycle by carrying a `meta` tag:

```go
type User struct {
    orm.FullAuditedModel

    Avatar   string            `json:"avatar" bun:"avatar" meta:"uploaded_file"`
    Gallery  []string          `json:"gallery" bun:"gallery,array" meta:"uploaded_file,category:gallery"`
    Profiles map[string]string `json:"profiles" bun:"profiles" meta:"uploaded_file"`
    Bio      string            `json:"bio" bun:"bio" meta:"rich_text"`
    Notes    string            `json:"notes" bun:"notes" meta:"markdown"`
}
```

| `meta` value | Field shape | Extraction strategy |
| --- | --- | --- |
| `uploaded_file` | `string` / `*string` / `[]string` / `map[string]string` | the value(s) are treated as file keys — for maps the **values** are the keys; the map's own keys are arbitrary labels |
| `rich_text` | `string` | scan HTML for embedded resource URLs and translate via `URLKeyMapper` |
| `markdown` | `string` | scan Markdown for embedded resource URLs and translate via `URLKeyMapper` |

> v0.21 added `map[string]string` support for `uploaded_file` fields (`feat(storage): support map[string]string for uploaded_file fields`).

Use `meta:"dive"` on a nested struct field when the file references live inside that nested struct; the scanner will recurse into the nested value and pick up its own `meta:"uploaded_file"`, `meta:"rich_text"`, and `meta:"markdown"` fields. Unsupported field shapes are ignored instead of producing refs.

`URLKeyMapper` translates rich-text/markdown URLs to storage keys during reconciliation. The framework DI graph supplies `storage.ProxyURLKeyMapper` by default, so content that embeds `/storage/files/<key>` is reconciled without extra wiring. If you call `storage.NewFiles(...)` directly, a nil mapper is normalized to `IdentityURLKeyMapper`; pass `&storage.IdentityURLKeyMapper{}` only when business content embeds bare keys directly.

The mapper surface is explicit in both directions: `URLToKey` consumes content
URLs during reconciliation, and `KeyToURL` is used when code needs to render
stored keys back into URLs.

Use `storage.ProxyURLKeyMapper{Prefix: storage.DefaultProxyPrefix}` when content
embeds the framework proxy URL form (`/storage/files/<key>`). The public helpers
`ReplaceHtmlURLs(content, replacements)` and `ReplaceMarkdownURLs(content,
replacements)` rewrite embedded URLs in rendered content, typically after mapping
storage keys through `URLKeyMapper.KeyToURL`.

## Storage Events

| Type constant / topic | Payload / constructor | JSON payload | Trigger |
| --- | --- | --- | --- |
| `EventTypeFileClaimed` / `vef.storage.file.claimed` | `FileClaimedEvent`; `NewFileClaimedEvent(key)` | `fileKey` | a previously pending claim was adopted by a business transaction (`Files.OnCreate` or update new-side) |
| `EventTypeFileDeleted` / `vef.storage.file.deleted` | `FileDeletedEvent`; `NewFileDeletedEvent(key, reason)` | `fileKey`, `reason` | the delete worker successfully removed an object from the backend |
| `EventTypeDeleteDeadLetter` / `vef.storage.delete.dead_letter` | `DeleteDeadLetterEvent`; `NewDeleteDeadLetterEvent(id, key, reason, attempts, lastErr)` | `pendingDeleteId`, `fileKey`, `reason`, `attempts`, optional `lastError` | the delete worker exhausted retries for a row; the queue row is removed after this event is published |

All three are published through the outbox transport with `event.WithTx(...)`. `FileClaimedEvent` shares the caller's business transaction; `FileDeletedEvent` and `DeleteDeadLetterEvent` share the delete worker's bookkeeping transaction. Subscribers attach with `event.WithGroup("...")` on the downstream sink transport and rely on the Inbox middleware for dedupe.

`DeleteReason` values forwarded onto the events:

| Reason | Wire value | Meaning |
| --- | --- | --- |
| `DeleteReasonReplaced` | `replaced` | an `uploaded_file` field was overwritten with a new key |
| `DeleteReasonDeleted` | `deleted` | the owning business row was deleted |
| `DeleteReasonClaimExpired` | `claim_expired` | a pending claim expired (framework-internal sweeper only) |

Dead-letter events carry a sanitized `lastError` classification rather than raw
backend errors. Current values are `access_denied`, `bucket_not_found`,
`session_not_found`, and `transient`.

Public supporting APIs:

| API group | Public surface |
| --- | --- |
| event constructors | `EventTypeFileClaimed`, `EventTypeFileDeleted`, `EventTypeDeleteDeadLetter`, `NewFileClaimedEvent`, `NewFileDeletedEvent`, `NewDeleteDeadLetterEvent` |
| facade constructors | `NewFiles`, `NewFilesFor`, `MultipartFor` |
| lifecycle services | `ClaimConsumer`, `DeleteEnqueuer`, `Files`, `FilesFor[T]` |
| storage interfaces | `Service`, `Multipart`, `FileACL`, `URLKeyMapper` |
| URL mappers | `DefaultFileACL`, `IdentityURLKeyMapper`, `ProxyURLKeyMapper`, `DefaultProxyPrefix` |
| option structs | `PutObjectOptions`, `GetObjectOptions`, `DeleteObjectOptions`, `DeleteObjectsOptions`, `CopyObjectOptions`, `StatObjectOptions`, `InitMultipartOptions`, `PutPartOptions`, `CompleteMultipartOptions`, `AbortMultipartOptions` |
| result structs | `ObjectInfo`, `MultipartSession`, `PartInfo`, `CompletedPart`, `FileRef` |
| meta constants | `MetaType`, `MetaTypeUploadedFile`, `MetaTypeRichText`, `MetaTypeMarkdown` |

The storage package audit currently locks **181 public storage entries** in the
generated API ledger. The grouped member surface covers **81 grouped storage
field/method entries** across **29 storage receiver/type families**: **50
exported storage field entries** and **31 exported storage method entries**.
The generated public API index remains the complete signature list; this page
documents the semantic families and user-facing runtime behavior.

## Error Sentinels

| Error | Cause |
| --- | --- |
| `storage.ErrInvalidFileKey` | malformed object key on a stat/download request |
| `storage.ErrFileNotFound` | object missing from the backend |
| `storage.ErrFailedToGetFile` | backend read failed |
| `storage.ErrUploadSessionNotFound` | multipart session already closed or never opened |
| `storage.ErrPartTooSmall` | non-final part smaller than `PartSize()` |
| `storage.ErrPartETagMismatch` | recorded part ETag disagrees with backend state during completion |
| `storage.ErrPartNumberOutOfRange` | parts don't cover `1..N` contiguously |
| `storage.ErrClaimNotFound` | a claim referenced by `Consume` doesn't exist or belongs to another principal |
| `storage.ErrAccessDenied` | anonymous / nil principal passed to a lifecycle method |

Upload API errors also expose matching `ErrCode*` constants in the `2200-2299`
range: `ErrCodeInvalidFileKey`, `ErrCodeFileNotFound`,
`ErrCodeFailedToGetFile`, `ErrCodeClaimNotPending`, `ErrCodeClaimExpired`,
`ErrCodeUploadSizeExceedsLimit`, `ErrCodeMultipartNotSupported`,
`ErrCodePublicUploadsNotAllowed`, `ErrCodeUploadTooManyParts`,
`ErrCodeTooManyPendingUploads`, `ErrCodeUploadRequiresMultipart`,
`ErrCodeUploadRequiresFile`, `ErrCodeClaimNotMultipart`,
`ErrCodeUploadPartNumberOutOfRange`, `ErrCodeUploadPartTooLarge`,
`ErrCodeUploadPartTooSmall`, `ErrCodeUploadPartsIncomplete`,
`ErrCodeUploadObjectNotFound`, `ErrCodeUploadSizeMismatch`, and
`ErrCodeAbortFailed`. Additional public sentinels include `ErrClaimNotPending`,
`ErrClaimExpired`, `ErrUploadSizeExceedsLimit`, `ErrMultipartNotSupported`,
`ErrPublicUploadsNotAllowed`, `ErrUploadTooManyParts`,
`ErrTooManyPendingUploads`, `ErrUploadRequiresMultipart`,
`ErrUploadRequiresFile`, `ErrClaimNotMultipart`,
`ErrUploadPartNumberOutOfRange`, `ErrUploadPartTooLarge`,
`ErrUploadPartTooSmall`, `ErrUploadPartsIncomplete`,
`ErrUploadObjectNotFound`, `ErrUploadSizeMismatch`, `ErrAbortFailed`,
`ErrBucketNotFound`, `ErrObjectNotFound`, and `ErrInvalidBucketName`.

## Minimal Service Example

```go
package avatars

import (
    "context"
    "strings"

    "github.com/coldsmirk/vef-framework-go/storage"
)

func SaveAvatar(ctx context.Context, svc storage.Service) error {
    _, err := svc.PutObject(ctx, storage.PutObjectOptions{
        Key:         "pub/avatars/user-1001.txt",
        Reader:      strings.NewReader("demo"),
        Size:        int64(len("demo")),
        ContentType: "text/plain",
    })

    return err
}
```

## CRUD Integration Pattern

For models with `meta`-tagged file fields, integrate via `FilesFor[T]` from a typed hook:

```go
filesUser := storage.NewFilesFor[User](filesFacade)

create := crud.NewCreate[User, UserParams]().
    AfterTx(func(ctx context.Context, tx orm.DB, principal *security.Principal, model *User) error {
        return filesUser.OnCreate(ctx, tx, principal, model)
    })
```

Generic CRUD already wires `FilesFor[T]` for the standard write builders (see [Hooks](../guide/hooks)); custom write paths should follow the same pattern.

## Practical Advice

- Depend on `storage.Service` and `storage.Multipart`, not provider types.
- Keep all `Files` / `FilesFor[T]` calls inside the business transaction — that is the whole point of the facade.
- Treat unconfirmed objects as quarantined: the claim sweeper will eventually evict them; relying on raw `PutObject` keys without a claim bypasses lifecycle tracking.
- Register a real `FileACL` once you store private files; the default denies every `priv/*` read.
- Subscribe to `vef.storage.delete.dead_letter` for ops dashboards — the queue row is already retired, and the event carries the details operators need.
- Extension group names used by the module are `vef:api:resources` and `vef:app:middlewares`; use `vef.SupplyURLKeyMapper(...)` when replacing URL mapping.

## Next Step

Read [Custom Handlers](../guide/custom-handlers) to combine direct `storage.Service` use with business workflows, or [Event Bus](./event-bus) for the outbox transport that backs the lifecycle events.
