# SPA 集成

VEF 可以在同一个 Fiber 服务里同时托管 API 和单页应用。公开入口是 `middleware.SPAConfig`，以及 `di.go` 中的两个 DI helper。

## 注册方式

你可以使用：

- `vef.ProvideSPAConfig(...)`
- `vef.SupplySPAConfigs(...)`

它们都会把配置写入 `vef:spa` 这个 group，随后内部 app middleware 会把这些配置转成静态资源服务和 SPA fallback 逻辑。

## 配置结构

`middleware.SPAConfig` 包含：

- `Path`
- `Fs`
- `ExcludePaths`

其中 `Path` 是挂载路径，`Fs` 是前端构建产物所在的文件系统。

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

它的中间件顺序是 `1000`，因此会在 API 路由之后执行。

## 当前限制

`ExcludePaths` 这个字段虽然定义在公共类型中，但当前内部 middleware 实现并没有真正使用它。因此现阶段不要把它写成“已生效的排除路由能力”。

## 下一步

继续阅读 [模块与依赖注入](../modules/overview)，把 SPA wiring 放回应用模块层面来看会更清楚。
