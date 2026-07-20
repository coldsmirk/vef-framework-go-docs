---
sidebar_position: 2
---

# 出站调用

出站是业务代码调用外部系统的流向：你调用一个契约，引擎解析出提供服务的
系统、运行其出站适配器脚本、返回经 Schema 校验的标准模型。

## 调用契约

```go
type PatientService struct {
    invoker integration.Invoker
}

func (s *PatientService) LoadPatient(ctx context.Context, hospitalArea, patientID string) (*Patient, error) {
    // 类型化便捷封装：把标准模型解码为 Patient。
    return integration.Call[*Patient](ctx, s.invoker, "patient.get",
        map[string]any{"patientId": patientID},
        integration.WithRoute(hospitalArea),
    )
}
```

启用 `vef.IntegrationModule` 后，`integration.Invoker` 即可从 DI 注入：

| API | 契约 |
| --- | --- |
| `Invoke(ctx, contract string, input any, opts ...InvokeOption) (*Result, error)` | 按契约输入 Schema 校验输入、运行适配器脚本、校验输出，返回 `Result` |
| `integration.Call[T](ctx, inv, contract, input, opts...) (T, error)` | `Invoke` 的类型化封装，经 JSON 往返把输出解码为 `T` |

### 调用选项

| 选项 | 行为 |
| --- | --- |
| `integration.WithSystem(code)` | 直接指定目标系统，绕过路由解析。与 `WithRoute` 互斥（`ErrTargetAmbiguous`） |
| `integration.WithRoute(key)` | 通过 `RouteResolver` 选择系统；不传目标选项则解析空路由键（默认路由） |
| `integration.WithTimeout(d)` | 覆盖本次调用的脚本运行超时 |
| `integration.WithCache(ttl)` | 将校验后的输出按（系统、契约、输入）为键缓存 `ttl`。默认关闭——只在业务能容忍该时效的数据处启用。缓存命中不记调用日志与统计 |

### Result

| 方法 | 含义 |
| --- | --- |
| `Output() any` | 标准模型，已通过输出 Schema 校验 |
| `Decode(v any) error` | 经 JSON 往返把输出解码到 `v` |
| `System() string` | 实际提供服务的系统编码 |
| `Duration() time.Duration` | 调用耗时 |
| `Cached() bool` | 输出是否来自响应缓存 |

## 适配器脚本环境

出站适配器脚本把契约输入翻译成对外部系统的调用，并把其响应翻译回契约输出。
脚本最后一个表达式的值即输出（会经契约输出 Schema 校验）。

所有脚本都能看到 JS 引擎基线（`BigNumber`、`dayjs`、`fxp`、`radashi`、`z`、
`URL` / `URLSearchParams` —— 见 [JS 引擎](../data-tools/js-engine)），外加以下
绑定：

| 绑定 | 何时存在 | 内容 |
| --- | --- | --- |
| `input` | 始终 | 经 Schema 校验的契约输入 |
| `system` | 始终 | 系统只读视图：`{ code, name, params }` —— 永远不含凭证 |
| `errors` | 始终 | 失败分类：`errors.upstream(message)` 抛出的异常被记录为上游失败而非脚本缺陷 |
| `codes` | 始终 | 码值映射翻译：`codes.toExternal` / `codes.toCanonical` / `codes.entries` —— 见[码值映射](./code-maps) |
| `http` | 系统配置了 `baseUrl` 时 | 系统作用域 HTTP 客户端（见下） |
| `sql` | 系统配置了 `dataSource` 时 | 系统作用域 SQL 访问（见下） |

脚本访问未配置的能力会得到普通的 `ReferenceError`。

### 作用域 `http` 客户端

```js
// 路径永远是相对的——客户端锁定在系统 baseUrl 上。
const res = http.get('/api/patients/' + input.patientId);
if (!res.ok) {
    errors.upstream('patient service returned HTTP ' + res.status);
}
const data = res.json();

({ patientId: data.id, name: data.patientName })
```

| 函数 | 行为 |
| --- | --- |
| `http.fetch(path, { method, headers, query, body, timeout, envelope })` | 完整请求形式 |
| `http.get(path, options?)` / `http.delete(path, options?)` | fetch 的语法糖 |
| `http.post(path, body, options?)` / `http.put(...)` / `http.patch(...)` | 带请求体的语法糖 |

- 认证由宿主按系统的 `outboundAuth` 注入；脚本永远看不到凭证，也无法访问
  其他主机（绝对 URL 被拒绝）。
- `body`：字符串和字节数组原样透传；其他值 JSON 编码并隐含
  `application/json` 内容类型。
- `timeout`（毫秒）只能收紧系统的调用超时。
- 每次调用返回 fetch `Response` 形态：`{ status, statusText, ok, url,
  headers, body, text(), json(), arrayBuffer() }`。响应头名小写化，
  多值以 `", "` 连接。
- 失败以可捕获异常抛出；脚本未捕获的传输失败将本次调用分类为 `transport`。
- 每次线上交换都会记入调用追踪（按 `vef.integration.log` 掩码、截断）。

### 作用域 `sql` 库

