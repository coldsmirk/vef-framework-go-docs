---
sidebar_position: 1
slug: /intro
---

# Introduction

VEF Framework is a resource-driven Go web framework for building internal platforms, admin systems, and service APIs on top of Uber FX, Fiber, and Bun.

The easiest way to understand VEF is to start from its runtime shape:

1. `vef.Run(...)` boots a fixed module pipeline.
2. You register your own resources, middleware, and behaviors through FX groups.
3. The API engine collects operations from those resources and mounts them as RPC or REST endpoints.
4. Handler parameters are injected automatically from the request context, decoded params, metadata, and container-managed services.

That means VEF is not primarily a “router helper” or a “CRUD library”. It is a framework for composing a Go application around explicit resources and predictable defaults.

## What VEF gives you by default

`vef.Run(...)` boots a fixed, ordered module pipeline; see [Application Lifecycle](./core-concepts/lifecycle) for the full boot sequence and what each stage guarantees.

Once the app starts, VEF already has opinions about:

- API versioning: default version is `v1`
- Authentication: default API auth strategy is Bearer token
- Request timeout: default is `30s`
- Rate limiting: default is `100` requests per `5m`
- Response envelope: success and error responses use `result.Result`
- Storage: memory storage is used when no storage provider is configured

These defaults are runtime behavior, not optional conventions.

## RPC and REST are both first-class

VEF supports two API styles side by side:

- RPC resources, mounted behind `POST /api`
- REST resources, mounted under `/api/<resource>`

They are declared explicitly:

```go
api.NewRPCResource("sys/user", ...)
api.NewRESTResource("users", ...)
```

VEF does **not** generate REST routes automatically from an RPC resource. If you want both styles, define both resources intentionally.

## What you write as an application developer

Most applications only touch a small set of public APIs:

- `vef.Run(...)`
- `vef.Module(...)`
- `vef.ProvideAPIResource(...)`
- `api.NewRPCResource(...)`
- `api.NewRESTResource(...)`
- `api.OperationSpec`
- `crud.NewCreate(...)`, `crud.NewFindPage(...)`, and other CRUD builders
- `orm.DB`
- `result.Ok(...)` and `result.Err(...)`
- `security` extension interfaces such as `UserLoader`, `PermissionChecker`, and `RolePermissionsLoader`

The rest of the framework exists mostly to support those user-facing entry points.

## Built-in resources you can use immediately

The framework also ships with several built-in resources and modules:

- `security/auth` for login, refresh, logout, challenge resolution, and optional user info loading
- `sys/storage` for multipart upload (init/part/list/complete/abort) plus a `/storage/files/<key>` download proxy
- `sys/schema` for schema inspection
- `sys/monitor` for runtime and host monitoring
- `sys/cron/*` for durable schedule management when the cron store is enabled
- `integration/*` for the integration engine when `vef.IntegrationModule` is enabled
- MCP middleware and server integration when enabled

You do not need to implement these from scratch unless your application requirements differ.

## Where to start

Most applications touch the framework in this order:

1. [Installation](./getting-started/installation) — environment and package setup
2. [Quick Start](./getting-started/quick-start) — a minimal app that actually boots and serves an endpoint
3. [Your First CRUD API](./getting-started/first-crud-api) — model, table, resource, and curl-verified CRUD endpoints end to end
4. [Core Concepts](./core-concepts/overview) — how modules, dependency injection, and the application lifecycle fit together
5. [Building APIs](./building-apis/api) — resources, operations, routing, and parameter binding
6. [Data Access](./data-access/models) — models, search filters, CRUD, the SQL builder, and transactions
7. [Security](./security/authentication) — authentication, authorization, and login hardening

From there, branch out by task:

- [Data Tools](./data-tools/expression) — expression engine, mold data cleansing, i18n, tabular import/export
- [Infrastructure](./infrastructure/cache) — cache, cron and durable schedules, sequence, event bus, server push, storage, schema, monitor
- [AI Integration](./ai-integration/ai) — AI helpers and MCP
- [Approval](./approval) — the workflow/approval engine
- [Integration](./integration/overview) — config- and script-driven integration with external systems
- [Advanced](./advanced/cqrs) — CQRS, custom parameter resolvers, CLI tooling
- [Utilities](./utilities/small-helpers) — small, focused helper packages
- [Conventions](./conventions/application-project-conventions) — project layout and database conventions
- [Reference](./reference/configuration-reference) — configuration keys, built-in resources, and API indexes

If you are new to the framework, go to [Installation](./getting-started/installation) next.
