---
sidebar_position: 5
---

# RPC Resources

When `vef.IntegrationModule` is enabled, the framework registers the
management resources below. All of them are RPC resources mounted under
`/api`, using the standard envelope (`resource`, `action`, `version`,
`params`, `meta`) documented in [API](../building-apis/api.md). None of the
operations are public: callers must be authenticated, and every operation
declares the permission listed in its table.

Conventions used on this page:

- CRUD read operations declare search structs embedding `crud.Sortable` — a
  `meta` struct — so their filter fields decode from the request's `meta`
  object; `find_page` additionally reads `meta.page` and `meta.size`
  (`page.Pageable`), and `meta.sort` carries sort specs.
- `find_page` responds with `page.Page[T]`: `page`, `size`, `total`,
  `items`.
- Mutations decode from `params`. Fields marked required are enforced by
  validation; the rest are optional.
- All definition models carry the standard audited-model columns
  (`id`, `createdAt`, `createdBy`, `updatedAt`, `updatedBy`) in responses;
  they are omitted from the field tables below.

## `integration/contract`

Contract definitions. Schemas are compiled at save time so a broken contract
never reaches an invocation.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `integration.contract.query` | `ContractSearch` + pageable meta | `page.Page[Contract]` |
| `find_all` | `integration.contract.query` | `ContractSearch` | `Contract[]` |
| `create` | `integration.contract.create` | `ContractParams` | created `Contract` |
| `update` | `integration.contract.update` | `ContractParams` | updated `Contract` |
| `delete` | `integration.contract.delete` | primary-key params (`params.id`) | success |

`ContractSearch` (query filters):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `code` | `string` | contains | filter by contract code fragment |
| `name` | `string` | contains | filter by name fragment |
| `isEnabled` | `bool` | equals | filter by enablement; omit to match both |
| `labels` | `object` (string→string) | equality on every pair | host-driven label filter (business-side contract pickers select by labels) |

`ContractParams` (create/update):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | update only | primary key of the row to update |
| `code` | `string` | Yes | unique contract code business code invokes |
| `name` | `string` | Yes | display name |
| `description` | `string` | No | free-text description |
| `labels` | `object` (string→string) | No | host-owned selection metadata; keys must not contain dots, keys/values size-bounded (`ErrInvalidLabel`) |
| `inputSchema` | JSON Schema object | No | self-contained draft 2020-12 schema for the invocation input; empty skips input validation |
| `outputSchema` | JSON Schema object | No | self-contained schema for the adapter's return value; empty skips output validation |
| `isEnabled` | `bool` | No | disabled contracts refuse invocation (`ErrContractDisabled`) |

Deleting a contract still referenced by routes fails with the standard
foreign-key violation error (the route table's contract column carries the
empty-string wildcard sentinel, so this check is enforced by the resource).

## `integration/system`

External system definitions. Writes encrypt sensitive auth parameters and
the data source password; reads always mask them as `"******"`. Submitting
the mask back keeps the stored value unchanged.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `integration.system.query` | `SystemSearch` + pageable meta | `page.Page[System]` (masked) |
| `find_all` | `integration.system.query` | `SystemSearch` | `System[]` (masked) |
| `create` | `integration.system.create` | `SystemParams` | created `System` |
| `update` | `integration.system.update` | `SystemParams` | updated `System` |
| `delete` | `integration.system.delete` | primary-key params (`params.id`) | success |

Deleting a system — or removing/renaming its data source on update —
releases its data source registry entry.

`SystemSearch` (query filters):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `code` | `string` | contains | filter by system code fragment |
| `name` | `string` | contains | filter by name fragment |
| `isEnabled` | `bool` | equals | filter by enablement |

`SystemParams` (create/update):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | update only | primary key |
| `code` | `string` | Yes | unique system code |
| `name` | `string` | Yes | display name |
| `baseUrl` | `string` | No | absolute base URL; enables the scoped `http` library for this system's scripts. Validated at save time (`ErrInvalidBaseURL`) |
| `outboundAuth` | `OutboundAuthConfig` | No | outbound authentication (below); `null` sends requests unauthenticated |
| `outboundEnvelope` | `OutboundEnvelopeConfig` | No | system-level request/response wrap scripts (below); `null` passes adapter requests through untouched |
| `inboundAuth` | `InboundAuthConfig` | No | inbound verification (below); `null` refuses inbound delivery entirely |
| `dataSource` | `DataSourceConfig` | No | direct database connection (below); enables the scoped `sql` library |
| `params` | `object` (string→string) | No | non-sensitive system-specific values, exposed to scripts as `system.params` |
| `timeoutMs` | `int` | No | per-HTTP-call bound; zero applies the framework default |
| `retry` | `RetryPolicy` | No | outbound retry policy (below) |
| `isEnabled` | `bool` | No | disabled systems refuse both flows (`ErrSystemDisabled`; inbound denies uniformly as auth failure) |

