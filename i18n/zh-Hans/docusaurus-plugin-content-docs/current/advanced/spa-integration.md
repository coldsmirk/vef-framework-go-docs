---
sidebar_position: 4
---

# SPA 集成

VEF 可以在同一个 Fiber 服务里同时托管 API 和单页应用。公开入口是 `middleware.SPAConfig`，以及 `di.go` 中的两个 DI helper。

## 注册方式

你可以使用：

- `vef.ProvideSPAConfig(...)`
- `vef.SupplySPAConfigs(...)`

它们都会把配置写入 `vef:spa` 这个 group，随后内部 app middleware 会把这些配置转成静态资源服务和 SPA fallback 逻辑。

## 配置结构

`github.com/coldsmirk/vef-framework-go/middleware` 的已审查公开表面是：

- 1 个 exported top-level symbol、3 个 exported fields、0 个 exported methods
- fingerprint `2b84c28d70bf6ca996dc263173670f3e39ddf36493968f0d68c951af3d6bdb8e`

`middleware.SPAConfig` 的 canonical public entries 是：

- `middleware.SPAConfig`
- `SPAConfig.Path`（`Path string`）：SPA 的挂载路径；空值会默认成 `/`。
- `SPAConfig.Fs`（`Fs fs.FS`）：包含 `index.html` 和 `/static/*` 资源的文件系统。
- `SPAConfig.ExcludePaths`（`ExcludePaths []string`）：SPA fallback 不应改写的路径前缀，例如 `/api` 或 `/ws`。

## 示例

```go
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

如果前端构建产物已经由别的包直接导出了文件系统，那么模块还能更简单：

```go
var Module = vef.Module(
  "app:web",
  vef.SupplySPAConfigs(&middleware.SPAConfig{
    Fs: dist.FS,
  }),
)
```

## 框架替你做了什么

内部 SPA middleware 会：

- 在指定路径提供 `index.html`
- 在 `Path` 前缀下提供 `/static/*` 静态资源
- 打开 `etag`
- 通过 `helmet` 添加安全头
- 对同一路径前缀下的未知 GET 路由自动 fallback 回 SPA 入口
- 在 fallback 之前遵守 `ExcludePaths`，让被排除的前缀保留正常路由或 404 行为

它的中间件顺序是 `1000`，因此会在 API 路由之后执行。

## 排除路径

`ExcludePaths` 按路径段边界匹配（path-segment boundaries）。例如 `/api` 会排除 `/api` 和 `/api/users`，但不会排除 `/apidocs`。空排除前缀会被忽略（empty exclusion prefixes are ignored），排除前缀末尾的 `/`（trailing slash）会被规范化。

## 下一步

继续阅读 [模块与依赖注入](../modules/overview)，把 SPA wiring 放回应用模块层面来看会更清楚。
