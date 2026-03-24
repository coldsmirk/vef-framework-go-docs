---
sidebar_position: 6
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

monitor 模块会注册：

| 资源 |
| --- |
| `sys/monitor` |

当前 action：

| Action | 输入参数 | 输出类型 | 说明 |
| --- | --- | --- | --- |
| `get_overview` | 无 | `monitor.SystemOverview` | 当部分探针不可用时，overview 仍可能返回部分数据 |
| `get_cpu` | 无 | `monitor.CPUInfo` | 当 CPU 采样缓存尚未准备好时会返回 `monitor not ready` |
| `get_memory` | 无 | `monitor.MemoryInfo` | 包含虚拟内存和 swap |
| `get_disk` | 无 | `monitor.DiskInfo` | 包含分区和 I/O 计数 |
| `get_network` | 无 | `monitor.NetworkInfo` | 包含网卡和 I/O 计数 |
| `get_host` | 无 | `monitor.HostInfo` | 返回主机静态信息 |
| `get_process` | 无 | `monitor.ProcessInfo` | 当进程采样缓存尚未准备好时会返回 `monitor not ready` |
| `get_load` | 无 | `monitor.LoadInfo` | 返回负载均值 |
| `get_build_info` | 无 | `monitor.BuildInfo` | 仅返回构建元数据 |

源码中的实现细节：

- 每个 action 当前都单独设置了 `Max = 60` 的限流上限
- `get_cpu` 和 `get_process` 在采样尚未准备好时会返回专门的 monitor-not-ready 业务错误

## 默认采样配置

当没有显式提供 monitor 配置时，模块默认使用：

| 配置项 | 默认值 |
| --- | --- |
| 采样间隔 | `10s` |
| 采样窗口 | `2s` |

这些默认值主要影响 CPU 与进程采样行为。

## 构建信息行为

monitor 模块会对构建信息做装饰，保证 `vefVersion` 一定存在，即使应用没有提供完整构建元数据对象。

回退行为如下：

| 字段 | 当应用没有提供构建信息时的回退值 |
| --- | --- |
| `appVersion` | `v0.0.0` |
| `buildTime` | `2022-08-08 01:00:00` |
| `gitCommit` | `-` |
| `vefVersion` | 当前框架版本 |

## 数据结构

### `monitor.SystemOverview`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `host` | `*monitor.HostSummary` | 简化主机信息 |
| `cpu` | `*monitor.CPUSummary` | 简化 CPU 信息 |
| `memory` | `*monitor.MemorySummary` | 简化内存使用情况 |
| `disk` | `*monitor.DiskSummary` | 简化磁盘使用情况 |
| `network` | `*monitor.NetworkSummary` | 简化网络活动 |
| `process` | `*monitor.ProcessSummary` | 简化当前进程指标 |
| `load` | `*monitor.LoadInfo` | 负载均值 |
| `build` | `*monitor.BuildInfo` | 构建元数据 |

### `monitor.HostSummary`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `hostname` | `string` | 主机名 |
| `os` | `string` | 操作系统 |
| `platform` | `string` | 平台名 |
| `platformVersion` | `string` | 平台版本 |
| `kernelVersion` | `string` | 内核版本 |
| `kernelArch` | `string` | 内核架构 |
| `uptime` | `uint64` | 运行时长（秒） |

### `monitor.HostInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `hostname` | `string` | 主机名 |
| `uptime` | `uint64` | 运行时长 |
| `bootTime` | `uint64` | 启动时间戳 |
| `processes` | `uint64` | 进程数 |
| `os` | `string` | 操作系统 |
| `platform` | `string` | 平台名 |
| `platformFamily` | `string` | 平台族 |
| `platformVersion` | `string` | 平台版本 |
| `kernelVersion` | `string` | 内核版本 |
| `kernelArch` | `string` | 内核架构 |
| `virtualizationSystem` | `string` | 虚拟化系统 |
| `virtualizationRole` | `string` | 虚拟化角色 |
| `hostId` | `string` | 主机标识 |

### `monitor.CPUSummary`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `physicalCores` | `int` | 物理核心数 |
| `logicalCores` | `int` | 逻辑核心数 |
| `usagePercent` | `float64` | 聚合 CPU 使用率 |

### `monitor.CPUInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `physicalCores` | `int` | 物理核心数 |
| `logicalCores` | `int` | 逻辑核心数 |
| `modelName` | `string` | CPU 型号 |
| `mhz` | `float64` | 主频 |
| `cacheSize` | `int32` | 缓存大小 |
| `usagePercent` | `[]float64` | 每核心使用率 |
| `totalPercent` | `float64` | 总使用率 |
| `vendorId` | `string` | vendor 标识 |
| `family` | `string` | CPU family |
| `model` | `string` | CPU model |
| `stepping` | `int32` | stepping |
| `microcode` | `string` | microcode 版本 |

### `monitor.MemorySummary`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `total` | `uint64` | 总内存 |
| `used` | `uint64` | 已用内存 |
| `usedPercent` | `float64` | 使用率 |

