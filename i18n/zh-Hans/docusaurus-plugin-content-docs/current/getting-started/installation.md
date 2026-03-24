---
sidebar_position: 1
---

# 安装

这一页只解决一件事：把一个 VEF 应用启动起来所需的最小环境准备好。

## 环境要求

当前框架版本要求：

- Go `1.26.0`

对 VEF 应用本身来说，Go 是唯一必须的运行时依赖。

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

下面这份配置已经足够启动一个使用 SQLite 和默认内存存储的应用：

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_source]
type = "sqlite"
```

这里没有写存储配置也没关系，因为框架会默认回退到内存存储。
只有当应用真的使用 Redis 相关能力时，再补 `vef.redis` 即可。

## 启动时实际会发生什么

当你调用 `vef.Run(...)` 时，框架会依次初始化配置、数据库、ORM、中间件、API、安全、事件、CQRS、定时任务、Redis、mold、存储、sequence、schema、监控、MCP 以及最终的 HTTP 服务。

所以对 VEF 来说，“安装完成”不只是把包 import 进来，而是要准备好能支撑这条启动链的最小配置。

## 常用环境变量

安装和本地调试阶段最常用的是这些：

- `VEF_CONFIG_PATH`：额外配置目录
- `VEF_LOG_LEVEL`：日志级别
- `VEF_NODE_ID`：节点 ID
- `VEF_I18N_LANGUAGE`：框架语言，默认是简体中文

## 下一步

接着看 [快速开始](./quick-start)，用一个最小资源把应用真正跑起来。
