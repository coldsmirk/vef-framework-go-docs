---
sidebar_position: 2
---

# JS Engine

The `js` package embeds a sandboxed JavaScript runtime powered by
[goja](https://github.com/dop251/goja). It is built around an
**Engine / Runtime / Lib** split: an immutable `Engine` holds a validated
set of libraries and stamps out single-use `Runtime`s; capabilities (HTTP,
SQL, cache, events, ...) are `Lib`s installed per runtime, so scripts can
touch nothing beyond what was installed.

The framework uses this engine for integration adapter scripts, signing and
verification scripts, and other script-execution seams; applications can use
it directly for their own scripting needs.

## Quick Start

```go
import "github.com/coldsmirk/vef-framework-go/js"

engine, err := js.NewEngine()
if err != nil {
    return err
}

// One engine serves the whole application; one runtime per execution.
rt, err := engine.NewRuntime(js.WithRunTimeout(5 * time.Second))
if err != nil {
    return err
}

value, err := rt.RunString(ctx, `1 + 2`)
fmt.Println(value.Export()) // 3
```

## Engine

`js.NewEngine(opts ...EngineOption)` builds an engine, validating the
library set eagerly: a nil library, an empty name, or a name collision
(across standard, always-on, and catalog libraries) fails construction.

| Option | Behavior |
| --- | --- |
| `js.WithBaseLibs(libs...)` | registers always-on libraries: installed into every runtime without opt-in. Reserve for safe, ubiquitous utilities |
| `js.WithLibs(libs...)` | registers catalog libraries: opt-in per runtime via `js.EnableLibs` |
| `js.WithoutStdLibs()` | builds a bare engine whose runtimes start without the standard library bundle |

Errors: `js.ErrInvalidLib` (nil lib / empty name), `js.ErrDuplicateLib`
(name collision), `js.ErrLibNotFound` (enabling an unregistered name).

The engine is safe for concurrent use. The framework provides a shared DI
engine with the built-in capability libraries pre-wired in two tiers:

- **always-on** (installed into every runtime): `console`, `crypto`, and
  `cache` (in-memory store, key prefix `js:`);
- **opt-in catalog** (activated per runtime via `EnableLibs`): `events`,
  `http` (30s default timeout), and `sql` (primary database, read-only).

Applications contribute or replace libraries with `vef.ProvideJSLib`: a lib
whose name matches a built-in replaces it in its tier, a new name joins the
opt-in catalog.

```go
vef.ProvideJSLib(func(db orm.DB) js.Lib {
    return jssql.New(db, config.Postgres, jssql.WithExecute())
})
```

## Runtime

`Engine.NewRuntime(opts ...RuntimeOption)` creates a fresh runtime carrying
the engine baseline plus the activated catalog libraries. A runtime is
**not** safe for concurrent use, only one Run call may be in flight at a
time, and it should be discarded after use — create one per goroutine per
execution.

| Option | Behavior |
| --- | --- |
| `js.EnableLibs(names...)` | activates catalog libraries by name; installation follows argument order; enabling twice is a no-op; unknown names fail with `ErrLibNotFound` |
| `js.WithRunTimeout(d)` | caps every Run call; combines with the caller's context — whichever deadline is earlier wins |
| `js.WithMaxCallStackSize(n)` | bounds the JS call stack depth, guarding against runaway recursion |

| Method | Contract |
| --- | --- |
| `RunProgram(ctx, program)` / `RunString(ctx, source)` | executes under `ctx`: cancellation interrupts the running script and, through `Context()`, any in-flight host library IO; the returned error is then the context's error |
| `Set(name, value)` | binds a global variable |
| `Context()` | the context of the in-flight Run call (or `context.Background()` when idle); host libraries must issue IO through it |
| `AsFunction(value)` | converts a value into a callable `js.Func` handle; calls must happen on the goroutine currently driving the runtime |
| `VM()` | the underlying `*goja.Runtime`, for advanced library authoring |

Every runtime is configured with JSON field-name mapping
(`goja.TagFieldNameMapper("json", true)`), so Go structs passed via `Set`
read naturally from scripts.

## The Standard Library Bundle

Unless the engine was built `WithoutStdLibs`, every runtime starts with one
vendored esbuild bundle installing each library under its ecosystem-native
global name:

| Global | Library | Purpose |
| --- | --- | --- |
| `BigNumber` | bignumber.js | arbitrary-precision decimals |
| `dayjs` | Day.js | date/time parsing, formatting, arithmetic |
| `fxp` | fast-xml-parser | XML parsing and building (`fxp.XMLParser`, `fxp.XMLBuilder`) |
| `radashi` | Radashi | functional utility helpers |
| `z` | Zod | schema validation, with the `en` and `zh-CN` locales bundled (default `zh-CN`) |
| `URL` / `URLSearchParams` | core-js polyfills | WHATWG URL handling |

## Built-in Capability Libraries

Capabilities with side effects are separate packages under `js/*`, wired
into the DI engine's catalog. Each installs one global object named after
the package; a runtime only sees the ones activated for it. All host IO
flows through `Runtime.Context()`, so cancellation reaches blocking Go
calls.

### `jssql` — global `sql`

```js
sql.queryList('SELECT name FROM users WHERE age > ?', 18)  // → [{...}, ...]; [] when no rows match
sql.queryOne('SELECT ... WHERE id = ?', id)                // → {...} | null
sql.execute('UPDATE ...', args)                            // → { rowsAffected }
```

- `jssql.New(db, kind, opts...)` builds the library over a chosen data
  source; the caller decides what scripts can reach.
- Only placeholder binding — deliberately no string interpolation helper.
- Read-only by default, enforced fail-closed by an AST-based guard
  (`sqlguard`): writing CTEs, stacked statements, and dialect-specific
  side-effecting functions are rejected; unparseable SQL is refused
  (`ErrQueryNotReadOnly`).
- `sql.execute` throws `ErrExecuteDisabled` unless built with
  `jssql.WithExecute()`.
- Result sets are capped (`jssql.WithMaxRows`, default 1000;
  `ErrTooManyRows` instructs the script to add `LIMIT`).

### `jshttp` — global `http`

A synchronous take on the fetch standard:

```js
http.fetch(url, { method, headers, query, body, redirect, timeout })
http.get(url, options?)          // sugar over fetch, likewise put/patch
http.post(url, body, options?)
http.delete(url, options?)
```

- Every call returns `{ status, statusText, ok, url, redirected, headers,
  body, text(), json(), arrayBuffer() }`; failures throw catchable
  exceptions. Response header names are lower-cased, multi-values joined
  with `", "`.
- Beyond fetch: `query` appends URL parameters, `timeout` (milliseconds)
  replaces AbortSignal, `redirect` supports `follow` / `error` / `manual`.
- By default nothing is restricted — timeouts, body size caps, host
  allowlists, and the private-network guard all activate through options
  (`jshttp.New(opts...)`).

### `jscache` — global `cache`

The only channel scripts have for keeping state across executions, since
runtimes are discarded after each run:

```js
cache.set('counter', { n: 1 })       // store with default TTL
cache.set('token', value, 60000)     // TTL in milliseconds
cache.get('counter')                 // → value | null
cache.has('counter')                 // → boolean
cache.delete('counter')
```

`jscache.New(store, opts...)` takes any `cache.Cache[any]` — memory for
single-node state, Redis for shared state — plus `WithKeyPrefix` for
namespacing.

### `jsevents` — global `events`

```js
events.publish('report.generated', { reportId: id })
```

The payload is JSON-encoded and published as an `event.RawPayload`, so Go
subscribers decode it with `event.SubscribeTyped` as usual. Publishing is
deliberately the only verb — subscriptions belong to the host, not to a
per-execution runtime. `jsevents.WithAllowedTypes(patterns...)` restricts
the publishable type namespace.

### `jscrypto` — global `crypto`

```js
crypto.md5(data)                    // likewise sha1 / sha256 / sha512 / sm3
crypto.hmac('sha256', key, data)    // hex digest
crypto.base64Encode(data) / crypto.base64Decode(encoded)
crypto.hexEncode(data) / crypto.hexDecode(encoded)
crypto.uuid()
```

Digests are lower-case hex; inputs are UTF-8 strings. The weak digests
(md5, sha1) exist for legacy API signature interop — password storage
belongs to the security module's encoders.

### `jsconsole` — global `console`

```js
console.info('processed', count, payload)   // also warn / error
```

Arguments are joined with a space; strings pass through, errors render
their message, everything else is JSON-encoded. Backed by a `logx.Logger`
(pass a named logger to distinguish script output).

## Writing Your Own Lib

```go
type Lib interface {
    Name() string          // unique key within an Engine; by convention the global it installs
    Install(rt *js.Runtime) error
}
```

- Pure-JS libraries: `js.SourceLib(name, source)` (compiles eagerly) or
  `js.ProgramLib(name, program)`.
- Host capabilities: implement `Install` with `rt.Set(...)`, holding only
  shared goroutine-safe dependencies — never per-runtime state — and issue
  IO through `rt.Context()` so cancellation propagates.

```go
lib, err := js.SourceLib("fmtx", `var fmtx = { pad: (s, n) => String(s).padStart(n, '0') };`)
engine, err := js.NewEngine(js.WithLibs(lib))
rt, err := engine.NewRuntime(js.EnableLibs("fmtx"))
```

## Compilation Helpers

| API | Contract |
| --- | --- |
| `js.Compile(name, source, strict)` / `js.MustCompile(...)` | pre-compile for repeated execution (goja aliases) |
| `js.Parse(name, source)` | returns `*js.AstProgram` (source maps disabled) |
| `js.Runtime`, `js.Value`, `js.Object`, `js.Program`, `js.AstProgram` | goja type aliases |
| `js.IsNaN`, `js.IsString`, `js.IsBigInt`, `js.IsNumber`, `js.IsInfinity`, `js.IsUndefined`, `js.IsNull` | goja helper aliases |

## Where the Framework Runs Scripts

| Seam | Libraries visible |
| --- | --- |
| [Integration adapter scripts](../integration/outbound#adapter-script-environment) | baseline + `errors`, `codes`, scoped `http` / `sql` per system |
| [Integration signing / verification scripts](../integration/outbound#outbound-authentication-schemes) | baseline only (zero IO) + `request`, `params` bindings |
| [Expression engine](./expression) | its own evaluator (not this engine) |

## Thread Safety

> **WARNING**: an `Engine` is safe for concurrent use; a `Runtime` is not.
> Create one runtime per goroutine per execution and discard it after use.
