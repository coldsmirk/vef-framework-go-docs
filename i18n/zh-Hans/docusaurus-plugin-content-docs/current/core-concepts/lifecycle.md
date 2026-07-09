---
sidebar_position: 2
---

# 应用生命周期

这一页解释从 `vef.Run(...)` 到 HTTP 服务真正开始监听，中间到底发生了什么。

## 启动顺序

这是 VEF 启动流水线的权威表述，来自 `bootstrap.go`（`vef.Run` 直接装配
`config` 和 `datasource`）以及 `internal/bootmodules.Core()`（业务模块的顺序
列表，由 `vef.Run` 和 `internal/apptest` 测试脚手架共用，保证两个 FX 图不会
出现分叉）：

`config -> datasource -> middleware -> api -> security -> event -> expression -> cqrs -> cron -> redis -> mold -> storage -> sequence -> outbox -> redis-stream -> inbox -> schema -> monitor -> mcp -> app`

`datasource` 是单独的一步：它在同一个模块里把 `*sql.DB` 连接起来（通过
`internal/database`）并包装成 `orm.DB`（通过 `internal/orm`）——启动流程里
并不存在独立的 `database` 或 `orm` 步骤。`outbox`、`redis-stream`、`inbox`
是事件传输子模块——即 outbox 传输模块、redis-stream 传输模块和 inbox
模块——排在 `sequence` 之后、`schema` 之前注册。

这个顺序不是随意排的，后面的模块依赖前面的模块：

- 配置必须最先可用
- datasource 必须先于需要 `orm.DB` 的 API 处理器
- 安全模块必须先于受保护请求
- 事件传输子模块（outbox、redis-stream、inbox）排在核心 `event` 模块之后，
  但先于 schema、monitor、MCP
- storage、monitor、schema、MCP 等能力要先注册，再启动 app

## `vef.Run(...)` 实际做的事情

`vef.Run(...)` 按下面的顺序组装 FX app：

1. 用 `fx.WithLogger(newFxLogger)` 安装框架 FX logger
2. 添加 internal config module（内部配置模块）
3. 添加 internal datasource module（内部数据源模块）
4. 追加 `bootmodules.Core()` 返回的全部 option
5. 追加用户传入的 `options...`
6. 追加 `fx.Invoke(startApp)`
7. 追加 `fx.StartTimeout(defaultTimeout)`
8. 追加 `fx.StopTimeout(defaultTimeout*2)`
9. 用 `fx.New(opts...)` 创建 app
10. 用 `app.Run()` 运行它

`defaultTimeout` 是 `30 * time.Second`，所以默认启动超时是 `30s`，
默认停止超时是 `60s`。

因为用户 option 是追加在 `bootmodules.Core()` 之后的，所以应用模块可以通过
`vef.ProvideAPIResource(...)` 这类 helper 继续追加 group 成员。如果要替换
core 已经提供的单例，通常要用 `vef.Decorate(...)`、`vef.Replace(...)`，
或者 `vef.SupplyFileACL(...)` 这类框架 replacement helper；为同一个 service
再注册一个普通 `vef.Provide(...)` 并不会形成 override。

高级模块可以接收 `vef.Lifecycle` 并调用 `Lifecycle.Append(...)` 直接注册
`fx.Hook`；`vef.StartHook`、`vef.StopHook` 和 `vef.StartStopHook` 是这些
hook 的便利构造函数。
`Lifecycle.Append` 的精确签名记录在 public API index 中。

内部的 `startApp` invoke 会在模块图构建完成后追加 HTTP server 的 lifecycle
hook。它的 `OnStart` 等待 `application.Start()` 或启动 context 超时；`OnStop`
调用 `application.Stop()`。

最小启动示例：

```go
func main() {
  vef.Run(
    ivef.Module,
    auth.Module,
    sys.Module,
    web.Module,
  )
}
```

## App 启动阶段

应用模块会先创建 Fiber app，然后按顺序：

1. 应用“前置”中间件
2. 挂载 API engine
3. 应用“后置”中间件

所以 VEF 的中间件顺序有两层：

- 包裹整个 Fiber app 的 app-level middleware
- 进入 API engine 后才运行的 api-level middleware

## App 级中间件顺序

根据当前实现，常见顺序大致是：

- compression
- headers
- CORS
- content type 检查
- request ID
- request logger 绑定
- panic recovery
- request record logging
- API 路由
- MCP 端点中间件
- storage 文件代理路由
- SPA fallback middleware

这意味着哪怕请求根本没有进入 API engine，某些 app 级中间件也一样会执行。

## API 级中间件顺序

进入 API engine 后，请求链当前是：

- auth
- contextual
- data permission
- rate limit
- audit
- handler

这个顺序决定了 handler 里能直接拿到：

- 已认证的 principal
- request-scoped `orm.DB`
- request-scoped logger
- 已解析好的数据权限 applier

## 内置模块自己的启动 hook

一些模块还会在生命周期里做额外事情：

- datasource 会先 ping 主数据源连接并输出数据库版本
- event bus 会启动内存事件分发器
- storage 会初始化需要启动动作的 provider
- app 会真正启动 HTTP server 并注册 stop hook

所以“应用成功启动”不只是 Fiber 起了，而是整条运行时依赖链都准备好了。

## 为什么这对排错很重要

如果 VEF 应用启动失败，常见原因通常在这些层级里：

- 配置文件没找到
- 数据库配置无效
- provider 配置不支持
- FX 构造函数注册错误
- handler factory 解析失败

理解启动顺序以后，能更快判断故障最可能出在哪一层。

## 下一步

接下来建议看 [路由](../building-apis/routing)，理解资源是如何真正变成 HTTP 端点的。
