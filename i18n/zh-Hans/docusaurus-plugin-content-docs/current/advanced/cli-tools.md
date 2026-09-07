---
sidebar_position: 3
---

# CLI 工具

VEF 提供的 CLI 有两类性质不同的命令，这个划分是与你之间的约定，不是命名上的修饰：

- **`new`** 只写一次文件，然后就把它交给你。文件里没有「生成文件」标记，永远不会被重新生成，你可以随意编辑。
- **`generate-*` 与 `export-api`** 永远拥有自己的产物，每次运行都会覆盖它。

判断一条新命令属于哪一类，只需要问一个问题：*会有人想去改它吗？* 会 ⇒ `new`；不会 ⇒ 一个带着 `// Code generated ... DO NOT EDIT.` 的派生产物。

## 当前有哪些命令

根命令是 `vef-cli`。`vef-cli --version` 会打印 CLI banner 和
`Version: ...`；当构建时间元数据可用时还会打印 `Built: ...`，dirty VCS
构建会在版本号后追加 `-dirty`。

CLI 注册了这些子命令：

- `new`——`project`、`resource`、`service`
- `generate-build-info`
- `generate-model-schema`
- `export-api`

## 最小命令示例

```bash
vef-cli --version
vef-cli new project acme-server
vef-cli new resource --table hr_employee --module hr
vef-cli new service --name Payroll --module hr --deps bus:event.Bus
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
vef-cli generate-model-schema -i models -o schemas -p schemas
vef-cli export-api --app ./cmd/server -o api-manifest.json
```

应用代码应该通过这些命令来使用 CLI，而不是直接 import
`cmd/vef-cli/cmd/*` 下的实现包。

## `new`

### `new project <name>`

写出框架自己的目录约定：`cmd/server` 下的入口、`configs` 下的运行时配置、`internal` 下的业务模块，以及一份记录代码生成约定、供其他 `new` 命令读取的 `vef.yml`。它会把框架版本钉在 CLI 自身的版本上——工具是对着它自己构建时的 API 生成代码的，这个版本不可能过期——并在收尾前运行 `generate-build-info`，因为 `internal/vef/module.go` 引用了生成出来的 `BuildInfo`。

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--name`, `-N` | 位置参数 | 项目名称 |
| `--module`, `-m` | 项目名称 | Go module 路径 |
| `--path`, `-p` | `./<name>` | 创建项目的目录 |
| `--with-example` | `true` | 生成一个起步业务模块 |
| `--skip-tidy` | `false` | 跳过 `go mod tidy` |
| `--skip-git` | `false` | 跳过 `git init` |
| `--dry-run` | `false` | 只打印将写入什么，不动文件系统 |

Dockerfile、git hook 和 CI 流水线是刻意排除掉的：它们编码的是部署与团队决策，框架对此没有意见。

### `new resource`

从一张真实的数据库表生成 model、payload 和 API 资源——**外加两处注册**：model 注册表的 `var` 块，以及模块的 fx option 列表。这两处注册是价值的一半：漏掉任何一处，得到的都是一个能编译、但一律返回 404 的应用。模块不存在时会被创建，并接入 `cmd/server/main.go` 的 `vef.Run` 调用；最后这一步是尽力而为的，因为入口文件的形态是约定而非契约，所以自行改过接线方式的项目会收到一行 `TODO`，而不是一条失败的命令。

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--table`, `-t` | 必填 | 用来推导实体的数据库表 |
| `--module`, `-m` | 必填 | 生成到哪个业务模块，可以是多级 |
| `--entity` | 去掉模块前缀后的表名 | 实体名覆盖，snake_case |
| `--alias` | 表名各单词的首字母 | 表别名覆盖 |
| `--ops` | 项目配置的集合 | 要内嵌的 CRUD 操作 |
| `--search` | 所有标量列 | 形如 `column:operator` 的检索条件；`none` 生成空的检索载荷 |
| `--source` | primary | 要检查的数据源 |
| `--config` | `<project>/configs/application.toml` | `application.toml` 的路径 |
| `--force` | `false` | 覆盖已存在的生成文件 |
| `--dry-run`, `-n` | `false` | 只打印将写入什么，不动文件系统 |

