# 扩展点

VEF 的扩展核心是 FX group。大多数框架级定制都不是通过修改运行时本身完成，而是通过把组件注册到合适的 group 中完成。

## API 资源

使用：

```go
vef.ProvideAPIResource(...)
```

对应 FX group：

```text
vef:api:resources
```

## 应用级 middleware

使用：

```go
vef.ProvideMiddleware(...)
```

对应 FX group：

```text
vef:app:middlewares
```

## SPA 配置

使用：

```go
vef.ProvideSPAConfig(...)
vef.SupplySPAConfigs(...)
```

对应 FX group：

```text
vef:spa
```

## CQRS behavior

使用：

```go
vef.ProvideCQRSBehavior(...)
```

对应 FX group：

```text
vef:cqrs:behaviors
```

## 安全 challenge provider

使用：

```go
vef.ProvideChallengeProvider(...)
```

对应 FX group：

```text
vef:security:challenge_providers
```

## MCP provider

使用：

```go
vef.ProvideMCPTools(...)
vef.ProvideMCPResources(...)
vef.ProvideMCPResourceTemplates(...)
vef.ProvideMCPPrompts(...)
vef.SupplyMCPServerInfo(...)
```

对应 FX group：

- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`

## API 参数注入解析器

这是更进阶的扩展点：

- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`

当内置 handler 参数集合不够时，可以往这里注册自定义 resolver。

## 事件中间件

对应 FX group：

```text
vef:event:middlewares
```

适合承载所有事件的统一逻辑。

## 一个简单判断原则

当你想扩展 VEF 时，优先判断：

- 框架是否已经为这个概念预留了 group
- 你的扩展是否应该纳入启动和生命周期管理
- 你是否希望模块之间依赖保持显式和可测试

如果答案是肯定的，那通常就应该走 FX group，而不是用隐式全局状态或手写单例。

## 延伸阅读

- [模块与依赖注入](../modules/overview)：这些 group 如何进入应用装配流程
- [自定义参数解析器](../advanced/custom-param-resolvers)：handler 注入扩展的具体做法
