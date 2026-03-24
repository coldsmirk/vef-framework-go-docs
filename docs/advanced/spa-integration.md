---
sidebar_position: 4
---

# SPA Integration

VEF can serve single-page applications through app middleware rather than requiring a separate frontend server in every deployment.

## The public config type

SPA integration is driven by `middleware.SPAConfig`:

- `Path`
- `Fs`
- `ExcludePaths`

## Registration helpers

Use:

- `vef.ProvideSPAConfig(...)`
- `vef.SupplySPAConfigs(...)`

## Minimal example

```go
package web

import (
  "embed"
  "io/fs"

  vef "github.com/coldsmirk/vef-framework-go"
  "github.com/coldsmirk/vef-framework-go/middleware"
)

//go:embed dist/*
var webFS embed.FS

func NewWebConfig() *middleware.SPAConfig {
  sub, _ := fs.Sub(webFS, "dist")

  return &middleware.SPAConfig{
    Path: "/",
    Fs:   sub,
  }
}

var Module = vef.Module(
  "app:web",
  vef.ProvideSPAConfig(NewWebConfig),
)
```

If your frontend build artifacts are already exposed as a file system from another package, the module can be even smaller:

```go
var Module = vef.Module(
  "app:web",
  vef.SupplySPAConfigs(&middleware.SPAConfig{
    Fs: dist.FS,
  }),
)
```

These register SPA configs into the `vef:spa` group, and the app middleware module picks them up automatically.

## Middleware behavior

The SPA middleware:

- serves the app entry at the configured path
- serves static assets under the configured path prefix at `/static/*`
- performs SPA-style fallback routing for non-API GET paths

This lets one VEF process own both the API and the SPA shell when that deployment model is useful.

## Current limitation

`ExcludePaths` exists on the public config type, but the current internal SPA middleware does not actively use that field yet. Do not document it as an already-enforced exclusion mechanism.

## Next step

Read [Modules & Dependency Injection](../modules/overview) if you want to place SPA wiring into its own application module cleanly.
