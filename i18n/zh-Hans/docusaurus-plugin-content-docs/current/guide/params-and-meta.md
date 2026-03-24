---
sidebar_position: 5
---

# 参数与元信息（Meta）

VEF 把请求输入拆成两个部分：

- `params`：业务输入
- `meta`：请求级控制信息

这种分层在 RPC 请求中是显式存在的，在 REST 请求中也会被框架内部保留下来。

## 请求模型总览

| 区段 | 作用 | 常见内容 |
| --- | --- | --- |
| `params` | 业务载荷 | 搜索字段、写入参数、上传文件、命令输入 |
| `meta` | 请求控制信息 | 分页、排序、导出格式、选项列映射 |

## 支持的 typed 目标

框架支持以下几类请求解码目标：

| 目标类型 | 解码来源 | 是否自动验证 | 常见用途 |
| --- | --- | --- | --- |
| 嵌入 `api.P` 的 typed struct | `params` | 是 | 业务 params |
| 嵌入 `api.M` 的 typed struct | `meta` | 是 | typed meta |
| `page.Pageable` | `meta` | 是 | 分页 |
| `api.Params` | `params` | 否 | 原始动态 payload |
| `api.Meta` | `meta` | 否 | 原始动态 meta |

## `api.P` 标记 Params 结构体

把 `api.P` 嵌入到应该从 `Request.Params` 解码的结构体里：

```go
type CreateUserParams struct {
	api.P

	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
```

当 handler 参数是 `CreateUserParams` 或 `*CreateUserParams` 时，框架会：

1. 解码 `params`
2. 验证结构体
3. 注入 typed value

## `api.M` 标记 Meta 结构体

把 `api.M` 嵌入到应该从 `Request.Meta` 解码的结构体里：

```go
type PageMeta struct {
	api.M
	page.Pageable
}
```

typed 请求控制信息就是这样注入的。

## 内置 Meta Helper

框架对以下 meta 相关 helper 类型有内置支持：

| 类型 | 含义 | 说明 |
| --- | --- | --- |
| `page.Pageable` | 页码与页大小 | 会被直接识别为内置 meta 类型 |
| `crud.Sortable` | 排序规则 | 通常通过嵌入 typed `api.M` 结构体使用 |

这里有一个容易忽略的区别：

- `page.Pageable` 是框架内置识别的 meta 目标类型
- `crud.Sortable` 并不是单独的内置 meta 类型，但它嵌入到 typed `api.M` 结构体里时会自然工作

## 原始访问

如果你不想使用 typed 解码，handler 可以直接接收：

| 类型 | 含义 |
| --- | --- |
| `api.Params` | 原始 params map |
| `api.Meta` | 原始 meta map |

原始访问适合动态代理类、半结构化或请求契约不稳定的场景。稳定的业务 API 仍应优先使用 typed 结构体。

## RPC 解码规则

对 RPC 请求来说，解码规则依赖请求的内容类型：

| RPC 请求类型 | `params` 来源 | `meta` 来源 | 说明 |
| --- | --- | --- | --- |
| JSON body | JSON 里的 `params` 对象 | JSON 里的 `meta` 对象 | 标准 RPC 形态 |
| form 请求 | form 字段 `params`，再按 JSON 字符串解析 | form 字段 `meta`，再按 JSON 字符串解析 | 适合表单风格客户端 |
| multipart form | form 字段 `params`，按 JSON 字符串解析，并把上传文件并入 params | form 字段 `meta`，按 JSON 字符串解析 | 文件字段会被塞进 params |

## REST 解码规则

对 REST 请求来说：

| 输入来源 | 最终落点 | 说明 |
| --- | --- | --- |
| query string | `params` | 读操作过滤条件或普通请求字段 |
| `POST` / `PUT` / `PATCH` 的 JSON body | `params` | 写入 payload |
| `POST` / `PUT` / `PATCH` 的 multipart 字段 | `params` | 包括上传文件 |
| `X-Meta-*` headers | `meta` | 请求级控制参数 |

这意味着分页和排序并不会自动从 query string 塞进内置 meta helper 里。如果 handler 期望的是 `page.Pageable` 这类 meta 目标，调用方应该通过 `X-Meta-*` headers 或显式 typed meta 契约来提供。

## Multipart 文件支持

multipart 上传可以填充这些 params 形态：

| 形态 | 说明 |
| --- | --- |
| `*multipart.FileHeader` | 标准单文件上传字段 |
| `api.Params` 里的原始文件条目 | 适合代理类或动态 handler |

内置的存储和导入接口就是通过这套机制接收上传文件的。

## 验证行为

typed params 和 typed meta 在解码后会自动验证。

| 目标类型 | 是否验证 |
| --- | --- |
| typed `api.P` struct | 是 |
| typed `api.M` struct | 是 |
| `page.Pageable` | 是 |
| `api.Params` | 否 |
| `api.Meta` | 否 |

验证发生在解码完成后，通过 `validator.Validate(...)` 执行。如果校验失败，框架会返回带翻译字段消息的 bad-request 风格结果。

## 常见模式

### 标准搜索请求

```go
type UserSearch struct {
	api.P
	Keyword string `json:"keyword" search:"contains,column=username|email"`
}

type UserMeta struct {
	api.M
	page.Pageable
	crud.Sortable
}
```

### 动态代理风格请求

```go
func (*ProxyResource) Forward(params api.Params, meta api.Meta) error {
	// handle raw data
	return nil
}
```

## 实践建议

- 业务字段放到 `params`
- 分页、排序、导出模式等请求控制信息放到 `meta`
- 稳定接口优先使用 typed 结构体，不要滥用原始 map
- 显式嵌入 `api.P` 和 `api.M`，让解码意图保持清晰
- 只有在请求契约真的动态时，才使用 `api.Params` / `api.Meta`

## 下一步

继续阅读 [自定义处理器](./custom-handlers)，看这些解码结果是如何注入到 handler 签名里的。
