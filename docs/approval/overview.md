---
sidebar_position: 1
slug: /approval
---

# Approval Module

The `approval` module provides a complete workflow engine for building approval-based business processes. It supports visual flow design (React Flow compatible), multi-level approval chains, conditional branching, parallel approval, delegation, rollbacks, and transactional event publishing.

This category splits the module across five pages: this overview (enabling, wiring, configuration), [RPC Resources](./resources.md), [Flow Design](./flow-design.md), [Instance Runtime](./runtime.md), and [Events & Integration](./integration.md).

## Enabling the Module

Approval is an optional feature module. It is intentionally absent from the
default `vef.Run(...)` boot graph, so applications that do not need workflow
support do not register its API resources, CQRS handlers, engine, binding
listener, or timeout scanners.

Enable it explicitly:

```go
vef.Run(
    vef.ApprovalModule,
    app.Module,
)
```

## Event Routing Prerequisite

Approval publishes `approval.*` events with `event.WithTx`, and its binding
listener subscribes to those events. The host application must route
`approval.*` to a transactional transport with a subscribable sink, such as an
outbox route whose sink is Redis Streams:

```toml
[vef.event]
default_transport = "memory"

[[vef.event.routing]]
pattern    = "approval.*"
transports = ["outbox", "redis_stream"]
```

The route must list the configured outbox `sink` (here `redis_stream`; the
in-process `memory` transport also qualifies) alongside `outbox` itself: the
outbox transport is publish-only, so subscribers — including the module's own
binding listener — attach to the sink transport. Approval's binding listener
and outbox publisher both assert routing at boot via `event.RouteInspector`, so
a misconfigured route fails the application instead of degrading silently. See
[Event Bus](../infrastructure/event-bus.md) for transports, routing semantics,
and the outbox relay.

`InstanceBindingFailedEvent` is the exception to the transactional-route
startup check: it is emitted by the asynchronous binding listener after the
approval transaction has already committed. `InstanceCompletedEvent` has the
strictest route requirement because the binding listener subscribes to it; the
route must include a subscribable sink transport such as `memory` or
`redis_stream` alongside the transactional outbox route.

## Architecture Overview

```
Flow Category → Flow → Flow Version → Nodes + Edges
                                        ↓
                                    Instance → Tasks → Action Logs
```

| Concept | Table | Description |
| --- | --- | --- |
| Flow Category | `apv_flow_category` | Hierarchical grouping of flows |
| Flow | `apv_flow` | A workflow definition (e.g., "Leave Request") |
| Flow Version | `apv_flow_version` | Versioned snapshot with nodes, edges, and form schema |
| Flow Node | `apv_flow_node` | A step in the workflow (approval, handle, condition, CC) |
| Flow Edge | `apv_flow_edge` | Directed connection between nodes |
| Instance | `apv_instance` | A running instance of a flow |
| Task | `apv_task` | An individual approval/handle task assigned to a user |
| Action Log | `apv_action_log` | Audit trail of all actions |

## Configuration

```toml
[vef.approval]
auto_migrate              = true
timeout_scan_interval     = "1m"
pre_warning_scan_interval = "5m"
cleanup_scan_interval     = "24h"
delegation_max_depth      = 10
form_snapshot_retention   = "2160h"  # 90 days
urge_record_retention     = "720h"   # 30 days
cc_record_retention       = "2160h"  # 90 days
```

`auto_migrate` is a plain boolean switch and is not set by
`ApprovalConfig.ApplyDefaults()`: enable it explicitly when the app should run
approval DDL on startup. `cc_record_retention` only prunes CC records that have
already been read.

> The outbox-related fields that previously lived under `[vef.approval]` (`outbox_relay_interval`, `outbox_max_retries`, `outbox_batch_size`) moved to `[vef.event.transports.outbox]` in v0.21. The approval module now publishes through the framework-wide outbox transport — see [Event Bus](../infrastructure/event-bus.md).

See [Configuration Reference](../reference/configuration-reference.md) for details.

## Binding Modes

| Mode | Constant | Wire value | Description |
| --- | --- | --- | --- |
| Standalone | `BindingStandalone` | `standalone` | Form data stored in the approval module's own tables |
| Business | `BindingBusiness` | `business` | Links to an existing business data table |

Business binding connects the approval flow to your domain tables via `BusinessTable`, `BusinessPKField`, and `BusinessStatusField`, plus the optional `BusinessInstanceIDField`, `BusinessStartedAtField`, and `BusinessFinishedAtField` linkage columns (see [Business Write-Back Linkage Matrix](./integration.md#business-write-back-linkage-matrix)).

---

Next: [RPC Resources](./resources.md) for the API surface, or [Flow Design](./flow-design.md) for node types and the designer wire shapes.
