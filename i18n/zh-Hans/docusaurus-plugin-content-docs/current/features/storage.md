---
sidebar_position: 4
---

# 文件存储

VEF 提供一套与 provider 无关的存储抽象、三种内置 provider、分片上传协议、用于 CRUD 流程同步模型与后端文件引用的类型化生命周期门面，以及用于下游清理的事务性 outbox。

> 存储模块在 v0.21 到 v0.26 期间经历了大量重构：原来的 `Promoter[T]` 已被 `Files` / `FilesFor[T]` 取代，上传协议统一为分片 multipart，引入 claim / queue 生命周期，把 principal 授权贯穿全流程，对外的 `Consume` / `Enqueue` API 也重命名精简。本页对应 v0.26 状态——老版本快照与现在已经不再兼容。

## 支持的 Provider

| 配置值 | 后端 |
| --- | --- |
| `memory` | 进程内 map；测试和短期演示 |
| `filesystem` | 本地文件系统 |
| `minio` | MinIO / S3 兼容对象存储 |

`storage.provider` 选择后端。未配置时默认 `memory`。

## `storage.Service` 接口

应用代码依赖 `storage.Service`，不依赖任何 provider 特定类型：

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

每个方法都使用对应的 option struct：`PutObjectOptions`、`GetObjectOptions`、`DeleteObjectOptions`、`DeleteObjectsOptions`、`CopyObjectOptions`、`StatObjectOptions`。框架刻意不支持位置参数 ——这样后续追加字段始终是可叠加变更。

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

> `sys/storage.list_parts` RPC action 用于让客户端恢复上传中断，它由框架的 part-store 表服务，**不**是 `Multipart.ListParts` 方法 —— 后端接口本身只暴露上面 5 个方法。

## 内置资源：`sys/storage`

存储模块注册了一个 RPC 资源，包含分片上传动作：

| Action | 作用 |
| --- | --- |
| `init_upload` | 开启一个新的 multipart session（返回不透明 `uploadId` 以及协商出的 `partSize`） |
| `upload_part` | 上传一片 |
| `list_parts` | 列出已上传的 parts |
| `complete_upload` | 用最终 part 清单完成上传 |
| `abort_upload` | 中止并释放 session |

下载通过下方代理中间件提供。

## 可见性前缀

object key 通过前缀表达可见意图：

| 常量 | 值 | 含义 |
| --- | --- | --- |
| `storage.PublicPrefix` | `pub/` | 世界可读；默认 ACL 直接放行 |
| `storage.PrivatePrefix` | `priv/` | 由业务侧 `FileACL` 控制 |

存储资源会根据上传时的 `public` 标志，把 key 落在 `pub/` 或 `priv/` 下。这只是约定，由 `FileACL` 强制——存储后端本身不做这种判断。

## FileACL

`storage.FileACL` 决定 principal 能否读取一个私有 key。

```go
type FileACL interface {
    CanRead(ctx context.Context, principal *security.Principal, key string) (bool, error)
}
```

默认实现 `storage.DefaultFileACL` 只允许 `pub/` 下的读。业务代码通过 `vef.SupplyFileACL(...)` 注入自己实现，按业务的归属表判断。默认安全——没有重写时不会暴露任何 `priv/*` 给认证调用方。

## 存储代理中间件

模块在 app 层挂了一个下载路由：

```
GET /storage/files/<key>
```

行为：

- 不是 RPC，是 app 中间件，不走 API engine
- 先向 `FileACL.CanRead` 询问授权
- 对 key 进行 URL 解码，从后端读取对象，设置 `Content-Type`、缓存 header 后流式输出

## 上传 Claim 与 Pending Delete（生命周期）

`complete_upload` 成功后，框架会写入一个属于上传调用方的 `upload_claim` 行。在业务模型真正引用（通过 `Files.OnCreate` / `OnUpdate`）这个 key 之前，对象处于隔离状态——周期性 sweeper 会把过期 claim 入队为异步删除（`DeleteReasonClaimExpired`）。

