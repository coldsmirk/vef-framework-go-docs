---
sidebar_position: 7
---

# 文件存储

VEF 提供一套与 provider 无关的存储抽象、三种内置 provider、分片上传协议、用于 CRUD 流程同步模型与后端文件引用的类型化生命周期门面，以及用于下游清理的事务性 outbox。

> 存储模块在 v0.21 之后经历了大量重构：原来的 `Promoter[T]` 已被 `Files` / `FilesFor[T]` 取代，上传协议统一为分片 multipart，引入 claim / queue 生命周期，把 principal 授权贯穿全流程，对外的 `Consume` / `Enqueue` API 也重命名精简。本页描述当前公开 API surface；老版本快照与现在已经不再兼容。

## 支持的 Provider

| 配置值 | 后端 |
| --- | --- |
| `memory` | 进程内 map；测试和短期演示 |
| `filesystem` | 本地文件系统 |
| `minio` | MinIO / S3 兼容对象存储 |

`storage.provider` 选择后端。未配置时默认 `memory` 并输出 warning；对象会在进程重启后丢失。

如果希望 storage 模块在启动时创建所需表，设置
`vef.storage.auto_migrate = true`。迁移是幂等的，会检查
`sys_storage_upload_claim`、`sys_storage_upload_part` 和
`sys_storage_pending_delete`。

## `storage.Service` 接口

应用代码依赖 `storage.Service`，不依赖任何 provider 特定类型：

```go
type Service interface {
    PutObject(ctx, opts PutObjectOptions) (*ObjectInfo, error)
    GetObject(ctx, opts GetObjectOptions) (io.ReadCloser, *ObjectInfo, error)
    DeleteObject(ctx, opts DeleteObjectOptions) error
    DeleteObjects(ctx, opts DeleteObjectsOptions) error
    CopyObject(ctx, opts CopyObjectOptions) (*ObjectInfo, error)
    StatObject(ctx, opts StatObjectOptions) (*ObjectInfo, error)
}
```

每个方法都使用对应的 option struct：`PutObjectOptions`、`GetObjectOptions`、`DeleteObjectOptions`、`DeleteObjectsOptions`、`CopyObjectOptions`、`StatObjectOptions`。框架刻意不支持位置参数 ——这样后续追加字段始终是可叠加变更。

`GetObject` 会返回 body reader 和 best-effort `ObjectInfo`。调用方必须关闭
reader，并对 `ObjectInfo` 做 nil-check。

## 分片上传

框架的上传协议**只接受分片 multipart**——原来的单次 PUT 上传在 v0.21 被移除（`refactor(storage): unify upload protocol on multipart`）。每个后端都实现 `storage.Multipart`：

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

用 `storage.MultipartFor(svc)` 拿到类型化句柄（后端不支持分片上传时返回 `nil`）。契约保证：

- 不同 part number 允许并发上传；相同 part number 的并发调用是 last-writer-wins。
- 除最后一片外，每片必须不小于 `PartSize()` 字节。
- `CompleteMultipart` 会校验每个已记录的 `(PartNumber, ETag)` 以及 parts 是否覆盖 `1..N` 连续区间。
- 调用 `CompleteMultipart` 或 `AbortMultipart` 后会话关闭，后续操作返回 `ErrUploadSessionNotFound`。`AbortMultipart` 幂等。

> `sys/storage.list_parts` RPC action 用于让客户端恢复上传中断，它由框架的 part-store 表服务，**不**是 `storage.Multipart` 上的 `ListParts` 方法 —— 后端接口本身只暴露上面 6 个方法。

## 内置资源：`sys/storage`

存储模块注册了一个 RPC 资源，包含分片上传动作：

| Action | 作用 |
| --- | --- |
| `init_upload` | 创建 pending claim、开启 multipart session，并返回不透明 `claimId` 与协商出的 `partSize` |
| `upload_part` | 上传一片 |
| `list_parts` | 列出已上传的 parts |
| `complete_upload` | 完成 session；服务端从已记录的 parts 组装最终清单 |
| `abort_upload` | 中止并释放 session |

下载通过下方代理中间件提供。