用于集成面即数据库（外部视图、交换表）的系统：

```js
const rows = sql.queryList('SELECT id, name FROM v_patients WHERE id = ?', input.patientId);
const one  = sql.queryOne('SELECT ... WHERE id = ?', input.id);  // 对象 | null
```

- 只提供占位符绑定——刻意不提供字符串拼接辅助。
- 默认只读，由基于 AST 的守卫 fail-closed 强制：写 CTE、堆叠语句、方言
  特有的副作用函数都被拒绝，无法解析的 SQL 一律拒绝。请编写可移植语句——
  守卫解析器拒绝方言独有语法，如 T-SQL 的 `SELECT TOP`、方括号引用、
  `WITH (NOLOCK)`，以及 Oracle 的 `(+)` 外连接与 `CONNECT BY`
  （`OFFSET/FETCH`、`ROWNUM`、`NVL`、`TO_CHAR` 等可移植写法可以通过）。
- `sql.execute('UPDATE ...', args)`（返回 `{ rowsAffected }`）仅当系统数据源
  声明 `mode = "read_write"` 时可用，否则抛错。
- 结果集有行数上限（超限报错并提示加 `LIMIT`）；无匹配行时返回 `[]` 而非
  `null`。

## 出站认证 Scheme

`system.outboundAuth` 决定框架如何向该系统证明身份：
`{ "scheme": "...", "params": { ... }, "script": "..." }`。scheme 声明的
敏感参数值加密存储，并在管理响应中掩码。

| Scheme | 参数 | 行为 |
| --- | --- | --- |
| `none` | — | 请求不带认证发出 |
| `http_basic` | `username`、`password`（敏感） | RFC 7617 Basic 凭证 |
| `bearer` | `token`（敏感） | 静态 `Authorization: Bearer` 令牌 |
| `header` | 每个凭证头一个条目（全部敏感） | 把每个配置对作为静态请求头发送 |
| `query` | 每个凭证参数一个条目（全部敏感） | 把每个配置对作为静态查询参数发送 |
| `signature` | `appId`、`secret`（hex，敏感） | 按框架 HMAC 签名约定为每个请求签名：对身份、方法、路径生成 `x-timestamp` / `x-nonce` / `x-signature` —— 两个 VEF 部署仅凭配置即可互相认证 |
| `script` | 自由参数（全部敏感）+ 签名脚本体 | 每个请求在零 IO 运行时中执行自定义签名体；可见已构建的 `request`（`{ method, url, path, query, headers, body }`）与解密后的 `params`，返回要添加的凭证头对象 |

应用可通过 `vef.ProvideIntegrationOutboundAuthScheme` 注册自定义 scheme；
与内置同名的 scheme 会替换内置实现。自定义 scheme 实现：

```go
type OutboundAuthScheme interface {
    Name() string
    Apply(cfg *OutboundAuthConfig) ([]httpx.Option, error)
    SensitiveParams() []string // integration.SensitiveAll 表示所有参数敏感
}
```

## 系统信封

多数厂商 API 在每个端点上重复同一种报文结构（`{code, msg, data}` 响应、
带签名的请求包装、SOAP 信封）。在系统级 `outboundEnvelope` 配置一次，而不是
在每个适配器里重复：

```js
// outboundEnvelope.request —— 包装适配器发出的请求。
// `request` = { method, path, headers, query, body }；返回对象省略的字段
// 保持适配器原值。
({ body: { reqData: request.body, ts: Date.now() } })
```

```js
// outboundEnvelope.response —— 解包完成的 HTTP 响应。
// `response` 是 fetch Response 形态；该脚本的返回值即适配器调用的所得。
const parsed = response.json();
if (parsed.code !== '0') {
    errors.upstream(parsed.msg);   // 厂商级错误在系统层面统一分类一次
}
parsed.data
```

两个脚本可以只配置其一（保存时校验至少存在一个），单次调用可用
`http.get(path, { envelope: false })` 绕过——适用于文件下载、健康检查等
异形端点。

## 失败分类与重试

- 幂等方法上的传输失败与 429/502/503/504 响应会按系统 `retry` 策略先重试再
  上浮。
- 失败词汇表（`input_invalid`、`output_invalid`、`upstream`、`transport`、
  `timeout`、`canceled`、`script`、`config`）与调用日志、统计共享；见
  [概览](./overview#失败分类)。
- 调用方取消与超时是区分开的：调用方离场分类为 `canceled`，超过期限分类为
  `timeout`。

## 先测试再启用

`integration/ops` 资源提供的测试台对已保存和未保存的定义都有效——测试先于
启用：

- `dry_run` 在契约下对系统执行脚本（可以是编辑器里未保存的内容），返回输出、
  失败分类与完整线上追踪。它发出的调用是真实的；不记入统计与调用日志。
- `test_connection` 探测已保存系统配置的每种传输（HTTP 基础 URL、数据库），
  以数据形式报告可达性。

两者的逐字段请求/响应文档见 [RPC 资源](./resources#integrationops)。

## 下一步

[入站投递](./inbound)介绍相反的流向——外部系统调入你的应用。