业务写入因此分裂为两组事务性表面：

- **Claim consumer**：在业务 insert 同一事务中删除 `upload_claim` 行。
- **Delete enqueuer**：把被替换的字段值、被删除的业务行对应的对象，写入 `pending_delete` 行做异步回收。

后台 `DeleteWorker` 会按重试/退避策略把 `pending_delete` 行打到后端；用尽重试的行会被 park 起来供人工排查。删除成功发出 `vef.storage.file.deleted`，park 的行发出 `vef.storage.delete.dead_letter`。

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

`URLKeyMapper` 用来把富文本/Markdown 中的 URL 翻译成存储 key。如果业务里直接嵌入 bare key，传入 `&storage.IdentityURLKeyMapper{}`（或 `nil`，会被规整为 identity）即可。

## 存储事件

| 事件类型 | 触发 |
| --- | --- |
| `vef.storage.file.claimed`（`storage.FileClaimedEvent`） | 某次业务事务采纳了一个之前 pending 的 claim（`Files.OnCreate` 或 update 新侧） |
| `vef.storage.file.deleted`（`storage.FileDeletedEvent`） | delete worker 成功把对象从后端删除 |
| `vef.storage.delete.dead_letter`（`storage.DeleteDeadLetterEvent`） | delete worker 重试用尽；该行被 park，未删除 |

三者都通过 outbox transport 在生产事务里发布。订阅者必须用 `event.WithGroup("...")` 挂上，并依赖 Inbox 中间件去重。

事件携带的 `DeleteReason`：

| Reason | 含义 |
| --- | --- |
| `DeleteReasonReplaced` | `uploaded_file` 字段被新值覆盖 |
| `DeleteReasonDeleted` | 拥有该 key 的业务行被删 |
| `DeleteReasonClaimExpired` | pending claim 超时（仅框架内 sweeper 使用） |

## 错误 Sentinel

| 错误 | 触发 |
| --- | --- |
| `storage.ErrInvalidFileKey` | stat / 下载请求中 key 不合法 |
| `storage.ErrFileNotFound` | 后端找不到对象 |
| `storage.ErrFailedToGetFile` | 后端读失败 |
| `storage.ErrUploadSessionNotFound` | multipart `UploadID` 已关闭或不存在 |
| `storage.ErrPartTooSmall` | 非最后一片小于 `PartSize()` |
| `storage.ErrPartETagMismatch` | `complete_upload` 中 part 列表 ETag 不匹配 |
| `storage.ErrPartNumberOutOfRange` | parts 没有覆盖 `1..N` 连续区间 |
| `storage.ErrClaimNotFound` | `Consume` 引用的 claim 不存在或归属其他 principal |
| `storage.ErrAccessDenied` | 生命周期方法接收到 nil / 匿名 principal |

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

泛型 CRUD 默认写路径已经接好了 `FilesFor[T]`（见 [Hooks](../guide/hooks)），自定义写路径照此写就行。

## 实践建议

- 依赖 `storage.Service` 和 `storage.Multipart`，不依赖具体 provider 类型。
- 所有 `Files` / `FilesFor[T]` 调用都放进业务事务里 —— 这就是这个门面的意义。
- 把未确认的对象当作隔离状态：claim sweeper 会最终回收，绕开 claim 直接调 `PutObject` 等于绕开生命周期。
- 一旦开始存私有文件，请实现一个真正的 `FileACL`，默认实现拒绝所有 `priv/*` 读。
- 把 `vef.storage.delete.dead_letter` 接入 ops 看板 —— 那些行需要人工介入。

## 下一步

参考 [自定义 Handler](../guide/custom-handlers) 把直接调用 `storage.Service` 与业务流程结合，或读 [事件总线](./event-bus) 了解生命周期事件背后的 outbox transport。
