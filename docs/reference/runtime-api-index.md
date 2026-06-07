---
sidebar_position: 91
---

# Runtime API Index

This page is generated from the current VEF Framework Go source tree. It covers runtime contracts users call, configure, send, receive, or match: HTTP/RPC protocol fields, built-in resources and actions, CLI commands and flags, configuration keys and defaults, events, error codes, wire JSON fields, tag grammars, MCP endpoints/tools/prompts, and runtime enum values.

It intentionally excludes test fixtures, internal log strings, and implementation-only literals. The exported Go import surface is tracked separately in [Public API Index](./public-api-index).

The complete public API audit is the union of this runtime index, the exported Go API index, and the package reviews in `scripts/api-contract-ledger.json`. A user-facing API change must update all affected audit artifacts before the docs review is complete.

Regenerate and verify this page whenever the framework runtime surface changes:

```bash
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-runtime-api-audit.go -source . -out ../vef-framework-go-docs -write)
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-runtime-api-audit.go -source . -out ../vef-framework-go-docs)
```

Fingerprint: `c202cd6432cd10026931c0e25691ad6358c8683e5e3a354331e8dd6cf7e99efe`
Entries: `2123`

## Coverage Evidence

| Category | Entries | Tier | Extractor | Method | Known residual |
| --- | ---: | --- | --- | --- | --- |
| `API default` | 4 | Tier 3 curated source references | `extractProtocolConstants` | Curated defaults from API engine call sites and protocol constants. | None in generated index; semantic behavior remains covered by guide pages. |
| `API version` | 9 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of api/version.go VersionV* string constants. | None. |
| `CLI command` | 4 | Tier 2 Cobra AST | `extractCLI` | AST scan of cobra.Command composites under cmd/vef-cli/cmd. | None in scanned CLI package. |
| `CLI flag` | 8 | Tier 2 Cobra AST | `extractCLI` | AST scan of String/Bool/Int flag helper families and MarkFlagRequired calls under cmd/vef-cli/cmd; unsupported flag definition helpers fail boundary verification. | None for current Cobra flag definition calls. |
| `CRUD REST action` | 14 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of CRUD REST action constants. | None. |
| `CRUD RPC action` | 14 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of CRUD RPC action constants. | None. |
| `HTTP endpoint` | 2 | Tier 2 source-derived constants | `extractProtocolConstants` | Source-derived REST/RPC/MCP endpoint constants and call-site evidence. | None for framework-owned default endpoints. |
| `HTTP header` | 5 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of api/header.go Header* constants. | None. |
| `HTTP wire field` | 8 | Tier 3 curated protocol fields | `extractProtocolConstants` | Curated source references for fundamental request/result fields shared by REST/RPC. | None in generated index; JSON DTO fields are covered separately. |
| `JSON wire field` | 1142 | Tier 2 scoped DTO AST with closed-world boundary check | `extractJSONFields` | AST scan of json tags on runtime DTO structs plus a boundary check over every non-test json-tagged struct field. | None for current non-test source; new json-tagged runtime fields must be indexed or explicitly excluded. |
| `MCP endpoint` | 1 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of the MCP Streamable HTTP endpoint constant. | None. |
| `MCP jsonschema tag` | 32 | Tier 2 pinned dependency parser catalog | `extractJSONSchemaTags` | Catalog of struct-tag keywords accepted by github.com/invopop/jsonschema v0.14.0, with boundary verification that fails on dependency-version drift and uncovered in-source jsonschema tags. | None for the pinned jsonschema parser version. |
| `MCP prompt` | 1 | Tier 2 MCP AST | `extractMCP` | AST scan of internal/mcp Prompt composites. | None in scanned MCP package. |
| `MCP tool` | 1 | Tier 2 MCP AST | `extractMCP` | AST scan of internal/mcp Tool composites. | None in scanned MCP package. |
| `REST action verb` | 10 | Tier 2 validator AST | `extractRESTVerbs` | AST scan of the REST action validator's allowed HTTP verb set. | None in current validator construction. |
| `RPC form key` | 2 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of FormKey* constants. | None. |
| `auth strategy` | 3 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of api/auth.go AuthStrategy* string constants. | None. |
| `auth type` | 4 | Tier 2 scoped AST constants | `extractAuthTypes` | AST scan of internal/security AuthType* constants that are sent through Authentication.Type. | None in known built-in authenticators. |
| `built-in resource` | 10 | Tier 2 scoped AST resources | `extractBuiltInResources` | AST scan of NewRPCResource/NewRESTResource calls in built-in runtime resource packages. | None in scanned built-in resource directories. |
| `built-in resource action` | 62 | Tier 2 scoped AST operations | `extractBuiltInResources` | AST scan of explicit OperationSpec values and CRUD builder defaults inside built-in runtime resource packages. | None in scanned built-in resource directories. |
| `config default` | 44 | Tier 3 mixed static extraction | `extractConfigDefaults` | AST extraction of Effective* accessors, ApplyDefaults assignments, monitor DefaultConfig values, and curated source references for defaults outside those named surfaces; boundary verification fails when a supported default surface is not indexed. | Defaults outside Effective*/ApplyDefaults/DefaultConfig and curated reviewed call sites require explicit review. |
| `config enum` | 8 | Tier 2 scoped AST constants | `extractProtocolConstants` | AST scan of storage and datasource enum constants used in configuration values. | None in current config enum files. |
| `config key` | 105 | Tier 2 config-tag AST | `extractConfigKeys` | AST walk of config structs rooted at known vef.* config roots plus vef.data_sources.&lt;name&gt;; verifier fails if a config/ struct with config tags is unreachable. | None for config/ structs with config tags. |
| `config reserved name` | 1 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of reserved configuration-name constants. | None. |
| `environment variable` | 5 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of Env* constants plus boundary checks for os.Getenv/os.LookupEnv call sites. | None for string-literal or const-backed environment lookups. |
| `event topic` | 33 | Tier 2 event constant/method scan | `extractProtocolConstants, extractMoldGrammar` | AST scan of EventType*/eventType* constants, EventType() return values, and built-in subscription/route-inspection topic call sites. | None for framework-owned non-test event topics. |
| `event transport contract` | 6 | Tier 2 scoped AST constants | `extractEventTransportContracts` | AST/source-derived extraction of outbox DLQ headers, topic prefix, retry backoff, and persisted-error bounds. | None for current built-in event transports. |
| `i18n key indirection` | 4 | Tier 2 AST call scan | `extractI18NMessageKeys` | AST scan of dynamic i18n.T call sites whose key source is another audited surface such as label_i18n tags, validator rules, or Fiber error mappings. | None for current dynamic i18n.T call sites. |
| `i18n message key` | 160 | Tier 2 AST call/tag scan | `extractI18NMessageKeys` | AST scan of literal or const-backed i18n.T calls, validator rule message keys, and label_i18n struct tags. | None for literal or const-backed keys; dynamic sources are tracked as i18n key indirections. |
| `meta tag grammar` | 7 | Tier 2 AST constants | `extractStructTagGrammars` | Catalog of storage meta tag name, dive value, file-reference kinds, and attribute grammar delimiters. | None for the current parser constants and tag parsing rules. |
| `mold tag grammar` | 9 | Tier 2 parser grammar scan | `extractMoldGrammar` | AST scan of the default mold tag name and restricted parser token constants, with boundary verification for parser token coverage. | None for current mold parser token constants. |
| `mold transformer tag` | 2 | Tier 2 transformer scan | `extractMoldGrammar` | AST scan of built-in FieldTransformer Tag() methods. | None for current built-in mold transformer Tag() methods. |
| `mold translate kind prefix` | 1 | Tier 2 translator scan | `extractMoldGrammar` | AST scan of built-in Translator Supports(kind) prefix checks. | None for current built-in translate kind prefixes. |
| `result error code` | 77 | Tier 1 AST constants | `extractErrorCodes` | AST scan of ErrCode* constants in api_errors.go and result/constants.go. | None for named error-code constants. |
| `result message key` | 33 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of ErrMessage* constants. | Inline i18n keys are covered by the i18n message key category. |
| `runtime enum value` | 234 | Tier 2 typed string constants | `extractRuntimeEnumValues` | AST scan of typed string constants in public packages plus runtime internal DTO/transport packages. | Integer/stringer enum renderings are covered by the generated public API index and package contract ledger. |
| `search tag grammar` | 40 | Tier 1 AST constants | `extractStructTagGrammars` | AST scan of search tag name, attributes, params, ignore marker, and operator/type tokens. | None for constants in search/constants.go. |
| `tabular tag grammar` | 10 | Tier 1 AST constants | `extractStructTagGrammars` | AST scan of tabular tag name, attributes, and ignore marker. | None for constants in tabular/constants.go. |
| `validator label tag` | 2 | Tier 2 validator tag scan | `extractValidatorRules` | AST scan of validator struct-tag key constants used by Field.Tag.Get. | None for current validator label tag lookups. |
| `validator tag` | 6 | Tier 2 validator AST | `extractValidatorRules` | AST scan of built-in validator rule registration calls. | None for current built-in validator registrations and ValidationRule composites. |

## API default

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `default auth strategy` | `bearer` |  | `api/auth.go:24`, `internal/api/engine.go:102` |
| `default rate limit` | `100 requests / 5m` |  | `internal/api/engine.go:104` |
| `default timeout` | `30s` |  | `internal/api/engine.go:100` |
| `default version` | `v1` |  | `api/version.go:4`, `internal/api/engine.go:101` |

## API version

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `VersionV1` | `v1` |  | `api/version.go:5` |
| `VersionV2` | `v2` |  | `api/version.go:6` |
| `VersionV3` | `v3` |  | `api/version.go:7` |
| `VersionV4` | `v4` |  | `api/version.go:8` |
| `VersionV5` | `v5` |  | `api/version.go:9` |
| `VersionV6` | `v6` |  | `api/version.go:10` |
| `VersionV7` | `v7` |  | `api/version.go:11` |
| `VersionV8` | `v8` |  | `api/version.go:12` |
| `VersionV9` | `v9` |  | `api/version.go:13` |

## CLI command

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `create` | `create` | Create a new VEF Framework project | `cmd/vef-cli/cmd/create/command.go:14` |
| `generate-build-info` | `generate-build-info` | Generate build information for the application | `cmd/vef-cli/cmd/buildinfo/command.go:14` |
| `generate-model-schema` | `generate-model-schema` | Generate schema structures from Go models | `cmd/vef-cli/cmd/modelschema/command.go:18` |
| `vef-cli` | `vef-cli` | VEF Framework CLI tool | `cmd/vef-cli/cmd/root.go:16` |

## CLI flag

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `create --module` | `--module` | Go module name (e.g., github.com/user/project)<br/>short: -m | `cmd/vef-cli/cmd/create/command.go:40` |
| `create --name` | `--name` | Project name (required)<br/>required<br/>short: -n | `cmd/vef-cli/cmd/create/command.go:38` |
| `create --path` | `--path` | Directory path to create the project<br/>default: .<br/>short: -p | `cmd/vef-cli/cmd/create/command.go:39` |
| `generate-build-info --output` | `--output` | Output file path<br/>default: build_info.go<br/>short: -o | `cmd/vef-cli/cmd/buildinfo/command.go:54` |
| `generate-build-info --package` | `--package` | Package name<br/>default: main<br/>short: -p | `cmd/vef-cli/cmd/buildinfo/command.go:55` |
| `generate-model-schema --input` | `--input` | Input model file or directory path<br/>required<br/>short: -i | `cmd/vef-cli/cmd/modelschema/command.go:54` |
| `generate-model-schema --output` | `--output` | Output schema file or directory path<br/>required<br/>short: -o | `cmd/vef-cli/cmd/modelschema/command.go:55` |
| `generate-model-schema --package` | `--package` | Package name for generated schemas<br/>default: schemas<br/>short: -p | `cmd/vef-cli/cmd/modelschema/command.go:56` |

## CRUD REST action

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `RESTActionCreate` | `post /` |  | `crud/constants.go:39` |
| `RESTActionCreateMany` | `post /many` |  | `crud/constants.go:42` |
| `RESTActionDelete` | `delete /:id` |  | `crud/constants.go:41` |
| `RESTActionDeleteMany` | `delete /many` |  | `crud/constants.go:44` |
| `RESTActionExport` | `get /export` |  | `crud/constants.go:52` |
| `RESTActionFindAll` | `get /` |  | `crud/constants.go:46` |
| `RESTActionFindOne` | `get /:id` |  | `crud/constants.go:45` |
| `RESTActionFindOptions` | `get /options` |  | `crud/constants.go:48` |
| `RESTActionFindPage` | `get /page` |  | `crud/constants.go:47` |
| `RESTActionFindTree` | `get /tree` |  | `crud/constants.go:49` |
| `RESTActionFindTreeOptions` | `get /tree/options` |  | `crud/constants.go:50` |
| `RESTActionImport` | `post /import` |  | `crud/constants.go:51` |
| `RESTActionUpdate` | `put /:id` |  | `crud/constants.go:40` |
| `RESTActionUpdateMany` | `put /many` |  | `crud/constants.go:43` |

## CRUD RPC action

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `RPCActionCreate` | `create` |  | `crud/constants.go:21` |
| `RPCActionCreateMany` | `create_many` |  | `crud/constants.go:24` |
| `RPCActionDelete` | `delete` |  | `crud/constants.go:23` |
| `RPCActionDeleteMany` | `delete_many` |  | `crud/constants.go:26` |
| `RPCActionExport` | `export` |  | `crud/constants.go:34` |
| `RPCActionFindAll` | `find_all` |  | `crud/constants.go:28` |
| `RPCActionFindOne` | `find_one` |  | `crud/constants.go:27` |
| `RPCActionFindOptions` | `find_options` |  | `crud/constants.go:30` |
| `RPCActionFindPage` | `find_page` |  | `crud/constants.go:29` |
| `RPCActionFindTree` | `find_tree` |  | `crud/constants.go:31` |
| `RPCActionFindTreeOptions` | `find_tree_options` |  | `crud/constants.go:32` |
| `RPCActionImport` | `import` |  | `crud/constants.go:33` |
| `RPCActionUpdate` | `update` |  | `crud/constants.go:22` |
| `RPCActionUpdateMany` | `update_many` |  | `crud/constants.go:25` |

## HTTP endpoint

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `REST base path` | `/api` |  | `internal/api/router/rest.go:19` |
| `RPC endpoint` | `/api` | POST endpoint for RPC requests | `internal/api/router/rpc.go:21` |

## HTTP header

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `HeaderXAppID` | `X-App-ID` |  | `api/header.go:5` |
| `HeaderXMetaPrefix` | `X-Meta-` |  | `api/header.go:9` |
| `HeaderXNonce` | `X-Nonce` |  | `api/header.go:7` |
| `HeaderXSignature` | `X-Signature` |  | `api/header.go:8` |
| `HeaderXTimestamp` | `X-Timestamp` |  | `api/header.go:6` |

## HTTP wire field

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `api.Identifier.action` | `action` | form:"action" | `api/request.go:15` |
| `api.Identifier.resource` | `resource` | form:"resource" | `api/request.go:14` |
| `api.Identifier.version` | `version` | form:"version" | `api/request.go:16` |
| `api.Request.meta` | `meta` |  | `api/request.go:59` |
| `api.Request.params` | `params` |  | `api/request.go:58` |
| `result.Result.code` | `code` |  | `result/result.go:11` |
| `result.Result.data` | `data` |  | `result/result.go:13` |
| `result.Result.message` | `message` |  | `result/result.go:12` |

