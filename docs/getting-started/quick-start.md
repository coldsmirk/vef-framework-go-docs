---
sidebar_position: 2
---

# Quick Start

This quick start builds the smallest VEF application that:

- boots successfully
- registers a resource
- serves an API request through the RPC endpoint

## 1. Create `main.go`

```go
package main

import (
	"github.com/gofiber/fiber/v3"

	"github.com/coldsmirk/vef-framework-go"
	"github.com/coldsmirk/vef-framework-go/api"
	"github.com/coldsmirk/vef-framework-go/result"
)

type PingResource struct {
	api.Resource
}

func NewPingResource() api.Resource {
	return &PingResource{
		Resource: api.NewRPCResource(
			"demo/ping",
			api.WithOperations(
				api.OperationSpec{
					Action: "hello",
					Public: true,
				},
			),
		),
	}
}

func (*PingResource) Hello(ctx fiber.Ctx) error {
	return result.Ok(map[string]any{
		"message": "hello from vef",
	}).Response(ctx)
}

func main() {
	vef.Run(
		vef.ProvideAPIResource(NewPingResource),
	)
}
```

## 2. Create `configs/application.toml`

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_source]
type = "sqlite"
```

This is enough for the quick start because:

- SQLite works without external services
- storage falls back to memory when not configured
- the resource is public, so no auth provider is needed yet
- this example does not use any Redis-backed capability

## 3. Start the app

```bash
go run .
```

If startup succeeds, VEF will print the application banner and begin listening on the configured port.

## 4. Call the RPC endpoint

Send a request to `POST /api`:

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "demo/ping",
    "action": "hello",
    "version": "v1",
    "params": {},
    "meta": {}
  }'
```

Expected response:

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "message": "hello from vef"
  }
}
```

With the framework's default language, the success message is usually `成功`. If you set `VEF_I18N_LANGUAGE=en`, the same response message becomes `Success`.

## What this example demonstrates

- `vef.Run(...)` boots the full runtime
- `vef.ProvideAPIResource(...)` registers your resource into the API engine
- `api.NewRPCResource(...)` defines an RPC resource
- `api.OperationSpec` declares a public operation
- `Hello` is discovered automatically because RPC actions fall back from `hello` to the `Hello` method
- `result.Ok(...)` wraps the response into the framework result envelope

## Why the handler stays small

You did not manually:

- mount a route
- parse the request body
- validate the RPC envelope
- wire the response shape
- attach middleware

Those are handled by the framework runtime.

## Where to go next

- [Project Structure](./project-structure): how to organize a real app
- [Modules & Dependency Injection](../modules/overview): how your own modules join the boot pipeline
- [Routing](../guide/routing): how RPC and REST resources differ
