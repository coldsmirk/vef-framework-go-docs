---
sidebar_position: 1
slug: /intro
---

# 简介

VEF Framework 是一个围绕“资源模型”组织应用的 Go Web 框架，底层基于 Uber FX、Fiber 和 Bun，适合搭建内部平台、管理后台和服务型 API。

理解 VEF 最有效的方式，不是先看路由，而是先看运行时形状：

1. `vef.Run(...)` 启动一条固定的模块装配链。
2. 你通过 FX group 注册自己的资源、中间件和行为。
3. API 引擎从这些资源里收集操作，再挂载成 RPC 或 REST 端点。
4. handler 参数由框架自动注入，来源包括请求上下文、解码后的 params/meta，以及容器里的服务。

所以，VEF 本质上不是一个“路由小工具”，也不只是一个“CRUD 库”，而是一套围绕显式资源和稳定默认值组织应用的框架。

## 框架默认替你做了什么

从当前框架实现来看，启动顺序是：

`config -> database -> orm -> middleware -> api -> security -> event -> cqrs -> cron -> redis -> mold -> storage -> sequence -> schema -> monitor -> mcp -> app`

这意味着框架默认已经对很多运行时行为做了约定：

- API 默认版本：`v1`
- API 默认认证：Bearer token
- 默认请求超时：`30s`
- 默认限流：`100` 次请求 / `5m`
- 默认响应信封：`result.Result`
- 默认存储提供者：未配置时回退到内存存储

这些结论来自运行时真实行为，而不是文档层面的约定描述。

## RPC 和 REST 都是一等能力

VEF 同时支持两种 API 风格：

- RPC 资源，统一走 `POST /api`
- REST 资源，挂载到 `/api/<resource>`

它们必须显式声明：

```go
api.NewRPCResource("sys/user", ...)
api.NewRESTResource("users", ...)
```

VEF **不会** 从一个 RPC 资源自动生成 REST 路由。如果你需要两套风格，就应当明确地定义两套资源。

## 开发者最常接触的 API

大多数业务开发只会频繁接触这几个公开入口：

- `vef.Run(...)`
- `vef.Module(...)`
- `vef.ProvideAPIResource(...)`
- `api.NewRPCResource(...)`
- `api.NewRESTResource(...)`
- `api.OperationSpec`
- `crud.NewCreate(...)`、`crud.NewFindPage(...)` 等 CRUD builder
- `orm.DB`
- `result.Ok(...)` 和 `result.Err(...)`
- `security` 包中的扩展接口，比如 `UserLoader`、`PermissionChecker`、`RolePermissionsLoader`

其余内部模块，大多是在为这些用户侧入口服务。

## 开箱即用的内置资源

框架当前已经内置了几类可直接使用的资源和模块：

- `security/auth`：登录、刷新令牌、登出、挑战流程解析、可选的用户信息读取
- `sys/storage`：上传、预签名 URL、临时文件清理、对象列表和对象元信息
- `sys/schema`：数据库 schema 检查
- `sys/monitor`：运行时与主机监控
- MCP 中间件与 MCP server 集成

如果这些能力符合你的需求，通常不需要从零重写。

## 这套文档怎么读

文档按开发者最常见的上手路径组织：

- [安装](./getting-started/installation)：环境和依赖准备
- [快速开始](./getting-started/quick-start)：真正跑起一个最小应用
- [配置](./getting-started/configuration)：`application.toml` 控制了什么
- [模块与依赖注入](./modules/overview)：你的代码如何接入运行时
- [路由](./guide/routing)：操作是如何变成 HTTP 端点的
- [模型](./guide/models)：Bun 模型、审计字段和标签如何配合
- [泛型 CRUD](./guide/crud)：如何最少样板暴露标准接口
- [认证](./security/authentication)：Bearer、Signature 与公开接口的真实行为

如果你是第一次接触 VEF，建议下一步直接看 [快速开始](./getting-started/quick-start)。
