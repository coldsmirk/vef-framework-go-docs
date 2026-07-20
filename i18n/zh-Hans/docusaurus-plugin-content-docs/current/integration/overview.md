---
sidebar_position: 1
---

# 集成引擎

集成模块是一个由配置和脚本驱动的引擎，用于与外部系统对接——HIS/ERP
厂商、合作方网关、省级平台——而无需把任何厂商的报文格式硬编码进业务代码。
它是可选模块：向 `vef.Run` 传入 `vef.IntegrationModule` 即可启用。

```go
vef.Run(
    vef.IntegrationModule,
    vef.Module("app", ...),
)
```

## 标准模型（Canonical Model）理念

业务代码只面向**契约（Contract）**编程——即你定义一次的标准输入/输出模型。
厂商差异（URL 形态、字段名、信封、签名、码值）全部收敛在按系统维护的
**适配器脚本**里，且可通过管理 API 在运行期编辑。更换厂商意味着编写一个新的
适配器，而不是修改业务代码。

引擎由四张定义表驱动：

| 定义 | 表 | 含义 |
| --- | --- | --- |
| `integration.Contract` | `itg_contract` | 一个标准操作：code、名称、宿主自有 `labels`、可选的输入/输出 JSON Schema（draft 2020-12） |
| `integration.System` | `itg_system` | 一个外部系统实例：基础 URL、出站/入站认证、可选出站信封、可选直连数据库、参数、超时、重试策略 |
| `integration.Adapter` | `itg_adapter` | 将一个系统绑定到一个契约的某个 `direction`（`outbound` / `inbound`），并携带翻译脚本 |
| `integration.Route` | `itg_route` | 将路由键（租户、分支机构、院区）映射到为契约提供服务的系统 |

两个流向共享这些定义：

- **出站（Outbound）** —— 业务代码调用 `integration.Invoker.Invoke`；引擎解析
  目标系统、运行出站适配器脚本、返回经 Schema 校验的标准模型。见
  [出站调用](./outbound)。
- **入站（Inbound）** —— 外部系统调用
  `POST /integration/inbound/:systemCode/:contractCode`；网关验证调用方身份，
  运行入站适配器脚本，脚本把标准输入分发给你注册的
  `integration.InboundHandler`，并按厂商期望的格式组装应答。见
  [入站投递](./inbound)。

值级差异（性别码、状态枚举）由按系统维护的[码值映射](./code-maps)完成翻译。

## 契约（Contract）

```go
type Contract struct {
    Code         string            // 业务代码调用的唯一编码
    Name         string
    Description  *string
    Labels       map[string]string // 宿主自有的筛选元数据，支持相等过滤
    InputSchema  json.RawMessage   // JSON Schema；为空则跳过输入校验
    OutputSchema json.RawMessage   // JSON Schema；为空则跳过输出校验
    IsEnabled    bool
}
```

Schema 必须是自包含的 JSON Schema 文档（draft 2020-12），在保存时编译
（`ErrInvalidSchema`），坏掉的契约永远不会进入调用。两个流向中，输入都在
适配器脚本运行前校验、脚本返回值都在运行后校验——适配器不可能把不符合契约的
模型交给业务代码。

`Labels` 由引擎存储和过滤，但从不解释。校验走共享的 `orm.ValidateLabels`
（与审批流程 labels 同一套规则）：键为字母数字加内部 `-`/`_`（不允许点号——
相等过滤中会被解析为 JSON 路径嵌套），最长 63 字符；值最长 256 字符，允许
空值（作为存在性标记）。违规返回 `ErrInvalidLabel`。

## 系统（System）

```go
type System struct {
    Code             string
    Name             string
    BaseURL          string                  // 启用作用域 http 库
    OutboundAuth     *OutboundAuthConfig     // nil 表示请求不带认证发出
    OutboundEnvelope *OutboundEnvelopeConfig // nil 表示适配器请求原样发出
    InboundAuth      *InboundAuthConfig      // nil 表示完全拒绝入站投递
    DataSource       *DataSourceConfig       // 启用作用域 sql 库
    Params           map[string]string       // 非敏感参数，脚本经 system.params 可见
    TimeoutMs        int                     // 单次 HTTP 调用上限；0 使用框架默认
    Retry            *RetryPolicy            // 基于 httpx 默认重试策略
    IsEnabled        bool
}
```

