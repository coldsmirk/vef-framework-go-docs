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

Fingerprint: `a2d38a3dad64968b03893a0d3ade08317c16e707124b59b57561bbbb9debbb1f`
Entries: `3023`

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
| `HTTP header` | 6 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of api/header.go Header* constants. | None. |
| `HTTP wire field` | 8 | Tier 3 curated protocol fields | `extractProtocolConstants` | Curated source references for fundamental request/result fields shared by REST/RPC. | None in generated index; JSON DTO fields are covered separately. |
| `JSON wire field` | 1602 | Tier 2 scoped DTO AST with closed-world boundary check | `extractJSONFields` | AST scan of json tags on runtime DTO structs plus a boundary check over every non-test json-tagged struct field. | None for current non-test source; new json-tagged runtime fields must be indexed or explicitly excluded. |
| `MCP endpoint` | 1 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of the MCP Streamable HTTP endpoint constant. | None. |
| `MCP jsonschema tag` | 32 | Tier 2 pinned dependency parser catalog | `extractJSONSchemaTags` | Catalog of struct-tag keywords accepted by github.com/invopop/jsonschema v0.14.0, with boundary verification that fails on dependency-version drift and uncovered in-source jsonschema tags. | None for the pinned jsonschema parser version. |
| `MCP prompt` | 1 | Tier 2 MCP AST | `extractMCP` | AST scan of internal/mcp Prompt composites. | None in scanned MCP package. |
| `MCP tool` | 1 | Tier 2 MCP AST | `extractMCP` | AST scan of internal/mcp Tool composites. | None in scanned MCP package. |
| `REST action verb` | 10 | Tier 2 validator AST | `extractRESTVerbs` | AST scan of the REST action validator's allowed HTTP verb set. | None in current validator construction. |
| `RPC form key` | 2 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of FormKey* constants. | None. |
| `auth strategy` | 6 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of api/auth.go AuthStrategy* string constants. | None. |
| `auth type` | 5 | Tier 2 scoped AST constants | `extractAuthTypes` | AST scan of internal/security AuthType* constants that are sent through Authentication.Type. | None in known built-in authenticators. |
| `built-in resource` | 18 | Tier 2 scoped AST resources | `extractBuiltInResources` | AST scan of NewRPCResource/NewRESTResource calls in built-in runtime resource packages. | None in scanned built-in resource directories. |
| `built-in resource action` | 95 | Tier 2 scoped AST operations | `extractBuiltInResources` | AST scan of explicit OperationSpec values and CRUD builder defaults inside built-in runtime resource packages. | None in scanned built-in resource directories. |
| `config default` | 81 | Tier 3 mixed static extraction | `extractConfigDefaults` | AST extraction of Effective* accessors, ApplyDefaults assignments, monitor DefaultConfig values, and curated source references for defaults outside those named surfaces; boundary verification fails when a supported default surface is not indexed. | Defaults outside Effective*/ApplyDefaults/DefaultConfig and curated reviewed call sites require explicit review. |
| `config enum` | 8 | Tier 2 scoped AST constants | `extractProtocolConstants` | AST scan of storage and datasource enum constants used in configuration values. | None in current config enum files. |
| `config key` | 186 | Tier 2 config-tag AST | `extractConfigKeys` | AST walk of config structs rooted at known vef.* config roots plus vef.data_sources.&lt;name&gt;; verifier fails if a config/ struct with config tags is unreachable. | None for config/ structs with config tags. |
| `config reserved name` | 1 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of reserved configuration-name constants. | None. |
| `environment variable` | 5 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of Env* constants plus boundary checks for os.Getenv/os.LookupEnv call sites. | None for string-literal or const-backed environment lookups. |
| `event topic` | 35 | Tier 2 event constant/method scan | `extractProtocolConstants, extractMoldGrammar` | AST scan of EventType*/eventType* constants, EventType() return values, and built-in subscription/route-inspection topic call sites. | None for framework-owned non-test event topics. |
| `event transport contract` | 6 | Tier 2 scoped AST constants | `extractEventTransportContracts` | AST/source-derived extraction of outbox DLQ headers, topic prefix, retry backoff, and persisted-error bounds. | None for current built-in event transports. |
| `i18n key indirection` | 4 | Tier 2 AST call scan | `extractI18NMessageKeys` | AST scan of dynamic i18n.T call sites whose key source is another audited surface such as label_i18n tags, validator rules, or Fiber error mappings. | None for current dynamic i18n.T call sites. |
| `i18n message key` | 254 | Tier 2 AST call/tag scan | `extractI18NMessageKeys` | AST scan of literal or const-backed i18n.T calls, validator rule message keys, and label_i18n struct tags. | None for literal or const-backed keys; dynamic sources are tracked as i18n key indirections. |
| `meta tag grammar` | 7 | Tier 2 AST constants | `extractStructTagGrammars` | Catalog of storage meta tag name, dive value, file-reference kinds, and attribute grammar delimiters. | None for the current parser constants and tag parsing rules. |
| `mold tag grammar` | 9 | Tier 2 parser grammar scan | `extractMoldGrammar` | AST scan of the default mold tag name and restricted parser token constants, with boundary verification for parser token coverage. | None for current mold parser token constants. |
| `mold transformer tag` | 2 | Tier 2 transformer scan | `extractMoldGrammar` | AST scan of built-in FieldTransformer Tag() methods. | None for current built-in mold transformer Tag() methods. |
| `result error code` | 119 | Tier 1 AST constants | `extractErrorCodes` | AST scan of ErrCode* constants in api_errors.go and result/constants.go. | None for named error-code constants. |
| `result message key` | 43 | Tier 1 AST constants | `extractProtocolConstants` | AST scan of ErrMessage* constants. | Inline i18n keys are covered by the i18n message key category. |
| `runtime enum value` | 365 | Tier 2 typed string constants | `extractRuntimeEnumValues` | AST scan of typed string constants in public packages plus runtime internal DTO/transport packages. | Integer/stringer enum renderings are covered by the generated public API index and package contract ledger. |
| `search tag grammar` | 38 | Tier 1 AST constants | `extractStructTagGrammars` | AST scan of search tag name, attributes, params, ignore marker, and operator/type tokens. | None for constants in search/constants.go. |
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
| `VersionV1` | `v1` |  | `api/version.go:9` |
| `VersionV2` | `v2` |  | `api/version.go:10` |
| `VersionV3` | `v3` |  | `api/version.go:11` |
| `VersionV4` | `v4` |  | `api/version.go:12` |
| `VersionV5` | `v5` |  | `api/version.go:13` |
| `VersionV6` | `v6` |  | `api/version.go:14` |
| `VersionV7` | `v7` |  | `api/version.go:15` |
| `VersionV8` | `v8` |  | `api/version.go:16` |
| `VersionV9` | `v9` |  | `api/version.go:17` |

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
| `generate-model-schema --input` | `--input` | Input model file or directory path<br/>required<br/>short: -i | `cmd/vef-cli/cmd/modelschema/command.go:55` |
| `generate-model-schema --output` | `--output` | Output schema file or directory path<br/>required<br/>short: -o | `cmd/vef-cli/cmd/modelschema/command.go:56` |
| `generate-model-schema --package` | `--package` | Package name for generated schemas<br/>default: schemas<br/>short: -p | `cmd/vef-cli/cmd/modelschema/command.go:57` |

## CRUD REST action

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `RESTActionCreate` | `post /` |  | `crud/constants.go:47` |
| `RESTActionCreateMany` | `post /many` |  | `crud/constants.go:50` |
| `RESTActionDelete` | `delete /:id` |  | `crud/constants.go:49` |
| `RESTActionDeleteMany` | `delete /many` |  | `crud/constants.go:52` |
| `RESTActionExport` | `get /export` |  | `crud/constants.go:60` |
| `RESTActionFindAll` | `get /` |  | `crud/constants.go:54` |
| `RESTActionFindOne` | `get /:id` |  | `crud/constants.go:53` |
| `RESTActionFindOptions` | `get /options` |  | `crud/constants.go:56` |
| `RESTActionFindPage` | `get /page` |  | `crud/constants.go:55` |
| `RESTActionFindTree` | `get /tree` |  | `crud/constants.go:57` |
| `RESTActionFindTreeOptions` | `get /tree/options` |  | `crud/constants.go:58` |
| `RESTActionImport` | `post /import` |  | `crud/constants.go:59` |
| `RESTActionUpdate` | `put /:id` |  | `crud/constants.go:48` |
| `RESTActionUpdateMany` | `put /many` |  | `crud/constants.go:51` |

## CRUD RPC action

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `RPCActionCreate` | `create` |  | `crud/constants.go:29` |
| `RPCActionCreateMany` | `create_many` |  | `crud/constants.go:32` |
| `RPCActionDelete` | `delete` |  | `crud/constants.go:31` |
| `RPCActionDeleteMany` | `delete_many` |  | `crud/constants.go:34` |
| `RPCActionExport` | `export` |  | `crud/constants.go:42` |
| `RPCActionFindAll` | `find_all` |  | `crud/constants.go:36` |
| `RPCActionFindOne` | `find_one` |  | `crud/constants.go:35` |
| `RPCActionFindOptions` | `find_options` |  | `crud/constants.go:38` |
| `RPCActionFindPage` | `find_page` |  | `crud/constants.go:37` |
| `RPCActionFindTree` | `find_tree` |  | `crud/constants.go:39` |
| `RPCActionFindTreeOptions` | `find_tree_options` |  | `crud/constants.go:40` |
| `RPCActionImport` | `import` |  | `crud/constants.go:41` |
| `RPCActionUpdate` | `update` |  | `crud/constants.go:30` |
| `RPCActionUpdateMany` | `update_many` |  | `crud/constants.go:33` |

## HTTP endpoint

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `REST base path` | `/api` |  | `internal/api/router/rest.go:19` |
| `RPC endpoint` | `/api` | POST endpoint for RPC requests | `internal/api/router/rpc.go:19` |

## HTTP header

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `HeaderXAPIKey` | `X-API-Key` |  | `api/header.go:9` |
| `HeaderXAppID` | `X-App-ID` |  | `api/header.go:5` |
| `HeaderXMetaPrefix` | `X-Meta-` |  | `api/header.go:10` |
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
| `AbortUploadParams.ClaimID` | `claimId` | Go field: AbortUploadParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:750` |
| `ActionLog.Action` | `action` | Go field: ActionLog.Action<br/>type: string | `approval/admin/instance_detail.go:51` |
| `ActionLog.Action` | `action` | Go field: ActionLog.Action<br/>type: ActionType | `approval/models.go:401` |
| `ActionLog.AddAssigneeType` | `addAssigneeType` | Go field: ActionLog.AddAssigneeType<br/>type: *AddAssigneeType | `approval/models.go:414` |
| `ActionLog.AddedAssignees` | `addedAssignees` | Go field: ActionLog.AddedAssignees<br/>type: []approval.UserInfo | `approval/admin/instance_detail.go:57` |
| `ActionLog.AddedAssignees` | `addedAssignees` | Go field: ActionLog.AddedAssignees<br/>type: []UserInfo | `approval/models.go:415` |
| `ActionLog.Attachments` | `attachments` | Go field: ActionLog.Attachments<br/>type: []string | `approval/admin/instance_detail.go:61` |
| `ActionLog.Attachments` | `attachments` | Go field: ActionLog.Attachments<br/>type: []string | `approval/models.go:418` |
| `ActionLog.CCUsers` | `ccUsers` | Go field: ActionLog.CCUsers<br/>type: []approval.UserInfo | `approval/admin/instance_detail.go:59` |
| `ActionLog.CCUsers` | `ccUsers` | Go field: ActionLog.CCUsers<br/>type: []UserInfo | `approval/models.go:417` |
| `ActionLog.CreatedAt` | `createdAt` | Go field: ActionLog.CreatedAt<br/>type: timex.DateTime | `approval/admin/instance_detail.go:62` |
| `ActionLog.IPAddress` | `ipAddress` | Go field: ActionLog.IPAddress<br/>type: *string | `approval/models.go:406` |
| `ActionLog.InstanceID` | `instanceId` | Go field: ActionLog.InstanceID<br/>type: string | `approval/models.go:398` |
| `ActionLog.LogID` | `logId` | Go field: ActionLog.LogID<br/>type: string | `approval/admin/instance_detail.go:50` |
| `ActionLog.Meta` | `meta` | Go field: ActionLog.Meta<br/>type: map[string]any | `approval/models.go:419` |
| `ActionLog.NodeID` | `nodeId` | Go field: ActionLog.NodeID<br/>type: *string | `approval/models.go:399` |
| `ActionLog.NodeID` | `nodeId` | Go field: ActionLog.NodeID<br/>type: *string | `approval/admin/instance_detail.go:52` |
| `ActionLog.Operator` | `operator` | Go field: ActionLog.Operator<br/>type: approval.UserInfo | `approval/admin/instance_detail.go:54` |
| `ActionLog.OperatorDepartmentID` | `operatorDepartmentId` | Go field: ActionLog.OperatorDepartmentID<br/>type: *string | `approval/models.go:404` |
| `ActionLog.OperatorDepartmentName` | `operatorDepartmentName` | Go field: ActionLog.OperatorDepartmentName<br/>type: *string | `approval/models.go:405` |
| `ActionLog.OperatorID` | `operatorId` | Go field: ActionLog.OperatorID<br/>type: string | `approval/models.go:402` |
| `ActionLog.OperatorName` | `operatorName` | Go field: ActionLog.OperatorName<br/>type: string | `approval/models.go:403` |
| `ActionLog.Opinion` | `opinion` | Go field: ActionLog.Opinion<br/>type: *string | `approval/models.go:408` |
| `ActionLog.Opinion` | `opinion` | Go field: ActionLog.Opinion<br/>type: *string | `approval/admin/instance_detail.go:60` |
| `ActionLog.RemovedAssignees` | `removedAssignees` | Go field: ActionLog.RemovedAssignees<br/>type: []approval.UserInfo | `approval/admin/instance_detail.go:58` |
| `ActionLog.RemovedAssignees` | `removedAssignees` | Go field: ActionLog.RemovedAssignees<br/>type: []UserInfo | `approval/models.go:416` |
| `ActionLog.RollbackToNodeID` | `rollbackToNodeId` | Go field: ActionLog.RollbackToNodeID<br/>type: *string | `approval/models.go:413` |
| `ActionLog.RollbackToNodeID` | `rollbackToNodeId` | Go field: ActionLog.RollbackToNodeID<br/>type: *string | `approval/admin/instance_detail.go:56` |
| `ActionLog.TaskID` | `taskId` | Go field: ActionLog.TaskID<br/>type: *string | `approval/admin/instance_detail.go:53` |
| `ActionLog.TaskID` | `taskId` | Go field: ActionLog.TaskID<br/>type: *string | `approval/models.go:400` |
| `ActionLog.TransferTo` | `transferTo` | Go field: ActionLog.TransferTo<br/>type: *approval.UserInfo | `approval/admin/instance_detail.go:55` |
| `ActionLog.TransferToDepartmentID` | `transferToDepartmentId` | Go field: ActionLog.TransferToDepartmentID<br/>type: *string | `approval/models.go:411` |
| `ActionLog.TransferToDepartmentName` | `transferToDepartmentName` | Go field: ActionLog.TransferToDepartmentName<br/>type: *string | `approval/models.go:412` |
| `ActionLog.TransferToID` | `transferToId` | Go field: ActionLog.TransferToID<br/>type: *string | `approval/models.go:409` |
| `ActionLog.TransferToName` | `transferToName` | Go field: ActionLog.TransferToName<br/>type: *string | `approval/models.go:410` |
| `ActionLog.UserAgent` | `userAgent` | Go field: ActionLog.UserAgent<br/>type: *string | `approval/models.go:407` |
| `Activity.Action` | `action` | Go field: Activity.Action<br/>type: string | `approval/node_view.go:42` |
| `Activity.AddedAssignees` | `addedAssignees` | Go field: Activity.AddedAssignees<br/>type: []UserInfo | `approval/node_view.go:50` |
| `Activity.Attachments` | `attachments` | Go field: Activity.Attachments<br/>type: []string | `approval/node_view.go:45` |
| `Activity.CCUsers` | `ccUsers` | Go field: Activity.CCUsers<br/>type: []UserInfo | `approval/node_view.go:52` |
| `Activity.CreatedAt` | `createdAt` | Go field: Activity.CreatedAt<br/>type: timex.DateTime | `approval/node_view.go:53` |
| `Activity.Operator` | `operator` | Go field: Activity.Operator<br/>type: UserInfo | `approval/node_view.go:43` |
| `Activity.Opinion` | `opinion` | Go field: Activity.Opinion<br/>type: *string | `approval/node_view.go:44` |
| `Activity.RemovedAssignees` | `removedAssignees` | Go field: Activity.RemovedAssignees<br/>type: []UserInfo | `approval/node_view.go:51` |
| `Activity.RollbackToNodeID` | `rollbackToNodeId` | Go field: Activity.RollbackToNodeID<br/>type: *string | `approval/node_view.go:48` |
| `Activity.RollbackToNodeName` | `rollbackToNodeName` | Go field: Activity.RollbackToNodeName<br/>type: *string | `approval/node_view.go:49` |
| `Activity.Target` | `target` | Go field: Activity.Target<br/>type: *UserInfo | `approval/node_view.go:47` |
| `Activity.TransferTo` | `transferTo` | Go field: Activity.TransferTo<br/>type: *UserInfo | `approval/node_view.go:46` |
| `Adapter.ContractID` | `contractId` | Go field: Adapter.ContractID<br/>type: string | `integration/models.go:233` |
| `Adapter.Direction` | `direction` | Go field: Adapter.Direction<br/>type: Direction | `integration/models.go:236` |
| `Adapter.IsEnabled` | `isEnabled` | Go field: Adapter.IsEnabled<br/>type: bool | `integration/models.go:242` |
| `Adapter.Script` | `script` | Go field: Adapter.Script<br/>type: string | `integration/models.go:237` |
| `Adapter.SystemID` | `systemId` | Go field: Adapter.SystemID<br/>type: string | `integration/models.go:232` |
| `Adapter.TimeoutMs` | `timeoutMs` | Go field: Adapter.TimeoutMs<br/>type: int | `integration/models.go:241` |
| `AdapterParams.ContractID` | `contractId` | Go field: AdapterParams.ContractID<br/>type: string<br/>validate: "required" | `internal/integration/resource/adapter.go:19` |
| `AdapterParams.Direction` | `direction` | Go field: AdapterParams.Direction<br/>type: integration.Direction | `internal/integration/resource/adapter.go:20` |
| `AdapterParams.ID` | `id` | Go field: AdapterParams.ID<br/>type: string | `internal/integration/resource/adapter.go:17` |
| `AdapterParams.IsEnabled` | `isEnabled` | Go field: AdapterParams.IsEnabled<br/>type: bool | `internal/integration/resource/adapter.go:23` |
| `AdapterParams.Script` | `script` | Go field: AdapterParams.Script<br/>type: string<br/>validate: "required" | `internal/integration/resource/adapter.go:21` |
| `AdapterParams.SystemID` | `systemId` | Go field: AdapterParams.SystemID<br/>type: string<br/>validate: "required" | `internal/integration/resource/adapter.go:18` |
| `AdapterParams.TimeoutMs` | `timeoutMs` | Go field: AdapterParams.TimeoutMs<br/>type: int | `internal/integration/resource/adapter.go:22` |
| `AdapterSearch.ContractID` | `contractId` | Go field: AdapterSearch.ContractID<br/>search: "eq,column=contract_id"<br/>type: string | `internal/integration/resource/adapter.go:31` |
| `AdapterSearch.Direction` | `direction` | Go field: AdapterSearch.Direction<br/>search: "eq,column=direction"<br/>type: integration.Direction | `internal/integration/resource/adapter.go:32` |
| `AdapterSearch.IsEnabled` | `isEnabled` | Go field: AdapterSearch.IsEnabled<br/>search: "eq,column=is_enabled"<br/>type: *bool | `internal/integration/resource/adapter.go:33` |
| `AdapterSearch.SystemID` | `systemId` | Go field: AdapterSearch.SystemID<br/>search: "eq,column=system_id"<br/>type: string | `internal/integration/resource/adapter.go:30` |
| `AddAssigneeParams.AddType` | `addType` | Go field: AddAssigneeParams.AddType<br/>type: string<br/>validate: "required,oneof=before after parallel" | `internal/approval/resource/instance.go:402` |
| `AddAssigneeParams.TaskID` | `taskId` | Go field: AddAssigneeParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:400` |
| `AddAssigneeParams.UserIDs` | `userIds` | Go field: AddAssigneeParams.UserIDs<br/>type: []string<br/>validate: "required,min=1,max=50" | `internal/approval/resource/instance.go:401` |
| `AddCCParams.CCUserIDs` | `ccUserIds` | Go field: AddCCParams.CCUserIDs<br/>type: []string<br/>validate: "required,min=1,max=50" | `internal/approval/resource/instance.go:349` |
| `AddCCParams.InstanceID` | `instanceId` | Go field: AddCCParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:348` |
| `AdminFindActionLogsParams.InstanceID` | `instanceId` | Go field: AdminFindActionLogsParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:173` |
| `AdminFindActionLogsParams.Page` | `page` | Go field: AdminFindActionLogsParams.Page<br/>type: int | `internal/approval/resource/admin.go:175` |
| `AdminFindActionLogsParams.PageSize` | `pageSize` | Go field: AdminFindActionLogsParams.PageSize<br/>type: int | `internal/approval/resource/admin.go:176` |
| `AdminFindActionLogsParams.TenantID` | `tenantId` | Go field: AdminFindActionLogsParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:174` |
| `AdminFindBusinessProjectionsParams.Page` | `page` | Go field: AdminFindBusinessProjectionsParams.Page<br/>type: int | `internal/approval/resource/admin.go:289` |
| `AdminFindBusinessProjectionsParams.PageSize` | `pageSize` | Go field: AdminFindBusinessProjectionsParams.PageSize<br/>type: int | `internal/approval/resource/admin.go:290` |
| `AdminFindBusinessProjectionsParams.Status` | `status` | Go field: AdminFindBusinessProjectionsParams.Status<br/>type: *approval.BindingProjectionStatus | `internal/approval/resource/admin.go:288` |
| `AdminFindBusinessProjectionsParams.TenantID` | `tenantId` | Go field: AdminFindBusinessProjectionsParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:287` |
| `AdminFindInstancesParams.ApplicantID` | `applicantId` | Go field: AdminFindInstancesParams.ApplicantID<br/>type: *string | `internal/approval/resource/admin.go:62` |
| `AdminFindInstancesParams.FlowID` | `flowId` | Go field: AdminFindInstancesParams.FlowID<br/>type: *string | `internal/approval/resource/admin.go:64` |
| `AdminFindInstancesParams.Keyword` | `keyword` | Go field: AdminFindInstancesParams.Keyword<br/>type: *string | `internal/approval/resource/admin.go:65` |
| `AdminFindInstancesParams.Page` | `page` | Go field: AdminFindInstancesParams.Page<br/>type: int | `internal/approval/resource/admin.go:66` |
| `AdminFindInstancesParams.PageSize` | `pageSize` | Go field: AdminFindInstancesParams.PageSize<br/>type: int | `internal/approval/resource/admin.go:67` |
| `AdminFindInstancesParams.Status` | `status` | Go field: AdminFindInstancesParams.Status<br/>type: *approval.InstanceStatus | `internal/approval/resource/admin.go:63` |
| `AdminFindInstancesParams.TenantID` | `tenantId` | Go field: AdminFindInstancesParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:61` |
| `AdminFindTasksParams.AssigneeID` | `assigneeId` | Go field: AdminFindTasksParams.AssigneeID<br/>type: *string | `internal/approval/resource/admin.go:116` |
| `AdminFindTasksParams.InstanceID` | `instanceId` | Go field: AdminFindTasksParams.InstanceID<br/>type: *string | `internal/approval/resource/admin.go:117` |
| `AdminFindTasksParams.Page` | `page` | Go field: AdminFindTasksParams.Page<br/>type: int | `internal/approval/resource/admin.go:119` |
| `AdminFindTasksParams.PageSize` | `pageSize` | Go field: AdminFindTasksParams.PageSize<br/>type: int | `internal/approval/resource/admin.go:120` |
| `AdminFindTasksParams.Status` | `status` | Go field: AdminFindTasksParams.Status<br/>type: *approval.TaskStatus | `internal/approval/resource/admin.go:118` |
| `AdminFindTasksParams.TenantID` | `tenantId` | Go field: AdminFindTasksParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:115` |
| `AdminGetInstanceDetailParams.InstanceID` | `instanceId` | Go field: AdminGetInstanceDetailParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:148` |
| `AdminGetMetricsParams.TenantID` | `tenantId` | Go field: AdminGetMetricsParams.TenantID<br/>type: *string | `internal/approval/resource/admin.go:258` |
| `AdminReassignTaskParams.NewAssigneeID` | `newAssigneeId` | Go field: AdminReassignTaskParams.NewAssigneeID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:230` |
| `AdminReassignTaskParams.Reason` | `reason` | Go field: AdminReassignTaskParams.Reason<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/admin.go:231` |
| `AdminReassignTaskParams.TaskID` | `taskId` | Go field: AdminReassignTaskParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:229` |
| `AdminRetryBusinessProjectionParams.ProjectionID` | `projectionId` | Go field: AdminRetryBusinessProjectionParams.ProjectionID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:324` |
| `AdminTerminateInstanceParams.InstanceID` | `instanceId` | Go field: AdminTerminateInstanceParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/admin.go:202` |
| `AdminTerminateInstanceParams.Reason` | `reason` | Go field: AdminTerminateInstanceParams.Reason<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/admin.go:203` |
| `ApprovalNodeData.AddAssigneeTypes` | `addAssigneeTypes` | Go field: ApprovalNodeData.AddAssigneeTypes<br/>type: []AddAssigneeType | `approval/node_data.go:181` |
| `ApprovalNodeData.ApprovalMethod` | `approvalMethod` | Go field: ApprovalNodeData.ApprovalMethod<br/>type: ApprovalMethod | `approval/node_data.go:171` |
| `ApprovalNodeData.ConsecutiveApproverAction` | `consecutiveApproverAction` | Go field: ApprovalNodeData.ConsecutiveApproverAction<br/>type: ConsecutiveApproverAction | `approval/node_data.go:175` |
| `ApprovalNodeData.IsAddAssigneeAllowed` | `isAddAssigneeAllowed` | Go field: ApprovalNodeData.IsAddAssigneeAllowed<br/>type: *bool | `approval/node_data.go:180` |
| `ApprovalNodeData.IsManualCCAllowed` | `isManualCcAllowed` | Go field: ApprovalNodeData.IsManualCCAllowed<br/>type: *bool | `approval/node_data.go:183` |
| `ApprovalNodeData.IsRemoveAssigneeAllowed` | `isRemoveAssigneeAllowed` | Go field: ApprovalNodeData.IsRemoveAssigneeAllowed<br/>type: *bool | `approval/node_data.go:182` |
| `ApprovalNodeData.IsRollbackAllowed` | `isRollbackAllowed` | Go field: ApprovalNodeData.IsRollbackAllowed<br/>type: *bool | `approval/node_data.go:179` |
| `ApprovalNodeData.PassRatio` | `passRatio` | Go field: ApprovalNodeData.PassRatio<br/>type: decimal.Decimal | `approval/node_data.go:173` |
| `ApprovalNodeData.PassRule` | `passRule` | Go field: ApprovalNodeData.PassRule<br/>type: PassRule | `approval/node_data.go:172` |
| `ApprovalNodeData.RollbackDataStrategy` | `rollbackDataStrategy` | Go field: ApprovalNodeData.RollbackDataStrategy<br/>type: RollbackDataStrategy | `approval/node_data.go:177` |
| `ApprovalNodeData.RollbackTargetKeys` | `rollbackTargetKeys` | Go field: ApprovalNodeData.RollbackTargetKeys<br/>type: []string | `approval/node_data.go:178` |
| `ApprovalNodeData.RollbackType` | `rollbackType` | Go field: ApprovalNodeData.RollbackType<br/>type: RollbackType | `approval/node_data.go:176` |
| `ApprovalNodeData.SameApplicantAction` | `sameApplicantAction` | Go field: ApprovalNodeData.SameApplicantAction<br/>type: SameApplicantAction | `approval/node_data.go:174` |
| `AssigneeDefinition.FormField` | `formField` | Go field: AssigneeDefinition.FormField<br/>type: *string | `approval/assignee.go:75` |
| `AssigneeDefinition.IDs` | `ids` | Go field: AssigneeDefinition.IDs<br/>type: []string | `approval/assignee.go:74` |
| `AssigneeDefinition.Kind` | `kind` | Go field: AssigneeDefinition.Kind<br/>type: AssigneeKind | `approval/assignee.go:73` |
| `AssigneeDefinition.SortOrder` | `sortOrder` | Go field: AssigneeDefinition.SortOrder<br/>type: int | `approval/assignee.go:76` |
| `AssigneesAddedEvent.AddType` | `addType` | Go field: AssigneesAddedEvent.AddType<br/>type: AddAssigneeType | `approval/events_task.go:175` |
| `AssigneesAddedEvent.Assignees` | `assignees` | Go field: AssigneesAddedEvent.Assignees<br/>type: []UserInfo | `approval/events_task.go:176` |
| `AssigneesRemovedEvent.Assignees` | `assignees` | Go field: AssigneesRemovedEvent.Assignees<br/>type: []UserInfo | `approval/events_task.go:194` |
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
| `AuthTokens.AccessToken` | `accessToken` | Go field: AuthTokens.AccessToken<br/>type: string | `security/security.go:16` |
| `AuthTokens.RefreshToken` | `refreshToken` | Go field: AuthTokens.RefreshToken<br/>type: string | `security/security.go:17` |
| `Authentication.Credentials` | `credentials` | Go field: Authentication.Credentials<br/>type: any | `security/security.go:24` |
| `Authentication.Principal` | `principal` | Go field: Authentication.Principal<br/>type: string | `security/security.go:23` |
| `Authentication.Type` | `type` | Go field: Authentication.Type<br/>type: string | `security/security.go:22` |
| `AvailableFlow.CategoryID` | `categoryId` | Go field: AvailableFlow.CategoryID<br/>type: string | `approval/my/available_flows.go:11` |
| `AvailableFlow.CategoryName` | `categoryName` | Go field: AvailableFlow.CategoryName<br/>type: string | `approval/my/available_flows.go:12` |
| `AvailableFlow.Description` | `description` | Go field: AvailableFlow.Description<br/>type: *string | `approval/my/available_flows.go:9` |
| `AvailableFlow.FlowCode` | `flowCode` | Go field: AvailableFlow.FlowCode<br/>type: string | `approval/my/available_flows.go:6` |
| `AvailableFlow.FlowID` | `flowId` | Go field: AvailableFlow.FlowID<br/>type: string | `approval/my/available_flows.go:5` |
| `AvailableFlow.FlowIcon` | `flowIcon` | Go field: AvailableFlow.FlowIcon<br/>type: *string | `approval/my/available_flows.go:8` |
| `AvailableFlow.FlowName` | `flowName` | Go field: AvailableFlow.FlowName<br/>type: string | `approval/my/available_flows.go:7` |
| `AvailableFlow.Labels` | `labels` | Go field: AvailableFlow.Labels<br/>type: map[string]string | `approval/my/available_flows.go:10` |
| `BaseNodeData.Description` | `description` | Go field: BaseNodeData.Description<br/>type: *string | `approval/node_data.go:65` |
| `BaseNodeData.Name` | `name` | Go field: BaseNodeData.Name<br/>type: string | `approval/node_data.go:64` |
| `BuildInfo.AppVersion` | `appVersion` | Go field: BuildInfo.AppVersion<br/>type: string | `monitor/service.go:270` |
| `BuildInfo.BuildTime` | `buildTime` | Go field: BuildInfo.BuildTime<br/>type: string | `monitor/service.go:271` |
| `BuildInfo.GitCommit` | `gitCommit` | Go field: BuildInfo.GitCommit<br/>type: string | `monitor/service.go:272` |
| `BuildInfo.VEFVersion` | `vefVersion` | Go field: BuildInfo.VEFVersion<br/>type: string | `monitor/service.go:269` |
| `BusinessBindingConfig.FinishedAtColumn` | `finishedAtColumn` | Go field: BusinessBindingConfig.FinishedAtColumn<br/>type: *string | `approval/binding.go:47` |
| `BusinessBindingConfig.InstanceIDColumn` | `instanceIdColumn` | Go field: BusinessBindingConfig.InstanceIDColumn<br/>type: *string | `approval/binding.go:45` |
| `BusinessBindingConfig.KeyColumns` | `keyColumns` | Go field: BusinessBindingConfig.KeyColumns<br/>type: []string | `approval/binding.go:40` |
| `BusinessBindingConfig.StartedAtColumn` | `startedAtColumn` | Go field: BusinessBindingConfig.StartedAtColumn<br/>type: *string | `approval/binding.go:46` |
| `BusinessBindingConfig.StatusColumn` | `statusColumn` | Go field: BusinessBindingConfig.StatusColumn<br/>type: string | `approval/binding.go:41` |
| `BusinessBindingConfig.StatusMapping` | `statusMapping` | Go field: BusinessBindingConfig.StatusMapping<br/>type: map[InstanceStatus]string | `approval/binding.go:50` |
| `BusinessBindingConfig.TableName` | `tableName` | Go field: BusinessBindingConfig.TableName<br/>type: string | `approval/binding.go:39` |
| `BusinessProjection.AppliedAt` | `appliedAt` | Go field: BusinessProjection.AppliedAt<br/>type: *timex.DateTime | `approval/admin/business_projection.go:33` |
| `BusinessProjection.AppliedAt` | `appliedAt` | Go field: BusinessProjection.AppliedAt<br/>type: *timex.DateTime | `approval/models.go:286` |
| `BusinessProjection.AppliedOwnerInstanceID` | `appliedOwnerInstanceId` | Go field: BusinessProjection.AppliedOwnerInstanceID<br/>type: *string | `approval/models.go:271` |
| `BusinessProjection.AppliedOwnerInstanceID` | `appliedOwnerInstanceId` | Go field: BusinessProjection.AppliedOwnerInstanceID<br/>type: *string | `approval/admin/business_projection.go:19` |
| `BusinessProjection.AppliedRevision` | `appliedRevision` | Go field: BusinessProjection.AppliedRevision<br/>type: int64 | `approval/models.go:280` |
| `BusinessProjection.AppliedRevision` | `appliedRevision` | Go field: BusinessProjection.AppliedRevision<br/>type: int64 | `approval/admin/business_projection.go:27` |
| `BusinessProjection.AttemptCount` | `attemptCount` | Go field: BusinessProjection.AttemptCount<br/>type: int | `approval/admin/business_projection.go:29` |
| `BusinessProjection.AttemptCount` | `attemptCount` | Go field: BusinessProjection.AttemptCount<br/>type: int | `approval/models.go:282` |
| `BusinessProjection.Binding` | `binding` | Go field: BusinessProjection.Binding<br/>type: *BusinessBindingConfig | `approval/models.go:274` |
| `BusinessProjection.BusinessTable` | `businessTable` | Go field: BusinessProjection.BusinessTable<br/>type: string | `approval/admin/business_projection.go:20` |
| `BusinessProjection.Consistency` | `consistency` | Go field: BusinessProjection.Consistency<br/>type: config.ApprovalBindingConsistency | `approval/models.go:273` |
| `BusinessProjection.Consistency` | `consistency` | Go field: BusinessProjection.Consistency<br/>type: config.ApprovalBindingConsistency | `approval/admin/business_projection.go:22` |
| `BusinessProjection.DesiredFinishedAt` | `desiredFinishedAt` | Go field: BusinessProjection.DesiredFinishedAt<br/>type: *timex.DateTime | `approval/admin/business_projection.go:25` |
| `BusinessProjection.DesiredFinishedAt` | `desiredFinishedAt` | Go field: BusinessProjection.DesiredFinishedAt<br/>type: *timex.DateTime | `approval/models.go:278` |
| `BusinessProjection.DesiredRevision` | `desiredRevision` | Go field: BusinessProjection.DesiredRevision<br/>type: int64 | `approval/models.go:279` |
| `BusinessProjection.DesiredRevision` | `desiredRevision` | Go field: BusinessProjection.DesiredRevision<br/>type: int64 | `approval/admin/business_projection.go:26` |
| `BusinessProjection.DesiredStartedAt` | `desiredStartedAt` | Go field: BusinessProjection.DesiredStartedAt<br/>type: timex.DateTime | `approval/models.go:277` |
| `BusinessProjection.DesiredStartedAt` | `desiredStartedAt` | Go field: BusinessProjection.DesiredStartedAt<br/>type: timex.DateTime | `approval/admin/business_projection.go:24` |
| `BusinessProjection.DesiredStatus` | `desiredStatus` | Go field: BusinessProjection.DesiredStatus<br/>type: approval.InstanceStatus | `approval/admin/business_projection.go:23` |
| `BusinessProjection.DesiredStatus` | `desiredStatus` | Go field: BusinessProjection.DesiredStatus<br/>type: InstanceStatus | `approval/models.go:276` |
| `BusinessProjection.FlowID` | `flowId` | Go field: BusinessProjection.FlowID<br/>type: string | `approval/admin/business_projection.go:16` |
| `BusinessProjection.FlowID` | `flowId` | Go field: BusinessProjection.FlowID<br/>type: string | `approval/models.go:268` |
| `BusinessProjection.FlowVersionID` | `flowVersionId` | Go field: BusinessProjection.FlowVersionID<br/>type: string | `approval/admin/business_projection.go:17` |
| `BusinessProjection.FlowVersionID` | `flowVersionId` | Go field: BusinessProjection.FlowVersionID<br/>type: string | `approval/models.go:269` |
| `BusinessProjection.LastError` | `lastError` | Go field: BusinessProjection.LastError<br/>type: *string | `approval/models.go:285` |
| `BusinessProjection.LastError` | `lastError` | Go field: BusinessProjection.LastError<br/>type: *string | `approval/admin/business_projection.go:32` |
| `BusinessProjection.LeaseUntil` | `leaseUntil` | Go field: BusinessProjection.LeaseUntil<br/>type: *timex.DateTime | `approval/admin/business_projection.go:31` |
| `BusinessProjection.LeaseUntil` | `leaseUntil` | Go field: BusinessProjection.LeaseUntil<br/>type: *timex.DateTime | `approval/models.go:284` |
| `BusinessProjection.NextAttemptAt` | `nextAttemptAt` | Go field: BusinessProjection.NextAttemptAt<br/>type: *timex.DateTime | `approval/models.go:283` |
| `BusinessProjection.NextAttemptAt` | `nextAttemptAt` | Go field: BusinessProjection.NextAttemptAt<br/>type: *timex.DateTime | `approval/admin/business_projection.go:30` |
| `BusinessProjection.OwnerInstanceID` | `ownerInstanceId` | Go field: BusinessProjection.OwnerInstanceID<br/>type: string | `approval/models.go:270` |
| `BusinessProjection.OwnerInstanceID` | `ownerInstanceId` | Go field: BusinessProjection.OwnerInstanceID<br/>type: string | `approval/admin/business_projection.go:18` |
| `BusinessProjection.ProjectionID` | `projectionId` | Go field: BusinessProjection.ProjectionID<br/>type: string | `approval/admin/business_projection.go:14` |
| `BusinessProjection.RecordKey` | `recordKey` | Go field: BusinessProjection.RecordKey<br/>type: json.RawMessage | `approval/admin/business_projection.go:21` |
| `BusinessProjection.RecordKey` | `recordKey` | Go field: BusinessProjection.RecordKey<br/>type: json.RawMessage | `approval/models.go:275` |
| `BusinessProjection.Status` | `status` | Go field: BusinessProjection.Status<br/>type: BindingProjectionStatus | `approval/models.go:281` |
| `BusinessProjection.Status` | `status` | Go field: BusinessProjection.Status<br/>type: approval.BindingProjectionStatus | `approval/admin/business_projection.go:28` |
| `BusinessProjection.TargetHash` | `targetHash` | Go field: BusinessProjection.TargetHash<br/>type: string | `approval/models.go:272` |
| `BusinessProjection.TenantID` | `tenantId` | Go field: BusinessProjection.TenantID<br/>type: string | `approval/admin/business_projection.go:15` |
| `BusinessProjection.TenantID` | `tenantId` | Go field: BusinessProjection.TenantID<br/>type: string | `approval/models.go:267` |
| `BusinessProjection.UpdatedAt` | `updatedAt` | Go field: BusinessProjection.UpdatedAt<br/>type: timex.DateTime | `approval/admin/business_projection.go:34` |
| `CCDefinition.FormField` | `formField` | Go field: CCDefinition.FormField<br/>type: *string | `approval/assignee.go:83` |
| `CCDefinition.IDs` | `ids` | Go field: CCDefinition.IDs<br/>type: []string | `approval/assignee.go:82` |
| `CCDefinition.Kind` | `kind` | Go field: CCDefinition.Kind<br/>type: CCKind | `approval/assignee.go:81` |
| `CCDefinition.Timing` | `timing` | Go field: CCDefinition.Timing<br/>type: CCTiming | `approval/assignee.go:84` |
| `CCNodeData.CCs` | `ccs` | Go field: CCNodeData.CCs<br/>type: []CCDefinition | `approval/node_data.go:239` |
| `CCNodeData.FieldPermissions` | `fieldPermissions` | Go field: CCNodeData.FieldPermissions<br/>type: map[string]Permission | `approval/node_data.go:241` |
| `CCNodeData.IsReadConfirmRequired` | `isReadConfirmRequired` | Go field: CCNodeData.IsReadConfirmRequired<br/>type: bool | `approval/node_data.go:240` |
| `CCNotifiedEvent.IsManual` | `isManual` | Go field: CCNotifiedEvent.IsManual<br/>type: bool | `approval/events_cc.go:11` |
| `CCNotifiedEvent.NodeID` | `nodeId` | Go field: CCNotifiedEvent.NodeID<br/>type: string | `approval/events_cc.go:8` |
| `CCNotifiedEvent.NodeName` | `nodeName` | Go field: CCNotifiedEvent.NodeName<br/>type: string | `approval/events_cc.go:9` |
| `CCNotifiedEvent.Recipients` | `recipients` | Go field: CCNotifiedEvent.Recipients<br/>type: []UserInfo | `approval/events_cc.go:10` |
| `CCRecipient.ReadAt` | `readAt` | Go field: CCRecipient.ReadAt<br/>type: *timex.DateTime | `approval/node_view.go:61` |
| `CCRecipient.User` | `user` | Go field: CCRecipient.User<br/>type: UserInfo | `approval/node_view.go:60` |
| `CCRecord.Applicant` | `applicant` | Go field: CCRecord.Applicant<br/>type: approval.UserInfo | `approval/my/cc_records.go:16` |
| `CCRecord.CCRecordID` | `ccRecordId` | Go field: CCRecord.CCRecordID<br/>type: string | `approval/my/cc_records.go:10` |
| `CCRecord.CCUserDepartmentID` | `ccUserDepartmentId` | Go field: CCRecord.CCUserDepartmentID<br/>type: *string | `approval/models.go:468` |
| `CCRecord.CCUserDepartmentName` | `ccUserDepartmentName` | Go field: CCRecord.CCUserDepartmentName<br/>type: *string | `approval/models.go:469` |
| `CCRecord.CCUserID` | `ccUserId` | Go field: CCRecord.CCUserID<br/>type: string | `approval/models.go:466` |
| `CCRecord.CCUserName` | `ccUserName` | Go field: CCRecord.CCUserName<br/>type: string | `approval/models.go:467` |
| `CCRecord.CreatedAt` | `createdAt` | Go field: CCRecord.CreatedAt<br/>type: timex.DateTime | `approval/my/cc_records.go:19` |
| `CCRecord.FlowIcon` | `flowIcon` | Go field: CCRecord.FlowIcon<br/>type: *string | `approval/my/cc_records.go:15` |
| `CCRecord.FlowName` | `flowName` | Go field: CCRecord.FlowName<br/>type: string | `approval/my/cc_records.go:14` |
| `CCRecord.InstanceID` | `instanceId` | Go field: CCRecord.InstanceID<br/>type: string | `approval/models.go:458` |
| `CCRecord.InstanceID` | `instanceId` | Go field: CCRecord.InstanceID<br/>type: string | `approval/my/cc_records.go:11` |
| `CCRecord.InstanceNo` | `instanceNo` | Go field: CCRecord.InstanceNo<br/>type: string | `approval/my/cc_records.go:13` |
| `CCRecord.InstanceTitle` | `instanceTitle` | Go field: CCRecord.InstanceTitle<br/>type: string | `approval/my/cc_records.go:12` |
| `CCRecord.IsManual` | `isManual` | Go field: CCRecord.IsManual<br/>type: bool | `approval/models.go:470` |
| `CCRecord.IsRead` | `isRead` | Go field: CCRecord.IsRead<br/>type: bool | `approval/my/cc_records.go:18` |
| `CCRecord.NodeID` | `nodeId` | Go field: CCRecord.NodeID<br/>type: *string | `approval/models.go:459` |
| `CCRecord.NodeName` | `nodeName` | Go field: CCRecord.NodeName<br/>type: *string | `approval/my/cc_records.go:17` |
| `CCRecord.ReadAt` | `readAt` | Go field: CCRecord.ReadAt<br/>type: *timex.DateTime | `approval/models.go:471` |
| `CCRecord.TaskID` | `taskId` | Go field: CCRecord.TaskID<br/>type: *string | `approval/models.go:465` |
| `CCRecord.VisitID` | `visitId` | Go field: CCRecord.VisitID<br/>type: *string | `approval/models.go:464` |
| `CPUInfo.CacheSize` | `cacheSize` | Go field: CPUInfo.CacheSize<br/>type: int32 | `monitor/service.go:64` |
| `CPUInfo.EffectiveCores` | `effectiveCores` | Go field: CPUInfo.EffectiveCores<br/>type: float64 | `monitor/service.go:72` |
| `CPUInfo.Family` | `family` | Go field: CPUInfo.Family<br/>type: string | `monitor/service.go:68` |
| `CPUInfo.LogicalCores` | `logicalCores` | Go field: CPUInfo.LogicalCores<br/>type: int | `monitor/service.go:61` |
| `CPUInfo.Mhz` | `mhz` | Go field: CPUInfo.Mhz<br/>type: float64 | `monitor/service.go:63` |
| `CPUInfo.Microcode` | `microcode` | Go field: CPUInfo.Microcode<br/>type: string | `monitor/service.go:71` |
| `CPUInfo.Model` | `model` | Go field: CPUInfo.Model<br/>type: string | `monitor/service.go:69` |
| `CPUInfo.ModelName` | `modelName` | Go field: CPUInfo.ModelName<br/>type: string | `monitor/service.go:62` |
| `CPUInfo.PhysicalCores` | `physicalCores` | Go field: CPUInfo.PhysicalCores<br/>type: int | `monitor/service.go:60` |
| `CPUInfo.Stepping` | `stepping` | Go field: CPUInfo.Stepping<br/>type: int32 | `monitor/service.go:70` |
| `CPUInfo.TotalPercent` | `totalPercent` | Go field: CPUInfo.TotalPercent<br/>type: float64 | `monitor/service.go:66` |
| `CPUInfo.UsagePercent` | `usagePercent` | Go field: CPUInfo.UsagePercent<br/>type: []float64 | `monitor/service.go:65` |
| `CPUInfo.VendorID` | `vendorId` | Go field: CPUInfo.VendorID<br/>type: string | `monitor/service.go:67` |
| `CPUSummary.EffectiveCores` | `effectiveCores` | Go field: CPUSummary.EffectiveCores<br/>type: float64 | `monitor/service.go:53` |
| `CPUSummary.LogicalCores` | `logicalCores` | Go field: CPUSummary.LogicalCores<br/>type: int | `monitor/service.go:51` |
| `CPUSummary.PhysicalCores` | `physicalCores` | Go field: CPUSummary.PhysicalCores<br/>type: int | `monitor/service.go:50` |
| `CPUSummary.UsagePercent` | `usagePercent` | Go field: CPUSummary.UsagePercent<br/>type: float64 | `monitor/service.go:52` |
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
| `Check.Expr` | `expr` | Go field: Check.Expr<br/>type: string | `schema/service.go:69` |
| `Check.Name` | `name` | Go field: Check.Name<br/>type: string | `schema/service.go:68` |
| `CodeCatalog.Codes` | `codes` | Go field: CodeCatalog.Codes<br/>type: []mold.CodeInfo | `internal/integration/resource/code_set_catalog.go:29` |
| `CodeCatalog.Supported` | `supported` | Go field: CodeCatalog.Supported<br/>type: bool | `internal/integration/resource/code_set_catalog.go:28` |
| `CodeInfo.Code` | `code` | Go field: CodeInfo.Code<br/>type: string | `mold/translator.go:39` |
| `CodeInfo.Label` | `label` | Go field: CodeInfo.Label<br/>type: string | `mold/translator.go:41` |
| `CodeMap.CodeSet` | `codeSet` | Go field: CodeMap.CodeSet<br/>type: string | `integration/code_map.go:60` |
| `CodeMap.Entries` | `entries` | Go field: CodeMap.Entries<br/>type: []CodeMapEntry | `integration/code_map.go:64` |
| `CodeMap.FallbackCanonical` | `fallbackCanonical` | Go field: CodeMap.FallbackCanonical<br/>type: any | `integration/code_map.go:70` |
| `CodeMap.FallbackExternal` | `fallbackExternal` | Go field: CodeMap.FallbackExternal<br/>type: any | `integration/code_map.go:73` |
| `CodeMap.IsEnabled` | `isEnabled` | Go field: CodeMap.IsEnabled<br/>type: bool | `integration/code_map.go:74` |
| `CodeMap.Name` | `name` | Go field: CodeMap.Name<br/>type: string | `integration/code_map.go:61` |
| `CodeMap.OnUnmapped` | `onUnmapped` | Go field: CodeMap.OnUnmapped<br/>type: UnmappedPolicy | `integration/code_map.go:67` |
| `CodeMap.SystemID` | `systemId` | Go field: CodeMap.SystemID<br/>type: string | `integration/code_map.go:56` |
| `CodeMapEntry.Canonical` | `canonical` | Go field: CodeMapEntry.Canonical<br/>type: any | `integration/code_map.go:37` |
| `CodeMapEntry.CanonicalAliases` | `canonicalAliases` | Go field: CodeMapEntry.CanonicalAliases<br/>type: []any | `integration/code_map.go:41` |
| `CodeMapEntry.External` | `external` | Go field: CodeMapEntry.External<br/>type: any | `integration/code_map.go:39` |
| `CodeMapEntry.ExternalAliases` | `externalAliases` | Go field: CodeMapEntry.ExternalAliases<br/>type: []any | `integration/code_map.go:43` |
| `CodeMapParams.CodeSet` | `codeSet` | Go field: CodeMapParams.CodeSet<br/>type: string<br/>validate: "required" | `internal/integration/resource/code_map.go:24` |
| `CodeMapParams.Entries` | `entries` | Go field: CodeMapParams.Entries<br/>type: []integration.CodeMapEntry | `internal/integration/resource/code_map.go:26` |
| `CodeMapParams.FallbackCanonical` | `fallbackCanonical` | Go field: CodeMapParams.FallbackCanonical<br/>type: any | `internal/integration/resource/code_map.go:28` |
| `CodeMapParams.FallbackExternal` | `fallbackExternal` | Go field: CodeMapParams.FallbackExternal<br/>type: any | `internal/integration/resource/code_map.go:29` |
| `CodeMapParams.ID` | `id` | Go field: CodeMapParams.ID<br/>type: string | `internal/integration/resource/code_map.go:22` |
| `CodeMapParams.IsEnabled` | `isEnabled` | Go field: CodeMapParams.IsEnabled<br/>type: bool | `internal/integration/resource/code_map.go:30` |
| `CodeMapParams.Name` | `name` | Go field: CodeMapParams.Name<br/>type: string<br/>validate: "required" | `internal/integration/resource/code_map.go:25` |
| `CodeMapParams.OnUnmapped` | `onUnmapped` | Go field: CodeMapParams.OnUnmapped<br/>type: integration.UnmappedPolicy | `internal/integration/resource/code_map.go:27` |
| `CodeMapParams.SystemID` | `systemId` | Go field: CodeMapParams.SystemID<br/>type: string<br/>validate: "required" | `internal/integration/resource/code_map.go:23` |
| `CodeMapSearch.CodeSet` | `codeSet` | Go field: CodeMapSearch.CodeSet<br/>search: "contains,column=code_set"<br/>type: string | `internal/integration/resource/code_map.go:38` |
| `CodeMapSearch.IsEnabled` | `isEnabled` | Go field: CodeMapSearch.IsEnabled<br/>search: "eq,column=is_enabled"<br/>type: *bool | `internal/integration/resource/code_map.go:40` |
| `CodeMapSearch.Name` | `name` | Go field: CodeMapSearch.Name<br/>search: "contains"<br/>type: string | `internal/integration/resource/code_map.go:39` |
| `CodeMapSearch.SystemID` | `systemId` | Go field: CodeMapSearch.SystemID<br/>search: "eq,column=system_id"<br/>type: string | `internal/integration/resource/code_map.go:37` |
| `CodeSetCatalog.CodeSets` | `codeSets` | Go field: CodeSetCatalog.CodeSets<br/>type: []mold.CodeSetInfo | `internal/integration/resource/code_set_catalog.go:23` |
| `CodeSetCatalog.Supported` | `supported` | Go field: CodeSetCatalog.Supported<br/>type: bool | `internal/integration/resource/code_set_catalog.go:22` |
| `CodeSetChangedEvent.Keys` | `keys` | Go field: CodeSetChangedEvent.Keys<br/>type: []string | `mold/cached_code_set_resolver.go:27` |
| `CodeSetInfo.CodeSet` | `codeSet` | Go field: CodeSetInfo.CodeSet<br/>type: string | `mold/translator.go:31` |
| `CodeSetInfo.Name` | `name` | Go field: CodeSetInfo.Name<br/>type: string | `mold/translator.go:33` |
| `Column.Comment` | `comment` | Go field: Column.Comment<br/>type: string | `schema/service.go:31` |
| `Column.Default` | `default` | Go field: Column.Default<br/>type: string | `schema/service.go:30` |
| `Column.IsAutoIncrement` | `isAutoIncrement` | Go field: Column.IsAutoIncrement<br/>type: bool | `schema/service.go:33` |
| `Column.IsPrimaryKey` | `isPrimaryKey` | Go field: Column.IsPrimaryKey<br/>type: bool | `schema/service.go:32` |
| `Column.Name` | `name` | Go field: Column.Name<br/>type: string | `schema/service.go:27` |
| `Column.Nullable` | `nullable` | Go field: Column.Nullable<br/>type: bool | `schema/service.go:29` |
| `Column.Type` | `type` | Go field: Column.Type<br/>type: string | `schema/service.go:28` |
| `CompleteUploadParams.ClaimID` | `claimId` | Go field: CompleteUploadParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:579` |
| `CompleteUploadResult.OriginalFilename` | `originalFilename` | Go field: CompleteUploadResult.OriginalFilename<br/>type: string | `internal/storage/resource.go:590` |
| `CompletedTask.Applicant` | `applicant` | Go field: CompletedTask.Applicant<br/>type: approval.UserInfo | `approval/my/completed_tasks.go:16` |
| `CompletedTask.FinishedAt` | `finishedAt` | Go field: CompletedTask.FinishedAt<br/>type: *timex.DateTime | `approval/my/completed_tasks.go:19` |
| `CompletedTask.FlowIcon` | `flowIcon` | Go field: CompletedTask.FlowIcon<br/>type: *string | `approval/my/completed_tasks.go:15` |
| `CompletedTask.FlowName` | `flowName` | Go field: CompletedTask.FlowName<br/>type: string | `approval/my/completed_tasks.go:14` |
| `CompletedTask.InstanceID` | `instanceId` | Go field: CompletedTask.InstanceID<br/>type: string | `approval/my/completed_tasks.go:11` |
| `CompletedTask.InstanceNo` | `instanceNo` | Go field: CompletedTask.InstanceNo<br/>type: string | `approval/my/completed_tasks.go:13` |
| `CompletedTask.InstanceTitle` | `instanceTitle` | Go field: CompletedTask.InstanceTitle<br/>type: string | `approval/my/completed_tasks.go:12` |
| `CompletedTask.NodeName` | `nodeName` | Go field: CompletedTask.NodeName<br/>type: string | `approval/my/completed_tasks.go:17` |
| `CompletedTask.Status` | `status` | Go field: CompletedTask.Status<br/>type: string | `approval/my/completed_tasks.go:18` |
| `CompletedTask.TaskID` | `taskId` | Go field: CompletedTask.TaskID<br/>type: string | `approval/my/completed_tasks.go:10` |
| `Condition.Aggregate` | `aggregate` | Go field: Condition.Aggregate<br/>type: AggregateKind | `approval/condition.go:97` |
| `Condition.Column` | `column` | Go field: Condition.Column<br/>type: string | `approval/condition.go:98` |
| `Condition.Expression` | `expression` | Go field: Condition.Expression<br/>type: string | `approval/condition.go:101` |
| `Condition.Kind` | `kind` | Go field: Condition.Kind<br/>type: ConditionKind | `approval/condition.go:90` |
| `Condition.Operator` | `operator` | Go field: Condition.Operator<br/>type: ConditionOperator | `approval/condition.go:99` |
| `Condition.Subject` | `subject` | Go field: Condition.Subject<br/>type: string | `approval/condition.go:91` |
| `Condition.Value` | `value` | Go field: Condition.Value<br/>type: any | `approval/condition.go:100` |
| `ConditionBranch.ConditionGroups` | `conditionGroups` | Go field: ConditionBranch.ConditionGroups<br/>type: []ConditionGroup | `approval/condition.go:115` |
| `ConditionBranch.ID` | `id` | Go field: ConditionBranch.ID<br/>type: string | `approval/condition.go:113` |
| `ConditionBranch.IsDefault` | `isDefault` | Go field: ConditionBranch.IsDefault<br/>type: bool | `approval/condition.go:116` |
| `ConditionBranch.Label` | `label` | Go field: ConditionBranch.Label<br/>type: string | `approval/condition.go:114` |
| `ConditionBranch.Priority` | `priority` | Go field: ConditionBranch.Priority<br/>type: int | `approval/condition.go:117` |
| `ConditionGroup.Conditions` | `conditions` | Go field: ConditionGroup.Conditions<br/>type: []Condition | `approval/condition.go:107` |
| `ConditionNodeData.Branches` | `branches` | Go field: ConditionNodeData.Branches<br/>type: []ConditionBranch | `approval/node_data.go:265` |
| `ConnectionCheck.Database` | `database` | Go field: ConnectionCheck.Database<br/>type: *DatabaseProbe | `internal/integration/exec/connection_check.go:17` |
| `ConnectionCheck.HTTP` | `http` | Go field: ConnectionCheck.HTTP<br/>type: *HTTPProbe | `internal/integration/exec/connection_check.go:16` |
| `Contract.Code` | `code` | Go field: Contract.Code<br/>type: string | `integration/models.go:17` |
| `Contract.Description` | `description` | Go field: Contract.Description<br/>type: *string | `integration/models.go:19` |
| `Contract.InputSchema` | `inputSchema` | Go field: Contract.InputSchema<br/>type: json.RawMessage | `integration/models.go:27` |
| `Contract.IsEnabled` | `isEnabled` | Go field: Contract.IsEnabled<br/>type: bool | `integration/models.go:31` |
| `Contract.Labels` | `labels` | Go field: Contract.Labels<br/>type: map[string]string | `integration/models.go:24` |
| `Contract.Name` | `name` | Go field: Contract.Name<br/>type: string | `integration/models.go:18` |
| `Contract.OutputSchema` | `outputSchema` | Go field: Contract.OutputSchema<br/>type: json.RawMessage | `integration/models.go:30` |
| `ContractParams.Code` | `code` | Go field: ContractParams.Code<br/>type: string<br/>validate: "required" | `internal/integration/resource/contract.go:21` |
| `ContractParams.Description` | `description` | Go field: ContractParams.Description<br/>type: *string | `internal/integration/resource/contract.go:23` |
| `ContractParams.ID` | `id` | Go field: ContractParams.ID<br/>type: string | `internal/integration/resource/contract.go:20` |
| `ContractParams.InputSchema` | `inputSchema` | Go field: ContractParams.InputSchema<br/>type: json.RawMessage | `internal/integration/resource/contract.go:25` |
| `ContractParams.IsEnabled` | `isEnabled` | Go field: ContractParams.IsEnabled<br/>type: bool | `internal/integration/resource/contract.go:27` |
| `ContractParams.Labels` | `labels` | Go field: ContractParams.Labels<br/>type: map[string]string | `internal/integration/resource/contract.go:24` |
| `ContractParams.Name` | `name` | Go field: ContractParams.Name<br/>type: string<br/>validate: "required" | `internal/integration/resource/contract.go:22` |
| `ContractParams.OutputSchema` | `outputSchema` | Go field: ContractParams.OutputSchema<br/>type: json.RawMessage | `internal/integration/resource/contract.go:26` |
| `ContractSearch.Code` | `code` | Go field: ContractSearch.Code<br/>search: "contains"<br/>type: string | `internal/integration/resource/contract.go:34` |
| `ContractSearch.IsEnabled` | `isEnabled` | Go field: ContractSearch.IsEnabled<br/>search: "eq,column=is_enabled"<br/>type: *bool | `internal/integration/resource/contract.go:36` |
| `ContractSearch.Labels` | `labels` | Go field: ContractSearch.Labels<br/>type: map[string]string | `internal/integration/resource/contract.go:40` |
| `ContractSearch.Name` | `name` | Go field: ContractSearch.Name<br/>search: "contains"<br/>type: string | `internal/integration/resource/contract.go:35` |
| `CreateFlowParams.AdminUserIDs` | `adminUserIds` | Go field: CreateFlowParams.AdminUserIDs<br/>type: []string | `internal/approval/resource/flow.go:64` |
| `CreateFlowParams.BindingMode` | `bindingMode` | Go field: CreateFlowParams.BindingMode<br/>type: approval.BindingMode<br/>validate: "required" | `internal/approval/resource/flow.go:62` |
| `CreateFlowParams.BusinessBinding` | `businessBinding` | Go field: CreateFlowParams.BusinessBinding<br/>type: *approval.BusinessBindingConfig | `internal/approval/resource/flow.go:63` |
| `CreateFlowParams.CategoryID` | `categoryId` | Go field: CreateFlowParams.CategoryID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:58` |
| `CreateFlowParams.Code` | `code` | Go field: CreateFlowParams.Code<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:56` |
| `CreateFlowParams.Description` | `description` | Go field: CreateFlowParams.Description<br/>type: *string | `internal/approval/resource/flow.go:60` |
| `CreateFlowParams.Icon` | `icon` | Go field: CreateFlowParams.Icon<br/>type: *string | `internal/approval/resource/flow.go:59` |
| `CreateFlowParams.Initiators` | `initiators` | Go field: CreateFlowParams.Initiators<br/>type: []CreateInitiatorParams | `internal/approval/resource/flow.go:67` |
| `CreateFlowParams.InstanceTitleTemplate` | `instanceTitleTemplate` | Go field: CreateFlowParams.InstanceTitleTemplate<br/>type: string | `internal/approval/resource/flow.go:66` |
| `CreateFlowParams.IsAllInitiationAllowed` | `isAllInitiationAllowed` | Go field: CreateFlowParams.IsAllInitiationAllowed<br/>type: bool | `internal/approval/resource/flow.go:65` |
| `CreateFlowParams.Labels` | `labels` | Go field: CreateFlowParams.Labels<br/>type: map[string]string | `internal/approval/resource/flow.go:61` |
| `CreateFlowParams.Name` | `name` | Go field: CreateFlowParams.Name<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:57` |
| `CreateFlowParams.TenantID` | `tenantId` | Go field: CreateFlowParams.TenantID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:55` |
| `CreateInitiatorParams.IDs` | `ids` | Go field: CreateInitiatorParams.IDs<br/>type: []string<br/>validate: "required" | `internal/approval/resource/flow.go:73` |
| `CreateInitiatorParams.Kind` | `kind` | Go field: CreateInitiatorParams.Kind<br/>type: approval.InitiatorKind<br/>validate: "required" | `internal/approval/resource/flow.go:72` |
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
| `DataSourceConfig.Database` | `database` | Go field: DataSourceConfig.Database<br/>type: string | `integration/models.go:108` |
| `DataSourceConfig.Host` | `host` | Go field: DataSourceConfig.Host<br/>type: string | `integration/models.go:104` |
| `DataSourceConfig.Kind` | `kind` | Go field: DataSourceConfig.Kind<br/>type: config.DBKind | `integration/models.go:101` |
| `DataSourceConfig.Mode` | `mode` | Go field: DataSourceConfig.Mode<br/>type: DataSourceMode | `integration/models.go:103` |
| `DataSourceConfig.Password` | `password` | Go field: DataSourceConfig.Password<br/>type: string | `integration/models.go:107` |
| `DataSourceConfig.Path` | `path` | Go field: DataSourceConfig.Path<br/>type: string | `integration/models.go:110` |
| `DataSourceConfig.Port` | `port` | Go field: DataSourceConfig.Port<br/>type: uint16 | `integration/models.go:105` |
| `DataSourceConfig.SSLMode` | `sslMode` | Go field: DataSourceConfig.SSLMode<br/>type: config.SSLMode | `integration/models.go:111` |
| `DataSourceConfig.SSLRootCert` | `sslRootCert` | Go field: DataSourceConfig.SSLRootCert<br/>type: string | `integration/models.go:112` |
| `DataSourceConfig.Schema` | `schema` | Go field: DataSourceConfig.Schema<br/>type: string | `integration/models.go:109` |
| `DataSourceConfig.User` | `user` | Go field: DataSourceConfig.User<br/>type: string | `integration/models.go:106` |
| `DatabaseProbe.DurationMs` | `durationMs` | Go field: DatabaseProbe.DurationMs<br/>type: int64 | `internal/integration/exec/connection_check.go:34` |
| `DatabaseProbe.Error` | `error` | Go field: DatabaseProbe.Error<br/>type: string | `internal/integration/exec/connection_check.go:35` |
| `DatabaseProbe.Reachable` | `reachable` | Go field: DatabaseProbe.Reachable<br/>type: bool | `internal/integration/exec/connection_check.go:32` |
| `DatabaseProbe.Version` | `version` | Go field: DatabaseProbe.Version<br/>type: string | `internal/integration/exec/connection_check.go:33` |
| `Delegation.DelegateeID` | `delegateeId` | Go field: Delegation.DelegateeID<br/>type: string | `approval/models.go:494` |
| `Delegation.DelegatorID` | `delegatorId` | Go field: Delegation.DelegatorID<br/>type: string | `approval/models.go:493` |
| `Delegation.EndTime` | `endTime` | Go field: Delegation.EndTime<br/>type: timex.DateTime | `approval/models.go:498` |
| `Delegation.FlowCategoryID` | `flowCategoryId` | Go field: Delegation.FlowCategoryID<br/>type: *string | `approval/models.go:495` |
| `Delegation.FlowID` | `flowId` | Go field: Delegation.FlowID<br/>type: *string | `approval/models.go:496` |
| `Delegation.IsActive` | `isActive` | Go field: Delegation.IsActive<br/>type: bool | `approval/models.go:499` |
| `Delegation.Reason` | `reason` | Go field: Delegation.Reason<br/>type: *string | `approval/models.go:500` |
| `Delegation.StartTime` | `startTime` | Go field: Delegation.StartTime<br/>type: timex.DateTime | `approval/models.go:497` |
| `DelegationParams.DelegateeID` | `delegateeId` | Go field: DelegationParams.DelegateeID<br/>type: string<br/>validate: "required" | `internal/approval/resource/delegation.go:20` |
| `DelegationParams.DelegatorID` | `delegatorId` | Go field: DelegationParams.DelegatorID<br/>type: string<br/>validate: "required" | `internal/approval/resource/delegation.go:19` |
| `DelegationParams.EndTime` | `endTime` | Go field: DelegationParams.EndTime<br/>type: *timex.DateTime<br/>validate: "required" | `internal/approval/resource/delegation.go:24` |
| `DelegationParams.FlowCategoryID` | `flowCategoryId` | Go field: DelegationParams.FlowCategoryID<br/>type: *string | `internal/approval/resource/delegation.go:21` |
| `DelegationParams.FlowID` | `flowId` | Go field: DelegationParams.FlowID<br/>type: *string | `internal/approval/resource/delegation.go:22` |
| `DelegationParams.ID` | `id` | Go field: DelegationParams.ID<br/>type: string | `internal/approval/resource/delegation.go:18` |
| `DelegationParams.IsActive` | `isActive` | Go field: DelegationParams.IsActive<br/>type: bool | `internal/approval/resource/delegation.go:25` |
| `DelegationParams.Reason` | `reason` | Go field: DelegationParams.Reason<br/>type: *string | `internal/approval/resource/delegation.go:26` |
| `DelegationParams.StartTime` | `startTime` | Go field: DelegationParams.StartTime<br/>type: *timex.DateTime<br/>validate: "required" | `internal/approval/resource/delegation.go:23` |
| `DelegationSearch.DelegateeID` | `delegateeId` | Go field: DelegationSearch.DelegateeID<br/>search: "eq,column=delegatee_id"<br/>type: string | `internal/approval/resource/delegation.go:34` |
| `DelegationSearch.DelegatorID` | `delegatorId` | Go field: DelegationSearch.DelegatorID<br/>search: "eq,column=delegator_id"<br/>type: string | `internal/approval/resource/delegation.go:33` |
| `DelegationSearch.IsActive` | `isActive` | Go field: DelegationSearch.IsActive<br/>search: "eq,column=is_active"<br/>type: *bool | `internal/approval/resource/delegation.go:35` |
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
| `DeployFlowParams.Description` | `description` | Go field: DeployFlowParams.Description<br/>type: *string | `internal/approval/resource/flow.go:125` |
| `DeployFlowParams.FlowDefinition` | `flowDefinition` | Go field: DeployFlowParams.FlowDefinition<br/>type: approval.FlowDefinition<br/>validate: "required" | `internal/approval/resource/flow.go:127` |
| `DeployFlowParams.FlowID` | `flowId` | Go field: DeployFlowParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:124` |
| `DeployFlowParams.FormSchema` | `formSchema` | Go field: DeployFlowParams.FormSchema<br/>type: json.RawMessage | `internal/approval/resource/flow.go:128` |
| `DeployFlowParams.StorageMode` | `storageMode` | Go field: DeployFlowParams.StorageMode<br/>type: approval.StorageMode | `internal/approval/resource/flow.go:126` |
| `DiskInfo.IOCounters` | `ioCounters` | Go field: DiskInfo.IOCounters<br/>type: map[string]*IOCounter | `monitor/service.go:155` |
| `DiskInfo.Partitions` | `partitions` | Go field: DiskInfo.Partitions<br/>type: []*PartitionInfo | `monitor/service.go:154` |
| `DiskSummary.Partitions` | `partitions` | Go field: DiskSummary.Partitions<br/>type: int | `monitor/service.go:149` |
| `DiskSummary.Total` | `total` | Go field: DiskSummary.Total<br/>type: uint64 | `monitor/service.go:146` |
| `DiskSummary.Used` | `used` | Go field: DiskSummary.Used<br/>type: uint64 | `monitor/service.go:147` |
| `DiskSummary.UsedPercent` | `usedPercent` | Go field: DiskSummary.UsedPercent<br/>type: float64 | `monitor/service.go:148` |
| `DryRunInboundParams.ContractCode` | `contractCode` | Go field: DryRunInboundParams.ContractCode<br/>type: string<br/>validate: "required" | `internal/integration/resource/ops.go:38` |
| `DryRunInboundParams.HandlerOutput` | `handlerOutput` | Go field: DryRunInboundParams.HandlerOutput<br/>type: json.RawMessage | `internal/integration/resource/ops.go:41` |
| `DryRunInboundParams.Request` | `request` | Go field: DryRunInboundParams.Request<br/>type: InboundRequestParams | `internal/integration/resource/ops.go:40` |
| `DryRunInboundParams.Script` | `script` | Go field: DryRunInboundParams.Script<br/>type: string | `internal/integration/resource/ops.go:39` |
| `DryRunInboundParams.SystemCode` | `systemCode` | Go field: DryRunInboundParams.SystemCode<br/>type: string<br/>validate: "required" | `internal/integration/resource/ops.go:37` |
| `DryRunParams.ContractCode` | `contractCode` | Go field: DryRunParams.ContractCode<br/>type: string<br/>validate: "required" | `internal/integration/resource/ops.go:25` |
| `DryRunParams.Input` | `input` | Go field: DryRunParams.Input<br/>type: json.RawMessage | `internal/integration/resource/ops.go:27` |
| `DryRunParams.Script` | `script` | Go field: DryRunParams.Script<br/>type: string | `internal/integration/resource/ops.go:26` |
| `DryRunParams.SystemCode` | `systemCode` | Go field: DryRunParams.SystemCode<br/>type: string<br/>validate: "required" | `internal/integration/resource/ops.go:24` |
| `DryRunResult.Error` | `error` | Go field: DryRunResult.Error<br/>type: string | `internal/integration/exec/invoker.go:169` |
| `DryRunResult.FailureKind` | `failureKind` | Go field: DryRunResult.FailureKind<br/>type: integration.FailureKind | `internal/integration/exec/invoker.go:168` |
| `DryRunResult.Output` | `output` | Go field: DryRunResult.Output<br/>type: any | `internal/integration/exec/invoker.go:166` |
| `DryRunResult.Trace` | `trace` | Go field: DryRunResult.Trace<br/>type: []integration.HTTPExchange | `internal/integration/exec/invoker.go:167` |
| `EdgeDefinition.Data` | `data` | Go field: EdgeDefinition.Data<br/>type: map[string]any | `approval/flow_definition.go:80` |
| `EdgeDefinition.ID` | `id` | Go field: EdgeDefinition.ID<br/>type: string | `approval/flow_definition.go:76` |
| `EdgeDefinition.Source` | `source` | Go field: EdgeDefinition.Source<br/>type: string | `approval/flow_definition.go:77` |
| `EdgeDefinition.SourceHandle` | `sourceHandle` | Go field: EdgeDefinition.SourceHandle<br/>type: *string | `approval/flow_definition.go:79` |
| `EdgeDefinition.Target` | `target` | Go field: EdgeDefinition.Target<br/>type: string | `approval/flow_definition.go:78` |
| `EventStreamsInfo.Enabled` | `enabled` | Go field: EventStreamsInfo.Enabled<br/>type: bool | `monitor/event_streams.go:12` |
| `EventStreamsInfo.Streams` | `streams` | Go field: EventStreamsInfo.Streams<br/>type: []event.StreamInfo | `monitor/event_streams.go:13` |
| `ExternalAppConfig.Enabled` | `enabled` | Go field: ExternalAppConfig.Enabled<br/>type: bool | `security/security.go:29` |
| `ExternalAppConfig.IPWhitelist` | `ipWhitelist` | Go field: ExternalAppConfig.IPWhitelist<br/>type: string | `security/security.go:30` |
| `FieldOption.Label` | `label` | Go field: FieldOption.Label<br/>type: string | `approval/form_field.go:43` |
| `FieldOption.Value` | `value` | Go field: FieldOption.Value<br/>type: any | `approval/form_field.go:44` |
| `FileClaimedEvent.FileKey` | `fileKey` | Go field: FileClaimedEvent.FileKey<br/>type: string | `storage/events.go:35` |
| `FileDeletedEvent.FileKey` | `fileKey` | Go field: FileDeletedEvent.FileKey<br/>type: string | `storage/events.go:51` |
| `FileDeletedEvent.Reason` | `reason` | Go field: FileDeletedEvent.Reason<br/>type: DeleteReason | `storage/events.go:53` |
| `FindAvailableFlowsParams.Keyword` | `keyword` | Go field: FindAvailableFlowsParams.Keyword<br/>type: *string | `internal/approval/resource/my.go:50` |
| `FindAvailableFlowsParams.Labels` | `labels` | Go field: FindAvailableFlowsParams.Labels<br/>type: map[string]string | `internal/approval/resource/my.go:51` |
| `FindAvailableFlowsParams.Page` | `page` | Go field: FindAvailableFlowsParams.Page<br/>type: int | `internal/approval/resource/my.go:52` |
| `FindAvailableFlowsParams.PageSize` | `pageSize` | Go field: FindAvailableFlowsParams.PageSize<br/>type: int | `internal/approval/resource/my.go:53` |
| `FindAvailableFlowsParams.TenantID` | `tenantId` | Go field: FindAvailableFlowsParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:49` |
| `FindCCRecordsParams.IsRead` | `isRead` | Go field: FindCCRecordsParams.IsRead<br/>type: *bool | `internal/approval/resource/my.go:185` |
| `FindCCRecordsParams.Page` | `page` | Go field: FindCCRecordsParams.Page<br/>type: int | `internal/approval/resource/my.go:186` |
| `FindCCRecordsParams.PageSize` | `pageSize` | Go field: FindCCRecordsParams.PageSize<br/>type: int | `internal/approval/resource/my.go:187` |
| `FindCCRecordsParams.TenantID` | `tenantId` | Go field: FindCCRecordsParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:184` |
| `FindCompletedTasksParams.Page` | `page` | Go field: FindCompletedTasksParams.Page<br/>type: int | `internal/approval/resource/my.go:162` |
| `FindCompletedTasksParams.PageSize` | `pageSize` | Go field: FindCompletedTasksParams.PageSize<br/>type: int | `internal/approval/resource/my.go:163` |
| `FindCompletedTasksParams.TenantID` | `tenantId` | Go field: FindCompletedTasksParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:161` |
| `FindFlowsParams.CategoryID` | `categoryId` | Go field: FindFlowsParams.CategoryID<br/>type: *string | `internal/approval/resource/flow.go:241` |
| `FindFlowsParams.IsActive` | `isActive` | Go field: FindFlowsParams.IsActive<br/>type: *bool | `internal/approval/resource/flow.go:243` |
| `FindFlowsParams.Keyword` | `keyword` | Go field: FindFlowsParams.Keyword<br/>type: *string | `internal/approval/resource/flow.go:242` |
| `FindFlowsParams.Labels` | `labels` | Go field: FindFlowsParams.Labels<br/>type: map[string]string | `internal/approval/resource/flow.go:244` |
| `FindFlowsParams.Page` | `page` | Go field: FindFlowsParams.Page<br/>type: int | `internal/approval/resource/flow.go:245` |
| `FindFlowsParams.PageSize` | `pageSize` | Go field: FindFlowsParams.PageSize<br/>type: int | `internal/approval/resource/flow.go:246` |
| `FindFlowsParams.TenantID` | `tenantId` | Go field: FindFlowsParams.TenantID<br/>type: *string | `internal/approval/resource/flow.go:240` |
| `FindInitiatedParams.Keyword` | `keyword` | Go field: FindInitiatedParams.Keyword<br/>type: *string | `internal/approval/resource/my.go:113` |
| `FindInitiatedParams.Page` | `page` | Go field: FindInitiatedParams.Page<br/>type: int | `internal/approval/resource/my.go:114` |
| `FindInitiatedParams.PageSize` | `pageSize` | Go field: FindInitiatedParams.PageSize<br/>type: int | `internal/approval/resource/my.go:115` |
| `FindInitiatedParams.Status` | `status` | Go field: FindInitiatedParams.Status<br/>type: *approval.InstanceStatus | `internal/approval/resource/my.go:112` |
| `FindInitiatedParams.TenantID` | `tenantId` | Go field: FindInitiatedParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:111` |
| `FindInitiatorsParams.FlowID` | `flowId` | Go field: FindInitiatorsParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:398` |
| `FindInitiatorsParams.TenantID` | `tenantId` | Go field: FindInitiatorsParams.TenantID<br/>type: *string | `internal/approval/resource/flow.go:399` |
| `FindPendingTasksParams.Page` | `page` | Go field: FindPendingTasksParams.Page<br/>type: int | `internal/approval/resource/my.go:139` |
| `FindPendingTasksParams.PageSize` | `pageSize` | Go field: FindPendingTasksParams.PageSize<br/>type: int | `internal/approval/resource/my.go:140` |
| `FindPendingTasksParams.TenantID` | `tenantId` | Go field: FindPendingTasksParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:138` |
| `FindVersionsParams.FlowID` | `flowId` | Go field: FindVersionsParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:367` |
| `FindVersionsParams.TenantID` | `tenantId` | Go field: FindVersionsParams.TenantID<br/>type: *string | `internal/approval/resource/flow.go:368` |
| `FiresPreview.NextFiresUnixMs` | `nextFiresUnixMs` | Go field: FiresPreview.NextFiresUnixMs<br/>type: []int64 | `internal/cron/store/schedule_resource.go:139` |
| `Flow.AdminUserIDs` | `adminUserIds` | Go field: Flow.AdminUserIDs<br/>type: []string | `approval/models.go:46` |
| `Flow.BindingMode` | `bindingMode` | Go field: Flow.BindingMode<br/>type: BindingMode | `approval/models.go:44` |
| `Flow.BusinessBinding` | `businessBinding` | Go field: Flow.BusinessBinding<br/>type: *BusinessBindingConfig | `approval/models.go:45` |
| `Flow.CategoryID` | `categoryId` | Go field: Flow.CategoryID<br/>type: string | `approval/models.go:33` |
| `Flow.Code` | `code` | Go field: Flow.Code<br/>type: string | `approval/models.go:34` |
| `Flow.CurrentVersion` | `currentVersion` | Go field: Flow.CurrentVersion<br/>type: int | `approval/models.go:50` |
| `Flow.Description` | `description` | Go field: Flow.Description<br/>type: *string | `approval/models.go:37` |
| `Flow.Icon` | `icon` | Go field: Flow.Icon<br/>type: *string | `approval/models.go:36` |
| `Flow.InstanceTitleTemplate` | `instanceTitleTemplate` | Go field: Flow.InstanceTitleTemplate<br/>type: string | `approval/models.go:48` |
| `Flow.IsActive` | `isActive` | Go field: Flow.IsActive<br/>type: bool | `approval/models.go:49` |
| `Flow.IsAllInitiationAllowed` | `isAllInitiationAllowed` | Go field: Flow.IsAllInitiationAllowed<br/>type: bool | `approval/models.go:47` |
| `Flow.Labels` | `labels` | Go field: Flow.Labels<br/>type: map[string]string | `approval/models.go:43` |
| `Flow.Name` | `name` | Go field: Flow.Name<br/>type: string | `approval/models.go:35` |
| `Flow.TenantID` | `tenantId` | Go field: Flow.TenantID<br/>type: string | `approval/models.go:32` |
| `FlowCategory.Children` | `children` | Go field: FlowCategory.Children<br/>type: []FlowCategory | `approval/models.go:66` |
| `FlowCategory.Code` | `code` | Go field: FlowCategory.Code<br/>type: string | `approval/models.go:59` |
| `FlowCategory.Icon` | `icon` | Go field: FlowCategory.Icon<br/>type: *string | `approval/models.go:61` |
| `FlowCategory.IsActive` | `isActive` | Go field: FlowCategory.IsActive<br/>type: bool | `approval/models.go:64` |
| `FlowCategory.Name` | `name` | Go field: FlowCategory.Name<br/>type: string | `approval/models.go:60` |
| `FlowCategory.ParentID` | `parentId` | Go field: FlowCategory.ParentID<br/>type: *string | `approval/models.go:62` |
| `FlowCategory.Remark` | `remark` | Go field: FlowCategory.Remark<br/>type: *string | `approval/models.go:65` |
| `FlowCategory.SortOrder` | `sortOrder` | Go field: FlowCategory.SortOrder<br/>type: int | `approval/models.go:63` |
| `FlowCategory.TenantID` | `tenantId` | Go field: FlowCategory.TenantID<br/>type: string | `approval/models.go:58` |
| `FlowCreatedEvent.CategoryID` | `categoryId` | Go field: FlowCreatedEvent.CategoryID<br/>type: string | `approval/events_flow.go:7` |
| `FlowDefinition.Edges` | `edges` | Go field: FlowDefinition.Edges<br/>type: []EdgeDefinition | `approval/flow_definition.go:23` |
| `FlowDefinition.Nodes` | `nodes` | Go field: FlowDefinition.Nodes<br/>type: []NodeDefinition | `approval/flow_definition.go:22` |
| `FlowDeployedEvent.Version` | `version` | Go field: FlowDeployedEvent.Version<br/>type: int | `approval/events_flow.go:35` |
| `FlowDeployedEvent.VersionID` | `versionId` | Go field: FlowDeployedEvent.VersionID<br/>type: string | `approval/events_flow.go:34` |
| `FlowEdge.FlowVersionID` | `flowVersionId` | Go field: FlowEdge.FlowVersionID<br/>type: string | `approval/models.go:174` |
| `FlowEdge.Key` | `key` | Go field: FlowEdge.Key<br/>type: string | `approval/models.go:175` |
| `FlowEdge.SourceHandle` | `sourceHandle` | Go field: FlowEdge.SourceHandle<br/>type: *string | `approval/models.go:180` |
| `FlowEdge.SourceNodeID` | `sourceNodeId` | Go field: FlowEdge.SourceNodeID<br/>type: string | `approval/models.go:176` |
| `FlowEdge.SourceNodeKey` | `sourceNodeKey` | Go field: FlowEdge.SourceNodeKey<br/>type: string | `approval/models.go:177` |
| `FlowEdge.TargetNodeID` | `targetNodeId` | Go field: FlowEdge.TargetNodeID<br/>type: string | `approval/models.go:178` |
| `FlowEdge.TargetNodeKey` | `targetNodeKey` | Go field: FlowEdge.TargetNodeKey<br/>type: string | `approval/models.go:179` |
| `FlowEventBase.Code` | `code` | Go field: FlowEventBase.Code<br/>type: string | `approval/events_base.go:72` |
| `FlowEventBase.FlowID` | `flowId` | Go field: FlowEventBase.FlowID<br/>type: string | `approval/events_base.go:70` |
| `FlowEventBase.Name` | `name` | Go field: FlowEventBase.Name<br/>type: string | `approval/events_base.go:73` |
| `FlowEventBase.OccurredTime` | `occurredTime` | Go field: FlowEventBase.OccurredTime<br/>type: timex.DateTime | `approval/events_base.go:74` |
| `FlowEventBase.TenantID` | `tenantId` | Go field: FlowEventBase.TenantID<br/>type: string | `approval/events_base.go:71` |
| `FlowGraph.Edges` | `edges` | Go field: FlowGraph.Edges<br/>type: []approval.FlowEdge | `internal/approval/shared/flow.go:13` |
| `FlowGraph.Flow` | `flow` | Go field: FlowGraph.Flow<br/>type: *approval.Flow | `internal/approval/shared/flow.go:10` |
| `FlowGraph.Nodes` | `nodes` | Go field: FlowGraph.Nodes<br/>type: []approval.FlowNode | `internal/approval/shared/flow.go:12` |
| `FlowGraph.Version` | `version` | Go field: FlowGraph.Version<br/>type: *approval.FlowVersion | `internal/approval/shared/flow.go:11` |
| `FlowGraphEdge.ID` | `id` | Go field: FlowGraphEdge.ID<br/>type: string | `approval/flow_graph_view.go:86` |
| `FlowGraphEdge.Source` | `source` | Go field: FlowGraphEdge.Source<br/>type: string | `approval/flow_graph_view.go:87` |
| `FlowGraphEdge.SourceHandle` | `sourceHandle` | Go field: FlowGraphEdge.SourceHandle<br/>type: *string | `approval/flow_graph_view.go:89` |
| `FlowGraphEdge.Target` | `target` | Go field: FlowGraphEdge.Target<br/>type: string | `approval/flow_graph_view.go:88` |
| `FlowGraphNode.Data` | `data` | Go field: FlowGraphNode.Data<br/>type: FlowGraphNodeData | `approval/flow_graph_view.go:55` |
| `FlowGraphNode.ID` | `id` | Go field: FlowGraphNode.ID<br/>type: string | `approval/flow_graph_view.go:51` |
| `FlowGraphNode.Kind` | `kind` | Go field: FlowGraphNode.Kind<br/>type: NodeKind | `approval/flow_graph_view.go:53` |
| `FlowGraphNode.NodeID` | `nodeId` | Go field: FlowGraphNode.NodeID<br/>type: string | `approval/flow_graph_view.go:52` |
| `FlowGraphNode.Position` | `position` | Go field: FlowGraphNode.Position<br/>type: Position | `approval/flow_graph_view.go:54` |
| `FlowGraphNodeData.Activities` | `activities` | Go field: FlowGraphNodeData.Activities<br/>type: []Activity | `approval/flow_graph_view.go:79` |
| `FlowGraphNodeData.ApprovalMethod` | `approvalMethod` | Go field: FlowGraphNodeData.ApprovalMethod<br/>type: string | `approval/flow_graph_view.go:74` |
| `FlowGraphNodeData.CCRecipients` | `ccRecipients` | Go field: FlowGraphNodeData.CCRecipients<br/>type: []CCRecipient | `approval/flow_graph_view.go:78` |
| `FlowGraphNodeData.ExecutionType` | `executionType` | Go field: FlowGraphNodeData.ExecutionType<br/>type: string | `approval/flow_graph_view.go:73` |
| `FlowGraphNodeData.FinishedAt` | `finishedAt` | Go field: FlowGraphNodeData.FinishedAt<br/>type: *timex.DateTime | `approval/flow_graph_view.go:81` |
| `FlowGraphNodeData.Name` | `name` | Go field: FlowGraphNodeData.Name<br/>type: string | `approval/flow_graph_view.go:71` |
| `FlowGraphNodeData.Participants` | `participants` | Go field: FlowGraphNodeData.Participants<br/>type: []NodeParticipant | `approval/flow_graph_view.go:77` |
| `FlowGraphNodeData.PassRatio` | `passRatio` | Go field: FlowGraphNodeData.PassRatio<br/>type: *decimal.Decimal | `approval/flow_graph_view.go:76` |
| `FlowGraphNodeData.PassRule` | `passRule` | Go field: FlowGraphNodeData.PassRule<br/>type: string | `approval/flow_graph_view.go:75` |
| `FlowGraphNodeData.StartedAt` | `startedAt` | Go field: FlowGraphNodeData.StartedAt<br/>type: *timex.DateTime | `approval/flow_graph_view.go:80` |
| `FlowGraphNodeData.Status` | `status` | Go field: FlowGraphNodeData.Status<br/>type: NodeProgressStatus | `approval/flow_graph_view.go:72` |
| `FlowInitiator.FlowID` | `flowId` | Go field: FlowInitiator.FlowID<br/>type: string | `approval/models.go:212` |
| `FlowInitiator.IDs` | `ids` | Go field: FlowInitiator.IDs<br/>type: []string | `approval/models.go:214` |
| `FlowInitiator.Kind` | `kind` | Go field: FlowInitiator.Kind<br/>type: InitiatorKind | `approval/models.go:213` |
| `FlowNode.AddAssigneeTypes` | `addAssigneeTypes` | Go field: FlowNode.AddAssigneeTypes<br/>type: []AddAssigneeType | `approval/models.go:154` |
| `FlowNode.AdminUserIDs` | `adminUserIds` | Go field: FlowNode.AdminUserIDs<br/>type: []string | `approval/models.go:147` |
| `FlowNode.ApprovalMethod` | `approvalMethod` | Go field: FlowNode.ApprovalMethod<br/>type: ApprovalMethod | `approval/models.go:142` |
| `FlowNode.Branches` | `branches` | Go field: FlowNode.Branches<br/>type: []ConditionBranch | `approval/models.go:166` |
| `FlowNode.ConsecutiveApproverAction` | `consecutiveApproverAction` | Go field: FlowNode.ConsecutiveApproverAction<br/>type: ConsecutiveApproverAction | `approval/models.go:164` |
| `FlowNode.Description` | `description` | Go field: FlowNode.Description<br/>type: *string | `approval/models.go:140` |
| `FlowNode.EmptyAssigneeAction` | `emptyAssigneeAction` | Go field: FlowNode.EmptyAssigneeAction<br/>type: EmptyAssigneeAction | `approval/models.go:145` |
| `FlowNode.ExecutionType` | `executionType` | Go field: FlowNode.ExecutionType<br/>type: ExecutionType | `approval/models.go:141` |
| `FlowNode.FallbackUserIDs` | `fallbackUserIds` | Go field: FlowNode.FallbackUserIDs<br/>type: []string | `approval/models.go:146` |
| `FlowNode.FieldPermissions` | `fieldPermissions` | Go field: FlowNode.FieldPermissions<br/>type: map[string]Permission | `approval/models.go:156` |
| `FlowNode.FlowVersionID` | `flowVersionId` | Go field: FlowNode.FlowVersionID<br/>type: string | `approval/models.go:136` |
| `FlowNode.IsAddAssigneeAllowed` | `isAddAssigneeAllowed` | Go field: FlowNode.IsAddAssigneeAllowed<br/>type: bool | `approval/models.go:153` |
| `FlowNode.IsManualCCAllowed` | `isManualCcAllowed` | Go field: FlowNode.IsManualCCAllowed<br/>type: bool | `approval/models.go:157` |
| `FlowNode.IsOpinionRequired` | `isOpinionRequired` | Go field: FlowNode.IsOpinionRequired<br/>type: bool | `approval/models.go:159` |
| `FlowNode.IsReadConfirmRequired` | `isReadConfirmRequired` | Go field: FlowNode.IsReadConfirmRequired<br/>type: bool | `approval/models.go:165` |
| `FlowNode.IsRemoveAssigneeAllowed` | `isRemoveAssigneeAllowed` | Go field: FlowNode.IsRemoveAssigneeAllowed<br/>type: bool | `approval/models.go:155` |
| `FlowNode.IsRollbackAllowed` | `isRollbackAllowed` | Go field: FlowNode.IsRollbackAllowed<br/>type: bool | `approval/models.go:149` |
| `FlowNode.IsTransferAllowed` | `isTransferAllowed` | Go field: FlowNode.IsTransferAllowed<br/>type: bool | `approval/models.go:158` |
| `FlowNode.Key` | `key` | Go field: FlowNode.Key<br/>type: string | `approval/models.go:137` |
| `FlowNode.Kind` | `kind` | Go field: FlowNode.Kind<br/>type: NodeKind | `approval/models.go:138` |
| `FlowNode.Name` | `name` | Go field: FlowNode.Name<br/>type: string | `approval/models.go:139` |
| `FlowNode.PassRatio` | `passRatio` | Go field: FlowNode.PassRatio<br/>type: decimal.Decimal | `approval/models.go:144` |
| `FlowNode.PassRule` | `passRule` | Go field: FlowNode.PassRule<br/>type: PassRule | `approval/models.go:143` |
| `FlowNode.RollbackDataStrategy` | `rollbackDataStrategy` | Go field: FlowNode.RollbackDataStrategy<br/>type: RollbackDataStrategy | `approval/models.go:151` |
| `FlowNode.RollbackTargetKeys` | `rollbackTargetKeys` | Go field: FlowNode.RollbackTargetKeys<br/>type: []string | `approval/models.go:152` |
| `FlowNode.RollbackType` | `rollbackType` | Go field: FlowNode.RollbackType<br/>type: RollbackType | `approval/models.go:150` |
| `FlowNode.SameApplicantAction` | `sameApplicantAction` | Go field: FlowNode.SameApplicantAction<br/>type: SameApplicantAction | `approval/models.go:148` |
| `FlowNode.TimeoutAction` | `timeoutAction` | Go field: FlowNode.TimeoutAction<br/>type: TimeoutAction | `approval/models.go:161` |
| `FlowNode.TimeoutHours` | `timeoutHours` | Go field: FlowNode.TimeoutHours<br/>type: int | `approval/models.go:160` |
| `FlowNode.TimeoutNotifyBeforeHours` | `timeoutNotifyBeforeHours` | Go field: FlowNode.TimeoutNotifyBeforeHours<br/>type: int | `approval/models.go:162` |
| `FlowNode.UrgeCooldownMinutes` | `urgeCooldownMinutes` | Go field: FlowNode.UrgeCooldownMinutes<br/>type: int | `approval/models.go:163` |
| `FlowNodeAssignee.FormField` | `formField` | Go field: FlowNodeAssignee.FormField<br/>type: *string | `approval/models.go:191` |
| `FlowNodeAssignee.IDs` | `ids` | Go field: FlowNodeAssignee.IDs<br/>type: []string | `approval/models.go:190` |
| `FlowNodeAssignee.Kind` | `kind` | Go field: FlowNodeAssignee.Kind<br/>type: AssigneeKind | `approval/models.go:189` |
| `FlowNodeAssignee.NodeID` | `nodeId` | Go field: FlowNodeAssignee.NodeID<br/>type: string | `approval/models.go:188` |
| `FlowNodeAssignee.SortOrder` | `sortOrder` | Go field: FlowNodeAssignee.SortOrder<br/>type: int | `approval/models.go:192` |
| `FlowNodeCC.FormField` | `formField` | Go field: FlowNodeCC.FormField<br/>type: *string | `approval/models.go:203` |
| `FlowNodeCC.IDs` | `ids` | Go field: FlowNodeCC.IDs<br/>type: []string | `approval/models.go:202` |
| `FlowNodeCC.Kind` | `kind` | Go field: FlowNodeCC.Kind<br/>type: CCKind | `approval/models.go:201` |
| `FlowNodeCC.NodeID` | `nodeId` | Go field: FlowNodeCC.NodeID<br/>type: string | `approval/models.go:200` |
| `FlowNodeCC.Timing` | `timing` | Go field: FlowNodeCC.Timing<br/>type: CCTiming | `approval/models.go:204` |
| `FlowPublishedEvent.VersionID` | `versionId` | Go field: FlowPublishedEvent.VersionID<br/>type: string | `approval/events_flow.go:68` |
| `FlowToggledEvent.IsActive` | `isActive` | Go field: FlowToggledEvent.IsActive<br/>type: bool | `approval/events_flow.go:52` |
| `FlowVersion.BusinessBinding` | `businessBinding` | Go field: FlowVersion.BusinessBinding<br/>type: *BusinessBindingConfig | `approval/models.go:92` |
| `FlowVersion.Description` | `description` | Go field: FlowVersion.Description<br/>type: *string | `approval/models.go:77` |
| `FlowVersion.FlowID` | `flowId` | Go field: FlowVersion.FlowID<br/>type: string | `approval/models.go:74` |
| `FlowVersion.FlowSchema` | `flowSchema` | Go field: FlowVersion.FlowSchema<br/>type: *FlowDefinition | `approval/models.go:79` |
| `FlowVersion.FormFields` | `formFields` | Go field: FlowVersion.FormFields<br/>type: []FormFieldDefinition | `approval/models.go:87` |
| `FlowVersion.FormSchema` | `formSchema` | Go field: FlowVersion.FormSchema<br/>type: json.RawMessage | `approval/models.go:84` |
| `FlowVersion.PublishedAt` | `publishedAt` | Go field: FlowVersion.PublishedAt<br/>type: *timex.DateTime | `approval/models.go:88` |
| `FlowVersion.PublishedBy` | `publishedBy` | Go field: FlowVersion.PublishedBy<br/>type: *string | `approval/models.go:89` |
| `FlowVersion.Status` | `status` | Go field: FlowVersion.Status<br/>type: VersionStatus | `approval/models.go:76` |
| `FlowVersion.StorageMode` | `storageMode` | Go field: FlowVersion.StorageMode<br/>type: StorageMode | `approval/models.go:78` |
| `FlowVersion.Version` | `version` | Go field: FlowVersion.Version<br/>type: int | `approval/models.go:75` |
| `FlowVersionSummary.CreatedAt` | `createdAt` | Go field: FlowVersionSummary.CreatedAt<br/>type: timex.DateTime | `internal/approval/shared/flow.go:30` |
| `FlowVersionSummary.CreatedBy` | `createdBy` | Go field: FlowVersionSummary.CreatedBy<br/>type: string | `internal/approval/shared/flow.go:31` |
| `FlowVersionSummary.Description` | `description` | Go field: FlowVersionSummary.Description<br/>type: *string | `internal/approval/shared/flow.go:26` |
| `FlowVersionSummary.FlowID` | `flowId` | Go field: FlowVersionSummary.FlowID<br/>type: string | `internal/approval/shared/flow.go:23` |
| `FlowVersionSummary.ID` | `id` | Go field: FlowVersionSummary.ID<br/>type: string | `internal/approval/shared/flow.go:22` |
| `FlowVersionSummary.PublishedAt` | `publishedAt` | Go field: FlowVersionSummary.PublishedAt<br/>type: *timex.DateTime | `internal/approval/shared/flow.go:28` |
| `FlowVersionSummary.PublishedBy` | `publishedBy` | Go field: FlowVersionSummary.PublishedBy<br/>type: *string | `internal/approval/shared/flow.go:29` |
| `FlowVersionSummary.Status` | `status` | Go field: FlowVersionSummary.Status<br/>type: approval.VersionStatus | `internal/approval/shared/flow.go:25` |
| `FlowVersionSummary.StorageMode` | `storageMode` | Go field: FlowVersionSummary.StorageMode<br/>type: approval.StorageMode | `internal/approval/shared/flow.go:27` |
| `FlowVersionSummary.Version` | `version` | Go field: FlowVersionSummary.Version<br/>type: int | `internal/approval/shared/flow.go:24` |
| `ForeignKey.Columns` | `columns` | Go field: ForeignKey.Columns<br/>type: []string | `schema/service.go:59` |
| `ForeignKey.Name` | `name` | Go field: ForeignKey.Name<br/>type: string | `schema/service.go:58` |
| `ForeignKey.OnDelete` | `onDelete` | Go field: ForeignKey.OnDelete<br/>type: string | `schema/service.go:63` |
| `ForeignKey.OnUpdate` | `onUpdate` | Go field: ForeignKey.OnUpdate<br/>type: string | `schema/service.go:62` |
| `ForeignKey.RefColumns` | `refColumns` | Go field: ForeignKey.RefColumns<br/>type: []string | `schema/service.go:61` |
| `ForeignKey.RefTable` | `refTable` | Go field: ForeignKey.RefTable<br/>type: string | `schema/service.go:60` |
| `FormFieldDefinition.ColumnType` | `columnType` | Go field: FormFieldDefinition.ColumnType<br/>type: ColumnDataType | `approval/form_field.go:28` |
| `FormFieldDefinition.Columns` | `columns` | Go field: FormFieldDefinition.Columns<br/>type: []FormFieldDefinition | `approval/form_field.go:38` |
| `FormFieldDefinition.DefaultValue` | `defaultValue` | Go field: FormFieldDefinition.DefaultValue<br/>type: any | `approval/form_field.go:14` |
| `FormFieldDefinition.IsRequired` | `isRequired` | Go field: FormFieldDefinition.IsRequired<br/>type: bool | `approval/form_field.go:16` |
| `FormFieldDefinition.Key` | `key` | Go field: FormFieldDefinition.Key<br/>type: string | `approval/form_field.go:6` |
| `FormFieldDefinition.Kind` | `kind` | Go field: FormFieldDefinition.Kind<br/>type: FieldKind | `approval/form_field.go:8` |
| `FormFieldDefinition.Label` | `label` | Go field: FormFieldDefinition.Label<br/>type: string | `approval/form_field.go:10` |
| `FormFieldDefinition.Options` | `options` | Go field: FormFieldDefinition.Options<br/>type: []FieldOption | `approval/form_field.go:18` |
| `FormFieldDefinition.Placeholder` | `placeholder` | Go field: FormFieldDefinition.Placeholder<br/>type: string | `approval/form_field.go:12` |
| `FormFieldDefinition.Props` | `props` | Go field: FormFieldDefinition.Props<br/>type: map[string]any | `approval/form_field.go:22` |
| `FormFieldDefinition.Scale` | `scale` | Go field: FormFieldDefinition.Scale<br/>type: *int | `approval/form_field.go:31` |
| `FormFieldDefinition.SortOrder` | `sortOrder` | Go field: FormFieldDefinition.SortOrder<br/>type: int | `approval/form_field.go:24` |
| `FormFieldDefinition.Validation` | `validation` | Go field: FormFieldDefinition.Validation<br/>type: *ValidationRule | `approval/form_field.go:20` |
| `FormSnapshot.FormData` | `formData` | Go field: FormSnapshot.FormData<br/>type: map[string]any | `approval/models.go:384` |
| `FormSnapshot.InstanceID` | `instanceId` | Go field: FormSnapshot.InstanceID<br/>type: string | `approval/models.go:382` |
| `FormSnapshot.NodeID` | `nodeId` | Go field: FormSnapshot.NodeID<br/>type: string | `approval/models.go:383` |
| `FormTable.FlowID` | `flowId` | Go field: FormTable.FlowID<br/>type: string | `approval/models.go:107` |
| `FormTable.PhysicalTableName` | `physicalTableName` | Go field: FormTable.PhysicalTableName<br/>type: string | `approval/models.go:109` |
| `FormTable.SourceFieldKey` | `sourceFieldKey` | Go field: FormTable.SourceFieldKey<br/>type: string | `approval/models.go:112` |
| `FormTable.VersionID` | `versionId` | Go field: FormTable.VersionID<br/>type: string | `approval/models.go:108` |
| `FormTableColumn.ColumnName` | `columnName` | Go field: FormTableColumn.ColumnName<br/>type: string | `approval/models.go:124` |
| `FormTableColumn.ColumnType` | `columnType` | Go field: FormTableColumn.ColumnType<br/>type: string | `approval/models.go:125` |
| `FormTableColumn.FormTableID` | `formTableId` | Go field: FormTableColumn.FormTableID<br/>type: string | `approval/models.go:123` |
| `FormTableColumn.IsNullable` | `isNullable` | Go field: FormTableColumn.IsNullable<br/>type: bool | `approval/models.go:126` |
| `FormTableColumn.SortOrder` | `sortOrder` | Go field: FormTableColumn.SortOrder<br/>type: int | `approval/models.go:128` |
| `FormTableColumn.SourceFieldKey` | `sourceFieldKey` | Go field: FormTableColumn.SourceFieldKey<br/>type: *string | `approval/models.go:127` |
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
| `GetGraphParams.FlowID` | `flowId` | Go field: GetGraphParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:190` |
| `GetGraphParams.TenantID` | `tenantId` | Go field: GetGraphParams.TenantID<br/>type: *string | `internal/approval/resource/flow.go:194` |
| `GetGraphParams.VersionID` | `versionId` | Go field: GetGraphParams.VersionID<br/>type: *string | `internal/approval/resource/flow.go:198` |
| `GetInstanceDetailParams.InstanceID` | `instanceId` | Go field: GetInstanceDetailParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/my.go:229` |
| `GetPendingCountsParams.TenantID` | `tenantId` | Go field: GetPendingCountsParams.TenantID<br/>type: *string | `internal/approval/resource/my.go:209` |
| `GetStartFormParams.FlowCode` | `flowCode` | Go field: GetStartFormParams.FlowCode<br/>type: string<br/>validate: "required" | `internal/approval/resource/my.go:83` |
| `GetStartFormParams.TenantID` | `tenantId` | Go field: GetStartFormParams.TenantID<br/>type: string<br/>validate: "required" | `internal/approval/resource/my.go:82` |
| `GetTableSchemaParams.Name` | `name` | Go field: GetTableSchemaParams.Name<br/>type: string<br/>validate: "required" | `internal/schema/resource.go:49` |
| `HTTPExchange.DurationMs` | `durationMs` | Go field: HTTPExchange.DurationMs<br/>type: int64 | `integration/models.go:291` |
| `HTTPExchange.Error` | `error` | Go field: HTTPExchange.Error<br/>type: string | `integration/models.go:292` |
| `HTTPExchange.Method` | `method` | Go field: HTTPExchange.Method<br/>type: string | `integration/models.go:284` |
| `HTTPExchange.RequestBody` | `requestBody` | Go field: HTTPExchange.RequestBody<br/>type: string | `integration/models.go:287` |
| `HTTPExchange.RequestHeaders` | `requestHeaders` | Go field: HTTPExchange.RequestHeaders<br/>type: map[string]string | `integration/models.go:286` |
| `HTTPExchange.ResponseBody` | `responseBody` | Go field: HTTPExchange.ResponseBody<br/>type: string | `integration/models.go:290` |
| `HTTPExchange.ResponseHeaders` | `responseHeaders` | Go field: HTTPExchange.ResponseHeaders<br/>type: map[string]string | `integration/models.go:289` |
| `HTTPExchange.Status` | `status` | Go field: HTTPExchange.Status<br/>type: int | `integration/models.go:288` |
| `HTTPExchange.URL` | `url` | Go field: HTTPExchange.URL<br/>type: string | `integration/models.go:285` |
| `HTTPProbe.DurationMs` | `durationMs` | Go field: HTTPProbe.DurationMs<br/>type: int64 | `internal/integration/exec/connection_check.go:25` |
| `HTTPProbe.Error` | `error` | Go field: HTTPProbe.Error<br/>type: string | `internal/integration/exec/connection_check.go:26` |
| `HTTPProbe.Reachable` | `reachable` | Go field: HTTPProbe.Reachable<br/>type: bool | `internal/integration/exec/connection_check.go:22` |
| `HTTPProbe.Status` | `status` | Go field: HTTPProbe.Status<br/>type: int | `internal/integration/exec/connection_check.go:23` |
| `HTTPProbe.StatusText` | `statusText` | Go field: HTTPProbe.StatusText<br/>type: string | `internal/integration/exec/connection_check.go:24` |
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
| `IOCounter.IOPSInProgress` | `iopsInProgress` | Go field: IOCounter.IOPSInProgress<br/>type: uint64 | `monitor/service.go:184` |
| `IOCounter.IOTime` | `ioTime` | Go field: IOCounter.IOTime<br/>type: uint64 | `monitor/service.go:185` |
| `IOCounter.Label` | `label` | Go field: IOCounter.Label<br/>type: string | `monitor/service.go:189` |
| `IOCounter.MergedReadCount` | `mergedReadCount` | Go field: IOCounter.MergedReadCount<br/>type: uint64 | `monitor/service.go:177` |
| `IOCounter.MergedWriteCount` | `mergedWriteCount` | Go field: IOCounter.MergedWriteCount<br/>type: uint64 | `monitor/service.go:179` |
| `IOCounter.Name` | `name` | Go field: IOCounter.Name<br/>type: string | `monitor/service.go:187` |
| `IOCounter.ReadBytes` | `readBytes` | Go field: IOCounter.ReadBytes<br/>type: uint64 | `monitor/service.go:180` |
| `IOCounter.ReadCount` | `readCount` | Go field: IOCounter.ReadCount<br/>type: uint64 | `monitor/service.go:176` |
| `IOCounter.ReadTime` | `readTime` | Go field: IOCounter.ReadTime<br/>type: uint64 | `monitor/service.go:182` |
| `IOCounter.SerialNumber` | `serialNumber` | Go field: IOCounter.SerialNumber<br/>type: string | `monitor/service.go:188` |
| `IOCounter.WeightedIO` | `weightedIo` | Go field: IOCounter.WeightedIO<br/>type: uint64 | `monitor/service.go:186` |
| `IOCounter.WriteBytes` | `writeBytes` | Go field: IOCounter.WriteBytes<br/>type: uint64 | `monitor/service.go:181` |
| `IOCounter.WriteCount` | `writeCount` | Go field: IOCounter.WriteCount<br/>type: uint64 | `monitor/service.go:178` |
| `IOCounter.WriteTime` | `writeTime` | Go field: IOCounter.WriteTime<br/>type: uint64 | `monitor/service.go:183` |
| `Identifier.Action` | `action` | Go field: Identifier.Action<br/>form: "action"<br/>type: string<br/>validate: "required" | `api/request.go:16` |
| `Identifier.Resource` | `resource` | Go field: Identifier.Resource<br/>form: "resource"<br/>type: string<br/>validate: "required,alphanum_us_slash" | `api/request.go:15` |
| `Identifier.Version` | `version` | Go field: Identifier.Version<br/>form: "version"<br/>type: string<br/>validate: "required,alphanum" | `api/request.go:17` |
| `InboundAuthConfig.Params` | `params` | Go field: InboundAuthConfig.Params<br/>type: map[string]string | `integration/models.go:184` |
| `InboundAuthConfig.Scheme` | `scheme` | Go field: InboundAuthConfig.Scheme<br/>type: string | `integration/models.go:183` |
| `InboundAuthConfig.Script` | `script` | Go field: InboundAuthConfig.Script<br/>type: string | `integration/models.go:188` |
| `InboundDryRunResult.DispatchedInput` | `dispatchedInput` | Go field: InboundDryRunResult.DispatchedInput<br/>type: any | `internal/integration/exec/receiver.go:331` |
| `InboundDryRunResult.Error` | `error` | Go field: InboundDryRunResult.Error<br/>type: string | `internal/integration/exec/receiver.go:333` |
| `InboundDryRunResult.FailureKind` | `failureKind` | Go field: InboundDryRunResult.FailureKind<br/>type: integration.FailureKind | `internal/integration/exec/receiver.go:332` |
| `InboundDryRunResult.Reply` | `reply` | Go field: InboundDryRunResult.Reply<br/>type: any | `internal/integration/exec/receiver.go:330` |
| `InboundRequestParams.Body` | `body` | Go field: InboundRequestParams.Body<br/>type: string | `internal/integration/resource/ops.go:50` |
| `InboundRequestParams.Headers` | `headers` | Go field: InboundRequestParams.Headers<br/>type: map[string]string | `internal/integration/resource/ops.go:48` |
| `InboundRequestParams.Method` | `method` | Go field: InboundRequestParams.Method<br/>type: string | `internal/integration/resource/ops.go:46` |
| `InboundRequestParams.Path` | `path` | Go field: InboundRequestParams.Path<br/>type: string | `internal/integration/resource/ops.go:47` |
| `InboundRequestParams.Query` | `query` | Go field: InboundRequestParams.Query<br/>type: map[string]string | `internal/integration/resource/ops.go:49` |
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
| `InitiatedInstance.CreatedAt` | `createdAt` | Go field: InitiatedInstance.CreatedAt<br/>type: timex.DateTime | `approval/my/initiated_instances.go:14` |
| `InitiatedInstance.CurrentNodeName` | `currentNodeName` | Go field: InitiatedInstance.CurrentNodeName<br/>type: *string | `approval/my/initiated_instances.go:13` |
| `InitiatedInstance.FinishedAt` | `finishedAt` | Go field: InitiatedInstance.FinishedAt<br/>type: *timex.DateTime | `approval/my/initiated_instances.go:15` |
| `InitiatedInstance.FlowIcon` | `flowIcon` | Go field: InitiatedInstance.FlowIcon<br/>type: *string | `approval/my/initiated_instances.go:11` |
| `InitiatedInstance.FlowName` | `flowName` | Go field: InitiatedInstance.FlowName<br/>type: string | `approval/my/initiated_instances.go:10` |
| `InitiatedInstance.InstanceID` | `instanceId` | Go field: InitiatedInstance.InstanceID<br/>type: string | `approval/my/initiated_instances.go:7` |
| `InitiatedInstance.InstanceNo` | `instanceNo` | Go field: InitiatedInstance.InstanceNo<br/>type: string | `approval/my/initiated_instances.go:8` |
| `InitiatedInstance.Status` | `status` | Go field: InitiatedInstance.Status<br/>type: string | `approval/my/initiated_instances.go:12` |
| `InitiatedInstance.Title` | `title` | Go field: InitiatedInstance.Title<br/>type: string | `approval/my/initiated_instances.go:9` |
| `Instance.Applicant` | `applicant` | Go field: Instance.Applicant<br/>type: approval.UserInfo | `approval/admin/instance.go:16` |
| `Instance.ApplicantDepartmentID` | `applicantDepartmentId` | Go field: Instance.ApplicantDepartmentID<br/>type: *string | `approval/models.go:236` |
| `Instance.ApplicantDepartmentName` | `applicantDepartmentName` | Go field: Instance.ApplicantDepartmentName<br/>type: *string | `approval/models.go:237` |
| `Instance.ApplicantID` | `applicantId` | Go field: Instance.ApplicantID<br/>type: string | `approval/models.go:234` |
| `Instance.ApplicantName` | `applicantName` | Go field: Instance.ApplicantName<br/>type: string | `approval/models.go:235` |
| `Instance.BusinessProjectionID` | `businessProjectionId` | Go field: Instance.BusinessProjectionID<br/>type: *string | `approval/models.go:256` |
| `Instance.BusinessRef` | `businessRef` | Go field: Instance.BusinessRef<br/>type: *string | `approval/models.go:245` |
| `Instance.CreatedAt` | `createdAt` | Go field: Instance.CreatedAt<br/>type: timex.DateTime | `approval/admin/instance.go:19` |
| `Instance.CurrentNodeID` | `currentNodeId` | Go field: Instance.CurrentNodeID<br/>type: *string | `approval/models.go:239` |
| `Instance.CurrentNodeName` | `currentNodeName` | Go field: Instance.CurrentNodeName<br/>type: *string | `approval/admin/instance.go:18` |
| `Instance.FinishedAt` | `finishedAt` | Go field: Instance.FinishedAt<br/>type: *timex.DateTime | `approval/models.go:240` |
| `Instance.FinishedAt` | `finishedAt` | Go field: Instance.FinishedAt<br/>type: *timex.DateTime | `approval/admin/instance.go:20` |
| `Instance.FlowCode` | `flowCode` | Go field: Instance.FlowCode<br/>type: string | `approval/models.go:230` |
| `Instance.FlowID` | `flowId` | Go field: Instance.FlowID<br/>type: string | `approval/models.go:225` |
| `Instance.FlowID` | `flowId` | Go field: Instance.FlowID<br/>type: string | `approval/admin/instance.go:14` |
| `Instance.FlowName` | `flowName` | Go field: Instance.FlowName<br/>type: string | `approval/admin/instance.go:15` |
| `Instance.FlowVersionID` | `flowVersionId` | Go field: Instance.FlowVersionID<br/>type: string | `approval/models.go:231` |
| `Instance.FormData` | `formData` | Go field: Instance.FormData<br/>type: map[string]any | `approval/models.go:246` |
| `Instance.Globals` | `globals` | Go field: Instance.Globals<br/>type: map[string]any | `approval/models.go:252` |
| `Instance.InstanceID` | `instanceId` | Go field: Instance.InstanceID<br/>type: string | `approval/admin/instance.go:10` |
| `Instance.InstanceNo` | `instanceNo` | Go field: Instance.InstanceNo<br/>type: string | `approval/admin/instance.go:11` |
| `Instance.InstanceNo` | `instanceNo` | Go field: Instance.InstanceNo<br/>type: string | `approval/models.go:233` |
| `Instance.Status` | `status` | Go field: Instance.Status<br/>type: InstanceStatus | `approval/models.go:238` |
| `Instance.Status` | `status` | Go field: Instance.Status<br/>type: string | `approval/admin/instance.go:17` |
| `Instance.TenantID` | `tenantId` | Go field: Instance.TenantID<br/>type: string | `approval/admin/instance.go:13` |
| `Instance.TenantID` | `tenantId` | Go field: Instance.TenantID<br/>type: string | `approval/models.go:224` |
| `Instance.Title` | `title` | Go field: Instance.Title<br/>type: string | `approval/admin/instance.go:12` |
| `Instance.Title` | `title` | Go field: Instance.Title<br/>type: string | `approval/models.go:232` |
| `InstanceBindingFailedEvent.BusinessTable` | `businessTable` | Go field: InstanceBindingFailedEvent.BusinessTable<br/>type: string | `approval/events_instance.go:158` |
| `InstanceBindingFailedEvent.ErrorMessage` | `errorMessage` | Go field: InstanceBindingFailedEvent.ErrorMessage<br/>type: string | `approval/events_instance.go:159` |
| `InstanceBindingFailedEvent.Status` | `status` | Go field: InstanceBindingFailedEvent.Status<br/>type: InstanceStatus | `approval/events_instance.go:157` |
| `InstanceBindingFailedEvent.Trigger` | `trigger` | Go field: InstanceBindingFailedEvent.Trigger<br/>type: BindingTrigger | `approval/events_instance.go:154` |
| `InstanceCompletedEvent.FinalStatus` | `finalStatus` | Go field: InstanceCompletedEvent.FinalStatus<br/>type: InstanceStatus | `approval/events_instance.go:23` |
| `InstanceCompletedEvent.FinishedAt` | `finishedAt` | Go field: InstanceCompletedEvent.FinishedAt<br/>type: timex.DateTime | `approval/events_instance.go:24` |
| `InstanceCompletedEvent.Reason` | `reason` | Go field: InstanceCompletedEvent.Reason<br/>type: *string | `approval/events_instance.go:29` |
| `InstanceDetail.AvailableActions` | `availableActions` | Go field: InstanceDetail.AvailableActions<br/>type: []string | `approval/my/instance_detail.go:22` |
| `InstanceDetail.FieldPermissions` | `fieldPermissions` | Go field: InstanceDetail.FieldPermissions<br/>type: map[string]approval.Permission | `approval/my/instance_detail.go:27` |
| `InstanceDetail.FlowGraph` | `flowGraph` | Go field: InstanceDetail.FlowGraph<br/>type: approval.InstanceFlowGraph | `approval/admin/instance_detail.go:20` |
| `InstanceDetail.FlowGraph` | `flowGraph` | Go field: InstanceDetail.FlowGraph<br/>type: approval.InstanceFlowGraph | `approval/my/instance_detail.go:21` |
| `InstanceDetail.FormSchema` | `formSchema` | Go field: InstanceDetail.FormSchema<br/>type: json.RawMessage | `approval/admin/instance_detail.go:18` |
| `InstanceDetail.FormSchema` | `formSchema` | Go field: InstanceDetail.FormSchema<br/>type: json.RawMessage | `approval/my/instance_detail.go:19` |
| `InstanceDetail.Instance` | `instance` | Go field: InstanceDetail.Instance<br/>type: InstanceDetailInfo | `approval/admin/instance_detail.go:17` |
| `InstanceDetail.Instance` | `instance` | Go field: InstanceDetail.Instance<br/>type: InstanceInfo | `approval/my/instance_detail.go:18` |
| `InstanceDetail.MyTask` | `myTask` | Go field: InstanceDetail.MyTask<br/>type: *ViewerTask | `approval/my/instance_detail.go:32` |
| `InstanceDetail.Timeline` | `timeline` | Go field: InstanceDetail.Timeline<br/>type: []approval.TimelineEntry | `approval/admin/instance_detail.go:19` |
| `InstanceDetail.Timeline` | `timeline` | Go field: InstanceDetail.Timeline<br/>type: []approval.TimelineEntry | `approval/my/instance_detail.go:20` |
| `InstanceDetailInfo.Applicant` | `applicant` | Go field: InstanceDetailInfo.Applicant<br/>type: approval.UserInfo | `approval/admin/instance_detail.go:37` |
| `InstanceDetailInfo.BusinessRef` | `businessRef` | Go field: InstanceDetailInfo.BusinessRef<br/>type: *string | `approval/admin/instance_detail.go:41` |
| `InstanceDetailInfo.CreatedAt` | `createdAt` | Go field: InstanceDetailInfo.CreatedAt<br/>type: timex.DateTime | `approval/admin/instance_detail.go:43` |
| `InstanceDetailInfo.CurrentNodeID` | `currentNodeId` | Go field: InstanceDetailInfo.CurrentNodeID<br/>type: *string | `approval/admin/instance_detail.go:39` |
| `InstanceDetailInfo.CurrentNodeName` | `currentNodeName` | Go field: InstanceDetailInfo.CurrentNodeName<br/>type: *string | `approval/admin/instance_detail.go:40` |
| `InstanceDetailInfo.FinishedAt` | `finishedAt` | Go field: InstanceDetailInfo.FinishedAt<br/>type: *timex.DateTime | `approval/admin/instance_detail.go:44` |
| `InstanceDetailInfo.FlowID` | `flowId` | Go field: InstanceDetailInfo.FlowID<br/>type: string | `approval/admin/instance_detail.go:30` |
| `InstanceDetailInfo.FlowName` | `flowName` | Go field: InstanceDetailInfo.FlowName<br/>type: string | `approval/admin/instance_detail.go:31` |
| `InstanceDetailInfo.FlowVersionID` | `flowVersionId` | Go field: InstanceDetailInfo.FlowVersionID<br/>type: string | `approval/admin/instance_detail.go:32` |
| `InstanceDetailInfo.FormData` | `formData` | Go field: InstanceDetailInfo.FormData<br/>type: map[string]any | `approval/admin/instance_detail.go:42` |
| `InstanceDetailInfo.InstanceID` | `instanceId` | Go field: InstanceDetailInfo.InstanceID<br/>type: string | `approval/admin/instance_detail.go:26` |
| `InstanceDetailInfo.InstanceNo` | `instanceNo` | Go field: InstanceDetailInfo.InstanceNo<br/>type: string | `approval/admin/instance_detail.go:27` |
| `InstanceDetailInfo.Labels` | `labels` | Go field: InstanceDetailInfo.Labels<br/>type: map[string]string | `approval/admin/instance_detail.go:36` |
| `InstanceDetailInfo.Status` | `status` | Go field: InstanceDetailInfo.Status<br/>type: string | `approval/admin/instance_detail.go:38` |
| `InstanceDetailInfo.TenantID` | `tenantId` | Go field: InstanceDetailInfo.TenantID<br/>type: string | `approval/admin/instance_detail.go:29` |
| `InstanceDetailInfo.Title` | `title` | Go field: InstanceDetailInfo.Title<br/>type: string | `approval/admin/instance_detail.go:28` |
| `InstanceEventBase.Applicant` | `applicant` | Go field: InstanceEventBase.Applicant<br/>type: UserInfo | `approval/events_base.go:20` |
| `InstanceEventBase.BusinessRef` | `businessRef` | Go field: InstanceEventBase.BusinessRef<br/>type: *string | `approval/events_base.go:19` |
| `InstanceEventBase.FlowCode` | `flowCode` | Go field: InstanceEventBase.FlowCode<br/>type: string | `approval/events_base.go:18` |
| `InstanceEventBase.FlowID` | `flowId` | Go field: InstanceEventBase.FlowID<br/>type: string | `approval/events_base.go:17` |
| `InstanceEventBase.InstanceID` | `instanceId` | Go field: InstanceEventBase.InstanceID<br/>type: string | `approval/events_base.go:13` |
| `InstanceEventBase.InstanceNo` | `instanceNo` | Go field: InstanceEventBase.InstanceNo<br/>type: string | `approval/events_base.go:14` |
| `InstanceEventBase.OccurredTime` | `occurredTime` | Go field: InstanceEventBase.OccurredTime<br/>type: timex.DateTime | `approval/events_base.go:21` |
| `InstanceEventBase.TenantID` | `tenantId` | Go field: InstanceEventBase.TenantID<br/>type: string | `approval/events_base.go:15` |
| `InstanceEventBase.Title` | `title` | Go field: InstanceEventBase.Title<br/>type: string | `approval/events_base.go:16` |
| `InstanceFlowGraph.Edges` | `edges` | Go field: InstanceFlowGraph.Edges<br/>type: []FlowGraphEdge | `approval/flow_graph_view.go:42` |
| `InstanceFlowGraph.Nodes` | `nodes` | Go field: InstanceFlowGraph.Nodes<br/>type: []FlowGraphNode | `approval/flow_graph_view.go:41` |
| `InstanceInfo.Applicant` | `applicant` | Go field: InstanceInfo.Applicant<br/>type: approval.UserInfo | `approval/my/instance_detail.go:82` |
| `InstanceInfo.BusinessRef` | `businessRef` | Go field: InstanceInfo.BusinessRef<br/>type: *string | `approval/my/instance_detail.go:86` |
| `InstanceInfo.CreatedAt` | `createdAt` | Go field: InstanceInfo.CreatedAt<br/>type: timex.DateTime | `approval/my/instance_detail.go:88` |
| `InstanceInfo.CurrentNodeID` | `currentNodeId` | Go field: InstanceInfo.CurrentNodeID<br/>type: *string | `approval/my/instance_detail.go:84` |
| `InstanceInfo.CurrentNodeName` | `currentNodeName` | Go field: InstanceInfo.CurrentNodeName<br/>type: *string | `approval/my/instance_detail.go:85` |
| `InstanceInfo.FinishedAt` | `finishedAt` | Go field: InstanceInfo.FinishedAt<br/>type: *timex.DateTime | `approval/my/instance_detail.go:89` |
| `InstanceInfo.FlowIcon` | `flowIcon` | Go field: InstanceInfo.FlowIcon<br/>type: *string | `approval/my/instance_detail.go:80` |
| `InstanceInfo.FlowName` | `flowName` | Go field: InstanceInfo.FlowName<br/>type: string | `approval/my/instance_detail.go:79` |
| `InstanceInfo.FormData` | `formData` | Go field: InstanceInfo.FormData<br/>type: map[string]any | `approval/my/instance_detail.go:87` |
| `InstanceInfo.InstanceID` | `instanceId` | Go field: InstanceInfo.InstanceID<br/>type: string | `approval/my/instance_detail.go:76` |
| `InstanceInfo.InstanceNo` | `instanceNo` | Go field: InstanceInfo.InstanceNo<br/>type: string | `approval/my/instance_detail.go:77` |
| `InstanceInfo.Labels` | `labels` | Go field: InstanceInfo.Labels<br/>type: map[string]string | `approval/my/instance_detail.go:81` |
| `InstanceInfo.Status` | `status` | Go field: InstanceInfo.Status<br/>type: string | `approval/my/instance_detail.go:83` |
| `InstanceInfo.Title` | `title` | Go field: InstanceInfo.Title<br/>type: string | `approval/my/instance_detail.go:78` |
| `InstanceResubmittedEvent.Operator` | `operator` | Go field: InstanceResubmittedEvent.Operator<br/>type: UserInfo | `approval/events_instance.go:133` |
| `InstanceReturnedEvent.FromNodeID` | `fromNodeId` | Go field: InstanceReturnedEvent.FromNodeID<br/>type: string | `approval/events_instance.go:105` |
| `InstanceReturnedEvent.FromNodeName` | `fromNodeName` | Go field: InstanceReturnedEvent.FromNodeName<br/>type: string | `approval/events_instance.go:106` |
| `InstanceReturnedEvent.Operator` | `operator` | Go field: InstanceReturnedEvent.Operator<br/>type: UserInfo | `approval/events_instance.go:109` |
| `InstanceReturnedEvent.Opinion` | `opinion` | Go field: InstanceReturnedEvent.Opinion<br/>type: *string | `approval/events_instance.go:112` |
| `InstanceReturnedEvent.ToNodeID` | `toNodeId` | Go field: InstanceReturnedEvent.ToNodeID<br/>type: string | `approval/events_instance.go:107` |
| `InstanceReturnedEvent.ToNodeName` | `toNodeName` | Go field: InstanceReturnedEvent.ToNodeName<br/>type: string | `approval/events_instance.go:108` |
| `InstanceRolledBackEvent.FromNodeID` | `fromNodeId` | Go field: InstanceRolledBackEvent.FromNodeID<br/>type: string | `approval/events_instance.go:76` |
| `InstanceRolledBackEvent.FromNodeName` | `fromNodeName` | Go field: InstanceRolledBackEvent.FromNodeName<br/>type: string | `approval/events_instance.go:77` |
| `InstanceRolledBackEvent.Operator` | `operator` | Go field: InstanceRolledBackEvent.Operator<br/>type: UserInfo | `approval/events_instance.go:80` |
| `InstanceRolledBackEvent.Opinion` | `opinion` | Go field: InstanceRolledBackEvent.Opinion<br/>type: *string | `approval/events_instance.go:83` |
| `InstanceRolledBackEvent.ToNodeID` | `toNodeId` | Go field: InstanceRolledBackEvent.ToNodeID<br/>type: string | `approval/events_instance.go:78` |
| `InstanceRolledBackEvent.ToNodeName` | `toNodeName` | Go field: InstanceRolledBackEvent.ToNodeName<br/>type: string | `approval/events_instance.go:79` |
| `InstanceWithdrawnEvent.Operator` | `operator` | Go field: InstanceWithdrawnEvent.Operator<br/>type: UserInfo | `approval/events_instance.go:54` |
| `InstanceWithdrawnEvent.Reason` | `reason` | Go field: InstanceWithdrawnEvent.Reason<br/>type: *string | `approval/events_instance.go:57` |
| `IntegrationStatsInfo.Enabled` | `enabled` | Go field: IntegrationStatsInfo.Enabled<br/>type: bool | `monitor/integration_stats.go:10` |
| `IntegrationStatsInfo.Stats` | `stats` | Go field: IntegrationStatsInfo.Stats<br/>type: []integration.InvocationStats | `monitor/integration_stats.go:11` |
| `InterfaceInfo.Addrs` | `addrs` | Go field: InterfaceInfo.Addrs<br/>type: []string | `monitor/service.go:214` |
| `InterfaceInfo.Flags` | `flags` | Go field: InterfaceInfo.Flags<br/>type: []string | `monitor/service.go:213` |
| `InterfaceInfo.HardwareAddr` | `hardwareAddr` | Go field: InterfaceInfo.HardwareAddr<br/>type: string | `monitor/service.go:212` |
| `InterfaceInfo.Index` | `index` | Go field: InterfaceInfo.Index<br/>type: int | `monitor/service.go:209` |
| `InterfaceInfo.MTU` | `mtu` | Go field: InterfaceInfo.MTU<br/>type: int | `monitor/service.go:210` |
| `InterfaceInfo.Name` | `name` | Go field: InterfaceInfo.Name<br/>type: string | `monitor/service.go:211` |
| `InvocationLog.ContractCode` | `contractCode` | Go field: InvocationLog.ContractCode<br/>type: string | `integration/models.go:266` |
| `InvocationLog.Direction` | `direction` | Go field: InvocationLog.Direction<br/>type: Direction | `integration/models.go:269` |
| `InvocationLog.DurationMs` | `durationMs` | Go field: InvocationLog.DurationMs<br/>type: int64 | `integration/models.go:272` |
| `InvocationLog.Error` | `error` | Go field: InvocationLog.Error<br/>type: *string | `integration/models.go:276` |
| `InvocationLog.FailureKind` | `failureKind` | Go field: InvocationLog.FailureKind<br/>type: FailureKind | `integration/models.go:271` |
| `InvocationLog.HTTPTrace` | `httpTrace` | Go field: InvocationLog.HTTPTrace<br/>type: []HTTPExchange | `integration/models.go:275` |
| `InvocationLog.Input` | `input` | Go field: InvocationLog.Input<br/>type: json.RawMessage | `integration/models.go:273` |
| `InvocationLog.Output` | `output` | Go field: InvocationLog.Output<br/>type: json.RawMessage | `integration/models.go:274` |
| `InvocationLog.RequestID` | `requestId` | Go field: InvocationLog.RequestID<br/>type: string | `integration/models.go:277` |
| `InvocationLog.SystemCode` | `systemCode` | Go field: InvocationLog.SystemCode<br/>type: string | `integration/models.go:265` |
| `InvocationStats.AvgDurationMs` | `avgDurationMs` | Go field: InvocationStats.AvgDurationMs<br/>type: int64 | `integration/stats.go:26` |
| `InvocationStats.Calls` | `calls` | Go field: InvocationStats.Calls<br/>type: int64 | `integration/stats.go:23` |
| `InvocationStats.Contract` | `contract` | Go field: InvocationStats.Contract<br/>type: string | `integration/stats.go:21` |
| `InvocationStats.Direction` | `direction` | Go field: InvocationStats.Direction<br/>type: Direction | `integration/stats.go:22` |
| `InvocationStats.Failures` | `failures` | Go field: InvocationStats.Failures<br/>type: map[FailureKind]int64 | `integration/stats.go:25` |
| `InvocationStats.LastError` | `lastError` | Go field: InvocationStats.LastError<br/>type: string | `integration/stats.go:28` |
| `InvocationStats.LastErrorAt` | `lastErrorAt` | Go field: InvocationStats.LastErrorAt<br/>type: time.Time | `integration/stats.go:29` |
| `InvocationStats.MaxDurationMs` | `maxDurationMs` | Go field: InvocationStats.MaxDurationMs<br/>type: int64 | `integration/stats.go:27` |
| `InvocationStats.Successes` | `successes` | Go field: InvocationStats.Successes<br/>type: int64 | `integration/stats.go:24` |
| `InvocationStats.System` | `system` | Go field: InvocationStats.System<br/>type: string | `integration/stats.go:20` |
| `ListCodesParams.CodeSet` | `codeSet` | Go field: ListCodesParams.CodeSet<br/>type: string<br/>validate: "required" | `internal/integration/resource/code_set_catalog.go:15` |
| `ListPartsParams.ClaimID` | `claimId` | Go field: ListPartsParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:519` |
| `ListPartsResult.Parts` | `parts` | Go field: ListPartsResult.Parts<br/>type: []ListedPart | `internal/storage/resource.go:533` |
| `ListedPart.PartNumber` | `partNumber` | Go field: ListedPart.PartNumber<br/>type: int | `internal/storage/resource.go:526` |
| `ListedPart.Size` | `size` | Go field: ListedPart.Size<br/>type: int64 | `internal/storage/resource.go:527` |
| `LoadInfo.Load1` | `load1` | Go field: LoadInfo.Load1<br/>type: float64 | `monitor/service.go:262` |
| `LoadInfo.Load15` | `load15` | Go field: LoadInfo.Load15<br/>type: float64 | `monitor/service.go:264` |
| `LoadInfo.Load5` | `load5` | Go field: LoadInfo.Load5<br/>type: float64 | `monitor/service.go:263` |
| `LogSearch.ContractCode` | `contractCode` | Go field: LogSearch.ContractCode<br/>search: "eq,column=contract_code"<br/>type: string | `internal/integration/resource/log.go:14` |
| `LogSearch.Direction` | `direction` | Go field: LogSearch.Direction<br/>search: "eq,column=direction"<br/>type: integration.Direction | `internal/integration/resource/log.go:15` |
| `LogSearch.FailureKind` | `failureKind` | Go field: LogSearch.FailureKind<br/>search: "eq,column=failure_kind"<br/>type: string | `internal/integration/resource/log.go:16` |
| `LogSearch.RequestID` | `requestId` | Go field: LogSearch.RequestID<br/>search: "eq,column=request_id"<br/>type: string | `internal/integration/resource/log.go:17` |
| `LogSearch.SystemCode` | `systemCode` | Go field: LogSearch.SystemCode<br/>search: "eq,column=system_code"<br/>type: string | `internal/integration/resource/log.go:13` |
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
| `LoginParams.Credentials` | `credentials` | Go field: LoginParams.Credentials<br/>type: any<br/>validate: "required" | `internal/security/auth_resource.go:117` |
| `LoginParams.Principal` | `principal` | Go field: LoginParams.Principal<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:116` |
| `LoginParams.Type` | `type` | Go field: LoginParams.Type<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:115` |
| `LoginResult.Challenge` | `challenge` | Go field: LoginResult.Challenge<br/>type: *LoginChallenge | `security/challenge.go:18` |
| `LoginResult.ChallengeToken` | `challengeToken` | Go field: LoginResult.ChallengeToken<br/>type: string | `security/challenge.go:17` |
| `LoginResult.Tokens` | `tokens` | Go field: LoginResult.Tokens<br/>type: *AuthTokens | `security/challenge.go:16` |
| `MarkCCReadParams.InstanceID` | `instanceId` | Go field: MarkCCReadParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:375` |
| `MemoryInfo.Swap` | `swap` | Go field: MemoryInfo.Swap<br/>type: *SwapMemory | `monitor/service.go:85` |
| `MemoryInfo.Virtual` | `virtual` | Go field: MemoryInfo.Virtual<br/>type: *VirtualMemory | `monitor/service.go:84` |
| `MemorySummary.Total` | `total` | Go field: MemorySummary.Total<br/>type: uint64 | `monitor/service.go:77` |
| `MemorySummary.Used` | `used` | Go field: MemorySummary.Used<br/>type: uint64 | `monitor/service.go:78` |
| `MemorySummary.UsedPercent` | `usedPercent` | Go field: MemorySummary.UsedPercent<br/>type: float64 | `monitor/service.go:79` |
| `Message.ID` | `id` | Go field: Message.ID<br/>type: string | `push/message.go:13` |
| `Message.Payload` | `payload` | Go field: Message.Payload<br/>type: any | `push/message.go:15` |
| `Message.Time` | `time` | Go field: Message.Time<br/>type: time.Time | `push/message.go:16` |
| `Message.Type` | `type` | Go field: Message.Type<br/>type: string | `push/message.go:14` |
| `Metrics.AvgCompletionSeconds` | `avgCompletionSeconds` | Go field: Metrics.AvgCompletionSeconds<br/>type: float64 | `approval/admin/metrics.go:26` |
| `Metrics.BusinessProjectionCounts` | `businessProjectionCounts` | Go field: Metrics.BusinessProjectionCounts<br/>type: map[string]int | `approval/admin/metrics.go:32` |
| `Metrics.CapturedAt` | `capturedAt` | Go field: Metrics.CapturedAt<br/>type: timex.DateTime | `approval/admin/metrics.go:14` |
| `Metrics.InstanceCounts` | `instanceCounts` | Go field: Metrics.InstanceCounts<br/>type: map[string]int | `approval/admin/metrics.go:17` |
| `Metrics.PendingBindingFailures` | `pendingBindingFailures` | Go field: Metrics.PendingBindingFailures<br/>type: int | `approval/admin/metrics.go:29` |
| `Metrics.PendingBusinessProjections` | `pendingBusinessProjections` | Go field: Metrics.PendingBusinessProjections<br/>type: int | `approval/admin/metrics.go:35` |
| `Metrics.TaskCounts` | `taskCounts` | Go field: Metrics.TaskCounts<br/>type: map[string]int | `approval/admin/metrics.go:19` |
| `Metrics.TenantID` | `tenantId` | Go field: Metrics.TenantID<br/>type: string | `approval/admin/metrics.go:12` |
| `Metrics.TimeoutTaskCount` | `timeoutTaskCount` | Go field: Metrics.TimeoutTaskCount<br/>type: int | `approval/admin/metrics.go:21` |
| `Model.ID` | `id` | Go field: Model.ID<br/>type: string | `internal/orm/model.go:7` |
| `NetIOCounter.BytesRecv` | `bytesRecv` | Go field: NetIOCounter.BytesRecv<br/>type: uint64 | `monitor/service.go:221` |
| `NetIOCounter.BytesSent` | `bytesSent` | Go field: NetIOCounter.BytesSent<br/>type: uint64 | `monitor/service.go:220` |
| `NetIOCounter.DroppedIn` | `droppedIn` | Go field: NetIOCounter.DroppedIn<br/>type: uint64 | `monitor/service.go:226` |
| `NetIOCounter.DroppedOut` | `droppedOut` | Go field: NetIOCounter.DroppedOut<br/>type: uint64 | `monitor/service.go:227` |
| `NetIOCounter.ErrorsIn` | `errorsIn` | Go field: NetIOCounter.ErrorsIn<br/>type: uint64 | `monitor/service.go:224` |
| `NetIOCounter.ErrorsOut` | `errorsOut` | Go field: NetIOCounter.ErrorsOut<br/>type: uint64 | `monitor/service.go:225` |
| `NetIOCounter.FIFOIn` | `fifoIn` | Go field: NetIOCounter.FIFOIn<br/>type: uint64 | `monitor/service.go:228` |
| `NetIOCounter.FIFOOut` | `fifoOut` | Go field: NetIOCounter.FIFOOut<br/>type: uint64 | `monitor/service.go:229` |
| `NetIOCounter.Name` | `name` | Go field: NetIOCounter.Name<br/>type: string | `monitor/service.go:219` |
| `NetIOCounter.PacketsRecv` | `packetsRecv` | Go field: NetIOCounter.PacketsRecv<br/>type: uint64 | `monitor/service.go:223` |
| `NetIOCounter.PacketsSent` | `packetsSent` | Go field: NetIOCounter.PacketsSent<br/>type: uint64 | `monitor/service.go:222` |
| `NetworkInfo.IOCounters` | `ioCounters` | Go field: NetworkInfo.IOCounters<br/>type: map[string]*NetIOCounter | `monitor/service.go:204` |
| `NetworkInfo.Interfaces` | `interfaces` | Go field: NetworkInfo.Interfaces<br/>type: []*InterfaceInfo | `monitor/service.go:203` |
| `NetworkSummary.BytesRecv` | `bytesRecv` | Go field: NetworkSummary.BytesRecv<br/>type: uint64 | `monitor/service.go:196` |
| `NetworkSummary.BytesSent` | `bytesSent` | Go field: NetworkSummary.BytesSent<br/>type: uint64 | `monitor/service.go:195` |
| `NetworkSummary.Interfaces` | `interfaces` | Go field: NetworkSummary.Interfaces<br/>type: int | `monitor/service.go:194` |
| `NetworkSummary.PacketsRecv` | `packetsRecv` | Go field: NetworkSummary.PacketsRecv<br/>type: uint64 | `monitor/service.go:198` |
| `NetworkSummary.PacketsSent` | `packetsSent` | Go field: NetworkSummary.PacketsSent<br/>type: uint64 | `monitor/service.go:197` |
| `NodeAutoPassedEvent.NodeID` | `nodeId` | Go field: NodeAutoPassedEvent.NodeID<br/>type: string | `approval/events_node.go:10` |
| `NodeAutoPassedEvent.NodeName` | `nodeName` | Go field: NodeAutoPassedEvent.NodeName<br/>type: string | `approval/events_node.go:11` |
| `NodeAutoPassedEvent.Reason` | `reason` | Go field: NodeAutoPassedEvent.Reason<br/>type: string | `approval/events_node.go:12` |
| `NodeDefinition.Data` | `data` | Go field: NodeDefinition.Data<br/>type: json.RawMessage | `approval/flow_definition.go:41` |
| `NodeDefinition.ID` | `id` | Go field: NodeDefinition.ID<br/>type: string | `approval/flow_definition.go:38` |
| `NodeDefinition.Kind` | `kind` | Go field: NodeDefinition.Kind<br/>type: NodeKind | `approval/flow_definition.go:39` |
| `NodeDefinition.Position` | `position` | Go field: NodeDefinition.Position<br/>type: Position | `approval/flow_definition.go:40` |
| `NodeParticipant.ActionTime` | `actionTime` | Go field: NodeParticipant.ActionTime<br/>type: *timex.DateTime | `approval/node_view.go:28` |
| `NodeParticipant.Attachments` | `attachments` | Go field: NodeParticipant.Attachments<br/>type: []string | `approval/node_view.go:27` |
| `NodeParticipant.Deadline` | `deadline` | Go field: NodeParticipant.Deadline<br/>type: *timex.DateTime | `approval/node_view.go:24` |
| `NodeParticipant.Delegator` | `delegator` | Go field: NodeParticipant.Delegator<br/>type: *UserInfo | `approval/node_view.go:22` |
| `NodeParticipant.IsTimeout` | `isTimeout` | Go field: NodeParticipant.IsTimeout<br/>type: bool | `approval/node_view.go:25` |
| `NodeParticipant.Opinion` | `opinion` | Go field: NodeParticipant.Opinion<br/>type: *string | `approval/node_view.go:26` |
| `NodeParticipant.Status` | `status` | Go field: NodeParticipant.Status<br/>type: string | `approval/node_view.go:23` |
| `NodeParticipant.TaskID` | `taskId` | Go field: NodeParticipant.TaskID<br/>type: string | `approval/node_view.go:20` |
| `NodeParticipant.TransferTo` | `transferTo` | Go field: NodeParticipant.TransferTo<br/>type: *UserInfo | `approval/node_view.go:29` |
| `NodeParticipant.User` | `user` | Go field: NodeParticipant.User<br/>type: UserInfo | `approval/node_view.go:21` |
| `NodeVisit.FinishedAt` | `finishedAt` | Go field: NodeVisit.FinishedAt<br/>type: *timex.DateTime | `approval/models.go:373` |
| `NodeVisit.InstanceID` | `instanceId` | Go field: NodeVisit.InstanceID<br/>type: string | `approval/models.go:369` |
| `NodeVisit.NodeID` | `nodeId` | Go field: NodeVisit.NodeID<br/>type: string | `approval/models.go:370` |
| `NodeVisit.Sequence` | `sequence` | Go field: NodeVisit.Sequence<br/>type: int | `approval/models.go:371` |
| `NodeVisit.Status` | `status` | Go field: NodeVisit.Status<br/>type: NodeVisitStatus | `approval/models.go:372` |
| `NodeVisit.TenantID` | `tenantId` | Go field: NodeVisit.TenantID<br/>type: string | `approval/models.go:368` |
| `OTPChallengeData.Destination` | `destination` | Go field: OTPChallengeData.Destination<br/>type: string | `security/otp.go:45` |
| `OTPChallengeData.Meta` | `meta` | Go field: OTPChallengeData.Meta<br/>type: map[string]any | `security/otp.go:46` |
| `ObjectInfo.Bucket` | `bucket` | Go field: ObjectInfo.Bucket<br/>type: string | `storage/service.go:125` |
| `ObjectInfo.ContentType` | `contentType` | Go field: ObjectInfo.ContentType<br/>type: string | `storage/service.go:133` |
| `ObjectInfo.ETag` | `eTag` | Go field: ObjectInfo.ETag<br/>type: string | `storage/service.go:129` |
| `ObjectInfo.Key` | `key` | Go field: ObjectInfo.Key<br/>type: string | `storage/service.go:127` |
| `ObjectInfo.LastModified` | `lastModified` | Go field: ObjectInfo.LastModified<br/>type: time.Time | `storage/service.go:135` |
| `ObjectInfo.Metadata` | `metadata` | Go field: ObjectInfo.Metadata<br/>type: map[string]string | `storage/service.go:141` |
| `ObjectInfo.Size` | `size` | Go field: ObjectInfo.Size<br/>type: int64 | `storage/service.go:131` |
| `OutboundAuthConfig.Params` | `params` | Go field: OutboundAuthConfig.Params<br/>type: map[string]string | `integration/models.go:149` |
| `OutboundAuthConfig.Scheme` | `scheme` | Go field: OutboundAuthConfig.Scheme<br/>type: string | `integration/models.go:148` |
| `OutboundAuthConfig.Script` | `script` | Go field: OutboundAuthConfig.Script<br/>type: string | `integration/models.go:153` |
| `OutboundEnvelopeConfig.Request` | `request` | Go field: OutboundEnvelopeConfig.Request<br/>type: string | `integration/models.go:168` |
| `OutboundEnvelopeConfig.Response` | `response` | Go field: OutboundEnvelopeConfig.Response<br/>type: string | `integration/models.go:175` |
| `Page.Items` | `items` | Go field: Page.Items<br/>type: []T | `page/page.go:50` |
| `Page.Page` | `page` | Go field: Page.Page<br/>type: int | `page/page.go:47` |
| `Page.Size` | `size` | Go field: Page.Size<br/>type: int | `page/page.go:48` |
| `Page.Total` | `total` | Go field: Page.Total<br/>type: int64 | `page/page.go:49` |
| `Pageable.Page` | `page` | Go field: Pageable.Page<br/>type: int | `page/page.go:14` |
| `Pageable.Size` | `size` | Go field: Pageable.Size<br/>type: int | `page/page.go:15` |
| `PartitionInfo.Device` | `device` | Go field: PartitionInfo.Device<br/>type: string | `monitor/service.go:160` |
| `PartitionInfo.FSType` | `fsType` | Go field: PartitionInfo.FSType<br/>type: string | `monitor/service.go:162` |
| `PartitionInfo.Free` | `free` | Go field: PartitionInfo.Free<br/>type: uint64 | `monitor/service.go:165` |
| `PartitionInfo.INodesFree` | `iNodesFree` | Go field: PartitionInfo.INodesFree<br/>type: uint64 | `monitor/service.go:170` |
| `PartitionInfo.INodesTotal` | `iNodesTotal` | Go field: PartitionInfo.INodesTotal<br/>type: uint64 | `monitor/service.go:168` |
| `PartitionInfo.INodesUsed` | `iNodesUsed` | Go field: PartitionInfo.INodesUsed<br/>type: uint64 | `monitor/service.go:169` |
| `PartitionInfo.INodesUsedPercent` | `iNodesUsedPercent` | Go field: PartitionInfo.INodesUsedPercent<br/>type: float64 | `monitor/service.go:171` |
| `PartitionInfo.MountPoint` | `mountPoint` | Go field: PartitionInfo.MountPoint<br/>type: string | `monitor/service.go:161` |
| `PartitionInfo.Options` | `options` | Go field: PartitionInfo.Options<br/>type: []string | `monitor/service.go:163` |
| `PartitionInfo.Total` | `total` | Go field: PartitionInfo.Total<br/>type: uint64 | `monitor/service.go:164` |
| `PartitionInfo.Used` | `used` | Go field: PartitionInfo.Used<br/>type: uint64 | `monitor/service.go:166` |
| `PartitionInfo.UsedPercent` | `usedPercent` | Go field: PartitionInfo.UsedPercent<br/>type: float64 | `monitor/service.go:167` |
| `PasswordChangeChallengeData.Meta` | `meta` | Go field: PasswordChangeChallengeData.Meta<br/>type: map[string]any | `security/password_change.go:19` |
| `PasswordChangeChallengeData.Reason` | `reason` | Go field: PasswordChangeChallengeData.Reason<br/>type: string | `security/password_change.go:18` |
| `PendingCounts.PendingTaskCount` | `pendingTaskCount` | Go field: PendingCounts.PendingTaskCount<br/>type: int | `approval/my/pending_counts.go:5` |
| `PendingCounts.UnreadCCCount` | `unreadCcCount` | Go field: PendingCounts.UnreadCCCount<br/>type: int | `approval/my/pending_counts.go:6` |
| `PendingDelete.Attempts` | `attempts` | Go field: PendingDelete.Attempts<br/>type: int | `internal/storage/store/delete.go:34` |
| `PendingDelete.CreatedAt` | `createdAt` | Go field: PendingDelete.CreatedAt<br/>type: timex.DateTime | `internal/storage/store/delete.go:36` |
| `PendingDelete.ID` | `id` | Go field: PendingDelete.ID<br/>type: string | `internal/storage/store/delete.go:30` |
| `PendingDelete.Key` | `key` | Go field: PendingDelete.Key<br/>type: string | `internal/storage/store/delete.go:31` |
| `PendingDelete.NextAttemptAt` | `nextAttemptAt` | Go field: PendingDelete.NextAttemptAt<br/>type: timex.DateTime | `internal/storage/store/delete.go:35` |
| `PendingDelete.Reason` | `reason` | Go field: PendingDelete.Reason<br/>type: storage.DeleteReason | `internal/storage/store/delete.go:33` |
| `PendingDelete.UploadID` | `uploadId` | Go field: PendingDelete.UploadID<br/>type: string | `internal/storage/store/delete.go:32` |
| `PendingTask.Applicant` | `applicant` | Go field: PendingTask.Applicant<br/>type: approval.UserInfo | `approval/my/pending_tasks.go:16` |
| `PendingTask.CreatedAt` | `createdAt` | Go field: PendingTask.CreatedAt<br/>type: timex.DateTime | `approval/my/pending_tasks.go:18` |
| `PendingTask.Deadline` | `deadline` | Go field: PendingTask.Deadline<br/>type: *timex.DateTime | `approval/my/pending_tasks.go:19` |
| `PendingTask.FlowIcon` | `flowIcon` | Go field: PendingTask.FlowIcon<br/>type: *string | `approval/my/pending_tasks.go:15` |
| `PendingTask.FlowName` | `flowName` | Go field: PendingTask.FlowName<br/>type: string | `approval/my/pending_tasks.go:14` |
| `PendingTask.InstanceID` | `instanceId` | Go field: PendingTask.InstanceID<br/>type: string | `approval/my/pending_tasks.go:11` |
| `PendingTask.InstanceNo` | `instanceNo` | Go field: PendingTask.InstanceNo<br/>type: string | `approval/my/pending_tasks.go:13` |
| `PendingTask.InstanceTitle` | `instanceTitle` | Go field: PendingTask.InstanceTitle<br/>type: string | `approval/my/pending_tasks.go:12` |
| `PendingTask.IsTimeout` | `isTimeout` | Go field: PendingTask.IsTimeout<br/>type: bool | `approval/my/pending_tasks.go:20` |
| `PendingTask.NodeName` | `nodeName` | Go field: PendingTask.NodeName<br/>type: string | `approval/my/pending_tasks.go:17` |
| `PendingTask.TaskID` | `taskId` | Go field: PendingTask.TaskID<br/>type: string | `approval/my/pending_tasks.go:10` |
| `Position.X` | `x` | Go field: Position.X<br/>type: float64 | `approval/flow_definition.go:28` |
| `Position.Y` | `y` | Go field: Position.Y<br/>type: float64 | `approval/flow_definition.go:29` |
| `PreviewFiresParams.EndsAtUnixMs` | `endsAtUnixMs` | Go field: PreviewFiresParams.EndsAtUnixMs<br/>type: *int64 | `internal/cron/store/schedule_resource.go:133` |
| `PreviewFiresParams.StartsAtUnixMs` | `startsAtUnixMs` | Go field: PreviewFiresParams.StartsAtUnixMs<br/>type: *int64 | `internal/cron/store/schedule_resource.go:132` |
| `PreviewFiresParams.Trigger` | `trigger` | Go field: PreviewFiresParams.Trigger<br/>type: TriggerParams | `internal/cron/store/schedule_resource.go:131` |
| `PrimaryKey.Columns` | `columns` | Go field: PrimaryKey.Columns<br/>type: []string | `schema/service.go:39` |
| `PrimaryKey.Name` | `name` | Go field: PrimaryKey.Name<br/>type: string | `schema/service.go:38` |
| `Principal.Details` | `details` | Go field: Principal.Details<br/>type: any | `security/principal.go:82` |
| `Principal.ID` | `id` | Go field: Principal.ID<br/>type: string | `security/principal.go:76` |
| `Principal.Name` | `name` | Go field: Principal.Name<br/>type: string | `security/principal.go:78` |
| `Principal.Roles` | `roles` | Go field: Principal.Roles<br/>type: []string | `security/principal.go:80` |
| `Principal.Type` | `type` | Go field: Principal.Type<br/>type: PrincipalType | `security/principal.go:74` |
| `ProcessInfo.CPUPercent` | `cpuPercent` | Go field: ProcessInfo.CPUPercent<br/>type: float64 | `monitor/service.go:253` |
| `ProcessInfo.CWD` | `cwd` | Go field: ProcessInfo.CWD<br/>type: string | `monitor/service.go:247` |
| `ProcessInfo.CommandLine` | `commandLine` | Go field: ProcessInfo.CommandLine<br/>type: string | `monitor/service.go:246` |
| `ProcessInfo.CreateTime` | `createTime` | Go field: ProcessInfo.CreateTime<br/>type: int64 | `monitor/service.go:250` |
| `ProcessInfo.Exe` | `exe` | Go field: ProcessInfo.Exe<br/>type: string | `monitor/service.go:245` |
| `ProcessInfo.MemoryPercent` | `memoryPercent` | Go field: ProcessInfo.MemoryPercent<br/>type: float32 | `monitor/service.go:254` |
| `ProcessInfo.MemoryRSS` | `memoryRss` | Go field: ProcessInfo.MemoryRSS<br/>type: uint64 | `monitor/service.go:255` |
| `ProcessInfo.MemorySwap` | `memorySwap` | Go field: ProcessInfo.MemorySwap<br/>type: uint64 | `monitor/service.go:257` |
| `ProcessInfo.MemoryVMS` | `memoryVms` | Go field: ProcessInfo.MemoryVMS<br/>type: uint64 | `monitor/service.go:256` |
| `ProcessInfo.Name` | `name` | Go field: ProcessInfo.Name<br/>type: string | `monitor/service.go:244` |
| `ProcessInfo.NumFDs` | `numFds` | Go field: ProcessInfo.NumFDs<br/>type: int32 | `monitor/service.go:252` |
| `ProcessInfo.NumThreads` | `numThreads` | Go field: ProcessInfo.NumThreads<br/>type: int32 | `monitor/service.go:251` |
| `ProcessInfo.PID` | `pid` | Go field: ProcessInfo.PID<br/>type: int32 | `monitor/service.go:242` |
| `ProcessInfo.ParentPID` | `parentPid` | Go field: ProcessInfo.ParentPID<br/>type: int32 | `monitor/service.go:243` |
| `ProcessInfo.Status` | `status` | Go field: ProcessInfo.Status<br/>type: string | `monitor/service.go:248` |
| `ProcessInfo.Username` | `username` | Go field: ProcessInfo.Username<br/>type: string | `monitor/service.go:249` |
| `ProcessSummary.CPUPercent` | `cpuPercent` | Go field: ProcessSummary.CPUPercent<br/>type: float64 | `monitor/service.go:236` |
| `ProcessSummary.MemoryPercent` | `memoryPercent` | Go field: ProcessSummary.MemoryPercent<br/>type: float32 | `monitor/service.go:237` |
| `ProcessSummary.Name` | `name` | Go field: ProcessSummary.Name<br/>type: string | `monitor/service.go:235` |
| `ProcessSummary.PID` | `pid` | Go field: ProcessSummary.PID<br/>type: int32 | `monitor/service.go:234` |
| `ProcessTaskParams.Action` | `action` | Go field: ProcessTaskParams.Action<br/>type: string<br/>validate: "required,oneof=approve reject transfer rollback handle" | `internal/approval/resource/instance.go:181` |
| `ProcessTaskParams.Attachments` | `attachments` | Go field: ProcessTaskParams.Attachments<br/>type: []string<br/>validate: "max=20,dive,max=512" | `internal/approval/resource/instance.go:184` |
| `ProcessTaskParams.FormData` | `formData` | Go field: ProcessTaskParams.FormData<br/>type: map[string]any | `internal/approval/resource/instance.go:183` |
| `ProcessTaskParams.Opinion` | `opinion` | Go field: ProcessTaskParams.Opinion<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/instance.go:182` |
| `ProcessTaskParams.TargetNodeID` | `targetNodeId` | Go field: ProcessTaskParams.TargetNodeID<br/>type: string | `internal/approval/resource/instance.go:186` |
| `ProcessTaskParams.TaskID` | `taskId` | Go field: ProcessTaskParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:180` |
| `ProcessTaskParams.TransferToID` | `transferToId` | Go field: ProcessTaskParams.TransferToID<br/>type: string | `internal/approval/resource/instance.go:185` |
| `PublishVersionParams.VersionID` | `versionId` | Go field: PublishVersionParams.VersionID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:161` |
| `QueryArgs.Params` | `params` | Go field: QueryArgs.Params<br/>jsonschema: "description=Parameters for the SQL query placeholders"<br/>type: []any | `internal/mcp/tools/query.go:17` |
| `QueryArgs.SQL` | `sql` | Go field: QueryArgs.SQL<br/>jsonschema: "required,description=The SQL query with placeholders (?) for parameters"<br/>type: string | `internal/mcp/tools/query.go:16` |
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
| `RefreshParams.RefreshToken` | `refreshToken` | Go field: RefreshParams.RefreshToken<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:190` |
| `RemovableAssignee.Assignee` | `assignee` | Go field: RemovableAssignee.Assignee<br/>type: approval.UserInfo | `approval/my/instance_detail.go:66` |
| `RemovableAssignee.Status` | `status` | Go field: RemovableAssignee.Status<br/>type: string | `approval/my/instance_detail.go:68` |
| `RemovableAssignee.TaskID` | `taskId` | Go field: RemovableAssignee.TaskID<br/>type: string | `approval/my/instance_detail.go:65` |
| `RemoveAssigneeParams.TaskID` | `taskId` | Go field: RemoveAssigneeParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:429` |
| `Request.Meta` | `meta` | Go field: Request.Meta<br/>type: Meta | `api/request.go:95` |
| `Request.Params` | `params` | Go field: Request.Params<br/>type: Params | `api/request.go:94` |
| `ResolveChallengeParams.ChallengeToken` | `challengeToken` | Go field: ResolveChallengeParams.ChallengeToken<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:249` |
| `ResolveChallengeParams.Response` | `response` | Go field: ResolveChallengeParams.Response<br/>type: any<br/>validate: "required" | `internal/security/auth_resource.go:251` |
| `ResolveChallengeParams.Type` | `type` | Go field: ResolveChallengeParams.Type<br/>type: string<br/>validate: "required" | `internal/security/auth_resource.go:250` |
| `ResubmitParams.FormData` | `formData` | Go field: ResubmitParams.FormData<br/>type: map[string]any | `internal/approval/resource/instance.go:322` |
| `ResubmitParams.InstanceID` | `instanceId` | Go field: ResubmitParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:321` |
| `Result.Code` | `code` | Go field: Result.Code<br/>type: int | `result/result.go:11` |
| `Result.Data` | `data` | Go field: Result.Data<br/>type: any | `result/result.go:13` |
| `Result.Message` | `message` | Go field: Result.Message<br/>type: string | `result/result.go:12` |
| `RetryPolicy.InitialBackoffMs` | `initialBackoffMs` | Go field: RetryPolicy.InitialBackoffMs<br/>type: int64 | `integration/models.go:199` |
| `RetryPolicy.MaxAttempts` | `maxAttempts` | Go field: RetryPolicy.MaxAttempts<br/>type: int | `integration/models.go:196` |
| `RetryPolicy.MaxBackoffMs` | `maxBackoffMs` | Go field: RetryPolicy.MaxBackoffMs<br/>type: int64 | `integration/models.go:202` |
| `RolePermissionsChangedEvent.Roles` | `roles` | Go field: RolePermissionsChangedEvent.Roles<br/>type: []string | `security/cached_role_permission_loader.go:18` |
| `RollbackTarget.Name` | `name` | Go field: RollbackTarget.Name<br/>type: string | `approval/my/instance_detail.go:60` |
| `RollbackTarget.NodeID` | `nodeId` | Go field: RollbackTarget.NodeID<br/>type: string | `approval/my/instance_detail.go:59` |
| `Route.ContractID` | `contractId` | Go field: Route.ContractID<br/>type: string | `integration/models.go:254` |
| `Route.IsEnabled` | `isEnabled` | Go field: Route.IsEnabled<br/>type: bool | `integration/models.go:256` |
| `Route.RouteKey` | `routeKey` | Go field: Route.RouteKey<br/>type: string | `integration/models.go:253` |
| `Route.SystemID` | `systemId` | Go field: Route.SystemID<br/>type: string | `integration/models.go:255` |
| `RouteDiagnostics.Findings` | `findings` | Go field: RouteDiagnostics.Findings<br/>type: []RouteFinding | `integration/diagnostics.go:50` |
| `RouteFinding.ContractCode` | `contractCode` | Go field: RouteFinding.ContractCode<br/>type: string | `integration/diagnostics.go:40` |
| `RouteFinding.ContractName` | `contractName` | Go field: RouteFinding.ContractName<br/>type: string | `integration/diagnostics.go:41` |
| `RouteFinding.Kind` | `kind` | Go field: RouteFinding.Kind<br/>type: RouteFindingKind | `integration/diagnostics.go:37` |
| `RouteFinding.RouteID` | `routeId` | Go field: RouteFinding.RouteID<br/>type: string | `integration/diagnostics.go:38` |
| `RouteFinding.RouteKey` | `routeKey` | Go field: RouteFinding.RouteKey<br/>type: string | `integration/diagnostics.go:39` |
| `RouteFinding.SystemCode` | `systemCode` | Go field: RouteFinding.SystemCode<br/>type: string | `integration/diagnostics.go:42` |
| `RouteFinding.SystemName` | `systemName` | Go field: RouteFinding.SystemName<br/>type: string | `integration/diagnostics.go:43` |
| `RouteParams.ContractID` | `contractId` | Go field: RouteParams.ContractID<br/>type: string | `internal/integration/resource/route.go:21` |
| `RouteParams.ID` | `id` | Go field: RouteParams.ID<br/>type: string | `internal/integration/resource/route.go:19` |
| `RouteParams.IsEnabled` | `isEnabled` | Go field: RouteParams.IsEnabled<br/>type: bool | `internal/integration/resource/route.go:23` |
| `RouteParams.RouteKey` | `routeKey` | Go field: RouteParams.RouteKey<br/>type: string | `internal/integration/resource/route.go:20` |
| `RouteParams.SystemID` | `systemId` | Go field: RouteParams.SystemID<br/>type: string<br/>validate: "required" | `internal/integration/resource/route.go:22` |
| `RouteSearch.ContractID` | `contractId` | Go field: RouteSearch.ContractID<br/>search: "eq,column=contract_id"<br/>type: string | `internal/integration/resource/route.go:31` |
| `RouteSearch.IsEnabled` | `isEnabled` | Go field: RouteSearch.IsEnabled<br/>search: "eq,column=is_enabled"<br/>type: *bool | `internal/integration/resource/route.go:33` |
| `RouteSearch.RouteKey` | `routeKey` | Go field: RouteSearch.RouteKey<br/>search: "contains,column=route_key"<br/>type: string | `internal/integration/resource/route.go:30` |
| `RouteSearch.SystemID` | `systemId` | Go field: RouteSearch.SystemID<br/>search: "eq,column=system_id"<br/>type: string | `internal/integration/resource/route.go:32` |
| `Run.ClaimedAtUnixMs` | `claimedAtUnixMs` | Go field: Run.ClaimedAtUnixMs<br/>type: int64 | `cron/run.go:48` |
| `Run.DurationMs` | `durationMs` | Go field: Run.DurationMs<br/>type: int64 | `cron/run.go:58` |
| `Run.Error` | `error` | Go field: Run.Error<br/>type: string | `cron/run.go:65` |
| `Run.FinishedAtUnixMs` | `finishedAtUnixMs` | Go field: Run.FinishedAtUnixMs<br/>type: *int64 | `cron/run.go:57` |
| `Run.HeartbeatAtUnixMs` | `heartbeatAtUnixMs` | Go field: Run.HeartbeatAtUnixMs<br/>type: *int64 | `cron/run.go:62` |
| `Run.JobName` | `jobName` | Go field: Run.JobName<br/>type: string | `cron/run.go:42` |
| `Run.MissedCount` | `missedCount` | Go field: Run.MissedCount<br/>type: int | `cron/run.go:68` |
| `Run.NodeID` | `nodeId` | Go field: Run.NodeID<br/>type: string | `cron/run.go:54` |
| `Run.ScheduleID` | `scheduleId` | Go field: Run.ScheduleID<br/>type: string | `cron/run.go:40` |
| `Run.ScheduleName` | `scheduleName` | Go field: Run.ScheduleName<br/>type: string | `cron/run.go:41` |
| `Run.ScheduledAtUnixMs` | `scheduledAtUnixMs` | Go field: Run.ScheduledAtUnixMs<br/>type: int64 | `cron/run.go:47` |
| `Run.StartedAtUnixMs` | `startedAtUnixMs` | Go field: Run.StartedAtUnixMs<br/>type: *int64 | `cron/run.go:56` |
| `Run.Status` | `status` | Go field: Run.Status<br/>type: RunStatus | `cron/run.go:50` |
| `RunAbandonedEvent.JobName` | `jobName` | Go field: RunAbandonedEvent.JobName<br/>type: string | `cron/events.go:49` |
| `RunAbandonedEvent.NodeID` | `nodeId` | Go field: RunAbandonedEvent.NodeID<br/>type: string | `cron/events.go:53` |
| `RunAbandonedEvent.RunID` | `runId` | Go field: RunAbandonedEvent.RunID<br/>type: string | `cron/events.go:47` |
| `RunAbandonedEvent.ScheduleName` | `scheduleName` | Go field: RunAbandonedEvent.ScheduleName<br/>type: string | `cron/events.go:48` |
| `RunAbandonedEvent.ScheduledAtUnixMs` | `scheduledAtUnixMs` | Go field: RunAbandonedEvent.ScheduledAtUnixMs<br/>type: int64 | `cron/events.go:51` |
| `RunFailedEvent.Error` | `error` | Go field: RunFailedEvent.Error<br/>type: string | `cron/events.go:26` |
| `RunFailedEvent.JobName` | `jobName` | Go field: RunFailedEvent.JobName<br/>type: string | `cron/events.go:20` |
| `RunFailedEvent.NodeID` | `nodeId` | Go field: RunFailedEvent.NodeID<br/>type: string | `cron/events.go:24` |
| `RunFailedEvent.RunID` | `runId` | Go field: RunFailedEvent.RunID<br/>type: string | `cron/events.go:18` |
| `RunFailedEvent.ScheduleName` | `scheduleName` | Go field: RunFailedEvent.ScheduleName<br/>type: string | `cron/events.go:19` |
| `RunFailedEvent.ScheduledAtUnixMs` | `scheduledAtUnixMs` | Go field: RunFailedEvent.ScheduledAtUnixMs<br/>type: int64 | `cron/events.go:22` |
| `RunSearch.ID` | `id` | Go field: RunSearch.ID<br/>search: "eq,column=id"<br/>type: string | `internal/cron/store/run_resource.go:17` |
| `RunSearch.JobName` | `jobName` | Go field: RunSearch.JobName<br/>search: "eq,column=job_name"<br/>type: string | `internal/cron/store/run_resource.go:19` |
| `RunSearch.NodeID` | `nodeId` | Go field: RunSearch.NodeID<br/>search: "eq,column=node_id"<br/>type: string | `internal/cron/store/run_resource.go:21` |
| `RunSearch.ScheduleName` | `scheduleName` | Go field: RunSearch.ScheduleName<br/>search: "eq,column=schedule_name"<br/>type: string | `internal/cron/store/run_resource.go:18` |
| `RunSearch.ScheduledAtFromUnixMs` | `scheduledAtFromUnixMs` | Go field: RunSearch.ScheduledAtFromUnixMs<br/>search: "gte,column=scheduled_at_unix_ms"<br/>type: *int64 | `internal/cron/store/run_resource.go:22` |
| `RunSearch.ScheduledAtToUnixMs` | `scheduledAtToUnixMs` | Go field: RunSearch.ScheduledAtToUnixMs<br/>search: "lte,column=scheduled_at_unix_ms"<br/>type: *int64 | `internal/cron/store/run_resource.go:23` |
| `RunSearch.Status` | `status` | Go field: RunSearch.Status<br/>search: "eq"<br/>type: string | `internal/cron/store/run_resource.go:20` |
| `Schedule.AnchorAtUnixMs` | `anchorAtUnixMs` | Go field: Schedule.AnchorAtUnixMs<br/>type: int64 | `cron/schedule.go:61` |
| `Schedule.ConcurrencyPolicy` | `concurrencyPolicy` | Go field: Schedule.ConcurrencyPolicy<br/>type: ConcurrencyPolicy | `cron/schedule.go:67` |
| `Schedule.EndsAtUnixMs` | `endsAtUnixMs` | Go field: Schedule.EndsAtUnixMs<br/>type: *int64 | `cron/schedule.go:60` |
| `Schedule.EveryMs` | `everyMs` | Go field: Schedule.EveryMs<br/>type: int64 | `cron/schedule.go:54` |
| `Schedule.Expr` | `expr` | Go field: Schedule.Expr<br/>type: string | `cron/schedule.go:52` |
| `Schedule.FireAtUnixMs` | `fireAtUnixMs` | Go field: Schedule.FireAtUnixMs<br/>type: *int64 | `cron/schedule.go:55` |
| `Schedule.IsEnabled` | `isEnabled` | Go field: Schedule.IsEnabled<br/>type: bool | `cron/schedule.go:79` |
| `Schedule.JobName` | `jobName` | Go field: Schedule.JobName<br/>type: string | `cron/schedule.go:49` |
| `Schedule.Kind` | `kind` | Go field: Schedule.Kind<br/>type: TriggerKind | `cron/schedule.go:51` |
| `Schedule.LastFireAtUnixMs` | `lastFireAtUnixMs` | Go field: Schedule.LastFireAtUnixMs<br/>type: *int64 | `cron/schedule.go:88` |
| `Schedule.MisfirePolicy` | `misfirePolicy` | Go field: Schedule.MisfirePolicy<br/>type: MisfirePolicy | `cron/schedule.go:66` |
| `Schedule.Name` | `name` | Go field: Schedule.Name<br/>type: string | `cron/schedule.go:47` |
| `Schedule.NextFireAtUnixMs` | `nextFireAtUnixMs` | Go field: Schedule.NextFireAtUnixMs<br/>type: *int64 | `cron/schedule.go:86` |
| `Schedule.Params` | `params` | Go field: Schedule.Params<br/>type: json.RawMessage | `cron/schedule.go:64` |
| `Schedule.Recover` | `recover` | Go field: Schedule.Recover<br/>type: bool | `cron/schedule.go:72` |
| `Schedule.StartsAtUnixMs` | `startsAtUnixMs` | Go field: Schedule.StartsAtUnixMs<br/>type: *int64 | `cron/schedule.go:59` |
| `Schedule.TimeoutMs` | `timeoutMs` | Go field: Schedule.TimeoutMs<br/>type: int64 | `cron/schedule.go:75` |
| `Schedule.Timezone` | `timezone` | Go field: Schedule.Timezone<br/>type: string | `cron/schedule.go:53` |
| `ScheduleDetail.NextFiresUnixMs` | `nextFiresUnixMs` | Go field: ScheduleDetail.NextFiresUnixMs<br/>type: []int64 | `internal/cron/store/schedule_resource.go:159` |
| `ScheduleDetail.Schedule` | `schedule` | Go field: ScheduleDetail.Schedule<br/>type: *cron.Schedule | `internal/cron/store/schedule_resource.go:155` |
| `ScheduleNameParams.Name` | `name` | Go field: ScheduleNameParams.Name<br/>type: string<br/>validate: "required" | `internal/cron/store/schedule_resource.go:123` |
| `ScheduleParams.ConcurrencyPolicy` | `concurrencyPolicy` | Go field: ScheduleParams.ConcurrencyPolicy<br/>type: cron.ConcurrencyPolicy | `internal/cron/store/schedule_resource.go:76` |
| `ScheduleParams.Enabled` | `enabled` | Go field: ScheduleParams.Enabled<br/>type: *bool | `internal/cron/store/schedule_resource.go:79` |
| `ScheduleParams.EndsAtUnixMs` | `endsAtUnixMs` | Go field: ScheduleParams.EndsAtUnixMs<br/>type: *int64 | `internal/cron/store/schedule_resource.go:74` |
| `ScheduleParams.JobName` | `jobName` | Go field: ScheduleParams.JobName<br/>type: string<br/>validate: "required" | `internal/cron/store/schedule_resource.go:70` |
| `ScheduleParams.MisfirePolicy` | `misfirePolicy` | Go field: ScheduleParams.MisfirePolicy<br/>type: cron.MisfirePolicy | `internal/cron/store/schedule_resource.go:75` |
| `ScheduleParams.Name` | `name` | Go field: ScheduleParams.Name<br/>type: string<br/>validate: "required" | `internal/cron/store/schedule_resource.go:68` |
| `ScheduleParams.NewName` | `newName` | Go field: ScheduleParams.NewName<br/>type: string | `internal/cron/store/schedule_resource.go:69` |
| `ScheduleParams.Params` | `params` | Go field: ScheduleParams.Params<br/>type: json.RawMessage | `internal/cron/store/schedule_resource.go:72` |
| `ScheduleParams.Recover` | `recover` | Go field: ScheduleParams.Recover<br/>type: bool | `internal/cron/store/schedule_resource.go:77` |
| `ScheduleParams.StartsAtUnixMs` | `startsAtUnixMs` | Go field: ScheduleParams.StartsAtUnixMs<br/>type: *int64 | `internal/cron/store/schedule_resource.go:73` |
| `ScheduleParams.TimeoutMs` | `timeoutMs` | Go field: ScheduleParams.TimeoutMs<br/>type: int64 | `internal/cron/store/schedule_resource.go:78` |
| `ScheduleParams.Trigger` | `trigger` | Go field: ScheduleParams.Trigger<br/>type: TriggerParams | `internal/cron/store/schedule_resource.go:71` |
| `ScheduleSearch.IsEnabled` | `isEnabled` | Go field: ScheduleSearch.IsEnabled<br/>search: "eq,column=is_enabled"<br/>type: *bool | `internal/cron/store/schedule_resource.go:149` |
| `ScheduleSearch.JobName` | `jobName` | Go field: ScheduleSearch.JobName<br/>search: "eq,column=job_name"<br/>type: string | `internal/cron/store/schedule_resource.go:147` |
| `ScheduleSearch.Kind` | `kind` | Go field: ScheduleSearch.Kind<br/>search: "eq"<br/>type: string | `internal/cron/store/schedule_resource.go:148` |
| `ScheduleSearch.Name` | `name` | Go field: ScheduleSearch.Name<br/>search: "contains"<br/>type: string | `internal/cron/store/schedule_resource.go:146` |
| `ScheduleSpec.ConcurrencyPolicy` | `concurrencyPolicy` | Go field: ScheduleSpec.ConcurrencyPolicy<br/>type: ConcurrencyPolicy | `cron/schedule.go:133` |
| `ScheduleSpec.Enabled` | `enabled` | Go field: ScheduleSpec.Enabled<br/>type: *bool | `cron/schedule.go:141` |
| `ScheduleSpec.EndsAt` | `endsAt` | Go field: ScheduleSpec.EndsAt<br/>type: *time.Time | `cron/schedule.go:130` |
| `ScheduleSpec.JobName` | `jobName` | Go field: ScheduleSpec.JobName<br/>type: string | `cron/schedule.go:122` |
| `ScheduleSpec.MisfirePolicy` | `misfirePolicy` | Go field: ScheduleSpec.MisfirePolicy<br/>type: MisfirePolicy | `cron/schedule.go:132` |
| `ScheduleSpec.Name` | `name` | Go field: ScheduleSpec.Name<br/>type: string | `cron/schedule.go:120` |
| `ScheduleSpec.Params` | `params` | Go field: ScheduleSpec.Params<br/>type: any | `cron/schedule.go:127` |
| `ScheduleSpec.Recover` | `recover` | Go field: ScheduleSpec.Recover<br/>type: bool | `cron/schedule.go:136` |
| `ScheduleSpec.StartsAt` | `startsAt` | Go field: ScheduleSpec.StartsAt<br/>type: *time.Time | `cron/schedule.go:129` |
| `ScheduleSpec.Timeout` | `timeout` | Go field: ScheduleSpec.Timeout<br/>type: time.Duration | `cron/schedule.go:139` |
| `ScheduleSpec.Trigger` | `trigger` | Go field: ScheduleSpec.Trigger<br/>type: TriggerSpec | `cron/schedule.go:124` |
| `Sortable.Sort` | `sort` | Go field: Sortable.Sort<br/>type: []sortx.OrderSpec | `crud/params.go:35` |
| `StartForm.Description` | `description` | Go field: StartForm.Description<br/>type: *string | `approval/my/start_form.go:16` |
| `StartForm.FlowCode` | `flowCode` | Go field: StartForm.FlowCode<br/>type: string | `approval/my/start_form.go:13` |
| `StartForm.FlowID` | `flowId` | Go field: StartForm.FlowID<br/>type: string | `approval/my/start_form.go:12` |
| `StartForm.FlowIcon` | `flowIcon` | Go field: StartForm.FlowIcon<br/>type: *string | `approval/my/start_form.go:15` |
| `StartForm.FlowName` | `flowName` | Go field: StartForm.FlowName<br/>type: string | `approval/my/start_form.go:14` |
| `StartForm.FormSchema` | `formSchema` | Go field: StartForm.FormSchema<br/>type: json.RawMessage | `approval/my/start_form.go:19` |
| `StartForm.Version` | `version` | Go field: StartForm.Version<br/>type: int | `approval/my/start_form.go:18` |
| `StartForm.VersionID` | `versionId` | Go field: StartForm.VersionID<br/>type: string | `approval/my/start_form.go:17` |
| `StartParams.BusinessRef` | `businessRef` | Go field: StartParams.BusinessRef<br/>type: *string<br/>validate: "omitempty,max=512" | `internal/approval/resource/instance.go:141` |
| `StartParams.FlowCode` | `flowCode` | Go field: StartParams.FlowCode<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:140` |
| `StartParams.FormData` | `formData` | Go field: StartParams.FormData<br/>type: map[string]any | `internal/approval/resource/instance.go:142` |
| `StartParams.TenantID` | `tenantId` | Go field: StartParams.TenantID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:139` |
| `StreamGroupInfo.Consumers` | `consumers` | Go field: StreamGroupInfo.Consumers<br/>type: int64 | `event/streams.go:38` |
| `StreamGroupInfo.Lag` | `lag` | Go field: StreamGroupInfo.Lag<br/>type: int64 | `event/streams.go:44` |
| `StreamGroupInfo.LastDeliveredID` | `lastDeliveredId` | Go field: StreamGroupInfo.LastDeliveredID<br/>type: string | `event/streams.go:47` |
| `StreamGroupInfo.Name` | `name` | Go field: StreamGroupInfo.Name<br/>type: string | `event/streams.go:35` |
| `StreamGroupInfo.Pending` | `pending` | Go field: StreamGroupInfo.Pending<br/>type: int64 | `event/streams.go:40` |
| `StreamInfo.Groups` | `groups` | Go field: StreamInfo.Groups<br/>type: []StreamGroupInfo | `event/streams.go:26` |
| `StreamInfo.Length` | `length` | Go field: StreamInfo.Length<br/>type: int64 | `event/streams.go:24` |
| `StreamInfo.Stream` | `stream` | Go field: StreamInfo.Stream<br/>type: string | `event/streams.go:22` |
| `SwapMemory.Free` | `free` | Go field: SwapMemory.Free<br/>type: uint64 | `monitor/service.go:134` |
| `SwapMemory.PageFault` | `pageFault` | Go field: SwapMemory.PageFault<br/>type: uint64 | `monitor/service.go:140` |
| `SwapMemory.PageIn` | `pageIn` | Go field: SwapMemory.PageIn<br/>type: uint64 | `monitor/service.go:138` |
| `SwapMemory.PageMajorFault` | `pageMajorFault` | Go field: SwapMemory.PageMajorFault<br/>type: uint64 | `monitor/service.go:141` |
| `SwapMemory.PageOut` | `pageOut` | Go field: SwapMemory.PageOut<br/>type: uint64 | `monitor/service.go:139` |
| `SwapMemory.SwapIn` | `swapIn` | Go field: SwapMemory.SwapIn<br/>type: uint64 | `monitor/service.go:136` |
| `SwapMemory.SwapOut` | `swapOut` | Go field: SwapMemory.SwapOut<br/>type: uint64 | `monitor/service.go:137` |
| `SwapMemory.Total` | `total` | Go field: SwapMemory.Total<br/>type: uint64 | `monitor/service.go:132` |
| `SwapMemory.Used` | `used` | Go field: SwapMemory.Used<br/>type: uint64 | `monitor/service.go:133` |
| `SwapMemory.UsedPercent` | `usedPercent` | Go field: SwapMemory.UsedPercent<br/>type: float64 | `monitor/service.go:135` |
| `System.BaseURL` | `baseUrl` | Go field: System.BaseURL<br/>type: string | `integration/models.go:45` |
| `System.Code` | `code` | Go field: System.Code<br/>type: string | `integration/models.go:43` |
| `System.DataSource` | `dataSource` | Go field: System.DataSource<br/>type: *DataSourceConfig | `integration/models.go:61` |
| `System.InboundAuth` | `inboundAuth` | Go field: System.InboundAuth<br/>type: *InboundAuthConfig | `integration/models.go:57` |
| `System.IsEnabled` | `isEnabled` | Go field: System.IsEnabled<br/>type: bool | `integration/models.go:69` |
| `System.Name` | `name` | Go field: System.Name<br/>type: string | `integration/models.go:44` |
| `System.OutboundAuth` | `outboundAuth` | Go field: System.OutboundAuth<br/>type: *OutboundAuthConfig | `integration/models.go:48` |
| `System.OutboundEnvelope` | `outboundEnvelope` | Go field: System.OutboundEnvelope<br/>type: *OutboundEnvelopeConfig | `integration/models.go:52` |
| `System.Params` | `params` | Go field: System.Params<br/>type: map[string]string | `integration/models.go:64` |
| `System.Retry` | `retry` | Go field: System.Retry<br/>type: *RetryPolicy | `integration/models.go:68` |
| `System.TimeoutMs` | `timeoutMs` | Go field: System.TimeoutMs<br/>type: int | `integration/models.go:67` |
| `SystemOverview.Build` | `build` | Go field: SystemOverview.Build<br/>type: *BuildInfo | `monitor/service.go:14` |
| `SystemOverview.CPU` | `cpu` | Go field: SystemOverview.CPU<br/>type: *CPUSummary | `monitor/service.go:8` |
| `SystemOverview.Disk` | `disk` | Go field: SystemOverview.Disk<br/>type: *DiskSummary | `monitor/service.go:10` |
| `SystemOverview.Host` | `host` | Go field: SystemOverview.Host<br/>type: *HostSummary | `monitor/service.go:7` |
| `SystemOverview.Load` | `load` | Go field: SystemOverview.Load<br/>type: *LoadInfo | `monitor/service.go:13` |
| `SystemOverview.Memory` | `memory` | Go field: SystemOverview.Memory<br/>type: *MemorySummary | `monitor/service.go:9` |
| `SystemOverview.Network` | `network` | Go field: SystemOverview.Network<br/>type: *NetworkSummary | `monitor/service.go:11` |
| `SystemOverview.Process` | `process` | Go field: SystemOverview.Process<br/>type: *ProcessSummary | `monitor/service.go:12` |
| `SystemParams.BaseURL` | `baseUrl` | Go field: SystemParams.BaseURL<br/>type: string | `internal/integration/resource/system.go:24` |
| `SystemParams.Code` | `code` | Go field: SystemParams.Code<br/>type: string<br/>validate: "required" | `internal/integration/resource/system.go:22` |
| `SystemParams.DataSource` | `dataSource` | Go field: SystemParams.DataSource<br/>type: *integration.DataSourceConfig | `internal/integration/resource/system.go:28` |
| `SystemParams.ID` | `id` | Go field: SystemParams.ID<br/>type: string | `internal/integration/resource/system.go:21` |
| `SystemParams.InboundAuth` | `inboundAuth` | Go field: SystemParams.InboundAuth<br/>type: *integration.InboundAuthConfig | `internal/integration/resource/system.go:27` |
| `SystemParams.IsEnabled` | `isEnabled` | Go field: SystemParams.IsEnabled<br/>type: bool | `internal/integration/resource/system.go:32` |
| `SystemParams.Name` | `name` | Go field: SystemParams.Name<br/>type: string<br/>validate: "required" | `internal/integration/resource/system.go:23` |
| `SystemParams.OutboundAuth` | `outboundAuth` | Go field: SystemParams.OutboundAuth<br/>type: *integration.OutboundAuthConfig | `internal/integration/resource/system.go:25` |
| `SystemParams.OutboundEnvelope` | `outboundEnvelope` | Go field: SystemParams.OutboundEnvelope<br/>type: *integration.OutboundEnvelopeConfig | `internal/integration/resource/system.go:26` |
| `SystemParams.Params` | `params` | Go field: SystemParams.Params<br/>type: map[string]string | `internal/integration/resource/system.go:29` |
| `SystemParams.Retry` | `retry` | Go field: SystemParams.Retry<br/>type: *integration.RetryPolicy | `internal/integration/resource/system.go:31` |
| `SystemParams.TimeoutMs` | `timeoutMs` | Go field: SystemParams.TimeoutMs<br/>type: int | `internal/integration/resource/system.go:30` |
| `SystemSearch.Code` | `code` | Go field: SystemSearch.Code<br/>search: "contains"<br/>type: string | `internal/integration/resource/system.go:39` |
| `SystemSearch.IsEnabled` | `isEnabled` | Go field: SystemSearch.IsEnabled<br/>search: "eq,column=is_enabled"<br/>type: *bool | `internal/integration/resource/system.go:41` |
| `SystemSearch.Name` | `name` | Go field: SystemSearch.Name<br/>search: "contains"<br/>type: string | `internal/integration/resource/system.go:40` |
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
| `Target.Kind` | `kind` | Go field: Target.Kind<br/>type: TargetKind | `push/target.go:18` |
| `Target.Values` | `values` | Go field: Target.Values<br/>type: []string | `push/target.go:19` |
| `Task.AddAssigneeType` | `addAssigneeType` | Go field: Task.AddAssigneeType<br/>type: *AddAssigneeType | `approval/models.go:320` |
| `Task.Assignee` | `assignee` | Go field: Task.Assignee<br/>type: approval.UserInfo | `approval/admin/task.go:15` |
| `Task.AssigneeDepartmentID` | `assigneeDepartmentId` | Go field: Task.AssigneeDepartmentID<br/>type: *string | `approval/models.go:310` |
| `Task.AssigneeDepartmentName` | `assigneeDepartmentName` | Go field: Task.AssigneeDepartmentName<br/>type: *string | `approval/models.go:311` |
| `Task.AssigneeID` | `assigneeId` | Go field: Task.AssigneeID<br/>type: string | `approval/models.go:308` |
| `Task.AssigneeName` | `assigneeName` | Go field: Task.AssigneeName<br/>type: string | `approval/models.go:309` |
| `Task.CreatedAt` | `createdAt` | Go field: Task.CreatedAt<br/>type: timex.DateTime | `approval/admin/task.go:17` |
| `Task.Deadline` | `deadline` | Go field: Task.Deadline<br/>type: *timex.DateTime | `approval/models.go:321` |
| `Task.Deadline` | `deadline` | Go field: Task.Deadline<br/>type: *timex.DateTime | `approval/admin/task.go:18` |
| `Task.DelegatorDepartmentID` | `delegatorDepartmentId` | Go field: Task.DelegatorDepartmentID<br/>type: *string | `approval/models.go:314` |
| `Task.DelegatorDepartmentName` | `delegatorDepartmentName` | Go field: Task.DelegatorDepartmentName<br/>type: *string | `approval/models.go:315` |
| `Task.DelegatorID` | `delegatorId` | Go field: Task.DelegatorID<br/>type: *string | `approval/models.go:312` |
| `Task.DelegatorName` | `delegatorName` | Go field: Task.DelegatorName<br/>type: *string | `approval/models.go:313` |
| `Task.FinishedAt` | `finishedAt` | Go field: Task.FinishedAt<br/>type: *timex.DateTime | `approval/admin/task.go:19` |
| `Task.FinishedAt` | `finishedAt` | Go field: Task.FinishedAt<br/>type: *timex.DateTime | `approval/models.go:324` |
| `Task.FlowName` | `flowName` | Go field: Task.FlowName<br/>type: string | `approval/admin/task.go:13` |
| `Task.InstanceID` | `instanceId` | Go field: Task.InstanceID<br/>type: string | `approval/models.go:305` |
| `Task.InstanceID` | `instanceId` | Go field: Task.InstanceID<br/>type: string | `approval/admin/task.go:11` |
| `Task.InstanceTitle` | `instanceTitle` | Go field: Task.InstanceTitle<br/>type: string | `approval/admin/task.go:12` |
| `Task.IsPreWarningSent` | `isPreWarningSent` | Go field: Task.IsPreWarningSent<br/>type: bool | `approval/models.go:323` |
| `Task.IsTimeout` | `isTimeout` | Go field: Task.IsTimeout<br/>type: bool | `approval/models.go:322` |
| `Task.NodeID` | `nodeId` | Go field: Task.NodeID<br/>type: string | `approval/models.go:306` |
| `Task.NodeName` | `nodeName` | Go field: Task.NodeName<br/>type: string | `approval/admin/task.go:14` |
| `Task.ParentTaskID` | `parentTaskId` | Go field: Task.ParentTaskID<br/>type: *string | `approval/models.go:319` |
| `Task.ReadAt` | `readAt` | Go field: Task.ReadAt<br/>type: *timex.DateTime | `approval/models.go:318` |
| `Task.SortOrder` | `sortOrder` | Go field: Task.SortOrder<br/>type: int | `approval/models.go:316` |
| `Task.Status` | `status` | Go field: Task.Status<br/>type: string | `approval/admin/task.go:16` |
| `Task.Status` | `status` | Go field: Task.Status<br/>type: TaskStatus | `approval/models.go:317` |
| `Task.TaskID` | `taskId` | Go field: Task.TaskID<br/>type: string | `approval/admin/task.go:10` |
| `Task.TenantID` | `tenantId` | Go field: Task.TenantID<br/>type: string | `approval/models.go:304` |
| `Task.VisitID` | `visitId` | Go field: Task.VisitID<br/>type: string | `approval/models.go:307` |
| `TaskApprovedEvent.Operator` | `operator` | Go field: TaskApprovedEvent.Operator<br/>type: UserInfo | `approval/events_task.go:35` |
| `TaskApprovedEvent.Opinion` | `opinion` | Go field: TaskApprovedEvent.Opinion<br/>type: *string | `approval/events_task.go:36` |
| `TaskCanceledEvent.Assignee` | `assignee` | Go field: TaskCanceledEvent.Assignee<br/>type: UserInfo | `approval/events_task.go:93` |
| `TaskCanceledEvent.Reason` | `reason` | Go field: TaskCanceledEvent.Reason<br/>type: string | `approval/events_task.go:94` |
| `TaskCreatedEvent.Assignee` | `assignee` | Go field: TaskCreatedEvent.Assignee<br/>type: UserInfo | `approval/events_task.go:17` |
| `TaskCreatedEvent.Deadline` | `deadline` | Go field: TaskCreatedEvent.Deadline<br/>type: *timex.DateTime | `approval/events_task.go:18` |
| `TaskDeadlineWarningEvent.Assignee` | `assignee` | Go field: TaskDeadlineWarningEvent.Assignee<br/>type: UserInfo | `approval/events_timeout.go:10` |
| `TaskDeadlineWarningEvent.Deadline` | `deadline` | Go field: TaskDeadlineWarningEvent.Deadline<br/>type: timex.DateTime | `approval/events_timeout.go:11` |
| `TaskDeadlineWarningEvent.HoursLeft` | `hoursLeft` | Go field: TaskDeadlineWarningEvent.HoursLeft<br/>type: int | `approval/events_timeout.go:12` |
| `TaskEventBase.NodeID` | `nodeId` | Go field: TaskEventBase.NodeID<br/>type: string | `approval/events_base.go:51` |
| `TaskEventBase.NodeName` | `nodeName` | Go field: TaskEventBase.NodeName<br/>type: string | `approval/events_base.go:52` |
| `TaskEventBase.TaskID` | `taskId` | Go field: TaskEventBase.TaskID<br/>type: string | `approval/events_base.go:50` |
| `TaskHandledEvent.Operator` | `operator` | Go field: TaskHandledEvent.Operator<br/>type: UserInfo | `approval/events_task.go:53` |
| `TaskHandledEvent.Opinion` | `opinion` | Go field: TaskHandledEvent.Opinion<br/>type: *string | `approval/events_task.go:54` |
| `TaskNodeData.AdminUserIDs` | `adminUserIds` | Go field: TaskNodeData.AdminUserIDs<br/>type: []string | `approval/node_data.go:92` |
| `TaskNodeData.Assignees` | `assignees` | Go field: TaskNodeData.Assignees<br/>type: []AssigneeDefinition | `approval/node_data.go:88` |
| `TaskNodeData.CCs` | `ccs` | Go field: TaskNodeData.CCs<br/>type: []CCDefinition | `approval/node_data.go:99` |
| `TaskNodeData.EmptyAssigneeAction` | `emptyAssigneeAction` | Go field: TaskNodeData.EmptyAssigneeAction<br/>type: EmptyAssigneeAction | `approval/node_data.go:90` |
| `TaskNodeData.ExecutionType` | `executionType` | Go field: TaskNodeData.ExecutionType<br/>type: ExecutionType | `approval/node_data.go:89` |
| `TaskNodeData.FallbackUserIDs` | `fallbackUserIds` | Go field: TaskNodeData.FallbackUserIDs<br/>type: []string | `approval/node_data.go:91` |
| `TaskNodeData.FieldPermissions` | `fieldPermissions` | Go field: TaskNodeData.FieldPermissions<br/>type: map[string]Permission | `approval/node_data.go:100` |
| `TaskNodeData.IsOpinionRequired` | `isOpinionRequired` | Go field: TaskNodeData.IsOpinionRequired<br/>type: bool | `approval/node_data.go:94` |
| `TaskNodeData.IsTransferAllowed` | `isTransferAllowed` | Go field: TaskNodeData.IsTransferAllowed<br/>type: *bool | `approval/node_data.go:93` |
| `TaskNodeData.TimeoutAction` | `timeoutAction` | Go field: TaskNodeData.TimeoutAction<br/>type: TimeoutAction | `approval/node_data.go:96` |
| `TaskNodeData.TimeoutHours` | `timeoutHours` | Go field: TaskNodeData.TimeoutHours<br/>type: int | `approval/node_data.go:95` |
| `TaskNodeData.TimeoutNotifyBeforeHours` | `timeoutNotifyBeforeHours` | Go field: TaskNodeData.TimeoutNotifyBeforeHours<br/>type: int | `approval/node_data.go:97` |
| `TaskNodeData.UrgeCooldownMinutes` | `urgeCooldownMinutes` | Go field: TaskNodeData.UrgeCooldownMinutes<br/>type: int | `approval/node_data.go:98` |
| `TaskReassignedEvent.From` | `from` | Go field: TaskReassignedEvent.From<br/>type: UserInfo | `approval/events_task.go:131` |
| `TaskReassignedEvent.Reason` | `reason` | Go field: TaskReassignedEvent.Reason<br/>type: *string | `approval/events_task.go:133` |
| `TaskReassignedEvent.To` | `to` | Go field: TaskReassignedEvent.To<br/>type: UserInfo | `approval/events_task.go:132` |
| `TaskRejectedEvent.Operator` | `operator` | Go field: TaskRejectedEvent.Operator<br/>type: UserInfo | `approval/events_task.go:71` |
| `TaskRejectedEvent.Opinion` | `opinion` | Go field: TaskRejectedEvent.Opinion<br/>type: *string | `approval/events_task.go:72` |
| `TaskTimedOutEvent.Assignee` | `assignee` | Go field: TaskTimedOutEvent.Assignee<br/>type: UserInfo | `approval/events_task.go:151` |
| `TaskTimedOutEvent.Deadline` | `deadline` | Go field: TaskTimedOutEvent.Deadline<br/>type: timex.DateTime | `approval/events_task.go:152` |
| `TaskTransferredEvent.From` | `from` | Go field: TaskTransferredEvent.From<br/>type: UserInfo | `approval/events_task.go:111` |
| `TaskTransferredEvent.Reason` | `reason` | Go field: TaskTransferredEvent.Reason<br/>type: *string | `approval/events_task.go:113` |
| `TaskTransferredEvent.To` | `to` | Go field: TaskTransferredEvent.To<br/>type: UserInfo | `approval/events_task.go:112` |
| `TaskUrgedEvent.Message` | `message` | Go field: TaskUrgedEvent.Message<br/>type: *string | `approval/events_timeout.go:37` |
| `TaskUrgedEvent.Target` | `target` | Go field: TaskUrgedEvent.Target<br/>type: UserInfo | `approval/events_timeout.go:36` |
| `TaskUrgedEvent.Urger` | `urger` | Go field: TaskUrgedEvent.Urger<br/>type: UserInfo | `approval/events_timeout.go:35` |
| `TestConnectionParams.Method` | `method` | Go field: TestConnectionParams.Method<br/>type: string | `internal/integration/resource/ops.go:59` |
| `TestConnectionParams.Path` | `path` | Go field: TestConnectionParams.Path<br/>type: string | `internal/integration/resource/ops.go:60` |
| `TestConnectionParams.SystemCode` | `systemCode` | Go field: TestConnectionParams.SystemCode<br/>type: string<br/>validate: "required" | `internal/integration/resource/ops.go:58` |
| `TimelineEntry.Activities` | `activities` | Go field: TimelineEntry.Activities<br/>type: []Activity | `approval/timeline_view.go:50` |
| `TimelineEntry.ApprovalMethod` | `approvalMethod` | Go field: TimelineEntry.ApprovalMethod<br/>type: string | `approval/timeline_view.go:45` |
| `TimelineEntry.CCRecipients` | `ccRecipients` | Go field: TimelineEntry.CCRecipients<br/>type: []CCRecipient | `approval/timeline_view.go:49` |
| `TimelineEntry.ExecutionType` | `executionType` | Go field: TimelineEntry.ExecutionType<br/>type: string | `approval/timeline_view.go:44` |
| `TimelineEntry.FinishedAt` | `finishedAt` | Go field: TimelineEntry.FinishedAt<br/>type: *timex.DateTime | `approval/timeline_view.go:52` |
| `TimelineEntry.Kind` | `kind` | Go field: TimelineEntry.Kind<br/>type: TimelineEntryKind | `approval/timeline_view.go:40` |
| `TimelineEntry.Name` | `name` | Go field: TimelineEntry.Name<br/>type: string | `approval/timeline_view.go:42` |
| `TimelineEntry.NodeID` | `nodeId` | Go field: TimelineEntry.NodeID<br/>type: *string | `approval/timeline_view.go:41` |
| `TimelineEntry.Participants` | `participants` | Go field: TimelineEntry.Participants<br/>type: []NodeParticipant | `approval/timeline_view.go:48` |
| `TimelineEntry.PassRatio` | `passRatio` | Go field: TimelineEntry.PassRatio<br/>type: *decimal.Decimal | `approval/timeline_view.go:47` |
| `TimelineEntry.PassRule` | `passRule` | Go field: TimelineEntry.PassRule<br/>type: string | `approval/timeline_view.go:46` |
| `TimelineEntry.StartedAt` | `startedAt` | Go field: TimelineEntry.StartedAt<br/>type: timex.DateTime | `approval/timeline_view.go:51` |
| `TimelineEntry.Status` | `status` | Go field: TimelineEntry.Status<br/>type: NodeVisitStatus | `approval/timeline_view.go:43` |
| `ToggleActiveParams.FlowID` | `flowId` | Go field: ToggleActiveParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:337` |
| `ToggleActiveParams.IsActive` | `isActive` | Go field: ToggleActiveParams.IsActive<br/>type: bool | `internal/approval/resource/flow.go:338` |
| `TreeDataOption.Children` | `children` | Go field: TreeDataOption.Children<br/>type: []TreeDataOption | `crud/option.go:57` |
| `TriggerParams.AtUnixMs` | `atUnixMs` | Go field: TriggerParams.AtUnixMs<br/>type: *int64 | `internal/cron/store/schedule_resource.go:28` |
| `TriggerParams.EveryMs` | `everyMs` | Go field: TriggerParams.EveryMs<br/>type: *int64 | `internal/cron/store/schedule_resource.go:27` |
| `TriggerParams.Expr` | `expr` | Go field: TriggerParams.Expr<br/>type: *string | `internal/cron/store/schedule_resource.go:25` |
| `TriggerParams.Kind` | `kind` | Go field: TriggerParams.Kind<br/>type: cron.TriggerKind<br/>validate: "required" | `internal/cron/store/schedule_resource.go:24` |
| `TriggerParams.Timezone` | `timezone` | Go field: TriggerParams.Timezone<br/>type: *string | `internal/cron/store/schedule_resource.go:26` |
| `TriggerSpec.At` | `at` | Go field: TriggerSpec.At<br/>type: *time.Time | `cron/trigger.go:62` |
| `TriggerSpec.EveryMs` | `everyMs` | Go field: TriggerSpec.EveryMs<br/>type: int64 | `cron/trigger.go:60` |
| `TriggerSpec.Expr` | `expr` | Go field: TriggerSpec.Expr<br/>type: string | `cron/trigger.go:55` |
| `TriggerSpec.Kind` | `kind` | Go field: TriggerSpec.Kind<br/>type: TriggerKind | `cron/trigger.go:52` |
| `TriggerSpec.Timezone` | `timezone` | Go field: TriggerSpec.Timezone<br/>type: string | `cron/trigger.go:58` |
| `UniqueKey.Columns` | `columns` | Go field: UniqueKey.Columns<br/>type: []string | `schema/service.go:51` |
| `UniqueKey.HasExpressions` | `hasExpressions` | Go field: UniqueKey.HasExpressions<br/>type: bool | `schema/service.go:53` |
| `UniqueKey.Name` | `name` | Go field: UniqueKey.Name<br/>type: string | `schema/service.go:50` |
| `UniqueKey.Predicate` | `predicate` | Go field: UniqueKey.Predicate<br/>type: string | `schema/service.go:52` |
| `UpdateFlowParams.AdminUserIDs` | `adminUserIds` | Go field: UpdateFlowParams.AdminUserIDs<br/>type: []string | `internal/approval/resource/flow.go:287` |
| `UpdateFlowParams.BindingMode` | `bindingMode` | Go field: UpdateFlowParams.BindingMode<br/>type: approval.BindingMode<br/>validate: "required" | `internal/approval/resource/flow.go:285` |
| `UpdateFlowParams.BusinessBinding` | `businessBinding` | Go field: UpdateFlowParams.BusinessBinding<br/>type: *approval.BusinessBindingConfig | `internal/approval/resource/flow.go:286` |
| `UpdateFlowParams.Description` | `description` | Go field: UpdateFlowParams.Description<br/>type: *string | `internal/approval/resource/flow.go:283` |
| `UpdateFlowParams.FlowID` | `flowId` | Go field: UpdateFlowParams.FlowID<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:280` |
| `UpdateFlowParams.Icon` | `icon` | Go field: UpdateFlowParams.Icon<br/>type: *string | `internal/approval/resource/flow.go:282` |
| `UpdateFlowParams.Initiators` | `initiators` | Go field: UpdateFlowParams.Initiators<br/>type: []CreateInitiatorParams | `internal/approval/resource/flow.go:290` |
| `UpdateFlowParams.InstanceTitleTemplate` | `instanceTitleTemplate` | Go field: UpdateFlowParams.InstanceTitleTemplate<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:289` |
| `UpdateFlowParams.IsAllInitiationAllowed` | `isAllInitiationAllowed` | Go field: UpdateFlowParams.IsAllInitiationAllowed<br/>type: bool | `internal/approval/resource/flow.go:288` |
| `UpdateFlowParams.Labels` | `labels` | Go field: UpdateFlowParams.Labels<br/>type: map[string]string | `internal/approval/resource/flow.go:284` |
| `UpdateFlowParams.Name` | `name` | Go field: UpdateFlowParams.Name<br/>type: string<br/>validate: "required" | `internal/approval/resource/flow.go:281` |
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
| `UploadPartParams.ClaimID` | `claimId` | Go field: UploadPartParams.ClaimID<br/>type: string<br/>validate: "required" | `internal/storage/resource.go:405` |
| `UploadPartParams.PartNumber` | `partNumber` | Go field: UploadPartParams.PartNumber<br/>type: int<br/>validate: "required,min=1" | `internal/storage/resource.go:406` |
| `UploadPartResult.PartNumber` | `partNumber` | Go field: UploadPartResult.PartNumber<br/>type: int | `internal/storage/resource.go:414` |
| `UploadPartResult.Size` | `size` | Go field: UploadPartResult.Size<br/>type: int64 | `internal/storage/resource.go:415` |
| `UrgeRecord.InstanceID` | `instanceId` | Go field: UrgeRecord.InstanceID<br/>type: string | `approval/models.go:509` |
| `UrgeRecord.Message` | `message` | Go field: UrgeRecord.Message<br/>type: string | `approval/models.go:520` |
| `UrgeRecord.NodeID` | `nodeId` | Go field: UrgeRecord.NodeID<br/>type: string | `approval/models.go:510` |
| `UrgeRecord.TargetUserDepartmentID` | `targetUserDepartmentId` | Go field: UrgeRecord.TargetUserDepartmentID<br/>type: *string | `approval/models.go:518` |
| `UrgeRecord.TargetUserDepartmentName` | `targetUserDepartmentName` | Go field: UrgeRecord.TargetUserDepartmentName<br/>type: *string | `approval/models.go:519` |
| `UrgeRecord.TargetUserID` | `targetUserId` | Go field: UrgeRecord.TargetUserID<br/>type: string | `approval/models.go:516` |
| `UrgeRecord.TargetUserName` | `targetUserName` | Go field: UrgeRecord.TargetUserName<br/>type: string | `approval/models.go:517` |
| `UrgeRecord.TaskID` | `taskId` | Go field: UrgeRecord.TaskID<br/>type: *string | `approval/models.go:511` |
| `UrgeRecord.UrgerDepartmentID` | `urgerDepartmentId` | Go field: UrgeRecord.UrgerDepartmentID<br/>type: *string | `approval/models.go:514` |
| `UrgeRecord.UrgerDepartmentName` | `urgerDepartmentName` | Go field: UrgeRecord.UrgerDepartmentName<br/>type: *string | `approval/models.go:515` |
| `UrgeRecord.UrgerID` | `urgerId` | Go field: UrgeRecord.UrgerID<br/>type: string | `approval/models.go:512` |
| `UrgeRecord.UrgerName` | `urgerName` | Go field: UrgeRecord.UrgerName<br/>type: string | `approval/models.go:513` |
| `UrgeTaskParams.Message` | `message` | Go field: UrgeTaskParams.Message<br/>type: string<br/>validate: "max=500" | `internal/approval/resource/instance.go:455` |
| `UrgeTaskParams.TaskID` | `taskId` | Go field: UrgeTaskParams.TaskID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:454` |
| `UserInfo.Avatar` | `avatar` | Go field: UserInfo.Avatar<br/>type: *string | `security/user_info.go:45` |
| `UserInfo.DepartmentID` | `departmentId` | Go field: UserInfo.DepartmentID<br/>type: *string | `approval/assignee.go:22` |
| `UserInfo.DepartmentName` | `departmentName` | Go field: UserInfo.DepartmentName<br/>type: *string | `approval/assignee.go:23` |
| `UserInfo.Details` | `details` | Go field: UserInfo.Details<br/>type: any | `security/user_info.go:48` |
| `UserInfo.Gender` | `gender` | Go field: UserInfo.Gender<br/>type: Gender | `security/user_info.go:44` |
| `UserInfo.ID` | `id` | Go field: UserInfo.ID<br/>type: string | `approval/assignee.go:20` |
| `UserInfo.ID` | `id` | Go field: UserInfo.ID<br/>type: string | `security/user_info.go:42` |
| `UserInfo.Menus` | `menus` | Go field: UserInfo.Menus<br/>type: []UserMenu | `security/user_info.go:47` |
| `UserInfo.Name` | `name` | Go field: UserInfo.Name<br/>type: string | `security/user_info.go:43` |
| `UserInfo.Name` | `name` | Go field: UserInfo.Name<br/>type: string | `approval/assignee.go:21` |
| `UserInfo.PermissionTokens` | `permissionTokens` | Go field: UserInfo.PermissionTokens<br/>type: []string | `security/user_info.go:46` |
| `UserMenu.Children` | `children` | Go field: UserMenu.Children<br/>type: []UserMenu | `security/user_info.go:38` |
| `UserMenu.Icon` | `icon` | Go field: UserMenu.Icon<br/>type: *string | `security/user_info.go:36` |
| `UserMenu.Meta` | `meta` | Go field: UserMenu.Meta<br/>type: map[string]any | `security/user_info.go:37` |
| `UserMenu.Name` | `name` | Go field: UserMenu.Name<br/>type: string | `security/user_info.go:35` |
| `UserMenu.Path` | `path` | Go field: UserMenu.Path<br/>type: string | `security/user_info.go:34` |
| `UserMenu.Type` | `type` | Go field: UserMenu.Type<br/>type: UserMenuType | `security/user_info.go:33` |
| `ValidationRule.Max` | `max` | Go field: ValidationRule.Max<br/>type: *float64 | `approval/form_field.go:52` |
| `ValidationRule.MaxLength` | `maxLength` | Go field: ValidationRule.MaxLength<br/>type: *int | `approval/form_field.go:50` |
| `ValidationRule.Message` | `message` | Go field: ValidationRule.Message<br/>type: string | `approval/form_field.go:54` |
| `ValidationRule.Min` | `min` | Go field: ValidationRule.Min<br/>type: *float64 | `approval/form_field.go:51` |
| `ValidationRule.MinLength` | `minLength` | Go field: ValidationRule.MinLength<br/>type: *int | `approval/form_field.go:49` |
| `ValidationRule.Pattern` | `pattern` | Go field: ValidationRule.Pattern<br/>type: string | `approval/form_field.go:53` |
| `View.Columns` | `columns` | Go field: View.Columns<br/>type: []string | `schema/service.go:78` |
| `View.Comment` | `comment` | Go field: View.Comment<br/>type: string | `schema/service.go:77` |
| `View.Definition` | `definition` | Go field: View.Definition<br/>type: string | `schema/service.go:76` |
| `View.Name` | `name` | Go field: View.Name<br/>type: string | `schema/service.go:74` |
| `View.Schema` | `schema` | Go field: View.Schema<br/>type: string | `schema/service.go:75` |
| `ViewerTask.AddAssigneeTypes` | `addAssigneeTypes` | Go field: ViewerTask.AddAssigneeTypes<br/>type: []approval.AddAssigneeType | `approval/my/instance_detail.go:45` |
| `ViewerTask.IsOpinionRequired` | `isOpinionRequired` | Go field: ViewerTask.IsOpinionRequired<br/>type: bool | `approval/my/instance_detail.go:42` |
| `ViewerTask.NodeID` | `nodeId` | Go field: ViewerTask.NodeID<br/>type: string | `approval/my/instance_detail.go:39` |
| `ViewerTask.RemovableAssignees` | `removableAssignees` | Go field: ViewerTask.RemovableAssignees<br/>type: []RemovableAssignee | `approval/my/instance_detail.go:54` |
| `ViewerTask.RollbackTargets` | `rollbackTargets` | Go field: ViewerTask.RollbackTargets<br/>type: []RollbackTarget | `approval/my/instance_detail.go:49` |
| `ViewerTask.TaskID` | `taskId` | Go field: ViewerTask.TaskID<br/>type: string | `approval/my/instance_detail.go:38` |
| `VirtualMemory.Active` | `active` | Go field: VirtualMemory.Active<br/>type: uint64 | `monitor/service.go:95` |
| `VirtualMemory.AnonHugePages` | `anonHugePages` | Go field: VirtualMemory.AnonHugePages<br/>type: uint64 | `monitor/service.go:127` |
| `VirtualMemory.Available` | `available` | Go field: VirtualMemory.Available<br/>type: uint64 | `monitor/service.go:91` |
| `VirtualMemory.Buffers` | `buffers` | Go field: VirtualMemory.Buffers<br/>type: uint64 | `monitor/service.go:99` |
| `VirtualMemory.Cached` | `cached` | Go field: VirtualMemory.Cached<br/>type: uint64 | `monitor/service.go:100` |
| `VirtualMemory.CommitLimit` | `commitLimit` | Go field: VirtualMemory.CommitLimit<br/>type: uint64 | `monitor/service.go:110` |
| `VirtualMemory.CommittedAs` | `committedAs` | Go field: VirtualMemory.CommittedAs<br/>type: uint64 | `monitor/service.go:111` |
| `VirtualMemory.Dirty` | `dirty` | Go field: VirtualMemory.Dirty<br/>type: uint64 | `monitor/service.go:102` |
| `VirtualMemory.Free` | `free` | Go field: VirtualMemory.Free<br/>type: uint64 | `monitor/service.go:94` |
| `VirtualMemory.HighFree` | `highFree` | Go field: VirtualMemory.HighFree<br/>type: uint64 | `monitor/service.go:113` |
| `VirtualMemory.HighTotal` | `highTotal` | Go field: VirtualMemory.HighTotal<br/>type: uint64 | `monitor/service.go:112` |
| `VirtualMemory.HugePageSize` | `hugePageSize` | Go field: VirtualMemory.HugePageSize<br/>type: uint64 | `monitor/service.go:126` |
| `VirtualMemory.HugePagesFree` | `hugePagesFree` | Go field: VirtualMemory.HugePagesFree<br/>type: uint64 | `monitor/service.go:123` |
| `VirtualMemory.HugePagesReserved` | `hugePagesReserved` | Go field: VirtualMemory.HugePagesReserved<br/>type: uint64 | `monitor/service.go:124` |
| `VirtualMemory.HugePagesSurplus` | `hugePagesSurplus` | Go field: VirtualMemory.HugePagesSurplus<br/>type: uint64 | `monitor/service.go:125` |
| `VirtualMemory.HugePagesTotal` | `hugePagesTotal` | Go field: VirtualMemory.HugePagesTotal<br/>type: uint64 | `monitor/service.go:122` |
| `VirtualMemory.Inactive` | `inactive` | Go field: VirtualMemory.Inactive<br/>type: uint64 | `monitor/service.go:96` |
| `VirtualMemory.Laundry` | `laundry` | Go field: VirtualMemory.Laundry<br/>type: uint64 | `monitor/service.go:98` |
| `VirtualMemory.LowFree` | `lowFree` | Go field: VirtualMemory.LowFree<br/>type: uint64 | `monitor/service.go:115` |
| `VirtualMemory.LowTotal` | `lowTotal` | Go field: VirtualMemory.LowTotal<br/>type: uint64 | `monitor/service.go:114` |
| `VirtualMemory.Mapped` | `mapped` | Go field: VirtualMemory.Mapped<br/>type: uint64 | `monitor/service.go:118` |
| `VirtualMemory.PageTables` | `pageTables` | Go field: VirtualMemory.PageTables<br/>type: uint64 | `monitor/service.go:108` |
| `VirtualMemory.Shared` | `shared` | Go field: VirtualMemory.Shared<br/>type: uint64 | `monitor/service.go:104` |
| `VirtualMemory.Slab` | `slab` | Go field: VirtualMemory.Slab<br/>type: uint64 | `monitor/service.go:105` |
| `VirtualMemory.SlabReclaimable` | `slabReclaimable` | Go field: VirtualMemory.SlabReclaimable<br/>type: uint64 | `monitor/service.go:106` |
| `VirtualMemory.SlabUnreclaimable` | `slabUnreclaimable` | Go field: VirtualMemory.SlabUnreclaimable<br/>type: uint64 | `monitor/service.go:107` |
| `VirtualMemory.SwapCached` | `swapCached` | Go field: VirtualMemory.SwapCached<br/>type: uint64 | `monitor/service.go:109` |
| `VirtualMemory.SwapFree` | `swapFree` | Go field: VirtualMemory.SwapFree<br/>type: uint64 | `monitor/service.go:117` |
| `VirtualMemory.SwapTotal` | `swapTotal` | Go field: VirtualMemory.SwapTotal<br/>type: uint64 | `monitor/service.go:116` |
| `VirtualMemory.Total` | `total` | Go field: VirtualMemory.Total<br/>type: uint64 | `monitor/service.go:90` |
| `VirtualMemory.Used` | `used` | Go field: VirtualMemory.Used<br/>type: uint64 | `monitor/service.go:92` |
| `VirtualMemory.UsedPercent` | `usedPercent` | Go field: VirtualMemory.UsedPercent<br/>type: float64 | `monitor/service.go:93` |
| `VirtualMemory.VMAllocChunk` | `vmAllocChunk` | Go field: VirtualMemory.VMAllocChunk<br/>type: uint64 | `monitor/service.go:121` |
| `VirtualMemory.VMAllocTotal` | `vmAllocTotal` | Go field: VirtualMemory.VMAllocTotal<br/>type: uint64 | `monitor/service.go:119` |
| `VirtualMemory.VMAllocUsed` | `vmAllocUsed` | Go field: VirtualMemory.VMAllocUsed<br/>type: uint64 | `monitor/service.go:120` |
| `VirtualMemory.Wired` | `wired` | Go field: VirtualMemory.Wired<br/>type: uint64 | `monitor/service.go:97` |
| `VirtualMemory.WriteBack` | `writeBack` | Go field: VirtualMemory.WriteBack<br/>type: uint64 | `monitor/service.go:101` |
| `VirtualMemory.WriteBackTmp` | `writeBackTmp` | Go field: VirtualMemory.WriteBackTmp<br/>type: uint64 | `monitor/service.go:103` |
| `WithdrawParams.InstanceID` | `instanceId` | Go field: WithdrawParams.InstanceID<br/>type: string<br/>validate: "required" | `internal/approval/resource/instance.go:294` |
| `WithdrawParams.Reason` | `reason` | Go field: WithdrawParams.Reason<br/>type: string<br/>validate: "max=2000" | `internal/approval/resource/instance.go:295` |
| `exportConfig.Format` | `format` | Go field: exportConfig.Format<br/>type: TabularFormat | `crud/export.go:87` |
| `fetchInit.Body` | `body` | Go field: fetchInit.Body<br/>type: any | `js/jshttp/jshttp.go:100` |
| `fetchInit.Body` | `body` | Go field: fetchInit.Body<br/>type: any | `internal/integration/exec/http_lib.go:68` |
| `fetchInit.Envelope` | `envelope` | Go field: fetchInit.Envelope<br/>type: *bool | `internal/integration/exec/http_lib.go:71` |
| `fetchInit.Headers` | `headers` | Go field: fetchInit.Headers<br/>type: map[string]string | `internal/integration/exec/http_lib.go:66` |
| `fetchInit.Headers` | `headers` | Go field: fetchInit.Headers<br/>type: map[string]string | `js/jshttp/jshttp.go:98` |
| `fetchInit.Method` | `method` | Go field: fetchInit.Method<br/>type: string | `js/jshttp/jshttp.go:97` |
| `fetchInit.Method` | `method` | Go field: fetchInit.Method<br/>type: string | `internal/integration/exec/http_lib.go:65` |
| `fetchInit.Query` | `query` | Go field: fetchInit.Query<br/>type: map[string]string | `js/jshttp/jshttp.go:99` |
| `fetchInit.Query` | `query` | Go field: fetchInit.Query<br/>type: map[string]string | `internal/integration/exec/http_lib.go:67` |
| `fetchInit.Redirect` | `redirect` | Go field: fetchInit.Redirect<br/>type: string | `internal/integration/exec/http_lib.go:69` |
| `fetchInit.Redirect` | `redirect` | Go field: fetchInit.Redirect<br/>type: string | `js/jshttp/jshttp.go:101` |
| `fetchInit.Timeout` | `timeout` | Go field: fetchInit.Timeout<br/>type: int64 | `js/jshttp/jshttp.go:102` |
| `fetchInit.Timeout` | `timeout` | Go field: fetchInit.Timeout<br/>type: int64 | `internal/integration/exec/http_lib.go:70` |
| `importConfig.Format` | `format` | Go field: importConfig.Format<br/>type: TabularFormat | `crud/import.go:92` |
| `importParams.File` | `file` | Go field: importParams.File<br/>type: *multipart.FileHeader | `crud/import.go:86` |
| `manifest.ContentType` | `contentType` | Go field: manifest.ContentType<br/>type: string | `internal/storage/filesystem/service.go:57` |
| `manifest.Key` | `key` | Go field: manifest.Key<br/>type: string | `internal/storage/filesystem/service.go:56` |
| `manifest.Metadata` | `metadata` | Go field: manifest.Metadata<br/>type: map[string]string | `internal/storage/filesystem/service.go:58` |
| `objectMeta.ContentType` | `contentType` | Go field: objectMeta.ContentType<br/>type: string | `internal/storage/filesystem/service.go:288` |
| `objectMeta.ETag` | `etag` | Go field: objectMeta.ETag<br/>type: string | `internal/storage/filesystem/service.go:287` |
| `objectMeta.Metadata` | `metadata` | Go field: objectMeta.Metadata<br/>type: map[string]string | `internal/storage/filesystem/service.go:289` |
| `redisSessionRecord.Session` | `session` | Go field: redisSessionRecord.Session<br/>type: Session | `security/redis_session_store.go:39` |
| `redisSessionRecord.TokenHash` | `tokenHash` | Go field: redisSessionRecord.TokenHash<br/>type: string | `security/redis_session_store.go:40` |
| `relayFrame.Envelope` | `envelope` | Go field: relayFrame.Envelope<br/>type: json.RawMessage | `internal/push/relay.go:44` |
| `relayFrame.Kind` | `kind` | Go field: relayFrame.Kind<br/>type: string | `internal/push/relay.go:42` |
| `relayFrame.Origin` | `origin` | Go field: relayFrame.Origin<br/>type: string | `internal/push/relay.go:43` |
| `relayFrame.SessionIDs` | `sessionIds` | Go field: relayFrame.SessionIDs<br/>type: []string | `internal/push/relay.go:46` |
| `relayFrame.Targets` | `targets` | Go field: relayFrame.Targets<br/>type: []push.Target | `internal/push/relay.go:45` |
| `requestOptions.Envelope` | `envelope` | Go field: requestOptions.Envelope<br/>type: *bool | `internal/integration/exec/http_lib.go:60` |
| `requestOptions.Headers` | `headers` | Go field: requestOptions.Headers<br/>type: map[string]string | `internal/integration/exec/http_lib.go:54` |
| `requestOptions.Headers` | `headers` | Go field: requestOptions.Headers<br/>type: map[string]string | `js/jshttp/jshttp.go:89` |
| `requestOptions.Query` | `query` | Go field: requestOptions.Query<br/>type: map[string]string | `internal/integration/exec/http_lib.go:55` |
| `requestOptions.Query` | `query` | Go field: requestOptions.Query<br/>type: map[string]string | `js/jshttp/jshttp.go:90` |
| `requestOptions.Redirect` | `redirect` | Go field: requestOptions.Redirect<br/>type: string | `js/jshttp/jshttp.go:91` |
| `requestOptions.Redirect` | `redirect` | Go field: requestOptions.Redirect<br/>type: string | `internal/integration/exec/http_lib.go:56` |
| `requestOptions.Timeout` | `timeout` | Go field: requestOptions.Timeout<br/>type: int64 | `internal/integration/exec/http_lib.go:57` |
| `requestOptions.Timeout` | `timeout` | Go field: requestOptions.Timeout<br/>type: int64 | `js/jshttp/jshttp.go:92` |
| `richBlock.Children` | `children` | Go field: richBlock.Children<br/>type: []richBlock | `internal/approval/formeditor/schema.go:59` |
| `richBlock.ColumnType` | `columnType` | Go field: richBlock.ColumnType<br/>type: string | `internal/approval/formeditor/schema.go:47` |
| `richBlock.DataSource` | `dataSource` | Go field: richBlock.DataSource<br/>type: *richOptionSource | `internal/approval/formeditor/schema.go:51` |
| `richBlock.Key` | `key` | Go field: richBlock.Key<br/>type: string | `internal/approval/formeditor/schema.go:42` |
| `richBlock.Label` | `label` | Go field: richBlock.Label<br/>type: *string | `internal/approval/formeditor/schema.go:46` |
| `richBlock.MaxRows` | `maxRows` | Go field: richBlock.MaxRows<br/>type: *int | `internal/approval/formeditor/schema.go:56` |
| `richBlock.MinRows` | `minRows` | Go field: richBlock.MinRows<br/>type: *int | `internal/approval/formeditor/schema.go:55` |
| `richBlock.Placeholder` | `placeholder` | Go field: richBlock.Placeholder<br/>type: string | `internal/approval/formeditor/schema.go:48` |
| `richBlock.Precision` | `precision` | Go field: richBlock.Precision<br/>type: *int | `internal/approval/formeditor/schema.go:49` |
| `richBlock.Tabs` | `tabs` | Go field: richBlock.Tabs<br/>type: []richTab | `internal/approval/formeditor/schema.go:60` |
| `richBlock.Template` | `template` | Go field: richBlock.Template<br/>type: []richBlock | `internal/approval/formeditor/schema.go:54` |
| `richBlock.Type` | `type` | Go field: richBlock.Type<br/>type: string | `internal/approval/formeditor/schema.go:41` |
| `richBlock.Validate` | `validate` | Go field: richBlock.Validate<br/>type: *richValidate | `internal/approval/formeditor/schema.go:50` |
| `richDataSource.ID` | `id` | Go field: richDataSource.ID<br/>type: string | `internal/approval/formeditor/schema.go:90` |
| `richDataSource.Kind` | `kind` | Go field: richDataSource.Kind<br/>type: string | `internal/approval/formeditor/schema.go:91` |
| `richDataSource.Options` | `options` | Go field: richDataSource.Options<br/>type: []approval.FieldOption | `internal/approval/formeditor/schema.go:92` |
| `richLayer.Children` | `children` | Go field: richLayer.Children<br/>type: []richBlock | `internal/approval/formeditor/schema.go:25` |
| `richOptionSource.DataSourceID` | `dataSourceId` | Go field: richOptionSource.DataSourceID<br/>type: string | `internal/approval/formeditor/schema.go:83` |
| `richOptionSource.Kind` | `kind` | Go field: richOptionSource.Kind<br/>type: string | `internal/approval/formeditor/schema.go:81` |
| `richOptionSource.Options` | `options` | Go field: richOptionSource.Options<br/>type: []approval.FieldOption | `internal/approval/formeditor/schema.go:82` |
| `richPresentations.Mobile` | `mobile` | Go field: richPresentations.Mobile<br/>type: *richLayer | `internal/approval/formeditor/schema.go:20` |
| `richPresentations.PC` | `pc` | Go field: richPresentations.PC<br/>type: *richLayer | `internal/approval/formeditor/schema.go:19` |
| `richSchema.DataSources` | `dataSources` | Go field: richSchema.DataSources<br/>type: []richDataSource | `internal/approval/formeditor/schema.go:12` |
| `richSchema.Presentations` | `presentations` | Go field: richSchema.Presentations<br/>type: richPresentations | `internal/approval/formeditor/schema.go:13` |
| `richTab.Children` | `children` | Go field: richTab.Children<br/>type: []richBlock | `internal/approval/formeditor/schema.go:31` |
| `richValidate.Max` | `max` | Go field: richValidate.Max<br/>type: *float64 | `internal/approval/formeditor/schema.go:71` |
| `richValidate.MaxLength` | `maxLength` | Go field: richValidate.MaxLength<br/>type: *int | `internal/approval/formeditor/schema.go:69` |
| `richValidate.Message` | `message` | Go field: richValidate.Message<br/>type: string | `internal/approval/formeditor/schema.go:73` |
| `richValidate.Min` | `min` | Go field: richValidate.Min<br/>type: *float64 | `internal/approval/formeditor/schema.go:70` |
| `richValidate.MinLength` | `minLength` | Go field: richValidate.MinLength<br/>type: *int | `internal/approval/formeditor/schema.go:68` |
| `richValidate.Pattern` | `pattern` | Go field: richValidate.Pattern<br/>type: string | `internal/approval/formeditor/schema.go:72` |
| `richValidate.Required` | `required` | Go field: richValidate.Required<br/>type: *bool | `internal/approval/formeditor/schema.go:67` |
| `storedRecordKeyPart.Column` | `column` | Go field: storedRecordKeyPart.Column<br/>type: string | `internal/approval/binding/record_key.go:16` |
| `storedRecordKeyPart.Kind` | `kind` | Go field: storedRecordKeyPart.Kind<br/>type: string | `internal/approval/binding/record_key.go:17` |
| `storedRecordKeyPart.Value` | `value` | Go field: storedRecordKeyPart.Value<br/>type: string | `internal/approval/binding/record_key.go:18` |

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
| `database_query` | `database_query` | Execute a read-only (SELECT) parameterized SQL query against the database. Returns query results as JSON array. | `internal/mcp/tools/query.go:36` |

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
| `FormKeyMeta` | `meta` |  | `internal/api/router/rpc.go:21` |
| `FormKeyParams` | `params` |  | `internal/api/router/rpc.go:20` |

## auth strategy

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AuthStrategyAPIKey` | `api_key` |  | `api/auth.go:36` |
| `AuthStrategyBearer` | `bearer` |  | `api/auth.go:33` |
| `AuthStrategyHTTPBasic` | `http_basic` |  | `api/auth.go:37` |
| `AuthStrategyIP` | `ip` |  | `api/auth.go:35` |
| `AuthStrategyNone` | `none` |  | `api/auth.go:32` |
| `AuthStrategySignature` | `signature` |  | `api/auth.go:34` |

## auth type

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AuthTypeJWTToken` | `jwt_token` |  | `internal/security/jwt_token_authenticator.go:11` |
| `AuthTypeOpaqueToken` | `opaque_token` |  | `internal/security/opaque_token_authenticator.go:11` |
| `AuthTypePassword` | `password` |  | `internal/security/password_authenticator.go:15` |
| `AuthTypeRefresh` | `refresh` |  | `internal/security/jwt_refresh_authenticator.go:13` |
| `AuthTypeSignature` | `signature` |  | `internal/security/signature_authenticator.go:14` |

## built-in resource

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `approval/admin` | `rpc` |  | `internal/approval/resource/admin.go:38` |
| `approval/category` | `rpc` |  | `internal/approval/resource/category.go:117` |
| `approval/delegation` | `rpc` |  | `internal/approval/resource/delegation.go:51` |
| `approval/flow` | `rpc` |  | `internal/approval/resource/flow.go:32` |
| `approval/instance` | `rpc` |  | `internal/approval/resource/instance.go:101` |
| `approval/my` | `rpc` |  | `internal/approval/resource/my.go:29` |
| `integration/adapter` | `rpc` |  | `internal/integration/resource/adapter.go:51` |
| `integration/code_map` | `rpc` |  | `internal/integration/resource/code_map.go:61` |
| `integration/code_set` | `rpc` |  | `internal/integration/resource/code_set_catalog.go:52` |
| `integration/contract` | `rpc` |  | `internal/integration/resource/contract.go:62` |
| `integration/log` | `rpc` |  | `internal/integration/resource/log.go:32` |
| `integration/ops` | `rpc` |  | `internal/integration/resource/ops.go:80` |
| `integration/route` | `rpc` |  | `internal/integration/resource/route.go:52` |
| `integration/system` | `rpc` |  | `internal/integration/resource/system.go:128` |
| `security/auth` | `rpc` |  | `internal/security/auth_resource.go:89` |
| `sys/monitor` | `rpc` |  | `internal/monitor/resource.go:25` |
| `sys/schema` | `rpc` |  | `internal/schema/resource.go:17` |
| `sys/storage` | `rpc` |  | `internal/storage/resource.go:176` |

## built-in resource action

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `approval/admin/find_action_logs` | `find_action_logs` | permission: approval.action_log.query<br/>resource kind: rpc | `internal/approval/resource/admin.go:44` |
| `approval/admin/find_business_projections` | `find_business_projections` | permission: approval.binding.query<br/>resource kind: rpc | `internal/approval/resource/admin.go:46` |
| `approval/admin/find_instances` | `find_instances` | permission: approval.instance.query<br/>resource kind: rpc | `internal/approval/resource/admin.go:41` |
| `approval/admin/find_tasks` | `find_tasks` | permission: approval.task.query<br/>resource kind: rpc | `internal/approval/resource/admin.go:42` |
| `approval/admin/get_instance_detail` | `get_instance_detail` | permission: approval.instance.detail<br/>resource kind: rpc | `internal/approval/resource/admin.go:43` |
| `approval/admin/get_metrics` | `get_metrics` | permission: approval.metrics.query<br/>resource kind: rpc | `internal/approval/resource/admin.go:45` |
| `approval/admin/reassign_task` | `reassign_task` | audit enabled<br/>permission: approval.task.reassign<br/>resource kind: rpc | `internal/approval/resource/admin.go:50` |
| `approval/admin/retry_business_projection` | `retry_business_projection` | audit enabled<br/>permission: approval.binding.retry<br/>resource kind: rpc | `internal/approval/resource/admin.go:51` |
| `approval/admin/terminate_instance` | `terminate_instance` | audit enabled<br/>permission: approval.instance.terminate<br/>resource kind: rpc | `internal/approval/resource/admin.go:49` |
| `approval/category/create` | `create` | permission: approval.category.create<br/>resource kind: rpc | `internal/approval/resource/category.go:121` |
| `approval/category/delete` | `delete` | permission: approval.category.delete<br/>resource kind: rpc | `internal/approval/resource/category.go:147` |
| `approval/category/find_tree` | `find_tree` | permission: approval.category.query<br/>resource kind: rpc | `internal/approval/resource/category.go:118` |
| `approval/category/update` | `update` | permission: approval.category.update<br/>resource kind: rpc | `internal/approval/resource/category.go:142` |
| `approval/delegation/create` | `create` | permission: approval.delegation.create<br/>resource kind: rpc | `internal/approval/resource/delegation.go:73` |
| `approval/delegation/delete` | `delete` | permission: approval.delegation.delete<br/>resource kind: rpc | `internal/approval/resource/delegation.go:110` |
| `approval/delegation/find_page` | `find_page` | permission: approval.delegation.query<br/>resource kind: rpc | `internal/approval/resource/delegation.go:52` |
| `approval/delegation/update` | `update` | permission: approval.delegation.update<br/>resource kind: rpc | `internal/approval/resource/delegation.go:93` |
| `approval/flow/create` | `create` | audit enabled<br/>permission: approval.flow.create<br/>resource kind: rpc | `internal/approval/resource/flow.go:37` |
| `approval/flow/deploy` | `deploy` | audit enabled<br/>permission: approval.flow.deploy<br/>resource kind: rpc | `internal/approval/resource/flow.go:38` |
| `approval/flow/find_flows` | `find_flows` | permission: approval.flow.query<br/>resource kind: rpc | `internal/approval/resource/flow.go:43` |
| `approval/flow/find_initiators` | `find_initiators` | permission: approval.flow.query<br/>resource kind: rpc | `internal/approval/resource/flow.go:45` |
| `approval/flow/find_versions` | `find_versions` | permission: approval.flow.query<br/>resource kind: rpc | `internal/approval/resource/flow.go:44` |
| `approval/flow/get_graph` | `get_graph` | permission: approval.flow.query<br/>resource kind: rpc | `internal/approval/resource/flow.go:42` |
| `approval/flow/publish_version` | `publish_version` | audit enabled<br/>permission: approval.flow.publish<br/>resource kind: rpc | `internal/approval/resource/flow.go:39` |
| `approval/flow/toggle_active` | `toggle_active` | audit enabled<br/>permission: approval.flow.update<br/>resource kind: rpc | `internal/approval/resource/flow.go:41` |
| `approval/flow/update` | `update` | audit enabled<br/>permission: approval.flow.update<br/>resource kind: rpc | `internal/approval/resource/flow.go:40` |
| `approval/instance/add_assignee` | `add_assignee` | audit enabled<br/>permission: approval.task.add_assignee<br/>resource kind: rpc | `internal/approval/resource/instance.go:121` |
| `approval/instance/add_cc` | `add_cc` | audit enabled<br/>permission: approval.instance.cc<br/>resource kind: rpc | `internal/approval/resource/instance.go:118` |
| `approval/instance/mark_cc_read` | `mark_cc_read` | permission: approval.instance.cc<br/>resource kind: rpc | `internal/approval/resource/instance.go:120` |
| `approval/instance/process_task` | `process_task` | audit enabled<br/>permission: approval.task.process<br/>resource kind: rpc | `internal/approval/resource/instance.go:115` |
| `approval/instance/remove_assignee` | `remove_assignee` | audit enabled<br/>permission: approval.task.remove_assignee<br/>resource kind: rpc | `internal/approval/resource/instance.go:122` |
| `approval/instance/resubmit` | `resubmit` | audit enabled<br/>permission: approval.instance.resubmit<br/>resource kind: rpc | `internal/approval/resource/instance.go:117` |
| `approval/instance/start` | `start` | audit enabled<br/>permission: approval.instance.start<br/>resource kind: rpc | `internal/approval/resource/instance.go:114` |
| `approval/instance/urge_task` | `urge_task` | permission: approval.task.urge<br/>resource kind: rpc | `internal/approval/resource/instance.go:125` |
| `approval/instance/withdraw` | `withdraw` | audit enabled<br/>permission: approval.instance.withdraw<br/>resource kind: rpc | `internal/approval/resource/instance.go:116` |
| `approval/my/find_available_flows` | `find_available_flows` | resource kind: rpc | `internal/approval/resource/my.go:32` |
| `approval/my/find_cc_records` | `find_cc_records` | resource kind: rpc | `internal/approval/resource/my.go:37` |
| `approval/my/find_completed_tasks` | `find_completed_tasks` | resource kind: rpc | `internal/approval/resource/my.go:36` |
| `approval/my/find_initiated` | `find_initiated` | resource kind: rpc | `internal/approval/resource/my.go:34` |
| `approval/my/find_pending_tasks` | `find_pending_tasks` | resource kind: rpc | `internal/approval/resource/my.go:35` |
| `approval/my/get_instance_detail` | `get_instance_detail` | resource kind: rpc | `internal/approval/resource/my.go:39` |
| `approval/my/get_pending_counts` | `get_pending_counts` | resource kind: rpc | `internal/approval/resource/my.go:38` |
| `approval/my/get_start_form` | `get_start_form` | resource kind: rpc | `internal/approval/resource/my.go:33` |
| `integration/adapter/create` | `create` | permission: integration.adapter.create<br/>resource kind: rpc | `internal/integration/resource/adapter.go:56` |
| `integration/adapter/delete` | `delete` | permission: integration.adapter.delete<br/>resource kind: rpc | `internal/integration/resource/adapter.go:66` |
| `integration/adapter/find_all` | `find_all` | permission: integration.adapter.query<br/>resource kind: rpc | `internal/integration/resource/adapter.go:54` |
| `integration/adapter/find_page` | `find_page` | permission: integration.adapter.query<br/>resource kind: rpc | `internal/integration/resource/adapter.go:52` |
| `integration/adapter/update` | `update` | permission: integration.adapter.update<br/>resource kind: rpc | `internal/integration/resource/adapter.go:61` |
| `integration/code_map/create` | `create` | permission: integration.code_map.create<br/>resource kind: rpc | `internal/integration/resource/code_map.go:66` |
| `integration/code_map/delete` | `delete` | permission: integration.code_map.delete<br/>resource kind: rpc | `internal/integration/resource/code_map.go:76` |
| `integration/code_map/find_all` | `find_all` | permission: integration.code_map.query<br/>resource kind: rpc | `internal/integration/resource/code_map.go:64` |
| `integration/code_map/find_page` | `find_page` | permission: integration.code_map.query<br/>resource kind: rpc | `internal/integration/resource/code_map.go:62` |
| `integration/code_map/update` | `update` | permission: integration.code_map.update<br/>resource kind: rpc | `internal/integration/resource/code_map.go:71` |
| `integration/code_set/list_code_sets` | `list_code_sets` | permission: integration.code_map.query<br/>resource kind: rpc | `internal/integration/resource/code_set_catalog.go:55` |
| `integration/code_set/list_codes` | `list_codes` | permission: integration.code_map.query<br/>resource kind: rpc | `internal/integration/resource/code_set_catalog.go:56` |
| `integration/contract/create` | `create` | permission: integration.contract.create<br/>resource kind: rpc | `internal/integration/resource/contract.go:69` |
| `integration/contract/delete` | `delete` | permission: integration.contract.delete<br/>resource kind: rpc | `internal/integration/resource/contract.go:79` |
| `integration/contract/find_all` | `find_all` | permission: integration.contract.query<br/>resource kind: rpc | `internal/integration/resource/contract.go:66` |
| `integration/contract/find_page` | `find_page` | permission: integration.contract.query<br/>resource kind: rpc | `internal/integration/resource/contract.go:63` |
| `integration/contract/update` | `update` | permission: integration.contract.update<br/>resource kind: rpc | `internal/integration/resource/contract.go:74` |
| `integration/log/find_one` | `find_one` | permission: integration.log.query<br/>resource kind: rpc | `internal/integration/resource/log.go:35` |
| `integration/log/find_page` | `find_page` | permission: integration.log.query<br/>resource kind: rpc | `internal/integration/resource/log.go:33` |
| `integration/ops/diagnose_routes` | `diagnose_routes` | permission: integration.ops.diagnose_routes<br/>resource kind: rpc | `internal/integration/resource/ops.go:86` |
| `integration/ops/dry_run` | `dry_run` | permission: integration.ops.dry_run<br/>resource kind: rpc | `internal/integration/resource/ops.go:83` |
| `integration/ops/dry_run_inbound` | `dry_run_inbound` | permission: integration.ops.dry_run_inbound<br/>resource kind: rpc | `internal/integration/resource/ops.go:84` |
| `integration/ops/test_connection` | `test_connection` | permission: integration.ops.test_connection<br/>resource kind: rpc | `internal/integration/resource/ops.go:85` |
| `integration/route/create` | `create` | permission: integration.route.create<br/>resource kind: rpc | `internal/integration/resource/route.go:57` |
| `integration/route/delete` | `delete` | permission: integration.route.delete<br/>resource kind: rpc | `internal/integration/resource/route.go:67` |
| `integration/route/find_all` | `find_all` | permission: integration.route.query<br/>resource kind: rpc | `internal/integration/resource/route.go:55` |
| `integration/route/find_page` | `find_page` | permission: integration.route.query<br/>resource kind: rpc | `internal/integration/resource/route.go:53` |
| `integration/route/update` | `update` | permission: integration.route.update<br/>resource kind: rpc | `internal/integration/resource/route.go:62` |
| `integration/system/create` | `create` | permission: integration.system.create<br/>resource kind: rpc | `internal/integration/resource/system.go:135` |
| `integration/system/delete` | `delete` | permission: integration.system.delete<br/>resource kind: rpc | `internal/integration/resource/system.go:152` |
| `integration/system/find_all` | `find_all` | permission: integration.system.query<br/>resource kind: rpc | `internal/integration/resource/system.go:132` |
| `integration/system/find_page` | `find_page` | permission: integration.system.query<br/>resource kind: rpc | `internal/integration/resource/system.go:129` |
| `integration/system/update` | `update` | permission: integration.system.update<br/>resource kind: rpc | `internal/integration/resource/system.go:140` |
| `sys/monitor/get_build_info` | `get_build_info` | resource kind: rpc | `internal/monitor/resource.go:36` |
| `sys/monitor/get_cpu` | `get_cpu` | resource kind: rpc | `internal/monitor/resource.go:29` |
| `sys/monitor/get_disk` | `get_disk` | resource kind: rpc | `internal/monitor/resource.go:31` |
| `sys/monitor/get_event_streams` | `get_event_streams` | resource kind: rpc | `internal/monitor/resource.go:37` |
| `sys/monitor/get_host` | `get_host` | resource kind: rpc | `internal/monitor/resource.go:33` |
| `sys/monitor/get_integration_stats` | `get_integration_stats` | resource kind: rpc | `internal/monitor/resource.go:38` |
| `sys/monitor/get_load` | `get_load` | resource kind: rpc | `internal/monitor/resource.go:35` |
| `sys/monitor/get_memory` | `get_memory` | resource kind: rpc | `internal/monitor/resource.go:30` |
| `sys/monitor/get_network` | `get_network` | resource kind: rpc | `internal/monitor/resource.go:32` |
| `sys/monitor/get_overview` | `get_overview` | resource kind: rpc | `internal/monitor/resource.go:28` |
| `sys/monitor/get_process` | `get_process` | resource kind: rpc | `internal/monitor/resource.go:34` |
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
| `vef.api.rate_limit.max` | `100 (100)` |  | `config/api.go:34` |
| `vef.api.rate_limit.period` | `5 * time.Minute` |  | `config/api.go:39` |
| `vef.approval.business_binding.batch_size` | `100` |  | `config/approval.go:57` |
| `vef.approval.business_binding.consistency` | `"synchronous" (synchronous)` |  | `config/approval.go:43` |
| `vef.approval.business_binding.scan_interval` | `10 * time.Second` |  | `config/approval.go:52` |
| `vef.approval.cc_record_retention` | `90 * 24 * time.Hour` |  | `config/approval.go:152` |
| `vef.approval.cleanup_scan_interval` | `24 * time.Hour` |  | `config/approval.go:136` |
| `vef.approval.delegation_max_depth` | `10` |  | `config/approval.go:140` |
| `vef.approval.form_snapshot_retention` | `90 * 24 * time.Hour` |  | `config/approval.go:144` |
| `vef.approval.pre_warning_scan_interval` | `5 * time.Minute` |  | `config/approval.go:132` |
| `vef.approval.timeout_scan_interval` | `time.Minute` |  | `config/approval.go:128` |
| `vef.approval.urge_record_retention` | `30 * 24 * time.Hour` |  | `config/approval.go:148` |
| `vef.cron.store.abandoned_after` | `time.Minute` |  | `config/cron.go:101` |
| `vef.cron.store.batch_size` | `32 (32)` |  | `config/cron.go:81` |
| `vef.cron.store.heartbeat_interval` | `10 * time.Second` |  | `config/cron.go:96` |
| `vef.cron.store.max_concurrent` | `16 (16)` |  | `config/cron.go:86` |
| `vef.cron.store.misfire_threshold` | `time.Minute` |  | `config/cron.go:91` |
| `vef.cron.store.poll_interval` | `5 * time.Second` |  | `config/cron.go:76` |
| `vef.event.async_queue_size` | `4096` |  | `config/event.go:133` |
| `vef.event.async_workers` | `4` |  | `config/event.go:138` |
| `vef.event.default_transport` | `memory` |  | `config/event.go:128` |
| `vef.event.inbox.cleanup_interval` | `time.Hour` |  | `config/event.go:168` |
| `vef.event.inbox.processing_lease` | `10 * time.Minute` |  | `config/event.go:163` |
| `vef.event.inbox.retention` | `7 * 24 * time.Hour` |  | `config/event.go:158` |
| `vef.event.publish_timeout` | `5 * time.Second` |  | `config/event.go:143` |
| `vef.event.transports.memory.full_policy` | `"error" (error)` |  | `event/transport/memory/memory.go:48` |
| `vef.event.transports.memory.queue_size` | `1024` |  | `event/transport/memory/memory.go:39` |
| `vef.event.transports.outbox.batch_size` | `100` |  | `event/transport/outbox/outbox.go:121` |
| `vef.event.transports.outbox.cleanup_interval` | `time.Hour` |  | `config/event.go:148` |
| `vef.event.transports.outbox.completed_ttl` | `7 * 24 * time.Hour` |  | `config/event.go:153` |
| `vef.event.transports.outbox.lease_multiplier` | `4` |  | `event/transport/outbox/outbox.go:130` |
| `vef.event.transports.outbox.max_retries` | `10` | EventConfig.Validate fallback when max_retries is unset | `config/event.go:175` |
| `vef.event.transports.outbox.max_retries` | `10` |  | `event/transport/outbox/outbox.go:112` |
| `vef.event.transports.outbox.min_lease` | `15 * time.Second` |  | `event/transport/outbox/outbox.go:139` |
| `vef.event.transports.outbox.relay_interval` | `10 * time.Second` |  | `event/transport/outbox/outbox.go:103` |
| `vef.event.transports.outbox.sink` | `memory` |  | `event/transport/outbox/outbox.go:148` |
| `vef.event.transports.redis_stream.block_timeout` | `5 * time.Second` |  | `event/transport/redisstream/redis_stream.go:80` |
| `vef.event.transports.redis_stream.claim_batch_size` | `64` |  | `event/transport/redisstream/redis_stream.go:107` |
| `vef.event.transports.redis_stream.claim_idle` | `60 * time.Second` |  | `event/transport/redisstream/redis_stream.go:89` |
| `vef.event.transports.redis_stream.claim_interval` | `30 * time.Second` |  | `event/transport/redisstream/redis_stream.go:98` |
| `vef.event.transports.redis_stream.idle_group_sweep_interval` | `10 * time.Minute` |  | `event/transport/redisstream/redis_stream.go:143` |
| `vef.event.transports.redis_stream.reaper_concurrency` | `4` |  | `event/transport/redisstream/redis_stream.go:116` |
| `vef.event.transports.redis_stream.setup_timeout` | `5 * time.Second` |  | `event/transport/redisstream/redis_stream.go:125` |
| `vef.event.transports.redis_stream.start_id` | `0` |  | `event/transport/redisstream/redis_stream.go:134` |
| `vef.event.transports.redis_stream.stream_prefix` | `vef:events:` |  | `event/transport/redisstream/redis_stream.go:71` |
| `vef.integration.inbound.rate_limit.max` | `120 (120)` |  | `config/integration.go:134` |
| `vef.integration.inbound.rate_limit.period` | `time.Minute` |  | `config/integration.go:139` |
| `vef.integration.log.capture_limit` | `4096` |  | `config/integration.go:51` |
| `vef.integration.log.mode` | `"errors" (errors)` |  | `config/integration.go:42` |
| `vef.integration.max_response_body` | `8 &lt;&lt; 20` |  | `config/integration.go:158` |
| `vef.integration.run_timeout` | `30 * time.Second` |  | `config/integration.go:153` |
| `vef.integration.secret_algorithm` | `"aes" (aes)` |  | `config/integration.go:144` |
| `vef.mcp.require_auth` | `true when unset` |  | `internal/mcp/handler.go:34` |
| `vef.monitor.sample_duration` | `2 * time.Second` |  | `internal/monitor/config.go:25` |
| `vef.monitor.sample_interval` | `10 * time.Second` |  | `internal/monitor/config.go:24` |
| `vef.push.path` | `/ws` |  | `config/push.go:35` |
| `vef.push.ping_interval` | `30 * time.Second` |  | `config/push.go:44` |
| `vef.push.send_buffer` | `32` |  | `config/push.go:54` |
| `vef.push.session_recheck_interval` | `60 * time.Second` |  | `config/push.go:60` |
| `vef.push.write_timeout` | `10 * time.Second` |  | `config/push.go:49` |
| `vef.security.lockout.backoff_base` | `1 * time.Second` |  | `config/security.go:309` |
| `vef.security.lockout.backoff_max` | `15 * time.Minute` |  | `config/security.go:314` |
| `vef.security.lockout.key` | `"user_ip" (user_ip)` |  | `config/security.go:319` |
| `vef.security.lockout.lock_duration` | `15 * time.Minute` |  | `config/security.go:295` |
| `vef.security.lockout.max_failures` | `10 (10)` |  | `config/security.go:285` |
| `vef.security.lockout.strategy` | `"lock" (lock)` |  | `config/security.go:300` |
| `vef.security.lockout.window` | `15 * time.Minute` |  | `config/security.go:290` |
| `vef.security.session.idle_ttl` | `30 * time.Minute` |  | `config/security.go:145` |
| `vef.security.session.max_lifetime` | `7 * 24 * time.Hour` |  | `config/security.go:150` |
| `vef.security.session.on_exceed` | `"evict_oldest" (evict_oldest)` |  | `config/security.go:136` |
| `vef.security.token_type` | `"jwt_token" (jwt_token)` |  | `config/security.go:88` |
| `vef.storage.claim_ttl` | `24 * time.Hour` |  | `config/storage.go:118` |
| `vef.storage.delete_batch_size` | `100 (100)` |  | `config/storage.go:143` |
| `vef.storage.delete_concurrency` | `8 (8)` |  | `config/storage.go:148` |
| `vef.storage.delete_lease_window` | `5 * time.Minute` |  | `config/storage.go:158` |
| `vef.storage.delete_max_attempts` | `12 (12)` |  | `config/storage.go:153` |
| `vef.storage.delete_worker_interval` | `5 * time.Minute` |  | `config/storage.go:138` |
| `vef.storage.max_pending_claims` | `100 (100)` |  | `config/storage.go:123` |
| `vef.storage.max_upload_size` | `1024 * 1024 * 1024` |  | `config/storage.go:113` |
| `vef.storage.sweep_batch_size` | `200 (200)` |  | `config/storage.go:133` |
| `vef.storage.sweep_interval` | `5 * time.Minute` |  | `config/storage.go:128` |

## config enum

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `MySQL` | `mysql` |  | `config/data_sources.go:19` |
| `Oracle` | `oracle` |  | `config/data_sources.go:16` |
| `Postgres` | `postgres` |  | `config/data_sources.go:18` |
| `SQLServer` | `sqlserver` |  | `config/data_sources.go:17` |
| `SQLite` | `sqlite` |  | `config/data_sources.go:20` |
| `StorageFilesystem` | `filesystem` |  | `config/storage.go:15` |
| `StorageMemory` | `memory` |  | `config/storage.go:14` |
| `StorageMinIO` | `minio` |  | `config/storage.go:13` |

## config key

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `vef.api.rate_limit` | `APIRateLimitConfig` | Go field: APIConfig.RateLimit | `config/api.go:17` |
| `vef.api.rate_limit.max` | `int` | Go field: APIRateLimitConfig.Max | `config/api.go:17`, `config/api.go:27` |
| `vef.api.rate_limit.period` | `time.Duration` | Go field: APIRateLimitConfig.Period | `config/api.go:17`, `config/api.go:30` |
| `vef.app.body_limit` | `string` | Go field: AppConfig.BodyLimit | `config/app.go:7` |
| `vef.app.name` | `string` | Go field: AppConfig.Name | `config/app.go:5` |
| `vef.app.port` | `uint16` | Go field: AppConfig.Port | `config/app.go:6` |
| `vef.app.trusted_proxies` | `[]string` | Go field: AppConfig.TrustedProxies | `config/app.go:11` |
| `vef.approval.auto_migrate` | `bool` | Go field: ApprovalConfig.AutoMigrate | `config/approval.go:89` |
| `vef.approval.business_binding` | `ApprovalBusinessBindingConfig` | Go field: ApprovalConfig.BusinessBinding | `config/approval.go:121` |
| `vef.approval.business_binding.batch_size` | `int` | Go field: ApprovalBusinessBindingConfig.BatchSize | `config/approval.go:121`, `config/approval.go:39` |
| `vef.approval.business_binding.consistency` | `ApprovalBindingConsistency` | Go field: ApprovalBusinessBindingConfig.Consistency | `config/approval.go:121`, `config/approval.go:34` |
| `vef.approval.business_binding.scan_interval` | `time.Duration` | Go field: ApprovalBusinessBindingConfig.ScanInterval | `config/approval.go:121`, `config/approval.go:36` |
| `vef.approval.cc_record_retention` | `time.Duration` | Go field: ApprovalConfig.CCRecordRetention | `config/approval.go:118` |
| `vef.approval.cleanup_scan_interval` | `time.Duration` | Go field: ApprovalConfig.CleanupScanInterval | `config/approval.go:102` |
| `vef.approval.delegation_max_depth` | `int` | Go field: ApprovalConfig.DelegationMaxDepth | `config/approval.go:106` |
| `vef.approval.form_snapshot_retention` | `time.Duration` | Go field: ApprovalConfig.FormSnapshotRetention | `config/approval.go:110` |
| `vef.approval.pre_warning_scan_interval` | `time.Duration` | Go field: ApprovalConfig.PreWarningScanInterval | `config/approval.go:97` |
| `vef.approval.timeout_scan_interval` | `time.Duration` | Go field: ApprovalConfig.TimeoutScanInterval | `config/approval.go:93` |
| `vef.approval.urge_record_retention` | `time.Duration` | Go field: ApprovalConfig.UrgeRecordRetention | `config/approval.go:114` |
| `vef.cors.allow_origins` | `[]string` | Go field: CORSConfig.AllowOrigins | `config/cors.go:6` |
| `vef.cors.enabled` | `bool` | Go field: CORSConfig.Enabled | `config/cors.go:5` |
| `vef.cron.store` | `CronStoreConfig` | Go field: CronConfig.Store | `config/cron.go:22` |
| `vef.cron.store.abandoned_after` | `time.Duration` | Go field: CronStoreConfig.AbandonedAfter | `config/cron.go:22`, `config/cron.go:62` |
| `vef.cron.store.auto_migrate` | `bool` | Go field: CronStoreConfig.AutoMigrate | `config/cron.go:22`, `config/cron.go:36` |
| `vef.cron.store.batch_size` | `int` | Go field: CronStoreConfig.BatchSize | `config/cron.go:22`, `config/cron.go:45` |
| `vef.cron.store.enabled` | `bool` | Go field: CronStoreConfig.Enabled | `config/cron.go:22`, `config/cron.go:33` |
| `vef.cron.store.heartbeat_interval` | `time.Duration` | Go field: CronStoreConfig.HeartbeatInterval | `config/cron.go:22`, `config/cron.go:57` |
| `vef.cron.store.max_concurrent` | `int` | Go field: CronStoreConfig.MaxConcurrent | `config/cron.go:22`, `config/cron.go:49` |
| `vef.cron.store.misfire_threshold` | `time.Duration` | Go field: CronStoreConfig.MisfireThreshold | `config/cron.go:22`, `config/cron.go:53` |
| `vef.cron.store.poll_interval` | `time.Duration` | Go field: CronStoreConfig.PollInterval | `config/cron.go:22`, `config/cron.go:42` |
| `vef.cron.store.run_retention` | `time.Duration` | Go field: CronStoreConfig.RunRetention | `config/cron.go:22`, `config/cron.go:72` |
| `vef.cron.store.run_timeout` | `time.Duration` | Go field: CronStoreConfig.RunTimeout | `config/cron.go:22`, `config/cron.go:67` |
| `vef.data_sources.&lt;name&gt;.database` | `string` | Go field: DataSourceConfig.Database | `config/data_sources.go:22`, `config/data_sources.go:57`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.enable_sql_guard` | `bool` | Go field: DataSourceConfig.EnableSQLGuard | `config/data_sources.go:22`, `config/data_sources.go:60`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.host` | `string` | Go field: DataSourceConfig.Host | `config/data_sources.go:22`, `config/data_sources.go:53`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.password` | `string` | Go field: DataSourceConfig.Password | `config/data_sources.go:22`, `config/data_sources.go:56`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.path` | `string` | Go field: DataSourceConfig.Path | `config/data_sources.go:22`, `config/data_sources.go:59`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.port` | `uint16` | Go field: DataSourceConfig.Port | `config/data_sources.go:22`, `config/data_sources.go:54`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.schema` | `string` | Go field: DataSourceConfig.Schema | `config/data_sources.go:22`, `config/data_sources.go:58`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.ssl_mode` | `SSLMode` | Go field: DataSourceConfig.SSLMode | `config/data_sources.go:22`, `config/data_sources.go:65`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.ssl_root_cert` | `string` | Go field: DataSourceConfig.SSLRootCert | `config/data_sources.go:22`, `config/data_sources.go:69`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.type` | `DBKind` | Go field: DataSourceConfig.Kind | `config/data_sources.go:22`, `config/data_sources.go:52`, `internal/config/data_sources.go:20` |
| `vef.data_sources.&lt;name&gt;.user` | `string` | Go field: DataSourceConfig.User | `config/data_sources.go:22`, `config/data_sources.go:55`, `internal/config/data_sources.go:20` |
| `vef.event.async_queue_size` | `int` | Go field: EventConfig.AsyncQueueSize | `config/event.go:18` |
| `vef.event.async_workers` | `int` | Go field: EventConfig.AsyncWorkers | `config/event.go:21` |
| `vef.event.default_transport` | `string` | Go field: EventConfig.DefaultTransport | `config/event.go:15` |
| `vef.event.inbox` | `EventInboxConfig` | Go field: EventConfig.Inbox | `config/event.go:27` |
| `vef.event.inbox.cleanup_interval` | `time.Duration` | Go field: EventInboxConfig.CleanupInterval | `config/event.go:115`, `config/event.go:27` |
| `vef.event.inbox.processing_lease` | `time.Duration` | Go field: EventInboxConfig.ProcessingLease | `config/event.go:114`, `config/event.go:27` |
| `vef.event.inbox.retention` | `time.Duration` | Go field: EventInboxConfig.Retention | `config/event.go:113`, `config/event.go:27` |
| `vef.event.middleware` | `EventMiddlewareConfig` | Go field: EventConfig.Middleware | `config/event.go:26` |
| `vef.event.middleware.inbox` | `bool` | Go field: EventMiddlewareConfig.Inbox | `config/event.go:107`, `config/event.go:26` |
| `vef.event.middleware.logging` | `bool` | Go field: EventMiddlewareConfig.Logging | `config/event.go:26`, `config/event.go:93` |
| `vef.event.middleware.metrics` | `bool` | Go field: EventMiddlewareConfig.Metrics | `config/event.go:102`, `config/event.go:26` |
| `vef.event.middleware.recover` | `bool` | Go field: EventMiddlewareConfig.Recover | `config/event.go:103`, `config/event.go:26` |
| `vef.event.middleware.tracing` | `bool` | Go field: EventMiddlewareConfig.Tracing | `config/event.go:26`, `config/event.go:94` |
| `vef.event.middleware.tracing_strict` | `bool` | Go field: EventMiddlewareConfig.TracingStrict | `config/event.go:101`, `config/event.go:26` |
| `vef.event.publish_timeout` | `time.Duration` | Go field: EventConfig.PublishTimeout | `config/event.go:23` |
| `vef.event.routing` | `[]EventRoutingRule` | Go field: EventConfig.Routing | `config/event.go:31` |
| `vef.event.routing.pattern` | `string` | Go field: EventRoutingRule.Pattern | `config/event.go:123`, `config/event.go:31` |
| `vef.event.routing.transports` | `[]string` | Go field: EventRoutingRule.Transports | `config/event.go:124`, `config/event.go:31` |
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
| `vef.event.transports.redis_stream.consumer_id` | `string` | Go field: EventRedisStreamTransportConfig.ConsumerID | `config/event.go:25`, `config/event.go:38`, `config/event.go:79` |
| `vef.event.transports.redis_stream.enabled` | `bool` | Go field: EventRedisStreamTransportConfig.Enabled | `config/event.go:25`, `config/event.go:38`, `config/event.go:63` |
| `vef.event.transports.redis_stream.handler_timeout` | `time.Duration` | Go field: EventRedisStreamTransportConfig.HandlerTimeout | `config/event.go:25`, `config/event.go:38`, `config/event.go:77` |
| `vef.event.transports.redis_stream.idle_group_retention` | `time.Duration` | Go field: EventRedisStreamTransportConfig.IdleGroupRetention | `config/event.go:25`, `config/event.go:38`, `config/event.go:85` |
| `vef.event.transports.redis_stream.idle_group_sweep_interval` | `time.Duration` | Go field: EventRedisStreamTransportConfig.IdleGroupSweepInterval | `config/event.go:25`, `config/event.go:38`, `config/event.go:88` |
| `vef.event.transports.redis_stream.max_len_approx` | `int64` | Go field: EventRedisStreamTransportConfig.MaxLenApprox | `config/event.go:25`, `config/event.go:38`, `config/event.go:65` |
| `vef.event.transports.redis_stream.reaper_concurrency` | `int` | Go field: EventRedisStreamTransportConfig.ReaperConcurrency | `config/event.go:25`, `config/event.go:38`, `config/event.go:70` |
| `vef.event.transports.redis_stream.setup_timeout` | `time.Duration` | Go field: EventRedisStreamTransportConfig.SetupTimeout | `config/event.go:25`, `config/event.go:38`, `config/event.go:78` |
| `vef.event.transports.redis_stream.start_id` | `string` | Go field: EventRedisStreamTransportConfig.StartID | `config/event.go:25`, `config/event.go:38`, `config/event.go:80` |
| `vef.event.transports.redis_stream.stream_prefix` | `string` | Go field: EventRedisStreamTransportConfig.StreamPrefix | `config/event.go:25`, `config/event.go:38`, `config/event.go:64` |
| `vef.integration.auto_migrate` | `bool` | Go field: IntegrationConfig.AutoMigrate | `config/integration.go:76` |
| `vef.integration.inbound` | `IntegrationInboundConfig` | Go field: IntegrationConfig.Inbound | `config/integration.go:103` |
| `vef.integration.inbound.rate_limit` | `IntegrationInboundRateLimitConfig` | Go field: IntegrationInboundConfig.RateLimit | `config/integration.go:103`, `config/integration.go:117` |
| `vef.integration.inbound.rate_limit.max` | `int` | Go field: IntegrationInboundRateLimitConfig.Max | `config/integration.go:103`, `config/integration.go:117`, `config/integration.go:127` |
| `vef.integration.inbound.rate_limit.period` | `time.Duration` | Go field: IntegrationInboundRateLimitConfig.Period | `config/integration.go:103`, `config/integration.go:117`, `config/integration.go:130` |
| `vef.integration.log` | `IntegrationLogConfig` | Go field: IntegrationConfig.Log | `config/integration.go:99` |
| `vef.integration.log.capture_limit` | `int` | Go field: IntegrationLogConfig.CaptureLimit | `config/integration.go:31`, `config/integration.go:99` |
| `vef.integration.log.mask_fields` | `[]string` | Go field: IntegrationLogConfig.MaskFields | `config/integration.go:34`, `config/integration.go:99` |
| `vef.integration.log.mode` | `IntegrationLogMode` | Go field: IntegrationLogConfig.Mode | `config/integration.go:28`, `config/integration.go:99` |
| `vef.integration.log.retention` | `time.Duration` | Go field: IntegrationLogConfig.Retention | `config/integration.go:38`, `config/integration.go:99` |
| `vef.integration.max_response_body` | `int64` | Go field: IntegrationConfig.MaxResponseBody | `config/integration.go:96` |
| `vef.integration.run_timeout` | `time.Duration` | Go field: IntegrationConfig.RunTimeout | `config/integration.go:92` |
| `vef.integration.secret_algorithm` | `IntegrationSecretAlgorithm` | Go field: IntegrationConfig.SecretAlgorithm | `config/integration.go:88` |
| `vef.integration.secret_key` | `string` | Go field: IntegrationConfig.SecretKey | `config/integration.go:82` |
| `vef.mcp.enabled` | `bool` | Go field: MCPConfig.Enabled | `config/mcp.go:5` |
| `vef.mcp.require_auth` | `*bool` | Go field: MCPConfig.RequireAuth | `config/mcp.go:9` |
| `vef.monitor.sample_duration` | `time.Duration` | Go field: MonitorConfig.SampleDuration | `config/monitor.go:8` |
| `vef.monitor.sample_interval` | `time.Duration` | Go field: MonitorConfig.SampleInterval | `config/monitor.go:7` |
| `vef.push.allowed_origins` | `[]string` | Go field: PushConfig.AllowedOrigins | `config/push.go:16` |
| `vef.push.enabled` | `bool` | Go field: PushConfig.Enabled | `config/push.go:10` |
| `vef.push.max_connections_per_user` | `int` | Go field: PushConfig.MaxConnectionsPerUser | `config/push.go:27` |
| `vef.push.path` | `string` | Go field: PushConfig.Path | `config/push.go:12` |
| `vef.push.ping_interval` | `time.Duration` | Go field: PushConfig.PingInterval | `config/push.go:19` |
| `vef.push.send_buffer` | `int` | Go field: PushConfig.SendBuffer | `config/push.go:24` |
| `vef.push.session_recheck_interval` | `time.Duration` | Go field: PushConfig.SessionRecheckInterval | `config/push.go:31` |
| `vef.push.write_timeout` | `time.Duration` | Go field: PushConfig.WriteTimeout | `config/push.go:21` |
| `vef.redis.database` | `uint8` | Go field: RedisConfig.Database | `config/redis.go:16` |
| `vef.redis.enabled` | `bool` | Go field: RedisConfig.Enabled | `config/redis.go:11` |
| `vef.redis.host` | `string` | Go field: RedisConfig.Host | `config/redis.go:12` |
| `vef.redis.network` | `string` | Go field: RedisConfig.Network | `config/redis.go:17` |
| `vef.redis.password` | `string` | Go field: RedisConfig.Password | `config/redis.go:15` |
| `vef.redis.port` | `uint16` | Go field: RedisConfig.Port | `config/redis.go:13` |
| `vef.redis.user` | `string` | Go field: RedisConfig.User | `config/redis.go:14` |
| `vef.security.api_keys` | `map[string]APIKeyConfig` | Go field: SecurityConfig.APIKeys | `config/security.go:38` |
| `vef.security.api_keys.&lt;name&gt;.key` | `string` | Go field: APIKeyConfig.Key | `config/security.go:38`, `config/security.go:62` |
| `vef.security.api_keys.&lt;name&gt;.roles` | `[]string` | Go field: APIKeyConfig.Roles | `config/security.go:38`, `config/security.go:64` |
| `vef.security.basic_accounts` | `map[string]BasicAccountConfig` | Go field: SecurityConfig.BasicAccounts | `config/security.go:44` |
| `vef.security.basic_accounts.&lt;name&gt;.password` | `string` | Go field: BasicAccountConfig.Password | `config/security.go:44`, `config/security.go:72` |
| `vef.security.basic_accounts.&lt;name&gt;.roles` | `[]string` | Go field: BasicAccountConfig.Roles | `config/security.go:44`, `config/security.go:74` |
| `vef.security.ip_whitelists` | `map[string][]string` | Go field: SecurityConfig.IPWhitelists | `config/security.go:33` |
| `vef.security.lockout` | `LockoutConfig` | Go field: SecurityConfig.Lockout | `config/security.go:46` |
| `vef.security.lockout.backoff_base` | `time.Duration` | Go field: LockoutConfig.BackoffBase | `config/security.go:271`, `config/security.go:46` |
| `vef.security.lockout.backoff_max` | `time.Duration` | Go field: LockoutConfig.BackoffMax | `config/security.go:273`, `config/security.go:46` |
| `vef.security.lockout.enabled` | `*bool` | Go field: LockoutConfig.Enabled | `config/security.go:256`, `config/security.go:46` |
| `vef.security.lockout.key` | `LockoutKey` | Go field: LockoutConfig.Key | `config/security.go:276`, `config/security.go:46` |
| `vef.security.lockout.lock_duration` | `time.Duration` | Go field: LockoutConfig.LockDuration | `config/security.go:265`, `config/security.go:46` |
| `vef.security.lockout.max_failures` | `int` | Go field: LockoutConfig.MaxFailures | `config/security.go:259`, `config/security.go:46` |
| `vef.security.lockout.strategy` | `LockoutStrategy` | Go field: LockoutConfig.Strategy | `config/security.go:268`, `config/security.go:46` |
| `vef.security.lockout.window` | `time.Duration` | Go field: LockoutConfig.Window | `config/security.go:262`, `config/security.go:46` |
| `vef.security.login_rate_limit` | `int` | Go field: SecurityConfig.LoginRateLimit | `config/security.go:25` |
| `vef.security.password_policy` | `PasswordPolicyConfig` | Go field: SecurityConfig.PasswordPolicy | `config/security.go:48` |
| `vef.security.password_policy.blocklist` | `[]string` | Go field: PasswordPolicyConfig.Blocklist | `config/security.go:200`, `config/security.go:48` |
| `vef.security.password_policy.disallow_username` | `bool` | Go field: PasswordPolicyConfig.DisallowUsername | `config/security.go:198`, `config/security.go:48` |
| `vef.security.password_policy.history_depth` | `int` | Go field: PasswordPolicyConfig.HistoryDepth | `config/security.go:203`, `config/security.go:48` |
| `vef.security.password_policy.max_age` | `time.Duration` | Go field: PasswordPolicyConfig.MaxAge | `config/security.go:207`, `config/security.go:48` |
| `vef.security.password_policy.max_length` | `int` | Go field: PasswordPolicyConfig.MaxLength | `config/security.go:184`, `config/security.go:48` |
| `vef.security.password_policy.min_char_classes` | `int` | Go field: PasswordPolicyConfig.MinCharClasses | `config/security.go:196`, `config/security.go:48` |
| `vef.security.password_policy.min_length` | `int` | Go field: PasswordPolicyConfig.MinLength | `config/security.go:182`, `config/security.go:48` |
| `vef.security.password_policy.require_digit` | `bool` | Go field: PasswordPolicyConfig.RequireDigit | `config/security.go:190`, `config/security.go:48` |
| `vef.security.password_policy.require_lower` | `bool` | Go field: PasswordPolicyConfig.RequireLower | `config/security.go:188`, `config/security.go:48` |
| `vef.security.password_policy.require_symbol` | `bool` | Go field: PasswordPolicyConfig.RequireSymbol | `config/security.go:193`, `config/security.go:48` |
| `vef.security.password_policy.require_upper` | `bool` | Go field: PasswordPolicyConfig.RequireUpper | `config/security.go:186`, `config/security.go:48` |
| `vef.security.refresh_not_before` | `time.Duration` | Go field: SecurityConfig.RefreshNotBefore | `config/security.go:24` |
| `vef.security.refresh_rate_limit` | `int` | Go field: SecurityConfig.RefreshRateLimit | `config/security.go:26` |
| `vef.security.secret` | `string` | Go field: SecurityConfig.Secret | `config/security.go:22` |
| `vef.security.session` | `SessionConfig` | Go field: SecurityConfig.Session | `config/security.go:55` |
| `vef.security.session.idle_ttl` | `time.Duration` | Go field: SessionConfig.IdleTTL | `config/security.go:126`, `config/security.go:55` |
| `vef.security.session.max_concurrent` | `int` | Go field: SessionConfig.MaxConcurrent | `config/security.go:120`, `config/security.go:55` |
| `vef.security.session.max_lifetime` | `time.Duration` | Go field: SessionConfig.MaxLifetime | `config/security.go:129`, `config/security.go:55` |
| `vef.security.session.on_exceed` | `SessionExceedPolicy` | Go field: SessionConfig.OnExceed | `config/security.go:123`, `config/security.go:55` |
| `vef.security.session.sliding` | `*bool` | Go field: SessionConfig.Sliding | `config/security.go:132`, `config/security.go:55` |
| `vef.security.token_expires` | `time.Duration` | Go field: SecurityConfig.TokenExpires | `config/security.go:23` |
| `vef.security.token_type` | `TokenType` | Go field: SecurityConfig.TokenType | `config/security.go:52` |
| `vef.storage.allow_public_uploads` | `bool` | Go field: StorageConfig.AllowPublicUploads | `config/storage.go:53` |
| `vef.storage.auto_migrate` | `bool` | Go field: StorageConfig.AutoMigrate | `config/storage.go:26` |
| `vef.storage.claim_ttl` | `time.Duration` | Go field: StorageConfig.ClaimTTL | `config/storage.go:37` |
| `vef.storage.delete_batch_size` | `int` | Go field: StorageConfig.DeleteBatchSize | `config/storage.go:67` |
| `vef.storage.delete_concurrency` | `int` | Go field: StorageConfig.DeleteConcurrency | `config/storage.go:70` |
| `vef.storage.delete_lease_window` | `time.Duration` | Go field: StorageConfig.DeleteLeaseWindow | `config/storage.go:77` |
| `vef.storage.delete_max_attempts` | `int` | Go field: StorageConfig.DeleteMaxAttempts | `config/storage.go:73` |
| `vef.storage.delete_worker_interval` | `time.Duration` | Go field: StorageConfig.DeleteWorkerInterval | `config/storage.go:64` |
| `vef.storage.filesystem` | `FilesystemConfig` | Go field: StorageConfig.Filesystem | `config/storage.go:28` |
| `vef.storage.filesystem.root` | `string` | Go field: FilesystemConfig.Root | `config/storage.go:174`, `config/storage.go:28` |
| `vef.storage.max_pending_claims` | `int` | Go field: StorageConfig.MaxPendingClaims | `config/storage.go:47` |
| `vef.storage.max_upload_size` | `int64` | Go field: StorageConfig.MaxUploadSize | `config/storage.go:33` |
| `vef.storage.minio` | `MinIOConfig` | Go field: StorageConfig.MinIO | `config/storage.go:27` |
| `vef.storage.minio.access_key` | `string` | Go field: MinIOConfig.AccessKey | `config/storage.go:165`, `config/storage.go:27` |
| `vef.storage.minio.bucket` | `string` | Go field: MinIOConfig.Bucket | `config/storage.go:167`, `config/storage.go:27` |
| `vef.storage.minio.endpoint` | `string` | Go field: MinIOConfig.Endpoint | `config/storage.go:164`, `config/storage.go:27` |
| `vef.storage.minio.region` | `string` | Go field: MinIOConfig.Region | `config/storage.go:168`, `config/storage.go:27` |
| `vef.storage.minio.secret_key` | `string` | Go field: MinIOConfig.SecretKey | `config/storage.go:166`, `config/storage.go:27` |
| `vef.storage.minio.use_ssl` | `bool` | Go field: MinIOConfig.UseSSL | `config/storage.go:169`, `config/storage.go:27` |
| `vef.storage.provider` | `StorageProvider` | Go field: StorageConfig.Provider | `config/storage.go:25` |
| `vef.storage.sweep_batch_size` | `int` | Go field: StorageConfig.SweepBatchSize | `config/storage.go:60` |
| `vef.storage.sweep_interval` | `time.Duration` | Go field: StorageConfig.SweepInterval | `config/storage.go:57` |

## config reserved name

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `PrimaryDataSourceName` | `primary` | used under vef.data_sources.&lt;name&gt; | `config/data_sources.go:7` |

## environment variable

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `EnvConfigPath` | `VEF_CONFIG_PATH` |  | `config/env.go:7` |
| `EnvI18NLanguage` | `VEF_I18N_LANGUAGE` |  | `config/env.go:8` |
| `EnvLogLevel` | `VEF_LOG_LEVEL` |  | `config/env.go:6` |
| `EnvPrefix` | `VEF` |  | `config/env.go:5` |
| `SystemDrive` | `SystemDrive` | OS-defined Windows variable read to locate the root filesystem for the disk overview; not framework configuration | `internal/monitor/service.go:189` |

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
| `EventTypeNodeAutoPassed` | `approval.node.auto_passed` |  | `approval/events.go:70` |
| `EventTypeRunAbandoned` | `vef.cron.run.abandoned` |  | `cron/events.go:13` |
| `EventTypeRunFailed` | `vef.cron.run.failed` |  | `cron/events.go:10` |
| `EventTypeTaskApproved` | `approval.task.approved` |  | `approval/events.go:73` |
| `EventTypeTaskCanceled` | `approval.task.canceled` |  | `approval/events.go:76` |
| `EventTypeTaskCreated` | `approval.task.created` |  | `approval/events.go:72` |
| `EventTypeTaskDeadlineWarning` | `approval.task.deadline_warning` |  | `approval/events.go:82` |
| `EventTypeTaskHandled` | `approval.task.handled` |  | `approval/events.go:74` |
| `EventTypeTaskReassigned` | `approval.task.reassigned` |  | `approval/events.go:78` |
| `EventTypeTaskRejected` | `approval.task.rejected` |  | `approval/events.go:75` |
| `EventTypeTaskTimedOut` | `approval.task.timed_out` |  | `approval/events.go:79` |
| `EventTypeTaskTransferred` | `approval.task.transferred` |  | `approval/events.go:77` |
| `EventTypeTaskUrged` | `approval.task.urged` |  | `approval/events.go:83` |
| `eventTypeAudit` | `vef.api.request.audit` |  | `api/audit.go:9` |
| `eventTypeCodeSetChanged` | `vef.translate.code_set.changed` |  | `mold/cached_code_set_resolver.go:14` |
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
| `validator/rule.go:66` | `dynamic key sourced from ValidationRule.ErrMessageI18nKey; built-in rule keys are indexed separately` |  | `validator/rule.go:66` |
| `validator/validator.go:71` | `dynamic key sourced from label_i18n struct tags; tag values are indexed separately` |  | `validator/validator.go:71` |

## i18n message key

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `access_denied` | `access_denied` | i18n.T call | `result/errors.go:20` |
| `api_request_action` | `api_request_action` | label_i18n struct tag | `api/request.go:16` |
| `api_request_meta_invalid_json` | `api_request_meta_invalid_json` | i18n.T call | `api/api_errors.go:22` |
| `api_request_params_invalid_json` | `api_request_params_invalid_json` | i18n.T call | `api/api_errors.go:17` |
| `api_request_resource` | `api_request_resource` | label_i18n struct tag | `api/request.go:15` |
| `api_request_version` | `api_request_version` | label_i18n struct tag | `api/request.go:17` |
| `approval_access_denied` | `approval_access_denied` | i18n.T call | `internal/approval/shared/api_errors.go:122` |
| `approval_add_assignee_not_allowed` | `approval_add_assignee_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:97` |
| `approval_assignee_resolve_failed` | `approval_assignee_resolve_failed` | i18n.T call | `internal/approval/shared/api_errors.go:110` |
| `approval_binding_columns_conflict` | `approval_binding_columns_conflict` | i18n.T call | `internal/approval/shared/api_errors.go:61` |
| `approval_binding_incomplete` | `approval_binding_incomplete` | i18n.T call | `internal/approval/shared/api_errors.go:44` |
| `approval_binding_key_not_unique` | `approval_binding_key_not_unique` | i18n.T call | `internal/approval/shared/api_errors.go:70` |
| `approval_binding_projection_not_found` | `approval_binding_projection_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:88` |
| `approval_binding_schema_invalid` | `approval_binding_schema_invalid` | i18n.T call | `internal/approval/shared/api_errors.go:67` |
| `approval_binding_status_mapping_invalid` | `approval_binding_status_mapping_invalid` | i18n.T call | `internal/approval/shared/api_errors.go:74` |
| `approval_binding_target_busy` | `approval_binding_target_busy` | i18n.T call | `internal/approval/shared/api_errors.go:85` |
| `approval_binding_unexpected` | `approval_binding_unexpected` | i18n.T call | `internal/approval/shared/api_errors.go:64` |
| `approval_business_ref_required` | `approval_business_ref_required` | i18n.T call | `internal/approval/shared/api_errors.go:84` |
| `approval_flow_code_exists` | `approval_flow_code_exists` | i18n.T call | `internal/approval/shared/api_errors.go:21` |
| `approval_flow_not_active` | `approval_flow_not_active` | i18n.T call | `internal/approval/shared/api_errors.go:17` |
| `approval_flow_not_found` | `approval_flow_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:16` |
| `approval_form_cross_device_kind_mismatch` | `approval_form_cross_device_kind_mismatch` | i18n.T call | `internal/approval/formeditor/api_errors.go:53` |
| `approval_form_cross_device_table_mismatch` | `approval_form_cross_device_table_mismatch` | i18n.T call | `internal/approval/formeditor/api_errors.go:63` |
| `approval_form_data_too_large` | `approval_form_data_too_large` | i18n.T call | `internal/approval/shared/api_errors.go:118` |
| `approval_form_field_empty` | `approval_form_field_empty` | i18n.T call | `internal/approval/service/validation.go:456` |
| `approval_form_field_empty` | `approval_form_field_empty` | i18n.T call | `internal/approval/service/validation.go:463` |
| `approval_form_field_empty` | `approval_form_field_empty` | i18n.T call | `internal/approval/service/validation.go:470` |
| `approval_form_field_invalid_file_item` | `approval_form_field_invalid_file_item` | i18n.T call | `internal/approval/service/validation.go:476` |
| `approval_form_field_invalid_validation` | `approval_form_field_invalid_validation` | i18n.T call | `internal/approval/service/validation.go:412` |
| `approval_form_field_invalid_value` | `approval_form_field_invalid_value` | i18n.T call | `internal/approval/service/validation.go:499` |
| `approval_form_field_max_length` | `approval_form_field_max_length` | i18n.T call | `internal/approval/service/validation.go:403` |
| `approval_form_field_max_rows` | `approval_form_field_max_rows` | i18n.T call | `internal/approval/service/validation.go:321` |
| `approval_form_field_max_value` | `approval_form_field_max_value` | i18n.T call | `internal/approval/service/validation.go:443` |
| `approval_form_field_min_length` | `approval_form_field_min_length` | i18n.T call | `internal/approval/service/validation.go:396` |
| `approval_form_field_min_rows` | `approval_form_field_min_rows` | i18n.T call | `internal/approval/service/validation.go:315` |
| `approval_form_field_min_value` | `approval_form_field_min_value` | i18n.T call | `internal/approval/service/validation.go:436` |
| `approval_form_field_must_be_file` | `approval_form_field_must_be_file` | i18n.T call | `internal/approval/service/validation.go:483` |
| `approval_form_field_must_be_integer` | `approval_form_field_must_be_integer` | i18n.T call | `internal/approval/service/validation.go:428` |
| `approval_form_field_must_be_number` | `approval_form_field_must_be_number` | i18n.T call | `internal/approval/service/validation.go:289` |
| `approval_form_field_must_be_row_list` | `approval_form_field_must_be_row_list` | i18n.T call | `internal/approval/service/validation.go:310` |
| `approval_form_field_must_be_row_object` | `approval_form_field_must_be_row_object` | i18n.T call | `internal/approval/service/validation.go:330` |
| `approval_form_field_must_be_string` | `approval_form_field_must_be_string` | i18n.T call | `internal/approval/service/validation.go:275` |
| `approval_form_field_not_defined` | `approval_form_field_not_defined` | i18n.T call | `internal/approval/service/validation.go:619` |
| `approval_form_field_not_defined` | `approval_form_field_not_defined` | i18n.T call | `internal/approval/service/validation.go:73` |
| `approval_form_field_not_defined` | `approval_form_field_not_defined` | i18n.T call | `internal/approval/service/validation.go:357` |
| `approval_form_field_pattern_mismatch` | `approval_form_field_pattern_mismatch` | i18n.T call | `internal/approval/service/validation.go:416` |
| `approval_form_field_required` | `approval_form_field_required` | i18n.T call | `internal/approval/service/validation.go:81` |
| `approval_form_field_required` | `approval_form_field_required` | i18n.T call | `internal/approval/service/validation.go:368` |
| `approval_form_field_required` | `approval_form_field_required` | i18n.T call | `internal/approval/service/validation.go:116` |
| `approval_form_field_table_cell` | `approval_form_field_table_cell` | i18n.T call | `internal/approval/service/validation.go:376` |
| `approval_form_field_table_cell` | `approval_form_field_table_cell` | i18n.T call | `internal/approval/service/validation.go:366` |
| `approval_form_field_table_cell` | `approval_form_field_table_cell` | i18n.T call | `internal/approval/service/validation.go:355` |
| `approval_form_nested_subform` | `approval_form_nested_subform` | i18n.T call | `internal/approval/formeditor/api_errors.go:42` |
| `approval_form_schema_malformed` | `approval_form_schema_malformed` | i18n.T call | `internal/approval/formeditor/api_errors.go:68` |
| `approval_form_table_columns_empty` | `approval_form_table_columns_empty` | i18n.T call | `internal/approval/formeditor/api_errors.go:47` |
| `approval_form_unknown_field_type` | `approval_form_unknown_field_type` | i18n.T call | `internal/approval/formeditor/api_errors.go:33` |
| `approval_form_unmappable_field_type` | `approval_form_unmappable_field_type` | i18n.T call | `internal/approval/formeditor/api_errors.go:24` |
| `approval_form_validation_failed` | `approval_form_validation_failed` | i18n.T call | `internal/approval/shared/api_errors.go:112` |
| `approval_instance_completed` | `approval_instance_completed` | i18n.T call | `internal/approval/shared/api_errors.go:79` |
| `approval_instance_not_found` | `approval_instance_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:78` |
| `approval_invalid_add_assignee_type` | `approval_invalid_add_assignee_type` | i18n.T call | `internal/approval/shared/api_errors.go:102` |
| `approval_invalid_binding_mode` | `approval_invalid_binding_mode` | i18n.T call | `internal/approval/shared/api_errors.go:49` |
| `approval_invalid_business_identifier` | `approval_invalid_business_identifier` | i18n.T call | `internal/approval/shared/api_errors.go:26` |
| `approval_invalid_business_ref` | `approval_invalid_business_ref` | i18n.T call | `internal/approval/shared/api_errors.go:86` |
| `approval_invalid_flow_design` | `approval_invalid_flow_design` | i18n.T call | `internal/approval/shared/api_errors.go:20` |
| `approval_invalid_flow_label` | `approval_invalid_flow_label` | i18n.T call | `internal/approval/shared/api_errors.go:36` |
| `approval_invalid_form_design` | `approval_invalid_form_design` | i18n.T call | `internal/approval/shared/api_errors.go:41` |
| `approval_invalid_initiator_kind` | `approval_invalid_initiator_kind` | i18n.T call | `internal/approval/shared/api_errors.go:53` |
| `approval_invalid_instance_transition` | `approval_invalid_instance_transition` | i18n.T call | `internal/approval/shared/api_errors.go:83` |
| `approval_invalid_rollback_target` | `approval_invalid_rollback_target` | i18n.T call | `internal/approval/shared/api_errors.go:104` |
| `approval_invalid_storage_mode` | `approval_invalid_storage_mode` | i18n.T call | `internal/approval/shared/api_errors.go:58` |
| `approval_invalid_task_transition` | `approval_invalid_task_transition` | i18n.T call | `internal/approval/shared/api_errors.go:95` |
| `approval_invalid_title_template` | `approval_invalid_title_template` | i18n.T call | `internal/approval/shared/api_errors.go:32` |
| `approval_invalid_transfer_target` | `approval_invalid_transfer_target` | i18n.T call | `internal/approval/shared/api_errors.go:106` |
| `approval_last_assignee_removal` | `approval_last_assignee_removal` | i18n.T call | `internal/approval/shared/api_errors.go:105` |
| `approval_manual_cc_not_allowed` | `approval_manual_cc_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:100` |
| `approval_no_assignee` | `approval_no_assignee` | i18n.T call | `internal/approval/shared/api_errors.go:109` |
| `approval_no_published_version` | `approval_no_published_version` | i18n.T call | `internal/approval/shared/api_errors.go:18` |
| `approval_no_users_specified` | `approval_no_users_specified` | i18n.T call | `internal/approval/shared/api_errors.go:107` |
| `approval_not_allowed_initiate` | `approval_not_allowed_initiate` | i18n.T call | `internal/approval/shared/api_errors.go:80` |
| `approval_not_applicant` | `approval_not_applicant` | i18n.T call | `internal/approval/shared/api_errors.go:103` |
| `approval_not_assignee` | `approval_not_assignee` | i18n.T call | `internal/approval/shared/api_errors.go:94` |
| `approval_opinion_required` | `approval_opinion_required` | i18n.T call | `internal/approval/shared/api_errors.go:99` |
| `approval_remove_assignee_not_allowed` | `approval_remove_assignee_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:101` |
| `approval_resubmit_not_allowed` | `approval_resubmit_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:82` |
| `approval_rollback_not_allowed` | `approval_rollback_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:96` |
| `approval_task_not_found` | `approval_task_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:92` |
| `approval_task_not_pending` | `approval_task_not_pending` | i18n.T call | `internal/approval/shared/api_errors.go:93` |
| `approval_terminate_not_allowed` | `approval_terminate_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:126` |
| `approval_transfer_not_allowed` | `approval_transfer_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:98` |
| `approval_urge_too_frequent` | `approval_urge_too_frequent` | i18n.T call | `internal/approval/command/urge_task.go:116` |
| `approval_version_not_draft` | `approval_version_not_draft` | i18n.T call | `internal/approval/shared/api_errors.go:19` |
| `approval_version_not_found` | `approval_version_not_found` | i18n.T call | `internal/approval/shared/api_errors.go:22` |
| `approval_withdraw_not_allowed` | `approval_withdraw_not_allowed` | i18n.T call | `internal/approval/shared/api_errors.go:81` |
| `auth_challenge_response` | `auth_challenge_response` | label_i18n struct tag | `internal/security/auth_resource.go:251` |
| `auth_challenge_token` | `auth_challenge_token` | label_i18n struct tag | `internal/security/auth_resource.go:249` |
| `auth_challenge_type` | `auth_challenge_type` | label_i18n struct tag | `internal/security/auth_resource.go:250` |
| `auth_credentials` | `auth_credentials` | label_i18n struct tag | `internal/security/auth_resource.go:117` |
| `auth_principal` | `auth_principal` | label_i18n struct tag | `internal/security/auth_resource.go:116` |
| `auth_refresh_token` | `auth_refresh_token` | label_i18n struct tag | `internal/security/auth_resource.go:190` |
| `auth_type` | `auth_type` | label_i18n struct tag | `internal/security/auth_resource.go:115` |
| `cron_job_not_registered` | `cron_job_not_registered` | i18n.T call | `cron/api_errors.go:35` |
| `cron_schedule_disabled` | `cron_schedule_disabled` | i18n.T call | `cron/api_errors.go:31` |
| `cron_schedule_exists` | `cron_schedule_exists` | i18n.T call | `cron/api_errors.go:27` |
| `cron_schedule_invalid` | `cron_schedule_invalid` | i18n.T call | `cron/api_errors.go:58` |
| `cron_schedule_not_found` | `cron_schedule_not_found` | i18n.T call | `cron/api_errors.go:23` |
| `cron_store_disabled` | `cron_store_disabled` | i18n.T call | `cron/api_errors.go:39` |
| `cron_trigger_invalid` | `cron_trigger_invalid` | i18n.T call | `cron/api_errors.go:48` |
| `crud_batch_create_list` | `crud_batch_create_list` | label_i18n struct tag | `crud/params.go:12` |
| `crud_batch_delete_pks` | `crud_batch_delete_pks` | label_i18n struct tag | `crud/params.go:28` |
| `crud_batch_update_list` | `crud_batch_update_list` | label_i18n struct tag | `crud/params.go:19` |
| `crud_composite_primary_key_requires_map` | `crud_composite_primary_key_requires_map` | i18n.T call | `crud/api_errors.go:27` |
| `crud_created` | `crud_created` | i18n.T call | `crud/create_many.go:103` |
| `crud_created` | `crud_created` | i18n.T call | `crud/create.go:88` |
| `crud_deleted` | `crud_deleted` | i18n.T call | `crud/delete.go:111` |
| `crud_deleted` | `crud_deleted` | i18n.T call | `crud/delete_many.go:128` |
| `crud_field_not_exist_in_model` | `crud_field_not_exist_in_model` | i18n.T call | `crud/api_errors.go:70` |
| `crud_file_open_failed` | `crud_file_open_failed` | i18n.T call | `crud/api_errors.go:47` |
| `crud_import_requires_file` | `crud_import_requires_file` | i18n.T call | `crud/api_errors.go:39` |
| `crud_import_requires_multipart` | `crud_import_requires_multipart` | i18n.T call | `crud/api_errors.go:35` |
| `crud_import_type_assertion_failed` | `crud_import_type_assertion_failed` | i18n.T call | `crud/api_errors.go:51` |
| `crud_import_validation_failed` | `crud_import_validation_failed` | i18n.T call | `crud/import.go:147` |
| `crud_imported` | `crud_imported` | i18n.T call | `crud/import.go:178` |
| `crud_primary_key_required` | `crud_primary_key_required` | i18n.T call | `crud/api_errors.go:60` |
| `crud_processor_must_return_slice` | `crud_processor_must_return_slice` | i18n.T call | `crud/find_page.go:88` |
| `crud_unsupported_export_format` | `crud_unsupported_export_format` | i18n.T call | `crud/api_errors.go:31` |
| `crud_unsupported_import_format` | `crud_unsupported_import_format` | i18n.T call | `crud/api_errors.go:43` |
| `crud_updated` | `crud_updated` | i18n.T call | `crud/update_many.go:138` |
| `crud_updated` | `crud_updated` | i18n.T call | `crud/update.go:126` |
| `dangerous_sql` | `dangerous_sql` | i18n.T call | `result/errors.go:58` |
| `error` | `error` | i18n.T call | `result/error.go:69` |
| `expression_evaluation_failed` | `expression_evaluation_failed` | i18n.T call | `expression/api_errors.go:15` |
| `foreign_key_violation` | `foreign_key_violation` | i18n.T call | `result/errors.go:54` |
| `integration_adapter_disabled` | `integration_adapter_disabled` | i18n.T call | `integration/api_errors.go:70` |
| `integration_adapter_not_found` | `integration_adapter_not_found` | i18n.T call | `integration/api_errors.go:66` |
| `integration_code_map_missing` | `integration_code_map_missing` | i18n.T call | `integration/api_errors.go:222` |
| `integration_contract_disabled` | `integration_contract_disabled` | i18n.T call | `integration/api_errors.go:54` |
| `integration_contract_not_found` | `integration_contract_not_found` | i18n.T call | `integration/api_errors.go:50` |
| `integration_inbound_auth_failed` | `integration_inbound_auth_failed` | i18n.T call | `integration/api_errors.go:250` |
| `integration_inbound_handler_missing` | `integration_inbound_handler_missing` | i18n.T call | `integration/api_errors.go:258` |
| `integration_input_invalid` | `integration_input_invalid` | i18n.T call | `integration/api_errors.go:117` |
| `integration_invalid_auth_params` | `integration_invalid_auth_params` | i18n.T call | `integration/api_errors.go:190` |
| `integration_invalid_base_url` | `integration_invalid_base_url` | i18n.T call | `integration/api_errors.go:109` |
| `integration_invalid_code_map` | `integration_invalid_code_map` | i18n.T call | `integration/api_errors.go:241` |
| `integration_invalid_data_source` | `integration_invalid_data_source` | i18n.T call | `integration/api_errors.go:199` |
| `integration_invalid_direction` | `integration_invalid_direction` | i18n.T call | `integration/api_errors.go:206` |
| `integration_invalid_envelope` | `integration_invalid_envelope` | i18n.T call | `integration/api_errors.go:180` |
| `integration_invalid_label` | `integration_invalid_label` | i18n.T call | `integration/api_errors.go:214` |
| `integration_invalid_route_ref` | `integration_invalid_route_ref` | i18n.T call | `integration/api_errors.go:103` |
| `integration_invalid_schema` | `integration_invalid_schema` | i18n.T call | `integration/api_errors.go:161` |
| `integration_invalid_script` | `integration_invalid_script` | i18n.T call | `integration/api_errors.go:170` |
| `integration_invocation_canceled` | `integration_invocation_canceled` | i18n.T call | `integration/api_errors.go:97` |
| `integration_invocation_timeout` | `integration_invocation_timeout` | i18n.T call | `integration/api_errors.go:91` |
| `integration_output_invalid` | `integration_output_invalid` | i18n.T call | `integration/api_errors.go:126` |
| `integration_route_not_found` | `integration_route_not_found` | i18n.T call | `integration/api_errors.go:74` |
| `integration_script_failed` | `integration_script_failed` | i18n.T call | `integration/api_errors.go:143` |
| `integration_system_disabled` | `integration_system_disabled` | i18n.T call | `integration/api_errors.go:62` |
| `integration_system_not_found` | `integration_system_not_found` | i18n.T call | `integration/api_errors.go:58` |
| `integration_target_ambiguous` | `integration_target_ambiguous` | i18n.T call | `integration/api_errors.go:80` |
| `integration_transport_failed` | `integration_transport_failed` | i18n.T call | `integration/api_errors.go:87` |
| `integration_unknown_auth_scheme` | `integration_unknown_auth_scheme` | i18n.T call | `integration/api_errors.go:152` |
| `integration_upstream_failed` | `integration_upstream_failed` | i18n.T call | `integration/api_errors.go:135` |
| `integration_value_unmapped` | `integration_value_unmapped` | i18n.T call | `integration/api_errors.go:231` |
| `monitor_collection_failed` | `monitor_collection_failed` | i18n.T call | `monitor/api_errors.go:25` |
| `monitor_not_ready` | `monitor_not_ready` | i18n.T call | `monitor/api_errors.go:19` |
| `ok` | `ok` | i18n.T call | `result/result.go:64` |
| `record_already_exists` | `record_already_exists` | i18n.T call | `result/errors.go:50` |
| `record_not_found` | `record_not_found` | i18n.T call | `result/errors.go:46` |
| `request_timeout` | `request_timeout` | i18n.T call | `result/errors.go:30` |
| `schema_table_not_found` | `schema_table_not_found` | i18n.T call | `schema/api_errors.go:16` |
| `security_account_locked` | `security_account_locked` | i18n.T call | `security/api_errors.go:321` |
| `security_api_key_invalid` | `security_api_key_invalid` | i18n.T call | `security/api_errors.go:167` |
| `security_app_id_required` | `security_app_id_required` | i18n.T call | `security/api_errors.go:122` |
| `security_auth_header_invalid` | `security_auth_header_invalid` | i18n.T call | `security/api_errors.go:197` |
| `security_auth_header_missing` | `security_auth_header_missing` | i18n.T call | `security/api_errors.go:192` |
| `security_basic_credentials_invalid` | `security_basic_credentials_invalid` | i18n.T call | `security/api_errors.go:172` |
| `security_challenge_resolve_failed` | `security_challenge_resolve_failed` | i18n.T call | `security/api_errors.go:221` |
| `security_challenge_token_invalid` | `security_challenge_token_invalid` | i18n.T call | `security/api_errors.go:211` |
| `security_challenge_type_invalid` | `security_challenge_type_invalid` | i18n.T call | `security/api_errors.go:216` |
| `security_credentials_format_invalid` | `security_credentials_format_invalid` | i18n.T call | `internal/security/signature_authenticator.go:81` |
| `security_department_required` | `security_department_required` | i18n.T call | `security/api_errors.go:241` |
| `security_external_app_disabled` | `security_external_app_disabled` | i18n.T call | `security/api_errors.go:157` |
| `security_external_app_loader_not_implemented` | `security_external_app_loader_not_implemented` | i18n.T call | `internal/security/signature_authenticator.go:71` |
| `security_external_app_not_found` | `security_external_app_not_found` | i18n.T call | `security/api_errors.go:152` |
| `security_invalid_credentials` | `security_invalid_credentials` | i18n.T call | `internal/security/password_authenticator.go:107` |
| `security_invalid_credentials` | `security_invalid_credentials` | i18n.T call | `internal/security/password_authenticator.go:111` |
| `security_invalid_credentials` | `security_invalid_credentials` | i18n.T call | `internal/security/password_authenticator.go:100` |
| `security_invalid_credentials` | `security_invalid_credentials` | i18n.T call | `internal/security/password_authenticator.go:82` |
| `security_ip_not_allowed` | `security_ip_not_allowed` | i18n.T call | `security/api_errors.go:162` |
| `security_new_password_required` | `security_new_password_required` | i18n.T call | `security/api_errors.go:236` |
| `security_nonce_already_used` | `security_nonce_already_used` | i18n.T call | `security/api_errors.go:187` |
| `security_nonce_invalid` | `security_nonce_invalid` | i18n.T call | `security/api_errors.go:182` |
| `security_nonce_required` | `security_nonce_required` | i18n.T call | `security/api_errors.go:177` |
| `security_otp_code_invalid` | `security_otp_code_invalid` | i18n.T call | `security/api_errors.go:231` |
| `security_otp_code_required` | `security_otp_code_required` | i18n.T call | `security/api_errors.go:226` |
| `security_password_blocked` | `security_password_blocked` | i18n.T call | `security/api_errors.go:276` |
| `security_password_contains_identity` | `security_password_contains_identity` | i18n.T call | `security/api_errors.go:271` |
| `security_password_missing_digit` | `security_password_missing_digit` | i18n.T call | `security/api_errors.go:261` |
| `security_password_missing_lowercase` | `security_password_missing_lowercase` | i18n.T call | `security/api_errors.go:256` |
| `security_password_missing_symbol` | `security_password_missing_symbol` | i18n.T call | `security/api_errors.go:266` |
| `security_password_missing_uppercase` | `security_password_missing_uppercase` | i18n.T call | `security/api_errors.go:251` |
| `security_password_required` | `security_password_required` | i18n.T call | `internal/security/password_authenticator.go:69` |
| `security_password_reused` | `security_password_reused` | i18n.T call | `security/api_errors.go:281` |
| `security_password_too_few_char_classes` | `security_password_too_few_char_classes` | i18n.T call | `security/api_errors.go:308` |
| `security_password_too_long` | `security_password_too_long` | i18n.T call | `security/api_errors.go:299` |
| `security_password_too_short` | `security_password_too_short` | i18n.T call | `security/api_errors.go:290` |
| `security_reserved_principal_forbidden` | `security_reserved_principal_forbidden` | i18n.T call | `security/api_errors.go:113` |
| `security_signature_expired` | `security_signature_expired` | i18n.T call | `security/api_errors.go:142` |
| `security_signature_invalid` | `security_signature_invalid` | i18n.T call | `security/api_errors.go:147` |
| `security_signature_required` | `security_signature_required` | i18n.T call | `security/api_errors.go:132` |
| `security_system_principal_login_forbidden` | `security_system_principal_login_forbidden` | i18n.T call | `internal/security/password_authenticator.go:64` |
| `security_timestamp_invalid` | `security_timestamp_invalid` | i18n.T call | `security/api_errors.go:137` |
| `security_timestamp_required` | `security_timestamp_required` | i18n.T call | `security/api_errors.go:127` |
| `security_token_expired` | `security_token_expired` | i18n.T call | `security/api_errors.go:84` |
| `security_token_invalid` | `security_token_invalid` | i18n.T call | `security/api_errors.go:89` |
| `security_token_invalid_audience` | `security_token_invalid_audience` | i18n.T call | `security/api_errors.go:104` |
| `security_token_invalid_issuer` | `security_token_invalid_issuer` | i18n.T call | `security/api_errors.go:99` |
| `security_token_not_valid_yet` | `security_token_not_valid_yet` | i18n.T call | `security/api_errors.go:94` |
| `security_too_many_concurrent_sessions` | `security_too_many_concurrent_sessions` | i18n.T call | `security/api_errors.go:202` |
| `security_unauthenticated` | `security_unauthenticated` | i18n.T call | `security/api_errors.go:79` |
| `security_unsupported_authentication_type` | `security_unsupported_authentication_type` | i18n.T call | `internal/security/auth_manager.go:69` |
| `security_user_info_loader_not_implemented` | `security_user_info_loader_not_implemented` | i18n.T call | `internal/security/auth_resource.go:356` |
| `security_user_loader_not_implemented` | `security_user_loader_not_implemented` | i18n.T call | `internal/security/jwt_refresh_authenticator.go:34` |
| `security_user_loader_not_implemented` | `security_user_loader_not_implemented` | i18n.T call | `internal/security/password_authenticator.go:55` |
| `security_username_required` | `security_username_required` | i18n.T call | `internal/security/password_authenticator.go:60` |
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
| `too_many_requests` | `too_many_requests` | i18n.T call | `result/errors.go:25` |
| `unknown_error` | `unknown_error` | i18n.T call | `result/errors.go:35` |
| `unknown_error` | `unknown_error` | i18n.T call | `internal/app/error.go:64` |
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
| `diveTag` | `dive` |  | `internal/mold/restricted.go:6` |
| `endKeysTag` | `endkeys` |  | `internal/mold/restricted.go:13` |
| `ignoreTag` | `-` |  | `internal/mold/restricted.go:9` |
| `keysTag` | `keys` |  | `internal/mold/restricted.go:12` |
| `restrictedTagChars` | `.[],\|=+()`~!@#$%^&amp;*\"/?&lt;&gt;{}` |  | `internal/mold/restricted.go:7` |
| `tag name` | `mold` |  | `internal/mold/mold.go:41` |
| `tagKeySeparator` | `=` |  | `internal/mold/restricted.go:10` |
| `tagSeparator` | `,` |  | `internal/mold/restricted.go:8` |
| `utf8HexComma` | `0x2C` |  | `internal/mold/restricted.go:11` |

## mold transformer tag

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `expr` | `expr` |  | `internal/expression/transformer.go:33` |
| `translate` | `translate` |  | `internal/mold/translate.go:190` |

## result error code

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `ErrCodeAPIKeyInvalid` | `1025` |  | `security/api_errors.go:59` |
| `ErrCodeAbortFailed` | `2219` |  | `storage/api_errors.go:30` |
| `ErrCodeAccessDenied` | `1100` |  | `result/constants.go:32` |
| `ErrCodeAccountLocked` | `1023` |  | `security/api_errors.go:57` |
| `ErrCodeAdapterDisabled` | `2605` |  | `integration/api_errors.go:19` |
| `ErrCodeAdapterNotFound` | `2604` |  | `integration/api_errors.go:18` |
| `ErrCodeAppIDRequired` | `1009` |  | `security/api_errors.go:43` |
| `ErrCodeAuthHeaderInvalid` | `1022` |  | `security/api_errors.go:56` |
| `ErrCodeAuthHeaderMissing` | `1021` |  | `security/api_errors.go:55` |
| `ErrCodeBadRequest` | `1400` |  | `result/constants.go:41` |
| `ErrCodeBasicCredentialsInvalid` | `1026` |  | `security/api_errors.go:60` |
| `ErrCodeChallengeResolveFailed` | `1034` |  | `security/api_errors.go:65` |
| `ErrCodeChallengeTokenInvalid` | `1031` |  | `security/api_errors.go:63` |
| `ErrCodeChallengeTypeInvalid` | `1033` |  | `security/api_errors.go:64` |
| `ErrCodeClaimExpired` | `2204` |  | `storage/api_errors.go:15` |
| `ErrCodeClaimNotMultipart` | `2212` |  | `storage/api_errors.go:23` |
| `ErrCodeClaimNotPending` | `2203` |  | `storage/api_errors.go:14` |
| `ErrCodeCollectionFailed` | `2101` |  | `monitor/api_errors.go:11` |
| `ErrCodeCompositePrimaryKeyRequiresMap` | `2403` |  | `crud/api_errors.go:13` |
| `ErrCodeContractDisabled` | `2601` |  | `integration/api_errors.go:15` |
| `ErrCodeContractNotFound` | `2600` |  | `integration/api_errors.go:14` |
| `ErrCodeCredentialsInvalid` | `1008` |  | `security/api_errors.go:42` |
| `ErrCodeDangerousSQL` | `1600` |  | `result/constants.go:49` |
| `ErrCodeDefault` | `2000` |  | `result/constants.go:55` |
| `ErrCodeDepartmentRequired` | `1038` |  | `security/api_errors.go:69` |
| `ErrCodeEvaluationFailed` | `2500` |  | `expression/api_errors.go:9` |
| `ErrCodeExternalAppDisabled` | `1015` |  | `security/api_errors.go:49` |
| `ErrCodeExternalAppNotFound` | `1014` |  | `security/api_errors.go:48` |
| `ErrCodeFailedToGetFile` | `2202` |  | `storage/api_errors.go:12` |
| `ErrCodeFieldNotExistInModel` | `2401` |  | `crud/api_errors.go:11` |
| `ErrCodeFileNotFound` | `2201` |  | `storage/api_errors.go:11` |
| `ErrCodeFileOpenFailed` | `2408` |  | `crud/api_errors.go:18` |
| `ErrCodeForeignKeyViolation` | `2003` |  | `result/constants.go:58` |
| `ErrCodeIPNotAllowed` | `1016` |  | `security/api_errors.go:50` |
| `ErrCodeImportRequiresFile` | `2406` |  | `crud/api_errors.go:16` |
| `ErrCodeImportRequiresMultipart` | `2405` |  | `crud/api_errors.go:15` |
| `ErrCodeImportTypeAssertionFailed` | `2409` |  | `crud/api_errors.go:19` |
| `ErrCodeImportValidationFailed` | `2410` |  | `crud/api_errors.go:20` |
| `ErrCodeInboundAuthFailed` | `2622` |  | `integration/api_errors.go:36` |
| `ErrCodeInboundHandlerMissing` | `2623` |  | `integration/api_errors.go:37` |
| `ErrCodeInputInvalid` | `2608` |  | `integration/api_errors.go:22` |
| `ErrCodeInvalidAuthParams` | `2617` |  | `integration/api_errors.go:31` |
| `ErrCodeInvalidBaseURL` | `2619` |  | `integration/api_errors.go:33` |
| `ErrCodeInvalidCodeMap` | `2629` |  | `integration/api_errors.go:43` |
| `ErrCodeInvalidDataSource` | `2620` |  | `integration/api_errors.go:34` |
| `ErrCodeInvalidDirection` | `2621` |  | `integration/api_errors.go:35` |
| `ErrCodeInvalidEnvelope` | `2625` |  | `integration/api_errors.go:39` |
| `ErrCodeInvalidFileKey` | `2200` |  | `storage/api_errors.go:10` |
| `ErrCodeInvalidLabel` | `2626` |  | `integration/api_errors.go:40` |
| `ErrCodeInvalidRouteRef` | `2618` |  | `integration/api_errors.go:32` |
| `ErrCodeInvalidSchema` | `2615` |  | `integration/api_errors.go:29` |
| `ErrCodeInvalidScript` | `2616` |  | `integration/api_errors.go:30` |
| `ErrCodeInvocationCanceled` | `2624` |  | `integration/api_errors.go:38` |
| `ErrCodeInvocationTimeout` | `2612` |  | `integration/api_errors.go:26` |
| `ErrCodeJobNotRegistered` | `2704` |  | `cron/api_errors.go:14` |
| `ErrCodeMissingCodeMap` | `2627` |  | `integration/api_errors.go:41` |
| `ErrCodeMultipartNotSupported` | `2206` |  | `storage/api_errors.go:17` |
| `ErrCodeNewPasswordRequired` | `1037` |  | `security/api_errors.go:68` |
| `ErrCodeNonceAlreadyUsed` | `1020` |  | `security/api_errors.go:54` |
| `ErrCodeNonceInvalid` | `1019` |  | `security/api_errors.go:53` |
| `ErrCodeNonceRequired` | `1018` |  | `security/api_errors.go:52` |
| `ErrCodeNotFound` | `1200` |  | `result/constants.go:35` |
| `ErrCodeNotImplemented` | `1500` |  | `result/constants.go:46` |
| `ErrCodeNotReady` | `2100` |  | `monitor/api_errors.go:10` |
| `ErrCodeOTPCodeInvalid` | `1036` |  | `security/api_errors.go:67` |
| `ErrCodeOTPCodeRequired` | `1035` |  | `security/api_errors.go:66` |
| `ErrCodeOutputInvalid` | `2609` |  | `integration/api_errors.go:23` |
| `ErrCodePasswordPolicyViolation` | `1050` |  | `security/api_errors.go:73` |
| `ErrCodePrimaryKeyRequired` | `2402` |  | `crud/api_errors.go:12` |
| `ErrCodePrincipalInvalid` | `1007` |  | `security/api_errors.go:41` |
| `ErrCodePublicUploadsNotAllowed` | `2207` |  | `storage/api_errors.go:18` |
| `ErrCodeRecordAlreadyExists` | `2002` |  | `result/constants.go:57` |
| `ErrCodeRecordNotFound` | `2001` |  | `result/constants.go:56` |
| `ErrCodeRequestTimeout` | `1402` |  | `result/constants.go:43` |
| `ErrCodeRouteNotFound` | `2606` |  | `integration/api_errors.go:20` |
| `ErrCodeScheduleDisabled` | `2702` |  | `cron/api_errors.go:12` |
| `ErrCodeScheduleExists` | `2701` |  | `cron/api_errors.go:11` |
| `ErrCodeScheduleInvalid` | `2706` |  | `cron/api_errors.go:16` |
| `ErrCodeScheduleNotFound` | `2700` |  | `cron/api_errors.go:10` |
| `ErrCodeScriptFailed` | `2613` |  | `integration/api_errors.go:27` |
| `ErrCodeSignatureExpired` | `1013` |  | `security/api_errors.go:47` |
| `ErrCodeSignatureInvalid` | `1017` |  | `security/api_errors.go:51` |
| `ErrCodeSignatureRequired` | `1011` |  | `security/api_errors.go:45` |
| `ErrCodeStoreDisabled` | `2705` |  | `cron/api_errors.go:15` |
| `ErrCodeSystemDisabled` | `2603` |  | `integration/api_errors.go:17` |
| `ErrCodeSystemNotFound` | `2602` |  | `integration/api_errors.go:16` |
| `ErrCodeTableNotFound` | `2300` |  | `schema/api_errors.go:10` |
| `ErrCodeTargetAmbiguous` | `2607` |  | `integration/api_errors.go:21` |
| `ErrCodeTimestampInvalid` | `1012` |  | `security/api_errors.go:46` |
| `ErrCodeTimestampRequired` | `1010` |  | `security/api_errors.go:44` |
| `ErrCodeTokenExpired` | `1002` |  | `security/api_errors.go:36` |
| `ErrCodeTokenInvalid` | `1003` |  | `security/api_errors.go:37` |
| `ErrCodeTokenInvalidAudience` | `1006` |  | `security/api_errors.go:40` |
| `ErrCodeTokenInvalidIssuer` | `1005` |  | `security/api_errors.go:39` |
| `ErrCodeTokenNotValidYet` | `1004` |  | `security/api_errors.go:38` |
| `ErrCodeTooManyConcurrentSessions` | `1024` |  | `security/api_errors.go:58` |
| `ErrCodeTooManyPendingUploads` | `2209` |  | `storage/api_errors.go:20` |
| `ErrCodeTooManyRequests` | `1401` |  | `result/constants.go:42` |
| `ErrCodeTransportFailed` | `2611` |  | `integration/api_errors.go:25` |
| `ErrCodeTriggerInvalid` | `2703` |  | `cron/api_errors.go:13` |
| `ErrCodeUnauthenticated` | `1000` |  | `security/api_errors.go:34` |
| `ErrCodeUnknown` | `1900` |  | `result/constants.go:52` |
| `ErrCodeUnknownAuthScheme` | `2614` |  | `integration/api_errors.go:28` |
| `ErrCodeUnmappedValue` | `2628` |  | `integration/api_errors.go:42` |
| `ErrCodeUnsupportedAuthenticationType` | `1001` |  | `security/api_errors.go:35` |
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
| `ErrCodeUpstreamFailed` | `2610` |  | `integration/api_errors.go:24` |

## result message key

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `ErrMessage` | `error` |  | `result/constants.go:8` |
| `ErrMessageAccessDenied` | `access_denied` |  | `result/constants.go:14` |
| `ErrMessageAccountLocked` | `security_account_locked` |  | `security/api_errors.go:25` |
| `ErrMessageChallengeResolveFailed` | `security_challenge_resolve_failed` |  | `security/api_errors.go:24` |
| `ErrMessageCredentialsFormatInvalid` | `security_credentials_format_invalid` |  | `security/api_errors.go:20` |
| `ErrMessageDangerousSQL` | `dangerous_sql` |  | `result/constants.go:22` |
| `ErrMessageExternalAppLoaderNotImplemented` | `security_external_app_loader_not_implemented` |  | `security/api_errors.go:19` |
| `ErrMessageForeignKeyViolation` | `foreign_key_violation` |  | `result/constants.go:21` |
| `ErrMessageFormFieldEmpty` | `approval_form_field_empty` |  | `internal/approval/shared/messages.go:27` |
| `ErrMessageFormFieldInvalidFileItem` | `approval_form_field_invalid_file_item` |  | `internal/approval/shared/messages.go:28` |
| `ErrMessageFormFieldInvalidValidation` | `approval_form_field_invalid_validation` |  | `internal/approval/shared/messages.go:23` |
| `ErrMessageFormFieldInvalidValue` | `approval_form_field_invalid_value` |  | `internal/approval/shared/messages.go:30` |
| `ErrMessageFormFieldMaxLength` | `approval_form_field_max_length` |  | `internal/approval/shared/messages.go:22` |
| `ErrMessageFormFieldMaxRows` | `approval_form_field_max_rows` |  | `internal/approval/shared/messages.go:34` |
| `ErrMessageFormFieldMaxValue` | `approval_form_field_max_value` |  | `internal/approval/shared/messages.go:26` |
| `ErrMessageFormFieldMinLength` | `approval_form_field_min_length` |  | `internal/approval/shared/messages.go:21` |
| `ErrMessageFormFieldMinRows` | `approval_form_field_min_rows` |  | `internal/approval/shared/messages.go:33` |
| `ErrMessageFormFieldMinValue` | `approval_form_field_min_value` |  | `internal/approval/shared/messages.go:25` |
| `ErrMessageFormFieldMustBeFile` | `approval_form_field_must_be_file` |  | `internal/approval/shared/messages.go:29` |
| `ErrMessageFormFieldMustBeInteger` | `approval_form_field_must_be_integer` |  | `internal/approval/shared/messages.go:20` |
| `ErrMessageFormFieldMustBeNumber` | `approval_form_field_must_be_number` |  | `internal/approval/shared/messages.go:19` |
| `ErrMessageFormFieldMustBeRowList` | `approval_form_field_must_be_row_list` |  | `internal/approval/shared/messages.go:31` |
| `ErrMessageFormFieldMustBeRowObject` | `approval_form_field_must_be_row_object` |  | `internal/approval/shared/messages.go:32` |
| `ErrMessageFormFieldMustBeString` | `approval_form_field_must_be_string` |  | `internal/approval/shared/messages.go:18` |
| `ErrMessageFormFieldNotDefined` | `approval_form_field_not_defined` |  | `internal/approval/shared/messages.go:16` |
| `ErrMessageFormFieldPatternMismatch` | `approval_form_field_pattern_mismatch` |  | `internal/approval/shared/messages.go:24` |
| `ErrMessageFormFieldRequired` | `approval_form_field_required` |  | `internal/approval/shared/messages.go:17` |
| `ErrMessageFormFieldTableCell` | `approval_form_field_table_cell` |  | `internal/approval/shared/messages.go:35` |
| `ErrMessageNotFound` | `not_found` |  | `result/constants.go:12` |
| `ErrMessagePasswordTooFewCharClasses` | `security_password_too_few_char_classes` |  | `security/api_errors.go:28` |
| `ErrMessagePasswordTooLong` | `security_password_too_long` |  | `security/api_errors.go:27` |
| `ErrMessagePasswordTooShort` | `security_password_too_short` |  | `security/api_errors.go:26` |
| `ErrMessageRecordAlreadyExists` | `record_already_exists` |  | `result/constants.go:20` |
| `ErrMessageRecordNotFound` | `record_not_found` |  | `result/constants.go:19` |
| `ErrMessageRequestTimeout` | `request_timeout` |  | `result/constants.go:16` |
| `ErrMessageTooManyRequests` | `too_many_requests` |  | `result/constants.go:13` |
| `ErrMessageUnauthenticated` | `security_unauthenticated` |  | `security/api_errors.go:18` |
| `ErrMessageUnknown` | `unknown_error` |  | `result/constants.go:11` |
| `ErrMessageUnsupportedAuthenticationType` | `security_unsupported_authentication_type` |  | `security/api_errors.go:21` |
| `ErrMessageUnsupportedMediaType` | `unsupported_media_type` |  | `result/constants.go:15` |
| `ErrMessageUrgeTooFrequent` | `approval_urge_too_frequent` |  | `internal/approval/shared/messages.go:10` |
| `ErrMessageUserInfoLoaderNotImplemented` | `security_user_info_loader_not_implemented` |  | `security/api_errors.go:23` |
| `ErrMessageUserLoaderNotImplemented` | `security_user_loader_not_implemented` |  | `security/api_errors.go:22` |

## runtime enum value

| Name | Value | Details | Source |
| --- | --- | --- | --- |
| `AcquireResultAcquired (AcquireResult)` | `acquired` |  | `event/inbox/inbox.go:29` |
| `AcquireResultCompleted (AcquireResult)` | `completed` |  | `event/inbox/inbox.go:32` |
| `AcquireResultInProgress (AcquireResult)` | `in_progress` |  | `event/inbox/inbox.go:36` |
| `ActionAddAssignee (ActionType)` | `add_assignee` |  | `approval/enums.go:353` |
| `ActionAddCC (ActionType)` | `add_cc` |  | `approval/enums.go:359` |
| `ActionApprove (ActionType)` | `approve` |  | `approval/enums.go:346` |
| `ActionCancel (ActionType)` | `cancel` |  | `approval/enums.go:351` |
| `ActionExecute (ActionType)` | `execute` |  | `approval/enums.go:355` |
| `ActionHandle (ActionType)` | `handle` |  | `approval/enums.go:347` |
| `ActionReassign (ActionType)` | `reassign` |  | `approval/enums.go:357` |
| `ActionReject (ActionType)` | `reject` |  | `approval/enums.go:348` |
| `ActionRemoveAssignee (ActionType)` | `remove_assignee` |  | `approval/enums.go:354` |
| `ActionResubmit (ActionType)` | `resubmit` |  | `approval/enums.go:356` |
| `ActionRollback (ActionType)` | `rollback` |  | `approval/enums.go:352` |
| `ActionSubmit (ActionType)` | `submit` |  | `approval/enums.go:345` |
| `ActionTerminate (ActionType)` | `terminate` |  | `approval/enums.go:358` |
| `ActionTransfer (ActionType)` | `transfer` |  | `approval/enums.go:349` |
| `ActionWithdraw (ActionType)` | `withdraw` |  | `approval/enums.go:350` |
| `AddAssigneeAfter (AddAssigneeType)` | `after` |  | `approval/enums.go:210` |
| `AddAssigneeBefore (AddAssigneeType)` | `before` |  | `approval/enums.go:209` |
| `AddAssigneeParallel (AddAssigneeType)` | `parallel` |  | `approval/enums.go:211` |
| `AesModeCbc (AESMode)` | `CBC` |  | `cryptox/aes_cipher.go:16` |
| `AesModeGcm (AESMode)` | `GCM` |  | `cryptox/aes_cipher.go:17` |
| `AggregateAvg (AggregateKind)` | `avg` |  | `approval/condition.go:57` |
| `AggregateCount (AggregateKind)` | `count` |  | `approval/condition.go:56` |
| `AggregateSum (AggregateKind)` | `sum` |  | `approval/condition.go:55` |
| `ApprovalBindingEventual (ApprovalBindingConsistency)` | `eventual` |  | `config/approval.go:19` |
| `ApprovalBindingSynchronous (ApprovalBindingConsistency)` | `synchronous` |  | `config/approval.go:16` |
| `ApprovalParallel (ApprovalMethod)` | `parallel` |  | `approval/enums.go:107` |
| `ApprovalSequential (ApprovalMethod)` | `sequential` |  | `approval/enums.go:106` |
| `AssigneeDepartment (AssigneeKind)` | `department` |  | `approval/enums.go:256` |
| `AssigneeDepartmentLeader (AssigneeKind)` | `department_leader` |  | `approval/enums.go:259` |
| `AssigneeFormField (AssigneeKind)` | `form_field` |  | `approval/enums.go:260` |
| `AssigneeRole (AssigneeKind)` | `role` |  | `approval/enums.go:255` |
| `AssigneeSelf (AssigneeKind)` | `self` |  | `approval/enums.go:257` |
| `AssigneeSuperior (AssigneeKind)` | `superior` |  | `approval/enums.go:258` |
| `AssigneeUser (AssigneeKind)` | `user` |  | `approval/enums.go:254` |
| `Between (Operator)` | `between` |  | `search/constants.go:13` |
| `BindingBusiness (BindingMode)` | `business` |  | `approval/enums.go:15` |
| `BindingProjectionApplied (BindingProjectionStatus)` | `applied` |  | `approval/binding.go:64` |
| `BindingProjectionFailed (BindingProjectionStatus)` | `failed` |  | `approval/binding.go:65` |
| `BindingProjectionPending (BindingProjectionStatus)` | `pending` |  | `approval/binding.go:62` |
| `BindingProjectionProcessing (BindingProjectionStatus)` | `processing` |  | `approval/binding.go:63` |
| `BindingStandalone (BindingMode)` | `standalone` |  | `approval/enums.go:14` |
| `BindingTriggerCompleted (BindingTrigger)` | `completed` |  | `approval/binding.go:90` |
| `BindingTriggerResubmitted (BindingTrigger)` | `resubmitted` |  | `approval/binding.go:93` |
| `BindingTriggerReturned (BindingTrigger)` | `returned` |  | `approval/binding.go:91` |
| `BindingTriggerStarted (BindingTrigger)` | `started` |  | `approval/binding.go:89` |
| `BindingTriggerWithdrawn (BindingTrigger)` | `withdrawn` |  | `approval/binding.go:92` |
| `CCDepartment (CCKind)` | `department` |  | `approval/enums.go:368` |
| `CCFormField (CCKind)` | `form_field` |  | `approval/enums.go:369` |
| `CCRole (CCKind)` | `role` |  | `approval/enums.go:367` |
| `CCTimingAlways (CCTiming)` | `always` |  | `approval/enums.go:381` |
| `CCTimingOnApprove (CCTiming)` | `on_approve` |  | `approval/enums.go:382` |
| `CCTimingOnReject (CCTiming)` | `on_reject` |  | `approval/enums.go:383` |
| `CCUser (CCKind)` | `user` |  | `approval/enums.go:366` |
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
| `ColumnBoolean (ColumnDataType)` | `boolean` |  | `approval/enums.go:433` |
| `ColumnDate (ColumnDataType)` | `date` |  | `approval/enums.go:434` |
| `ColumnDatetime (ColumnDataType)` | `datetime` |  | `approval/enums.go:435` |
| `ColumnDecimal (ColumnDataType)` | `decimal` |  | `approval/enums.go:432` |
| `ColumnInteger (ColumnDataType)` | `integer` |  | `approval/enums.go:431` |
| `ColumnJSON (ColumnDataType)` | `json` |  | `approval/enums.go:436` |
| `ColumnString (ColumnDataType)` | `string` |  | `approval/enums.go:429` |
| `ColumnText (ColumnDataType)` | `text` |  | `approval/enums.go:430` |
| `ConcurrencyAllow (ConcurrencyPolicy)` | `allow` |  | `cron/schedule.go:34` |
| `ConcurrencyForbid (ConcurrencyPolicy)` | `forbid` |  | `cron/schedule.go:32` |
| `ConditionExpression (ConditionKind)` | `expression` |  | `approval/enums.go:338` |
| `ConditionField (ConditionKind)` | `field` |  | `approval/enums.go:337` |
| `ConsecutiveApproverAutoPass (ConsecutiveApproverAction)` | `auto_pass` |  | `approval/enums.go:242` |
| `ConsecutiveApproverNone (ConsecutiveApproverAction)` | `none` |  | `approval/enums.go:241` |
| `Contains (Operator)` | `contains` |  | `search/constants.go:22` |
| `ContainsIgnoreCase (Operator)` | `iContains` |  | `search/constants.go:29` |
| `DataSourceModeReadOnly (DataSourceMode)` | `read_only` |  | `integration/models.go:80` |
| `DataSourceModeReadWrite (DataSourceMode)` | `read_write` |  | `integration/models.go:83` |
| `DeleteReasonClaimExpired (DeleteReason)` | `claim_expired` |  | `storage/delete_enqueuer.go:26` |
| `DeleteReasonDeleted (DeleteReason)` | `deleted` |  | `storage/delete_enqueuer.go:21` |
| `DeleteReasonReplaced (DeleteReason)` | `replaced` |  | `storage/delete_enqueuer.go:18` |
| `DirectionInbound (Direction)` | `inbound` |  | `integration/models.go:217` |
| `DirectionOutbound (Direction)` | `outbound` |  | `integration/models.go:213` |
| `EcdsaCurveP224 (ECDSACurve)` | `P224` |  | `cryptox/ecdsa_cipher.go:20` |
| `EcdsaCurveP256 (ECDSACurve)` | `P256` |  | `cryptox/ecdsa_cipher.go:21` |
| `EcdsaCurveP384 (ECDSACurve)` | `P384` |  | `cryptox/ecdsa_cipher.go:22` |
| `EcdsaCurveP521 (ECDSACurve)` | `P521` |  | `cryptox/ecdsa_cipher.go:23` |
| `EciesCurveP256 (ECIESCurve)` | `P256` |  | `cryptox/ecies_cipher.go:20` |
| `EciesCurveP384 (ECIESCurve)` | `P384` |  | `cryptox/ecies_cipher.go:21` |
| `EciesCurveP521 (ECIESCurve)` | `P521` |  | `cryptox/ecies_cipher.go:22` |
| `EciesCurveX25519 (ECIESCurve)` | `X25519` |  | `cryptox/ecies_cipher.go:23` |
| `EmptyAssigneeAutoPass (EmptyAssigneeAction)` | `auto_pass` |  | `approval/enums.go:136` |
| `EmptyAssigneeTransferAdmin (EmptyAssigneeAction)` | `transfer_admin` |  | `approval/enums.go:137` |
| `EmptyAssigneeTransferApplicant (EmptyAssigneeAction)` | `transfer_applicant` |  | `approval/enums.go:139` |
| `EmptyAssigneeTransferSpecified (EmptyAssigneeAction)` | `transfer_specified` |  | `approval/enums.go:140` |
| `EmptyAssigneeTransferSuperior (EmptyAssigneeAction)` | `transfer_superior` |  | `approval/enums.go:138` |
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
| `ExecutionAutoPass (ExecutionType)` | `auto_pass` |  | `approval/enums.go:92` |
| `ExecutionAutoReject (ExecutionType)` | `auto_reject` |  | `approval/enums.go:93` |
| `ExecutionManual (ExecutionType)` | `manual` |  | `approval/enums.go:91` |
| `FailureAuth (FailureKind)` | `auth` |  | `integration/failure.go:37` |
| `FailureCanceled (FailureKind)` | `canceled` |  | `integration/failure.go:26` |
| `FailureConfig (FailureKind)` | `config` |  | `integration/failure.go:34` |
| `FailureHandler (FailureKind)` | `handler` |  | `integration/failure.go:41` |
| `FailureInputInvalid (FailureKind)` | `input_invalid` |  | `integration/failure.go:11` |
| `FailureOutputInvalid (FailureKind)` | `output_invalid` |  | `integration/failure.go:14` |
| `FailureScript (FailureKind)` | `script` |  | `integration/failure.go:30` |
| `FailureTimeout (FailureKind)` | `timeout` |  | `integration/failure.go:22` |
| `FailureTransport (FailureKind)` | `transport` |  | `integration/failure.go:20` |
| `FailureUpstream (FailureKind)` | `upstream` |  | `integration/failure.go:17` |
| `FieldDate (FieldKind)` | `date` |  | `approval/enums.go:399` |
| `FieldInput (FieldKind)` | `input` |  | `approval/enums.go:395` |
| `FieldNumber (FieldKind)` | `number` |  | `approval/enums.go:398` |
| `FieldSelect (FieldKind)` | `select` |  | `approval/enums.go:397` |
| `FieldTable (FieldKind)` | `table` |  | `approval/enums.go:405` |
| `FieldTextarea (FieldKind)` | `textarea` |  | `approval/enums.go:396` |
| `FieldUpload (FieldKind)` | `upload` |  | `approval/enums.go:400` |
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
| `InitiatorDepartment (InitiatorKind)` | `department` |  | `approval/enums.go:38` |
| `InitiatorRole (InitiatorKind)` | `role` |  | `approval/enums.go:37` |
| `InitiatorUser (InitiatorKind)` | `user` |  | `approval/enums.go:36` |
| `InstanceApproved (InstanceStatus)` | `approved` |  | `approval/enums.go:279` |
| `InstanceRejected (InstanceStatus)` | `rejected` |  | `approval/enums.go:280` |
| `InstanceReturned (InstanceStatus)` | `returned` |  | `approval/enums.go:282` |
| `InstanceRunning (InstanceStatus)` | `running` |  | `approval/enums.go:278` |
| `InstanceTerminated (InstanceStatus)` | `terminated` |  | `approval/enums.go:283` |
| `InstanceWithdrawn (InstanceStatus)` | `withdrawn` |  | `approval/enums.go:281` |
| `IntegrationLogAll (IntegrationLogMode)` | `all` |  | `config/integration.go:19` |
| `IntegrationLogErrors (IntegrationLogMode)` | `errors` |  | `config/integration.go:17` |
| `IntegrationLogOff (IntegrationLogMode)` | `off` |  | `config/integration.go:15` |
| `IntegrationSecretAlgorithmAES (IntegrationSecretAlgorithm)` | `aes` |  | `config/integration.go:63` |
| `IntegrationSecretAlgorithmSM4 (IntegrationSecretAlgorithm)` | `sm4` |  | `config/integration.go:66` |
| `IsNotNull (Operator)` | `isNotNull` |  | `search/constants.go:20` |
| `IsNull (Operator)` | `isNull` |  | `search/constants.go:19` |
| `LessThan (Operator)` | `lt` |  | `search/constants.go:10` |
| `LessThanOrEqual (Operator)` | `lte` |  | `search/constants.go:11` |
| `LockoutKeyIP (LockoutKey)` | `ip` |  | `security/login_guard.go:28` |
| `LockoutKeyIP (LockoutKey)` | `ip` |  | `config/security.go:231` |
| `LockoutKeyUser (LockoutKey)` | `user` |  | `security/login_guard.go:26` |
| `LockoutKeyUser (LockoutKey)` | `user` |  | `config/security.go:229` |
| `LockoutKeyUserIP (LockoutKey)` | `user_ip` |  | `security/login_guard.go:30` |
| `LockoutKeyUserIP (LockoutKey)` | `user_ip` |  | `config/security.go:235` |
| `LockoutStrategyBackoff (LockoutStrategy)` | `backoff` |  | `security/login_guard.go:18` |
| `LockoutStrategyBackoff (LockoutStrategy)` | `backoff` |  | `config/security.go:221` |
| `LockoutStrategyLock (LockoutStrategy)` | `lock` |  | `config/security.go:216` |
| `LockoutStrategyLock (LockoutStrategy)` | `lock` |  | `security/login_guard.go:15` |
| `MetaTypeMarkdown (MetaType)` | `markdown` |  | `storage/file_refs.go:23` |
| `MetaTypeRichText (MetaType)` | `rich_text` |  | `storage/file_refs.go:21` |
| `MetaTypeUploadedFile (MetaType)` | `uploaded_file` |  | `storage/file_refs.go:19` |
| `MisfireFireNow (MisfirePolicy)` | `fire_now` |  | `cron/schedule.go:19` |
| `MisfireSkip (MisfirePolicy)` | `skip` |  | `cron/schedule.go:21` |
| `MySQL (DBKind)` | `mysql` |  | `config/data_sources.go:19` |
| `NodeApproval (NodeKind)` | `approval` |  | `approval/enums.go:78` |
| `NodeCC (NodeKind)` | `cc` |  | `approval/enums.go:82` |
| `NodeCondition (NodeKind)` | `condition` |  | `approval/enums.go:80` |
| `NodeEnd (NodeKind)` | `end` |  | `approval/enums.go:81` |
| `NodeHandle (NodeKind)` | `handle` |  | `approval/enums.go:79` |
| `NodeProgressActive (NodeProgressStatus)` | `active` |  | `approval/flow_graph_view.go:19` |
| `NodeProgressCanceled (NodeProgressStatus)` | `canceled` |  | `approval/flow_graph_view.go:28` |
| `NodeProgressPassed (NodeProgressStatus)` | `passed` |  | `approval/flow_graph_view.go:21` |
| `NodeProgressPending (NodeProgressStatus)` | `pending` |  | `approval/flow_graph_view.go:17` |
| `NodeProgressRejected (NodeProgressStatus)` | `rejected` |  | `approval/flow_graph_view.go:23` |
| `NodeProgressReturned (NodeProgressStatus)` | `returned` |  | `approval/flow_graph_view.go:25` |
| `NodeStart (NodeKind)` | `start` |  | `approval/enums.go:77` |
| `NodeVisitActive (NodeVisitStatus)` | `active` |  | `approval/enums.go:324` |
| `NodeVisitCanceled (NodeVisitStatus)` | `canceled` |  | `approval/enums.go:328` |
| `NodeVisitPassed (NodeVisitStatus)` | `passed` |  | `approval/enums.go:325` |
| `NodeVisitRejected (NodeVisitStatus)` | `rejected` |  | `approval/enums.go:326` |
| `NodeVisitReturned (NodeVisitStatus)` | `returned` |  | `approval/enums.go:327` |
| `NotBetween (Operator)` | `notBetween` |  | `search/constants.go:14` |
| `NotContains (Operator)` | `notContains` |  | `search/constants.go:23` |
| `NotContainsIgnoreCase (Operator)` | `iNotContains` |  | `search/constants.go:30` |
| `NotEndsWith (Operator)` | `notEndsWith` |  | `search/constants.go:27` |
| `NotEndsWithIgnoreCase (Operator)` | `iNotEndsWith` |  | `search/constants.go:34` |
| `NotEquals (Operator)` | `neq` |  | `search/constants.go:7` |
| `NotIn (Operator)` | `notIn` |  | `search/constants.go:17` |
| `NotStartsWith (Operator)` | `notStartsWith` |  | `search/constants.go:25` |
| `NotStartsWithIgnoreCase (Operator)` | `iNotStartsWith` |  | `search/constants.go:32` |
| `OperatorContains (ConditionOperator)` | `contains` |  | `approval/condition.go:24` |
| `OperatorEndsWith (ConditionOperator)` | `ends_with` |  | `approval/condition.go:27` |
| `OperatorEquals (ConditionOperator)` | `eq` |  | `approval/condition.go:16` |
| `OperatorGreater (ConditionOperator)` | `gt` |  | `approval/condition.go:18` |
| `OperatorGreaterOrEq (ConditionOperator)` | `gte` |  | `approval/condition.go:19` |
| `OperatorIn (ConditionOperator)` | `in` |  | `approval/condition.go:22` |
| `OperatorIsEmpty (ConditionOperator)` | `is_empty` |  | `approval/condition.go:28` |
| `OperatorIsNotEmpty (ConditionOperator)` | `is_not_empty` |  | `approval/condition.go:29` |
| `OperatorLess (ConditionOperator)` | `lt` |  | `approval/condition.go:20` |
| `OperatorLessOrEq (ConditionOperator)` | `lte` |  | `approval/condition.go:21` |
| `OperatorNotContains (ConditionOperator)` | `not_contains` |  | `approval/condition.go:25` |
| `OperatorNotEquals (ConditionOperator)` | `ne` |  | `approval/condition.go:17` |
| `OperatorNotIn (ConditionOperator)` | `not_in` |  | `approval/condition.go:23` |
| `OperatorStartsWith (ConditionOperator)` | `starts_with` |  | `approval/condition.go:26` |
| `Oracle (DBKind)` | `oracle` |  | `config/data_sources.go:16` |
| `OverflowError (OverflowStrategy)` | `error` |  | `sequence/rule.go:22` |
| `OverflowExtend (OverflowStrategy)` | `extend` |  | `sequence/rule.go:27` |
| `OverflowReset (OverflowStrategy)` | `reset` |  | `sequence/rule.go:24` |
| `PassAll (PassRule)` | `all` |  | `approval/enums.go:122` |
| `PassAny (PassRule)` | `any` |  | `approval/enums.go:123` |
| `PassRatio (PassRule)` | `ratio` |  | `approval/enums.go:124` |
| `PermissionEditable (Permission)` | `editable` |  | `approval/enums.go:477` |
| `PermissionHidden (Permission)` | `hidden` |  | `approval/enums.go:478` |
| `PermissionRequired (Permission)` | `required` |  | `approval/enums.go:479` |
| `PermissionVisible (Permission)` | `visible` |  | `approval/enums.go:476` |
| `Postgres (DBKind)` | `postgres` |  | `config/data_sources.go:18` |
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
| `RollbackAny (RollbackType)` | `any` |  | `approval/enums.go:175` |
| `RollbackDataClear (RollbackDataStrategy)` | `clear` |  | `approval/enums.go:193` |
| `RollbackDataKeep (RollbackDataStrategy)` | `keep` |  | `approval/enums.go:194` |
| `RollbackNone (RollbackType)` | `none` |  | `approval/enums.go:172` |
| `RollbackPrevious (RollbackType)` | `previous` |  | `approval/enums.go:173` |
| `RollbackSpecified (RollbackType)` | `specified` |  | `approval/enums.go:176` |
| `RollbackStart (RollbackType)` | `start` |  | `approval/enums.go:174` |
| `RouteFindingDanglingAdapter (RouteFindingKind)` | `dangling_adapter` |  | `integration/diagnostics.go:13` |
| `RouteFindingDisabledContract (RouteFindingKind)` | `disabled_contract` |  | `integration/diagnostics.go:24` |
| `RouteFindingDisabledSystem (RouteFindingKind)` | `disabled_system` |  | `integration/diagnostics.go:21` |
| `RouteFindingUncoveredContract (RouteFindingKind)` | `uncovered_contract` |  | `integration/diagnostics.go:29` |
| `RouteFindingWildcardGap (RouteFindingKind)` | `wildcard_gap` |  | `integration/diagnostics.go:18` |
| `RsaModeOAEP (RSAMode)` | `OAEP` |  | `cryptox/rsa_cipher.go:18` |
| `RsaModePKCS1v15 (RSAMode)` | `PKCS1v15` |  | `cryptox/rsa_cipher.go:24` |
| `RsaSignModePKCS1v15 (RSASignMode)` | `PKCS1v15` |  | `cryptox/rsa_cipher.go:31` |
| `RsaSignModePSS (RSASignMode)` | `PSS` |  | `cryptox/rsa_cipher.go:30` |
| `RunAbandoned (RunStatus)` | `abandoned` |  | `cron/run.go:23` |
| `RunCanceled (RunStatus)` | `canceled` |  | `cron/run.go:25` |
| `RunFailed (RunStatus)` | `failed` |  | `cron/run.go:15` |
| `RunMissed (RunStatus)` | `missed` |  | `cron/run.go:18` |
| `RunRunning (RunStatus)` | `running` |  | `cron/run.go:10` |
| `RunSkipped (RunStatus)` | `skipped` |  | `cron/run.go:20` |
| `RunSucceeded (RunStatus)` | `succeeded` |  | `cron/run.go:12` |
| `SQLServer (DBKind)` | `sqlserver` |  | `config/data_sources.go:17` |
| `SQLite (DBKind)` | `sqlite` |  | `config/data_sources.go:20` |
| `SSLDisable (SSLMode)` | `disable` |  | `config/data_sources.go:35` |
| `SSLRequire (SSLMode)` | `require` |  | `config/data_sources.go:38` |
| `SSLVerifyCA (SSLMode)` | `verify-ca` |  | `config/data_sources.go:41` |
| `SSLVerifyFull (SSLMode)` | `verify-full` |  | `config/data_sources.go:44` |
| `SameApplicantAutoPass (SameApplicantAction)` | `auto_pass` |  | `approval/enums.go:158` |
| `SameApplicantSelfApprove (SameApplicantAction)` | `self_approve` |  | `approval/enums.go:159` |
| `SameApplicantTransferSuperior (SameApplicantAction)` | `transfer_superior` |  | `approval/enums.go:160` |
| `SessionExceedEvictOldest (SessionExceedPolicy)` | `evict_oldest` |  | `config/security.go:104` |
| `SessionExceedEvictOldest (SessionExceedPolicy)` | `evict_oldest` |  | `security/session.go:24` |
| `SessionExceedReject (SessionExceedPolicy)` | `reject` |  | `config/security.go:102` |
| `SessionExceedReject (SessionExceedPolicy)` | `reject` |  | `security/session.go:21` |
| `SignatureAlgHmacSHA256 (SignatureAlgorithm)` | `HMAC-SHA256` |  | `security/signature.go:32` |
| `SignatureAlgHmacSHA512 (SignatureAlgorithm)` | `HMAC-SHA512` |  | `security/signature.go:33` |
| `SignatureAlgHmacSM3 (SignatureAlgorithm)` | `HMAC-SM3` |  | `security/signature.go:34` |
| `Sm4ModeCbc (SM4Mode)` | `CBC` |  | `cryptox/sm4_cipher.go:17` |
| `Sm4ModeGcm (SM4Mode)` | `GCM` |  | `cryptox/sm4_cipher.go:18` |
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
| `StorageJSON (StorageMode)` | `json` |  | `approval/enums.go:58` |
| `StorageMemory (StorageProvider)` | `memory` |  | `config/storage.go:14` |
| `StorageMinIO (StorageProvider)` | `minio` |  | `config/storage.go:13` |
| `StorageTable (StorageMode)` | `table` |  | `approval/enums.go:64` |
| `TargetBroadcast (TargetKind)` | `broadcast` |  | `push/target.go:12` |
| `TargetRoles (TargetKind)` | `roles` |  | `push/target.go:10` |
| `TargetUsers (TargetKind)` | `users` |  | `push/target.go:8` |
| `TaskApproved (TaskStatus)` | `approved` |  | `approval/enums.go:297` |
| `TaskCanceled (TaskStatus)` | `canceled` |  | `approval/enums.go:302` |
| `TaskHandled (TaskStatus)` | `handled` |  | `approval/enums.go:299` |
| `TaskPending (TaskStatus)` | `pending` |  | `approval/enums.go:296` |
| `TaskRejected (TaskStatus)` | `rejected` |  | `approval/enums.go:298` |
| `TaskRemoved (TaskStatus)` | `removed` |  | `approval/enums.go:303` |
| `TaskRolledBack (TaskStatus)` | `rolled_back` |  | `approval/enums.go:301` |
| `TaskSkipped (TaskStatus)` | `skipped` |  | `approval/enums.go:304` |
| `TaskTransferred (TaskStatus)` | `transferred` |  | `approval/enums.go:300` |
| `TaskWaiting (TaskStatus)` | `waiting` |  | `approval/enums.go:295` |
| `TimelineEntryApproval (TimelineEntryKind)` | `approval` |  | `approval/timeline_view.go:17` |
| `TimelineEntryCC (TimelineEntryKind)` | `cc` |  | `approval/timeline_view.go:19` |
| `TimelineEntryHandle (TimelineEntryKind)` | `handle` |  | `approval/timeline_view.go:18` |
| `TimelineEntryStart (TimelineEntryKind)` | `start` |  | `approval/timeline_view.go:16` |
| `TimelineEntryTerminate (TimelineEntryKind)` | `terminate` |  | `approval/timeline_view.go:21` |
| `TimelineEntryWithdraw (TimelineEntryKind)` | `withdraw` |  | `approval/timeline_view.go:20` |
| `TimeoutActionAutoPass (TimeoutAction)` | `auto_pass` |  | `approval/enums.go:455` |
| `TimeoutActionAutoReject (TimeoutAction)` | `auto_reject` |  | `approval/enums.go:456` |
| `TimeoutActionNone (TimeoutAction)` | `none` |  | `approval/enums.go:454` |
| `TimeoutActionNotify (TimeoutAction)` | `notify` |  | `approval/enums.go:457` |
| `TimeoutActionTransferAdmin (TimeoutAction)` | `transfer_admin` |  | `approval/enums.go:458` |
| `TokenTypeJWT (TokenType)` | `jwt_token` |  | `config/security.go:82` |
| `TokenTypeOpaque (TokenType)` | `opaque_token` |  | `config/security.go:84` |
| `TriggerCron (TriggerKind)` | `cron` |  | `cron/trigger.go:22` |
| `TriggerInterval (TriggerKind)` | `interval` |  | `cron/trigger.go:25` |
| `TriggerOnce (TriggerKind)` | `once` |  | `cron/trigger.go:27` |
| `UnmappedPolicyFallback (UnmappedPolicy)` | `fallback` |  | `integration/code_map.go:19` |
| `UnmappedPolicyPassthrough (UnmappedPolicy)` | `passthrough` |  | `integration/code_map.go:16` |
| `UnmappedPolicyReject (UnmappedPolicy)` | `reject` |  | `integration/code_map.go:14` |
| `UserMenuTypeDashboard (UserMenuType)` | `dashboard` |  | `security/user_info.go:28` |
| `UserMenuTypeDirectory (UserMenuType)` | `directory` |  | `security/user_info.go:25` |
| `UserMenuTypeMenu (UserMenuType)` | `menu` |  | `security/user_info.go:26` |
| `UserMenuTypeReport (UserMenuType)` | `report` |  | `security/user_info.go:29` |
| `UserMenuTypeView (UserMenuType)` | `view` |  | `security/user_info.go:27` |
| `VersionArchived (VersionStatus)` | `archived` |  | `approval/enums.go:29` |
| `VersionDraft (VersionStatus)` | `draft` |  | `approval/enums.go:27` |
| `VersionPublished (VersionStatus)` | `published` |  | `approval/enums.go:28` |
| `actionApprove (processTaskAction)` | `approve` |  | `internal/approval/resource/instance.go:199` |
| `actionHandle (processTaskAction)` | `handle` |  | `internal/approval/resource/instance.go:200` |
| `actionReject (processTaskAction)` | `reject` |  | `internal/approval/resource/instance.go:201` |
| `actionRollback (processTaskAction)` | `rollback` |  | `internal/approval/resource/instance.go:203` |
| `actionTransfer (processTaskAction)` | `transfer` |  | `internal/approval/resource/instance.go:202` |
| `columnBool (schemaColumnKind)` | `bool` |  | `internal/cron/store/migration/schema.go:29` |
| `columnInt32 (schemaColumnKind)` | `int32` |  | `internal/cron/store/migration/schema.go:28` |
| `columnInt64 (schemaColumnKind)` | `int64` |  | `internal/cron/store/migration/schema.go:27` |
| `columnJSON (schemaColumnKind)` | `json` |  | `internal/cron/store/migration/schema.go:30` |
| `columnText (schemaColumnKind)` | `text` |  | `internal/cron/store/migration/schema.go:31` |
| `columnTimestamp (schemaColumnKind)` | `timestamp` |  | `internal/cron/store/migration/schema.go:26` |
| `columnVarchar (schemaColumnKind)` | `varchar` |  | `internal/cron/store/migration/schema.go:25` |
| `fireRequestManual (fireRequestKind)` | `manual` |  | `internal/cron/store/fire_request.go:12` |
| `fireRequestRecovery (fireRequestKind)` | `recovery` |  | `internal/cron/store/fire_request.go:13` |

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
| `TypeDate` | `date` |  | `search/constants.go:52` |
| `TypeDateTime` | `datetime` |  | `search/constants.go:53` |
| `TypeDecimal` | `dec` |  | `search/constants.go:51` |
| `TypeInt` | `int` |  | `search/constants.go:50` |
| `TypeTime` | `time` |  | `search/constants.go:54` |
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