## JSON wire field

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AbortUploadParams.ClaimID` | `claimId` | Go field: AbortUploadParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:746` |
| `ActionLog.Action` | `action` | Go field: ActionLog.Action<br/>type: string | `approval/admin/types.go:83` |
| `ActionLog.Action` | `action` | Go field: ActionLog.Action<br/>type: ActionType | `approval/models.go:240` |
| `ActionLog.AddAssigneeType` | `addAssigneeType` | Go field: ActionLog.AddAssigneeType<br/>type: *AddAssigneeType | `approval/models.go:251` |
| `ActionLog.AddedAssigneeIDs` | `addedAssigneeIds` | Go field: ActionLog.AddedAssigneeIDs<br/>type: []string | `approval/models.go:252` |
| `ActionLog.Attachments` | `attachments` | Go field: ActionLog.Attachments<br/>type: []string | `approval/models.go:255` |
| `ActionLog.CCUserIDs` | `ccUserIds` | Go field: ActionLog.CCUserIDs<br/>type: []string | `approval/models.go:254` |
| `ActionLog.CreatedAt` | `createdAt` | Go field: ActionLog.CreatedAt<br/>type: timex.DateTime | `approval/admin/types.go:90` |
| `ActionLog.IPAddress` | `ipAddress` | Go field: ActionLog.IPAddress<br/>type: *string | `approval/models.go:245` |
| `ActionLog.InstanceID` | `instanceId` | Go field: ActionLog.InstanceID<br/>type: string | `approval/models.go:237` |
| `ActionLog.LogID` | `logId` | Go field: ActionLog.LogID<br/>type: string | `approval/admin/types.go:82` |
| `ActionLog.Meta` | `meta` | Go field: ActionLog.Meta<br/>type: map[string]any | `approval/models.go:256` |
| `ActionLog.NodeID` | `nodeId` | Go field: ActionLog.NodeID<br/>type: *string | `approval/models.go:238` |
| `ActionLog.OperatorDepartmentID` | `operatorDepartmentId` | Go field: ActionLog.OperatorDepartmentID<br/>type: *string | `approval/models.go:243` |
| `ActionLog.OperatorDepartmentName` | `operatorDepartmentName` | Go field: ActionLog.OperatorDepartmentName<br/>type: *string | `approval/admin/types.go:86` |
| `ActionLog.OperatorDepartmentName` | `operatorDepartmentName` | Go field: ActionLog.OperatorDepartmentName<br/>type: *string | `approval/models.go:244` |
| `ActionLog.OperatorID` | `operatorId` | Go field: ActionLog.OperatorID<br/>type: string | `approval/admin/types.go:84` |
| `ActionLog.OperatorID` | `operatorId` | Go field: ActionLog.OperatorID<br/>type: string | `approval/models.go:241` |
| `ActionLog.OperatorName` | `operatorName` | Go field: ActionLog.OperatorName<br/>type: string | `approval/models.go:242` |
| `ActionLog.OperatorName` | `operatorName` | Go field: ActionLog.OperatorName<br/>type: string | `approval/admin/types.go:85` |
| `ActionLog.Opinion` | `opinion` | Go field: ActionLog.Opinion<br/>type: *string | `approval/models.go:247` |
| `ActionLog.Opinion` | `opinion` | Go field: ActionLog.Opinion<br/>type: *string | `approval/admin/types.go:89` |
| `ActionLog.RemovedAssigneeIDs` | `removedAssigneeIds` | Go field: ActionLog.RemovedAssigneeIDs<br/>type: []string | `approval/models.go:253` |
| `ActionLog.RollbackToNodeID` | `rollbackToNodeId` | Go field: ActionLog.RollbackToNodeID<br/>type: *string | `approval/models.go:250` |
| `ActionLog.TaskID` | `taskId` | Go field: ActionLog.TaskID<br/>type: *string | `approval/models.go:239` |
| `ActionLog.TransferToID` | `transferToId` | Go field: ActionLog.TransferToID<br/>type: *string | `approval/models.go:248` |
| `ActionLog.TransferToID` | `transferToId` | Go field: ActionLog.TransferToID<br/>type: *string | `approval/admin/types.go:87` |
| `ActionLog.TransferToName` | `transferToName` | Go field: ActionLog.TransferToName<br/>type: *string | `approval/models.go:249` |
| `ActionLog.TransferToName` | `transferToName` | Go field: ActionLog.TransferToName<br/>type: *string | `approval/admin/types.go:88` |
| `ActionLog.UserAgent` | `userAgent` | Go field: ActionLog.UserAgent<br/>type: *string | `approval/models.go:246` |
| `ActionLogInfo.Action` | `action` | Go field: ActionLogInfo.Action<br/>type: string | `approval/my/types.go:118` |
| `ActionLogInfo.CreatedAt` | `createdAt` | Go field: ActionLogInfo.CreatedAt<br/>type: timex.DateTime | `approval/my/types.go:121` |
| `ActionLogInfo.OperatorName` | `operatorName` | Go field: ActionLogInfo.OperatorName<br/>type: string | `approval/my/types.go:119` |
| `ActionLogInfo.Opinion` | `opinion` | Go field: ActionLogInfo.Opinion<br/>type: *string | `approval/my/types.go:120` |
| `AddAssigneeParams.AddType` | `addType` | Go field: AddAssigneeParams.AddType<br/>type: string<br/>validate: "required,oneof=before after parallel" | `internal/approval/resource/instance.go:335` |
| `AddAssigneeParams.TaskID` | `taskId` | Go field: AddAssigneeParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:333` |
| `AddAssigneeParams.UserIDs` | `userIds` | Go field: AddAssigneeParams.UserIDs<br/>type: []string<br/>validate: "required,min=1,max=50" | `internal/approval/resource/instance.go:334` |
| `AddCCParams.CCUserIDs` | `ccUserIds` | Go field: AddCCParams.CCUserIDs<br/>type: []string<br/>validate: "required,min=1,max=50" | `internal/approval/resource/instance.go:282` |
| `AddCCParams.InstanceID` | `instanceId` | Go field: AddCCParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:281` |
| `AdminFindActionLogsParams.InstanceID` | `instanceId` | Go field: AdminFindActionLogsParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:176` |
| `AdminFindActionLogsParams.Page` | `page` | Go field: AdminFindActionLogsParams.Page<br/>type: int | `internal/approval/resource/admin.go:178` |
| `AdminFindActionLogsParams.PageSize` | `pageSize` | Go field: AdminFindActionLogsParams.PageSize<br/>type: int | `internal/approval/resource/admin.go:179` |
| `AdminFindActionLogsParams.TenantID` | `tenantId` | Go field: AdminFindActionLogsParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:177` |
| `AdminFindInstancesParams.ApplicantID` | `applicantId` | Go field: AdminFindInstancesParams.ApplicantID<br/>type: *string | `internal/approval/resource/admin.go:60` |
| `AdminFindInstancesParams.FlowID` | `flowId` | Go field: AdminFindInstancesParams.FlowID<br/>type: *string | `internal/approval/resource/admin.go:62` |
| `AdminFindInstancesParams.Keyword` | `keyword` | Go field: AdminFindInstancesParams.Keyword<br/>type: *string | `internal/approval/resource/admin.go:63` |
| `AdminFindInstancesParams.Page` | `page` | Go field: AdminFindInstancesParams.Page<br/>type: int | `internal/approval/resource/admin.go:64` |
| `AdminFindInstancesParams.PageSize` | `pageSize` | Go field: AdminFindInstancesParams.PageSize<br/>type: int | `internal/approval/resource/admin.go:65` |
| `AdminFindInstancesParams.Status` | `status` | Go field: AdminFindInstancesParams.Status<br/>type: *approval.InstanceStatus | `internal/approval/resource/admin.go:61` |
| `AdminFindInstancesParams.TenantID` | `tenantId` | Go field: AdminFindInstancesParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:59` |
| `AdminFindTasksParams.AssigneeID` | `assigneeId` | Go field: AdminFindTasksParams.AssigneeID<br/>type: *string | `internal/approval/resource/admin.go:119` |
| `AdminFindTasksParams.InstanceID` | `instanceId` | Go field: AdminFindTasksParams.InstanceID<br/>type: *string | `internal/approval/resource/admin.go:120` |
| `AdminFindTasksParams.Page` | `page` | Go field: AdminFindTasksParams.Page<br/>type: int | `internal/approval/resource/admin.go:122` |
| `AdminFindTasksParams.PageSize` | `pageSize` | Go field: AdminFindTasksParams.PageSize<br/>type: int | `internal/approval/resource/admin.go:123` |
| `AdminFindTasksParams.Status` | `status` | Go field: AdminFindTasksParams.Status<br/>type: *approval.TaskStatus | `internal/approval/resource/admin.go:121` |
| `AdminFindTasksParams.TenantID` | `tenantId` | Go field: AdminFindTasksParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:118` |
| `AdminGetInstanceDetailParams.InstanceID` | `instanceId` | Go field: AdminGetInstanceDetailParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:151` |
| `AdminGetMetricsParams.TenantID` | `tenantId` | Go field: AdminGetMetricsParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:261` |
| `AdminReassignTaskParams.NewAssigneeID` | `newAssigneeId` | Go field: AdminReassignTaskParams.NewAssigneeID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:233` |
| `AdminReassignTaskParams.Reason` | `reason` | Go field: AdminReassignTaskParams.Reason<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/admin.go:234` |
| `AdminReassignTaskParams.TaskID` | `taskId` | Go field: AdminReassignTaskParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:232` |
| `AdminTerminateInstanceParams.InstanceID` | `instanceId` | Go field: AdminTerminateInstanceParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:205` |
| `AdminTerminateInstanceParams.Reason` | `reason` | Go field: AdminTerminateInstanceParams.Reason<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/admin.go:206` |
| `ApprovalNodeData.AddAssigneeTypes` | `addAssigneeTypes` | Go field: ApprovalNodeData.AddAssigneeTypes<br/>type: []AddAssigneeType | `approval/node_data.go:131` |
| `ApprovalNodeData.ApprovalMethod` | `approvalMethod` | Go field: ApprovalNodeData.ApprovalMethod<br/>type: ApprovalMethod | `approval/node_data.go:121` |
| `ApprovalNodeData.ConsecutiveApproverAction` | `consecutiveApproverAction` | Go field: ApprovalNodeData.ConsecutiveApproverAction<br/>type: ConsecutiveApproverAction | `approval/node_data.go:125` |
| `ApprovalNodeData.IsAddAssigneeAllowed` | `isAddAssigneeAllowed` | Go field: ApprovalNodeData.IsAddAssigneeAllowed<br/>type: bool | `approval/node_data.go:130` |
| `ApprovalNodeData.IsManualCCAllowed` | `isManualCcAllowed` | Go field: ApprovalNodeData.IsManualCCAllowed<br/>type: bool | `approval/node_data.go:133` |
| `ApprovalNodeData.IsRemoveAssigneeAllowed` | `isRemoveAssigneeAllowed` | Go field: ApprovalNodeData.IsRemoveAssigneeAllowed<br/>type: bool | `approval/node_data.go:132` |
| `ApprovalNodeData.IsRollbackAllowed` | `isRollbackAllowed` | Go field: ApprovalNodeData.IsRollbackAllowed<br/>type: bool | `approval/node_data.go:129` |
| `ApprovalNodeData.PassRatio` | `passRatio` | Go field: ApprovalNodeData.PassRatio<br/>type: decimal.Decimal | `approval/node_data.go:123` |
| `ApprovalNodeData.PassRule` | `passRule` | Go field: ApprovalNodeData.PassRule<br/>type: PassRule | `approval/node_data.go:122` |
| `ApprovalNodeData.RollbackDataStrategy` | `rollbackDataStrategy` | Go field: ApprovalNodeData.RollbackDataStrategy<br/>type: RollbackDataStrategy | `approval/node_data.go:127` |
| `ApprovalNodeData.RollbackTargetKeys` | `rollbackTargetKeys` | Go field: ApprovalNodeData.RollbackTargetKeys<br/>type: []string | `approval/node_data.go:128` |
| `ApprovalNodeData.RollbackType` | `rollbackType` | Go field: ApprovalNodeData.RollbackType<br/>type: RollbackType | `approval/node_data.go:126` |
| `ApprovalNodeData.SameApplicantAction` | `sameApplicantAction` | Go field: ApprovalNodeData.SameApplicantAction<br/>type: SameApplicantAction | `approval/node_data.go:124` |
| `AssigneeDefinition.FormField` | `formField` | Go field: AssigneeDefinition.FormField<br/>type: *string | `approval/assignee.go:51` |
| `AssigneeDefinition.IDs` | `ids` | Go field: AssigneeDefinition.IDs<br/>type: []string | `approval/assignee.go:50` |
| `AssigneeDefinition.Kind` | `kind` | Go field: AssigneeDefinition.Kind<br/>type: AssigneeKind | `approval/assignee.go:49` |
| `AssigneeDefinition.SortOrder` | `sortOrder` | Go field: AssigneeDefinition.SortOrder<br/>type: int | `approval/assignee.go:52` |
| `AssigneesAddedEvent.AddType` | `addType` | Go field: AssigneesAddedEvent.AddType<br/>type: AddAssigneeType | `approval/events_task.go:212` |
| `AssigneesAddedEvent.AssigneeIDs` | `assigneeIds` | Go field: AssigneesAddedEvent.AssigneeIDs<br/>type: []string | `approval/events_task.go:213` |
| `AssigneesAddedEvent.AssigneeNames` | `assigneeNames` | Go field: AssigneesAddedEvent.AssigneeNames<br/>type: map[string]string | `approval/events_task.go:214` |
| `AssigneesAddedEvent.InstanceID` | `instanceId` | Go field: AssigneesAddedEvent.InstanceID<br/>type: string | `approval/events_task.go:208` |
| `AssigneesAddedEvent.NodeID` | `nodeId` | Go field: AssigneesAddedEvent.NodeID<br/>type: string | `approval/events_task.go:210` |
| `AssigneesAddedEvent.OccurredTime` | `occurredTime` | Go field: AssigneesAddedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:215` |
| `AssigneesAddedEvent.TaskID` | `taskId` | Go field: AssigneesAddedEvent.TaskID<br/>type: string | `approval/events_task.go:211` |
| `AssigneesAddedEvent.TenantID` | `tenantId` | Go field: AssigneesAddedEvent.TenantID<br/>type: string | `approval/events_task.go:209` |
| `AssigneesRemovedEvent.AssigneeIDs` | `assigneeIds` | Go field: AssigneesRemovedEvent.AssigneeIDs<br/>type: []string | `approval/events_task.go:239` |
| `AssigneesRemovedEvent.AssigneeNames` | `assigneeNames` | Go field: AssigneesRemovedEvent.AssigneeNames<br/>type: map[string]string | `approval/events_task.go:240` |
| `AssigneesRemovedEvent.InstanceID` | `instanceId` | Go field: AssigneesRemovedEvent.InstanceID<br/>type: string | `approval/events_task.go:235` |
| `AssigneesRemovedEvent.NodeID` | `nodeId` | Go field: AssigneesRemovedEvent.NodeID<br/>type: string | `approval/events_task.go:237` |
| `AssigneesRemovedEvent.OccurredTime` | `occurredTime` | Go field: AssigneesRemovedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:241` |
| `AssigneesRemovedEvent.TaskID` | `taskId` | Go field: AssigneesRemovedEvent.TaskID<br/>type: string | `approval/events_task.go:238` |
| `AssigneesRemovedEvent.TenantID` | `tenantId` | Go field: AssigneesRemovedEvent.TenantID<br/>type: string | `approval/events_task.go:236` |
| `AuditEvent.Action` | `action` | Go field: AuditEvent.Action<br/>type: string | `api/audit.go:15` |
| `AuditEvent.ElapsedTime` | `elapsedTime` | Go field: AuditEvent.ElapsedTime<br/>type: int64 | `api/audit.go:34` |
| `AuditEvent.RequestID` | `requestId` | Go field: AuditEvent.RequestID<br/>type: string | `api/audit.go:23` |
| `AuditEvent.RequestIP` | `requestIp` | Go field: AuditEvent.RequestIP<br/>type: string | `api/audit.go:24` |
| `AuditEvent.RequestMeta` | `requestMeta` | Go field: AuditEvent.RequestMeta<br/>type: map[string]any | `api/audit.go:26` |
| `AuditEvent.RequestParams` | `requestParams` | Go field: AuditEvent.RequestParams<br/>type: map[string]any | `api/audit.go:25` |
| `AuditEvent.Resource` | `resource` | Go field: AuditEvent.Resource<br/>type: string | `api/audit.go:14` |
| `AuditEvent.ResultCode` | `resultCode` | Go field: AuditEvent.ResultCode<br/>type: int | `api/audit.go:29` |
| `AuditEvent.ResultData` | `resultData` | Go field: AuditEvent.ResultData<br/>type: any | `api/audit.go:31` |
| `AuditEvent.ResultMessage` | `resultMessage` | Go field: AuditEvent.ResultMessage<br/>type: string | `api/audit.go:30` |
| `AuditEvent.UserAgent` | `userAgent` | Go field: AuditEvent.UserAgent<br/>type: string | `api/audit.go:20` |
| `AuditEvent.UserID` | `userId` | Go field: AuditEvent.UserID<br/>type: string | `api/audit.go:19` |
| `AuditEvent.Version` | `version` | Go field: AuditEvent.Version<br/>type: string | `api/audit.go:16` |
| `AuthTokens.AccessToken` | `accessToken` | Go field: AuthTokens.AccessToken<br/>type: string | `security/security.go:14` |
| `AuthTokens.RefreshToken` | `refreshToken` | Go field: AuthTokens.RefreshToken<br/>type: string | `security/security.go:15` |
| `Authentication.Credentials` | `credentials` | Go field: Authentication.Credentials<br/>type: any | `security/security.go:22` |
| `Authentication.Principal` | `principal` | Go field: Authentication.Principal<br/>type: string | `security/security.go:21` |
| `Authentication.Type` | `type` | Go field: Authentication.Type<br/>type: string | `security/security.go:20` |
| `AvailableFlow.CategoryID` | `categoryId` | Go field: AvailableFlow.CategoryID<br/>type: string | `approval/my/types.go:12` |
| `AvailableFlow.CategoryName` | `categoryName` | Go field: AvailableFlow.CategoryName<br/>type: string | `approval/my/types.go:13` |
| `AvailableFlow.Description` | `description` | Go field: AvailableFlow.Description<br/>type: *string | `approval/my/types.go:11` |
| `AvailableFlow.FlowCode` | `flowCode` | Go field: AvailableFlow.FlowCode<br/>type: string | `approval/my/types.go:8` |
| `AvailableFlow.FlowID` | `flowId` | Go field: AvailableFlow.FlowID<br/>type: string | `approval/my/types.go:7` |
| `AvailableFlow.FlowIcon` | `flowIcon` | Go field: AvailableFlow.FlowIcon<br/>type: *string | `approval/my/types.go:10` |
| `AvailableFlow.FlowName` | `flowName` | Go field: AvailableFlow.FlowName<br/>type: string | `approval/my/types.go:9` |
| `BaseNodeData.Description` | `description` | Go field: BaseNodeData.Description<br/>type: *string | `approval/node_data.go:22` |
| `BaseNodeData.Name` | `name` | Go field: BaseNodeData.Name<br/>type: string | `approval/node_data.go:21` |
| `BuildInfo.AppVersion` | `appVersion` | Go field: BuildInfo.AppVersion<br/>type: string | `monitor/service.go:263` |
| `BuildInfo.BuildTime` | `buildTime` | Go field: BuildInfo.BuildTime<br/>type: string | `monitor/service.go:264` |
| `BuildInfo.GitCommit` | `gitCommit` | Go field: BuildInfo.GitCommit<br/>type: string | `monitor/service.go:265` |
| `BuildInfo.VEFVersion` | `vefVersion` | Go field: BuildInfo.VEFVersion<br/>type: string | `monitor/service.go:262` |
| `CCDefinition.FormField` | `formField` | Go field: CCDefinition.FormField<br/>type: *string | `approval/assignee.go:59` |
| `CCDefinition.IDs` | `ids` | Go field: CCDefinition.IDs<br/>type: []string | `approval/assignee.go:58` |
| `CCDefinition.Kind` | `kind` | Go field: CCDefinition.Kind<br/>type: CCKind | `approval/assignee.go:57` |
| `CCDefinition.Timing` | `timing` | Go field: CCDefinition.Timing<br/>type: CCTiming | `approval/assignee.go:60` |
| `CCNodeData.CCs` | `ccs` | Go field: CCNodeData.CCs<br/>type: []CCDefinition | `approval/node_data.go:194` |
| `CCNodeData.FieldPermissions` | `fieldPermissions` | Go field: CCNodeData.FieldPermissions<br/>type: map[string]Permission | `approval/node_data.go:196` |
| `CCNodeData.IsReadConfirmRequired` | `isReadConfirmRequired` | Go field: CCNodeData.IsReadConfirmRequired<br/>type: bool | `approval/node_data.go:195` |
| `CCNotifiedEvent.CCUserIDs` | `ccUserIds` | Go field: CCNotifiedEvent.CCUserIDs<br/>type: []string | `approval/events_cc.go:10` |
| `CCNotifiedEvent.CCUserNames` | `ccUserNames` | Go field: CCNotifiedEvent.CCUserNames<br/>type: map[string]string | `approval/events_cc.go:11` |
| `CCNotifiedEvent.InstanceID` | `instanceId` | Go field: CCNotifiedEvent.InstanceID<br/>type: string | `approval/events_cc.go:7` |
| `CCNotifiedEvent.IsManual` | `isManual` | Go field: CCNotifiedEvent.IsManual<br/>type: bool | `approval/events_cc.go:12` |
| `CCNotifiedEvent.NodeID` | `nodeId` | Go field: CCNotifiedEvent.NodeID<br/>type: string | `approval/events_cc.go:9` |
| `CCNotifiedEvent.OccurredTime` | `occurredTime` | Go field: CCNotifiedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_cc.go:13` |
| `CCNotifiedEvent.TenantID` | `tenantId` | Go field: CCNotifiedEvent.TenantID<br/>type: string | `approval/events_cc.go:8` |
| `CCRecord.ApplicantName` | `applicantName` | Go field: CCRecord.ApplicantName<br/>type: string | `approval/my/types.go:66` |
| `CCRecord.CCRecordID` | `ccRecordId` | Go field: CCRecord.CCRecordID<br/>type: string | `approval/my/types.go:60` |
| `CCRecord.CCUserID` | `ccUserId` | Go field: CCRecord.CCUserID<br/>type: string | `approval/models.go:268` |
| `CCRecord.CCUserName` | `ccUserName` | Go field: CCRecord.CCUserName<br/>type: string | `approval/models.go:269` |
| `CCRecord.CreatedAt` | `createdAt` | Go field: CCRecord.CreatedAt<br/>type: timex.DateTime | `approval/my/types.go:69` |
| `CCRecord.FlowIcon` | `flowIcon` | Go field: CCRecord.FlowIcon<br/>type: *string | `approval/my/types.go:65` |
| `CCRecord.FlowName` | `flowName` | Go field: CCRecord.FlowName<br/>type: string | `approval/my/types.go:64` |
| `CCRecord.InstanceID` | `instanceId` | Go field: CCRecord.InstanceID<br/>type: string | `approval/my/types.go:61` |
| `CCRecord.InstanceID` | `instanceId` | Go field: CCRecord.InstanceID<br/>type: string | `approval/models.go:265` |
| `CCRecord.InstanceNo` | `instanceNo` | Go field: CCRecord.InstanceNo<br/>type: string | `approval/my/types.go:63` |
| `CCRecord.InstanceTitle` | `instanceTitle` | Go field: CCRecord.InstanceTitle<br/>type: string | `approval/my/types.go:62` |
| `CCRecord.IsManual` | `isManual` | Go field: CCRecord.IsManual<br/>type: bool | `approval/models.go:270` |
| `CCRecord.IsRead` | `isRead` | Go field: CCRecord.IsRead<br/>type: bool | `approval/my/types.go:68` |
| `CCRecord.NodeID` | `nodeId` | Go field: CCRecord.NodeID<br/>type: *string | `approval/models.go:266` |
| `CCRecord.NodeName` | `nodeName` | Go field: CCRecord.NodeName<br/>type: *string | `approval/my/types.go:67` |
| `CCRecord.ReadAt` | `readAt` | Go field: CCRecord.ReadAt<br/>type: *timex.DateTime | `approval/models.go:271` |
| `CCRecord.TaskID` | `taskId` | Go field: CCRecord.TaskID<br/>type: *string | `approval/models.go:267` |
| `CPUInfo.CacheSize` | `cacheSize` | Go field: CPUInfo.CacheSize<br/>type: int32 | `monitor/service.go:58` |
| `CPUInfo.Family` | `family` | Go field: CPUInfo.Family<br/>type: string | `monitor/service.go:62` |
| `CPUInfo.LogicalCores` | `logicalCores` | Go field: CPUInfo.LogicalCores<br/>type: int | `monitor/service.go:55` |
| `CPUInfo.Mhz` | `mhz` | Go field: CPUInfo.Mhz<br/>type: float64 | `monitor/service.go:57` |
| `CPUInfo.Microcode` | `microcode` | Go field: CPUInfo.Microcode<br/>type: string | `monitor/service.go:65` |
| `CPUInfo.Model` | `model` | Go field: CPUInfo.Model<br/>type: string | `monitor/service.go:63` |
| `CPUInfo.ModelName` | `modelName` | Go field: CPUInfo.ModelName<br/>type: string | `monitor/service.go:56` |
| `CPUInfo.PhysicalCores` | `physicalCores` | Go field: CPUInfo.PhysicalCores<br/>type: int | `monitor/service.go:54` |
| `CPUInfo.Stepping` | `stepping` | Go field: CPUInfo.Stepping<br/>type: int32 | `monitor/service.go:64` |
| `CPUInfo.TotalPercent` | `totalPercent` | Go field: CPUInfo.TotalPercent<br/>type: float64 | `monitor/service.go:60` |
| `CPUInfo.UsagePercent` | `usagePercent` | Go field: CPUInfo.UsagePercent<br/>type: []float64 | `monitor/service.go:59` |
| `CPUInfo.VendorID` | `vendorId` | Go field: CPUInfo.VendorID<br/>type: string | `monitor/service.go:61` |
| `CPUSummary.LogicalCores` | `logicalCores` | Go field: CPUSummary.LogicalCores<br/>type: int | `monitor/service.go:48` |
| `CPUSummary.PhysicalCores` | `physicalCores` | Go field: CPUSummary.PhysicalCores<br/>type: int | `monitor/service.go:47` |
| `CPUSummary.UsagePercent` | `usagePercent` | Go field: CPUSummary.UsagePercent<br/>type: float64 | `monitor/service.go:49` |
| `CategoryParams.Code` | `code` | Go field: CategoryParams.Code<br/>type: string<br/>validate: "required" | `internal/approval/resource/category.go:20` |
| `CategoryParams.ID` | `id` | Go field: CategoryParams.ID<br/>type: string | `internal/approval/resource/category.go:18` |
| `CategoryParams.Icon` | `icon` | Go field: CategoryParams.Icon<br/>type: *string | `internal/approval/resource/category.go:22` |
| `CategoryParams.IsActive` | `isActive` | Go field: CategoryParams.IsActive<br/>type: bool | `internal/approval/resource/category.go:25` |
| `CategoryParams.Name` | `name` | Go field: CategoryParams.Name<br/>type: string<br/>validate: "required" | `internal/approval/resource/category.go:21` |
| `CategoryParams.ParentID` | `parentId` | Go field: CategoryParams.ParentID<br/>type: *string | `internal/approval/resource/category.go:23` |
| `CategoryParams.Remark` | `remark` | Go field: CategoryParams.Remark<br/>type: *string | `internal/approval/resource/category.go:26` |
| `CategoryParams.SortOrder` | `sortOrder` | Go field: CategoryParams.SortOrder<br/>type: int | `internal/approval/resource/category.go:24` |
| `CategoryParams.TenantID` | `tenantId` | Go field: CategoryParams.TenantID<br/>type: string<br/>validate: "required" | `internal/approval/resource/category.go:19` |
| `CategorySearch.IsActive` | `isActive` | Go field: CategorySearch.IsActive<br/>search: "eq,column=is_active"<br/>type: *bool | `internal/approval/resource/category.go:34` |
| `CategorySearch.Name` | `name` | Go field: CategorySearch.Name<br/>search: "contains"<br/>type: string | `internal/approval/resource/category.go:33` |
| `Check.Expr` | `expr` | Go field: Check.Expr<br/>type: string | `schema/service.go:67` |
| `Check.Name` | `name` | Go field: Check.Name<br/>type: string | `schema/service.go:66` |
| `Column.Comment` | `comment` | Go field: Column.Comment<br/>type: string | `schema/service.go:31` |
| `Column.Default` | `default` | Go field: Column.Default<br/>type: string | `schema/service.go:30` |
| `Column.IsAutoIncrement` | `isAutoIncrement` | Go field: Column.IsAutoIncrement<br/>type: bool | `schema/service.go:33` |
| `Column.IsPrimaryKey` | `isPrimaryKey` | Go field: Column.IsPrimaryKey<br/>type: bool | `schema/service.go:32` |
| `Column.Name` | `name` | Go field: Column.Name<br/>type: string | `schema/service.go:27` |
| `Column.Nullable` | `nullable` | Go field: Column.Nullable<br/>type: bool | `schema/service.go:29` |
| `Column.Type` | `type` | Go field: Column.Type<br/>type: string | `schema/service.go:28` |
| `CompleteUploadParams.ClaimID` | `claimId` | Go field: CompleteUploadParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:575` |
| `CompleteUploadResult.OriginalFilename` | `originalFilename` | Go field: CompleteUploadResult.OriginalFilename<br/>type: string | `internal/storage/resource.go:586` |
| `CompletedTask.ApplicantName` | `applicantName` | Go field: CompletedTask.ApplicantName<br/>type: string | `approval/my/types.go:52` |
| `CompletedTask.FinishedAt` | `finishedAt` | Go field: CompletedTask.FinishedAt<br/>type: *timex.DateTime | `approval/my/types.go:55` |
| `CompletedTask.FlowIcon` | `flowIcon` | Go field: CompletedTask.FlowIcon<br/>type: *string | `approval/my/types.go:51` |
| `CompletedTask.FlowName` | `flowName` | Go field: CompletedTask.FlowName<br/>type: string | `approval/my/types.go:50` |
| `CompletedTask.InstanceID` | `instanceId` | Go field: CompletedTask.InstanceID<br/>type: string | `approval/my/types.go:47` |
| `CompletedTask.InstanceNo` | `instanceNo` | Go field: CompletedTask.InstanceNo<br/>type: string | `approval/my/types.go:49` |
| `CompletedTask.InstanceTitle` | `instanceTitle` | Go field: CompletedTask.InstanceTitle<br/>type: string | `approval/my/types.go:48` |
| `CompletedTask.NodeName` | `nodeName` | Go field: CompletedTask.NodeName<br/>type: string | `approval/my/types.go:53` |
| `CompletedTask.Status` | `status` | Go field: CompletedTask.Status<br/>type: string | `approval/my/types.go:54` |
| `CompletedTask.TaskID` | `taskId` | Go field: CompletedTask.TaskID<br/>type: string | `approval/my/types.go:46` |
| `Condition.Expression` | `expression` | Go field: Condition.Expression<br/>type: string | `approval/condition.go:11` |
| `Condition.Kind` | `kind` | Go field: Condition.Kind<br/>type: ConditionKind | `approval/condition.go:7` |
| `Condition.Operator` | `operator` | Go field: Condition.Operator<br/>type: string | `approval/condition.go:9` |
| `Condition.Subject` | `subject` | Go field: Condition.Subject<br/>type: string | `approval/condition.go:8` |
| `Condition.Value` | `value` | Go field: Condition.Value<br/>type: any | `approval/condition.go:10` |
| `ConditionBranch.ConditionGroups` | `conditionGroups` | Go field: ConditionBranch.ConditionGroups<br/>type: []ConditionGroup | `approval/condition.go:25` |
| `ConditionBranch.ID` | `id` | Go field: ConditionBranch.ID<br/>type: string | `approval/condition.go:23` |
| `ConditionBranch.IsDefault` | `isDefault` | Go field: ConditionBranch.IsDefault<br/>type: bool | `approval/condition.go:26` |
| `ConditionBranch.Label` | `label` | Go field: ConditionBranch.Label<br/>type: string | `approval/condition.go:24` |
| `ConditionBranch.Priority` | `priority` | Go field: ConditionBranch.Priority<br/>type: int | `approval/condition.go:27` |
| `ConditionGroup.Conditions` | `conditions` | Go field: ConditionGroup.Conditions<br/>type: []Condition | `approval/condition.go:17` |
| `ConditionNodeData.Branches` | `branches` | Go field: ConditionNodeData.Branches<br/>type: []ConditionBranch | `approval/node_data.go:220` |
| `CreateFlowParams.AdminUserIDs` | `adminUserIds` | Go field: CreateFlowParams.AdminUserIDs<br/>type: []string | `internal/approval/resource/flow.go:63` |
| `CreateFlowParams.BindingMode` | `bindingMode` | Go field: CreateFlowParams.BindingMode<br/>type: approval.BindingMode<br/>validate: "required" | `internal/approval/resource/flow.go:58` |
| `CreateFlowParams.BusinessPkField` | `businessPkField` | Go field: CreateFlowParams.BusinessPkField<br/>type: *string | `internal/approval/resource/flow.go:60` |
| `CreateFlowParams.BusinessStatusField` | `businessStatusField` | Go field: CreateFlowParams.BusinessStatusField<br/>type: *string | `internal/approval/resource/flow.go:62` |
| `CreateFlowParams.BusinessTable` | `businessTable` | Go field: CreateFlowParams.BusinessTable<br/>type: *string | `internal/approval/resource/flow.go:59` |
| `CreateFlowParams.BusinessTitleField` | `businessTitleField` | Go field: CreateFlowParams.BusinessTitleField<br/>type: *string | `internal/approval/resource/flow.go:61` |
| `CreateFlowParams.CategoryID` | `categoryId` | Go field: CreateFlowParams.CategoryID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:55` |
| `CreateFlowParams.Code` | `code` | Go field: CreateFlowParams.Code<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:53` |
| `CreateFlowParams.Description` | `description` | Go field: CreateFlowParams.Description<br/>type: *string | `internal/approval/resource/flow.go:57` |
| `CreateFlowParams.Icon` | `icon` | Go field: CreateFlowParams.Icon<br/>type: *string | `internal/approval/resource/flow.go:56` |
| `CreateFlowParams.Initiators` | `initiators` | Go field: CreateFlowParams.Initiators<br/>type: []CreateInitiatorParams | `internal/approval/resource/flow.go:66` |
| `CreateFlowParams.InstanceTitleTemplate` | `instanceTitleTemplate` | Go field: CreateFlowParams.InstanceTitleTemplate<br/>type: string | `internal/approval/resource/flow.go:65` |
| `CreateFlowParams.IsAllInitiationAllowed` | `isAllInitiationAllowed` | Go field: CreateFlowParams.IsAllInitiationAllowed<br/>type: bool | `internal/approval/resource/flow.go:64` |
| `CreateFlowParams.Name` | `name` | Go field: CreateFlowParams.Name<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:54` |
| `CreateFlowParams.TenantID` | `tenantId` | Go field: CreateFlowParams.TenantID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:52` |
| `CreateInitiatorParams.IDs` | `ids` | Go field: CreateInitiatorParams.IDs<br/>type: []string<br/>validate: "required" | `internal/approval/resource/flow.go:72` |
| `CreateInitiatorParams.Kind` | `kind` | Go field: CreateInitiatorParams.Kind<br/>type: approval.InitiatorKind<br/>validate: "required" | `internal/approval/resource/flow.go:71` |
| `CreateManyParams.List` | `list` | Go field: CreateManyParams.List<br/>type: []TParams<br/>validate: "required,min=1,dive" | `crud/params.go:12` |
| `CreationAuditedModel.CreatedAt` | `createdAt` | Go field: CreationAuditedModel.CreatedAt<br/>type: timex.DateTime | `internal/orm/model.go:32` |
| `CreationAuditedModel.CreatedBy` | `createdBy` | Go field: CreationAuditedModel.CreatedBy<br/>mold: "translate=user?"<br/>type: string | `internal/orm/model.go:33` |
| `CreationAuditedModel.CreatedByName` | `createdByName` | Go field: CreationAuditedModel.CreatedByName<br/>type: string | `internal/orm/model.go:34` |
| `CreationAuditedModel.ID` | `id` | Go field: CreationAuditedModel.ID<br/>type: string | `internal/orm/model.go:31` |
| `CreationTrackedModel.CreatedAt` | `createdAt` | Go field: CreationTrackedModel.CreatedAt<br/>type: timex.DateTime | `internal/orm/model.go:13` |
| `CreationTrackedModel.CreatedBy` | `createdBy` | Go field: CreationTrackedModel.CreatedBy<br/>mold: "translate=user?"<br/>type: string | `internal/orm/model.go:14` |
| `CreationTrackedModel.CreatedByName` | `createdByName` | Go field: CreationTrackedModel.CreatedByName<br/>type: string | `internal/orm/model.go:15` |
| `DataOption.Description` | `description` | Go field: DataOption.Description<br/>type: string | `crud/option.go:21` |
| `DataOption.Label` | `label` | Go field: DataOption.Label<br/>type: string | `crud/option.go:17` |
| `DataOption.Meta` | `meta` | Go field: DataOption.Meta<br/>type: map[string]any | `crud/option.go:23` |
| `DataOption.Value` | `value` | Go field: DataOption.Value<br/>type: string | `crud/option.go:19` |
| `DataOptionColumnMapping.DescriptionColumn` | `descriptionColumn` | Go field: DataOptionColumnMapping.DescriptionColumn<br/>type: string | `crud/option.go:33` |
| `DataOptionColumnMapping.LabelColumn` | `labelColumn` | Go field: DataOptionColumnMapping.LabelColumn<br/>type: string | `crud/option.go:29` |
| `DataOptionColumnMapping.MetaColumns` | `metaColumns` | Go field: DataOptionColumnMapping.MetaColumns<br/>type: []string | `crud/option.go:37` |
| `DataOptionColumnMapping.ValueColumn` | `valueColumn` | Go field: DataOptionColumnMapping.ValueColumn<br/>type: string | `crud/option.go:31` |
| `Delegation.DelegateeID` | `delegateeId` | Go field: Delegation.DelegateeID<br/>type: string | `approval/models.go:280` |
| `Delegation.DelegatorID` | `delegatorId` | Go field: Delegation.DelegatorID<br/>type: string | `approval/models.go:279` |
| `Delegation.EndTime` | `endTime` | Go field: Delegation.EndTime<br/>type: timex.DateTime | `approval/models.go:284` |
| `Delegation.FlowCategoryID` | `flowCategoryId` | Go field: Delegation.FlowCategoryID<br/>type: *string | `approval/models.go:281` |
| `Delegation.FlowID` | `flowId` | Go field: Delegation.FlowID<br/>type: *string | `approval/models.go:282` |
| `Delegation.IsActive` | `isActive` | Go field: Delegation.IsActive<br/>type: bool | `approval/models.go:285` |
| `Delegation.Reason` | `reason` | Go field: Delegation.Reason<br/>type: *string | `approval/models.go:286` |
| `Delegation.StartTime` | `startTime` | Go field: Delegation.StartTime<br/>type: timex.DateTime | `approval/models.go:283` |
| `DelegationParams.DelegateeID` | `delegateeId` | Go field: DelegationParams.DelegateeID<br/>type: string<br/>validate: "required" | `internal/approval/resource/delegation.go:21` |
| `DelegationParams.DelegatorID` | `delegatorId` | Go field: DelegationParams.DelegatorID<br/>type: string<br/>validate: "required" | `internal/approval/resource/delegation.go:20` |
| `DelegationParams.EndTime` | `endTime` | Go field: DelegationParams.EndTime<br/>type: *time.Time | `internal/approval/resource/delegation.go:25` |
| `DelegationParams.FlowCategoryID` | `flowCategoryId` | Go field: DelegationParams.FlowCategoryID<br/>type: *string | `internal/approval/resource/delegation.go:22` |
| `DelegationParams.FlowID` | `flowId` | Go field: DelegationParams.FlowID<br/>type: *string | `internal/approval/resource/delegation.go:23` |
| `DelegationParams.ID` | `id` | Go field: DelegationParams.ID<br/>type: string | `internal/approval/resource/delegation.go:19` |
| `DelegationParams.IsActive` | `isActive` | Go field: DelegationParams.IsActive<br/>type: bool | `internal/approval/resource/delegation.go:26` |
| `DelegationParams.Reason` | `reason` | Go field: DelegationParams.Reason<br/>type: *string | `internal/approval/resource/delegation.go:27` |
| `DelegationParams.StartTime` | `startTime` | Go field: DelegationParams.StartTime<br/>type: *time.Time | `internal/approval/resource/delegation.go:24` |
| `DelegationSearch.DelegateeID` | `delegateeId` | Go field: DelegationSearch.DelegateeID<br/>search: "eq,column=delegatee_id"<br/>type: string | `internal/approval/resource/delegation.go:35` |
| `DelegationSearch.DelegatorID` | `delegatorId` | Go field: DelegationSearch.DelegatorID<br/>search: "eq,column=delegator_id"<br/>type: string | `internal/approval/resource/delegation.go:34` |
| `DelegationSearch.IsActive` | `isActive` | Go field: DelegationSearch.IsActive<br/>search: "eq,column=is_active"<br/>type: *bool | `internal/approval/resource/delegation.go:36` |
| `DeleteDeadLetterEvent.Attempts` | `attempts` | Go field: DeleteDeadLetterEvent.Attempts<br/>type: int | `storage/events.go:75` |
| `DeleteDeadLetterEvent.FileKey` | `fileKey` | Go field: DeleteDeadLetterEvent.FileKey<br/>type: string | `storage/events.go:71` |
| `DeleteDeadLetterEvent.LastError` | `lastError` | Go field: DeleteDeadLetterEvent.LastError<br/>type: string | `storage/events.go:77` |
| `DeleteDeadLetterEvent.PendingDeleteID` | `pendingDeleteId` | Go field: DeleteDeadLetterEvent.PendingDeleteID<br/>type: string | `storage/events.go:69` |
| `DeleteDeadLetterEvent.Reason` | `reason` | Go field: DeleteDeadLetterEvent.Reason<br/>type: DeleteReason | `storage/events.go:73` |
| `DeleteManyParams.PKs` | `pks` | Go field: DeleteManyParams.PKs<br/>type: []any<br/>validate: "required,min=1" | `crud/params.go:28` |
| `DepartmentOption.ID` | `id` | Go field: DepartmentOption.ID<br/>type: string | `security/department_selection.go:12` |
| `DepartmentOption.Name` | `name` | Go field: DepartmentOption.Name<br/>type: string | `security/department_selection.go:13` |
| `DepartmentSelectionChallengeData.Departments` | `departments` | Go field: DepartmentSelectionChallengeData.Departments<br/>type: []DepartmentOption | `security/department_selection.go:18` |
| `DepartmentSelectionChallengeData.Meta` | `meta` | Go field: DepartmentSelectionChallengeData.Meta<br/>type: map[string]any | `security/department_selection.go:19` |
| `DeployFlowParams.Description` | `description` | Go field: DeployFlowParams.Description<br/>type: *string | `internal/approval/resource/flow.go:124` |
| `DeployFlowParams.FlowDefinition` | `flowDefinition` | Go field: DeployFlowParams.FlowDefinition<br/>type: approval.FlowDefinition<br/>validate: "required" | `internal/approval/resource/flow.go:125` |
| `DeployFlowParams.FlowID` | `flowId` | Go field: DeployFlowParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:123` |
| `DeployFlowParams.FormDefinition` | `formDefinition` | Go field: DeployFlowParams.FormDefinition<br/>type: *approval.FormDefinition | `internal/approval/resource/flow.go:126` |
| `DictionaryChangedEvent.Keys` | `keys` | Go field: DictionaryChangedEvent.Keys<br/>type: []string | `mold/cached_dictionary_resolver.go:27` |
| `DiskInfo.IOCounters` | `ioCounters` | Go field: DiskInfo.IOCounters<br/>type: map[string]*IOCounter | `monitor/service.go:148` |
| `DiskInfo.Partitions` | `partitions` | Go field: DiskInfo.Partitions<br/>type: []*PartitionInfo | `monitor/service.go:147` |
| `DiskSummary.Partitions` | `partitions` | Go field: DiskSummary.Partitions<br/>type: int | `monitor/service.go:142` |
| `DiskSummary.Total` | `total` | Go field: DiskSummary.Total<br/>type: uint64 | `monitor/service.go:139` |
| `DiskSummary.Used` | `used` | Go field: DiskSummary.Used<br/>type: uint64 | `monitor/service.go:140` |
| `DiskSummary.UsedPercent` | `usedPercent` | Go field: DiskSummary.UsedPercent<br/>type: float64 | `monitor/service.go:141` |
| `EdgeDefinition.Data` | `data` | Go field: EdgeDefinition.Data<br/>type: map[string]any | `approval/flow_definition.go:71` |
| `EdgeDefinition.ID` | `id` | Go field: EdgeDefinition.ID<br/>type: string | `approval/flow_definition.go:67` |
| `EdgeDefinition.Source` | `source` | Go field: EdgeDefinition.Source<br/>type: string | `approval/flow_definition.go:68` |
| `EdgeDefinition.SourceHandle` | `sourceHandle` | Go field: EdgeDefinition.SourceHandle<br/>type: *string | `approval/flow_definition.go:70` |
| `EdgeDefinition.Target` | `target` | Go field: EdgeDefinition.Target<br/>type: string | `approval/flow_definition.go:69` |
| `ExternalAppConfig.Enabled` | `enabled` | Go field: ExternalAppConfig.Enabled<br/>type: bool | `security/security.go:27` |
| `ExternalAppConfig.IPWhitelist` | `ipWhitelist` | Go field: ExternalAppConfig.IPWhitelist<br/>type: string | `security/security.go:28` |
| `FieldOption.Label` | `label` | Go field: FieldOption.Label<br/>type: string | `approval/form_definition.go:34` |
| `FieldOption.Value` | `value` | Go field: FieldOption.Value<br/>type: any | `approval/form_definition.go:35` |
| `FileClaimedEvent.FileKey` | `fileKey` | Go field: FileClaimedEvent.FileKey<br/>type: string | `storage/events.go:35` |
| `FileDeletedEvent.FileKey` | `fileKey` | Go field: FileDeletedEvent.FileKey<br/>type: string | `storage/events.go:51` |
| `FileDeletedEvent.Reason` | `reason` | Go field: FileDeletedEvent.Reason<br/>type: DeleteReason | `storage/events.go:53` |
| `FindAvailableFlowsParams.Keyword` | `keyword` | Go field: FindAvailableFlowsParams.Keyword<br/>type: *string | `internal/approval/resource/my.go:49` |
| `FindAvailableFlowsParams.Page` | `page` | Go field: FindAvailableFlowsParams.Page<br/>type: int | `internal/approval/resource/my.go:50` |
| `FindAvailableFlowsParams.PageSize` | `pageSize` | Go field: FindAvailableFlowsParams.PageSize<br/>type: int | `internal/approval/resource/my.go:51` |
| `FindAvailableFlowsParams.TenantID` | `tenantId` | Go field: FindAvailableFlowsParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:48` |
| `FindCCRecordsParams.IsRead` | `isRead` | Go field: FindCCRecordsParams.IsRead<br/>type: *bool | `internal/approval/resource/my.go:153` |
| `FindCCRecordsParams.Page` | `page` | Go field: FindCCRecordsParams.Page<br/>type: int | `internal/approval/resource/my.go:154` |
| `FindCCRecordsParams.PageSize` | `pageSize` | Go field: FindCCRecordsParams.PageSize<br/>type: int | `internal/approval/resource/my.go:155` |
| `FindCCRecordsParams.TenantID` | `tenantId` | Go field: FindCCRecordsParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:152` |
| `FindCompletedTasksParams.Page` | `page` | Go field: FindCompletedTasksParams.Page<br/>type: int | `internal/approval/resource/my.go:130` |
| `FindCompletedTasksParams.PageSize` | `pageSize` | Go field: FindCompletedTasksParams.PageSize<br/>type: int | `internal/approval/resource/my.go:131` |
| `FindCompletedTasksParams.TenantID` | `tenantId` | Go field: FindCompletedTasksParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:129` |
| `FindFlowsParams.CategoryID` | `categoryId` | Go field: FindFlowsParams.CategoryID<br/>type: *string | `internal/approval/resource/flow.go:219` |
| `FindFlowsParams.IsActive` | `isActive` | Go field: FindFlowsParams.IsActive<br/>type: *bool | `internal/approval/resource/flow.go:221` |
| `FindFlowsParams.Keyword` | `keyword` | Go field: FindFlowsParams.Keyword<br/>type: *string | `internal/approval/resource/flow.go:220` |
| `FindFlowsParams.Page` | `page` | Go field: FindFlowsParams.Page<br/>type: int | `internal/approval/resource/flow.go:222` |
| `FindFlowsParams.PageSize` | `pageSize` | Go field: FindFlowsParams.PageSize<br/>type: int | `internal/approval/resource/flow.go:223` |
| `FindFlowsParams.TenantID` | `tenantId` | Go field: FindFlowsParams.TenantID<br/>type: *string | `internal/approval/resource/flow.go:218` |
| `FindInitiatedParams.Keyword` | `keyword` | Go field: FindInitiatedParams.Keyword<br/>type: *string | `internal/approval/resource/my.go:81` |
| `FindInitiatedParams.Page` | `page` | Go field: FindInitiatedParams.Page<br/>type: int | `internal/approval/resource/my.go:82` |
| `FindInitiatedParams.PageSize` | `pageSize` | Go field: FindInitiatedParams.PageSize<br/>type: int | `internal/approval/resource/my.go:83` |
| `FindInitiatedParams.Status` | `status` | Go field: FindInitiatedParams.Status<br/>type: *approval.InstanceStatus | `internal/approval/resource/my.go:80` |
| `FindInitiatedParams.TenantID` | `tenantId` | Go field: FindInitiatedParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:79` |
| `FindPendingTasksParams.Page` | `page` | Go field: FindPendingTasksParams.Page<br/>type: int | `internal/approval/resource/my.go:107` |
| `FindPendingTasksParams.PageSize` | `pageSize` | Go field: FindPendingTasksParams.PageSize<br/>type: int | `internal/approval/resource/my.go:108` |
| `FindPendingTasksParams.TenantID` | `tenantId` | Go field: FindPendingTasksParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:106` |
| `FindVersionsParams.FlowID` | `flowId` | Go field: FindVersionsParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:337` |
| `FindVersionsParams.TenantID` | `tenantId` | Go field: FindVersionsParams.TenantID<br/>type: *string | `internal/approval/resource/flow.go:338` |
| `Flow.AdminUserIDs` | `adminUserIds` | Go field: Flow.AdminUserIDs<br/>type: []string | `approval/models.go:47` |
| `Flow.BindingMode` | `bindingMode` | Go field: Flow.BindingMode<br/>type: BindingMode | `approval/models.go:42` |
| `Flow.BusinessPkField` | `businessPkField` | Go field: Flow.BusinessPkField<br/>type: *string | `approval/models.go:44` |
| `Flow.BusinessStatusField` | `businessStatusField` | Go field: Flow.BusinessStatusField<br/>type: *string | `approval/models.go:46` |
| `Flow.BusinessTable` | `businessTable` | Go field: Flow.BusinessTable<br/>type: *string | `approval/models.go:43` |
| `Flow.BusinessTitleField` | `businessTitleField` | Go field: Flow.BusinessTitleField<br/>type: *string | `approval/models.go:45` |
| `Flow.CategoryID` | `categoryId` | Go field: Flow.CategoryID<br/>type: string | `approval/models.go:37` |
| `Flow.Code` | `code` | Go field: Flow.Code<br/>type: string | `approval/models.go:38` |
| `Flow.CurrentVersion` | `currentVersion` | Go field: Flow.CurrentVersion<br/>type: int | `approval/models.go:51` |
| `Flow.Description` | `description` | Go field: Flow.Description<br/>type: *string | `approval/models.go:41` |
| `Flow.Icon` | `icon` | Go field: Flow.Icon<br/>type: *string | `approval/models.go:40` |
| `Flow.InstanceTitleTemplate` | `instanceTitleTemplate` | Go field: Flow.InstanceTitleTemplate<br/>type: string | `approval/models.go:49` |
| `Flow.IsActive` | `isActive` | Go field: Flow.IsActive<br/>type: bool | `approval/models.go:50` |
| `Flow.IsAllInitiationAllowed` | `isAllInitiationAllowed` | Go field: Flow.IsAllInitiationAllowed<br/>type: bool | `approval/models.go:48` |
| `Flow.Name` | `name` | Go field: Flow.Name<br/>type: string | `approval/models.go:39` |
| `Flow.TenantID` | `tenantId` | Go field: Flow.TenantID<br/>type: string | `approval/models.go:36` |
| `FlowCategory.Children` | `children` | Go field: FlowCategory.Children<br/>type: []FlowCategory | `approval/models.go:67` |
| `FlowCategory.Code` | `code` | Go field: FlowCategory.Code<br/>type: string | `approval/models.go:60` |
| `FlowCategory.Icon` | `icon` | Go field: FlowCategory.Icon<br/>type: *string | `approval/models.go:62` |
| `FlowCategory.IsActive` | `isActive` | Go field: FlowCategory.IsActive<br/>type: bool | `approval/models.go:65` |
| `FlowCategory.Name` | `name` | Go field: FlowCategory.Name<br/>type: string | `approval/models.go:61` |
| `FlowCategory.ParentID` | `parentId` | Go field: FlowCategory.ParentID<br/>type: *string | `approval/models.go:63` |
| `FlowCategory.Remark` | `remark` | Go field: FlowCategory.Remark<br/>type: *string | `approval/models.go:66` |
| `FlowCategory.SortOrder` | `sortOrder` | Go field: FlowCategory.SortOrder<br/>type: int | `approval/models.go:64` |
| `FlowCategory.TenantID` | `tenantId` | Go field: FlowCategory.TenantID<br/>type: string | `approval/models.go:59` |
| `FlowCreatedEvent.CategoryID` | `categoryId` | Go field: FlowCreatedEvent.CategoryID<br/>type: string | `approval/events_flow.go:11` |
| `FlowCreatedEvent.Code` | `code` | Go field: FlowCreatedEvent.Code<br/>type: string | `approval/events_flow.go:9` |
| `FlowCreatedEvent.FlowID` | `flowId` | Go field: FlowCreatedEvent.FlowID<br/>type: string | `approval/events_flow.go:7` |
| `FlowCreatedEvent.Name` | `name` | Go field: FlowCreatedEvent.Name<br/>type: string | `approval/events_flow.go:10` |
| `FlowCreatedEvent.OccurredTime` | `occurredTime` | Go field: FlowCreatedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_flow.go:12` |
| `FlowCreatedEvent.TenantID` | `tenantId` | Go field: FlowCreatedEvent.TenantID<br/>type: string | `approval/events_flow.go:8` |
| `FlowDefinition.Edges` | `edges` | Go field: FlowDefinition.Edges<br/>type: []EdgeDefinition | `approval/flow_definition.go:18` |
| `FlowDefinition.Nodes` | `nodes` | Go field: FlowDefinition.Nodes<br/>type: []NodeDefinition | `approval/flow_definition.go:17` |
| `FlowDeployedEvent.FlowID` | `flowId` | Go field: FlowDeployedEvent.FlowID<br/>type: string | `approval/events_flow.go:49` |
| `FlowDeployedEvent.OccurredTime` | `occurredTime` | Go field: FlowDeployedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_flow.go:53` |
| `FlowDeployedEvent.TenantID` | `tenantId` | Go field: FlowDeployedEvent.TenantID<br/>type: string | `approval/events_flow.go:50` |
| `FlowDeployedEvent.Version` | `version` | Go field: FlowDeployedEvent.Version<br/>type: int | `approval/events_flow.go:52` |
| `FlowDeployedEvent.VersionID` | `versionId` | Go field: FlowDeployedEvent.VersionID<br/>type: string | `approval/events_flow.go:51` |
| `FlowEdge.FlowVersionID` | `flowVersionId` | Go field: FlowEdge.FlowVersionID<br/>type: string | `approval/models.go:129` |
| `FlowEdge.Key` | `key` | Go field: FlowEdge.Key<br/>type: string | `approval/models.go:130` |
| `FlowEdge.SourceHandle` | `sourceHandle` | Go field: FlowEdge.SourceHandle<br/>type: *string | `approval/models.go:135` |
| `FlowEdge.SourceNodeID` | `sourceNodeId` | Go field: FlowEdge.SourceNodeID<br/>type: string | `approval/models.go:131` |
| `FlowEdge.SourceNodeKey` | `sourceNodeKey` | Go field: FlowEdge.SourceNodeKey<br/>type: string | `approval/models.go:132` |
| `FlowEdge.TargetNodeID` | `targetNodeId` | Go field: FlowEdge.TargetNodeID<br/>type: string | `approval/models.go:133` |
| `FlowEdge.TargetNodeKey` | `targetNodeKey` | Go field: FlowEdge.TargetNodeKey<br/>type: string | `approval/models.go:134` |
| `FlowGraph.Edges` | `edges` | Go field: FlowGraph.Edges<br/>type: []approval.FlowEdge | `internal/approval/shared/flow.go:10` |
| `FlowGraph.Flow` | `flow` | Go field: FlowGraph.Flow<br/>type: *approval.Flow | `internal/approval/shared/flow.go:7` |
| `FlowGraph.Nodes` | `nodes` | Go field: FlowGraph.Nodes<br/>type: []approval.FlowNode | `internal/approval/shared/flow.go:9` |
| `FlowGraph.Version` | `version` | Go field: FlowGraph.Version<br/>type: *approval.FlowVersion | `internal/approval/shared/flow.go:8` |
| `FlowInitiator.FlowID` | `flowId` | Go field: FlowInitiator.FlowID<br/>type: string | `approval/models.go:167` |
| `FlowInitiator.IDs` | `ids` | Go field: FlowInitiator.IDs<br/>type: []string | `approval/models.go:169` |
| `FlowInitiator.Kind` | `kind` | Go field: FlowInitiator.Kind<br/>type: InitiatorKind | `approval/models.go:168` |
| `FlowNode.AddAssigneeTypes` | `addAssigneeTypes` | Go field: FlowNode.AddAssigneeTypes<br/>type: []AddAssigneeType | `approval/models.go:109` |
| `FlowNode.AdminUserIDs` | `adminUserIds` | Go field: FlowNode.AdminUserIDs<br/>type: []string | `approval/models.go:102` |
| `FlowNode.ApprovalMethod` | `approvalMethod` | Go field: FlowNode.ApprovalMethod<br/>type: ApprovalMethod | `approval/models.go:97` |
| `FlowNode.Branches` | `branches` | Go field: FlowNode.Branches<br/>type: []ConditionBranch | `approval/models.go:121` |
| `FlowNode.ConsecutiveApproverAction` | `consecutiveApproverAction` | Go field: FlowNode.ConsecutiveApproverAction<br/>type: ConsecutiveApproverAction | `approval/models.go:119` |
| `FlowNode.Description` | `description` | Go field: FlowNode.Description<br/>type: *string | `approval/models.go:95` |
| `FlowNode.EmptyAssigneeAction` | `emptyAssigneeAction` | Go field: FlowNode.EmptyAssigneeAction<br/>type: EmptyAssigneeAction | `approval/models.go:100` |
| `FlowNode.ExecutionType` | `executionType` | Go field: FlowNode.ExecutionType<br/>type: ExecutionType | `approval/models.go:96` |
| `FlowNode.FallbackUserIDs` | `fallbackUserIds` | Go field: FlowNode.FallbackUserIDs<br/>type: []string | `approval/models.go:101` |
| `FlowNode.FieldPermissions` | `fieldPermissions` | Go field: FlowNode.FieldPermissions<br/>type: map[string]Permission | `approval/models.go:111` |
| `FlowNode.FlowVersionID` | `flowVersionId` | Go field: FlowNode.FlowVersionID<br/>type: string | `approval/models.go:91` |
| `FlowNode.IsAddAssigneeAllowed` | `isAddAssigneeAllowed` | Go field: FlowNode.IsAddAssigneeAllowed<br/>type: bool | `approval/models.go:108` |
| `FlowNode.IsManualCCAllowed` | `isManualCcAllowed` | Go field: FlowNode.IsManualCCAllowed<br/>type: bool | `approval/models.go:112` |
| `FlowNode.IsOpinionRequired` | `isOpinionRequired` | Go field: FlowNode.IsOpinionRequired<br/>type: bool | `approval/models.go:114` |
| `FlowNode.IsReadConfirmRequired` | `isReadConfirmRequired` | Go field: FlowNode.IsReadConfirmRequired<br/>type: bool | `approval/models.go:120` |
| `FlowNode.IsRemoveAssigneeAllowed` | `isRemoveAssigneeAllowed` | Go field: FlowNode.IsRemoveAssigneeAllowed<br/>type: bool | `approval/models.go:110` |
| `FlowNode.IsRollbackAllowed` | `isRollbackAllowed` | Go field: FlowNode.IsRollbackAllowed<br/>type: bool | `approval/models.go:104` |
| `FlowNode.IsTransferAllowed` | `isTransferAllowed` | Go field: FlowNode.IsTransferAllowed<br/>type: bool | `approval/models.go:113` |
| `FlowNode.Key` | `key` | Go field: FlowNode.Key<br/>type: string | `approval/models.go:92` |
| `FlowNode.Kind` | `kind` | Go field: FlowNode.Kind<br/>type: NodeKind | `approval/models.go:93` |
| `FlowNode.Name` | `name` | Go field: FlowNode.Name<br/>type: string | `approval/models.go:94` |
| `FlowNode.PassRatio` | `passRatio` | Go field: FlowNode.PassRatio<br/>type: decimal.Decimal | `approval/models.go:99` |
| `FlowNode.PassRule` | `passRule` | Go field: FlowNode.PassRule<br/>type: PassRule | `approval/models.go:98` |
| `FlowNode.RollbackDataStrategy` | `rollbackDataStrategy` | Go field: FlowNode.RollbackDataStrategy<br/>type: RollbackDataStrategy | `approval/models.go:106` |
| `FlowNode.RollbackTargetKeys` | `rollbackTargetKeys` | Go field: FlowNode.RollbackTargetKeys<br/>type: []string | `approval/models.go:107` |
| `FlowNode.RollbackType` | `rollbackType` | Go field: FlowNode.RollbackType<br/>type: RollbackType | `approval/models.go:105` |
| `FlowNode.SameApplicantAction` | `sameApplicantAction` | Go field: FlowNode.SameApplicantAction<br/>type: SameApplicantAction | `approval/models.go:103` |
| `FlowNode.TimeoutAction` | `timeoutAction` | Go field: FlowNode.TimeoutAction<br/>type: TimeoutAction | `approval/models.go:116` |
| `FlowNode.TimeoutHours` | `timeoutHours` | Go field: FlowNode.TimeoutHours<br/>type: int | `approval/models.go:115` |
| `FlowNode.TimeoutNotifyBeforeHours` | `timeoutNotifyBeforeHours` | Go field: FlowNode.TimeoutNotifyBeforeHours<br/>type: int | `approval/models.go:117` |
| `FlowNode.UrgeCooldownMinutes` | `urgeCooldownMinutes` | Go field: FlowNode.UrgeCooldownMinutes<br/>type: int | `approval/models.go:118` |
| `FlowNodeAssignee.FormField` | `formField` | Go field: FlowNodeAssignee.FormField<br/>type: *string | `approval/models.go:146` |
| `FlowNodeAssignee.IDs` | `ids` | Go field: FlowNodeAssignee.IDs<br/>type: []string | `approval/models.go:145` |
| `FlowNodeAssignee.Kind` | `kind` | Go field: FlowNodeAssignee.Kind<br/>type: AssigneeKind | `approval/models.go:144` |
| `FlowNodeAssignee.NodeID` | `nodeId` | Go field: FlowNodeAssignee.NodeID<br/>type: string | `approval/models.go:143` |
| `FlowNodeAssignee.SortOrder` | `sortOrder` | Go field: FlowNodeAssignee.SortOrder<br/>type: int | `approval/models.go:147` |
| `FlowNodeCC.FormField` | `formField` | Go field: FlowNodeCC.FormField<br/>type: *string | `approval/models.go:158` |
| `FlowNodeCC.IDs` | `ids` | Go field: FlowNodeCC.IDs<br/>type: []string | `approval/models.go:157` |
| `FlowNodeCC.Kind` | `kind` | Go field: FlowNodeCC.Kind<br/>type: CCKind | `approval/models.go:156` |
| `FlowNodeCC.NodeID` | `nodeId` | Go field: FlowNodeCC.NodeID<br/>type: string | `approval/models.go:155` |
| `FlowNodeCC.Timing` | `timing` | Go field: FlowNodeCC.Timing<br/>type: CCTiming | `approval/models.go:159` |
| `FlowNodeInfo.ExecutionType` | `executionType` | Go field: FlowNodeInfo.ExecutionType<br/>type: string | `approval/admin/types.go:99` |
| `FlowNodeInfo.Key` | `key` | Go field: FlowNodeInfo.Key<br/>type: string | `approval/my/types.go:127` |
| `FlowNodeInfo.Key` | `key` | Go field: FlowNodeInfo.Key<br/>type: string | `approval/admin/types.go:96` |
| `FlowNodeInfo.Kind` | `kind` | Go field: FlowNodeInfo.Kind<br/>type: string | `approval/admin/types.go:97` |
| `FlowNodeInfo.Kind` | `kind` | Go field: FlowNodeInfo.Kind<br/>type: string | `approval/my/types.go:128` |
| `FlowNodeInfo.Name` | `name` | Go field: FlowNodeInfo.Name<br/>type: string | `approval/admin/types.go:98` |
| `FlowNodeInfo.Name` | `name` | Go field: FlowNodeInfo.Name<br/>type: string | `approval/my/types.go:129` |
| `FlowNodeInfo.NodeID` | `nodeId` | Go field: FlowNodeInfo.NodeID<br/>type: string | `approval/admin/types.go:95` |
| `FlowNodeInfo.NodeID` | `nodeId` | Go field: FlowNodeInfo.NodeID<br/>type: string | `approval/my/types.go:126` |
| `FlowPublishedEvent.FlowID` | `flowId` | Go field: FlowPublishedEvent.FlowID<br/>type: string | `approval/events_flow.go:89` |
| `FlowPublishedEvent.OccurredTime` | `occurredTime` | Go field: FlowPublishedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_flow.go:92` |
| `FlowPublishedEvent.TenantID` | `tenantId` | Go field: FlowPublishedEvent.TenantID<br/>type: string | `approval/events_flow.go:90` |
| `FlowPublishedEvent.VersionID` | `versionId` | Go field: FlowPublishedEvent.VersionID<br/>type: string | `approval/events_flow.go:91` |
| `FlowToggledEvent.FlowID` | `flowId` | Go field: FlowToggledEvent.FlowID<br/>type: string | `approval/events_flow.go:70` |
| `FlowToggledEvent.IsActive` | `isActive` | Go field: FlowToggledEvent.IsActive<br/>type: bool | `approval/events_flow.go:72` |
| `FlowToggledEvent.OccurredTime` | `occurredTime` | Go field: FlowToggledEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_flow.go:73` |
| `FlowToggledEvent.TenantID` | `tenantId` | Go field: FlowToggledEvent.TenantID<br/>type: string | `approval/events_flow.go:71` |
| `FlowUpdatedEvent.FlowID` | `flowId` | Go field: FlowUpdatedEvent.FlowID<br/>type: string | `approval/events_flow.go:31` |
| `FlowUpdatedEvent.OccurredTime` | `occurredTime` | Go field: FlowUpdatedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_flow.go:33` |
| `FlowUpdatedEvent.TenantID` | `tenantId` | Go field: FlowUpdatedEvent.TenantID<br/>type: string | `approval/events_flow.go:32` |
| `FlowVersion.Description` | `description` | Go field: FlowVersion.Description<br/>type: *string | `approval/models.go:78` |
| `FlowVersion.FlowID` | `flowId` | Go field: FlowVersion.FlowID<br/>type: string | `approval/models.go:75` |
| `FlowVersion.FlowSchema` | `flowSchema` | Go field: FlowVersion.FlowSchema<br/>type: *FlowDefinition | `approval/models.go:80` |
| `FlowVersion.FormSchema` | `formSchema` | Go field: FlowVersion.FormSchema<br/>type: *FormDefinition | `approval/models.go:81` |
| `FlowVersion.PublishedAt` | `publishedAt` | Go field: FlowVersion.PublishedAt<br/>type: *timex.DateTime | `approval/models.go:82` |
| `FlowVersion.PublishedBy` | `publishedBy` | Go field: FlowVersion.PublishedBy<br/>type: *string | `approval/models.go:83` |
| `FlowVersion.Status` | `status` | Go field: FlowVersion.Status<br/>type: VersionStatus | `approval/models.go:77` |
| `FlowVersion.StorageMode` | `storageMode` | Go field: FlowVersion.StorageMode<br/>type: StorageMode | `approval/models.go:79` |
| `FlowVersion.Version` | `version` | Go field: FlowVersion.Version<br/>type: int | `approval/models.go:76` |
| `ForeignKey.Columns` | `columns` | Go field: ForeignKey.Columns<br/>type: []string | `schema/service.go:57` |
| `ForeignKey.Name` | `name` | Go field: ForeignKey.Name<br/>type: string | `schema/service.go:56` |
| `ForeignKey.OnDelete` | `onDelete` | Go field: ForeignKey.OnDelete<br/>type: string | `schema/service.go:61` |
| `ForeignKey.OnUpdate` | `onUpdate` | Go field: ForeignKey.OnUpdate<br/>type: string | `schema/service.go:60` |
| `ForeignKey.RefColumns` | `refColumns` | Go field: ForeignKey.RefColumns<br/>type: []string | `schema/service.go:59` |
| `ForeignKey.RefTable` | `refTable` | Go field: ForeignKey.RefTable<br/>type: string | `schema/service.go:58` |
| `FormDefinition.Fields` | `fields` | Go field: FormDefinition.Fields<br/>type: []FormFieldDefinition | `approval/form_definition.go:5` |
| `FormFieldDefinition.DefaultValue` | `defaultValue` | Go field: FormFieldDefinition.DefaultValue<br/>type: any | `approval/form_definition.go:19` |
| `FormFieldDefinition.IsRequired` | `isRequired` | Go field: FormFieldDefinition.IsRequired<br/>type: bool | `approval/form_definition.go:21` |
| `FormFieldDefinition.Key` | `key` | Go field: FormFieldDefinition.Key<br/>type: string | `approval/form_definition.go:11` |
| `FormFieldDefinition.Kind` | `kind` | Go field: FormFieldDefinition.Kind<br/>type: FieldKind | `approval/form_definition.go:13` |
| `FormFieldDefinition.Label` | `label` | Go field: FormFieldDefinition.Label<br/>type: string | `approval/form_definition.go:15` |
| `FormFieldDefinition.Options` | `options` | Go field: FormFieldDefinition.Options<br/>type: []FieldOption | `approval/form_definition.go:23` |
| `FormFieldDefinition.Placeholder` | `placeholder` | Go field: FormFieldDefinition.Placeholder<br/>type: string | `approval/form_definition.go:17` |
| `FormFieldDefinition.Props` | `props` | Go field: FormFieldDefinition.Props<br/>type: map[string]any | `approval/form_definition.go:27` |
| `FormFieldDefinition.SortOrder` | `sortOrder` | Go field: FormFieldDefinition.SortOrder<br/>type: int | `approval/form_definition.go:29` |
| `FormFieldDefinition.Validation` | `validation` | Go field: FormFieldDefinition.Validation<br/>type: *ValidationRule | `approval/form_definition.go:25` |
| `FormSnapshot.FormData` | `formData` | Go field: FormSnapshot.FormData<br/>type: map[string]any | `approval/models.go:226` |
| `FormSnapshot.InstanceID` | `instanceId` | Go field: FormSnapshot.InstanceID<br/>type: string | `approval/models.go:224` |
| `FormSnapshot.NodeID` | `nodeId` | Go field: FormSnapshot.NodeID<br/>type: string | `approval/models.go:225` |
| `FullAuditedModel.CreatedAt` | `createdAt` | Go field: FullAuditedModel.CreatedAt<br/>type: timex.DateTime | `internal/orm/model.go:40` |
| `FullAuditedModel.CreatedBy` | `createdBy` | Go field: FullAuditedModel.CreatedBy<br/>mold: "translate=user?"<br/>type: string | `internal/orm/model.go:41` |
| `FullAuditedModel.CreatedByName` | `createdByName` | Go field: FullAuditedModel.CreatedByName<br/>type: string | `internal/orm/model.go:42` |
| `FullAuditedModel.ID` | `id` | Go field: FullAuditedModel.ID<br/>type: string | `internal/orm/model.go:39` |
| `FullAuditedModel.UpdatedAt` | `updatedAt` | Go field: FullAuditedModel.UpdatedAt<br/>type: timex.DateTime | `internal/orm/model.go:43` |
| `FullAuditedModel.UpdatedBy` | `updatedBy` | Go field: FullAuditedModel.UpdatedBy<br/>mold: "translate=user?"<br/>type: string | `internal/orm/model.go:44` |
| `FullAuditedModel.UpdatedByName` | `updatedByName` | Go field: FullAuditedModel.UpdatedByName<br/>type: string | `internal/orm/model.go:45` |
| `FullTrackedModel.CreatedAt` | `createdAt` | Go field: FullTrackedModel.CreatedAt<br/>type: timex.DateTime | `internal/orm/model.go:21` |
| `FullTrackedModel.CreatedBy` | `createdBy` | Go field: FullTrackedModel.CreatedBy<br/>mold: "translate=user?"<br/>type: string | `internal/orm/model.go:22` |
| `FullTrackedModel.CreatedByName` | `createdByName` | Go field: FullTrackedModel.CreatedByName<br/>type: string | `internal/orm/model.go:23` |
| `FullTrackedModel.UpdatedAt` | `updatedAt` | Go field: FullTrackedModel.UpdatedAt<br/>type: timex.DateTime | `internal/orm/model.go:24` |
| `FullTrackedModel.UpdatedBy` | `updatedBy` | Go field: FullTrackedModel.UpdatedBy<br/>mold: "translate=user?"<br/>type: string | `internal/orm/model.go:25` |
| `FullTrackedModel.UpdatedByName` | `updatedByName` | Go field: FullTrackedModel.UpdatedByName<br/>type: string | `internal/orm/model.go:26` |
| `GetGraphParams.FlowID` | `flowId` | Go field: GetGraphParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:187` |
| `GetGraphParams.TenantID` | `tenantId` | Go field: GetGraphParams.TenantID<br/>type: string | `internal/approval/resource/flow.go:188` |
| `GetInstanceDetailParams.InstanceID` | `instanceId` | Go field: GetInstanceDetailParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/my.go:197` |
| `GetPendingCountsParams.TenantID` | `tenantId` | Go field: GetPendingCountsParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:177` |
| `GetTableSchemaParams.Name` | `name` | Go field: GetTableSchemaParams.Name<br/>type: string<br/>validate: "required" | `internal/schema/resource.go:49` |
| `HostInfo.BootTime` | `bootTime` | Go field: HostInfo.BootTime<br/>type: uint64 | `monitor/service.go:32` |
| `HostInfo.HostID` | `hostId` | Go field: HostInfo.HostID<br/>type: string | `monitor/service.go:42` |
| `HostInfo.Hostname` | `hostname` | Go field: HostInfo.Hostname<br/>type: string | `monitor/service.go:30` |
| `HostInfo.KernelArch` | `kernelArch` | Go field: HostInfo.KernelArch<br/>type: string | `monitor/service.go:39` |
| `HostInfo.KernelVersion` | `kernelVersion` | Go field: HostInfo.KernelVersion<br/>type: string | `monitor/service.go:38` |
| `HostInfo.OS` | `os` | Go field: HostInfo.OS<br/>type: string | `monitor/service.go:34` |
| `HostInfo.Platform` | `platform` | Go field: HostInfo.Platform<br/>type: string | `monitor/service.go:35` |
| `HostInfo.PlatformFamily` | `platformFamily` | Go field: HostInfo.PlatformFamily<br/>type: string | `monitor/service.go:36` |
| `HostInfo.PlatformVersion` | `platformVersion` | Go field: HostInfo.PlatformVersion<br/>type: string | `monitor/service.go:37` |
| `HostInfo.Processes` | `processes` | Go field: HostInfo.Processes<br/>type: uint64 | `monitor/service.go:33` |
| `HostInfo.Uptime` | `uptime` | Go field: HostInfo.Uptime<br/>type: uint64 | `monitor/service.go:31` |
| `HostInfo.VirtualizationRole` | `virtualizationRole` | Go field: HostInfo.VirtualizationRole<br/>type: string | `monitor/service.go:41` |
| `HostInfo.VirtualizationSystem` | `virtualizationSystem` | Go field: HostInfo.VirtualizationSystem<br/>type: string | `monitor/service.go:40` |
| `HostSummary.Hostname` | `hostname` | Go field: HostSummary.Hostname<br/>type: string | `monitor/service.go:19` |
| `HostSummary.KernelArch` | `kernelArch` | Go field: HostSummary.KernelArch<br/>type: string | `monitor/service.go:24` |
| `HostSummary.KernelVersion` | `kernelVersion` | Go field: HostSummary.KernelVersion<br/>type: string | `monitor/service.go:23` |
| `HostSummary.OS` | `os` | Go field: HostSummary.OS<br/>type: string | `monitor/service.go:20` |
| `HostSummary.Platform` | `platform` | Go field: HostSummary.Platform<br/>type: string | `monitor/service.go:21` |
| `HostSummary.PlatformVersion` | `platformVersion` | Go field: HostSummary.PlatformVersion<br/>type: string | `monitor/service.go:22` |
| `HostSummary.Uptime` | `uptime` | Go field: HostSummary.Uptime<br/>type: uint64 | `monitor/service.go:25` |
| `IOCounter.IOPSInProgress` | `iopsInProgress` | Go field: IOCounter.IOPSInProgress<br/>type: uint64 | `monitor/service.go:177` |
| `IOCounter.IOTime` | `ioTime` | Go field: IOCounter.IOTime<br/>type: uint64 | `monitor/service.go:178` |
| `IOCounter.Label` | `label` | Go field: IOCounter.Label<br/>type: string | `monitor/service.go:182` |
| `IOCounter.MergedReadCount` | `mergedReadCount` | Go field: IOCounter.MergedReadCount<br/>type: uint64 | `monitor/service.go:170` |
| `IOCounter.MergedWriteCount` | `mergedWriteCount` | Go field: IOCounter.MergedWriteCount<br/>type: uint64 | `monitor/service.go:172` |
| `IOCounter.Name` | `name` | Go field: IOCounter.Name<br/>type: string | `monitor/service.go:180` |
| `IOCounter.ReadBytes` | `readBytes` | Go field: IOCounter.ReadBytes<br/>type: uint64 | `monitor/service.go:173` |
| `IOCounter.ReadCount` | `readCount` | Go field: IOCounter.ReadCount<br/>type: uint64 | `monitor/service.go:169` |
| `IOCounter.ReadTime` | `readTime` | Go field: IOCounter.ReadTime<br/>type: uint64 | `monitor/service.go:175` |
| `IOCounter.SerialNumber` | `serialNumber` | Go field: IOCounter.SerialNumber<br/>type: string | `monitor/service.go:181` |
| `IOCounter.WeightedIO` | `weightedIo` | Go field: IOCounter.WeightedIO<br/>type: uint64 | `monitor/service.go:179` |
| `IOCounter.WriteBytes` | `writeBytes` | Go field: IOCounter.WriteBytes<br/>type: uint64 | `monitor/service.go:174` |
| `IOCounter.WriteCount` | `writeCount` | Go field: IOCounter.WriteCount<br/>type: uint64 | `monitor/service.go:171` |
| `IOCounter.WriteTime` | `writeTime` | Go field: IOCounter.WriteTime<br/>type: uint64 | `monitor/service.go:176` |
| `Identifier.Action` | `action` | Go field: Identifier.Action<br/>form: "action"<br/>type: string<br/>validate: "required" | `api/request.go:15` |
| `Identifier.Resource` | `resource` | Go field: Identifier.Resource<br/>form: "resource"<br/>type: string<br/>validate: "required,alphanum_us_slash" | `api/request.go:14` |
| `Identifier.Version` | `version` | Go field: Identifier.Version<br/>form: "version"<br/>type: string<br/>validate: "required,alphanum" | `api/request.go:16` |
| `Index.Columns` | `columns` | Go field: Index.Columns<br/>type: []string | `schema/service.go:45` |
| `Index.Name` | `name` | Go field: Index.Name<br/>type: string | `schema/service.go:44` |
| `InitUploadParams.ContentType` | `contentType` | Go field: InitUploadParams.ContentType<br/>type: string<br/>validate: "max=127" | `internal/storage/resource.go:241` |
| `InitUploadParams.Filename` | `filename` | Go field: InitUploadParams.Filename<br/>type: string<br/>validate: "required,max=255" | `internal/storage/resource.go:239` |
| `InitUploadParams.Public` | `public` | Go field: InitUploadParams.Public<br/>type: bool | `internal/storage/resource.go:242` |
| `InitUploadParams.Size` | `size` | Go field: InitUploadParams.Size<br/>type: int64<br/>validate: "required,min=1" | `internal/storage/resource.go:240` |
| `InitUploadResult.ClaimID` | `claimId` | Go field: InitUploadResult.ClaimID<br/>type: string | `internal/storage/resource.go:258` |
| `InitUploadResult.ExpiresAt` | `expiresAt` | Go field: InitUploadResult.ExpiresAt<br/>type: time.Time | `internal/storage/resource.go:262` |
| `InitUploadResult.Key` | `key` | Go field: InitUploadResult.Key<br/>type: string | `internal/storage/resource.go:257` |
| `InitUploadResult.OriginalFilename` | `originalFilename` | Go field: InitUploadResult.OriginalFilename<br/>type: string | `internal/storage/resource.go:259` |
| `InitUploadResult.PartCount` | `partCount` | Go field: InitUploadResult.PartCount<br/>type: int | `internal/storage/resource.go:261` |
| `InitUploadResult.PartSize` | `partSize` | Go field: InitUploadResult.PartSize<br/>type: int64 | `internal/storage/resource.go:260` |
| `InitiatedInstance.CreatedAt` | `createdAt` | Go field: InitiatedInstance.CreatedAt<br/>type: timex.DateTime | `approval/my/types.go:25` |
| `InitiatedInstance.CurrentNodeName` | `currentNodeName` | Go field: InitiatedInstance.CurrentNodeName<br/>type: *string | `approval/my/types.go:24` |
| `InitiatedInstance.FinishedAt` | `finishedAt` | Go field: InitiatedInstance.FinishedAt<br/>type: *timex.DateTime | `approval/my/types.go:26` |
| `InitiatedInstance.FlowIcon` | `flowIcon` | Go field: InitiatedInstance.FlowIcon<br/>type: *string | `approval/my/types.go:22` |
| `InitiatedInstance.FlowName` | `flowName` | Go field: InitiatedInstance.FlowName<br/>type: string | `approval/my/types.go:21` |
| `InitiatedInstance.InstanceID` | `instanceId` | Go field: InitiatedInstance.InstanceID<br/>type: string | `approval/my/types.go:18` |
| `InitiatedInstance.InstanceNo` | `instanceNo` | Go field: InitiatedInstance.InstanceNo<br/>type: string | `approval/my/types.go:19` |
| `InitiatedInstance.Status` | `status` | Go field: InitiatedInstance.Status<br/>type: string | `approval/my/types.go:23` |
| `InitiatedInstance.Title` | `title` | Go field: InitiatedInstance.Title<br/>type: string | `approval/my/types.go:20` |
| `Instance.ApplicantDepartmentID` | `applicantDepartmentId` | Go field: Instance.ApplicantDepartmentID<br/>type: *string | `approval/models.go:186` |
| `Instance.ApplicantDepartmentName` | `applicantDepartmentName` | Go field: Instance.ApplicantDepartmentName<br/>type: *string | `approval/models.go:187` |
| `Instance.ApplicantID` | `applicantId` | Go field: Instance.ApplicantID<br/>type: string | `approval/models.go:184` |
| `Instance.ApplicantID` | `applicantId` | Go field: Instance.ApplicantID<br/>type: string | `approval/admin/types.go:13` |
| `Instance.ApplicantName` | `applicantName` | Go field: Instance.ApplicantName<br/>type: string | `approval/admin/types.go:14` |
| `Instance.ApplicantName` | `applicantName` | Go field: Instance.ApplicantName<br/>type: string | `approval/models.go:185` |
| `Instance.BusinessRecordID` | `businessRecordId` | Go field: Instance.BusinessRecordID<br/>type: *string | `approval/models.go:191` |
| `Instance.CreatedAt` | `createdAt` | Go field: Instance.CreatedAt<br/>type: timex.DateTime | `approval/admin/types.go:17` |
| `Instance.CurrentNodeID` | `currentNodeId` | Go field: Instance.CurrentNodeID<br/>type: *string | `approval/models.go:189` |
| `Instance.CurrentNodeName` | `currentNodeName` | Go field: Instance.CurrentNodeName<br/>type: *string | `approval/admin/types.go:16` |
| `Instance.FinishedAt` | `finishedAt` | Go field: Instance.FinishedAt<br/>type: *timex.DateTime | `approval/admin/types.go:18` |
| `Instance.FinishedAt` | `finishedAt` | Go field: Instance.FinishedAt<br/>type: *timex.DateTime | `approval/models.go:190` |
| `Instance.FlowID` | `flowId` | Go field: Instance.FlowID<br/>type: string | `approval/models.go:180` |
| `Instance.FlowID` | `flowId` | Go field: Instance.FlowID<br/>type: string | `approval/admin/types.go:11` |
| `Instance.FlowName` | `flowName` | Go field: Instance.FlowName<br/>type: string | `approval/admin/types.go:12` |
| `Instance.FlowVersionID` | `flowVersionId` | Go field: Instance.FlowVersionID<br/>type: string | `approval/models.go:181` |
| `Instance.FormData` | `formData` | Go field: Instance.FormData<br/>type: map[string]any | `approval/models.go:192` |
| `Instance.InstanceID` | `instanceId` | Go field: Instance.InstanceID<br/>type: string | `approval/admin/types.go:7` |
| `Instance.InstanceNo` | `instanceNo` | Go field: Instance.InstanceNo<br/>type: string | `approval/admin/types.go:8` |
| `Instance.InstanceNo` | `instanceNo` | Go field: Instance.InstanceNo<br/>type: string | `approval/models.go:183` |
| `Instance.Status` | `status` | Go field: Instance.Status<br/>type: InstanceStatus | `approval/models.go:188` |
| `Instance.Status` | `status` | Go field: Instance.Status<br/>type: string | `approval/admin/types.go:15` |
| `Instance.TenantID` | `tenantId` | Go field: Instance.TenantID<br/>type: string | `approval/admin/types.go:10` |
| `Instance.TenantID` | `tenantId` | Go field: Instance.TenantID<br/>type: string | `approval/models.go:179` |
| `Instance.Title` | `title` | Go field: Instance.Title<br/>type: string | `approval/models.go:182` |
| `Instance.Title` | `title` | Go field: Instance.Title<br/>type: string | `approval/admin/types.go:9` |
| `InstanceBindingFailedEvent.BusinessTable` | `businessTable` | Go field: InstanceBindingFailedEvent.BusinessTable<br/>type: string | `approval/events_instance.go:146` |
| `InstanceBindingFailedEvent.ErrorMessage` | `errorMessage` | Go field: InstanceBindingFailedEvent.ErrorMessage<br/>type: string | `approval/events_instance.go:147` |
| `InstanceBindingFailedEvent.FinalStatus` | `finalStatus` | Go field: InstanceBindingFailedEvent.FinalStatus<br/>type: InstanceStatus | `approval/events_instance.go:145` |
| `InstanceBindingFailedEvent.FlowID` | `flowId` | Go field: InstanceBindingFailedEvent.FlowID<br/>type: string | `approval/events_instance.go:144` |
| `InstanceBindingFailedEvent.InstanceID` | `instanceId` | Go field: InstanceBindingFailedEvent.InstanceID<br/>type: string | `approval/events_instance.go:142` |
| `InstanceBindingFailedEvent.OccurredTime` | `occurredTime` | Go field: InstanceBindingFailedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_instance.go:148` |
| `InstanceBindingFailedEvent.TenantID` | `tenantId` | Go field: InstanceBindingFailedEvent.TenantID<br/>type: string | `approval/events_instance.go:143` |
| `InstanceCompletedEvent.FinalStatus` | `finalStatus` | Go field: InstanceCompletedEvent.FinalStatus<br/>type: InstanceStatus | `approval/events_instance.go:34` |
| `InstanceCompletedEvent.FinishedAt` | `finishedAt` | Go field: InstanceCompletedEvent.FinishedAt<br/>type: timex.DateTime | `approval/events_instance.go:35` |
| `InstanceCompletedEvent.InstanceID` | `instanceId` | Go field: InstanceCompletedEvent.InstanceID<br/>type: string | `approval/events_instance.go:32` |
| `InstanceCompletedEvent.OccurredTime` | `occurredTime` | Go field: InstanceCompletedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_instance.go:36` |
| `InstanceCompletedEvent.TenantID` | `tenantId` | Go field: InstanceCompletedEvent.TenantID<br/>type: string | `approval/events_instance.go:33` |
| `InstanceCreatedEvent.ApplicantID` | `applicantId` | Go field: InstanceCreatedEvent.ApplicantID<br/>type: string | `approval/events_instance.go:11` |
| `InstanceCreatedEvent.ApplicantName` | `applicantName` | Go field: InstanceCreatedEvent.ApplicantName<br/>type: string | `approval/events_instance.go:12` |
| `InstanceCreatedEvent.FlowID` | `flowId` | Go field: InstanceCreatedEvent.FlowID<br/>type: string | `approval/events_instance.go:9` |
| `InstanceCreatedEvent.InstanceID` | `instanceId` | Go field: InstanceCreatedEvent.InstanceID<br/>type: string | `approval/events_instance.go:7` |
| `InstanceCreatedEvent.OccurredTime` | `occurredTime` | Go field: InstanceCreatedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_instance.go:13` |
| `InstanceCreatedEvent.TenantID` | `tenantId` | Go field: InstanceCreatedEvent.TenantID<br/>type: string | `approval/events_instance.go:8` |
| `InstanceCreatedEvent.Title` | `title` | Go field: InstanceCreatedEvent.Title<br/>type: string | `approval/events_instance.go:10` |
| `InstanceDetail.ActionLogs` | `actionLogs` | Go field: InstanceDetail.ActionLogs<br/>type: []ActionLogInfo | `approval/my/types.go:82` |
| `InstanceDetail.ActionLogs` | `actionLogs` | Go field: InstanceDetail.ActionLogs<br/>type: []ActionLog | `approval/admin/types.go:40` |
| `InstanceDetail.AvailableActions` | `availableActions` | Go field: InstanceDetail.AvailableActions<br/>type: []string | `approval/my/types.go:84` |
| `InstanceDetail.FlowNodes` | `flowNodes` | Go field: InstanceDetail.FlowNodes<br/>type: []FlowNodeInfo | `approval/admin/types.go:41` |
| `InstanceDetail.FlowNodes` | `flowNodes` | Go field: InstanceDetail.FlowNodes<br/>type: []FlowNodeInfo | `approval/my/types.go:83` |
| `InstanceDetail.Instance` | `instance` | Go field: InstanceDetail.Instance<br/>type: InstanceInfo | `approval/my/types.go:80` |
| `InstanceDetail.Instance` | `instance` | Go field: InstanceDetail.Instance<br/>type: InstanceDetailInfo | `approval/admin/types.go:38` |
| `InstanceDetail.Tasks` | `tasks` | Go field: InstanceDetail.Tasks<br/>type: []TaskDetailInfo | `approval/admin/types.go:39` |
| `InstanceDetail.Tasks` | `tasks` | Go field: InstanceDetail.Tasks<br/>type: []TaskInfo | `approval/my/types.go:81` |
| `InstanceDetailInfo.ApplicantID` | `applicantId` | Go field: InstanceDetailInfo.ApplicantID<br/>type: string | `approval/admin/types.go:53` |
| `InstanceDetailInfo.ApplicantName` | `applicantName` | Go field: InstanceDetailInfo.ApplicantName<br/>type: string | `approval/admin/types.go:54` |
| `InstanceDetailInfo.BusinessRecordID` | `businessRecordId` | Go field: InstanceDetailInfo.BusinessRecordID<br/>type: *string | `approval/admin/types.go:57` |
| `InstanceDetailInfo.CreatedAt` | `createdAt` | Go field: InstanceDetailInfo.CreatedAt<br/>type: timex.DateTime | `approval/admin/types.go:59` |
| `InstanceDetailInfo.CurrentNodeName` | `currentNodeName` | Go field: InstanceDetailInfo.CurrentNodeName<br/>type: *string | `approval/admin/types.go:56` |
| `InstanceDetailInfo.FinishedAt` | `finishedAt` | Go field: InstanceDetailInfo.FinishedAt<br/>type: *timex.DateTime | `approval/admin/types.go:60` |
| `InstanceDetailInfo.FlowID` | `flowId` | Go field: InstanceDetailInfo.FlowID<br/>type: string | `approval/admin/types.go:50` |
| `InstanceDetailInfo.FlowName` | `flowName` | Go field: InstanceDetailInfo.FlowName<br/>type: string | `approval/admin/types.go:51` |
| `InstanceDetailInfo.FlowVersionID` | `flowVersionId` | Go field: InstanceDetailInfo.FlowVersionID<br/>type: string | `approval/admin/types.go:52` |
| `InstanceDetailInfo.FormData` | `formData` | Go field: InstanceDetailInfo.FormData<br/>type: map[string]any | `approval/admin/types.go:58` |
| `InstanceDetailInfo.InstanceID` | `instanceId` | Go field: InstanceDetailInfo.InstanceID<br/>type: string | `approval/admin/types.go:46` |
| `InstanceDetailInfo.InstanceNo` | `instanceNo` | Go field: InstanceDetailInfo.InstanceNo<br/>type: string | `approval/admin/types.go:47` |
| `InstanceDetailInfo.Status` | `status` | Go field: InstanceDetailInfo.Status<br/>type: string | `approval/admin/types.go:55` |
| `InstanceDetailInfo.TenantID` | `tenantId` | Go field: InstanceDetailInfo.TenantID<br/>type: string | `approval/admin/types.go:49` |
| `InstanceDetailInfo.Title` | `title` | Go field: InstanceDetailInfo.Title<br/>type: string | `approval/admin/types.go:48` |
| `InstanceInfo.ApplicantID` | `applicantId` | Go field: InstanceInfo.ApplicantID<br/>type: string | `approval/my/types.go:94` |
| `InstanceInfo.ApplicantName` | `applicantName` | Go field: InstanceInfo.ApplicantName<br/>type: string | `approval/my/types.go:95` |
| `InstanceInfo.BusinessRecordID` | `businessRecordId` | Go field: InstanceInfo.BusinessRecordID<br/>type: *string | `approval/my/types.go:98` |
| `InstanceInfo.CreatedAt` | `createdAt` | Go field: InstanceInfo.CreatedAt<br/>type: timex.DateTime | `approval/my/types.go:100` |
| `InstanceInfo.CurrentNodeName` | `currentNodeName` | Go field: InstanceInfo.CurrentNodeName<br/>type: *string | `approval/my/types.go:97` |
| `InstanceInfo.FinishedAt` | `finishedAt` | Go field: InstanceInfo.FinishedAt<br/>type: *timex.DateTime | `approval/my/types.go:101` |
| `InstanceInfo.FlowIcon` | `flowIcon` | Go field: InstanceInfo.FlowIcon<br/>type: *string | `approval/my/types.go:93` |
| `InstanceInfo.FlowName` | `flowName` | Go field: InstanceInfo.FlowName<br/>type: string | `approval/my/types.go:92` |
| `InstanceInfo.FormData` | `formData` | Go field: InstanceInfo.FormData<br/>type: map[string]any | `approval/my/types.go:99` |
| `InstanceInfo.InstanceID` | `instanceId` | Go field: InstanceInfo.InstanceID<br/>type: string | `approval/my/types.go:89` |
| `InstanceInfo.InstanceNo` | `instanceNo` | Go field: InstanceInfo.InstanceNo<br/>type: string | `approval/my/types.go:90` |
| `InstanceInfo.Status` | `status` | Go field: InstanceInfo.Status<br/>type: string | `approval/my/types.go:96` |
| `InstanceInfo.Title` | `title` | Go field: InstanceInfo.Title<br/>type: string | `approval/my/types.go:91` |
| `InstanceResubmittedEvent.InstanceID` | `instanceId` | Go field: InstanceResubmittedEvent.InstanceID<br/>type: string | `approval/events_instance.go:120` |
| `InstanceResubmittedEvent.OccurredTime` | `occurredTime` | Go field: InstanceResubmittedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_instance.go:123` |
| `InstanceResubmittedEvent.OperatorID` | `operatorId` | Go field: InstanceResubmittedEvent.OperatorID<br/>type: string | `approval/events_instance.go:122` |
| `InstanceResubmittedEvent.TenantID` | `tenantId` | Go field: InstanceResubmittedEvent.TenantID<br/>type: string | `approval/events_instance.go:121` |
| `InstanceReturnedEvent.FromNodeID` | `fromNodeId` | Go field: InstanceReturnedEvent.FromNodeID<br/>type: string | `approval/events_instance.go:99` |
| `InstanceReturnedEvent.InstanceID` | `instanceId` | Go field: InstanceReturnedEvent.InstanceID<br/>type: string | `approval/events_instance.go:97` |
| `InstanceReturnedEvent.OccurredTime` | `occurredTime` | Go field: InstanceReturnedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_instance.go:102` |
| `InstanceReturnedEvent.OperatorID` | `operatorId` | Go field: InstanceReturnedEvent.OperatorID<br/>type: string | `approval/events_instance.go:101` |
| `InstanceReturnedEvent.TenantID` | `tenantId` | Go field: InstanceReturnedEvent.TenantID<br/>type: string | `approval/events_instance.go:98` |
| `InstanceReturnedEvent.ToNodeID` | `toNodeId` | Go field: InstanceReturnedEvent.ToNodeID<br/>type: string | `approval/events_instance.go:100` |
| `InstanceRolledBackEvent.FromNodeID` | `fromNodeId` | Go field: InstanceRolledBackEvent.FromNodeID<br/>type: string | `approval/events_instance.go:76` |
| `InstanceRolledBackEvent.InstanceID` | `instanceId` | Go field: InstanceRolledBackEvent.InstanceID<br/>type: string | `approval/events_instance.go:74` |
| `InstanceRolledBackEvent.OccurredTime` | `occurredTime` | Go field: InstanceRolledBackEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_instance.go:79` |
| `InstanceRolledBackEvent.OperatorID` | `operatorId` | Go field: InstanceRolledBackEvent.OperatorID<br/>type: string | `approval/events_instance.go:78` |
| `InstanceRolledBackEvent.TenantID` | `tenantId` | Go field: InstanceRolledBackEvent.TenantID<br/>type: string | `approval/events_instance.go:75` |
| `InstanceRolledBackEvent.ToNodeID` | `toNodeId` | Go field: InstanceRolledBackEvent.ToNodeID<br/>type: string | `approval/events_instance.go:77` |
| `InstanceWithdrawnEvent.InstanceID` | `instanceId` | Go field: InstanceWithdrawnEvent.InstanceID<br/>type: string | `approval/events_instance.go:55` |
| `InstanceWithdrawnEvent.OccurredTime` | `occurredTime` | Go field: InstanceWithdrawnEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_instance.go:58` |
| `InstanceWithdrawnEvent.OperatorID` | `operatorId` | Go field: InstanceWithdrawnEvent.OperatorID<br/>type: string | `approval/events_instance.go:57` |
| `InstanceWithdrawnEvent.TenantID` | `tenantId` | Go field: InstanceWithdrawnEvent.TenantID<br/>type: string | `approval/events_instance.go:56` |
| `InterfaceInfo.Addrs` | `addrs` | Go field: InterfaceInfo.Addrs<br/>type: []string | `monitor/service.go:207` |
| `InterfaceInfo.Flags` | `flags` | Go field: InterfaceInfo.Flags<br/>type: []string | `monitor/service.go:206` |
| `InterfaceInfo.HardwareAddr` | `hardwareAddr` | Go field: InterfaceInfo.HardwareAddr<br/>type: string | `monitor/service.go:205` |
| `InterfaceInfo.Index` | `index` | Go field: InterfaceInfo.Index<br/>type: int | `monitor/service.go:202` |
| `InterfaceInfo.MTU` | `mtu` | Go field: InterfaceInfo.MTU<br/>type: int | `monitor/service.go:203` |
| `InterfaceInfo.Name` | `name` | Go field: InterfaceInfo.Name<br/>type: string | `monitor/service.go:204` |
| `ListPartsParams.ClaimID` | `claimId` | Go field: ListPartsParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:515` |
| `ListPartsResult.Parts` | `parts` | Go field: ListPartsResult.Parts<br/>type: []ListedPart | `internal/storage/resource.go:529` |
| `ListedPart.PartNumber` | `partNumber` | Go field: ListedPart.PartNumber<br/>type: int | `internal/storage/resource.go:522` |
| `ListedPart.Size` | `size` | Go field: ListedPart.Size<br/>type: int64 | `internal/storage/resource.go:523` |
| `LoadInfo.Load1` | `load1` | Go field: LoadInfo.Load1<br/>type: float64 | `monitor/service.go:255` |
| `LoadInfo.Load15` | `load15` | Go field: LoadInfo.Load15<br/>type: float64 | `monitor/service.go:257` |
| `LoadInfo.Load5` | `load5` | Go field: LoadInfo.Load5<br/>type: float64 | `monitor/service.go:256` |
| `LoginChallenge.Data` | `data` | Go field: LoginChallenge.Data<br/>type: any | `security/challenge.go:8` |
| `LoginChallenge.Required` | `required` | Go field: LoginChallenge.Required<br/>type: bool | `security/challenge.go:9` |
| `LoginChallenge.Type` | `type` | Go field: LoginChallenge.Type<br/>type: string | `security/challenge.go:7` |
| `LoginEvent.AuthType` | `authType` | Go field: LoginEvent.AuthType<br/>type: string | `security/login_event.go:13` |
| `LoginEvent.ErrorCode` | `errorCode` | Go field: LoginEvent.ErrorCode<br/>type: int | `security/login_event.go:21` |
| `LoginEvent.FailReason` | `failReason` | Go field: LoginEvent.FailReason<br/>type: string | `security/login_event.go:20` |
| `LoginEvent.IsOk` | `isOk` | Go field: LoginEvent.IsOk<br/>type: bool | `security/login_event.go:19` |
| `LoginEvent.LoginIP` | `loginIp` | Go field: LoginEvent.LoginIP<br/>type: string | `security/login_event.go:16` |
| `LoginEvent.TraceID` | `traceId` | Go field: LoginEvent.TraceID<br/>type: string | `security/login_event.go:18` |
| `LoginEvent.UserAgent` | `userAgent` | Go field: LoginEvent.UserAgent<br/>type: string | `security/login_event.go:17` |
| `LoginEvent.UserID` | `userId` | Go field: LoginEvent.UserID<br/>type: *string | `security/login_event.go:14` |
| `LoginEvent.Username` | `username` | Go field: LoginEvent.Username<br/>type: string | `security/login_event.go:15` |
| `LoginParams.Credentials` | `credentials` | Go field: LoginParams.Credentials<br/>type: any<br/>validate: "required" | `internal/security/auth_resource.go:96` |
| `LoginParams.Principal` | `principal` | Go field: LoginParams.Principal<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:95` |
| `LoginParams.Type` | `type` | Go field: LoginParams.Type<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:94` |
| `LoginResult.Challenge` | `challenge` | Go field: LoginResult.Challenge<br/>type: *LoginChallenge | `security/challenge.go:18` |
| `LoginResult.ChallengeToken` | `challengeToken` | Go field: LoginResult.ChallengeToken<br/>type: string | `security/challenge.go:17` |
| `LoginResult.Tokens` | `tokens` | Go field: LoginResult.Tokens<br/>type: *AuthTokens | `security/challenge.go:16` |
| `MarkCCReadParams.InstanceID` | `instanceId` | Go field: MarkCCReadParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:308` |
| `MemoryInfo.Swap` | `swap` | Go field: MemoryInfo.Swap<br/>type: *SwapMemory | `monitor/service.go:78` |
| `MemoryInfo.Virtual` | `virtual` | Go field: MemoryInfo.Virtual<br/>type: *VirtualMemory | `monitor/service.go:77` |
| `MemorySummary.Total` | `total` | Go field: MemorySummary.Total<br/>type: uint64 | `monitor/service.go:70` |
| `MemorySummary.Used` | `used` | Go field: MemorySummary.Used<br/>type: uint64 | `monitor/service.go:71` |
| `MemorySummary.UsedPercent` | `usedPercent` | Go field: MemorySummary.UsedPercent<br/>type: float64 | `monitor/service.go:72` |
| `Metrics.AvgCompletionSeconds` | `avgCompletionSeconds` | Go field: Metrics.AvgCompletionSeconds<br/>type: float64 | `approval/admin/metrics.go:25` |
| `Metrics.CapturedAt` | `capturedAt` | Go field: Metrics.CapturedAt<br/>type: timex.DateTime | `approval/admin/metrics.go:14` |
| `Metrics.InstanceCounts` | `instanceCounts` | Go field: Metrics.InstanceCounts<br/>type: map[string]int | `approval/admin/metrics.go:17` |
| `Metrics.PendingBindingFailures` | `pendingBindingFailures` | Go field: Metrics.PendingBindingFailures<br/>type: int | `approval/admin/metrics.go:28` |
| `Metrics.TaskCounts` | `taskCounts` | Go field: Metrics.TaskCounts<br/>type: map[string]int | `approval/admin/metrics.go:19` |
| `Metrics.TenantID` | `tenantId` | Go field: Metrics.TenantID<br/>type: string | `approval/admin/metrics.go:12` |
| `Metrics.TimeoutTaskCount` | `timeoutTaskCount` | Go field: Metrics.TimeoutTaskCount<br/>type: int | `approval/admin/metrics.go:21` |
| `Model.ID` | `id` | Go field: Model.ID<br/>type: string | `internal/orm/model.go:7` |
| `NetIOCounter.BytesRecv` | `bytesRecv` | Go field: NetIOCounter.BytesRecv<br/>type: uint64 | `monitor/service.go:214` |
| `NetIOCounter.BytesSent` | `bytesSent` | Go field: NetIOCounter.BytesSent<br/>type: uint64 | `monitor/service.go:213` |
| `NetIOCounter.DroppedIn` | `droppedIn` | Go field: NetIOCounter.DroppedIn<br/>type: uint64 | `monitor/service.go:219` |
| `NetIOCounter.DroppedOut` | `droppedOut` | Go field: NetIOCounter.DroppedOut<br/>type: uint64 | `monitor/service.go:220` |
| `NetIOCounter.ErrorsIn` | `errorsIn` | Go field: NetIOCounter.ErrorsIn<br/>type: uint64 | `monitor/service.go:217` |
| `NetIOCounter.ErrorsOut` | `errorsOut` | Go field: NetIOCounter.ErrorsOut<br/>type: uint64 | `monitor/service.go:218` |
| `NetIOCounter.FIFOIn` | `fifoIn` | Go field: NetIOCounter.FIFOIn<br/>type: uint64 | `monitor/service.go:221` |
| `NetIOCounter.FIFOOut` | `fifoOut` | Go field: NetIOCounter.FIFOOut<br/>type: uint64 | `monitor/service.go:222` |
| `NetIOCounter.Name` | `name` | Go field: NetIOCounter.Name<br/>type: string | `monitor/service.go:212` |
| `NetIOCounter.PacketsRecv` | `packetsRecv` | Go field: NetIOCounter.PacketsRecv<br/>type: uint64 | `monitor/service.go:216` |
| `NetIOCounter.PacketsSent` | `packetsSent` | Go field: NetIOCounter.PacketsSent<br/>type: uint64 | `monitor/service.go:215` |
| `NetworkInfo.IOCounters` | `ioCounters` | Go field: NetworkInfo.IOCounters<br/>type: map[string]*NetIOCounter | `monitor/service.go:197` |
| `NetworkInfo.Interfaces` | `interfaces` | Go field: NetworkInfo.Interfaces<br/>type: []*InterfaceInfo | `monitor/service.go:196` |
| `NetworkSummary.BytesRecv` | `bytesRecv` | Go field: NetworkSummary.BytesRecv<br/>type: uint64 | `monitor/service.go:189` |
| `NetworkSummary.BytesSent` | `bytesSent` | Go field: NetworkSummary.BytesSent<br/>type: uint64 | `monitor/service.go:188` |
| `NetworkSummary.Interfaces` | `interfaces` | Go field: NetworkSummary.Interfaces<br/>type: int | `monitor/service.go:187` |
| `NetworkSummary.PacketsRecv` | `packetsRecv` | Go field: NetworkSummary.PacketsRecv<br/>type: uint64 | `monitor/service.go:191` |
| `NetworkSummary.PacketsSent` | `packetsSent` | Go field: NetworkSummary.PacketsSent<br/>type: uint64 | `monitor/service.go:190` |
| `NodeAutoPassedEvent.InstanceID` | `instanceId` | Go field: NodeAutoPassedEvent.InstanceID<br/>type: string | `approval/events_node.go:28` |
| `NodeAutoPassedEvent.NodeID` | `nodeId` | Go field: NodeAutoPassedEvent.NodeID<br/>type: string | `approval/events_node.go:30` |
| `NodeAutoPassedEvent.OccurredTime` | `occurredTime` | Go field: NodeAutoPassedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_node.go:32` |
| `NodeAutoPassedEvent.Reason` | `reason` | Go field: NodeAutoPassedEvent.Reason<br/>type: string | `approval/events_node.go:31` |
| `NodeAutoPassedEvent.TenantID` | `tenantId` | Go field: NodeAutoPassedEvent.TenantID<br/>type: string | `approval/events_node.go:29` |
| `NodeDefinition.Data` | `data` | Go field: NodeDefinition.Data<br/>type: json.RawMessage | `approval/flow_definition.go:32` |
| `NodeDefinition.ID` | `id` | Go field: NodeDefinition.ID<br/>type: string | `approval/flow_definition.go:29` |
| `NodeDefinition.Kind` | `kind` | Go field: NodeDefinition.Kind<br/>type: NodeKind | `approval/flow_definition.go:30` |
| `NodeDefinition.Position` | `position` | Go field: NodeDefinition.Position<br/>type: Position | `approval/flow_definition.go:31` |
| `NodeEnteredEvent.InstanceID` | `instanceId` | Go field: NodeEnteredEvent.InstanceID<br/>type: string | `approval/events_node.go:7` |
| `NodeEnteredEvent.NodeID` | `nodeId` | Go field: NodeEnteredEvent.NodeID<br/>type: string | `approval/events_node.go:9` |
| `NodeEnteredEvent.NodeName` | `nodeName` | Go field: NodeEnteredEvent.NodeName<br/>type: string | `approval/events_node.go:10` |
| `NodeEnteredEvent.OccurredTime` | `occurredTime` | Go field: NodeEnteredEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_node.go:11` |
| `NodeEnteredEvent.TenantID` | `tenantId` | Go field: NodeEnteredEvent.TenantID<br/>type: string | `approval/events_node.go:8` |
| `OTPChallengeData.Destination` | `destination` | Go field: OTPChallengeData.Destination<br/>type: string | `security/otp.go:45` |
| `OTPChallengeData.Meta` | `meta` | Go field: OTPChallengeData.Meta<br/>type: map[string]any | `security/otp.go:46` |
| `ObjectInfo.Bucket` | `bucket` | Go field: ObjectInfo.Bucket<br/>type: string | `storage/service.go:119` |
| `ObjectInfo.ContentType` | `contentType` | Go field: ObjectInfo.ContentType<br/>type: string | `storage/service.go:127` |
| `ObjectInfo.ETag` | `eTag` | Go field: ObjectInfo.ETag<br/>type: string | `storage/service.go:123` |
| `ObjectInfo.Key` | `key` | Go field: ObjectInfo.Key<br/>type: string | `storage/service.go:121` |
| `ObjectInfo.LastModified` | `lastModified` | Go field: ObjectInfo.LastModified<br/>type: time.Time | `storage/service.go:129` |
| `ObjectInfo.Metadata` | `metadata` | Go field: ObjectInfo.Metadata<br/>type: map[string]string | `storage/service.go:131` |
| `ObjectInfo.Size` | `size` | Go field: ObjectInfo.Size<br/>type: int64 | `storage/service.go:125` |
| `OperatorInfo.DepartmentID` | `departmentId` | Go field: OperatorInfo.DepartmentID<br/>type: *string | `approval/models.go:13` |
| `OperatorInfo.DepartmentName` | `departmentName` | Go field: OperatorInfo.DepartmentName<br/>type: *string | `approval/models.go:14` |
| `OperatorInfo.ID` | `id` | Go field: OperatorInfo.ID<br/>type: string | `approval/models.go:11` |
| `OperatorInfo.Name` | `name` | Go field: OperatorInfo.Name<br/>type: string | `approval/models.go:12` |
| `Page.Items` | `items` | Go field: Page.Items<br/>type: []T | `page/page.go:50` |
| `Page.Page` | `page` | Go field: Page.Page<br/>type: int | `page/page.go:47` |
| `Page.Size` | `size` | Go field: Page.Size<br/>type: int | `page/page.go:48` |
| `Page.Total` | `total` | Go field: Page.Total<br/>type: int64 | `page/page.go:49` |
| `Pageable.Page` | `page` | Go field: Pageable.Page<br/>type: int | `page/page.go:14` |
| `Pageable.Size` | `size` | Go field: Pageable.Size<br/>type: int | `page/page.go:15` |
| `PartitionInfo.Device` | `device` | Go field: PartitionInfo.Device<br/>type: string | `monitor/service.go:153` |
| `PartitionInfo.FSType` | `fsType` | Go field: PartitionInfo.FSType<br/>type: string | `monitor/service.go:155` |
| `PartitionInfo.Free` | `free` | Go field: PartitionInfo.Free<br/>type: uint64 | `monitor/service.go:158` |
| `PartitionInfo.INodesFree` | `iNodesFree` | Go field: PartitionInfo.INodesFree<br/>type: uint64 | `monitor/service.go:163` |
| `PartitionInfo.INodesTotal` | `iNodesTotal` | Go field: PartitionInfo.INodesTotal<br/>type: uint64 | `monitor/service.go:161` |
| `PartitionInfo.INodesUsed` | `iNodesUsed` | Go field: PartitionInfo.INodesUsed<br/>type: uint64 | `monitor/service.go:162` |
| `PartitionInfo.INodesUsedPercent` | `iNodesUsedPercent` | Go field: PartitionInfo.INodesUsedPercent<br/>type: float64 | `monitor/service.go:164` |
| `PartitionInfo.MountPoint` | `mountPoint` | Go field: PartitionInfo.MountPoint<br/>type: string | `monitor/service.go:154` |
| `PartitionInfo.Options` | `options` | Go field: PartitionInfo.Options<br/>type: []string | `monitor/service.go:156` |
| `PartitionInfo.Total` | `total` | Go field: PartitionInfo.Total<br/>type: uint64 | `monitor/service.go:157` |
| `PartitionInfo.Used` | `used` | Go field: PartitionInfo.Used<br/>type: uint64 | `monitor/service.go:159` |
| `PartitionInfo.UsedPercent` | `usedPercent` | Go field: PartitionInfo.UsedPercent<br/>type: float64 | `monitor/service.go:160` |
| `PasswordChangeChallengeData.Meta` | `meta` | Go field: PasswordChangeChallengeData.Meta<br/>type: map[string]any | `security/password_change.go:19` |
| `PasswordChangeChallengeData.Reason` | `reason` | Go field: PasswordChangeChallengeData.Reason<br/>type: string | `security/password_change.go:18` |
| `PendingCounts.PendingTaskCount` | `pendingTaskCount` | Go field: PendingCounts.PendingTaskCount<br/>type: int | `approval/my/types.go:74` |
| `PendingCounts.UnreadCCCount` | `unreadCcCount` | Go field: PendingCounts.UnreadCCCount<br/>type: int | `approval/my/types.go:75` |
| `PendingDelete.Attempts` | `attempts` | Go field: PendingDelete.Attempts<br/>type: int | `internal/storage/store/delete.go:34` |
| `PendingDelete.CreatedAt` | `createdAt` | Go field: PendingDelete.CreatedAt<br/>type: timex.DateTime | `internal/storage/store/delete.go:36` |
| `PendingDelete.ID` | `id` | Go field: PendingDelete.ID<br/>type: string | `internal/storage/store/delete.go:30` |
| `PendingDelete.Key` | `key` | Go field: PendingDelete.Key<br/>type: string | `internal/storage/store/delete.go:31` |
| `PendingDelete.NextAttemptAt` | `nextAttemptAt` | Go field: PendingDelete.NextAttemptAt<br/>type: timex.DateTime | `internal/storage/store/delete.go:35` |
| `PendingDelete.Reason` | `reason` | Go field: PendingDelete.Reason<br/>type: storage.DeleteReason | `internal/storage/store/delete.go:33` |
| `PendingDelete.UploadID` | `uploadId` | Go field: PendingDelete.UploadID<br/>type: string | `internal/storage/store/delete.go:32` |
| `PendingTask.ApplicantName` | `applicantName` | Go field: PendingTask.ApplicantName<br/>type: string | `approval/my/types.go:37` |
| `PendingTask.CreatedAt` | `createdAt` | Go field: PendingTask.CreatedAt<br/>type: timex.DateTime | `approval/my/types.go:39` |
| `PendingTask.Deadline` | `deadline` | Go field: PendingTask.Deadline<br/>type: *timex.DateTime | `approval/my/types.go:40` |
| `PendingTask.FlowIcon` | `flowIcon` | Go field: PendingTask.FlowIcon<br/>type: *string | `approval/my/types.go:36` |
| `PendingTask.FlowName` | `flowName` | Go field: PendingTask.FlowName<br/>type: string | `approval/my/types.go:35` |
| `PendingTask.InstanceID` | `instanceId` | Go field: PendingTask.InstanceID<br/>type: string | `approval/my/types.go:32` |
| `PendingTask.InstanceNo` | `instanceNo` | Go field: PendingTask.InstanceNo<br/>type: string | `approval/my/types.go:34` |
| `PendingTask.InstanceTitle` | `instanceTitle` | Go field: PendingTask.InstanceTitle<br/>type: string | `approval/my/types.go:33` |
| `PendingTask.IsTimeout` | `isTimeout` | Go field: PendingTask.IsTimeout<br/>type: bool | `approval/my/types.go:41` |
| `PendingTask.NodeName` | `nodeName` | Go field: PendingTask.NodeName<br/>type: string | `approval/my/types.go:38` |
| `PendingTask.TaskID` | `taskId` | Go field: PendingTask.TaskID<br/>type: string | `approval/my/types.go:31` |
| `Position.X` | `x` | Go field: Position.X<br/>type: float64 | `approval/flow_definition.go:23` |
| `Position.Y` | `y` | Go field: Position.Y<br/>type: float64 | `approval/flow_definition.go:24` |
| `PrimaryKey.Columns` | `columns` | Go field: PrimaryKey.Columns<br/>type: []string | `schema/service.go:39` |
| `PrimaryKey.Name` | `name` | Go field: PrimaryKey.Name<br/>type: string | `schema/service.go:38` |
| `Principal.Details` | `details` | Go field: Principal.Details<br/>type: any | `security/principal.go:82` |
| `Principal.ID` | `id` | Go field: Principal.ID<br/>type: string | `security/principal.go:76` |
| `Principal.Name` | `name` | Go field: Principal.Name<br/>type: string | `security/principal.go:78` |
| `Principal.Roles` | `roles` | Go field: Principal.Roles<br/>type: []string | `security/principal.go:80` |
| `Principal.Type` | `type` | Go field: Principal.Type<br/>type: PrincipalType | `security/principal.go:74` |
| `ProcessInfo.CPUPercent` | `cpuPercent` | Go field: ProcessInfo.CPUPercent<br/>type: float64 | `monitor/service.go:246` |
| `ProcessInfo.CWD` | `cwd` | Go field: ProcessInfo.CWD<br/>type: string | `monitor/service.go:240` |
| `ProcessInfo.CommandLine` | `commandLine` | Go field: ProcessInfo.CommandLine<br/>type: string | `monitor/service.go:239` |
| `ProcessInfo.CreateTime` | `createTime` | Go field: ProcessInfo.CreateTime<br/>type: int64 | `monitor/service.go:243` |
| `ProcessInfo.Exe` | `exe` | Go field: ProcessInfo.Exe<br/>type: string | `monitor/service.go:238` |
| `ProcessInfo.MemoryPercent` | `memoryPercent` | Go field: ProcessInfo.MemoryPercent<br/>type: float32 | `monitor/service.go:247` |
| `ProcessInfo.MemoryRSS` | `memoryRss` | Go field: ProcessInfo.MemoryRSS<br/>type: uint64 | `monitor/service.go:248` |
| `ProcessInfo.MemorySwap` | `memorySwap` | Go field: ProcessInfo.MemorySwap<br/>type: uint64 | `monitor/service.go:250` |
| `ProcessInfo.MemoryVMS` | `memoryVms` | Go field: ProcessInfo.MemoryVMS<br/>type: uint64 | `monitor/service.go:249` |
| `ProcessInfo.Name` | `name` | Go field: ProcessInfo.Name<br/>type: string | `monitor/service.go:237` |
| `ProcessInfo.NumFDs` | `numFds` | Go field: ProcessInfo.NumFDs<br/>type: int32 | `monitor/service.go:245` |
| `ProcessInfo.NumThreads` | `numThreads` | Go field: ProcessInfo.NumThreads<br/>type: int32 | `monitor/service.go:244` |
| `ProcessInfo.PID` | `pid` | Go field: ProcessInfo.PID<br/>type: int32 | `monitor/service.go:235` |
| `ProcessInfo.ParentPID` | `parentPid` | Go field: ProcessInfo.ParentPID<br/>type: int32 | `monitor/service.go:236` |
| `ProcessInfo.Status` | `status` | Go field: ProcessInfo.Status<br/>type: string | `monitor/service.go:241` |
| `ProcessInfo.Username` | `username` | Go field: ProcessInfo.Username<br/>type: string | `monitor/service.go:242` |
| `ProcessSummary.CPUPercent` | `cpuPercent` | Go field: ProcessSummary.CPUPercent<br/>type: float64 | `monitor/service.go:229` |
| `ProcessSummary.MemoryPercent` | `memoryPercent` | Go field: ProcessSummary.MemoryPercent<br/>type: float32 | `monitor/service.go:230` |
| `ProcessSummary.Name` | `name` | Go field: ProcessSummary.Name<br/>type: string | `monitor/service.go:228` |
| `ProcessSummary.PID` | `pid` | Go field: ProcessSummary.PID<br/>type: int32 | `monitor/service.go:227` |
| `ProcessTaskParams.Action` | `action` | Go field: ProcessTaskParams.Action<br/>type: string<br/>validate: "required,oneof=approve reject transfer rollback handle" | `internal/approval/resource/instance.go:162` |
| `ProcessTaskParams.FormData` | `formData` | Go field: ProcessTaskParams.FormData<br/>type: map[string]any | `internal/approval/resource/instance.go:164` |
| `ProcessTaskParams.Opinion` | `opinion` | Go field: ProcessTaskParams.Opinion<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/instance.go:163` |
| `ProcessTaskParams.TargetNodeID` | `targetNodeId` | Go field: ProcessTaskParams.TargetNodeID<br/>type: string | `internal/approval/resource/instance.go:166` |
| `ProcessTaskParams.TaskID` | `taskId` | Go field: ProcessTaskParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:161` |
| `ProcessTaskParams.TransferToID` | `transferToId` | Go field: ProcessTaskParams.TransferToID<br/>type: string | `internal/approval/resource/instance.go:165` |
| `PublishVersionParams.VersionID` | `versionId` | Go field: PublishVersionParams.VersionID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:158` |
| `QueryArgs.Params` | `params` | Go field: QueryArgs.Params<br/>jsonschema: "description=Parameters for the SQL query placeholders"<br/>type: []any | `internal/mcp/tools/query.go:16` |
| `QueryArgs.SQL` | `sql` | Go field: QueryArgs.SQL<br/>jsonschema: "required,description=The SQL query with placeholders (?) for parameters"<br/>type: string | `internal/mcp/tools/query.go:15` |
| `Record.CompletedAt` | `completedAt` | Go field: Record.CompletedAt<br/>type: *timex.DateTime | `event/inbox/inbox.go:59` |
| `Record.ConsumerGroup` | `consumerGroup` | Go field: Record.ConsumerGroup<br/>type: string | `event/inbox/inbox.go:50` |
| `Record.CorrelationID` | `correlationId` | Go field: Record.CorrelationID<br/>type: string | `event/transport/outbox/outbox.go:48` |
| `Record.EventID` | `eventId` | Go field: Record.EventID<br/>type: string | `event/inbox/inbox.go:47` |
| `Record.EventID` | `eventId` | Go field: Record.EventID<br/>type: string | `event/transport/outbox/outbox.go:43` |
| `Record.EventType` | `eventType` | Go field: Record.EventType<br/>type: string | `event/transport/outbox/outbox.go:44` |
| `Record.Headers` | `headers` | Go field: Record.Headers<br/>type: map[string]string | `event/transport/outbox/outbox.go:49` |
| `Record.LastError` | `lastError` | Go field: Record.LastError<br/>type: *string | `event/transport/outbox/outbox.go:53` |
| `Record.LockID` | `lockId` | Go field: Record.LockID<br/>type: string | `event/inbox/inbox.go:54` |
| `Record.LockedUntil` | `lockedUntil` | Go field: Record.LockedUntil<br/>type: *timex.DateTime | `event/inbox/inbox.go:57` |
| `Record.OccurredAt` | `occurredAt` | Go field: Record.OccurredAt<br/>type: timex.DateTime | `event/transport/outbox/outbox.go:56` |
| `Record.Payload` | `payload` | Go field: Record.Payload<br/>type: json.RawMessage | `event/transport/outbox/outbox.go:50` |
| `Record.ProcessedAt` | `processedAt` | Go field: Record.ProcessedAt<br/>type: *timex.DateTime | `event/transport/outbox/outbox.go:54` |
| `Record.RetryAfter` | `retryAfter` | Go field: Record.RetryAfter<br/>type: *timex.DateTime | `event/transport/outbox/outbox.go:55` |
| `Record.RetryCount` | `retryCount` | Go field: Record.RetryCount<br/>type: int | `event/transport/outbox/outbox.go:52` |
| `Record.Source` | `source` | Go field: Record.Source<br/>type: string | `event/transport/outbox/outbox.go:45` |
| `Record.SpanID` | `spanId` | Go field: Record.SpanID<br/>type: string | `event/transport/outbox/outbox.go:47` |
| `Record.Status` | `status` | Go field: Record.Status<br/>type: Status | `event/transport/outbox/outbox.go:51` |
| `Record.Status` | `status` | Go field: Record.Status<br/>type: Status | `event/inbox/inbox.go:52` |
| `Record.TraceID` | `traceId` | Go field: Record.TraceID<br/>type: string | `event/transport/outbox/outbox.go:46` |
| `RefreshParams.RefreshToken` | `refreshToken` | Go field: RefreshParams.RefreshToken<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:182` |
| `RemoveAssigneeParams.TaskID` | `taskId` | Go field: RemoveAssigneeParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:362` |
| `Request.Meta` | `meta` | Go field: Request.Meta<br/>type: Meta | `api/request.go:59` |
| `Request.Params` | `params` | Go field: Request.Params<br/>type: Params | `api/request.go:58` |
| `ResolveChallengeParams.ChallengeToken` | `challengeToken` | Go field: ResolveChallengeParams.ChallengeToken<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:214` |
| `ResolveChallengeParams.Response` | `response` | Go field: ResolveChallengeParams.Response<br/>type: any<br/>validate: "required" | `internal/security/auth_resource.go:216` |
| `ResolveChallengeParams.Type` | `type` | Go field: ResolveChallengeParams.Type<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:215` |
| `ResubmitParams.FormData` | `formData` | Go field: ResubmitParams.FormData<br/>type: map[string]any | `internal/approval/resource/instance.go:255` |
| `ResubmitParams.InstanceID` | `instanceId` | Go field: ResubmitParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:254` |
| `Result.Code` | `code` | Go field: Result.Code<br/>type: int | `result/result.go:11` |
| `Result.Data` | `data` | Go field: Result.Data<br/>type: any | `result/result.go:13` |
| `Result.Message` | `message` | Go field: Result.Message<br/>type: string | `result/result.go:12` |
| `RolePermissionsChangedEvent.Roles` | `roles` | Go field: RolePermissionsChangedEvent.Roles<br/>type: []string | `security/cached_role_permission_loader.go:18` |
| `Sortable.Sort` | `sort` | Go field: Sortable.Sort<br/>type: []sortx.OrderSpec | `crud/params.go:35` |
| `StartInstanceParams.BusinessRecordID` | `businessRecordId` | Go field: StartInstanceParams.BusinessRecordID<br/>type: *string | `internal/approval/resource/instance.go:131` |
| `StartInstanceParams.FlowCode` | `flowCode` | Go field: StartInstanceParams.FlowCode<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:130` |
| `StartInstanceParams.FormData` | `formData` | Go field: StartInstanceParams.FormData<br/>type: map[string]any | `internal/approval/resource/instance.go:132` |
| `StartInstanceParams.TenantID` | `tenantId` | Go field: StartInstanceParams.TenantID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:129` |
| `SwapMemory.Free` | `free` | Go field: SwapMemory.Free<br/>type: uint64 | `monitor/service.go:127` |
| `SwapMemory.PageFault` | `pageFault` | Go field: SwapMemory.PageFault<br/>type: uint64 | `monitor/service.go:133` |
| `SwapMemory.PageIn` | `pageIn` | Go field: SwapMemory.PageIn<br/>type: uint64 | `monitor/service.go:131` |
| `SwapMemory.PageMajorFault` | `pageMajorFault` | Go field: SwapMemory.PageMajorFault<br/>type: uint64 | `monitor/service.go:134` |
| `SwapMemory.PageOut` | `pageOut` | Go field: SwapMemory.PageOut<br/>type: uint64 | `monitor/service.go:132` |
| `SwapMemory.SwapIn` | `swapIn` | Go field: SwapMemory.SwapIn<br/>type: uint64 | `monitor/service.go:129` |
| `SwapMemory.SwapOut` | `swapOut` | Go field: SwapMemory.SwapOut<br/>type: uint64 | `monitor/service.go:130` |
| `SwapMemory.Total` | `total` | Go field: SwapMemory.Total<br/>type: uint64 | `monitor/service.go:125` |
| `SwapMemory.Used` | `used` | Go field: SwapMemory.Used<br/>type: uint64 | `monitor/service.go:126` |
| `SwapMemory.UsedPercent` | `usedPercent` | Go field: SwapMemory.UsedPercent<br/>type: float64 | `monitor/service.go:128` |
| `SystemOverview.Build` | `build` | Go field: SystemOverview.Build<br/>type: *BuildInfo | `monitor/service.go:14` |
| `SystemOverview.CPU` | `cpu` | Go field: SystemOverview.CPU<br/>type: *CPUSummary | `monitor/service.go:8` |
| `SystemOverview.Disk` | `disk` | Go field: SystemOverview.Disk<br/>type: *DiskSummary | `monitor/service.go:10` |
| `SystemOverview.Host` | `host` | Go field: SystemOverview.Host<br/>type: *HostSummary | `monitor/service.go:7` |
| `SystemOverview.Load` | `load` | Go field: SystemOverview.Load<br/>type: *LoadInfo | `monitor/service.go:13` |
| `SystemOverview.Memory` | `memory` | Go field: SystemOverview.Memory<br/>type: *MemorySummary | `monitor/service.go:9` |
| `SystemOverview.Network` | `network` | Go field: SystemOverview.Network<br/>type: *NetworkSummary | `monitor/service.go:11` |
| `SystemOverview.Process` | `process` | Go field: SystemOverview.Process<br/>type: *ProcessSummary | `monitor/service.go:12` |
| `Table.Comment` | `comment` | Go field: Table.Comment<br/>type: string | `schema/service.go:9` |
| `Table.Name` | `name` | Go field: Table.Name<br/>type: string | `schema/service.go:7` |
| `Table.Schema` | `schema` | Go field: Table.Schema<br/>type: string | `schema/service.go:8` |
| `TableSchema.Checks` | `checks` | Go field: TableSchema.Checks<br/>type: []Check | `schema/service.go:22` |
| `TableSchema.Columns` | `columns` | Go field: TableSchema.Columns<br/>type: []Column | `schema/service.go:17` |
| `TableSchema.Comment` | `comment` | Go field: TableSchema.Comment<br/>type: string | `schema/service.go:16` |
| `TableSchema.ForeignKeys` | `foreignKeys` | Go field: TableSchema.ForeignKeys<br/>type: []ForeignKey | `schema/service.go:21` |
| `TableSchema.Indexes` | `indexes` | Go field: TableSchema.Indexes<br/>type: []Index | `schema/service.go:19` |
| `TableSchema.Name` | `name` | Go field: TableSchema.Name<br/>type: string | `schema/service.go:14` |
| `TableSchema.PrimaryKey` | `primaryKey` | Go field: TableSchema.PrimaryKey<br/>type: *PrimaryKey | `schema/service.go:18` |
| `TableSchema.Schema` | `schema` | Go field: TableSchema.Schema<br/>type: string | `schema/service.go:15` |
| `TableSchema.UniqueKeys` | `uniqueKeys` | Go field: TableSchema.UniqueKeys<br/>type: []UniqueKey | `schema/service.go:20` |
| `Task.AddAssigneeType` | `addAssigneeType` | Go field: Task.AddAssigneeType<br/>type: *AddAssigneeType | `approval/models.go:211` |
| `Task.AssigneeID` | `assigneeId` | Go field: Task.AssigneeID<br/>type: string | `approval/admin/types.go:28` |
| `Task.AssigneeID` | `assigneeId` | Go field: Task.AssigneeID<br/>type: string | `approval/models.go:203` |
| `Task.AssigneeName` | `assigneeName` | Go field: Task.AssigneeName<br/>type: string | `approval/admin/types.go:29` |
| `Task.AssigneeName` | `assigneeName` | Go field: Task.AssigneeName<br/>type: string | `approval/models.go:204` |
| `Task.CreatedAt` | `createdAt` | Go field: Task.CreatedAt<br/>type: timex.DateTime | `approval/admin/types.go:31` |
| `Task.Deadline` | `deadline` | Go field: Task.Deadline<br/>type: *timex.DateTime | `approval/models.go:212` |
| `Task.Deadline` | `deadline` | Go field: Task.Deadline<br/>type: *timex.DateTime | `approval/admin/types.go:32` |
| `Task.DelegatorID` | `delegatorId` | Go field: Task.DelegatorID<br/>type: *string | `approval/models.go:205` |
| `Task.DelegatorName` | `delegatorName` | Go field: Task.DelegatorName<br/>type: *string | `approval/models.go:206` |
| `Task.FinishedAt` | `finishedAt` | Go field: Task.FinishedAt<br/>type: *timex.DateTime | `approval/models.go:215` |
| `Task.FinishedAt` | `finishedAt` | Go field: Task.FinishedAt<br/>type: *timex.DateTime | `approval/admin/types.go:33` |
| `Task.FlowName` | `flowName` | Go field: Task.FlowName<br/>type: string | `approval/admin/types.go:26` |
| `Task.InstanceID` | `instanceId` | Go field: Task.InstanceID<br/>type: string | `approval/models.go:201` |
| `Task.InstanceID` | `instanceId` | Go field: Task.InstanceID<br/>type: string | `approval/admin/types.go:24` |
| `Task.InstanceTitle` | `instanceTitle` | Go field: Task.InstanceTitle<br/>type: string | `approval/admin/types.go:25` |
| `Task.IsPreWarningSent` | `isPreWarningSent` | Go field: Task.IsPreWarningSent<br/>type: bool | `approval/models.go:214` |
| `Task.IsTimeout` | `isTimeout` | Go field: Task.IsTimeout<br/>type: bool | `approval/models.go:213` |
| `Task.NodeID` | `nodeId` | Go field: Task.NodeID<br/>type: string | `approval/models.go:202` |
| `Task.NodeName` | `nodeName` | Go field: Task.NodeName<br/>type: string | `approval/admin/types.go:27` |
| `Task.ParentTaskID` | `parentTaskId` | Go field: Task.ParentTaskID<br/>type: *string | `approval/models.go:210` |
| `Task.ReadAt` | `readAt` | Go field: Task.ReadAt<br/>type: *timex.DateTime | `approval/models.go:209` |
| `Task.SortOrder` | `sortOrder` | Go field: Task.SortOrder<br/>type: int | `approval/models.go:207` |
| `Task.Status` | `status` | Go field: Task.Status<br/>type: string | `approval/admin/types.go:30` |
| `Task.Status` | `status` | Go field: Task.Status<br/>type: TaskStatus | `approval/models.go:208` |
| `Task.TaskID` | `taskId` | Go field: Task.TaskID<br/>type: string | `approval/admin/types.go:23` |
| `Task.TenantID` | `tenantId` | Go field: Task.TenantID<br/>type: string | `approval/models.go:200` |
| `TaskApprovedEvent.InstanceID` | `instanceId` | Go field: TaskApprovedEvent.InstanceID<br/>type: string | `approval/events_task.go:44` |
| `TaskApprovedEvent.NodeID` | `nodeId` | Go field: TaskApprovedEvent.NodeID<br/>type: string | `approval/events_task.go:45` |
| `TaskApprovedEvent.OccurredTime` | `occurredTime` | Go field: TaskApprovedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:48` |
| `TaskApprovedEvent.OperatorID` | `operatorId` | Go field: TaskApprovedEvent.OperatorID<br/>type: string | `approval/events_task.go:46` |
| `TaskApprovedEvent.Opinion` | `opinion` | Go field: TaskApprovedEvent.Opinion<br/>type: *string | `approval/events_task.go:47` |
| `TaskApprovedEvent.TaskID` | `taskId` | Go field: TaskApprovedEvent.TaskID<br/>type: string | `approval/events_task.go:42` |
| `TaskApprovedEvent.TenantID` | `tenantId` | Go field: TaskApprovedEvent.TenantID<br/>type: string | `approval/events_task.go:43` |
| `TaskCreatedEvent.AssigneeID` | `assigneeId` | Go field: TaskCreatedEvent.AssigneeID<br/>type: string | `approval/events_task.go:19` |
| `TaskCreatedEvent.AssigneeName` | `assigneeName` | Go field: TaskCreatedEvent.AssigneeName<br/>type: string | `approval/events_task.go:20` |
| `TaskCreatedEvent.Deadline` | `deadline` | Go field: TaskCreatedEvent.Deadline<br/>type: *timex.DateTime | `approval/events_task.go:21` |
| `TaskCreatedEvent.InstanceID` | `instanceId` | Go field: TaskCreatedEvent.InstanceID<br/>type: string | `approval/events_task.go:17` |
| `TaskCreatedEvent.NodeID` | `nodeId` | Go field: TaskCreatedEvent.NodeID<br/>type: string | `approval/events_task.go:18` |
| `TaskCreatedEvent.OccurredTime` | `occurredTime` | Go field: TaskCreatedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:22` |
| `TaskCreatedEvent.TaskID` | `taskId` | Go field: TaskCreatedEvent.TaskID<br/>type: string | `approval/events_task.go:15` |
| `TaskCreatedEvent.TenantID` | `tenantId` | Go field: TaskCreatedEvent.TenantID<br/>type: string | `approval/events_task.go:16` |
| `TaskDeadlineWarningEvent.AssigneeID` | `assigneeId` | Go field: TaskDeadlineWarningEvent.AssigneeID<br/>type: string | `approval/events_timeout.go:11` |
| `TaskDeadlineWarningEvent.AssigneeName` | `assigneeName` | Go field: TaskDeadlineWarningEvent.AssigneeName<br/>type: string | `approval/events_timeout.go:12` |
| `TaskDeadlineWarningEvent.Deadline` | `deadline` | Go field: TaskDeadlineWarningEvent.Deadline<br/>type: timex.DateTime | `approval/events_timeout.go:13` |
| `TaskDeadlineWarningEvent.HoursLeft` | `hoursLeft` | Go field: TaskDeadlineWarningEvent.HoursLeft<br/>type: int | `approval/events_timeout.go:14` |
| `TaskDeadlineWarningEvent.InstanceID` | `instanceId` | Go field: TaskDeadlineWarningEvent.InstanceID<br/>type: string | `approval/events_timeout.go:9` |
| `TaskDeadlineWarningEvent.NodeID` | `nodeId` | Go field: TaskDeadlineWarningEvent.NodeID<br/>type: string | `approval/events_timeout.go:10` |
| `TaskDeadlineWarningEvent.OccurredTime` | `occurredTime` | Go field: TaskDeadlineWarningEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_timeout.go:15` |
| `TaskDeadlineWarningEvent.TaskID` | `taskId` | Go field: TaskDeadlineWarningEvent.TaskID<br/>type: string | `approval/events_timeout.go:7` |
| `TaskDeadlineWarningEvent.TenantID` | `tenantId` | Go field: TaskDeadlineWarningEvent.TenantID<br/>type: string | `approval/events_timeout.go:8` |
| `TaskDetailInfo.AssigneeID` | `assigneeId` | Go field: TaskDetailInfo.AssigneeID<br/>type: string | `approval/admin/types.go:68` |
| `TaskDetailInfo.AssigneeName` | `assigneeName` | Go field: TaskDetailInfo.AssigneeName<br/>type: string | `approval/admin/types.go:69` |
| `TaskDetailInfo.CreatedAt` | `createdAt` | Go field: TaskDetailInfo.CreatedAt<br/>type: timex.DateTime | `approval/admin/types.go:76` |
| `TaskDetailInfo.Deadline` | `deadline` | Go field: TaskDetailInfo.Deadline<br/>type: *timex.DateTime | `approval/admin/types.go:74` |
| `TaskDetailInfo.DelegatorID` | `delegatorId` | Go field: TaskDetailInfo.DelegatorID<br/>type: *string | `approval/admin/types.go:70` |
| `TaskDetailInfo.DelegatorName` | `delegatorName` | Go field: TaskDetailInfo.DelegatorName<br/>type: *string | `approval/admin/types.go:71` |
| `TaskDetailInfo.FinishedAt` | `finishedAt` | Go field: TaskDetailInfo.FinishedAt<br/>type: *timex.DateTime | `approval/admin/types.go:77` |
| `TaskDetailInfo.IsTimeout` | `isTimeout` | Go field: TaskDetailInfo.IsTimeout<br/>type: bool | `approval/admin/types.go:75` |
| `TaskDetailInfo.NodeID` | `nodeId` | Go field: TaskDetailInfo.NodeID<br/>type: string | `approval/admin/types.go:66` |
| `TaskDetailInfo.NodeName` | `nodeName` | Go field: TaskDetailInfo.NodeName<br/>type: string | `approval/admin/types.go:67` |
| `TaskDetailInfo.SortOrder` | `sortOrder` | Go field: TaskDetailInfo.SortOrder<br/>type: int | `approval/admin/types.go:73` |
| `TaskDetailInfo.Status` | `status` | Go field: TaskDetailInfo.Status<br/>type: string | `approval/admin/types.go:72` |
| `TaskDetailInfo.TaskID` | `taskId` | Go field: TaskDetailInfo.TaskID<br/>type: string | `approval/admin/types.go:65` |
| `TaskHandledEvent.InstanceID` | `instanceId` | Go field: TaskHandledEvent.InstanceID<br/>type: string | `approval/events_task.go:69` |
| `TaskHandledEvent.NodeID` | `nodeId` | Go field: TaskHandledEvent.NodeID<br/>type: string | `approval/events_task.go:70` |
| `TaskHandledEvent.OccurredTime` | `occurredTime` | Go field: TaskHandledEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:73` |
| `TaskHandledEvent.OperatorID` | `operatorId` | Go field: TaskHandledEvent.OperatorID<br/>type: string | `approval/events_task.go:71` |
| `TaskHandledEvent.Opinion` | `opinion` | Go field: TaskHandledEvent.Opinion<br/>type: *string | `approval/events_task.go:72` |
| `TaskHandledEvent.TaskID` | `taskId` | Go field: TaskHandledEvent.TaskID<br/>type: string | `approval/events_task.go:67` |
| `TaskHandledEvent.TenantID` | `tenantId` | Go field: TaskHandledEvent.TenantID<br/>type: string | `approval/events_task.go:68` |
| `TaskInfo.AssigneeID` | `assigneeId` | Go field: TaskInfo.AssigneeID<br/>type: string | `approval/my/types.go:108` |
| `TaskInfo.AssigneeName` | `assigneeName` | Go field: TaskInfo.AssigneeName<br/>type: string | `approval/my/types.go:109` |
| `TaskInfo.CreatedAt` | `createdAt` | Go field: TaskInfo.CreatedAt<br/>type: timex.DateTime | `approval/my/types.go:112` |
| `TaskInfo.FinishedAt` | `finishedAt` | Go field: TaskInfo.FinishedAt<br/>type: *timex.DateTime | `approval/my/types.go:113` |
| `TaskInfo.NodeName` | `nodeName` | Go field: TaskInfo.NodeName<br/>type: string | `approval/my/types.go:107` |
| `TaskInfo.SortOrder` | `sortOrder` | Go field: TaskInfo.SortOrder<br/>type: int | `approval/my/types.go:111` |
| `TaskInfo.Status` | `status` | Go field: TaskInfo.Status<br/>type: string | `approval/my/types.go:110` |
| `TaskInfo.TaskID` | `taskId` | Go field: TaskInfo.TaskID<br/>type: string | `approval/my/types.go:106` |
| `TaskNodeData.AdminUserIDs` | `adminUserIds` | Go field: TaskNodeData.AdminUserIDs<br/>type: []string | `approval/node_data.go:45` |
| `TaskNodeData.Assignees` | `assignees` | Go field: TaskNodeData.Assignees<br/>type: []AssigneeDefinition | `approval/node_data.go:41` |
| `TaskNodeData.CCs` | `ccs` | Go field: TaskNodeData.CCs<br/>type: []CCDefinition | `approval/node_data.go:52` |
| `TaskNodeData.EmptyAssigneeAction` | `emptyAssigneeAction` | Go field: TaskNodeData.EmptyAssigneeAction<br/>type: EmptyAssigneeAction | `approval/node_data.go:43` |
| `TaskNodeData.ExecutionType` | `executionType` | Go field: TaskNodeData.ExecutionType<br/>type: ExecutionType | `approval/node_data.go:42` |
| `TaskNodeData.FallbackUserIDs` | `fallbackUserIds` | Go field: TaskNodeData.FallbackUserIDs<br/>type: []string | `approval/node_data.go:44` |
| `TaskNodeData.FieldPermissions` | `fieldPermissions` | Go field: TaskNodeData.FieldPermissions<br/>type: map[string]Permission | `approval/node_data.go:53` |
| `TaskNodeData.IsOpinionRequired` | `isOpinionRequired` | Go field: TaskNodeData.IsOpinionRequired<br/>type: bool | `approval/node_data.go:47` |
| `TaskNodeData.IsTransferAllowed` | `isTransferAllowed` | Go field: TaskNodeData.IsTransferAllowed<br/>type: bool | `approval/node_data.go:46` |
| `TaskNodeData.TimeoutAction` | `timeoutAction` | Go field: TaskNodeData.TimeoutAction<br/>type: TimeoutAction | `approval/node_data.go:49` |
| `TaskNodeData.TimeoutHours` | `timeoutHours` | Go field: TaskNodeData.TimeoutHours<br/>type: int | `approval/node_data.go:48` |
| `TaskNodeData.TimeoutNotifyBeforeHours` | `timeoutNotifyBeforeHours` | Go field: TaskNodeData.TimeoutNotifyBeforeHours<br/>type: int | `approval/node_data.go:50` |
| `TaskNodeData.UrgeCooldownMinutes` | `urgeCooldownMinutes` | Go field: TaskNodeData.UrgeCooldownMinutes<br/>type: int | `approval/node_data.go:51` |
| `TaskReassignedEvent.FromUserID` | `fromUserId` | Go field: TaskReassignedEvent.FromUserID<br/>type: string | `approval/events_task.go:153` |
| `TaskReassignedEvent.FromUserName` | `fromUserName` | Go field: TaskReassignedEvent.FromUserName<br/>type: string | `approval/events_task.go:154` |
| `TaskReassignedEvent.InstanceID` | `instanceId` | Go field: TaskReassignedEvent.InstanceID<br/>type: string | `approval/events_task.go:151` |
| `TaskReassignedEvent.NodeID` | `nodeId` | Go field: TaskReassignedEvent.NodeID<br/>type: string | `approval/events_task.go:152` |
| `TaskReassignedEvent.OccurredTime` | `occurredTime` | Go field: TaskReassignedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:158` |
| `TaskReassignedEvent.Reason` | `reason` | Go field: TaskReassignedEvent.Reason<br/>type: *string | `approval/events_task.go:157` |
| `TaskReassignedEvent.TaskID` | `taskId` | Go field: TaskReassignedEvent.TaskID<br/>type: string | `approval/events_task.go:149` |
| `TaskReassignedEvent.TenantID` | `tenantId` | Go field: TaskReassignedEvent.TenantID<br/>type: string | `approval/events_task.go:150` |
| `TaskReassignedEvent.ToUserID` | `toUserId` | Go field: TaskReassignedEvent.ToUserID<br/>type: string | `approval/events_task.go:155` |
| `TaskReassignedEvent.ToUserName` | `toUserName` | Go field: TaskReassignedEvent.ToUserName<br/>type: string | `approval/events_task.go:156` |
| `TaskRejectedEvent.InstanceID` | `instanceId` | Go field: TaskRejectedEvent.InstanceID<br/>type: string | `approval/events_task.go:94` |
| `TaskRejectedEvent.NodeID` | `nodeId` | Go field: TaskRejectedEvent.NodeID<br/>type: string | `approval/events_task.go:95` |
| `TaskRejectedEvent.OccurredTime` | `occurredTime` | Go field: TaskRejectedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:98` |
| `TaskRejectedEvent.OperatorID` | `operatorId` | Go field: TaskRejectedEvent.OperatorID<br/>type: string | `approval/events_task.go:96` |
| `TaskRejectedEvent.Opinion` | `opinion` | Go field: TaskRejectedEvent.Opinion<br/>type: *string | `approval/events_task.go:97` |
| `TaskRejectedEvent.TaskID` | `taskId` | Go field: TaskRejectedEvent.TaskID<br/>type: string | `approval/events_task.go:92` |
| `TaskRejectedEvent.TenantID` | `tenantId` | Go field: TaskRejectedEvent.TenantID<br/>type: string | `approval/events_task.go:93` |
| `TaskTimedOutEvent.AssigneeID` | `assigneeId` | Go field: TaskTimedOutEvent.AssigneeID<br/>type: string | `approval/events_task.go:185` |
| `TaskTimedOutEvent.AssigneeName` | `assigneeName` | Go field: TaskTimedOutEvent.AssigneeName<br/>type: string | `approval/events_task.go:186` |
| `TaskTimedOutEvent.Deadline` | `deadline` | Go field: TaskTimedOutEvent.Deadline<br/>type: timex.DateTime | `approval/events_task.go:187` |
| `TaskTimedOutEvent.InstanceID` | `instanceId` | Go field: TaskTimedOutEvent.InstanceID<br/>type: string | `approval/events_task.go:183` |
| `TaskTimedOutEvent.NodeID` | `nodeId` | Go field: TaskTimedOutEvent.NodeID<br/>type: string | `approval/events_task.go:184` |
| `TaskTimedOutEvent.OccurredTime` | `occurredTime` | Go field: TaskTimedOutEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:188` |
| `TaskTimedOutEvent.TaskID` | `taskId` | Go field: TaskTimedOutEvent.TaskID<br/>type: string | `approval/events_task.go:181` |
| `TaskTimedOutEvent.TenantID` | `tenantId` | Go field: TaskTimedOutEvent.TenantID<br/>type: string | `approval/events_task.go:182` |
| `TaskTransferredEvent.FromUserID` | `fromUserId` | Go field: TaskTransferredEvent.FromUserID<br/>type: string | `approval/events_task.go:121` |
| `TaskTransferredEvent.FromUserName` | `fromUserName` | Go field: TaskTransferredEvent.FromUserName<br/>type: string | `approval/events_task.go:122` |
| `TaskTransferredEvent.InstanceID` | `instanceId` | Go field: TaskTransferredEvent.InstanceID<br/>type: string | `approval/events_task.go:119` |
| `TaskTransferredEvent.NodeID` | `nodeId` | Go field: TaskTransferredEvent.NodeID<br/>type: string | `approval/events_task.go:120` |
| `TaskTransferredEvent.OccurredTime` | `occurredTime` | Go field: TaskTransferredEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_task.go:126` |
| `TaskTransferredEvent.Reason` | `reason` | Go field: TaskTransferredEvent.Reason<br/>type: *string | `approval/events_task.go:125` |
| `TaskTransferredEvent.TaskID` | `taskId` | Go field: TaskTransferredEvent.TaskID<br/>type: string | `approval/events_task.go:117` |
| `TaskTransferredEvent.TenantID` | `tenantId` | Go field: TaskTransferredEvent.TenantID<br/>type: string | `approval/events_task.go:118` |
| `TaskTransferredEvent.ToUserID` | `toUserId` | Go field: TaskTransferredEvent.ToUserID<br/>type: string | `approval/events_task.go:123` |
| `TaskTransferredEvent.ToUserName` | `toUserName` | Go field: TaskTransferredEvent.ToUserName<br/>type: string | `approval/events_task.go:124` |
| `TaskUrgedEvent.InstanceID` | `instanceId` | Go field: TaskUrgedEvent.InstanceID<br/>type: string | `approval/events_timeout.go:36` |
| `TaskUrgedEvent.Message` | `message` | Go field: TaskUrgedEvent.Message<br/>type: *string | `approval/events_timeout.go:44` |
| `TaskUrgedEvent.NodeID` | `nodeId` | Go field: TaskUrgedEvent.NodeID<br/>type: string | `approval/events_timeout.go:38` |
| `TaskUrgedEvent.OccurredTime` | `occurredTime` | Go field: TaskUrgedEvent.OccurredTime<br/>type: timex.DateTime | `approval/events_timeout.go:45` |
| `TaskUrgedEvent.TargetUserID` | `targetUserId` | Go field: TaskUrgedEvent.TargetUserID<br/>type: string | `approval/events_timeout.go:42` |
| `TaskUrgedEvent.TargetUserName` | `targetUserName` | Go field: TaskUrgedEvent.TargetUserName<br/>type: string | `approval/events_timeout.go:43` |
| `TaskUrgedEvent.TaskID` | `taskId` | Go field: TaskUrgedEvent.TaskID<br/>type: string | `approval/events_timeout.go:39` |
| `TaskUrgedEvent.TenantID` | `tenantId` | Go field: TaskUrgedEvent.TenantID<br/>type: string | `approval/events_timeout.go:37` |
| `TaskUrgedEvent.UrgerID` | `urgerId` | Go field: TaskUrgedEvent.UrgerID<br/>type: string | `approval/events_timeout.go:40` |
| `TaskUrgedEvent.UrgerName` | `urgerName` | Go field: TaskUrgedEvent.UrgerName<br/>type: string | `approval/events_timeout.go:41` |
| `ToggleActiveParams.FlowID` | `flowId` | Go field: ToggleActiveParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:307` |
| `ToggleActiveParams.IsActive` | `isActive` | Go field: ToggleActiveParams.IsActive<br/>type: bool | `internal/approval/resource/flow.go:308` |
| `TreeDataOption.Children` | `children` | Go field: TreeDataOption.Children<br/>type: []TreeDataOption | `crud/option.go:57` |
| `UniqueKey.Columns` | `columns` | Go field: UniqueKey.Columns<br/>type: []string | `schema/service.go:51` |
| `UniqueKey.Name` | `name` | Go field: UniqueKey.Name<br/>type: string | `schema/service.go:50` |
| `UpdateFlowParams.AdminUserIDs` | `adminUserIds` | Go field: UpdateFlowParams.AdminUserIDs<br/>type: []string | `internal/approval/resource/flow.go:260` |
| `UpdateFlowParams.Description` | `description` | Go field: UpdateFlowParams.Description<br/>type: *string | `internal/approval/resource/flow.go:259` |
| `UpdateFlowParams.FlowID` | `flowId` | Go field: UpdateFlowParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:256` |
| `UpdateFlowParams.Icon` | `icon` | Go field: UpdateFlowParams.Icon<br/>type: *string | `internal/approval/resource/flow.go:258` |
| `UpdateFlowParams.Initiators` | `initiators` | Go field: UpdateFlowParams.Initiators<br/>type: []CreateInitiatorParams | `internal/approval/resource/flow.go:263` |
| `UpdateFlowParams.InstanceTitleTemplate` | `instanceTitleTemplate` | Go field: UpdateFlowParams.InstanceTitleTemplate<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:262` |
| `UpdateFlowParams.IsAllInitiationAllowed` | `isAllInitiationAllowed` | Go field: UpdateFlowParams.IsAllInitiationAllowed<br/>type: bool | `internal/approval/resource/flow.go:261` |
| `UpdateFlowParams.Name` | `name` | Go field: UpdateFlowParams.Name<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:257` |
| `UpdateManyParams.List` | `list` | Go field: UpdateManyParams.List<br/>type: []TParams<br/>validate: "required,min=1,dive" | `crud/params.go:19` |
| `UploadClaim.ContentType` | `contentType` | Go field: UploadClaim.ContentType<br/>type: string | `internal/storage/store/claim.go:48` |
| `UploadClaim.CreatedAt` | `createdAt` | Go field: UploadClaim.CreatedAt<br/>type: timex.DateTime | `internal/storage/store/claim.go:43` |
| `UploadClaim.CreatedBy` | `createdBy` | Go field: UploadClaim.CreatedBy<br/>type: string | `internal/storage/store/claim.go:44` |
| `UploadClaim.ExpiresAt` | `expiresAt` | Go field: UploadClaim.ExpiresAt<br/>type: timex.DateTime | `internal/storage/store/claim.go:54` |
| `UploadClaim.ID` | `id` | Go field: UploadClaim.ID<br/>type: string | `internal/storage/store/claim.go:42` |
| `UploadClaim.Key` | `key` | Go field: UploadClaim.Key<br/>type: string | `internal/storage/store/claim.go:45` |
| `UploadClaim.OriginalFilename` | `originalFilename` | Go field: UploadClaim.OriginalFilename<br/>type: string | `internal/storage/store/claim.go:49` |
| `UploadClaim.PartCount` | `partCount` | Go field: UploadClaim.PartCount<br/>type: int | `internal/storage/store/claim.go:53` |
| `UploadClaim.PartSize` | `partSize` | Go field: UploadClaim.PartSize<br/>type: int64 | `internal/storage/store/claim.go:52` |
| `UploadClaim.Public` | `public` | Go field: UploadClaim.Public<br/>type: bool | `internal/storage/store/claim.go:51` |
| `UploadClaim.Size` | `size` | Go field: UploadClaim.Size<br/>type: int64 | `internal/storage/store/claim.go:47` |
| `UploadClaim.Status` | `status` | Go field: UploadClaim.Status<br/>type: ClaimStatus | `internal/storage/store/claim.go:50` |
| `UploadClaim.UploadID` | `uploadId` | Go field: UploadClaim.UploadID<br/>type: string | `internal/storage/store/claim.go:46` |
| `UploadPart.ClaimID` | `claimId` | Go field: UploadPart.ClaimID<br/>type: string | `internal/storage/store/part.go:21` |
| `UploadPart.CreatedAt` | `createdAt` | Go field: UploadPart.CreatedAt<br/>type: timex.DateTime | `internal/storage/store/part.go:25` |
| `UploadPart.ETag` | `eTag` | Go field: UploadPart.ETag<br/>type: string | `internal/storage/store/part.go:23` |
| `UploadPart.ID` | `id` | Go field: UploadPart.ID<br/>type: string | `internal/storage/store/part.go:20` |
| `UploadPart.PartNumber` | `partNumber` | Go field: UploadPart.PartNumber<br/>type: int | `internal/storage/store/part.go:22` |
| `UploadPart.Size` | `size` | Go field: UploadPart.Size<br/>type: int64 | `internal/storage/store/part.go:24` |
| `UploadPartParams.ClaimID` | `claimId` | Go field: UploadPartParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:401` |
| `UploadPartParams.PartNumber` | `partNumber` | Go field: UploadPartParams.PartNumber<br/>type: int<br/>validate: "required,min=1" | `internal/storage/resource.go:402` |
| `UploadPartResult.PartNumber` | `partNumber` | Go field: UploadPartResult.PartNumber<br/>type: int | `internal/storage/resource.go:410` |
| `UploadPartResult.Size` | `size` | Go field: UploadPartResult.Size<br/>type: int64 | `internal/storage/resource.go:411` |
| `UrgeRecord.InstanceID` | `instanceId` | Go field: UrgeRecord.InstanceID<br/>type: string | `approval/models.go:295` |
| `UrgeRecord.Message` | `message` | Go field: UrgeRecord.Message<br/>type: string | `approval/models.go:302` |
| `UrgeRecord.NodeID` | `nodeId` | Go field: UrgeRecord.NodeID<br/>type: string | `approval/models.go:296` |
| `UrgeRecord.TargetUserID` | `targetUserId` | Go field: UrgeRecord.TargetUserID<br/>type: string | `approval/models.go:300` |
| `UrgeRecord.TargetUserName` | `targetUserName` | Go field: UrgeRecord.TargetUserName<br/>type: string | `approval/models.go:301` |
| `UrgeRecord.TaskID` | `taskId` | Go field: UrgeRecord.TaskID<br/>type: *string | `approval/models.go:297` |
| `UrgeRecord.UrgerID` | `urgerId` | Go field: UrgeRecord.UrgerID<br/>type: string | `approval/models.go:298` |
| `UrgeRecord.UrgerName` | `urgerName` | Go field: UrgeRecord.UrgerName<br/>type: string | `approval/models.go:299` |
| `UrgeTaskParams.Message` | `message` | Go field: UrgeTaskParams.Message<br/>type: string<br/>validate: "max=500" | `internal/approval/resource/instance.go:388` |
| `UrgeTaskParams.TaskID` | `taskId` | Go field: UrgeTaskParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:387` |
| `UserInfo.Avatar` | `avatar` | Go field: UserInfo.Avatar<br/>type: *string | `security/user_info.go:45` |
| `UserInfo.Details` | `details` | Go field: UserInfo.Details<br/>type: any | `security/user_info.go:48` |
| `UserInfo.Gender` | `gender` | Go field: UserInfo.Gender<br/>type: Gender | `security/user_info.go:44` |
| `UserInfo.ID` | `id` | Go field: UserInfo.ID<br/>type: string | `approval/assignee.go:11` |
| `UserInfo.ID` | `id` | Go field: UserInfo.ID<br/>type: string | `security/user_info.go:42` |
| `UserInfo.Menus` | `menus` | Go field: UserInfo.Menus<br/>type: []UserMenu | `security/user_info.go:47` |
| `UserInfo.Name` | `name` | Go field: UserInfo.Name<br/>type: string | `security/user_info.go:43` |
| `UserInfo.Name` | `name` | Go field: UserInfo.Name<br/>type: string | `approval/assignee.go:12` |
| `UserInfo.PermissionTokens` | `permissionTokens` | Go field: UserInfo.PermissionTokens<br/>type: []string | `security/user_info.go:46` |
| `UserMenu.Children` | `children` | Go field: UserMenu.Children<br/>type: []UserMenu | `security/user_info.go:38` |
| `UserMenu.Icon` | `icon` | Go field: UserMenu.Icon<br/>type: *string | `security/user_info.go:36` |
| `UserMenu.Meta` | `metadata` | Go field: UserMenu.Meta<br/>type: map[string]any | `security/user_info.go:37` |
| `UserMenu.Name` | `name` | Go field: UserMenu.Name<br/>type: string | `security/user_info.go:35` |
| `UserMenu.Path` | `path` | Go field: UserMenu.Path<br/>type: string | `security/user_info.go:34` |
| `UserMenu.Type` | `type` | Go field: UserMenu.Type<br/>type: UserMenuType | `security/user_info.go:33` |
| `ValidationRule.Max` | `max` | Go field: ValidationRule.Max<br/>type: *float64 | `approval/form_definition.go:43` |
| `ValidationRule.MaxLength` | `maxLength` | Go field: ValidationRule.MaxLength<br/>type: *int | `approval/form_definition.go:41` |
| `ValidationRule.Message` | `message` | Go field: ValidationRule.Message<br/>type: string | `approval/form_definition.go:45` |
| `ValidationRule.Min` | `min` | Go field: ValidationRule.Min<br/>type: *float64 | `approval/form_definition.go:42` |
| `ValidationRule.MinLength` | `minLength` | Go field: ValidationRule.MinLength<br/>type: *int | `approval/form_definition.go:40` |
| `ValidationRule.Pattern` | `pattern` | Go field: ValidationRule.Pattern<br/>type: string | `approval/form_definition.go:44` |
| `View.Columns` | `columns` | Go field: View.Columns<br/>type: []string | `schema/service.go:76` |
| `View.Comment` | `comment` | Go field: View.Comment<br/>type: string | `schema/service.go:75` |
| `View.Definition` | `definition` | Go field: View.Definition<br/>type: string | `schema/service.go:74` |
| `View.Name` | `name` | Go field: View.Name<br/>type: string | `schema/service.go:72` |
| `View.Schema` | `schema` | Go field: View.Schema<br/>type: string | `schema/service.go:73` |
| `VirtualMemory.Active` | `active` | Go field: VirtualMemory.Active<br/>type: uint64 | `monitor/service.go:88` |
| `VirtualMemory.AnonHugePages` | `anonHugePages` | Go field: VirtualMemory.AnonHugePages<br/>type: uint64 | `monitor/service.go:120` |
| `VirtualMemory.Available` | `available` | Go field: VirtualMemory.Available<br/>type: uint64 | `monitor/service.go:84` |
| `VirtualMemory.Buffers` | `buffers` | Go field: VirtualMemory.Buffers<br/>type: uint64 | `monitor/service.go:92` |
| `VirtualMemory.Cached` | `cached` | Go field: VirtualMemory.Cached<br/>type: uint64 | `monitor/service.go:93` |
| `VirtualMemory.CommitLimit` | `commitLimit` | Go field: VirtualMemory.CommitLimit<br/>type: uint64 | `monitor/service.go:103` |
| `VirtualMemory.CommittedAs` | `committedAs` | Go field: VirtualMemory.CommittedAs<br/>type: uint64 | `monitor/service.go:104` |
| `VirtualMemory.Dirty` | `dirty` | Go field: VirtualMemory.Dirty<br/>type: uint64 | `monitor/service.go:95` |
| `VirtualMemory.Free` | `free` | Go field: VirtualMemory.Free<br/>type: uint64 | `monitor/service.go:87` |
| `VirtualMemory.HighFree` | `highFree` | Go field: VirtualMemory.HighFree<br/>type: uint64 | `monitor/service.go:106` |
| `VirtualMemory.HighTotal` | `highTotal` | Go field: VirtualMemory.HighTotal<br/>type: uint64 | `monitor/service.go:105` |
| `VirtualMemory.HugePageSize` | `hugePageSize` | Go field: VirtualMemory.HugePageSize<br/>type: uint64 | `monitor/service.go:119` |
| `VirtualMemory.HugePagesFree` | `hugePagesFree` | Go field: VirtualMemory.HugePagesFree<br/>type: uint64 | `monitor/service.go:116` |
| `VirtualMemory.HugePagesReserved` | `hugePagesReserved` | Go field: VirtualMemory.HugePagesReserved<br/>type: uint64 | `monitor/service.go:117` |
| `VirtualMemory.HugePagesSurplus` | `hugePagesSurplus` | Go field: VirtualMemory.HugePagesSurplus<br/>type: uint64 | `monitor/service.go:118` |
| `VirtualMemory.HugePagesTotal` | `hugePagesTotal` | Go field: VirtualMemory.HugePagesTotal<br/>type: uint64 | `monitor/service.go:115` |
| `VirtualMemory.Inactive` | `inactive` | Go field: VirtualMemory.Inactive<br/>type: uint64 | `monitor/service.go:89` |
| `VirtualMemory.Laundry` | `laundry` | Go field: VirtualMemory.Laundry<br/>type: uint64 | `monitor/service.go:91` |
| `VirtualMemory.LowFree` | `lowFree` | Go field: VirtualMemory.LowFree<br/>type: uint64 | `monitor/service.go:108` |
| `VirtualMemory.LowTotal` | `lowTotal` | Go field: VirtualMemory.LowTotal<br/>type: uint64 | `monitor/service.go:107` |
| `VirtualMemory.Mapped` | `mapped` | Go field: VirtualMemory.Mapped<br/>type: uint64 | `monitor/service.go:111` |
| `VirtualMemory.PageTables` | `pageTables` | Go field: VirtualMemory.PageTables<br/>type: uint64 | `monitor/service.go:101` |
| `VirtualMemory.Shared` | `shared` | Go field: VirtualMemory.Shared<br/>type: uint64 | `monitor/service.go:97` |
| `VirtualMemory.Slab` | `slab` | Go field: VirtualMemory.Slab<br/>type: uint64 | `monitor/service.go:98` |
| `VirtualMemory.SlabReclaimable` | `slabReclaimable` | Go field: VirtualMemory.SlabReclaimable<br/>type: uint64 | `monitor/service.go:99` |
| `VirtualMemory.SlabUnreclaimable` | `slabUnreclaimable` | Go field: VirtualMemory.SlabUnreclaimable<br/>type: uint64 | `monitor/service.go:100` |
| `VirtualMemory.SwapCached` | `swapCached` | Go field: VirtualMemory.SwapCached<br/>type: uint64 | `monitor/service.go:102` |
| `VirtualMemory.SwapFree` | `swapFree` | Go field: VirtualMemory.SwapFree<br/>type: uint64 | `monitor/service.go:110` |
| `VirtualMemory.SwapTotal` | `swapTotal` | Go field: VirtualMemory.SwapTotal<br/>type: uint64 | `monitor/service.go:109` |
| `VirtualMemory.Total` | `total` | Go field: VirtualMemory.Total<br/>type: uint64 | `monitor/service.go:83` |
| `VirtualMemory.Used` | `used` | Go field: VirtualMemory.Used<br/>type: uint64 | `monitor/service.go:85` |
| `VirtualMemory.UsedPercent` | `usedPercent` | Go field: VirtualMemory.UsedPercent<br/>type: float64 | `monitor/service.go:86` |
| `VirtualMemory.VMAllocChunk` | `vmAllocChunk` | Go field: VirtualMemory.VMAllocChunk<br/>type: uint64 | `monitor/service.go:114` |
| `VirtualMemory.VMAllocTotal` | `vmAllocTotal` | Go field: VirtualMemory.VMAllocTotal<br/>type: uint64 | `monitor/service.go:112` |
| `VirtualMemory.VMAllocUsed` | `vmAllocUsed` | Go field: VirtualMemory.VMAllocUsed<br/>type: uint64 | `monitor/service.go:113` |
| `VirtualMemory.Wired` | `wired` | Go field: VirtualMemory.Wired<br/>type: uint64 | `monitor/service.go:90` |
| `VirtualMemory.WriteBack` | `writeBack` | Go field: VirtualMemory.WriteBack<br/>type: uint64 | `monitor/service.go:94` |
| `VirtualMemory.WriteBackTmp` | `writeBackTmp` | Go field: VirtualMemory.WriteBackTmp<br/>type: uint64 | `monitor/service.go:96` |
| `WithdrawParams.InstanceID` | `instanceId` | Go field: WithdrawParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:227` |
| `WithdrawParams.Reason` | `reason` | Go field: WithdrawParams.Reason<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/instance.go:228` |
| `exportConfig.Format` | `format` | Go field: exportConfig.Format<br/>type: TabularFormat | `crud/export.go:87` |
| `importConfig.Format` | `format` | Go field: importConfig.Format<br/>type: TabularFormat | `crud/import.go:92` |
| `importParams.File` | `file` | Go field: importParams.File<br/>type: *multipart.FileHeader | `crud/import.go:86` |
| `manifest.ContentType` | `contentType` | Go field: manifest.ContentType<br/>type: string | `internal/storage/filesystem/service.go:60` |
| `manifest.Key` | `key` | Go field: manifest.Key<br/>type: string | `internal/storage/filesystem/service.go:59` |
| `objectMeta.ContentType` | `contentType` | Go field: objectMeta.ContentType<br/>type: string | `internal/storage/filesystem/service.go:183` |
| `objectMeta.ETag` | `etag` | Go field: objectMeta.ETag<br/>type: string | `internal/storage/filesystem/service.go:182` |

## MCP endpoint

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `MCP Streamable HTTP endpoint` | `/mcp` | all HTTP methods | `internal/mcp/middleware.go:10` |

## MCP jsonschema tag

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `-` | `-` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `anchor` | `anchor` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `anyof_ref` | `anyof_ref` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `anyof_required` | `anyof_required` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `anyof_type` | `anyof_type` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `default` | `default` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `description` | `description` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `enum` | `enum` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `example` | `example` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `exclusiveMaximum` | `exclusiveMaximum` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `exclusiveMinimum` | `exclusiveMinimum` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `format` | `format` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `jsonschema_description` | `jsonschema_description` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `jsonschema_extras` | `jsonschema_extras` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `maxItems` | `maxItems` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `maxLength` | `maxLength` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `maximum` | `maximum` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `minItems` | `minItems` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `minLength` | `minLength` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `minimum` | `minimum` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `multipleOf` | `multipleOf` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `nullable` | `nullable` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `oneof_ref` | `oneof_ref` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `oneof_required` | `oneof_required` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `oneof_type` | `oneof_type` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `pattern` | `pattern` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `readOnly` | `readOnly` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `required` | `required` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `title` | `title` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `type` | `type` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `uniqueItems` | `uniqueItems` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |
| `writeOnly` | `writeOnly` |  | `github.com/invopop/jsonschema@v0.14.0/reflect.go:613` |

## MCP prompt

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `naming-master` | `naming-master` | Senior IT naming expert for code identifiers and database objects. Provides professional naming schemes following industry standards for multiple languages (Java, TypeScript, Go, Rust, Python) and databases (PostgreSQL, MySQL, SQLite). Includes database design guidance on audit fields, indexes, constraints, and foreign key strategies. | `internal/mcp/prompts/naming_master.go:26` |

## MCP tool

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `database_query` | `database_query` | Execute a read-only (SELECT) parameterized SQL query against the database. Returns query results as JSON array. | `internal/mcp/tools/query.go:33` |

## REST action verb

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `ALL` | `all` |  | `api/resource.go:86` |
| `CONNECT` | `connect` |  | `api/resource.go:85` |
| `DELETE` | `delete` |  | `api/resource.go:80` |
| `GET` | `get` |  | `api/resource.go:77` |
| `HEAD` | `head` |  | `api/resource.go:82` |
| `OPTIONS` | `options` |  | `api/resource.go:83` |
| `PATCH` | `patch` |  | `api/resource.go:81` |
| `POST` | `post` |  | `api/resource.go:78` |
| `PUT` | `put` |  | `api/resource.go:79` |
| `TRACE` | `trace` |  | `api/resource.go:84` |

## RPC form key

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `FormKeyMeta` | `meta` |  | `internal/api/router/rpc.go:23` |
| `FormKeyParams` | `params` |  | `internal/api/router/rpc.go:22` |

## auth strategy

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AuthStrategyBearer` | `bearer` |  | `api/auth.go:32` |
| `AuthStrategyNone` | `none` |  | `api/auth.go:31` |
| `AuthStrategySignature` | `signature` |  | `api/auth.go:33` |

