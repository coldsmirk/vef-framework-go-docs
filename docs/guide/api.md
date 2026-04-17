---
sidebar_position: 1
---

# API Package

The `api` package is the foundation of VEF's request handling layer. It defines the core abstractions — resources, operations, request/response types, and handler resolution — that all other packages build upon.

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

When using CRUD builders, you typically embed `api.Resource` and CRUD providers. The framework automatically collects operations from all embedded `OperationsProvider` implementations.

## OperationSpec

`OperationSpec` is the static definition of an API endpoint:

```go
type OperationSpec struct {
    Action      string             // Action name (e.g., "create", "find_page")
    EnableAudit bool               // Enable audit logging
    Timeout     time.Duration      // Request timeout
    Public      bool               // No authentication required
    PermToken   string             // Required permission token
    RateLimit   *RateLimitConfig   // Rate limiting config
    Handler     any                // Business logic handler
}
```

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
        FindPage: crud.NewFindPage[User, UserSearch]().PermToken("sys:user:query"),
        Create:   crud.NewCreate[User, UserParams]().PermToken("sys:user:create"),
        Update:   crud.NewUpdate[User, UserParams]().PermToken("sys:user:update"),
        Delete:   crud.NewDelete[User]().PermToken("sys:user:delete"),
    }
}
```

### Using with Custom Handlers

For non-CRUD operations, use `WithOperations` or implement `OperationsProvider`:

```go
resource := api.NewRPCResource("sys/user",
    api.WithOperations(
        api.OperationSpec{
            Action:  "reset_password",
            Handler: resetPasswordHandler,
            PermToken: "sys:user:reset_password",
        },
    ),
)
```

## Identifier

Every operation has a unique `Identifier`:

```go
type Identifier struct {
    Resource string  // e.g., "sys/user"
    Action   string  // e.g., "create"
    Version  string  // e.g., "v1"
}

// String format: "sys/user:create:v1"
id.String()
```

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

### Params

`Params` holds the business data of a request (the "what"):

```go
type Params map[string]any

// Decode into a typed struct
var userParams UserParams
err := request.Params.Decode(&userParams)

// Access individual values
value, exists := request.GetParam("username")
```

### Meta

`Meta` holds request metadata (the "how" — pagination, sorting, format):

```go
type Meta map[string]any

// Decode into a typed struct
var pageable page.Pageable
err := request.Meta.Decode(&pageable)

// Access individual values
value, exists := request.GetMeta("format")
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
    Strategy string         // "none", "bearer", "signature", or custom
    Options  map[string]any // Strategy-specific options
}
```

### Built-in Auth Strategies

| Strategy | Constant | Description |
| --- | --- | --- |
| None | `api.AuthStrategyNone` | No authentication (public) |
| Bearer | `api.AuthStrategyBearer` | Bearer token authentication |
| Signature | `api.AuthStrategySignature` | Request signature authentication |

### Helper Functions

```go
api.Public()        // AuthConfig with strategy "none"
api.BearerAuth()    // AuthConfig with strategy "bearer"
api.SignatureAuth() // AuthConfig with strategy "signature"
```

### Auth at Resource vs Operation Level

```go
// Resource-level: all operations use signature auth
api.NewRPCResource("external/webhook", api.WithAuth(api.SignatureAuth()))

// Operation-level: override per operation
crud.NewCreate[User, UserParams]().Public()                    // No auth
crud.NewFindPage[User, UserSearch]().PermToken("sys:user:query") // Bearer + permission
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

## Next Step

- Read [Generic CRUD](./crud) to learn how CRUD builders generate operations automatically
- Read [Custom Handlers](./custom-handlers) to create non-CRUD operations
- Read [Routing](./routing) for HTTP routing details
- Read [Params and Meta](./params-and-meta) for request data contracts