- 配置了 `BaseURL` 的系统为其适配器脚本提供作用域 `http` 客户端；配置了
  `DataSource` 的系统提供作用域 `sql` 库（默认只读，除非
  `dataSource.mode = "read_write"`）；一个系统可以同时具备两者。
- `Retry` 只对幂等方法在传输错误和 429/502/503/504 响应时重试：
  `maxAttempts`（总尝试次数，含首个调用）、`initialBackoffMs`、
  `maxBackoffMs`。
- `TimeoutMs` 约束每次 HTTP 调用；适配器上的 `timeoutMs` 约束整个脚本运行
  （另一个维度，为 0 时继承 `vef.integration.run_timeout`）。

### 静态存储的密文

敏感认证参数值与数据源密码使用 `vef.integration.secret_key` 中的密钥加密
（base64；默认 AES-GCM，配置 `vef.integration.secret_algorithm = "sm4"` 使用
SM4-GCM）。管理 API 响应始终将其掩码为 `"******"`
（`integration.MaskedSecret`）；更新时提交该占位符将保持已存值不变。
不设置 `secret_key` 时以明文存储并在启动时打印警告。用一种算法封存的值无法
在另一种算法下读取——切换算法需要重新录入已存密文。

## 适配器（Adapter）

```go
type Adapter struct {
    SystemID   string
    ContractID string
    Direction  Direction // "outbound"（默认）或 "inbound"
    Script     string    // 保存时做编译检查
    TimeoutMs  int       // 脚本运行超时；0 继承 vef.integration.run_timeout
    IsEnabled  bool
}
```

