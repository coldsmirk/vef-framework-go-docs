---
sidebar_position: 1
---

# Integration Engine

The integration module (v0.39) is a config- and script-driven engine for
talking to external systems — HIS/ERP vendors, partner gateways, provincial
platforms — without hard-coding any vendor's wire format into business code.
It is optional: enable it by passing `vef.IntegrationModule` to `vef.Run`.

```go
vef.Run(
    vef.IntegrationModule,
    vef.Module("app", ...),
)
```

## The Canonical-Model Idea

Business code programs against **contracts** — standard input/output models
you define once. Vendor differences (URL shapes, field names, envelopes,
signatures, code values) live in per-system **adapter scripts**, editable at
runtime through management APIs. Swapping a vendor means writing a new
adapter, not touching business code.

Four definition tables drive the engine:

| Definition | Table | Meaning |
| --- | --- | --- |
| `integration.Contract` | `itg_contract` | One standard operation: code, name, host-owned `labels`, optional JSON Schema (draft 2020-12) for input and output |
| `integration.System` | `itg_system` | One external system instance: base URL, outbound/inbound auth, optional outbound envelope, optional direct database, params, timeout, retry policy |
| `integration.Adapter` | `itg_adapter` | Binds one system to one contract in one `direction` (`outbound` / `inbound`) with a translation script |
| `integration.Route` | `itg_route` | Maps a route key (tenant, branch, hospital area) to the system serving a contract |

Two flows share those definitions:

- **Outbound** — business code calls `integration.Invoker.Invoke`; the engine
  resolves the target system, runs the outbound adapter script, and returns
  the schema-validated standard model. See [Outbound Calls](./outbound).
- **Inbound** — the external system calls
  `POST /integration/inbound/:systemCode/:contractCode`; the gateway verifies
  the caller, runs the inbound adapter script, which dispatches the standard
  input to your registered `integration.InboundHandler` and shapes the reply
  the vendor expects. See [Inbound Delivery](./inbound).

Value-level differences (gender codes, status enums) are translated by
per-system [Code Maps](./code-maps).

## Contract

```go
type Contract struct {
    Code         string            // unique code business code invokes
    Name         string
    Description  *string
    Labels       map[string]string // host-owned selection metadata, equality-filterable
    InputSchema  json.RawMessage   // JSON Schema; empty skips input validation
    OutputSchema json.RawMessage   // JSON Schema; empty skips output validation
    IsEnabled    bool
}
```

Schemas must be self-contained JSON Schema documents (draft 2020-12); they
are compiled at save time (`ErrInvalidSchema`) so a broken contract never
reaches an invocation. Input is validated before the adapter script runs and
the script's return value is validated after, in both flows — an adapter can
never hand business code an out-of-contract model.

`Labels` are stored and filtered by the engine but never interpreted. The
shared `orm.ValidateLabels` rule applies (the same one behind approval flow
labels): keys are alphanumeric with inner `-`/`_` (no dots — they would read
as JSON-path nesting in the equality filter), at most 63 characters; values
are at most 256 characters, and an empty value is legal (a presence flag).
Violations fail with `ErrInvalidLabel`.

## System

```go
type System struct {
    Code             string
    Name             string
    BaseURL          string                  // enables the scoped http library
    OutboundAuth     *OutboundAuthConfig     // nil sends requests unauthenticated
    OutboundEnvelope *OutboundEnvelopeConfig // nil sends adapter requests untouched
    InboundAuth      *InboundAuthConfig      // nil refuses inbound delivery entirely
    DataSource       *DataSourceConfig       // enables the scoped sql library
    Params           map[string]string       // non-sensitive values, visible as system.params
    TimeoutMs        int                     // per-HTTP-call bound; zero = framework default
    Retry            *RetryPolicy            // rides the httpx default retry policy
    IsEnabled        bool
}
```

- A system with `BaseURL` gives its adapter scripts the scoped `http` client;
  one with `DataSource` gives them the scoped `sql` library (read-only unless
  `dataSource.mode = "read_write"`); a system may carry both.
