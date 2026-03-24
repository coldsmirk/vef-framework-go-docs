---
sidebar_position: 1
---

# Modules & Dependency Injection

VEF is built on Uber FX. The public `vef` package re-exports the core FX helpers so that most applications can stay inside one consistent API surface.

## The key idea

You do not bootstrap subsystems manually. Instead, you compose them through FX options:

```go
vef.Run(
  user.Module,
  auth.Module,
  vef.ProvideAPIResource(resources.NewHealthResource),
)
```

Internally, `vef.Run(...)` appends your options to the framework’s own module list and starts the FX application.

In a real application, `main.go` often looks more like this:

```go
vef.Run(
  ivef.Module,  // framework-facing integration
  mcp.Module,   // custom MCP providers
  web.Module,   // SPA hosting
  auth.Module,  // auth loaders
  sys.Module,   // system/admin resources
  md.Module,    // master data resources
  pmr.Module,   // business resources
)
```

That structure reflects an important pattern: VEF apps usually compose several small domain or integration modules rather than one giant app module.

## Helpers re-exported by `vef`

The `vef` package re-exports the FX primitives you use most often:

- `vef.Run`
- `vef.Module`
- `vef.Provide`
- `vef.Supply`
- `vef.Invoke`
- `vef.Decorate`
- `vef.Replace`

This keeps most application code from importing `fx` directly unless you need something more specific.

## Group-based extension points

Several framework features are connected through FX groups. These are the most important ones for application developers:

- `vef:api:resources`
- `vef:app:middlewares`
- `vef:cqrs:behaviors`
- `vef:security:challenge_providers`
- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`

The helper functions in `di.go` exist mainly to register values into those groups safely.

## API resources

The most common helper is:

```go
vef.ProvideAPIResource(NewUserResource)
```

It annotates the constructor result into the API resource group. During startup, the API module collects every resource from that group and registers their operations into the engine.

## Middleware and other providers

The same pattern is used for other extension points:

```go
vef.ProvideMiddleware(NewAuditTrailMiddleware)
vef.ProvideCQRSBehavior(NewTracingBehavior)
vef.ProvideChallengeProvider(NewTOTPChallengeProvider)
vef.ProvideMCPTools(NewToolProvider)
```

Each helper hides the FX group tag so that application code stays easier to read.

## Module roles in a larger app

A production-style VEF app often has a few recurring module types:

- an `internal/vef` module for build info, shared framework-facing services, and event subscribers
- an `internal/auth` module for `UserLoader`, `UserInfoLoader`, and auth-specific setup
- one or more business-domain modules that register API resources
- optional `web` and `mcp` modules for SPA and MCP integration

This keeps responsibilities obvious and prevents domain modules from becoming catch-all wiring buckets.

## Invoke-driven integration modules

In larger applications, a dedicated integration module often uses `vef.Invoke(...)` for startup-time wiring that does not belong to any one business resource.

Example:

```go
var Module = vef.Module(
  "app:vef",
  vef.Supply(BuildInfo),
  vef.Provide(NewDataDictLoader, password.NewBcryptEncoder),
  vef.Invoke(registerEventSubscribers),
)
```

This is a good home for build metadata, shared framework-facing services, and event-subscriber registration.

## Why resources are discovered automatically

The API engine does not need you to mount routes by hand. Instead it:

1. collects resources from the DI container
2. collects operation specs from the resource itself and from embedded CRUD providers
3. resolves handlers
4. adapts handlers into Fiber handlers
5. mounts them into the RPC or REST router

That design is why VEF code usually looks like “define resource + register constructor” instead of “declare router + bind handler + wire middleware”.

## When to use plain `fx`

Most apps can stay inside the `vef` wrapper, but using plain `fx` is still valid when needed, for example:

- advanced annotations
- optional dependencies
- direct access to lifecycle hooks
- test-only overrides

VEF does not prevent direct FX usage. It just tries to make the common paths shorter.

## Next step

Continue to [Application Lifecycle](./lifecycle) to see what happens when `vef.Run(...)` starts the system.