一个系统对一个契约在每个方向上恰好用一个适配器实现。两个方向的脚本环境不同，
分别见[出站调用](./outbound#适配器脚本环境)与
[入站投递](./inbound#适配器脚本环境)。

## 路由（Route）

```go
type Route struct {
    RouteKey   string // "" 为默认路由
    ContractID string // "" 适用于所有契约
    SystemID   string
    IsEnabled  bool
}
```

`integration.RouteResolver` 将 `(contract, routeKey)` 解析为系统编码。框架
默认实现读取 `itg_route` 表：精确 `(key, contract)` 规则优先于契约通配规则，
空键是默认路由。当路由信息在别处维护时（租户注册表、配置中心），可通过
`fx.Decorate` 整体替换。没有任何规则匹配的键返回 `ErrRouteNotFound`。

路由健康状况可通过 `diagnose_routes` 操作在运行期检查（悬挂的适配器、被禁用
的目标、未覆盖的契约）；见 [RPC 资源](./resources#integrationops)。

## 配置

```toml
[vef.integration]
auto_migrate = true          # 启动时执行集成模块 DDL 迁移
secret_key = "base64-key"    # 敏感值静态加密密钥
secret_algorithm = "aes"     # "aes"（AES-GCM，默认）或 "sm4"（SM4-GCM）
run_timeout = "30s"          # 单次脚本运行上限，含线上调用
max_response_body = 8388608  # 脚本读取的单个 HTTP 响应体上限（8 MiB）

[vef.integration.log]
mode = "errors"              # "off" | "errors"（默认）| "all"
capture_limit = 4096         # 每个捕获负载截断前的字节数
mask_fields = ["idCard"]     # 捕获中额外掩码的 JSON 字段名
retention = "720h"           # 调用日志保留窗口；0 表示永久保留

[vef.integration.inbound.rate_limit]
max = 120                    # 每窗口每 (系统, 客户端 IP) 允许的投递数
period = "1m"
```

调用日志写入 `itg_invocation_log`，包含失败分类、耗时、输入/输出捕获与完整
HTTP 线上追踪——按 `vef.integration.log` 掩码（凭证头始终掩码；`mask_fields`
额外掩码）并截断到 `capture_limit`。设置 `retention` 后，每小时的清扫任务会
删除超过窗口的行。

## 失败分类

`integration.FailureKind` 是调用日志、统计与 API 错误共享的唯一失败词汇表；
空值表示成功。

| 类别 | 含义 |
| --- | --- |
| `input_invalid` | 输入在脚本运行前被契约输入 Schema 拒绝 |
| `output_invalid` | 脚本返回值被输出 Schema 拒绝 |
| `upstream` | 外部系统自身报告的失败（`errors.upstream(...)`） |
| `transport` | 线上调用未完成（连接拒绝、TLS 失败） |
| `timeout` | 调用超过运行超时 |
| `canceled` | 调用方取消了调用 |
| `script` | 未捕获的脚本异常或编译错误——适配器缺陷 |
| `config` | 认证 scheme 未注册、凭证无法解密等配置故障 |
| `auth` | 入站投递被入站认证验证拒绝 |
| `handler` | 入站业务处理器在分发成功后返回错误——业务失败 |

## 统计

Invoker 实现 `integration.StatsInspector`：进程启动以来按
`(系统, 契约, 方向)` 维度的本节点计数——调用数、成功数、按类别的失败数、
平均/最大耗时、最近错误。监控模块通过 `sys/monitor.get_integration_stats`
读取（见[内置资源](../reference/built-in-resources)）。每个
`InvocationStats` 条目包含 `system`、`contract`、`direction`、`calls`、
`successes`、`failures`（`FailureKind` 到计数的映射）、`avgDurationMs`、
`maxDurationMs`、`lastError`、`lastErrorAt`。被验证拒绝的入站投递聚合在空
`contract` 下——拒绝时契约编码还只是未验证的调用方输入。

## 错误码

集成 API 错误使用响应码 `2600`–`2699`，除注明外均以 HTTP 200 承载、失败由
响应体 code 表达。

| 码 | 错误 | 含义 |
| --- | --- | --- |
| `2600` | `ErrContractNotFound` | 契约不存在 |
| `2601` | `ErrContractDisabled` | 契约被禁用 |
| `2602` | `ErrSystemNotFound` | 系统不存在 |
| `2603` | `ErrSystemDisabled` | 系统被禁用 |
| `2604` | `ErrAdapterNotFound` | 该方向上没有适配器将系统绑定到契约 |
| `2605` | `ErrAdapterDisabled` | 适配器被禁用 |
| `2606` | `ErrRouteNotFound` | 路由键未匹配任何规则 |
| `2607` | `ErrTargetAmbiguous` | 同时传入了 `WithSystem` 与 `WithRoute` |
| `2608` | `ErrInputInvalid(detail)` | 输入被输入 Schema 拒绝 |
| `2609` | `ErrOutputInvalid(detail)` | 脚本返回值被输出 Schema 拒绝 |
| `2610` | `ErrUpstreamFailed(message)` | 外部系统报告的失败 |
| `2611` | `ErrTransportFailed` | 线上调用未完成 |
| `2612` | `ErrInvocationTimeout` | 超过运行超时 |
| `2613` | `ErrScriptFailed(detail)` | 脚本抛出异常或编译失败 |
| `2614` | `ErrUnknownAuthScheme(scheme)` | 系统引用了未注册的认证 scheme |
| `2615` | `ErrInvalidSchema(detail)` | 契约 Schema 保存时被拒绝 |
| `2616` | `ErrInvalidScript(detail)` | 适配器脚本保存时被拒绝 |
| `2617` | `ErrInvalidAuthParams(detail)` | 认证配置被其 scheme 拒绝 |
| `2618` | `ErrInvalidRouteRef` | 路由引用了不存在的契约或系统 |
| `2619` | `ErrInvalidBaseURL` | 系统基础 URL 不是绝对 URL |
| `2620` | `ErrInvalidDataSource(detail)` | 系统数据源不完整或凭证无法处理 |
| `2621` | `ErrInvalidDirection` | 适配器方向不在已知流向内 |
| `2622` | `ErrInboundAuthFailed` | 入站投递验证失败（HTTP 401，刻意统一） |
| `2623` | `ErrInboundHandlerMissing` | 入站契约没有注册处理器（HTTP 501） |
| `2624` | `ErrInvocationCanceled` | 调用方取消了调用 |
| `2625` | `ErrInvalidEnvelope(detail)` | 出站信封配置保存时被拒绝 |
| `2626` | `ErrInvalidLabel` | 契约 label 键/值校验失败 |
| `2627` | `ErrMissingCodeMap(codeSet)` | 目标系统没有该码值集的已启用映射 |
| `2628` | `ErrUnmappedValue(codeSet, value)` | reject 策略下的未映射值 |
| `2629` | `ErrInvalidCodeMap(detail)` | 码值映射定义保存时被拒绝 |

## 下一步

- [出站调用](./outbound) —— `Invoker`、适配器脚本环境、认证 scheme、信封
- [入站投递](./inbound) —— HTTP 网关、验证 scheme、业务处理器
- [码值映射](./code-maps) —— 按系统的值翻译与宿主码值目录
- [RPC 资源](./resources) —— 管理 API 逐字段参考
