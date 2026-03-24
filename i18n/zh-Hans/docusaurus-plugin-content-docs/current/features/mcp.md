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

启用后，HTTP 端点挂载在：

```text
/mcp
```

## 模块提供了什么

启用后，MCP 模块会提供：

| 运行时组件 | 含义 |
| --- | --- |
| MCP server 构造 | 创建 MCP server 实例 |
| HTTP handler | 把 MCP server 适配成 HTTP |
| app middleware | 把 `/mcp` 注册到 Fiber 应用 |
| 内置 tool | `database_query` |
| 内置 prompt | `data-dict-assistant`、`naming-master` |

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

- 查询结果会以 JSON 文本内容返回
- UTF-8 的 `[]byte` 会在 JSON 编码前自动转成字符串
- 查询通过 `mcp.DBWithOperator(...)` 执行，因此当前 MCP principal 会在可用时绑定为 operator

## 内置 Prompt

框架当前注册了这些内置 prompt：

| Prompt 名称 | 作用 |
| --- | --- |
| `data-dict-assistant` | 数据字典管理助手 |
| `naming-master` | 代码与数据库命名助手 |

## Schema Helper

公共 `mcp` 包还提供 JSON Schema helper，用于构造 MCP tool / prompt 契约：

| Helper | 作用 |
| --- | --- |
| `mcp.SchemaFor[T]()` | 从泛型类型生成 schema |
| `mcp.SchemaOf(v)` | 从运行时值生成 schema |
| `mcp.MustSchemaFor[T]()` | schema 生成失败时 panic |
| `mcp.MustSchemaOf(v)` | schema 生成失败时 panic |

这些 helper 很适合用于 tool 输入 schema。

## 认证

当 `vef.mcp.require_auth = true` 时，HTTP handler 会通过应用 auth manager 加 Bearer token 校验。

这点很重要，因为框架自带的 `database_query` 已经可以读取实时数据库状态。

## 最小示例

```go
package appmcp

import (
  vef "github.com/coldsmirk/vef-framework-go"
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