## auth type

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AuthTypePassword` | `password` |  | `internal/security/password_authenticator.go:15` |
| `AuthTypeRefresh` | `refresh` |  | `internal/security/jwt_refresh_authenticator.go:13` |
| `AuthTypeSignature` | `signature` |  | `internal/security/signature_authenticator.go:14` |
| `AuthTypeToken` | `token` |  | `internal/security/jwt_token_authenticator.go:11` |

## built-in resource

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `approval/admin` | `rpc` |  | `internal/approval/resource/admin.go:38` |
| `approval/category` | `rpc` |  | `internal/approval/resource/category.go:98` |
| `approval/delegation` | `rpc` |  | `internal/approval/resource/delegation.go:52` |
| `approval/flow` | `rpc` |  | `internal/approval/resource/flow.go:30` |
| `approval/instance` | `rpc` |  | `internal/approval/resource/instance.go:98` |
| `approval/my` | `rpc` |  | `internal/approval/resource/my.go:29` |
| `security/auth` | `rpc` |  | `internal/security/auth_resource.go:49` |
| `sys/monitor` | `rpc` |  | `internal/monitor/resource.go:18` |
| `sys/schema` | `rpc` |  | `internal/schema/resource.go:17` |
| `sys/storage` | `rpc` |  | `internal/storage/resource.go:176` |

## built-in resource action

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `approval/admin/find_action_logs` | `find_action_logs` | permission: approval:log:query<br/>resource kind: rpc | `internal/approval/resource/admin.go:44` |
| `approval/admin/find_instances` | `find_instances` | permission: approval:instance:query<br/>resource kind: rpc | `internal/approval/resource/admin.go:41` |
| `approval/admin/find_tasks` | `find_tasks` | permission: approval:task:query<br/>resource kind: rpc | `internal/approval/resource/admin.go:42` |
| `approval/admin/get_instance_detail` | `get_instance_detail` | permission: approval:instance:detail<br/>resource kind: rpc | `internal/approval/resource/admin.go:43` |
| `approval/admin/get_metrics` | `get_metrics` | permission: approval:metrics:query<br/>resource kind: rpc | `internal/approval/resource/admin.go:45` |
| `approval/admin/reassign_task` | `reassign_task` | audit enabled<br/>permission: approval:task:reassign<br/>resource kind: rpc | `internal/approval/resource/admin.go:49` |
| `approval/admin/terminate_instance` | `terminate_instance` | audit enabled<br/>permission: approval:instance:terminate<br/>resource kind: rpc | `internal/approval/resource/admin.go:48` |
| `approval/category/create` | `create` | permission: approval:category:create<br/>resource kind: rpc | `internal/approval/resource/category.go:105` |
| `approval/category/delete` | `delete` | permission: approval:category:delete<br/>resource kind: rpc | `internal/approval/resource/category.go:142` |
| `approval/category/find_tree` | `find_tree` | permission: approval:category:query<br/>resource kind: rpc | `internal/approval/resource/category.go:99` |
| `approval/category/find_tree_options` | `find_tree_options` | permission: approval:category:query<br/>resource kind: rpc | `internal/approval/resource/category.go:102` |
| `approval/category/update` | `update` | permission: approval:category:update<br/>resource kind: rpc | `internal/approval/resource/category.go:127` |
| `approval/delegation/create` | `create` | permission: approval:delegation:create<br/>resource kind: rpc | `internal/approval/resource/delegation.go:71` |
| `approval/delegation/delete` | `delete` | permission: approval:delegation:delete<br/>resource kind: rpc | `internal/approval/resource/delegation.go:105` |
| `approval/delegation/find_page` | `find_page` | permission: approval:delegation:query<br/>resource kind: rpc | `internal/approval/resource/delegation.go:53` |
| `approval/delegation/update` | `update` | permission: approval:delegation:update<br/>resource kind: rpc | `internal/approval/resource/delegation.go:88` |
| `approval/flow/create` | `create` | audit enabled<br/>permission: approval:flow:create<br/>resource kind: rpc | `internal/approval/resource/flow.go:35` |
| `approval/flow/deploy` | `deploy` | audit enabled<br/>permission: approval:flow:deploy<br/>resource kind: rpc | `internal/approval/resource/flow.go:36` |
| `approval/flow/find_flows` | `find_flows` | permission: approval:flow:query<br/>resource kind: rpc | `internal/approval/resource/flow.go:41` |
| `approval/flow/find_versions` | `find_versions` | permission: approval:flow:query<br/>resource kind: rpc | `internal/approval/resource/flow.go:42` |
| `approval/flow/get_graph` | `get_graph` | permission: approval:flow:query<br/>resource kind: rpc | `internal/approval/resource/flow.go:40` |
| `approval/flow/publish_version` | `publish_version` | audit enabled<br/>permission: approval:flow:publish<br/>resource kind: rpc | `internal/approval/resource/flow.go:37` |
| `approval/flow/toggle_active` | `toggle_active` | audit enabled<br/>permission: approval:flow:update<br/>resource kind: rpc | `internal/approval/resource/flow.go:39` |
| `approval/flow/update_flow` | `update_flow` | audit enabled<br/>permission: approval:flow:update<br/>resource kind: rpc | `internal/approval/resource/flow.go:38` |
| `approval/instance/add_assignee` | `add_assignee` | audit enabled<br/>permission: approval:task:add_assignee<br/>resource kind: rpc | `internal/approval/resource/instance.go:111` |
| `approval/instance/add_cc` | `add_cc` | audit enabled<br/>permission: approval:instance:cc<br/>resource kind: rpc | `internal/approval/resource/instance.go:108` |
| `approval/instance/mark_cc_read` | `mark_cc_read` | permission: approval:instance:cc<br/>resource kind: rpc | `internal/approval/resource/instance.go:110` |
| `approval/instance/process_task` | `process_task` | audit enabled<br/>permission: approval:task:process<br/>resource kind: rpc | `internal/approval/resource/instance.go:105` |
| `approval/instance/remove_assignee` | `remove_assignee` | audit enabled<br/>permission: approval:task:remove_assignee<br/>resource kind: rpc | `internal/approval/resource/instance.go:112` |
| `approval/instance/resubmit` | `resubmit` | audit enabled<br/>permission: approval:instance:resubmit<br/>resource kind: rpc | `internal/approval/resource/instance.go:107` |
| `approval/instance/start` | `start` | audit enabled<br/>permission: approval:instance:start<br/>resource kind: rpc | `internal/approval/resource/instance.go:104` |
| `approval/instance/urge_task` | `urge_task` | permission: approval:task:urge<br/>resource kind: rpc | `internal/approval/resource/instance.go:115` |
| `approval/instance/withdraw` | `withdraw` | audit enabled<br/>permission: approval:instance:withdraw<br/>resource kind: rpc | `internal/approval/resource/instance.go:106` |
| `approval/my/find_available_flows` | `find_available_flows` | resource kind: rpc | `internal/approval/resource/my.go:32` |
| `approval/my/find_cc_records` | `find_cc_records` | resource kind: rpc | `internal/approval/resource/my.go:36` |
| `approval/my/find_completed_tasks` | `find_completed_tasks` | resource kind: rpc | `internal/approval/resource/my.go:35` |
| `approval/my/find_initiated` | `find_initiated` | resource kind: rpc | `internal/approval/resource/my.go:33` |
| `approval/my/find_pending_tasks` | `find_pending_tasks` | resource kind: rpc | `internal/approval/resource/my.go:34` |
| `approval/my/get_instance_detail` | `get_instance_detail` | resource kind: rpc | `internal/approval/resource/my.go:38` |
| `approval/my/get_pending_counts` | `get_pending_counts` | resource kind: rpc | `internal/approval/resource/my.go:37` |
| `security/auth/get_user_info` | `get_user_info` | resource kind: rpc | `internal/security/auth_resource.go:70` |
| `security/auth/login` | `login` | public<br/>resource kind: rpc | `internal/security/auth_resource.go:52` |
| `security/auth/logout` | `logout` | resource kind: rpc | `internal/security/auth_resource.go:62` |
| `security/auth/refresh` | `refresh` | public<br/>resource kind: rpc | `internal/security/auth_resource.go:57` |
| `security/auth/resolve_challenge` | `resolve_challenge` | public<br/>resource kind: rpc | `internal/security/auth_resource.go:65` |
| `sys/monitor/get_build_info` | `get_build_info` | resource kind: rpc | `internal/monitor/resource.go:29` |
| `sys/monitor/get_cpu` | `get_cpu` | resource kind: rpc | `internal/monitor/resource.go:22` |
| `sys/monitor/get_disk` | `get_disk` | resource kind: rpc | `internal/monitor/resource.go:24` |
| `sys/monitor/get_host` | `get_host` | resource kind: rpc | `internal/monitor/resource.go:26` |
| `sys/monitor/get_load` | `get_load` | resource kind: rpc | `internal/monitor/resource.go:28` |
| `sys/monitor/get_memory` | `get_memory` | resource kind: rpc | `internal/monitor/resource.go:23` |
| `sys/monitor/get_network` | `get_network` | resource kind: rpc | `internal/monitor/resource.go:25` |
| `sys/monitor/get_overview` | `get_overview` | resource kind: rpc | `internal/monitor/resource.go:21` |
| `sys/monitor/get_process` | `get_process` | resource kind: rpc | `internal/monitor/resource.go:27` |
| `sys/schema/get_table_schema` | `get_table_schema` | resource kind: rpc | `internal/schema/resource.go:21` |
| `sys/schema/list_tables` | `list_tables` | resource kind: rpc | `internal/schema/resource.go:20` |
| `sys/schema/list_views` | `list_views` | resource kind: rpc | `internal/schema/resource.go:22` |
| `sys/storage/abort_upload` | `abort_upload` | resource kind: rpc | `internal/storage/resource.go:183` |
| `sys/storage/complete_upload` | `complete_upload` | resource kind: rpc | `internal/storage/resource.go:182` |
| `sys/storage/init_upload` | `init_upload` | resource kind: rpc | `internal/storage/resource.go:179` |
| `sys/storage/list_parts` | `list_parts` | resource kind: rpc | `internal/storage/resource.go:181` |
| `sys/storage/upload_part` | `upload_part` | resource kind: rpc | `internal/storage/resource.go:180` |

## config default

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `vef.approval.cc_record_retention` | `90 * 24 * time.Hour` |  | `config/approval.go:72` |
| `vef.approval.cleanup_scan_interval` | `24 * time.Hour` |  | `config/approval.go:56` |
| `vef.approval.delegation_max_depth` | `10` |  | `config/approval.go:60` |
| `vef.approval.form_snapshot_retention` | `90 * 24 * time.Hour` |  | `config/approval.go:64` |
| `vef.approval.pre_warning_scan_interval` | `5 * time.Minute` |  | `config/approval.go:52` |
| `vef.approval.timeout_scan_interval` | `time.Minute` |  | `config/approval.go:48` |
| `vef.approval.urge_record_retention` | `30 * 24 * time.Hour` |  | `config/approval.go:68` |
| `vef.event.async_queue_size` | `4096` |  | `config/event.go:116` |
| `vef.event.async_workers` | `4` |  | `config/event.go:121` |
| `vef.event.default_transport` | `memory` |  | `config/event.go:111` |
| `vef.event.inbox.cleanup_interval` | `time.Hour` |  | `config/event.go:151` |
| `vef.event.inbox.processing_lease` | `10 * time.Minute` |  | `config/event.go:146` |
| `vef.event.inbox.retention` | `7 * 24 * time.Hour` |  | `config/event.go:141` |
| `vef.event.publish_timeout` | `5 * time.Second` |  | `config/event.go:126` |
| `vef.event.transports.memory.full_policy` | `"error" (error)` |  | `event/transport/memory/memory.go:48` |
| `vef.event.transports.memory.queue_size` | `1024` |  | `event/transport/memory/memory.go:39` |
| `vef.event.transports.outbox.batch_size` | `100` |  | `event/transport/outbox/outbox.go:121` |
| `vef.event.transports.outbox.cleanup_interval` | `time.Hour` |  | `config/event.go:131` |
| `vef.event.transports.outbox.completed_ttl` | `7 * 24 * time.Hour` |  | `config/event.go:136` |
| `vef.event.transports.outbox.lease_multiplier` | `4` |  | `event/transport/outbox/outbox.go:130` |
| `vef.event.transports.outbox.max_retries` | `10` | EventConfig.Validate fallback when max_retries is unset | `config/event.go:175` |
| `vef.event.transports.outbox.max_retries` | `10` |  | `event/transport/outbox/outbox.go:112` |
| `vef.event.transports.outbox.min_lease` | `15 * time.Second` |  | `event/transport/outbox/outbox.go:139` |
| `vef.event.transports.outbox.relay_interval` | `10 * time.Second` |  | `event/transport/outbox/outbox.go:103` |
| `vef.event.transports.outbox.sink` | `memory` |  | `event/transport/outbox/outbox.go:148` |
| `vef.event.transports.redis_stream.block_timeout` | `5 * time.Second` |  | `event/transport/redisstream/redis_stream.go:50` |
| `vef.event.transports.redis_stream.claim_batch_size` | `64` |  | `event/transport/redisstream/redis_stream.go:77` |
| `vef.event.transports.redis_stream.claim_idle` | `60 * time.Second` |  | `event/transport/redisstream/redis_stream.go:59` |
| `vef.event.transports.redis_stream.claim_interval` | `30 * time.Second` |  | `event/transport/redisstream/redis_stream.go:68` |
| `vef.event.transports.redis_stream.start_id` | `0` |  | `event/transport/redisstream/redis_stream.go:86` |
| `vef.event.transports.redis_stream.stream_prefix` | `vef:events:` |  | `event/transport/redisstream/redis_stream.go:41` |
| `vef.mcp.require_auth` | `true when unset` |  | `internal/mcp/handler.go:34` |
| `vef.monitor.sample_duration` | `2 * time.Second` |  | `internal/monitor/config.go:25` |
| `vef.monitor.sample_interval` | `10 * time.Second` |  | `internal/monitor/config.go:24` |
| `vef.storage.claim_ttl` | `24 * time.Hour` |  | `config/storage.go:114` |
| `vef.storage.delete_batch_size` | `100 (100)` |  | `config/storage.go:139` |
| `vef.storage.delete_concurrency` | `8 (8)` |  | `config/storage.go:144` |
| `vef.storage.delete_lease_window` | `5 * time.Minute` |  | `config/storage.go:154` |
| `vef.storage.delete_max_attempts` | `12 (12)` |  | `config/storage.go:149` |
| `vef.storage.delete_worker_interval` | `5 * time.Minute` |  | `config/storage.go:134` |
| `vef.storage.max_pending_claims` | `100 (100)` |  | `config/storage.go:119` |
| `vef.storage.max_upload_size` | `1024 * 1024 * 1024` |  | `config/storage.go:109` |
| `vef.storage.sweep_batch_size` | `200 (200)` |  | `config/storage.go:129` |
| `vef.storage.sweep_interval` | `5 * time.Minute` |  | `config/storage.go:124` |

## config enum

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `MySQL` | `mysql` |  | `config/data_sources.go:17` |
| `Oracle` | `oracle` |  | `config/data_sources.go:14` |
| `Postgres` | `postgres` |  | `config/data_sources.go:16` |
| `SQLServer` | `sqlserver` |  | `config/data_sources.go:15` |
| `SQLite` | `sqlite` |  | `config/data_sources.go:18` |
| `StorageFilesystem` | `filesystem` |  | `config/storage.go:15` |
| `StorageMemory` | `memory` |  | `config/storage.go:14` |
| `StorageMinIO` | `minio` |  | `config/storage.go:13` |

## config key

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `vef.app.body_limit` | `string` | Go field: AppConfig.BodyLimit | `config/app.go:13` |
| `vef.app.name` | `string` | Go field: AppConfig.Name | `config/app.go:11` |
| `vef.app.port` | `uint16` | Go field: AppConfig.Port | `config/app.go:12` |
| `vef.app.trusted_proxies` | `[]string` | Go field: AppConfig.TrustedProxies | `config/app.go:17` |
| `vef.approval.auto_migrate` | `bool` | Go field: ApprovalConfig.AutoMigrate | `config/approval.go:12` |
| `vef.approval.cc_record_retention` | `time.Duration` | Go field: ApprovalConfig.CCRecordRetention | `config/approval.go:41` |
| `vef.approval.cleanup_scan_interval` | `time.Duration` | Go field: ApprovalConfig.CleanupScanInterval | `config/approval.go:25` |
| `vef.approval.delegation_max_depth` | `int` | Go field: ApprovalConfig.DelegationMaxDepth | `config/approval.go:29` |
| `vef.approval.form_snapshot_retention` | `time.Duration` | Go field: ApprovalConfig.FormSnapshotRetention | `config/approval.go:33` |
| `vef.approval.pre_warning_scan_interval` | `time.Duration` | Go field: ApprovalConfig.PreWarningScanInterval | `config/approval.go:20` |
| `vef.approval.timeout_scan_interval` | `time.Duration` | Go field: ApprovalConfig.TimeoutScanInterval | `config/approval.go:16` |
| `vef.approval.urge_record_retention` | `time.Duration` | Go field: ApprovalConfig.UrgeRecordRetention | `config/approval.go:37` |
| `vef.cors.allow_origins` | `[]string` | Go field: CorsConfig.AllowOrigins | `config/cors.go:6` |
| `vef.cors.enabled` | `bool` | Go field: CorsConfig.Enabled | `config/cors.go:5` |
| `vef.data_sources.&lt;name&gt;.database` | `string` | Go field: DataSourceConfig.Database | `config/data_sources.go:22`, `config/data_sources.go:31`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.enable_sql_guard` | `bool` | Go field: DataSourceConfig.EnableSQLGuard | `config/data_sources.go:22`, `config/data_sources.go:34`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.host` | `string` | Go field: DataSourceConfig.Host | `config/data_sources.go:22`, `config/data_sources.go:27`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.password` | `string` | Go field: DataSourceConfig.Password | `config/data_sources.go:22`, `config/data_sources.go:30`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.path` | `string` | Go field: DataSourceConfig.Path | `config/data_sources.go:22`, `config/data_sources.go:33`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.port` | `uint16` | Go field: DataSourceConfig.Port | `config/data_sources.go:22`, `config/data_sources.go:28`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.schema` | `string` | Go field: DataSourceConfig.Schema | `config/data_sources.go:22`, `config/data_sources.go:32`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.type` | `DBKind` | Go field: DataSourceConfig.Kind | `config/data_sources.go:22`, `config/data_sources.go:26`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.user` | `string` | Go field: DataSourceConfig.User | `config/data_sources.go:22`, `config/data_sources.go:29`, `internal/config/data_sources.go:20` |
| `vef.event.async_queue_size` | `int` | Go field: EventConfig.AsyncQueueSize | `config/event.go:18` |
| `vef.event.async_workers` | `int` | Go field: EventConfig.AsyncWorkers | `config/event.go:21` |
| `vef.event.default_transport` | `string` | Go field: EventConfig.DefaultTransport | `config/event.go:15` |
| `vef.event.inbox` | `EventInboxConfig` | Go field: EventConfig.Inbox | `config/event.go:27` |
| `vef.event.inbox.cleanup_interval` | `time.Duration` | Go field: EventInboxConfig.CleanupInterval | `config/event.go:27`, `config/event.go:98` |
| `vef.event.inbox.processing_lease` | `time.Duration` | Go field: EventInboxConfig.ProcessingLease | `config/event.go:27`, `config/event.go:97` |
| `vef.event.inbox.retention` | `time.Duration` | Go field: EventInboxConfig.Retention | `config/event.go:27`, `config/event.go:96` |
| `vef.event.middleware` | `EventMiddlewareConfig` | Go field: EventConfig.Middleware | `config/event.go:26` |
| `vef.event.middleware.inbox` | `bool` | Go field: EventMiddlewareConfig.Inbox | `config/event.go:26`, `config/event.go:90` |
| `vef.event.middleware.logging` | `bool` | Go field: EventMiddlewareConfig.Logging | `config/event.go:26`, `config/event.go:76` |
| `vef.event.middleware.metrics` | `bool` | Go field: EventMiddlewareConfig.Metrics | `config/event.go:26`, `config/event.go:85` |
| `vef.event.middleware.recover` | `bool` | Go field: EventMiddlewareConfig.Recover | `config/event.go:26`, `config/event.go:86` |
| `vef.event.middleware.tracing` | `bool` | Go field: EventMiddlewareConfig.Tracing | `config/event.go:26`, `config/event.go:77` |
| `vef.event.middleware.tracing_strict` | `bool` | Go field: EventMiddlewareConfig.TracingStrict | `config/event.go:26`, `config/event.go:84` |
| `vef.event.publish_timeout` | `time.Duration` | Go field: EventConfig.PublishTimeout | `config/event.go:23` |
| `vef.event.routing` | `[]EventRoutingRule` | Go field: EventConfig.Routing | `config/event.go:31` |
| `vef.event.routing.pattern` | `string` | Go field: EventRoutingRule.Pattern | `config/event.go:106`, `config/event.go:31` |
| `vef.event.routing.transports` | `[]string` | Go field: EventRoutingRule.Transports | `config/event.go:107`, `config/event.go:31` |
| `vef.event.transports` | `EventTransportsConfig` | Go field: EventConfig.Transports | `config/event.go:25` |
| `vef.event.transports.memory` | `EventMemoryTransportConfig` | Go field: EventTransportsConfig.Memory | `config/event.go:25`, `config/event.go:36` |
| `vef.event.transports.memory.full_policy` | `string` | Go field: EventMemoryTransportConfig.FullPolicy | `config/event.go:25`, `config/event.go:36`, `config/event.go:44` |
| `vef.event.transports.memory.publish_timeout` | `time.Duration` | Go field: EventMemoryTransportConfig.PublishTimeout | `config/event.go:25`, `config/event.go:36`, `config/event.go:45` |
| `vef.event.transports.memory.queue_size` | `int` | Go field: EventMemoryTransportConfig.QueueSize | `config/event.go:25`, `config/event.go:36`, `config/event.go:43` |
| `vef.event.transports.outbox` | `EventOutboxTransportConfig` | Go field: EventTransportsConfig.Outbox | `config/event.go:25`, `config/event.go:37` |
| `vef.event.transports.outbox.batch_size` | `int` | Go field: EventOutboxTransportConfig.BatchSize | `config/event.go:25`, `config/event.go:37`, `config/event.go:53` |
| `vef.event.transports.outbox.cleanup_interval` | `time.Duration` | Go field: EventOutboxTransportConfig.CleanupInterval | `config/event.go:25`, `config/event.go:37`, `config/event.go:57` |
| `vef.event.transports.outbox.completed_ttl` | `time.Duration` | Go field: EventOutboxTransportConfig.CompletedTTL | `config/event.go:25`, `config/event.go:37`, `config/event.go:58` |
| `vef.event.transports.outbox.enabled` | `bool` | Go field: EventOutboxTransportConfig.Enabled | `config/event.go:25`, `config/event.go:37`, `config/event.go:50` |
| `vef.event.transports.outbox.lease_multiplier` | `int` | Go field: EventOutboxTransportConfig.LeaseMultiplier | `config/event.go:25`, `config/event.go:37`, `config/event.go:54` |
| `vef.event.transports.outbox.max_retries` | `int` | Go field: EventOutboxTransportConfig.MaxRetries | `config/event.go:25`, `config/event.go:37`, `config/event.go:52` |
| `vef.event.transports.outbox.min_lease` | `time.Duration` | Go field: EventOutboxTransportConfig.MinLease | `config/event.go:25`, `config/event.go:37`, `config/event.go:55` |
| `vef.event.transports.outbox.relay_interval` | `time.Duration` | Go field: EventOutboxTransportConfig.RelayInterval | `config/event.go:25`, `config/event.go:37`, `config/event.go:51` |
| `vef.event.transports.outbox.sink` | `string` | Go field: EventOutboxTransportConfig.SinkName | `config/event.go:25`, `config/event.go:37`, `config/event.go:56` |
| `vef.event.transports.redis_stream` | `EventRedisStreamTransportConfig` | Go field: EventTransportsConfig.RedisStream | `config/event.go:25`, `config/event.go:38` |
| `vef.event.transports.redis_stream.block_timeout` | `time.Duration` | Go field: EventRedisStreamTransportConfig.BlockTimeout | `config/event.go:25`, `config/event.go:38`, `config/event.go:66` |
| `vef.event.transports.redis_stream.claim_batch_size` | `int64` | Go field: EventRedisStreamTransportConfig.ClaimBatchSize | `config/event.go:25`, `config/event.go:38`, `config/event.go:69` |
| `vef.event.transports.redis_stream.claim_idle` | `time.Duration` | Go field: EventRedisStreamTransportConfig.ClaimIdle | `config/event.go:25`, `config/event.go:38`, `config/event.go:67` |
| `vef.event.transports.redis_stream.claim_interval` | `time.Duration` | Go field: EventRedisStreamTransportConfig.ClaimInterval | `config/event.go:25`, `config/event.go:38`, `config/event.go:68` |
| `vef.event.transports.redis_stream.consumer_id` | `string` | Go field: EventRedisStreamTransportConfig.ConsumerID | `config/event.go:25`, `config/event.go:38`, `config/event.go:70` |
| `vef.event.transports.redis_stream.enabled` | `bool` | Go field: EventRedisStreamTransportConfig.Enabled | `config/event.go:25`, `config/event.go:38`, `config/event.go:63` |
| `vef.event.transports.redis_stream.max_len_approx` | `int64` | Go field: EventRedisStreamTransportConfig.MaxLenApprox | `config/event.go:25`, `config/event.go:38`, `config/event.go:65` |
| `vef.event.transports.redis_stream.start_id` | `string` | Go field: EventRedisStreamTransportConfig.StartID | `config/event.go:25`, `config/event.go:38`, `config/event.go:71` |
| `vef.event.transports.redis_stream.stream_prefix` | `string` | Go field: EventRedisStreamTransportConfig.StreamPrefix | `config/event.go:25`, `config/event.go:38`, `config/event.go:64` |
| `vef.mcp.enabled` | `bool` | Go field: MCPConfig.Enabled | `config/mcp.go:5` |
| `vef.mcp.require_auth` | `*bool` | Go field: MCPConfig.RequireAuth | `config/mcp.go:9` |
| `vef.monitor.excluded_mounts` | `[]string` | Go field: MonitorConfig.ExcludedMounts | `config/monitor.go:13` |
| `vef.monitor.sample_duration` | `time.Duration` | Go field: MonitorConfig.SampleDuration | `config/monitor.go:8` |
| `vef.monitor.sample_interval` | `time.Duration` | Go field: MonitorConfig.SampleInterval | `config/monitor.go:7` |
| `vef.redis.database` | `uint8` | Go field: RedisConfig.Database | `config/redis.go:16` |
| `vef.redis.enabled` | `bool` | Go field: RedisConfig.Enabled | `config/redis.go:11` |
| `vef.redis.host` | `string` | Go field: RedisConfig.Host | `config/redis.go:12` |
| `vef.redis.network` | `string` | Go field: RedisConfig.Network | `config/redis.go:17` |
| `vef.redis.password` | `string` | Go field: RedisConfig.Password | `config/redis.go:15` |
| `vef.redis.port` | `uint16` | Go field: RedisConfig.Port | `config/redis.go:13` |
| `vef.redis.user` | `string` | Go field: RedisConfig.User | `config/redis.go:14` |
| `vef.security.login_rate_limit` | `int` | Go field: SecurityConfig.LoginRateLimit | `config/security.go:13` |
| `vef.security.refresh_not_before` | `time.Duration` | Go field: SecurityConfig.RefreshNotBefore | `config/security.go:12` |
| `vef.security.refresh_rate_limit` | `int` | Go field: SecurityConfig.RefreshRateLimit | `config/security.go:14` |
| `vef.security.secret` | `string` | Go field: SecurityConfig.Secret | `config/security.go:10` |
| `vef.security.token_expires` | `time.Duration` | Go field: SecurityConfig.TokenExpires | `config/security.go:11` |
| `vef.storage.allow_public_uploads` | `bool` | Go field: StorageConfig.AllowPublicUploads | `config/storage.go:49` |
| `vef.storage.auto_migrate` | `bool` | Go field: StorageConfig.AutoMigrate | `config/storage.go:26` |
| `vef.storage.claim_ttl` | `time.Duration` | Go field: StorageConfig.ClaimTTL | `config/storage.go:37` |
| `vef.storage.delete_batch_size` | `int` | Go field: StorageConfig.DeleteBatchSize | `config/storage.go:63` |
| `vef.storage.delete_concurrency` | `int` | Go field: StorageConfig.DeleteConcurrency | `config/storage.go:66` |
| `vef.storage.delete_lease_window` | `time.Duration` | Go field: StorageConfig.DeleteLeaseWindow | `config/storage.go:73` |
| `vef.storage.delete_max_attempts` | `int` | Go field: StorageConfig.DeleteMaxAttempts | `config/storage.go:69` |
| `vef.storage.delete_worker_interval` | `time.Duration` | Go field: StorageConfig.DeleteWorkerInterval | `config/storage.go:60` |
| `vef.storage.filesystem` | `FilesystemConfig` | Go field: StorageConfig.Filesystem | `config/storage.go:28` |
| `vef.storage.filesystem.root` | `string` | Go field: FilesystemConfig.Root | `config/storage.go:170`, `config/storage.go:28` |
| `vef.storage.max_pending_claims` | `int` | Go field: StorageConfig.MaxPendingClaims | `config/storage.go:43` |
| `vef.storage.max_upload_size` | `int64` | Go field: StorageConfig.MaxUploadSize | `config/storage.go:33` |
| `vef.storage.minio` | `MinIOConfig` | Go field: StorageConfig.MinIO | `config/storage.go:27` |
| `vef.storage.minio.access_key` | `string` | Go field: MinIOConfig.AccessKey | `config/storage.go:161`, `config/storage.go:27` |
| `vef.storage.minio.bucket` | `string` | Go field: MinIOConfig.Bucket | `config/storage.go:163`, `config/storage.go:27` |
| `vef.storage.minio.endpoint` | `string` | Go field: MinIOConfig.Endpoint | `config/storage.go:160`, `config/storage.go:27` |
| `vef.storage.minio.region` | `string` | Go field: MinIOConfig.Region | `config/storage.go:164`, `config/storage.go:27` |
| `vef.storage.minio.secret_key` | `string` | Go field: MinIOConfig.SecretKey | `config/storage.go:162`, `config/storage.go:27` |
| `vef.storage.minio.use_ssl` | `bool` | Go field: MinIOConfig.UseSSL | `config/storage.go:165`, `config/storage.go:27` |
| `vef.storage.provider` | `StorageProvider` | Go field: StorageConfig.Provider | `config/storage.go:25` |
| `vef.storage.sweep_batch_size` | `int` | Go field: StorageConfig.SweepBatchSize | `config/storage.go:56` |
| `vef.storage.sweep_interval` | `time.Duration` | Go field: StorageConfig.SweepInterval | `config/storage.go:53` |

## config reserved name

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `PrimaryDataSourceName` | `primary` | used under vef.data_sources.&lt;name&gt; | `config/data_sources.go:7` |

## environment variable

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `EnvConfigPath` | `VEF_CONFIG_PATH` |  | `config/env.go:8` |
| `EnvI18NLanguage` | `VEF_I18N_LANGUAGE` |  | `config/env.go:9` |
| `EnvKeyPrefix` | `VEF` |  | `config/env.go:5` |
| `EnvLogLevel` | `VEF_LOG_LEVEL` |  | `config/env.go:7` |
| `EnvNodeID` | `VEF_NODE_ID` |  | `config/env.go:6` |

## event topic

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `EventTypeAssigneesAdded` | `approval.task.assignees_added` |  | `approval/events.go:80` |
| `EventTypeAssigneesRemoved` | `approval.task.assignees_removed` |  | `approval/events.go:81` |
| `EventTypeCCNotified` | `approval.cc.notified` |  | `approval/events.go:85` |
| `EventTypeDeleteDeadLetter` | `vef.storage.delete.dead_letter` |  | `storage/events.go:24` |
| `EventTypeFileClaimed` | `vef.storage.file.claimed` |  | `storage/events.go:16` |
| `EventTypeFileDeleted` | `vef.storage.file.deleted` |  | `storage/events.go:20` |
| `EventTypeFlowCreated` | `approval.flow.created` |  | `approval/events.go:87` |
| `EventTypeFlowDeployed` | `approval.flow.deployed` |  | `approval/events.go:89` |
| `EventTypeFlowPublished` | `approval.flow.published` |  | `approval/events.go:91` |
| `EventTypeFlowToggled` | `approval.flow.toggled` |  | `approval/events.go:90` |
| `EventTypeFlowUpdated` | `approval.flow.updated` |  | `approval/events.go:88` |
| `EventTypeInstanceBindingFailed` | `approval.instance.binding_failed` |  | `approval/events.go:68` |
| `EventTypeInstanceCompleted` | `approval.instance.completed` |  | `approval/events.go:63` |
| `EventTypeInstanceCreated` | `approval.instance.created` |  | `approval/events.go:62` |
| `EventTypeInstanceResubmitted` | `approval.instance.resubmitted` |  | `approval/events.go:67` |
| `EventTypeInstanceReturned` | `approval.instance.returned` |  | `approval/events.go:66` |
| `EventTypeInstanceRolledBack` | `approval.instance.rolled_back` |  | `approval/events.go:65` |
| `EventTypeInstanceWithdrawn` | `approval.instance.withdrawn` |  | `approval/events.go:64` |
| `EventTypeNodeAutoPassed` | `approval.node.auto_passed` |  | `approval/events.go:71` |
| `EventTypeNodeEntered` | `approval.node.entered` |  | `approval/events.go:70` |
| `EventTypeTaskApproved` | `approval.task.approved` |  | `approval/events.go:74` |
| `EventTypeTaskCreated` | `approval.task.created` |  | `approval/events.go:73` |
| `EventTypeTaskDeadlineWarning` | `approval.task.deadline_warning` |  | `approval/events.go:82` |
| `EventTypeTaskHandled` | `approval.task.handled` |  | `approval/events.go:75` |
| `EventTypeTaskReassigned` | `approval.task.reassigned` |  | `approval/events.go:78` |
| `EventTypeTaskRejected` | `approval.task.rejected` |  | `approval/events.go:76` |
| `EventTypeTaskTimedOut` | `approval.task.timed_out` |  | `approval/events.go:79` |
| `EventTypeTaskTransferred` | `approval.task.transferred` |  | `approval/events.go:77` |
| `EventTypeTaskUrged` | `approval.task.urged` |  | `approval/events.go:83` |
| `eventTypeAudit` | `vef.api.request.audit` |  | `api/audit.go:9` |
| `eventTypeDictionaryChanged` | `vef.translate.dictionary.changed` |  | `mold/cached_dictionary_resolver.go:14` |
| `eventTypeLogin` | `vef.security.login` |  | `security/login_event.go:9` |
| `eventTypeRolePermissionsChanged` | `vef.security.role_permissions.changed` |  | `security/cached_role_permission_loader.go:14` |

## event transport contract

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `outbox DLQ header` | `vef.dlq` |  | `internal/event/transport/outbox/relay.go:37` |
| `outbox DLQ header value` | `1` |  | `internal/event/transport/outbox/relay.go:178` |
| `outbox DLQ topic prefix` | `vef-dlq.` |  | `internal/event/transport/outbox/relay.go:225` |
| `outbox persisted error max bytes` | `256` |  | `internal/event/transport/outbox/relay.go:20` |
| `outbox retry backoff cap` | `1h` |  | `internal/event/transport/outbox/relay.go:211` |
| `outbox retry backoff formula` | `2^retryCount seconds capped at 1h` |  | `internal/event/transport/outbox/relay.go:214` |

## i18n key indirection

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `internal/api/middleware/audit.go:139` | `dynamic key sourced from app.MapFiberError; mapped result/security message constants are indexed separately` |  | `internal/api/middleware/audit.go:139` |
| `internal/app/error.go:54` | `dynamic key sourced from fiberErrorMappings; mapped result/security message constants are indexed separately` |  | `internal/app/error.go:54` |
| `validator/rule.go:45` | `dynamic key sourced from ValidationRule.ErrMessageI18nKey; built-in rule keys are indexed separately` |  | `validator/rule.go:45` |
| `validator/validator.go:71` | `dynamic key sourced from label_i18n struct tags; tag values are indexed separately` |  | `validator/validator.go:71` |

## i18n message key

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `access_denied` | `access_denied` | i18n.T call | `result/errors.go:12` |
| `api_request_action` | `api_request_action` | label_i18n struct tag | `api/request.go:15` |
| `api_request_meta_invalid_json` | `api_request_meta_invalid_json` | i18n.T call | `api/api_errors.go:22` |
| `api_request_params_invalid_json` | `api_request_params_invalid_json` | i18n.T call | `api/api_errors.go:17` |
| `api_request_resource` | `api_request_resource` | label_i18n struct tag | `api/request.go:14` |
| `api_request_version` | `api_request_version` | label_i18n struct tag | `api/request.go:16` |
| `approval_access_denied` | `approval_access_denied` | i18n.T call | `internal/approval/shared/api_errors.go:72` |
| `approval_add_assignee_not_allowed` | `approval_add_assignee_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:44` |
| `approval_assignee_resolve_failed` | `approval_assignee_resolve_failed` | i18n.T call | `internal/approval/shared/api_errors.go:56` |
| `approval_delegation_conflict` | `approval_delegation_conflict` | i18n.T call | `internal/approval/shared/api_errors.go:70` |
| `approval_delegation_not_found` | `approval_delegation_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:69` |
| `approval_field_not_editable` | `approval_field_not_editable` | i18n.T call | `internal/approval/shared/api_errors.go:59` |
| `approval_flow_code_exists` | `approval_flow_code_exists` | i18n.T call | `internal/approval/shared/api_errors.go:21` |
| `approval_flow_not_active` | `approval_flow_not_active` | i18n.T call | `internal/approval/shared/api_errors.go:17` |
| `approval_flow_not_found` | `approval_flow_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:16` |
| `approval_form_data_too_large` | `approval_form_data_too_large` | i18n.T call | `internal/approval/shared/api_errors.go:65` |
| `approval_form_field_empty` | `approval_form_field_empty` | i18n.T call | `internal/approval/service/validation.go:292` |
| `approval_form_field_empty` | `approval_form_field_empty` | i18n.T call | `internal/approval/service/validation.go:285` |
| `approval_form_field_empty` | `approval_form_field_empty` | i18n.T call | `internal/approval/service/validation.go:278` |
| `approval_form_field_invalid_file_item` | `approval_form_field_invalid_file_item` | i18n.T call | `internal/approval/service/validation.go:298` |
| `approval_form_field_invalid_validation` | `approval_form_field_invalid_validation` | i18n.T call | `internal/approval/service/validation.go:241` |
| `approval_form_field_invalid_value` | `approval_form_field_invalid_value` | i18n.T call | `internal/approval/service/validation.go:321` |
| `approval_form_field_max_length` | `approval_form_field_max_length` | i18n.T call | `internal/approval/service/validation.go:232` |
| `approval_form_field_max_value` | `approval_form_field_max_value` | i18n.T call | `internal/approval/service/validation.go:265` |
| `approval_form_field_min_length` | `approval_form_field_min_length` | i18n.T call | `internal/approval/service/validation.go:225` |
| `approval_form_field_min_value` | `approval_form_field_min_value` | i18n.T call | `internal/approval/service/validation.go:258` |
| `approval_form_field_must_be_file` | `approval_form_field_must_be_file` | i18n.T call | `internal/approval/service/validation.go:305` |
| `approval_form_field_must_be_number` | `approval_form_field_must_be_number` | i18n.T call | `internal/approval/service/validation.go:206` |
| `approval_form_field_must_be_string` | `approval_form_field_must_be_string` | i18n.T call | `internal/approval/service/validation.go:195` |
| `approval_form_field_not_defined` | `approval_form_field_not_defined` | i18n.T call | `internal/approval/service/validation.go:77` |
| `approval_form_field_pattern_mismatch` | `approval_form_field_pattern_mismatch` | i18n.T call | `internal/approval/service/validation.go:245` |
| `approval_form_field_required` | `approval_form_field_required` | i18n.T call | `internal/approval/service/validation.go:85` |
| `approval_form_validation_failed` | `approval_form_validation_failed` | i18n.T call | `internal/approval/shared/api_errors.go:58` |
| `approval_instance_completed` | `approval_instance_completed` | i18n.T call | `internal/approval/shared/api_errors.go:33` |
| `approval_instance_not_found` | `approval_instance_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:32` |
| `approval_instance_not_running` | `approval_instance_not_running` | i18n.T call | `internal/approval/shared/api_errors.go:73` |
| `approval_invalid_add_assignee_type` | `approval_invalid_add_assignee_type` | i18n.T call | `internal/approval/shared/api_errors.go:49` |
| `approval_invalid_business_identifier` | `approval_invalid_business_identifier` | i18n.T call | `internal/approval/shared/api_errors.go:28` |
| `approval_invalid_flow_design` | `approval_invalid_flow_design` | i18n.T call | `internal/approval/shared/api_errors.go:20` |
| `approval_invalid_instance_transition` | `approval_invalid_instance_transition` | i18n.T call | `internal/approval/shared/api_errors.go:37` |
| `approval_invalid_rollback_target` | `approval_invalid_rollback_target` | i18n.T call | `internal/approval/shared/api_errors.go:51` |
| `approval_invalid_task_transition` | `approval_invalid_task_transition` | i18n.T call | `internal/approval/shared/api_errors.go:42` |
| `approval_invalid_transfer_target` | `approval_invalid_transfer_target` | i18n.T call | `internal/approval/shared/api_errors.go:53` |
| `approval_last_assignee_removal` | `approval_last_assignee_removal` | i18n.T call | `internal/approval/shared/api_errors.go:52` |
| `approval_manual_cc_not_allowed` | `approval_manual_cc_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:47` |
| `approval_no_assignee` | `approval_no_assignee` | i18n.T call | `internal/approval/shared/api_errors.go:55` |
| `approval_no_published_version` | `approval_no_published_version` | i18n.T call | `internal/approval/shared/api_errors.go:18` |
| `approval_not_allowed_initiate` | `approval_not_allowed_initiate` | i18n.T call | `internal/approval/shared/api_errors.go:34` |
| `approval_not_applicant` | `approval_not_applicant` | i18n.T call | `internal/approval/shared/api_errors.go:50` |
| `approval_not_assignee` | `approval_not_assignee` | i18n.T call | `internal/approval/shared/api_errors.go:41` |
| `approval_opinion_required` | `approval_opinion_required` | i18n.T call | `internal/approval/shared/api_errors.go:46` |
| `approval_remove_assignee_not_allowed` | `approval_remove_assignee_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:48` |
| `approval_resubmit_not_allowed` | `approval_resubmit_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:36` |
| `approval_rollback_not_allowed` | `approval_rollback_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:43` |
| `approval_task_not_found` | `approval_task_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:39` |
| `approval_task_not_pending` | `approval_task_not_pending` | i18n.T call | `internal/approval/shared/api_errors.go:40` |
| `approval_transfer_not_allowed` | `approval_transfer_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:45` |
| `approval_urge_too_frequent` | `approval_urge_too_frequent` | i18n.T call | `internal/approval/command/urge_task.go:115` |
| `approval_version_not_draft` | `approval_version_not_draft` | i18n.T call | `internal/approval/shared/api_errors.go:19` |
| `approval_version_not_found` | `approval_version_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:22` |
| `approval_withdraw_not_allowed` | `approval_withdraw_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:35` |
| `auth_challenge_response` | `auth_challenge_response` | label_i18n struct tag | `internal/security/auth_resource.go:216` |
| `auth_challenge_token` | `auth_challenge_token` | label_i18n struct tag | `internal/security/auth_resource.go:214` |
| `auth_challenge_type` | `auth_challenge_type` | label_i18n struct tag | `internal/security/auth_resource.go:215` |
| `auth_credentials` | `auth_credentials` | label_i18n struct tag | `internal/security/auth_resource.go:96` |
| `auth_principal` | `auth_principal` | label_i18n struct tag | `internal/security/auth_resource.go:95` |
| `auth_refresh_token` | `auth_refresh_token` | label_i18n struct tag | `internal/security/auth_resource.go:182` |
| `auth_type` | `auth_type` | label_i18n struct tag | `internal/security/auth_resource.go:94` |
| `crud_batch_create_list` | `crud_batch_create_list` | label_i18n struct tag | `crud/params.go:12` |
| `crud_batch_delete_pks` | `crud_batch_delete_pks` | label_i18n struct tag | `crud/params.go:28` |
| `crud_batch_update_list` | `crud_batch_update_list` | label_i18n struct tag | `crud/params.go:19` |
| `crud_composite_primary_key_requires_map` | `crud_composite_primary_key_requires_map` | i18n.T call | `crud/api_errors.go:27` |
| `crud_field_not_exist_in_model` | `crud_field_not_exist_in_model` | i18n.T call | `crud/api_errors.go:70` |
| `crud_file_open_failed` | `crud_file_open_failed` | i18n.T call | `crud/api_errors.go:47` |
| `crud_import_requires_file` | `crud_import_requires_file` | i18n.T call | `crud/api_errors.go:39` |
| `crud_import_requires_multipart` | `crud_import_requires_multipart` | i18n.T call | `crud/api_errors.go:35` |
| `crud_import_type_assertion_failed` | `crud_import_type_assertion_failed` | i18n.T call | `crud/api_errors.go:51` |
| `crud_import_validation_failed` | `crud_import_validation_failed` | i18n.T call | `crud/import.go:147` |
| `crud_primary_key_required` | `crud_primary_key_required` | i18n.T call | `crud/api_errors.go:60` |
| `crud_processor_must_return_slice` | `crud_processor_must_return_slice` | i18n.T call | `crud/find_page.go:88` |
| `crud_unsupported_export_format` | `crud_unsupported_export_format` | i18n.T call | `crud/api_errors.go:31` |
| `crud_unsupported_import_format` | `crud_unsupported_import_format` | i18n.T call | `crud/api_errors.go:43` |
| `dangerous_sql` | `dangerous_sql` | i18n.T call | `result/errors.go:49` |
| `error` | `error` | i18n.T call | `result/error.go:69` |
| `expression_evaluation_failed` | `expression_evaluation_failed` | i18n.T call | `expression/api_errors.go:15` |
| `foreign_key_violation` | `foreign_key_violation` | i18n.T call | `result/errors.go:45` |
| `monitor_collection_failed` | `monitor_collection_failed` | i18n.T call | `monitor/api_errors.go:25` |
| `monitor_not_ready` | `monitor_not_ready` | i18n.T call | `monitor/api_errors.go:19` |
| `ok` | `ok` | i18n.T call | `result/result.go:64` |
| `record_already_exists` | `record_already_exists` | i18n.T call | `result/errors.go:41` |
| `record_not_found` | `record_not_found` | i18n.T call | `result/errors.go:37` |
| `request_timeout` | `request_timeout` | i18n.T call | `result/errors.go:22` |
| `schema_table_not_found` | `schema_table_not_found` | i18n.T call | `schema/api_errors.go:16` |
| `security_app_id_required` | `security_app_id_required` | i18n.T call | `security/api_errors.go:98` |
| `security_auth_header_invalid` | `security_auth_header_invalid` | i18n.T call | `security/api_errors.go:163` |
| `security_auth_header_missing` | `security_auth_header_missing` | i18n.T call | `security/api_errors.go:158` |
| `security_challenge_token_invalid` | `security_challenge_token_invalid` | i18n.T call | `security/api_errors.go:172` |
| `security_challenge_type_invalid` | `security_challenge_type_invalid` | i18n.T call | `security/api_errors.go:177` |
| `security_credentials_format_invalid` | `security_credentials_format_invalid` | i18n.T call | `internal/security/signature_authenticator.go:54` |
| `security_department_required` | `security_department_required` | i18n.T call | `security/api_errors.go:197` |
| `security_external_app_disabled` | `security_external_app_disabled` | i18n.T call | `security/api_errors.go:133` |
| `security_external_app_loader_not_implemented` | `security_external_app_loader_not_implemented` | i18n.T call | `internal/security/signature_authenticator.go:44` |
| `security_external_app_not_found` | `security_external_app_not_found` | i18n.T call | `security/api_errors.go:128` |
| `security_invalid_credentials` | `security_invalid_credentials` | i18n.T call | `internal/security/password_authenticator.go:88` |
| `security_invalid_credentials` | `security_invalid_credentials` | i18n.T call | `internal/security/password_authenticator.go:77` |
| `security_invalid_credentials` | `security_invalid_credentials` | i18n.T call | `internal/security/password_authenticator.go:84` |
| `security_ip_not_allowed` | `security_ip_not_allowed` | i18n.T call | `security/api_errors.go:138` |
| `security_new_password_required` | `security_new_password_required` | i18n.T call | `security/api_errors.go:192` |
| `security_nonce_already_used` | `security_nonce_already_used` | i18n.T call | `security/api_errors.go:153` |
| `security_nonce_invalid` | `security_nonce_invalid` | i18n.T call | `security/api_errors.go:148` |
| `security_nonce_required` | `security_nonce_required` | i18n.T call | `security/api_errors.go:143` |
| `security_otp_code_invalid` | `security_otp_code_invalid` | i18n.T call | `security/api_errors.go:187` |
| `security_otp_code_required` | `security_otp_code_required` | i18n.T call | `security/api_errors.go:182` |
| `security_password_required` | `security_password_required` | i18n.T call | `internal/security/password_authenticator.go:62` |
| `security_signature_expired` | `security_signature_expired` | i18n.T call | `security/api_errors.go:118` |
| `security_signature_invalid` | `security_signature_invalid` | i18n.T call | `security/api_errors.go:123` |
| `security_signature_required` | `security_signature_required` | i18n.T call | `security/api_errors.go:108` |
| `security_system_principal_login_forbidden` | `security_system_principal_login_forbidden` | i18n.T call | `internal/security/password_authenticator.go:57` |
| `security_timestamp_invalid` | `security_timestamp_invalid` | i18n.T call | `security/api_errors.go:113` |
| `security_timestamp_required` | `security_timestamp_required` | i18n.T call | `security/api_errors.go:103` |
| `security_token_expired` | `security_token_expired` | i18n.T call | `security/api_errors.go:69` |
| `security_token_invalid` | `security_token_invalid` | i18n.T call | `security/api_errors.go:74` |
| `security_token_invalid_audience` | `security_token_invalid_audience` | i18n.T call | `security/api_errors.go:89` |
| `security_token_invalid_issuer` | `security_token_invalid_issuer` | i18n.T call | `security/api_errors.go:84` |
| `security_token_not_valid_yet` | `security_token_not_valid_yet` | i18n.T call | `security/api_errors.go:79` |
| `security_unauthenticated` | `security_unauthenticated` | i18n.T call | `security/api_errors.go:64` |
| `security_unsupported_authentication_type` | `security_unsupported_authentication_type` | i18n.T call | `internal/security/auth_manager.go:29` |
| `security_user_info_loader_not_implemented` | `security_user_info_loader_not_implemented` | i18n.T call | `internal/security/auth_resource.go:285` |
| `security_user_loader_not_implemented` | `security_user_loader_not_implemented` | i18n.T call | `internal/security/password_authenticator.go:48` |
| `security_user_loader_not_implemented` | `security_user_loader_not_implemented` | i18n.T call | `internal/security/jwt_refresh_authenticator.go:34` |
| `security_username_required` | `security_username_required` | i18n.T call | `internal/security/password_authenticator.go:53` |
| `storage_abort_failed` | `storage_abort_failed` | i18n.T call | `storage/api_errors.go:114` |
| `storage_claim_expired` | `storage_claim_expired` | i18n.T call | `storage/api_errors.go:54` |
| `storage_claim_not_multipart` | `storage_claim_not_multipart` | i18n.T call | `storage/api_errors.go:86` |
| `storage_claim_not_pending` | `storage_claim_not_pending` | i18n.T call | `storage/api_errors.go:50` |
| `storage_failed_to_get_file` | `storage_failed_to_get_file` | i18n.T call | `storage/api_errors.go:45` |
| `storage_file_not_found` | `storage_file_not_found` | i18n.T call | `storage/api_errors.go:41` |
| `storage_invalid_file_key` | `storage_invalid_file_key` | i18n.T call | `storage/api_errors.go:37` |
| `storage_multipart_not_supported` | `storage_multipart_not_supported` | i18n.T call | `storage/api_errors.go:62` |
| `storage_object_not_found` | `storage_object_not_found` | i18n.T call | `storage/api_errors.go:106` |
| `storage_part_number_out_of_range` | `storage_part_number_out_of_range` | i18n.T call | `storage/api_errors.go:90` |
| `storage_public_uploads_not_allowed` | `storage_public_uploads_not_allowed` | i18n.T call | `storage/api_errors.go:66` |
| `storage_too_many_pending_uploads` | `storage_too_many_pending_uploads` | i18n.T call | `storage/api_errors.go:74` |
| `storage_upload_part_too_large` | `storage_upload_part_too_large` | i18n.T call | `storage/api_errors.go:94` |
| `storage_upload_part_too_small` | `storage_upload_part_too_small` | i18n.T call | `storage/api_errors.go:98` |
| `storage_upload_parts_incomplete` | `storage_upload_parts_incomplete` | i18n.T call | `storage/api_errors.go:102` |
| `storage_upload_requires_file` | `storage_upload_requires_file` | i18n.T call | `storage/api_errors.go:82` |
| `storage_upload_requires_multipart` | `storage_upload_requires_multipart` | i18n.T call | `storage/api_errors.go:78` |
| `storage_upload_size_exceeds_limit` | `storage_upload_size_exceeds_limit` | i18n.T call | `storage/api_errors.go:58` |
| `storage_upload_size_mismatch` | `storage_upload_size_mismatch` | i18n.T call | `storage/api_errors.go:110` |
| `storage_upload_too_many_parts` | `storage_upload_too_many_parts` | i18n.T call | `storage/api_errors.go:70` |
| `too_many_requests` | `too_many_requests` | i18n.T call | `result/errors.go:17` |
| `unknown_error` | `unknown_error` | i18n.T call | `internal/app/error.go:64` |
| `unknown_error` | `unknown_error` | i18n.T call | `result/errors.go:27` |
| `validator_alphanum_us` | `validator_alphanum_us` | validator rule message key | `validator/alphanum.go:16` |
| `validator_alphanum_us_dot` | `validator_alphanum_us_dot` | validator rule message key | `validator/alphanum.go:24` |
| `validator_alphanum_us_slash` | `validator_alphanum_us_slash` | validator rule message key | `validator/alphanum.go:20` |
| `validator_decimal_max` | `validator_decimal_max` | validator rule message key | `validator/decimal.go:16` |
| `validator_decimal_min` | `validator_decimal_min` | validator rule message key | `validator/decimal.go:10` |
| `validator_phone_number` | `validator_phone_number` | validator rule message key | `validator/phone_number.go:8` |

## meta tag grammar

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `MetaTypeMarkdown` | `markdown` |  | `storage/file_refs.go:23` |
| `MetaTypeRichText` | `rich_text` |  | `storage/file_refs.go:21` |
| `MetaTypeUploadedFile` | `uploaded_file` |  | `storage/file_refs.go:19` |
| `attribute key/value delimiter` | `:` |  | `storage/file_refs.go:214` |
| `attribute pair delimiter` | `space` |  | `storage/file_refs.go:212` |
| `dive` | `dive` |  | `storage/file_refs.go:239` |
| `tag name` | `meta` |  | `storage/file_refs.go:26` |

## mold tag grammar

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `diveTag` | `dive` |  | `internal/mold/restricted.go:4` |
| `endKeysTag` | `endkeys` |  | `internal/mold/restricted.go:11` |
| `ignoreTag` | `-` |  | `internal/mold/restricted.go:7` |
| `keysTag` | `keys` |  | `internal/mold/restricted.go:10` |
| `restrictedTagChars` | `.[],\|=+()`~!@#$%^&amp;*\"/?&lt;&gt;{}` |  | `internal/mold/restricted.go:5` |
| `tag name` | `mold` |  | `internal/mold/mold.go:41` |
| `tagKeySeparator` | `=` |  | `internal/mold/restricted.go:8` |
| `tagSeparator` | `,` |  | `internal/mold/restricted.go:6` |
| `utf8HexComma` | `0x2C` |  | `internal/mold/restricted.go:9` |

