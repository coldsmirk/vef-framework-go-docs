---
sidebar_position: 5
---

# RPC 资源

启用 `vef.IntegrationModule` 后，框架注册以下管理资源。它们都是挂载在
`/api` 下的 RPC 资源，使用[API](../building-apis/api.md)文档中的标准信封
（`resource`、`action`、`version`、`params`、`meta`）。所有操作都不是公开
的：调用方必须已认证，且每个操作声明其表格所列的权限。

本页使用的约定：

- CRUD 读操作从 `params` 解码其查询结构体；`find_page` 额外从 `meta.page`
  与 `meta.size` 读取分页（`page.Pageable`），可排序资源读取 `meta.sort`。
- `find_page` 的响应是 `page.Page[T]`：`page`、`size`、`total`、`items`。
- 变更操作从 `params` 解码。标记必填的字段由校验强制；其余为可选。
- 所有定义模型的响应都携带标准审计字段（`id`、`createdAt`、`createdBy`、
  `updatedAt`、`updatedBy`），下方字段表中不再重复。

## `integration/contract`

契约定义。Schema 在保存时编译，坏掉的契约永远不会进入调用。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `integration.contract.query` | `ContractSearch` + 分页 meta | `page.Page[Contract]` |
| `find_all` | `integration.contract.query` | `ContractSearch` | `Contract[]` |
| `create` | `integration.contract.create` | `ContractParams` | 创建后的 `Contract` |
| `update` | `integration.contract.update` | `ContractParams` | 更新后的 `Contract` |
| `delete` | `integration.contract.delete` | 主键参数（`params.id`） | 成功 |

`ContractSearch`（查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `code` | `string` | contains | 按契约编码片段过滤 |
| `name` | `string` | contains | 按名称片段过滤 |
| `isEnabled` | `bool` | equals | 按启用状态过滤；不传则两者皆匹配 |
| `labels` | `object`（string→string） | 每对相等 | 宿主驱动的 label 过滤（业务侧契约选择器按 label 选取） |

`ContractParams`（create/update）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | 仅 update | 要更新记录的主键 |
| `code` | `string` | 是 | 业务代码调用的唯一契约编码 |
| `name` | `string` | 是 | 显示名 |
| `description` | `string` | 否 | 自由描述 |
| `labels` | `object`（string→string） | 否 | 宿主自有的筛选元数据；键不允许含点号，键/值有长度上限（`ErrInvalidLabel`） |
| `inputSchema` | JSON Schema 对象 | 否 | 自包含的 draft 2020-12 输入 Schema；为空则跳过输入校验 |
| `outputSchema` | JSON Schema 对象 | 否 | 适配器返回值的自包含 Schema；为空则跳过输出校验 |
| `isEnabled` | `bool` | 否 | 禁用的契约拒绝调用（`ErrContractDisabled`） |

删除仍被路由引用的契约会返回标准外键冲突错误（路由表的契约列携带空字符串
通配哨兵，因此该检查由资源自身强制）。

## `integration/system`

外部系统定义。写入时加密敏感认证参数与数据源密码；读取时始终掩码为
`"******"`。更新时提交掩码保持已存值不变。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `integration.system.query` | `SystemSearch` + 分页 meta | `page.Page[System]`（已掩码） |
| `find_all` | `integration.system.query` | `SystemSearch` | `System[]`（已掩码） |
| `create` | `integration.system.create` | `SystemParams` | 创建后的 `System` |
| `update` | `integration.system.update` | `SystemParams` | 更新后的 `System` |
| `delete` | `integration.system.delete` | 主键参数（`params.id`） | 成功 |

删除系统——或在更新中移除/改名其数据源——会释放其数据源注册表条目。

`SystemSearch`（查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `code` | `string` | contains | 按系统编码片段过滤 |
| `name` | `string` | contains | 按名称片段过滤 |
| `isEnabled` | `bool` | equals | 按启用状态过滤 |

