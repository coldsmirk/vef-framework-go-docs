---
sidebar_position: 1
---

# API Package

The `api` package is the foundation of VEF's request handling layer. It defines the core abstractions — resources, operations, request/response types, and handler resolution — that all other packages build upon.

## API Reference

The generated [Public API Index](../reference/public-api-index) is the
no-omissions checklist for every exported symbol, exported field, and exported
method. This guide explains the supported behavior and runtime contracts behind
that surface.

API groups covered in this guide:

| Group | Public APIs |
| --- | --- |
| resource and kind | `api.Resource`, `api.Kind`, `api.KindRPC`, `api.KindREST`, `api.ValidateActionName(action, kind) error`, `api.NewRPCResource(name, opts...)`, `api.NewRESTResource(name, opts...)`, `api.WithVersion(v)`, `api.WithAuth(config)`, `api.WithOperations(specs...)` |
| engine and routing extension | `api.Engine`, `api.RouterStrategy`, `api.Middleware` |
| operations | `api.OperationSpec`, `api.Operation`, `api.RateLimitConfig`, `api.OperationsProvider`, `api.OperationsCollector` |
| request model | `api.Identifier`, `api.Request`, `api.Params`, `api.Meta`, `api.P`, `api.StrictP`, `api.M` |
| auth | `api.AuthConfig`, `api.Public()`, `api.BearerAuth()`, `api.SignatureAuth()`, `api.IPAuth(...)`, `api.APIKeyAuth(...)`, `api.HTTPBasicAuth()`, `api.AuthStrategy`, `api.AuthStrategyRegistry` |
| handler extension | `api.HandlerResolver`, `api.HandlerAdapter`, `api.HandlerParamResolver`, `api.FactoryParamResolver` |
| audit, headers, versions, errors | `api.AuditEvent`, `api.SubscribeAuditEvent`, `api.HeaderXMetaPrefix`, `api.HeaderXTimestamp`, `api.HeaderXNonce`, `api.HeaderXSignature`, `api.HeaderXAppID`, `api.VersionV1`, `api.VersionV9`, `api.ErrInvalidRequestParams`, `api.ErrInvalidRequestMeta`, `api.ErrInvalidParamsType`, `api.ErrInvalidMetaType` |

## Architecture

```
api.Engine
├── Register(resources...)       — register resources
├── Mount(router)                — attach to Fiber
└── Lookup(id)                   — find operation at runtime

api.Resource
├── Kind()        — RPC or REST
├── Name()        — resource path (e.g., "sys/user")
├── Version()     — API version (e.g., "v1")
├── Auth()        — authentication config
└── Operations()  — list of OperationSpec

api.OperationSpec → api.Operation (runtime)
```

## Resource

A `Resource` groups related API operations under a common path. VEF provides two resource kinds:

### Creating Resources

```go
// RPC resource — uses snake_case actions
resource := api.NewRPCResource("sys/user")

// REST resource — uses HTTP verbs
resource := api.NewRESTResource("sys/user")

// With options
resource := api.NewRPCResource("sys/user",
    api.WithVersion("v2"),
    api.WithAuth(api.BearerAuth()),
    api.WithOperations(
        api.OperationSpec{Action: "create", Handler: createHandler},
        api.OperationSpec{Action: "find_page", Handler: findPageHandler},
    ),
)
```

### Resource Kind

| Kind | Constant | Name Format | Action Format | Example |
| --- | --- | --- | --- | --- |
| RPC | `api.KindRPC` | `snake_case` with `/` separators | `snake_case` | `sys/user` → `create`, `find_page` |
| REST | `api.KindREST` | `kebab-case` with `/` separators | `<verb>` or `<verb> <sub-resource>` | `sys/user` → `get`, `post`, `get user-friends` |

`Kind.String()` returns `rpc` for `KindRPC`, `rest` for `KindREST`, and
`unknown` for any other value.

### Resource Name Rules

| Rule | Valid | Invalid |
| --- | --- | --- |
| Must start with lowercase letter | `user`, `sys/user` | `User`, `1user` |
| No leading/trailing slashes | `sys/user` | `/sys/user/` |
| No consecutive slashes | `sys/user` | `sys//user` |
| RPC: snake_case segments | `sys/data_dict` | `sys/data-dict` |
| REST: kebab-case segments | `sys/data-dict` | `sys/data_dict` |

### Resource Options

| Option | Description |
| --- | --- |
| `api.WithVersion(v)` | Override the engine's default version (e.g., `"v2"`) |
| `api.WithAuth(config)` | Set resource-level authentication |
| `api.WithOperations(specs...)` | Provide operation specs directly |

