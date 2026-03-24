# VEF Framework Docs Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 把当前 Docusaurus 骨架升级为与 `vef-framework-go` 源码实现一致、可持续维护、可双语发布的官方文档站点。

**Architecture:** 以 `../vef-framework-go` 源码为唯一事实来源，以 `README.md` / `README.zh-CN.md` 作为第一批内容母本，把单体 README 拆分为任务导向页面；保留现有站点的大体入口，补充缺失的 `reference` 与 `integrations` 类别，并把尚未进入 README 的能力明确标记为 `Preview` 或第二阶段内容。

**Tech Stack:** Docusaurus 3.9.2, React 19, TypeScript, Shiki 4, pnpm, Docusaurus i18n (`en`, `zh-Hans`)

---

## Analysis Snapshot

### Relationship

- `vef-framework-go-docs` 是独立的 Docusaurus 文档站仓库。
- `vef-framework-go` 是真实源码仓库，当前文档站应当以它为事实来源。
- 当前站点已具备主页、导航、分类和多语言框架，但内容基本仍是占位骨架。

### Current Documentation State

- `docs/` 下除 `intro.md` 外，几乎所有页面都是 7 行左右的占位文本，没有真正源码级说明。
- 当前侧边栏完全由目录自动生成，因此新增内部计划文件到 `docs/` 会被直接发布；执行计划应保存在仓库根目录 `plans/`。
- `zh-Hans` 已启用，但当前只有 `i18n/zh-Hans/docusaurus-plugin-content-docs/current.json` 这类 UI 翻译文件，没有文档正文翻译目录。

### Source-Backed First-Wave Topics

这些主题已经有较完整的源码与 README 支撑，应优先进入正式文档：

- 应用启动与模块装配：`vef.Run`, `di.go`, `bootstrap.go`, `start.go`
- API 设计：`api`, `crud`, `result`, `page`, `search`
- 数据访问：`orm`, `db.RunInTx(...)`, 搜索标签、分页、树查询
- 安全：`security`, JWT, signature auth, password auth, RBAC, data scopes
- 内置能力：`cache`, `event`, `cron`, `storage`, `validator`, `contextx`, `i18n`
- 开发工具：`cmd/vef-cli`, `generate-build-info`, `generate-model-schema`

### Source Topics Missing From Current IA

这些能力在源码中已有公开包，但当前站点没有入口页：

- `monitor`
- `sequence`
- `mcp`
- `approval`
- `ai` / `ai/stream`
- 参考型内容：RPC payload 结构、可注入参数、配置字段表、CRUD 动作矩阵、搜索标签矩阵

### Known Mismatches To Fix Early

- 首页文案把 RPC 描述成“自动生成 RESTful endpoints”，源码实际上是 `api.NewRPCResource(...)` 和 `api.NewRESTResource(...)` 两套资源类型，不应误导。
- `docs/advanced/cli-tools.md` 的占位说明提到“scaffold resources, models, and more”，但源码里的 `vef-cli create` 目前明确返回 `not implemented`，该页必须收敛到真实支持的命令。

### Assumptions

- 第一阶段先把稳定核心能力文档化，再决定是否提升 `approval` / `ai` 到首页级入口。
- 英文文档作为 source of truth；中文在英文稳定后通过 Docusaurus 翻译流程生成并补译。
- 尽量保留现有 URL 路径与分类目录，减少站点结构抖动和未来重定向成本。

## Recommended Information Architecture

保留现有分类，并做最小增量扩展：

- `docs/intro.md`
- `docs/getting-started/*`
- `docs/guide/*`
- `docs/modules/*`
- `docs/security/*`
- `docs/features/*`
- `docs/advanced/*`
- `docs/reference/*` 新增
- `docs/integrations/*` 新增

推荐页面编排：

- `getting-started`
  - `installation.md`
  - `quick-start.md`
  - `project-structure.md`
  - `configuration.md`
- `guide`
  - `models.md`
  - `routing.md`
  - `crud.md`
  - `custom-handlers.md`
  - `query-builder.md`
  - `hooks.md`
  - `validation.md`
  - `error-handling.md`
