---
sidebar_position: 8
---

# MCP

VEF 对 MCP（Model Context Protocol）提供了一等支持。

## 运行时行为

MCP 模块始终在启动链中，但 MCP server 只有在满足下面条件时才真正启用：

| 配置 | 含义 |
| --- | --- |
| `vef.mcp.enabled = true` | MCP server 和对应 HTTP 端点才会激活 |
| `vef.mcp.require_auth = false` | 显式允许匿名访问 MCP 端点 |

启用后，HTTP 端点挂载在：

```text
/mcp
```

这个端点是 Streamable HTTP MCP endpoint。应用 middleware 会以 order `500`
把它注册到所有 HTTP method。

Server identity 优先来自显式提供的 `mcp.ServerInfo`。如果没有提供，server
name 会依次回退到 `vef.app.name` 和 `vef-mcp-server`；默认 version 是
`v1.0.0`，默认 instructions 为空。

## 模块提供了什么

启用后，MCP 模块会提供：

| 运行时组件 | 含义 |
| --- | --- |
| MCP server 构造 | 创建 MCP server 实例 |
| HTTP handler | 把 MCP server 适配成 HTTP |
| app middleware | 以 order `500` 和所有 HTTP method 把 `/mcp` 注册到 Fiber 应用 |
| 内置 tool | `database_query` |
| 内置 prompt | `naming-master` |

当前模块没有内置的静态 MCP resource，也没有内置 resource template。

## 公共 Provider 接口

公开的 `mcp` 包暴露了这些 provider 接口：

| 接口 | 作用 |
| --- | --- |
| `mcp.ToolProvider` | 注册 MCP 工具 |
| `mcp.ResourceProvider` | 注册静态 MCP 资源 |
| `mcp.ResourceTemplateProvider` | 注册 MCP resource template |
| `mcp.PromptProvider` | 注册 MCP prompt |

配套的 definition 类型：

| 类型 | 作用 |
| --- | --- |
| `mcp.ToolDefinition` | tool + handler |
| `mcp.ResourceDefinition` | 静态 resource + handler |
| `mcp.ResourceTemplateDefinition` | resource template + handler |
| `mcp.PromptDefinition` | prompt + handler |
| `mcp.ServerInfo` | server 名称、版本、说明 |

VEF 自有 definition fields 是显式公开 surface：`ToolDefinition.Tool`、
`ToolDefinition.Handler`、`ResourceDefinition.Resource`、
`ResourceDefinition.Handler`、`ResourceTemplateDefinition.Template`、
`ResourceTemplateDefinition.Handler`、`PromptDefinition.Prompt`、
`PromptDefinition.Handler`、`ServerInfo.Name`、`ServerInfo.Version` 和
`ServerInfo.Instructions`。

这个包也重新导出了应用在自定义 handler 中会用到的 MCP SDK protocol 类型：

| 类型组 | 别名 |
| --- | --- |
| server/session | `Server`, `ServerOptions`, `ServerSession`, `Implementation` |
| content | `Content`, `TextContent`, `ImageContent`, `AudioContent` |
| tools | `Tool`, `CallToolRequest`, `CallToolResult`, `ToolHandler`, `Annotations` |
| resources | `Resource`, `ReadResourceRequest`, `ReadResourceResult`, `ResourceHandler`, `ResourceTemplate` |
| prompts | `Prompt`, `PromptArgument`, `PromptMessage`, `Role`, `GetPromptParams`, `GetPromptRequest`, `GetPromptResult`, `PromptHandler` |

MCP SDK pass-through surface 是有意暴露的便利层。这些 aliases 来自
`github.com/modelcontextprotocol/go-sdk/mcp`，其 promoted SDK methods 保持
上游签名和行为。public API index 列出精确 method signatures；VEF-specific
行为仅限下面说明的 provider interfaces、schema helpers、principal/database
helpers 和 tool-result helpers。

当前 MCP 包审计在生成的 API ledger 中锁定 **119 public MCP entries**。分组
member surface 覆盖 **75 grouped MCP field/method entries**，分布在 **27
MCP receiver/type families** 中：其中包含 **11 exported MCP field entries**
和 **64 exported MCP method entries**。生成的公开 API 索引仍是完整签名清单。

辅助 helper：

| Helper | 作用 |
| --- | --- |
| `mcp.NewToolResultText(text)` | 返回文本 tool result |
| `mcp.NewToolResultError(message)` | 返回错误 tool result |
| `mcp.GetPrincipalFromContext(ctx)` | 从 MCP request context 读取认证后的 VEF principal；不存在时返回 `security.PrincipalAnonymous` |
| `mcp.DBWithOperator(ctx, db)` | 把 MCP principal ID 绑定为 `orm.PlaceholderKeyOperator` |
| `mcp.ResourceNotFoundError` | SDK resource-not-found error 别名 |

## 依赖注入扩展点

应用可以通过这些 helper 注册自定义 MCP 能力：

| Helper | 注册到 |
| --- | --- |
| `vef.ProvideMCPTools(...)` | `vef:mcp:tools` |
| `vef.ProvideMCPResources(...)` | `vef:mcp:resources` |
| `vef.ProvideMCPResourceTemplates(...)` | `vef:mcp:templates` |
| `vef.ProvideMCPPrompts(...)` | `vef:mcp:prompts` |
| `vef.SupplyMCPServerInfo(...)` | 可选 server info 覆盖 |

