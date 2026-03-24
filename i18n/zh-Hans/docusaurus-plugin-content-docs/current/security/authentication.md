---
sidebar_position: 1
---

# 认证

VEF 的认证发生在 API 操作层。每个操作都有自己的 auth 配置，API 中间件会在 handler 执行前先解析出当前 principal。

## 默认行为

如果你不做额外配置：

- 操作默认使用 Bearer 认证
- 显式 `Public` 的操作则不要求认证

这个默认值来自 API 引擎，而不是来自你的应用配置文件。

## 内置认证策略

公开的 `api` 包暴露了三种策略 helper：

- `api.Public()`
- `api.BearerAuth()`
- `api.SignatureAuth()`

实际使用中，你通常通过操作配置来控制：

```go
api.OperationSpec{
	Action: "login",
	Public: true,
}
```

或者通过资源级别的 auth 默认值来设置。

## Bearer 认证

Bearer 认证支持两种 token 来源：

- `Authorization: Bearer <token>`
- 查询参数 `__accessToken`

真正的 token 校验逻辑由安全模块中的 auth manager 负责。

## Signature 认证

Signature 认证主要用于外部应用和请求签名场景。

它要求这些 header：

- `X-App-ID`
- `X-Timestamp`
- `X-Nonce`
- `X-Signature`

校验逻辑由安全模块的 signature authenticator 执行。

## 公开操作

公开操作会得到一个匿名 principal，而不是直接被拒绝。

适合标记为 `Public` 的接口包括：

- 登录
- 刷新 token
- 某些匿名健康检查或回调入口

## 内置认证资源

安全模块会自动注册一个内置 RPC 资源：

```text
security/auth
```

主要 actions 包括：

- `login`
- `refresh`
- `logout`
- `resolve_challenge`
- `get_user_info`

## 登录流程

内置认证资源支持两阶段模型：

1. 先校验凭证
2. 如有需要，再进入 challenge 流程

如果没有 challenge，`login` 会直接返回 token。

如果 challenge provider 已配置且当前用户需要额外挑战，`login` 会返回：

- challenge token
- 下一步 challenge 描述

客户端之后继续调用 `resolve_challenge`，直到所有挑战都完成。

## 应用通常还需要提供什么

这里要分场景来看：

- `security.UserLoader` 通常是用户登录和 refresh 流程的前提
- `security.ExternalAppLoader` 只在你使用签名认证的外部应用场景时需要
- challenge provider 是可选项，只有在你启用了挑战式登录流时才相关
- `security.UserInfoLoader` 只在你希望 `security/auth.get_user_info` 返回应用自定义用户信息时需要

框架提供的是认证流程和中间件，而不是你的身份源本身。

## 真实项目里常见的 auth 模块形态

很多 VEF 项目里的 auth 模块其实都很小，常见模式是：

- 在 `init()` 里设置用户详情类型
- 提供 `UserLoader`
- 提供 `UserInfoLoader`

概念上通常像这样：

```go
func init() {
  security.SetUserDetailsType[*UserDetails]()
}

var Module = vef.Module(
  "app:auth",
  vef.Provide(
    NewUserLoader,
    NewUserInfoLoader,
  ),
)
```

这样认证接入代码就能和业务资源模块保持分离。

## 实践建议

- `Public` 只用于明确需要匿名访问的操作
- 普通用户认证优先保持在 Bearer
- Signature 更适合系统对系统集成，而不是替代普通用户会话

## 下一步

继续阅读 [授权](./authorization)，看认证之后权限检查是如何继续发生的。
