---
sidebar_position: 3
---

# CLI Tools

VEF ships a CLI with two distinct families of commands, and the split is a
contract with you rather than cosmetics:

- **`new`** writes a file once and hands it over. There is no generated-file
  header, it is never regenerated, and it is yours to edit.
- **`generate-*` and `export-api`** own their output forever and overwrite it on
  every run.

The test for which family a command belongs to is one question: *will anyone
want to edit this?* Yes ⇒ `new`; no ⇒ a derived artifact carrying
`// Code generated ... DO NOT EDIT.`.

## Current commands

The root command is `vef-cli`. `vef-cli --version` prints the CLI banner plus
`Version: ...`; when build-date metadata is available it also prints
`Built: ...`, and dirty VCS builds append `-dirty` to the version string.

The CLI registers these subcommands:

- `new` — `project`, `resource`, `service`
- `generate-build-info`
- `generate-model-schema`
- `export-api`

## Minimal command examples

```bash
vef-cli --version
vef-cli new project acme-server
vef-cli new resource --table hr_employee --module hr
vef-cli new service --name Payroll --module hr --deps bus:event.Bus
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
vef-cli generate-model-schema -i models -o schemas -p schemas
vef-cli export-api --app ./cmd/server -o api-manifest.json
```

Application code should consume the CLI through these commands instead of
importing the `cmd/vef-cli/cmd/*` implementation packages directly.

## `new`

### `new project <name>`

Writes the framework's own layout: an entry point under `cmd/server`, runtime
configuration under `configs`, business modules under `internal`, and a
`vef.yml` recording the code-generation conventions the other `new` commands
read. It pins the framework to the CLI's own version — the tool generates
against the API it was built from, which cannot go stale — and runs
`generate-build-info` before finishing, because `internal/vef/module.go`
references the generated `BuildInfo`.

| Flag | Default | Purpose |
| --- | --- | --- |
| `--name`, `-N` | the positional argument | project name |
| `--module`, `-m` | the project name | Go module path |
| `--path`, `-p` | `./<name>` | directory to create the project in |
| `--with-example` | `true` | generate a starter business module |
| `--skip-tidy` | `false` | skip `go mod tidy` |
| `--skip-git` | `false` | skip `git init` |
| `--dry-run` | `false` | print what would be written without touching the filesystem |

Dockerfiles, git hooks and CI pipelines are deliberately excluded: those encode
deployment and team decisions the framework has no opinion about.

### `new resource`

Generates a model, payload and API resource from a live database table — **plus
both registrations**, the model registry's `var` block and the module's fx
option list. The registrations are half the value: forgetting either produces
an application that compiles and answers 404. A module that does not exist yet
is created and wired into `cmd/server/main.go`'s `vef.Run` call; that last step
is best-effort, since the entry point's shape is a convention rather than a
contract, so a project that wires itself differently gets a `TODO` line instead
of a failed command.

| Flag | Default | Purpose |
| --- | --- | --- |
| `--table`, `-t` | required | database table to derive the entity from |
| `--module`, `-m` | required | business module to generate into, may be nested |
| `--entity` | table name without the module prefix | entity name override, snake_case |
| `--alias` | the initials of the table's words | table alias override |
| `--ops` | the project's configured set | CRUD operations to embed |
| `--search` | every scalar column | search criteria as `column:operator`; `none` generates an empty search payload |
| `--source` | primary | data source to inspect |
| `--config` | `<project>/configs/application.toml` | path to `application.toml` |
| `--force` | `false` | overwrite generated files that already exist |
| `--dry-run`, `-n` | `false` | print what would be written without touching the filesystem |

The table is inspected through the framework's own schema service, so a column
comment becomes a `label` tag, nullability becomes a pointer plus `omitempty`, a
declared character bound becomes `max=N`, and a table carrying the framework's
audit columns embeds the matching `orm` mixin instead of re-declaring them. The
mixin match checks column **types**, not just names — `orm.Model` declares
`ID string` and the ORM writes a generated XID into any zero-valued string
primary key, so matching a `BIGINT` identity column on its name alone would push
a 20-character string into a numeric column on every create.

`find_tree` is refused rather than rendered: `crud.NewFindTree` takes a
tree-building function no generator can supply.

### `new service`

The convention-shaped skeleton — a struct of injected dependencies only,
methods taking `(ctx, db orm.DB, …)` — registered with its module. Framework
dependency types resolve their own imports.