- `Retry` retries idempotent methods on transport errors and 429/502/503/504
  responses: `maxAttempts` (total attempts, first call included),
  `initialBackoffMs`, `maxBackoffMs`.
- `TimeoutMs` bounds each HTTP call; the adapter's `timeoutMs` bounds the
  whole script run (a different axis, inherited from
  `vef.integration.run_timeout` when zero).

### Secrets at rest

Sensitive auth parameter values and the data-source password are encrypted
with the key in `vef.integration.secret_key` (base64; AES-GCM by default,
SM4-GCM via `vef.integration.secret_algorithm = "sm4"`). Management API
responses always mask them as `"******"` (`integration.MaskedSecret`);
an update that submits the placeholder keeps the stored value unchanged.
Leaving `secret_key` unset stores secrets in plaintext and logs a startup
warning. Values sealed with one algorithm are not readable under the other —
switching requires re-entering stored secrets.

## Adapter

```go
type Adapter struct {
    SystemID   string
    ContractID string
    Direction  Direction // "outbound" (default) or "inbound"
    Script     string    // compile-checked at save time
    TimeoutMs  int       // script run timeout; zero inherits vef.integration.run_timeout
    IsEnabled  bool
}
```

A system implements a contract with exactly one adapter per direction. The
script environments differ per direction and are documented in
[Outbound Calls](./outbound#adapter-script-environment) and
[Inbound Delivery](./inbound#adapter-script-environment).

## Route

```go
type Route struct {
    RouteKey   string // "" is the default route
    ContractID string // "" applies to every contract
    SystemID   string
    IsEnabled  bool
}
```

`integration.RouteResolver` resolves `(contract, routeKey)` to a system code.
The framework default reads `itg_route`: exact `(key, contract)` rules win
over contract-wildcard rules, and the empty key is the default route. Replace
it via `fx.Decorate` when routing lives elsewhere (a tenant registry, a
config center). A key matching no rule fails with `ErrRouteNotFound`.

Route health is inspectable at runtime through the `diagnose_routes`
operation (dangling adapters, disabled targets, uncovered contracts); see
[RPC Resources](./resources#integrationops).

## Configuration

```toml
[vef.integration]
auto_migrate = true          # run the integration DDL migration at startup
secret_key = "base64-key"    # encrypts sensitive values at rest
secret_algorithm = "aes"     # "aes" (AES-GCM, default) or "sm4" (SM4-GCM)
run_timeout = "30s"          # per script run, wire calls included
max_response_body = 8388608  # cap per HTTP response body read by scripts (8 MiB)

[vef.integration.log]
mode = "errors"              # "off" | "errors" (default) | "all"
capture_limit = 4096         # bytes per captured payload before truncation
mask_fields = ["idCard"]     # extra JSON field names masked in captures
retention = "720h"           # prune invocation logs; zero keeps forever

[vef.integration.inbound.rate_limit]
max = 120                    # deliveries per window per (system, client IP)
period = "1m"
```

Invocation logging writes `itg_invocation_log` rows carrying the failure
classification, timing, input/output captures, and the full HTTP wire trace
— masked (credential headers always; `mask_fields` additionally) and
truncated to `capture_limit`. An hourly sweep prunes rows older than
`retention` when set.

## Failure Vocabulary

`integration.FailureKind` is the single classification shared by invocation
logs, statistics, and API errors; an empty value means success.

| Kind | Meaning |
| --- | --- |
| `input_invalid` | input rejected by the contract's input schema before the script ran |
| `output_invalid` | script return value rejected by the output schema |
| `upstream` | failure the external system itself signaled (`errors.upstream(...)`) |
| `transport` | wire call never completed (connection refused, TLS failure) |
| `timeout` | invocation exceeded its run timeout |
| `canceled` | caller canceled the invocation |
| `script` | uncaught script exception or compile error — an adapter bug |
| `config` | auth scheme unregistered, credential undecryptable, and similar |
| `auth` | inbound delivery rejected by inbound auth verification |
| `handler` | inbound business handler returned an error after a successful dispatch |

## Statistics

The invoker implements `integration.StatsInspector`: per-node counters per
`(system, contract, direction)` tuple since process start — calls, successes,
failures by kind, average/max duration, last error. The monitor module reads
it through `sys/monitor.get_integration_stats`
(see [Built-in Resources](../reference/built-in-resources)). Each
`InvocationStats` entry carries `system`, `contract`, `direction`, `calls`,
`successes`, `failures` (map of `FailureKind` to count), `avgDurationMs`,
`maxDurationMs`, `lastError`, and `lastErrorAt`. Inbound deliveries rejected
by verification aggregate under an empty `contract` — the contract code is
unvalidated caller input at rejection time.

## Error Codes

Integration API errors use response codes `2600`–`2699` and ride HTTP 200
with the failure in the body code, except where noted.

| Code | Error | Meaning |
| --- | --- | --- |
| `2600` | `ErrContractNotFound` | contract lookup failed |
| `2601` | `ErrContractDisabled` | contract is disabled |
| `2602` | `ErrSystemNotFound` | system lookup failed |
| `2603` | `ErrSystemDisabled` | system is disabled |
| `2604` | `ErrAdapterNotFound` | no adapter binds the system to the contract in that direction |
| `2605` | `ErrAdapterDisabled` | adapter is disabled |
| `2606` | `ErrRouteNotFound` | route key matches no rule |
| `2607` | `ErrTargetAmbiguous` | both `WithSystem` and `WithRoute` were passed |
| `2608` | `ErrInputInvalid(detail)` | input rejected by the input schema |
| `2609` | `ErrOutputInvalid(detail)` | script return rejected by the output schema |
| `2610` | `ErrUpstreamFailed(message)` | failure signaled by the external system |
| `2611` | `ErrTransportFailed` | wire call never completed |
| `2612` | `ErrInvocationTimeout` | run timeout exceeded |
| `2613` | `ErrScriptFailed(detail)` | script threw or failed to compile |
| `2614` | `ErrUnknownAuthScheme(scheme)` | system references an unregistered auth scheme |
| `2615` | `ErrInvalidSchema(detail)` | contract schema rejected at save time |
| `2616` | `ErrInvalidScript(detail)` | adapter script rejected at save time |
| `2617` | `ErrInvalidAuthParams(detail)` | auth configuration refused by its scheme |
| `2618` | `ErrInvalidRouteRef` | route references a missing contract or system |
| `2619` | `ErrInvalidBaseURL` | system base URL is not an absolute URL |
| `2620` | `ErrInvalidDataSource(detail)` | system data source incomplete or credential unprocessable |
| `2621` | `ErrInvalidDirection` | adapter direction outside the known flows |
| `2622` | `ErrInboundAuthFailed` | inbound delivery failed verification (HTTP 401, deliberately uniform) |
| `2623` | `ErrInboundHandlerMissing` | inbound contract has no registered handler (HTTP 501) |
| `2624` | `ErrInvocationCanceled` | caller canceled the invocation |
| `2625` | `ErrInvalidEnvelope(detail)` | outbound envelope config rejected at save time |
| `2626` | `ErrInvalidLabel` | contract label key/value failed validation |
| `2627` | `ErrMissingCodeMap(codeSet)` | codes lookup against a code set with no enabled map |
| `2628` | `ErrUnmappedValue(codeSet, value)` | value unmapped under the reject policy |
| `2629` | `ErrInvalidCodeMap(detail)` | code map definition rejected at save time |

## Next Steps

- [Outbound Calls](./outbound) — the `Invoker`, adapter script environment, auth schemes, envelopes
- [Inbound Delivery](./inbound) — the HTTP gateway, verification schemes, business handlers
- [Code Maps](./code-maps) — per-system value translation and the host code set catalog
- [RPC Resources](./resources) — field-by-field reference of the management APIs
