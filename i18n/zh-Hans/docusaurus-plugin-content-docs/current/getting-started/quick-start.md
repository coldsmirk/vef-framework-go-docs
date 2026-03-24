---
sidebar_position: 2
---

# 快速开始

这个快速开始会搭一个最小但真实可运行的 VEF 应用，它能做到：

- 正常启动
- 注册一个资源
- 通过 RPC 入口返回结果

## 1. 创建 `main.go`

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

## 2. 创建 `configs/application.toml`

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_source]
type = "sqlite"
```

这份配置之所以够用，是因为：

- SQLite 不需要外部服务
- storage 没配置时默认走内存实现
- 当前资源是 `Public`，所以不需要先接认证加载器
- 这个示例本身没有使用任何 Redis 相关能力

## 3. 启动应用

```bash
go run .
```

只要启动成功，VEF 就会打印应用 banner 并开始监听端口。

## 4. 调用 RPC 入口

发一个 `POST /api` 请求：

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

预期响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "message": "hello from vef"
  }
}
```

在框架默认语言下，成功消息通常会是 `成功`。如果你设置 `VEF_I18N_LANGUAGE=en`，同一条响应消息会变成 `Success`。

## 这个例子到底说明了什么

- `vef.Run(...)` 会启动整套运行时
- `vef.ProvideAPIResource(...)` 会把资源注册到 API engine
- `api.NewRPCResource(...)` 声明了一个 RPC 资源
- `api.OperationSpec` 定义了一个公开操作
- RPC 下 `hello` 会自动回退到 `Hello` 方法
- `result.Ok(...)` 负责生成框架统一响应结构

## 为什么 handler 这么短

你并没有手写下面这些事情：

- 挂载路由
- 解析请求体
- 验证 RPC 信封
- 组织统一响应格式
- 绑定中间件

这些都已经交给框架运行时处理了。

## 下一步

- [项目结构](./project-structure)：如何组织真实应用
- [模块与依赖注入](../modules/overview)：你的模块如何接到启动链里
- [路由](../guide/routing)：RPC 和 REST 的真实差异
