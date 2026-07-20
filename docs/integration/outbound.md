---
sidebar_position: 2
---

# Outbound Calls

Outbound is the flow where business code invokes an external system: you call
a contract, the engine resolves the serving system, runs its outbound adapter
script, and hands back the schema-validated standard model.

## Invoking a Contract

```go
type PatientService struct {
    invoker integration.Invoker
}

func (s *PatientService) LoadPatient(ctx context.Context, hospitalArea, patientID string) (*Patient, error) {
    // Typed convenience wrapper: decodes the standard model into Patient.
    return integration.Call[*Patient](ctx, s.invoker, "patient.get",
        map[string]any{"patientId": patientID},
        integration.WithRoute(hospitalArea),
    )
}
```

`integration.Invoker` is available from DI whenever `vef.IntegrationModule`
is enabled:

| API | Contract |
| --- | --- |
| `Invoke(ctx, contract string, input any, opts ...InvokeOption) (*Result, error)` | validates input against the contract's input schema, runs the adapter script, validates the output, returns the `Result` |
| `integration.Call[T](ctx, inv, contract, input, opts...) (T, error)` | typed wrapper over `Invoke` that decodes the output into `T` through a JSON round-trip |

### Invoke options

| Option | Behavior |
| --- | --- |
| `integration.WithSystem(code)` | targets a system directly, bypassing route resolution. Mutually exclusive with `WithRoute` (`ErrTargetAmbiguous`) |
| `integration.WithRoute(key)` | selects the system through the `RouteResolver`; no target option resolves the empty route key (the default route) |
| `integration.WithTimeout(d)` | overrides the script run timeout for this invocation |
| `integration.WithCache(ttl)` | caches the validated output for `ttl`, keyed by system, contract, and input. Off by default — opt in only where the business tolerates data of that age. Cache hits bypass the invocation log and statistics |

### Result

| Method | Meaning |
| --- | --- |
| `Output() any` | the standard model, already validated against the output schema |
| `Decode(v any) error` | unmarshals the output into `v` through a JSON round-trip |
| `System() string` | code of the system that served the invocation |
| `Duration() time.Duration` | wall time of the invocation |
| `Cached() bool` | whether the output came from the response cache |

## Adapter Script Environment

The outbound adapter script translates contract input into calls on the
external system and its responses into the contract output. The value of the
script's final expression is the output (validated against the contract's
output schema).

Every script sees the JS engine baseline (`BigNumber`, `dayjs`, `fxp`,
`radashi`, `z`, `URL` / `URLSearchParams` — see [JS Engine](../data-tools/js-engine))
plus these bindings:

| Binding | Present | Contents |
| --- | --- | --- |
| `input` | always | the schema-validated contract input |
| `system` | always | read-only system view: `{ code, name, params }` — never credentials |
| `errors` | always | failure classification: `errors.upstream(message)` throws an exception recorded as an upstream failure instead of a script bug |
| `codes` | always | code map translation: `codes.toExternal` / `codes.toCanonical` / `codes.entries` — see [Code Maps](./code-maps) |
| `http` | when the system has a `baseUrl` | system-scoped HTTP client (below) |
| `sql` | when the system has a `dataSource` | system-scoped SQL access (below) |

A script reaching for an unconfigured capability fails with a plain
`ReferenceError`.

### The scoped `http` client

```js
// Paths are always relative — the client is locked to the system's base URL.
const res = http.get('/api/patients/' + input.patientId);
if (!res.ok) {
    errors.upstream('patient service returned HTTP ' + res.status);
}
const data = res.json();

({ patientId: data.id, name: data.patientName })
```

| Function | Behavior |
| --- | --- |
| `http.fetch(path, { method, headers, query, body, timeout, envelope })` | full request form |
| `http.get(path, options?)` / `http.delete(path, options?)` | sugar over fetch |
| `http.post(path, body, options?)` / `http.put(...)` / `http.patch(...)` | sugar over fetch with a body |

- Authentication is injected by the host per the system's `outboundAuth`;
  scripts never see credentials and cannot reach other hosts (absolute URLs
  are rejected).
- `body`: strings and byte arrays pass through verbatim; any other value is
  JSON-encoded with an implied `application/json` content type.
- `timeout` (milliseconds) may only shorten the system's call timeout.
- Every call returns the fetch `Response` shape: `{ status, statusText, ok,
  url, headers, body, text(), json(), arrayBuffer() }`. Header names are
  lower-cased, multi-values joined with `", "`.
- Failures throw catchable exceptions; a transport failure that the script
  does not catch classifies the invocation as `transport`.
- Every wire exchange is recorded into the invocation trace (masked and
  truncated per `vef.integration.log`).

### The scoped `sql` library

For systems whose integration surface is a database (external views,
exchange tables):

