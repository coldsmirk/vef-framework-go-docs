---
sidebar_position: 2
---

# 应用生命周期

这一页解释从 `vef.Run(...)` 到 HTTP 服务真正开始监听，中间到底发生了什么。

## 启动顺序

当前框架在 `bootstrap.go` 里定义的启动顺序是：

`config -> database -> orm -> middleware -> api -> security -> event -> cqrs -> cron -> redis -> mold -> storage -> sequence -> schema -> monitor -> mcp -> app`

这个顺序不是随意排的，后面的模块依赖前面的模块：

- 配置必须最先可用
- 数据库和 ORM 必须先于需要 `orm.DB` 的 API 处理器
- 安全模块必须先于受保护请求
- storage、monitor、schema、MCP 等能力要先注册，再启动 app

## `vef.Run(...)` 实际做的事情

从行为上看，它会：

1. 构建框架默认模块列表
2. 追加你自己的 FX option
3. 附加 `startApp`
4. 创建 FX app
5. 运行它

默认启动超时是 `30s`，默认停止超时是 `60s`。

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

- database 会先 ping 连接并输出数据库版本
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

接下来建议看 [路由](../guide/routing)，理解资源是如何真正变成 HTTP 端点的。
