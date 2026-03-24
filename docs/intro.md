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

The framework boot pipeline is:

`config -> database -> orm -> middleware -> api -> security -> event -> cqrs -> cron -> redis -> mold -> storage -> sequence -> schema -> monitor -> mcp -> app`

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
- `sys/storage` for upload, presigned URLs, temporary file cleanup, object listing, and object metadata
- `sys/schema` for schema inspection
- `sys/monitor` for runtime and host monitoring
- MCP middleware and server integration when enabled

You do not need to implement these from scratch unless your application requirements differ.

## How to read this documentation

The docs are organized around the order in which most users encounter the framework:

- [Installation](./getting-started/installation): environment and package setup
- [Quick Start](./getting-started/quick-start): a minimal app that actually boots and serves an endpoint
- [Configuration](./getting-started/configuration): what `application.toml` controls
- [Modules & Dependency Injection](./modules/overview): how your code joins the runtime
- [Routing](./guide/routing): how operations become HTTP endpoints
- [Models](./guide/models): how Bun models, audit fields, and tags work together
- [Generic CRUD](./guide/crud): how to expose typed CRUD operations with minimal glue code
- [Authentication](./security/authentication): how Bearer, Signature, and public endpoints work

If you are new to the framework, go to [Quick Start](./getting-started/quick-start) next.