`SystemParams`（create/update）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | 仅 update | 主键 |
| `code` | `string` | 是 | 唯一系统编码 |
| `name` | `string` | 是 | 显示名 |
| `baseUrl` | `string` | 否 | 绝对基础 URL；为该系统脚本启用作用域 `http` 库。保存时校验（`ErrInvalidBaseURL`） |
| `outboundAuth` | `OutboundAuthConfig` | 否 | 出站认证（见下）；`null` 表示请求不带认证发出 |
| `outboundEnvelope` | `OutboundEnvelopeConfig` | 否 | 系统级请求/响应包装脚本（见下）；`null` 表示适配器请求原样透传 |
| `inboundAuth` | `InboundAuthConfig` | 否 | 入站验证（见下）；`null` 表示完全拒绝入站投递 |
| `dataSource` | `DataSourceConfig` | 否 | 直连数据库（见下）；启用作用域 `sql` 库 |
| `params` | `object`（string→string） | 否 | 非敏感的系统级参数，脚本经 `system.params` 可见 |
| `timeoutMs` | `int` | 否 | 单次 HTTP 调用上限；0 使用框架默认 |
| `retry` | `RetryPolicy` | 否 | 出站重试策略（见下） |
| `isEnabled` | `bool` | 否 | 禁用的系统拒绝两个流向（出站 `ErrSystemDisabled`；入站统一按认证失败拒绝） |

`OutboundAuthConfig` / `InboundAuthConfig`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `scheme` | `string` | scheme 名。出站：`none`、`http_basic`、`bearer`、`header`、`query`、`signature`、`script` 或自定义。入站额外支持 `ip` |
| `params` | `object`（string→string） | scheme 参数；scheme 声明的敏感参数值加密存储、响应中掩码 |
| `script` | `string` | `script` scheme 的自定义签名/验证脚本体；在零 IO 运行时中执行 |

