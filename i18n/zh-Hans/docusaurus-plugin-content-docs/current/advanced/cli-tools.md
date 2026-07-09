---
sidebar_position: 3
---

# CLI 工具

VEF 提供了一个 CLI，但目前的范围比完整的项目脚手架工具要窄。

## 当前有哪些命令

根命令是 `vef-cli`。`vef-cli --version` 会打印 CLI banner 和
`Version: ...`；当构建时间元数据可用时还会打印 `Built: ...`，dirty VCS
构建会在版本号后追加 `-dirty`。

CLI 当前注册了这些子命令：

- `create`
- `generate-build-info`
- `generate-model-schema`

## 最小命令示例

```bash
vef-cli --version
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
vef-cli generate-model-schema -i models -o schemas -p schemas
```

应用代码应该通过这些命令来使用 CLI，而不是直接 import
`cmd/vef-cli/cmd/*` 下的实现包。

## 重要现状说明

`vef-cli create` 作为命令存在，但目前**尚未实现**。

不要把它当作一个可用的项目生成器来使用。

该命令会返回这个错误：

```text
vef-cli create is not implemented yet, please generate the project manually
```

因为计划中的命令形态已经存在，命令仍然定义了这些 flags：

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--name`, `-n` | 必填 | 项目名称 |
| `--path`, `-p` | `.` | 项目将要创建到的目录路径 |
| `--module`, `-m` | 空 | Go module 路径 |

## `generate-build-info`

这个命令会生成一个包含构建元数据的 Go 源文件，例如：

- 应用版本
- 构建时间
- git commit

它适合放在 `go:generate` 或构建流水线里使用。

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

## `generate-model-schema`

这个命令会检查 model 文件，并为 ORM 使用生成类型安全的 schema 辅助代码。

它支持：

- 文件到文件的生成
- 目录到目录的生成

目标是减少查询代码里硬编码的列名字符串。

Flags:

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--input`, `-i` | 必填 | 输入 model 文件或目录 |
| `--output`, `-o` | 必填 | 输出 schema 文件或目录 |
| `--package`, `-p` | `schemas` | 生成 schema 文件的 package 名称 |

目录输入会按输入文件逐一生成 schema 文件。目录模式只处理输入目录直属的
`*.go` 文件，不会递归子目录。目录输入可以指向一个已存在的输出目录，也可以
指向一个尚不存在的目录路径。如果输出路径已经作为文件存在，目录到单文件的
生成会被拒绝。

生成器会读取目标文件中嵌入 `orm.BaseModel` 的 struct。表元数据来自嵌入的
`orm.BaseModel` 字段上的 `bun` tag：`table:...` 设置表名，`alias:...` 设置
默认 alias。缺少这些 tag 部分时，表名默认为 model 名称的复数 snake_case，
alias 默认为 model 名称的单数 snake_case。

字段处理遵循这些规则：

- 只有 exported 字段会生成 accessor
- `bun:"-"` 字段会被跳过
- `bun:"rel:*"` 和 `bun:"m2m:*"` 关系字段会被跳过
- 类似 `bun:"user_name"` 这样的第一个 `bun` tag 片段会设置列名
- 没有列名 tag 的字段使用字段名的 snake_case 形式
- embedded struct 会被展开
- `bun:"embed:prefix_"` 会用给定的前缀展开嵌套字段
- `label:"..."` 会变成生成代码里的方法注释
- `bun:",scanonly"` 字段仍然会有 accessor，但会被排除在 `Columns()` 之外

生成的公开 API 会暴露一个以 model 命名的 exported schema 变量，例如 `User`，
其背后是一个 unexported schema 类型，例如 `userSchema`。每个 schema 都有
字段 accessor，以及 `Table()`、`Alias()`、`As(alias)`、`Columns()`。

字段 accessor 默认通过 `dbx.ColumnWithAlias` 返回带 alias 限定的列名。传入
`raw=true` 会返回原始列名：

```go
schemas.User.Name()     // 例如 "u.name"
schemas.User.Name(true) // "name"
```

如果 model 字段会和 `Table`、`Alias`、`As` 或 `Columns` 冲突，生成的
accessor 会加上 `Col` 前缀，例如 `ColTable`。生成的 struct 字段标识符如果
会撞上 Go 关键字，会加上 `__` 前缀。

## 常见的 `go:generate` 用法

在真实的 VEF 应用里，这些命令通常直接写在 `module.go` 上方：

```go
//go:generate vef-cli generate-model-schema -i ./models -o ./schemas -p schemas
package sys
```

以及面向框架的构建元数据：

```go
//go:generate vef-cli generate-build-info -o ./build_info.go -p vef
package vef
```

这样可以让 schema 辅助代码和构建元数据在物理位置上贴近使用它们的模块。

## 现阶段的合理预期

目前，这个 CLI 最好被当作：

- 生成构建元数据的辅助工具
- 生成 model schema 的辅助工具

它**还不**适合作为「一条命令搭好整个项目」这类 onboarding 文档的基础。

## 下一步

如果你希望生成的构建信息通过 `sys/monitor` 展示出来，继续阅读 [监控](../infrastructure/monitor)。
