---
sidebar_position: 6
---

# Monitor

VEF includes a monitor service and a built-in resource for runtime inspection.

## Module Outputs

The monitor module provides:

| Output | Meaning |
| --- | --- |
| `monitor.Service` | runtime monitoring service |
| `sys/monitor` | built-in RPC resource |

The service is initialized and closed through lifecycle hooks when needed.

## `monitor.Service` Interface

The public monitoring service exposes:

| Method | Return type | Purpose |
| --- | --- | --- |
| `Overview(ctx)` | `*monitor.SystemOverview` | combined overview snapshot |
| `CPU(ctx)` | `*monitor.CPUInfo` | CPU detail and usage |
| `Memory(ctx)` | `*monitor.MemoryInfo` | virtual and swap memory detail |
| `Disk(ctx)` | `*monitor.DiskInfo` | partitions and disk I/O detail |
| `Network(ctx)` | `*monitor.NetworkInfo` | interfaces and network I/O detail |
| `Host(ctx)` | `*monitor.HostInfo` | static host metadata |
| `Process(ctx)` | `*monitor.ProcessInfo` | current process detail |
| `Load(ctx)` | `*monitor.LoadInfo` | load averages |
| `BuildInfo()` | `*monitor.BuildInfo` | build metadata |

## Built-In Resource

The monitor module registers:

| Resource |
| --- |
| `sys/monitor` |

Current actions:

| Action | Input params | Output type | Notes |
| --- | --- | --- | --- |
| `get_overview` | none | `monitor.SystemOverview` | overview can contain partial data when some probes are unavailable |
| `get_cpu` | none | `monitor.CPUInfo` | returns `monitor not ready` when CPU sample cache is not ready |
| `get_memory` | none | `monitor.MemoryInfo` | includes virtual memory and swap |
| `get_disk` | none | `monitor.DiskInfo` | includes partitions and I/O counters |
| `get_network` | none | `monitor.NetworkInfo` | includes interfaces and I/O counters |
| `get_host` | none | `monitor.HostInfo` | static host metadata |
| `get_process` | none | `monitor.ProcessInfo` | returns `monitor not ready` when process sample cache is not ready |
| `get_load` | none | `monitor.LoadInfo` | load averages |
| `get_build_info` | none | `monitor.BuildInfo` | build metadata only |

Implementation details visible in source:

- each action currently sets a per-operation rate-limit max of `60`
- `get_cpu` and `get_process` use a specific monitor-not-ready business error when samples are not ready

## Default Sampling Configuration

When no explicit monitor config is supplied, the module uses:

| Setting | Default |
| --- | --- |
| sample interval | `10s` |
| sample duration | `2s` |

These defaults primarily affect CPU and process sampling behavior.

## Build Info Behavior

The monitor module decorates build info so that `vefVersion` is always present, even when the application does not provide a complete build metadata object.

Fallback behavior:

| Field | Fallback value when app does not supply build info |
| --- | --- |
| `appVersion` | `v0.0.0` |
| `buildTime` | `2022-08-08 01:00:00` |
| `gitCommit` | `-` |
| `vefVersion` | current framework version |

## Data Shapes

### `monitor.SystemOverview`

| Field | Type | Meaning |
| --- | --- | --- |
| `host` | `*monitor.HostSummary` | condensed host information |
| `cpu` | `*monitor.CPUSummary` | condensed CPU information |
| `memory` | `*monitor.MemorySummary` | condensed memory usage |
| `disk` | `*monitor.DiskSummary` | condensed disk usage |
| `network` | `*monitor.NetworkSummary` | condensed network activity |
| `process` | `*monitor.ProcessSummary` | condensed current process metrics |
| `load` | `*monitor.LoadInfo` | load averages |
| `build` | `*monitor.BuildInfo` | build metadata |

### `monitor.HostSummary`

| Field | Type | Meaning |
| --- | --- | --- |
| `hostname` | `string` | host name |
| `os` | `string` | operating system |
| `platform` | `string` | platform name |
| `platformVersion` | `string` | platform version |
| `kernelVersion` | `string` | kernel version |
| `kernelArch` | `string` | kernel architecture |
| `uptime` | `uint64` | uptime in seconds |

### `monitor.HostInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `hostname` | `string` | host name |
| `uptime` | `uint64` | uptime in seconds |
| `bootTime` | `uint64` | boot timestamp |
| `processes` | `uint64` | number of processes |
| `os` | `string` | operating system |
| `platform` | `string` | platform name |
| `platformFamily` | `string` | platform family |
| `platformVersion` | `string` | platform version |
| `kernelVersion` | `string` | kernel version |
| `kernelArch` | `string` | kernel architecture |
| `virtualizationSystem` | `string` | virtualization system |
| `virtualizationRole` | `string` | virtualization role |
| `hostId` | `string` | host identifier |

