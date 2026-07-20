---
sidebar_position: 10
---

# 多数据源

大多数应用只需要和一个数据库打交道：**主数据源**（primary），它在框架各处都以 `orm.DB` 的形式被注入。这篇文档是给"其余情况"看的——报表仓库、按租户拆分的数据库、你只读不写的遗留系统——它们都通过 `datasource.Registry` 访问。

## 主数据源与附加数据源

主数据源在 TOML 中的 `vef.data_sources.primary` 下声明。它是必填项，也是框架全局以 `orm.DB` 形式暴露的数据源，并且不能通过动态 registry API 修改——`Register`、`Update`、`Unregister` 都会拒绝 `datasource.PrimaryName`（即 `"primary"`），返回 `datasource.ErrPrimaryReserved`。

除主数据源之外的都是"附加数据源"：既可以在 TOML 里以另一个名字静态声明，也可以在运行时动态注册。要访问附加数据源，就在需要的地方注入 `datasource.Registry`。

框架内部模块——CRUD、审批（approval）、storage、event inbox/outbox、schema 反射——都只操作主数据源。是否要访问附加数据源、以及如何使用它，都是应用层自己的事。

## 静态数据源：TOML

添加数据源最简单的方式，就是在 `primary` 旁边再写一个 `vef.data_sources.<name>` 表：

```toml
[vef.data_sources.primary]
type = "postgres"
host = "127.0.0.1"
port = 5432
user = "postgres"
password = "postgres"
database = "my_app"
schema = "public"

[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

每一项都使用同样的 `config.DataSourceConfig` 结构：

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `type` | `postgres \| mysql \| sqlite \| sqlserver \| oracle` | 数据库类型 |
| `host` | `string` | 数据库网络地址 |
| `port` | `uint16` | 数据库端口 |
| `user` | `string` | 数据库用户名 |
| `password` | `string` | 数据库密码 |
| `database` | `string` | 数据库名 |
| `schema` | `string` | 支持 schema 的驱动使用的 schema 名 |
| `path` | `string` | SQLite 文件路径 |
| `enable_sql_guard` | `bool` | 是否为原生 SQL 通道启用 SQL guard |
| `ssl_mode` | `disable \| require \| verify-ca \| verify-full` | 网络型驱动的 TLS 策略，默认 `disable` |
| `ssl_root_cert` | `string` | `verify-ca` / `verify-full` 模式下可选的 PEM 路径 |

`vef.data_sources` 下除 `primary` 之外的每一项，都会在应用开始处理请求之前，以其在 map 中的键名注册进 `datasource.Registry`。完整字段列表见 [Configuration Reference](../reference/configuration-reference)。

SQL Server 与 Oracle 的方言说明：

- `sqlserver`（驱动 `microsoft/go-mssqldb`）：默认端口 `1433`；`database`
  为空时落在登录账号的默认目录上。任何 TLS `ssl_mode` 都会强制开启 TDS
  加密；`verify-ca`/`verify-full` 语义（自定义 PEM 根证书、主机名固定）与
  Postgres/MySQL 一致。
- `oracle`（驱动 `sijms/go-ora`）：默认端口 `1521`；`database` 填**服务名**
  且必填；`schema` 不映射——Oracle 按连接用户自身的 schema 解析未限定名。
  `verify-ca` 以 `verify-full` 语义（更强的校验）提供，驱动不支持 PEM
  `ssl_root_cert`，配置了会快速失败。

## 注入 Registry

`datasource.Registry` 在 FX 容器中随处可用，直接注入即可：

```go
package report

import (
	"context"

	"github.com/coldsmirk/vef-framework-go/datasource"
)

type Service struct {
	sources datasource.Registry
}

func NewService(sources datasource.Registry) *Service {
	return &Service{sources: sources}
}