### `monitor.MemoryInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `virtual` | `*monitor.VirtualMemory` | 虚拟/物理内存详情 |
| `swap` | `*monitor.SwapMemory` | swap 详情 |

### `monitor.VirtualMemory`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `total` | `uint64` | 总虚拟内存 |
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
| `hugePagesTotal` | `uint64` | huge page 总量 |
| `hugePagesFree` | `uint64` | huge page 空闲量 |
| `hugePagesReserved` | `uint64` | huge page 预留量 |
| `hugePagesSurplus` | `uint64` | huge page surplus |
| `hugePageSize` | `uint64` | huge page 大小 |
| `anonHugePages` | `uint64` | 匿名 huge page 数量 |

### `monitor.SwapMemory`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `total` | `uint64` | swap 总量 |
| `used` | `uint64` | 已用 swap |
| `free` | `uint64` | 空闲 swap |
| `usedPercent` | `float64` | swap 使用率 |
| `swapIn` | `uint64` | swap-in 次数 |
| `swapOut` | `uint64` | swap-out 次数 |
| `pageIn` | `uint64` | page-in 次数 |
| `pageOut` | `uint64` | page-out 次数 |
| `pageFault` | `uint64` | page fault 数 |
| `pageMajorFault` | `uint64` | major page fault 数 |

### `monitor.DiskSummary`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `total` | `uint64` | 统计到的总磁盘大小 |
| `used` | `uint64` | 已用磁盘大小 |
| `usedPercent` | `float64` | 使用率 |
| `partitions` | `int` | 分区数量 |

### `monitor.DiskInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `partitions` | `[]*monitor.PartitionInfo` | 分区详情 |
| `ioCounters` | `map[string]*monitor.IOCounter` | 每设备 I/O 计数 |

### `monitor.PartitionInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `device` | `string` | 设备名 |
| `mountPoint` | `string` | 挂载点 |
| `fsType` | `string` | 文件系统类型 |
| `options` | `[]string` | 挂载选项 |
| `total` | `uint64` | 总大小 |
| `free` | `uint64` | 空闲大小 |
| `used` | `uint64` | 已用大小 |
| `usedPercent` | `float64` | 使用率 |
| `iNodesTotal` | `uint64` | inode 总量 |
| `iNodesUsed` | `uint64` | 已用 inode |
| `iNodesFree` | `uint64` | 空闲 inode |
| `iNodesUsedPercent` | `float64` | inode 使用率 |

### `monitor.IOCounter`

| 字段 | 类型 | 含义 |
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

### `monitor.NetworkSummary`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `interfaces` | `int` | 网卡数量 |
| `bytesSent` | `uint64` | 发送字节总量 |
| `bytesRecv` | `uint64` | 接收字节总量 |
| `packetsSent` | `uint64` | 发送包总量 |
| `packetsRecv` | `uint64` | 接收包总量 |

### `monitor.NetworkInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `interfaces` | `[]*monitor.InterfaceInfo` | 网卡元数据 |
| `ioCounters` | `map[string]*monitor.NetIOCounter` | 每网卡 I/O 计数 |

### `monitor.InterfaceInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `index` | `int` | 接口索引 |
| `mtu` | `int` | MTU |
| `name` | `string` | 接口名 |
| `hardwareAddr` | `string` | MAC 地址 |
| `flags` | `[]string` | 接口 flags |
| `addrs` | `[]string` | 绑定地址 |

### `monitor.NetIOCounter`

| 字段 | 类型 | 含义 |
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

### `monitor.ProcessSummary`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `pid` | `int32` | 进程 ID |
| `name` | `string` | 进程名 |
| `cpuPercent` | `float64` | CPU 使用率 |
| `memoryPercent` | `float32` | 内存使用率 |

### `monitor.ProcessInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `pid` | `int32` | 进程 ID |
| `parentPid` | `int32` | 父进程 ID |
| `name` | `string` | 进程名 |
| `exe` | `string` | 可执行文件路径 |
| `commandLine` | `string` | 完整命令行 |
| `cwd` | `string` | 当前工作目录 |
| `status` | `string` | 进程状态 |
| `username` | `string` | 所属用户名 |
| `createTime` | `int64` | 创建时间戳 |
| `numThreads` | `int32` | 线程数 |
| `numFds` | `int32` | 打开文件描述符数 |
| `cpuPercent` | `float64` | CPU 使用率 |
| `memoryPercent` | `float32` | 内存使用率 |
| `memoryRss` | `uint64` | RSS 内存 |
| `memoryVms` | `uint64` | 虚拟内存大小 |
| `memorySwap` | `uint64` | swap 使用量 |

### `monitor.LoadInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `load1` | `float64` | 1 分钟负载均值 |
| `load5` | `float64` | 5 分钟负载均值 |
| `load15` | `float64` | 15 分钟负载均值 |

### `monitor.BuildInfo`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `vefVersion` | `string` | 框架版本 |
| `appVersion` | `string` | 应用版本 |
| `buildTime` | `string` | 构建时间 |
| `gitCommit` | `string` | Git 提交号 |

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
