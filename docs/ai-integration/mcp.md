---
sidebar_position: 2
---

# MCP

VEF has first-class support for MCP (Model Context Protocol) server integration.

## Runtime Behavior

The MCP module is always part of the boot sequence, but the MCP server only activates when:

| Config | Meaning |
| --- | --- |
| `vef.mcp.enabled = true` | MCP server and endpoint become active |
| `vef.mcp.require_auth = false` | explicitly allow anonymous MCP access |

When enabled, the HTTP endpoint is mounted at:

```text
/mcp
```

The endpoint is a Streamable HTTP MCP endpoint. The app middleware registers it
for all HTTP methods at middleware order `500`.

Server identity comes from `mcp.ServerInfo` when supplied. Otherwise the server
name falls back to `vef.app.name`, then to `vef-mcp-server`; the default version
is `version.VEFVersion`, and default instructions are empty.

## What The Module Provides

When enabled, the MCP module provides:

| Runtime piece | Meaning |
| --- | --- |
| MCP server construction | builds the MCP server instance |
| HTTP handler | adapts the MCP server to HTTP |
| app middleware | mounts `/mcp` into the Fiber app with order `500` and all HTTP methods |
| built-in tool | `database_query` |
| built-in prompts | `naming-master` |

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

The VEF-owned definition fields are explicit: `ToolDefinition.Tool`,
`ToolDefinition.Handler`, `ResourceDefinition.Resource`,
`ResourceDefinition.Handler`, `ResourceTemplateDefinition.Template`,
`ResourceTemplateDefinition.Handler`, `PromptDefinition.Prompt`,
`PromptDefinition.Handler`, `ServerInfo.Name`, `ServerInfo.Version`, and
`ServerInfo.Instructions`.

The package also re-exports the MCP SDK protocol types that application code
uses when building custom handlers:

| Type group | Aliases |
| --- | --- |
| server/session | `Server`, `ServerOptions`, `ServerSession`, `Implementation` |
| content | `Content`, `TextContent`, `ImageContent`, `AudioContent` |
| tools | `Tool`, `CallToolRequest`, `CallToolResult`, `ToolHandler`, `Annotations` |
| resources | `Resource`, `ReadResourceRequest`, `ReadResourceResult`, `ResourceHandler`, `ResourceTemplate` |
| prompts | `Prompt`, `PromptArgument`, `PromptMessage`, `Role`, `GetPromptParams`, `GetPromptRequest`, `GetPromptResult`, `PromptHandler` |

The MCP SDK pass-through surface is intentional. These aliases come from
`github.com/modelcontextprotocol/go-sdk/mcp`, and their promoted SDK methods keep
the upstream signatures and behavior. The public API index lists the exact
method signatures; VEF-specific behavior is limited to the provider interfaces,
schema helpers, principal/database helpers, and tool-result helpers documented
below.

Utility helpers:

| Helper | Purpose |
| --- | --- |
| `mcp.NewToolResultText(text)` | return text content from a tool |
| `mcp.NewToolResultError(message)` | return an error tool result |
| `mcp.GetPrincipalFromContext(ctx)` | read the authenticated VEF principal from an MCP request, or return `security.PrincipalAnonymous` |
| `mcp.DBWithOperator(ctx, db)` | bind the MCP principal ID as `orm.PlaceholderKeyOperator` |
| `mcp.ResourceNotFoundError` | SDK resource-not-found error alias |

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

- only read-only `SELECT` statements are permitted; the SQL is checked before
  execution and data-changing CTEs or side-effect statements are rejected
- `sql` is required and must not be empty
- results are returned as JSON text content
- tool failures are returned as MCP tool error results instead of Go handler errors
- UTF-8 byte slices are converted to strings before JSON encoding, including
  values nested inside maps or slices
- non-UTF-8 byte slices are preserved as binary values and therefore become
  Base64 strings in the JSON output
- the query runs through `mcp.DBWithOperator(...)`, so the current MCP principal is bound as the operator when possible

The read-only guard is fail-closed: parse errors, empty statements,
non-read statements, data-modifying CTEs, multi-statement write attempts, and
side-effecting functions such as `pg_read_file`, `pg_sleep`, `nextval`, and
`setval` are rejected before execution.

## Built-In Prompts

The framework currently registers these built-in prompts:

| Prompt name | Purpose |
| --- | --- |
| `naming-master` | naming assistant for code identifiers, database objects, audit fields, indexes, constraints, and foreign key strategy |

## Schema Helpers

The public `mcp` package also provides JSON Schema helpers for MCP tool and prompt contracts:

| Helper | Purpose |
| --- | --- |
| `mcp.SchemaFor[T]()` | generate schema from a generic type |
| `mcp.SchemaOf(v)` | generate schema from a runtime value |
| `mcp.MustSchemaFor[T]()` | panic-on-error schema generation |
| `mcp.MustSchemaOf(v)` | panic-on-error schema generation |

These helpers are suitable for tool input schemas.

Schema generation follows the framework's MCP-specific reflector settings:

- `jsonschema:"required"` marks a field as required
- nested Go types are inlined instead of emitted as `$ref`
- package-derived `$id` values are not generated
- the generated `$schema` property is removed because MCP input schemas do not use it
- `mcp.SchemaOf(nil)` returns `nil`
- `mcp.MustSchemaOf(nil)` panics with `mcp: cannot generate schema for nil value`
- schema-generation failures in `mcp.MustSchemaFor[T]()` and
  `mcp.MustSchemaOf(v)` panic with `mcp: generate schema: ...`, preserving the
  underlying cause
- `mcp.MustSchemaFor[T]()` and `mcp.MustSchemaOf(v)` are convenience helpers for places where schema generation failure should be a boot-time error

Supported `jsonschema` tag keywords are:

| Area | Tags |
| --- | --- |
| generic | `required`, `nullable`, `title=...`, `description=...`, `type=...`, `anchor=...`, `default=...`, `example=...`, `enum=...` |
| union helpers | `oneof_required=...`, `anyof_required=...`, `oneof_ref=...`, `oneof_type=...`, `anyof_ref=...`, `anyof_type=...` |
| string | `minLength=...`, `maxLength=...`, `pattern=...`, `format=...`, `readOnly=true`, `writeOnly=true` |
| number/integer | `minimum=...`, `maximum=...`, `exclusiveMinimum=...`, `exclusiveMaximum=...`, `multipleOf=...` |
| array | `minItems=...`, `maxItems=...`, `uniqueItems=true` |
| standalone tags | `jsonschema_description:"..."`, `jsonschema_extras:"a=b,c=d"` |

## Authentication

The MCP endpoint is secure by default. If `vef.mcp.require_auth` is omitted or
set to `true`, the HTTP handler applies Bearer-token verification through the
application auth manager using the framework token authentication path.
The SDK accepts both `Bearer <token>` and `bearer <token>` header prefixes; a
raw token without a Bearer prefix is rejected.

Set `vef.mcp.require_auth = false` only for deliberately anonymous MCP surfaces.

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
