---
sidebar_position: 3
---

# 入站投递

入站是外部系统调入你应用的流向：HTTP 网关接收厂商请求，按目标系统的入站
认证验证调用方，运行入站适配器脚本——脚本把线上请求翻译成契约的标准输入、
分发给你的业务处理器，并组装厂商期望的应答。

## HTTP 网关

启用 `vef.IntegrationModule` 后，框架注册：

```
POST /integration/inbound/:systemCode/:contractCode
```

- 该端点是真实路由（注册在 API 引擎之后、SPA 回退之前）；非集成流量对它
  零开销。
- 它刻意绕过 `/api` 分发模型：对外应答需要对状态码和响应体的原始控制，
  永远不套标准 result 信封。
- 投递按 `(系统, 客户端 IP)` 在滑动窗口内限流
  （`vef.integration.inbound.rate_limit`，默认每分钟 120 次，按节点独立
  计数）。

网关把每个 HTTP 请求翻译成协议中立的 `integration.InboundRequest` 信封：

| 字段 | 内容 |
| --- | --- |
| `SystemCode` / `ContractCode` | 从 URL 路径参数解析 |
| `Protocol` | `"http"` |
| `Method` / `Path` / `Query` | HTTP 原生请求数据 |
| `Headers` | 头名小写化，多值以 `", "` 连接 |
| `Body` | 原始请求负载 |
| `ClientAddr` | 网络对端地址，用于 IP 验证 |

## 入站认证

验证是 fail-closed 的：没有配置 `inboundAuth` 的系统完全拒绝入站投递——
`none` scheme 是刻意的显式放行。所有验证失败统一返回
`ErrInboundAuthFailed`（HTTP 401）；缺配置、缺凭证、凭证错误对调用方
不可区分，未知或被禁用的系统也以完全相同的方式拒绝，因此系统编码无法被
枚举。

内置 scheme 与出站词汇表互为镜像——同名 scheme 验证其出站对应方发送的
报文格式；`ip` 仅入站可用：

| Scheme | 参数 | 行为 |
| --- | --- | --- |
| `none` | — | 接受所有调用方；请配合网络层控制使用 |
| `ip` | `whitelist`（逗号分隔的 IP/CIDR） | 验证调用方网络地址；空白名单按配置故障处理，而非放行 |
| `http_basic` | `username`、`password`（敏感） | 常数时间验证 RFC 7617 Basic 凭证 |
| `bearer` | `token`（敏感） | 常数时间验证静态 `Authorization: Bearer` 令牌 |
| `header` | 每个凭证头一个条目（全部敏感） | 请求必须携带每个配置对（AND）；配置为空值的凭证永不通过 |
| `query` | 每个凭证参数一个条目（全部敏感） | 请求必须携带每个配置对（AND） |
| `signature` | `secret`（hex，敏感） | 验证框架 HMAC 签名约定：对系统编码、方法、路径签名的 `x-timestamp` / `x-nonce` / `x-signature` —— 经共享 nonce 存储防重放 |
| `script` | 自由参数（全部敏感）+ 验证脚本体 | 在零 IO 运行时中执行自定义验证体；可见 `request` 与解密后的 `params`，返回真值即放行 |

应用可通过 `vef.ProvideIntegrationInboundAuthScheme` 注册自定义 scheme：

```go
type InboundAuthScheme interface {
    Name() string
    Verify(ctx context.Context, req *InboundRequest, auth *InboundAuthConfig) error
    SensitiveParams() []string
}
```

## 业务处理器

入站契约的业务侧是 `integration.InboundHandler` —— 每个契约编码一个处理器，
启动时注册：

```go
vef.ProvideIntegrationInboundHandler(func(db orm.DB) integration.InboundHandler {
    return integration.NewInboundHandler("patient.sync",
        func(ctx context.Context, input PatientSyncInput) (PatientSyncOutput, error) {
            // input 已通过契约输入 Schema 校验
            return doSync(ctx, db, input)
        })
})
```

- `integration.NewInboundHandler[I, O](contract, handle)` 适配类型化函数；
  经 Schema 校验的输入通过 JSON 往返解码为 `I`。
- 处理器必须幂等：外部系统按 at-least-once 投递。
- 契约编码为空或重复注册在启动时报错。
- 契约没有注册处理器时投递失败并返回 `ErrInboundHandlerMissing`
  （HTTP 501）——这是部署故障，不是调用方错误。

## 适配器脚本环境

入站适配器脚本（方向 `inbound`）可见引擎基线，外加：

| 绑定 | 内容 |
| --- | --- |
| `request` | 只读线上请求：`{ protocol, method, path, headers, query, body, clientAddr }` —— 头名小写化，`body` 为字符串 |
| `system` | 系统只读视图：`{ code, name, params }` |
| `dispatch(input)` | 按契约输入 Schema 校验 `input`、运行业务处理器、按输出 Schema 校验其输出并返回；失败以可捕获异常上浮 |
| `codes` | 码值映射翻译（翻译外部系统码值正是入站脚本的核心工作）—— 见[码值映射](./code-maps) |

典型脚本：

```js
// 1. 把厂商线上请求翻译为标准输入。
const body = JSON.parse(request.body);
const input = {
    patientId: body.patient_no,
    gender: codes.toCanonical('gender', body.sex),
};

// 2. 分发给业务处理器（两侧都有 Schema 强制）。
const output = dispatch(input);

// 3. 组装厂商期望的应答。
({ code: '0', msg: 'ok', data: { syncId: output.syncId } })
```

批量负载可以多次调用 `dispatch`；每次分发都被记录，最后一次分发失败即使被
脚本捕获（用于组装部分成功应答）也会保留其失败分类。

## 组装应答

脚本的最终值成为 HTTP 响应：

- `null` / `undefined` → `200`，空响应体。
- 普通值 → `200`，值作为 JSON。
- 形如 `{ $response: { status, headers, body } }` 的对象获得原始控制：
  `status` 默认 200（越界值回退），`headers` 原样设置，字符串 `body` 以
  `text/plain` 发送（除非脚本自行设置了内容类型），其他 body 为 JSON。
  `$` 前缀标记不可能与真实业务负载冲突。

```js
({ $response: {
    status: 200,
    headers: { 'content-type': 'application/xml' },
    body: '<Response><Code>0</Code></Response>',
}})
```

未捕获的管线错误会重映射为外部调用方重试逻辑可以理解的状态码（`/api` 面上
业务错误以 HTTP 200 承载，但外部调用方不读信封）：

| 条件 | 状态码 |
| --- | --- |
| 验证失败、未知/禁用的系统 | `401` |
| 超出限流 | `429` |
| 契约/适配器缺失或禁用（已通过验证的调用方） | `404`（统一，调用方无法探测缺的是哪一块） |
| 输入被契约 Schema 拒绝 | `400` |
| 没有注册业务处理器 | `501` |
| 脚本故障、渲染故障及其他 | `500` |

## 可观测性

- 通过验证的投递按 `vef.integration.log` 记入调用日志（方向 `inbound`）；
  被拒绝的投递刻意不入日志（未认证流量不得增长持久证据链），只以空契约聚合
  进统计。
- `dry_run_inbound` 是入站测试台：对合成的外部请求执行入站脚本，业务处理器
  被桩替换为返回给定样例输出。验证被绕过——测试台测的是翻译，不是凭证；
  不触碰业务代码、不做任何记录。见 [RPC 资源](./resources#integrationops)。

## 下一步

[码值映射](./code-maps)介绍两个流向共用的按系统值翻译。