func (s *Service) RunReport(ctx context.Context) error {
	analytics, err := s.sources.Get("analytics")
	if err != nil {
		return err
	}

	var count int
	return analytics.NewSelect().
		ColumnExpr("count(*)").
		Table("events").
		Scan(ctx, &count)
}
```

`Get` 返回一个 `orm.DB`，因此它和主数据源暴露的是同一套查询构建 API（`NewSelect`、`NewInsert`、`NewRaw` 等）——参见 [Query Builder](./query-builder)。未注册或已注销的名字会返回 `datasource.ErrNotFound`。`sources.Primary()` 返回的 `orm.DB` 和直接注入 `orm.DB` 得到的是同一个，且永远不会出错。

`datasource.Registry` 同时也是内置的 API handler 参数——你可以像 `orm.DB`、`fiber.Ctx` 一样，直接把它作为 resource handler 的参数请求——参见 [Custom Handlers](../building-apis/custom-handlers)。

## 运行时数据源：`datasource.Provider`

对于部署时还不知道的数据源——比如主数据库里的一张租户表，每一行描述一个附加数据库——实现 `datasource.Provider`：

```go
type Provider interface {
	Name() string
	Load(ctx context.Context) ([]Spec, error)
}

type Spec struct {
	Name   string
	Config config.DataSourceConfig
}
```

框架在启动期间调用一次 `Load`，时机在主数据源和 TOML 静态数据源**都已经**注册完成**之后**，并把每一个返回的 `Spec` `Register` 进去。名字与 TOML 或其他 provider 冲突会导致启动失败。`Provider.Name` 只用于诊断信息中标识该 provider，并不是数据源名字本身。

用 `vef.ProvideDataSourceProvider` 注册这个 provider:

```go
func NewTenantSourceProvider(primary orm.DB) datasource.Provider {
	return &tenantSourceProvider{primary: primary}
}