所有 HTTP 上传都使用同一套协议：`init_upload -> upload_part ->
complete_upload`。小文件也会返回 `partCount = 1`；不存在单次 PUT HTTP
action。`public` 标志默认按私有处理，只有 `vef.storage.allow_public_uploads`
为 true 时，客户端才可以请求 `pub/` key。

客户端传入的 `contentType` 会被清洗。安全的二进制、图片、音频、视频、
字体、压缩包和 PDF 类型会被接受；`text/html`、`application/javascript`
这类同源不安全类型会被扩展名探测结果或 `application/octet-stream` 替换。

## 客户端演练：分片上传

以下是客户端对 `POST /api` 实现的完整线上（wire）序列，示例为在 MinIO 后端上传一个 40 MiB 的 `report.pdf`（`partSize` 16 MiB，共 3 片）。五个 action 都要求认证（默认 Bearer），后续每个调用都通过 `init_upload` 返回的 `claimId` 路由 —— 后端的 multipart `UploadID` 永远不会离开服务端。成功响应使用标准信封（`message` 文案随语言变化）；失败复用同一信封，`code` 取存储错误码区间（2200–2299，见下方错误表）。

### 1. `init_upload`

```bash
curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "sys/storage",
    "action": "init_upload",
    "version": "v1",
    "params": {
      "filename": "report.pdf",
      "size": 41943040,
      "contentType": "application/pdf",
      "public": false
    }
  }'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "key": "priv/2026/07/09/6c9e6f0e-8d5a-4d5e-9a3b-2f4a1c7e9b21.pdf",
    "claimId": "b3a2c1d0-4e5f-47a9-8bcd-ef0123456789",
    "originalFilename": "report.pdf",
    "partSize": 16777216,
    "partCount": 3,
    "expiresAt": "2026-07-10T12:04:05Z"
  }
}
```

`size` 必须是精确的字节数：如果实际上传总量与声明不符，`complete_upload` 会删除已组装的对象并返回 size mismatch 错误。`partSize` 是后端权威的切片大小 —— 除最后一片外每片必须恰好是 `partSize` 字节，最后一片承载余量。

### 2. `upload_part`（× `partCount`）

`upload_part` 拒绝 JSON body。请发送 `multipart/form-data`：`resource`、`action`、`version` 是普通表单字段，`params` 是 JSON 字符串，文件字节放在名为 `file` 的表单分片里：

```bash
split -b 16777216 report.pdf part-   # part-aa, part-ab, part-ac

curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -F 'resource=sys/storage' \
  -F 'action=upload_part' \
  -F 'version=v1' \
  -F 'params={"claimId":"b3a2c1d0-4e5f-47a9-8bcd-ef0123456789","partNumber":1}' \
  -F 'file=@part-aa'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "partNumber": 1,
    "size": 16777216
  }
}
```

对 `partNumber` 2 和 3 重复此调用（`part-ac` 是 8388608 字节的余量）。不同 part number 可以并发上传；重发同一个 part number 会覆盖之前的字节（last-writer-wins）。后端 ETag 由服务端记录且刻意不返回 —— 客户端永远不需要回传 ETag。

### 3. `complete_upload`

```bash
curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "sys/storage",
    "action": "complete_upload",
    "version": "v1",
    "params": { "claimId": "b3a2c1d0-4e5f-47a9-8bcd-ef0123456789" }
  }'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "bucket": "app-files",
    "key": "priv/2026/07/09/6c9e6f0e-8d5a-4d5e-9a3b-2f4a1c7e9b21.pdf",
    "eTag": "9b2cf535f27731c974343645a3985328-3",
    "size": 41943040,
    "contentType": "application/pdf",
    "lastModified": "2026-07-09T12:08:15Z",
    "originalFilename": "report.pdf"
  }
}
```

服务端从自己的表组装 parts 清单；已记录的 parts 少于 `partCount` 时调用失败，返回 `ErrCodeUploadPartsIncomplete`。重试是幂等的 —— 后端 session 已关闭之后到达的重试会重新 stat 对象并返回相同形状的响应。此时 claim 处于 `uploaded` 状态，等待业务采纳（见下方 `Files`）。

