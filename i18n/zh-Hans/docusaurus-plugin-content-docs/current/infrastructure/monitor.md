---
sidebar_position: 9
---

# 监控

VEF 内置了一个监控 service，以及一个用于运行时检查的内置资源。

## 模块输出

监控模块会提供：

| 输出 | 含义 |
| --- | --- |
| `monitor.Service` | 运行时监控服务 |
| `sys/monitor` | 内置 RPC 资源 |

当需要时，service 会通过生命周期 hook 自动初始化与关闭。

## `monitor.Service` 接口

公开监控 service 暴露的方法如下：

| 方法 | 返回类型 | 作用 |
| --- | --- | --- |
| `Overview(ctx)` | `*monitor.SystemOverview` | 返回综合概览快照 |
| `CPU(ctx)` | `*monitor.CPUInfo` | 返回 CPU 详情与使用率 |
| `Memory(ctx)` | `*monitor.MemoryInfo` | 返回虚拟内存与 swap 详情 |
| `Disk(ctx)` | `*monitor.DiskInfo` | 返回磁盘分区与 I/O 详情 |
| `Network(ctx)` | `*monitor.NetworkInfo` | 返回网络接口与 I/O 详情 |
| `Host(ctx)` | `*monitor.HostInfo` | 返回主机静态元数据 |
| `Process(ctx)` | `*monitor.ProcessInfo` | 返回当前进程详情 |
| `Load(ctx)` | `*monitor.LoadInfo` | 返回系统负载 |
| `BuildInfo()` | `*monitor.BuildInfo` | 返回构建元数据 |

## 内置资源

monitor 模块注册 `sys/monitor` RPC 资源，挂载在 `/api` 下，使用标准请求
envelope（`resource`、`action`、`version`、`params`、`meta`）。没有任何
操作是公开的，也没有声明专门的权限点：每个 action 都继承 API 引擎默认的
Bearer 认证。

每个 action 都单独设置了 `max 60` 的限流上限。窗口长度未覆写，因此继承
`vef.api.rate_limit.period`（默认 `5m`）；限流按「操作 + 客户端 IP +
principal」计数，每个节点在进程内存中独立执行。

这些 action 都没有定义框架级入参：`params` 会被忽略，可以完全省略。

| Action | 访问 | 限流 | 入参 | 出参 |
| --- | --- | --- | --- | --- |
| `get_overview` | Bearer 认证 | `max 60` | 无 | `monitor.SystemOverview` |
| `get_cpu` | Bearer 认证 | `max 60` | 无 | `monitor.CPUInfo` |
| `get_memory` | Bearer 认证 | `max 60` | 无 | `monitor.MemoryInfo` |
| `get_disk` | Bearer 认证 | `max 60` | 无 | `monitor.DiskInfo` |
| `get_network` | Bearer 认证 | `max 60` | 无 | `monitor.NetworkInfo` |
| `get_host` | Bearer 认证 | `max 60` | 无 | `monitor.HostInfo` |
| `get_process` | Bearer 认证 | `max 60` | 无 | `monitor.ProcessInfo` |
| `get_load` | Bearer 认证 | `max 60` | 无 | `monitor.LoadInfo` |
| `get_build_info` | Bearer 认证 | `max 60` | 无 | `monitor.BuildInfo` |
| `get_event_streams` | Bearer 认证 | `max 60` | 无 | `monitor.EventStreamsInfo` |
| `get_integration_stats` | Bearer 认证 | `max 60` | 无 | `monitor.IntegrationStatsInfo` |

源码中可见的行为语义：

- `get_overview` 是尽力而为的，整体从不失败：某个子探针出错时只记日志，
  对应的 overview 字段留为 `null`，单个损坏的采集器不会掩盖其余数据。
- `get_cpu` 和 `get_process` 从后台采样缓存读取，在第一次采样落地前返回
  monitor-not-ready 业务错误（`monitor.ErrNotReady`）。
- `get_memory`、`get_disk`、`get_network`、`get_host`、`get_load` 实时读取
  探针；探针失败映射为 `monitor.ErrCollectionFailed`。
