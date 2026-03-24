# VEF Framework Docs Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 基于 `vef-framework-go` 当前源码与测试，替换文档站点中的占位页面，产出一套面向使用者、与实现一致、可持续维护的官方文档。

**Architecture:** 以“源码与测试是真相，README 仅作历史线索”为原则重建文档。信息架构按开发者上手路径组织：启动应用、定义模型、注册资源、处理请求、配置安全、扩展功能、查询参考。文档正文优先使用来自 `bootstrap.go`、`di.go`、`api/`、`crud/`、`orm/`、`security/`、`internal/api/*` 与集成测试的真实用法。

**Tech Stack:** Docusaurus 3.9.2、Markdown/MDX、TypeScript 站点壳、VEF Framework Go v0.20.0 源码、Go 测试用例。

---

## Source-Backed Findings

- 当前文档站点的正文几乎全部是占位内容，现有价值主要在目录骨架、首页和 Docusaurus 配置，而不是页面内容。
- 框架的核心心智模型不是“工具包合集”，而是“基于 Uber FX + Fiber 的 API 资源框架”：用户通过 `vef.Run(...)` 组合模块，通过 `vef.ProvideAPIResource(...)` 注册资源，通过 `api.NewRPCResource(...)` / `api.NewRESTResource(...)` 暴露接口。
- 请求生命周期分两层：
  - 全局 Fiber 中间件：压缩、响应头、CORS、Content-Type、Request ID、日志、恢复、请求记录、SPA。
  - API 中间件：鉴权、上下文注入、数据权限、限流、审计。
- API 默认行为由源码明确给出：
  - 默认版本 `v1`
  - 默认认证 `Bearer`
  - 默认超时 `30s`
  - 默认限流 `100 requests / 5 minutes`
- RPC 与 REST 都是一级能力，不应把 REST 只当作补充说明：
  - RPC 入口固定为 `POST /api`
  - REST 入口为 `/api/<resource>`
  - REST 支持 path/query/body 合并到 `params`，并通过 `X-Meta-*` 请求头传递 `meta`
- CRUD 是用户最常用的编程接口，且真实能力比当前骨架暗示的更强：除了 `Create/Update/Delete/FindPage/FindAll`，还包括 `FindOne/FindOptions/FindTree/FindTreeOptions/Export/Import/CreateMany/UpdateMany/DeleteMany`，并支持链式配置、前后置 Hook、事务、数据权限、文件提升。
- handler 注入模型是文档必须讲清的重点：框架会自动注入 `fiber.Ctx`、`orm.DB`、`log.Logger`、`*security.Principal`、`event.Publisher`、`cron.Scheduler`、`mold.Transformer`、`storage.Service`、`api.Params`、`api.Meta`、嵌入 `api.P` / `api.M` 的结构体，以及 handler factory 的启动期依赖。
- 安全能力不只是“JWT 鉴权”：
  - 默认内置 `security/auth` 资源
  - 支持 Bearer、Signature、Public 三种 API 鉴权策略
  - 登录支持 challenge 流程
  - RBAC 权限检查与数据权限解析均通过扩展接口完成
- 文档骨架遗漏了几类默认可见能力：
  - `sys/storage`
  - `sys/schema`
  - `sys/monitor`
  - MCP 集成
  - Sequence/CQRS/SPA 集成等扩展点
- 像 `approval`、`cryptox`、`password`、`decimal`、`timex` 这类公共包虽然存在，但不应直接进入第一版“框架使用文档”主线；除非后续明确要扩展为“VEF 全家桶文档”，否则第一版只在参考页或 API Reference 链接中提及。

## Documentation IA Decisions

- 保留现有主分类，但重写其定位：
  - `getting-started/`：安装、最小应用、配置、项目结构
  - `modules/`：运行时装配、生命周期、模块组合与 DI
  - `guide/`：日常开发主线
  - `security/`：认证、授权、数据权限
  - `features/`：内置基础设施能力
  - `advanced/`：扩展点、非默认路径、跨模块能力
- 在 `guide/` 下新增并前置一个页面，专门解释参数注入与 `api.P/api.M`，因为这是理解 handler 签名的关键。
- 在 `features/` 下补出当前源码已默认启用、但文档骨架缺失的页面：
  - `monitor`
  - `schema`
  - `mcp`
  - 视篇幅决定是否加入 `sequence`
- 新增 `reference/` 分类，承接不适合做成长教程、但开发者经常需要查阅的内容：
  - 配置字段表
  - 内置资源与默认端点
  - 扩展点总览
