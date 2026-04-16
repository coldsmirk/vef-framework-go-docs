---
sidebar_position: 1
---

# ID 生成

`id` 包提供可插拔的唯一标识符生成，内置两种策略。

## 快速开始

```go
import "github.com/coldsmirk/vef-framework-go/id"

// 默认：XID（大多数场景推荐）
xid := id.Generate()
// → "9m4e2mr0ui3e8a215n4g"（20 字符，base32 编码）

// UUID v7（需要 RFC 4122 兼容时使用）
uuid := id.GenerateUUID()
// → "018f4e42-832a-7123-9abc-def012345678"（36 字符）
```

## 内置生成器

### XID（默认）

XID 是框架的默认 ID 生成策略，在性能和唯一性之间取得了最佳平衡。

| 属性 | 值 |
| --- | --- |
| 格式 | 20 字符 base32 字符串（`0-9, a-v`）|
| 可排序 | ✅ 基于时间排序 |
| 全局唯一 | ✅ 机器 ID + 计数器 |
| 性能 | 所有策略中最优 |

```go
xid := id.Generate()
```

### UUID v7

UUID v7 提供基于时间的排序，遵循 RFC 4122 标准。

| 属性 | 值 |
| --- | --- |
| 格式 | 36 字符 UUID（`xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx`）|
| 可排序 | ✅ 基于时间排序 |
| RFC 兼容 | ✅ RFC 4122 |
| 使用场景 | 外部系统需要 UUID 时 |

```go
uuid := id.GenerateUUID()
```

## IDGenerator 接口

通过实现 `IDGenerator` 接口来自定义 ID 生成器：

```go
type IDGenerator interface {
    Generate() string
}
```

框架使用默认的 XID 生成器（`id.DefaultXIDGenerator`）来生成模型主键。`orm` 包在插入 ID 为空的记录时会自动调用 `id.Generate()`。

## 预构建生成器实例

```go
id.DefaultXIDGenerator  // *XIDGenerator 单例
id.DefaultUUIDGenerator // *UUIDGenerator 单例
```

## 何时使用哪种

| 场景 | 建议 |
| --- | --- |
| 通用应用 ID | `id.Generate()`（XID）|
| 外部 API 集成 | `id.GenerateUUID()`（UUID v7）|
| 需要自定义格式 | 实现 `IDGenerator` 接口 |