### 恢复中断的上传：`list_parts`

只要 claim 仍处于 pending 且未过期（`expiresAt`），已被接受的 parts 在客户端重启后依然有效。先问服务端已经持有哪些 parts，跳过它们，只上传剩余部分：

```bash
curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "sys/storage",
    "action": "list_parts",
    "version": "v1",
    "params": { "claimId": "b3a2c1d0-4e5f-47a9-8bcd-ef0123456789" }
  }'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "parts": [
      { "partNumber": 1, "size": 16777216 },
      { "partNumber": 2, "size": 16777216 }
    ]
  }
}
```

这里缺第 3 片：把它传上去，然后调用 `complete_upload`。列表按 `partNumber` 升序排列，每个列出的 part 都已连同 ETag 记录在服务端，`complete_upload` 会原样采信。

### `abort_upload`

同样的 JSON 信封，`"action": "abort_upload"`，`params` 里带 `claimId`：

```json
{ "code": 0, "message": "Success", "data": null }
```

Abort 是幂等的：未知或已中止的 `claimId` 仍返回 `code: 0`。只有 pending 状态的 claim 会被中止 —— 对 `uploaded` claim 调用是 no-op，绝不会删除已完成的对象。

### 通过代理下载

下载是对代理路由的普通 HTTP `GET`，不是 RPC action：

```bash
curl -O http://localhost:8080/storage/files/pub/2026/07/09/6c9e6f0e-8d5a-4d5e-9a3b-2f4a1c7e9b21.pdf
```

`pub/*` key 匿名可读。其他 key 由代理带着请求 principal 调用 `FileACL.CanRead`；框架自身不会在这个路由上解析 token，所以私有下载既需要注册 `FileACL`，也需要在 app 层挂设置 principal 的中间件（见下方代理中间件小节）。

## 可见性前缀

object key 通过前缀表达可见意图：

| 常量 | 值 | 含义 |
| --- | --- | --- |
| `storage.PublicPrefix` | `pub/` | 世界可读；默认 ACL 直接放行 |
| `storage.PrivatePrefix` | `priv/` | 由业务侧 `FileACL` 控制 |

存储资源会根据上传时的 `public` 标志，把 key 落在 `pub/` 或 `priv/` 下。代理下载会匿名放行 `pub/*`，其他 key 才调用 `FileACL`；存储后端本身不做这种判断。

## FileACL

`storage.FileACL` 决定 principal 能否读取一个私有 key。

```go
type FileACL interface {
    CanRead(ctx context.Context, principal *security.Principal, key string) (bool, error)
}
```

默认实现 `storage.DefaultFileACL` 只允许 `pub/` 下的读。代理会在调用 ACL 前直接放行 `pub/*`，让公共文件无需 auth token 也可访问；业务代码通过 `vef.SupplyFileACL(...)` 注入自己的私有文件授权逻辑。

## 存储代理中间件

模块在 app 层挂了一个下载路由：

```
GET /storage/files/<key>
```

行为：

| 表面 | 行为 |
| --- | --- |
| 路由 | 名为 `storage_proxy`、顺序为 `900` 的 app 中间件；不是 RPC action，也不由 API engine 分发 |
| key 校验 | 对 `<key>` 做一次 URL 解码；拒绝空 key、绝对路径、`..` 片段、反斜杠、NUL 字节、重复斜杠和结尾斜杠 |
| 访问控制 | `pub/*` 匿名可读；其他 key 都会带着请求 principal 调用 `FileACL.CanRead` |
| Content-Type | 优先使用后端 metadata 或扩展名探测结果，再把不安全类型清洗为 `application/octet-stream`；总是发送 `X-Content-Type-Options: nosniff` |
| 缓存 header | `pub/*` 返回 `Cache-Control: public, max-age=3600, immutable`，如果 stat 数据有 ETag 也会返回 `ETag`；非公开 key 返回 `Cache-Control: private, no-store` 且不返回 `ETag` |

## 上传 Claim 与 Pending Delete（生命周期）