- `modules`
  - `overview.md`
  - `boot-sequence.md` 新增
  - `request-lifecycle.md` 新增
  - `lifecycle.md`
- `security`
  - `authentication.md`
  - `authorization.md`
  - `data-permissions.md`
  - `login-challenges.md` 新增
- `features`
  - `cache.md`
  - `event-bus.md`
  - `cron.md`
  - `storage.md`
  - `i18n.md`
  - `monitoring.md` 新增
- `advanced`
  - `transactions.md`
  - `context-helpers.md`
  - `cli-tools.md`
- `reference`
  - `rpc-request-format.md`
  - `handler-injection.md`
  - `search-tags.md`
  - `configuration-reference.md`
  - `crud-operation-matrix.md`
  - `result-and-error-model.md`
- `integrations`
  - `mcp.md`
  - `sequence.md`
  - `approval.md`
  - `ai.md`

## Content Strategy

### Recommended Approach

采用“先拆 README，再补源码缺口”的策略，而不是直接从零写整站：

- 先把 `../vef-framework-go/README.md` 和 `../vef-framework-go/README.zh-CN.md` 中已经成型的教程段落拆到对应页面。
- 再用源码把 README 没覆盖到的细节补成 reference 页面。
- 最后再处理 `approval` / `mcp` / `ai` / `sequence` 这类不在当前 README 首页承诺中的能力。

### Source Files To Use As Canonical Inputs

- `../vef-framework-go/README.md`
- `../vef-framework-go/README.zh-CN.md`
- `../vef-framework-go/bootstrap.go`
- `../vef-framework-go/start.go`
- `../vef-framework-go/di.go`
- `../vef-framework-go/api/*.go`
- `../vef-framework-go/crud/*.go`
- `../vef-framework-go/orm/*.go`
- `../vef-framework-go/search/*.go`
- `../vef-framework-go/security/*.go`
- `../vef-framework-go/cache/*.go`
- `../vef-framework-go/event/*.go`
- `../vef-framework-go/cron/*.go`
- `../vef-framework-go/storage/*.go`
- `../vef-framework-go/contextx/*.go`
- `../vef-framework-go/cmd/vef-cli/cmd/**/*.go`

## Task 1: Lock IA And Correct Site Messaging

**Files:**
- Modify: `docusaurus.config.ts`
- Modify: `src/pages/index.tsx`
- Modify: `docs/modules/_category_.json`
- Modify: `docs/advanced/_category_.json`
- Create: `docs/reference/_category_.json`
- Create: `docs/integrations/_category_.json`

**Step 1:** 调整首页 copy，使其描述“RPC 和 REST 并列支持”，不要暗示 RPC 自动生成 REST。

**Step 2:** 在首页或首屏下方增加 “Getting Started / Guide / Security / Reference” 的入口，降低第一次进入的路径选择成本。

**Step 3:** 新增 `reference` 与 `integrations` 分类，保留现有大类 URL，不做大规模改名。

**Step 4:** 校正 navbar / footer 文案，使入口指向真实存在且优先级最高的页面。

**Step 5:** 运行验证。

Run: `pnpm build`

Expected: 构建成功，无 broken links，无新分类导致的侧边栏错误。

## Task 2: Replace Placeholders In Intro And Getting Started

**Files:**
- Modify: `docs/intro.md`
- Modify: `docs/getting-started/installation.md`
- Modify: `docs/getting-started/quick-start.md`
- Modify: `docs/getting-started/project-structure.md`
- Modify: `docs/getting-started/configuration.md`

**Step 1:** 把 `README` 的安装、最小示例、目录结构、配置章节拆成独立页面。

**Step 2:** 在 `quick-start.md` 中明确最小运行路径：`go get` → `main.go` → `configs/application.toml` → `go run`.

**Step 3:** 在 `project-structure.md` 中给出推荐目录树，并解释 `internal/auth`, `internal/sys`, `internal/vef`, `internal/web` 的角色。

**Step 4:** 在 `configuration.md` 中只保留新手必需项，详细字段表移到 reference 页面，避免入门页变成配置字典。