- 第一批正文优先级按“用户第一次真正写业务代码的顺序”排列，而不是按包名罗列。

## Writing Rules For This Doc Set

- 每一页都要同时回答三件事：
  - 这个能力解决什么问题
  - 用户最少要写什么代码
  - 框架默认替你做了什么
- 每一页至少包含一段来自真实源码模式的例子，优先取材于：
  - `internal/api/*_test.go`
  - `crud/*_test.go`
  - `internal/security/auth_resource.go`
  - `internal/storage/storage_resource.go`
  - `internal/schema/schema_resource.go`
  - `internal/monitor/monitor_resource.go`
- 页面中凡是涉及默认值、顺序、限流、超时、路径、header、配置键名，必须以源码常量或构造函数为准，不引用 README 旧描述。
- 第一版不追求把所有 public package 都解释完，优先覆盖“搭应用会碰到的 API”。

### Task 1: Build The Narrative Foundation

**Files:**
- Modify: `docs/intro.md`
- Modify: `docs/getting-started/installation.md`
- Modify: `docs/getting-started/quick-start.md`
- Modify: `docs/getting-started/configuration.md`
- Modify: `docs/getting-started/project-structure.md`

**Step 1: Rewrite `docs/intro.md`**

写清框架定位：FX 模块装配、Fiber HTTP 服务、统一 API 资源模型、内置认证/存储/监控等能力。

**Step 2: Rewrite `docs/getting-started/installation.md`**

说明 Go 版本、安装命令、配置文件搜索路径、最小必需配置项，避免引用过时 README 里的不准确信息。

**Step 3: Rewrite `docs/getting-started/quick-start.md`**

提供一个真正可运行的最小示例：`main.go` + `application.toml` + 一个简单 RPC 资源。

**Step 4: Rewrite `docs/getting-started/configuration.md`**

按 `vef.app`、`vef.data_source`、`vef.security`、`vef.storage`、`vef.mcp` 等小节解释配置。

**Step 5: Rewrite `docs/getting-started/project-structure.md`**

不要再照搬 README 的旧结构；改为文档化推荐结构，并明确“推荐组织方式”和“框架强制要求”的区别。

**Step 6: Verify**

Run: `pnpm build`
Expected: 站点构建成功，无新的 markdown 解析错误。

### Task 2: Explain Runtime Assembly And Request Lifecycle

**Files:**
- Modify: `docs/modules/overview.md`
- Modify: `docs/modules/lifecycle.md`
- Modify: `docs/guide/routing.md`
- Create: `docs/guide/params-and-meta.md`

**Step 1: Rewrite `docs/modules/overview.md`**

解释 `vef.Run(...)` 默认装配的模块顺序，以及 `vef.ProvideAPIResource(...)`、`vef.ProvideMiddleware(...)`、`vef.ProvideCQRSBehavior(...)` 等入口。

**Step 2: Rewrite `docs/modules/lifecycle.md`**

说明 app 启动/停止、FX 生命周期、数据库与事件总线等模块在启动期做什么。

**Step 3: Rewrite `docs/guide/routing.md`**

分别解释 RPC 和 REST 的路由模型、请求格式、默认路径、handler 解析规则和 action 命名规则。

**Step 4: Create `docs/guide/params-and-meta.md`**

说明 `api.P`、`api.M`、`page.Pageable`、`api.Params`、`api.Meta`、`X-Meta-*` header、multipart 文件字段如何被绑定。

**Step 5: Add lifecycle diagram**

在 `docs/modules/lifecycle.md` 或 `docs/guide/routing.md` 中加入请求生命周期图：
全局中间件 -> Router 解析 -> Auth -> Contextual -> DataPermission -> RateLimit -> Audit -> Handler。

**Step 6: Verify**

Run: `pnpm build`
Expected: 新增页面出现在 sidebar 中，内部链接无 broken link。

### Task 3: Document Models, Search, CRUD, And Hooks

**Files:**
- Modify: `docs/guide/models.md`
- Modify: `docs/guide/crud.md`
- Modify: `docs/guide/query-builder.md`
- Modify: `docs/guide/hooks.md`
- Modify: `docs/advanced/transactions.md`

**Step 1: Rewrite `docs/guide/models.md`**

说明 `bun.BaseModel`、`orm.Model/IDModel/CreatedModel/AuditedModel`、审计字段、主键规则、JSON/validate/search 标签的协作方式。

**Step 2: Rewrite `docs/guide/query-builder.md`**

