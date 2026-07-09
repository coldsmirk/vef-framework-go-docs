---
sidebar_position: 3
---

# CLI Tools

VEF includes a CLI, but its current scope is narrower than a full project scaffolder.

## Current commands

The root command is `vef-cli`. `vef-cli --version` prints the CLI banner plus
`Version: ...`; when build-date metadata is available it also prints
`Built: ...`, and dirty VCS builds append `-dirty` to the version string.

The CLI currently registers these subcommands:

- `create`
- `generate-build-info`
- `generate-model-schema`

## Minimal command examples

```bash
vef-cli --version
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
vef-cli generate-model-schema -i models -o schemas -p schemas
```

Application code should consume the CLI through these commands instead of
importing the `cmd/vef-cli/cmd/*` implementation packages directly.

## Important reality check

`vef-cli create` exists as a command, but it is currently **not implemented**.

Do not treat it as a working project generator yet.

The command returns this error:

```text
vef-cli create is not implemented yet, please generate the project manually
```

The command still defines these flags because the planned command shape exists:

| Flag | Default | Purpose |
| --- | --- | --- |
| `--name`, `-n` | required | project name |
| `--path`, `-p` | `.` | directory path where the project would be created |
| `--module`, `-m` | empty | Go module path |

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
only `*.go` files directly inside the input directory; it is not recursive. A
directory input may point at an existing output directory or a directory path
that does not exist yet. If the output path already exists as a file,
directory-to-file generation is rejected.

The generator reads structs in the target file that embed `orm.BaseModel`.
Table metadata comes from the embedded `orm.BaseModel` field's `bun` tag:
`table:...` sets the table name and `alias:...` sets the default alias. Without
those tag parts, the table defaults to the pluralized snake_case model name and
the alias defaults to the singular snake_case model name.

Field handling is source-compatible with these rules:

- only exported fields generate accessors
- `bun:"-"` fields are skipped
- `bun:"rel:*"` and `bun:"m2m:*"` relationship fields are skipped
- a first `bun` tag component such as `bun:"user_name"` sets the column name
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

## Recommended expectation

Today, the CLI is best treated as:

- a helper for build metadata generation
- a helper for model schema generation

It is **not** yet the right foundation for onboarding docs that promise one-command project scaffolding.

## Next step

Read [Monitor](../infrastructure/monitor) if you want the generated build info to show up through `sys/monitor`.