`api.NewRPCResource` and `api.NewRESTResource` validate the resource name,
version, and any specs passed through direct `api.WithOperations(...)` at
construction time. They panic when validation fails. `api.ValidateActionName(action, kind)`
is public for code that builds resources dynamically and wants to apply the same
RPC/REST action validation before constructing an `OperationSpec`.

REST action validation accepts these lowercase method tokens: `get`, `post`,
`put`, `delete`, `patch`, `head`, `options`, `trace`, `connect`, and `all`.
Sub-resource paths may contain `/`, but each segment must use kebab-case;
dynamic Fiber params such as `/:id` are not accepted by the public validator.

## Resource Interface

```go
type Resource interface {
    Kind() Kind
    Name() string
    Version() string
    Auth() *AuthConfig
    Operations() []OperationSpec
}
```

When using CRUD builders, you typically embed `api.Resource` and CRUD providers.
The built-in collectors read direct `Resource.Operations()` specs and anonymous
embedded `OperationsProvider` values. Direct `WithOperations(...)` specs are
validated by the resource constructor. During engine registration, collected
operation specs must have a non-empty action; custom `OperationsProvider`
implementations should still produce action strings that already satisfy
`api.ValidateActionName(...)`.

## OperationSpec

`OperationSpec` is the static definition of an API endpoint:

```go
type OperationSpec struct {
    Action             string            // Action name (e.g., "create", "find_page")
    EnableAudit        bool              // Enable audit logging
    Timeout            time.Duration     // Request timeout
    Public             bool              // No authentication required
    RequiredPermission string            // Required permission token
    RateLimit          *RateLimitConfig  // Rate limiting config
    Handler            any               // Business logic handler
}
```

At runtime the engine materializes an `Operation` with final `Auth` and
`RateLimit` pointers for router strategies, middleware, diagnostics, and tests.
`Operation` and `Request` both embed `Identifier`, so `Identifier.String()` is
promoted to `Operation.String()` and `Request.String()`.

Operation defaulting rules:

| Field | Runtime behavior |
| --- | --- |
| `Action` | required; direct `WithOperations(...)` specs are validated against the resource kind by the resource constructor; engine registration rejects an empty action |
| `EnableAudit` | copied directly to the runtime operation |
| `Timeout` | non-positive values use the engine default, which is `30s` unless overridden |
| `Public` | `true` resolves auth to `api.Public()` before resource/default auth |
| `RequiredPermission` | copied into auth options as the required permission token when non-empty |
| `RateLimit` | nil uses the engine default — `vef.api.rate_limit` when configured, else `Max=100`, `Period=5m`; a custom `RateLimitConfig` replaces the default. An explicit `Max <= 0` does **not** disable limiting — the middleware falls back to the engine/config default for non-positive values |
| `Handler` | RPC may infer from action when omitted; REST requires an explicit handler |

Permission declarations are validated at registration and a violation fails
startup:

- `RequiredPermission` must be a **dot-separated** token matching
  `^[A-Za-z0-9_]+(\.[A-Za-z0-9_]+)*$` (e.g. `sys.user.query`); colons,
  slashes, dashes, empty segments, and whitespace are rejected with
  `ErrPermissionTokenInvalid`.
- Declaring `RequiredPermission` on an operation whose resolved auth strategy
  is `none` (`Public: true`, or a resource-level `api.Public()`) is
  contradictory — the anonymous principal can never satisfy it — and is
  rejected with `ErrPermissionOnPublicOp`.

### Using with CRUD Builders

Most operations are defined through CRUD builders rather than raw `OperationSpec`:

```go
type UserResource struct {
    api.Resource

    crud.FindPage[User, UserSearch]
    crud.Create[User, UserParams]
    crud.Update[User, UserParams]
    crud.Delete[User]
}

func NewUserResource() *UserResource {
    return &UserResource{
        Resource: api.NewRPCResource("sys/user"),
        FindPage: crud.NewFindPage[User, UserSearch]().RequiredPermission("sys.user.query"),
        Create:   crud.NewCreate[User, UserParams]().RequiredPermission("sys.user.create"),
        Update:   crud.NewUpdate[User, UserParams]().RequiredPermission("sys.user.update"),
        Delete:   crud.NewDelete[User]().RequiredPermission("sys.user.delete"),
    }
}
```

### Using with Custom Handlers

For non-CRUD operations, use `WithOperations` or implement `OperationsProvider`:

```go
resource := api.NewRPCResource("sys/user",
    api.WithOperations(
        api.OperationSpec{
            Action:             "reset_password",
            Handler:            resetPasswordHandler,
            RequiredPermission: "sys.user.reset_password",
        },
    ),
)
```

## Identifier