解释 `search:"..."` 标签、`search.Applier[T]()`、`page.Pageable` 和 `crud.Sortable` 的联动关系。

**Step 3: Rewrite `docs/guide/crud.md`**

按“定义资源结构体 -> 嵌入 CRUD provider -> `crud.NewXxx()` 链式配置 -> 注册资源”这个顺序写。
必须覆盖：
- `Create`
- `Update`
- `Delete`
- `FindPage`
- `FindAll`
- `FindOne`
- `FindOptions`
- `FindTree`
- `Export`
- `Import`
- Many variants

**Step 4: Rewrite `docs/guide/hooks.md`**

说明 `WithPreCreate`、`WithPostCreate` 等 hook 的事务边界、典型用途和失败行为。

**Step 5: Rewrite `docs/advanced/transactions.md`**

说明 `db.RunInTX(...)`、CRUD 默认事务、何时自己开事务、事务里如何继续使用 `orm.DB`。

**Step 6: Verify**

Run: `pnpm build`
Expected: CRUD、模型、查询相关页面示例代码可读，构建通过。

### Task 4: Document Custom Handlers And Extension Signatures

**Files:**
- Modify: `docs/guide/custom-handlers.md`
- Modify: `docs/advanced/context-helpers.md`
- Create: `docs/advanced/custom-param-resolvers.md`

**Step 1: Rewrite `docs/guide/custom-handlers.md`**

讲清三件事：
- RPC 默认按 `Action -> PascalCase` 解析方法
- REST 必须显式提供 handler
- handler 可以是普通函数，也可以是 factory function

**Step 2: Add signature matrix**

在 `docs/guide/custom-handlers.md` 中加入可注入参数矩阵：
`fiber.Ctx`、`orm.DB`、`*security.Principal`、`log.Logger`、`event.Publisher`、`cron.Scheduler`、`mold.Transformer`、`storage.Service`、`api.Params`、`api.Meta`、`api.P`/`api.M` 结构体。

**Step 3: Rewrite `docs/advanced/context-helpers.md`**

把 `contextx` 放回真实语境：何时应该直接依赖 handler 参数注入，何时需要从 context 手动取值。

**Step 4: Create `docs/advanced/custom-param-resolvers.md`**

说明如何通过 `group:"vef:api:handler_param_resolvers"` 和 `group:"vef:api:factory_param_resolvers"` 扩展注入。

**Step 5: Verify**

Run: `pnpm build`
Expected: 自定义 handler 与扩展点页面均能正常出现在站点。

### Task 5: Rebuild Security Documentation Around Actual Flows

**Files:**
- Modify: `docs/security/authentication.md`
- Modify: `docs/security/authorization.md`
- Modify: `docs/security/data-permissions.md`
- Create: `docs/reference/built-in-auth-resource.md`

**Step 1: Rewrite `docs/security/authentication.md`**

解释 `Public/Bearer/Signature` 三种 API 认证策略、默认行为、相关 header 和登录资源的角色。

**Step 2: Add login flow section**

基于 `internal/security/auth_resource.go` 写出 `login -> challenge -> refresh -> logout -> get_user_info` 的顺序和返回结构。

**Step 3: Rewrite `docs/security/authorization.md`**

说明 `PermToken`、`PermissionChecker`、`RolePermissionsLoader` 的职责分界，强调框架不内置你的角色来源。

**Step 4: Rewrite `docs/security/data-permissions.md`**

说明 `DataPermissionResolver`、`DataScope`、`RequestScopedDataPermApplier` 与 CRUD 查询的关系。

**Step 5: Create `docs/reference/built-in-auth-resource.md`**

整理 `security/auth` 资源的 actions、入参、限流、认证要求。

**Step 6: Verify**

Run: `pnpm build`
Expected: 安全相关页面链接正确，概念与 built-in resource 页面互相可跳转。

### Task 6: Fill In Built-In Infrastructure Features

**Files:**
- Modify: `docs/features/cache.md`
- Modify: `docs/features/event-bus.md`
- Modify: `docs/features/cron.md`
- Modify: `docs/features/storage.md`
- Modify: `docs/features/i18n.md`
- Create: `docs/features/monitor.md`
- Create: `docs/features/schema.md`
- Create: `docs/features/mcp.md`

**Step 1: Rewrite `docs/features/storage.md`**

覆盖 `storage.Service`、临时文件上传、预签名 URL、文件提升、`sys/storage` 资源。

**Step 2: Rewrite `docs/features/cache.md`**