## mold transformer tag

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `expr` | `expr` |  | `internal/expression/transformer.go:33` |
| `translate` | `translate` |  | `internal/mold/translate.go:190` |

## mold translate kind prefix

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `dict:` | `dict:` |  | `internal/mold/dictionary_translator.go:26` |

## result error code

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `ErrCodeAbortFailed` | `2219` |  | `storage/api_errors.go:30` |
| `ErrCodeAccessDenied` | `1100` |  | `result/constants.go:32` |
| `ErrCodeAppIDRequired` | `1009` |  | `security/api_errors.go:36` |
| `ErrCodeAuthHeaderInvalid` | `1022` |  | `security/api_errors.go:49` |
| `ErrCodeAuthHeaderMissing` | `1021` |  | `security/api_errors.go:48` |
| `ErrCodeBadRequest` | `1400` |  | `result/constants.go:41` |
| `ErrCodeChallengeResolveFailed` | `1034` |  | `security/api_errors.go:54` |
| `ErrCodeChallengeTokenInvalid` | `1031` |  | `security/api_errors.go:52` |
| `ErrCodeChallengeTypeInvalid` | `1033` |  | `security/api_errors.go:53` |
| `ErrCodeClaimExpired` | `2204` |  | `storage/api_errors.go:15` |
| `ErrCodeClaimNotMultipart` | `2212` |  | `storage/api_errors.go:23` |
| `ErrCodeClaimNotPending` | `2203` |  | `storage/api_errors.go:14` |
| `ErrCodeCollectionFailed` | `2101` |  | `monitor/api_errors.go:11` |
| `ErrCodeCompositePrimaryKeyRequiresMap` | `2403` |  | `crud/api_errors.go:13` |
| `ErrCodeCredentialsInvalid` | `1008` |  | `security/api_errors.go:35` |
| `ErrCodeDangerousSQL` | `1600` |  | `result/constants.go:49` |
| `ErrCodeDefault` | `2000` |  | `result/constants.go:55` |
| `ErrCodeDepartmentRequired` | `1038` |  | `security/api_errors.go:58` |
| `ErrCodeEvaluationFailed` | `2500` |  | `expression/api_errors.go:9` |
| `ErrCodeExternalAppDisabled` | `1015` |  | `security/api_errors.go:42` |
| `ErrCodeExternalAppNotFound` | `1014` |  | `security/api_errors.go:41` |
| `ErrCodeFailedToGetFile` | `2202` |  | `storage/api_errors.go:12` |
| `ErrCodeFieldNotExistInModel` | `2401` |  | `crud/api_errors.go:11` |
| `ErrCodeFileNotFound` | `2201` |  | `storage/api_errors.go:11` |
| `ErrCodeFileOpenFailed` | `2408` |  | `crud/api_errors.go:18` |
| `ErrCodeForeignKeyViolation` | `2003` |  | `result/constants.go:58` |
| `ErrCodeIPNotAllowed` | `1016` |  | `security/api_errors.go:43` |
| `ErrCodeImportRequiresFile` | `2406` |  | `crud/api_errors.go:16` |
| `ErrCodeImportRequiresMultipart` | `2405` |  | `crud/api_errors.go:15` |
| `ErrCodeImportTypeAssertionFailed` | `2409` |  | `crud/api_errors.go:19` |
| `ErrCodeImportValidationFailed` | `2410` |  | `crud/api_errors.go:20` |
| `ErrCodeInvalidFileKey` | `2200` |  | `storage/api_errors.go:10` |
| `ErrCodeMultipartNotSupported` | `2206` |  | `storage/api_errors.go:17` |
| `ErrCodeNewPasswordRequired` | `1037` |  | `security/api_errors.go:57` |
| `ErrCodeNonceAlreadyUsed` | `1020` |  | `security/api_errors.go:47` |
| `ErrCodeNonceInvalid` | `1019` |  | `security/api_errors.go:46` |
| `ErrCodeNonceRequired` | `1018` |  | `security/api_errors.go:45` |
| `ErrCodeNotFound` | `1200` |  | `result/constants.go:35` |
| `ErrCodeNotImplemented` | `1500` |  | `result/constants.go:46` |
| `ErrCodeNotReady` | `2100` |  | `monitor/api_errors.go:10` |
| `ErrCodeOTPCodeInvalid` | `1036` |  | `security/api_errors.go:56` |
| `ErrCodeOTPCodeRequired` | `1035` |  | `security/api_errors.go:55` |
| `ErrCodePrimaryKeyRequired` | `2402` |  | `crud/api_errors.go:12` |
| `ErrCodePrincipalInvalid` | `1007` |  | `security/api_errors.go:34` |
| `ErrCodePublicUploadsNotAllowed` | `2207` |  | `storage/api_errors.go:18` |
| `ErrCodeRecordAlreadyExists` | `2002` |  | `result/constants.go:57` |
| `ErrCodeRecordNotFound` | `2001` |  | `result/constants.go:56` |
| `ErrCodeRequestTimeout` | `1402` |  | `result/constants.go:43` |
| `ErrCodeSignatureExpired` | `1013` |  | `security/api_errors.go:40` |
| `ErrCodeSignatureInvalid` | `1017` |  | `security/api_errors.go:44` |
| `ErrCodeSignatureRequired` | `1011` |  | `security/api_errors.go:38` |
| `ErrCodeTableNotFound` | `2300` |  | `schema/api_errors.go:10` |
| `ErrCodeTimestampInvalid` | `1012` |  | `security/api_errors.go:39` |
| `ErrCodeTimestampRequired` | `1010` |  | `security/api_errors.go:37` |
| `ErrCodeTokenExpired` | `1002` |  | `security/api_errors.go:29` |
| `ErrCodeTokenInvalid` | `1003` |  | `security/api_errors.go:30` |
| `ErrCodeTokenInvalidAudience` | `1006` |  | `security/api_errors.go:33` |
| `ErrCodeTokenInvalidIssuer` | `1005` |  | `security/api_errors.go:32` |
| `ErrCodeTokenNotValidYet` | `1004` |  | `security/api_errors.go:31` |
| `ErrCodeTooManyPendingUploads` | `2209` |  | `storage/api_errors.go:20` |
| `ErrCodeTooManyRequests` | `1401` |  | `result/constants.go:42` |
| `ErrCodeUnauthenticated` | `1000` |  | `security/api_errors.go:27` |
| `ErrCodeUnknown` | `1900` |  | `result/constants.go:52` |
| `ErrCodeUnsupportedAuthenticationType` | `1001` |  | `security/api_errors.go:28` |
| `ErrCodeUnsupportedExportFormat` | `2404` |  | `crud/api_errors.go:14` |
| `ErrCodeUnsupportedImportFormat` | `2407` |  | `crud/api_errors.go:17` |
| `ErrCodeUnsupportedMediaType` | `1300` |  | `result/constants.go:38` |
| `ErrCodeUploadObjectNotFound` | `2217` |  | `storage/api_errors.go:28` |
| `ErrCodeUploadPartNumberOutOfRange` | `2213` |  | `storage/api_errors.go:24` |
| `ErrCodeUploadPartTooLarge` | `2214` |  | `storage/api_errors.go:25` |
| `ErrCodeUploadPartTooSmall` | `2215` |  | `storage/api_errors.go:26` |
| `ErrCodeUploadPartsIncomplete` | `2216` |  | `storage/api_errors.go:27` |
| `ErrCodeUploadRequiresFile` | `2211` |  | `storage/api_errors.go:22` |
| `ErrCodeUploadRequiresMultipart` | `2210` |  | `storage/api_errors.go:21` |
| `ErrCodeUploadSizeExceedsLimit` | `2205` |  | `storage/api_errors.go:16` |
| `ErrCodeUploadSizeMismatch` | `2218` |  | `storage/api_errors.go:29` |
| `ErrCodeUploadTooManyParts` | `2208` |  | `storage/api_errors.go:19` |

