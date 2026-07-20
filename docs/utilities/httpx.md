---
sidebar_position: 9
---

# httpx (Outbound HTTP Client)

`httpx` is the framework's outbound HTTP client with a fluent
request API — construct one `Client` per third-party system and build
per-call `Request`s from it. The integration engine's scoped `http` script
library rides on it; application Go code can use it directly.

> Not to be confused with [`fiberx`](./small-helpers#fiberx), the package of
> Fiber request helpers — the outbound client is `httpx`.

## Quick Start

```go
import "github.com/coldsmirk/vef-framework-go/httpx"

client, err := httpx.New(
    httpx.WithBaseURL("https://api.example.com"),
    httpx.WithTimeout(10*time.Second),
    httpx.WithBearerToken(token),
    httpx.WithRetry(httpx.RetryConfig{}), // defaults: 3 attempts, 100ms→2s backoff
)

var out struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

resp, err := client.NewRequest().
    SetPathParam("id", "42").
    SetQuery("expand", "profile").
    Get(ctx, "/users/:id")
if err != nil {
    return err
}
if !resp.IsSuccess() {
    return fmt.Errorf("upstream returned %s", resp.Status())
}
if err := resp.JSON(&out); err != nil {
    return err
}
```

## Client

`httpx.New(opts ...Option)` validates options eagerly: a malformed base or
proxy URL and conflicting transport-level options fail construction. The
zero-option client is ready to use — no base URL, a 30s call timeout
(retries included), no retries. A `Client` is immutable after `New` and safe
for concurrent use; per-call state lives in the `Request`.

| Option | Behavior |
| --- | --- |
| `WithBaseURL(url)` | absolute URL every relative request URL joins onto; absolute request URLs bypass it |
| `WithTimeout(d)` | bounds a whole call, retries included (default 30s) |
| `WithHeader(k, v)` / `WithQuery(k, v)` | default header / query pair on every request |
| `WithBasicAuth(user, pass)` / `WithBearerToken(token)` | default `Authorization` header |
| `WithRetry(cfg)` | enables automatic retries (below) |
| `WithProxy(url)` | outbound proxy |
| `WithTLSConfig(cfg)` | custom TLS configuration |
| `WithCookieJar(jar)` | cookie persistence across calls |
| `WithMaxRedirects(n)` | redirect cap (default 10; exceeding fails with `ErrTooManyRedirects`) |
| `WithMaxResponseBody(n)` | response body byte cap (`ErrResponseTooLarge`) |
| `WithRequestHook(hooks...)` | runs after a request is fully built and before it is sent — the hook point for signing, audit, logging; a returned error aborts the call |
| `WithResponseHook(hooks...)` | runs after a response arrives and its body is buffered |
| `WithTransport(rt)` / `WithHTTPClient(hc)` | custom transport / fully custom `http.Client` (mutually exclusive with transport-level options — `ErrConflictingOptions`) |

The client sends `User-Agent: vef/<version>` unless the application sets its
own.

## Request

`client.NewRequest()` starts a fluent, single-use request builder
(re-executing one fails with `ErrRequestReused`):

| Group | Methods |
| --- | --- |
| Headers | `SetHeader`, `AddHeader`, `SetHeaders` |
| Query | `SetQuery`, `AddQuery`, `SetQueries` |
| Path params | `SetPathParam`, `SetPathParams` — substitute `:name` segments; unresolved segments fail with `ErrMissingPathParam` |
| Cookies / auth | `SetCookie`, `SetBasicAuth`, `SetBearerToken` |
| Body | `SetJSON(v)`, `SetXML(v)`, `SetBody(bytes, contentType)`, `SetBodyReader(r, contentType)`, `SetForm(map)`, `AddFormField(k, v)`, `AddFile(field, path)`, `AddFileReader(field, filename, r)` |
| Timeout | `SetTimeout(d)` — per-request override of the client timeout |
| Execute | `Get`, `Post`, `Put`, `Patch`, `Delete`, `Head`, `Options`, or `Do(ctx, method, url)` |
| Introspection | `Method()`, `URL()`, `Header(k)`, `Headers()`, `Body()`, `Context()` — the read surface request hooks use |

`SetForm`/`AddFormField` produce URL-encoded forms; adding files upgrades
the body to multipart automatically.

## Response

| Method | Contract |
| --- | --- |
| `StatusCode()` / `Status()` / `IsSuccess()` | status introspection; `IsSuccess` is 2xx |
| `Header(k)` / `Headers()` / `Cookies()` | response metadata |
| `Body()` / `String()` | the buffered body (always fully read and buffered) |
| `JSON(v)` / `XML(v)` | decode the body |
| `Duration()` | wall time of the call |
| `Attempts()` | attempts made, first call included |
| `Request()` | the originating request |

Non-2xx responses are **not** errors: the call succeeded at the transport
level, and the application decides what statuses mean. Errors are reserved
for transport failures, timeouts, and policy violations.

## Retries

`WithRetry(httpx.RetryConfig{...})` enables automatic retries. Zero fields
resolve to defaults:

| Field | Default | Meaning |
| --- | --- | --- |
| `MaxAttempts` | `3` | total attempts, the first call included |
| `InitialBackoff` | `100ms` | base delay before the first retry; doubles per retry, with full jitter |
| `MaxBackoff` | `2s` | cap on the delay between attempts, a server-sent `Retry-After` included |
| `RetryIf` | see below | custom predicate replacing the default policy entirely |

The default policy retries a transport error or a `429`/`502`/`503`/`504`
response, and **only for idempotent methods** (GET, HEAD, PUT, DELETE,
OPTIONS, TRACE) — a POST is never retried unless `RetryIf` allows it.

## Error Sentinels

| Error | Trigger |
| --- | --- |
| `ErrInvalidOption` | malformed base/proxy URL or other invalid option value |
| `ErrConflictingOptions` | `WithHTTPClient` combined with transport-level options |
| `ErrInvalidRequestURL` | unparsable request URL, or a relative URL without a base URL |
| `ErrMissingPathParam` | a `:name` segment left unresolved |
| `ErrRequestReused` | second execution of a single-use request |
| `ErrTooManyRedirects` | redirect cap exceeded |
| `ErrResponseTooLarge` | response body over the configured cap |

## See also

- [Integration Engine](../integration/overview) — systems configure `httpx` clients declaratively (auth schemes, retry policy, timeouts)
- [Small Helpers](./small-helpers) — `fiberx`, the inbound Fiber request helpers formerly named `httpx`