```js
const rows = sql.queryList('SELECT id, name FROM v_patients WHERE id = ?', input.patientId);
const one  = sql.queryOne('SELECT ... WHERE id = ?', input.id);  // object | null
```

- Only placeholder binding is offered — deliberately no string interpolation.
- Read-only by default, enforced fail-closed by an AST-based guard: writing
  CTEs, stacked statements, and dialect-specific side-effecting functions are
  rejected, and unparseable SQL is refused. Write portable statements — the
  guard's parser refuses dialect-only syntax such as T-SQL `SELECT TOP`,
  bracket quoting, `WITH (NOLOCK)`, or Oracle `(+)` joins and `CONNECT BY`
  (portable forms like `OFFSET/FETCH`, `ROWNUM`, `NVL`, `TO_CHAR` pass).
- `sql.execute('UPDATE ...', args)` (returns `{ rowsAffected }`) throws unless
  the system's data source declares `mode = "read_write"`.
- Result sets are capped (over-limit queries fail with an instruction to add
  `LIMIT`); no matching rows yield `[]`, not `null`.

## Outbound Authentication Schemes

`system.outboundAuth` selects how the framework authenticates against the
system: `{ "scheme": "...", "params": { ... }, "script": "..." }`. Values of
scheme-declared sensitive parameters are stored encrypted and masked in
management responses.

| Scheme | Params | Behavior |
| --- | --- | --- |
| `none` | — | requests go out unauthenticated |
| `http_basic` | `username`, `password` (sensitive) | RFC 7617 Basic credentials |
| `bearer` | `token` (sensitive) | static `Authorization: Bearer` token |
| `header` | one entry per credential header (all sensitive) | sends every configured pair as static headers |
| `query` | one entry per credential parameter (all sensitive) | sends every configured pair as static query parameters |
| `signature` | `appId`, `secret` (hex, sensitive) | signs every request with the framework HMAC convention: `x-timestamp` / `x-nonce` / `x-signature` over the identity, method, and path — two VEF deployments authenticate each other with configuration only |
| `script` | free-form params (all sensitive) + signing body | runs the custom signing body per request in a zero-IO runtime; it sees the built `request` (`{ method, url, path, query, headers, body }`) and decrypted `params`, and returns an object of credential headers to add |

Applications register custom schemes with
`vef.ProvideIntegrationOutboundAuthScheme`; a scheme whose name matches a
built-in replaces it. A custom scheme implements:

```go
type OutboundAuthScheme interface {
    Name() string
    Apply(cfg *OutboundAuthConfig) ([]httpx.Option, error)
    SensitiveParams() []string // integration.SensitiveAll marks every param sensitive
}
```

## System Envelopes

Most vendor APIs repeat one wire structure on every endpoint (`{code, msg,
data}` responses, signed request wrappers, SOAP envelopes). Configure it once
per system in `outboundEnvelope` instead of repeating it in every adapter:

```js
// outboundEnvelope.request — wraps the request the adapter issued.
// `request` = { method, path, headers, query, body }; fields the returned
// object omits keep the adapter's values.
({ body: { reqData: request.body, ts: Date.now() } })
```

```js
// outboundEnvelope.response — unwraps the completed HTTP response.
// `response` is the fetch Response shape; whatever this returns is what the
// adapter's call yields.
const parsed = response.json();
if (parsed.code !== '0') {
    errors.upstream(parsed.msg);   // vendor-level errors classified once for the whole system
}
parsed.data
```

Either script may be empty (save-time validation requires at least one), and
individual calls opt out with `http.get(path, { envelope: false })` — for
deviant endpoints such as file downloads or health checks.

## Failure Classification and Retries

- Transport failures and 429/502/503/504 responses on idempotent methods are
  retried per the system's `retry` policy before they surface.
- The failure vocabulary (`input_invalid`, `output_invalid`, `upstream`,
  `transport`, `timeout`, `canceled`, `script`, `config`) is shared with the
  invocation log and statistics; see [Overview](./overview#failure-vocabulary).
- Caller cancellation is distinguished from timeout: a caller that walked
  away classifies as `canceled`, an exceeded deadline as `timeout`.

## Testing Before Enabling

The `integration/ops` resource ships a test console that operates on saved
and unsaved definitions alike — testing precedes enabling:

- `dry_run` executes a script (possibly unsaved editor content) against a
  system under a contract and returns the output, failure classification,
  and the full wire trace. The calls it makes are real; nothing is recorded
  to statistics or the invocation log.
- `test_connection` probes a saved system on every transport it configures
  (HTTP base URL, database) and reports reachability as data.

Field-by-field request/response documentation for both lives in
[RPC Resources](./resources#integrationops).

## Next Step

[Inbound Delivery](./inbound) covers the opposite flow — external systems
calling into your application.