各 scheme 的参数参考见[出站调用](./outbound#出站认证-scheme)与
[入站投递](./inbound#入站认证)。

`OutboundEnvelopeConfig`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `request` | `string` | 包装脚本：以 `request`（`{ method, path, headers, query, body }`）接收适配器发出的请求，返回真正上线的请求；省略的字段保持适配器原值 |
| `response` | `string` | 解包脚本：以 `response`（fetch Response 形态）接收完成的 HTTP 响应；其返回值即适配器调用的所得 |

信封存在时至少要配置两者之一，且系统必须具备 HTTP 传输
（`ErrInvalidEnvelope`）。

`DataSourceConfig`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `kind` | `string` | 数据库类型（与 `vef.data_sources.type` 相同词汇：`postgres`、`mysql`、`sqlite`、`sqlserver`、`oracle`） |
| `mode` | `string` | 脚本写权限：`read_only`（默认；`sql.execute` 抛错）或 `read_write`（启用 `sql.execute`） |
| `host` | `string` | 服务器主机 |
| `port` | `int` | 服务器端口 |
| `user` | `string` | 登录用户 |
| `password` | `string` | 登录密码——加密存储、响应中掩码 |
| `database` | `string` | 数据库名 |
| `schema` | `string` | Schema 名（支持的方言） |
| `path` | `string` | 文件路径（sqlite） |
| `sslMode` | `string` | SSL 模式（与 `vef.data_sources.ssl_mode` 相同词汇） |
| `sslRootCert` | `string` | CA 证书路径 |

`RetryPolicy`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `maxAttempts` | `int` | 总尝试次数，含首个调用 |
| `initialBackoffMs` | `int` | 首次重试前的基础延迟；0 使用 httpx 默认 |
| `maxBackoffMs` | `int` | 尝试间延迟上限；0 使用 httpx 默认 |

## `integration/adapter`

适配器绑定。脚本在保存时做编译检查；绑定关系本身由数据库唯一键与外键守护。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `integration.adapter.query` | `AdapterSearch` + 分页 meta | `page.Page[Adapter]` |
| `find_all` | `integration.adapter.query` | `AdapterSearch` | `Adapter[]` |
| `create` | `integration.adapter.create` | `AdapterParams` | 创建后的 `Adapter` |
| `update` | `integration.adapter.update` | `AdapterParams` | 更新后的 `Adapter` |
| `delete` | `integration.adapter.delete` | 主键参数（`params.id`） | 成功 |

`AdapterSearch`（查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `systemId` | `string` | equals | 按所属系统过滤 |
| `contractId` | `string` | equals | 按绑定契约过滤 |
| `direction` | `string` | equals | `outbound` 或 `inbound` |
| `isEnabled` | `bool` | equals | 按启用状态过滤 |

`AdapterParams`（create/update）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | 仅 update | 主键 |
| `systemId` | `string` | 是 | 适配器所属系统 |
| `contractId` | `string` | 是 | 适配器实现的契约 |
| `direction` | `string` | 否 | `outbound`（省略时默认）或 `inbound`；其他值返回 `ErrInvalidDirection` |
| `script` | `string` | 是 | 翻译脚本；必须可编译（`ErrInvalidScript`） |
| `timeoutMs` | `int` | 否 | 脚本运行超时覆盖；0 继承 `vef.integration.run_timeout` |
| `isEnabled` | `bool` | 否 | 禁用的适配器拒绝调用（`ErrAdapterDisabled`） |

## `integration/route`

路由规则。契约引用在保存时校验，因为契约列携带空字符串通配哨兵、没有外键
（`ErrInvalidRouteRef`）。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `integration.route.query` | `RouteSearch` + 分页 meta | `page.Page[Route]` |
| `find_all` | `integration.route.query` | `RouteSearch` | `Route[]` |
| `create` | `integration.route.create` | `RouteParams` | 创建后的 `Route` |
| `update` | `integration.route.update` | `RouteParams` | 更新后的 `Route` |
| `delete` | `integration.route.delete` | 主键参数（`params.id`） | 成功 |

`RouteSearch`（查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `routeKey` | `string` | contains | 按路由键片段过滤 |
| `contractId` | `string` | equals | 按作用契约过滤 |
| `systemId` | `string` | equals | 按目标系统过滤 |
| `isEnabled` | `bool` | equals | 按启用状态过滤 |

`RouteParams`（create/update）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | 仅 update | 主键 |
| `routeKey` | `string` | 否 | 该规则服务的键（租户、分支机构、院区）；空为默认路由 |
| `contractId` | `string` | 否 | 将规则限定到一个契约；空表示适用于所有契约。精确 `(key, contract)` 匹配优先于契约通配匹配 |
| `systemId` | `string` | 是 | 命中后提供服务的系统 |
| `isEnabled` | `bool` | 否 | 禁用的规则永不命中 |

## `integration/code_map`

按系统的值翻译表。条目在保存时构建索引，冲突或格式错误的映射永远不会进入
查找；宿主注册了可枚举码值目录时，`codeSet` 标识还必须是目录中已注册的
集合。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `integration.code_map.query` | `CodeMapSearch` + 分页 meta | `page.Page[CodeMap]` |
| `find_all` | `integration.code_map.query` | `CodeMapSearch` | `CodeMap[]` |
| `create` | `integration.code_map.create` | `CodeMapParams` | 创建后的 `CodeMap` |
| `update` | `integration.code_map.update` | `CodeMapParams` | 更新后的 `CodeMap` |
| `delete` | `integration.code_map.delete` | 主键参数（`params.id`） | 成功 |

`CodeMapSearch`（查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `systemId` | `string` | equals | 按所属系统过滤 |
| `codeSet` | `string` | contains | 按码值集标识片段过滤 |
| `name` | `string` | contains | 按名称片段过滤 |
| `isEnabled` | `bool` | equals | 按启用状态过滤 |

`CodeMapParams`（create/update）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | 仅 update | 主键 |
| `systemId` | `string` | 是 | 所属系统 |
| `codeSet` | `string` | 是 | 被翻译码值集标识（如 `gender`）；宿主注册目录时受目录约束 |
| `name` | `string` | 是 | 显示名 |
| `entries` | `CodeMapEntry[]` | 否 | 映射对（见下）；任一侧重复查找值被拒绝（`ErrInvalidCodeMap`） |
| `onUnmapped` | `string` | 否 | `reject`（省略时默认——fail closed）、`passthrough` 或 `fallback` |
| `fallbackCanonical` | 任意 JSON 值 | 否 | `fallback` 策略下 `toCanonical` 对未映射输入返回的值 |
| `fallbackExternal` | 任意 JSON 值 | 否 | `fallback` 策略下 `toExternal` 对未映射输入返回的值 |
| `isEnabled` | `bool` | 否 | 禁用的映射视同不存在（`ErrMissingCodeMap`） |

`CodeMapEntry`：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `canonical` | 字符串 / 数字 / 布尔 | 是 | 宿主侧主值，`toCanonical` 查找的输出 |
| `external` | 字符串 / 数字 / 布尔 | 是 | 外部侧主值，`toExternal` 查找的输出 |
| `canonicalAliases` | 数组 | 否 | 匹配该条目的额外宿主侧值（只匹配，不输出） |
| `externalAliases` | 数组 | 否 | 匹配该条目的额外外部侧值 |

## `integration/code_set`

宿主标准码值目录的只读视图，服务于映射编辑器的选择器。它按能力降级：没有
`mold.CodeSetInspector` 时两个操作都返回 `supported: false`。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `list_code_sets` | `integration.code_map.query` | 无 | `CodeSetCatalog` |
| `list_codes` | `integration.code_map.query` | `ListCodesParams` | `CodeCatalog` |

`ListCodesParams`：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `codeSet` | `string` | 是 | 要枚举的码值集 |

`CodeSetCatalog` 响应：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `supported` | `bool` | 宿主未注册可枚举目录时为 `false`（编辑器退化为自由文本输入） |
| `codeSets` | `CodeSetInfo[]` | 条目含 `codeSet`（标识）与 `name`（显示名） |

`CodeCatalog` 响应：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `supported` | `bool` | 同上 |
| `codes` | `CodeInfo[]` | 条目含 `code`（标准值）与 `label`（显示名） |

## `integration/log`

只读调用日志：分页视图用于浏览，单条视图用于查看完整捕获。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `integration.log.query` | `LogSearch` + 分页 meta | `page.Page[InvocationLog]` |
| `find_one` | `integration.log.query` | `LogSearch` | 一条 `InvocationLog` |

`LogSearch`（查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `systemCode` | `string` | equals | 按系统编码过滤 |
| `contractCode` | `string` | equals | 按契约编码过滤 |
| `direction` | `string` | equals | `outbound` 或 `inbound` |
| `failureKind` | `string` | equals | [失败分类](./overview#失败分类)之一；空行为成功 |
| `requestId` | `string` | equals | 关联触发调用的 API 请求 |

`InvocationLog` 响应字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 日志行 ID |
| `systemCode` | `string` | 提供服务（或拒绝）的系统 |
| `contractCode` | `string` | 被调用的契约 |
| `direction` | `string` | `outbound` 或 `inbound` |
| `failureKind` | `string` | 失败分类；成功为空 |
| `durationMs` | `int` | 调用耗时 |
| `input` | JSON | 捕获的标准输入（按 `vef.integration.log` 掩码、截断） |
| `output` | JSON | 捕获的标准输出（掩码、截断） |
| `httpTrace` | `HTTPExchange[]` | 脚本运行期间捕获的线上交换（见下） |
| `error` | `string` | 失败消息；成功时缺省 |
| `requestId` | `string` | 来源 API 请求 ID |
| `createdAt` / `createdBy` | 时间戳 / `string` | 创建审计字段 |

`HTTPExchange`（日志与 dry-run 追踪共用）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `method` | `string` | HTTP 方法 |
| `url` | `string` | 请求 URL（已掩码） |
| `requestHeaders` | `object` | 请求头（凭证头始终掩码） |
| `requestBody` | `string` | 捕获的请求体（掩码、截断） |
| `status` | `int` | 响应状态；调用未完成时为 `0` |
| `responseHeaders` | `object` | 响应头 |
| `responseBody` | `string` | 捕获的响应体（掩码、截断） |
| `durationMs` | `int` | 交换耗时 |
| `error` | `string` | 调用失败时的传输错误消息 |

## `integration/ops`

运维端点：脚本测试台、连接探测与路由诊断。dry run 与探测对禁用的定义同样
有效——测试先于启用。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `dry_run` | `integration.ops.dry_run` | `DryRunParams` | `DryRunResult` |
| `dry_run_inbound` | `integration.ops.dry_run_inbound` | `DryRunInboundParams` | `InboundDryRunResult` |
| `test_connection` | `integration.ops.test_connection` | `TestConnectionParams` | `ConnectionCheck` |
| `diagnose_routes` | `integration.ops.diagnose_routes` | 无 | `RouteDiagnostics` |

### `dry_run`

在契约下对系统执行脚本，返回输出、失败分类与完整线上追踪。它发出的调用是
真实的；不记入统计与调用日志。

请求（`DryRunParams`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `systemCode` | `string` | 是 | 目标系统（允许禁用状态） |
| `contractCode` | `string` | 是 | 以其 Schema 约束本次运行的契约 |
| `script` | `string` | 否 | 编辑器中未保存的内容；为空回退到已保存的出站适配器脚本（不存在时 `ErrAdapterNotFound`） |
| `input` | 任意 JSON 值 | 否 | 调用输入，经契约输入 Schema 校验 |

响应（`DryRunResult`）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `output` | 任意 JSON 值 | 脚本返回值（已通过 Schema 校验）；失败时为 `null` |
| `trace` | `HTTPExchange[]` | 线上交换，即使运行失败也会填充——运维可看到脚本走到了哪里 |
| `failureKind` | `string` | 失败分类；成功时缺省 |
| `error` | `string` | 失败消息；成功时缺省 |

### `dry_run_inbound`

对合成的外部请求执行入站脚本，业务处理器被桩替换为返回给定的样例输出。
验证被绕过（测试台测的是翻译，不是凭证），不触碰业务代码、不做任何记录；
契约 Schema 在分发两侧真实强制。

请求（`DryRunInboundParams`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `systemCode` | `string` | 是 | 目标系统 |
| `contractCode` | `string` | 是 | 以其 Schema 约束分发的契约 |
| `script` | `string` | 否 | 编辑器中未保存的内容；为空回退到已保存的入站适配器脚本 |
| `request` | `InboundRequestParams` | 否 | 合成的外部请求（见下） |
| `handlerOutput` | 任意 JSON 值 | 否 | 桩业务处理器返回的样例；经输出 Schema 校验 |

`InboundRequestParams`（全部可选）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `method` | `string` | 合成请求的 HTTP 方法 |
| `path` | `string` | 请求路径 |
| `headers` | `object`（string→string） | 头名会归一化为小写，与真实网关投递一致 |
| `query` | `object`（string→string） | 查询参数 |
| `body` | `string` | 原始请求负载 |

响应（`InboundDryRunResult`）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `reply` | 任意 JSON 值 | 外部系统将收到的应答（脚本返回 `$response` 信封时含该信封） |
| `dispatchedInput` | 任意 JSON 值 | 脚本分发给（桩）处理器的内容——单值，或脚本多次分发时的数组 |
| `failureKind` | `string` | 失败分类；成功时缺省 |
| `error` | `string` | 失败消息；成功时缺省 |

### `test_connection`

探测已保存系统配置的每种传输。探测失败是数据（`reachable: false`）而非
错误——探测已经回答了问题。配置故障（未知认证 scheme、凭证无法解密）则返回
错误。

请求（`TestConnectionParams`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `systemCode` | `string` | 是 | 要探测的系统 |
| `method` | `string` | 否 | 探测 HTTP 方法；默认 `GET` |
| `path` | `string` | 否 | 相对基础 URL 的探测路径；默认 `/` |

响应（`ConnectionCheck`；系统配置了某传输时对应探测才出现）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `http` | `HTTPProbe` | 系统配置了 `baseUrl` 时出现 |
| `http.reachable` | `bool` | 请求是否完成 |
| `http.status` | `int` | 可达时的响应状态 |
| `http.statusText` | `string` | 可达时的状态文本 |
| `http.durationMs` | `int` | 探测耗时 |
| `http.error` | `string` | 不可达时的传输错误 |
| `database` | `DatabaseProbe` | 系统配置了 `dataSource` 时出现 |
| `database.reachable` | `bool` | 一次性连接是否成功 |
| `database.version` | `string` | 成功时的服务器版本 |
| `database.durationMs` | `int` | 探测耗时 |
| `database.error` | `string` | 不可达时的连接错误 |

### `diagnose_routes`

在配置缺口变成运行期错误之前报告路由表的问题——悬挂的适配器、被禁用的
目标、未覆盖的契约。按需计算；无参数。

响应（`RouteDiagnostics`）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `findings` | `RouteFinding[]` | 每个缺口一条（见下）；空表示路由表一致 |

`RouteFinding`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `kind` | `string` | 发现分类（见下） |
| `routeId` | `string` | 涉及的路由行（存在时） |
| `routeKey` | `string` | 始终有意义——`""` 为默认路由 |
| `contractCode` / `contractName` | `string` | 涉及契约的编码与显示名 |
| `systemCode` / `systemName` | `string` | 涉及系统的编码与显示名 |

`kind` 词汇表：

| 类别 | 含义 |
| --- | --- |
| `dangling_adapter` | 契约限定路由的目标系统对该契约没有已启用适配器——经此规则调用将得到 `ErrAdapterNotFound` |
| `wildcard_gap` | 某个已启用契约无法由通配（或默认）路由服务，因为目标系统对它没有已启用适配器。提示性 |
| `disabled_system` | 已启用路由指向被禁用的系统——经此调用将得到 `ErrSystemDisabled` |
| `disabled_contract` | 已启用路由限定到被禁用的契约——该规则永远无法命中成功调用 |
| `uncovered_contract` | 某个已启用契约在路由表现有的某个键下解析不到任何规则——用该键调用将得到 `ErrRouteNotFound`。当该键有意只路由子集时为提示性 |

## 另请参阅

- [概览](./overview) —— 定义模型与错误码
- [出站调用](./outbound)与[入站投递](./inbound) —— `dry_run` / `dry_run_inbound` 背后的流向
- [内置资源](../reference/built-in-resources) —— 框架级资源索引
