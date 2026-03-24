---
sidebar_position: 1
---

# Installation

This page covers the minimum environment and project setup required to boot a VEF application.

## Requirements

The current framework version requires:

- Go `1.26.0`

For application development, Go is the only mandatory runtime dependency.

## Runtime prerequisites

If you use the default `vef.Run(...)` boot path, the current framework runtime always boots the database module.

That means the true minimum startup prerequisite is:

- a reachable database
- a valid `application.toml`

The default boot graph also includes the Redis module, but Redis only becomes a practical prerequisite when your application or an enabled feature actually consumes `*redis.Client` or another Redis-backed capability.

For the smallest local setup, use SQLite.

## Add the framework

Install the package in your Go module:

```bash
go get github.com/coldsmirk/vef-framework-go
```

If you are starting from an empty directory:

```bash
go mod init example.com/my-app
go get github.com/coldsmirk/vef-framework-go
```

## Pick a database early

VEF boots its database module during startup, so your application configuration needs a valid data source from the start. The framework supports:

- PostgreSQL
- MySQL
- SQLite

For local exploration and small demos, SQLite is the simplest choice because it avoids external infrastructure.

## Add Redis only when needed

If you later enable Redis-backed features, the default client settings are:

- host: `127.0.0.1`
- port: `6379`
- network: `tcp`

One easy local option is:

```bash
docker run --name vef-redis -p 6379:6379 -d redis:7-alpine
```

## Create the config file

By default VEF looks for `application.toml` in:

- `./configs`
- `$VEF_CONFIG_PATH`
- `.`
- `../configs`

The most common layout is:

```text
my-app/
├── configs/
│   └── application.toml
└── main.go
```

## Minimal configuration

This is enough to boot an application with SQLite and the default in-memory storage provider:

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_source]
type = "sqlite"
```

If you omit `vef.storage.provider`, the framework falls back to memory storage.
Add `vef.redis` only when the application really uses Redis-backed features.

## What happens during startup

When you call `vef.Run(...)`, the framework initializes configuration, database access, ORM, middleware, API routing, security, events, CQRS, cron, Redis, mold, storage, sequence, schema, monitoring, MCP, and finally the HTTP application.

That is why installation in VEF is not just “import the package”. A valid config file is part of installation.

## Optional environment variables

The following environment variables are especially useful during setup:

- `VEF_CONFIG_PATH`: add an extra config search directory
- `VEF_LOG_LEVEL`: adjust log verbosity
- `VEF_NODE_ID`: provide a node identifier for distributed ID scenarios
- `VEF_I18N_LANGUAGE`: switch the framework language, defaulting to Simplified Chinese

## Next step

Continue to [Quick Start](./quick-start) to build a minimal resource and confirm the app is serving requests.