**Step 5:** 运行验证。

Run: `pnpm build`

Expected: 入门页都有有效内容，且首页、footer、sidebar 指向的页面都能正常构建。

## Task 3: Write Core API Guide Pages

**Files:**
- Modify: `docs/guide/models.md`
- Modify: `docs/guide/routing.md`
- Modify: `docs/guide/crud.md`
- Modify: `docs/guide/custom-handlers.md`
- Modify: `docs/guide/query-builder.md`
- Modify: `docs/guide/hooks.md`
- Modify: `docs/guide/validation.md`
- Modify: `docs/guide/error-handling.md`

**Step 1:** `models.md` 说明 `orm.BaseModel`, `orm.Model`, 审计字段、`bun` / `json` / `validate` / `label` 标签的组合方式。

**Step 2:** `routing.md` 说明 RPC 与 REST 的资源命名、动作命名、请求路径、请求体结构和 `api.P` / `api.M` 绑定方式。

**Step 3:** `crud.md` 说明 `crud.NewFindAll`, `crud.NewFindPage`, `crud.NewCreate`, `crud.NewUpdate`, `crud.NewDelete` 及其链式 builder 配置。

**Step 4:** `custom-handlers.md` 说明 `api.WithOperations(...)`、RPC 方法名映射规则、REST `OperationSpec.Handler` 的显式绑定方式。

**Step 5:** `query-builder.md` 说明 `orm.ConditionBuilder`, `search` 标签, `search.Applier[T]`, 树查询 `QueryPart` 系统。

**Step 6:** `hooks.md` 说明 CRUD 的 `WithPre*` / `WithPost*`，区分单条、批量、导入导出钩子，并明确 `Post*` 在事务内运行。

**Step 7:** `validation.md` 说明 `validator.Validate`, 自定义 rule 与 `null.*` 类型支持。

**Step 8:** `error-handling.md` 说明 `result.Ok(...)`, `result.Err(...)`, HTTP status 与业务 code 的关系。

**Step 9:** 运行验证。

Run: `pnpm build`

Expected: `guide` 分类可独立支撑“如何定义一个资源”的完整学习路径。

## Task 4: Expand Modules Documentation Around Runtime Mental Model

**Files:**
- Modify: `docs/modules/overview.md`
- Create: `docs/modules/boot-sequence.md`
- Create: `docs/modules/request-lifecycle.md`
- Modify: `docs/modules/lifecycle.md`

**Step 1:** `overview.md` 解释 VEF 的核心心智模型：FX module + grouped providers + public package / internal module 分层。

**Step 2:** `boot-sequence.md` 明确 `config → database → orm → middleware → api → security → event → cqrs → cron → redis → mold → storage → sequence → schema → monitor → mcp → app`。

**Step 3:** `request-lifecycle.md` 解释 `/api` 请求从解析、认证、上下文注入、授权、限流、handler dispatch 到 response 的流程。

**Step 4:** `lifecycle.md` 说明 `vef.Lifecycle`, `vef.StartHook`, `vef.StopHook`, 以及资源清理模式。

**Step 5:** 运行验证。

Run: `pnpm build`

Expected: “框架如何启动、如何装配、一次请求如何流转” 可以在 `modules` 分类中闭环。

## Task 5: Write Security Pages Against Real Source Behavior

**Files:**
- Modify: `docs/security/authentication.md`
- Modify: `docs/security/authorization.md`
- Modify: `docs/security/data-permissions.md`
- Create: `docs/security/login-challenges.md`

**Step 1:** `authentication.md` 解释 JWT、signature auth、password auth 的适用场景与配置入口。

**Step 2:** 把 `security.UserLoader`, `ExternalAppLoader`, `PasswordDecryptor`, `NonceStore` 的实现点写清楚，不只写概念。

**Step 3:** `authorization.md` 说明 `PermToken`, `PermissionChecker`, `RolePermissionsLoader`, `CachedRolePermissionsLoader` 的配合关系。

**Step 4:** `data-permissions.md` 说明 `AllDataScope`, `SelfDataScope`, 优先级规则和自定义 `DataScope`。