## result message key

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `ErrMessage` | `error` |  | `result/constants.go:8` |
| `ErrMessageAccessDenied` | `access_denied` |  | `result/constants.go:14` |
| `ErrMessageChallengeResolveFailed` | `security_challenge_resolve_failed` |  | `security/api_errors.go:21` |
| `ErrMessageCredentialsFormatInvalid` | `security_credentials_format_invalid` |  | `security/api_errors.go:17` |
| `ErrMessageDangerousSQL` | `dangerous_sql` |  | `result/constants.go:22` |
| `ErrMessageExternalAppLoaderNotImplemented` | `security_external_app_loader_not_implemented` |  | `security/api_errors.go:16` |
| `ErrMessageForeignKeyViolation` | `foreign_key_violation` |  | `result/constants.go:21` |
| `ErrMessageFormFieldEmpty` | `approval_form_field_empty` |  | `internal/approval/shared/messages.go:26` |
| `ErrMessageFormFieldInvalidFileItem` | `approval_form_field_invalid_file_item` |  | `internal/approval/shared/messages.go:27` |
| `ErrMessageFormFieldInvalidValidation` | `approval_form_field_invalid_validation` |  | `internal/approval/shared/messages.go:22` |
| `ErrMessageFormFieldInvalidValue` | `approval_form_field_invalid_value` |  | `internal/approval/shared/messages.go:29` |
| `ErrMessageFormFieldMaxLength` | `approval_form_field_max_length` |  | `internal/approval/shared/messages.go:21` |
| `ErrMessageFormFieldMaxValue` | `approval_form_field_max_value` |  | `internal/approval/shared/messages.go:25` |
| `ErrMessageFormFieldMinLength` | `approval_form_field_min_length` |  | `internal/approval/shared/messages.go:20` |
| `ErrMessageFormFieldMinValue` | `approval_form_field_min_value` |  | `internal/approval/shared/messages.go:24` |
| `ErrMessageFormFieldMustBeFile` | `approval_form_field_must_be_file` |  | `internal/approval/shared/messages.go:28` |
| `ErrMessageFormFieldMustBeNumber` | `approval_form_field_must_be_number` |  | `internal/approval/shared/messages.go:19` |
| `ErrMessageFormFieldMustBeString` | `approval_form_field_must_be_string` |  | `internal/approval/shared/messages.go:18` |
| `ErrMessageFormFieldNotDefined` | `approval_form_field_not_defined` |  | `internal/approval/shared/messages.go:16` |
| `ErrMessageFormFieldPatternMismatch` | `approval_form_field_pattern_mismatch` |  | `internal/approval/shared/messages.go:23` |
| `ErrMessageFormFieldRequired` | `approval_form_field_required` |  | `internal/approval/shared/messages.go:17` |
| `ErrMessageNotFound` | `not_found` |  | `result/constants.go:12` |
| `ErrMessageRecordAlreadyExists` | `record_already_exists` |  | `result/constants.go:20` |
| `ErrMessageRecordNotFound` | `record_not_found` |  | `result/constants.go:19` |
| `ErrMessageRequestTimeout` | `request_timeout` |  | `result/constants.go:16` |
| `ErrMessageTooManyRequests` | `too_many_requests` |  | `result/constants.go:13` |
| `ErrMessageUnauthenticated` | `security_unauthenticated` |  | `security/api_errors.go:15` |
| `ErrMessageUnknown` | `unknown_error` |  | `result/constants.go:11` |
| `ErrMessageUnsupportedAuthenticationType` | `security_unsupported_authentication_type` |  | `security/api_errors.go:18` |
| `ErrMessageUnsupportedMediaType` | `unsupported_media_type` |  | `result/constants.go:15` |
| `ErrMessageUrgeTooFrequent` | `approval_urge_too_frequent` |  | `internal/approval/shared/messages.go:10` |
| `ErrMessageUserInfoLoaderNotImplemented` | `security_user_info_loader_not_implemented` |  | `security/api_errors.go:20` |
| `ErrMessageUserLoaderNotImplemented` | `security_user_loader_not_implemented` |  | `security/api_errors.go:19` |

