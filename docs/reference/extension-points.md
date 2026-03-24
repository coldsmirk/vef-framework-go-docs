---
sidebar_position: 3
---

# Extension Points

Most VEF extension points are explicit FX groups.

## API and app groups

- `vef:api:resources`
- `vef:app:middlewares`

Helpers:

- `vef.ProvideAPIResource(...)`
- `vef.ProvideMiddleware(...)`

## Minimal module example

```go
var Module = vef.Module(
  "app:user",
  vef.ProvideAPIResource(NewUserResource),
  vef.ProvideMiddleware(NewAuditMiddleware),
)
```

## API parameter injection

- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`

These extend request-time and startup-time handler injection.

## CQRS

- `vef:cqrs:behaviors`

Helper:

- `vef.ProvideCQRSBehavior(...)`

## Security

- `vef:security:challenge_providers`

Helper:

- `vef.ProvideChallengeProvider(...)`

## SPA

- `vef:spa`

Helpers:

- `vef.ProvideSPAConfig(...)`
- `vef.SupplySPAConfigs(...)`

## MCP

- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`

Helpers:

- `vef.ProvideMCPTools(...)`
- `vef.ProvideMCPResources(...)`
- `vef.ProvideMCPResourceTemplates(...)`
- `vef.ProvideMCPPrompts(...)`
- `vef.SupplyMCPServerInfo(...)`

## See also

- [Modules & Dependency Injection](../modules/overview) for how these groups fit into app composition
- [Custom Param Resolvers](../advanced/custom-param-resolvers) for handler injection extension
