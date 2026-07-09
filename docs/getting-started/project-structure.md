---
sidebar_position: 4
---

# Project Structure

VEF does not force a single application layout, but the framework strongly encourages a modular structure that mirrors how the runtime is assembled.

## A practical starting layout

```text
my-app/
├── configs/
│   └── application.toml
├── internal/
│   ├── user/
│   │   ├── module.go
│   │   ├── model/
│   │   ├── payload/
│   │   └── resource/
│   ├── auth/
│   │   ├── module.go
│   │   ├── user_loader.go
│   │   └── permission_loader.go
│   └── app/
│       └── module.go
└── main.go
```

This layout works well because it matches the way you register things with FX.

## A production-style VEF layout

In a larger real-world application, the module graph often becomes more explicit. A structure like this is common:

```text
my-app/
├── cmd/
│   └── server/
│       └── main.go
├── configs/
│   └── application.toml
├── internal/
│   ├── vef/    # framework-facing integration module
│   ├── auth/   # UserLoader, UserInfoLoader, auth-specific setup
│   ├── web/    # SPA hosting
│   ├── mcp/    # MCP tool/resource providers
│   ├── sys/    # system/admin domain resources
│   ├── md/     # master data resources
│   └── pmr/    # business domain resources
└── go.mod
```

This pattern is useful because it separates:

- framework integration concerns
- auth and identity loading
- frontend hosting
- business domains

## Organize by domain, not by framework layer

A good rule is:

- keep model, payload, service, and resource packages close to the business domain they belong to
- expose one `Module` per domain or subdomain
- let `main.go` compose those modules

For example:

```go
package main

import (
	"github.com/coldsmirk/vef-framework-go"

	"example.com/my-app/internal/auth"
	"example.com/my-app/internal/user"
)

func main() {
	vef.Run(
		auth.Module,
		user.Module,
	)
}
```

## What belongs in a module

A typical module registers one or more of the following:

- API resources
- domain services
- middleware
- permission loaders or auth-related providers
- CQRS handlers or behaviors

Example:

```go
package user

import (
	"github.com/coldsmirk/vef-framework-go"

	"example.com/my-app/internal/user/resource"
)

var Module = vef.Module(
	"app:user",
	vef.ProvideAPIResource(resource.NewUserResource),
)
```

In larger apps, a dedicated integration module is also common. That module may:

- `vef.Supply(...)` build info
- `vef.Provide(...)` framework-facing loaders or shared services
- `vef.Invoke(...)` event subscriber registration

That keeps framework integration code out of domain modules.

## Recommended subpackages

These package names are common and easy to scale:

- `model`: Bun models and persistence-facing types
- `payload`: request params, search params, and transport-facing DTOs
- `resource`: VEF API resources
- `service`: application or domain services
- `query` / `command`: CQRS-focused code if you adopt that style

Prefer singular package names by default. In Go, package names are usually read as a namespace rather than a collection, so `model`, `payload`, `resource`, and `service` are generally more idiomatic than `models`, `payloads`, `resources`, and `services`.

You do not need every subpackage on day one. Start small and split only when the domain grows.

## Where auth integrations usually live

Authentication and authorization often need application-specific loaders:

- `security.UserLoader`
- `security.UserInfoLoader`
- `security.RolePermissionsLoader`
- `security.ExternalAppLoader`

Put these in an `auth` or `security` domain package that your application owns.

## Where frontend assets fit

If your application serves a single-page app through VEF, keep those assets or the embedded file system adapter in a dedicated frontend or `web` module rather than mixing them into API resource packages.

The same idea applies to MCP. If your app provides custom MCP tools or resources, a dedicated `internal/mcp` module keeps those providers separate from ordinary API resources.

## What to avoid

Avoid organizing the entire application as one giant framework bucket:

- one package for every model in the system
- one package for every resource in the system
- one package for every service in the system

That structure scales poorly because VEF applications are usually extended by feature modules, not by framework layer buckets.

## Next step

Continue to [Modules & Dependency Injection](../core-concepts/overview) to see how those modules are composed at runtime.
