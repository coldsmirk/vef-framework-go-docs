---
sidebar_position: 4
---

# File Storage

VEF includes a storage abstraction, configurable providers, a built-in resource, a proxy middleware, and automatic file promotion helpers for CRUD flows.

## Supported Providers

The storage module currently supports:

| Provider value | Meaning |
| --- | --- |
| `memory` | in-memory storage |
| `filesystem` | local filesystem storage |
| `minio` | MinIO-backed object storage |

If you do not configure a provider, the module defaults to `memory`.

## `storage.Service` Interface

Application code works against `storage.Service`, not a provider-specific implementation.

The full service surface includes:

| Method | Purpose |
| --- | --- |
| `PutObject(ctx, opts)` | upload one object |
| `GetObject(ctx, opts)` | read one object |
| `DeleteObject(ctx, opts)` | delete one object |
| `DeleteObjects(ctx, opts)` | batch delete |
| `ListObjects(ctx, opts)` | list objects |
| `GetPresignedURL(ctx, opts)` | generate a temporary presigned URL |
| `CopyObject(ctx, opts)` | copy an object |
| `MoveObject(ctx, opts)` | move an object |
| `StatObject(ctx, opts)` | inspect object metadata |
| `PromoteObject(ctx, tempKey)` | move a `temp/` object into permanent storage |

## Option Types

The storage package exposes option structs for each operation:

| Option type | Used by |
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

## Built-In Resource: `sys/storage`

The storage module also registers an RPC resource:

| Resource |
| --- |
| `sys/storage` |

Current actions:

| Action | Purpose |
| --- | --- |
| `upload` | upload one file |
| `get_presigned_url` | generate a presigned URL |
| `delete_temp` | delete a temporary object only |
| `stat` | inspect object metadata |
| `list` | list objects |

See [Built-in Resources](../reference/built-in-resources) for the exact action contracts.

## Temporary Upload Model

The built-in upload path generates temporary object keys under the `temp/` prefix using a date-partitioned layout.

Related public constants:

| Constant | Meaning |
| --- | --- |
| `storage.TempPrefix` | `temp/` prefix for temporary uploads |
| `storage.MetadataKeyOriginalFilename` | metadata key storing the original uploaded filename |

## File Promotion

The public promotion abstraction is:

| Type | Purpose |
| --- | --- |
| `storage.Promoter[T]` | promote and clean up file references across CRUD lifecycle transitions |

Promoter scenario matrix:

| Call pattern | Meaning |
| --- | --- |
| `newModel != nil && oldModel == nil` | create flow |
| `newModel != nil && oldModel != nil` | update flow |
| `newModel == nil && oldModel != nil` | delete flow |

Supported metadata field types:

| Meta type | Meaning |
| --- | --- |
| `uploaded_file` | direct file field |
| `richtext` | HTML content containing resource references |
| `markdown` | Markdown content containing resource references |

## Storage Events

The storage module publishes these file lifecycle events:

| Event type | Meaning |
| --- | --- |
| `vef.storage.file.promoted` | a file was promoted from temporary to permanent storage |
| `vef.storage.file.deleted` | a file was deleted |

## Storage Proxy Middleware

The storage module contributes an app-level download proxy route:

| Route |
| --- |
| `/storage/files/<key>` |

Important behavior:

- this route is normal app middleware, not an RPC action
- it does not automatically inherit API Bearer authentication
- it URL-decodes the key path, fetches the object, sets `Content-Type`, emits cache headers, and streams the file

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
    Key:         "avatars/user-1001.txt",
    Reader:      strings.NewReader("demo"),
    Size:        int64(len("demo")),
    ContentType: "text/plain",
  })

  return err
}
```

## Practical Advice

- use `storage.Service` in application code, not provider-specific implementations
- treat `temp/` uploads as a staging area, not as a permanent object location
- document the `/storage/files/<key>` route separately from RPC resources
- if your model carries file references, prefer integrating through `storage.NewPromoter(...)` instead of ad hoc cleanup logic

## Next Step

Read [Custom Handlers](../guide/custom-handlers) if you want to combine direct `storage.Service` usage with business-specific upload workflows.