| Flag | Default | Purpose |
| --- | --- | --- |
| `--name`, `-n` | required | service name in PascalCase, with or without the `Service` suffix |
| `--module`, `-m` | required | business module to generate into, may be nested |
| `--deps` | none | injected dependencies as `field:Type`, for example `bus:event.Bus` |
| `--force` | `false` | overwrite the service file if it already exists |
| `--dry-run` | `false` | print what would be written without touching the filesystem |

### Conventions come from `vef.yml`

`vef.yml` at the project root records `module_root` and, under `resource`, the
`name` / `permission` templates (placeholders `{module}`, `{domain}`,
`{entity}`, `{action}`), `ops`, `audit` and `audit_user_model`. Every key has a
default, so the file is optional; `new project` writes one so the conventions
are visible rather than implied.

Re-running a generator over existing code reports what already exists instead of
overwriting it, unless `--force` says otherwise. Every command plans all its
effects first, so `--dry-run` shows the whole thing and a mid-planning failure
leaves nothing half-written.

## `generate-build-info`

This command generates a Go source file containing build metadata such as:

- app version
- build time
- git commit

It is designed to be used from `go:generate` or from your build pipeline.

Flags:

| Flag | Default | Purpose |
| --- | --- | --- |
| `--output`, `-o` | `build_info.go` | output Go file |
| `--package`, `-p` | `main` | package name for the generated file |

The generated file exports `BuildInfo = &monitor.BuildInfo{...}` and fills:

- `AppVersion` from `git describe --tags --always --dirty`, falling back to `dev`
- `BuildTime` from `timex.Now().String()`
- `GitCommit` from `git rev-parse HEAD`, falling back to `none`

The generator creates the output directory when needed. The public shape of the
generated file is:

```go
var BuildInfo = &monitor.BuildInfo{
	AppVersion: "...",
	BuildTime:  "...",
	GitCommit:  "...",
}
```

## `generate-model-schema`

This command inspects model files and generates type-safe schema helpers for ORM usage.

It supports:

- file-to-file generation
- directory-to-directory generation

The goal is to reduce hard-coded column-name strings in query code.

Flags:

| Flag | Default | Purpose |
| --- | --- | --- |
| `--input`, `-i` | required | input model file or directory |
| `--output`, `-o` | required | output schema file or directory |
| `--package`, `-p` | `schemas` | package name for generated schema files |

Directory input writes one schema file per input file. Directory mode processes
only `*.go` files directly inside the input directory; it is not recursive. Test
files (`_test.go`) and files excluded by build constraints are skipped. For
directory input, the output may be an existing directory or a directory path that
does not exist yet (it is created as needed). If the output path already exists
as a file, directory-to-file generation is rejected.

The generator reads structs in the target file that embed `orm.BaseModel`.
Table metadata comes from the embedded `orm.BaseModel` field's `bun` tag: the
bare name segment (e.g. `bun:"users"`) or the `table:...` option sets the table
name (`table:` wins), and `alias:...` sets the default alias. Without those tag
parts, the table defaults to the pluralized snake_case model name and the alias
defaults to the singular snake_case model name.

Field handling is source-compatible with these rules:

- only exported fields generate accessors
- `bun:"-"` fields are skipped
- `bun:"rel:*"` and `bun:"m2m:*"` relationship fields are skipped
- a first `bun` tag component such as `bun:"user_name"` sets the column name
- an explicit `column:...` option overrides both the bare name segment and the field name
- fields without a column tag use the field name in snake_case
- embedded structs are expanded
- `bun:"embed:prefix_"` expands nested fields with the prefix
- `label:"..."` becomes a method comment in generated code
- `bun:",scanonly"` fields still get accessors but are excluded from `Columns()`

The generated public API exposes an exported schema variable named after the
model, for example `User`, backed by an unexported schema type such as
`userSchema`. Each schema has field accessors plus `Table()`, `Alias()`,
`As(alias)`, and `Columns()`.

Field accessors return alias-qualified columns with `dbx.ColumnWithAlias` by
default. Passing `raw=true` returns the raw column name:

```go
schemas.User.Name()     // e.g. "u.name"
schemas.User.Name(true) // "name"
```

If a model field would collide with `Table`, `Alias`, `As`, or `Columns`, the
generated accessor is prefixed with `Col`, for example `ColTable`. Generated
struct-field identifiers that would be Go keywords are prefixed with `__`.

## `export-api`

`export-api` describes the application's API surface as data. The manifest is
produced by the **application itself**: the command runs the main package with
`VEF_EXPORT_API` set, so what it reports is what the binary really registers —
resource names, actions, auth strategies, permission tokens, audit flags,
effective timeouts and rate limits, and the full shape of each request payload —
rather than what a source scan guesses.

