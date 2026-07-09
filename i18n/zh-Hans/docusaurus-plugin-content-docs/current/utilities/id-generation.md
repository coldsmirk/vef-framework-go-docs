---
sidebar_position: 1
---

# ID 生成

`id` 包提供可插拔的唯一标识符生成。框架内置 3 种策略：XID、UUID v7、随机/Nano。

## API 参考

| API | Contract |
| --- | --- |
| `id.IDGenerator` | 所有内置 generator 都实现的 interface |
| `IDGenerator.Generate()` | 返回下一个 string ID；具体格式取决于 generator |
| `id.Generate()` | 委托 `DefaultXIDGenerator.Generate()`，返回 20 字符 XID |
| `id.GenerateUUID()` | 委托 `DefaultUUIDGenerator.Generate()`，返回 UUID v7 string |
| `id.DefaultXIDGenerator` | 由 `NewXIDGenerator()` 创建的包级 XID singleton |
| `id.DefaultUUIDGenerator` | 由 `NewUUIDGenerator()` 创建的包级 UUID v7 singleton |
| `id.NewXIDGenerator()` | 返回一个包装 `xid.New().String()` 的 `IDGenerator` |
| `id.NewUUIDGenerator()` | 返回使用 `uuid.NewV7()` 的 `IDGenerator`；UUID 创建失败时 panic |
| `id.NewRandomIDGenerator(opts...)` | 使用默认值创建 random/Nano-style generator，然后按顺序应用 options |
| `id.RandomIDGeneratorOption` | random generator constructor 使用的 function option 类型 |
| `id.WithAlphabet(alphabet)` | 设置 random generator alphabet |
| `id.WithLength(length)` | 设置 random generator output length |
| `id.DefaultRandomIDGeneratorAlphabet` | 默认 random alphabet：`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ` |
| `id.DefaultRandomIDGeneratorLength` | 默认 random output length：`32` |

## 快速开始

```go
import "github.com/coldsmirk/vef-framework-go/id"

// XID（模型主键的默认值）
xid := id.Generate()
// → "9m4e2mr0ui3e8a215n4g"（20 字符 base32）

// UUID v7（需要 RFC 4122 兼容时使用）
uuid := id.GenerateUUID()
// → "018f4e42-832a-7123-9abc-def012345678"
```

## 内置生成器

### XID

XID 是框架模型主键的默认值——在性能与唯一性之间取得最佳平衡。

| 属性 | 值 |
| --- | --- |
| 格式 | 20 字符 base32 字符串（`0-9, a-v`） |
| 可排序 | 时间序 |
| 全局唯一 | 机器 ID + 计数器 |
| 性能 | 三种策略中最快 |

```go
xid := id.Generate()
// 或
xid := id.DefaultXIDGenerator.Generate()
// 或
xid = id.NewXIDGenerator().Generate()
```

### UUID v7

时间序、符合 RFC 4122 的 UUID —— 对接外部系统时使用。

| 属性 | 值 |
| --- | --- |
| 格式 | 36 字符 UUID（`xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx`） |
| 可排序 | 时间序 |
| RFC 兼容 | RFC 4122 |

```go
uuid := id.GenerateUUID()
// 或
uuid := id.DefaultUUIDGenerator.Generate()
// 或
uuid = id.NewUUIDGenerator().Generate()
```

### 随机 / Nano 风格

可配置字符表的密码学随机 ID —— 适合短小、不透明的令牌。

| 属性 | 值 |
| --- | --- |
| 默认字符表 | `0-9 a-z A-Z`（62 字符，`id.DefaultRandomIDGeneratorAlphabet`） |
| 默认长度 | 32（`id.DefaultRandomIDGeneratorLength`） |

```go
// 默认 32 字符的字母数字 token
gen := id.NewRandomIDGenerator()
token := gen.Generate()

// 自定义：16 位纯数字
gen = id.NewRandomIDGenerator(
    id.WithAlphabet("0123456789"),
    id.WithLength(16),
)
```

`RandomIDGeneratorOption` 是 `WithAlphabet(...)` 和 `WithLength(...)` 使用的
option 类型。

options 会按传给 `NewRandomIDGenerator(...)` 的顺序应用。constructor 不校验
自定义 alphabet 或 length；生成时使用 `go-nanoid/v2` 的 `MustGenerate`，
因此 empty alphabet 或 zero length 会在调用 `Generate()` 时 panic。

## IDGenerator 接口

所有内置生成器都实现同一个接口：

```go
type IDGenerator interface {
    Generate() string
}
```

预构建单例：

```go
id.DefaultXIDGenerator   // IDGenerator
id.DefaultUUIDGenerator  // IDGenerator
```

`orm` 包在插入主键为空的记录时会自动调用 `id.Generate()`（即 XID）。

## 何时使用哪种

| 场景 | 建议 |
| --- | --- |
| 普通应用 ID（主键） | `id.Generate()`（XID） |
| 对外 API 需要 UUID | `id.GenerateUUID()` |
| 短令牌 / 邀请码 / 分享链接 | `id.NewRandomIDGenerator(...)` |
| 自定义格式 | 实现 `IDGenerator` 接口 |