`init_upload` 会写入一个属于上传调用方、状态为 `pending` 的
`upload_claim` 行。`complete_upload` 会把同一条 claim 标记为 `uploaded`。
在业务模型真正引用（通过 `Files.OnCreate` / `OnUpdate`）这个 key 之前，
对象处于隔离状态——周期性 sweeper 会把已过期但对象已经完成的 multipart
claim 恢复为 uploaded，或者把废弃对象入队为异步删除
（`DeleteReasonClaimExpired`）。

业务写入因此分裂为两组事务性表面：

- **Claim consumer**：在业务 insert 同一事务中删除 `upload_claim` 行。
- **Delete enqueuer**：把被替换的字段值、被删除的业务行对应的对象，写入 `pending_delete` 行做异步回收。

后台 `DeleteWorker` 会按重试/退避策略处理 `pending_delete` 行。删除成功会发出 `vef.storage.file.deleted`；用尽重试预算的行会在发出 `vef.storage.delete.dead_letter` 后从队列移除，事件本身就是供人工排查的持久信号。

Storage 会在启动时检查 `vef.storage.file.claimed`、
`vef.storage.file.deleted` 和 `vef.storage.delete.dead_letter` 是否路由到
事务性 event transport；否则直接启动失败。实践中应启用 outbox transport，
并把 `vef.storage.*` 路由到 `outbox`，或者把默认 event transport 设为
`outbox`。

## `Files` 与 `FilesFor[T]`

CRUD 生命周期门面——已经取代了旧的 `Promoter[T]`：

```go
type Files interface {
    OnCreate(ctx, tx orm.DB, principal *security.Principal, model any) error
    OnUpdate(ctx, tx orm.DB, principal *security.Principal, oldModel, newModel any) error
    OnDelete(ctx, tx orm.DB, model any) error
}
```

关键语义：

- 三个方法**都必须在业务事务里调用**（`orm.DB.RunInTx`）。传入的 `tx` 就是业务事务 DB 实例，所以 claim 消费和 pending-delete 记账会和业务写入一起提交或回滚。
- `OnCreate` / `OnUpdate` 接受 `*security.Principal`——只有归属该 principal 的 claim 才会被采纳。nil / 匿名 principal 直接 `ErrAccessDenied`。背景任务如果合法地"代表系统"操作，需要显式传入一个合成的系统 principal。
- `OnDelete` 不消费 claim，所以不需要 principal；调用前请先在 CRUD 层验证行所有权。
- `FileClaimedEvent` 通过**outbox transport 在调用方事务里发布**（`event.WithTx`）——订阅者只有在业务事务提交后才能看到事件。

### 类型化版本

`storage.FilesFor[T]` 把 meta 解析提前到构造时完成，调用时不再做反射查找：

```go
files := storage.NewFilesFor[User](filesFacade)
err := files.OnCreate(ctx, tx, principal, &user)
```

CRUD 的生命周期 hook 在 v0.22 全面迁到了 `FilesFor[T]`（`refactor(crud): use FilesFor[T] for typed file lifecycle hooks`），自定义 hook 也应该照此写法。

## 认领文件的两种方式

用户上传的文件，在你的业务代码"认领"之前一直处于待定状态。认领有两种方式，按你手里有什么选：

- **手上有模型 struct？** 把它交给 `Files` / `FilesFor[T]`，框架会自己识别其中的文件字段。
- **手上只有一个文件 key（或一组 key）？** 直接调用 `ClaimConsumer.Consume(...)`。

两种方式最终做的事情一样——后者只是前者的手动版。哪个更符合调用现场就用哪个。

### 方式一：传结构体（推荐写法）

