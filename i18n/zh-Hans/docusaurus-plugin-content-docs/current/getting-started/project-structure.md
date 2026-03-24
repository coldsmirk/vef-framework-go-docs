---
sidebar_position: 3
---

# 项目结构

VEF 不强制单一目录结构，但它天然鼓励一种“贴近模块装配方式”的组织形式。

## 一个实用的起步结构

```text
my-app/
├── configs/
│   └── application.toml
├── internal/
│   ├── user/
│   │   ├── module.go
│   │   ├── model/
│   │   ├── payload/
│   │   └── resource/
│   ├── auth/
│   │   ├── module.go
│   │   ├── user_loader.go
│   │   └── permission_loader.go
│   └── app/
│       └── module.go
└── main.go
```

这种结构的好处在于，它和你最终注册到 FX 的方式是对应的。

## 更贴近生产项目的结构

在规模更大的 VEF 应用里，模块边界通常会更明确。很常见的一种组织方式是：

```text
my-app/
├── cmd/
│   └── server/
│       └── main.go
├── configs/
│   └── application.toml
├── internal/
│   ├── vef/    # 面向框架集成的模块
│   ├── auth/   # UserLoader、UserInfoLoader 等认证接入
│   ├── web/    # SPA 托管
│   ├── mcp/    # MCP tool/resource provider
│   ├── sys/    # 系统/管理域资源
│   ├── md/     # 主数据资源
│   └── pmr/    # 业务域资源
└── go.mod
```

这种写法的好处是把这几类职责拆开了：

- 框架集成代码
- 身份与认证加载
- 前端托管
- 业务域本身

## 按业务域组织，而不是按技术层切碎

更推荐的原则是：

- 同一业务域下的 model、payload、service、resource 包放在一起
- 每个业务域或子域导出一个 `Module`
- 在 `main.go` 里组合这些模块

例如：

```go
package main

import (
	vef "github.com/coldsmirk/vef-framework-go"

	"example.com/my-app/internal/auth"
	"example.com/my-app/internal/user"
)

func main() {
	vef.Run(
		auth.Module,
		user.Module,
	)
}
```

## 一个模块里通常放什么

一个典型模块通常会注册这些内容中的若干项：

- API 资源
- 领域服务
- 中间件
- 权限加载器或认证相关 provider
- CQRS handler 或 behavior

示例：

```go
package user

import (
	vef "github.com/coldsmirk/vef-framework-go"

	"example.com/my-app/internal/user/resource"
)

var Module = vef.Module(
	"app:user",
	vef.ProvideAPIResource(resource.NewUserResource),
)
```

在更大的应用里，还常见一个专门的集成模块，用来做这些事情：

- `vef.Supply(...)` build info
- `vef.Provide(...)` 共享的框架侧 service 或 loader
- `vef.Invoke(...)` 事件订阅器注册

这样业务域模块就不用夹带太多框架粘合代码。

## 推荐的子包划分

下面这些名字在 VEF 应用里比较常见，也便于后续扩展：

- `model`：Bun 模型和持久化类型
- `payload`：请求参数、搜索参数、传输层 DTO
- `resource`：VEF API 资源
- `service`：应用服务或领域服务
- `query` / `command`：如果你采用 CQRS，可以在这里拆分

默认优先推荐单数包名。Go 里的 package 名更像命名空间，而不是集合名，所以 `model`、`payload`、`resource`、`service` 往往比 `models`、`payloads`、`resources`、`services` 更符合 Go 社区习惯。

不是一上来就必须全有。规模小的时候完全可以先少量文件起步，再按增长拆分。

## 认证相关代码放哪里

认证和授权通常需要应用自己提供这些实现：

- `security.UserLoader`
- `security.UserInfoLoader`
- `security.RolePermissionsLoader`
- `security.ExternalAppLoader`

最常见做法是把这些实现集中在你自己的 `auth` 或 `security` 业务模块里。

## 前端资源放哪里

如果应用还要通过 VEF 托管单页前端，建议把前端构建产物或嵌入文件系统适配器放进独立的 `web` / `frontend` 模块，不要和 API 资源包混在一起。

同样地，如果应用要扩展 MCP，也更适合放在单独的 `internal/mcp` 模块里，而不是散落到普通业务资源包中。

## 应该避免什么结构

不建议把整个系统按“框架层”切成几个超大桶：

- 一个包装全站所有 model
- 一个包装全站所有 resource
- 一个包装全站所有 service

这种结构早期看似简单，后面会因为业务耦合越来越难维护。VEF 更适合按 feature module 扩展，而不是按技术层堆仓库。

## 下一步

接下来建议看 [模块与依赖注入](../modules/overview)，理解这些模块如何真正进入运行时。
