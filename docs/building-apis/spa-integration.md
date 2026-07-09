---
sidebar_position: 6
---

# SPA Integration

VEF can serve single-page applications through app middleware rather than requiring a separate frontend server in every deployment.

## The public config type

SPA integration is driven by `middleware.SPAConfig`:

- `middleware.SPAConfig`
- `SPAConfig.Path` (`Path string`): mount path for the SPA. An empty value defaults to `/`.
- `SPAConfig.Fs` (`Fs fs.FS`): file system that contains `index.html` and the `/static/*` assets.
- `SPAConfig.ExcludePaths` (`ExcludePaths []string`): path prefixes that the SPA fallback should not rewrite, such as `/api` or `/ws`.

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

  "github.com/coldsmirk/vef-framework-go"
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
- enables `etag` caching and adds security headers via `helmet`
- performs SPA-style fallback routing for non-API GET paths to `index.html`
- honors `ExcludePaths` before fallback routing so excluded prefixes keep their normal route or 404 behavior

Its `Order()` is `1000`, so it runs after API routing. This lets one VEF process own both the API and the SPA shell when that deployment model is useful.

## Exclusions

`ExcludePaths` entries are matched on path-segment boundaries. For example, `/api` excludes `/api` and `/api/users`, but it does not exclude `/apidocs`; empty exclusion prefixes are ignored, and a trailing slash in an exclusion prefix is normalized.

## Next step

Read [Modules & Dependency Injection](../core-concepts/overview) if you want to place SPA wiring into its own application module cleanly.