### `monitor.CPUSummary`

| Field | Type | Meaning |
| --- | --- | --- |
| `physicalCores` | `int` | number of physical cores |
| `logicalCores` | `int` | number of logical cores |
| `usagePercent` | `float64` | aggregated CPU usage percent |

### `monitor.CPUInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `physicalCores` | `int` | number of physical cores |
| `logicalCores` | `int` | number of logical cores |
| `modelName` | `string` | CPU model name |
| `mhz` | `float64` | clock frequency |
| `cacheSize` | `int32` | cache size |
| `usagePercent` | `[]float64` | per-core usage percentages |
| `totalPercent` | `float64` | total usage percent |
| `vendorId` | `string` | vendor identifier |
| `family` | `string` | CPU family |
| `model` | `string` | CPU model |
| `stepping` | `int32` | CPU stepping |
| `microcode` | `string` | microcode version |

### `monitor.MemorySummary`

| Field | Type | Meaning |
| --- | --- | --- |
| `total` | `uint64` | total memory |
| `used` | `uint64` | used memory |
| `usedPercent` | `float64` | memory usage percentage |

### `monitor.MemoryInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `virtual` | `*monitor.VirtualMemory` | physical or virtual memory detail |
| `swap` | `*monitor.SwapMemory` | swap detail |

### `monitor.VirtualMemory`

| Field | Type | Meaning |
| --- | --- | --- |
| `total` | `uint64` | total virtual memory |
| `available` | `uint64` | available memory |
| `used` | `uint64` | used memory |
| `usedPercent` | `float64` | used percentage |
| `free` | `uint64` | free memory |
| `active` | `uint64` | active memory |
| `inactive` | `uint64` | inactive memory |
| `wired` | `uint64` | wired memory |
| `laundry` | `uint64` | laundry pages |
| `buffers` | `uint64` | buffer memory |
| `cached` | `uint64` | cached memory |
| `writeBack` | `uint64` | write-back pages |
| `dirty` | `uint64` | dirty pages |
| `writeBackTmp` | `uint64` | temporary write-back pages |
| `shared` | `uint64` | shared memory |
| `slab` | `uint64` | slab memory |
| `slabReclaimable` | `uint64` | reclaimable slab |
| `slabUnreclaimable` | `uint64` | unreclaimable slab |
| `pageTables` | `uint64` | page table usage |
| `swapCached` | `uint64` | cached swap |
| `commitLimit` | `uint64` | commit limit |
| `committedAs` | `uint64` | committed memory |
| `highTotal` | `uint64` | high memory total |
| `highFree` | `uint64` | high memory free |
| `lowTotal` | `uint64` | low memory total |
| `lowFree` | `uint64` | low memory free |
| `swapTotal` | `uint64` | swap total |
| `swapFree` | `uint64` | swap free |
| `mapped` | `uint64` | mapped memory |
| `vmAllocTotal` | `uint64` | VM allocated total |
| `vmAllocUsed` | `uint64` | VM allocated used |
| `vmAllocChunk` | `uint64` | VM allocation chunk |
| `hugePagesTotal` | `uint64` | huge pages total |
| `hugePagesFree` | `uint64` | huge pages free |
| `hugePagesReserved` | `uint64` | huge pages reserved |
| `hugePagesSurplus` | `uint64` | huge pages surplus |
| `hugePageSize` | `uint64` | huge page size |
| `anonHugePages` | `uint64` | anonymous huge pages |

### `monitor.SwapMemory`

| Field | Type | Meaning |
| --- | --- | --- |
| `total` | `uint64` | total swap |
| `used` | `uint64` | used swap |
| `free` | `uint64` | free swap |
| `usedPercent` | `float64` | swap usage percentage |
| `swapIn` | `uint64` | swap-in count |
| `swapOut` | `uint64` | swap-out count |
| `pageIn` | `uint64` | page-in count |
| `pageOut` | `uint64` | page-out count |
| `pageFault` | `uint64` | page faults |
| `pageMajorFault` | `uint64` | major page faults |

### `monitor.DiskSummary`

| Field | Type | Meaning |
| --- | --- | --- |
| `total` | `uint64` | total disk size across counted partitions |
| `used` | `uint64` | used disk size |
| `usedPercent` | `float64` | disk usage percentage |
| `partitions` | `int` | partition count |

### `monitor.DiskInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `partitions` | `[]*monitor.PartitionInfo` | partition details |
| `ioCounters` | `map[string]*monitor.IOCounter` | per-device I/O counters |