**Step 5:** `login-challenges.md` 说明 `ChallengeProvider`, OTP, TOTP, department selection 这类多阶段登录挑战能力。

**Step 6:** 运行验证。

Run: `pnpm build`

Expected: 安全章节覆盖框架默认行为与常见扩展点，不再只停留在目录名级别。

## Task 6: Fill Built-In Features Pages

**Files:**
- Modify: `docs/features/cache.md`
- Modify: `docs/features/event-bus.md`
- Modify: `docs/features/cron.md`
- Modify: `docs/features/storage.md`
- Modify: `docs/features/i18n.md`
- Create: `docs/features/monitoring.md`

**Step 1:** `cache.md` 说明 `cache.NewMemory[T]`, `cache.NewRedis[T]`, TTL, eviction policy, key builder。

**Step 2:** `event-bus.md` 说明 `event.Bus`, `Publish`, `Subscribe`, middleware 和审计事件的订阅方式。

**Step 3:** `cron.md` 说明 `cron.NewCronJob`, `NewDurationJob`, `WithTask`, `WithTags`, `WithConcurrent` 等调度能力。

**Step 4:** `storage.md` 说明 `storage.Service`, `PutObject`, `GetPresignedURL`, `PromoteObject`, 以及 `sys/storage` 内置资源。

**Step 5:** `i18n.md` 说明 `i18n.T`, `SetLanguage`, 默认语言与验证错误国际化。

**Step 6:** `monitoring.md` 说明 `monitor.Service` 提供的系统概览与构建信息能力。

**Step 7:** 运行验证。

Run: `pnpm build`

Expected: `features` 分类能真实反映“batteries included”的范围，而不是只剩标题。

## Task 7: Correct Advanced Pages And Add Reference Material

**Files:**
- Modify: `docs/advanced/transactions.md`
- Modify: `docs/advanced/context-helpers.md`
- Modify: `docs/advanced/cli-tools.md`
- Create: `docs/reference/_category_.json`
- Create: `docs/reference/rpc-request-format.md`
- Create: `docs/reference/handler-injection.md`
- Create: `docs/reference/search-tags.md`
- Create: `docs/reference/configuration-reference.md`
- Create: `docs/reference/crud-operation-matrix.md`
- Create: `docs/reference/result-and-error-model.md`

**Step 1:** `transactions.md` 说明 `db.RunInTx(...)`、CRUD 钩子事务边界和常见误区；不要写成“全局自动事务”。

**Step 2:** `context-helpers.md` 说明 `contextx.DB`, `contextx.Principal`, `contextx.Logger`, `contextx.DataPermApplier` 的边界和推荐使用场景。

**Step 3:** `cli-tools.md` 只写真实可用命令：
`generate-build-info`、`generate-model-schema`，并明确 `create` 目前未实现。

**Step 4:** 新增 `reference` 分类，把表格型和协议型内容从教程页抽离。

**Step 5:** `rpc-request-format.md` 给出 `api.Request`, `Identifier`, `params`, `meta` 的完整格式。

**Step 6:** `handler-injection.md` 列出 handler 可注入参数：`fiber.Ctx`, `orm.DB`, `log.Logger`, `mold.Transformer`, `*security.Principal`, `page.Pageable`, `api.P`, `api.M`。

**Step 7:** `search-tags.md` 给出支持的 `search` 操作符矩阵和示例。

**Step 8:** `configuration-reference.md` 给出 `vef.app`, `vef.data_source`, `vef.security`, `vef.storage`, `vef.redis`, `vef.cors`, `vef.monitor`, `vef.mcp`, `vef.approval` 的字段参考。

**Step 9:** `crud-operation-matrix.md` 给出 RPC action 与 REST method/path 的对应表。

**Step 10:** `result-and-error-model.md` 说明 `Result`, `Error`, `WithCode`, `WithStatus`, 常见预定义错误。

**Step 11:** 运行验证。

Run: `pnpm build`

Expected: 教程页更聚焦，复杂表格与协议细节有独立 reference 落点。

