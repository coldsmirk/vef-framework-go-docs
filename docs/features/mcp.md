---
sidebar_position: 8
---

# MCP

VEF has first-class support for MCP (Model Context Protocol) server integration.

## Runtime Behavior

The MCP module is always part of the boot sequence, but the MCP server only activates when:

| Config | Meaning |
| --- | --- |
| `vef.mcp.enabled = true` | MCP server and endpoint become active |

When enabled, the HTTP endpoint is mounted at:

```text
/mcp
```

## What The Module Provides

When enabled, the MCP module provides:

| Runtime piece | Meaning |
| --- | --- |
| MCP server construction | builds the MCP server instance |
| HTTP handler | adapts the MCP server to HTTP |
| app middleware | mounts `/mcp` into the Fiber app |
| built-in tool | `database_query` |
| built-in prompts | `data-dict-assistant`, `naming-master` |

The module does not currently ship built-in static MCP resources or built-in resource templates.

## Public Provider Interfaces

The public `mcp` package exposes these provider interfaces:

| Interface | Purpose |
| --- | --- |
| `mcp.ToolProvider` | register MCP tools |
| `mcp.ResourceProvider` | register static MCP resources |
| `mcp.ResourceTemplateProvider` | register MCP resource templates |
| `mcp.PromptProvider` | register MCP prompts |

Supporting definition types:

| Type | Purpose |
| --- | --- |
| `mcp.ToolDefinition` | tool + handler |
| `mcp.ResourceDefinition` | static resource + handler |
| `mcp.ResourceTemplateDefinition` | resource template + handler |
| `mcp.PromptDefinition` | prompt + handler |
| `mcp.ServerInfo` | server name, version, and instructions |

## Dependency-Injection Extension Points

Applications can register custom MCP capabilities through:

| Helper | Registers into |
| --- | --- |
| `vef.ProvideMCPTools(...)` | `vef:mcp:tools` |
| `vef.ProvideMCPResources(...)` | `vef:mcp:resources` |
| `vef.ProvideMCPResourceTemplates(...)` | `vef:mcp:templates` |
| `vef.ProvideMCPPrompts(...)` | `vef:mcp:prompts` |
| `vef.SupplyMCPServerInfo(...)` | optional server info override |

## Built-In Tool: `database_query`

The built-in MCP tool currently exposed by the framework is:

| Tool name | Purpose |
| --- | --- |
| `database_query` | execute a parameterized SQL query and return JSON rows |

Input arguments:

| Argument | Meaning |
| --- | --- |
| `sql` | SQL query string using `?` placeholders |
| `params` | optional positional parameters |

Behavior notes:

- results are returned as JSON text content
- UTF-8 byte slices are converted to strings before JSON encoding
- the query runs through `mcp.DBWithOperator(...)`, so the current MCP principal is bound as the operator when possible

## Built-In Prompts

The framework currently registers these built-in prompts:

| Prompt name | Purpose |
| --- | --- |
| `data-dict-assistant` | data dictionary management assistant |
| `naming-master` | naming assistant for code and database naming |

## Schema Helpers

The public `mcp` package also provides JSON Schema helpers for MCP tool and prompt contracts:

| Helper | Purpose |
| --- | --- |
| `mcp.SchemaFor[T]()` | generate schema from a generic type |
| `mcp.SchemaOf(v)` | generate schema from a runtime value |
| `mcp.MustSchemaFor[T]()` | panic-on-error schema generation |
| `mcp.MustSchemaOf(v)` | panic-on-error schema generation |

These helpers are suitable for tool input schemas.

## Authentication

When `vef.mcp.require_auth = true`, the HTTP handler applies Bearer-token verification through the application auth manager.

That matters because the built-in `database_query` tool can already inspect live database state.

## Minimal Example

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

In many apps, the first step is even smaller:

```go
var Module = vef.Module(
  "app:mcp",
  vef.ProvideMCPTools(NewToolProvider),
)
```

## What To Document For Your Own App

If you expose custom MCP features, document:

- which tools are registered
- whether auth is required
- which domain data your resources or prompts expose
- whether any prompt or tool performs side effects

## Next Step

Read [Extension Points](../reference/extension-points) for the full list of MCP-related DI groups and helpers.