// in your fx.Module:
vef.ProvideDataSourceProvider(NewTenantSourceProvider)
```

多个 provider 之间的执行顺序是不确定的，所以不同 provider 返回的 spec 之间不能在名字上冲突。

## 运行时数据源：直接调用 `Register`

在启动期的 `Provider` 钩子之外——比如某个管理端点允许运维人员按需添加一个数据源——可以直接调用 `Registry.Register`：

```go
db, err := sources.Register(ctx, "tenant-42", config.DataSourceConfig{
	Kind:     config.Postgres,
	Host:     "tenant-42.internal",
	Port:     5432,
	User:     "app",
	Password: pw,
	Database: "tenant_42",
})
```

`Register` 会先打开并 ping 通新连接，再把它插入 registry;一旦发生冲突（重名返回 `datasource.ErrExists`，名字为 `"primary"` 返回 `datasource.ErrPrimaryReserved`，名字为空或包含空白/控制字符返回 `datasource.ErrNameInvalid`），刚打开的连接会被关闭，不会插入任何内容。`Register` 从不关闭已有连接，所以它不接受任何 option。

## Registry 的完整方法

| 方法 | 行为 |
| --- | --- |
| `Primary()` | 返回主数据源的 `orm.DB`。等价于 `Get(datasource.PrimaryName)`，但永远不会出错。 |
| `Get(name)` | 返回已注册的 `orm.DB`。名字未注册或已被注销时返回 `datasource.ErrNotFound`。 |
| `Has(name)` | 判断 `name` 当前是否已注册且未关闭。 |
| `Names()` | 返回全部已注册的名字（含 `primary`）,按稳定的字典序排列。 |
| `Kind(name)` | 返回 `name` 对应的 `config.DBKind`。未找到时的语义与 `Get` 相同。 |
| `Register(ctx, name, cfg)` | 打开并 ping 通一个新的非主数据源，然后插入。 |
| `Update(ctx, name, cfg, opts...)` | 原子地替换一个已有数据源的连接。旧连接池异步关闭。 |
| `Unregister(ctx, name, opts...)` | 移除一个非主数据源。旧连接池异步关闭。 |
| `Reconcile(ctx, specs, opts...)` | 一次调用，把 registry 驱动到指定的非主数据源目标集合。 |
| `TestConnection(ctx, cfg)` | 打开一个临时连接、验证、再关闭。不会修改 registry。 |
| `HealthCheck(ctx)` | 并行 ping 所有已注册数据源；返回 `name -> error` 的映射。 |

所有只读方法（`Get`、`Has`、`Names`、`Kind`、`Primary`）都可以并发安全地调用。`Register`、`Update`、`Unregister` 会原子地修改 registry。

## 对齐目标集合：Reconcile

当数据源列表来自某张外部表时（和上面 `Provider` 例子里的租户表场景一样）,表内容和 registry 之间随着应用重启发生漂移是很常见的：行被增加、更新或删除。`Reconcile` 用一次调用就能补齐这个差距，不需要你自己手写 diff 逻辑：

```go
report, err := sources.Reconcile(ctx, specs)
```

给定一个目标 `[]datasource.Spec`，`Reconcile` 会计算出三类差异并驱动 registry 达成：

- spec 存在但 registry 中没有对应项 → `Register`
- spec 与 registry 中当前配置不同 → `Update`
- registry 中存在但 spec 里没有对应项 → `Unregister`

引用了主数据源名字的 spec 会被忽略。单个名字失败会被记录到 `ReconcileReport.Errors`（以名字为 key）,不会中断本批次其余的处理——十个里有一个配置有问题，不会连累另外九个。多次 `Reconcile` 调用是串行执行的：同一个定时刷新任务（通常是按计划调用 `Reconcile` 的 cron job）的两次执行不会互相交叉、产生竞态。但直接调用 `Register` / `Update` / `Unregister`，**不会**与正在执行中的 `Reconcile` 同步。

用 `datasource.WithReconcileDryRun()` 只计算报告、不实际打开或关闭任何连接——适合预览一个刷新任务将要做什么：

```go
preview, _ := sources.Reconcile(ctx, specs, datasource.WithReconcileDryRun())
```

## 更新与移除数据源

`Update` 和 `Unregister` 是唯二会关闭已有连接的操作，二者都接受 `datasource.RegisterOption`。默认情况下，被替换或移除的连接池会在后台 goroutine 中立即关闭；`WithCloseGrace(d)` 可以延迟这次关闭，给正在进行中的查询留出时间完成：

```go
_, err := sources.Update(ctx, "analytics", newCfg, datasource.WithCloseGrace(10*time.Second))
```

```go
err := sources.Unregister(ctx, "analytics", datasource.WithCloseGrace(10*time.Second))
```

无论哪种情况，调用一旦返回，`Get("analytics")` 都会立即反映新状态（`Update` 返回新配置对应的 `orm.DB`，`Unregister` 返回 `datasource.ErrNotFound`）——宽限期只影响**旧的** `*sql.DB` 何时真正关闭，所以在替换之前已经持有 `orm.DB` 引用的调用方，可以借着这段时间完成正在进行的查询。

## 测试连接与健康检查

`TestConnection` 是一个纯粹的连通性探测：它用候选配置打开一个临时连接、通过查询服务器版本来确认可用性、然后立即关闭——全程不触碰 registry。它天然适合作为管理后台"测试连接"按钮的后端实现，在调用 `Register` 或 `Update` 之前先跑一遍：

```go
info, err := sources.TestConnection(ctx, candidateCfg)
if err != nil {
	// unreachable or unusable
}
// info.Version, e.g. "PostgreSQL 16.2 on x86_64-pc-linux-gnu"
```

`HealthCheck` 则是并行 ping 当前所有已注册的数据源（含主数据源）,返回一个 `name -> error` 的映射，适合用在需要一次性掌握所有数据源状态的存活/就绪探针里。

## 需要记住的约束

- **框架内部模块只用主数据源。** CRUD、approval、storage、event inbox/outbox、schema 反射读写的都是主数据源。附加数据源永远不会被框架内部代码触碰——只有你自己的代码会用到它们。
- **不支持跨数据源事务。** `orm.DB.RunInTx` 只针对单个数据源开启事务。如果你用 `event.WithTx(tx)` 发布事件，`tx` 必须来自主数据源上开启的事务——参见 [Transactions](./transactions)。
- **`Reconcile` 只管理非主数据源。** 引用了 `datasource.PrimaryName` 的 spec 会被静默忽略，所以一个 Provider 或定时对账任务即使在输入里带上了主数据源的名字，也不会导致启动失败。

## 下一步

拿到 `Get` 或 `Primary` 返回的 `orm.DB` 之后，关于模型和查询构建部分，可以阅读 [Query Builder](./query-builder) 和 [Transactions](./transactions)。想了解包括 `vef.ProvideDataSourceProvider` 在内的其他框架扩展点，参见 [Extension Points](../reference/extension-points)。
