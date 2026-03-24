---
sidebar_position: 3
---

# Custom Handlers

CRUD builders cover a lot of standard endpoints, but real applications always need business-specific operations. VEF lets you expose custom handlers for both RPC and REST resources and inject a wide range of request-scoped or startup-scoped dependencies directly into handler signatures.

## Handler Resolution Overview

| Resource kind | Handler source | Resolution rule |
| --- | --- | --- |
| RPC | omitted `Handler` | action name is converted from `snake_case` to a PascalCase method name on the resource |
| RPC | explicit `Handler: "MethodName"` | the named method is looked up on the resource |
| RPC | explicit function value | the function is used directly |
| REST | explicit `Handler` required | REST does not infer handler methods from the action string |

## RPC Handler Resolution

For RPC resources, you can omit the explicit handler and let the framework resolve it from the action name.

Example:

```go
type UserResource struct {
	api.Resource
}

func NewUserResource() api.Resource {
	return &UserResource{
		Resource: api.NewRPCResource(
			"sys/user",
			api.WithOperations(
				api.OperationSpec{Action: "ping", Public: true},
			),
		),
	}
}

func (*UserResource) Ping(ctx fiber.Ctx) error {
	return result.Ok("pong").Response(ctx)
}
```

RPC resolution rules:

| Action name | Resolved method |
| --- | --- |
| `ping` | `Ping` |
| `find_page` | `FindPage` |
| `get_user_info` | `GetUserInfo` |
| `resolve_challenge` | `ResolveChallenge` |

Requirements:

- RPC action names must use `snake_case`
- inferred handler methods must use PascalCase
- if you want to bypass inference, set `Handler` explicitly in `api.OperationSpec`

## REST Handler Resolution

For REST resources, you must provide the handler explicitly. The action string defines the HTTP method and optional sub-path.

```go
api.NewRESTResource(
	"users",
	api.WithOperations(
		api.OperationSpec{
			Action:  "get",
			Public:  true,
			Handler: "List",
		},
		api.OperationSpec{
			Action:  "post admin",
			Handler: "CreateAdmin",
		},
	),
)
```

REST action format:

| Format | Meaning | Examples |
| --- | --- | --- |
| `<method>` | root resource route | `get`, `post`, `delete` |
| `<method> <sub-resource>` | sub-resource route | `get profile`, `post admin`, `get user-friends` |

Rules:

- HTTP verbs stay lowercase
- sub-resources use kebab-case
- REST action strings do not infer handler method names

## Supported Handler Return Shapes

Handlers may return either nothing or one `error`.

| Shape | Meaning |
| --- | --- |
| `func(...)` | no explicit error result |
| `func(...) error` | standard handler shape |

Any other return shape is invalid for a direct handler.

## Supported Factory Shapes

VEF also supports startup-time handler factories. A factory returns a handler closure and optionally an `error`.

| Shape | Meaning |
| --- | --- |
| `func(...) func(...) error` | factory returning a handler |
| `func(...) (func(...) error, error)` | factory returning a handler plus startup error |
| `func(...) func(...)` | factory returning a no-error handler |
| `func(...) (func(...), error)` | factory returning a no-error handler plus startup error |

Factories are useful when some dependencies should be resolved once at startup instead of on every request.

## Built-In Handler Parameters

The framework ships these built-in handler parameter resolvers:

| Parameter type | Source | Typical use |
| --- | --- | --- |
| `fiber.Ctx` | request context | direct access to Fiber request/response APIs |
| `orm.DB` | request context | query and transaction entry point |
| `log.Logger` | request context | request-aware logging |
| `*security.Principal` | request context | current authenticated principal |
| `api.Params` | request body/query/form | raw params map |
| `api.Meta` | request meta | raw meta map |
| typed struct embedding `api.P` | request params | strongly typed params decoding + validation |
| typed struct embedding `api.M` | request meta | strongly typed meta decoding + validation |
| `page.Pageable` | request meta | paging metadata helper |
| `cron.Scheduler` | DI container | scheduler access |
| `event.Publisher` | DI container | event publishing |
| `mold.Transformer` | DI container | output transformation |
| `storage.Service` | DI container | storage access |

Important decoding rules:

- typed params and typed meta structs are automatically validated with `validator.Validate(...)`
- `page.Pageable` is treated as a built-in meta helper
- `api.Params` and `api.Meta` bypass typed decoding and validation

## Field And Resource Injection

If a handler parameter type is not covered by built-in resolvers, VEF tries to resolve it from the resource struct itself.

Resolution order:

| Step | How it works |
| --- | --- |
| direct field match | finds a non-embedded field on the resource with a compatible type |
| tagged dive field match | searches nested fields under `api:"dive"` |
| embedded field match | searches anonymous embedded fields |

That means custom services, repositories, or helpers can be injected by storing them as fields on the resource.

## Supported Factory Parameters

Startup-time factory resolvers are narrower than request-time handler resolvers. Built-in factory parameter types include:

| Parameter type | Source |
| --- | --- |
| `orm.DB` | DI container |
| `cron.Scheduler` | DI container |
| `event.Publisher` | DI container |
| `mold.Transformer` | DI container |
| `storage.Service` | DI container |

Factory parameter resolution also supports compatible fields on the resource struct.

## Common Handler Shapes

### Thin RPC handler

```go
func (*UserResource) ResetPassword(
	ctx fiber.Ctx,
	db orm.DB,
	principal *security.Principal,
	params *ResetPasswordParams,
) error {
	// business logic
	return result.Ok().Response(ctx)
}
```

### Raw proxy-style handler

```go
func (*DebugResource) Echo(params api.Params, meta api.Meta) error {
	return result.Ok(fiber.Map{
		"params": params,
		"meta":   meta,
	}).Response(ctx)
}
```

### Startup-time factory

```go
func (*UserResource) BuildReport(
	db orm.DB,
) (func(ctx fiber.Ctx, params ReportParams) error, error) {
	return func(ctx fiber.Ctx, params ReportParams) error {
		return result.Ok().Response(ctx)
	}, nil
}
```

## When To Use Custom Handlers

Use custom handlers when:

- the operation is not a standard CRUD action
- the endpoint coordinates multiple services
- the response shape is domain-specific
- the action triggers workflows, events, or external integrations
- the request contract is simpler to express directly than through a CRUD builder

## Mixing CRUD And Custom Actions

A common real-world pattern is to combine:

- embedded CRUD builders for standard actions
- explicit `api.OperationSpec` entries for the few domain-specific actions that do not fit CRUD

This keeps one resource aligned with one business area while still allowing extra actions such as “save permissions”, “publish version”, or “find related users”.

## Practical Advice

- keep handlers thin and business-focused
- prefer typed params structs over raw `api.Params`
- use raw maps only for dynamic or proxy-style endpoints
- use factories only when startup-time dependency shaping is actually useful
- for REST resources, always be explicit about handler mapping
- if a parameter can be modeled as typed `api.P` / `api.M`, prefer that over manual extraction from `fiber.Ctx`

## Next Step

Read [Routing](./routing) and [Parameters And Metadata](./params-and-meta) for exact request shapes and decoding rules.
