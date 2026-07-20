---
sidebar_position: 9
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

The monitor module registers the `sys/monitor` RPC resource, mounted under
`/api` with the standard envelope (`resource`, `action`, `version`,
`params`, `meta`). No operation is public and none declares a dedicated
permission token: every action inherits the API engine's default Bearer
authentication.

Every action sets a custom per-operation rate limit of `max 60`. The window
length is not overridden, so it inherits `vef.api.rate_limit.period`
(default `5m`); the limiter counts per operation + client IP + principal,
in process memory on each node.

None of the actions define framework-level input parameters: `params` is
ignored and may be omitted entirely.

| Action | Access | Rate limit | Input | Output |
| --- | --- | --- | --- | --- |
| `get_overview` | Bearer auth | `max 60` | none | `monitor.SystemOverview` |
| `get_cpu` | Bearer auth | `max 60` | none | `monitor.CPUInfo` |
| `get_memory` | Bearer auth | `max 60` | none | `monitor.MemoryInfo` |
| `get_disk` | Bearer auth | `max 60` | none | `monitor.DiskInfo` |
| `get_network` | Bearer auth | `max 60` | none | `monitor.NetworkInfo` |
| `get_host` | Bearer auth | `max 60` | none | `monitor.HostInfo` |
| `get_process` | Bearer auth | `max 60` | none | `monitor.ProcessInfo` |
| `get_load` | Bearer auth | `max 60` | none | `monitor.LoadInfo` |
| `get_build_info` | Bearer auth | `max 60` | none | `monitor.BuildInfo` |
| `get_event_streams` | Bearer auth | `max 60` | none | `monitor.EventStreamsInfo` |
| `get_integration_stats` | Bearer auth | `max 60` | none | `monitor.IntegrationStatsInfo` |

Behavior visible in source:

- `get_overview` is best-effort and never fails as a whole: a sub-probe that
  errors is logged and its overview field is left `null`, so one broken
  collector does not mask the rest.
- `get_cpu` and `get_process` are served from the background sample cache and
  return the monitor-not-ready business error (`monitor.ErrNotReady`) until
  the first sample lands.
- `get_memory`, `get_disk`, `get_network`, `get_host`, and `get_load` read
  live probes; a probe failure maps to `monitor.ErrCollectionFailed`.
