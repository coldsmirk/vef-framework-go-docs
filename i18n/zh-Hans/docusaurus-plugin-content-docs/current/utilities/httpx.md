---
sidebar_position: 9
---

# httpx（出站 HTTP 客户端）

`httpx`（v0.39）是框架的出站 HTTP 客户端，提供流式请求 API——为每个第三方
系统构建一个 `Client`，再由它派生按调用的 `Request`。集成引擎的作用域
`http` 脚本库即构建在它之上；应用 Go 代码也可以直接使用。

> 不要与旧的 `httpx` Fiber 辅助函数混淆——它们在 v0.39 更名为
> [`fiberx`](./small-helpers#fiberx)，把这个名字让给了出站客户端。

## 快速开始

```go
import "github.com/coldsmirk/vef-framework-go/httpx"

client, err := httpx.New(
    httpx.WithBaseURL("https://api.example.com"),
    httpx.WithTimeout(10*time.Second),
    httpx.WithBearerToken(token),
    httpx.WithRetry(httpx.RetryConfig{}), // 默认：3 次尝试，100ms→2s 退避
)

var out struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

resp, err := client.NewRequest().
    SetPathParam("id", "42").
    SetQuery("expand", "profile").
    Get(ctx, "/users/:id")
if err != nil {
    return err
}
if !resp.IsSuccess() {
    return fmt.Errorf("upstream returned %s", resp.Status())
}
if err := resp.JSON(&out); err != nil {
    return err
}
```

## Client

`httpx.New(opts ...Option)` 急切校验选项：畸形的 base/proxy URL 与互相
冲突的传输层选项都会使构建失败。零选项客户端开箱即用——无 base URL、整个
调用（含重试）30s 超时、无重试。`Client` 在 `New` 后不可变、可并发使用；
按调用状态位于 `Request`。

| 选项 | 行为 |
| --- | --- |
| `WithBaseURL(url)` | 相对请求 URL 拼接的绝对基址；绝对请求 URL 绕过它 |
| `WithTimeout(d)` | 约束整个调用（含重试），默认 30s |
| `WithHeader(k, v)` / `WithQuery(k, v)` | 每个请求的默认头 / 查询对 |
| `WithBasicAuth(user, pass)` / `WithBearerToken(token)` | 默认 `Authorization` 头 |
| `WithRetry(cfg)` | 启用自动重试（见下） |
| `WithProxy(url)` | 出站代理 |
| `WithTLSConfig(cfg)` | 自定义 TLS 配置 |
| `WithCookieJar(jar)` | 跨调用 Cookie 持久化 |
| `WithMaxRedirects(n)` | 重定向上限（默认 10；超出返回 `ErrTooManyRedirects`） |
| `WithMaxResponseBody(n)` | 响应体字节上限（`ErrResponseTooLarge`） |
| `WithRequestHook(hooks...)` | 在请求完全构建后、发送前执行——签名、审计、日志的挂点；返回错误则中止调用 |
| `WithResponseHook(hooks...)` | 在响应到达且响应体缓冲后执行 |
| `WithTransport(rt)` / `WithHTTPClient(hc)` | 自定义传输 / 完全自定义 `http.Client`（与传输层选项互斥——`ErrConflictingOptions`） |

除非应用自行设置，客户端发送 `User-Agent: vef/<version>`。

## Request

`client.NewRequest()` 开始一个流式、一次性的请求构建器（重复执行返回
`ErrRequestReused`）：

| 分组 | 方法 |
| --- | --- |
| 请求头 | `SetHeader`、`AddHeader`、`SetHeaders` |
| 查询 | `SetQuery`、`AddQuery`、`SetQueries` |
| 路径参数 | `SetPathParam`、`SetPathParams` —— 替换 `:name` 片段；未解析的片段返回 `ErrMissingPathParam` |
| Cookie / 认证 | `SetCookie`、`SetBasicAuth`、`SetBearerToken` |
| 请求体 | `SetJSON(v)`、`SetXML(v)`、`SetBody(bytes, contentType)`、`SetBodyReader(r, contentType)`、`SetForm(map)`、`AddFormField(k, v)`、`AddFile(field, path)`、`AddFileReader(field, filename, r)` |
| 超时 | `SetTimeout(d)` —— 按请求覆盖客户端超时 |
| 执行 | `Get`、`Post`、`Put`、`Patch`、`Delete`、`Head`、`Options` 或 `Do(ctx, method, url)` |
| 自省 | `Method()`、`URL()`、`Header(k)`、`Headers()`、`Body()`、`Context()` —— 请求钩子使用的读取面 |

`SetForm`/`AddFormField` 生成 URL 编码表单；添加文件自动升级为 multipart。

## Response

| 方法 | 契约 |
| --- | --- |
| `StatusCode()` / `Status()` / `IsSuccess()` | 状态自省；`IsSuccess` 为 2xx |
| `Header(k)` / `Headers()` / `Cookies()` | 响应元数据 |
| `Body()` / `String()` | 已缓冲的响应体（始终完整读取并缓冲） |
| `JSON(v)` / `XML(v)` | 解码响应体 |
| `Duration()` | 调用耗时 |
| `Attempts()` | 尝试次数，含首个调用 |
| `Request()` | 来源请求 |

非 2xx 响应**不是**错误：传输层调用已成功，各状态码的含义由应用决定。
错误只保留给传输失败、超时与策略违规。

## 重试

`WithRetry(httpx.RetryConfig{...})` 启用自动重试。零值字段解析为默认值：

| 字段 | 默认 | 含义 |
| --- | --- | --- |
| `MaxAttempts` | `3` | 总尝试次数，含首个调用 |
| `InitialBackoff` | `100ms` | 首次重试前的基础延迟；每次重试翻倍并施加全抖动 |
| `MaxBackoff` | `2s` | 尝试间延迟上限，含服务端 `Retry-After` |
| `RetryIf` | 见下 | 自定义谓词，整体替换默认策略 |

默认策略在传输错误或 `429`/`502`/`503`/`504` 响应时重试，且**仅限幂等
方法**（GET、HEAD、PUT、DELETE、OPTIONS、TRACE）——除非 `RetryIf` 放行，
POST 永不重试。

## 错误哨兵

| 错误 | 触发 |
| --- | --- |
| `ErrInvalidOption` | 畸形的 base/proxy URL 或其他非法选项值 |
| `ErrConflictingOptions` | `WithHTTPClient` 与传输层选项同时使用 |
| `ErrInvalidRequestURL` | 无法解析的请求 URL，或没有 base URL 时使用相对 URL |
| `ErrMissingPathParam` | `:name` 片段未被解析 |
| `ErrRequestReused` | 一次性请求被二次执行 |
| `ErrTooManyRedirects` | 超出重定向上限 |
| `ErrResponseTooLarge` | 响应体超出配置上限 |

## 另请参阅

- [集成引擎](../integration/overview) —— 系统以声明方式配置 `httpx` 客户端（认证 scheme、重试策略、超时）
- [小工具集](./small-helpers) —— `fiberx`，即原名 `httpx` 的入站 Fiber 请求辅助