### `monitor.PartitionInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `device` | `string` | device name |
| `mountPoint` | `string` | mount path |
| `fsType` | `string` | filesystem type |
| `options` | `[]string` | mount options |
| `total` | `uint64` | total size |
| `free` | `uint64` | free size |
| `used` | `uint64` | used size |
| `usedPercent` | `float64` | usage percentage |
| `iNodesTotal` | `uint64` | total inodes |
| `iNodesUsed` | `uint64` | used inodes |
| `iNodesFree` | `uint64` | free inodes |
| `iNodesUsedPercent` | `float64` | inode usage percentage |

### `monitor.IOCounter`

| Field | Type | Meaning |
| --- | --- | --- |
| `readCount` | `uint64` | read operation count |
| `mergedReadCount` | `uint64` | merged read operation count |
| `writeCount` | `uint64` | write operation count |
| `mergedWriteCount` | `uint64` | merged write operation count |
| `readBytes` | `uint64` | bytes read |
| `writeBytes` | `uint64` | bytes written |
| `readTime` | `uint64` | read time |
| `writeTime` | `uint64` | write time |
| `iopsInProgress` | `uint64` | I/O operations in progress |
| `ioTime` | `uint64` | total I/O time |
| `weightedIo` | `uint64` | weighted I/O time |
| `name` | `string` | device name |
| `serialNumber` | `string` | device serial number |
| `label` | `string` | device label |

### `monitor.NetworkSummary`

| Field | Type | Meaning |
| --- | --- | --- |
| `interfaces` | `int` | interface count |
| `bytesSent` | `uint64` | total bytes sent |
| `bytesRecv` | `uint64` | total bytes received |
| `packetsSent` | `uint64` | total packets sent |
| `packetsRecv` | `uint64` | total packets received |

### `monitor.NetworkInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `interfaces` | `[]*monitor.InterfaceInfo` | interface metadata |
| `ioCounters` | `map[string]*monitor.NetIOCounter` | per-interface counters |

### `monitor.InterfaceInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `index` | `int` | interface index |
| `mtu` | `int` | MTU |
| `name` | `string` | interface name |
| `hardwareAddr` | `string` | MAC address |
| `flags` | `[]string` | interface flags |
| `addrs` | `[]string` | bound addresses |

### `monitor.NetIOCounter`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | interface name |
| `bytesSent` | `uint64` | bytes sent |
| `bytesRecv` | `uint64` | bytes received |
| `packetsSent` | `uint64` | packets sent |
| `packetsRecv` | `uint64` | packets received |
| `errorsIn` | `uint64` | inbound errors |
| `errorsOut` | `uint64` | outbound errors |
| `droppedIn` | `uint64` | inbound drops |
| `droppedOut` | `uint64` | outbound drops |
| `fifoIn` | `uint64` | inbound FIFO count |
| `fifoOut` | `uint64` | outbound FIFO count |

### `monitor.ProcessSummary`

| Field | Type | Meaning |
| --- | --- | --- |
| `pid` | `int32` | process ID |
| `name` | `string` | process name |
| `cpuPercent` | `float64` | CPU usage percent |
| `memoryPercent` | `float32` | memory usage percent |

### `monitor.ProcessInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `pid` | `int32` | process ID |
| `parentPid` | `int32` | parent process ID |
| `name` | `string` | process name |
| `exe` | `string` | executable path |
| `commandLine` | `string` | full command line |
| `cwd` | `string` | working directory |
| `status` | `string` | process status |
| `username` | `string` | owner username |
| `createTime` | `int64` | creation timestamp |
| `numThreads` | `int32` | thread count |
| `numFds` | `int32` | open file-descriptor count |
| `cpuPercent` | `float64` | CPU usage percent |
| `memoryPercent` | `float32` | memory usage percent |
| `memoryRss` | `uint64` | RSS memory |
| `memoryVms` | `uint64` | virtual memory size |
| `memorySwap` | `uint64` | swap usage |

### `monitor.LoadInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `load1` | `float64` | 1-minute load average |
| `load5` | `float64` | 5-minute load average |
| `load15` | `float64` | 15-minute load average |

### `monitor.BuildInfo`

| Field | Type | Meaning |
| --- | --- | --- |
| `vefVersion` | `string` | framework version |
| `appVersion` | `string` | application version |
| `buildTime` | `string` | build time |
| `gitCommit` | `string` | git commit |

## Minimal Request Example

```json
{
  "resource": "sys/monitor",
  "action": "get_overview",
  "version": "v1"
}
```

## Practical Use

- admin or ops dashboards
- health and diagnostics surfaces
- internal tooling
- build metadata exposure

## Next Step

Read [CLI Tools](../advanced/cli-tools) if you want `generate-build-info` to populate richer build metadata.