## 内置 Tool：`database_query`

框架当前内置暴露的 MCP tool 是：

| Tool 名称 | 作用 |
| --- | --- |
| `database_query` | 执行参数化 SQL 查询并返回 JSON 行数据 |

输入参数：

| 参数 | 含义 |
| --- | --- |
| `sql` | 使用 `?` 占位符的 SQL 语句 |
| `params` | 可选的位置参数列表 |

行为说明：

- 只允许 read-only `SELECT` 语句；SQL 会在执行前检查，带数据修改语义的 CTE 或有副作用的语句会被拒绝
- `sql` 必填，且不能为空
- 查询结果会以 JSON 文本内容返回
- tool 失败会作为 MCP tool error result 返回，而不是作为 Go handler error 返回
- UTF-8 的 `[]byte` 会在 JSON 编码前自动转成字符串，包括嵌套在 map 或 slice
  里的值
- 非 UTF-8 的 `[]byte` 会保留为 binary value，因此在 JSON output 中会变成
  Base64 string
- 查询通过 `mcp.DBWithOperator(...)` 执行，因此当前 MCP principal 会在可用时绑定为 operator

read-only guard 是 fail-closed：parse errors、空 statement、非 read statement、
data-modifying CTE、多语句写入尝试，以及 `pg_read_file`、`pg_sleep`、
`nextval`、`setval` 这类有副作用的函数都会在执行前被拒绝。

## 内置 Prompt

框架当前注册了这些内置 prompt：

| Prompt 名称 | 作用 |
| --- | --- |
| `naming-master` | 面向代码 identifier、数据库对象、审计字段、索引、约束和外键策略的命名助手 |

## Schema Helper

公共 `mcp` 包还提供 JSON Schema helper，用于构造 MCP tool / prompt 契约：

| Helper | 作用 |
| --- | --- |
| `mcp.SchemaFor[T]()` | 从泛型类型生成 schema |
| `mcp.SchemaOf(v)` | 从运行时值生成 schema |
| `mcp.MustSchemaFor[T]()` | schema 生成失败时 panic |
| `mcp.MustSchemaOf(v)` | schema 生成失败时 panic |

这些 helper 很适合用于 tool 输入 schema。

Schema 生成会使用框架针对 MCP 调整过的 reflector 设置：

- `jsonschema:"required"` 会把字段标记为 required
- 嵌套 Go 类型会内联，而不是生成 `$ref`
- 不会生成来自包路径的 `$id`
- 生成结果中的 `$schema` 字段会被移除，因为 MCP input schema 不使用它
- `mcp.SchemaOf(nil)` 返回 `nil`
- `mcp.MustSchemaOf(nil)` 会 panic，消息为 `mcp: failed to generate schema`
- `mcp.MustSchemaFor[T]()` 和 `mcp.MustSchemaOf(v)` 适合用于希望 schema 生成失败直接变成启动期错误的场景

支持的 `jsonschema` tag 关键字包括：

| 范围 | Tags |
| --- | --- |
| 通用 | `required`, `nullable`, `title=...`, `description=...`, `type=...`, `anchor=...`, `default=...`, `example=...`, `enum=...` |
| union helper | `oneof_required=...`, `anyof_required=...`, `oneof_ref=...`, `oneof_type=...`, `anyof_ref=...`, `anyof_type=...` |
| string | `minLength=...`, `maxLength=...`, `pattern=...`, `format=...`, `readOnly=true`, `writeOnly=true` |
| number/integer | `minimum=...`, `maximum=...`, `exclusiveMinimum=...`, `exclusiveMaximum=...`, `multipleOf=...` |
| array | `minItems=...`, `maxItems=...`, `uniqueItems=true` |
| 独立 tag | `jsonschema_description:"..."`, `jsonschema_extras:"a=b,c=d"` |

## 认证

MCP 端点默认是安全的。如果 `vef.mcp.require_auth` 未配置或设置为 `true`，HTTP handler 会通过应用 auth manager 加 Bearer token 校验。
这里使用框架的 token authentication path。
SDK 接受 `Bearer <token>` 和 `bearer <token>` 两种 header prefix；没有 Bearer
prefix 的裸 token 会被拒绝。

只有在你明确需要匿名 MCP surface 时，才设置 `vef.mcp.require_auth = false`。

这点很重要，因为框架自带的 `database_query` 已经可以读取实时数据库状态。

## 最小示例

```go
package appmcp

import (
  "github.com/coldsmirk/vef-framework-go"
  "github.com/coldsmirk/vef-framework-go/mcp"
)

var Module = vef.Module(
  "app:mcp",
  vef.ProvideMCPTools(NewToolProvider),
  vef.SupplyMCPServerInfo(&mcp.ServerInfo{
    Name:         "my-app",
    Version:      "v1.0.0",
    Instructions: "Internal assistant surface for My App",
  }),
)
```

很多应用的第一步甚至更小：

```go
var Module = vef.Module(
  "app:mcp",
  vef.ProvideMCPTools(NewToolProvider),
)
```

## 你自己的应用应该补充说明什么

如果你暴露了自定义 MCP 能力，文档里至少应写清楚：

- 注册了哪些 tool
- 是否要求认证
- resource / prompt 会暴露哪些领域数据
- 哪些 prompt 或 tool 会产生副作用

## 下一步

继续阅读 [扩展点](../reference/extension-points)，看 MCP 相关 DI group 和 helper 的完整列表。