## Task 8: Add Ecosystem / Preview Pages For Missing Public Packages

**Files:**
- Create: `docs/integrations/_category_.json`
- Create: `docs/integrations/mcp.md`
- Create: `docs/integrations/sequence.md`
- Create: `docs/integrations/approval.md`
- Create: `docs/integrations/ai.md`

**Step 1:** 给 `integrations` 分类加清晰说明：这些页面覆盖源码已有公开包，其中部分能力仍属于 Preview。

**Step 2:** `mcp.md` 说明 `mcp.ToolProvider`, `ResourceProvider`, `PromptProvider`, `ProvideMCPTools(...)` 等注册方式。

**Step 3:** `sequence.md` 说明序列号规则、`Generator`, `Store`, 内存 / Redis / DB 存储实现。

**Step 4:** `approval.md` 只先写高层能力图和核心模型，不承诺完整业务教程；必要时标记为 Preview。

**Step 5:** `ai.md` 说明 `ai.Agent`, `ToolableChatModel`, `Tool`, `ai/stream` 的用途与适用边界，同样标记为 Preview。

**Step 6:** 运行验证。

Run: `pnpm build`

Expected: 源码里已存在的公开能力不再完全消失于站点结构之外。

## Task 9: Generate And Fill Chinese Documentation

**Files:**
- Create: `i18n/zh-Hans/docusaurus-plugin-content-docs/current/**`
- Modify: `i18n/zh-Hans/docusaurus-plugin-content-docs/current.json`

**Step 1:** 在英文页稳定后生成文档翻译骨架。

Run: `pnpm write-translations -- --locale zh-Hans`

Expected: 生成 `i18n/zh-Hans/docusaurus-plugin-content-docs/current/` 目录。

**Step 2:** 优先翻译首页路径和 first-wave 页面：
`intro`, `getting-started/*`, `guide/crud`, `guide/models`, `security/*`, `features/cache`, `features/storage`。

**Step 3:** 对代码块保持英文标识符与注释风格，只翻译自然语言。

**Step 4:** 用 `pnpm build` 验证双语构建。

Run: `pnpm build`

Expected: `en` 与 `zh-Hans` 同时成功输出，无缺失翻译导致的构建错误。

## Task 10: Final QA And Publication Readiness

**Files:**
- Modify as needed across `docs/**`, `i18n/**`, `src/pages/index.tsx`, `docusaurus.config.ts`

**Step 1:** 逐页核对所有代码示例都能在 `../vef-framework-go` 中找到对应 API，避免“文档发明接口”。

**Step 2:** 检查所有 cross-links、frontmatter、sidebar labels、slug 是否稳定。

**Step 3:** 在英文与中文两个入口分别抽查以下路径：
`/docs/intro`
`/docs/getting-started/quick-start`
`/docs/guide/crud`
`/docs/security/authentication`
`/docs/reference/rpc-request-format`

**Step 4:** 最终构建。

Run: `pnpm build`

Expected: 全站构建成功，首页文案、目录结构、核心页面、参考页面、多语言路径均可访问。

## Delivery Order

推荐执行顺序：

1. `Task 1`
2. `Task 2`
3. `Task 3`
4. `Task 4`
5. `Task 5`
6. `Task 6`
7. `Task 7`
8. `Task 8`
9. `Task 9`
10. `Task 10`

## Risks And Controls

- 风险：直接把 `approval` / `ai` 写成正式承诺页面，后续源码 API 变动会导致文档快速失效。
  - 控制：这两页先标记 `Preview`，只覆盖稳定接口和适用场景。
- 风险：在 `docs/` 下存放内部计划或迁移记录，会被 Docusaurus 自动发布。
  - 控制：所有执行计划、审计记录、迁移脚本说明统一放仓库根目录 `plans/`。
- 风险：英文和中文同时改会导致翻译返工。
  - 控制：先稳定英文结构，再生成并翻译 `zh-Hans`。
- 风险：README 与站点内容双写后失去同步。
  - 控制：把 README 视为“项目首页摘要”，站点为完整文档；未来 README 只保留压缩版并链接到站点。