- `get_build_info` 不会失败：service 始终持有非 nil 的构建信息对象（见
  [构建信息行为](#构建信息行为)）。
- `get_event_streams` 依赖可选的 `event.StreamInspector`。inspector 为 nil
  （redis_stream transport 未启用）时仍返回 `200 OK`，只是 `enabled: false`
  且 `streams` 为空列表，不会报错；inspector 读取出错则映射为
  `monitor.ErrCollectionFailed`。
- `get_integration_stats` 对可选的 `integration.StatsInspector`（集成模块
  未启用时为 nil）采用相同的降级模式：`enabled: false` 且 `stats` 为空
  列表。读取进程内快照本身不会失败。
- 业务错误使用标准 result envelope：HTTP 状态保持 `200`，失败通过 body 的
  `code` 传递。

## 错误 API

| API | 含义 |
| --- | --- |
| `monitor.ErrNotReady` / `ErrCodeNotReady`（`2100`） | CPU 或 process 这类依赖采样的数据尚未准备好 |
| `monitor.ErrCollectionFailed` / `ErrCodeCollectionFailed`（`2101`） | 某个 runtime 探针采集数据失败 |

## 默认采样配置

当没有显式提供 monitor 配置时，模块默认使用：

| 配置项 | 默认值 |
| --- | --- |
| `vef.monitor.sample_interval` | `10s` |
| `vef.monitor.sample_duration` | `2s` |

这些配置驱动 `get_cpu` 与 `get_process` 背后的后台采样器：启动时立即采样
一次，之后每个采样间隔采样一次，每次采样在一个采样窗口内测量使用率。在
第一次采样完成前（大约是启动后的第一个窗口），这两个 action 都会返回
`monitor.ErrNotReady`。

## 构建信息行为

monitor 模块会对构建信息做装饰，保证 `vefVersion` 一定存在，即使应用没有提供完整构建元数据对象。

回退行为如下：

| 字段 | 当应用没有提供构建信息时的回退值 |
| --- | --- |
| `appVersion` | `unknown` |
| `buildTime` | `unknown` |
| `gitCommit` | `unknown` |
| `vefVersion` | 当前框架版本 |

## 按 Action 划分的响应结构

以下字段名即 JSON wire 名称（Go struct 的 json tag）。字节量都是纯字节
数（无单位换算），百分比范围 `0`–`100`，计数器除特别说明外都是自启动以来
的累计值。平台不提供的字段报告为 `0` 或空。

### `get_overview` — `monitor.SystemOverview`

由全部探针组装的综合快照。某个探针失败时对应字段为 `null`；`build`
始终存在。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `host` | `*monitor.HostSummary` | 简化主机信息 |
| `cpu` | `*monitor.CPUSummary` | 简化 CPU 信息 |
| `memory` | `*monitor.MemorySummary` | 简化内存使用情况 |
| `disk` | `*monitor.DiskSummary` | 简化磁盘使用情况 |
| `network` | `*monitor.NetworkSummary` | 简化网络活动 |
| `process` | `*monitor.ProcessSummary` | 简化当前进程指标 |
| `load` | `*monitor.LoadInfo` | 负载均值（与 `get_load` 同构） |
| `build` | `*monitor.BuildInfo` | 构建元数据（与 `get_build_info` 同构） |

#### `monitor.HostSummary`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `hostname` | `string` | 主机名 |
| `os` | `string` | 操作系统 |
| `platform` | `string` | 平台名 |
| `platformVersion` | `string` | 平台版本 |
| `kernelVersion` | `string` | 内核版本 |
| `kernelArch` | `string` | 内核架构 |
| `uptime` | `uint64` | 主机运行时长（秒） |

#### `monitor.CPUSummary`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `physicalCores` | `int` | 物理核心数（宿主机拓扑） |
| `logicalCores` | `int` | 逻辑核心数（宿主机拓扑） |
| `usagePercent` | `float64` | 最近一个采样窗口的聚合 CPU 使用率，按 `effectiveCores` 归一化 |
| `effectiveCores` | `float64` | 归一化使用率所用的算力容量：容器内是 cgroup CPU 配额（支持 v1 与 v2），无法一致采样受限用量时回退为 `logicalCores` |

#### `monitor.MemorySummary`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `total` | `uint64` | 总内存（字节） |
| `used` | `uint64` | 已用内存（字节） |
| `usedPercent` | `float64` | 内存使用率 |

监控是容器感知的：当进程运行在实际限制内存的 cgroup（v2 或
v1）下时，头部指标（`total`、`used`、`usedPercent` 以及 `VirtualMemory`
的可用/空闲）反映 cgroup 限额与 cgroup 自身用量，而不是宿主机全量——64
GiB 宿主机上的 512 MiB 容器按 512 MiB 报告。没有限额时仍报告宿主机数据。

#### `monitor.DiskSummary`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `total` | `uint64` | 根文件系统总大小（字节） |
| `used` | `uint64` | 根文件系统已用大小（字节） |
| `usedPercent` | `float64` | 根文件系统使用率 |
| `partitions` | `int` | 恒为 `1`（摘要只覆盖单个文件系统） |

overview 的磁盘摘要报告**进程根路径所在的文件系统**，而不是累加
所有挂载分区——远程挂载、磁盘镜像和并列卷不会虚增宿主机容量，也不存在
`vef.monitor.excluded_mounts` 配置（本就不做累加，无需排除）。完整挂载
清单仍可通过 `DiskInfo.partitions` 获取。

#### `monitor.NetworkSummary`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `interfaces` | `int` | 网卡数量 |
| `bytesSent` | `uint64` | 发送字节总量（跨网卡累加） |
| `bytesRecv` | `uint64` | 接收字节总量（跨网卡累加） |
| `packetsSent` | `uint64` | 发送包总量（跨网卡累加） |
| `packetsRecv` | `uint64` | 接收包总量（跨网卡累加） |

#### `monitor.ProcessSummary`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `pid` | `int32` | 进程 ID |
| `name` | `string` | 进程名 |
| `cpuPercent` | `float64` | 最近一个采样窗口的进程 CPU 使用率；以单核为基准，多核机器上可以超过 `100` |
| `memoryPercent` | `float32` | 进程占宿主机总内存的百分比 |

### `get_cpu` — `monitor.CPUInfo`

从后台采样缓存读取：每个采样间隔（默认 `10s`）刷新一次，每次刷新在一个
采样窗口（默认 `2s`）内测量。清单类字段（`modelName`、`vendorId`、
`family`、`model`、`stepping`、`microcode`、`mhz`、`cacheSize`）描述第一
颗 CPU。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `physicalCores` | `int` | 物理核心数（宿主机拓扑） |
| `logicalCores` | `int` | 逻辑核心数（宿主机拓扑） |
| `modelName` | `string` | CPU 型号 |
| `mhz` | `float64` | 标称主频（MHz） |
| `cacheSize` | `int32` | 缓存大小（KB） |
| `usagePercent` | `[]float64` | 采样窗口内的每核心使用率，每个逻辑核心一项；在 CPU 受限的容器内为 `null`（cgroup 测量取代每核心采样） |
| `totalPercent` | `float64` | 聚合使用率：每核心采样的均值；在 CPU 受限的容器内则是窗口内消耗的 cgroup 算力份额，上限 `100` |
| `vendorId` | `string` | vendor 标识 |
| `family` | `string` | CPU family |
| `model` | `string` | CPU model |
| `stepping` | `int32` | stepping |
| `microcode` | `string` | microcode 版本 |
| `effectiveCores` | `float64` | 归一化使用率所用的算力容量；见 `CPUSummary.effectiveCores` |

### `get_memory` — `monitor.MemoryInfo`

每次调用实时读取。`MemorySummary` 描述的容器感知头部指标行为同样适用于
`virtual`。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `virtual` | `*monitor.VirtualMemory` | 虚拟/物理内存详情 |
| `swap` | `*monitor.SwapMemory` | swap 详情；swap 探针失败时为 `null` |

#### `monitor.VirtualMemory`

除 `usedPercent`（百分比）和 huge page 计数器外，所有字段都是字节量：
`hugePagesTotal`、`hugePagesFree`、`hugePagesReserved`、`hugePagesSurplus`
是页数，`hugePageSize` 和 `anonHugePages` 是字节。即使在内存受限的容器
内，明细字段仍保持宿主机含义。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `total` | `uint64` | 总内存 |
| `available` | `uint64` | 可用内存 |
| `used` | `uint64` | 已用内存 |
| `usedPercent` | `float64` | 使用率 |
| `free` | `uint64` | 空闲内存 |
| `active` | `uint64` | 活跃内存 |
| `inactive` | `uint64` | 非活跃内存 |
| `wired` | `uint64` | wired 内存 |
| `laundry` | `uint64` | laundry 页数 |
| `buffers` | `uint64` | buffer 内存 |
| `cached` | `uint64` | 缓存内存 |
| `writeBack` | `uint64` | write-back 页数 |
| `dirty` | `uint64` | dirty 页数 |
| `writeBackTmp` | `uint64` | 临时 write-back 页数 |
| `shared` | `uint64` | 共享内存 |
| `slab` | `uint64` | slab 内存 |
| `slabReclaimable` | `uint64` | 可回收 slab |
| `slabUnreclaimable` | `uint64` | 不可回收 slab |
| `pageTables` | `uint64` | 页表占用 |
| `swapCached` | `uint64` | swap 缓存 |
| `commitLimit` | `uint64` | commit 上限 |
| `committedAs` | `uint64` | committed 内存 |
| `highTotal` | `uint64` | high memory 总量 |
| `highFree` | `uint64` | high memory 空闲量 |
| `lowTotal` | `uint64` | low memory 总量 |
| `lowFree` | `uint64` | low memory 空闲量 |
| `swapTotal` | `uint64` | swap 总量 |
| `swapFree` | `uint64` | swap 空闲量 |
| `mapped` | `uint64` | mapped 内存 |
| `vmAllocTotal` | `uint64` | VM 分配总量 |
| `vmAllocUsed` | `uint64` | VM 已用分配量 |
| `vmAllocChunk` | `uint64` | VM 分配块 |
| `hugePagesTotal` | `uint64` | huge page 总量（页数） |
| `hugePagesFree` | `uint64` | huge page 空闲量（页数） |
| `hugePagesReserved` | `uint64` | huge page 预留量（页数） |
| `hugePagesSurplus` | `uint64` | huge page surplus（页数） |
| `hugePageSize` | `uint64` | huge page 大小（字节） |
| `anonHugePages` | `uint64` | 匿名 huge page（字节） |

#### `monitor.SwapMemory`

`total`、`used`、`free` 是字节。`swapIn`、`swapOut`、`pageIn`、`pageOut`
是由内核页计数换算的累计字节量；`pageFault`、`pageMajorFault` 是累计事件
次数。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `total` | `uint64` | swap 总量 |
| `used` | `uint64` | 已用 swap |
| `free` | `uint64` | 空闲 swap |
| `usedPercent` | `float64` | swap 使用率 |
| `swapIn` | `uint64` | swap 换入量 |
| `swapOut` | `uint64` | swap 换出量 |
| `pageIn` | `uint64` | page 换入量 |
| `pageOut` | `uint64` | page 换出量 |
| `pageFault` | `uint64` | page fault 数 |
| `pageMajorFault` | `uint64` | major page fault 数 |

### `get_disk` — `monitor.DiskInfo`

每次调用实时读取。用量探针失败的分区会被跳过、不出现在 `partitions`
中；I/O 计数探针失败时 `ioCounters` 为 `null`。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `partitions` | `[]*monitor.PartitionInfo` | 每个挂载点的分区详情 |
| `ioCounters` | `map[string]*monitor.IOCounter` | 每设备 I/O 计数，key 为设备名 |

#### `monitor.PartitionInfo`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `device` | `string` | 设备名 |
| `mountPoint` | `string` | 挂载点 |
| `fsType` | `string` | 文件系统类型 |
| `options` | `[]string` | 挂载选项 |
| `total` | `uint64` | 总大小（字节） |
| `free` | `uint64` | 空闲大小（字节） |
| `used` | `uint64` | 已用大小（字节） |
| `usedPercent` | `float64` | 使用率 |
| `iNodesTotal` | `uint64` | inode 总量 |
| `iNodesUsed` | `uint64` | 已用 inode |
| `iNodesFree` | `uint64` | 空闲 inode |
| `iNodesUsedPercent` | `float64` | inode 使用率 |

#### `monitor.IOCounter`

计数器为自启动以来的累计值；`readTime`、`writeTime`、`ioTime`、
`weightedIo` 单位为毫秒。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `readCount` | `uint64` | 读操作次数 |
| `mergedReadCount` | `uint64` | 合并读次数 |
| `writeCount` | `uint64` | 写操作次数 |
| `mergedWriteCount` | `uint64` | 合并写次数 |
| `readBytes` | `uint64` | 读取字节数 |
| `writeBytes` | `uint64` | 写入字节数 |
| `readTime` | `uint64` | 读耗时 |
| `writeTime` | `uint64` | 写耗时 |
| `iopsInProgress` | `uint64` | 正在进行的 I/O 数 |
| `ioTime` | `uint64` | I/O 总耗时 |
| `weightedIo` | `uint64` | 加权 I/O 时间 |
| `name` | `string` | 设备名 |
| `serialNumber` | `string` | 设备序列号 |
| `label` | `string` | 设备标签 |

### `get_network` — `monitor.NetworkInfo`

每次调用实时读取。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `interfaces` | `[]*monitor.InterfaceInfo` | 网卡元数据 |
| `ioCounters` | `map[string]*monitor.NetIOCounter` | 每网卡 I/O 计数，key 为网卡名 |

#### `monitor.InterfaceInfo`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `index` | `int` | 接口索引 |
| `mtu` | `int` | MTU |
| `name` | `string` | 接口名 |
| `hardwareAddr` | `string` | MAC 地址 |
| `flags` | `[]string` | 接口 flags |
| `addrs` | `[]string` | 绑定地址 |

#### `monitor.NetIOCounter`

按网卡统计，自启动以来累计。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 接口名 |
| `bytesSent` | `uint64` | 发送字节数 |
| `bytesRecv` | `uint64` | 接收字节数 |
| `packetsSent` | `uint64` | 发送包数 |
| `packetsRecv` | `uint64` | 接收包数 |
| `errorsIn` | `uint64` | 入站错误数 |
| `errorsOut` | `uint64` | 出站错误数 |
| `droppedIn` | `uint64` | 入站丢包数 |
| `droppedOut` | `uint64` | 出站丢包数 |
| `fifoIn` | `uint64` | 入站 FIFO 计数 |
| `fifoOut` | `uint64` | 出站 FIFO 计数 |

### `get_host` — `monitor.HostInfo`

主机静态元数据，每次调用实时读取。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `hostname` | `string` | 主机名 |
| `uptime` | `uint64` | 主机运行时长（秒） |
| `bootTime` | `uint64` | 启动时间（Unix 时间戳，秒） |
| `processes` | `uint64` | 宿主机上的进程数 |
| `os` | `string` | 操作系统 |
| `platform` | `string` | 平台名 |
| `platformFamily` | `string` | 平台族 |
| `platformVersion` | `string` | 平台版本 |
| `kernelVersion` | `string` | 内核版本 |
| `kernelArch` | `string` | 内核架构 |
| `virtualizationSystem` | `string` | 虚拟化系统 |
| `virtualizationRole` | `string` | 虚拟化角色 |
| `hostId` | `string` | 主机标识 |

### `get_process` — `monitor.ProcessInfo`

描述应用自身进程。与 `get_cpu` 相同节奏，从后台采样缓存读取。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `pid` | `int32` | 进程 ID |
| `parentPid` | `int32` | 父进程 ID |
| `name` | `string` | 进程名 |
| `exe` | `string` | 可执行文件路径 |
| `commandLine` | `string` | 完整命令行 |
| `cwd` | `string` | 当前工作目录 |
| `status` | `string` | 进程状态 |
| `username` | `string` | 所属用户名 |
| `createTime` | `int64` | 进程创建时间，自 Unix 纪元以来的毫秒数（UTC） |
| `numThreads` | `int32` | 线程数 |
| `numFds` | `int32` | 打开文件描述符数 |
| `cpuPercent` | `float64` | 采样窗口内的进程 CPU 使用率；以单核为基准，多核机器上可以超过 `100` |
| `memoryPercent` | `float32` | 进程占宿主机总内存的百分比 |
| `memoryRss` | `uint64` | 常驻内存 RSS（字节） |
| `memoryVms` | `uint64` | 虚拟内存大小（字节） |
| `memorySwap` | `uint64` | swap 使用量（字节） |

### `get_load` — `monitor.LoadInfo`

每次调用实时读取。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `load1` | `float64` | 1 分钟负载均值 |
| `load5` | `float64` | 5 分钟负载均值 |
| `load15` | `float64` | 15 分钟负载均值 |

### `get_build_info` — `monitor.BuildInfo`

仅返回构建元数据；回退值见[构建信息行为](#构建信息行为)。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `vefVersion` | `string` | 框架版本，模块始终会盖章写入 |
| `appVersion` | `string` | 应用版本 |
| `buildTime` | `string` | 构建时间 |
| `gitCommit` | `string` | Git 提交号 |

### `get_event_streams` — `monitor.EventStreamsInfo`

通过可选的 `event.StreamInspector`（由 redis_stream transport 提供）报告
跨进程 event stream 与 consumer group 状态。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `enabled` | `bool` | 是否有可用的 `event.StreamInspector`（redis_stream transport 已启用）；`false` 表示这是一次空降级返回，不是错误 |
| `streams` | `[]event.StreamInfo` | 每个 transport stream 对应一条记录；`enabled` 为 `false` 时为空 |

#### `event.StreamInfo`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `stream` | `string` | 完整 transport 级 stream key（prefix + event type） |
| `length` | `int64` | stream 当前条目数（trim 之后） |
| `groups` | `[]event.StreamGroupInfo` | 挂在该 stream 上的 consumer group |

#### `event.StreamGroupInfo`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | consumer group 名称（订阅方 `WithGroup` 的值，或其派生默认值） |
| `consumers` | `int64` | 该 group 内已注册的 consumer 记录数，包括已重启进程留下的历史 consumer |
| `pending` | `int64` | 已投递但未 ack 的条目数 |
| `lag` | `int64` | 尚未投递给该 group 的 stream 条目数（trim 后为近似值；部分 Redis server 版本不上报 lag 时为 0） |
| `lastDeliveredId` | `string` | 该 group 最后一次收到投递的 stream ID |

`lag` 持续增长而 consumer 都处于空闲状态的 group，很可能是一个已下线或改名、却没有清理 consumer group 的订阅者遗留下来的孤儿。transport 层细节见 [事件总线](./event-bus)。

### `get_integration_stats` — `monitor.IntegrationStatsInfo`

通过可选的 `integration.StatsInspector` 报告本节点的集成调用统计。
数字是进程启动以来的内存计数——持久记录以调用日志为准。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `enabled` | `bool` | 是否有可用的 `integration.StatsInspector`（集成模块已启用）；`false` 表示这是一次空降级返回，不是错误 |
| `stats` | `[]integration.InvocationStats` | 进程启动以来观察到的每个（系统、契约、方向）组合一条记录，按系统、契约、方向排序；`enabled` 为 `false` 时为空 |

#### `integration.InvocationStats`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `system` | `string` | 服务（或拒绝）该调用的系统 code |
| `contract` | `string` | 被调用的契约 code；被验证拒绝的入站投递聚合在空 `contract` 下——拒绝发生时契约 code 还是未经校验的调用方输入 |
| `direction` | `string` | `outbound` 或 `inbound` |
| `calls` | `int64` | 观察到的调用总数 |
| `successes` | `int64` | 成功完成的调用数 |
| `failures` | `map[string]int64` | 按[失败分类](../integration/overview#失败分类)（`input_invalid`、`output_invalid`、`upstream`、`transport`、`timeout`、`canceled`、`script`、`config`、`auth`、`handler`）计数的失败数；为空时省略 |
| `avgDurationMs` | `int64` | 平均调用耗时（毫秒） |
| `maxDurationMs` | `int64` | 最大调用耗时（毫秒） |
| `lastError` | `string` | 最近一次失败的信息；从未失败时省略 |
| `lastErrorAt` | 时间戳 | 最近一次失败发生的时间；从未失败时省略 |

这些计数如何被记录，见[集成引擎](../integration/overview#统计)。

## 最小请求示例

```json
{
  "resource": "sys/monitor",
  "action": "get_overview",
  "version": "v1"
}
```

## 典型用途

- 运维或后台监控面板
- 健康检查与诊断界面
- 内部开发者工具
- 构建元信息暴露

## 下一步

继续阅读 [CLI 工具](../advanced/cli-tools)，如果你想用 `generate-build-info` 提供更丰富的构建信息，就会接到那里。
