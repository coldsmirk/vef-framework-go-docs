---
sidebar_position: 1
---

# CLI 工具

VEF 在 `cmd/vef-cli` 下提供了一个 CLI 入口，但当前可用能力是有限的。这一页只记录**源码里真实已经实现**的部分。

## 当前有哪些命令

根命令是 `vef-cli`。`vef-cli --version` 会打印 CLI banner 和
`Version: ...`；如果构建时间元数据存在，还会打印 `Built: ...`，dirty VCS
构建会在版本号后追加 `-dirty`。

当前 CLI 根命令注册了三个子命令：

- `create`
- `generate-build-info`
- `generate-model-schema`

`cmd/vef-cli/**` 会出现在 [公开 API 索引](../reference/public-api-index)
里，这是为了导出审计完整性；应用代码应通过这些 CLI 命令来消费能力，不应
import 命令实现包。

已审查的 CLI implementation packages：

| Package | Public entries | Fingerprint | 用户可见 contract |
| --- | ---: | --- | --- |
| `github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd` | 10 | `6a01b8fdcb43f6842164be353432a6dbc7849601835c454228aab6cb5ef046ef` | 根命令、已注册子命令、`--version` 输出 |
| `github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/buildinfo` | 2 | `a9f40a22aaf4f4e6313cea5a7fcd439a5dcde2d0b13f977e954753c1317ab33e` | `generate-build-info` |
| `github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/create` | 2 | `26171a8454bd55208efc47d3ba16ce5744a971956d17bfac4972c1468619cd3b` | `create` placeholder 和 not-implemented error |
| `github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/modelschema` | 21 | `19164973da27a846f72a4df3b55d320998b55e57ee2b3dc40dd7abc4868e8735` | `generate-model-schema` |

最小示例：

```bash
vef-cli --version
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
vef-cli generate-model-schema -i models -o schemas -p schemas
```

## 重要限制：`create` 还没有实现

`create` 命令当前会直接返回一个明确错误：

```text
vef-cli create is not implemented yet, please generate the project manually
```

所以不要把它当成一个可用的项目脚手架来文档化或依赖。

这个命令仍然定义了以下 flags，因为计划中的命令形状已经存在：

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--name`, `-n` | 必填 | 项目名称 |
| `--path`, `-p` | `.` | 项目将要创建到的目录 |
| `--module`, `-m` | 空 | Go module 路径 |

## `generate-build-info`

这个命令用于生成一个包含构建信息的 Go 文件，信息包括：

- 应用版本
- 构建时间
- Git commit

常见用法：

```bash
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
```

Flags:

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--output`, `-o` | `build_info.go` | 输出 Go 文件 |
| `--package`, `-p` | `main` | 生成文件的 package 名称 |

生成文件会导出 `BuildInfo = &monitor.BuildInfo{...}`，并填充：

- `AppVersion` 来自 `git describe --tags --always --dirty`，失败时回退到 `dev`
- `BuildTime` 来自 `timex.Now().String()`
- `GitCommit` 来自 `git rev-parse HEAD`，失败时回退到 `none`

生成器会按需创建输出目录。生成文件的公开形状是：

```go
var BuildInfo = &monitor.BuildInfo{
	AppVersion: "...",
	BuildTime:  "...",
	GitCommit:  "...",
}
```

如果你想让 `sys/monitor.get_build_info` 返回应用自己的构建元数据，这个命令是当前最实用的一项 CLI 能力。

## `generate-model-schema`

这个命令会读取符合 VEF ORM 模型约定的 Go model，并生成 schema 辅助代码，避免在 ORM 查询中到处硬编码列名。

常见用法：

```bash
vef-cli generate-model-schema -i models -o schemas -p schemas
```

Flags:

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--input`, `-i` | 必填 | 输入 model 文件或目录 |
| `--output`, `-o` | 必填 | 输出 schema 文件或目录 |
| `--package`, `-p` | `schemas` | 生成 schema 文件的 package 名称 |

目录输入会按输入文件一一生成 schema 文件。目录模式只处理输入目录直属的
`*.go` 文件，不递归子目录。目录输入可以指向一个已存在的输出目录，也可以
指向一个尚不存在的目录路径；如果输出路径已经存在并且是文件，目录到单文件
生成会被拒绝。

生成器会读取目标文件中嵌入 `orm.BaseModel` 的 struct。表元数据来自嵌入
`orm.BaseModel` 字段上的 `bun` tag：`table:...` 设置表名，`alias:...` 设置
默认 alias。缺少这些 tag 时，表名默认是 model 名称的 pluralized snake_case，
alias 默认是 model 名称的 singular snake_case。

字段处理规则如下：

- 只有 exported fields 会生成 accessors
- `bun:"-"` 字段会被跳过
- `bun:"rel:*"` 和 `bun:"m2m:*"` 关系字段会被跳过
- `bun:"user_name"` 这类第一个 `bun` tag 片段会设置列名
- 没有列名 tag 的字段使用字段名的 snake_case
- embedded structs 会被展开
- `bun:"embed:prefix_"` 会用 prefix 展开嵌套字段
- `label:"..."` 会变成生成代码里的方法注释
- `bun:",scanonly"` 字段仍然有 accessor，但会从 `Columns()` 中排除

生成代码的公开 API 是一个以 model 名命名的 exported schema 变量，例如
`User`，背后对应 unexported schema 类型，例如 `userSchema`。每个 schema
都有字段 accessor，以及 `Table()`、`Alias()`、`As(alias)`、`Columns()`。

字段 accessor 默认通过 `dbx.ColumnWithAlias` 返回带 alias 的列名。传入
`raw=true` 会返回 raw column name：

```go
schemas.User.Name()     // 例如 "u.name"
schemas.User.Name(true) // "name"
```

如果 model 字段名会和 `Table`、`Alias`、`As` 或 `Columns` 冲突，生成的
accessor 会加 `Col` 前缀，例如 `ColTable`。生成的 struct-field identifier
如果会和 Go keyword 冲突，会加 `__` 前缀。

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