表是通过框架自己的 schema 服务检查的，因此列注释会变成 `label` tag，可空会变成指针加 `omitempty`，声明的字符长度上限会变成 `max=N`，带有框架审计列的表会内嵌对应的 `orm` mixin 而不是重新声明这些字段。mixin 的匹配检查列的**类型**，而不只是名字——`orm.Model` 声明的是 `ID string`，而 ORM 会往任何零值字符串主键里写一个生成的 XID，所以仅凭名字去匹配一个 `BIGINT` 自增列，会在每次创建时把一个 20 字符的字符串塞进数字列里。

`find_tree` 会被拒绝而不是生成：`crud.NewFindTree` 需要一个任何生成器都无法提供的建树函数。

### `new service`

符合约定的骨架——一个只装注入依赖的结构体，方法签名为 `(ctx, db orm.DB, …)`——并注册到它所属的模块。框架依赖类型会自行解析所需的 import。

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--name`, `-n` | 必填 | PascalCase 的服务名，带不带 `Service` 后缀都可以 |
| `--module`, `-m` | 必填 | 生成到哪个业务模块，可以是多级 |
| `--deps` | 无 | 形如 `field:Type` 的注入依赖，例如 `bus:event.Bus` |
| `--force` | `false` | 服务文件已存在时覆盖它 |
| `--dry-run` | `false` | 只打印将写入什么，不动文件系统 |

### 约定来自 `vef.yml`

项目根目录的 `vef.yml` 记录 `module_root`，以及 `resource` 下的 `name` / `permission` 模板（占位符 `{module}`、`{domain}`、`{entity}`、`{action}`）、`ops`、`audit` 和 `audit_user_model`。每个键都有默认值，所以这个文件是可选的；`new project` 会写一份出来，让约定是看得见的，而不是隐含的。

对已有代码重复运行生成器，会报告哪些东西已经存在，而不是覆盖它们，除非 `--force` 另有指示。每条命令都会先把全部改动规划出来，因此 `--dry-run` 能展示完整结果，而规划中途失败也不会留下写了一半的文件。

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
`*.go` 文件，不会递归子目录。测试文件（`_test.go`）和被构建约束排除的文件会被
跳过。目录输入时，输出可以是一个已存在的目录，也可以是一个尚不存在的目录路径
（会按需创建）。如果输出路径已经作为文件存在，目录到单文件的生成会被拒绝。

生成器会读取目标文件中嵌入 `orm.BaseModel` 的 struct。表元数据来自嵌入的
`orm.BaseModel` 字段上的 `bun` tag：裸名称段（如 `bun:"users"`）或
`table:...` 选项设置表名（`table:` 优先），`alias:...` 设置默认 alias。
缺少这些 tag 部分时，表名默认为 model 名称的复数 snake_case，alias 默认为
model 名称的单数 snake_case。

字段处理遵循这些规则：

- 只有 exported 字段会生成 accessor
- `bun:"-"` 字段会被跳过
- `bun:"rel:*"` 和 `bun:"m2m:*"` 关系字段会被跳过
- 类似 `bun:"user_name"` 这样的第一个 `bun` tag 片段会设置列名
- 显式的 `column:...` 选项会同时覆盖裸名称段和字段名
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

## `export-api`

`export-api` 把应用的 API 接口面描述成数据。清单是由**应用自己**产出的：命令会带着 `VEF_EXPORT_API` 运行你的 main 包，所以它报告的是这个二进制真正注册了什么——资源名、action、认证策略、权限令牌、审计开关、生效的超时与限流，以及每个请求载荷的完整形状——而不是源码扫描猜出来的东西。

| Flag | 默认值 | 用途 |
| --- | --- | --- |
| `--app` | `./cmd/server` | 要描述的应用的 main 包 |
| `--output`, `-o` | `api-manifest.json` | 写入清单的文件，`-` 表示标准输出 |
| `--check` | `false` | 已有文件不一致时失败，而不是重写它 |

```bash
vef-cli export-api --app ./cmd/server -o api-manifest.json
vef-cli export-api --app ./cmd/server -o api-manifest.json --check
```

输出是排过序的，也不带时间戳，因此它就是拿来**提交进仓库**的：它的 diff 就是你 API 契约的 diff，而 `--check` 是让两者保持同步的 CI 闸门。有两条规则让这条命令可组合，而且都是踩过坑才有的：进度信息一律走 stderr 而不是 stdout，因此 `-o -` 可以直接管道给 `jq`；`--check` 与 `-o -` 同时使用会被拒绝而不是被忽略，因为那样写出来的 CI 任务什么都没比较，却在任何漂移下都返回 0。

### 清单的形状

```json
{
  "framework": "v0.51.0",
  "resources": [
    {
      "name": "security/auth",
      "kind": "rpc",
      "version": "v1",
      "operations": [
        {
          "action": "login",
          "auth": "none",
          "timeoutMs": 30000,
          "rateLimit": { "max": 6, "periodMs": 300000 },
          "params": "github.com/…/security.LoginParams"
        }
      ]
    }
  ],
  "types": {
    "github.com/…/security.LoginParams": {
      "fields": [
        { "name": "type", "type": "string", "validate": "required" }
      ]
    }
  }
}
```

字段带有 `name`（取自 json tag）、`type`，以及可选的 `optional` / `label` / `validate`；`permission` 与 `audit` 只在设置过时才出现在 operation 上。具名结构体会渲染成它在 `types` 中占据的键，因此一张载荷关系图只描述一次，之后按名字引用。

返回类型是刻意缺席的：crud handler 只返回 `error`，响应是通过 context 写出去的，所以模型类型根本不出现在签名里。补上这个缺口只需要在 `api.OperationSpec` 上追加一个字段，那是计划中的下一步而非疏漏——同一个缺口也吞掉了 `delete` / `delete_many`，它们的 handler 接收的是无类型的参数袋。

### 从你自己的入口导出

当 `vef.EnvExportAPI`（`VEF_EXPORT_API`）环境变量指定了目标时，`vef.Run` 会转交给 `vef.ExportAPI`，因此宿主不需要第二个入口。想在进程内拿到清单时直接调用它：

```go
if err := vef.ExportAPI(os.Stdout, options...); err != nil {
    return err
}
```

用环境变量而不是命令行 flag，是因为 `Run` 并不拥有进程的 flag set——宿主的 `main` 可能已经定义了自己的。

**什么都不会启动，也什么都不会连接。** 资源在构造期就完成注册，数据源是懒打开的，而调用 crud handler 工厂的挂载动作发生在之后——这正是让导出能安全跑在容器构建或 CI 任务里的原因。它是给「跑完就退出」的进程用的：有少数构造函数会自己起 goroutine 而不是注册生命周期钩子（会话存储和登录守卫背后的内存缓存各跑一个 GC ticker），而由于依赖图从未启动，也就从未停止，所以在长期运行的进程里反复调用 `ExportAPI` 会让它们不断累积。

接口面是通过 `api.EngineInspector` 读取的，这是一个**可选**接口（`Operations() []*api.Operation`，按标识符排序），框架自己的 engine 实现了它。做成 inspector 而不是 `api.Engine` 上的方法——与 `event.StreamInspector` 及其同类一致——是为了让新增内省能力不破坏任何自行实现 `Engine` 的宿主。没有实现它的 engine 会得到 `vef.ErrEngineNotInspectable`，这意味着有东西替换掉了框架自己的实现。

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

## 每条命令各归其位

| 命令 | 拥有产物吗？ | 用来做什么 |
| --- | --- | --- |
| `new project` / `new resource` / `new service` | 不拥有——写一次，之后归你 | 起一个项目、加一个 CRUD 资源或一个服务 |
| `generate-build-info` | 拥有——会被重新生成 | 通过 `sys/monitor` 暴露的构建元数据 |
| `generate-model-schema` | 拥有——会被重新生成 | 从 model 派生的 schema 辅助代码 |
| `export-api` | 拥有——会被重新生成 | 提交进仓库的 API 清单及其 CI 漂移闸门 |

把这条线模糊掉，生成器就会变成一条没人敢运行的命令。

## 下一步

如果你希望生成的构建信息通过 `sys/monitor` 展示出来，继续阅读 [监控](../infrastructure/monitor)。
