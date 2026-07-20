---
sidebar_position: 3
---

# Inbound Delivery

Inbound is the flow where an external system calls into your application:
the HTTP gateway receives the vendor's request, verifies it against the
system's inbound auth, runs the inbound adapter script — which translates the
wire request into the contract's standard input, dispatches it to your
business handler, and shapes the reply the vendor expects.

## The HTTP Gateway

When `vef.IntegrationModule` is enabled, the framework registers:

```
POST /integration/inbound/:systemCode/:contractCode
```

- The endpoint is a real route (registered after the API engine, before the
  SPA fallback); non-integration traffic pays nothing for it.
- It deliberately bypasses the `/api` dispatch model: external replies need
  raw control of status and body, never the standard result envelope.
- Deliveries are rate-limited per `(system, client IP)` in a sliding window
  (`vef.integration.inbound.rate_limit`, default 120 per minute, enforced
  per node).

The gateway translates each HTTP request into the protocol-neutral
`integration.InboundRequest` envelope:

| Field | Contents |
| --- | --- |
| `SystemCode` / `ContractCode` | resolved from the URL path parameters |
| `Protocol` | `"http"` |
| `Method` / `Path` / `Query` | HTTP-native request data |
| `Headers` | header names lowercased, multi-values joined with `", "` |
| `Body` | the raw request payload |
| `ClientAddr` | network peer address, for IP-based verification |

## Inbound Authentication

Verification is fail-closed: a system without `inboundAuth` refuses inbound
delivery entirely — the `none` scheme opens it up deliberately. Every
verification failure returns the uniform `ErrInboundAuthFailed` (HTTP 401);
missing configuration, missing credentials, and wrong credentials are
indistinguishable to the caller, and an unknown or disabled system denies
exactly the same way so system codes cannot be enumerated.

Built-in schemes mirror the outbound vocabulary — the same name verifies the
wire format its outbound counterpart sends; `ip` is inbound-only:

| Scheme | Params | Behavior |
| --- | --- | --- |
| `none` | — | accepts every caller; pair with network-level controls |
| `ip` | `whitelist` (comma-separated IP/CIDR entries) | verifies the caller's network address; an empty whitelist is treated as a config fault, not open access |
| `http_basic` | `username`, `password` (sensitive) | verifies RFC 7617 Basic credentials in constant time |
| `bearer` | `token` (sensitive) | verifies a static `Authorization: Bearer` token in constant time |
| `header` | one entry per credential header (all sensitive) | request must present every configured pair (AND); blank configured values never authenticate |
| `query` | one entry per credential parameter (all sensitive) | request must present every configured pair (AND) |
| `signature` | `secret` (hex, sensitive) | verifies the framework HMAC convention: `x-timestamp` / `x-nonce` / `x-signature` signed over the system code, method, and path — replay-protected through the shared nonce store |
| `script` | free-form params (all sensitive) + verification body | runs the custom verification body in a zero-IO runtime; it sees `request` and decrypted `params` and grants access by returning a truthy value |

Applications register custom schemes with
`vef.ProvideIntegrationInboundAuthScheme`:

```go
type InboundAuthScheme interface {
    Name() string
    Verify(ctx context.Context, req *InboundRequest, auth *InboundAuthConfig) error
    SensitiveParams() []string
}
```

## Business Handlers

The business side of an inbound contract is an `integration.InboundHandler` —
one handler per contract code, registered at boot:

```go
vef.ProvideIntegrationInboundHandler(func(db orm.DB) integration.InboundHandler {
    return integration.NewInboundHandler("patient.sync",
        func(ctx context.Context, input PatientSyncInput) (PatientSyncOutput, error) {
            // input already validated against the contract's input schema
            return doSync(ctx, db, input)
        })
})
```

- `integration.NewInboundHandler[I, O](contract, handle)` adapts a typed
  function; the schema-validated input is decoded into `I` through a JSON
  round-trip.
- Handlers must be idempotent: external systems deliver at-least-once.
- Blank or duplicate contract registrations fail at boot.
- A delivery for a contract with no registered handler fails with
  `ErrInboundHandlerMissing` (HTTP 501) — a deployment fault, not a caller
  error.

## Adapter Script Environment

The inbound adapter script (direction `inbound`) sees the engine baseline
plus:

| Binding | Contents |
| --- | --- |
| `request` | read-only wire request: `{ protocol, method, path, headers, query, body, clientAddr }` — header names lowercased, `body` as string |
| `system` | read-only system view: `{ code, name, params }` |
| `dispatch(input)` | validates `input` against the contract's input schema, runs the business handler, validates its output against the output schema, and returns it; failures surface as catchable exceptions |
| `codes` | code map translation (external code values are the inbound script's core job) — see [Code Maps](./code-maps) |

A typical script:

```js
// 1. Translate the vendor's wire request into the standard input.
const body = JSON.parse(request.body);
const input = {
    patientId: body.patient_no,
    gender: codes.toCanonical('gender', body.sex),
};

// 2. Dispatch to the business handler (schema-enforced both ways).
const output = dispatch(input);

// 3. Shape the reply the vendor expects.
({ code: '0', msg: 'ok', data: { syncId: output.syncId } })
```

Batch payloads may call `dispatch` multiple times; every dispatch is
recorded, and the last dispatch failure stays sticky for classification even
when the script catches it to shape a partial-success reply.

## Shaping the Reply

The script's final value becomes the HTTP response:

- `null` / `undefined` → `200` with an empty body.
- Any ordinary value → `200` with the value as JSON.
- An object of the form `{ $response: { status, headers, body } }` takes raw
  control: `status` defaults to 200 (out-of-range values fall back),
  `headers` are set verbatim, a string `body` is sent as `text/plain` unless
  the script set its own content type, any other body is JSON. The
  `$`-prefixed marker cannot collide with a plausible business payload.

```js
({ $response: {
    status: 200,
    headers: { 'content-type': 'application/xml' },
    body: '<Response><Code>0</Code></Response>',
}})
```

Uncaught pipeline errors are remapped to statuses an external caller's retry
logic can act on (business errors on the `/api` surface ride HTTP 200, but
external callers do not read the envelope):

| Condition | Status |
| --- | --- |
| verification failed, unknown/disabled system | `401` |
| rate limit exceeded | `429` |
| contract/adapter missing or disabled (behind a verified caller) | `404` (uniform, so callers cannot probe which piece exists) |
| input rejected by the contract schema | `400` |
| no registered business handler | `501` |
| script fault, render fault, anything else | `500` |

## Observability

- Verified deliveries are recorded to the invocation log per
  `vef.integration.log` with direction `inbound`; rejected deliveries stay
  out of the log deliberately (unauthenticated traffic must not grow the
  durable evidence trail) and fold into statistics under an empty contract.
- The `dry_run_inbound` operation is the inbound test console: it executes an
  inbound script against a synthetic external request with the business
  handler stubbed to return a supplied sample output. Verification is
  bypassed — the console tests translation, not credentials; nothing runs
  against business code and nothing is recorded. See
  [RPC Resources](./resources#integrationops).

## Next Step

[Code Maps](./code-maps) covers per-system value translation used by both
flows.