| Flag | Default | Purpose |
| --- | --- | --- |
| `--app` | `./cmd/server` | main package of the application to describe |
| `--output`, `-o` | `api-manifest.json` | file to write the manifest to, or `-` for stdout |
| `--check` | `false` | fail if the existing file differs instead of rewriting it |

```bash
vef-cli export-api --app ./cmd/server -o api-manifest.json
vef-cli export-api --app ./cmd/server -o api-manifest.json --check
```

The output is sorted and carries no timestamp, so it is meant to be
**committed**: its diff is the diff of your API contract, and `--check` is the
CI gate that keeps the two in step. Two rules make the command composable, and
both were learned the hard way: progress goes to stderr, never stdout, so
`-o -` can be piped into `jq`; and `--check` with `-o -` is refused rather than
ignored, because a CI job written that way would compare nothing and exit 0 on
every drift.

### The manifest shape

```json
{
  "framework": "v0.51.0",
  "resources": [
    {
      "name": "security/auth",
      "kind": "rpc",
      "version": "v1",
      "operations": [
        {
          "action": "login",
          "auth": "none",
          "timeoutMs": 30000,
          "rateLimit": { "max": 6, "periodMs": 300000 },
          "params": "github.com/…/security.LoginParams"
        }
      ]
    }
  ],
  "types": {
    "github.com/…/security.LoginParams": {
      "fields": [
        { "name": "type", "type": "string", "validate": "required" }
      ]
    }
  }
}
```

A field carries `name` (from its json tag), `type`, and the optional
`optional` / `label` / `validate`; `permission` and `audit` appear on an
operation only when set. A named struct is rendered as the key it occupies in
`types`, so a payload graph is described once and referenced by name.

Result types are deliberately absent: crud handlers return only `error` and
write the response through the context, so the model type is nowhere in the
signature. Closing that gap needs one appended field on `api.OperationSpec`,
which is the planned next step rather than an oversight — the same gap
swallows `delete` / `delete_many`, whose handler takes the untyped parameter
bag.

### Exporting from your own entry point

`vef.Run` delegates to `vef.ExportAPI` when the `vef.EnvExportAPI`
(`VEF_EXPORT_API`) environment variable names a destination, so a host needs no
second entry point. Call it directly when you want the manifest in-process:

```go
if err := vef.ExportAPI(os.Stdout, options...); err != nil {
    return err
}
```

It is an environment variable rather than a flag because `Run` does not own the
process's flag set — a host's `main` may already define its own.

**Nothing starts and nothing connects.** Resources register during
construction, the data source opens lazily, and the operation mount that
invokes a crud handler's factory only runs afterwards — which is what makes the
export safe in a container build or a CI job. It is meant for a process that
exits afterwards: a few constructors start their own goroutines instead of
registering a lifecycle hook (the in-memory caches behind the session store and
the login guard each run a GC ticker), and since the graph is never started it
is never stopped either, so calling `ExportAPI` repeatedly inside a
long-running process accumulates them.

The surface is read through `api.EngineInspector`, an **optional** interface
(`Operations() []*api.Operation`, ordered by identifier) that the framework's
own engine implements. It is an inspector rather than a method on `api.Engine`
— matching `event.StreamInspector` and its siblings — so adding introspection
breaks no host that implements `Engine` itself. An engine that does not
implement it yields `vef.ErrEngineNotInspectable`, which means something
replaced the framework's own.

## Common `go:generate` pattern

In real VEF apps, these commands are often placed directly above `module.go`:

```go
//go:generate vef-cli generate-model-schema -i ./models -o ./schemas -p schemas
package sys
```

and for framework-facing build metadata:

```go
//go:generate vef-cli generate-build-info -o ./build_info.go -p vef
package vef
```

That keeps schema helpers and build metadata physically close to the module that uses them.

## Where each command belongs

| Command | Owns its output? | Use it for |
| --- | --- | --- |
| `new project` / `new resource` / `new service` | no — written once, then yours | starting a project, adding a CRUD resource or a service |
| `generate-build-info` | yes — regenerated | build metadata surfaced through `sys/monitor` |
| `generate-model-schema` | yes — regenerated | schema helpers derived from your models |
| `export-api` | yes — regenerated | the committed API manifest and its CI drift gate |

Blurring that line is how a generator becomes a command nobody dares to run.

## Next step

Read [Monitor](../infrastructure/monitor) if you want the generated build info to show up through `sys/monitor`.
