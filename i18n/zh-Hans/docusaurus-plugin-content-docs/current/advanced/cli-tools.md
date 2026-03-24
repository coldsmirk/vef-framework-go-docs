# CLI 工具

VEF 在 `cmd/vef-cli` 下提供了一个 CLI 入口，但当前可用能力是有限的。这一页只记录**源码里真实已经实现**的部分。

## 当前有哪些命令

当前 CLI 根命令注册了三个子命令：

- `create`
- `generate-build-info`
- `generate-model-schema`

## 重要限制：`create` 还没有实现

`create` 命令当前会直接返回一个明确错误：

```text
vef-cli create is not implemented yet, please generate the project manually
```

所以不要把它当成一个可用的项目脚手架来文档化或依赖。

## `generate-build-info`

这个命令用于生成一个包含构建信息的 Go 文件，信息包括：

- 应用版本
- 构建时间
- Git commit

常见用法：

```bash
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
```

如果你想让 `sys/monitor.get_build_info` 返回应用自己的构建元数据，这个命令是当前最实用的一项 CLI 能力。

## `generate-model-schema`

这个命令会读取符合 VEF ORM 模型约定的 Go model，并生成 schema 辅助代码，避免在 ORM 查询中到处硬编码列名。

常见用法：

```bash
vef-cli generate-model-schema -i models -o schemas -p schemas
```

## 常见的 `go:generate` 用法

在真实 VEF 项目里，这两个命令通常直接写在 `module.go` 顶部：

```go
//go:generate vef-cli generate-model-schema -i ./models -o ./schemas -p schemas
package sys
```

以及：

```go
//go:generate vef-cli generate-build-info -o ./build_info.go -p vef
package vef
```

这样 schema helper 和 build info 会和真正使用它们的模块放在一起。

## 当前真实工作流

在现阶段，比较现实的做法是：

1. 手动创建项目结构
2. 手动编写模块与资源
3. 需要时使用 `generate-build-info`
4. 需要时使用 `generate-model-schema`

## 下一步

如果你想让生成的构建信息最终出现在运行时接口里，继续阅读 [监控](../features/monitor)。