覆盖 `cache.NewMemory[T]()`、`cache.NewRedis[T]()`、TTL、命名空间、singleflight mixin 场景。

**Step 3: Rewrite `docs/features/event-bus.md`**

解释内存事件总线、发布/订阅、事件中间件组以及 audit/login/file 等事件的文档示例。

**Step 4: Rewrite `docs/features/cron.md`**

说明 `cron.Scheduler` 注入、job 定义器、任务注册方式与适用边界。

**Step 5: Rewrite `docs/features/i18n.md`**

说明默认语言、消息键、`i18n.T(...)` 与验证/错误响应的关系。

**Step 6: Create `docs/features/monitor.md`, `docs/features/schema.md`, `docs/features/mcp.md`**

这些页面至少要回答：
- 是否默认启用
- 暴露了什么 endpoint / middleware
- 用户如何扩展
- 适合什么场景

**Step 7: Verify**

Run: `pnpm build`
Expected: features 分类不再只有占位页，新增页面能被 sidebar 自动收录。

### Task 7: Add Advanced Integration And Reference Layers

**Files:**
- Modify: `docs/advanced/cli-tools.md`
- Create: `docs/advanced/spa-integration.md`
- Create: `docs/advanced/cqrs.md`
- Create: `docs/reference/_category_.json`
- Create: `docs/reference/configuration-reference.md`
- Create: `docs/reference/built-in-resources.md`
- Create: `docs/reference/extension-points.md`

**Step 1: Rewrite `docs/advanced/cli-tools.md`**

如实说明当前 CLI 的状态，不夸大“脚手架能力”；把已实现与未实现分开写。

**Step 2: Create `docs/advanced/spa-integration.md`**

解释 `middleware.SPAConfig`、`vef.ProvideSPAConfig(...)`、`vef.SupplySPAConfigs(...)`、SPA fallback 和排除路径。

**Step 3: Create `docs/advanced/cqrs.md`**

说明 `cqrs.Register(...)`、`cqrs.Send(...)`、`vef.ProvideCQRSBehavior(...)` 的接入方式。

**Step 4: Create `docs/reference/configuration-reference.md`**

列出 `vef.app`、`vef.data_source`、`vef.security`、`vef.redis`、`vef.storage`、`vef.monitor`、`vef.mcp` 等配置字段。

**Step 5: Create `docs/reference/built-in-resources.md`**

收拢 `security/auth`、`sys/storage`、`sys/schema`、`sys/monitor` 等默认资源入口。

**Step 6: Create `docs/reference/extension-points.md`**

列出最关键 DI group 和 helper：
- `vef:api:resources`
- `vef:app:middlewares`
- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`
- `vef:cqrs:behaviors`
- `vef:mcp:*`

**Step 7: Verify**

Run: `pnpm build`
Expected: reference 分类出现且结构清晰。

### Task 8: Align Homepage And Navigation With Real Capabilities

**Files:**
- Modify: `src/pages/index.tsx`
- Modify: `i18n/zh-Hans/code.json`
- Modify: `docs/intro.md`

**Step 1: Update homepage feature copy**

把首页文案从泛泛描述改成与真实能力一致的表述，特别是：
- 统一资源模型
- RPC/REST 双路由
- 泛型 CRUD
- 安全与数据权限
- 内置 sys 资源和扩展点

**Step 2: Update Chinese UI translations**

同步首页 feature 文案的中文翻译，避免英文准确、中文仍是旧说法。

**Step 3: Verify**

Run: `pnpm build`
Expected: 首页、文档导航与正文表述一致。

## Recommended Execution Order

1. Task 1
2. Task 2
3. Task 3
4. Task 4
5. Task 5
6. Task 6
7. Task 7
8. Task 8

## Scope Guardrails

- 第一轮不做以下事情：
  - 把所有 utility package 都写成教程
  - 为 `approval` 建完整产品文档
  - 复制 README 的旧段落作为正文
  - 发明源码里没有的最佳实践
- 第一轮必须做到：
  - 让新用户能从 0 写出一个最小资源
  - 让现有用户能查到默认行为、默认资源、扩展点和请求生命周期
  - 所有关键默认值都有源码依据

## Verification Checklist

- `docs/` 中不再有只有标题的一屏占位页
- `pnpm build` 通过
- 站点首页与正文描述一致
- RPC、REST、CRUD、Security、Storage 五条主线都有完整可运行示例
- 至少一页 reference 文档明确列出默认资源与配置字段
