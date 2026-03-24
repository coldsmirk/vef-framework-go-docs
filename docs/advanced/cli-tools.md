---
sidebar_position: 1
---

# CLI Tools

VEF includes a CLI, but its current scope is narrower than a full project scaffolder.

## Current commands

The CLI currently registers these commands:

- `create`
- `generate-build-info`
- `generate-model-schema`

## Minimal command examples

```bash
vef-cli generate-build-info -o internal/vef/build_info.go -p vef
vef-cli generate-model-schema -i models -o schemas -p schemas
```

## Important reality check

`vef-cli create` exists as a command, but it is currently **not implemented**.

Do not treat it as a working project generator yet.

## `generate-build-info`

This command generates a Go source file containing build metadata such as:

- app version
- build time
- git commit

It is designed to be used from `go:generate` or from your build pipeline.

## `generate-model-schema`

This command inspects model files and generates type-safe schema helpers for ORM usage.

It supports:

- file-to-file generation
- directory-to-directory generation

The goal is to reduce hard-coded column-name strings in query code.

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

Read [Monitor](../features/monitor) if you want the generated build info to show up through `sys/monitor`.
