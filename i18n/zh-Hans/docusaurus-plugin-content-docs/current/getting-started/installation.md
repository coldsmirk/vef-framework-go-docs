---
sidebar_position: 1
---

# 安装

这一页只解决一件事：把一个 VEF 应用启动起来所需的最小环境准备好。

## 环境要求

当前框架版本要求：

- Go `1.26.4`

内置表达式引擎是 pure Go。expression 模块本身不要求 CGO；只有当你选择
的数据库驱动或其他 native integration 需要时，才需要启用 CGO。

## 运行前提

如果你直接使用默认的 `vef.Run(...)` 启动路径，框架一定会先装配数据库模块。

因此真正的最小启动前提是：

- 一个可达的数据源
- 一份可读的 `application.toml`

默认启动图里虽然也包含 Redis 模块，但只有当你的应用或某个启用能力真的消费 `*redis.Client` 或其他 Redis-backed 能力时，Redis 才会成为实际前提。

如果只是本地最小化跑通流程，SQLite 就已经足够。

## 安装框架

在你的 Go module 里安装：

```bash
go get github.com/coldsmirk/vef-framework-go
```

如果目录还是空的，可以这样开始：

```bash
go mod init example.com/my-app
go get github.com/coldsmirk/vef-framework-go
```

## 先选好数据库

VEF 在启动阶段就会装配数据库模块，所以配置文件里必须从一开始就有可用的数据源。当前框架支持：

- PostgreSQL
- MySQL
- SQLite
- SQL Server
- Oracle

如果你只是想先跑通流程，SQLite 是最省事的，因为它不需要额外服务。

## 需要 Redis 时再接入

如果你后续启用了 Redis 相关能力，默认客户端参数是：

- host：`127.0.0.1`
- port：`6379`
- network：`tcp`

本地最简单的方式之一是：

```bash
docker run --name vef-redis -p 6379:6379 -d redis:7-alpine
```

## 创建配置文件

默认情况下，VEF 会在这些位置查找 `application.toml`：

- `./configs`
- `$VEF_CONFIG_PATH`
- `.`
- `../configs`

最常见的布局是：

```text
my-app/
├── configs/
│   └── application.toml
└── main.go
```

## 最小配置

下面是能够启动的最小配置——SQLite、默认内存存储，以及存储模块要求的
outbox 事件传输：

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_sources.primary]
type = "sqlite"

[vef.event.transports.outbox]
enabled = true

[[vef.event.routing]]
pattern    = "vef.storage.*"
transports = ["outbox"]
```

`primary` 数据源是必填项，它为全框架注入的 `orm.DB` 提供来源；其他命名数据源也放在同一个 `vef.data_sources` map 下。

两个事件配置块不是可有可无的点缀：存储模块以事务方式发布领域事件，当
`vef.storage.*` 没有事务性路由时（默认的 `memory` 传输不具备事务性），
应用**启动即失败**。改设 `vef.event.default_transport = "outbox"` 是免路由
规则的替代方案。这一失败的完整演示见[快速开始](./quick-start)。

这里没有写存储配置也没关系，因为框架会默认回退到内存存储。
只有当应用真的使用 Redis 相关能力时，再补 `vef.redis` 即可。

## 启动时实际会发生什么

当你调用 `vef.Run(...)` 时，框架会依次初始化配置、数据源 registry 和 primary `orm.DB`、中间件、API、安全、事件（含 outbox / redis-stream / inbox 传输子模块）、表达式引擎、JS 引擎、CQRS、定时任务、Redis、分布式锁、mold、存储、sequence、schema、监控、MCP、服务端推送以及最终的 HTTP 服务。

所以对 VEF 来说，“安装完成”不只是把包 import 进来，而是要准备好能支撑这条启动链的最小配置。

## 常用环境变量

框架只读取以下三个环境变量，安装和本地调试阶段都会用到：

- `VEF_CONFIG_PATH`：额外配置目录
- `VEF_LOG_LEVEL`：日志级别
- `VEF_I18N_LANGUAGE`：框架语言，默认是简体中文

## 下一步

接着看 [快速开始](./quick-start)，用一个最小资源把应用真正跑起来。