## runtime enum value

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AcquireResultAcquired (AcquireResult)` | `acquired` |  | `event/inbox/inbox.go:29` |
| `AcquireResultCompleted (AcquireResult)` | `completed` |  | `event/inbox/inbox.go:32` |
| `AcquireResultInProgress (AcquireResult)` | `in_progress` |  | `event/inbox/inbox.go:36` |
| `ActionAddAssignee (ActionType)` | `add_assignee` |  | `approval/enums.go:254` |
| `ActionAddCC (ActionType)` | `add_cc` |  | `approval/enums.go:260` |
| `ActionApprove (ActionType)` | `approve` |  | `approval/enums.go:247` |
| `ActionCancel (ActionType)` | `cancel` |  | `approval/enums.go:252` |
| `ActionExecute (ActionType)` | `execute` |  | `approval/enums.go:256` |
| `ActionHandle (ActionType)` | `handle` |  | `approval/enums.go:248` |
| `ActionReassign (ActionType)` | `reassign` |  | `approval/enums.go:258` |
| `ActionReject (ActionType)` | `reject` |  | `approval/enums.go:249` |
| `ActionRemoveAssignee (ActionType)` | `remove_assignee` |  | `approval/enums.go:255` |
| `ActionResubmit (ActionType)` | `resubmit` |  | `approval/enums.go:257` |
| `ActionRollback (ActionType)` | `rollback` |  | `approval/enums.go:253` |
| `ActionSubmit (ActionType)` | `submit` |  | `approval/enums.go:246` |
| `ActionTerminate (ActionType)` | `terminate` |  | `approval/enums.go:259` |
| `ActionTransfer (ActionType)` | `transfer` |  | `approval/enums.go:250` |
| `ActionWithdraw (ActionType)` | `withdraw` |  | `approval/enums.go:251` |
| `AddAssigneeAfter (AddAssigneeType)` | `after` |  | `approval/enums.go:141` |
| `AddAssigneeBefore (AddAssigneeType)` | `before` |  | `approval/enums.go:140` |
| `AddAssigneeParallel (AddAssigneeType)` | `parallel` |  | `approval/enums.go:142` |
| `AesModeCbc (AESMode)` | `CBC` |  | `cryptox/aes_cipher.go:16` |
| `AesModeGcm (AESMode)` | `GCM` |  | `cryptox/aes_cipher.go:17` |
| `ApprovalParallel (ApprovalMethod)` | `parallel` |  | `approval/enums.go:81` |
| `ApprovalSequential (ApprovalMethod)` | `sequential` |  | `approval/enums.go:80` |
| `AssigneeDepartment (AssigneeKind)` | `department` |  | `approval/enums.go:182` |
| `AssigneeDepartmentLeader (AssigneeKind)` | `department_leader` |  | `approval/enums.go:185` |
| `AssigneeFormField (AssigneeKind)` | `form_field` |  | `approval/enums.go:186` |
| `AssigneeRole (AssigneeKind)` | `role` |  | `approval/enums.go:181` |
| `AssigneeSelf (AssigneeKind)` | `self` |  | `approval/enums.go:183` |
| `AssigneeSuperior (AssigneeKind)` | `superior` |  | `approval/enums.go:184` |
| `AssigneeUser (AssigneeKind)` | `user` |  | `approval/enums.go:180` |
| `Between (Operator)` | `between` |  | `search/constants.go:13` |
| `BindingBusiness (BindingMode)` | `business` |  | `approval/enums.go:15` |
| `BindingStandalone (BindingMode)` | `standalone` |  | `approval/enums.go:14` |
| `CCDepartment (CCKind)` | `department` |  | `approval/enums.go:269` |
| `CCFormField (CCKind)` | `form_field` |  | `approval/enums.go:270` |
| `CCRole (CCKind)` | `role` |  | `approval/enums.go:268` |
| `CCTimingAlways (CCTiming)` | `always` |  | `approval/enums.go:277` |
| `CCTimingOnApprove (CCTiming)` | `on_approve` |  | `approval/enums.go:278` |
| `CCTimingOnReject (CCTiming)` | `on_reject` |  | `approval/enums.go:279` |
| `CCUser (CCKind)` | `user` |  | `approval/enums.go:267` |
| `ChunkTypeError (ChunkType)` | `error` |  | `ai/stream/chunk.go:12` |
| `ChunkTypeFile (ChunkType)` | `file` |  | `ai/stream/chunk.go:30` |
| `ChunkTypeFinish (ChunkType)` | `finish` |  | `ai/stream/chunk.go:9` |
| `ChunkTypeFinishStep (ChunkType)` | `finish-step` |  | `ai/stream/chunk.go:11` |
| `ChunkTypeReasoningDelta (ChunkType)` | `reasoning-delta` |  | `ai/stream/chunk.go:19` |
| `ChunkTypeReasoningEnd (ChunkType)` | `reasoning-end` |  | `ai/stream/chunk.go:20` |
| `ChunkTypeReasoningStart (ChunkType)` | `reasoning-start` |  | `ai/stream/chunk.go:18` |
| `ChunkTypeSourceDocument (ChunkType)` | `source-document` |  | `ai/stream/chunk.go:28` |
| `ChunkTypeSourceURL (ChunkType)` | `source-url` |  | `ai/stream/chunk.go:27` |
| `ChunkTypeStart (ChunkType)` | `start` |  | `ai/stream/chunk.go:8` |
| `ChunkTypeStartStep (ChunkType)` | `start-step` |  | `ai/stream/chunk.go:10` |
| `ChunkTypeTextDelta (ChunkType)` | `text-delta` |  | `ai/stream/chunk.go:15` |
| `ChunkTypeTextEnd (ChunkType)` | `text-end` |  | `ai/stream/chunk.go:16` |
| `ChunkTypeTextStart (ChunkType)` | `text-start` |  | `ai/stream/chunk.go:14` |
| `ChunkTypeToolInputAvailable (ChunkType)` | `tool-input-available` |  | `ai/stream/chunk.go:24` |
| `ChunkTypeToolInputDelta (ChunkType)` | `tool-input-delta` |  | `ai/stream/chunk.go:23` |
| `ChunkTypeToolInputStart (ChunkType)` | `tool-input-start` |  | `ai/stream/chunk.go:22` |
| `ChunkTypeToolOutputAvailable (ChunkType)` | `tool-output-available` |  | `ai/stream/chunk.go:25` |
| `ClaimStatusPending (ClaimStatus)` | `pending` |  | `internal/storage/store/claim.go:21` |
| `ClaimStatusUploaded (ClaimStatus)` | `uploaded` |  | `internal/storage/store/claim.go:22` |
| `ConditionExpression (ConditionKind)` | `expression` |  | `approval/enums.go:239` |
| `ConditionField (ConditionKind)` | `field` |  | `approval/enums.go:238` |
| `ConsecutiveApproverAutoPass (ConsecutiveApproverAction)` | `auto_pass` |  | `approval/enums.go:173` |
| `ConsecutiveApproverNone (ConsecutiveApproverAction)` | `none` |  | `approval/enums.go:172` |
| `Contains (Operator)` | `contains` |  | `search/constants.go:22` |
| `ContainsIgnoreCase (Operator)` | `iContains` |  | `search/constants.go:29` |
| `DeleteReasonClaimExpired (DeleteReason)` | `claim_expired` |  | `storage/delete_enqueuer.go:26` |
| `DeleteReasonDeleted (DeleteReason)` | `deleted` |  | `storage/delete_enqueuer.go:21` |
| `DeleteReasonReplaced (DeleteReason)` | `replaced` |  | `storage/delete_enqueuer.go:18` |
| `EcdsaCurveP224 (ECDSACurve)` | `P224` |  | `cryptox/ecdsa_cipher.go:20` |
| `EcdsaCurveP256 (ECDSACurve)` | `P256` |  | `cryptox/ecdsa_cipher.go:21` |
| `EcdsaCurveP384 (ECDSACurve)` | `P384` |  | `cryptox/ecdsa_cipher.go:22` |
| `EcdsaCurveP521 (ECDSACurve)` | `P521` |  | `cryptox/ecdsa_cipher.go:23` |
| `EciesCurveP256 (ECIESCurve)` | `P256` |  | `cryptox/ecies_cipher.go:20` |
| `EciesCurveP384 (ECIESCurve)` | `P384` |  | `cryptox/ecies_cipher.go:21` |
| `EciesCurveP521 (ECIESCurve)` | `P521` |  | `cryptox/ecies_cipher.go:22` |
| `EciesCurveX25519 (ECIESCurve)` | `X25519` |  | `cryptox/ecies_cipher.go:23` |
| `EmptyAssigneeAutoPass (EmptyAssigneeAction)` | `auto_pass` |  | `approval/enums.go:98` |
| `EmptyAssigneeTransferAdmin (EmptyAssigneeAction)` | `transfer_admin` |  | `approval/enums.go:99` |
| `EmptyAssigneeTransferApplicant (EmptyAssigneeAction)` | `transfer_applicant` |  | `approval/enums.go:101` |
| `EmptyAssigneeTransferSpecified (EmptyAssigneeAction)` | `transfer_specified` |  | `approval/enums.go:102` |
| `EmptyAssigneeTransferSuperior (EmptyAssigneeAction)` | `transfer_superior` |  | `approval/enums.go:100` |
| `EncoderArgon2 (EncoderID)` | `argon2` |  | `password/password.go:8` |
| `EncoderBcrypt (EncoderID)` | `bcrypt` |  | `password/password.go:7` |
| `EncoderMd5 (EncoderID)` | `md5` |  | `password/password.go:11` |
| `EncoderPbkdf2 (EncoderID)` | `pbkdf2` |  | `password/password.go:10` |
| `EncoderPlaintext (EncoderID)` | `plaintext` |  | `password/password.go:13` |
| `EncoderScrypt (EncoderID)` | `scrypt` |  | `password/password.go:9` |
| `EncoderSha256 (EncoderID)` | `sha256` |  | `password/password.go:12` |
| `EndsWith (Operator)` | `endsWith` |  | `search/constants.go:26` |
| `EndsWithIgnoreCase (Operator)` | `iEndsWith` |  | `search/constants.go:33` |
| `Equals (Operator)` | `eq` |  | `search/constants.go:6` |
| `ExecutionAuto (ExecutionType)` | `auto` |  | `approval/enums.go:70` |
| `ExecutionAutoPass (ExecutionType)` | `auto_pass` |  | `approval/enums.go:71` |
| `ExecutionAutoReject (ExecutionType)` | `auto_reject` |  | `approval/enums.go:72` |
| `ExecutionManual (ExecutionType)` | `manual` |  | `approval/enums.go:69` |
| `FieldDate (FieldKind)` | `date` |  | `approval/enums.go:290` |
| `FieldInput (FieldKind)` | `input` |  | `approval/enums.go:286` |
| `FieldNumber (FieldKind)` | `number` |  | `approval/enums.go:289` |
| `FieldSelect (FieldKind)` | `select` |  | `approval/enums.go:288` |
| `FieldTextarea (FieldKind)` | `textarea` |  | `approval/enums.go:287` |
| `FieldUpload (FieldKind)` | `upload` |  | `approval/enums.go:291` |
| `FormatCsv (TabularFormat)` | `csv` |  | `crud/constants.go:11` |
| `FormatExcel (TabularFormat)` | `excel` |  | `crud/constants.go:10` |
| `FullPolicyBlock (FullPolicy)` | `block` |  | `event/transport/memory/memory.go:18` |
| `FullPolicyDropOldest (FullPolicy)` | `drop_oldest` |  | `event/transport/memory/memory.go:22` |
| `FullPolicyError (FullPolicy)` | `error` |  | `event/transport/memory/memory.go:15` |
| `GenderFemale (Gender)` | `female` |  | `security/user_info.go:18` |
| `GenderMale (Gender)` | `male` |  | `security/user_info.go:17` |
| `GenderUnknown (Gender)` | `unknown` |  | `security/user_info.go:19` |
| `GreaterThan (Operator)` | `gt` |  | `search/constants.go:8` |
| `GreaterThanOrEqual (Operator)` | `gte` |  | `search/constants.go:9` |
| `In (Operator)` | `in` |  | `search/constants.go:16` |
| `InitiatorDepartment (InitiatorKind)` | `department` |  | `approval/enums.go:33` |
| `InitiatorRole (InitiatorKind)` | `role` |  | `approval/enums.go:32` |
| `InitiatorUser (InitiatorKind)` | `user` |  | `approval/enums.go:31` |
| `InstanceApproved (InstanceStatus)` | `approved` |  | `approval/enums.go:194` |
| `InstanceRejected (InstanceStatus)` | `rejected` |  | `approval/enums.go:195` |
| `InstanceReturned (InstanceStatus)` | `returned` |  | `approval/enums.go:197` |
| `InstanceRunning (InstanceStatus)` | `running` |  | `approval/enums.go:193` |
| `InstanceTerminated (InstanceStatus)` | `terminated` |  | `approval/enums.go:198` |
| `InstanceWithdrawn (InstanceStatus)` | `withdrawn` |  | `approval/enums.go:196` |
| `IsNotNull (Operator)` | `isNotNull` |  | `search/constants.go:20` |
| `IsNull (Operator)` | `isNull` |  | `search/constants.go:19` |
| `LessThan (Operator)` | `lt` |  | `search/constants.go:10` |
| `LessThanOrEqual (Operator)` | `lte` |  | `search/constants.go:11` |
| `MetaTypeMarkdown (MetaType)` | `markdown` |  | `storage/file_refs.go:23` |
| `MetaTypeRichText (MetaType)` | `rich_text` |  | `storage/file_refs.go:21` |
| `MetaTypeUploadedFile (MetaType)` | `uploaded_file` |  | `storage/file_refs.go:19` |
| `MySQL (DBKind)` | `mysql` |  | `config/data_sources.go:17` |
| `NodeApproval (NodeKind)` | `approval` |  | `approval/enums.go:57` |
| `NodeCC (NodeKind)` | `cc` |  | `approval/enums.go:61` |
| `NodeCondition (NodeKind)` | `condition` |  | `approval/enums.go:59` |
| `NodeEnd (NodeKind)` | `end` |  | `approval/enums.go:60` |
| `NodeHandle (NodeKind)` | `handle` |  | `approval/enums.go:58` |
| `NodeStart (NodeKind)` | `start` |  | `approval/enums.go:56` |
| `NotBetween (Operator)` | `notBetween` |  | `search/constants.go:14` |
| `NotContains (Operator)` | `notContains` |  | `search/constants.go:23` |
| `NotContainsIgnoreCase (Operator)` | `iNotContains` |  | `search/constants.go:30` |
| `NotEndsWith (Operator)` | `notEndsWith` |  | `search/constants.go:27` |
| `NotEndsWithIgnoreCase (Operator)` | `iNotEndsWith` |  | `search/constants.go:34` |
| `NotEquals (Operator)` | `neq` |  | `search/constants.go:7` |
| `NotIn (Operator)` | `notIn` |  | `search/constants.go:17` |
| `NotStartsWith (Operator)` | `notStartsWith` |  | `search/constants.go:25` |
| `NotStartsWithIgnoreCase (Operator)` | `iNotStartsWith` |  | `search/constants.go:32` |
| `Oracle (DBKind)` | `oracle` |  | `config/data_sources.go:14` |
| `OverflowError (OverflowStrategy)` | `error` |  | `sequence/rule.go:22` |
| `OverflowExtend (OverflowStrategy)` | `extend` |  | `sequence/rule.go:27` |
| `OverflowReset (OverflowStrategy)` | `reset` |  | `sequence/rule.go:24` |
| `PassAll (PassRule)` | `all` |  | `approval/enums.go:88` |
| `PassAny (PassRule)` | `any` |  | `approval/enums.go:89` |
| `PassAnyReject (PassRule)` | `any_reject` |  | `approval/enums.go:91` |
| `PassRatio (PassRule)` | `ratio` |  | `approval/enums.go:90` |
| `PermissionEditable (Permission)` | `editable` |  | `approval/enums.go:310` |
| `PermissionHidden (Permission)` | `hidden` |  | `approval/enums.go:311` |
| `PermissionRequired (Permission)` | `required` |  | `approval/enums.go:312` |
| `PermissionVisible (Permission)` | `visible` |  | `approval/enums.go:309` |
| `Postgres (DBKind)` | `postgres` |  | `config/data_sources.go:16` |
| `PrincipalTypeExternalApp (PrincipalType)` | `external_app` |  | `security/principal.go:20` |
| `PrincipalTypeSystem (PrincipalType)` | `system` |  | `security/principal.go:22` |
| `PrincipalTypeUser (PrincipalType)` | `user` |  | `security/principal.go:18` |
| `ResetDaily (ResetCycle)` | `D` |  | `sequence/rule.go:10` |
| `ResetMonthly (ResetCycle)` | `M` |  | `sequence/rule.go:12` |
| `ResetNone (ResetCycle)` | `N` |  | `sequence/rule.go:9` |
| `ResetQuarterly (ResetCycle)` | `Q` |  | `sequence/rule.go:13` |
| `ResetWeekly (ResetCycle)` | `W` |  | `sequence/rule.go:11` |
| `ResetYearly (ResetCycle)` | `Y` |  | `sequence/rule.go:14` |
| `RoleAssistant (Role)` | `assistant` |  | `ai/stream/adapters.go:14` |
| `RoleAssistant (Role)` | `assistant` |  | `ai/message.go:12` |
| `RoleSystem (Role)` | `system` |  | `ai/stream/adapters.go:16` |
| `RoleSystem (Role)` | `system` |  | `ai/message.go:8` |
| `RoleTool (Role)` | `tool` |  | `ai/stream/adapters.go:15` |
| `RoleTool (Role)` | `tool` |  | `ai/message.go:14` |
| `RoleUser (Role)` | `user` |  | `ai/message.go:10` |
| `RoleUser (Role)` | `user` |  | `ai/stream/adapters.go:13` |
| `RollbackAny (RollbackType)` | `any` |  | `approval/enums.go:121` |
| `RollbackDataClear (RollbackDataStrategy)` | `clear` |  | `approval/enums.go:129` |
| `RollbackDataKeep (RollbackDataStrategy)` | `keep` |  | `approval/enums.go:130` |
| `RollbackNone (RollbackType)` | `none` |  | `approval/enums.go:118` |
| `RollbackPrevious (RollbackType)` | `previous` |  | `approval/enums.go:119` |
| `RollbackSpecified (RollbackType)` | `specified` |  | `approval/enums.go:122` |
| `RollbackStart (RollbackType)` | `start` |  | `approval/enums.go:120` |
| `RsaModeOAEP (RSAMode)` | `OAEP` |  | `cryptox/rsa_cipher.go:18` |
| `RsaModePKCS1v15 (RSAMode)` | `PKCS1v15` |  | `cryptox/rsa_cipher.go:19` |
| `RsaSignModePKCS1v15 (RSASignMode)` | `PKCS1v15` |  | `cryptox/rsa_cipher.go:26` |
| `RsaSignModePSS (RSASignMode)` | `PSS` |  | `cryptox/rsa_cipher.go:25` |
| `SM4ModeCBC (SM4Mode)` | `CBC` |  | `cryptox/sm4_cipher.go:15` |
| `SM4ModeECB (SM4Mode)` | `ECB` |  | `cryptox/sm4_cipher.go:16` |
| `SQLServer (DBKind)` | `sqlserver` |  | `config/data_sources.go:15` |
| `SQLite (DBKind)` | `sqlite` |  | `config/data_sources.go:18` |
| `SameApplicantAutoPass (SameApplicantAction)` | `auto_pass` |  | `approval/enums.go:109` |
| `SameApplicantSelfApprove (SameApplicantAction)` | `self_approve` |  | `approval/enums.go:110` |
| `SameApplicantTransferSuperior (SameApplicantAction)` | `transfer_superior` |  | `approval/enums.go:111` |
| `SignatureAlgHmacSHA256 (SignatureAlgorithm)` | `HMAC-SHA256` |  | `security/signature.go:31` |
| `SignatureAlgHmacSHA512 (SignatureAlgorithm)` | `HMAC-SHA512` |  | `security/signature.go:32` |
| `SignatureAlgHmacSM3 (SignatureAlgorithm)` | `HMAC-SM3` |  | `security/signature.go:33` |
| `StartsWith (Operator)` | `startsWith` |  | `search/constants.go:24` |
| `StartsWithIgnoreCase (Operator)` | `iStartsWith` |  | `search/constants.go:31` |
| `StatusCompleted (Status)` | `completed` |  | `event/inbox/inbox.go:20` |
| `StatusCompleted (Status)` | `completed` |  | `event/transport/outbox/outbox.go:25` |
| `StatusDead (Status)` | `dead` |  | `event/transport/outbox/outbox.go:32` |
| `StatusFailed (Status)` | `failed` |  | `event/transport/outbox/outbox.go:28` |
| `StatusPending (Status)` | `pending` |  | `event/transport/outbox/outbox.go:20` |
| `StatusProcessing (Status)` | `processing` |  | `event/transport/outbox/outbox.go:23` |
| `StatusProcessing (Status)` | `processing` |  | `event/inbox/inbox.go:16` |
| `StorageFilesystem (StorageProvider)` | `filesystem` |  | `config/storage.go:15` |
| `StorageJSON (StorageMode)` | `json` |  | `approval/enums.go:48` |
| `StorageMemory (StorageProvider)` | `memory` |  | `config/storage.go:14` |
| `StorageMinIO (StorageProvider)` | `minio` |  | `config/storage.go:13` |
| `TaskApproved (TaskStatus)` | `approved` |  | `approval/enums.go:212` |
| `TaskCanceled (TaskStatus)` | `canceled` |  | `approval/enums.go:217` |
| `TaskHandled (TaskStatus)` | `handled` |  | `approval/enums.go:214` |
| `TaskPending (TaskStatus)` | `pending` |  | `approval/enums.go:211` |
| `TaskRejected (TaskStatus)` | `rejected` |  | `approval/enums.go:213` |
| `TaskRemoved (TaskStatus)` | `removed` |  | `approval/enums.go:218` |
| `TaskRolledBack (TaskStatus)` | `rolled_back` |  | `approval/enums.go:216` |
| `TaskSkipped (TaskStatus)` | `skipped` |  | `approval/enums.go:219` |
| `TaskTransferred (TaskStatus)` | `transferred` |  | `approval/enums.go:215` |
| `TaskWaiting (TaskStatus)` | `waiting` |  | `approval/enums.go:210` |
| `TimeoutActionAutoPass (TimeoutAction)` | `auto_pass` |  | `approval/enums.go:299` |
| `TimeoutActionAutoReject (TimeoutAction)` | `auto_reject` |  | `approval/enums.go:300` |
| `TimeoutActionNone (TimeoutAction)` | `none` |  | `approval/enums.go:298` |
| `TimeoutActionNotify (TimeoutAction)` | `notify` |  | `approval/enums.go:301` |
| `TimeoutActionTransferAdmin (TimeoutAction)` | `transfer_admin` |  | `approval/enums.go:302` |
| `UserMenuTypeDashboard (UserMenuType)` | `dashboard` |  | `security/user_info.go:28` |
| `UserMenuTypeDirectory (UserMenuType)` | `directory` |  | `security/user_info.go:25` |
| `UserMenuTypeMenu (UserMenuType)` | `menu` |  | `security/user_info.go:26` |
| `UserMenuTypeReport (UserMenuType)` | `report` |  | `security/user_info.go:29` |
| `UserMenuTypeView (UserMenuType)` | `view` |  | `security/user_info.go:27` |
| `VersionArchived (VersionStatus)` | `archived` |  | `approval/enums.go:24` |
| `VersionDraft (VersionStatus)` | `draft` |  | `approval/enums.go:22` |
| `VersionPublished (VersionStatus)` | `published` |  | `approval/enums.go:23` |

## search tag grammar

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AttrAlias` | `alias` |  | `search/constants.go:39` |
| `AttrColumn` | `column` |  | `search/constants.go:40` |
| `AttrDive` | `dive` |  | `search/constants.go:38` |
| `AttrOperator` | `operator` |  | `search/constants.go:41` |
| `AttrParams` | `params` |  | `search/constants.go:42` |
| `IgnoreField` | `-` |  | `search/constants.go:47` |
| `ParamDelimiter` | `delimiter` |  | `search/constants.go:44` |
| `ParamType` | `type` |  | `search/constants.go:45` |
| `TypeBool` | `bool` |  | `search/constants.go:52` |
| `TypeDate` | `date` |  | `search/constants.go:54` |
| `TypeDateTime` | `datetime` |  | `search/constants.go:55` |
| `TypeDecimal` | `dec` |  | `search/constants.go:53` |
| `TypeInt` | `int` |  | `search/constants.go:50` |
| `TypeString` | `str` |  | `search/constants.go:51` |
| `TypeTime` | `time` |  | `search/constants.go:56` |
| `operator Between` | `between` |  | `search/constants.go:13` |
| `operator Contains` | `contains` |  | `search/constants.go:22` |
| `operator ContainsIgnoreCase` | `iContains` |  | `search/constants.go:29` |
| `operator EndsWith` | `endsWith` |  | `search/constants.go:26` |
| `operator EndsWithIgnoreCase` | `iEndsWith` |  | `search/constants.go:33` |
| `operator Equals` | `eq` |  | `search/constants.go:6` |
| `operator GreaterThan` | `gt` |  | `search/constants.go:8` |
| `operator GreaterThanOrEqual` | `gte` |  | `search/constants.go:9` |
| `operator In` | `in` |  | `search/constants.go:16` |
| `operator IsNotNull` | `isNotNull` |  | `search/constants.go:20` |
| `operator IsNull` | `isNull` |  | `search/constants.go:19` |
| `operator LessThan` | `lt` |  | `search/constants.go:10` |
| `operator LessThanOrEqual` | `lte` |  | `search/constants.go:11` |
| `operator NotBetween` | `notBetween` |  | `search/constants.go:14` |
| `operator NotContains` | `notContains` |  | `search/constants.go:23` |
| `operator NotContainsIgnoreCase` | `iNotContains` |  | `search/constants.go:30` |
| `operator NotEndsWith` | `notEndsWith` |  | `search/constants.go:27` |
| `operator NotEndsWithIgnoreCase` | `iNotEndsWith` |  | `search/constants.go:34` |
| `operator NotEquals` | `neq` |  | `search/constants.go:7` |
| `operator NotIn` | `notIn` |  | `search/constants.go:17` |
| `operator NotStartsWith` | `notStartsWith` |  | `search/constants.go:25` |
| `operator NotStartsWithIgnoreCase` | `iNotStartsWith` |  | `search/constants.go:32` |
| `operator StartsWith` | `startsWith` |  | `search/constants.go:24` |
| `operator StartsWithIgnoreCase` | `iStartsWith` |  | `search/constants.go:31` |
| `tag name` | `search` |  | `search/constants.go:36` |

