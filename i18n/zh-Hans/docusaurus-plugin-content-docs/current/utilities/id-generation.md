---
sidebar_position: 1
---

# ID 生成

`id` 包提供可插拔的唯一标识符生成。框架内置 4 种策略：XID、UUID v7、Snowflake、随机/Nano。

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
| 性能 | 四种策略中最快 |

```go
xid := id.Generate()
// 或
xid := id.DefaultXIDGenerator.Generate()
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
```

### Snowflake

Twitter 风格 Snowflake ID —— 64 位整数编码为十进制字符串。需要分布式、有序的整数 ID 时使用。

| 属性 | 值 |
| --- | --- |
| 编码 | 自定义：6 位节点（0-63）、12 位步进（每毫秒每节点 4096 个 ID） |
| Epoch | `1754582400000`（包内固化的自定义起点） |
| 节点 ID | 启动时读取 `VEF_NODE_ID` 环境变量 |
| 默认实例 | `id.DefaultSnowflakeIDGenerator` |

```go
snow := id.DefaultSnowflakeIDGenerator.Generate()
// → "7234567890123456789"
```

如需自定义节点 ID，构造一个新生成器：

```go
gen, err := id.NewSnowflakeIDGenerator(int64(42))
if err != nil {
    return err
}
sid := gen.Generate()
```

> Snowflake 最多支持 64 个节点、每毫秒每节点 4096 个 ID。每个进程必须配置唯一的 `VEF_NODE_ID`，避免冲突。

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

## IDGenerator 接口

所有内置生成器都实现同一个接口：

```go
type IDGenerator interface {
    Generate() string
}
```

预构建单例：

```go
id.DefaultXIDGenerator         // IDGenerator
id.DefaultUUIDGenerator        // IDGenerator
id.DefaultSnowflakeIDGenerator // IDGenerator
```

`orm` 包在插入主键为空的记录时会自动调用 `id.Generate()`（即 XID）。

## 何时使用哪种

| 场景 | 建议 |
| --- | --- |
| 普通应用 ID（主键） | `id.Generate()`（XID） |
| 对外 API 需要 UUID | `id.GenerateUUID()` |
| 分布式、有序的整数 ID | `id.DefaultSnowflakeIDGenerator` |
| 短令牌 / 邀请码 / 分享链接 | `id.NewRandomIDGenerator(...)` |
| 自定义格式 | 实现 `IDGenerator` 接口 |