Every operation has a unique `Identifier`:

```go
type Identifier struct {
    Resource string `json:"resource" form:"resource" validate:"required,alphanum_us_slash" label_i18n:"api_request_resource"`
    Action   string `json:"action" form:"action" validate:"required" label_i18n:"api_request_action"`
    Version  string `json:"version" form:"version" validate:"required,alphanum" label_i18n:"api_request_version"`
}

// String format: "sys/user:create:v1"
id.String()
```

The same `Identifier` fields are used by JSON RPC bodies and form RPC
requests. `Identifier.String()` always formats as
`{resource}:{action}:{version}`.

## Request / Params / Meta

### Request

The unified API request structure:

```go
type Request struct {
    Identifier
    Params Params `json:"params"`
    Meta   Meta   `json:"meta"`
}
```

Embed `api.P` in request parameter structs and `api.M` in metadata structs when
you want handler injection to decode from `params` or `meta` explicitly.

### Params

`Params` holds the business data of a request (the "what"):

```go
type Params map[string]any

// Decode into a typed struct
var userParams UserParams
err := request.Params.Decode(&userParams)

// Access individual values
value, exists := request.Params["username"]
```

### Meta

`Meta` holds request metadata (the "how" — pagination, sorting, format):

```go
type Meta map[string]any

// Decode into a typed struct
var pageable page.Pageable
err := request.Meta.Decode(&pageable)

// Access individual values
value, exists := request.Meta["format"]
```

### Params vs Meta

| Aspect | `Params` | `Meta` |
| --- | --- | --- |
| Purpose | Business data | Request control |
| Decoded into | `TParams` or `TSearch` | `page.Pageable`, `DataOptionConfig`, etc. |
| Examples | `username`, `email`, `departmentId` | `page`, `size`, `sort`, `format` |

## Authentication

### AuthConfig

```go
type AuthConfig struct {
    Strategy string         // "none", "bearer", "signature", "ip", or custom
    Options  map[string]any // Strategy-specific options
}
```

`AuthConfig.Clone()` returns a copy of the strategy/options pair. Use it when a
resource or operation customizes auth without mutating shared config.

### Built-in Auth Strategies

| Strategy | Constant | Description |
| --- | --- | --- |
| None | `api.AuthStrategyNone` | No authentication (public) |
| Bearer | `api.AuthStrategyBearer` | Bearer token authentication |
| Signature | `api.AuthStrategySignature` | Request signature authentication |
| IP | `api.AuthStrategyIP` | Source-IP whitelist authentication |
| API key | `api.AuthStrategyAPIKey` | Static API key authentication |
| HTTP Basic | `api.AuthStrategyHTTPBasic` | RFC 7617 Basic authentication |

### Helper Functions

```go
api.Public()               // AuthConfig with strategy "none"
api.BearerAuth()           // AuthConfig with strategy "bearer"
api.SignatureAuth()        // AuthConfig with strategy "signature"
api.IPAuth()               // AuthConfig with strategy "ip" and whitelist "default"
api.IPAuth("ops")          // AuthConfig with strategy "ip" and whitelist "ops"
api.APIKeyAuth()           // AuthConfig with strategy "api_key", header X-API-Key
api.APIKeyAuth("X-My-Key") // custom key header; more than one name panics
api.HTTPBasicAuth()        // AuthConfig with strategy "http_basic"
```

`api.IPAuth(...)` accepts zero or one whitelist name. With no argument it uses
`api.DefaultIPWhitelist` (`"default"`); the selected name is stored under
`api.AuthOptionWhitelist` in `AuthConfig.Options`. Passing more than one name
panics. The built-in IP strategy resolves the named list through
`security.IPWhitelistLoader`; the default loader reads
`vef.security.ip_whitelists`. All auth failures deny with
`security.ErrIPNotAllowed`, and an empty or missing named whitelist is
fail-closed rather than treated as public access. Behind a reverse proxy,
configure `vef.app.trusted_proxies` so Fiber resolves the real client IP.

`api.APIKeyAuth(...)` reads the key from `api.HeaderXAPIKey` (`X-API-Key`) by
default; one custom header name may be passed and is stored under
`api.AuthOptionAPIKeyHeader`. Keys resolve through `security.APIKeyLoader`
(default: `vef.security.api_keys`). `api.HTTPBasicAuth()` verifies
`Authorization: Basic` credentials through `security.BasicAccountLoader`
(default: `vef.security.basic_accounts`) in constant time. Both reject
uniformly with 401 (`ErrAPIKeyInvalid` / `ErrBasicCredentialsInvalid`); see
[Authentication](../security/authentication) for the loader contracts.

