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

Approval publishes `approval.*` events with `event.WithTx`, so every event
type except `approval.instance.binding_failed` must resolve to a
**transactional** transport. The module asserts this at boot via
`event.RouteInspector` — a misconfigured route fails the application instead
of degrading silently:

```toml
[vef.event]
default_transport = "memory"

[[vef.event.routing]]
pattern    = "approval.*"
transports = ["outbox", "redis_stream"]
```

The business projection does not consume lifecycle events (it
converges from the durable `apv_business_projection` table — see
[Business-State Projection](./integration.md#business-state-projection)), so
the module itself requires no subscribable sink in the route: a bare
`["outbox"]` route satisfies the startup check. Add the configured outbox
`sink` transport (`memory` single-node, `redis_stream` cross-node) alongside
`outbox` — as in the example above — whenever the **host** subscribes to
approval events via `approval.SubscribeInstance` or `approval.BindCommand`,
because the outbox transport is publish-only and subscribers attach to the
sink. See [Event Bus](../infrastructure/event-bus.md) for transports, routing
semantics, and the outbox relay.

`InstanceBindingFailedEvent` is the exception to the transactional-route
startup check: it is published by the eventual projection worker after the
approval transaction has already committed.

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
| Business Projection | `apv_business_projection` | Durable desired-state convergence for business-bound flows |

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

[vef.approval.business_binding]
consistency   = "synchronous"  # or "eventual"
scan_interval = "10s"          # eventual worker cadence
batch_size    = 100            # projections per scan
```

`auto_migrate` is a plain boolean switch and is not set by
`ApprovalConfig.ApplyDefaults()`: enable it explicitly when the app should run
approval DDL on startup. `cc_record_retention` only prunes CC records that have
already been read.

`business_binding` controls how approval state is projected onto bound
business tables: `synchronous` (default) writes the business row inside the
approval transaction, `eventual` commits the desired state and lets a
background worker converge the row (see
[Business-State Projection](./integration.md#business-state-projection)).
An out-of-enum `consistency` or a negative worker setting fails config
validation at startup (`config.ErrInvalidApprovalBindingConsistency` /
`ErrInvalidApprovalBusinessBindingWorkerConfig`).

> Outbox tuning (`relay_interval`, `max_retries`, `batch_size`) lives under `[vef.event.transports.outbox]`, not `[vef.approval]`. The approval module publishes through the framework-wide outbox transport — see [Event Bus](../infrastructure/event-bus.md).

See [Configuration Reference](../reference/configuration-reference.md) for details.

## Binding Modes

| Mode | Constant | Wire value | Description |
| --- | --- | --- | --- |
| Standalone | `BindingStandalone` | `standalone` | Form data stored in the approval module's own tables |
| Business | `BindingBusiness` | `business` | Links to an existing business data table |

Business binding connects the approval flow to your domain tables via a single `Flow.BusinessBinding` document (`approval.BusinessBindingConfig`): `tableName`, composite `keyColumns`, `statusColumn`, the mandatory `instanceIdColumn` CAS fence, optional `startedAtColumn` / `finishedAtColumn`, and an optional `statusMapping` (see [Business-State Projection](./integration.md#business-state-projection)). The binding is snapshotted onto each deployed flow version, and its state converges through the durable `apv_business_projection` table.

---

Next: [RPC Resources](./resources.md) for the API surface, or [Flow Design](./flow-design.md) for node types and the designer wire shapes.