给文件字段加上 `meta:"uploaded_file"` 标签，把整个 struct 交给 `FilesFor[T]`，就完了。

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
    // 一行就认领 article 引用的所有文件
    return files.OnCreate(ctx, tx, principal, article)
})
```

更新时把新旧两份模型都传进去——框架会认领新文件、把被替换掉的旧文件入队删除：

```go
err := files.OnUpdate(ctx, tx, principal, oldArticle, newArticle)
```

删除时把模型传进去，里面所有文件都会入队删除：

```go
err := files.OnDelete(ctx, tx, article)
```

普通 CRUD 默认就走这条路径。只要结构体合适，优先选它。

### 方式二：传文件 key（没有结构体时）

有时候你手上并没有模型——可能是后台任务、自定义上传流程，或者只想认领某个具体 key。注入 `storage.ClaimConsumer`，把 key 切片传给 `Consume`：

```go
err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
    if _, err := tx.NewInsert().Model(report).Exec(ctx); err != nil {
        return err
    }
    // 直接按 key 认领文件
    return claims.Consume(ctx, tx, principal, []string{report.FileKey})
})
```

如果还需要删除文件（比如上一版的附件），用 `storage.DeleteEnqueuer`：

```go
err := deletes.Enqueue(ctx, tx,
    []string{oldKey},
    storage.DeleteReasonReplaced, // 或 DeleteReasonDeleted
)
```

几条要记的：

- 一定要放在 `RunInTx` 里调用，并且传入的 `tx` 必须是同一个 —— 这样认领和业务写入才会一起提交。
- `Consume` 只能认领**当前 principal 自己上传**的文件。试图认领别人的文件会返回 `storage.ErrClaimNotFound`。
- nil / 匿名 principal 直接返回 `storage.ErrAccessDenied`；后台任务需要先构造一个真正的系统 principal。
- 传空 / nil 切片没问题，就是什么都不做。
- 字段被覆盖时用 `DeleteReasonReplaced`，整条记录被删时用 `DeleteReasonDeleted`。`DeleteReasonClaimExpired` 是框架内部专用，不要传。

> 如果你发现自己在 `Consume` / `Enqueue` 之上手写反射扫结构体，就停下来 —— `FilesFor[T]` 做的就是这件事，切回方式一。

## Meta 标签字段

字段通过 `meta` tag 加入生命周期管理：

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

| `meta` 取值 | 字段形态 | 抽取策略 |
| --- | --- | --- |
| `uploaded_file` | `string` / `*string` / `[]string` / `map[string]string` | 字段值就是 file key — 对 map 而言取的是 **value**，map 的 key 只是自定义标签 |
| `rich_text` | `string` | 扫描 HTML 中嵌入的资源 URL，再通过 `URLKeyMapper` 翻译 |
| `markdown` | `string` | 扫描 Markdown 中嵌入的资源 URL，再通过 `URLKeyMapper` 翻译 |

> v0.21 增加了对 `uploaded_file` 字段使用 `map[string]string` 的支持（`feat(storage): support map[string]string for uploaded_file fields`）。

如果文件引用位于嵌套 struct 里，可以在外层字段上使用 `meta:"dive"`；扫描器会递归进入这个嵌套值，并识别其中自己的 `meta:"uploaded_file"`、`meta:"rich_text"` 和 `meta:"markdown"` 字段。不支持的字段形态会被忽略，不会生成文件引用。

`URLKeyMapper` 用来把富文本/Markdown 中的 URL 翻译成存储 key。框架 DI 图默认提供 `storage.ProxyURLKeyMapper`，因此内容里嵌入 `/storage/files/<key>` 时无需额外配置即可对账。如果你直接调用 `storage.NewFiles(...)`，nil mapper 会被规整为 `IdentityURLKeyMapper`；只有业务内容直接嵌入 bare key 时才传 `&storage.IdentityURLKeyMapper{}`。

这个 mapper 的方向是显式的：`URLToKey` 在对账时消费内容里的 URL，
`KeyToURL` 在代码需要把已存储 key 渲染回 URL 时使用。

当内容中嵌入的是框架代理 URL 形态（`/storage/files/<key>`）时，可以使用
`storage.ProxyURLKeyMapper{Prefix: storage.DefaultProxyPrefix}`。公开 helper
`ReplaceHtmlURLs(content, replacements)` 与 `ReplaceMarkdownURLs(content,
replacements)` 用于重写已渲染内容里的嵌入 URL，通常和
`URLKeyMapper.KeyToURL` 配合使用。

## 存储事件

| 类型常量 / topic | Payload / 构造器 | JSON payload | 触发 |
| --- | --- | --- | --- |
| `EventTypeFileClaimed` / `vef.storage.file.claimed` | `FileClaimedEvent`; `NewFileClaimedEvent(key)` | `fileKey` | 某次业务事务采纳了一个之前 pending 的 claim（`Files.OnCreate` 或 update 新侧） |
| `EventTypeFileDeleted` / `vef.storage.file.deleted` | `FileDeletedEvent`; `NewFileDeletedEvent(key, reason)` | `fileKey`, `reason` | delete worker 成功把对象从后端删除 |
| `EventTypeDeleteDeadLetter` / `vef.storage.delete.dead_letter` | `DeleteDeadLetterEvent`; `NewDeleteDeadLetterEvent(id, key, reason, attempts, lastErr)` | `pendingDeleteId`, `fileKey`, `reason`, `attempts`，可选 `lastError` | delete worker 重试用尽；该队列行会在事件发布后被删除 |

三者都通过 outbox transport 搭配 `event.WithTx(...)` 发布。`FileClaimedEvent` 共享调用方业务事务；`FileDeletedEvent` 与 `DeleteDeadLetterEvent` 共享 delete worker 的记账事务。订阅者应挂在下游 sink transport 上，使用 `event.WithGroup("...")`，并依赖 Inbox 中间件去重。

事件携带的 `DeleteReason`：

| Reason | Wire 值 | 含义 |
| --- | --- | --- |
| `DeleteReasonReplaced` | `replaced` | `uploaded_file` 字段被新值覆盖 |
| `DeleteReasonDeleted` | `deleted` | 拥有该 key 的业务行被删 |
| `DeleteReasonClaimExpired` | `claim_expired` | pending claim 超时（仅框架内 sweeper 使用） |

dead-letter 事件携带清洗后的 `lastError` 分类，而不是原始后端错误。
当前取值是 `access_denied`、`bucket_not_found`、`session_not_found` 和
`transient`。

公开 supporting API：

| API 组 | 公开 surface |
| --- | --- |
| 事件构造器 | `EventTypeFileClaimed`, `EventTypeFileDeleted`, `EventTypeDeleteDeadLetter`, `NewFileClaimedEvent`, `NewFileDeletedEvent`, `NewDeleteDeadLetterEvent` |
| 门面构造器 | `NewFiles`, `NewFilesFor`, `MultipartFor` |
| 生命周期服务 | `ClaimConsumer`, `DeleteEnqueuer`, `Files`, `FilesFor[T]` |
| 存储接口 | `Service`, `Multipart`, `FileACL`, `URLKeyMapper` |
| URL mapper | `DefaultFileACL`, `IdentityURLKeyMapper`, `ProxyURLKeyMapper`, `DefaultProxyPrefix` |
| metadata helper | `CanonicalizeMetadataKeys` |
| option struct | `PutObjectOptions`, `GetObjectOptions`, `DeleteObjectOptions`, `DeleteObjectsOptions`, `CopyObjectOptions`, `StatObjectOptions`, `InitMultipartOptions`, `PutPartOptions`, `CompleteMultipartOptions`, `AbortMultipartOptions` |
| 结果结构 | `ObjectInfo`, `MultipartSession`, `PartInfo`, `CompletedPart`, `FileRef` |
| meta 常量 | `MetaType`, `MetaTypeUploadedFile`, `MetaTypeRichText`, `MetaTypeMarkdown` |

`storage.CanonicalizeMetadataKeys(m)` 会返回一个新 metadata map，key 使用
S3/HTTP-header canonical form，例如 `author` 会变成 `Author`；nil 或空输入返回
nil。每个 backend 都在 store boundary 应用这个 helper，因此 metadata 会以统一的
provider-neutral 形状 round-trip。

## 错误 Sentinel

| 错误 | 触发 |
| --- | --- |
| `storage.ErrInvalidFileKey` | stat / 下载请求中 key 不合法 |
| `storage.ErrFileNotFound` | 后端找不到对象 |
| `storage.ErrFailedToGetFile` | 后端读失败 |
| `storage.ErrUploadSessionNotFound` | multipart session 已关闭或不存在 |
| `storage.ErrPartTooSmall` | 非最后一片小于 `PartSize()` |
| `storage.ErrPartETagMismatch` | completion 时记录的 part ETag 与后端状态不匹配 |
| `storage.ErrPartNumberOutOfRange` | parts 没有覆盖 `1..N` 连续区间 |
| `storage.ErrClaimNotFound` | `Consume` 引用的 claim 不存在或归属其他 principal |
| `storage.ErrAccessDenied` | 生命周期方法接收到 nil / 匿名 principal |

上传 API 错误也暴露对应的 `ErrCode*` 常量，编号范围是 `2200-2299`：
`ErrCodeInvalidFileKey`, `ErrCodeFileNotFound`,
`ErrCodeFailedToGetFile`, `ErrCodeClaimNotPending`, `ErrCodeClaimExpired`,
`ErrCodeUploadSizeExceedsLimit`, `ErrCodeMultipartNotSupported`,
`ErrCodePublicUploadsNotAllowed`, `ErrCodeUploadTooManyParts`,
`ErrCodeTooManyPendingUploads`, `ErrCodeUploadRequiresMultipart`,
`ErrCodeUploadRequiresFile`, `ErrCodeClaimNotMultipart`,
`ErrCodeUploadPartNumberOutOfRange`, `ErrCodeUploadPartTooLarge`,
`ErrCodeUploadPartTooSmall`, `ErrCodeUploadPartsIncomplete`,
`ErrCodeUploadObjectNotFound`, `ErrCodeUploadSizeMismatch` 和
`ErrCodeAbortFailed`。其他公开 sentinel 还包括 `ErrClaimNotPending`、`ErrClaimExpired`、
`ErrUploadSizeExceedsLimit`、`ErrMultipartNotSupported`、
`ErrPublicUploadsNotAllowed`、`ErrUploadTooManyParts`、
`ErrTooManyPendingUploads`、`ErrUploadRequiresMultipart`、
`ErrUploadRequiresFile`、`ErrClaimNotMultipart`、
`ErrUploadPartNumberOutOfRange`、`ErrUploadPartTooLarge`、
`ErrUploadPartTooSmall`、`ErrUploadPartsIncomplete`、
`ErrUploadObjectNotFound`、`ErrUploadSizeMismatch`、`ErrAbortFailed`、
`ErrBucketNotFound`、`ErrObjectNotFound` 和 `ErrInvalidBucketName`。

## 最小服务示例

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

## CRUD 集成模式

对于带有 `meta` 文件字段的模型，建议用类型化 hook 集成 `FilesFor[T]`：

```go
filesUser := storage.NewFilesFor[User](filesFacade)