## tabular tag grammar

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AttrDefault` | `default` |  | `tabular/constants.go:12` |
| `AttrDive` | `dive` |  | `tabular/constants.go:8` |
| `AttrFormat` | `format` |  | `tabular/constants.go:13` |
| `AttrFormatter` | `formatter` |  | `tabular/constants.go:14` |
| `AttrName` | `name` |  | `tabular/constants.go:9` |
| `AttrOrder` | `order` |  | `tabular/constants.go:11` |
| `AttrParser` | `parser` |  | `tabular/constants.go:15` |
| `AttrWidth` | `width` |  | `tabular/constants.go:10` |
| `IgnoreField` | `-` |  | `tabular/constants.go:18` |
| `tag name` | `tabular` |  | `tabular/constants.go:5` |

## validator label tag

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `label` | `label` |  | `validator/validator.go:64` |
| `label_i18n` | `label_i18n` |  | `validator/validator.go:69` |

## validator tag

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `alphanum_us` | `alphanum_us` |  | `validator/alphanum.go:16` |
| `alphanum_us_dot` | `alphanum_us_dot` |  | `validator/alphanum.go:24` |
| `alphanum_us_slash` | `alphanum_us_slash` |  | `validator/alphanum.go:20` |
| `dec_max` | `dec_max` |  | `validator/decimal.go:16` |
| `dec_min` | `dec_min` |  | `validator/decimal.go:10` |
| `phone_number` | `phone_number` |  | `validator/phone_number.go:8` |