- `get_build_info` cannot fail: the service always holds a non-nil build-info
  object (see [Build Info Behavior](#build-info-behavior)).
- `get_event_streams` is gated by the optional `event.StreamInspector`
  dependency. A nil inspector (the redis_stream transport is off) still
  returns `200 OK` with `enabled: false` and an empty `streams` list instead
  of failing; an inspector read error maps to `monitor.ErrCollectionFailed`.
- `get_integration_stats` mirrors the same degradation over the optional
  `integration.StatsInspector` (nil when the integration module is off):
  `enabled: false` with an empty `stats` list. Reading the in-memory snapshot
  itself cannot fail.
- business errors ride the standard result envelope: the HTTP status stays
  `200` and the failure is carried by the body `code`.

## Error API

| API | Meaning |
| --- | --- |
| `monitor.ErrNotReady` / `ErrCodeNotReady` (`2100`) | sample-backed data such as CPU or process metrics is not ready yet |
| `monitor.ErrCollectionFailed` / `ErrCodeCollectionFailed` (`2101`) | a monitor probe failed while collecting runtime data |

## Default Sampling Configuration

When no explicit monitor config is supplied, the module uses:

| Setting | Default |
| --- | --- |
| `vef.monitor.sample_interval` | `10s` |
| `vef.monitor.sample_duration` | `2s` |

These settings drive the background sampler behind `get_cpu` and
`get_process`: a sample is taken immediately at startup and then once per
sample interval, and each sample measures utilization over one
sample-duration window. Until the first sample completes (roughly the first
window after startup), both actions answer with `monitor.ErrNotReady`.

## Build Info Behavior

The monitor module decorates build info so that `vefVersion` is always present, even when the application does not provide a complete build metadata object.

Fallback behavior:

| Field | Fallback value when app does not supply build info |
| --- | --- |
| `appVersion` | `unknown` |
| `buildTime` | `unknown` |
| `gitCommit` | `unknown` |
| `vefVersion` | current framework version |

## Responses by Action

Field names below are the JSON wire names (the Go structs' json tags). Byte
quantities are plain byte counts, percentages range `0`–`100`, and counters
are cumulative since boot unless noted otherwise. Fields a platform does not
expose are reported as `0` or empty.

### `get_overview` — `monitor.SystemOverview`

One combined snapshot assembled from every probe. Each field is `null` when
its probe failed; `build` is always present.

| Field | Type | Description |
| --- | --- | --- |
| `host` | `*monitor.HostSummary` | condensed host information |
| `cpu` | `*monitor.CPUSummary` | condensed CPU information |
| `memory` | `*monitor.MemorySummary` | condensed memory usage |
| `disk` | `*monitor.DiskSummary` | condensed disk usage |
| `network` | `*monitor.NetworkSummary` | condensed network activity |
| `process` | `*monitor.ProcessSummary` | condensed current process metrics |
| `load` | `*monitor.LoadInfo` | load averages (same shape as `get_load`) |
| `build` | `*monitor.BuildInfo` | build metadata (same shape as `get_build_info`) |

#### `monitor.HostSummary`

| Field | Type | Description |
| --- | --- | --- |
| `hostname` | `string` | host name |
| `os` | `string` | operating system |
| `platform` | `string` | platform name |
| `platformVersion` | `string` | platform version |
| `kernelVersion` | `string` | kernel version |
| `kernelArch` | `string` | kernel architecture |
| `uptime` | `uint64` | host uptime in seconds |

#### `monitor.CPUSummary`

| Field | Type | Description |
| --- | --- | --- |
| `physicalCores` | `int` | number of physical cores (host topology) |
| `logicalCores` | `int` | number of logical cores (host topology) |
| `usagePercent` | `float64` | aggregated CPU usage percent over the last sampling window, normalized by `effectiveCores` |
| `effectiveCores` | `float64` | the capacity used to normalize utilization: inside a container this is the cgroup CPU quota (v1 and v2 supported), falling back to `logicalCores` when constrained usage cannot be sampled coherently |

#### `monitor.MemorySummary`

| Field | Type | Description |
| --- | --- | --- |
| `total` | `uint64` | total memory in bytes |
| `used` | `uint64` | used memory in bytes |
| `usedPercent` | `float64` | memory usage percentage |

The monitor is container-aware: when the process runs under a
cgroup (v2 or v1) that actually limits memory, the headline figures (`total`,
`used`, `usedPercent`, and `VirtualMemory`'s available/free) reflect the
cgroup limit and the cgroup's own usage instead of host-wide numbers — a
512 MiB container on a 64 GiB host reports against 512 MiB. Without a limit,
host-wide figures are reported unchanged.

#### `monitor.DiskSummary`

| Field | Type | Description |
| --- | --- | --- |
| `total` | `uint64` | total size of the root filesystem in bytes |
| `used` | `uint64` | used size of the root filesystem in bytes |
| `usedPercent` | `float64` | root filesystem usage percentage |
| `partitions` | `int` | always `1` (the summary covers a single filesystem) |

The overview's disk summary reports **the filesystem that bounds
the process's root path** rather than summing every mounted partition — remote
mounts, disk images, and sibling volumes do not inflate host capacity, and
there is no `vef.monitor.excluded_mounts` config (nothing is summed, so
nothing needs excluding). The raw mount inventory remains available through
`DiskInfo.partitions`.

#### `monitor.NetworkSummary`

| Field | Type | Description |
| --- | --- | --- |
| `interfaces` | `int` | interface count |
| `bytesSent` | `uint64` | total bytes sent, summed across interfaces |
| `bytesRecv` | `uint64` | total bytes received, summed across interfaces |
| `packetsSent` | `uint64` | total packets sent, summed across interfaces |
| `packetsRecv` | `uint64` | total packets received, summed across interfaces |

#### `monitor.ProcessSummary`

| Field | Type | Description |
| --- | --- | --- |
| `pid` | `int32` | process ID |
| `name` | `string` | process name |
| `cpuPercent` | `float64` | process CPU usage percent over the last sampling window; expressed against one CPU, so it can exceed `100` on multi-core hosts |
| `memoryPercent` | `float32` | share of total host RAM used by the process, percent |

### `get_cpu` — `monitor.CPUInfo`

Served from the background sample cache: refreshed once per sample interval
(default `10s`), each refresh measuring one sample-duration window (default
`2s`). Inventory fields (`modelName`, `vendorId`, `family`, `model`,
`stepping`, `microcode`, `mhz`, `cacheSize`) describe the first CPU package.

| Field | Type | Description |
| --- | --- | --- |
| `physicalCores` | `int` | number of physical cores (host topology) |
| `logicalCores` | `int` | number of logical cores (host topology) |
| `modelName` | `string` | CPU model name |
| `mhz` | `float64` | nominal clock frequency in MHz |
| `cacheSize` | `int32` | cache size in KB |
| `usagePercent` | `[]float64` | per-core busy percentage over the sampling window, one entry per logical core; `null` inside a CPU-limited container (the cgroup measurement replaces the per-core sample) |
| `totalPercent` | `float64` | aggregate usage percent: the mean of the per-core sample, or — inside a CPU-limited container — the share of the cgroup capacity consumed over the window, capped at `100` |
| `vendorId` | `string` | vendor identifier |
| `family` | `string` | CPU family |
| `model` | `string` | CPU model |
| `stepping` | `int32` | CPU stepping |
| `microcode` | `string` | microcode version |
| `effectiveCores` | `float64` | capacity used to normalize utilization; see `CPUSummary.effectiveCores` |

### `get_memory` — `monitor.MemoryInfo`

Read live on every call. The container-aware headline behavior described
under `MemorySummary` applies to `virtual` as well.

| Field | Type | Description |
| --- | --- | --- |
| `virtual` | `*monitor.VirtualMemory` | physical or virtual memory detail |
| `swap` | `*monitor.SwapMemory` | swap detail; `null` when the swap probe fails |

#### `monitor.VirtualMemory`

All fields are byte quantities except `usedPercent` (percent) and the
huge-page counters: `hugePagesTotal`, `hugePagesFree`, `hugePagesReserved`,
and `hugePagesSurplus` are page counts, while `hugePageSize` and
`anonHugePages` are bytes. Detail fields keep their host meaning even inside
a memory-limited container.

| Field | Type | Description |
| --- | --- | --- |
| `total` | `uint64` | total memory |
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
| `hugePagesTotal` | `uint64` | huge pages total (count) |
| `hugePagesFree` | `uint64` | huge pages free (count) |
| `hugePagesReserved` | `uint64` | huge pages reserved (count) |
| `hugePagesSurplus` | `uint64` | huge pages surplus (count) |
| `hugePageSize` | `uint64` | huge page size in bytes |
| `anonHugePages` | `uint64` | anonymous huge pages in bytes |

#### `monitor.SwapMemory`

`total`, `used`, and `free` are bytes. `swapIn`, `swapOut`, `pageIn`, and
`pageOut` are cumulative byte volumes converted from kernel page counters;
`pageFault` and `pageMajorFault` are cumulative event counts.

| Field | Type | Description |
| --- | --- | --- |
| `total` | `uint64` | total swap |
| `used` | `uint64` | used swap |
| `free` | `uint64` | free swap |
| `usedPercent` | `float64` | swap usage percentage |
| `swapIn` | `uint64` | swapped-in volume |
| `swapOut` | `uint64` | swapped-out volume |
| `pageIn` | `uint64` | paged-in volume |
| `pageOut` | `uint64` | paged-out volume |
| `pageFault` | `uint64` | page faults |
| `pageMajorFault` | `uint64` | major page faults |

### `get_disk` — `monitor.DiskInfo`

Read live on every call. A partition whose usage probe fails is skipped from
`partitions`; `ioCounters` is `null` when the I/O counter probe fails.

| Field | Type | Description |
| --- | --- | --- |
| `partitions` | `[]*monitor.PartitionInfo` | per-mount partition details |
| `ioCounters` | `map[string]*monitor.IOCounter` | per-device I/O counters, keyed by device name |

#### `monitor.PartitionInfo`

| Field | Type | Description |
| --- | --- | --- |
| `device` | `string` | device name |
| `mountPoint` | `string` | mount path |
| `fsType` | `string` | filesystem type |
| `options` | `[]string` | mount options |
| `total` | `uint64` | total size in bytes |
| `free` | `uint64` | free size in bytes |
| `used` | `uint64` | used size in bytes |
| `usedPercent` | `float64` | usage percentage |
| `iNodesTotal` | `uint64` | total inodes |
| `iNodesUsed` | `uint64` | used inodes |
| `iNodesFree` | `uint64` | free inodes |
| `iNodesUsedPercent` | `float64` | inode usage percentage |

#### `monitor.IOCounter`

Counters are cumulative since boot; `readTime`, `writeTime`, `ioTime`, and
`weightedIo` are milliseconds.

| Field | Type | Description |
| --- | --- | --- |
| `readCount` | `uint64` | read operation count |
| `mergedReadCount` | `uint64` | merged read operation count |
| `writeCount` | `uint64` | write operation count |
| `mergedWriteCount` | `uint64` | merged write operation count |
| `readBytes` | `uint64` | bytes read |
| `writeBytes` | `uint64` | bytes written |
| `readTime` | `uint64` | time spent reading |
| `writeTime` | `uint64` | time spent writing |
| `iopsInProgress` | `uint64` | I/O operations in progress |
| `ioTime` | `uint64` | total time spent on I/O |
| `weightedIo` | `uint64` | weighted I/O time |
| `name` | `string` | device name |
| `serialNumber` | `string` | device serial number |
| `label` | `string` | device label |

### `get_network` — `monitor.NetworkInfo`

Read live on every call.

| Field | Type | Description |
| --- | --- | --- |
| `interfaces` | `[]*monitor.InterfaceInfo` | interface metadata |
| `ioCounters` | `map[string]*monitor.NetIOCounter` | per-interface counters, keyed by interface name |

#### `monitor.InterfaceInfo`

| Field | Type | Description |
| --- | --- | --- |
| `index` | `int` | interface index |
| `mtu` | `int` | MTU |
| `name` | `string` | interface name |
| `hardwareAddr` | `string` | MAC address |
| `flags` | `[]string` | interface flags |
| `addrs` | `[]string` | bound addresses |

#### `monitor.NetIOCounter`

Counters are cumulative since boot, per interface.

| Field | Type | Description |
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

### `get_host` — `monitor.HostInfo`

Static host metadata, read live on every call.

| Field | Type | Description |
| --- | --- | --- |
| `hostname` | `string` | host name |
| `uptime` | `uint64` | host uptime in seconds |
| `bootTime` | `uint64` | boot time as a Unix timestamp (seconds) |
| `processes` | `uint64` | number of processes on the host |
| `os` | `string` | operating system |
| `platform` | `string` | platform name |
| `platformFamily` | `string` | platform family |
| `platformVersion` | `string` | platform version |
| `kernelVersion` | `string` | kernel version |
| `kernelArch` | `string` | kernel architecture |
| `virtualizationSystem` | `string` | virtualization system |
| `virtualizationRole` | `string` | virtualization role |
| `hostId` | `string` | host identifier |

### `get_process` — `monitor.ProcessInfo`

Describes the application's own process. Served from the background sample
cache on the same cadence as `get_cpu`.

| Field | Type | Description |
| --- | --- | --- |
| `pid` | `int32` | process ID |
| `parentPid` | `int32` | parent process ID |
| `name` | `string` | process name |
| `exe` | `string` | executable path |
| `commandLine` | `string` | full command line |
| `cwd` | `string` | working directory |
| `status` | `string` | process status |
| `username` | `string` | owner username |
| `createTime` | `int64` | process creation time, milliseconds since the Unix epoch (UTC) |
| `numThreads` | `int32` | thread count |
| `numFds` | `int32` | open file-descriptor count |
| `cpuPercent` | `float64` | process CPU usage percent over the sampling window; expressed against one CPU, so it can exceed `100` on multi-core hosts |
| `memoryPercent` | `float32` | share of total host RAM used by the process, percent |
| `memoryRss` | `uint64` | resident set size in bytes |
| `memoryVms` | `uint64` | virtual memory size in bytes |
| `memorySwap` | `uint64` | swap usage in bytes |

### `get_load` — `monitor.LoadInfo`

Read live on every call.

| Field | Type | Description |
| --- | --- | --- |
| `load1` | `float64` | 1-minute load average |
| `load5` | `float64` | 5-minute load average |
| `load15` | `float64` | 15-minute load average |

### `get_build_info` — `monitor.BuildInfo`

Build metadata only; see [Build Info Behavior](#build-info-behavior) for the
fallback values.

| Field | Type | Description |
| --- | --- | --- |
| `vefVersion` | `string` | framework version, always stamped by the module |
| `appVersion` | `string` | application version |
| `buildTime` | `string` | build time |
| `gitCommit` | `string` | git commit |

### `get_event_streams` — `monitor.EventStreamsInfo`

Reports cross-process event stream and consumer-group state through the
optional `event.StreamInspector` (provided by the redis_stream transport).

| Field | Type | Description |
| --- | --- | --- |
| `enabled` | `bool` | whether an `event.StreamInspector` is available (the redis_stream transport is on); `false` means the report is an empty degradation, not an error |
| `streams` | `[]event.StreamInfo` | one entry per transport stream; empty when `enabled` is `false` |

#### `event.StreamInfo`

| Field | Type | Description |
| --- | --- | --- |
| `stream` | `string` | full transport-level stream key (prefix + event type) |
| `length` | `int64` | current number of entries in the stream (post-trim) |
| `groups` | `[]event.StreamGroupInfo` | consumer groups attached to the stream |

#### `event.StreamGroupInfo`

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | consumer group name (the subscription's `WithGroup` value or its derived default) |
| `consumers` | `int64` | number of consumer records registered in the group, including historical consumers of restarted processes |
| `pending` | `int64` | number of delivered-but-unacknowledged entries |
| `lag` | `int64` | number of stream entries not yet delivered to this group (approximate after trimming; zero on server versions that do not report lag) |
| `lastDeliveredId` | `string` | stream ID of the last entry delivered to the group |

A group with growing `lag` and only idle consumers is an orphan candidate — a subscriber that was removed or renamed without decommissioning its consumer group. See the [Event Bus](./event-bus) page for the transport-level detail.

### `get_integration_stats` — `monitor.IntegrationStatsInfo`

Reports per-node integration invocation statistics through the
optional `integration.StatsInspector`. Numbers are in-memory counters held
since process start — the invocation log is the durable record.

| Field | Type | Description |
| --- | --- | --- |
| `enabled` | `bool` | whether an `integration.StatsInspector` is available (the integration module is on); `false` means the report is an empty degradation, not an error |
| `stats` | `[]integration.InvocationStats` | one entry per `(system, contract, direction)` tuple observed since process start, ordered by system, contract, then direction; empty when `enabled` is `false` |

#### `integration.InvocationStats`

| Field | Type | Description |
| --- | --- | --- |
| `system` | `string` | system code that served (or rejected) the invocation |
| `contract` | `string` | invoked contract code; empty for inbound deliveries rejected by verification — the contract code is unvalidated caller input at rejection time |
| `direction` | `string` | `outbound` or `inbound` |
| `calls` | `int64` | total invocations observed |
| `successes` | `int64` | invocations that completed successfully |
| `failures` | `map[string]int64` | failure counts keyed by [failure kind](../integration/overview#failure-vocabulary) (`input_invalid`, `output_invalid`, `upstream`, `transport`, `timeout`, `canceled`, `script`, `config`, `auth`, `handler`); omitted when empty |
| `avgDurationMs` | `int64` | average invocation duration in milliseconds |
| `maxDurationMs` | `int64` | maximum invocation duration in milliseconds |
| `lastError` | `string` | most recent failure message; omitted when no failure occurred |
| `lastErrorAt` | timestamp | time of the most recent failure; omitted when no failure occurred |

See [Integration Engine](../integration/overview#statistics) for how these
counters are recorded.

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