create := crud.NewCreate[User, UserParams]().
    AfterTx(func(ctx context.Context, tx orm.DB, principal *security.Principal, model *User) error {
        return filesUser.OnCreate(ctx, tx, principal, model)
    })
```

泛型 CRUD 默认写路径已经接好了 `FilesFor[T]`（见 [Hooks](../data-access/hooks)），自定义写路径照此写就行。

## 实践建议

- 依赖 `storage.Service` 和 `storage.Multipart`，不依赖具体 provider 类型。
- 所有 `Files` / `FilesFor[T]` 调用都放进业务事务里 —— 这就是这个门面的意义。
- 把未确认的对象当作隔离状态：claim sweeper 会最终回收，绕开 claim 直接调 `PutObject` 等于绕开生命周期。
- 一旦开始存私有文件，请实现一个真正的 `FileACL`，默认实现拒绝所有 `priv/*` 读。
- 把 `vef.storage.delete.dead_letter` 接入 ops 看板 —— 队列行已经退役，事件里带着排查所需的信息。
- 模块使用的 extension group 名称是 `vef:api:resources` 和 `vef:app:middlewares`；需要替换 URL 映射时使用 `vef.SupplyURLKeyMapper(...)`。

## 下一步

参考 [自定义 Handler](../building-apis/custom-handlers) 把直接调用 `storage.Service` 与业务流程结合，或读 [事件总线](./event-bus) 了解生命周期事件背后的 outbox transport。