`OutboundAuthConfig` / `InboundAuthConfig`:

| Field | Type | Description |
| --- | --- | --- |
| `scheme` | `string` | scheme name. Outbound: `none`, `http_basic`, `bearer`, `header`, `query`, `signature`, `script`, or custom. Inbound additionally: `ip` |
| `params` | `object` (string→string) | scheme parameters; values of scheme-declared sensitive parameters are stored encrypted and masked in responses |
| `script` | `string` | custom signing/verification body for the `script` scheme; runs in a zero-IO runtime |

See [Outbound Calls](./outbound#outbound-authentication-schemes) and
[Inbound Delivery](./inbound#inbound-authentication) for each scheme's
parameter reference.

`OutboundEnvelopeConfig`:

| Field | Type | Description |
| --- | --- | --- |
| `request` | `string` | wrap script: receives the adapter's request as `request` (`{ method, path, headers, query, body }`) and returns the request to put on the wire; omitted fields keep the adapter's values |
| `response` | `string` | unwrap script: receives the completed HTTP response as `response` (fetch Response shape); its return value is what the adapter's call yields |

At least one of the two scripts is required when the envelope is present, and
the system must have an HTTP transport (`ErrInvalidEnvelope`).

`DataSourceConfig`:

| Field | Type | Description |
| --- | --- | --- |
| `kind` | `string` | database kind (same vocabulary as `vef.data_sources.type`: `postgres`, `mysql`, `sqlite`, `sqlserver`, `oracle`) |
| `mode` | `string` | script write access: `read_only` (default; `sql.execute` throws) or `read_write` (enables `sql.execute`) |
| `host` | `string` | server host |
| `port` | `int` | server port |
| `user` | `string` | login user |
| `password` | `string` | login password — stored encrypted, masked in responses |
| `database` | `string` | database name |
| `schema` | `string` | schema name (where the kind supports it) |
| `path` | `string` | file path (sqlite) |
| `sslMode` | `string` | SSL mode (same vocabulary as `vef.data_sources.ssl_mode`) |
| `sslRootCert` | `string` | CA certificate path |

`RetryPolicy`:

| Field | Type | Description |
| --- | --- | --- |
| `maxAttempts` | `int` | total number of attempts, the first call included |
| `initialBackoffMs` | `int` | base delay before the first retry; zero applies the httpx default |
| `maxBackoffMs` | `int` | cap on the delay between attempts; zero applies the httpx default |

## `integration/adapter`

Adapter bindings. Scripts are compile-checked at save time; the database's
unique and foreign keys guard the binding itself.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `integration.adapter.query` | `AdapterSearch` + pageable meta | `page.Page[Adapter]` |
| `find_all` | `integration.adapter.query` | `AdapterSearch` | `Adapter[]` |
| `create` | `integration.adapter.create` | `AdapterParams` | created `Adapter` |
| `update` | `integration.adapter.update` | `AdapterParams` | updated `Adapter` |
| `delete` | `integration.adapter.delete` | primary-key params (`params.id`) | success |

`AdapterSearch` (query filters):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `systemId` | `string` | equals | filter by owning system |
| `contractId` | `string` | equals | filter by bound contract |
| `direction` | `string` | equals | `outbound` or `inbound` |
| `isEnabled` | `bool` | equals | filter by enablement |

`AdapterParams` (create/update):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | update only | primary key |
| `systemId` | `string` | Yes | the system this adapter belongs to |
| `contractId` | `string` | Yes | the contract this adapter implements |
| `direction` | `string` | No | `outbound` (default when omitted) or `inbound`; anything else fails with `ErrInvalidDirection` |
| `script` | `string` | Yes | the translation script; must compile (`ErrInvalidScript`) |
| `timeoutMs` | `int` | No | script run timeout override; zero inherits `vef.integration.run_timeout` |
| `isEnabled` | `bool` | No | disabled adapters refuse invocation (`ErrAdapterDisabled`) |

## `integration/route`

Routing rules. The contract reference is validated at save time because the
contract column carries the empty-string wildcard sentinel and has no
foreign key (`ErrInvalidRouteRef`).

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `integration.route.query` | `RouteSearch` + pageable meta | `page.Page[Route]` |
| `find_all` | `integration.route.query` | `RouteSearch` | `Route[]` |
| `create` | `integration.route.create` | `RouteParams` | created `Route` |
| `update` | `integration.route.update` | `RouteParams` | updated `Route` |
| `delete` | `integration.route.delete` | primary-key params (`params.id`) | success |

`RouteSearch` (query filters):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `routeKey` | `string` | contains | filter by route key fragment |
| `contractId` | `string` | equals | filter by scoped contract |
| `systemId` | `string` | equals | filter by target system |
| `isEnabled` | `bool` | equals | filter by enablement |

`RouteParams` (create/update):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | update only | primary key |
| `routeKey` | `string` | No | the key (tenant, branch, hospital area) this rule serves; empty is the default route |
| `contractId` | `string` | No | scopes the rule to one contract; empty applies to every contract. Exact `(key, contract)` matches win over contract-wildcard matches |
| `systemId` | `string` | Yes | the system serving matched invocations |
| `isEnabled` | `bool` | No | disabled rules never match |

## `integration/code_map`

Per-system value translation tables. Entries are index-built at save time so
a colliding or malformed map never reaches a lookup; when the host registers
an enumerable code set catalog, the `codeSet` identifier must also be one of
its registered sets.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `integration.code_map.query` | `CodeMapSearch` + pageable meta | `page.Page[CodeMap]` |
| `find_all` | `integration.code_map.query` | `CodeMapSearch` | `CodeMap[]` |
| `create` | `integration.code_map.create` | `CodeMapParams` | created `CodeMap` |
| `update` | `integration.code_map.update` | `CodeMapParams` | updated `CodeMap` |
| `delete` | `integration.code_map.delete` | primary-key params (`params.id`) | success |

`CodeMapSearch` (query filters):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `systemId` | `string` | equals | filter by owning system |
| `codeSet` | `string` | contains | filter by code set identifier fragment |
| `name` | `string` | contains | filter by name fragment |
| `isEnabled` | `bool` | equals | filter by enablement |

`CodeMapParams` (create/update):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | update only | primary key |
| `systemId` | `string` | Yes | the owning system |
| `codeSet` | `string` | Yes | translated code set identifier (e.g. `gender`); constrained to the host catalog when one is registered |
| `name` | `string` | Yes | display name |
| `entries` | `CodeMapEntry[]` | No | mapping pairs (below); duplicate lookup values per side are rejected (`ErrInvalidCodeMap`) |
| `onUnmapped` | `string` | No | `reject` (default when omitted — fail closed), `passthrough`, or `fallback` |
| `fallbackCanonical` | any JSON value | No | value `toCanonical` yields for unmapped input under the `fallback` policy |
| `fallbackExternal` | any JSON value | No | value `toExternal` yields for unmapped input under the `fallback` policy |
| `isEnabled` | `bool` | No | disabled maps behave as missing (`ErrMissingCodeMap`) |

`CodeMapEntry`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `canonical` | string / number / boolean | Yes | host-side primary value, emitted by `toCanonical` lookups |
| `external` | string / number / boolean | Yes | external-side primary value, emitted by `toExternal` lookups |
| `canonicalAliases` | array | No | additional host-side values matching this entry (matched, never emitted) |
| `externalAliases` | array | No | additional external-side values matching this entry |

## `integration/code_set`

Read-only view of the host's canonical code catalog, for the mapping
editor's pickers. Present only in the sense that it degrades: without a
`mold.CodeSetInspector` both operations answer `supported: false`.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `list_code_sets` | `integration.code_map.query` | none | `CodeSetCatalog` |
| `list_codes` | `integration.code_map.query` | `ListCodesParams` | `CodeCatalog` |

`ListCodesParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `codeSet` | `string` | Yes | the code set to enumerate |

`CodeSetCatalog` response:

| Field | Type | Description |
| --- | --- | --- |
| `supported` | `bool` | `false` when the host registered no enumerable catalog (editor falls back to free-text input) |
| `codeSets` | `CodeSetInfo[]` | entries with `codeSet` (identifier) and `name` (display name) |

`CodeCatalog` response:

| Field | Type | Description |
| --- | --- | --- |
| `supported` | `bool` | as above |
| `codes` | `CodeInfo[]` | entries with `code` (canonical value) and `label` (display name) |

## `integration/log`

Read-only invocation log: the paged view for browsing and the single-record
view for the full captures.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `integration.log.query` | `LogSearch` + pageable meta | `page.Page[InvocationLog]` |
| `find_one` | `integration.log.query` | `LogSearch` | one `InvocationLog` |

`LogSearch` (query filters):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `systemCode` | `string` | equals | filter by system code |
| `contractCode` | `string` | equals | filter by contract code |
| `direction` | `string` | equals | `outbound` or `inbound` |
| `failureKind` | `string` | equals | one of the [failure kinds](./overview#failure-vocabulary); empty rows are successes |
| `requestId` | `string` | equals | correlate with the API request that triggered the invocation |

`InvocationLog` response fields:

| Field | Type | Description |
| --- | --- | --- |
| `id` | `string` | log row ID |
| `systemCode` | `string` | system that served (or rejected) the invocation |
| `contractCode` | `string` | invoked contract |
| `direction` | `string` | `outbound` or `inbound` |
| `failureKind` | `string` | failure classification; empty for success |
| `durationMs` | `int` | wall time of the invocation |
| `input` | JSON | captured standard input (masked, truncated per `vef.integration.log`) |
| `output` | JSON | captured standard output (masked, truncated) |
| `httpTrace` | `HTTPExchange[]` | wire exchanges captured while the script ran (below) |
| `error` | `string` | failure message; absent on success |
| `requestId` | `string` | originating API request ID |
| `createdAt` / `createdBy` | timestamp / `string` | creation audit columns |

`HTTPExchange` (shared by the log and the dry-run trace):

| Field | Type | Description |
| --- | --- | --- |
| `method` | `string` | HTTP method |
| `url` | `string` | request URL (masked) |
| `requestHeaders` | `object` | request headers (credential headers always masked) |
| `requestBody` | `string` | captured request body (masked, truncated) |
| `status` | `int` | response status; `0` when the call never completed |
| `responseHeaders` | `object` | response headers |
| `responseBody` | `string` | captured response body (masked, truncated) |
| `durationMs` | `int` | exchange duration |
| `error` | `string` | transport error message when the call failed |

## `integration/ops`

Operational endpoints: the script test consoles, the connection probe, and
the routing diagnosis. Dry run and probing operate on disabled definitions
too — testing precedes enabling.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `dry_run` | `integration.ops.dry_run` | `DryRunParams` | `DryRunResult` |
| `dry_run_inbound` | `integration.ops.dry_run_inbound` | `DryRunInboundParams` | `InboundDryRunResult` |
| `test_connection` | `integration.ops.test_connection` | `TestConnectionParams` | `ConnectionCheck` |
| `diagnose_routes` | `integration.ops.diagnose_routes` | none | `RouteDiagnostics` |

### `dry_run`

Executes a script against a system under a contract and returns the output,
the failure classification, and the full wire trace. The calls it makes are
real; nothing is recorded to statistics or the invocation log.

Request (`DryRunParams`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `systemCode` | `string` | Yes | target system (disabled systems are allowed) |
| `contractCode` | `string` | Yes | contract whose schemas gate the run |
| `script` | `string` | No | unsaved editor content; empty falls back to the saved outbound adapter script (`ErrAdapterNotFound` when none exists) |
| `input` | any JSON value | No | invocation input, validated against the contract's input schema |

Response (`DryRunResult`):

| Field | Type | Description |
| --- | --- | --- |
| `output` | any JSON value | the script's return value (schema-validated); `null` when the run failed |
| `trace` | `HTTPExchange[]` | wire exchanges, populated even when the run failed — operators see how far the script got |
| `failureKind` | `string` | failure classification; absent on success |
| `error` | `string` | failure message; absent on success |

### `dry_run_inbound`

Executes an inbound script against a synthetic external request with the
business handler stubbed to return the supplied sample output. Verification
is bypassed (the console tests translation, not credentials), no business
code runs, and nothing is recorded; the contract schemas are enforced for
real on both sides of the dispatch.

Request (`DryRunInboundParams`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `systemCode` | `string` | Yes | target system |
| `contractCode` | `string` | Yes | contract whose schemas gate the dispatch |
| `script` | `string` | No | unsaved editor content; empty falls back to the saved inbound adapter script |
| `request` | `InboundRequestParams` | No | the synthetic external request (below) |
| `handlerOutput` | any JSON value | No | the sample the stubbed business handler returns; validated against the output schema |

`InboundRequestParams` (all optional):

| Field | Type | Description |
| --- | --- | --- |
| `method` | `string` | HTTP method of the synthetic request |
| `path` | `string` | request path |
| `headers` | `object` (string→string) | header names are normalized to lowercase, as a real gateway would deliver them |
| `query` | `object` (string→string) | query parameters |
| `body` | `string` | raw request payload |

Response (`InboundDryRunResult`):

| Field | Type | Description |
| --- | --- | --- |
| `reply` | any JSON value | the reply the external system would receive (including a `$response` envelope when the script returns one) |
| `dispatchedInput` | any JSON value | what the script dispatched to the (stubbed) handler — one value, or an array when the script dispatched multiple times |
| `failureKind` | `string` | failure classification; absent on success |
| `error` | `string` | failure message; absent on success |

### `test_connection`

Probes a saved system on every transport it configures. Probe failures are
data (`reachable: false`), not errors — the probe answered the question.
Configuration faults (unknown auth scheme, undecryptable credential) return
an error instead.

Request (`TestConnectionParams`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `systemCode` | `string` | Yes | system to probe |
| `method` | `string` | No | probe HTTP method; defaults to `GET` |
| `path` | `string` | No | probe path relative to the base URL; defaults to `/` |

Response (`ConnectionCheck`; each probe present iff the system configures
that transport):

| Field | Type | Description |
| --- | --- | --- |
| `http` | `HTTPProbe` | present when the system has a `baseUrl` |
| `http.reachable` | `bool` | whether the request completed |
| `http.status` | `int` | response status when reachable |
| `http.statusText` | `string` | status text when reachable |
| `http.durationMs` | `int` | probe duration |
| `http.error` | `string` | transport error when unreachable |
| `database` | `DatabaseProbe` | present when the system has a `dataSource` |
| `database.reachable` | `bool` | whether a throwaway connection succeeded |
| `database.version` | `string` | server version on success |
| `database.durationMs` | `int` | probe duration |
| `database.error` | `string` | connection error when unreachable |

### `diagnose_routes`

Reports the routing table's configuration gaps — dangling adapters, disabled
targets, uncovered contracts — before they surface as runtime errors.
Computed on demand; takes no parameters.

Response (`RouteDiagnostics`):

| Field | Type | Description |
| --- | --- | --- |
| `findings` | `RouteFinding[]` | one entry per gap (below); empty means the table is coherent |

`RouteFinding`:

| Field | Type | Description |
| --- | --- | --- |
| `kind` | `string` | finding classification (below) |
| `routeId` | `string` | involved route row, when one exists |
| `routeKey` | `string` | always meaningful — `""` is the default route |
| `contractCode` / `contractName` | `string` | involved contract, by code and display name |
| `systemCode` / `systemName` | `string` | involved system, by code and display name |

`kind` vocabulary:

| Kind | Meaning |
| --- | --- |
| `dangling_adapter` | a contract-scoped route whose target system has no enabled adapter for that contract — invoking through this rule fails with `ErrAdapterNotFound` |
| `wildcard_gap` | an enabled contract a wildcard (or default) route cannot serve because its target system has no enabled adapter for it. Informational |
| `disabled_system` | an enabled route targeting a disabled system — invocations through it fail with `ErrSystemDisabled` |
| `disabled_contract` | an enabled route scoped to a disabled contract — the rule can never match a successful invocation |
| `uncovered_contract` | an enabled contract that resolves to no rule under a route key present in the table — invoking it with that key fails with `ErrRouteNotFound`. Informational when the key intentionally routes a subset |

## See also

- [Overview](./overview) for the definition model and error codes
- [Outbound Calls](./outbound) and [Inbound Delivery](./inbound) for the flows behind `dry_run` / `dry_run_inbound`
- [Built-in Resources](../reference/built-in-resources) for the framework-wide resource index