Custom authentication strategies implement `api.AuthStrategy` and register with
`vef.ProvideAuthStrategy(...)` into `vef:api:auth_strategies`.

### Auth at Resource vs Operation Level

```go
// Resource-level: all operations use signature auth
api.NewRPCResource("external/webhook", api.WithAuth(api.SignatureAuth()))

// Operation-level: override per operation
crud.NewCreate[User, UserParams]().Public()                    // No auth
crud.NewFindPage[User, UserSearch]().RequiredPermission("sys.user.query") // Bearer + permission
```

## Rate Limiting

```go
type RateLimitConfig struct {
    Max    int           // Maximum requests allowed
    Period time.Duration // Time window
    Key    string        // Custom rate limit key (optional)
}
```

Usage via CRUD builder:

```go
crud.NewCreate[User, UserParams]().RateLimit(100, time.Minute)
```

The built-in rate limiter uses a sliding window, counted per node. An
operation that declares no `RateLimit` of its own receives the engine default,
which is user-configurable:

```toml
[vef.api.rate_limit]
max    = 100   # default when omitted
period = "5m"  # default when omitted
```

An explicit `RateLimitConfig` with `Max <= 0` does **not** disable limiting:
the middleware only honors positive values and falls back to the
engine/config default (which is always positive) for anything else — there is
no per-operation off switch. The framework's default key includes resource,
version, action, resolved client IP, and the principal ID; anonymous requests
use the anonymous principal.

Operation auth is carried by `Operation.Auth`; a public operation resolves to
`api.AuthStrategyNone`, while protected operations carry the selected auth
strategy and options.

## Engine

The `Engine` manages resource registration and HTTP routing:

```go
type Engine interface {
    Register(resources ...Resource) error  // Add resources
    Lookup(id Identifier) *Operation       // Find operation at runtime
    Mount(router fiber.Router) error       // Attach to Fiber router
}
```

## Handler Resolution

VEF supports flexible handler signatures through parameter injection:

```go
// Minimal handler
func (r *UserResource) Create(ctx fiber.Ctx) error { ... }

// With auto-injected parameters
func (r *UserResource) Create(ctx fiber.Ctx, db orm.DB, principal *security.Principal) error { ... }

// With typed params
func (r *UserResource) Create(ctx fiber.Ctx, params UserParams, db orm.DB) error { ... }
```

### Parameter Resolution Interfaces

| Interface | Purpose |
| --- | --- |
| `HandlerParamResolver` | Resolves handler params from request context at runtime |
| `FactoryParamResolver` | Resolves handler params once at startup (dependency injection) |
| `HandlerAdapter` | Converts any handler signature to `fiber.Handler` |
| `HandlerResolver` | Finds the handler function on a resource |

`RouterStrategy` is also public for custom HTTP exposure. A strategy declares
which `api.Kind` values it can handle, receives the Fiber router in `Setup`, and
registers each resolved operation in `Route`.

## Error Types

| Error | Meaning |
| --- | --- |
| `ErrEmptyResourceName` | Resource name is empty |
| `ErrInvalidResourceName` | Resource name doesn't match naming rules |
| `ErrResourceNameSlash` | Resource name starts or ends with `/` |
| `ErrResourceNameDoubleSlash` | Resource name contains `//` |
| `ErrInvalidResourceKind` | Invalid resource kind value |
| `ErrInvalidVersionFormat` | Version doesn't match `v\d+` pattern |
| `ErrEmptyActionName` | Action name is empty |
| `ErrInvalidActionName` | Action doesn't match kind-specific rules |
| `ErrInvalidParamsType` | Params.Decode target is not a pointer to struct |
| `ErrInvalidMetaType` | Meta.Decode target is not a pointer to struct |

`Params.Decode` and `Meta.Decode` require a pointer to a struct. Passing any
other target returns `ErrInvalidParamsType` or `ErrInvalidMetaType`.

`ErrInvalidRequestParams` and `ErrInvalidRequestMeta` use
`result.ErrCodeBadRequest` (`1400`) and HTTP status `400`. They are returned
when RPC form `params`/`meta` JSON or REST JSON body parsing fails.

The remaining public surface of the `api` package — version constants, request
headers, the audit event, the auth strategy registry, operation collection
types, request helpers, and handler/router extension interfaces — is ledgered
symbol by symbol in the [Public API Index](../reference/public-api-index).

## Next Step

- Read [Generic CRUD](../data-access/crud) to learn how CRUD builders generate operations automatically
- Read [Custom Handlers](./custom-handlers) to create non-CRUD operations
- Read [Routing](./routing) for HTTP routing details
- Read [Params and Meta](./params-and-meta) for request data contracts
