---
sidebar_position: 4
---

# 文件存储

VEF 内置了存储抽象、可配置 provider、内置资源、代理中间件，以及用于 CRUD 流程的自动文件提升能力。

## 支持的 Provider

storage 模块当前支持：

| Provider 值 | 含义 |
| --- | --- |
| `memory` | 内存存储 |
| `filesystem` | 本地文件系统存储 |
| `minio` | MinIO 对象存储 |

如果没有显式配置 provider，模块默认使用 `memory`。

## `storage.Service` 接口

应用代码应依赖 `storage.Service`，而不是某个具体 provider 实现。

完整公开接口包括：

| 方法 | 作用 |
| --- | --- |
| `PutObject(ctx, opts)` | 上传单个对象 |
| `GetObject(ctx, opts)` | 读取单个对象 |
| `DeleteObject(ctx, opts)` | 删除单个对象 |
| `DeleteObjects(ctx, opts)` | 批量删除 |
| `ListObjects(ctx, opts)` | 列出对象 |
| `GetPresignedURL(ctx, opts)` | 生成临时预签名 URL |
| `CopyObject(ctx, opts)` | 复制对象 |
| `MoveObject(ctx, opts)` | 移动对象 |
| `StatObject(ctx, opts)` | 检查对象元数据 |
| `PromoteObject(ctx, tempKey)` | 把 `temp/` 对象提升为永久对象 |

## Option 类型

storage 包为每个操作都暴露了对应 option struct：

| Option 类型 | 对应方法 |
| --- | --- |
| `storage.PutObjectOptions` | `PutObject` |
| `storage.GetObjectOptions` | `GetObject` |
| `storage.DeleteObjectOptions` | `DeleteObject` |
| `storage.DeleteObjectsOptions` | `DeleteObjects` |
| `storage.ListObjectsOptions` | `ListObjects` |
| `storage.PresignedURLOptions` | `GetPresignedURL` |
| `storage.CopyObjectOptions` | `CopyObject` |
| `storage.MoveObjectOptions` | `MoveObject` |
| `storage.StatObjectOptions` | `StatObject` |

## 内置资源：`sys/storage`

storage 模块还会注册一个内置 RPC 资源：

| 资源 |
| --- |
| `sys/storage` |

当前 action：

| Action | 作用 |
| --- | --- |
| `upload` | 上传单个文件 |
| `get_presigned_url` | 生成预签名 URL |
| `delete_temp` | 只允许删除临时对象 |
| `stat` | 查询对象元数据 |
| `list` | 列出对象 |

精确 action 契约见 [内置资源](../reference/built-in-resources)。

## 临时上传模型

内置上传路径会把对象 key 生成为带 `temp/` 前缀的日期分区路径。

相关公共常量：

| 常量 | 含义 |
| --- | --- |
| `storage.TempPrefix` | 临时上传的 `temp/` 前缀 |
| `storage.MetadataKeyOriginalFilename` | 保存原始文件名的 metadata key |

## 文件提升（Promoter）

公开的提升抽象是：

| 类型 | 作用 |
| --- | --- |
| `storage.Promoter[T]` | 在 CRUD 生命周期中自动提升和清理文件引用 |

Promoter 的场景矩阵：

| 调用模式 | 含义 |
| --- | --- |
| `newModel != nil && oldModel == nil` | create 场景 |
| `newModel != nil && oldModel != nil` | update 场景 |
| `newModel == nil && oldModel != nil` | delete 场景 |

支持的 metadata 字段类型：

| Meta 类型 | 含义 |
| --- | --- |
| `uploaded_file` | 直接文件字段 |
| `richtext` | 含资源引用的 HTML 富文本 |
| `markdown` | 含资源引用的 Markdown 文本 |

## 存储事件

storage 模块会发布这些文件生命周期事件：

| 事件类型 | 含义 |
| --- | --- |
| `vef.storage.file.promoted` | 文件从临时存储提升到永久存储 |
| `vef.storage.file.deleted` | 文件被删除 |

## 存储代理中间件

storage 模块还会挂一个 app 级下载代理路由：

| 路由 |
| --- |
| `/storage/files/<key>` |

关键行为：

- 这个路由是普通 app middleware，不是 RPC action
- 它不会自动继承 API Bearer 认证
- 它会对 key 做 URL 解码、读取对象、设置 `Content-Type`、写缓存头，并以流方式输出文件

## 最小 Service 示例

```go
package avatars

import (
  "context"
  "strings"

  "github.com/coldsmirk/vef-framework-go/storage"
)

func SaveAvatar(ctx context.Context, svc storage.Service) error {
  _, err := svc.PutObject(ctx, storage.PutObjectOptions{
    Key:         "avatars/user-1001.txt",
    Reader:      strings.NewReader("demo"),
    Size:        int64(len("demo")),
    ContentType: "text/plain",
  })

  return err
}
```

## 实践建议

- 应用代码里依赖 `storage.Service`，不要直接依赖 provider 实现
- 把 `temp/` 上传当作中间态，而不是永久对象位置
- `/storage/files/<key>` 这条路由应与 RPC 资源分开单独说明
- 如果模型里有文件引用，优先用 `storage.NewPromoter(...)` 接入，而不是自己手写清理逻辑

## 下一步

继续阅读 [自定义处理器](../guide/custom-handlers)，如果你要把 `storage.Service` 和业务化上传流程结合起来，就会接到那里。
