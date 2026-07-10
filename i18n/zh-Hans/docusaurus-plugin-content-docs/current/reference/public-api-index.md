---
sidebar_position: 4
---

# 公开 API 索引

这一页是从当前 VEF Framework Go 源码生成的审计索引。它列出所有非 `internal/` 包中的 exported symbols，以及 exported 类型上的公开字段和方法，用来给文档审查提供一个可核对的无遗漏清单。

各专题文档仍然负责说明推荐用法和示例。出现在这个索引中并不代表稳定性承诺；支持的 API contract 以专题文档为准。

本次文档审计采用的覆盖规则：

- 所有非 `internal/` exported symbol 都列在本页，包括 exported 字段和方法
- exported 常量值、exported struct 字段顺序、exported struct 字段 tag，以及无歧义的 promoted exported fields 都会进入生成签名
- 除 `cmd/vef-cli/**` 外，所有 exported top-level symbol 都必须在专题或参考文档中出现
- 每一个 exported top-level symbol、field 和 method 都必须在 `scripts/api-audit-ledger.json` 中有 documented、grouped 或 excluded 审计处置
- `scripts/api-contract-ledger.json` 中的每个 package review 都会固定当前源码已审查的 top-level/field/method 数量和 fingerprint；一旦漂移，审计失败
- 不属于普通 exported Go symbol 的用户可见运行时 contract 单独由 [Runtime API Index](./runtime-api-index) 跟踪
- `cmd/vef-cli/**` 条目只为导出审计完整性保留；它们不是受支持的 import API，用户应按 [CLI 工具](../advanced/cli-tools) 记录的命令面使用

当框架源码的公开面发生变化时，请同步重新生成索引和审计账本：

```bash
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/gen-api-index.go -source . -out ../vef-framework-go-docs)
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-api-audit.go -source . -manifest ../vef-framework-go-docs/scripts/api-audit-manifest.json -ledger ../vef-framework-go-docs/scripts/api-audit-ledger.json -write-ledger)
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-api-audit.go -source . -manifest ../vef-framework-go-docs/scripts/api-audit-manifest.json -ledger ../vef-framework-go-docs/scripts/api-audit-ledger.json -contract-ledger ../vef-framework-go-docs/scripts/api-contract-ledger.json)
```

```text

## github.com/coldsmirk/vef-framework-go
VAR Annotate : func(t interface{}, anns ...go.uber.org/fx.Annotation) interface{}
VAR ApprovalModule : go.uber.org/fx.Option
VAR As : func(interfaces ...interface{}) go.uber.org/fx.Annotation
VAR Decorate : func(decorators ...interface{}) go.uber.org/fx.Option
VAR From : func(interfaces ...interface{}) go.uber.org/fx.Annotation
TYPE Hook : github.com/coldsmirk/vef-framework-go.Hook
TYPE HookFunc : github.com/coldsmirk/vef-framework-go.HookFunc
TYPE In : github.com/coldsmirk/vef-framework-go.In
VAR Invoke : func(funcs ...interface{}) go.uber.org/fx.Option
TYPE Lifecycle : github.com/coldsmirk/vef-framework-go.Lifecycle
  METHOD Append : func(go.uber.org/fx.Hook)
VAR Module : func(name string, opts ...go.uber.org/fx.Option) go.uber.org/fx.Option
FUNC NamedLogger : func(name string) github.com/coldsmirk/vef-framework-go/logx.Logger
VAR OnStart : func(onStart interface{}) go.uber.org/fx.Annotation
VAR OnStop : func(onStop interface{}) go.uber.org/fx.Annotation
TYPE Out : github.com/coldsmirk/vef-framework-go.Out
VAR ParamTags : func(tags ...string) go.uber.org/fx.Annotation
VAR Populate : func(targets ...interface{}) go.uber.org/fx.Option
VAR Private : go.uber.org/fx.privateOption
VAR Provide : func(constructors ...interface{}) go.uber.org/fx.Option
FUNC ProvideAPIResource : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideApprovalAggregator : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideApprovalFormSchemaParser : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideApprovalLifecycleHook : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideAuthStrategy : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideCQRSBehavior : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideChallengeProvider : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideDataSourceProvider : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideEventConsumeMiddleware : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideEventErrorSink : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideEventMetricsRecorder : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideEventPublishMiddleware : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideEventTransport : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideMCPPrompts : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideMCPResourceTemplates : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideMCPResources : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideMCPTools : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideMiddleware : func(constructor any, paramTags ...string) go.uber.org/fx.Option
FUNC ProvideSPAConfig : func(constructor any, paramTags ...string) go.uber.org/fx.Option
VAR Replace : func(values ...interface{}) go.uber.org/fx.Option
VAR ResultTags : func(tags ...string) go.uber.org/fx.Annotation
FUNC Run : func(options ...go.uber.org/fx.Option)
VAR Self : func() any
FUNC StartHook : func[T github.com/coldsmirk/vef-framework-go.HookFunc](start T) github.com/coldsmirk/vef-framework-go.Hook
FUNC StartStopHook : func[T1, T2 github.com/coldsmirk/vef-framework-go.HookFunc](start T1, stop T2) github.com/coldsmirk/vef-framework-go.Hook
FUNC StopHook : func[T github.com/coldsmirk/vef-framework-go.HookFunc](stop T) github.com/coldsmirk/vef-framework-go.Hook
VAR Supply : func(values ...interface{}) go.uber.org/fx.Option
FUNC SupplyBusinessRefProvider : func(constructor any) go.uber.org/fx.Option
FUNC SupplyBusinessRefResolver : func(constructor any) go.uber.org/fx.Option
FUNC SupplyFileACL : func(constructor any) go.uber.org/fx.Option
FUNC SupplyMCPServerInfo : func(info *github.com/coldsmirk/vef-framework-go/mcp.ServerInfo) go.uber.org/fx.Option
FUNC SupplySPAConfigs : func(config *github.com/coldsmirk/vef-framework-go/middleware.SPAConfig, configs ...*github.com/coldsmirk/vef-framework-go/middleware.SPAConfig) go.uber.org/fx.Option
FUNC SupplyURLKeyMapper : func(constructor any) go.uber.org/fx.Option

## github.com/coldsmirk/vef-framework-go/ai
TYPE Agent : github.com/coldsmirk/vef-framework-go/ai.Agent
  METHOD Run : func(ctx context.Context, input string, opts ...github.com/coldsmirk/vef-framework-go/ai.Option) (*github.com/coldsmirk/vef-framework-go/ai.Message, error)
  METHOD Stream : func(ctx context.Context, input string, opts ...github.com/coldsmirk/vef-framework-go/ai.Option) (github.com/coldsmirk/vef-framework-go/ai.MessageStream, error)
TYPE AgentBuilder : github.com/coldsmirk/vef-framework-go/ai.AgentBuilder
  METHOD Build : func(ctx context.Context) (github.com/coldsmirk/vef-framework-go/ai.Agent, error)
  METHOD WithMaxIterations : func(n int) github.com/coldsmirk/vef-framework-go/ai.AgentBuilder
  METHOD WithModel : func(model github.com/coldsmirk/vef-framework-go/ai.ToolableChatModel) github.com/coldsmirk/vef-framework-go/ai.AgentBuilder
  METHOD WithSystemPrompt : func(prompt string) github.com/coldsmirk/vef-framework-go/ai.AgentBuilder
  METHOD WithTools : func(tools ...github.com/coldsmirk/vef-framework-go/ai.Tool) github.com/coldsmirk/vef-framework-go/ai.AgentBuilder
TYPE AgentConfig : github.com/coldsmirk/vef-framework-go/ai.AgentConfig
  FIELD Model : github.com/coldsmirk/vef-framework-go/ai.ToolableChatModel [field_order=1 tag=""]
  FIELD Tools : []github.com/coldsmirk/vef-framework-go/ai.Tool [field_order=2 tag=""]
  FIELD SystemPrompt : string [field_order=3 tag=""]
  FIELD MaxIterations : int [field_order=4 tag=""]
TYPE AgentFactory : github.com/coldsmirk/vef-framework-go/ai.AgentFactory
  METHOD CreateBuilder : func() github.com/coldsmirk/vef-framework-go/ai.AgentBuilder
  METHOD Name : func() string
TYPE ChatModel : github.com/coldsmirk/vef-framework-go/ai.ChatModel
  METHOD Generate : func(ctx context.Context, messages []*github.com/coldsmirk/vef-framework-go/ai.Message, opts ...github.com/coldsmirk/vef-framework-go/ai.Option) (*github.com/coldsmirk/vef-framework-go/ai.Message, error)
  METHOD Stream : func(ctx context.Context, messages []*github.com/coldsmirk/vef-framework-go/ai.Message, opts ...github.com/coldsmirk/vef-framework-go/ai.Option) (github.com/coldsmirk/vef-framework-go/ai.MessageStream, error)
VAR ErrAgentNotFound : error
VAR ErrInvalidArguments : error
VAR ErrMaxIterationsReached : error
VAR ErrModelNotSupported : error
VAR ErrNoContent : error
VAR ErrProviderNotFound : error
VAR ErrStreamClosed : error
VAR ErrToolNotFound : error
FUNC ListAgentFactories : func() []string
FUNC ListModelProviders : func() []string
TYPE Message : github.com/coldsmirk/vef-framework-go/ai.Message
  FIELD Role : github.com/coldsmirk/vef-framework-go/ai.Role [field_order=1 tag=""]
  FIELD Content : string [field_order=2 tag=""]
  FIELD ToolCalls : []github.com/coldsmirk/vef-framework-go/ai.ToolCall [field_order=3 tag=""]
  FIELD ToolResult : *github.com/coldsmirk/vef-framework-go/ai.ToolResult [field_order=4 tag=""]
  FIELD Usage : *github.com/coldsmirk/vef-framework-go/ai.TokenUsage [field_order=5 tag=""]
  METHOD HasToolCalls : func() bool
  METHOD IsAssistant : func() bool
  METHOD IsSystem : func() bool
  METHOD IsTool : func() bool
  METHOD IsUser : func() bool
TYPE MessageChunk : github.com/coldsmirk/vef-framework-go/ai.MessageChunk
  FIELD Content : string [field_order=1 tag=""]
  FIELD ToolCalls : []github.com/coldsmirk/vef-framework-go/ai.ToolCall [field_order=2 tag=""]
  FIELD Done : bool [field_order=3 tag=""]
TYPE MessageStream : github.com/coldsmirk/vef-framework-go/ai.MessageStream
  METHOD Close : func() error
  METHOD Collect : func() (*github.com/coldsmirk/vef-framework-go/ai.Message, error)
  METHOD Recv : func() (*github.com/coldsmirk/vef-framework-go/ai.MessageChunk, error)
TYPE ModelConfig : github.com/coldsmirk/vef-framework-go/ai.ModelConfig
  FIELD Provider : string [field_order=1 tag=""]
  FIELD Model : string [field_order=2 tag=""]
  FIELD APIKey : string [field_order=3 tag=""]
  FIELD BaseURL : string [field_order=4 tag=""]
  FIELD Temperature : float64 [field_order=5 tag=""]
  FIELD MaxTokens : int [field_order=6 tag=""]
  FIELD Timeout : time.Duration [field_order=7 tag=""]
TYPE ModelError : github.com/coldsmirk/vef-framework-go/ai.ModelError
  FIELD Provider : string [field_order=1 tag=""]
  FIELD StatusCode : int [field_order=2 tag=""]
  FIELD Message : string [field_order=3 tag=""]
  METHOD Error : func() string
TYPE ModelInfo : github.com/coldsmirk/vef-framework-go/ai.ModelInfo
  FIELD Provider : string [field_order=1 tag=""]
  FIELD Model : string [field_order=2 tag=""]
  FIELD MaxTokens : int [field_order=3 tag=""]
  FIELD Temperature : float64 [field_order=4 tag=""]
TYPE ModelProvider : github.com/coldsmirk/vef-framework-go/ai.ModelProvider
  METHOD CreateModel : func(ctx context.Context, cfg *github.com/coldsmirk/vef-framework-go/ai.ModelConfig) (github.com/coldsmirk/vef-framework-go/ai.ToolableChatModel, error)
  METHOD Name : func() string
FUNC NewAgentBuilder : func(agentType string) (github.com/coldsmirk/vef-framework-go/ai.AgentBuilder, error)
FUNC NewAssistantMessage : func(content string) *github.com/coldsmirk/vef-framework-go/ai.Message
FUNC NewAssistantMessageWithToolCalls : func(content string, toolCalls []github.com/coldsmirk/vef-framework-go/ai.ToolCall) *github.com/coldsmirk/vef-framework-go/ai.Message
FUNC NewChatModel : func(ctx context.Context, cfg *github.com/coldsmirk/vef-framework-go/ai.ModelConfig) (github.com/coldsmirk/vef-framework-go/ai.ToolableChatModel, error)
FUNC NewModelError : func(provider string, statusCode int, message string) *github.com/coldsmirk/vef-framework-go/ai.ModelError
FUNC NewOptions : func() *github.com/coldsmirk/vef-framework-go/ai.Options
FUNC NewSystemMessage : func(content string) *github.com/coldsmirk/vef-framework-go/ai.Message
FUNC NewToolError : func(toolName string, err error) *github.com/coldsmirk/vef-framework-go/ai.ToolError
FUNC NewToolMessage : func(callID string, content string) *github.com/coldsmirk/vef-framework-go/ai.Message
FUNC NewUserMessage : func(content string) *github.com/coldsmirk/vef-framework-go/ai.Message
TYPE Option : github.com/coldsmirk/vef-framework-go/ai.Option
TYPE Options : github.com/coldsmirk/vef-framework-go/ai.Options
  FIELD Temperature : *float64 [field_order=1 tag=""]
  FIELD MaxTokens : *int [field_order=2 tag=""]
  FIELD StopSequences : []string [field_order=3 tag=""]
  FIELD Meta : map[string]string [field_order=4 tag=""]
  METHOD Apply : func(opts ...github.com/coldsmirk/vef-framework-go/ai.Option) *github.com/coldsmirk/vef-framework-go/ai.Options
TYPE ParameterSchema : github.com/coldsmirk/vef-framework-go/ai.ParameterSchema
  FIELD Type : string [field_order=1 tag=""]
  FIELD Properties : map[string]*github.com/coldsmirk/vef-framework-go/ai.PropertySchema [field_order=2 tag=""]
  FIELD Required : []string [field_order=3 tag=""]
TYPE PropertySchema : github.com/coldsmirk/vef-framework-go/ai.PropertySchema
  FIELD Type : string [field_order=1 tag=""]
  FIELD Description : string [field_order=2 tag=""]
  FIELD Enum : []string [field_order=3 tag=""]
  FIELD Items : *github.com/coldsmirk/vef-framework-go/ai.PropertySchema [field_order=4 tag=""]
FUNC RegisterAgentFactory : func(f github.com/coldsmirk/vef-framework-go/ai.AgentFactory)
FUNC RegisterModelProvider : func(p github.com/coldsmirk/vef-framework-go/ai.ModelProvider)
TYPE Role : github.com/coldsmirk/vef-framework-go/ai.Role
CONST RoleAssistant : github.com/coldsmirk/vef-framework-go/ai.Role = "assistant"
CONST RoleSystem : github.com/coldsmirk/vef-framework-go/ai.Role = "system"
CONST RoleTool : github.com/coldsmirk/vef-framework-go/ai.Role = "tool"
CONST RoleUser : github.com/coldsmirk/vef-framework-go/ai.Role = "user"
TYPE StreamableTool : github.com/coldsmirk/vef-framework-go/ai.StreamableTool
  METHOD Info : func() *github.com/coldsmirk/vef-framework-go/ai.ToolInfo
  METHOD Invoke : func(ctx context.Context, arguments string) (string, error)
  METHOD InvokeStream : func(ctx context.Context, arguments string) (github.com/coldsmirk/vef-framework-go/ai.StringStream, error)
TYPE StringStream : github.com/coldsmirk/vef-framework-go/ai.StringStream
  METHOD Close : func() error
  METHOD Collect : func() (string, error)
  METHOD Recv : func() (string, error)
TYPE TokenUsage : github.com/coldsmirk/vef-framework-go/ai.TokenUsage
  FIELD PromptTokens : int [field_order=1 tag=""]
  FIELD CompletionTokens : int [field_order=2 tag=""]
  FIELD TotalTokens : int [field_order=3 tag=""]
TYPE Tool : github.com/coldsmirk/vef-framework-go/ai.Tool
  METHOD Info : func() *github.com/coldsmirk/vef-framework-go/ai.ToolInfo
  METHOD Invoke : func(ctx context.Context, arguments string) (string, error)
TYPE ToolCall : github.com/coldsmirk/vef-framework-go/ai.ToolCall
  FIELD ID : string [field_order=1 tag=""]
  FIELD Name : string [field_order=2 tag=""]
  FIELD Arguments : string [field_order=3 tag=""]
TYPE ToolError : github.com/coldsmirk/vef-framework-go/ai.ToolError
  FIELD ToolName : string [field_order=1 tag=""]
  FIELD Err : error [field_order=2 tag=""]
  METHOD Error : func() string
  METHOD Unwrap : func() error
TYPE ToolInfo : github.com/coldsmirk/vef-framework-go/ai.ToolInfo
  FIELD Name : string [field_order=1 tag=""]
  FIELD Description : string [field_order=2 tag=""]
  FIELD Parameters : *github.com/coldsmirk/vef-framework-go/ai.ParameterSchema [field_order=3 tag=""]
TYPE ToolResult : github.com/coldsmirk/vef-framework-go/ai.ToolResult
  FIELD CallID : string [field_order=1 tag=""]
  FIELD Content : string [field_order=2 tag=""]
TYPE ToolableChatModel : github.com/coldsmirk/vef-framework-go/ai.ToolableChatModel
  METHOD Generate : func(ctx context.Context, messages []*github.com/coldsmirk/vef-framework-go/ai.Message, opts ...github.com/coldsmirk/vef-framework-go/ai.Option) (*github.com/coldsmirk/vef-framework-go/ai.Message, error)
  METHOD Stream : func(ctx context.Context, messages []*github.com/coldsmirk/vef-framework-go/ai.Message, opts ...github.com/coldsmirk/vef-framework-go/ai.Option) (github.com/coldsmirk/vef-framework-go/ai.MessageStream, error)
  METHOD WithTools : func(tools ...github.com/coldsmirk/vef-framework-go/ai.Tool) github.com/coldsmirk/vef-framework-go/ai.ToolableChatModel
FUNC WithMaxTokens : func(n int) github.com/coldsmirk/vef-framework-go/ai.Option
FUNC WithMeta : func(key string, value string) github.com/coldsmirk/vef-framework-go/ai.Option
FUNC WithStopSequences : func(seqs ...string) github.com/coldsmirk/vef-framework-go/ai.Option
FUNC WithTemperature : func(t float64) github.com/coldsmirk/vef-framework-go/ai.Option

## github.com/coldsmirk/vef-framework-go/ai/stream
TYPE Builder : github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD OnError : func(handler func(err error) string) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD OnFinish : func(handler func(content string)) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD Stream : func(ctx github.com/gofiber/fiber/v3.Ctx) error
  METHOD StreamToWriter : func(w *bufio.Writer)
  METHOD WithFinish : func(enabled bool) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD WithHeader : func(key string, value string) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD WithIDGenerator : func(gen func(prefix string) string) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD WithMessageID : func(id string) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD WithReasoning : func(enabled bool) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD WithSource : func(source github.com/coldsmirk/vef-framework-go/ai/stream.MessageSource) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD WithSources : func(enabled bool) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
  METHOD WithStart : func(enabled bool) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
TYPE CallbackWriter : github.com/coldsmirk/vef-framework-go/ai/stream.CallbackWriter
  METHOD WriteData : func(dataType string, data any)
  METHOD WriteMessage : func(msg github.com/coldsmirk/vef-framework-go/ai/stream.Message)
  METHOD WriteReasoning : func(reasoning string)
  METHOD WriteText : func(content string)
  METHOD WriteToolCall : func(id string, name string, arguments string)
  METHOD WriteToolResult : func(toolCallID string, content string)
TYPE Chunk : github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
TYPE ChunkType : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType
CONST ChunkTypeError : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "error"
CONST ChunkTypeFile : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "file"
CONST ChunkTypeFinish : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "finish"
CONST ChunkTypeFinishStep : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "finish-step"
CONST ChunkTypeReasoningDelta : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "reasoning-delta"
CONST ChunkTypeReasoningEnd : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "reasoning-end"
CONST ChunkTypeReasoningStart : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "reasoning-start"
CONST ChunkTypeSourceDocument : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "source-document"
CONST ChunkTypeSourceURL : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "source-url"
CONST ChunkTypeStart : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "start"
CONST ChunkTypeStartStep : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "start-step"
CONST ChunkTypeTextDelta : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "text-delta"
CONST ChunkTypeTextEnd : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "text-end"
CONST ChunkTypeTextStart : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "text-start"
CONST ChunkTypeToolInputAvailable : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "tool-input-available"
CONST ChunkTypeToolInputDelta : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "tool-input-delta"
CONST ChunkTypeToolInputStart : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "tool-input-start"
CONST ChunkTypeToolOutputAvailable : github.com/coldsmirk/vef-framework-go/ai/stream.ChunkType = "tool-output-available"
FUNC DefaultOptions : func() github.com/coldsmirk/vef-framework-go/ai/stream.Options
VAR ErrSourceClosed : error
VAR ErrSourceRequired : error
FUNC FromAiMessageStream : func(stream github.com/coldsmirk/vef-framework-go/ai.MessageStream) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
FUNC FromCallback : func(execute func(writer github.com/coldsmirk/vef-framework-go/ai/stream.CallbackWriter) error) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
FUNC FromChannel : func(ch <-chan github.com/coldsmirk/vef-framework-go/ai/stream.Message) *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
TYPE Message : github.com/coldsmirk/vef-framework-go/ai/stream.Message
  FIELD Role : github.com/coldsmirk/vef-framework-go/ai/stream.Role [field_order=1 tag=""]
  FIELD Content : string [field_order=2 tag=""]
  FIELD ToolCalls : []github.com/coldsmirk/vef-framework-go/ai/stream.ToolCall [field_order=3 tag=""]
  FIELD ToolCallID : string [field_order=4 tag=""]
  FIELD Reasoning : string [field_order=5 tag=""]
  FIELD Data : map[string]any [field_order=6 tag=""]
TYPE MessageSource : github.com/coldsmirk/vef-framework-go/ai/stream.MessageSource
  METHOD Close : func() error
  METHOD Recv : func() (github.com/coldsmirk/vef-framework-go/ai/stream.Message, error)
FUNC New : func() *github.com/coldsmirk/vef-framework-go/ai/stream.Builder
FUNC NewAiMessageStreamSource : func(stream github.com/coldsmirk/vef-framework-go/ai.MessageStream) github.com/coldsmirk/vef-framework-go/ai/stream.MessageSource
FUNC NewCallbackSource : func(execute func(writer github.com/coldsmirk/vef-framework-go/ai/stream.CallbackWriter) error) github.com/coldsmirk/vef-framework-go/ai/stream.MessageSource
FUNC NewChannelSource : func(ch <-chan github.com/coldsmirk/vef-framework-go/ai/stream.Message) github.com/coldsmirk/vef-framework-go/ai/stream.MessageSource
FUNC NewDataChunk : func(dataType string, data any) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewErrorChunk : func(errorText string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewFileChunk : func(fileID string, mediaType string, url string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewFinishChunk : func() github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewFinishStepChunk : func() github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewReasoningDeltaChunk : func(id string, delta string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewReasoningEndChunk : func(id string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewReasoningStartChunk : func(id string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewSourceDocumentChunk : func(sourceID string, mediaType string, title string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewSourceURLChunk : func(sourceID string, url string, title string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewStartChunk : func(messageID string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewStartStepChunk : func() github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewTextDeltaChunk : func(id string, delta string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewTextEndChunk : func(id string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewTextStartChunk : func(id string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewToolInputAvailableChunk : func(toolCallID string, toolName string, input any) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewToolInputDeltaChunk : func(toolCallID string, delta string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewToolInputStartChunk : func(toolCallID string, toolName string) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
FUNC NewToolOutputAvailableChunk : func(toolCallID string, output any) github.com/coldsmirk/vef-framework-go/ai/stream.Chunk
TYPE Options : github.com/coldsmirk/vef-framework-go/ai/stream.Options
  FIELD SendReasoning : bool [field_order=1 tag=""]
  FIELD SendSources : bool [field_order=2 tag=""]
  FIELD SendStart : bool [field_order=3 tag=""]
  FIELD SendFinish : bool [field_order=4 tag=""]
  FIELD OnError : func(err error) string [field_order=5 tag=""]
  FIELD OnFinish : func(content string) [field_order=6 tag=""]
  FIELD GenerateID : func(prefix string) string [field_order=7 tag=""]
TYPE ResponseWriter : github.com/coldsmirk/vef-framework-go/ai/stream.ResponseWriter
  METHOD Write : func(p []byte) (n int, err error)
TYPE Role : github.com/coldsmirk/vef-framework-go/ai/stream.Role
CONST RoleAssistant : github.com/coldsmirk/vef-framework-go/ai/stream.Role = "assistant"
CONST RoleSystem : github.com/coldsmirk/vef-framework-go/ai/stream.Role = "system"
CONST RoleTool : github.com/coldsmirk/vef-framework-go/ai/stream.Role = "tool"
CONST RoleUser : github.com/coldsmirk/vef-framework-go/ai/stream.Role = "user"
TYPE Source : github.com/coldsmirk/vef-framework-go/ai/stream.Source
  FIELD Type : string [field_order=1 tag=""]
  FIELD ID : string [field_order=2 tag=""]
  FIELD URL : string [field_order=3 tag=""]
  FIELD Title : string [field_order=4 tag=""]
  FIELD MediaType : string [field_order=5 tag=""]
VAR SseHeaders : map[string]string
TYPE StreamWriter : github.com/coldsmirk/vef-framework-go/ai/stream.StreamWriter
  METHOD Flush : func() error
  METHOD WriteChunk : func(chunk github.com/coldsmirk/vef-framework-go/ai/stream.Chunk) error
TYPE ToolCall : github.com/coldsmirk/vef-framework-go/ai/stream.ToolCall
  FIELD ID : string [field_order=1 tag=""]
  FIELD Name : string [field_order=2 tag=""]
  FIELD Arguments : string [field_order=3 tag=""]

## github.com/coldsmirk/vef-framework-go/api
TYPE AuditEvent : github.com/coldsmirk/vef-framework-go/api.AuditEvent
  FIELD Resource : string [field_order=1 tag="json:\"resource\""]
  FIELD Action : string [field_order=2 tag="json:\"action\""]
  FIELD Version : string [field_order=3 tag="json:\"version\""]
  FIELD UserID : string [field_order=4 tag="json:\"userId\""]
  FIELD UserAgent : string [field_order=5 tag="json:\"userAgent\""]
  FIELD RequestID : string [field_order=6 tag="json:\"requestId\""]
  FIELD RequestIP : string [field_order=7 tag="json:\"requestIp\""]
  FIELD RequestParams : map[string]any [field_order=8 tag="json:\"requestParams\""]
  FIELD RequestMeta : map[string]any [field_order=9 tag="json:\"requestMeta\""]
  FIELD ResultCode : int [field_order=10 tag="json:\"resultCode\""]
  FIELD ResultMessage : string [field_order=11 tag="json:\"resultMessage\""]
  FIELD ResultData : any [field_order=12 tag="json:\"resultData\""]
  FIELD ElapsedTime : int64 [field_order=13 tag="json:\"elapsedTime\""]
  METHOD EventType : func() string
TYPE AuthConfig : github.com/coldsmirk/vef-framework-go/api.AuthConfig
  FIELD Strategy : string [field_order=1 tag=""]
  FIELD Options : map[string]any [field_order=2 tag=""]
  METHOD Clone : func() *github.com/coldsmirk/vef-framework-go/api.AuthConfig
CONST AuthOptionWhitelist : untyped string = "whitelist"
TYPE AuthStrategy : github.com/coldsmirk/vef-framework-go/api.AuthStrategy
  METHOD Authenticate : func(ctx github.com/gofiber/fiber/v3.Ctx, options map[string]any) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
  METHOD Name : func() string
CONST AuthStrategyBearer : untyped string = "bearer"
CONST AuthStrategyIP : untyped string = "ip"
CONST AuthStrategyNone : untyped string = "none"
TYPE AuthStrategyRegistry : github.com/coldsmirk/vef-framework-go/api.AuthStrategyRegistry
  METHOD Get : func(name string) (github.com/coldsmirk/vef-framework-go/api.AuthStrategy, bool)
  METHOD Names : func() []string
  METHOD Register : func(strategy github.com/coldsmirk/vef-framework-go/api.AuthStrategy)
CONST AuthStrategySignature : untyped string = "signature"
FUNC BearerAuth : func() *github.com/coldsmirk/vef-framework-go/api.AuthConfig
CONST DefaultIPWhitelist : untyped string = "default"
TYPE Engine : github.com/coldsmirk/vef-framework-go/api.Engine
  METHOD Lookup : func(id github.com/coldsmirk/vef-framework-go/api.Identifier) *github.com/coldsmirk/vef-framework-go/api.Operation
  METHOD Mount : func(router github.com/gofiber/fiber/v3.Router) error
  METHOD Register : func(resources ...github.com/coldsmirk/vef-framework-go/api.Resource) error
VAR ErrEmptyActionName : error
VAR ErrEmptyResourceName : error
VAR ErrInvalidActionName : error
VAR ErrInvalidMetaType : error
VAR ErrInvalidParamsType : error
VAR ErrInvalidRequestMeta : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrInvalidRequestParams : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrInvalidResourceKind : error
VAR ErrInvalidResourceName : error
VAR ErrInvalidVersionFormat : error
VAR ErrResourceNameDoubleSlash : error
VAR ErrResourceNameSlash : error
TYPE FactoryParamResolver : github.com/coldsmirk/vef-framework-go/api.FactoryParamResolver
  METHOD Resolve : func() (reflect.Value, error)
  METHOD Type : func() reflect.Type
TYPE HandlerAdapter : github.com/coldsmirk/vef-framework-go/api.HandlerAdapter
  METHOD Adapt : func(handler any, op *github.com/coldsmirk/vef-framework-go/api.Operation) (github.com/gofiber/fiber/v3.Handler, error)
TYPE HandlerParamResolver : github.com/coldsmirk/vef-framework-go/api.HandlerParamResolver
  METHOD Resolve : func(ctx github.com/gofiber/fiber/v3.Ctx) (reflect.Value, error)
  METHOD Type : func() reflect.Type
TYPE HandlerResolver : github.com/coldsmirk/vef-framework-go/api.HandlerResolver
  METHOD Resolve : func(resource github.com/coldsmirk/vef-framework-go/api.Resource, spec github.com/coldsmirk/vef-framework-go/api.OperationSpec) (any, error)
CONST HeaderXAppID : untyped string = "X-App-ID"
CONST HeaderXMetaPrefix : untyped string = "X-Meta-"
CONST HeaderXNonce : untyped string = "X-Nonce"
CONST HeaderXSignature : untyped string = "X-Signature"
CONST HeaderXTimestamp : untyped string = "X-Timestamp"
FUNC IPAuth : func(whitelistName ...string) *github.com/coldsmirk/vef-framework-go/api.AuthConfig
TYPE Identifier : github.com/coldsmirk/vef-framework-go/api.Identifier
  FIELD Resource : string [field_order=1 tag="json:\"resource\" form:\"resource\" validate:\"required,alphanum_us_slash\" label_i18n:\"api_request_resource\""]
  FIELD Action : string [field_order=2 tag="json:\"action\" form:\"action\" validate:\"required\" label_i18n:\"api_request_action\""]
  FIELD Version : string [field_order=3 tag="json:\"version\" form:\"version\" validate:\"required,alphanum\" label_i18n:\"api_request_version\""]
  METHOD String : func() string
TYPE Kind : github.com/coldsmirk/vef-framework-go/api.Kind
  METHOD String : func() string
CONST KindREST : github.com/coldsmirk/vef-framework-go/api.Kind = 2
CONST KindRPC : github.com/coldsmirk/vef-framework-go/api.Kind = 1
TYPE M : github.com/coldsmirk/vef-framework-go/api.M
TYPE Meta : github.com/coldsmirk/vef-framework-go/api.Meta
  METHOD Decode : func(out any) error
TYPE Middleware : github.com/coldsmirk/vef-framework-go/api.Middleware
  METHOD Name : func() string
  METHOD Order : func() int
  METHOD Process : func(ctx github.com/gofiber/fiber/v3.Ctx) error
FUNC NewRESTResource : func(name string, opts ...github.com/coldsmirk/vef-framework-go/api.ResourceOption) github.com/coldsmirk/vef-framework-go/api.Resource
FUNC NewRPCResource : func(name string, opts ...github.com/coldsmirk/vef-framework-go/api.ResourceOption) github.com/coldsmirk/vef-framework-go/api.Resource
TYPE Operation : github.com/coldsmirk/vef-framework-go/api.Operation
  FIELD Identifier : github.com/coldsmirk/vef-framework-go/api.Identifier [field_order=1 tag=""]
  FIELD EnableAudit : bool [field_order=2 tag=""]
  FIELD Timeout : time.Duration [field_order=3 tag=""]
  FIELD Auth : *github.com/coldsmirk/vef-framework-go/api.AuthConfig [field_order=4 tag=""]
  FIELD RateLimit : *github.com/coldsmirk/vef-framework-go/api.RateLimitConfig [field_order=5 tag=""]
  FIELD Handler : any [field_order=6 tag=""]
  FIELD Meta : map[string]any [field_order=7 tag=""]
  FIELD Action : string [promoted_from=Identifier depth=1 field_order=2 tag="json:\"action\" form:\"action\" validate:\"required\" label_i18n:\"api_request_action\""]
  FIELD Resource : string [promoted_from=Identifier depth=1 field_order=1 tag="json:\"resource\" form:\"resource\" validate:\"required,alphanum_us_slash\" label_i18n:\"api_request_resource\""]
  FIELD Version : string [promoted_from=Identifier depth=1 field_order=3 tag="json:\"version\" form:\"version\" validate:\"required,alphanum\" label_i18n:\"api_request_version\""]
  METHOD String : func() string
TYPE OperationSpec : github.com/coldsmirk/vef-framework-go/api.OperationSpec
  FIELD Action : string [field_order=1 tag=""]
  FIELD EnableAudit : bool [field_order=2 tag=""]
  FIELD Timeout : time.Duration [field_order=3 tag=""]
  FIELD Public : bool [field_order=4 tag=""]
  FIELD RequiredPermission : string [field_order=5 tag=""]
  FIELD RateLimit : *github.com/coldsmirk/vef-framework-go/api.RateLimitConfig [field_order=6 tag=""]
  FIELD Handler : any [field_order=7 tag=""]
TYPE OperationsCollector : github.com/coldsmirk/vef-framework-go/api.OperationsCollector
  METHOD Collect : func(resource github.com/coldsmirk/vef-framework-go/api.Resource) []github.com/coldsmirk/vef-framework-go/api.OperationSpec
TYPE OperationsProvider : github.com/coldsmirk/vef-framework-go/api.OperationsProvider
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
TYPE P : github.com/coldsmirk/vef-framework-go/api.P
TYPE Params : github.com/coldsmirk/vef-framework-go/api.Params
  METHOD Decode : func(out any) error
FUNC Public : func() *github.com/coldsmirk/vef-framework-go/api.AuthConfig
TYPE RateLimitConfig : github.com/coldsmirk/vef-framework-go/api.RateLimitConfig
  FIELD Max : int [field_order=1 tag=""]
  FIELD Period : time.Duration [field_order=2 tag=""]
  FIELD Key : string [field_order=3 tag=""]
TYPE Request : github.com/coldsmirk/vef-framework-go/api.Request
  FIELD Identifier : github.com/coldsmirk/vef-framework-go/api.Identifier [field_order=1 tag=""]
  FIELD Params : github.com/coldsmirk/vef-framework-go/api.Params [field_order=2 tag="json:\"params\""]
  FIELD Meta : github.com/coldsmirk/vef-framework-go/api.Meta [field_order=3 tag="json:\"meta\""]
  FIELD Action : string [promoted_from=Identifier depth=1 field_order=2 tag="json:\"action\" form:\"action\" validate:\"required\" label_i18n:\"api_request_action\""]
  FIELD Resource : string [promoted_from=Identifier depth=1 field_order=1 tag="json:\"resource\" form:\"resource\" validate:\"required,alphanum_us_slash\" label_i18n:\"api_request_resource\""]
  FIELD Version : string [promoted_from=Identifier depth=1 field_order=3 tag="json:\"version\" form:\"version\" validate:\"required,alphanum\" label_i18n:\"api_request_version\""]
  METHOD String : func() string
TYPE Resource : github.com/coldsmirk/vef-framework-go/api.Resource
  METHOD Auth : func() *github.com/coldsmirk/vef-framework-go/api.AuthConfig
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/api.Kind
  METHOD Name : func() string
  METHOD Operations : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Version : func() string
TYPE ResourceOption : github.com/coldsmirk/vef-framework-go/api.ResourceOption
TYPE RouterStrategy : github.com/coldsmirk/vef-framework-go/api.RouterStrategy
  METHOD CanHandle : func(kind github.com/coldsmirk/vef-framework-go/api.Kind) bool
  METHOD Name : func() string
  METHOD Route : func(handler github.com/gofiber/fiber/v3.Handler, op *github.com/coldsmirk/vef-framework-go/api.Operation)
  METHOD Setup : func(router github.com/gofiber/fiber/v3.Router) error
FUNC SignatureAuth : func() *github.com/coldsmirk/vef-framework-go/api.AuthConfig
FUNC SubscribeAuditEvent : func(bus github.com/coldsmirk/vef-framework-go/event.Bus, handler func(context.Context, *github.com/coldsmirk/vef-framework-go/api.AuditEvent) error, opts ...github.com/coldsmirk/vef-framework-go/event.SubscribeOption) (github.com/coldsmirk/vef-framework-go/event.Unsubscribe, error)
FUNC ValidateActionName : func(action string, kind github.com/coldsmirk/vef-framework-go/api.Kind) error
CONST VersionV1 : untyped string = "v1"
CONST VersionV2 : untyped string = "v2"
CONST VersionV3 : untyped string = "v3"
CONST VersionV4 : untyped string = "v4"
CONST VersionV5 : untyped string = "v5"
CONST VersionV6 : untyped string = "v6"
CONST VersionV7 : untyped string = "v7"
CONST VersionV8 : untyped string = "v8"
CONST VersionV9 : untyped string = "v9"
FUNC WithAuth : func(auth *github.com/coldsmirk/vef-framework-go/api.AuthConfig) github.com/coldsmirk/vef-framework-go/api.ResourceOption
FUNC WithOperations : func(ops ...github.com/coldsmirk/vef-framework-go/api.OperationSpec) github.com/coldsmirk/vef-framework-go/api.ResourceOption
FUNC WithVersion : func(v string) github.com/coldsmirk/vef-framework-go/api.ResourceOption

## github.com/coldsmirk/vef-framework-go/approval
CONST ActionAddAssignee : github.com/coldsmirk/vef-framework-go/approval.ActionType = "add_assignee"
CONST ActionAddCC : github.com/coldsmirk/vef-framework-go/approval.ActionType = "add_cc"
CONST ActionApprove : github.com/coldsmirk/vef-framework-go/approval.ActionType = "approve"
CONST ActionCancel : github.com/coldsmirk/vef-framework-go/approval.ActionType = "cancel"
CONST ActionExecute : github.com/coldsmirk/vef-framework-go/approval.ActionType = "execute"
CONST ActionHandle : github.com/coldsmirk/vef-framework-go/approval.ActionType = "handle"
TYPE ActionLog : github.com/coldsmirk/vef-framework-go/approval.ActionLog
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_action_log,alias:aal\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD CreationTrackedModel : github.com/coldsmirk/vef-framework-go/orm.CreationTrackedModel [field_order=3 tag=""]
  FIELD InstanceID : string [field_order=4 tag="json:\"instanceId\" bun:\"instance_id\""]
  FIELD NodeID : *string [field_order=5 tag="json:\"nodeId\" bun:\"node_id,nullzero\""]
  FIELD TaskID : *string [field_order=6 tag="json:\"taskId\" bun:\"task_id,nullzero\""]
  FIELD Action : github.com/coldsmirk/vef-framework-go/approval.ActionType [field_order=7 tag="json:\"action\" bun:\"action\""]
  FIELD OperatorID : string [field_order=8 tag="json:\"operatorId\" bun:\"operator_id\""]
  FIELD OperatorName : string [field_order=9 tag="json:\"operatorName\" bun:\"operator_name\""]
  FIELD OperatorDepartmentID : *string [field_order=10 tag="json:\"operatorDepartmentId\" bun:\"operator_department_id,nullzero\""]
  FIELD OperatorDepartmentName : *string [field_order=11 tag="json:\"operatorDepartmentName\" bun:\"operator_department_name,nullzero\""]
  FIELD IPAddress : *string [field_order=12 tag="json:\"ipAddress\" bun:\"ip_address,nullzero\""]
  FIELD UserAgent : *string [field_order=13 tag="json:\"userAgent\" bun:\"user_agent,nullzero\""]
  FIELD Opinion : *string [field_order=14 tag="json:\"opinion\" bun:\"opinion,nullzero\""]
  FIELD TransferToID : *string [field_order=15 tag="json:\"transferToId\" bun:\"transfer_to_id,nullzero\""]
  FIELD TransferToName : *string [field_order=16 tag="json:\"transferToName\" bun:\"transfer_to_name,nullzero\""]
  FIELD TransferToDepartmentID : *string [field_order=17 tag="json:\"transferToDepartmentId\" bun:\"transfer_to_department_id,nullzero\""]
  FIELD TransferToDepartmentName : *string [field_order=18 tag="json:\"transferToDepartmentName\" bun:\"transfer_to_department_name,nullzero\""]
  FIELD RollbackToNodeID : *string [field_order=19 tag="json:\"rollbackToNodeId\" bun:\"rollback_to_node_id,nullzero\""]
  FIELD AddAssigneeType : *github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType [field_order=20 tag="json:\"addAssigneeType\" bun:\"add_assignee_type,nullzero\""]
  FIELD AddedAssignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=21 tag="json:\"addedAssignees\" bun:\"added_assignees,type:jsonb\""]
  FIELD RemovedAssignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=22 tag="json:\"removedAssignees\" bun:\"removed_assignees,type:jsonb\""]
  FIELD CCUsers : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=23 tag="json:\"ccUsers\" bun:\"cc_users,type:jsonb\""]
  FIELD Attachments : []string [field_order=24 tag="json:\"attachments\" bun:\"attachments,type:jsonb,nullzero\""]
  FIELD Meta : map[string]any [field_order=25 tag="json:\"meta\" bun:\"meta,type:jsonb,nullzero\""]
  METHOD Operator : func() github.com/coldsmirk/vef-framework-go/approval.UserInfo
  METHOD TransferTo : func() *github.com/coldsmirk/vef-framework-go/approval.UserInfo
CONST ActionReassign : github.com/coldsmirk/vef-framework-go/approval.ActionType = "reassign"
CONST ActionReject : github.com/coldsmirk/vef-framework-go/approval.ActionType = "reject"
CONST ActionRemoveAssignee : github.com/coldsmirk/vef-framework-go/approval.ActionType = "remove_assignee"
CONST ActionResubmit : github.com/coldsmirk/vef-framework-go/approval.ActionType = "resubmit"
CONST ActionRollback : github.com/coldsmirk/vef-framework-go/approval.ActionType = "rollback"
CONST ActionSubmit : github.com/coldsmirk/vef-framework-go/approval.ActionType = "submit"
CONST ActionTerminate : github.com/coldsmirk/vef-framework-go/approval.ActionType = "terminate"
CONST ActionTransfer : github.com/coldsmirk/vef-framework-go/approval.ActionType = "transfer"
TYPE ActionType : github.com/coldsmirk/vef-framework-go/approval.ActionType
CONST ActionWithdraw : github.com/coldsmirk/vef-framework-go/approval.ActionType = "withdraw"
TYPE Activity : github.com/coldsmirk/vef-framework-go/approval.Activity
  FIELD Action : string [field_order=1 tag="json:\"action\""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"operator\""]
  FIELD Opinion : *string [field_order=3 tag="json:\"opinion,omitempty\""]
  FIELD Attachments : []string [field_order=4 tag="json:\"attachments,omitempty\""]
  FIELD TransferTo : *github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=5 tag="json:\"transferTo,omitempty\""]
  FIELD Target : *github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=6 tag="json:\"target,omitempty\""]
  FIELD RollbackToNodeID : *string [field_order=7 tag="json:\"rollbackToNodeId,omitempty\""]
  FIELD RollbackToNodeName : *string [field_order=8 tag="json:\"rollbackToNodeName,omitempty\""]
  FIELD AddedAssignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=9 tag="json:\"addedAssignees,omitempty\""]
  FIELD RemovedAssignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=10 tag="json:\"removedAssignees,omitempty\""]
  FIELD CCUsers : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=11 tag="json:\"ccUsers,omitempty\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=12 tag="json:\"createdAt\""]
CONST ActivityUrge : untyped string = "urge"
CONST AddAssigneeAfter : github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType = "after"
CONST AddAssigneeBefore : github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType = "before"
CONST AddAssigneeParallel : github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType = "parallel"
TYPE AddAssigneeType : github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType
  METHOD IsValid : func() bool
  METHOD UnmarshalJSON : func(data []byte) error
CONST AggregateAvg : github.com/coldsmirk/vef-framework-go/approval.AggregateKind = "avg"
CONST AggregateCount : github.com/coldsmirk/vef-framework-go/approval.AggregateKind = "count"
TYPE AggregateKind : github.com/coldsmirk/vef-framework-go/approval.AggregateKind
  METHOD FoldsColumn : func() bool
  METHOD IsValid : func() bool
CONST AggregateSum : github.com/coldsmirk/vef-framework-go/approval.AggregateKind = "sum"
TYPE Aggregator : github.com/coldsmirk/vef-framework-go/approval.Aggregator
  METHOD Fold : func(values []float64, rowCount int) (result float64, matchable bool)
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.AggregateKind
FUNC AllEventTypes : func() []string
TYPE ApprovalMethod : github.com/coldsmirk/vef-framework-go/approval.ApprovalMethod
  METHOD IsValid : func() bool
TYPE ApprovalNodeData : github.com/coldsmirk/vef-framework-go/approval.ApprovalNodeData
  FIELD BaseNodeData : github.com/coldsmirk/vef-framework-go/approval.BaseNodeData [field_order=1 tag=""]
  FIELD TaskNodeData : github.com/coldsmirk/vef-framework-go/approval.TaskNodeData [field_order=2 tag=""]
  FIELD ApprovalMethod : github.com/coldsmirk/vef-framework-go/approval.ApprovalMethod [field_order=3 tag="json:\"approvalMethod,omitempty\""]
  FIELD PassRule : github.com/coldsmirk/vef-framework-go/approval.PassRule [field_order=4 tag="json:\"passRule,omitempty\""]
  FIELD PassRatio : github.com/coldsmirk/vef-framework-go/decimal.Decimal [field_order=5 tag="json:\"passRatio\""]
  FIELD SameApplicantAction : github.com/coldsmirk/vef-framework-go/approval.SameApplicantAction [field_order=6 tag="json:\"sameApplicantAction,omitempty\""]
  FIELD ConsecutiveApproverAction : github.com/coldsmirk/vef-framework-go/approval.ConsecutiveApproverAction [field_order=7 tag="json:\"consecutiveApproverAction,omitempty\""]
  FIELD RollbackType : github.com/coldsmirk/vef-framework-go/approval.RollbackType [field_order=8 tag="json:\"rollbackType,omitempty\""]
  FIELD RollbackDataStrategy : github.com/coldsmirk/vef-framework-go/approval.RollbackDataStrategy [field_order=9 tag="json:\"rollbackDataStrategy,omitempty\""]
  FIELD RollbackTargetKeys : []string [field_order=10 tag="json:\"rollbackTargetKeys,omitempty\""]
  FIELD IsRollbackAllowed : *bool [field_order=11 tag="json:\"isRollbackAllowed,omitempty\""]
  FIELD IsAddAssigneeAllowed : *bool [field_order=12 tag="json:\"isAddAssigneeAllowed,omitempty\""]
  FIELD AddAssigneeTypes : []github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType [field_order=13 tag="json:\"addAssigneeTypes,omitempty\""]
  FIELD IsRemoveAssigneeAllowed : *bool [field_order=14 tag="json:\"isRemoveAssigneeAllowed,omitempty\""]
  FIELD IsManualCCAllowed : *bool [field_order=15 tag="json:\"isManualCcAllowed,omitempty\""]
  FIELD AdminUserIDs : []string [promoted_from=TaskNodeData depth=1 field_order=5 tag="json:\"adminUserIds,omitempty\""]
  FIELD Assignees : []github.com/coldsmirk/vef-framework-go/approval.AssigneeDefinition [promoted_from=TaskNodeData depth=1 field_order=1 tag="json:\"assignees,omitempty\""]
  FIELD CCs : []github.com/coldsmirk/vef-framework-go/approval.CCDefinition [promoted_from=TaskNodeData depth=1 field_order=12 tag="json:\"ccs,omitempty\""]
  FIELD Description : *string [promoted_from=BaseNodeData depth=1 field_order=2 tag="json:\"description,omitempty\""]
  FIELD EmptyAssigneeAction : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction [promoted_from=TaskNodeData depth=1 field_order=3 tag="json:\"emptyAssigneeAction,omitempty\""]
  FIELD ExecutionType : github.com/coldsmirk/vef-framework-go/approval.ExecutionType [promoted_from=TaskNodeData depth=1 field_order=2 tag="json:\"executionType,omitempty\""]
  FIELD FallbackUserIDs : []string [promoted_from=TaskNodeData depth=1 field_order=4 tag="json:\"fallbackUserIds,omitempty\""]
  FIELD FieldPermissions : map[string]github.com/coldsmirk/vef-framework-go/approval.Permission [promoted_from=TaskNodeData depth=1 field_order=13 tag="json:\"fieldPermissions,omitempty\""]
  FIELD IsOpinionRequired : bool [promoted_from=TaskNodeData depth=1 field_order=7 tag="json:\"isOpinionRequired,omitempty\""]
  FIELD IsTransferAllowed : *bool [promoted_from=TaskNodeData depth=1 field_order=6 tag="json:\"isTransferAllowed,omitempty\""]
  FIELD Name : string [promoted_from=BaseNodeData depth=1 field_order=1 tag="json:\"name,omitempty\""]
  FIELD TimeoutAction : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction [promoted_from=TaskNodeData depth=1 field_order=9 tag="json:\"timeoutAction,omitempty\""]
  FIELD TimeoutHours : int [promoted_from=TaskNodeData depth=1 field_order=8 tag="json:\"timeoutHours,omitempty\""]
  FIELD TimeoutNotifyBeforeHours : int [promoted_from=TaskNodeData depth=1 field_order=10 tag="json:\"timeoutNotifyBeforeHours,omitempty\""]
  FIELD UrgeCooldownMinutes : int [promoted_from=TaskNodeData depth=1 field_order=11 tag="json:\"urgeCooldownMinutes,omitempty\""]
  METHOD ApplyTo : func(node *github.com/coldsmirk/vef-framework-go/approval.FlowNode)
  METHOD GetAssignees : func() []github.com/coldsmirk/vef-framework-go/approval.AssigneeDefinition
  METHOD GetCCs : func() []github.com/coldsmirk/vef-framework-go/approval.CCDefinition
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.NodeKind
CONST ApprovalParallel : github.com/coldsmirk/vef-framework-go/approval.ApprovalMethod = "parallel"
CONST ApprovalSequential : github.com/coldsmirk/vef-framework-go/approval.ApprovalMethod = "sequential"
TYPE AssigneeDefinition : github.com/coldsmirk/vef-framework-go/approval.AssigneeDefinition
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind [field_order=1 tag="json:\"kind\""]
  FIELD IDs : []string [field_order=2 tag="json:\"ids,omitempty\""]
  FIELD FormField : *string [field_order=3 tag="json:\"formField,omitempty\""]
  FIELD SortOrder : int [field_order=4 tag="json:\"sortOrder\""]
CONST AssigneeDepartment : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind = "department"
CONST AssigneeDepartmentLeader : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind = "department_leader"
CONST AssigneeFormField : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind = "form_field"
TYPE AssigneeKind : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind
  METHOD IsValid : func() bool
CONST AssigneeRole : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind = "role"
CONST AssigneeSelf : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind = "self"
TYPE AssigneeService : github.com/coldsmirk/vef-framework-go/approval.AssigneeService
  METHOD GetDepartmentLeaders : func(ctx context.Context, departmentID string) ([]github.com/coldsmirk/vef-framework-go/approval.UserInfo, error)
  METHOD GetRoleUsers : func(ctx context.Context, roleID string) ([]github.com/coldsmirk/vef-framework-go/approval.UserInfo, error)
  METHOD GetSuperior : func(ctx context.Context, userID string) (*github.com/coldsmirk/vef-framework-go/approval.UserInfo, error)
CONST AssigneeSuperior : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind = "superior"
CONST AssigneeUser : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind = "user"
TYPE AssigneesAddedEvent : github.com/coldsmirk/vef-framework-go/approval.AssigneesAddedEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD AddType : github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType [field_order=2 tag="json:\"addType\""]
  FIELD Assignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=3 tag="json:\"assignees\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE AssigneesRemovedEvent : github.com/coldsmirk/vef-framework-go/approval.AssigneesRemovedEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Assignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"assignees\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE BaseNodeData : github.com/coldsmirk/vef-framework-go/approval.BaseNodeData
  FIELD Name : string [field_order=1 tag="json:\"name,omitempty\""]
  FIELD Description : *string [field_order=2 tag="json:\"description,omitempty\""]
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
CONST BindingBusiness : github.com/coldsmirk/vef-framework-go/approval.BindingMode = "business"
TYPE BindingMode : github.com/coldsmirk/vef-framework-go/approval.BindingMode
  METHOD IsValid : func() bool
CONST BindingStandalone : github.com/coldsmirk/vef-framework-go/approval.BindingMode = "standalone"
TYPE BindingTrigger : github.com/coldsmirk/vef-framework-go/approval.BindingTrigger
CONST BindingTriggerCompleted : github.com/coldsmirk/vef-framework-go/approval.BindingTrigger = "completed"
CONST BindingTriggerResubmitted : github.com/coldsmirk/vef-framework-go/approval.BindingTrigger = "resubmitted"
CONST BindingTriggerReturned : github.com/coldsmirk/vef-framework-go/approval.BindingTrigger = "returned"
CONST BindingTriggerStarted : github.com/coldsmirk/vef-framework-go/approval.BindingTrigger = "started"
CONST BindingTriggerWithdrawn : github.com/coldsmirk/vef-framework-go/approval.BindingTrigger = "withdrawn"
TYPE BusinessRefProvider : github.com/coldsmirk/vef-framework-go/approval.BusinessRefProvider
  METHOD OnInstanceCreated : func(ctx context.Context, db github.com/coldsmirk/vef-framework-go/orm.DB, flow *github.com/coldsmirk/vef-framework-go/approval.Flow, instance *github.com/coldsmirk/vef-framework-go/approval.Instance) (businessRef string, err error)
TYPE BusinessRefResolver : github.com/coldsmirk/vef-framework-go/approval.BusinessRefResolver
  METHOD ResolveRecordID : func(ctx context.Context, flow *github.com/coldsmirk/vef-framework-go/approval.Flow, businessRef string) (string, error)
TYPE CCDefinition : github.com/coldsmirk/vef-framework-go/approval.CCDefinition
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.CCKind [field_order=1 tag="json:\"kind\""]
  FIELD IDs : []string [field_order=2 tag="json:\"ids,omitempty\""]
  FIELD FormField : *string [field_order=3 tag="json:\"formField,omitempty\""]
  FIELD Timing : github.com/coldsmirk/vef-framework-go/approval.CCTiming [field_order=4 tag="json:\"timing,omitempty\""]
CONST CCDepartment : github.com/coldsmirk/vef-framework-go/approval.CCKind = "department"
CONST CCFormField : github.com/coldsmirk/vef-framework-go/approval.CCKind = "form_field"
TYPE CCKind : github.com/coldsmirk/vef-framework-go/approval.CCKind
  METHOD IsValid : func() bool
TYPE CCNodeData : github.com/coldsmirk/vef-framework-go/approval.CCNodeData
  FIELD BaseNodeData : github.com/coldsmirk/vef-framework-go/approval.BaseNodeData [field_order=1 tag=""]
  FIELD CCs : []github.com/coldsmirk/vef-framework-go/approval.CCDefinition [field_order=2 tag="json:\"ccs,omitempty\""]
  FIELD IsReadConfirmRequired : bool [field_order=3 tag="json:\"isReadConfirmRequired,omitempty\""]
  FIELD FieldPermissions : map[string]github.com/coldsmirk/vef-framework-go/approval.Permission [field_order=4 tag="json:\"fieldPermissions,omitempty\""]
  FIELD Description : *string [promoted_from=BaseNodeData depth=1 field_order=2 tag="json:\"description,omitempty\""]
  FIELD Name : string [promoted_from=BaseNodeData depth=1 field_order=1 tag="json:\"name,omitempty\""]
  METHOD ApplyTo : func(node *github.com/coldsmirk/vef-framework-go/approval.FlowNode)
  METHOD GetCCs : func() []github.com/coldsmirk/vef-framework-go/approval.CCDefinition
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.NodeKind
TYPE CCNotifiedEvent : github.com/coldsmirk/vef-framework-go/approval.CCNotifiedEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD NodeID : string [field_order=2 tag="json:\"nodeId\""]
  FIELD NodeName : string [field_order=3 tag="json:\"nodeName\""]
  FIELD Recipients : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=4 tag="json:\"recipients\""]
  FIELD IsManual : bool [field_order=5 tag="json:\"isManual\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE CCRecipient : github.com/coldsmirk/vef-framework-go/approval.CCRecipient
  FIELD User : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=1 tag="json:\"user\""]
  FIELD ReadAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=2 tag="json:\"readAt,omitempty\""]
TYPE CCRecord : github.com/coldsmirk/vef-framework-go/approval.CCRecord
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_cc_record,alias:acr\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD CreationTrackedModel : github.com/coldsmirk/vef-framework-go/orm.CreationTrackedModel [field_order=3 tag=""]
  FIELD InstanceID : string [field_order=4 tag="json:\"instanceId\" bun:\"instance_id\""]
  FIELD NodeID : *string [field_order=5 tag="json:\"nodeId\" bun:\"node_id,nullzero\""]
  FIELD VisitID : *string [field_order=6 tag="json:\"visitId\" bun:\"visit_id,nullzero\""]
  FIELD TaskID : *string [field_order=7 tag="json:\"taskId\" bun:\"task_id,nullzero\""]
  FIELD CCUserID : string [field_order=8 tag="json:\"ccUserId\" bun:\"cc_user_id\""]
  FIELD CCUserName : string [field_order=9 tag="json:\"ccUserName\" bun:\"cc_user_name\""]
  FIELD CCUserDepartmentID : *string [field_order=10 tag="json:\"ccUserDepartmentId\" bun:\"cc_user_department_id,nullzero\""]
  FIELD CCUserDepartmentName : *string [field_order=11 tag="json:\"ccUserDepartmentName\" bun:\"cc_user_department_name,nullzero\""]
  FIELD IsManual : bool [field_order=12 tag="json:\"isManual\" bun:\"is_manual\""]
  FIELD ReadAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=13 tag="json:\"readAt\" bun:\"read_at,nullzero\""]
  METHOD Recipient : func() github.com/coldsmirk/vef-framework-go/approval.CCRecipient
CONST CCRole : github.com/coldsmirk/vef-framework-go/approval.CCKind = "role"
TYPE CCTiming : github.com/coldsmirk/vef-framework-go/approval.CCTiming
  METHOD IsValid : func() bool
CONST CCTimingAlways : github.com/coldsmirk/vef-framework-go/approval.CCTiming = "always"
CONST CCTimingOnApprove : github.com/coldsmirk/vef-framework-go/approval.CCTiming = "on_approve"
CONST CCTimingOnReject : github.com/coldsmirk/vef-framework-go/approval.CCTiming = "on_reject"
CONST CCUser : github.com/coldsmirk/vef-framework-go/approval.CCKind = "user"
TYPE CallerContext : github.com/coldsmirk/vef-framework-go/approval.CallerContext
  FIELD TenantID : string [field_order=1 tag=""]
  FIELD IsSuperAdmin : bool [field_order=2 tag=""]
  FIELD IsSystemInternal : bool [field_order=3 tag=""]
  METHOD Allows : func(entityTenantID string) bool
  METHOD Authorize : func(entityTenantID string) error
  METHOD ResolveWriteTenant : func(clientTenant string) (string, error)
  METHOD TenantScopeFilter : func(override string) (*string, error)
CONST ColumnBoolean : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "boolean"
TYPE ColumnDataType : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType
  METHOD IsValid : func() bool
CONST ColumnDate : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "date"
CONST ColumnDatetime : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "datetime"
CONST ColumnDecimal : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "decimal"
CONST ColumnInteger : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "integer"
CONST ColumnJSON : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "json"
CONST ColumnString : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "string"
CONST ColumnText : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType = "text"
TYPE Condition : github.com/coldsmirk/vef-framework-go/approval.Condition
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.ConditionKind [field_order=1 tag="json:\"kind\""]
  FIELD Subject : string [field_order=2 tag="json:\"subject\""]
  FIELD Aggregate : github.com/coldsmirk/vef-framework-go/approval.AggregateKind [field_order=3 tag="json:\"aggregate,omitempty\""]
  FIELD Column : string [field_order=4 tag="json:\"column,omitempty\""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator [field_order=5 tag="json:\"operator\""]
  FIELD Value : any [field_order=6 tag="json:\"value\""]
  FIELD Expression : string [field_order=7 tag="json:\"expression\""]
TYPE ConditionBranch : github.com/coldsmirk/vef-framework-go/approval.ConditionBranch
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD Label : string [field_order=2 tag="json:\"label\""]
  FIELD ConditionGroups : []github.com/coldsmirk/vef-framework-go/approval.ConditionGroup [field_order=3 tag="json:\"conditionGroups,omitempty\""]
  FIELD IsDefault : bool [field_order=4 tag="json:\"isDefault,omitempty\""]
  FIELD Priority : int [field_order=5 tag="json:\"priority\""]
TYPE ConditionEvaluator : github.com/coldsmirk/vef-framework-go/approval.ConditionEvaluator
  METHOD Evaluate : func(ctx context.Context, cond github.com/coldsmirk/vef-framework-go/approval.Condition, ec *github.com/coldsmirk/vef-framework-go/approval.EvaluationContext) (bool, error)
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.ConditionKind
CONST ConditionExpression : github.com/coldsmirk/vef-framework-go/approval.ConditionKind = "expression"
CONST ConditionField : github.com/coldsmirk/vef-framework-go/approval.ConditionKind = "field"
TYPE ConditionGroup : github.com/coldsmirk/vef-framework-go/approval.ConditionGroup
  FIELD Conditions : []github.com/coldsmirk/vef-framework-go/approval.Condition [field_order=1 tag="json:\"conditions\""]
TYPE ConditionKind : github.com/coldsmirk/vef-framework-go/approval.ConditionKind
TYPE ConditionNodeData : github.com/coldsmirk/vef-framework-go/approval.ConditionNodeData
  FIELD BaseNodeData : github.com/coldsmirk/vef-framework-go/approval.BaseNodeData [field_order=1 tag=""]
  FIELD Branches : []github.com/coldsmirk/vef-framework-go/approval.ConditionBranch [field_order=2 tag="json:\"branches,omitempty\""]
  FIELD Description : *string [promoted_from=BaseNodeData depth=1 field_order=2 tag="json:\"description,omitempty\""]
  FIELD Name : string [promoted_from=BaseNodeData depth=1 field_order=1 tag="json:\"name,omitempty\""]
  METHOD ApplyTo : func(node *github.com/coldsmirk/vef-framework-go/approval.FlowNode)
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.NodeKind
TYPE ConditionOperator : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator
  METHOD IsValid : func() bool
TYPE ConsecutiveApproverAction : github.com/coldsmirk/vef-framework-go/approval.ConsecutiveApproverAction
  METHOD IsValid : func() bool
CONST ConsecutiveApproverAutoPass : github.com/coldsmirk/vef-framework-go/approval.ConsecutiveApproverAction = "auto_pass"
CONST ConsecutiveApproverNone : github.com/coldsmirk/vef-framework-go/approval.ConsecutiveApproverAction = "none"
CONST DefaultApprovalMethod : github.com/coldsmirk/vef-framework-go/approval.ApprovalMethod = "parallel"
CONST DefaultCCTiming : github.com/coldsmirk/vef-framework-go/approval.CCTiming = "always"
CONST DefaultConsecutiveApproverAction : github.com/coldsmirk/vef-framework-go/approval.ConsecutiveApproverAction = "none"
CONST DefaultEmptyAssigneeAction : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction = "auto_pass"
CONST DefaultExecutionType : github.com/coldsmirk/vef-framework-go/approval.ExecutionType = "manual"
CONST DefaultHandleApprovalMethod : github.com/coldsmirk/vef-framework-go/approval.ApprovalMethod = "sequential"
CONST DefaultHandlePassRule : github.com/coldsmirk/vef-framework-go/approval.PassRule = "any"
CONST DefaultPassRule : github.com/coldsmirk/vef-framework-go/approval.PassRule = "all"
CONST DefaultRollbackDataStrategy : github.com/coldsmirk/vef-framework-go/approval.RollbackDataStrategy = "keep"
CONST DefaultRollbackType : github.com/coldsmirk/vef-framework-go/approval.RollbackType = "previous"
CONST DefaultSameApplicantAction : github.com/coldsmirk/vef-framework-go/approval.SameApplicantAction = "self_approve"
CONST DefaultTenantID : untyped string = "default"
CONST DefaultTimeoutAction : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction = "none"
CONST DefaultUrgeCooldownMinutes : untyped int = 30
TYPE Delegation : github.com/coldsmirk/vef-framework-go/approval.Delegation
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_delegation,alias:ad\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD DelegatorID : string [field_order=3 tag="json:\"delegatorId\" bun:\"delegator_id\""]
  FIELD DelegateeID : string [field_order=4 tag="json:\"delegateeId\" bun:\"delegatee_id\""]
  FIELD FlowCategoryID : *string [field_order=5 tag="json:\"flowCategoryId\" bun:\"flow_category_id,nullzero\""]
  FIELD FlowID : *string [field_order=6 tag="json:\"flowId\" bun:\"flow_id,nullzero\""]
  FIELD StartTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=7 tag="json:\"startTime\" bun:\"start_time\""]
  FIELD EndTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=8 tag="json:\"endTime\" bun:\"end_time\""]
  FIELD IsActive : bool [field_order=9 tag="json:\"isActive\" bun:\"is_active\""]
  FIELD Reason : *string [field_order=10 tag="json:\"reason\" bun:\"reason,nullzero\""]
TYPE DomainEvent : github.com/coldsmirk/vef-framework-go/approval.DomainEvent
  METHOD EventType : func() string
TYPE EdgeDefinition : github.com/coldsmirk/vef-framework-go/approval.EdgeDefinition
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD Source : string [field_order=2 tag="json:\"source\""]
  FIELD Target : string [field_order=3 tag="json:\"target\""]
  FIELD SourceHandle : *string [field_order=4 tag="json:\"sourceHandle,omitempty\""]
  FIELD Data : map[string]any [field_order=5 tag="json:\"data,omitempty\""]
TYPE EmptyAssigneeAction : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction
  METHOD IsValid : func() bool
CONST EmptyAssigneeAutoPass : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction = "auto_pass"
CONST EmptyAssigneeTransferAdmin : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction = "transfer_admin"
CONST EmptyAssigneeTransferApplicant : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction = "transfer_applicant"
CONST EmptyAssigneeTransferSpecified : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction = "transfer_specified"
CONST EmptyAssigneeTransferSuperior : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction = "transfer_superior"
TYPE EndNodeData : github.com/coldsmirk/vef-framework-go/approval.EndNodeData
  FIELD BaseNodeData : github.com/coldsmirk/vef-framework-go/approval.BaseNodeData [field_order=1 tag=""]
  FIELD Description : *string [promoted_from=BaseNodeData depth=1 field_order=2 tag="json:\"description,omitempty\""]
  FIELD Name : string [promoted_from=BaseNodeData depth=1 field_order=1 tag="json:\"name,omitempty\""]
  METHOD ApplyTo : func(node *github.com/coldsmirk/vef-framework-go/approval.FlowNode)
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.NodeKind
VAR ErrAnonymousSubscriberGroup : error
VAR ErrCrossTenantAccess : error
VAR ErrDerivedGroupConflict : error
VAR ErrInvalidBusinessIdentifier : error
VAR ErrNodeDataUnmarshal : error
VAR ErrUnknownNodeKind : error
TYPE EvaluationContext : github.com/coldsmirk/vef-framework-go/approval.EvaluationContext
  FIELD FormData : github.com/coldsmirk/vef-framework-go/approval.FormData [field_order=1 tag=""]
  FIELD ApplicantID : string [field_order=2 tag=""]
  FIELD ApplicantDepartmentID : *string [field_order=3 tag=""]
  FIELD Globals : map[string]any [field_order=4 tag=""]
CONST EventTypeAssigneesAdded : untyped string = "approval.task.assignees_added"
CONST EventTypeAssigneesRemoved : untyped string = "approval.task.assignees_removed"
CONST EventTypeCCNotified : untyped string = "approval.cc.notified"
CONST EventTypeFlowCreated : untyped string = "approval.flow.created"
CONST EventTypeFlowDeployed : untyped string = "approval.flow.deployed"
CONST EventTypeFlowPublished : untyped string = "approval.flow.published"
CONST EventTypeFlowToggled : untyped string = "approval.flow.toggled"
CONST EventTypeFlowUpdated : untyped string = "approval.flow.updated"
CONST EventTypeInstanceBindingFailed : untyped string = "approval.instance.binding_failed"
CONST EventTypeInstanceCompleted : untyped string = "approval.instance.completed"
CONST EventTypeInstanceCreated : untyped string = "approval.instance.created"
CONST EventTypeInstanceResubmitted : untyped string = "approval.instance.resubmitted"
CONST EventTypeInstanceReturned : untyped string = "approval.instance.returned"
CONST EventTypeInstanceRolledBack : untyped string = "approval.instance.rolled_back"
CONST EventTypeInstanceWithdrawn : untyped string = "approval.instance.withdrawn"
CONST EventTypeNodeAutoPassed : untyped string = "approval.node.auto_passed"
CONST EventTypeTaskApproved : untyped string = "approval.task.approved"
CONST EventTypeTaskCanceled : untyped string = "approval.task.canceled"
CONST EventTypeTaskCreated : untyped string = "approval.task.created"
CONST EventTypeTaskDeadlineWarning : untyped string = "approval.task.deadline_warning"
CONST EventTypeTaskHandled : untyped string = "approval.task.handled"
CONST EventTypeTaskReassigned : untyped string = "approval.task.reassigned"
CONST EventTypeTaskRejected : untyped string = "approval.task.rejected"
CONST EventTypeTaskTimedOut : untyped string = "approval.task.timed_out"
CONST EventTypeTaskTransferred : untyped string = "approval.task.transferred"
CONST EventTypeTaskUrged : untyped string = "approval.task.urged"
CONST ExecutionAutoPass : github.com/coldsmirk/vef-framework-go/approval.ExecutionType = "auto_pass"
CONST ExecutionAutoReject : github.com/coldsmirk/vef-framework-go/approval.ExecutionType = "auto_reject"
CONST ExecutionManual : github.com/coldsmirk/vef-framework-go/approval.ExecutionType = "manual"
TYPE ExecutionType : github.com/coldsmirk/vef-framework-go/approval.ExecutionType
  METHOD IsValid : func() bool
CONST FieldDate : github.com/coldsmirk/vef-framework-go/approval.FieldKind = "date"
CONST FieldInput : github.com/coldsmirk/vef-framework-go/approval.FieldKind = "input"
TYPE FieldKind : github.com/coldsmirk/vef-framework-go/approval.FieldKind
  METHOD IsValid : func() bool
CONST FieldNumber : github.com/coldsmirk/vef-framework-go/approval.FieldKind = "number"
TYPE FieldOption : github.com/coldsmirk/vef-framework-go/approval.FieldOption
  FIELD Label : string [field_order=1 tag="json:\"label\""]
  FIELD Value : any [field_order=2 tag="json:\"value\""]
CONST FieldSelect : github.com/coldsmirk/vef-framework-go/approval.FieldKind = "select"
CONST FieldTable : github.com/coldsmirk/vef-framework-go/approval.FieldKind = "table"
CONST FieldTextarea : github.com/coldsmirk/vef-framework-go/approval.FieldKind = "textarea"
CONST FieldUpload : github.com/coldsmirk/vef-framework-go/approval.FieldKind = "upload"
TYPE Flow : github.com/coldsmirk/vef-framework-go/approval.Flow
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow,alias:af\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD TenantID : string [field_order=3 tag="json:\"tenantId\" bun:\"tenant_id\""]
  FIELD CategoryID : string [field_order=4 tag="json:\"categoryId\" bun:\"category_id\""]
  FIELD Code : string [field_order=5 tag="json:\"code\" bun:\"code\""]
  FIELD Name : string [field_order=6 tag="json:\"name\" bun:\"name\""]
  FIELD Icon : *string [field_order=7 tag="json:\"icon\" bun:\"icon,nullzero\""]
  FIELD Description : *string [field_order=8 tag="json:\"description\" bun:\"description,nullzero\""]
  FIELD BindingMode : github.com/coldsmirk/vef-framework-go/approval.BindingMode [field_order=9 tag="json:\"bindingMode\" bun:\"binding_mode\""]
  FIELD BusinessTable : *string [field_order=10 tag="json:\"businessTable\" bun:\"business_table,nullzero\""]
  FIELD BusinessPKField : *string [field_order=11 tag="json:\"businessPkField\" bun:\"business_pk_field,nullzero\""]
  FIELD BusinessStatusField : *string [field_order=12 tag="json:\"businessStatusField\" bun:\"business_status_field,nullzero\""]
  FIELD BusinessInstanceIDField : *string [field_order=13 tag="json:\"businessInstanceIdField\" bun:\"business_instance_id_field,nullzero\""]
  FIELD BusinessStartedAtField : *string [field_order=14 tag="json:\"businessStartedAtField\" bun:\"business_started_at_field,nullzero\""]
  FIELD BusinessFinishedAtField : *string [field_order=15 tag="json:\"businessFinishedAtField\" bun:\"business_finished_at_field,nullzero\""]
  FIELD AdminUserIDs : []string [field_order=16 tag="json:\"adminUserIds\" bun:\"admin_user_ids,type:jsonb\""]
  FIELD IsAllInitiationAllowed : bool [field_order=17 tag="json:\"isAllInitiationAllowed\" bun:\"is_all_initiation_allowed\""]
  FIELD InstanceTitleTemplate : string [field_order=18 tag="json:\"instanceTitleTemplate\" bun:\"instance_title_template\""]
  FIELD IsActive : bool [field_order=19 tag="json:\"isActive\" bun:\"is_active\""]
  FIELD CurrentVersion : int [field_order=20 tag="json:\"currentVersion\" bun:\"current_version\""]
TYPE FlowCategory : github.com/coldsmirk/vef-framework-go/approval.FlowCategory
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow_category,alias:afc\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD TenantID : string [field_order=3 tag="json:\"tenantId\" bun:\"tenant_id\""]
  FIELD Code : string [field_order=4 tag="json:\"code\" bun:\"code\""]
  FIELD Name : string [field_order=5 tag="json:\"name\" bun:\"name\""]
  FIELD Icon : *string [field_order=6 tag="json:\"icon\" bun:\"icon,nullzero\""]
  FIELD ParentID : *string [field_order=7 tag="json:\"parentId\" bun:\"parent_id,nullzero\""]
  FIELD SortOrder : int [field_order=8 tag="json:\"sortOrder\" bun:\"sort_order\""]
  FIELD IsActive : bool [field_order=9 tag="json:\"isActive\" bun:\"is_active\""]
  FIELD Remark : *string [field_order=10 tag="json:\"remark\" bun:\"remark,nullzero\""]
  FIELD Children : []github.com/coldsmirk/vef-framework-go/approval.FlowCategory [field_order=11 tag="json:\"children,omitempty\" bun:\"-\""]
TYPE FlowCreatedEvent : github.com/coldsmirk/vef-framework-go/approval.FlowCreatedEvent
  FIELD FlowEventBase : github.com/coldsmirk/vef-framework-go/approval.FlowEventBase [field_order=1 tag=""]
  FIELD CategoryID : string [field_order=2 tag="json:\"categoryId\""]
  FIELD Code : string [promoted_from=FlowEventBase depth=1 field_order=3 tag="json:\"code\""]
  FIELD FlowID : string [promoted_from=FlowEventBase depth=1 field_order=1 tag="json:\"flowId\""]
  FIELD Name : string [promoted_from=FlowEventBase depth=1 field_order=4 tag="json:\"name\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=FlowEventBase depth=1 field_order=5 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=FlowEventBase depth=1 field_order=2 tag="json:\"tenantId\""]
  METHOD EventType : func() string
TYPE FlowDefinition : github.com/coldsmirk/vef-framework-go/approval.FlowDefinition
  FIELD Nodes : []github.com/coldsmirk/vef-framework-go/approval.NodeDefinition [field_order=1 tag="json:\"nodes\""]
  FIELD Edges : []github.com/coldsmirk/vef-framework-go/approval.EdgeDefinition [field_order=2 tag="json:\"edges\""]
TYPE FlowDeployedEvent : github.com/coldsmirk/vef-framework-go/approval.FlowDeployedEvent
  FIELD FlowEventBase : github.com/coldsmirk/vef-framework-go/approval.FlowEventBase [field_order=1 tag=""]
  FIELD VersionID : string [field_order=2 tag="json:\"versionId\""]
  FIELD Version : int [field_order=3 tag="json:\"version\""]
  FIELD Code : string [promoted_from=FlowEventBase depth=1 field_order=3 tag="json:\"code\""]
  FIELD FlowID : string [promoted_from=FlowEventBase depth=1 field_order=1 tag="json:\"flowId\""]
  FIELD Name : string [promoted_from=FlowEventBase depth=1 field_order=4 tag="json:\"name\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=FlowEventBase depth=1 field_order=5 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=FlowEventBase depth=1 field_order=2 tag="json:\"tenantId\""]
  METHOD EventType : func() string
TYPE FlowEdge : github.com/coldsmirk/vef-framework-go/approval.FlowEdge
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow_edge,alias:afe\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD FlowVersionID : string [field_order=3 tag="json:\"flowVersionId\" bun:\"flow_version_id\""]
  FIELD Key : string [field_order=4 tag="json:\"key\" bun:\"key,nullzero\""]
  FIELD SourceNodeID : string [field_order=5 tag="json:\"sourceNodeId\" bun:\"source_node_id\""]
  FIELD SourceNodeKey : string [field_order=6 tag="json:\"sourceNodeKey\" bun:\"source_node_key\""]
  FIELD TargetNodeID : string [field_order=7 tag="json:\"targetNodeId\" bun:\"target_node_id\""]
  FIELD TargetNodeKey : string [field_order=8 tag="json:\"targetNodeKey\" bun:\"target_node_key\""]
  FIELD SourceHandle : *string [field_order=9 tag="json:\"sourceHandle\" bun:\"source_handle,nullzero\""]
TYPE FlowEventBase : github.com/coldsmirk/vef-framework-go/approval.FlowEventBase
  FIELD FlowID : string [field_order=1 tag="json:\"flowId\""]
  FIELD TenantID : string [field_order=2 tag="json:\"tenantId\""]
  FIELD Code : string [field_order=3 tag="json:\"code\""]
  FIELD Name : string [field_order=4 tag="json:\"name\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=5 tag="json:\"occurredTime\""]
TYPE FlowGraphEdge : github.com/coldsmirk/vef-framework-go/approval.FlowGraphEdge
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD Source : string [field_order=2 tag="json:\"source\""]
  FIELD Target : string [field_order=3 tag="json:\"target\""]
  FIELD SourceHandle : *string [field_order=4 tag="json:\"sourceHandle,omitempty\""]
TYPE FlowGraphNode : github.com/coldsmirk/vef-framework-go/approval.FlowGraphNode
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD NodeID : string [field_order=2 tag="json:\"nodeId\""]
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.NodeKind [field_order=3 tag="json:\"kind\""]
  FIELD Position : github.com/coldsmirk/vef-framework-go/approval.Position [field_order=4 tag="json:\"position\""]
  FIELD Data : github.com/coldsmirk/vef-framework-go/approval.FlowGraphNodeData [field_order=5 tag="json:\"data\""]
TYPE FlowGraphNodeData : github.com/coldsmirk/vef-framework-go/approval.FlowGraphNodeData
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus [field_order=2 tag="json:\"status\""]
  FIELD ExecutionType : string [field_order=3 tag="json:\"executionType,omitempty\""]
  FIELD ApprovalMethod : string [field_order=4 tag="json:\"approvalMethod,omitempty\""]
  FIELD PassRule : string [field_order=5 tag="json:\"passRule,omitempty\""]
  FIELD PassRatio : *github.com/coldsmirk/vef-framework-go/decimal.Decimal [field_order=6 tag="json:\"passRatio,omitempty\""]
  FIELD Participants : []github.com/coldsmirk/vef-framework-go/approval.NodeParticipant [field_order=7 tag="json:\"participants,omitempty\""]
  FIELD CCRecipients : []github.com/coldsmirk/vef-framework-go/approval.CCRecipient [field_order=8 tag="json:\"ccRecipients,omitempty\""]
  FIELD Activities : []github.com/coldsmirk/vef-framework-go/approval.Activity [field_order=9 tag="json:\"activities,omitempty\""]
  FIELD StartedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=10 tag="json:\"startedAt,omitempty\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=11 tag="json:\"finishedAt,omitempty\""]
TYPE FlowInitiator : github.com/coldsmirk/vef-framework-go/approval.FlowInitiator
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow_initiator,alias:afi\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD FlowID : string [field_order=3 tag="json:\"flowId\" bun:\"flow_id\""]
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.InitiatorKind [field_order=4 tag="json:\"kind\" bun:\"kind\""]
  FIELD IDs : []string [field_order=5 tag="json:\"ids\" bun:\"ids,type:jsonb\""]
TYPE FlowNode : github.com/coldsmirk/vef-framework-go/approval.FlowNode
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow_node,alias:afn\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD FlowVersionID : string [field_order=3 tag="json:\"flowVersionId\" bun:\"flow_version_id\""]
  FIELD Key : string [field_order=4 tag="json:\"key\" bun:\"key\""]
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.NodeKind [field_order=5 tag="json:\"kind\" bun:\"kind\""]
  FIELD Name : string [field_order=6 tag="json:\"name\" bun:\"name\""]
  FIELD Description : *string [field_order=7 tag="json:\"description\" bun:\"description,nullzero\""]
  FIELD ExecutionType : github.com/coldsmirk/vef-framework-go/approval.ExecutionType [field_order=8 tag="json:\"executionType\" bun:\"execution_type\""]
  FIELD ApprovalMethod : github.com/coldsmirk/vef-framework-go/approval.ApprovalMethod [field_order=9 tag="json:\"approvalMethod\" bun:\"approval_method\""]
  FIELD PassRule : github.com/coldsmirk/vef-framework-go/approval.PassRule [field_order=10 tag="json:\"passRule\" bun:\"pass_rule\""]
  FIELD PassRatio : github.com/coldsmirk/vef-framework-go/decimal.Decimal [field_order=11 tag="json:\"passRatio\" bun:\"pass_ratio,type:numeric(5,2)\""]
  FIELD EmptyAssigneeAction : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction [field_order=12 tag="json:\"emptyAssigneeAction\" bun:\"empty_assignee_action\""]
  FIELD FallbackUserIDs : []string [field_order=13 tag="json:\"fallbackUserIds\" bun:\"fallback_user_ids,type:jsonb\""]
  FIELD AdminUserIDs : []string [field_order=14 tag="json:\"adminUserIds\" bun:\"admin_user_ids,type:jsonb\""]
  FIELD SameApplicantAction : github.com/coldsmirk/vef-framework-go/approval.SameApplicantAction [field_order=15 tag="json:\"sameApplicantAction\" bun:\"same_applicant_action\""]
  FIELD IsRollbackAllowed : bool [field_order=16 tag="json:\"isRollbackAllowed\" bun:\"is_rollback_allowed\""]
  FIELD RollbackType : github.com/coldsmirk/vef-framework-go/approval.RollbackType [field_order=17 tag="json:\"rollbackType\" bun:\"rollback_type\""]
  FIELD RollbackDataStrategy : github.com/coldsmirk/vef-framework-go/approval.RollbackDataStrategy [field_order=18 tag="json:\"rollbackDataStrategy\" bun:\"rollback_data_strategy\""]
  FIELD RollbackTargetKeys : []string [field_order=19 tag="json:\"rollbackTargetKeys\" bun:\"rollback_target_keys,type:jsonb,nullzero\""]
  FIELD IsAddAssigneeAllowed : bool [field_order=20 tag="json:\"isAddAssigneeAllowed\" bun:\"is_add_assignee_allowed\""]
  FIELD AddAssigneeTypes : []github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType [field_order=21 tag="json:\"addAssigneeTypes\" bun:\"add_assignee_types,type:jsonb\""]
  FIELD IsRemoveAssigneeAllowed : bool [field_order=22 tag="json:\"isRemoveAssigneeAllowed\" bun:\"is_remove_assignee_allowed\""]
  FIELD FieldPermissions : map[string]github.com/coldsmirk/vef-framework-go/approval.Permission [field_order=23 tag="json:\"fieldPermissions\" bun:\"field_permissions,type:jsonb\""]
  FIELD IsManualCCAllowed : bool [field_order=24 tag="json:\"isManualCcAllowed\" bun:\"is_manual_cc_allowed\""]
  FIELD IsTransferAllowed : bool [field_order=25 tag="json:\"isTransferAllowed\" bun:\"is_transfer_allowed\""]
  FIELD IsOpinionRequired : bool [field_order=26 tag="json:\"isOpinionRequired\" bun:\"is_opinion_required\""]
  FIELD TimeoutHours : int [field_order=27 tag="json:\"timeoutHours\" bun:\"timeout_hours\""]
  FIELD TimeoutAction : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction [field_order=28 tag="json:\"timeoutAction\" bun:\"timeout_action\""]
  FIELD TimeoutNotifyBeforeHours : int [field_order=29 tag="json:\"timeoutNotifyBeforeHours\" bun:\"timeout_notify_before_hours\""]
  FIELD UrgeCooldownMinutes : int [field_order=30 tag="json:\"urgeCooldownMinutes\" bun:\"urge_cooldown_minutes\""]
  FIELD ConsecutiveApproverAction : github.com/coldsmirk/vef-framework-go/approval.ConsecutiveApproverAction [field_order=31 tag="json:\"consecutiveApproverAction\" bun:\"consecutive_approver_action\""]
  FIELD IsReadConfirmRequired : bool [field_order=32 tag="json:\"isReadConfirmRequired\" bun:\"is_read_confirm_required\""]
  FIELD Branches : []github.com/coldsmirk/vef-framework-go/approval.ConditionBranch [field_order=33 tag="json:\"branches\" bun:\"branches,type:jsonb,nullzero\""]
TYPE FlowNodeAssignee : github.com/coldsmirk/vef-framework-go/approval.FlowNodeAssignee
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow_node_assignee,alias:afna\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD NodeID : string [field_order=3 tag="json:\"nodeId\" bun:\"node_id\""]
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.AssigneeKind [field_order=4 tag="json:\"kind\" bun:\"kind\""]
  FIELD IDs : []string [field_order=5 tag="json:\"ids\" bun:\"ids,type:jsonb\""]
  FIELD FormField : *string [field_order=6 tag="json:\"formField\" bun:\"form_field,nullzero\""]
  FIELD SortOrder : int [field_order=7 tag="json:\"sortOrder\" bun:\"sort_order\""]
TYPE FlowNodeCC : github.com/coldsmirk/vef-framework-go/approval.FlowNodeCC
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow_node_cc,alias:afnc\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD NodeID : string [field_order=3 tag="json:\"nodeId\" bun:\"node_id\""]
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.CCKind [field_order=4 tag="json:\"kind\" bun:\"kind\""]
  FIELD IDs : []string [field_order=5 tag="json:\"ids\" bun:\"ids,type:jsonb\""]
  FIELD FormField : *string [field_order=6 tag="json:\"formField\" bun:\"form_field,nullzero\""]
  FIELD Timing : github.com/coldsmirk/vef-framework-go/approval.CCTiming [field_order=7 tag="json:\"timing\" bun:\"timing\""]
TYPE FlowPublishedEvent : github.com/coldsmirk/vef-framework-go/approval.FlowPublishedEvent
  FIELD FlowEventBase : github.com/coldsmirk/vef-framework-go/approval.FlowEventBase [field_order=1 tag=""]
  FIELD VersionID : string [field_order=2 tag="json:\"versionId\""]
  FIELD Code : string [promoted_from=FlowEventBase depth=1 field_order=3 tag="json:\"code\""]
  FIELD FlowID : string [promoted_from=FlowEventBase depth=1 field_order=1 tag="json:\"flowId\""]
  FIELD Name : string [promoted_from=FlowEventBase depth=1 field_order=4 tag="json:\"name\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=FlowEventBase depth=1 field_order=5 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=FlowEventBase depth=1 field_order=2 tag="json:\"tenantId\""]
  METHOD EventType : func() string
TYPE FlowToggledEvent : github.com/coldsmirk/vef-framework-go/approval.FlowToggledEvent
  FIELD FlowEventBase : github.com/coldsmirk/vef-framework-go/approval.FlowEventBase [field_order=1 tag=""]
  FIELD IsActive : bool [field_order=2 tag="json:\"isActive\""]
  FIELD Code : string [promoted_from=FlowEventBase depth=1 field_order=3 tag="json:\"code\""]
  FIELD FlowID : string [promoted_from=FlowEventBase depth=1 field_order=1 tag="json:\"flowId\""]
  FIELD Name : string [promoted_from=FlowEventBase depth=1 field_order=4 tag="json:\"name\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=FlowEventBase depth=1 field_order=5 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=FlowEventBase depth=1 field_order=2 tag="json:\"tenantId\""]
  METHOD EventType : func() string
TYPE FlowUpdatedEvent : github.com/coldsmirk/vef-framework-go/approval.FlowUpdatedEvent
  FIELD FlowEventBase : github.com/coldsmirk/vef-framework-go/approval.FlowEventBase [field_order=1 tag=""]
  FIELD Code : string [promoted_from=FlowEventBase depth=1 field_order=3 tag="json:\"code\""]
  FIELD FlowID : string [promoted_from=FlowEventBase depth=1 field_order=1 tag="json:\"flowId\""]
  FIELD Name : string [promoted_from=FlowEventBase depth=1 field_order=4 tag="json:\"name\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=FlowEventBase depth=1 field_order=5 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=FlowEventBase depth=1 field_order=2 tag="json:\"tenantId\""]
  METHOD EventType : func() string
TYPE FlowVersion : github.com/coldsmirk/vef-framework-go/approval.FlowVersion
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_flow_version,alias:afv\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD FlowID : string [field_order=3 tag="json:\"flowId\" bun:\"flow_id\""]
  FIELD Version : int [field_order=4 tag="json:\"version\" bun:\"version\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/approval.VersionStatus [field_order=5 tag="json:\"status\" bun:\"status\""]
  FIELD Description : *string [field_order=6 tag="json:\"description\" bun:\"description,nullzero\""]
  FIELD StorageMode : github.com/coldsmirk/vef-framework-go/approval.StorageMode [field_order=7 tag="json:\"storageMode\" bun:\"storage_mode\""]
  FIELD FlowSchema : *github.com/coldsmirk/vef-framework-go/approval.FlowDefinition [field_order=8 tag="json:\"flowSchema\" bun:\"flow_schema,type:jsonb,nullzero\""]
  FIELD FormSchema : encoding/json.RawMessage [field_order=9 tag="json:\"formSchema\" bun:\"form_schema,type:jsonb,nullzero\""]
  FIELD FormFields : []github.com/coldsmirk/vef-framework-go/approval.FormFieldDefinition [field_order=10 tag="json:\"formFields\" bun:\"form_fields,type:jsonb,nullzero\""]
  FIELD PublishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=11 tag="json:\"publishedAt\" bun:\"published_at,nullzero\""]
  FIELD PublishedBy : *string [field_order=12 tag="json:\"publishedBy\" bun:\"published_by,nullzero\""]
FUNC ForFlows : func(codes ...string) github.com/coldsmirk/vef-framework-go/approval.InstanceFilter
FUNC ForTenants : func(ids ...string) github.com/coldsmirk/vef-framework-go/approval.InstanceFilter
TYPE FormData : github.com/coldsmirk/vef-framework-go/approval.FormData
  METHOD Clone : func() (github.com/coldsmirk/vef-framework-go/approval.FormData, error)
  METHOD Get : func(key string) any
  METHOD Set : func(key string, val any)
  METHOD ToMap : func() map[string]any
TYPE FormFieldDefinition : github.com/coldsmirk/vef-framework-go/approval.FormFieldDefinition
  FIELD Key : string [field_order=1 tag="json:\"key\""]
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.FieldKind [field_order=2 tag="json:\"kind\""]
  FIELD Label : string [field_order=3 tag="json:\"label\""]
  FIELD Placeholder : string [field_order=4 tag="json:\"placeholder,omitempty\""]
  FIELD DefaultValue : any [field_order=5 tag="json:\"defaultValue,omitempty\""]
  FIELD IsRequired : bool [field_order=6 tag="json:\"isRequired,omitempty\""]
  FIELD Options : []github.com/coldsmirk/vef-framework-go/approval.FieldOption [field_order=7 tag="json:\"options,omitempty\""]
  FIELD Validation : *github.com/coldsmirk/vef-framework-go/approval.ValidationRule [field_order=8 tag="json:\"validation,omitempty\""]
  FIELD Props : map[string]any [field_order=9 tag="json:\"props,omitempty\""]
  FIELD SortOrder : int [field_order=10 tag="json:\"sortOrder\""]
  FIELD ColumnType : github.com/coldsmirk/vef-framework-go/approval.ColumnDataType [field_order=11 tag="json:\"columnType,omitempty\""]
  FIELD Scale : *int [field_order=12 tag="json:\"scale,omitempty\""]
  FIELD Columns : []github.com/coldsmirk/vef-framework-go/approval.FormFieldDefinition [field_order=13 tag="json:\"columns,omitempty\""]
TYPE FormSchemaParser : github.com/coldsmirk/vef-framework-go/approval.FormSchemaParser
  METHOD ParseFormFields : func(schema encoding/json.RawMessage) ([]github.com/coldsmirk/vef-framework-go/approval.FormFieldDefinition, error)
TYPE FormSnapshot : github.com/coldsmirk/vef-framework-go/approval.FormSnapshot
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_form_snapshot,alias:afs\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD CreationTrackedModel : github.com/coldsmirk/vef-framework-go/orm.CreationTrackedModel [field_order=3 tag=""]
  FIELD InstanceID : string [field_order=4 tag="json:\"instanceId\" bun:\"instance_id\""]
  FIELD NodeID : string [field_order=5 tag="json:\"nodeId\" bun:\"node_id\""]
  FIELD FormData : map[string]any [field_order=6 tag="json:\"formData\" bun:\"form_data,type:jsonb\""]
TYPE FormTable : github.com/coldsmirk/vef-framework-go/approval.FormTable
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_form_table,alias:aft\""]
  FIELD CreationAuditedModel : github.com/coldsmirk/vef-framework-go/orm.CreationAuditedModel [field_order=2 tag=""]
  FIELD FlowID : string [field_order=3 tag="json:\"flowId\" bun:\"flow_id\""]
  FIELD VersionID : string [field_order=4 tag="json:\"versionId\" bun:\"version_id\""]
  FIELD PhysicalTableName : string [field_order=5 tag="json:\"physicalTableName\" bun:\"physical_table_name\""]
  FIELD SourceFieldKey : string [field_order=6 tag="json:\"sourceFieldKey\" bun:\"source_field_key\""]
TYPE FormTableColumn : github.com/coldsmirk/vef-framework-go/approval.FormTableColumn
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_form_table_column,alias:aftc\""]
  FIELD CreationAuditedModel : github.com/coldsmirk/vef-framework-go/orm.CreationAuditedModel [field_order=2 tag=""]
  FIELD FormTableID : string [field_order=3 tag="json:\"formTableId\" bun:\"form_table_id\""]
  FIELD ColumnName : string [field_order=4 tag="json:\"columnName\" bun:\"column_name\""]
  FIELD ColumnType : string [field_order=5 tag="json:\"columnType\" bun:\"column_type\""]
  FIELD IsNullable : bool [field_order=6 tag="json:\"isNullable\" bun:\"is_nullable\""]
  FIELD SourceFieldKey : *string [field_order=7 tag="json:\"sourceFieldKey\" bun:\"source_field_key,nullzero\""]
  FIELD SortOrder : int [field_order=8 tag="json:\"sortOrder\" bun:\"sort_order\""]
TYPE HandleNodeData : github.com/coldsmirk/vef-framework-go/approval.HandleNodeData
  FIELD BaseNodeData : github.com/coldsmirk/vef-framework-go/approval.BaseNodeData [field_order=1 tag=""]
  FIELD TaskNodeData : github.com/coldsmirk/vef-framework-go/approval.TaskNodeData [field_order=2 tag=""]
  FIELD AdminUserIDs : []string [promoted_from=TaskNodeData depth=1 field_order=5 tag="json:\"adminUserIds,omitempty\""]
  FIELD Assignees : []github.com/coldsmirk/vef-framework-go/approval.AssigneeDefinition [promoted_from=TaskNodeData depth=1 field_order=1 tag="json:\"assignees,omitempty\""]
  FIELD CCs : []github.com/coldsmirk/vef-framework-go/approval.CCDefinition [promoted_from=TaskNodeData depth=1 field_order=12 tag="json:\"ccs,omitempty\""]
  FIELD Description : *string [promoted_from=BaseNodeData depth=1 field_order=2 tag="json:\"description,omitempty\""]
  FIELD EmptyAssigneeAction : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction [promoted_from=TaskNodeData depth=1 field_order=3 tag="json:\"emptyAssigneeAction,omitempty\""]
  FIELD ExecutionType : github.com/coldsmirk/vef-framework-go/approval.ExecutionType [promoted_from=TaskNodeData depth=1 field_order=2 tag="json:\"executionType,omitempty\""]
  FIELD FallbackUserIDs : []string [promoted_from=TaskNodeData depth=1 field_order=4 tag="json:\"fallbackUserIds,omitempty\""]
  FIELD FieldPermissions : map[string]github.com/coldsmirk/vef-framework-go/approval.Permission [promoted_from=TaskNodeData depth=1 field_order=13 tag="json:\"fieldPermissions,omitempty\""]
  FIELD IsOpinionRequired : bool [promoted_from=TaskNodeData depth=1 field_order=7 tag="json:\"isOpinionRequired,omitempty\""]
  FIELD IsTransferAllowed : *bool [promoted_from=TaskNodeData depth=1 field_order=6 tag="json:\"isTransferAllowed,omitempty\""]
  FIELD Name : string [promoted_from=BaseNodeData depth=1 field_order=1 tag="json:\"name,omitempty\""]
  FIELD TimeoutAction : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction [promoted_from=TaskNodeData depth=1 field_order=9 tag="json:\"timeoutAction,omitempty\""]
  FIELD TimeoutHours : int [promoted_from=TaskNodeData depth=1 field_order=8 tag="json:\"timeoutHours,omitempty\""]
  FIELD TimeoutNotifyBeforeHours : int [promoted_from=TaskNodeData depth=1 field_order=10 tag="json:\"timeoutNotifyBeforeHours,omitempty\""]
  FIELD UrgeCooldownMinutes : int [promoted_from=TaskNodeData depth=1 field_order=11 tag="json:\"urgeCooldownMinutes,omitempty\""]
  METHOD ApplyTo : func(node *github.com/coldsmirk/vef-framework-go/approval.FlowNode)
  METHOD GetAssignees : func() []github.com/coldsmirk/vef-framework-go/approval.AssigneeDefinition
  METHOD GetCCs : func() []github.com/coldsmirk/vef-framework-go/approval.CCDefinition
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.NodeKind
CONST InitiatorDepartment : github.com/coldsmirk/vef-framework-go/approval.InitiatorKind = "department"
TYPE InitiatorKind : github.com/coldsmirk/vef-framework-go/approval.InitiatorKind
  METHOD IsValid : func() bool
CONST InitiatorRole : github.com/coldsmirk/vef-framework-go/approval.InitiatorKind = "role"
CONST InitiatorUser : github.com/coldsmirk/vef-framework-go/approval.InitiatorKind = "user"
TYPE Instance : github.com/coldsmirk/vef-framework-go/approval.Instance
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_instance,alias:ai\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD TenantID : string [field_order=3 tag="json:\"tenantId\" bun:\"tenant_id\""]
  FIELD FlowID : string [field_order=4 tag="json:\"flowId\" bun:\"flow_id\""]
  FIELD FlowCode : string [field_order=5 tag="json:\"flowCode\" bun:\"flow_code\""]
  FIELD FlowVersionID : string [field_order=6 tag="json:\"flowVersionId\" bun:\"flow_version_id\""]
  FIELD Title : string [field_order=7 tag="json:\"title\" bun:\"title\""]
  FIELD InstanceNo : string [field_order=8 tag="json:\"instanceNo\" bun:\"instance_no\""]
  FIELD ApplicantID : string [field_order=9 tag="json:\"applicantId\" bun:\"applicant_id\""]
  FIELD ApplicantName : string [field_order=10 tag="json:\"applicantName\" bun:\"applicant_name\""]
  FIELD ApplicantDepartmentID : *string [field_order=11 tag="json:\"applicantDepartmentId\" bun:\"applicant_department_id,nullzero\""]
  FIELD ApplicantDepartmentName : *string [field_order=12 tag="json:\"applicantDepartmentName\" bun:\"applicant_department_name,nullzero\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus [field_order=13 tag="json:\"status\" bun:\"status\""]
  FIELD CurrentNodeID : *string [field_order=14 tag="json:\"currentNodeId\" bun:\"current_node_id,nullzero\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=15 tag="json:\"finishedAt\" bun:\"finished_at,nullzero\""]
  FIELD BusinessRef : *string [field_order=16 tag="json:\"businessRef\" bun:\"business_ref,nullzero\""]
  FIELD FormData : map[string]any [field_order=17 tag="json:\"formData\" bun:\"form_data,type:jsonb,nullzero\""]
  FIELD Globals : map[string]any [field_order=18 tag="json:\"globals\" bun:\"globals,type:jsonb,nullzero\""]
  METHOD Applicant : func() github.com/coldsmirk/vef-framework-go/approval.UserInfo
CONST InstanceApproved : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus = "approved"
TYPE InstanceBindingFailedEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceBindingFailedEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD Trigger : github.com/coldsmirk/vef-framework-go/approval.BindingTrigger [field_order=2 tag="json:\"trigger\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus [field_order=3 tag="json:\"status\""]
  FIELD BusinessTable : string [field_order=4 tag="json:\"businessTable\""]
  FIELD ErrorMessage : string [field_order=5 tag="json:\"errorMessage\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE InstanceCompletedEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceCompletedEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD FinalStatus : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus [field_order=2 tag="json:\"finalStatus\""]
  FIELD FinishedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=3 tag="json:\"finishedAt\""]
  FIELD Reason : *string [field_order=4 tag="json:\"reason,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE InstanceCreatedEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceCreatedEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE InstanceEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceEvent
  METHOD EventType : func() string
TYPE InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase
  FIELD InstanceID : string [field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [field_order=2 tag="json:\"instanceNo\""]
  FIELD TenantID : string [field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [field_order=4 tag="json:\"title\""]
  FIELD FlowID : string [field_order=5 tag="json:\"flowId\""]
  FIELD FlowCode : string [field_order=6 tag="json:\"flowCode\""]
  FIELD BusinessRef : *string [field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=8 tag="json:\"applicant\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=9 tag="json:\"occurredTime\""]
TYPE InstanceFilter : github.com/coldsmirk/vef-framework-go/approval.InstanceFilter
  FIELD FlowCodes : []string [field_order=1 tag=""]
  FIELD TenantIDs : []string [field_order=2 tag=""]
  METHOD Matches : func(flowCode string, tenantID string) bool
TYPE InstanceFlowGraph : github.com/coldsmirk/vef-framework-go/approval.InstanceFlowGraph
  FIELD Nodes : []github.com/coldsmirk/vef-framework-go/approval.FlowGraphNode [field_order=1 tag="json:\"nodes\""]
  FIELD Edges : []github.com/coldsmirk/vef-framework-go/approval.FlowGraphEdge [field_order=2 tag="json:\"edges\""]
TYPE InstanceGlobalsResolver : github.com/coldsmirk/vef-framework-go/approval.InstanceGlobalsResolver
  METHOD Resolve : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, flowCode string) (map[string]any, error)
TYPE InstanceLifecycleHook : github.com/coldsmirk/vef-framework-go/approval.InstanceLifecycleHook
  METHOD OnInstanceCompleted : func(ctx context.Context, db github.com/coldsmirk/vef-framework-go/orm.DB, instance *github.com/coldsmirk/vef-framework-go/approval.Instance, finalStatus github.com/coldsmirk/vef-framework-go/approval.InstanceStatus) error
  METHOD OnInstanceCreated : func(ctx context.Context, db github.com/coldsmirk/vef-framework-go/orm.DB, instance *github.com/coldsmirk/vef-framework-go/approval.Instance) error
TYPE InstanceNoGenerator : github.com/coldsmirk/vef-framework-go/approval.InstanceNoGenerator
  METHOD Generate : func(ctx context.Context, flowCode string) (string, error)
CONST InstanceRejected : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus = "rejected"
TYPE InstanceResubmittedEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceResubmittedEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"operator\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST InstanceReturned : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus = "returned"
TYPE InstanceReturnedEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceReturnedEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD FromNodeID : string [field_order=2 tag="json:\"fromNodeId\""]
  FIELD FromNodeName : string [field_order=3 tag="json:\"fromNodeName\""]
  FIELD ToNodeID : string [field_order=4 tag="json:\"toNodeId\""]
  FIELD ToNodeName : string [field_order=5 tag="json:\"toNodeName\""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=6 tag="json:\"operator\""]
  FIELD Opinion : *string [field_order=7 tag="json:\"opinion,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE InstanceRolledBackEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceRolledBackEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD FromNodeID : string [field_order=2 tag="json:\"fromNodeId\""]
  FIELD FromNodeName : string [field_order=3 tag="json:\"fromNodeName\""]
  FIELD ToNodeID : string [field_order=4 tag="json:\"toNodeId\""]
  FIELD ToNodeName : string [field_order=5 tag="json:\"toNodeName\""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=6 tag="json:\"operator\""]
  FIELD Opinion : *string [field_order=7 tag="json:\"opinion,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST InstanceRunning : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus = "running"
TYPE InstanceStatus : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus
  METHOD IsFinal : func() bool
  METHOD String : func() string
TYPE InstanceSubscribeOption : github.com/coldsmirk/vef-framework-go/approval.InstanceSubscribeOption
CONST InstanceTerminated : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus = "terminated"
CONST InstanceWithdrawn : github.com/coldsmirk/vef-framework-go/approval.InstanceStatus = "withdrawn"
TYPE InstanceWithdrawnEvent : github.com/coldsmirk/vef-framework-go/approval.InstanceWithdrawnEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"operator\""]
  FIELD Reason : *string [field_order=3 tag="json:\"reason,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
FUNC IsSuperAdmin : func(p *github.com/coldsmirk/vef-framework-go/security.Principal) bool
FUNC NewAssigneesAddedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, addType github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType, assignees []github.com/coldsmirk/vef-framework-go/approval.UserInfo) *github.com/coldsmirk/vef-framework-go/approval.AssigneesAddedEvent
FUNC NewAssigneesRemovedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, assignees []github.com/coldsmirk/vef-framework-go/approval.UserInfo) *github.com/coldsmirk/vef-framework-go/approval.AssigneesRemovedEvent
FUNC NewCCNotifiedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, recipients []github.com/coldsmirk/vef-framework-go/approval.UserInfo, isManual bool) *github.com/coldsmirk/vef-framework-go/approval.CCNotifiedEvent
FUNC NewFilteredLifecycleHook : func(hook github.com/coldsmirk/vef-framework-go/approval.InstanceLifecycleHook, filters ...github.com/coldsmirk/vef-framework-go/approval.InstanceFilter) github.com/coldsmirk/vef-framework-go/approval.InstanceLifecycleHook
FUNC NewFlowCreatedEvent : func(flow *github.com/coldsmirk/vef-framework-go/approval.Flow) *github.com/coldsmirk/vef-framework-go/approval.FlowCreatedEvent
FUNC NewFlowDeployedEvent : func(flow *github.com/coldsmirk/vef-framework-go/approval.Flow, versionID string, version int) *github.com/coldsmirk/vef-framework-go/approval.FlowDeployedEvent
FUNC NewFlowEventBase : func(flow *github.com/coldsmirk/vef-framework-go/approval.Flow) github.com/coldsmirk/vef-framework-go/approval.FlowEventBase
FUNC NewFlowPublishedEvent : func(flow *github.com/coldsmirk/vef-framework-go/approval.Flow, versionID string) *github.com/coldsmirk/vef-framework-go/approval.FlowPublishedEvent
FUNC NewFlowToggledEvent : func(flow *github.com/coldsmirk/vef-framework-go/approval.Flow, isActive bool) *github.com/coldsmirk/vef-framework-go/approval.FlowToggledEvent
FUNC NewFlowUpdatedEvent : func(flow *github.com/coldsmirk/vef-framework-go/approval.Flow) *github.com/coldsmirk/vef-framework-go/approval.FlowUpdatedEvent
FUNC NewFormData : func(data map[string]any) github.com/coldsmirk/vef-framework-go/approval.FormData
FUNC NewInstanceBindingFailedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, trigger github.com/coldsmirk/vef-framework-go/approval.BindingTrigger, status github.com/coldsmirk/vef-framework-go/approval.InstanceStatus, businessTable string, errorMessage string) *github.com/coldsmirk/vef-framework-go/approval.InstanceBindingFailedEvent
FUNC NewInstanceCompletedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, finalStatus github.com/coldsmirk/vef-framework-go/approval.InstanceStatus) *github.com/coldsmirk/vef-framework-go/approval.InstanceCompletedEvent
FUNC NewInstanceCreatedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance) *github.com/coldsmirk/vef-framework-go/approval.InstanceCreatedEvent
FUNC NewInstanceEventBase : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance) github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase
FUNC NewInstanceResubmittedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, operator github.com/coldsmirk/vef-framework-go/approval.UserInfo) *github.com/coldsmirk/vef-framework-go/approval.InstanceResubmittedEvent
FUNC NewInstanceReturnedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, fromNode *github.com/coldsmirk/vef-framework-go/approval.FlowNode, toNode *github.com/coldsmirk/vef-framework-go/approval.FlowNode, operator github.com/coldsmirk/vef-framework-go/approval.UserInfo, opinion *string) *github.com/coldsmirk/vef-framework-go/approval.InstanceReturnedEvent
FUNC NewInstanceRolledBackEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, fromNode *github.com/coldsmirk/vef-framework-go/approval.FlowNode, toNode *github.com/coldsmirk/vef-framework-go/approval.FlowNode, operator github.com/coldsmirk/vef-framework-go/approval.UserInfo, opinion *string) *github.com/coldsmirk/vef-framework-go/approval.InstanceRolledBackEvent
FUNC NewInstanceWithdrawnEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, operator github.com/coldsmirk/vef-framework-go/approval.UserInfo, reason *string) *github.com/coldsmirk/vef-framework-go/approval.InstanceWithdrawnEvent
FUNC NewNodeAutoPassedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, reason string) *github.com/coldsmirk/vef-framework-go/approval.NodeAutoPassedEvent
FUNC NewTaskApprovedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, operator github.com/coldsmirk/vef-framework-go/approval.UserInfo, opinion string) *github.com/coldsmirk/vef-framework-go/approval.TaskApprovedEvent
FUNC NewTaskCanceledEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, reason string) *github.com/coldsmirk/vef-framework-go/approval.TaskCanceledEvent
FUNC NewTaskCreatedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode) *github.com/coldsmirk/vef-framework-go/approval.TaskCreatedEvent
FUNC NewTaskDeadlineWarningEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, hoursLeft int) *github.com/coldsmirk/vef-framework-go/approval.TaskDeadlineWarningEvent
FUNC NewTaskEventBase : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode) github.com/coldsmirk/vef-framework-go/approval.TaskEventBase
FUNC NewTaskHandledEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, operator github.com/coldsmirk/vef-framework-go/approval.UserInfo, opinion string) *github.com/coldsmirk/vef-framework-go/approval.TaskHandledEvent
FUNC NewTaskReassignedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, from github.com/coldsmirk/vef-framework-go/approval.UserInfo, to github.com/coldsmirk/vef-framework-go/approval.UserInfo, reason string) *github.com/coldsmirk/vef-framework-go/approval.TaskReassignedEvent
FUNC NewTaskRejectedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, operator github.com/coldsmirk/vef-framework-go/approval.UserInfo, opinion string) *github.com/coldsmirk/vef-framework-go/approval.TaskRejectedEvent
FUNC NewTaskTimedOutEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode) *github.com/coldsmirk/vef-framework-go/approval.TaskTimedOutEvent
FUNC NewTaskTransferredEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, from github.com/coldsmirk/vef-framework-go/approval.UserInfo, to github.com/coldsmirk/vef-framework-go/approval.UserInfo, reason string) *github.com/coldsmirk/vef-framework-go/approval.TaskTransferredEvent
FUNC NewTaskUrgedEvent : func(instance *github.com/coldsmirk/vef-framework-go/approval.Instance, task *github.com/coldsmirk/vef-framework-go/approval.Task, node *github.com/coldsmirk/vef-framework-go/approval.FlowNode, urger github.com/coldsmirk/vef-framework-go/approval.UserInfo, message string) *github.com/coldsmirk/vef-framework-go/approval.TaskUrgedEvent
CONST NodeApproval : github.com/coldsmirk/vef-framework-go/approval.NodeKind = "approval"
TYPE NodeAutoPassedEvent : github.com/coldsmirk/vef-framework-go/approval.NodeAutoPassedEvent
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD NodeID : string [field_order=2 tag="json:\"nodeId\""]
  FIELD NodeName : string [field_order=3 tag="json:\"nodeName\""]
  FIELD Reason : string [field_order=4 tag="json:\"reason\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST NodeCC : github.com/coldsmirk/vef-framework-go/approval.NodeKind = "cc"
CONST NodeCondition : github.com/coldsmirk/vef-framework-go/approval.NodeKind = "condition"
TYPE NodeData : github.com/coldsmirk/vef-framework-go/approval.NodeData
  METHOD ApplyTo : func(node *github.com/coldsmirk/vef-framework-go/approval.FlowNode)
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.NodeKind
TYPE NodeDefinition : github.com/coldsmirk/vef-framework-go/approval.NodeDefinition
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.NodeKind [field_order=2 tag="json:\"kind\""]
  FIELD Position : github.com/coldsmirk/vef-framework-go/approval.Position [field_order=3 tag="json:\"position\""]
  FIELD Data : encoding/json.RawMessage [field_order=4 tag="json:\"data,omitempty\""]
  METHOD ParseData : func() (github.com/coldsmirk/vef-framework-go/approval.NodeData, error)
CONST NodeEnd : github.com/coldsmirk/vef-framework-go/approval.NodeKind = "end"
CONST NodeHandle : github.com/coldsmirk/vef-framework-go/approval.NodeKind = "handle"
TYPE NodeKind : github.com/coldsmirk/vef-framework-go/approval.NodeKind
TYPE NodeParticipant : github.com/coldsmirk/vef-framework-go/approval.NodeParticipant
  FIELD TaskID : string [field_order=1 tag="json:\"taskId\""]
  FIELD User : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"user\""]
  FIELD Delegator : *github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=3 tag="json:\"delegator,omitempty\""]
  FIELD Status : string [field_order=4 tag="json:\"status\""]
  FIELD Deadline : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=5 tag="json:\"deadline,omitempty\""]
  FIELD IsTimeout : bool [field_order=6 tag="json:\"isTimeout,omitempty\""]
  FIELD Opinion : *string [field_order=7 tag="json:\"opinion,omitempty\""]
  FIELD Attachments : []string [field_order=8 tag="json:\"attachments,omitempty\""]
  FIELD ActionTime : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=9 tag="json:\"actionTime,omitempty\""]
  FIELD TransferTo : *github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=10 tag="json:\"transferTo,omitempty\""]
CONST NodeProgressActive : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus = "active"
CONST NodeProgressCanceled : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus = "canceled"
CONST NodeProgressPassed : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus = "passed"
CONST NodeProgressPending : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus = "pending"
CONST NodeProgressRejected : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus = "rejected"
CONST NodeProgressReturned : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus = "returned"
TYPE NodeProgressStatus : github.com/coldsmirk/vef-framework-go/approval.NodeProgressStatus
CONST NodeStart : github.com/coldsmirk/vef-framework-go/approval.NodeKind = "start"
TYPE NodeVisit : github.com/coldsmirk/vef-framework-go/approval.NodeVisit
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_node_visit,alias:anv\""]
  FIELD CreationAuditedModel : github.com/coldsmirk/vef-framework-go/orm.CreationAuditedModel [field_order=2 tag=""]
  FIELD TenantID : string [field_order=3 tag="json:\"tenantId\" bun:\"tenant_id\""]
  FIELD InstanceID : string [field_order=4 tag="json:\"instanceId\" bun:\"instance_id\""]
  FIELD NodeID : string [field_order=5 tag="json:\"nodeId\" bun:\"node_id\""]
  FIELD Sequence : int [field_order=6 tag="json:\"sequence\" bun:\"sequence\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus [field_order=7 tag="json:\"status\" bun:\"status\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=8 tag="json:\"finishedAt\" bun:\"finished_at,nullzero\""]
CONST NodeVisitActive : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus = "active"
CONST NodeVisitCanceled : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus = "canceled"
CONST NodeVisitPassed : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus = "passed"
CONST NodeVisitRejected : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus = "rejected"
CONST NodeVisitReturned : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus = "returned"
TYPE NodeVisitStatus : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus
  METHOD String : func() string
CONST OperatorContains : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "contains"
CONST OperatorEndsWith : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "ends_with"
CONST OperatorEquals : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "eq"
CONST OperatorGreater : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "gt"
CONST OperatorGreaterOrEq : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "gte"
CONST OperatorIn : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "in"
CONST OperatorIsEmpty : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "is_empty"
CONST OperatorIsNotEmpty : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "is_not_empty"
CONST OperatorLess : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "lt"
CONST OperatorLessOrEq : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "lte"
CONST OperatorNotContains : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "not_contains"
CONST OperatorNotEquals : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "ne"
CONST OperatorNotIn : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "not_in"
CONST OperatorStartsWith : github.com/coldsmirk/vef-framework-go/approval.ConditionOperator = "starts_with"
CONST PassAll : github.com/coldsmirk/vef-framework-go/approval.PassRule = "all"
CONST PassAny : github.com/coldsmirk/vef-framework-go/approval.PassRule = "any"
CONST PassRatio : github.com/coldsmirk/vef-framework-go/approval.PassRule = "ratio"
TYPE PassRule : github.com/coldsmirk/vef-framework-go/approval.PassRule
  METHOD IsValid : func() bool
TYPE PassRuleContext : github.com/coldsmirk/vef-framework-go/approval.PassRuleContext
  FIELD ApprovedCount : int [field_order=1 tag=""]
  FIELD RejectedCount : int [field_order=2 tag=""]
  FIELD TotalCount : int [field_order=3 tag=""]
  FIELD PassRatio : float64 [field_order=4 tag=""]
CONST PassRulePassed : github.com/coldsmirk/vef-framework-go/approval.PassRuleResult = 1
CONST PassRulePending : github.com/coldsmirk/vef-framework-go/approval.PassRuleResult = 0
CONST PassRuleRejected : github.com/coldsmirk/vef-framework-go/approval.PassRuleResult = 2
TYPE PassRuleResult : github.com/coldsmirk/vef-framework-go/approval.PassRuleResult
TYPE PassRuleStrategy : github.com/coldsmirk/vef-framework-go/approval.PassRuleStrategy
  METHOD Evaluate : func(ctx github.com/coldsmirk/vef-framework-go/approval.PassRuleContext) github.com/coldsmirk/vef-framework-go/approval.PassRuleResult
  METHOD Rule : func() github.com/coldsmirk/vef-framework-go/approval.PassRule
FUNC PayloadOccurredAt : func(e github.com/coldsmirk/vef-framework-go/approval.DomainEvent) github.com/coldsmirk/vef-framework-go/timex.DateTime
TYPE Permission : github.com/coldsmirk/vef-framework-go/approval.Permission
  METHOD IsValid : func() bool
CONST PermissionEditable : github.com/coldsmirk/vef-framework-go/approval.Permission = "editable"
CONST PermissionHidden : github.com/coldsmirk/vef-framework-go/approval.Permission = "hidden"
CONST PermissionRequired : github.com/coldsmirk/vef-framework-go/approval.Permission = "required"
CONST PermissionVisible : github.com/coldsmirk/vef-framework-go/approval.Permission = "visible"
TYPE Position : github.com/coldsmirk/vef-framework-go/approval.Position
  FIELD X : float64 [field_order=1 tag="json:\"x\""]
  FIELD Y : float64 [field_order=2 tag="json:\"y\""]
TYPE PrincipalDepartmentResolver : github.com/coldsmirk/vef-framework-go/approval.PrincipalDepartmentResolver
  METHOD Resolve : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (departmentID *string, departmentName *string, err error)
TYPE PrincipalTenantResolver : github.com/coldsmirk/vef-framework-go/approval.PrincipalTenantResolver
  METHOD Resolve : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (string, error)
TYPE ResolvedAssignee : github.com/coldsmirk/vef-framework-go/approval.ResolvedAssignee
  FIELD User : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=1 tag=""]
  FIELD Delegator : *github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag=""]
TYPE RoleMembershipChecker : github.com/coldsmirk/vef-framework-go/approval.RoleMembershipChecker
  METHOD UserHasRole : func(ctx context.Context, userID string, roleID string) (bool, error)
CONST RollbackAny : github.com/coldsmirk/vef-framework-go/approval.RollbackType = "any"
CONST RollbackDataClear : github.com/coldsmirk/vef-framework-go/approval.RollbackDataStrategy = "clear"
CONST RollbackDataKeep : github.com/coldsmirk/vef-framework-go/approval.RollbackDataStrategy = "keep"
TYPE RollbackDataStrategy : github.com/coldsmirk/vef-framework-go/approval.RollbackDataStrategy
  METHOD IsValid : func() bool
CONST RollbackNone : github.com/coldsmirk/vef-framework-go/approval.RollbackType = "none"
CONST RollbackPrevious : github.com/coldsmirk/vef-framework-go/approval.RollbackType = "previous"
CONST RollbackSpecified : github.com/coldsmirk/vef-framework-go/approval.RollbackType = "specified"
CONST RollbackStart : github.com/coldsmirk/vef-framework-go/approval.RollbackType = "start"
TYPE RollbackType : github.com/coldsmirk/vef-framework-go/approval.RollbackType
  METHOD IsValid : func() bool
TYPE SameApplicantAction : github.com/coldsmirk/vef-framework-go/approval.SameApplicantAction
  METHOD IsValid : func() bool
CONST SameApplicantAutoPass : github.com/coldsmirk/vef-framework-go/approval.SameApplicantAction = "auto_pass"
CONST SameApplicantSelfApprove : github.com/coldsmirk/vef-framework-go/approval.SameApplicantAction = "self_approve"
CONST SameApplicantTransferSuperior : github.com/coldsmirk/vef-framework-go/approval.SameApplicantAction = "transfer_superior"
TYPE StartNodeData : github.com/coldsmirk/vef-framework-go/approval.StartNodeData
  FIELD BaseNodeData : github.com/coldsmirk/vef-framework-go/approval.BaseNodeData [field_order=1 tag=""]
  FIELD Description : *string [promoted_from=BaseNodeData depth=1 field_order=2 tag="json:\"description,omitempty\""]
  FIELD Name : string [promoted_from=BaseNodeData depth=1 field_order=1 tag="json:\"name,omitempty\""]
  METHOD ApplyTo : func(node *github.com/coldsmirk/vef-framework-go/approval.FlowNode)
  METHOD GetDescription : func() *string
  METHOD GetName : func() string
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/approval.NodeKind
CONST StorageJSON : github.com/coldsmirk/vef-framework-go/approval.StorageMode = "json"
TYPE StorageMode : github.com/coldsmirk/vef-framework-go/approval.StorageMode
  METHOD IsValid : func() bool
CONST StorageTable : github.com/coldsmirk/vef-framework-go/approval.StorageMode = "table"
FUNC SubscribeInstance : func[T github.com/coldsmirk/vef-framework-go/approval.InstanceEvent](bus github.com/coldsmirk/vef-framework-go/event.Bus, handler func(ctx context.Context, evt T) error, opts ...github.com/coldsmirk/vef-framework-go/approval.InstanceSubscribeOption) (github.com/coldsmirk/vef-framework-go/event.Unsubscribe, error)
CONST SuperAdminRole : untyped string = "approval:super_admin"
VAR SystemCaller : github.com/coldsmirk/vef-framework-go/approval.CallerContext
TYPE Task : github.com/coldsmirk/vef-framework-go/approval.Task
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_task,alias:at\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD TenantID : string [field_order=3 tag="json:\"tenantId\" bun:\"tenant_id\""]
  FIELD InstanceID : string [field_order=4 tag="json:\"instanceId\" bun:\"instance_id\""]
  FIELD NodeID : string [field_order=5 tag="json:\"nodeId\" bun:\"node_id\""]
  FIELD VisitID : string [field_order=6 tag="json:\"visitId\" bun:\"visit_id\""]
  FIELD AssigneeID : string [field_order=7 tag="json:\"assigneeId\" bun:\"assignee_id\""]
  FIELD AssigneeName : string [field_order=8 tag="json:\"assigneeName\" bun:\"assignee_name\""]
  FIELD AssigneeDepartmentID : *string [field_order=9 tag="json:\"assigneeDepartmentId\" bun:\"assignee_department_id,nullzero\""]
  FIELD AssigneeDepartmentName : *string [field_order=10 tag="json:\"assigneeDepartmentName\" bun:\"assignee_department_name,nullzero\""]
  FIELD DelegatorID : *string [field_order=11 tag="json:\"delegatorId\" bun:\"delegator_id,nullzero\""]
  FIELD DelegatorName : *string [field_order=12 tag="json:\"delegatorName\" bun:\"delegator_name,nullzero\""]
  FIELD DelegatorDepartmentID : *string [field_order=13 tag="json:\"delegatorDepartmentId\" bun:\"delegator_department_id,nullzero\""]
  FIELD DelegatorDepartmentName : *string [field_order=14 tag="json:\"delegatorDepartmentName\" bun:\"delegator_department_name,nullzero\""]
  FIELD SortOrder : int [field_order=15 tag="json:\"sortOrder\" bun:\"sort_order\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/approval.TaskStatus [field_order=16 tag="json:\"status\" bun:\"status\""]
  FIELD ReadAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=17 tag="json:\"readAt\" bun:\"read_at,nullzero\""]
  FIELD ParentTaskID : *string [field_order=18 tag="json:\"parentTaskId\" bun:\"parent_task_id,nullzero\""]
  FIELD AddAssigneeType : *github.com/coldsmirk/vef-framework-go/approval.AddAssigneeType [field_order=19 tag="json:\"addAssigneeType\" bun:\"add_assignee_type,nullzero\""]
  FIELD Deadline : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=20 tag="json:\"deadline\" bun:\"deadline,nullzero\""]
  FIELD IsTimeout : bool [field_order=21 tag="json:\"isTimeout\" bun:\"is_timeout\""]
  FIELD IsPreWarningSent : bool [field_order=22 tag="json:\"isPreWarningSent\" bun:\"is_pre_warning_sent\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=23 tag="json:\"finishedAt\" bun:\"finished_at,nullzero\""]
  METHOD Assignee : func() github.com/coldsmirk/vef-framework-go/approval.UserInfo
  METHOD Delegator : func() *github.com/coldsmirk/vef-framework-go/approval.UserInfo
CONST TaskApproved : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "approved"
TYPE TaskApprovedEvent : github.com/coldsmirk/vef-framework-go/approval.TaskApprovedEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"operator\""]
  FIELD Opinion : *string [field_order=3 tag="json:\"opinion,omitempty\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST TaskCanceled : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "canceled"
TYPE TaskCanceledEvent : github.com/coldsmirk/vef-framework-go/approval.TaskCanceledEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Assignee : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"assignee\""]
  FIELD Reason : string [field_order=3 tag="json:\"reason\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE TaskCreatedEvent : github.com/coldsmirk/vef-framework-go/approval.TaskCreatedEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Assignee : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"assignee\""]
  FIELD Deadline : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=3 tag="json:\"deadline,omitempty\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE TaskDeadlineWarningEvent : github.com/coldsmirk/vef-framework-go/approval.TaskDeadlineWarningEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Assignee : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"assignee\""]
  FIELD Deadline : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=3 tag="json:\"deadline\""]
  FIELD HoursLeft : int [field_order=4 tag="json:\"hoursLeft\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [field_order=1 tag=""]
  FIELD TaskID : string [field_order=2 tag="json:\"taskId\""]
  FIELD NodeID : string [field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [field_order=4 tag="json:\"nodeName\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=InstanceEventBase depth=1 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=InstanceEventBase depth=1 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=InstanceEventBase depth=1 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=InstanceEventBase depth=1 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=InstanceEventBase depth=1 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=InstanceEventBase depth=1 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=InstanceEventBase depth=1 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=InstanceEventBase depth=1 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=InstanceEventBase depth=1 field_order=4 tag="json:\"title\""]
CONST TaskHandled : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "handled"
TYPE TaskHandledEvent : github.com/coldsmirk/vef-framework-go/approval.TaskHandledEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"operator\""]
  FIELD Opinion : *string [field_order=3 tag="json:\"opinion,omitempty\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE TaskNodeData : github.com/coldsmirk/vef-framework-go/approval.TaskNodeData
  FIELD Assignees : []github.com/coldsmirk/vef-framework-go/approval.AssigneeDefinition [field_order=1 tag="json:\"assignees,omitempty\""]
  FIELD ExecutionType : github.com/coldsmirk/vef-framework-go/approval.ExecutionType [field_order=2 tag="json:\"executionType,omitempty\""]
  FIELD EmptyAssigneeAction : github.com/coldsmirk/vef-framework-go/approval.EmptyAssigneeAction [field_order=3 tag="json:\"emptyAssigneeAction,omitempty\""]
  FIELD FallbackUserIDs : []string [field_order=4 tag="json:\"fallbackUserIds,omitempty\""]
  FIELD AdminUserIDs : []string [field_order=5 tag="json:\"adminUserIds,omitempty\""]
  FIELD IsTransferAllowed : *bool [field_order=6 tag="json:\"isTransferAllowed,omitempty\""]
  FIELD IsOpinionRequired : bool [field_order=7 tag="json:\"isOpinionRequired,omitempty\""]
  FIELD TimeoutHours : int [field_order=8 tag="json:\"timeoutHours,omitempty\""]
  FIELD TimeoutAction : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction [field_order=9 tag="json:\"timeoutAction,omitempty\""]
  FIELD TimeoutNotifyBeforeHours : int [field_order=10 tag="json:\"timeoutNotifyBeforeHours,omitempty\""]
  FIELD UrgeCooldownMinutes : int [field_order=11 tag="json:\"urgeCooldownMinutes,omitempty\""]
  FIELD CCs : []github.com/coldsmirk/vef-framework-go/approval.CCDefinition [field_order=12 tag="json:\"ccs,omitempty\""]
  FIELD FieldPermissions : map[string]github.com/coldsmirk/vef-framework-go/approval.Permission [field_order=13 tag="json:\"fieldPermissions,omitempty\""]
  METHOD GetAssignees : func() []github.com/coldsmirk/vef-framework-go/approval.AssigneeDefinition
  METHOD GetCCs : func() []github.com/coldsmirk/vef-framework-go/approval.CCDefinition
CONST TaskPending : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "pending"
TYPE TaskReassignedEvent : github.com/coldsmirk/vef-framework-go/approval.TaskReassignedEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD From : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"from\""]
  FIELD To : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=3 tag="json:\"to\""]
  FIELD Reason : *string [field_order=4 tag="json:\"reason,omitempty\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST TaskRejected : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "rejected"
TYPE TaskRejectedEvent : github.com/coldsmirk/vef-framework-go/approval.TaskRejectedEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"operator\""]
  FIELD Opinion : *string [field_order=3 tag="json:\"opinion,omitempty\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST TaskRemoved : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "removed"
CONST TaskRolledBack : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "rolled_back"
CONST TaskSkipped : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "skipped"
TYPE TaskStatus : github.com/coldsmirk/vef-framework-go/approval.TaskStatus
  METHOD IsFinal : func() bool
  METHOD String : func() string
TYPE TaskTimedOutEvent : github.com/coldsmirk/vef-framework-go/approval.TaskTimedOutEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Assignee : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"assignee\""]
  FIELD Deadline : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=3 tag="json:\"deadline\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST TaskTransferred : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "transferred"
TYPE TaskTransferredEvent : github.com/coldsmirk/vef-framework-go/approval.TaskTransferredEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD From : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"from\""]
  FIELD To : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=3 tag="json:\"to\""]
  FIELD Reason : *string [field_order=4 tag="json:\"reason,omitempty\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
TYPE TaskUrgedEvent : github.com/coldsmirk/vef-framework-go/approval.TaskUrgedEvent
  FIELD TaskEventBase : github.com/coldsmirk/vef-framework-go/approval.TaskEventBase [field_order=1 tag=""]
  FIELD Urger : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=2 tag="json:\"urger\""]
  FIELD Target : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=3 tag="json:\"target\""]
  FIELD Message : *string [field_order=4 tag="json:\"message,omitempty\""]
  FIELD InstanceEventBase : github.com/coldsmirk/vef-framework-go/approval.InstanceEventBase [promoted_from=TaskEventBase depth=1 field_order=1 tag=""]
  FIELD NodeID : string [promoted_from=TaskEventBase depth=1 field_order=3 tag="json:\"nodeId\""]
  FIELD NodeName : string [promoted_from=TaskEventBase depth=1 field_order=4 tag="json:\"nodeName\""]
  FIELD TaskID : string [promoted_from=TaskEventBase depth=1 field_order=2 tag="json:\"taskId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=8 tag="json:\"applicant\""]
  FIELD BusinessRef : *string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=7 tag="json:\"businessRef,omitempty\""]
  FIELD FlowCode : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=6 tag="json:\"flowCode\""]
  FIELD FlowID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=5 tag="json:\"flowId\""]
  FIELD InstanceID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=2 tag="json:\"instanceNo\""]
  FIELD OccurredTime : github.com/coldsmirk/vef-framework-go/timex.DateTime [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=9 tag="json:\"occurredTime\""]
  FIELD TenantID : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=3 tag="json:\"tenantId\""]
  FIELD Title : string [promoted_from=TaskEventBase.InstanceEventBase depth=2 field_order=4 tag="json:\"title\""]
  METHOD EventType : func() string
CONST TaskWaiting : github.com/coldsmirk/vef-framework-go/approval.TaskStatus = "waiting"
TYPE TimelineEntry : github.com/coldsmirk/vef-framework-go/approval.TimelineEntry
  FIELD Kind : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind [field_order=1 tag="json:\"kind\""]
  FIELD NodeID : *string [field_order=2 tag="json:\"nodeId,omitempty\""]
  FIELD Name : string [field_order=3 tag="json:\"name,omitempty\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/approval.NodeVisitStatus [field_order=4 tag="json:\"status,omitempty\""]
  FIELD ExecutionType : string [field_order=5 tag="json:\"executionType,omitempty\""]
  FIELD ApprovalMethod : string [field_order=6 tag="json:\"approvalMethod,omitempty\""]
  FIELD PassRule : string [field_order=7 tag="json:\"passRule,omitempty\""]
  FIELD PassRatio : *github.com/coldsmirk/vef-framework-go/decimal.Decimal [field_order=8 tag="json:\"passRatio,omitempty\""]
  FIELD Participants : []github.com/coldsmirk/vef-framework-go/approval.NodeParticipant [field_order=9 tag="json:\"participants,omitempty\""]
  FIELD CCRecipients : []github.com/coldsmirk/vef-framework-go/approval.CCRecipient [field_order=10 tag="json:\"ccRecipients,omitempty\""]
  FIELD Activities : []github.com/coldsmirk/vef-framework-go/approval.Activity [field_order=11 tag="json:\"activities,omitempty\""]
  FIELD StartedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=12 tag="json:\"startedAt\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=13 tag="json:\"finishedAt,omitempty\""]
CONST TimelineEntryApproval : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind = "approval"
CONST TimelineEntryCC : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind = "cc"
CONST TimelineEntryHandle : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind = "handle"
TYPE TimelineEntryKind : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind
CONST TimelineEntryStart : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind = "start"
CONST TimelineEntryTerminate : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind = "terminate"
CONST TimelineEntryWithdraw : github.com/coldsmirk/vef-framework-go/approval.TimelineEntryKind = "withdraw"
TYPE TimeoutAction : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction
  METHOD IsValid : func() bool
CONST TimeoutActionAutoPass : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction = "auto_pass"
CONST TimeoutActionAutoReject : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction = "auto_reject"
CONST TimeoutActionNone : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction = "none"
CONST TimeoutActionNotify : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction = "notify"
CONST TimeoutActionTransferAdmin : github.com/coldsmirk/vef-framework-go/approval.TimeoutAction = "transfer_admin"
TYPE UrgeRecord : github.com/coldsmirk/vef-framework-go/approval.UrgeRecord
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:apv_urge_record,alias:aur\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD CreationTrackedModel : github.com/coldsmirk/vef-framework-go/orm.CreationTrackedModel [field_order=3 tag=""]
  FIELD InstanceID : string [field_order=4 tag="json:\"instanceId\" bun:\"instance_id\""]
  FIELD NodeID : string [field_order=5 tag="json:\"nodeId\" bun:\"node_id\""]
  FIELD TaskID : *string [field_order=6 tag="json:\"taskId\" bun:\"task_id,nullzero\""]
  FIELD UrgerID : string [field_order=7 tag="json:\"urgerId\" bun:\"urger_id\""]
  FIELD UrgerName : string [field_order=8 tag="json:\"urgerName\" bun:\"urger_name\""]
  FIELD UrgerDepartmentID : *string [field_order=9 tag="json:\"urgerDepartmentId\" bun:\"urger_department_id,nullzero\""]
  FIELD UrgerDepartmentName : *string [field_order=10 tag="json:\"urgerDepartmentName\" bun:\"urger_department_name,nullzero\""]
  FIELD TargetUserID : string [field_order=11 tag="json:\"targetUserId\" bun:\"target_user_id\""]
  FIELD TargetUserName : string [field_order=12 tag="json:\"targetUserName\" bun:\"target_user_name\""]
  FIELD TargetUserDepartmentID : *string [field_order=13 tag="json:\"targetUserDepartmentId\" bun:\"target_user_department_id,nullzero\""]
  FIELD TargetUserDepartmentName : *string [field_order=14 tag="json:\"targetUserDepartmentName\" bun:\"target_user_department_name,nullzero\""]
  FIELD Message : string [field_order=15 tag="json:\"message\" bun:\"message\""]
  METHOD Target : func() github.com/coldsmirk/vef-framework-go/approval.UserInfo
  METHOD Urger : func() github.com/coldsmirk/vef-framework-go/approval.UserInfo
TYPE UserInfo : github.com/coldsmirk/vef-framework-go/approval.UserInfo
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD Name : string [field_order=2 tag="json:\"name\""]
  FIELD DepartmentID : *string [field_order=3 tag="json:\"departmentId,omitempty\""]
  FIELD DepartmentName : *string [field_order=4 tag="json:\"departmentName,omitempty\""]
  METHOD NewActionLog : func(instanceID string, action github.com/coldsmirk/vef-framework-go/approval.ActionType) *github.com/coldsmirk/vef-framework-go/approval.ActionLog
TYPE UserInfoResolver : github.com/coldsmirk/vef-framework-go/approval.UserInfoResolver
  METHOD ResolveUsers : func(ctx context.Context, userIDs []string) (map[string]github.com/coldsmirk/vef-framework-go/approval.UserInfo, error)
FUNC ValidateBusinessIdentifier : func(id string) error
TYPE ValidationRule : github.com/coldsmirk/vef-framework-go/approval.ValidationRule
  FIELD MinLength : *int [field_order=1 tag="json:\"minLength,omitempty\""]
  FIELD MaxLength : *int [field_order=2 tag="json:\"maxLength,omitempty\""]
  FIELD Min : *float64 [field_order=3 tag="json:\"min,omitempty\""]
  FIELD Max : *float64 [field_order=4 tag="json:\"max,omitempty\""]
  FIELD Pattern : string [field_order=5 tag="json:\"pattern,omitempty\""]
  FIELD Message : string [field_order=6 tag="json:\"message,omitempty\""]
CONST VersionArchived : github.com/coldsmirk/vef-framework-go/approval.VersionStatus = "archived"
CONST VersionDraft : github.com/coldsmirk/vef-framework-go/approval.VersionStatus = "draft"
CONST VersionPublished : github.com/coldsmirk/vef-framework-go/approval.VersionStatus = "published"
TYPE VersionStatus : github.com/coldsmirk/vef-framework-go/approval.VersionStatus
FUNC WithConcurrency : func(n int) github.com/coldsmirk/vef-framework-go/approval.InstanceSubscribeOption
FUNC WithGroup : func(name string) github.com/coldsmirk/vef-framework-go/approval.InstanceSubscribeOption

## github.com/coldsmirk/vef-framework-go/approval/admin
TYPE ActionLog : github.com/coldsmirk/vef-framework-go/approval/admin.ActionLog
  FIELD LogID : string [field_order=1 tag="json:\"logId\""]
  FIELD Action : string [field_order=2 tag="json:\"action\""]
  FIELD NodeID : *string [field_order=3 tag="json:\"nodeId,omitempty\""]
  FIELD TaskID : *string [field_order=4 tag="json:\"taskId,omitempty\""]
  FIELD Operator : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=5 tag="json:\"operator\""]
  FIELD TransferTo : *github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=6 tag="json:\"transferTo,omitempty\""]
  FIELD RollbackToNodeID : *string [field_order=7 tag="json:\"rollbackToNodeId,omitempty\""]
  FIELD AddedAssignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=8 tag="json:\"addedAssignees,omitempty\""]
  FIELD RemovedAssignees : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=9 tag="json:\"removedAssignees,omitempty\""]
  FIELD CCUsers : []github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=10 tag="json:\"ccUsers,omitempty\""]
  FIELD Opinion : *string [field_order=11 tag="json:\"opinion,omitempty\""]
  FIELD Attachments : []string [field_order=12 tag="json:\"attachments,omitempty\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=13 tag="json:\"createdAt\""]
TYPE Instance : github.com/coldsmirk/vef-framework-go/approval/admin.Instance
  FIELD InstanceID : string [field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [field_order=2 tag="json:\"instanceNo\""]
  FIELD Title : string [field_order=3 tag="json:\"title\""]
  FIELD TenantID : string [field_order=4 tag="json:\"tenantId\""]
  FIELD FlowID : string [field_order=5 tag="json:\"flowId\""]
  FIELD FlowName : string [field_order=6 tag="json:\"flowName\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=7 tag="json:\"applicant\""]
  FIELD Status : string [field_order=8 tag="json:\"status\""]
  FIELD CurrentNodeName : *string [field_order=9 tag="json:\"currentNodeName,omitempty\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=10 tag="json:\"createdAt\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=11 tag="json:\"finishedAt,omitempty\""]
TYPE InstanceDetail : github.com/coldsmirk/vef-framework-go/approval/admin.InstanceDetail
  FIELD Instance : github.com/coldsmirk/vef-framework-go/approval/admin.InstanceDetailInfo [field_order=1 tag="json:\"instance\""]
  FIELD FormSchema : encoding/json.RawMessage [field_order=2 tag="json:\"formSchema,omitempty\""]
  FIELD Timeline : []github.com/coldsmirk/vef-framework-go/approval.TimelineEntry [field_order=3 tag="json:\"timeline\""]
  FIELD FlowGraph : github.com/coldsmirk/vef-framework-go/approval.InstanceFlowGraph [field_order=4 tag="json:\"flowGraph\""]
TYPE InstanceDetailInfo : github.com/coldsmirk/vef-framework-go/approval/admin.InstanceDetailInfo
  FIELD InstanceID : string [field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [field_order=2 tag="json:\"instanceNo\""]
  FIELD Title : string [field_order=3 tag="json:\"title\""]
  FIELD TenantID : string [field_order=4 tag="json:\"tenantId\""]
  FIELD FlowID : string [field_order=5 tag="json:\"flowId\""]
  FIELD FlowName : string [field_order=6 tag="json:\"flowName\""]
  FIELD FlowVersionID : string [field_order=7 tag="json:\"flowVersionId\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=8 tag="json:\"applicant\""]
  FIELD Status : string [field_order=9 tag="json:\"status\""]
  FIELD CurrentNodeID : *string [field_order=10 tag="json:\"currentNodeId,omitempty\""]
  FIELD CurrentNodeName : *string [field_order=11 tag="json:\"currentNodeName,omitempty\""]
  FIELD BusinessRef : *string [field_order=12 tag="json:\"businessRef,omitempty\""]
  FIELD FormData : map[string]any [field_order=13 tag="json:\"formData,omitempty\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=14 tag="json:\"createdAt\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=15 tag="json:\"finishedAt,omitempty\""]
TYPE Metrics : github.com/coldsmirk/vef-framework-go/approval/admin.Metrics
  FIELD TenantID : string [field_order=1 tag="json:\"tenantId\""]
  FIELD CapturedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=2 tag="json:\"capturedAt\""]
  FIELD InstanceCounts : map[string]int [field_order=3 tag="json:\"instanceCounts\""]
  FIELD TaskCounts : map[string]int [field_order=4 tag="json:\"taskCounts\""]
  FIELD TimeoutTaskCount : int [field_order=5 tag="json:\"timeoutTaskCount\""]
  FIELD AvgCompletionSeconds : float64 [field_order=6 tag="json:\"avgCompletionSeconds\""]
  FIELD PendingBindingFailures : int [field_order=7 tag="json:\"pendingBindingFailures\""]
TYPE Task : github.com/coldsmirk/vef-framework-go/approval/admin.Task
  FIELD TaskID : string [field_order=1 tag="json:\"taskId\""]
  FIELD InstanceID : string [field_order=2 tag="json:\"instanceId\""]
  FIELD InstanceTitle : string [field_order=3 tag="json:\"instanceTitle\""]
  FIELD FlowName : string [field_order=4 tag="json:\"flowName\""]
  FIELD NodeName : string [field_order=5 tag="json:\"nodeName\""]
  FIELD Assignee : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=6 tag="json:\"assignee\""]
  FIELD Status : string [field_order=7 tag="json:\"status\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=8 tag="json:\"createdAt\""]
  FIELD Deadline : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=9 tag="json:\"deadline,omitempty\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=10 tag="json:\"finishedAt,omitempty\""]

## github.com/coldsmirk/vef-framework-go/approval/my
TYPE AvailableFlow : github.com/coldsmirk/vef-framework-go/approval/my.AvailableFlow
  FIELD FlowID : string [field_order=1 tag="json:\"flowId\""]
  FIELD FlowCode : string [field_order=2 tag="json:\"flowCode\""]
  FIELD FlowName : string [field_order=3 tag="json:\"flowName\""]
  FIELD FlowIcon : *string [field_order=4 tag="json:\"flowIcon,omitempty\""]
  FIELD Description : *string [field_order=5 tag="json:\"description,omitempty\""]
  FIELD CategoryID : string [field_order=6 tag="json:\"categoryId\""]
  FIELD CategoryName : string [field_order=7 tag="json:\"categoryName\""]
TYPE CCRecord : github.com/coldsmirk/vef-framework-go/approval/my.CCRecord
  FIELD CCRecordID : string [field_order=1 tag="json:\"ccRecordId\""]
  FIELD InstanceID : string [field_order=2 tag="json:\"instanceId\""]
  FIELD InstanceTitle : string [field_order=3 tag="json:\"instanceTitle\""]
  FIELD InstanceNo : string [field_order=4 tag="json:\"instanceNo\""]
  FIELD FlowName : string [field_order=5 tag="json:\"flowName\""]
  FIELD FlowIcon : *string [field_order=6 tag="json:\"flowIcon,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=7 tag="json:\"applicant\""]
  FIELD NodeName : *string [field_order=8 tag="json:\"nodeName,omitempty\""]
  FIELD IsRead : bool [field_order=9 tag="json:\"isRead\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=10 tag="json:\"createdAt\""]
TYPE CompletedTask : github.com/coldsmirk/vef-framework-go/approval/my.CompletedTask
  FIELD TaskID : string [field_order=1 tag="json:\"taskId\""]
  FIELD InstanceID : string [field_order=2 tag="json:\"instanceId\""]
  FIELD InstanceTitle : string [field_order=3 tag="json:\"instanceTitle\""]
  FIELD InstanceNo : string [field_order=4 tag="json:\"instanceNo\""]
  FIELD FlowName : string [field_order=5 tag="json:\"flowName\""]
  FIELD FlowIcon : *string [field_order=6 tag="json:\"flowIcon,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=7 tag="json:\"applicant\""]
  FIELD NodeName : string [field_order=8 tag="json:\"nodeName\""]
  FIELD Status : string [field_order=9 tag="json:\"status\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=10 tag="json:\"finishedAt,omitempty\""]
TYPE InitiatedInstance : github.com/coldsmirk/vef-framework-go/approval/my.InitiatedInstance
  FIELD InstanceID : string [field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [field_order=2 tag="json:\"instanceNo\""]
  FIELD Title : string [field_order=3 tag="json:\"title\""]
  FIELD FlowName : string [field_order=4 tag="json:\"flowName\""]
  FIELD FlowIcon : *string [field_order=5 tag="json:\"flowIcon,omitempty\""]
  FIELD Status : string [field_order=6 tag="json:\"status\""]
  FIELD CurrentNodeName : *string [field_order=7 tag="json:\"currentNodeName,omitempty\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=8 tag="json:\"createdAt\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=9 tag="json:\"finishedAt,omitempty\""]
TYPE InstanceDetail : github.com/coldsmirk/vef-framework-go/approval/my.InstanceDetail
  FIELD Instance : github.com/coldsmirk/vef-framework-go/approval/my.InstanceInfo [field_order=1 tag="json:\"instance\""]
  FIELD FormSchema : encoding/json.RawMessage [field_order=2 tag="json:\"formSchema,omitempty\""]
  FIELD Timeline : []github.com/coldsmirk/vef-framework-go/approval.TimelineEntry [field_order=3 tag="json:\"timeline\""]
  FIELD FlowGraph : github.com/coldsmirk/vef-framework-go/approval.InstanceFlowGraph [field_order=4 tag="json:\"flowGraph\""]
  FIELD AvailableActions : []string [field_order=5 tag="json:\"availableActions\""]
  FIELD FieldPermissions : map[string]github.com/coldsmirk/vef-framework-go/approval.Permission [field_order=6 tag="json:\"fieldPermissions,omitempty\""]
TYPE InstanceInfo : github.com/coldsmirk/vef-framework-go/approval/my.InstanceInfo
  FIELD InstanceID : string [field_order=1 tag="json:\"instanceId\""]
  FIELD InstanceNo : string [field_order=2 tag="json:\"instanceNo\""]
  FIELD Title : string [field_order=3 tag="json:\"title\""]
  FIELD FlowName : string [field_order=4 tag="json:\"flowName\""]
  FIELD FlowIcon : *string [field_order=5 tag="json:\"flowIcon,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=6 tag="json:\"applicant\""]
  FIELD Status : string [field_order=7 tag="json:\"status\""]
  FIELD CurrentNodeID : *string [field_order=8 tag="json:\"currentNodeId,omitempty\""]
  FIELD CurrentNodeName : *string [field_order=9 tag="json:\"currentNodeName,omitempty\""]
  FIELD BusinessRef : *string [field_order=10 tag="json:\"businessRef,omitempty\""]
  FIELD FormData : map[string]any [field_order=11 tag="json:\"formData,omitempty\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=12 tag="json:\"createdAt\""]
  FIELD FinishedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=13 tag="json:\"finishedAt,omitempty\""]
TYPE PendingCounts : github.com/coldsmirk/vef-framework-go/approval/my.PendingCounts
  FIELD PendingTaskCount : int [field_order=1 tag="json:\"pendingTaskCount\""]
  FIELD UnreadCCCount : int [field_order=2 tag="json:\"unreadCcCount\""]
TYPE PendingTask : github.com/coldsmirk/vef-framework-go/approval/my.PendingTask
  FIELD TaskID : string [field_order=1 tag="json:\"taskId\""]
  FIELD InstanceID : string [field_order=2 tag="json:\"instanceId\""]
  FIELD InstanceTitle : string [field_order=3 tag="json:\"instanceTitle\""]
  FIELD InstanceNo : string [field_order=4 tag="json:\"instanceNo\""]
  FIELD FlowName : string [field_order=5 tag="json:\"flowName\""]
  FIELD FlowIcon : *string [field_order=6 tag="json:\"flowIcon,omitempty\""]
  FIELD Applicant : github.com/coldsmirk/vef-framework-go/approval.UserInfo [field_order=7 tag="json:\"applicant\""]
  FIELD NodeName : string [field_order=8 tag="json:\"nodeName\""]
  FIELD CreatedAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=9 tag="json:\"createdAt\""]
  FIELD Deadline : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=10 tag="json:\"deadline,omitempty\""]
  FIELD IsTimeout : bool [field_order=11 tag="json:\"isTimeout\""]

## github.com/coldsmirk/vef-framework-go/cache
TYPE Cache : github.com/coldsmirk/vef-framework-go/cache.Cache[T any]
  METHOD Clear : func(ctx context.Context) error
  METHOD Close : func() error
  METHOD Contains : func(ctx context.Context, key string) bool
  METHOD Delete : func(ctx context.Context, key string) error
  METHOD ForEach : func(ctx context.Context, callback func(key string, value T) bool, prefix ...string) error
  METHOD Get : func(ctx context.Context, key string) (T, bool)
  METHOD GetOrLoad : func(ctx context.Context, key string, loader github.com/coldsmirk/vef-framework-go/cache.LoaderFunc[T], ttl ...time.Duration) (T, error)
  METHOD Keys : func(ctx context.Context, prefix ...string) ([]string, error)
  METHOD Set : func(ctx context.Context, key string, value T, ttl ...time.Duration) error
  METHOD Size : func(ctx context.Context) (int64, error)
VAR ErrCacheClosed : error
VAR ErrLoaderRequired : error
VAR ErrMemoryLimitExceeded : error
VAR ErrTypeAssertionFailed : error
TYPE EvictionPolicy : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy
CONST EvictionPolicyFIFO : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 3
CONST EvictionPolicyLFU : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 2
CONST EvictionPolicyLRU : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 1
CONST EvictionPolicyNone : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 0
TYPE GetFunc : github.com/coldsmirk/vef-framework-go/cache.GetFunc[T any]
TYPE Invalidating : github.com/coldsmirk/vef-framework-go/cache.Invalidating[T any]
  METHOD Get : func(ctx context.Context, key string) (T, error)
  METHOD Invalidate : func(ctx context.Context, keys ...string) error
FUNC Key : func(keyParts ...string) string
TYPE KeyBuilder : github.com/coldsmirk/vef-framework-go/cache.KeyBuilder
  METHOD Build : func(keyParts ...string) string
TYPE KeyedLoaderFunc : github.com/coldsmirk/vef-framework-go/cache.KeyedLoaderFunc[T any]
TYPE LoaderFunc : github.com/coldsmirk/vef-framework-go/cache.LoaderFunc[T any]
TYPE MemoryOption : github.com/coldsmirk/vef-framework-go/cache.MemoryOption
FUNC NewInvalidating : func[T any](loader github.com/coldsmirk/vef-framework-go/cache.KeyedLoaderFunc[T], logger github.com/coldsmirk/vef-framework-go/logx.Logger, opts ...github.com/coldsmirk/vef-framework-go/cache.MemoryOption) *github.com/coldsmirk/vef-framework-go/cache.Invalidating[T]
FUNC NewMemory : func[T any](opts ...github.com/coldsmirk/vef-framework-go/cache.MemoryOption) github.com/coldsmirk/vef-framework-go/cache.Cache[T]
FUNC NewPrefixKeyBuilder : func(prefix string) *github.com/coldsmirk/vef-framework-go/cache.PrefixKeyBuilder
FUNC NewPrefixKeyBuilderWithSeparator : func(prefix string, separator string) *github.com/coldsmirk/vef-framework-go/cache.PrefixKeyBuilder
FUNC NewRedis : func[T any](client *github.com/redis/go-redis/v9.Client, namespace string, opts ...github.com/coldsmirk/vef-framework-go/cache.RedisOption) github.com/coldsmirk/vef-framework-go/cache.Cache[T]
TYPE PrefixKeyBuilder : github.com/coldsmirk/vef-framework-go/cache.PrefixKeyBuilder
  METHOD Build : func(keyParts ...string) string
TYPE RedisOption : github.com/coldsmirk/vef-framework-go/cache.RedisOption
TYPE SetFunc : github.com/coldsmirk/vef-framework-go/cache.SetFunc[T any]
TYPE SingleflightMixin : github.com/coldsmirk/vef-framework-go/cache.SingleflightMixin[T any]
  METHOD GetOrLoad : func(ctx context.Context, cacheKey string, loader github.com/coldsmirk/vef-framework-go/cache.LoaderFunc[T], ttl []time.Duration, getFn github.com/coldsmirk/vef-framework-go/cache.GetFunc[T], setFn github.com/coldsmirk/vef-framework-go/cache.SetFunc[T]) (value T, _ error)
FUNC WithMemDefaultTTL : func(ttl time.Duration) github.com/coldsmirk/vef-framework-go/cache.MemoryOption
FUNC WithMemEvictionPolicy : func(policy github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy) github.com/coldsmirk/vef-framework-go/cache.MemoryOption
FUNC WithMemGCInterval : func(interval time.Duration) github.com/coldsmirk/vef-framework-go/cache.MemoryOption
FUNC WithMemMaxSize : func(size int64) github.com/coldsmirk/vef-framework-go/cache.MemoryOption
FUNC WithRdsDefaultTTL : func(ttl time.Duration) github.com/coldsmirk/vef-framework-go/cache.RedisOption

## github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd
CONST Banner : untyped string = "\n██╗   ██╗███████╗███████╗     ██████╗██╗     ██╗\n██║   ██║██╔════╝██╔════╝    ██╔════╝██║     ██║\n██║   ██║█████╗  █████╗      ██║     ██║     ██║\n╚██╗ ██╔╝██╔══╝  ██╔══╝      ██║     ██║     ██║\n ╚████╔╝ ███████╗██║         ╚██████╗███████╗██║\n  ╚═══╝  ╚══════╝╚═╝          ╚═════╝╚══════╝╚═╝\n"
FUNC Execute : func() error
FUNC GetVersionInfo : func(ldflagsVersion string, ldflagsDate string) github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd.VersionInfo
FUNC Init : func(ldflagsVersion string, ldflagsDate string)
FUNC PrintBanner : func()
TYPE VersionInfo : github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd.VersionInfo
  FIELD Version : string [field_order=1 tag=""]
  FIELD Date : string [field_order=2 tag=""]
  FIELD Dirty : bool [field_order=3 tag=""]
  METHOD String : func() string

## github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/buildinfo
FUNC Command : func() *github.com/spf13/cobra.Command
FUNC Generate : func(outputPath string, packageName string) error

## github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/create
FUNC Command : func() *github.com/spf13/cobra.Command
VAR ErrNotImplemented : error

## github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/modelschema
FUNC Command : func() *github.com/spf13/cobra.Command
VAR ErrFileNotFoundInPackage : error
VAR ErrMultiplePackages : error
VAR ErrNoGoFilesFound : error
VAR ErrNoPackagesFound : error
FUNC GenerateDirectory : func(inputDir string, outputDir string, packageName string) error
FUNC GenerateFile : func(inputFile string, outputFile string, packageName string) error
TYPE ModelField : github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/modelschema.ModelField
  FIELD GoName : string [field_order=1 tag=""]
  FIELD ColumnName : string [field_order=2 tag=""]
  FIELD MethodName : string [field_order=3 tag=""]
  FIELD Label : string [field_order=4 tag=""]
  FIELD Scanonly : bool [field_order=5 tag=""]
TYPE ModelSchemaInfo : github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/modelschema.ModelSchemaInfo
  FIELD PackageName : string [field_order=1 tag=""]
  FIELD ModelName : string [field_order=2 tag=""]
  FIELD SchemaTypeName : string [field_order=3 tag=""]
  FIELD VarName : string [field_order=4 tag=""]
  FIELD TableName : string [field_order=5 tag=""]
  FIELD AliasName : string [field_order=6 tag=""]
  FIELD Fields : []github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/modelschema.ModelField [field_order=7 tag=""]

## github.com/coldsmirk/vef-framework-go/config
TYPE AppConfig : github.com/coldsmirk/vef-framework-go/config.AppConfig
  FIELD Name : string [field_order=1 tag="config:\"name\""]
  FIELD Port : uint16 [field_order=2 tag="config:\"port\""]
  FIELD BodyLimit : string [field_order=3 tag="config:\"body_limit\""]
  FIELD TrustedProxies : []string [field_order=4 tag="config:\"trusted_proxies\""]
TYPE ApprovalConfig : github.com/coldsmirk/vef-framework-go/config.ApprovalConfig
  FIELD AutoMigrate : bool [field_order=1 tag="config:\"auto_migrate\""]
  FIELD TimeoutScanInterval : time.Duration [field_order=2 tag="config:\"timeout_scan_interval\""]
  FIELD PreWarningScanInterval : time.Duration [field_order=3 tag="config:\"pre_warning_scan_interval\""]
  FIELD CleanupScanInterval : time.Duration [field_order=4 tag="config:\"cleanup_scan_interval\""]
  FIELD DelegationMaxDepth : int [field_order=5 tag="config:\"delegation_max_depth\""]
  FIELD FormSnapshotRetention : time.Duration [field_order=6 tag="config:\"form_snapshot_retention\""]
  FIELD UrgeRecordRetention : time.Duration [field_order=7 tag="config:\"urge_record_retention\""]
  FIELD CCRecordRetention : time.Duration [field_order=8 tag="config:\"cc_record_retention\""]
  METHOD ApplyDefaults : func()
TYPE CORSConfig : github.com/coldsmirk/vef-framework-go/config.CORSConfig
  FIELD Enabled : bool [field_order=1 tag="config:\"enabled\""]
  FIELD AllowOrigins : []string [field_order=2 tag="config:\"allow_origins\""]
TYPE Config : github.com/coldsmirk/vef-framework-go/config.Config
  METHOD Unmarshal : func(key string, target any) error
TYPE DBKind : github.com/coldsmirk/vef-framework-go/config.DBKind
TYPE DataSourceConfig : github.com/coldsmirk/vef-framework-go/config.DataSourceConfig
  FIELD Kind : github.com/coldsmirk/vef-framework-go/config.DBKind [field_order=1 tag="config:\"type\""]
  FIELD Host : string [field_order=2 tag="config:\"host\""]
  FIELD Port : uint16 [field_order=3 tag="config:\"port\""]
  FIELD User : string [field_order=4 tag="config:\"user\""]
  FIELD Password : string [field_order=5 tag="config:\"password\""]
  FIELD Database : string [field_order=6 tag="config:\"database\""]
  FIELD Schema : string [field_order=7 tag="config:\"schema\""]
  FIELD Path : string [field_order=8 tag="config:\"path\""]
  FIELD EnableSQLGuard : bool [field_order=9 tag="config:\"enable_sql_guard\""]
  FIELD SSLMode : github.com/coldsmirk/vef-framework-go/config.SSLMode [field_order=10 tag="config:\"ssl_mode\""]
  FIELD SSLRootCert : string [field_order=11 tag="config:\"ssl_root_cert\""]
TYPE DataSourcesConfig : github.com/coldsmirk/vef-framework-go/config.DataSourcesConfig
  FIELD Map : map[string]github.com/coldsmirk/vef-framework-go/config.DataSourceConfig [field_order=1 tag=""]
  METHOD Primary : func() github.com/coldsmirk/vef-framework-go/config.DataSourceConfig
CONST DefaultClaimTTL : time.Duration = 86400000000000
CONST DefaultDeleteBatchSize : int = 100
CONST DefaultDeleteConcurrency : int = 8
CONST DefaultDeleteLeaseWindow : time.Duration = 300000000000
CONST DefaultDeleteMaxAttempts : int = 12
CONST DefaultDeleteWorkerInterval : time.Duration = 300000000000
CONST DefaultLockoutBackoffBase : time.Duration = 1000000000
CONST DefaultLockoutBackoffMax : time.Duration = 900000000000
CONST DefaultLockoutLockDuration : time.Duration = 900000000000
CONST DefaultLockoutMaxFailures : untyped int = 10
CONST DefaultLockoutWindow : time.Duration = 900000000000
CONST DefaultMaxPendingClaims : int = 100
CONST DefaultMaxUploadSize : int64 = 1073741824
CONST DefaultSessionIdleTTL : time.Duration = 1800000000000
CONST DefaultSessionMaxLifetime : time.Duration = 604800000000000
CONST DefaultSweepBatchSize : int = 200
CONST DefaultSweepInterval : time.Duration = 300000000000
CONST EnvConfigPath : untyped string = "VEF_CONFIG_PATH"
CONST EnvI18NLanguage : untyped string = "VEF_I18N_LANGUAGE"
CONST EnvLogLevel : untyped string = "VEF_LOG_LEVEL"
CONST EnvPrefix : untyped string = "VEF"
VAR ErrInboxRetentionTooShort : error
VAR ErrInvalidLockoutKey : error
VAR ErrInvalidLockoutStrategy : error
VAR ErrInvalidSessionOnExceed : error
VAR ErrInvalidTokenType : error
TYPE EventConfig : github.com/coldsmirk/vef-framework-go/config.EventConfig
  FIELD DefaultTransport : string [field_order=1 tag="config:\"default_transport\""]
  FIELD AsyncQueueSize : int [field_order=2 tag="config:\"async_queue_size\""]
  FIELD AsyncWorkers : int [field_order=3 tag="config:\"async_workers\""]
  FIELD PublishTimeout : time.Duration [field_order=4 tag="config:\"publish_timeout\""]
  FIELD Transports : github.com/coldsmirk/vef-framework-go/config.EventTransportsConfig [field_order=5 tag="config:\"transports\""]
  FIELD Middleware : github.com/coldsmirk/vef-framework-go/config.EventMiddlewareConfig [field_order=6 tag="config:\"middleware\""]
  FIELD Inbox : github.com/coldsmirk/vef-framework-go/config.EventInboxConfig [field_order=7 tag="config:\"inbox\""]
  FIELD Routing : []github.com/coldsmirk/vef-framework-go/config.EventRoutingRule [field_order=8 tag="config:\"routing\""]
  METHOD EffectiveAsyncQueueSize : func() int
  METHOD EffectiveAsyncWorkers : func() int
  METHOD EffectiveDefaultTransport : func() string
  METHOD EffectivePublishTimeout : func() time.Duration
  METHOD Validate : func() error
TYPE EventInboxConfig : github.com/coldsmirk/vef-framework-go/config.EventInboxConfig
  FIELD Retention : time.Duration [field_order=1 tag="config:\"retention\""]
  FIELD ProcessingLease : time.Duration [field_order=2 tag="config:\"processing_lease\""]
  FIELD CleanupInterval : time.Duration [field_order=3 tag="config:\"cleanup_interval\""]
  METHOD EffectiveCleanupInterval : func() time.Duration
  METHOD EffectiveProcessingLease : func() time.Duration
  METHOD EffectiveRetention : func() time.Duration
TYPE EventMemoryTransportConfig : github.com/coldsmirk/vef-framework-go/config.EventMemoryTransportConfig
  FIELD QueueSize : int [field_order=1 tag="config:\"queue_size\""]
  FIELD FullPolicy : string [field_order=2 tag="config:\"full_policy\""]
  FIELD PublishTimeout : time.Duration [field_order=3 tag="config:\"publish_timeout\""]
TYPE EventMiddlewareConfig : github.com/coldsmirk/vef-framework-go/config.EventMiddlewareConfig
  FIELD Logging : bool [field_order=1 tag="config:\"logging\""]
  FIELD Tracing : bool [field_order=2 tag="config:\"tracing\""]
  FIELD TracingStrict : bool [field_order=3 tag="config:\"tracing_strict\""]
  FIELD Metrics : bool [field_order=4 tag="config:\"metrics\""]
  FIELD Recover : bool [field_order=5 tag="config:\"recover\""]
  FIELD Inbox : bool [field_order=6 tag="config:\"inbox\""]
TYPE EventOutboxTransportConfig : github.com/coldsmirk/vef-framework-go/config.EventOutboxTransportConfig
  FIELD Enabled : bool [field_order=1 tag="config:\"enabled\""]
  FIELD RelayInterval : time.Duration [field_order=2 tag="config:\"relay_interval\""]
  FIELD MaxRetries : int [field_order=3 tag="config:\"max_retries\""]
  FIELD BatchSize : int [field_order=4 tag="config:\"batch_size\""]
  FIELD LeaseMultiplier : int [field_order=5 tag="config:\"lease_multiplier\""]
  FIELD MinLease : time.Duration [field_order=6 tag="config:\"min_lease\""]
  FIELD SinkName : string [field_order=7 tag="config:\"sink\""]
  FIELD CleanupInterval : time.Duration [field_order=8 tag="config:\"cleanup_interval\""]
  FIELD CompletedTTL : time.Duration [field_order=9 tag="config:\"completed_ttl\""]
  METHOD EffectiveCleanupInterval : func() time.Duration
  METHOD EffectiveCompletedTTL : func() time.Duration
TYPE EventRedisStreamTransportConfig : github.com/coldsmirk/vef-framework-go/config.EventRedisStreamTransportConfig
  FIELD Enabled : bool [field_order=1 tag="config:\"enabled\""]
  FIELD StreamPrefix : string [field_order=2 tag="config:\"stream_prefix\""]
  FIELD MaxLenApprox : int64 [field_order=3 tag="config:\"max_len_approx\""]
  FIELD BlockTimeout : time.Duration [field_order=4 tag="config:\"block_timeout\""]
  FIELD ClaimIdle : time.Duration [field_order=5 tag="config:\"claim_idle\""]
  FIELD ClaimInterval : time.Duration [field_order=6 tag="config:\"claim_interval\""]
  FIELD ClaimBatchSize : int64 [field_order=7 tag="config:\"claim_batch_size\""]
  FIELD ReaperConcurrency : int [field_order=8 tag="config:\"reaper_concurrency\""]
  FIELD HandlerTimeout : time.Duration [field_order=9 tag="config:\"handler_timeout\""]
  FIELD SetupTimeout : time.Duration [field_order=10 tag="config:\"setup_timeout\""]
  FIELD ConsumerID : string [field_order=11 tag="config:\"consumer_id\""]
  FIELD StartID : string [field_order=12 tag="config:\"start_id\""]
  FIELD IdleGroupRetention : time.Duration [field_order=13 tag="config:\"idle_group_retention\""]
  FIELD IdleGroupSweepInterval : time.Duration [field_order=14 tag="config:\"idle_group_sweep_interval\""]
TYPE EventRoutingRule : github.com/coldsmirk/vef-framework-go/config.EventRoutingRule
  FIELD Pattern : string [field_order=1 tag="config:\"pattern\""]
  FIELD Transports : []string [field_order=2 tag="config:\"transports\""]
TYPE EventTransportsConfig : github.com/coldsmirk/vef-framework-go/config.EventTransportsConfig
  FIELD Memory : github.com/coldsmirk/vef-framework-go/config.EventMemoryTransportConfig [field_order=1 tag="config:\"memory\""]
  FIELD Outbox : github.com/coldsmirk/vef-framework-go/config.EventOutboxTransportConfig [field_order=2 tag="config:\"outbox\""]
  FIELD RedisStream : github.com/coldsmirk/vef-framework-go/config.EventRedisStreamTransportConfig [field_order=3 tag="config:\"redis_stream\""]
TYPE FilesystemConfig : github.com/coldsmirk/vef-framework-go/config.FilesystemConfig
  FIELD Root : string [field_order=1 tag="config:\"root\""]
TYPE LockoutConfig : github.com/coldsmirk/vef-framework-go/config.LockoutConfig
  FIELD Enabled : *bool [field_order=1 tag="config:\"enabled\""]
  FIELD MaxFailures : int [field_order=2 tag="config:\"max_failures\""]
  FIELD Window : time.Duration [field_order=3 tag="config:\"window\""]
  FIELD LockDuration : time.Duration [field_order=4 tag="config:\"lock_duration\""]
  FIELD Strategy : github.com/coldsmirk/vef-framework-go/config.LockoutStrategy [field_order=5 tag="config:\"strategy\""]
  FIELD BackoffBase : time.Duration [field_order=6 tag="config:\"backoff_base\""]
  FIELD BackoffMax : time.Duration [field_order=7 tag="config:\"backoff_max\""]
  FIELD Key : github.com/coldsmirk/vef-framework-go/config.LockoutKey [field_order=8 tag="config:\"key\""]
  METHOD EffectiveBackoffBase : func() time.Duration
  METHOD EffectiveBackoffMax : func() time.Duration
  METHOD EffectiveKey : func() github.com/coldsmirk/vef-framework-go/config.LockoutKey
  METHOD EffectiveLockDuration : func() time.Duration
  METHOD EffectiveMaxFailures : func() int
  METHOD EffectiveStrategy : func() github.com/coldsmirk/vef-framework-go/config.LockoutStrategy
  METHOD EffectiveWindow : func() time.Duration
  METHOD IsEnabled : func() bool
  METHOD Validate : func() error
TYPE LockoutKey : github.com/coldsmirk/vef-framework-go/config.LockoutKey
CONST LockoutKeyIP : github.com/coldsmirk/vef-framework-go/config.LockoutKey = "ip"
CONST LockoutKeyUser : github.com/coldsmirk/vef-framework-go/config.LockoutKey = "user"
CONST LockoutKeyUserIP : github.com/coldsmirk/vef-framework-go/config.LockoutKey = "user_ip"
TYPE LockoutStrategy : github.com/coldsmirk/vef-framework-go/config.LockoutStrategy
CONST LockoutStrategyBackoff : github.com/coldsmirk/vef-framework-go/config.LockoutStrategy = "backoff"
CONST LockoutStrategyLock : github.com/coldsmirk/vef-framework-go/config.LockoutStrategy = "lock"
TYPE MCPConfig : github.com/coldsmirk/vef-framework-go/config.MCPConfig
  FIELD Enabled : bool [field_order=1 tag="config:\"enabled\""]
  FIELD RequireAuth : *bool [field_order=2 tag="config:\"require_auth\""]
TYPE MinIOConfig : github.com/coldsmirk/vef-framework-go/config.MinIOConfig
  FIELD Endpoint : string [field_order=1 tag="config:\"endpoint\""]
  FIELD AccessKey : string [field_order=2 tag="config:\"access_key\""]
  FIELD SecretKey : string [field_order=3 tag="config:\"secret_key\""]
  FIELD Bucket : string [field_order=4 tag="config:\"bucket\""]
  FIELD Region : string [field_order=5 tag="config:\"region\""]
  FIELD UseSSL : bool [field_order=6 tag="config:\"use_ssl\""]
TYPE MonitorConfig : github.com/coldsmirk/vef-framework-go/config.MonitorConfig
  FIELD SampleInterval : time.Duration [field_order=1 tag="config:\"sample_interval\""]
  FIELD SampleDuration : time.Duration [field_order=2 tag="config:\"sample_duration\""]
  FIELD ExcludedMounts : []string [field_order=3 tag="config:\"excluded_mounts\""]
CONST MySQL : github.com/coldsmirk/vef-framework-go/config.DBKind = "mysql"
CONST Oracle : github.com/coldsmirk/vef-framework-go/config.DBKind = "oracle"
TYPE PasswordPolicyConfig : github.com/coldsmirk/vef-framework-go/config.PasswordPolicyConfig
  FIELD MinLength : int [field_order=1 tag="config:\"min_length\""]
  FIELD MaxLength : int [field_order=2 tag="config:\"max_length\""]
  FIELD RequireUpper : bool [field_order=3 tag="config:\"require_upper\""]
  FIELD RequireLower : bool [field_order=4 tag="config:\"require_lower\""]
  FIELD RequireDigit : bool [field_order=5 tag="config:\"require_digit\""]
  FIELD RequireSymbol : bool [field_order=6 tag="config:\"require_symbol\""]
  FIELD MinCharClasses : int [field_order=7 tag="config:\"min_char_classes\""]
  FIELD DisallowUsername : bool [field_order=8 tag="config:\"disallow_username\""]
  FIELD Blocklist : []string [field_order=9 tag="config:\"blocklist\""]
  FIELD HistoryDepth : int [field_order=10 tag="config:\"history_depth\""]
  FIELD MaxAge : time.Duration [field_order=11 tag="config:\"max_age\""]
CONST Postgres : github.com/coldsmirk/vef-framework-go/config.DBKind = "postgres"
CONST PrimaryDataSourceName : untyped string = "primary"
TYPE RedisConfig : github.com/coldsmirk/vef-framework-go/config.RedisConfig
  FIELD Enabled : bool [field_order=1 tag="config:\"enabled\""]
  FIELD Host : string [field_order=2 tag="config:\"host\""]
  FIELD Port : uint16 [field_order=3 tag="config:\"port\""]
  FIELD User : string [field_order=4 tag="config:\"user\""]
  FIELD Password : string [field_order=5 tag="config:\"password\""]
  FIELD Database : uint8 [field_order=6 tag="config:\"database\""]
  FIELD Network : string [field_order=7 tag="config:\"network\""]
CONST SQLServer : github.com/coldsmirk/vef-framework-go/config.DBKind = "sqlserver"
CONST SQLite : github.com/coldsmirk/vef-framework-go/config.DBKind = "sqlite"
CONST SSLDisable : github.com/coldsmirk/vef-framework-go/config.SSLMode = "disable"
TYPE SSLMode : github.com/coldsmirk/vef-framework-go/config.SSLMode
CONST SSLRequire : github.com/coldsmirk/vef-framework-go/config.SSLMode = "require"
CONST SSLVerifyCA : github.com/coldsmirk/vef-framework-go/config.SSLMode = "verify-ca"
CONST SSLVerifyFull : github.com/coldsmirk/vef-framework-go/config.SSLMode = "verify-full"
TYPE SecurityConfig : github.com/coldsmirk/vef-framework-go/config.SecurityConfig
  FIELD Secret : string [field_order=1 tag="config:\"secret\""]
  FIELD TokenExpires : time.Duration [field_order=2 tag="config:\"token_expires\""]
  FIELD RefreshNotBefore : time.Duration [field_order=3 tag="config:\"refresh_not_before\""]
  FIELD LoginRateLimit : int [field_order=4 tag="config:\"login_rate_limit\""]
  FIELD RefreshRateLimit : int [field_order=5 tag="config:\"refresh_rate_limit\""]
  FIELD IPWhitelists : map[string][]string [field_order=6 tag="config:\"ip_whitelists\""]
  FIELD Lockout : github.com/coldsmirk/vef-framework-go/config.LockoutConfig [field_order=7 tag="config:\"lockout\""]
  FIELD PasswordPolicy : github.com/coldsmirk/vef-framework-go/config.PasswordPolicyConfig [field_order=8 tag="config:\"password_policy\""]
  FIELD TokenType : github.com/coldsmirk/vef-framework-go/config.TokenType [field_order=9 tag="config:\"token_type\""]
  FIELD Session : github.com/coldsmirk/vef-framework-go/config.SessionConfig [field_order=10 tag="config:\"session\""]
  METHOD EffectiveTokenType : func() github.com/coldsmirk/vef-framework-go/config.TokenType
  METHOD Validate : func() error
TYPE SessionConfig : github.com/coldsmirk/vef-framework-go/config.SessionConfig
  FIELD MaxConcurrent : int [field_order=1 tag="config:\"max_concurrent\""]
  FIELD OnExceed : github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy [field_order=2 tag="config:\"on_exceed\""]
  FIELD IdleTTL : time.Duration [field_order=3 tag="config:\"idle_ttl\""]
  FIELD MaxLifetime : time.Duration [field_order=4 tag="config:\"max_lifetime\""]
  FIELD Sliding : *bool [field_order=5 tag="config:\"sliding\""]
  METHOD EffectiveIdleTTL : func() time.Duration
  METHOD EffectiveMaxLifetime : func() time.Duration
  METHOD EffectiveOnExceed : func() github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy
  METHOD IsSliding : func() bool
CONST SessionExceedEvictOldest : github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy = "evict_oldest"
TYPE SessionExceedPolicy : github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy
CONST SessionExceedReject : github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy = "reject"
TYPE StorageConfig : github.com/coldsmirk/vef-framework-go/config.StorageConfig
  FIELD Provider : github.com/coldsmirk/vef-framework-go/config.StorageProvider [field_order=1 tag="config:\"provider\""]
  FIELD AutoMigrate : bool [field_order=2 tag="config:\"auto_migrate\""]
  FIELD MinIO : github.com/coldsmirk/vef-framework-go/config.MinIOConfig [field_order=3 tag="config:\"minio\""]
  FIELD Filesystem : github.com/coldsmirk/vef-framework-go/config.FilesystemConfig [field_order=4 tag="config:\"filesystem\""]
  FIELD MaxUploadSize : int64 [field_order=5 tag="config:\"max_upload_size\""]
  FIELD ClaimTTL : time.Duration [field_order=6 tag="config:\"claim_ttl\""]
  FIELD MaxPendingClaims : int [field_order=7 tag="config:\"max_pending_claims\""]
  FIELD AllowPublicUploads : bool [field_order=8 tag="config:\"allow_public_uploads\""]
  FIELD SweepInterval : time.Duration [field_order=9 tag="config:\"sweep_interval\""]
  FIELD SweepBatchSize : int [field_order=10 tag="config:\"sweep_batch_size\""]
  FIELD DeleteWorkerInterval : time.Duration [field_order=11 tag="config:\"delete_worker_interval\""]
  FIELD DeleteBatchSize : int [field_order=12 tag="config:\"delete_batch_size\""]
  FIELD DeleteConcurrency : int [field_order=13 tag="config:\"delete_concurrency\""]
  FIELD DeleteMaxAttempts : int [field_order=14 tag="config:\"delete_max_attempts\""]
  FIELD DeleteLeaseWindow : time.Duration [field_order=15 tag="config:\"delete_lease_window\""]
  METHOD EffectiveClaimTTL : func() time.Duration
  METHOD EffectiveDeleteBatchSize : func() int
  METHOD EffectiveDeleteConcurrency : func() int
  METHOD EffectiveDeleteLeaseWindow : func() time.Duration
  METHOD EffectiveDeleteMaxAttempts : func() int
  METHOD EffectiveDeleteWorkerInterval : func() time.Duration
  METHOD EffectiveMaxPendingClaims : func() int
  METHOD EffectiveMaxUploadSize : func() int64
  METHOD EffectiveSweepBatchSize : func() int
  METHOD EffectiveSweepInterval : func() time.Duration
CONST StorageFilesystem : github.com/coldsmirk/vef-framework-go/config.StorageProvider = "filesystem"
CONST StorageMemory : github.com/coldsmirk/vef-framework-go/config.StorageProvider = "memory"
CONST StorageMinIO : github.com/coldsmirk/vef-framework-go/config.StorageProvider = "minio"
TYPE StorageProvider : github.com/coldsmirk/vef-framework-go/config.StorageProvider
TYPE TokenType : github.com/coldsmirk/vef-framework-go/config.TokenType
CONST TokenTypeJWT : github.com/coldsmirk/vef-framework-go/config.TokenType = "jwt_token"
CONST TokenTypeOpaque : github.com/coldsmirk/vef-framework-go/config.TokenType = "opaque_token"

## github.com/coldsmirk/vef-framework-go/contextx
FUNC DB : func(ctx context.Context, fallbacks ...github.com/coldsmirk/vef-framework-go/orm.DB) github.com/coldsmirk/vef-framework-go/orm.DB
FUNC DataPermApplier : func(ctx context.Context) github.com/coldsmirk/vef-framework-go/security.DataPermissionApplier
CONST KeyDB : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 5
CONST KeyDataPermApplier : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 6
CONST KeyLogger : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 4
CONST KeyPrincipal : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 3
CONST KeyRequest : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 0
CONST KeyRequestID : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 1
CONST KeyRequestIP : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 2
CONST KeyRequestMethod : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 7
CONST KeyRequestPath : github.com/coldsmirk/vef-framework-go/contextx.contextKey = 8
FUNC Logger : func(ctx context.Context, fallbacks ...github.com/coldsmirk/vef-framework-go/logx.Logger) github.com/coldsmirk/vef-framework-go/logx.Logger
FUNC Principal : func(ctx context.Context) *github.com/coldsmirk/vef-framework-go/security.Principal
FUNC RequestID : func(ctx context.Context) string
FUNC RequestIP : func(ctx context.Context) string
FUNC RequestMethod : func(ctx context.Context) string
FUNC RequestPath : func(ctx context.Context) string
FUNC SetDB : func(ctx context.Context, db github.com/coldsmirk/vef-framework-go/orm.DB) context.Context
FUNC SetDataPermApplier : func(ctx context.Context, applier github.com/coldsmirk/vef-framework-go/security.DataPermissionApplier) context.Context
FUNC SetLogger : func(ctx context.Context, logger github.com/coldsmirk/vef-framework-go/logx.Logger) context.Context
FUNC SetPrincipal : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) context.Context
FUNC SetRequestID : func(ctx context.Context, requestID string) context.Context
FUNC SetRequestIP : func(ctx context.Context, ip string) context.Context
FUNC SetRequestMethod : func(ctx context.Context, method string) context.Context
FUNC SetRequestPath : func(ctx context.Context, path string) context.Context

## github.com/coldsmirk/vef-framework-go/copier
FUNC Copy : func(src any, dst any, options ...github.com/coldsmirk/vef-framework-go/copier.CopyOption) error
TYPE CopyOption : github.com/coldsmirk/vef-framework-go/copier.CopyOption
TYPE FieldNameMapping : github.com/coldsmirk/vef-framework-go/copier.FieldNameMapping
TYPE TypeConverter : github.com/coldsmirk/vef-framework-go/copier.TypeConverter
FUNC WithCaseInsensitive : func() github.com/coldsmirk/vef-framework-go/copier.CopyOption
FUNC WithDeepCopy : func() github.com/coldsmirk/vef-framework-go/copier.CopyOption
FUNC WithFieldNameMapping : func(mappings ...github.com/coldsmirk/vef-framework-go/copier.FieldNameMapping) github.com/coldsmirk/vef-framework-go/copier.CopyOption
FUNC WithIgnoreEmpty : func() github.com/coldsmirk/vef-framework-go/copier.CopyOption
FUNC WithTypeConverters : func(converters ...github.com/coldsmirk/vef-framework-go/copier.TypeConverter) github.com/coldsmirk/vef-framework-go/copier.CopyOption

## github.com/coldsmirk/vef-framework-go/cqrs
TYPE Action : github.com/coldsmirk/vef-framework-go/cqrs.Action
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/internal/cqrs.ActionKind
TYPE ActionKind : github.com/coldsmirk/vef-framework-go/cqrs.ActionKind
TYPE BaseCommand : github.com/coldsmirk/vef-framework-go/cqrs.BaseCommand
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/internal/cqrs.ActionKind
TYPE BaseQuery : github.com/coldsmirk/vef-framework-go/cqrs.BaseQuery
  METHOD Kind : func() github.com/coldsmirk/vef-framework-go/internal/cqrs.ActionKind
TYPE Behavior : github.com/coldsmirk/vef-framework-go/cqrs.Behavior
  METHOD Handle : func(ctx context.Context, action github.com/coldsmirk/vef-framework-go/internal/cqrs.Action, next func(ctx context.Context) (any, error)) (any, error)
TYPE BehaviorFunc : github.com/coldsmirk/vef-framework-go/cqrs.BehaviorFunc
  METHOD Handle : func(ctx context.Context, action github.com/coldsmirk/vef-framework-go/internal/cqrs.Action, next func(ctx context.Context) (any, error)) (any, error)
TYPE Bus : github.com/coldsmirk/vef-framework-go/cqrs.Bus
CONST Command : github.com/coldsmirk/vef-framework-go/internal/cqrs.ActionKind = 0
VAR ErrHandlerNotFound : error
VAR ErrResultTypeMismatch : error
TYPE Handler : github.com/coldsmirk/vef-framework-go/cqrs.Handler[TAction github.com/coldsmirk/vef-framework-go/internal/cqrs.Action, TResult any]
  METHOD Handle : func(ctx context.Context, action TAction) (TResult, error)
TYPE HandlerFunc : github.com/coldsmirk/vef-framework-go/cqrs.HandlerFunc[TAction github.com/coldsmirk/vef-framework-go/internal/cqrs.Action, TResult any]
  METHOD Handle : func(ctx context.Context, action TAction) (TResult, error)
FUNC NewBus : func(behaviors []github.com/coldsmirk/vef-framework-go/cqrs.Behavior) github.com/coldsmirk/vef-framework-go/cqrs.Bus
TYPE Ordered : github.com/coldsmirk/vef-framework-go/cqrs.Ordered
  METHOD Order : func() int
CONST Query : github.com/coldsmirk/vef-framework-go/internal/cqrs.ActionKind = 1
FUNC Register : func[TAction github.com/coldsmirk/vef-framework-go/internal/cqrs.Action, TResult any](bus github.com/coldsmirk/vef-framework-go/cqrs.Bus, handler github.com/coldsmirk/vef-framework-go/cqrs.Handler[TAction, TResult])
FUNC Send : func[TAction github.com/coldsmirk/vef-framework-go/internal/cqrs.Action, TResult any](ctx context.Context, bus github.com/coldsmirk/vef-framework-go/cqrs.Bus, action TAction) (TResult, error)
TYPE Unit : github.com/coldsmirk/vef-framework-go/cqrs.Unit

## github.com/coldsmirk/vef-framework-go/cron
TYPE CronJobDefinition : github.com/coldsmirk/vef-framework-go/cron.CronJobDefinition
TYPE DurationJobDefinition : github.com/coldsmirk/vef-framework-go/cron.DurationJobDefinition
TYPE DurationRandomJobDefinition : github.com/coldsmirk/vef-framework-go/cron.DurationRandomJobDefinition
VAR ErrJobNameRequired : error
VAR ErrJobTaskHandlerMustFunc : error
VAR ErrJobTaskHandlerRequired : error
TYPE Job : github.com/coldsmirk/vef-framework-go/cron.Job
  METHOD ID : func() string
  METHOD LastRun : func() (time.Time, error)
  METHOD Name : func() string
  METHOD NextRun : func() (time.Time, error)
  METHOD NextRuns : func(count int) ([]time.Time, error)
  METHOD RunNow : func() error
  METHOD Tags : func() []string
TYPE JobDefinition : github.com/coldsmirk/vef-framework-go/cron.JobDefinition
TYPE JobDescriptorOption : github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC NewCronJob : func(expression string, withSeconds bool, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.CronJobDefinition
FUNC NewDurationJob : func(interval time.Duration, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.DurationJobDefinition
FUNC NewDurationRandomJob : func(minInterval time.Duration, maxInterval time.Duration, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.DurationRandomJobDefinition
FUNC NewOneTimeJob : func(times []time.Time, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.OneTimeJobDefinition
FUNC NewScheduler : func(scheduler github.com/go-co-op/gocron/v2.Scheduler) github.com/coldsmirk/vef-framework-go/cron.Scheduler
TYPE OneTimeJobDefinition : github.com/coldsmirk/vef-framework-go/cron.OneTimeJobDefinition
TYPE Scheduler : github.com/coldsmirk/vef-framework-go/cron.Scheduler
  METHOD Jobs : func() []github.com/coldsmirk/vef-framework-go/cron.Job
  METHOD JobsWaitingInQueue : func() int
  METHOD NewJob : func(definition github.com/coldsmirk/vef-framework-go/cron.JobDefinition) (github.com/coldsmirk/vef-framework-go/cron.Job, error)
  METHOD RemoveByTags : func(tags ...string)
  METHOD RemoveJob : func(id string) error
  METHOD Start : func()
  METHOD StopJobs : func() error
  METHOD Update : func(id string, definition github.com/coldsmirk/vef-framework-go/cron.JobDefinition) (github.com/coldsmirk/vef-framework-go/cron.Job, error)
FUNC WithConcurrent : func() github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithContext : func(ctx context.Context) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithLimitedRuns : func(limitedRuns uint) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithName : func(name string) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithStartAt : func(startAt time.Time) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithStartImmediately : func() github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithStopAt : func(stopAt time.Time) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithTags : func(tags ...string) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption
FUNC WithTask : func(handler any, params ...any) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption

## github.com/coldsmirk/vef-framework-go/crud
FUNC ApplyDataPermission : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, ctx github.com/gofiber/fiber/v3.Ctx) error
TYPE Builder : github.com/coldsmirk/vef-framework-go/crud.Builder[T any]
  METHOD Action : func(action string) T
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD EnableAudit : func() T
  METHOD Public : func() T
  METHOD RateLimit : func(maxRequests int, period time.Duration) T
  METHOD RequiredPermission : func(token string) T
  METHOD Timeout : func(timeout time.Duration) T
TYPE Create : github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
  METHOD WithPostCreate : func(processor github.com/coldsmirk/vef-framework-go/crud.PostCreateProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
  METHOD WithPreCreate : func(processor github.com/coldsmirk/vef-framework-go/crud.PreCreateProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
TYPE CreateMany : github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
  METHOD WithPostCreateMany : func(processor github.com/coldsmirk/vef-framework-go/crud.PostCreateManyProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
  METHOD WithPreCreateMany : func(processor github.com/coldsmirk/vef-framework-go/crud.PreCreateManyProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
TYPE CreateManyParams : github.com/coldsmirk/vef-framework-go/crud.CreateManyParams[TParams any]
  FIELD P : github.com/coldsmirk/vef-framework-go/api.P [field_order=1 tag=""]
  FIELD List : []TParams [field_order=2 tag="json:\"list\" validate:\"required,min=1,dive\" label_i18n:\"crud_batch_create_list\""]
TYPE DataOption : github.com/coldsmirk/vef-framework-go/crud.DataOption
  FIELD Label : string [field_order=1 tag="json:\"label\" bun:\"label\""]
  FIELD Value : string [field_order=2 tag="json:\"value\" bun:\"value\""]
  FIELD Description : string [field_order=3 tag="json:\"description,omitempty\" bun:\"description\""]
  FIELD Meta : map[string]any [field_order=4 tag="json:\"meta,omitempty\" bun:\"meta\""]
TYPE DataOptionColumnMapping : github.com/coldsmirk/vef-framework-go/crud.DataOptionColumnMapping
  FIELD LabelColumn : string [field_order=1 tag="json:\"labelColumn\""]
  FIELD ValueColumn : string [field_order=2 tag="json:\"valueColumn\""]
  FIELD DescriptionColumn : string [field_order=3 tag="json:\"descriptionColumn\""]
  FIELD MetaColumns : []string [field_order=4 tag="json:\"metaColumns\""]
TYPE DataOptionConfig : github.com/coldsmirk/vef-framework-go/crud.DataOptionConfig
  FIELD M : github.com/coldsmirk/vef-framework-go/api.M [field_order=1 tag=""]
  FIELD DataOptionColumnMapping : github.com/coldsmirk/vef-framework-go/crud.DataOptionColumnMapping [field_order=2 tag=""]
  FIELD DescriptionColumn : string [promoted_from=DataOptionColumnMapping depth=1 field_order=3 tag="json:\"descriptionColumn\""]
  FIELD LabelColumn : string [promoted_from=DataOptionColumnMapping depth=1 field_order=1 tag="json:\"labelColumn\""]
  FIELD MetaColumns : []string [promoted_from=DataOptionColumnMapping depth=1 field_order=4 tag="json:\"metaColumns\""]
  FIELD ValueColumn : string [promoted_from=DataOptionColumnMapping depth=1 field_order=2 tag="json:\"valueColumn\""]
TYPE Delete : github.com/coldsmirk/vef-framework-go/crud.Delete[TModel any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD WithPostDelete : func(processor github.com/coldsmirk/vef-framework-go/crud.PostDeleteProcessor[TModel]) github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
  METHOD WithPreDelete : func(processor github.com/coldsmirk/vef-framework-go/crud.PreDeleteProcessor[TModel]) github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
TYPE DeleteMany : github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD WithPostDeleteMany : func(processor github.com/coldsmirk/vef-framework-go/crud.PostDeleteManyProcessor[TModel]) github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
  METHOD WithPreDeleteMany : func(processor github.com/coldsmirk/vef-framework-go/crud.PreDeleteManyProcessor[TModel]) github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
TYPE DeleteManyParams : github.com/coldsmirk/vef-framework-go/crud.DeleteManyParams
  FIELD P : github.com/coldsmirk/vef-framework-go/api.P [field_order=1 tag=""]
  FIELD PKs : []any [field_order=2 tag="json:\"pks\" validate:\"required,min=1\" label_i18n:\"crud_batch_delete_pks\""]
CONST DescriptionColumn : untyped string = "description"
VAR ErrAuditUserCompositePK : error
CONST ErrCodeCompositePrimaryKeyRequiresMap : untyped int = 2403
CONST ErrCodeFieldNotExistInModel : untyped int = 2401
CONST ErrCodeFileOpenFailed : untyped int = 2408
CONST ErrCodeImportRequiresFile : untyped int = 2406
CONST ErrCodeImportRequiresMultipart : untyped int = 2405
CONST ErrCodeImportTypeAssertionFailed : untyped int = 2409
CONST ErrCodeImportValidationFailed : untyped int = 2410
CONST ErrCodePrimaryKeyRequired : untyped int = 2402
CONST ErrCodeProcessorInvalidReturn : untyped int = 2400
CONST ErrCodeUnsupportedExportFormat : untyped int = 2404
CONST ErrCodeUnsupportedImportFormat : untyped int = 2407
VAR ErrCompositePrimaryKeyRequiresMap : github.com/coldsmirk/vef-framework-go/result.Error
FUNC ErrFieldNotExistInModel : func(field string, name string, model string) github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrFileOpenFailed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrImportRequiresFile : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrImportRequiresMultipart : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrImportTypeAssertionFailed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrModelNoPrimaryKey : error
FUNC ErrPrimaryKeyRequired : func(field string) github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUnsupportedExportFormat : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUnsupportedImportFormat : github.com/coldsmirk/vef-framework-go/result.Error
TYPE Export : github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD Process : func(input []TModel, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithCsvOptions : func(opts ...github.com/coldsmirk/vef-framework-go/csv.ExportOption) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithDefaultFormat : func(format github.com/coldsmirk/vef-framework-go/crud.TabularFormat) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithExcelOptions : func(opts ...github.com/coldsmirk/vef-framework-go/excel.ExportOption) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithFilenameBuilder : func(builder github.com/coldsmirk/vef-framework-go/crud.FilenameBuilder[TSearch]) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithPreExport : func(processor github.com/coldsmirk/vef-framework-go/crud.PreExportProcessor[TModel, TSearch]) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[[]TModel, TSearch]) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
TYPE FilenameBuilder : github.com/coldsmirk/vef-framework-go/crud.FilenameBuilder[TSearch any]
TYPE Find : github.com/coldsmirk/vef-framework-go/crud.Find[TModel, TSearch, TProcessorIn, TOperation any]
  METHOD Action : func(action string) TOperation
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() TOperation
  METHOD EnableAudit : func() TOperation
  METHOD Process : func(input TProcessorIn, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Public : func() TOperation
  METHOD RateLimit : func(maxRequests int, period time.Duration) TOperation
  METHOD RequiredPermission : func(token string) TOperation
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) TOperation
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) TOperation
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) TOperation
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) TOperation
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) TOperation
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[TProcessorIn, TSearch]) TOperation
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) TOperation
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) TOperation
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) TOperation
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) TOperation
TYPE FindAll : github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD Process : func(input []TModel, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[[]TModel, TSearch]) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
TYPE FindOne : github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD Process : func(input TModel, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[TModel, TSearch]) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
TYPE FindOperationConfig : github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig
  FIELD QueryParts : *github.com/coldsmirk/vef-framework-go/crud.QueryPartsConfig [field_order=1 tag=""]
TYPE FindOperationOption : github.com/coldsmirk/vef-framework-go/crud.FindOperationOption
  FIELD Parts : []github.com/coldsmirk/vef-framework-go/crud.QueryPart [field_order=1 tag=""]
  FIELD Applier : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search any, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx) error [field_order=2 tag=""]
TYPE FindOptions : github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD Process : func(input []github.com/coldsmirk/vef-framework-go/crud.DataOption, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithDefaultColumnMapping : func(mapping *github.com/coldsmirk/vef-framework-go/crud.DataOptionColumnMapping) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[[]github.com/coldsmirk/vef-framework-go/crud.DataOption, TSearch]) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
TYPE FindPage : github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD Process : func(input []TModel, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithDefaultPageSize : func(size int) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[[]TModel, TSearch]) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
TYPE FindTree : github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD Process : func(input []TModel, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithIDColumn : func(name string) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithParentIDColumn : func(name string) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[[]TModel, TSearch]) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
TYPE FindTreeOptions : github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD ConfigureQuery : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, meta github.com/coldsmirk/vef-framework-go/api.Meta, ctx github.com/gofiber/fiber/v3.Ctx, part github.com/coldsmirk/vef-framework-go/crud.QueryPart) error
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD Process : func(input []github.com/coldsmirk/vef-framework-go/crud.TreeDataOption, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) any
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD Setup : func(db github.com/coldsmirk/vef-framework-go/orm.DB, config *github.com/coldsmirk/vef-framework-go/crud.FindOperationConfig, opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) error
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithAuditUserNames : func(userModel any, nameColumn ...string) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithCondition : func(fn func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder), parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithDefaultColumnMapping : func(mapping *github.com/coldsmirk/vef-framework-go/crud.DataOptionColumnMapping) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithDefaultSort : func(sort ...*github.com/coldsmirk/vef-framework-go/sortx.OrderSpec) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithIDColumn : func(name string) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithOptions : func(opts ...*github.com/coldsmirk/vef-framework-go/crud.FindOperationOption) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithParentIDColumn : func(name string) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithProcessor : func(processor github.com/coldsmirk/vef-framework-go/crud.Processor[[]github.com/coldsmirk/vef-framework-go/crud.TreeDataOption, TSearch]) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithQueryApplier : func(applier func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery, search TSearch, ctx github.com/gofiber/fiber/v3.Ctx) error, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithRelation : func(relation *github.com/coldsmirk/vef-framework-go/orm.RelationSpec, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithSelect : func(column string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
  METHOD WithSelectAs : func(column string, alias string, parts ...github.com/coldsmirk/vef-framework-go/crud.QueryPart) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
CONST FormatCsv : github.com/coldsmirk/vef-framework-go/crud.TabularFormat = "csv"
CONST FormatExcel : github.com/coldsmirk/vef-framework-go/crud.TabularFormat = "excel"
FUNC GetAuditUserNameRelations : func(userModel any, nameColumn ...string) []*github.com/coldsmirk/vef-framework-go/orm.RelationSpec
CONST IDColumn : untyped string = "id"
TYPE Import : github.com/coldsmirk/vef-framework-go/crud.Import[TModel any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD WithCsvOptions : func(opts ...github.com/coldsmirk/vef-framework-go/csv.ImportOption) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD WithDefaultFormat : func(format github.com/coldsmirk/vef-framework-go/crud.TabularFormat) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD WithExcelOptions : func(opts ...github.com/coldsmirk/vef-framework-go/excel.ImportOption) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD WithPostImport : func(processor github.com/coldsmirk/vef-framework-go/crud.PostImportProcessor[TModel]) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
  METHOD WithPreImport : func(processor github.com/coldsmirk/vef-framework-go/crud.PreImportProcessor[TModel]) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
CONST LabelColumn : untyped string = "label"
CONST MessageCreated : untyped string = "crud_created"
CONST MessageDeleted : untyped string = "crud_deleted"
CONST MessageImported : untyped string = "crud_imported"
CONST MessageUpdated : untyped string = "crud_updated"
FUNC NewBuilder : func[T any](self T, kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.Builder[T]
FUNC NewCreate : func[TModel, TParams any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.Create[TModel, TParams]
FUNC NewCreateMany : func[TModel, TParams any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.CreateMany[TModel, TParams]
FUNC NewDelete : func[TModel any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.Delete[TModel]
FUNC NewDeleteMany : func[TModel any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.DeleteMany[TModel]
FUNC NewExport : func[TModel, TSearch any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.Export[TModel, TSearch]
FUNC NewFind : func[TModel, TSearch, TProcessor, TOperation any](self TOperation, kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.Find[TModel, TSearch, TProcessor, TOperation]
FUNC NewFindAll : func[TModel, TSearch any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.FindAll[TModel, TSearch]
FUNC NewFindOne : func[TModel, TSearch any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.FindOne[TModel, TSearch]
FUNC NewFindOptions : func[TModel, TSearch any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.FindOptions[TModel, TSearch]
FUNC NewFindPage : func[TModel, TSearch any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.FindPage[TModel, TSearch]
FUNC NewFindTree : func[TModel, TSearch any](treeBuilder func(flatModels []TModel) []TModel, kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.FindTree[TModel, TSearch]
FUNC NewFindTreeOptions : func[TModel, TSearch any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.FindTreeOptions[TModel, TSearch]
FUNC NewImport : func[TModel any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.Import[TModel]
FUNC NewUpdate : func[TModel, TParams any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
FUNC NewUpdateMany : func[TModel, TParams any](kind ...github.com/coldsmirk/vef-framework-go/api.Kind) github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
CONST ParentIDColumn : untyped string = "parent_id"
TYPE PostCreateManyProcessor : github.com/coldsmirk/vef-framework-go/crud.PostCreateManyProcessor[TModel, TParams any]
TYPE PostCreateProcessor : github.com/coldsmirk/vef-framework-go/crud.PostCreateProcessor[TModel, TParams any]
TYPE PostDeleteManyProcessor : github.com/coldsmirk/vef-framework-go/crud.PostDeleteManyProcessor[TModel any]
TYPE PostDeleteProcessor : github.com/coldsmirk/vef-framework-go/crud.PostDeleteProcessor[TModel any]
TYPE PostImportProcessor : github.com/coldsmirk/vef-framework-go/crud.PostImportProcessor[TModel any]
TYPE PostUpdateManyProcessor : github.com/coldsmirk/vef-framework-go/crud.PostUpdateManyProcessor[TModel, TParams any]
TYPE PostUpdateProcessor : github.com/coldsmirk/vef-framework-go/crud.PostUpdateProcessor[TModel, TParams any]
TYPE PreCreateManyProcessor : github.com/coldsmirk/vef-framework-go/crud.PreCreateManyProcessor[TModel, TParams any]
TYPE PreCreateProcessor : github.com/coldsmirk/vef-framework-go/crud.PreCreateProcessor[TModel, TParams any]
TYPE PreDeleteManyProcessor : github.com/coldsmirk/vef-framework-go/crud.PreDeleteManyProcessor[TModel any]
TYPE PreDeleteProcessor : github.com/coldsmirk/vef-framework-go/crud.PreDeleteProcessor[TModel any]
TYPE PreExportProcessor : github.com/coldsmirk/vef-framework-go/crud.PreExportProcessor[TModel, TSearch any]
TYPE PreImportProcessor : github.com/coldsmirk/vef-framework-go/crud.PreImportProcessor[TModel any]
TYPE PreUpdateManyProcessor : github.com/coldsmirk/vef-framework-go/crud.PreUpdateManyProcessor[TModel, TParams any]
TYPE PreUpdateProcessor : github.com/coldsmirk/vef-framework-go/crud.PreUpdateProcessor[TModel, TParams any]
TYPE Processor : github.com/coldsmirk/vef-framework-go/crud.Processor[TIn, TSearch any]
CONST QueryAll : github.com/coldsmirk/vef-framework-go/crud.QueryPart = 3
CONST QueryBase : github.com/coldsmirk/vef-framework-go/crud.QueryPart = 1
TYPE QueryPart : github.com/coldsmirk/vef-framework-go/crud.QueryPart
TYPE QueryPartsConfig : github.com/coldsmirk/vef-framework-go/crud.QueryPartsConfig
  FIELD Condition : []github.com/coldsmirk/vef-framework-go/crud.QueryPart [field_order=1 tag=""]
  FIELD Sort : []github.com/coldsmirk/vef-framework-go/crud.QueryPart [field_order=2 tag=""]
  FIELD AuditUserRelation : []github.com/coldsmirk/vef-framework-go/crud.QueryPart [field_order=3 tag=""]
CONST QueryRecursive : github.com/coldsmirk/vef-framework-go/crud.QueryPart = 2
CONST QueryRoot : github.com/coldsmirk/vef-framework-go/crud.QueryPart = 0
CONST RESTActionCreate : untyped string = "post /"
CONST RESTActionCreateMany : untyped string = "post /many"
CONST RESTActionDelete : untyped string = "delete /:id"
CONST RESTActionDeleteMany : untyped string = "delete /many"
CONST RESTActionExport : untyped string = "get /export"
CONST RESTActionFindAll : untyped string = "get /"
CONST RESTActionFindOne : untyped string = "get /:id"
CONST RESTActionFindOptions : untyped string = "get /options"
CONST RESTActionFindPage : untyped string = "get /page"
CONST RESTActionFindTree : untyped string = "get /tree"
CONST RESTActionFindTreeOptions : untyped string = "get /tree/options"
CONST RESTActionImport : untyped string = "post /import"
CONST RESTActionUpdate : untyped string = "put /:id"
CONST RESTActionUpdateMany : untyped string = "put /many"
CONST RPCActionCreate : untyped string = "create"
CONST RPCActionCreateMany : untyped string = "create_many"
CONST RPCActionDelete : untyped string = "delete"
CONST RPCActionDeleteMany : untyped string = "delete_many"
CONST RPCActionExport : untyped string = "export"
CONST RPCActionFindAll : untyped string = "find_all"
CONST RPCActionFindOne : untyped string = "find_one"
CONST RPCActionFindOptions : untyped string = "find_options"
CONST RPCActionFindPage : untyped string = "find_page"
CONST RPCActionFindTree : untyped string = "find_tree"
CONST RPCActionFindTreeOptions : untyped string = "find_tree_options"
CONST RPCActionImport : untyped string = "import"
CONST RPCActionUpdate : untyped string = "update"
CONST RPCActionUpdateMany : untyped string = "update_many"
TYPE Sortable : github.com/coldsmirk/vef-framework-go/crud.Sortable
  FIELD M : github.com/coldsmirk/vef-framework-go/api.M [field_order=1 tag=""]
  FIELD Sort : []github.com/coldsmirk/vef-framework-go/sortx.OrderSpec [field_order=2 tag="json:\"sort\""]
TYPE TabularFormat : github.com/coldsmirk/vef-framework-go/crud.TabularFormat
TYPE TreeDataOption : github.com/coldsmirk/vef-framework-go/crud.TreeDataOption
  FIELD DataOption : github.com/coldsmirk/vef-framework-go/crud.DataOption [field_order=1 tag=""]
  FIELD ID : string [field_order=2 tag="json:\"-\" bun:\"id\""]
  FIELD ParentID : *string [field_order=3 tag="json:\"-\" bun:\"parent_id\""]
  FIELD Children : []github.com/coldsmirk/vef-framework-go/crud.TreeDataOption [field_order=4 tag="json:\"children,omitempty\" bun:\"-\""]
  FIELD Description : string [promoted_from=DataOption depth=1 field_order=3 tag="json:\"description,omitempty\" bun:\"description\""]
  FIELD Label : string [promoted_from=DataOption depth=1 field_order=1 tag="json:\"label\" bun:\"label\""]
  FIELD Meta : map[string]any [promoted_from=DataOption depth=1 field_order=4 tag="json:\"meta,omitempty\" bun:\"meta\""]
  FIELD Value : string [promoted_from=DataOption depth=1 field_order=2 tag="json:\"value\" bun:\"value\""]
TYPE Update : github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD WithPostUpdate : func(processor github.com/coldsmirk/vef-framework-go/crud.PostUpdateProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
  METHOD WithPreUpdate : func(processor github.com/coldsmirk/vef-framework-go/crud.PreUpdateProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.Update[TModel, TParams]
TYPE UpdateMany : github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams any]
  METHOD Action : func(action string) github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD Build : func(handler any) github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD DisableDataPerm : func() github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD EnableAudit : func() github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD Provide : func() []github.com/coldsmirk/vef-framework-go/api.OperationSpec
  METHOD Public : func() github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD RateLimit : func(maxRequests int, period time.Duration) github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD RequiredPermission : func(token string) github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD Timeout : func(timeout time.Duration) github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD WithPostUpdateMany : func(processor github.com/coldsmirk/vef-framework-go/crud.PostUpdateManyProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
  METHOD WithPreUpdateMany : func(processor github.com/coldsmirk/vef-framework-go/crud.PreUpdateManyProcessor[TModel, TParams]) github.com/coldsmirk/vef-framework-go/crud.UpdateMany[TModel, TParams]
TYPE UpdateManyParams : github.com/coldsmirk/vef-framework-go/crud.UpdateManyParams[TParams any]
  FIELD P : github.com/coldsmirk/vef-framework-go/api.P [field_order=1 tag=""]
  FIELD List : []TParams [field_order=2 tag="json:\"list\" validate:\"required,min=1,dive\" label_i18n:\"crud_batch_update_list\""]
CONST ValueColumn : untyped string = "value"

## github.com/coldsmirk/vef-framework-go/cryptox
TYPE AESMode : github.com/coldsmirk/vef-framework-go/cryptox.AESMode
TYPE AESOption : github.com/coldsmirk/vef-framework-go/cryptox.AESOption
CONST AesModeCbc : github.com/coldsmirk/vef-framework-go/cryptox.AESMode = "CBC"
CONST AesModeGcm : github.com/coldsmirk/vef-framework-go/cryptox.AESMode = "GCM"
TYPE Cipher : github.com/coldsmirk/vef-framework-go/cryptox.Cipher
  METHOD Decrypt : func(ciphertext string) (string, error)
  METHOD Encrypt : func(plaintext string) (string, error)
TYPE CipherSigner : github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner
  METHOD Decrypt : func(ciphertext string) (string, error)
  METHOD Encrypt : func(plaintext string) (string, error)
  METHOD Sign : func(data string) (signature string, err error)
  METHOD Verify : func(data string, signature string) (bool, error)
TYPE ECDSACurve : github.com/coldsmirk/vef-framework-go/cryptox.ECDSACurve
TYPE ECIESCurve : github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve
CONST EcdsaCurveP224 : github.com/coldsmirk/vef-framework-go/cryptox.ECDSACurve = "P224"
CONST EcdsaCurveP256 : github.com/coldsmirk/vef-framework-go/cryptox.ECDSACurve = "P256"
CONST EcdsaCurveP384 : github.com/coldsmirk/vef-framework-go/cryptox.ECDSACurve = "P384"
CONST EcdsaCurveP521 : github.com/coldsmirk/vef-framework-go/cryptox.ECDSACurve = "P521"
CONST EciesCurveP256 : github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve = "P256"
CONST EciesCurveP384 : github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve = "P384"
CONST EciesCurveP521 : github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve = "P521"
CONST EciesCurveX25519 : github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve = "X25519"
VAR ErrAtLeastOneKeyRequired : error
VAR ErrCiphertextNotMultipleOfBlock : error
VAR ErrCiphertextTooShort : error
VAR ErrDataEmpty : error
VAR ErrFailedDecodePEMBlock : error
VAR ErrInvalidAESKeySize : error
VAR ErrInvalidIVSizeCBC : error
VAR ErrInvalidPadding : error
VAR ErrInvalidSM4KeySize : error
VAR ErrInvalidSignature : error
VAR ErrNotECDSAPrivateKey : error
VAR ErrNotECDSAPublicKey : error
VAR ErrNotRSAPrivateKey : error
VAR ErrNotRSAPublicKey : error
VAR ErrPrivateKeyRequiredForDecrypt : error
VAR ErrPrivateKeyRequiredForSign : error
VAR ErrPublicKeyRequiredForEncrypt : error
VAR ErrPublicKeyRequiredForVerify : error
VAR ErrUnsupportedPEMType : error
TYPE FixedIVDecrypter : github.com/coldsmirk/vef-framework-go/cryptox.FixedIVDecrypter
  METHOD DecryptWithFixedIV : func(ciphertext string) (string, error)
FUNC GenerateECDSAKey : func(curve github.com/coldsmirk/vef-framework-go/cryptox.ECDSACurve) (*crypto/ecdsa.PrivateKey, error)
FUNC GenerateECIESKey : func(curve github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve) (*crypto/ecdh.PrivateKey, error)
FUNC NewAES : func(key []byte, opts ...github.com/coldsmirk/vef-framework-go/cryptox.AESOption) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewAESFromBase64 : func(keyBase64 string, opts ...github.com/coldsmirk/vef-framework-go/cryptox.AESOption) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewAESFromHex : func(keyHex string, opts ...github.com/coldsmirk/vef-framework-go/cryptox.AESOption) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewECDSA : func(privateKey *crypto/ecdsa.PrivateKey, publicKey *crypto/ecdsa.PublicKey) (github.com/coldsmirk/vef-framework-go/cryptox.Signer, error)
FUNC NewECDSAFromBase64 : func(privateKeyBase64 string, publicKeyBase64 string) (github.com/coldsmirk/vef-framework-go/cryptox.Signer, error)
FUNC NewECDSAFromHex : func(privateKeyHex string, publicKeyHex string) (github.com/coldsmirk/vef-framework-go/cryptox.Signer, error)
FUNC NewECDSAFromPEM : func(privatePEM []byte, publicPEM []byte) (github.com/coldsmirk/vef-framework-go/cryptox.Signer, error)
FUNC NewECIES : func(privateKey *crypto/ecdh.PrivateKey, publicKey *crypto/ecdh.PublicKey) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewECIESFromBase64 : func(privateKeyBase64 string, publicKeyBase64 string, curve github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewECIESFromBytes : func(privateKeyBytes []byte, publicKeyBytes []byte, curve github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewECIESFromHex : func(privateKeyHex string, publicKeyHex string, curve github.com/coldsmirk/vef-framework-go/cryptox.ECIESCurve) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewRSA : func(privateKey *crypto/rsa.PrivateKey, publicKey *crypto/rsa.PublicKey, opts ...github.com/coldsmirk/vef-framework-go/cryptox.RSAOption) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewRSAFromBase64 : func(privateKeyBase64 string, publicKeyBase64 string, opts ...github.com/coldsmirk/vef-framework-go/cryptox.RSAOption) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewRSAFromHex : func(privateKeyHex string, publicKeyHex string, opts ...github.com/coldsmirk/vef-framework-go/cryptox.RSAOption) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewRSAFromPEM : func(privatePEM []byte, publicPEM []byte, opts ...github.com/coldsmirk/vef-framework-go/cryptox.RSAOption) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewSM2 : func(privateKey *github.com/tjfoc/gmsm/sm2.PrivateKey, publicKey *github.com/tjfoc/gmsm/sm2.PublicKey) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewSM2FromBase64 : func(privateKeyBase64 string, publicKeyBase64 string) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewSM2FromHex : func(privateKeyHex string, publicKeyHex string) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewSM2FromPEM : func(privatePEM []byte, publicPEM []byte) (github.com/coldsmirk/vef-framework-go/cryptox.CipherSigner, error)
FUNC NewSM4 : func(key []byte, opts ...github.com/coldsmirk/vef-framework-go/cryptox.SM4Option) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewSM4FromBase64 : func(keyBase64 string, opts ...github.com/coldsmirk/vef-framework-go/cryptox.SM4Option) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
FUNC NewSM4FromHex : func(keyHex string, opts ...github.com/coldsmirk/vef-framework-go/cryptox.SM4Option) (github.com/coldsmirk/vef-framework-go/cryptox.Cipher, error)
TYPE RSAMode : github.com/coldsmirk/vef-framework-go/cryptox.RSAMode
TYPE RSAOption : github.com/coldsmirk/vef-framework-go/cryptox.RSAOption
TYPE RSASignMode : github.com/coldsmirk/vef-framework-go/cryptox.RSASignMode
CONST RsaModeOAEP : github.com/coldsmirk/vef-framework-go/cryptox.RSAMode = "OAEP"
CONST RsaModePKCS1v15 : github.com/coldsmirk/vef-framework-go/cryptox.RSAMode = "PKCS1v15"
CONST RsaSignModePKCS1v15 : github.com/coldsmirk/vef-framework-go/cryptox.RSASignMode = "PKCS1v15"
CONST RsaSignModePSS : github.com/coldsmirk/vef-framework-go/cryptox.RSASignMode = "PSS"
TYPE SM4Option : github.com/coldsmirk/vef-framework-go/cryptox.SM4Option
TYPE Signer : github.com/coldsmirk/vef-framework-go/cryptox.Signer
  METHOD Sign : func(data string) (signature string, err error)
  METHOD Verify : func(data string, signature string) (bool, error)
FUNC WithAESIv : func(iv []byte) github.com/coldsmirk/vef-framework-go/cryptox.AESOption
FUNC WithAESMode : func(mode github.com/coldsmirk/vef-framework-go/cryptox.AESMode) github.com/coldsmirk/vef-framework-go/cryptox.AESOption
FUNC WithRSAMode : func(mode github.com/coldsmirk/vef-framework-go/cryptox.RSAMode) github.com/coldsmirk/vef-framework-go/cryptox.RSAOption
FUNC WithRSASignMode : func(signMode github.com/coldsmirk/vef-framework-go/cryptox.RSASignMode) github.com/coldsmirk/vef-framework-go/cryptox.RSAOption
FUNC WithSM4Iv : func(iv []byte) github.com/coldsmirk/vef-framework-go/cryptox.SM4Option

## github.com/coldsmirk/vef-framework-go/csv
TYPE ExportOption : github.com/coldsmirk/vef-framework-go/csv.ExportOption
TYPE ImportOption : github.com/coldsmirk/vef-framework-go/csv.ImportOption
FUNC NewExporter : func(adapter github.com/coldsmirk/vef-framework-go/tabular.RowAdapter, opts ...github.com/coldsmirk/vef-framework-go/csv.ExportOption) github.com/coldsmirk/vef-framework-go/tabular.Exporter
FUNC NewExporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/csv.ExportOption) github.com/coldsmirk/vef-framework-go/tabular.Exporter
FUNC NewImporter : func(adapter github.com/coldsmirk/vef-framework-go/tabular.RowAdapter, opts ...github.com/coldsmirk/vef-framework-go/csv.ImportOption) github.com/coldsmirk/vef-framework-go/tabular.Importer
FUNC NewImporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/csv.ImportOption) github.com/coldsmirk/vef-framework-go/tabular.Importer
FUNC NewMapExporter : func(specs []github.com/coldsmirk/vef-framework-go/tabular.ColumnSpec, opts ...github.com/coldsmirk/vef-framework-go/csv.ExportOption) (github.com/coldsmirk/vef-framework-go/tabular.Exporter, error)
FUNC NewMapImporter : func(specs []github.com/coldsmirk/vef-framework-go/tabular.ColumnSpec, mapOpts []github.com/coldsmirk/vef-framework-go/tabular.MapOption, opts ...github.com/coldsmirk/vef-framework-go/csv.ImportOption) (github.com/coldsmirk/vef-framework-go/tabular.Importer, error)
FUNC NewTypedExporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/csv.ExportOption) github.com/coldsmirk/vef-framework-go/tabular.TypedExporter[T]
FUNC NewTypedImporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/csv.ImportOption) github.com/coldsmirk/vef-framework-go/tabular.TypedImporter[T]
FUNC WithCRLF : func() github.com/coldsmirk/vef-framework-go/csv.ExportOption
FUNC WithComment : func(comment rune) github.com/coldsmirk/vef-framework-go/csv.ImportOption
FUNC WithExportDelimiter : func(delimiter rune) github.com/coldsmirk/vef-framework-go/csv.ExportOption
FUNC WithImportDelimiter : func(delimiter rune) github.com/coldsmirk/vef-framework-go/csv.ImportOption
FUNC WithSkipRows : func(rows int) github.com/coldsmirk/vef-framework-go/csv.ImportOption
FUNC WithoutHeader : func() github.com/coldsmirk/vef-framework-go/csv.ImportOption
FUNC WithoutTrimSpace : func() github.com/coldsmirk/vef-framework-go/csv.ImportOption
FUNC WithoutWriteHeader : func() github.com/coldsmirk/vef-framework-go/csv.ExportOption

## github.com/coldsmirk/vef-framework-go/datasource
TYPE ConnectionInfo : github.com/coldsmirk/vef-framework-go/datasource.ConnectionInfo
  FIELD Version : string [field_order=1 tag=""]
VAR ErrClosed : error
VAR ErrExists : error
VAR ErrNameInvalid : error
VAR ErrNotFound : error
VAR ErrPrimaryReserved : error
CONST PrimaryName : untyped string = "primary"
TYPE Provider : github.com/coldsmirk/vef-framework-go/datasource.Provider
  METHOD Load : func(ctx context.Context) ([]github.com/coldsmirk/vef-framework-go/datasource.Spec, error)
  METHOD Name : func() string
TYPE ReconcileOption : github.com/coldsmirk/vef-framework-go/datasource.ReconcileOption
TYPE ReconcileOptions : github.com/coldsmirk/vef-framework-go/datasource.ReconcileOptions
  FIELD DryRun : bool [field_order=1 tag=""]
TYPE ReconcileReport : github.com/coldsmirk/vef-framework-go/datasource.ReconcileReport
  FIELD Added : []string [field_order=1 tag=""]
  FIELD Updated : []string [field_order=2 tag=""]
  FIELD Removed : []string [field_order=3 tag=""]
  FIELD Errors : map[string]error [field_order=4 tag=""]
TYPE RegisterOption : github.com/coldsmirk/vef-framework-go/datasource.RegisterOption
TYPE RegisterOptions : github.com/coldsmirk/vef-framework-go/datasource.RegisterOptions
  FIELD CloseGrace : time.Duration [field_order=1 tag=""]
TYPE Registry : github.com/coldsmirk/vef-framework-go/datasource.Registry
  METHOD Get : func(name string) (github.com/coldsmirk/vef-framework-go/orm.DB, error)
  METHOD Has : func(name string) bool
  METHOD HealthCheck : func(ctx context.Context) map[string]error
  METHOD Kind : func(name string) (github.com/coldsmirk/vef-framework-go/config.DBKind, error)
  METHOD Names : func() []string
  METHOD Primary : func() github.com/coldsmirk/vef-framework-go/orm.DB
  METHOD Reconcile : func(ctx context.Context, specs []github.com/coldsmirk/vef-framework-go/datasource.Spec, opts ...github.com/coldsmirk/vef-framework-go/datasource.ReconcileOption) (github.com/coldsmirk/vef-framework-go/datasource.ReconcileReport, error)
  METHOD Register : func(ctx context.Context, name string, cfg github.com/coldsmirk/vef-framework-go/config.DataSourceConfig) (github.com/coldsmirk/vef-framework-go/orm.DB, error)
  METHOD TestConnection : func(ctx context.Context, cfg github.com/coldsmirk/vef-framework-go/config.DataSourceConfig) (github.com/coldsmirk/vef-framework-go/datasource.ConnectionInfo, error)
  METHOD Unregister : func(ctx context.Context, name string, opts ...github.com/coldsmirk/vef-framework-go/datasource.RegisterOption) error
  METHOD Update : func(ctx context.Context, name string, cfg github.com/coldsmirk/vef-framework-go/config.DataSourceConfig, opts ...github.com/coldsmirk/vef-framework-go/datasource.RegisterOption) (github.com/coldsmirk/vef-framework-go/orm.DB, error)
TYPE Spec : github.com/coldsmirk/vef-framework-go/datasource.Spec
  FIELD Name : string [field_order=1 tag=""]
  FIELD Config : github.com/coldsmirk/vef-framework-go/config.DataSourceConfig [field_order=2 tag=""]
FUNC WithCloseGrace : func(d time.Duration) github.com/coldsmirk/vef-framework-go/datasource.RegisterOption
FUNC WithReconcileDryRun : func() github.com/coldsmirk/vef-framework-go/datasource.ReconcileOption

## github.com/coldsmirk/vef-framework-go/dbx
FUNC ColumnWithAlias : func(column string, alias ...string) string
FUNC IsDuplicateKeyError : func(err error) bool
FUNC IsForeignKeyError : func(err error) bool

## github.com/coldsmirk/vef-framework-go/decimal
VAR Avg : func(first github.com/shopspring/decimal.Decimal, rest ...github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
TYPE Decimal : github.com/coldsmirk/vef-framework-go/decimal.Decimal
  METHOD Abs : func() github.com/shopspring/decimal.Decimal
  METHOD Add : func(d2 github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
  METHOD Atan : func() github.com/shopspring/decimal.Decimal
  METHOD BigFloat : func() *math/big.Float
  METHOD BigInt : func() *math/big.Int
  METHOD Ceil : func() github.com/shopspring/decimal.Decimal
  METHOD Cmp : func(d2 github.com/shopspring/decimal.Decimal) int
  METHOD Coefficient : func() *math/big.Int
  METHOD CoefficientInt64 : func() int64
  METHOD Compare : func(d2 github.com/shopspring/decimal.Decimal) int
  METHOD Copy : func() github.com/shopspring/decimal.Decimal
  METHOD Cos : func() github.com/shopspring/decimal.Decimal
  METHOD Div : func(d2 github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
  METHOD DivRound : func(d2 github.com/shopspring/decimal.Decimal, precision int32) github.com/shopspring/decimal.Decimal
  METHOD Equal : func(d2 github.com/shopspring/decimal.Decimal) bool
  METHOD Equals : func(d2 github.com/shopspring/decimal.Decimal) bool
  METHOD ExpHullAbrham : func(overallPrecision uint32) (github.com/shopspring/decimal.Decimal, error)
  METHOD ExpTaylor : func(precision int32) (github.com/shopspring/decimal.Decimal, error)
  METHOD Exponent : func() int32
  METHOD Float64 : func() (f float64, exact bool)
  METHOD Floor : func() github.com/shopspring/decimal.Decimal
  METHOD GobDecode : func(data []byte) error
  METHOD GobEncode : func() ([]byte, error)
  METHOD GreaterThan : func(d2 github.com/shopspring/decimal.Decimal) bool
  METHOD GreaterThanOrEqual : func(d2 github.com/shopspring/decimal.Decimal) bool
  METHOD InexactFloat64 : func() float64
  METHOD IntPart : func() int64
  METHOD IsInteger : func() bool
  METHOD IsNegative : func() bool
  METHOD IsPositive : func() bool
  METHOD IsZero : func() bool
  METHOD LessThan : func(d2 github.com/shopspring/decimal.Decimal) bool
  METHOD LessThanOrEqual : func(d2 github.com/shopspring/decimal.Decimal) bool
  METHOD Ln : func(precision int32) (github.com/shopspring/decimal.Decimal, error)
  METHOD MarshalBinary : func() (data []byte, err error)
  METHOD MarshalJSON : func() ([]byte, error)
  METHOD MarshalText : func() (text []byte, err error)
  METHOD Mod : func(d2 github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
  METHOD Mul : func(d2 github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
  METHOD Neg : func() github.com/shopspring/decimal.Decimal
  METHOD NumDigits : func() int
  METHOD Pow : func(d2 github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
  METHOD PowBigInt : func(exp *math/big.Int) (github.com/shopspring/decimal.Decimal, error)
  METHOD PowInt32 : func(exp int32) (github.com/shopspring/decimal.Decimal, error)
  METHOD PowWithPrecision : func(d2 github.com/shopspring/decimal.Decimal, precision int32) (github.com/shopspring/decimal.Decimal, error)
  METHOD QuoRem : func(d2 github.com/shopspring/decimal.Decimal, precision int32) (github.com/shopspring/decimal.Decimal, github.com/shopspring/decimal.Decimal)
  METHOD Rat : func() *math/big.Rat
  METHOD Round : func(places int32) github.com/shopspring/decimal.Decimal
  METHOD RoundBank : func(places int32) github.com/shopspring/decimal.Decimal
  METHOD RoundCash : func(interval uint8) github.com/shopspring/decimal.Decimal
  METHOD RoundCeil : func(places int32) github.com/shopspring/decimal.Decimal
  METHOD RoundDown : func(places int32) github.com/shopspring/decimal.Decimal
  METHOD RoundFloor : func(places int32) github.com/shopspring/decimal.Decimal
  METHOD RoundUp : func(places int32) github.com/shopspring/decimal.Decimal
  METHOD Scan : func(value interface{}) error
  METHOD Shift : func(shift int32) github.com/shopspring/decimal.Decimal
  METHOD Sign : func() int
  METHOD Sin : func() github.com/shopspring/decimal.Decimal
  METHOD String : func() string
  METHOD StringFixed : func(places int32) string
  METHOD StringFixedBank : func(places int32) string
  METHOD StringFixedCash : func(interval uint8) string
  METHOD StringScaled : func(exp int32) string
  METHOD Sub : func(d2 github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
  METHOD Tan : func() github.com/shopspring/decimal.Decimal
  METHOD Truncate : func(precision int32) github.com/shopspring/decimal.Decimal
  METHOD UnmarshalBinary : func(data []byte) error
  METHOD UnmarshalJSON : func(decimalBytes []byte) error
  METHOD UnmarshalText : func(text []byte) error
  METHOD Value : func() (database/sql/driver.Value, error)
VAR Max : func(first github.com/shopspring/decimal.Decimal, rest ...github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
VAR Min : func(first github.com/shopspring/decimal.Decimal, rest ...github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
FUNC MustFromAny : func(v any) github.com/coldsmirk/vef-framework-go/decimal.Decimal
VAR New : func(value int64, exp int32) github.com/shopspring/decimal.Decimal
FUNC NewFromAny : func(v any) (github.com/coldsmirk/vef-framework-go/decimal.Decimal, error)
VAR NewFromBigInt : func(value *math/big.Int, exp int32) github.com/shopspring/decimal.Decimal
VAR NewFromBigRat : func(value *math/big.Rat, precision int32) github.com/shopspring/decimal.Decimal
VAR NewFromFloat : func(value float64) github.com/shopspring/decimal.Decimal
VAR NewFromFloat32 : func(value float32) github.com/shopspring/decimal.Decimal
VAR NewFromFloatWithExponent : func(value float64, exp int32) github.com/shopspring/decimal.Decimal
VAR NewFromFormattedString : func(value string, replRegexp *regexp.Regexp) (github.com/shopspring/decimal.Decimal, error)
VAR NewFromInt : func(value int64) github.com/shopspring/decimal.Decimal
VAR NewFromInt32 : func(value int32) github.com/shopspring/decimal.Decimal
VAR NewFromString : func(value string) (github.com/shopspring/decimal.Decimal, error)
VAR NewFromUint64 : func(value uint64) github.com/shopspring/decimal.Decimal
VAR One : github.com/shopspring/decimal.Decimal
VAR RequireFromString : func(value string) github.com/shopspring/decimal.Decimal
VAR RescalePair : func(d1 github.com/shopspring/decimal.Decimal, d2 github.com/shopspring/decimal.Decimal) (github.com/shopspring/decimal.Decimal, github.com/shopspring/decimal.Decimal)
VAR Sum : func(first github.com/shopspring/decimal.Decimal, rest ...github.com/shopspring/decimal.Decimal) github.com/shopspring/decimal.Decimal
VAR Zero : github.com/shopspring/decimal.Decimal

## github.com/coldsmirk/vef-framework-go/event
FUNC ApplyPublishOptions : func(opts []github.com/coldsmirk/vef-framework-go/event.PublishOption) github.com/coldsmirk/vef-framework-go/event.PublishConfig
FUNC ApplySubscribeOptions : func(opts []github.com/coldsmirk/vef-framework-go/event.SubscribeOption) github.com/coldsmirk/vef-framework-go/event.SubscribeConfig
FUNC AsEvents : func[T github.com/coldsmirk/vef-framework-go/event.Event](items []T) []github.com/coldsmirk/vef-framework-go/event.Event
TYPE Bus : github.com/coldsmirk/vef-framework-go/event.Bus
  METHOD Publish : func(ctx context.Context, evt github.com/coldsmirk/vef-framework-go/event.Event, opts ...github.com/coldsmirk/vef-framework-go/event.PublishOption) error
  METHOD PublishBatch : func(ctx context.Context, evts []github.com/coldsmirk/vef-framework-go/event.Event, opts ...github.com/coldsmirk/vef-framework-go/event.PublishOption) error
  METHOD Subscribe : func(eventType string, h github.com/coldsmirk/vef-framework-go/event.Handler, opts ...github.com/coldsmirk/vef-framework-go/event.SubscribeOption) (github.com/coldsmirk/vef-framework-go/event.Unsubscribe, error)
TYPE Envelope : github.com/coldsmirk/vef-framework-go/event.Envelope
  FIELD ID : string [field_order=1 tag=""]
  FIELD Type : string [field_order=2 tag=""]
  FIELD Source : string [field_order=3 tag=""]
  FIELD OccurredAt : time.Time [field_order=4 tag=""]
  FIELD PublishedAt : time.Time [field_order=5 tag=""]
  FIELD TraceID : string [field_order=6 tag=""]
  FIELD SpanID : string [field_order=7 tag=""]
  FIELD CorrelationID : string [field_order=8 tag=""]
  FIELD Headers : map[string]string [field_order=9 tag=""]
  FIELD Payload : github.com/coldsmirk/vef-framework-go/event.Event [field_order=10 tag=""]
VAR ErrAsyncQueueFull : error
VAR ErrBusAlreadyStarted : error
VAR ErrBusNotStarted : error
VAR ErrGroupRequired : error
VAR ErrHandlerPanic : error
VAR ErrInvalidEventType : error
VAR ErrNilTypeParameter : error
VAR ErrNoRouteMatched : error
VAR ErrPayloadTooLarge : error
VAR ErrQueueFull : error
VAR ErrShutdownTimeout : error
VAR ErrTransportNotFound : error
VAR ErrTxAsyncMutex : error
VAR ErrTxRequired : error
VAR ErrUnknownPayload : error
TYPE ErrorSink : github.com/coldsmirk/vef-framework-go/event.ErrorSink
TYPE Event : github.com/coldsmirk/vef-framework-go/event.Event
  METHOD EventType : func() string
TYPE Handler : github.com/coldsmirk/vef-framework-go/event.Handler
TYPE MetricsRecorder : github.com/coldsmirk/vef-framework-go/event.MetricsRecorder
  METHOD ConsumeObserved : func(eventType string, elapsed time.Duration, err error)
  METHOD PublishObserved : func(eventType string, err error)
TYPE PublishConfig : github.com/coldsmirk/vef-framework-go/event.PublishConfig
  FIELD Tx : github.com/coldsmirk/vef-framework-go/orm.DB [field_order=1 tag=""]
  FIELD Async : bool [field_order=2 tag=""]
  FIELD Source : string [field_order=3 tag=""]
  FIELD OccurredAt : time.Time [field_order=4 tag=""]
  FIELD CorrelationID : string [field_order=5 tag=""]
  FIELD Headers : map[string]string [field_order=6 tag=""]
TYPE PublishOption : github.com/coldsmirk/vef-framework-go/event.PublishOption
TYPE RawPayload : github.com/coldsmirk/vef-framework-go/event.RawPayload
  FIELD Type : string [field_order=1 tag=""]
  FIELD Body : []byte [field_order=2 tag=""]
  METHOD EventType : func() string
TYPE RouteInspector : github.com/coldsmirk/vef-framework-go/event.RouteInspector
  METHOD HasSubscribableTransport : func(eventType string) bool
  METHOD HasTransactionalRoute : func(eventType string) bool
TYPE StreamGroupInfo : github.com/coldsmirk/vef-framework-go/event.StreamGroupInfo
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Consumers : int64 [field_order=2 tag="json:\"consumers\""]
  FIELD Pending : int64 [field_order=3 tag="json:\"pending\""]
  FIELD Lag : int64 [field_order=4 tag="json:\"lag\""]
  FIELD LastDeliveredID : string [field_order=5 tag="json:\"lastDeliveredId\""]
TYPE StreamInfo : github.com/coldsmirk/vef-framework-go/event.StreamInfo
  FIELD Stream : string [field_order=1 tag="json:\"stream\""]
  FIELD Length : int64 [field_order=2 tag="json:\"length\""]
  FIELD Groups : []github.com/coldsmirk/vef-framework-go/event.StreamGroupInfo [field_order=3 tag="json:\"groups\""]
TYPE StreamInspector : github.com/coldsmirk/vef-framework-go/event.StreamInspector
  METHOD Streams : func(ctx context.Context) ([]github.com/coldsmirk/vef-framework-go/event.StreamInfo, error)
TYPE SubscribeConfig : github.com/coldsmirk/vef-framework-go/event.SubscribeConfig
  FIELD Group : string [field_order=1 tag=""]
  FIELD Concurrency : int [field_order=2 tag=""]
TYPE SubscribeOption : github.com/coldsmirk/vef-framework-go/event.SubscribeOption
FUNC SubscribeTyped : func[T github.com/coldsmirk/vef-framework-go/event.Event](b github.com/coldsmirk/vef-framework-go/event.Bus, h github.com/coldsmirk/vef-framework-go/event.TypedHandler[T], opts ...github.com/coldsmirk/vef-framework-go/event.SubscribeOption) (github.com/coldsmirk/vef-framework-go/event.Unsubscribe, error)
TYPE TypedHandler : github.com/coldsmirk/vef-framework-go/event.TypedHandler[T github.com/coldsmirk/vef-framework-go/event.Event]
TYPE Unsubscribe : github.com/coldsmirk/vef-framework-go/event.Unsubscribe
FUNC WithAsync : func() github.com/coldsmirk/vef-framework-go/event.PublishOption
FUNC WithConcurrency : func(n int) github.com/coldsmirk/vef-framework-go/event.SubscribeOption
FUNC WithCorrelationID : func(id string) github.com/coldsmirk/vef-framework-go/event.PublishOption
FUNC WithGroup : func(name string) github.com/coldsmirk/vef-framework-go/event.SubscribeOption
FUNC WithHeaders : func(h map[string]string) github.com/coldsmirk/vef-framework-go/event.PublishOption
FUNC WithOccurredAt : func(t time.Time) github.com/coldsmirk/vef-framework-go/event.PublishOption
FUNC WithSource : func(src string) github.com/coldsmirk/vef-framework-go/event.PublishOption
FUNC WithTx : func(tx github.com/coldsmirk/vef-framework-go/orm.DB) github.com/coldsmirk/vef-framework-go/event.PublishOption

## github.com/coldsmirk/vef-framework-go/event/inbox
TYPE AcquireResult : github.com/coldsmirk/vef-framework-go/event/inbox.AcquireResult
CONST AcquireResultAcquired : github.com/coldsmirk/vef-framework-go/event/inbox.AcquireResult = "acquired"
CONST AcquireResultCompleted : github.com/coldsmirk/vef-framework-go/event/inbox.AcquireResult = "completed"
CONST AcquireResultInProgress : github.com/coldsmirk/vef-framework-go/event/inbox.AcquireResult = "in_progress"
VAR ErrInProgress : error
VAR ErrLockLost : error
VAR ErrMissingLockID : error
VAR ErrUnknownAcquireResult : error
TYPE Record : github.com/coldsmirk/vef-framework-go/event/inbox.Record
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:sys_event_inbox,alias:sei\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD CreationTrackedModel : github.com/coldsmirk/vef-framework-go/orm.CreationTrackedModel [field_order=3 tag=""]
  FIELD EventID : string [field_order=4 tag="json:\"eventId\" bun:\"event_id\""]
  FIELD ConsumerGroup : string [field_order=5 tag="json:\"consumerGroup\" bun:\"consumer_group\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/event/inbox.Status [field_order=6 tag="json:\"status\" bun:\"status\""]
  FIELD LockID : string [field_order=7 tag="json:\"lockId,omitempty\" bun:\"lock_id,nullzero\""]
  FIELD LockedUntil : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=8 tag="json:\"lockedUntil,omitempty\" bun:\"locked_until,nullzero\""]
  FIELD CompletedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=9 tag="json:\"completedAt,omitempty\" bun:\"completed_at,nullzero\""]
TYPE Repository : github.com/coldsmirk/vef-framework-go/event/inbox.Repository
  METHOD Acquire : func(ctx context.Context, consumerGroup string, eventID string, lockUntil github.com/coldsmirk/vef-framework-go/timex.DateTime) (github.com/coldsmirk/vef-framework-go/event/inbox.AcquireResult, string, error)
  METHOD DeleteOlderThan : func(ctx context.Context, cutoff github.com/coldsmirk/vef-framework-go/timex.DateTime) (int64, error)
  METHOD MarkCompleted : func(ctx context.Context, consumerGroup string, eventID string, lockID string) error
  METHOD Release : func(ctx context.Context, consumerGroup string, eventID string, lockID string) error
TYPE Status : github.com/coldsmirk/vef-framework-go/event/inbox.Status
CONST StatusCompleted : github.com/coldsmirk/vef-framework-go/event/inbox.Status = "completed"
CONST StatusProcessing : github.com/coldsmirk/vef-framework-go/event/inbox.Status = "processing"

## github.com/coldsmirk/vef-framework-go/event/middleware
FUNC ChainConsume : func(mws []github.com/coldsmirk/vef-framework-go/event/middleware.ConsumeMiddleware, caps github.com/coldsmirk/vef-framework-go/event/transport.Capabilities, base github.com/coldsmirk/vef-framework-go/event/middleware.ConsumeHandler) github.com/coldsmirk/vef-framework-go/event/middleware.ConsumeHandler
FUNC ChainPublish : func(mws []github.com/coldsmirk/vef-framework-go/event/middleware.PublishMiddleware, base github.com/coldsmirk/vef-framework-go/event/middleware.PublishHandler) github.com/coldsmirk/vef-framework-go/event/middleware.PublishHandler
TYPE ConsumeHandler : github.com/coldsmirk/vef-framework-go/event/middleware.ConsumeHandler
TYPE ConsumeMiddleware : github.com/coldsmirk/vef-framework-go/event/middleware.ConsumeMiddleware
  METHOD Applies : func(caps github.com/coldsmirk/vef-framework-go/event/transport.Capabilities) bool
  METHOD Name : func() string
  METHOD Order : func() int
  METHOD WrapConsume : func(next github.com/coldsmirk/vef-framework-go/event/middleware.ConsumeHandler) github.com/coldsmirk/vef-framework-go/event/middleware.ConsumeHandler
FUNC IncomingTraceIDFromContext : func(ctx context.Context) string
CONST OrderInbox : untyped int = 100
CONST OrderLogging : untyped int = -25
CONST OrderMetrics : untyped int = 0
CONST OrderRecover : untyped int = -100
CONST OrderTracing : untyped int = -50
TYPE PublishHandler : github.com/coldsmirk/vef-framework-go/event/middleware.PublishHandler
TYPE PublishMiddleware : github.com/coldsmirk/vef-framework-go/event/middleware.PublishMiddleware
  METHOD Name : func() string
  METHOD Order : func() int
  METHOD WrapPublish : func(next github.com/coldsmirk/vef-framework-go/event/middleware.PublishHandler) github.com/coldsmirk/vef-framework-go/event/middleware.PublishHandler
FUNC TraceIDFromContext : func(ctx context.Context) string
FUNC WithIncomingTraceID : func(ctx context.Context, traceID string) context.Context
FUNC WithTraceID : func(ctx context.Context, traceID string) context.Context

## github.com/coldsmirk/vef-framework-go/event/transport
TYPE Capabilities : github.com/coldsmirk/vef-framework-go/event/transport.Capabilities
  FIELD Durable : bool [field_order=1 tag=""]
  FIELD Transactional : bool [field_order=2 tag=""]
  FIELD Ordered : bool [field_order=3 tag=""]
  FIELD AtLeastOnce : bool [field_order=4 tag=""]
  FIELD SupportsGroups : bool [field_order=5 tag=""]
  FIELD PublishOnly : bool [field_order=6 tag=""]
TYPE ConsumeFunc : github.com/coldsmirk/vef-framework-go/event/transport.ConsumeFunc
TYPE Delivery : github.com/coldsmirk/vef-framework-go/event/transport.Delivery
  METHOD Ack : func(ctx context.Context) error
  METHOD Attempt : func() int
  METHOD Frame : func() github.com/coldsmirk/vef-framework-go/event/transport.Frame
  METHOD Nack : func(ctx context.Context, retryAfter time.Duration, err error) error
VAR ErrSubscribeUnsupported : error
VAR EventTypePattern : *regexp.Regexp
TYPE Frame : github.com/coldsmirk/vef-framework-go/event/transport.Frame
  FIELD ID : string [field_order=1 tag=""]
  FIELD Type : string [field_order=2 tag=""]
  FIELD Source : string [field_order=3 tag=""]
  FIELD OccurredAt : time.Time [field_order=4 tag=""]
  FIELD PublishedAt : time.Time [field_order=5 tag=""]
  FIELD TraceID : string [field_order=6 tag=""]
  FIELD SpanID : string [field_order=7 tag=""]
  FIELD CorrelationID : string [field_order=8 tag=""]
  FIELD Headers : map[string]string [field_order=9 tag=""]
  FIELD Body : []byte [field_order=10 tag=""]
TYPE SubscribeConfig : github.com/coldsmirk/vef-framework-go/event/transport.SubscribeConfig
  FIELD Group : string [field_order=1 tag=""]
  FIELD Concurrency : int [field_order=2 tag=""]
TYPE Transport : github.com/coldsmirk/vef-framework-go/event/transport.Transport
  METHOD Capabilities : func() github.com/coldsmirk/vef-framework-go/event/transport.Capabilities
  METHOD Name : func() string
  METHOD Publish : func(ctx context.Context, frames []github.com/coldsmirk/vef-framework-go/event/transport.Frame) error
  METHOD Start : func(ctx context.Context) error
  METHOD Stop : func(ctx context.Context) error
  METHOD Subscribe : func(eventType string, group string, fn github.com/coldsmirk/vef-framework-go/event/transport.ConsumeFunc, cfg github.com/coldsmirk/vef-framework-go/event/transport.SubscribeConfig) (github.com/coldsmirk/vef-framework-go/event/transport.Unsubscribe, error)
TYPE TxTransport : github.com/coldsmirk/vef-framework-go/event/transport.TxTransport
  METHOD Capabilities : func() github.com/coldsmirk/vef-framework-go/event/transport.Capabilities
  METHOD Name : func() string
  METHOD Publish : func(ctx context.Context, frames []github.com/coldsmirk/vef-framework-go/event/transport.Frame) error
  METHOD PublishTx : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, frames []github.com/coldsmirk/vef-framework-go/event/transport.Frame) error
  METHOD Start : func(ctx context.Context) error
  METHOD Stop : func(ctx context.Context) error
  METHOD Subscribe : func(eventType string, group string, fn github.com/coldsmirk/vef-framework-go/event/transport.ConsumeFunc, cfg github.com/coldsmirk/vef-framework-go/event/transport.SubscribeConfig) (github.com/coldsmirk/vef-framework-go/event/transport.Unsubscribe, error)
TYPE Unsubscribe : github.com/coldsmirk/vef-framework-go/event/transport.Unsubscribe

## github.com/coldsmirk/vef-framework-go/event/transport/memory
TYPE Config : github.com/coldsmirk/vef-framework-go/event/transport/memory.Config
  FIELD QueueSize : int [field_order=1 tag=""]
  FIELD FullPolicy : github.com/coldsmirk/vef-framework-go/event/transport/memory.FullPolicy [field_order=2 tag=""]
  FIELD PublishTimeout : time.Duration [field_order=3 tag=""]
  METHOD EffectiveFullPolicy : func() github.com/coldsmirk/vef-framework-go/event/transport/memory.FullPolicy
  METHOD EffectiveQueueSize : func() int
TYPE FullPolicy : github.com/coldsmirk/vef-framework-go/event/transport/memory.FullPolicy
CONST FullPolicyBlock : github.com/coldsmirk/vef-framework-go/event/transport/memory.FullPolicy = "block"
CONST FullPolicyDropOldest : github.com/coldsmirk/vef-framework-go/event/transport/memory.FullPolicy = "drop_oldest"
CONST FullPolicyError : github.com/coldsmirk/vef-framework-go/event/transport/memory.FullPolicy = "error"
CONST Name : untyped string = "memory"

## github.com/coldsmirk/vef-framework-go/event/transport/outbox
TYPE Config : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Config
  FIELD RelayInterval : time.Duration [field_order=1 tag=""]
  FIELD MaxRetries : int [field_order=2 tag=""]
  FIELD BatchSize : int [field_order=3 tag=""]
  FIELD LeaseMultiplier : int [field_order=4 tag=""]
  FIELD MinLease : time.Duration [field_order=5 tag=""]
  FIELD SinkName : string [field_order=6 tag=""]
  METHOD EffectiveBatchSize : func() int
  METHOD EffectiveLeaseMultiplier : func() int
  METHOD EffectiveMaxRetries : func() int
  METHOD EffectiveMinLease : func() time.Duration
  METHOD EffectiveRelayInterval : func() time.Duration
  METHOD EffectiveSinkName : func() string
CONST Name : untyped string = "outbox"
TYPE Record : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Record
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:sys_event_outbox,alias:seo\""]
  FIELD Model : github.com/coldsmirk/vef-framework-go/orm.Model [field_order=2 tag=""]
  FIELD CreationTrackedModel : github.com/coldsmirk/vef-framework-go/orm.CreationTrackedModel [field_order=3 tag=""]
  FIELD EventID : string [field_order=4 tag="json:\"eventId\" bun:\"event_id\""]
  FIELD EventType : string [field_order=5 tag="json:\"eventType\" bun:\"event_type\""]
  FIELD Source : string [field_order=6 tag="json:\"source\" bun:\"source\""]
  FIELD TraceID : string [field_order=7 tag="json:\"traceId,omitempty\" bun:\"trace_id,nullzero\""]
  FIELD SpanID : string [field_order=8 tag="json:\"spanId,omitempty\" bun:\"span_id,nullzero\""]
  FIELD CorrelationID : string [field_order=9 tag="json:\"correlationId,omitempty\" bun:\"correlation_id,nullzero\""]
  FIELD Headers : map[string]string [field_order=10 tag="json:\"headers,omitempty\" bun:\"headers,type:jsonb,nullzero\""]
  FIELD Payload : encoding/json.RawMessage [field_order=11 tag="json:\"payload\" bun:\"payload,type:jsonb\""]
  FIELD Status : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Status [field_order=12 tag="json:\"status\" bun:\"status\""]
  FIELD RetryCount : int [field_order=13 tag="json:\"retryCount\" bun:\"retry_count\""]
  FIELD LastError : *string [field_order=14 tag="json:\"lastError,omitempty\" bun:\"last_error,nullzero\""]
  FIELD ProcessedAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=15 tag="json:\"processedAt,omitempty\" bun:\"processed_at,nullzero\""]
  FIELD RetryAfter : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=16 tag="json:\"retryAfter,omitempty\" bun:\"retry_after,nullzero\""]
  FIELD OccurredAt : github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=17 tag="json:\"occurredAt\" bun:\"occurred_at\""]
TYPE Repository : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Repository
  METHOD ClaimBatch : func(ctx context.Context, batchSize int, maxRetries int, leaseUntil github.com/coldsmirk/vef-framework-go/timex.DateTime) ([]github.com/coldsmirk/vef-framework-go/event/transport/outbox.Record, error)
  METHOD DeleteCompletedOlderThan : func(ctx context.Context, cutoff github.com/coldsmirk/vef-framework-go/timex.DateTime) (int64, error)
  METHOD InsertBatch : func(ctx context.Context, records []github.com/coldsmirk/vef-framework-go/event/transport/outbox.Record) error
  METHOD InsertBatchTx : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, records []github.com/coldsmirk/vef-framework-go/event/transport/outbox.Record) error
  METHOD MarkCompleted : func(ctx context.Context, id string) error
  METHOD MarkFailed : func(ctx context.Context, id string, errMsg string, retryCount int, retryAfter github.com/coldsmirk/vef-framework-go/timex.DateTime, maxRetries int) error
TYPE Status : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Status
CONST StatusCompleted : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Status = "completed"
CONST StatusDead : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Status = "dead"
CONST StatusFailed : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Status = "failed"
CONST StatusPending : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Status = "pending"
CONST StatusProcessing : github.com/coldsmirk/vef-framework-go/event/transport/outbox.Status = "processing"

## github.com/coldsmirk/vef-framework-go/event/transport/redisstream
TYPE Config : github.com/coldsmirk/vef-framework-go/event/transport/redisstream.Config
  FIELD StreamPrefix : string [field_order=1 tag=""]
  FIELD MaxLenApprox : int64 [field_order=2 tag=""]
  FIELD BlockTimeout : time.Duration [field_order=3 tag=""]
  FIELD ClaimIdle : time.Duration [field_order=4 tag=""]
  FIELD ClaimInterval : time.Duration [field_order=5 tag=""]
  FIELD ClaimBatchSize : int64 [field_order=6 tag=""]
  FIELD ReaperConcurrency : int [field_order=7 tag=""]
  FIELD HandlerTimeout : time.Duration [field_order=8 tag=""]
  FIELD SetupTimeout : time.Duration [field_order=9 tag=""]
  FIELD ConsumerID : string [field_order=10 tag=""]
  FIELD StartID : string [field_order=11 tag=""]
  FIELD IdleGroupRetention : time.Duration [field_order=12 tag=""]
  FIELD IdleGroupSweepInterval : time.Duration [field_order=13 tag=""]
  METHOD EffectiveBlockTimeout : func() time.Duration
  METHOD EffectiveClaimBatchSize : func() int64
  METHOD EffectiveClaimIdle : func() time.Duration
  METHOD EffectiveClaimInterval : func() time.Duration
  METHOD EffectiveIdleGroupSweepInterval : func() time.Duration
  METHOD EffectiveReaperConcurrency : func() int
  METHOD EffectiveSetupTimeout : func() time.Duration
  METHOD EffectiveStartID : func() string
  METHOD EffectiveStreamPrefix : func() string
  METHOD StreamKey : func(eventType string) string
CONST Name : untyped string = "redis_stream"

## github.com/coldsmirk/vef-framework-go/excel
VAR ErrSheetIndexOutOfRange : error
TYPE ExportOption : github.com/coldsmirk/vef-framework-go/excel.ExportOption
TYPE ImportOption : github.com/coldsmirk/vef-framework-go/excel.ImportOption
FUNC NewExporter : func(adapter github.com/coldsmirk/vef-framework-go/tabular.RowAdapter, opts ...github.com/coldsmirk/vef-framework-go/excel.ExportOption) github.com/coldsmirk/vef-framework-go/tabular.Exporter
FUNC NewExporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/excel.ExportOption) github.com/coldsmirk/vef-framework-go/tabular.Exporter
FUNC NewImporter : func(adapter github.com/coldsmirk/vef-framework-go/tabular.RowAdapter, opts ...github.com/coldsmirk/vef-framework-go/excel.ImportOption) github.com/coldsmirk/vef-framework-go/tabular.Importer
FUNC NewImporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/excel.ImportOption) github.com/coldsmirk/vef-framework-go/tabular.Importer
FUNC NewMapExporter : func(specs []github.com/coldsmirk/vef-framework-go/tabular.ColumnSpec, opts ...github.com/coldsmirk/vef-framework-go/excel.ExportOption) (github.com/coldsmirk/vef-framework-go/tabular.Exporter, error)
FUNC NewMapImporter : func(specs []github.com/coldsmirk/vef-framework-go/tabular.ColumnSpec, mapOpts []github.com/coldsmirk/vef-framework-go/tabular.MapOption, opts ...github.com/coldsmirk/vef-framework-go/excel.ImportOption) (github.com/coldsmirk/vef-framework-go/tabular.Importer, error)
FUNC NewTypedExporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/excel.ExportOption) github.com/coldsmirk/vef-framework-go/tabular.TypedExporter[T]
FUNC NewTypedImporterFor : func[T any](opts ...github.com/coldsmirk/vef-framework-go/excel.ImportOption) github.com/coldsmirk/vef-framework-go/tabular.TypedImporter[T]
FUNC WithImportSheetIndex : func(index int) github.com/coldsmirk/vef-framework-go/excel.ImportOption
FUNC WithImportSheetName : func(name string) github.com/coldsmirk/vef-framework-go/excel.ImportOption
FUNC WithSheetName : func(name string) github.com/coldsmirk/vef-framework-go/excel.ExportOption
FUNC WithSkipRows : func(rows int) github.com/coldsmirk/vef-framework-go/excel.ImportOption
FUNC WithoutHeader : func() github.com/coldsmirk/vef-framework-go/excel.ImportOption
FUNC WithoutTrimSpace : func() github.com/coldsmirk/vef-framework-go/excel.ImportOption

## github.com/coldsmirk/vef-framework-go/expression
FUNC AsPredicate : func() github.com/coldsmirk/vef-framework-go/expression.CompileOption
TYPE CompileOption : github.com/coldsmirk/vef-framework-go/expression.CompileOption
TYPE CompileOptions : github.com/coldsmirk/vef-framework-go/expression.CompileOptions
  FIELD Predicate : bool [field_order=1 tag=""]
FUNC DecodeValue : func[T any](v github.com/coldsmirk/vef-framework-go/expression.Value) (T, error)
TYPE Engine : github.com/coldsmirk/vef-framework-go/expression.Engine
  METHOD Compile : func(source string, opts ...github.com/coldsmirk/vef-framework-go/expression.CompileOption) (github.com/coldsmirk/vef-framework-go/expression.Program, error)
  METHOD Evaluate : func(ctx context.Context, source string, env any) (github.com/coldsmirk/vef-framework-go/expression.Value, error)
CONST ErrCodeEvaluationFailed : untyped int = 2500
VAR ErrEvaluationFailed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUnexpectedType : error
FUNC EvaluateAs : func[T any](ctx context.Context, e github.com/coldsmirk/vef-framework-go/expression.Engine, source string, env any) (T, error)
FUNC Match : func(ctx context.Context, e github.com/coldsmirk/vef-framework-go/expression.Engine, source string, env any) (bool, error)
FUNC NewValue : func(raw any) github.com/coldsmirk/vef-framework-go/expression.Value
TYPE Program : github.com/coldsmirk/vef-framework-go/expression.Program
  METHOD Run : func(ctx context.Context, env any) (github.com/coldsmirk/vef-framework-go/expression.Value, error)
  METHOD Source : func() string
TYPE Value : github.com/coldsmirk/vef-framework-go/expression.Value
  METHOD Bool : func() (bool, error)
  METHOD Decode : func(target any) error
  METHOD Interface : func() any
  METHOD IsNil : func() bool

## github.com/coldsmirk/vef-framework-go/hashx
FUNC HmacMD5 : func(key []byte, data []byte) string
FUNC HmacSHA1 : func(key []byte, data []byte) string
FUNC HmacSHA256 : func(key []byte, data []byte) string
FUNC HmacSHA512 : func(key []byte, data []byte) string
FUNC HmacSM3 : func(key []byte, data []byte) string
FUNC MD5 : func(data string) string
FUNC MD5Bytes : func(data []byte) string
FUNC SHA1 : func(data string) string
FUNC SHA1Bytes : func(data []byte) string
FUNC SHA256 : func(data string) string
FUNC SHA256Bytes : func(data []byte) string
FUNC SHA512 : func(data string) string
FUNC SHA512Bytes : func(data []byte) string
FUNC SM3 : func(data string) string
FUNC SM3Bytes : func(data []byte) string

## github.com/coldsmirk/vef-framework-go/httpx
FUNC GetIP : func(ctx github.com/gofiber/fiber/v3.Ctx) string
FUNC IsJSON : func(ctx github.com/gofiber/fiber/v3.Ctx) bool
FUNC IsMultipart : func(ctx github.com/gofiber/fiber/v3.Ctx) bool

## github.com/coldsmirk/vef-framework-go/i18n
TYPE Config : github.com/coldsmirk/vef-framework-go/i18n.Config
  FIELD Locales : embed.FS [field_order=1 tag=""]
FUNC CurrentLanguage : func() string
CONST DefaultLanguage : untyped string = "zh-CN"
VAR ErrMessageIDEmpty : error
VAR ErrUnsupportedLanguage : error
FUNC GetSupportedLanguages : func() []string
FUNC IsLanguageSupported : func(languageCode string) bool
FUNC New : func(config github.com/coldsmirk/vef-framework-go/i18n.Config) (github.com/coldsmirk/vef-framework-go/i18n.Translator, error)
FUNC SetLanguage : func(languageCode string) error
FUNC T : func(messageID string, templateData ...map[string]any) string
FUNC Te : func(messageID string, templateData ...map[string]any) (string, error)
TYPE Translator : github.com/coldsmirk/vef-framework-go/i18n.Translator
  METHOD T : func(messageID string, templateData ...map[string]any) string
  METHOD Te : func(messageID string, templateData ...map[string]any) (string, error)

## github.com/coldsmirk/vef-framework-go/i18n/locales
VAR EmbedLocales : embed.FS

## github.com/coldsmirk/vef-framework-go/id
CONST DefaultRandomIDGeneratorAlphabet : untyped string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
CONST DefaultRandomIDGeneratorLength : untyped int = 32
VAR DefaultUUIDGenerator : github.com/coldsmirk/vef-framework-go/id.IDGenerator
VAR DefaultXIDGenerator : github.com/coldsmirk/vef-framework-go/id.IDGenerator
FUNC Generate : func() string
FUNC GenerateUUID : func() string
TYPE IDGenerator : github.com/coldsmirk/vef-framework-go/id.IDGenerator
  METHOD Generate : func() string
FUNC NewRandomIDGenerator : func(opts ...github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption) github.com/coldsmirk/vef-framework-go/id.IDGenerator
FUNC NewUUIDGenerator : func() github.com/coldsmirk/vef-framework-go/id.IDGenerator
FUNC NewXIDGenerator : func() github.com/coldsmirk/vef-framework-go/id.IDGenerator
TYPE RandomIDGeneratorOption : github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption
FUNC WithAlphabet : func(alphabet string) github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption
FUNC WithLength : func(length int) github.com/coldsmirk/vef-framework-go/id.RandomIDGeneratorOption

## github.com/coldsmirk/vef-framework-go/js
TYPE AstProgram : github.com/coldsmirk/vef-framework-go/js.AstProgram
  METHOD Idx0 : func() github.com/dop251/goja/file.Idx
  METHOD Idx1 : func() github.com/dop251/goja/file.Idx
VAR Compile : func(name string, src string, strict bool) (*github.com/dop251/goja.Program, error)
VAR IsBigInt : func(v github.com/dop251/goja.Value) bool
VAR IsInfinity : func(v github.com/dop251/goja.Value) bool
VAR IsNaN : func(v github.com/dop251/goja.Value) bool
VAR IsNull : func(v github.com/dop251/goja.Value) bool
VAR IsNumber : func(v github.com/dop251/goja.Value) bool
VAR IsString : func(v github.com/dop251/goja.Value) bool
VAR IsUndefined : func(v github.com/dop251/goja.Value) bool
VAR MustCompile : func(name string, src string, strict bool) *github.com/dop251/goja.Program
FUNC New : func() (*github.com/coldsmirk/vef-framework-go/js.Runtime, error)
TYPE Object : github.com/coldsmirk/vef-framework-go/js.Object
  METHOD ClassName : func() string
  METHOD DefineAccessorProperty : func(name string, getter github.com/dop251/goja.Value, setter github.com/dop251/goja.Value, configurable github.com/dop251/goja.Flag, enumerable github.com/dop251/goja.Flag) error
  METHOD DefineAccessorPropertySymbol : func(name *github.com/dop251/goja.Symbol, getter github.com/dop251/goja.Value, setter github.com/dop251/goja.Value, configurable github.com/dop251/goja.Flag, enumerable github.com/dop251/goja.Flag) error
  METHOD DefineDataProperty : func(name string, value github.com/dop251/goja.Value, writable github.com/dop251/goja.Flag, configurable github.com/dop251/goja.Flag, enumerable github.com/dop251/goja.Flag) error
  METHOD DefineDataPropertySymbol : func(name *github.com/dop251/goja.Symbol, value github.com/dop251/goja.Value, writable github.com/dop251/goja.Flag, configurable github.com/dop251/goja.Flag, enumerable github.com/dop251/goja.Flag) error
  METHOD Delete : func(name string) error
  METHOD DeleteSymbol : func(name *github.com/dop251/goja.Symbol) error
  METHOD Equals : func(other github.com/dop251/goja.Value) bool
  METHOD Export : func() interface{}
  METHOD ExportType : func() reflect.Type
  METHOD Get : func(name string) github.com/dop251/goja.Value
  METHOD GetOwnPropertyNames : func() (keys []string)
  METHOD GetSymbol : func(sym *github.com/dop251/goja.Symbol) github.com/dop251/goja.Value
  METHOD Keys : func() (keys []string)
  METHOD MarshalJSON : func() ([]byte, error)
  METHOD Prototype : func() *github.com/dop251/goja.Object
  METHOD SameAs : func(other github.com/dop251/goja.Value) bool
  METHOD Set : func(name string, value interface{}) error
  METHOD SetPrototype : func(proto *github.com/dop251/goja.Object) error
  METHOD SetSymbol : func(name *github.com/dop251/goja.Symbol, value interface{}) error
  METHOD StrictEquals : func(other github.com/dop251/goja.Value) bool
  METHOD String : func() string
  METHOD Symbols : func() []*github.com/dop251/goja.Symbol
  METHOD ToBoolean : func() bool
  METHOD ToFloat : func() float64
  METHOD ToInteger : func() int64
  METHOD ToNumber : func() github.com/dop251/goja.Value
  METHOD ToObject : func(*github.com/dop251/goja.Runtime) *github.com/dop251/goja.Object
  METHOD ToString : func() github.com/dop251/goja.Value
  METHOD UnmarshalJSON : func([]byte) error
FUNC Parse : func(name string, src string) (*github.com/coldsmirk/vef-framework-go/js.AstProgram, error)
TYPE Program : github.com/coldsmirk/vef-framework-go/js.Program
TYPE Runtime : github.com/coldsmirk/vef-framework-go/js.Runtime
  METHOD CaptureCallStack : func(depth int, stack []github.com/dop251/goja.StackFrame) []github.com/dop251/goja.StackFrame
  METHOD ClearInterrupt : func()
  METHOD CreateObject : func(proto *github.com/dop251/goja.Object) *github.com/dop251/goja.Object
  METHOD ExportTo : func(v github.com/dop251/goja.Value, target interface{}) error
  METHOD ForOf : func(iterable github.com/dop251/goja.Value, step func(curValue github.com/dop251/goja.Value) (continueIteration bool))
  METHOD Get : func(name string) github.com/dop251/goja.Value
  METHOD GlobalObject : func() *github.com/dop251/goja.Object
  METHOD InstanceOf : func(left github.com/dop251/goja.Value, right *github.com/dop251/goja.Object) (res bool)
  METHOD Interrupt : func(v interface{})
  METHOD New : func(construct github.com/dop251/goja.Value, args ...github.com/dop251/goja.Value) (o *github.com/dop251/goja.Object, err error)
  METHOD NewArray : func(items ...interface{}) *github.com/dop251/goja.Object
  METHOD NewArrayBuffer : func(data []byte) github.com/dop251/goja.ArrayBuffer
  METHOD NewDynamicArray : func(a github.com/dop251/goja.DynamicArray) *github.com/dop251/goja.Object
  METHOD NewDynamicObject : func(d github.com/dop251/goja.DynamicObject) *github.com/dop251/goja.Object
  METHOD NewGoError : func(err error) *github.com/dop251/goja.Object
  METHOD NewObject : func() (v *github.com/dop251/goja.Object)
  METHOD NewPromise : func() (promise *github.com/dop251/goja.Promise, resolve func(reason interface{}) error, reject func(reason interface{}) error)
  METHOD NewProxy : func(target *github.com/dop251/goja.Object, nativeHandler *github.com/dop251/goja.ProxyTrapConfig) github.com/dop251/goja.Proxy
  METHOD NewTypeError : func(args ...interface{}) *github.com/dop251/goja.Object
  METHOD RunProgram : func(p *github.com/dop251/goja.Program) (result github.com/dop251/goja.Value, err error)
  METHOD RunScript : func(name string, src string) (github.com/dop251/goja.Value, error)
  METHOD RunString : func(str string) (github.com/dop251/goja.Value, error)
  METHOD Set : func(name string, value interface{}) error
  METHOD SetAsyncContextTracker : func(tracker github.com/dop251/goja.AsyncContextTracker)
  METHOD SetFieldNameMapper : func(mapper github.com/dop251/goja.FieldNameMapper)
  METHOD SetGlobalObject : func(object *github.com/dop251/goja.Object)
  METHOD SetMaxCallStackSize : func(size int)
  METHOD SetParserOptions : func(opts ...github.com/dop251/goja/parser.Option)
  METHOD SetPromiseRejectionTracker : func(tracker github.com/dop251/goja.PromiseRejectionTracker)
  METHOD SetRandSource : func(source github.com/dop251/goja.RandSource)
  METHOD SetTimeSource : func(now github.com/dop251/goja.Now)
  METHOD ToValue : func(i interface{}) github.com/dop251/goja.Value
  METHOD Try : func(f func()) *github.com/dop251/goja.Exception
TYPE Value : github.com/coldsmirk/vef-framework-go/js.Value
  METHOD Equals : func(github.com/dop251/goja.Value) bool
  METHOD Export : func() interface{}
  METHOD ExportType : func() reflect.Type
  METHOD SameAs : func(github.com/dop251/goja.Value) bool
  METHOD StrictEquals : func(github.com/dop251/goja.Value) bool
  METHOD String : func() string
  METHOD ToBoolean : func() bool
  METHOD ToFloat : func() float64
  METHOD ToInteger : func() int64
  METHOD ToNumber : func() github.com/dop251/goja.Value
  METHOD ToObject : func(*github.com/dop251/goja.Runtime) *github.com/dop251/goja.Object
  METHOD ToString : func() github.com/dop251/goja.Value

## github.com/coldsmirk/vef-framework-go/logx
TYPE Level : github.com/coldsmirk/vef-framework-go/logx.Level
  METHOD String : func() string
CONST LevelDebug : github.com/coldsmirk/vef-framework-go/logx.Level = 1
CONST LevelError : github.com/coldsmirk/vef-framework-go/logx.Level = 4
CONST LevelInfo : github.com/coldsmirk/vef-framework-go/logx.Level = 2
CONST LevelPanic : github.com/coldsmirk/vef-framework-go/logx.Level = 5
CONST LevelWarn : github.com/coldsmirk/vef-framework-go/logx.Level = 3
TYPE Logger : github.com/coldsmirk/vef-framework-go/logx.Logger
  METHOD Debug : func(message string)
  METHOD Debugf : func(template string, args ...any)
  METHOD Enabled : func(level github.com/coldsmirk/vef-framework-go/logx.Level) bool
  METHOD Error : func(message string)
  METHOD Errorf : func(template string, args ...any)
  METHOD Info : func(message string)
  METHOD Infof : func(template string, args ...any)
  METHOD Named : func(name string) github.com/coldsmirk/vef-framework-go/logx.Logger
  METHOD Panic : func(message string)
  METHOD Panicf : func(template string, args ...any)
  METHOD Sync : func()
  METHOD Warn : func(message string)
  METHOD Warnf : func(template string, args ...any)
  METHOD WithCallerSkip : func(skip int) github.com/coldsmirk/vef-framework-go/logx.Logger
TYPE LoggerConfigurable : github.com/coldsmirk/vef-framework-go/logx.LoggerConfigurable[T any]
  METHOD WithLogger : func(logger github.com/coldsmirk/vef-framework-go/logx.Logger) T

## github.com/coldsmirk/vef-framework-go/mapx
VAR DecoderHook : github.com/go-viper/mapstructure/v2.DecodeHookFunc
TYPE DecoderOption : github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
VAR ErrCollectionSetIncompatibleKind : error
VAR ErrCollectionSetNegative : error
VAR ErrCollectionSetNilElement : error
VAR ErrCollectionSetNonInteger : error
VAR ErrCollectionSetNotFinite : error
VAR ErrCollectionSetOverflow : error
VAR ErrCollectionSetUnsupportedTarget : error
VAR ErrInvalidFromMapType : error
VAR ErrInvalidToMapValue : error
FUNC FromMap : func[T any](value map[string]any, options ...github.com/coldsmirk/vef-framework-go/mapx.DecoderOption) (*T, error)
TYPE Metadata : github.com/coldsmirk/vef-framework-go/mapx.Metadata
FUNC NewDecoder : func(result any, options ...github.com/coldsmirk/vef-framework-go/mapx.DecoderOption) (*github.com/go-viper/mapstructure/v2.Decoder, error)
FUNC ToMap : func(value any, options ...github.com/coldsmirk/vef-framework-go/mapx.DecoderOption) (map[string]any, error)
FUNC WithAllowUnsetPointer : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithDecodeHook : func(decodeHook github.com/go-viper/mapstructure/v2.DecodeHookFunc) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithDecodeNil : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithErrorUnset : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithErrorUnused : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithIgnoreUntaggedFields : func(ignoreUntaggedFields bool) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithMatchName : func(matchName func(mapKey string, fieldName string) bool) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithMetadata : func(metadata *github.com/coldsmirk/vef-framework-go/mapx.Metadata) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithTagName : func(tagName string) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithWeaklyTypedInput : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption
FUNC WithZeroFields : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption

## github.com/coldsmirk/vef-framework-go/mcp
TYPE Annotations : github.com/coldsmirk/vef-framework-go/mcp.Annotations
TYPE AudioContent : github.com/coldsmirk/vef-framework-go/mcp.AudioContent
  METHOD MarshalJSON : func() ([]byte, error)
TYPE CallToolRequest : github.com/coldsmirk/vef-framework-go/mcp.CallToolRequest
  METHOD GetExtra : func() *github.com/modelcontextprotocol/go-sdk/mcp.RequestExtra
  METHOD GetParams : func() github.com/modelcontextprotocol/go-sdk/mcp.Params
  METHOD GetSession : func() github.com/modelcontextprotocol/go-sdk/mcp.Session
TYPE CallToolResult : github.com/coldsmirk/vef-framework-go/mcp.CallToolResult
  METHOD GetError : func() error
  METHOD GetMeta : func() map[string]any
  METHOD SetError : func(err error)
  METHOD SetMeta : func(x map[string]any)
  METHOD UnmarshalJSON : func(data []byte) error
TYPE Content : github.com/coldsmirk/vef-framework-go/mcp.Content
  METHOD MarshalJSON : func() ([]byte, error)
FUNC DBWithOperator : func(ctx context.Context, db github.com/coldsmirk/vef-framework-go/orm.DB) github.com/coldsmirk/vef-framework-go/orm.DB
FUNC GetPrincipalFromContext : func(ctx context.Context) *github.com/coldsmirk/vef-framework-go/security.Principal
TYPE GetPromptParams : github.com/coldsmirk/vef-framework-go/mcp.GetPromptParams
  METHOD GetMeta : func() map[string]any
  METHOD GetProgressToken : func() any
  METHOD SetMeta : func(x map[string]any)
  METHOD SetProgressToken : func(t any)
TYPE GetPromptRequest : github.com/coldsmirk/vef-framework-go/mcp.GetPromptRequest
  METHOD GetExtra : func() *github.com/modelcontextprotocol/go-sdk/mcp.RequestExtra
  METHOD GetParams : func() github.com/modelcontextprotocol/go-sdk/mcp.Params
  METHOD GetSession : func() github.com/modelcontextprotocol/go-sdk/mcp.Session
TYPE GetPromptResult : github.com/coldsmirk/vef-framework-go/mcp.GetPromptResult
  METHOD GetMeta : func() map[string]any
  METHOD SetMeta : func(x map[string]any)
TYPE ImageContent : github.com/coldsmirk/vef-framework-go/mcp.ImageContent
  METHOD MarshalJSON : func() ([]byte, error)
TYPE Implementation : github.com/coldsmirk/vef-framework-go/mcp.Implementation
FUNC MustSchemaFor : func[T any]() map[string]any
FUNC MustSchemaOf : func(v any) map[string]any
FUNC NewToolResultError : func(errMsg string) *github.com/coldsmirk/vef-framework-go/mcp.CallToolResult
FUNC NewToolResultText : func(text string) *github.com/coldsmirk/vef-framework-go/mcp.CallToolResult
TYPE Prompt : github.com/coldsmirk/vef-framework-go/mcp.Prompt
  METHOD GetMeta : func() map[string]any
  METHOD SetMeta : func(x map[string]any)
TYPE PromptArgument : github.com/coldsmirk/vef-framework-go/mcp.PromptArgument
TYPE PromptDefinition : github.com/coldsmirk/vef-framework-go/mcp.PromptDefinition
  FIELD Prompt : *github.com/coldsmirk/vef-framework-go/mcp.Prompt [field_order=1 tag=""]
  FIELD Handler : github.com/coldsmirk/vef-framework-go/mcp.PromptHandler [field_order=2 tag=""]
TYPE PromptHandler : github.com/coldsmirk/vef-framework-go/mcp.PromptHandler
TYPE PromptMessage : github.com/coldsmirk/vef-framework-go/mcp.PromptMessage
  METHOD UnmarshalJSON : func(data []byte) error
TYPE PromptProvider : github.com/coldsmirk/vef-framework-go/mcp.PromptProvider
  METHOD Prompts : func() []github.com/coldsmirk/vef-framework-go/mcp.PromptDefinition
TYPE ReadResourceRequest : github.com/coldsmirk/vef-framework-go/mcp.ReadResourceRequest
  METHOD GetExtra : func() *github.com/modelcontextprotocol/go-sdk/mcp.RequestExtra
  METHOD GetParams : func() github.com/modelcontextprotocol/go-sdk/mcp.Params
  METHOD GetSession : func() github.com/modelcontextprotocol/go-sdk/mcp.Session
TYPE ReadResourceResult : github.com/coldsmirk/vef-framework-go/mcp.ReadResourceResult
  METHOD GetMeta : func() map[string]any
  METHOD SetMeta : func(x map[string]any)
TYPE Resource : github.com/coldsmirk/vef-framework-go/mcp.Resource
  METHOD GetMeta : func() map[string]any
  METHOD SetMeta : func(x map[string]any)
TYPE ResourceDefinition : github.com/coldsmirk/vef-framework-go/mcp.ResourceDefinition
  FIELD Resource : *github.com/coldsmirk/vef-framework-go/mcp.Resource [field_order=1 tag=""]
  FIELD Handler : github.com/coldsmirk/vef-framework-go/mcp.ResourceHandler [field_order=2 tag=""]
TYPE ResourceHandler : github.com/coldsmirk/vef-framework-go/mcp.ResourceHandler
VAR ResourceNotFoundError : func(uri string) error
TYPE ResourceProvider : github.com/coldsmirk/vef-framework-go/mcp.ResourceProvider
  METHOD Resources : func() []github.com/coldsmirk/vef-framework-go/mcp.ResourceDefinition
TYPE ResourceTemplate : github.com/coldsmirk/vef-framework-go/mcp.ResourceTemplate
  METHOD GetMeta : func() map[string]any
  METHOD SetMeta : func(x map[string]any)
TYPE ResourceTemplateDefinition : github.com/coldsmirk/vef-framework-go/mcp.ResourceTemplateDefinition
  FIELD Template : *github.com/coldsmirk/vef-framework-go/mcp.ResourceTemplate [field_order=1 tag=""]
  FIELD Handler : github.com/coldsmirk/vef-framework-go/mcp.ResourceHandler [field_order=2 tag=""]
TYPE ResourceTemplateProvider : github.com/coldsmirk/vef-framework-go/mcp.ResourceTemplateProvider
  METHOD ResourceTemplates : func() []github.com/coldsmirk/vef-framework-go/mcp.ResourceTemplateDefinition
TYPE Role : github.com/coldsmirk/vef-framework-go/mcp.Role
FUNC SchemaFor : func[T any]() map[string]any
FUNC SchemaOf : func(v any) map[string]any
TYPE Server : github.com/coldsmirk/vef-framework-go/mcp.Server
  METHOD AddPrompt : func(p *github.com/modelcontextprotocol/go-sdk/mcp.Prompt, h github.com/modelcontextprotocol/go-sdk/mcp.PromptHandler)
  METHOD AddReceivingMiddleware : func(middleware ...github.com/modelcontextprotocol/go-sdk/mcp.Middleware)
  METHOD AddResource : func(r *github.com/modelcontextprotocol/go-sdk/mcp.Resource, h github.com/modelcontextprotocol/go-sdk/mcp.ResourceHandler)
  METHOD AddResourceTemplate : func(t *github.com/modelcontextprotocol/go-sdk/mcp.ResourceTemplate, h github.com/modelcontextprotocol/go-sdk/mcp.ResourceHandler)
  METHOD AddSendingMiddleware : func(middleware ...github.com/modelcontextprotocol/go-sdk/mcp.Middleware)
  METHOD AddTool : func(t *github.com/modelcontextprotocol/go-sdk/mcp.Tool, h github.com/modelcontextprotocol/go-sdk/mcp.ToolHandler)
  METHOD Connect : func(ctx context.Context, t github.com/modelcontextprotocol/go-sdk/mcp.Transport, opts *github.com/modelcontextprotocol/go-sdk/mcp.ServerSessionOptions) (*github.com/modelcontextprotocol/go-sdk/mcp.ServerSession, error)
  METHOD RemovePrompts : func(names ...string)
  METHOD RemoveResourceTemplates : func(uriTemplates ...string)
  METHOD RemoveResources : func(uris ...string)
  METHOD RemoveTools : func(names ...string)
  METHOD ResourceUpdated : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.ResourceUpdatedNotificationParams) error
  METHOD Run : func(ctx context.Context, t github.com/modelcontextprotocol/go-sdk/mcp.Transport) error
  METHOD Sessions : func() iter.Seq[*github.com/modelcontextprotocol/go-sdk/mcp.ServerSession]
TYPE ServerInfo : github.com/coldsmirk/vef-framework-go/mcp.ServerInfo
  FIELD Name : string [field_order=1 tag=""]
  FIELD Version : string [field_order=2 tag=""]
  FIELD Instructions : string [field_order=3 tag=""]
TYPE ServerOptions : github.com/coldsmirk/vef-framework-go/mcp.ServerOptions
TYPE ServerSession : github.com/coldsmirk/vef-framework-go/mcp.ServerSession
  METHOD Close : func() error
  METHOD CreateMessage : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.CreateMessageParams) (*github.com/modelcontextprotocol/go-sdk/mcp.CreateMessageResult, error)
  METHOD CreateMessageWithTools : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.CreateMessageWithToolsParams) (*github.com/modelcontextprotocol/go-sdk/mcp.CreateMessageWithToolsResult, error)
  METHOD Elicit : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.ElicitParams) (*github.com/modelcontextprotocol/go-sdk/mcp.ElicitResult, error)
  METHOD ID : func() string
  METHOD InitializeParams : func() *github.com/modelcontextprotocol/go-sdk/mcp.InitializeParams
  METHOD ListRoots : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.ListRootsParams) (*github.com/modelcontextprotocol/go-sdk/mcp.ListRootsResult, error)
  METHOD Log : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.LoggingMessageParams) error
  METHOD NotifyProgress : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.ProgressNotificationParams) error
  METHOD Ping : func(ctx context.Context, params *github.com/modelcontextprotocol/go-sdk/mcp.PingParams) error
  METHOD Wait : func() error
TYPE TextContent : github.com/coldsmirk/vef-framework-go/mcp.TextContent
  METHOD MarshalJSON : func() ([]byte, error)
TYPE Tool : github.com/coldsmirk/vef-framework-go/mcp.Tool
  METHOD GetMeta : func() map[string]any
  METHOD SetMeta : func(x map[string]any)
TYPE ToolDefinition : github.com/coldsmirk/vef-framework-go/mcp.ToolDefinition
  FIELD Tool : *github.com/coldsmirk/vef-framework-go/mcp.Tool [field_order=1 tag=""]
  FIELD Handler : github.com/coldsmirk/vef-framework-go/mcp.ToolHandler [field_order=2 tag=""]
TYPE ToolHandler : github.com/coldsmirk/vef-framework-go/mcp.ToolHandler
TYPE ToolProvider : github.com/coldsmirk/vef-framework-go/mcp.ToolProvider
  METHOD Tools : func() []github.com/coldsmirk/vef-framework-go/mcp.ToolDefinition

## github.com/coldsmirk/vef-framework-go/middleware
TYPE SPAConfig : github.com/coldsmirk/vef-framework-go/middleware.SPAConfig
  FIELD Path : string [field_order=1 tag=""]
  FIELD Fs : io/fs.FS [field_order=2 tag=""]
  FIELD ExcludePaths : []string [field_order=3 tag=""]

## github.com/coldsmirk/vef-framework-go/mold
TYPE CachedDictionaryResolver : github.com/coldsmirk/vef-framework-go/mold.CachedDictionaryResolver
  METHOD Resolve : func(ctx context.Context, key string, code string) (string, error)
TYPE DictionaryChangedEvent : github.com/coldsmirk/vef-framework-go/mold.DictionaryChangedEvent
  FIELD Keys : []string [field_order=1 tag="json:\"keys\""]
  METHOD EventType : func() string
TYPE DictionaryLoader : github.com/coldsmirk/vef-framework-go/mold.DictionaryLoader
  METHOD Load : func(ctx context.Context, key string) (map[string]string, error)
TYPE DictionaryLoaderFunc : github.com/coldsmirk/vef-framework-go/mold.DictionaryLoaderFunc
  METHOD Load : func(ctx context.Context, key string) (map[string]string, error)
TYPE DictionaryResolver : github.com/coldsmirk/vef-framework-go/mold.DictionaryResolver
  METHOD Resolve : func(ctx context.Context, key string, code string) (string, error)
TYPE FieldLevel : github.com/coldsmirk/vef-framework-go/mold.FieldLevel
  METHOD Field : func() reflect.Value
  METHOD Name : func() string
  METHOD Param : func() string
  METHOD Parent : func() reflect.Value
  METHOD SiblingField : func(name string) (reflect.Value, bool)
  METHOD Struct : func() reflect.Value
  METHOD Transformer : func() github.com/coldsmirk/vef-framework-go/mold.Transformer
TYPE FieldTransformer : github.com/coldsmirk/vef-framework-go/mold.FieldTransformer
  METHOD Tag : func() string
  METHOD Transform : func(ctx context.Context, fl github.com/coldsmirk/vef-framework-go/mold.FieldLevel) error
TYPE Func : github.com/coldsmirk/vef-framework-go/mold.Func
TYPE Interceptor : github.com/coldsmirk/vef-framework-go/mold.Interceptor
  METHOD Intercept : func(current reflect.Value) (inner reflect.Value)
TYPE InterceptorFunc : github.com/coldsmirk/vef-framework-go/mold.InterceptorFunc
FUNC NewCachedDictionaryResolver : func(loader github.com/coldsmirk/vef-framework-go/mold.DictionaryLoader, bus github.com/coldsmirk/vef-framework-go/event.Bus) github.com/coldsmirk/vef-framework-go/mold.DictionaryResolver
FUNC PublishDictionaryChangedEvent : func(ctx context.Context, bus github.com/coldsmirk/vef-framework-go/event.Bus, keys ...string) error
TYPE StructLevel : github.com/coldsmirk/vef-framework-go/mold.StructLevel
  METHOD Parent : func() reflect.Value
  METHOD Struct : func() reflect.Value
  METHOD Transformer : func() github.com/coldsmirk/vef-framework-go/mold.Transformer
TYPE StructLevelFunc : github.com/coldsmirk/vef-framework-go/mold.StructLevelFunc
TYPE StructTransformer : github.com/coldsmirk/vef-framework-go/mold.StructTransformer
  METHOD Transform : func(ctx context.Context, sl github.com/coldsmirk/vef-framework-go/mold.StructLevel) error
TYPE Transformer : github.com/coldsmirk/vef-framework-go/mold.Transformer
  METHOD Field : func(ctx context.Context, value any, tags string) error
  METHOD Struct : func(ctx context.Context, value any) error
TYPE Translator : github.com/coldsmirk/vef-framework-go/mold.Translator
  METHOD Supports : func(kind string) bool
  METHOD Translate : func(ctx context.Context, kind string, value string) (string, error)

## github.com/coldsmirk/vef-framework-go/monad
FUNC NewRange : func[T cmp.Ordered](start T, end T) github.com/coldsmirk/vef-framework-go/monad.Range[T]
TYPE Range : github.com/coldsmirk/vef-framework-go/monad.Range[T cmp.Ordered]
  FIELD Start : T [field_order=1 tag=""]
  FIELD End : T [field_order=2 tag=""]
  METHOD Contains : func(value T) bool
  METHOD Intersection : func(other github.com/coldsmirk/vef-framework-go/monad.Range[T]) (github.com/coldsmirk/vef-framework-go/monad.Range[T], bool)
  METHOD IsEmpty : func() bool
  METHOD IsNotEmpty : func() bool
  METHOD IsValid : func() bool
  METHOD Overlaps : func(other github.com/coldsmirk/vef-framework-go/monad.Range[T]) bool

## github.com/coldsmirk/vef-framework-go/monitor
TYPE BuildInfo : github.com/coldsmirk/vef-framework-go/monitor.BuildInfo
  FIELD VEFVersion : string [field_order=1 tag="json:\"vefVersion\""]
  FIELD AppVersion : string [field_order=2 tag="json:\"appVersion\""]
  FIELD BuildTime : string [field_order=3 tag="json:\"buildTime\""]
  FIELD GitCommit : string [field_order=4 tag="json:\"gitCommit\""]
TYPE CPUInfo : github.com/coldsmirk/vef-framework-go/monitor.CPUInfo
  FIELD PhysicalCores : int [field_order=1 tag="json:\"physicalCores\""]
  FIELD LogicalCores : int [field_order=2 tag="json:\"logicalCores\""]
  FIELD ModelName : string [field_order=3 tag="json:\"modelName\""]
  FIELD Mhz : float64 [field_order=4 tag="json:\"mhz\""]
  FIELD CacheSize : int32 [field_order=5 tag="json:\"cacheSize\""]
  FIELD UsagePercent : []float64 [field_order=6 tag="json:\"usagePercent\""]
  FIELD TotalPercent : float64 [field_order=7 tag="json:\"totalPercent\""]
  FIELD VendorID : string [field_order=8 tag="json:\"vendorId\""]
  FIELD Family : string [field_order=9 tag="json:\"family\""]
  FIELD Model : string [field_order=10 tag="json:\"model\""]
  FIELD Stepping : int32 [field_order=11 tag="json:\"stepping\""]
  FIELD Microcode : string [field_order=12 tag="json:\"microcode\""]
TYPE CPUSummary : github.com/coldsmirk/vef-framework-go/monitor.CPUSummary
  FIELD PhysicalCores : int [field_order=1 tag="json:\"physicalCores\""]
  FIELD LogicalCores : int [field_order=2 tag="json:\"logicalCores\""]
  FIELD UsagePercent : float64 [field_order=3 tag="json:\"usagePercent\""]
TYPE DiskInfo : github.com/coldsmirk/vef-framework-go/monitor.DiskInfo
  FIELD Partitions : []*github.com/coldsmirk/vef-framework-go/monitor.PartitionInfo [field_order=1 tag="json:\"partitions\""]
  FIELD IOCounters : map[string]*github.com/coldsmirk/vef-framework-go/monitor.IOCounter [field_order=2 tag="json:\"ioCounters\""]
TYPE DiskSummary : github.com/coldsmirk/vef-framework-go/monitor.DiskSummary
  FIELD Total : uint64 [field_order=1 tag="json:\"total\""]
  FIELD Used : uint64 [field_order=2 tag="json:\"used\""]
  FIELD UsedPercent : float64 [field_order=3 tag="json:\"usedPercent\""]
  FIELD Partitions : int [field_order=4 tag="json:\"partitions\""]
CONST ErrCodeCollectionFailed : untyped int = 2101
CONST ErrCodeNotReady : untyped int = 2100
VAR ErrCollectionFailed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrNotReady : github.com/coldsmirk/vef-framework-go/result.Error
TYPE EventStreamsInfo : github.com/coldsmirk/vef-framework-go/monitor.EventStreamsInfo
  FIELD Enabled : bool [field_order=1 tag="json:\"enabled\""]
  FIELD Streams : []github.com/coldsmirk/vef-framework-go/event.StreamInfo [field_order=2 tag="json:\"streams\""]
TYPE HostInfo : github.com/coldsmirk/vef-framework-go/monitor.HostInfo
  FIELD Hostname : string [field_order=1 tag="json:\"hostname\""]
  FIELD Uptime : uint64 [field_order=2 tag="json:\"uptime\""]
  FIELD BootTime : uint64 [field_order=3 tag="json:\"bootTime\""]
  FIELD Processes : uint64 [field_order=4 tag="json:\"processes\""]
  FIELD OS : string [field_order=5 tag="json:\"os\""]
  FIELD Platform : string [field_order=6 tag="json:\"platform\""]
  FIELD PlatformFamily : string [field_order=7 tag="json:\"platformFamily\""]
  FIELD PlatformVersion : string [field_order=8 tag="json:\"platformVersion\""]
  FIELD KernelVersion : string [field_order=9 tag="json:\"kernelVersion\""]
  FIELD KernelArch : string [field_order=10 tag="json:\"kernelArch\""]
  FIELD VirtualizationSystem : string [field_order=11 tag="json:\"virtualizationSystem\""]
  FIELD VirtualizationRole : string [field_order=12 tag="json:\"virtualizationRole\""]
  FIELD HostID : string [field_order=13 tag="json:\"hostId\""]
TYPE HostSummary : github.com/coldsmirk/vef-framework-go/monitor.HostSummary
  FIELD Hostname : string [field_order=1 tag="json:\"hostname\""]
  FIELD OS : string [field_order=2 tag="json:\"os\""]
  FIELD Platform : string [field_order=3 tag="json:\"platform\""]
  FIELD PlatformVersion : string [field_order=4 tag="json:\"platformVersion\""]
  FIELD KernelVersion : string [field_order=5 tag="json:\"kernelVersion\""]
  FIELD KernelArch : string [field_order=6 tag="json:\"kernelArch\""]
  FIELD Uptime : uint64 [field_order=7 tag="json:\"uptime\""]
TYPE IOCounter : github.com/coldsmirk/vef-framework-go/monitor.IOCounter
  FIELD ReadCount : uint64 [field_order=1 tag="json:\"readCount\""]
  FIELD MergedReadCount : uint64 [field_order=2 tag="json:\"mergedReadCount\""]
  FIELD WriteCount : uint64 [field_order=3 tag="json:\"writeCount\""]
  FIELD MergedWriteCount : uint64 [field_order=4 tag="json:\"mergedWriteCount\""]
  FIELD ReadBytes : uint64 [field_order=5 tag="json:\"readBytes\""]
  FIELD WriteBytes : uint64 [field_order=6 tag="json:\"writeBytes\""]
  FIELD ReadTime : uint64 [field_order=7 tag="json:\"readTime\""]
  FIELD WriteTime : uint64 [field_order=8 tag="json:\"writeTime\""]
  FIELD IOPSInProgress : uint64 [field_order=9 tag="json:\"iopsInProgress\""]
  FIELD IOTime : uint64 [field_order=10 tag="json:\"ioTime\""]
  FIELD WeightedIO : uint64 [field_order=11 tag="json:\"weightedIo\""]
  FIELD Name : string [field_order=12 tag="json:\"name\""]
  FIELD SerialNumber : string [field_order=13 tag="json:\"serialNumber\""]
  FIELD Label : string [field_order=14 tag="json:\"label\""]
TYPE InterfaceInfo : github.com/coldsmirk/vef-framework-go/monitor.InterfaceInfo
  FIELD Index : int [field_order=1 tag="json:\"index\""]
  FIELD MTU : int [field_order=2 tag="json:\"mtu\""]
  FIELD Name : string [field_order=3 tag="json:\"name\""]
  FIELD HardwareAddr : string [field_order=4 tag="json:\"hardwareAddr\""]
  FIELD Flags : []string [field_order=5 tag="json:\"flags\""]
  FIELD Addrs : []string [field_order=6 tag="json:\"addrs\""]
TYPE LoadInfo : github.com/coldsmirk/vef-framework-go/monitor.LoadInfo
  FIELD Load1 : float64 [field_order=1 tag="json:\"load1\""]
  FIELD Load5 : float64 [field_order=2 tag="json:\"load5\""]
  FIELD Load15 : float64 [field_order=3 tag="json:\"load15\""]
TYPE MemoryInfo : github.com/coldsmirk/vef-framework-go/monitor.MemoryInfo
  FIELD Virtual : *github.com/coldsmirk/vef-framework-go/monitor.VirtualMemory [field_order=1 tag="json:\"virtual\""]
  FIELD Swap : *github.com/coldsmirk/vef-framework-go/monitor.SwapMemory [field_order=2 tag="json:\"swap\""]
TYPE MemorySummary : github.com/coldsmirk/vef-framework-go/monitor.MemorySummary
  FIELD Total : uint64 [field_order=1 tag="json:\"total\""]
  FIELD Used : uint64 [field_order=2 tag="json:\"used\""]
  FIELD UsedPercent : float64 [field_order=3 tag="json:\"usedPercent\""]
TYPE NetIOCounter : github.com/coldsmirk/vef-framework-go/monitor.NetIOCounter
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD BytesSent : uint64 [field_order=2 tag="json:\"bytesSent\""]
  FIELD BytesRecv : uint64 [field_order=3 tag="json:\"bytesRecv\""]
  FIELD PacketsSent : uint64 [field_order=4 tag="json:\"packetsSent\""]
  FIELD PacketsRecv : uint64 [field_order=5 tag="json:\"packetsRecv\""]
  FIELD ErrorsIn : uint64 [field_order=6 tag="json:\"errorsIn\""]
  FIELD ErrorsOut : uint64 [field_order=7 tag="json:\"errorsOut\""]
  FIELD DroppedIn : uint64 [field_order=8 tag="json:\"droppedIn\""]
  FIELD DroppedOut : uint64 [field_order=9 tag="json:\"droppedOut\""]
  FIELD FIFOIn : uint64 [field_order=10 tag="json:\"fifoIn\""]
  FIELD FIFOOut : uint64 [field_order=11 tag="json:\"fifoOut\""]
TYPE NetworkInfo : github.com/coldsmirk/vef-framework-go/monitor.NetworkInfo
  FIELD Interfaces : []*github.com/coldsmirk/vef-framework-go/monitor.InterfaceInfo [field_order=1 tag="json:\"interfaces\""]
  FIELD IOCounters : map[string]*github.com/coldsmirk/vef-framework-go/monitor.NetIOCounter [field_order=2 tag="json:\"ioCounters\""]
TYPE NetworkSummary : github.com/coldsmirk/vef-framework-go/monitor.NetworkSummary
  FIELD Interfaces : int [field_order=1 tag="json:\"interfaces\""]
  FIELD BytesSent : uint64 [field_order=2 tag="json:\"bytesSent\""]
  FIELD BytesRecv : uint64 [field_order=3 tag="json:\"bytesRecv\""]
  FIELD PacketsSent : uint64 [field_order=4 tag="json:\"packetsSent\""]
  FIELD PacketsRecv : uint64 [field_order=5 tag="json:\"packetsRecv\""]
TYPE PartitionInfo : github.com/coldsmirk/vef-framework-go/monitor.PartitionInfo
  FIELD Device : string [field_order=1 tag="json:\"device\""]
  FIELD MountPoint : string [field_order=2 tag="json:\"mountPoint\""]
  FIELD FSType : string [field_order=3 tag="json:\"fsType\""]
  FIELD Options : []string [field_order=4 tag="json:\"options\""]
  FIELD Total : uint64 [field_order=5 tag="json:\"total\""]
  FIELD Free : uint64 [field_order=6 tag="json:\"free\""]
  FIELD Used : uint64 [field_order=7 tag="json:\"used\""]
  FIELD UsedPercent : float64 [field_order=8 tag="json:\"usedPercent\""]
  FIELD INodesTotal : uint64 [field_order=9 tag="json:\"iNodesTotal\""]
  FIELD INodesUsed : uint64 [field_order=10 tag="json:\"iNodesUsed\""]
  FIELD INodesFree : uint64 [field_order=11 tag="json:\"iNodesFree\""]
  FIELD INodesUsedPercent : float64 [field_order=12 tag="json:\"iNodesUsedPercent\""]
TYPE ProcessInfo : github.com/coldsmirk/vef-framework-go/monitor.ProcessInfo
  FIELD PID : int32 [field_order=1 tag="json:\"pid\""]
  FIELD ParentPID : int32 [field_order=2 tag="json:\"parentPid\""]
  FIELD Name : string [field_order=3 tag="json:\"name\""]
  FIELD Exe : string [field_order=4 tag="json:\"exe\""]
  FIELD CommandLine : string [field_order=5 tag="json:\"commandLine\""]
  FIELD CWD : string [field_order=6 tag="json:\"cwd\""]
  FIELD Status : string [field_order=7 tag="json:\"status\""]
  FIELD Username : string [field_order=8 tag="json:\"username\""]
  FIELD CreateTime : int64 [field_order=9 tag="json:\"createTime\""]
  FIELD NumThreads : int32 [field_order=10 tag="json:\"numThreads\""]
  FIELD NumFDs : int32 [field_order=11 tag="json:\"numFds\""]
  FIELD CPUPercent : float64 [field_order=12 tag="json:\"cpuPercent\""]
  FIELD MemoryPercent : float32 [field_order=13 tag="json:\"memoryPercent\""]
  FIELD MemoryRSS : uint64 [field_order=14 tag="json:\"memoryRss\""]
  FIELD MemoryVMS : uint64 [field_order=15 tag="json:\"memoryVms\""]
  FIELD MemorySwap : uint64 [field_order=16 tag="json:\"memorySwap\""]
TYPE ProcessSummary : github.com/coldsmirk/vef-framework-go/monitor.ProcessSummary
  FIELD PID : int32 [field_order=1 tag="json:\"pid\""]
  FIELD Name : string [field_order=2 tag="json:\"name\""]
  FIELD CPUPercent : float64 [field_order=3 tag="json:\"cpuPercent\""]
  FIELD MemoryPercent : float32 [field_order=4 tag="json:\"memoryPercent\""]
TYPE Service : github.com/coldsmirk/vef-framework-go/monitor.Service
  METHOD BuildInfo : func() *github.com/coldsmirk/vef-framework-go/monitor.BuildInfo
  METHOD CPU : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.CPUInfo, error)
  METHOD Disk : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.DiskInfo, error)
  METHOD Host : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.HostInfo, error)
  METHOD Load : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.LoadInfo, error)
  METHOD Memory : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.MemoryInfo, error)
  METHOD Network : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.NetworkInfo, error)
  METHOD Overview : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.SystemOverview, error)
  METHOD Process : func(ctx context.Context) (*github.com/coldsmirk/vef-framework-go/monitor.ProcessInfo, error)
TYPE SwapMemory : github.com/coldsmirk/vef-framework-go/monitor.SwapMemory
  FIELD Total : uint64 [field_order=1 tag="json:\"total\""]
  FIELD Used : uint64 [field_order=2 tag="json:\"used\""]
  FIELD Free : uint64 [field_order=3 tag="json:\"free\""]
  FIELD UsedPercent : float64 [field_order=4 tag="json:\"usedPercent\""]
  FIELD SwapIn : uint64 [field_order=5 tag="json:\"swapIn\""]
  FIELD SwapOut : uint64 [field_order=6 tag="json:\"swapOut\""]
  FIELD PageIn : uint64 [field_order=7 tag="json:\"pageIn\""]
  FIELD PageOut : uint64 [field_order=8 tag="json:\"pageOut\""]
  FIELD PageFault : uint64 [field_order=9 tag="json:\"pageFault\""]
  FIELD PageMajorFault : uint64 [field_order=10 tag="json:\"pageMajorFault\""]
TYPE SystemOverview : github.com/coldsmirk/vef-framework-go/monitor.SystemOverview
  FIELD Host : *github.com/coldsmirk/vef-framework-go/monitor.HostSummary [field_order=1 tag="json:\"host\""]
  FIELD CPU : *github.com/coldsmirk/vef-framework-go/monitor.CPUSummary [field_order=2 tag="json:\"cpu\""]
  FIELD Memory : *github.com/coldsmirk/vef-framework-go/monitor.MemorySummary [field_order=3 tag="json:\"memory\""]
  FIELD Disk : *github.com/coldsmirk/vef-framework-go/monitor.DiskSummary [field_order=4 tag="json:\"disk\""]
  FIELD Network : *github.com/coldsmirk/vef-framework-go/monitor.NetworkSummary [field_order=5 tag="json:\"network\""]
  FIELD Process : *github.com/coldsmirk/vef-framework-go/monitor.ProcessSummary [field_order=6 tag="json:\"process\""]
  FIELD Load : *github.com/coldsmirk/vef-framework-go/monitor.LoadInfo [field_order=7 tag="json:\"load\""]
  FIELD Build : *github.com/coldsmirk/vef-framework-go/monitor.BuildInfo [field_order=8 tag="json:\"build\""]
TYPE VirtualMemory : github.com/coldsmirk/vef-framework-go/monitor.VirtualMemory
  FIELD Total : uint64 [field_order=1 tag="json:\"total\""]
  FIELD Available : uint64 [field_order=2 tag="json:\"available\""]
  FIELD Used : uint64 [field_order=3 tag="json:\"used\""]
  FIELD UsedPercent : float64 [field_order=4 tag="json:\"usedPercent\""]
  FIELD Free : uint64 [field_order=5 tag="json:\"free\""]
  FIELD Active : uint64 [field_order=6 tag="json:\"active\""]
  FIELD Inactive : uint64 [field_order=7 tag="json:\"inactive\""]
  FIELD Wired : uint64 [field_order=8 tag="json:\"wired\""]
  FIELD Laundry : uint64 [field_order=9 tag="json:\"laundry\""]
  FIELD Buffers : uint64 [field_order=10 tag="json:\"buffers\""]
  FIELD Cached : uint64 [field_order=11 tag="json:\"cached\""]
  FIELD WriteBack : uint64 [field_order=12 tag="json:\"writeBack\""]
  FIELD Dirty : uint64 [field_order=13 tag="json:\"dirty\""]
  FIELD WriteBackTmp : uint64 [field_order=14 tag="json:\"writeBackTmp\""]
  FIELD Shared : uint64 [field_order=15 tag="json:\"shared\""]
  FIELD Slab : uint64 [field_order=16 tag="json:\"slab\""]
  FIELD SlabReclaimable : uint64 [field_order=17 tag="json:\"slabReclaimable\""]
  FIELD SlabUnreclaimable : uint64 [field_order=18 tag="json:\"slabUnreclaimable\""]
  FIELD PageTables : uint64 [field_order=19 tag="json:\"pageTables\""]
  FIELD SwapCached : uint64 [field_order=20 tag="json:\"swapCached\""]
  FIELD CommitLimit : uint64 [field_order=21 tag="json:\"commitLimit\""]
  FIELD CommittedAs : uint64 [field_order=22 tag="json:\"committedAs\""]
  FIELD HighTotal : uint64 [field_order=23 tag="json:\"highTotal\""]
  FIELD HighFree : uint64 [field_order=24 tag="json:\"highFree\""]
  FIELD LowTotal : uint64 [field_order=25 tag="json:\"lowTotal\""]
  FIELD LowFree : uint64 [field_order=26 tag="json:\"lowFree\""]
  FIELD SwapTotal : uint64 [field_order=27 tag="json:\"swapTotal\""]
  FIELD SwapFree : uint64 [field_order=28 tag="json:\"swapFree\""]
  FIELD Mapped : uint64 [field_order=29 tag="json:\"mapped\""]
  FIELD VMAllocTotal : uint64 [field_order=30 tag="json:\"vmAllocTotal\""]
  FIELD VMAllocUsed : uint64 [field_order=31 tag="json:\"vmAllocUsed\""]
  FIELD VMAllocChunk : uint64 [field_order=32 tag="json:\"vmAllocChunk\""]
  FIELD HugePagesTotal : uint64 [field_order=33 tag="json:\"hugePagesTotal\""]
  FIELD HugePagesFree : uint64 [field_order=34 tag="json:\"hugePagesFree\""]
  FIELD HugePagesReserved : uint64 [field_order=35 tag="json:\"hugePagesReserved\""]
  FIELD HugePagesSurplus : uint64 [field_order=36 tag="json:\"hugePagesSurplus\""]
  FIELD HugePageSize : uint64 [field_order=37 tag="json:\"hugePageSize\""]
  FIELD AnonHugePages : uint64 [field_order=38 tag="json:\"anonHugePages\""]

## github.com/coldsmirk/vef-framework-go/orm
TYPE AddColumnQuery : github.com/coldsmirk/vef-framework-go/orm.AddColumnQuery
  METHOD Column : func(name string, dataType github.com/coldsmirk/vef-framework-go/internal/orm.DataTypeDef, constraints ...github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint) github.com/coldsmirk/vef-framework-go/internal/orm.AddColumnQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD IfNotExists : func() github.com/coldsmirk/vef-framework-go/internal/orm.AddColumnQuery
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.AddColumnQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) github.com/coldsmirk/vef-framework-go/internal/orm.AddColumnQuery
TYPE AfterDeleteHook : github.com/coldsmirk/vef-framework-go/orm.AfterDeleteHook
  METHOD AfterDelete : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunDeleteQuery) error
TYPE AfterInsertHook : github.com/coldsmirk/vef-framework-go/orm.AfterInsertHook
  METHOD AfterInsert : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunInsertQuery) error
TYPE AfterScanRowHook : github.com/coldsmirk/vef-framework-go/orm.AfterScanRowHook
  METHOD AfterScanRow : func(context.Context) error
TYPE AfterSelectHook : github.com/coldsmirk/vef-framework-go/orm.AfterSelectHook
  METHOD AfterSelect : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunSelectQuery) error
TYPE AfterUpdateHook : github.com/coldsmirk/vef-framework-go/orm.AfterUpdateHook
  METHOD AfterUpdate : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunUpdateQuery) error
TYPE Applier : github.com/coldsmirk/vef-framework-go/orm.Applier[T any]
  METHOD Apply : func(fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[T]) T
  METHOD ApplyIf : func(condition bool, fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[T]) T
TYPE ApplyFunc : github.com/coldsmirk/vef-framework-go/orm.ApplyFunc[T any]
VAR ApplySort : func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery, orders []github.com/coldsmirk/vef-framework-go/sortx.OrderSpec)
TYPE ArrayAggBuilder : github.com/coldsmirk/vef-framework-go/orm.ArrayAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD IgnoreNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
  METHOD RespectNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder
VAR AutoIncrement : func() github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
TYPE AvgBuilder : github.com/coldsmirk/vef-framework-go/orm.AvgBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.AvgBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.AvgBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.AvgBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.AvgBuilder
TYPE BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel
TYPE BeforeDeleteHook : github.com/coldsmirk/vef-framework-go/orm.BeforeDeleteHook
  METHOD BeforeDelete : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunDeleteQuery) error
TYPE BeforeInsertHook : github.com/coldsmirk/vef-framework-go/orm.BeforeInsertHook
  METHOD BeforeInsert : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunInsertQuery) error
TYPE BeforeScanRowHook : github.com/coldsmirk/vef-framework-go/orm.BeforeScanRowHook
  METHOD BeforeScanRow : func(context.Context) error
TYPE BeforeSelectHook : github.com/coldsmirk/vef-framework-go/orm.BeforeSelectHook
  METHOD BeforeSelect : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunSelectQuery) error
TYPE BeforeUpdateHook : github.com/coldsmirk/vef-framework-go/orm.BeforeUpdateHook
  METHOD BeforeUpdate : func(ctx context.Context, query *github.com/coldsmirk/vef-framework-go/orm.BunUpdateQuery) error
TYPE BitAndBuilder : github.com/coldsmirk/vef-framework-go/orm.BitAndBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.BitAndBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.BitAndBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.BitAndBuilder
TYPE BitOrBuilder : github.com/coldsmirk/vef-framework-go/orm.BitOrBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.BitOrBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.BitOrBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.BitOrBuilder
TYPE BoolAndBuilder : github.com/coldsmirk/vef-framework-go/orm.BoolAndBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.BoolAndBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.BoolAndBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.BoolAndBuilder
TYPE BoolOrBuilder : github.com/coldsmirk/vef-framework-go/orm.BoolOrBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.BoolOrBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.BoolOrBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.BoolOrBuilder
TYPE BunDeleteQuery : github.com/coldsmirk/vef-framework-go/orm.BunDeleteQuery
  METHOD AppendNamedArg : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte, name string) ([]byte, bool)
  METHOD AppendQuery : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte) (_ []byte, err error)
  METHOD Apply : func(fns ...func(*github.com/uptrace/bun.DeleteQuery) *github.com/uptrace/bun.DeleteQuery) *github.com/uptrace/bun.DeleteQuery
  METHOD ApplyQueryBuilder : func(fn func(github.com/uptrace/bun.QueryBuilder) github.com/uptrace/bun.QueryBuilder) *github.com/uptrace/bun.DeleteQuery
  METHOD Comment : func(comment string) *github.com/uptrace/bun.DeleteQuery
  METHOD Conn : func(db github.com/uptrace/bun.IConn) *github.com/uptrace/bun.DeleteQuery
  METHOD DB : func() *github.com/uptrace/bun.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Err : func(err error) *github.com/uptrace/bun.DeleteQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD ForceDelete : func() *github.com/uptrace/bun.DeleteQuery
  METHOD GetModel : func() github.com/uptrace/bun.Model
  METHOD GetTableName : func() string
  METHOD Limit : func(n int) *github.com/uptrace/bun.DeleteQuery
  METHOD Model : func(model any) *github.com/uptrace/bun.DeleteQuery
  METHOD ModelTableExpr : func(query string, args ...any) *github.com/uptrace/bun.DeleteQuery
  METHOD NewAddColumn : func() *github.com/uptrace/bun.AddColumnQuery
  METHOD NewCreateIndex : func() *github.com/uptrace/bun.CreateIndexQuery
  METHOD NewCreateTable : func() *github.com/uptrace/bun.CreateTableQuery
  METHOD NewDelete : func() *github.com/uptrace/bun.DeleteQuery
  METHOD NewDropColumn : func() *github.com/uptrace/bun.DropColumnQuery
  METHOD NewDropIndex : func() *github.com/uptrace/bun.DropIndexQuery
  METHOD NewDropTable : func() *github.com/uptrace/bun.DropTableQuery
  METHOD NewInsert : func() *github.com/uptrace/bun.InsertQuery
  METHOD NewRaw : func(query string, args ...any) *github.com/uptrace/bun.RawQuery
  METHOD NewSelect : func() *github.com/uptrace/bun.SelectQuery
  METHOD NewTruncateTable : func() *github.com/uptrace/bun.TruncateTableQuery
  METHOD NewUpdate : func() *github.com/uptrace/bun.UpdateQuery
  METHOD NewValues : func(model any) *github.com/uptrace/bun.ValuesQuery
  METHOD Operation : func() string
  METHOD Order : func(orders ...string) *github.com/uptrace/bun.DeleteQuery
  METHOD OrderExpr : func(query string, args ...any) *github.com/uptrace/bun.DeleteQuery
  METHOD QueryBuilder : func() github.com/uptrace/bun.QueryBuilder
  METHOD Returning : func(query string, args ...any) *github.com/uptrace/bun.DeleteQuery
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD String : func() string
  METHOD Table : func(tables ...string) *github.com/uptrace/bun.DeleteQuery
  METHOD TableExpr : func(query string, args ...any) *github.com/uptrace/bun.DeleteQuery
  METHOD Where : func(query string, args ...any) *github.com/uptrace/bun.DeleteQuery
  METHOD WhereAllWithDeleted : func() *github.com/uptrace/bun.DeleteQuery
  METHOD WhereDeleted : func() *github.com/uptrace/bun.DeleteQuery
  METHOD WhereGroup : func(sep string, fn func(*github.com/uptrace/bun.DeleteQuery) *github.com/uptrace/bun.DeleteQuery) *github.com/uptrace/bun.DeleteQuery
  METHOD WhereOr : func(query string, args ...any) *github.com/uptrace/bun.DeleteQuery
  METHOD WherePK : func(cols ...string) *github.com/uptrace/bun.DeleteQuery
  METHOD With : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.DeleteQuery
  METHOD WithQuery : func(query *github.com/uptrace/bun.WithQuery) *github.com/uptrace/bun.DeleteQuery
  METHOD WithRecursive : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.DeleteQuery
TYPE BunInsertQuery : github.com/coldsmirk/vef-framework-go/orm.BunInsertQuery
  METHOD AppendNamedArg : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte, name string) ([]byte, bool)
  METHOD AppendQuery : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte) (_ []byte, err error)
  METHOD Apply : func(fns ...func(*github.com/uptrace/bun.InsertQuery) *github.com/uptrace/bun.InsertQuery) *github.com/uptrace/bun.InsertQuery
  METHOD Column : func(columns ...string) *github.com/uptrace/bun.InsertQuery
  METHOD ColumnExpr : func(query string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD Comment : func(comment string) *github.com/uptrace/bun.InsertQuery
  METHOD Conn : func(db github.com/uptrace/bun.IConn) *github.com/uptrace/bun.InsertQuery
  METHOD DB : func() *github.com/uptrace/bun.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Err : func(err error) *github.com/uptrace/bun.InsertQuery
  METHOD ExcludeColumn : func(columns ...string) *github.com/uptrace/bun.InsertQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD GetModel : func() github.com/uptrace/bun.Model
  METHOD GetTableName : func() string
  METHOD Ignore : func() *github.com/uptrace/bun.InsertQuery
  METHOD Model : func(model any) *github.com/uptrace/bun.InsertQuery
  METHOD ModelTableExpr : func(query string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD NewAddColumn : func() *github.com/uptrace/bun.AddColumnQuery
  METHOD NewCreateIndex : func() *github.com/uptrace/bun.CreateIndexQuery
  METHOD NewCreateTable : func() *github.com/uptrace/bun.CreateTableQuery
  METHOD NewDelete : func() *github.com/uptrace/bun.DeleteQuery
  METHOD NewDropColumn : func() *github.com/uptrace/bun.DropColumnQuery
  METHOD NewDropIndex : func() *github.com/uptrace/bun.DropIndexQuery
  METHOD NewDropTable : func() *github.com/uptrace/bun.DropTableQuery
  METHOD NewInsert : func() *github.com/uptrace/bun.InsertQuery
  METHOD NewRaw : func(query string, args ...any) *github.com/uptrace/bun.RawQuery
  METHOD NewSelect : func() *github.com/uptrace/bun.SelectQuery
  METHOD NewTruncateTable : func() *github.com/uptrace/bun.TruncateTableQuery
  METHOD NewUpdate : func() *github.com/uptrace/bun.UpdateQuery
  METHOD NewValues : func(model any) *github.com/uptrace/bun.ValuesQuery
  METHOD On : func(s string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD Operation : func() string
  METHOD Replace : func() *github.com/uptrace/bun.InsertQuery
  METHOD Returning : func(query string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD Set : func(query string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD SetValues : func(values *github.com/uptrace/bun.ValuesQuery) *github.com/uptrace/bun.InsertQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) *github.com/uptrace/bun.InsertQuery
  METHOD TableExpr : func(query string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD Value : func(column string, expr string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD Where : func(query string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD WhereOr : func(query string, args ...any) *github.com/uptrace/bun.InsertQuery
  METHOD With : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.InsertQuery
  METHOD WithQuery : func(query *github.com/uptrace/bun.WithQuery) *github.com/uptrace/bun.InsertQuery
  METHOD WithRecursive : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.InsertQuery
TYPE BunSelectQuery : github.com/coldsmirk/vef-framework-go/orm.BunSelectQuery
  METHOD AppendNamedArg : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte, name string) ([]byte, bool)
  METHOD AppendQuery : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte) (_ []byte, err error)
  METHOD Apply : func(fns ...func(*github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD ApplyQueryBuilder : func(fn func(github.com/uptrace/bun.QueryBuilder) github.com/uptrace/bun.QueryBuilder) *github.com/uptrace/bun.SelectQuery
  METHOD Clone : func() *github.com/uptrace/bun.SelectQuery
  METHOD Column : func(columns ...string) *github.com/uptrace/bun.SelectQuery
  METHOD ColumnExpr : func(query string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD Comment : func(comment string) *github.com/uptrace/bun.SelectQuery
  METHOD Conn : func(db github.com/uptrace/bun.IConn) *github.com/uptrace/bun.SelectQuery
  METHOD Count : func(ctx context.Context) (int, error)
  METHOD DB : func() *github.com/uptrace/bun.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Distinct : func() *github.com/uptrace/bun.SelectQuery
  METHOD DistinctOn : func(query string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD Err : func(err error) *github.com/uptrace/bun.SelectQuery
  METHOD Except : func(other *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD ExceptAll : func(other *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD ExcludeColumn : func(columns ...string) *github.com/uptrace/bun.SelectQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (res database/sql.Result, err error)
  METHOD Exists : func(ctx context.Context) (bool, error)
  METHOD For : func(s string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD ForceIndex : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD ForceIndexForGroupBy : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD ForceIndexForJoin : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD ForceIndexForOrderBy : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD GetModel : func() github.com/uptrace/bun.Model
  METHOD GetTableName : func() string
  METHOD Group : func(columns ...string) *github.com/uptrace/bun.SelectQuery
  METHOD GroupExpr : func(group string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD Having : func(having string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD IgnoreIndex : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD IgnoreIndexForGroupBy : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD IgnoreIndexForJoin : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD IgnoreIndexForOrderBy : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD Intersect : func(other *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD IntersectAll : func(other *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD Join : func(join string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD JoinOn : func(cond string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD JoinOnOr : func(cond string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD Limit : func(n int) *github.com/uptrace/bun.SelectQuery
  METHOD Model : func(model any) *github.com/uptrace/bun.SelectQuery
  METHOD ModelTableExpr : func(query string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD NewAddColumn : func() *github.com/uptrace/bun.AddColumnQuery
  METHOD NewCreateIndex : func() *github.com/uptrace/bun.CreateIndexQuery
  METHOD NewCreateTable : func() *github.com/uptrace/bun.CreateTableQuery
  METHOD NewDelete : func() *github.com/uptrace/bun.DeleteQuery
  METHOD NewDropColumn : func() *github.com/uptrace/bun.DropColumnQuery
  METHOD NewDropIndex : func() *github.com/uptrace/bun.DropIndexQuery
  METHOD NewDropTable : func() *github.com/uptrace/bun.DropTableQuery
  METHOD NewInsert : func() *github.com/uptrace/bun.InsertQuery
  METHOD NewRaw : func(query string, args ...any) *github.com/uptrace/bun.RawQuery
  METHOD NewSelect : func() *github.com/uptrace/bun.SelectQuery
  METHOD NewTruncateTable : func() *github.com/uptrace/bun.TruncateTableQuery
  METHOD NewUpdate : func() *github.com/uptrace/bun.UpdateQuery
  METHOD NewValues : func(model any) *github.com/uptrace/bun.ValuesQuery
  METHOD Offset : func(n int) *github.com/uptrace/bun.SelectQuery
  METHOD Operation : func() string
  METHOD Order : func(orders ...string) *github.com/uptrace/bun.SelectQuery
  METHOD OrderBy : func(colName string, sortDir github.com/uptrace/bun.Order) *github.com/uptrace/bun.SelectQuery
  METHOD OrderExpr : func(query string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD QueryBuilder : func() github.com/uptrace/bun.QueryBuilder
  METHOD Relation : func(name string, apply ...func(*github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD RelationWithOpts : func(name string, opts github.com/uptrace/bun.RelationOpts) *github.com/uptrace/bun.SelectQuery
  METHOD Rows : func(ctx context.Context) (*database/sql.Rows, error)
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD ScanAndCount : func(ctx context.Context, dest ...any) (int, error)
  METHOD String : func() string
  METHOD Table : func(tables ...string) *github.com/uptrace/bun.SelectQuery
  METHOD TableExpr : func(query string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD Union : func(other *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD UnionAll : func(other *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD UseIndex : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD UseIndexForGroupBy : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD UseIndexForJoin : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD UseIndexForOrderBy : func(indexes ...string) *github.com/uptrace/bun.SelectQuery
  METHOD Where : func(query string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD WhereAllWithDeleted : func() *github.com/uptrace/bun.SelectQuery
  METHOD WhereDeleted : func() *github.com/uptrace/bun.SelectQuery
  METHOD WhereGroup : func(sep string, fn func(*github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery) *github.com/uptrace/bun.SelectQuery
  METHOD WhereOr : func(query string, args ...any) *github.com/uptrace/bun.SelectQuery
  METHOD WherePK : func(cols ...string) *github.com/uptrace/bun.SelectQuery
  METHOD With : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.SelectQuery
  METHOD WithQuery : func(query *github.com/uptrace/bun.WithQuery) *github.com/uptrace/bun.SelectQuery
  METHOD WithRecursive : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.SelectQuery
TYPE BunUpdateQuery : github.com/coldsmirk/vef-framework-go/orm.BunUpdateQuery
  METHOD AppendNamedArg : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte, name string) ([]byte, bool)
  METHOD AppendQuery : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte) (_ []byte, err error)
  METHOD Apply : func(fns ...func(*github.com/uptrace/bun.UpdateQuery) *github.com/uptrace/bun.UpdateQuery) *github.com/uptrace/bun.UpdateQuery
  METHOD ApplyQueryBuilder : func(fn func(github.com/uptrace/bun.QueryBuilder) github.com/uptrace/bun.QueryBuilder) *github.com/uptrace/bun.UpdateQuery
  METHOD Bulk : func() *github.com/uptrace/bun.UpdateQuery
  METHOD Column : func(columns ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD Comment : func(comment string) *github.com/uptrace/bun.UpdateQuery
  METHOD Conn : func(db github.com/uptrace/bun.IConn) *github.com/uptrace/bun.UpdateQuery
  METHOD DB : func() *github.com/uptrace/bun.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Err : func(err error) *github.com/uptrace/bun.UpdateQuery
  METHOD ExcludeColumn : func(columns ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD FQN : func(column string) github.com/uptrace/bun.Ident
  METHOD ForceIndex : func(indexes ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD GetModel : func() github.com/uptrace/bun.Model
  METHOD GetTableName : func() string
  METHOD IgnoreIndex : func(indexes ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD Join : func(join string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD JoinOn : func(cond string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD JoinOnOr : func(cond string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD Limit : func(n int) *github.com/uptrace/bun.UpdateQuery
  METHOD Model : func(model any) *github.com/uptrace/bun.UpdateQuery
  METHOD ModelTableExpr : func(query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD NewAddColumn : func() *github.com/uptrace/bun.AddColumnQuery
  METHOD NewCreateIndex : func() *github.com/uptrace/bun.CreateIndexQuery
  METHOD NewCreateTable : func() *github.com/uptrace/bun.CreateTableQuery
  METHOD NewDelete : func() *github.com/uptrace/bun.DeleteQuery
  METHOD NewDropColumn : func() *github.com/uptrace/bun.DropColumnQuery
  METHOD NewDropIndex : func() *github.com/uptrace/bun.DropIndexQuery
  METHOD NewDropTable : func() *github.com/uptrace/bun.DropTableQuery
  METHOD NewInsert : func() *github.com/uptrace/bun.InsertQuery
  METHOD NewRaw : func(query string, args ...any) *github.com/uptrace/bun.RawQuery
  METHOD NewSelect : func() *github.com/uptrace/bun.SelectQuery
  METHOD NewTruncateTable : func() *github.com/uptrace/bun.TruncateTableQuery
  METHOD NewUpdate : func() *github.com/uptrace/bun.UpdateQuery
  METHOD NewValues : func(model any) *github.com/uptrace/bun.ValuesQuery
  METHOD OmitZero : func() *github.com/uptrace/bun.UpdateQuery
  METHOD Operation : func() string
  METHOD Order : func(orders ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD OrderExpr : func(query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD QueryBuilder : func() github.com/uptrace/bun.QueryBuilder
  METHOD Returning : func(query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD Set : func(query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD SetColumn : func(column string, query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD TableExpr : func(query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD UseIndex : func(indexes ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD Value : func(column string, query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD Where : func(query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD WhereAllWithDeleted : func() *github.com/uptrace/bun.UpdateQuery
  METHOD WhereDeleted : func() *github.com/uptrace/bun.UpdateQuery
  METHOD WhereGroup : func(sep string, fn func(*github.com/uptrace/bun.UpdateQuery) *github.com/uptrace/bun.UpdateQuery) *github.com/uptrace/bun.UpdateQuery
  METHOD WhereOr : func(query string, args ...any) *github.com/uptrace/bun.UpdateQuery
  METHOD WherePK : func(cols ...string) *github.com/uptrace/bun.UpdateQuery
  METHOD With : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.UpdateQuery
  METHOD WithQuery : func(query *github.com/uptrace/bun.WithQuery) *github.com/uptrace/bun.UpdateQuery
  METHOD WithRecursive : func(name string, query github.com/uptrace/bun.Query) *github.com/uptrace/bun.UpdateQuery
TYPE CaseBuilder : github.com/coldsmirk/vef-framework-go/orm.CaseBuilder
  METHOD Case : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.CaseBuilder
  METHOD CaseColumn : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.CaseBuilder
  METHOD CaseSubQuery : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.CaseBuilder
  METHOD Else : func(expr any)
  METHOD ElseSubQuery : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery))
  METHOD When : func(func(cb github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CaseWhenBuilder
  METHOD WhenExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.CaseWhenBuilder
  METHOD WhenSubQuery : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.CaseWhenBuilder
TYPE CaseWhenBuilder : github.com/coldsmirk/vef-framework-go/orm.CaseWhenBuilder
  METHOD Then : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.CaseBuilder
  METHOD ThenSubQuery : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.CaseBuilder
VAR Check : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
TYPE CheckBuilder : github.com/coldsmirk/vef-framework-go/orm.CheckBuilder
  METHOD Condition : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CheckBuilder
  METHOD Name : func(name string) github.com/coldsmirk/vef-framework-go/internal/orm.CheckBuilder
TYPE ColumnConstraint : github.com/coldsmirk/vef-framework-go/orm.ColumnConstraint
CONST ColumnCreatedAt : untyped string = "created_at"
CONST ColumnCreatedBy : untyped string = "created_by"
CONST ColumnCreatedByName : untyped string = "created_by_name"
CONST ColumnID : untyped string = "id"
TYPE ColumnInfo : github.com/coldsmirk/vef-framework-go/orm.ColumnInfo
CONST ColumnUpdatedAt : untyped string = "updated_at"
CONST ColumnUpdatedBy : untyped string = "updated_by"
CONST ColumnUpdatedByName : untyped string = "updated_by_name"
TYPE ConditionBuilder : github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder
  METHOD Apply : func(fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder]) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD ApplyIf : func(condition bool, fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder]) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD Between : func(column string, start any, end any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD BetweenExpr : func(column string, startB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, endB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD Contains : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD ContainsAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD ContainsAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD ContainsIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedAtBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedAtGreaterThan : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedAtGreaterThanOrEqual : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedAtLessThan : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedAtLessThanOrEqual : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedAtNotBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByEquals : func(createdBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByIn : func(createdBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByNotEquals : func(createdBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByNotEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByNotEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByNotEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByNotEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByNotIn : func(createdBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD CreatedByNotInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EndsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EndsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EndsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EndsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD Equals : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EqualsAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EqualsAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EqualsColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EqualsExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD EqualsSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD Expr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThan : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanOrEqual : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanOrEqualAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanOrEqualAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanOrEqualColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanOrEqualExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanOrEqualSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD GreaterThanSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD Group : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD In : func(column string, values any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD InExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD InSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsFalse : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsFalseExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsFalseSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsNotNull : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsNotNullExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsNotNullSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsNull : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsNullExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsNullSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsTrue : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsTrueExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD IsTrueSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThan : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanOrEqual : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanOrEqualAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanOrEqualAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanOrEqualColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanOrEqualExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanOrEqualSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD LessThanSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotBetween : func(column string, start any, end any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotBetweenExpr : func(column string, startB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, endB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotContains : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotContainsAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotContainsAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotContainsIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEndsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEndsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEndsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEndsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEquals : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEqualsAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEqualsAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEqualsColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEqualsExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotEqualsSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotIn : func(column string, values any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotInExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotInSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotStartsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotStartsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotStartsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD NotStartsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrBetween : func(column string, start any, end any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrBetweenExpr : func(column string, startB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, endB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrContains : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrContainsAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrContainsAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrContainsIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedAtBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedAtGreaterThan : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedAtGreaterThanOrEqual : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedAtLessThan : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedAtLessThanOrEqual : func(createdAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedAtNotBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByEquals : func(createdBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByIn : func(createdBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByNotEquals : func(createdBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByNotEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByNotEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByNotEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByNotEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByNotIn : func(createdBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrCreatedByNotInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEndsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEndsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEndsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEndsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEquals : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEqualsAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEqualsAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEqualsColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEqualsExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrEqualsSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThan : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanOrEqual : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanOrEqualAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanOrEqualAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanOrEqualColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanOrEqualExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanOrEqualSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGreaterThanSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrGroup : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIn : func(column string, values any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrInExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrInSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsFalse : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsFalseExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsFalseSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsNotNull : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsNotNullExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsNotNullSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsNull : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsNullExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsNullSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsTrue : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsTrueExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrIsTrueSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThan : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanOrEqual : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanOrEqualAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanOrEqualAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanOrEqualColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanOrEqualExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanOrEqualSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrLessThanSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotBetween : func(column string, start any, end any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotBetweenExpr : func(column string, startB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, endB func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotContains : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotContainsAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotContainsAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotContainsIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEndsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEndsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEndsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEndsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEquals : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEqualsAll : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEqualsAny : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEqualsColumn : func(column1 string, column2 string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEqualsExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotEqualsSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotIn : func(column string, values any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotInExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotInSubQuery : func(column string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotStartsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotStartsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotStartsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrNotStartsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrPKEquals : func(pk any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrPKIn : func(pks any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrPKNotEquals : func(pk any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrPKNotIn : func(pks any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrStartsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrStartsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrStartsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrStartsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedAtBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedAtGreaterThan : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedAtGreaterThanOrEqual : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedAtLessThan : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedAtLessThanOrEqual : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedAtNotBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByEquals : func(updatedBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByIn : func(updatedBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByNotEquals : func(updatedBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByNotEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByNotEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByNotEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByNotEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByNotIn : func(updatedBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD OrUpdatedByNotInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD PKEquals : func(pk any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD PKIn : func(pks any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD PKNotEquals : func(pk any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD PKNotIn : func(pks any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD StartsWith : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD StartsWithAny : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD StartsWithAnyIgnoreCase : func(column string, values []string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD StartsWithIgnoreCase : func(column string, value string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedAtBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedAtGreaterThan : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedAtGreaterThanOrEqual : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedAtLessThan : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedAtLessThanOrEqual : func(updatedAt time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedAtNotBetween : func(start time.Time, end time.Time, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByEquals : func(updatedBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByIn : func(updatedBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByNotEquals : func(updatedBy string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByNotEqualsAll : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByNotEqualsAny : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByNotEqualsCurrent : func(alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByNotEqualsSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByNotIn : func(updatedBys []string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
  METHOD UpdatedByNotInSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder
TYPE ConflictAction : github.com/coldsmirk/vef-framework-go/orm.ConflictAction
  METHOD String : func() string
TYPE ConflictBuilder : github.com/coldsmirk/vef-framework-go/orm.ConflictBuilder
  METHOD Columns : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ConflictBuilder
  METHOD Constraint : func(name string) github.com/coldsmirk/vef-framework-go/internal/orm.ConflictBuilder
  METHOD DoNothing : func()
  METHOD DoUpdate : func() github.com/coldsmirk/vef-framework-go/internal/orm.ConflictUpdateBuilder
  METHOD Where : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.ConflictBuilder
CONST ConflictDoNothing : github.com/coldsmirk/vef-framework-go/internal/orm.ConflictAction = 0
CONST ConflictDoUpdate : github.com/coldsmirk/vef-framework-go/internal/orm.ConflictAction = 1
TYPE ConflictUpdateBuilder : github.com/coldsmirk/vef-framework-go/orm.ConflictUpdateBuilder
  METHOD Set : func(column string, value ...any) github.com/coldsmirk/vef-framework-go/internal/orm.ConflictUpdateBuilder
  METHOD SetExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.ConflictUpdateBuilder
  METHOD Where : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.ConflictUpdateBuilder
TYPE CountBuilder : github.com/coldsmirk/vef-framework-go/orm.CountBuilder
  METHOD All : func() github.com/coldsmirk/vef-framework-go/internal/orm.CountBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.CountBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.CountBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.CountBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CountBuilder
TYPE CreateIndexQuery : github.com/coldsmirk/vef-framework-go/orm.CreateIndexQuery
  METHOD Column : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD ColumnExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Concurrently : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD ExcludeColumn : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD IfNotExists : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Include : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Index : func(name string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Unique : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Using : func(method github.com/coldsmirk/vef-framework-go/internal/orm.IndexMethod) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD Where : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
TYPE CreateTableQuery : github.com/coldsmirk/vef-framework-go/orm.CreateTableQuery
  METHOD Check : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.CheckBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD Column : func(name string, dataType github.com/coldsmirk/vef-framework-go/internal/orm.DataTypeDef, constraints ...github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD DefaultVarChar : func(n int) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD ForeignKey : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ForeignKeyBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD IfNotExists : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD PartitionBy : func(strategy github.com/coldsmirk/vef-framework-go/internal/orm.PartitionStrategy, columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD PrimaryKey : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.PrimaryKeyBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD TableSpace : func(tablespace string) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD Temp : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD Unique : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.UniqueBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD WithForeignKeys : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
TYPE CreationAuditedModel : github.com/coldsmirk/vef-framework-go/orm.CreationAuditedModel
TYPE CreationTrackedModel : github.com/coldsmirk/vef-framework-go/orm.CreationTrackedModel
TYPE CumeDistBuilder : github.com/coldsmirk/vef-framework-go/orm.CumeDistBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE DB : github.com/coldsmirk/vef-framework-go/orm.DB
  METHOD BeginTx : func(ctx context.Context, opts *database/sql.TxOptions) (github.com/coldsmirk/vef-framework-go/internal/orm.Tx, error)
  METHOD Connection : func(ctx context.Context) (*database/sql.Conn, error)
  METHOD ModelPKFields : func(model any) []*github.com/coldsmirk/vef-framework-go/internal/orm.PKField
  METHOD ModelPKs : func(model any) (map[string]any, error)
  METHOD NewAddColumn : func() github.com/coldsmirk/vef-framework-go/internal/orm.AddColumnQuery
  METHOD NewCreateIndex : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD NewCreateTable : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD NewDelete : func() github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD NewDropColumn : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropColumnQuery
  METHOD NewDropIndex : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropIndexQuery
  METHOD NewDropTable : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropTableQuery
  METHOD NewInsert : func() github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD NewMerge : func() github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD NewRaw : func(query string, args ...any) github.com/coldsmirk/vef-framework-go/internal/orm.RawQuery
  METHOD NewSelect : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD NewTruncateTable : func() github.com/coldsmirk/vef-framework-go/internal/orm.TruncateTableQuery
  METHOD NewUpdate : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD RegisterModel : func(models ...any)
  METHOD ResetModel : func(ctx context.Context, models ...any) error
  METHOD RunInReadOnlyTx : func(ctx context.Context, fn func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/internal/orm.DB) error) error
  METHOD RunInTx : func(ctx context.Context, fn func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/internal/orm.DB) error) error
  METHOD ScanRow : func(ctx context.Context, rows *database/sql.Rows, dest ...any) error
  METHOD ScanRows : func(ctx context.Context, rows *database/sql.Rows, dest ...any) error
  METHOD TableOf : func(model any) *github.com/uptrace/bun/schema.Table
  METHOD WithNamedArg : func(name string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.DB
VAR DataType : github.com/coldsmirk/vef-framework-go/internal/orm.DataTypeFactory
TYPE DataTypeDef : github.com/coldsmirk/vef-framework-go/orm.DataTypeDef
TYPE DateTimeUnit : github.com/coldsmirk/vef-framework-go/orm.DateTimeUnit
  METHOD ForDateTrunc : func() string
  METHOD ForMySQL : func() string
  METHOD ForPostgres : func() string
  METHOD ForSQLite : func() string
  METHOD String : func() string
VAR Default : func(value any) github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
TYPE DeleteQuery : github.com/coldsmirk/vef-framework-go/orm.DeleteQuery
  METHOD Apply : func(fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD ApplyIf : func(condition bool, fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD BuildCondition : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) interface{github.com/uptrace/bun/schema.QueryAppender; github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder}
  METHOD BuildSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) *github.com/uptrace/bun.SelectQuery
  METHOD CreateSubQuery : func(subQuery *github.com/uptrace/bun.SelectQuery) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DB : func() github.com/coldsmirk/vef-framework-go/internal/orm.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD ExprBuilder : func() github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder
  METHOD ForceDelete : func() github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD GetTable : func() *github.com/uptrace/bun/schema.Table
  METHOD IncludeDeleted : func() github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD Limit : func(limit int) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD ModelTable : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD OrderByExpr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD Query : func() github.com/uptrace/bun.Query
  METHOD Returning : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD ReturningAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD ReturningNone : func() github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD String : func() string
  METHOD Table : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD TableExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD TableFrom : func(model any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD TableSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD Where : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD WhereDeleted : func() github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD WherePK : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD With : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD WithOrderedValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD WithRecursive : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD WithValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
TYPE DenseRankBuilder : github.com/coldsmirk/vef-framework-go/orm.DenseRankBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE Dialect : github.com/coldsmirk/vef-framework-go/orm.Dialect
  METHOD AppendBool : func(b []byte, v bool) []byte
  METHOD AppendBytes : func(b []byte, bs []byte) []byte
  METHOD AppendJSON : func(b []byte, jsonb []byte) []byte
  METHOD AppendSequence : func(b []byte, t *github.com/uptrace/bun/schema.Table, f *github.com/uptrace/bun/schema.Field) []byte
  METHOD AppendString : func(b []byte, s string) []byte
  METHOD AppendTime : func(b []byte, tm time.Time) []byte
  METHOD AppendUint32 : func(b []byte, n uint32) []byte
  METHOD AppendUint64 : func(b []byte, n uint64) []byte
  METHOD DefaultSchema : func() string
  METHOD DefaultVarcharLen : func() int
  METHOD Features : func() github.com/uptrace/bun/dialect/feature.Feature
  METHOD IdentQuote : func() byte
  METHOD Init : func(db *database/sql.DB)
  METHOD Name : func() github.com/uptrace/bun/dialect.Name
  METHOD OnTable : func(table *github.com/uptrace/bun/schema.Table)
  METHOD Tables : func() *github.com/uptrace/bun/schema.Tables
TYPE DropColumnQuery : github.com/coldsmirk/vef-framework-go/orm.DropColumnQuery
  METHOD Column : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DropColumnQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.DropColumnQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DropColumnQuery
TYPE DropIndexQuery : github.com/coldsmirk/vef-framework-go/orm.DropIndexQuery
  METHOD Cascade : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropIndexQuery
  METHOD Concurrently : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropIndexQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD IfExists : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropIndexQuery
  METHOD Index : func(name string) github.com/coldsmirk/vef-framework-go/internal/orm.DropIndexQuery
  METHOD Restrict : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropIndexQuery
  METHOD String : func() string
TYPE DropTableQuery : github.com/coldsmirk/vef-framework-go/orm.DropTableQuery
  METHOD Cascade : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropTableQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD IfExists : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropTableQuery
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.DropTableQuery
  METHOD Restrict : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropTableQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) github.com/coldsmirk/vef-framework-go/internal/orm.DropTableQuery
TYPE Executor : github.com/coldsmirk/vef-framework-go/orm.Executor
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
TYPE ExprBuilder : github.com/coldsmirk/vef-framework-go/orm.ExprBuilder
  METHOD Abs : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Acos : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Add : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Age : func(start any, end any) github.com/uptrace/bun/schema.QueryAppender
  METHOD All : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/uptrace/bun/schema.QueryAppender
  METHOD AllColumns : func(tableAlias ...string) github.com/uptrace/bun/schema.QueryAppender
  METHOD Any : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/uptrace/bun/schema.QueryAppender
  METHOD ArrayAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ArrayAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Asin : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Atan : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Avg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.AvgBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD AvgColumn : func(column string, distinct ...bool) github.com/uptrace/bun/schema.QueryAppender
  METHOD Between : func(expr any, lower any, upper any) github.com/uptrace/bun/schema.QueryAppender
  METHOD BitAnd : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.BitAndBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD BitOr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.BitOrBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD BoolAnd : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.BoolAndBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD BoolOr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.BoolOrBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Case : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.CaseBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Ceil : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD CharLength : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Coalesce : func(args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Column : func(column string, withTableAlias ...bool) github.com/uptrace/bun/schema.QueryAppender
  METHOD Concat : func(args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ConcatWithSep : func(separator any, args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Contains : func(expr any, substr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ContainsIgnoreCase : func(expr any, substr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Cos : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Count : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.CountBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD CountAll : func(distinct ...bool) github.com/uptrace/bun/schema.QueryAppender
  METHOD CountColumn : func(column string, distinct ...bool) github.com/uptrace/bun/schema.QueryAppender
  METHOD CumeDist : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.CumeDistBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD CurrentDate : func() github.com/uptrace/bun/schema.QueryAppender
  METHOD CurrentTime : func() github.com/uptrace/bun/schema.QueryAppender
  METHOD CurrentTimestamp : func() github.com/uptrace/bun/schema.QueryAppender
  METHOD DateAdd : func(expr any, interval any, unit github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit) github.com/uptrace/bun/schema.QueryAppender
  METHOD DateDiff : func(start any, end any, unit github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit) github.com/uptrace/bun/schema.QueryAppender
  METHOD DateSubtract : func(expr any, interval any, unit github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit) github.com/uptrace/bun/schema.QueryAppender
  METHOD DateTrunc : func(unit github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit, expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Decode : func(args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD DenseRank : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.DenseRankBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Divide : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD EndsWith : func(expr any, suffix any) github.com/uptrace/bun/schema.QueryAppender
  METHOD EndsWithIgnoreCase : func(expr any, suffix any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Equals : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExecByDialect : func(execs github.com/coldsmirk/vef-framework-go/internal/orm.DialectExecs)
  METHOD ExecByDialectWithErr : func(execs github.com/coldsmirk/vef-framework-go/internal/orm.DialectExecsWithErr) error
  METHOD Exists : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Exp : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Expr : func(expr string, args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExprByDialect : func(exprs github.com/coldsmirk/vef-framework-go/internal/orm.DialectExprs) github.com/uptrace/bun/schema.QueryAppender
  METHOD Exprs : func(exprs ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExprsWithSep : func(separator any, exprs ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExtractDay : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExtractHour : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExtractMinute : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExtractMonth : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExtractSecond : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ExtractYear : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD FirstValue : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.FirstValueBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Floor : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD FragmentByDialect : func(fragments github.com/coldsmirk/vef-framework-go/internal/orm.DialectFragments) ([]byte, error)
  METHOD GreaterThan : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD GreaterThanOrEqual : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Greatest : func(args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD IfNull : func(expr any, defaultValue any) github.com/uptrace/bun/schema.QueryAppender
  METHOD In : func(expr any, values ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD IsFalse : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD IsNotNull : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD IsNull : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD IsTrue : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONArray : func(args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONArrayAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONArrayAppend : func(json any, path any, value any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONContains : func(json any, value any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONContainsPath : func(json any, path any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONExtract : func(json any, path any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONInsert : func(json any, path any, value any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONKeys : func(json any, path ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONLength : func(json any, path ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONObject : func(keyValues ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONObjectAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONReplace : func(json any, path any, value any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONSet : func(json any, path any, value any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONType : func(json any, path ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONUnquote : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD JSONValid : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Lag : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.LagBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD LastValue : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.LastValueBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Lead : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.LeadBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Least : func(args ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Left : func(expr any, length any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Length : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD LessThan : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD LessThanOrEqual : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Literal : func(value any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Ln : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Log : func(expr any, base ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Lower : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Max : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.MaxBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD MaxColumn : func(column string) github.com/uptrace/bun/schema.QueryAppender
  METHOD Min : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.MinBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD MinColumn : func(column string) github.com/uptrace/bun/schema.QueryAppender
  METHOD Mod : func(dividend any, divisor any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Multiply : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD NTile : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.NTileBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Not : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD NotBetween : func(expr any, lower any, upper any) github.com/uptrace/bun/schema.QueryAppender
  METHOD NotEquals : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD NotExists : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/uptrace/bun/schema.QueryAppender
  METHOD NotIn : func(expr any, values ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Now : func() github.com/uptrace/bun/schema.QueryAppender
  METHOD NthValue : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Null : func() github.com/uptrace/bun/schema.QueryAppender
  METHOD NullIf : func(expr1 any, expr2 any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Order : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.OrderBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Paren : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD PercentRank : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.PercentRankBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Pi : func() github.com/uptrace/bun/schema.QueryAppender
  METHOD Position : func(substring any, str any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Power : func(base any, exponent any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Random : func() github.com/uptrace/bun/schema.QueryAppender
  METHOD Rank : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.RankBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Repeat : func(expr any, count any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Replace : func(expr any, search any, replacement any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Reverse : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Right : func(expr any, length any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Round : func(expr any, precision ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD RowNumber : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.RowNumberBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD Sign : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Sin : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Sqrt : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD StartsWith : func(expr any, prefix any) github.com/uptrace/bun/schema.QueryAppender
  METHOD StartsWithIgnoreCase : func(expr any, prefix any) github.com/uptrace/bun/schema.QueryAppender
  METHOD StdDev : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.StdDevBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD StringAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD SubQuery : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/uptrace/bun/schema.QueryAppender
  METHOD SubString : func(expr any, start any, length ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Subtract : func(left any, right any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Sum : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.SumBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD SumColumn : func(column string, distinct ...bool) github.com/uptrace/bun/schema.QueryAppender
  METHOD TableColumns : func(withTableAlias ...bool) github.com/uptrace/bun/schema.QueryAppender
  METHOD Tan : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToBool : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToDate : func(expr any, format ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToDecimal : func(expr any, precision ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToFloat : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToInteger : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToJSON : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToString : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToTime : func(expr any, format ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD ToTimestamp : func(expr any, format ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Trim : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD TrimLeft : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD TrimRight : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Trunc : func(expr any, precision ...any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Upper : func(expr any) github.com/uptrace/bun/schema.QueryAppender
  METHOD Variance : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.VarianceBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinArrayAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinAvg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowAvgBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinBitAnd : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitAndBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinBitOr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitOrBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinBoolAnd : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolAndBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinBoolOr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolOrBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinCount : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowCountBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinJSONArrayAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinJSONObjectAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinMax : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowMaxBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinMin : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowMinBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinStdDev : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowStdDevBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinStringAgg : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinSum : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowSumBuilder)) github.com/uptrace/bun/schema.QueryAppender
  METHOD WinVariance : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.WindowVarianceBuilder)) github.com/uptrace/bun/schema.QueryAppender
CONST ExprColumns : untyped string = "?Columns"
CONST ExprOperator : untyped string = "?Operator"
CONST ExprPKs : untyped string = "?PKs"
CONST ExprTableAlias : untyped string = "?TableAlias"
CONST ExprTableColumns : untyped string = "?TableColumns"
CONST ExprTableName : untyped string = "?TableName"
CONST ExprTablePKs : untyped string = "?TablePKs"
TYPE Field : github.com/coldsmirk/vef-framework-go/orm.Field
  METHOD AppendValue : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte, strct reflect.Value) []byte
  METHOD AppendValueOrDefault : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte, strct reflect.Value) []byte
  METHOD Clone : func() *github.com/uptrace/bun/schema.Field
  METHOD HasNilValue : func(v reflect.Value) bool
  METHOD HasZeroValue : func(v reflect.Value) bool
  METHOD ScanValue : func(strct reflect.Value, src any) error
  METHOD ScanWithCheck : func(fv reflect.Value, src any) error
  METHOD SkipUpdate : func() bool
  METHOD String : func() string
  METHOD Value : func(strct reflect.Value) reflect.Value
  METHOD WithIndex : func(path []int) *github.com/uptrace/bun/schema.Field
CONST FieldCreatedAt : untyped string = "CreatedAt"
CONST FieldCreatedBy : untyped string = "CreatedBy"
CONST FieldCreatedByName : untyped string = "CreatedByName"
CONST FieldID : untyped string = "ID"
CONST FieldUpdatedAt : untyped string = "UpdatedAt"
CONST FieldUpdatedBy : untyped string = "UpdatedBy"
CONST FieldUpdatedByName : untyped string = "UpdatedByName"
TYPE FirstValueBuilder : github.com/coldsmirk/vef-framework-go/orm.FirstValueBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.FirstValueBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.FirstValueBuilder
  METHOD IgnoreNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.FirstValueBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
  METHOD RespectNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.FirstValueBuilder
TYPE ForeignKeyBuilder : github.com/coldsmirk/vef-framework-go/orm.ForeignKeyBuilder
  METHOD Columns : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ForeignKeyBuilder
  METHOD Name : func(name string) github.com/coldsmirk/vef-framework-go/internal/orm.ForeignKeyBuilder
  METHOD OnDelete : func(action github.com/coldsmirk/vef-framework-go/internal/orm.ReferenceAction) github.com/coldsmirk/vef-framework-go/internal/orm.ForeignKeyBuilder
  METHOD OnUpdate : func(action github.com/coldsmirk/vef-framework-go/internal/orm.ReferenceAction) github.com/coldsmirk/vef-framework-go/internal/orm.ForeignKeyBuilder
  METHOD References : func(table string, columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ForeignKeyBuilder
CONST FrameBoundCurrentRow : github.com/coldsmirk/vef-framework-go/internal/orm.FrameBoundKind = 3
CONST FrameBoundFollowing : github.com/coldsmirk/vef-framework-go/internal/orm.FrameBoundKind = 5
TYPE FrameBoundKind : github.com/coldsmirk/vef-framework-go/orm.FrameBoundKind
  METHOD String : func() string
CONST FrameBoundNone : github.com/coldsmirk/vef-framework-go/internal/orm.FrameBoundKind = 0
CONST FrameBoundPreceding : github.com/coldsmirk/vef-framework-go/internal/orm.FrameBoundKind = 4
CONST FrameBoundUnboundedFollowing : github.com/coldsmirk/vef-framework-go/internal/orm.FrameBoundKind = 2
CONST FrameBoundUnboundedPreceding : github.com/coldsmirk/vef-framework-go/internal/orm.FrameBoundKind = 1
CONST FrameDefault : github.com/coldsmirk/vef-framework-go/internal/orm.FrameType = 0
CONST FrameGroups : github.com/coldsmirk/vef-framework-go/internal/orm.FrameType = 3
CONST FrameRange : github.com/coldsmirk/vef-framework-go/internal/orm.FrameType = 2
CONST FrameRows : github.com/coldsmirk/vef-framework-go/internal/orm.FrameType = 1
TYPE FrameType : github.com/coldsmirk/vef-framework-go/orm.FrameType
  METHOD String : func() string
CONST FromDefault : github.com/coldsmirk/vef-framework-go/internal/orm.FromDirection = 0
TYPE FromDirection : github.com/coldsmirk/vef-framework-go/orm.FromDirection
  METHOD String : func() string
CONST FromFirst : github.com/coldsmirk/vef-framework-go/internal/orm.FromDirection = 1
CONST FromLast : github.com/coldsmirk/vef-framework-go/internal/orm.FromDirection = 2
TYPE FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel
TYPE FullTrackedModel : github.com/coldsmirk/vef-framework-go/orm.FullTrackedModel
CONST FuzzyContains : github.com/coldsmirk/vef-framework-go/internal/orm.FuzzyKind = 2
CONST FuzzyEnds : github.com/coldsmirk/vef-framework-go/internal/orm.FuzzyKind = 1
TYPE FuzzyKind : github.com/coldsmirk/vef-framework-go/orm.FuzzyKind
  METHOD BuildPattern : func(value string) string
CONST FuzzyStarts : github.com/coldsmirk/vef-framework-go/internal/orm.FuzzyKind = 0
CONST IndexBRIN : github.com/coldsmirk/vef-framework-go/internal/orm.IndexMethod = 5
CONST IndexBTree : github.com/coldsmirk/vef-framework-go/internal/orm.IndexMethod = 0
CONST IndexGIN : github.com/coldsmirk/vef-framework-go/internal/orm.IndexMethod = 2
CONST IndexGiST : github.com/coldsmirk/vef-framework-go/internal/orm.IndexMethod = 3
CONST IndexHash : github.com/coldsmirk/vef-framework-go/internal/orm.IndexMethod = 1
TYPE IndexMethod : github.com/coldsmirk/vef-framework-go/orm.IndexMethod
  METHOD String : func() string
CONST IndexSPGiST : github.com/coldsmirk/vef-framework-go/internal/orm.IndexMethod = 4
TYPE InsertQuery : github.com/coldsmirk/vef-framework-go/orm.InsertQuery
  METHOD Apply : func(fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD ApplyIf : func(condition bool, fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD BuildCondition : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) interface{github.com/uptrace/bun/schema.QueryAppender; github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder}
  METHOD BuildSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) *github.com/uptrace/bun.SelectQuery
  METHOD Column : func(name string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD ColumnExpr : func(name string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD CreateSubQuery : func(subQuery *github.com/uptrace/bun.SelectQuery) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DB : func() github.com/coldsmirk/vef-framework-go/internal/orm.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Exclude : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD ExcludeAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD ExprBuilder : func() github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder
  METHOD GetTable : func() *github.com/uptrace/bun/schema.Table
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD ModelTable : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD OnConflict : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConflictBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD Query : func() github.com/uptrace/bun.Query
  METHOD Returning : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD ReturningAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD ReturningNone : func() github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD Select : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD SelectAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD String : func() string
  METHOD Table : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD TableExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD TableFrom : func(model any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD TableSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD With : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD WithOrderedValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD WithRecursive : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD WithValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
TYPE JSONArrayAggBuilder : github.com/coldsmirk/vef-framework-go/orm.JSONArrayAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.JSONArrayAggBuilder
TYPE JSONObjectAggBuilder : github.com/coldsmirk/vef-framework-go/orm.JSONObjectAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD KeyColumn : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD KeyExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.JSONObjectAggBuilder
CONST JoinCross : github.com/coldsmirk/vef-framework-go/internal/orm.JoinType = 5
CONST JoinDefault : github.com/coldsmirk/vef-framework-go/internal/orm.JoinType = 0
CONST JoinFull : github.com/coldsmirk/vef-framework-go/internal/orm.JoinType = 4
CONST JoinInner : github.com/coldsmirk/vef-framework-go/internal/orm.JoinType = 1
CONST JoinLeft : github.com/coldsmirk/vef-framework-go/internal/orm.JoinType = 2
CONST JoinRight : github.com/coldsmirk/vef-framework-go/internal/orm.JoinType = 3
TYPE JoinType : github.com/coldsmirk/vef-framework-go/orm.JoinType
  METHOD String : func() string
TYPE LagBuilder : github.com/coldsmirk/vef-framework-go/orm.LagBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.LagBuilder
  METHOD DefaultValue : func(value any) github.com/coldsmirk/vef-framework-go/internal/orm.LagBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.LagBuilder
  METHOD Offset : func(offset int) github.com/coldsmirk/vef-framework-go/internal/orm.LagBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowPartitionBuilder
TYPE LastValueBuilder : github.com/coldsmirk/vef-framework-go/orm.LastValueBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.LastValueBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.LastValueBuilder
  METHOD IgnoreNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.LastValueBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
  METHOD RespectNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.LastValueBuilder
TYPE LeadBuilder : github.com/coldsmirk/vef-framework-go/orm.LeadBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.LeadBuilder
  METHOD DefaultValue : func(value any) github.com/coldsmirk/vef-framework-go/internal/orm.LeadBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.LeadBuilder
  METHOD Offset : func(offset int) github.com/coldsmirk/vef-framework-go/internal/orm.LeadBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowPartitionBuilder
TYPE MaxBuilder : github.com/coldsmirk/vef-framework-go/orm.MaxBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.MaxBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.MaxBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MaxBuilder
TYPE MergeInsertBuilder : github.com/coldsmirk/vef-framework-go/orm.MergeInsertBuilder
  METHOD Value : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.MergeInsertBuilder
  METHOD ValueExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.MergeInsertBuilder
  METHOD Values : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeInsertBuilder
  METHOD ValuesAll : func(excludedColumns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeInsertBuilder
TYPE MergeQuery : github.com/coldsmirk/vef-framework-go/orm.MergeQuery
  METHOD Apply : func(fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD ApplyIf : func(condition bool, fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD BuildCondition : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) interface{github.com/uptrace/bun/schema.QueryAppender; github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder}
  METHOD BuildSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) *github.com/uptrace/bun.SelectQuery
  METHOD CreateSubQuery : func(subQuery *github.com/uptrace/bun.SelectQuery) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DB : func() github.com/coldsmirk/vef-framework-go/internal/orm.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD ExprBuilder : func() github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder
  METHOD GetTable : func() *github.com/uptrace/bun/schema.Table
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD ModelTable : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD On : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD Query : func() github.com/uptrace/bun.Query
  METHOD Returning : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD ReturningAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD ReturningNone : func() github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD String : func() string
  METHOD Table : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD TableExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD TableFrom : func(model any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD TableSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD Using : func(model any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD UsingExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD UsingSubQuery : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD UsingTable : func(table string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD WhenMatched : func(builder ...func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeWhenBuilder
  METHOD WhenNotMatched : func(builder ...func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeWhenBuilder
  METHOD WhenNotMatchedBySource : func(builder ...func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeWhenBuilder
  METHOD WhenNotMatchedByTarget : func(builder ...func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeWhenBuilder
  METHOD With : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD WithOrderedValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD WithRecursive : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD WithValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
TYPE MergeUpdateBuilder : github.com/coldsmirk/vef-framework-go/orm.MergeUpdateBuilder
  METHOD Set : func(column string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.MergeUpdateBuilder
  METHOD SetAll : func(excludedColumns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeUpdateBuilder
  METHOD SetColumns : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.MergeUpdateBuilder
  METHOD SetExpr : func(column string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.MergeUpdateBuilder
TYPE MergeWhenBuilder : github.com/coldsmirk/vef-framework-go/orm.MergeWhenBuilder
  METHOD ThenDelete : func() github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD ThenDoNothing : func() github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD ThenInsert : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.MergeInsertBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD ThenUpdate : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.MergeUpdateBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
TYPE MinBuilder : github.com/coldsmirk/vef-framework-go/orm.MinBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.MinBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.MinBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.MinBuilder
TYPE Model : github.com/coldsmirk/vef-framework-go/orm.Model
TYPE NTileBuilder : github.com/coldsmirk/vef-framework-go/orm.NTileBuilder
  METHOD Buckets : func(n int) github.com/coldsmirk/vef-framework-go/internal/orm.NTileBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
VAR NotNull : func() github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
TYPE NthValueBuilder : github.com/coldsmirk/vef-framework-go/orm.NthValueBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder
  METHOD FromFirst : func() github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder
  METHOD FromLast : func() github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder
  METHOD IgnoreNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder
  METHOD N : func(n int) github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
  METHOD RespectNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.NthValueBuilder
VAR Nullable : func() github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
CONST NullsDefault : github.com/coldsmirk/vef-framework-go/internal/orm.NullsMode = 0
CONST NullsIgnore : github.com/coldsmirk/vef-framework-go/internal/orm.NullsMode = 2
TYPE NullsMode : github.com/coldsmirk/vef-framework-go/orm.NullsMode
  METHOD String : func() string
CONST NullsRespect : github.com/coldsmirk/vef-framework-go/internal/orm.NullsMode = 1
CONST OperatorAnonymous : untyped string = "anonymous"
CONST OperatorCronJob : untyped string = "cron_job"
CONST OperatorSystem : untyped string = "system"
TYPE OrderBuilder : github.com/coldsmirk/vef-framework-go/orm.OrderBuilder
  METHOD Asc : func() github.com/coldsmirk/vef-framework-go/internal/orm.OrderBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.OrderBuilder
  METHOD Desc : func() github.com/coldsmirk/vef-framework-go/internal/orm.OrderBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.OrderBuilder
  METHOD NullsFirst : func() github.com/coldsmirk/vef-framework-go/internal/orm.OrderBuilder
  METHOD NullsLast : func() github.com/coldsmirk/vef-framework-go/internal/orm.OrderBuilder
TYPE PKField : github.com/coldsmirk/vef-framework-go/orm.PKField
  METHOD Set : func(model any, value any) error
  METHOD Value : func(model any) (any, error)
CONST PartitionHash : github.com/coldsmirk/vef-framework-go/internal/orm.PartitionStrategy = 2
CONST PartitionList : github.com/coldsmirk/vef-framework-go/internal/orm.PartitionStrategy = 1
CONST PartitionRange : github.com/coldsmirk/vef-framework-go/internal/orm.PartitionStrategy = 0
TYPE PartitionStrategy : github.com/coldsmirk/vef-framework-go/orm.PartitionStrategy
  METHOD String : func() string
TYPE PercentRankBuilder : github.com/coldsmirk/vef-framework-go/orm.PercentRankBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
CONST PlaceholderKeyOperator : untyped string = "Operator"
VAR PrimaryKey : func() github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
TYPE PrimaryKeyBuilder : github.com/coldsmirk/vef-framework-go/orm.PrimaryKeyBuilder
  METHOD Columns : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.PrimaryKeyBuilder
  METHOD Name : func(name string) github.com/coldsmirk/vef-framework-go/internal/orm.PrimaryKeyBuilder
TYPE QueryBuilder : github.com/coldsmirk/vef-framework-go/orm.QueryBuilder
  METHOD BuildCondition : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) interface{github.com/uptrace/bun/schema.QueryAppender; github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder}
  METHOD BuildSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) *github.com/uptrace/bun.SelectQuery
  METHOD CreateSubQuery : func(subQuery *github.com/uptrace/bun.SelectQuery) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DB : func() github.com/coldsmirk/vef-framework-go/internal/orm.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD ExprBuilder : func() github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder
  METHOD GetTable : func() *github.com/uptrace/bun/schema.Table
  METHOD Query : func() github.com/uptrace/bun.Query
  METHOD String : func() string
TYPE RankBuilder : github.com/coldsmirk/vef-framework-go/orm.RankBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE RawDefault : github.com/coldsmirk/vef-framework-go/orm.RawDefault
TYPE RawQuery : github.com/coldsmirk/vef-framework-go/orm.RawQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD Scan : func(ctx context.Context, dest ...any) error
TYPE ReferenceAction : github.com/coldsmirk/vef-framework-go/orm.ReferenceAction
  METHOD String : func() string
CONST ReferenceCascade : github.com/coldsmirk/vef-framework-go/internal/orm.ReferenceAction = 0
CONST ReferenceNoAction : github.com/coldsmirk/vef-framework-go/internal/orm.ReferenceAction = 4
CONST ReferenceRestrict : github.com/coldsmirk/vef-framework-go/internal/orm.ReferenceAction = 1
CONST ReferenceSetDefault : github.com/coldsmirk/vef-framework-go/internal/orm.ReferenceAction = 3
CONST ReferenceSetNull : github.com/coldsmirk/vef-framework-go/internal/orm.ReferenceAction = 2
VAR References : func(table string, columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
TYPE Relation : github.com/coldsmirk/vef-framework-go/orm.Relation
  METHOD References : func() bool
  METHOD String : func() string
TYPE RelationSpec : github.com/coldsmirk/vef-framework-go/orm.RelationSpec
TYPE RowNumberBuilder : github.com/coldsmirk/vef-framework-go/orm.RowNumberBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE SelectQuery : github.com/coldsmirk/vef-framework-go/orm.SelectQuery
  METHOD Apply : func(fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ApplyIf : func(condition bool, fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD BuildCondition : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) interface{github.com/uptrace/bun/schema.QueryAppender; github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder}
  METHOD BuildSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) *github.com/uptrace/bun.SelectQuery
  METHOD Count : func(ctx context.Context) (int64, error)
  METHOD CreateSubQuery : func(subQuery *github.com/uptrace/bun.SelectQuery) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD CrossJoin : func(model any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD CrossJoinExpr : func(eBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD CrossJoinSubQuery : func(sqBuilder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD CrossJoinTable : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DB : func() github.com/coldsmirk/vef-framework-go/internal/orm.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DistinctOnColumns : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DistinctOnExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Except : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ExceptAll : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Exclude : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ExcludeAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD Exists : func(ctx context.Context) (bool, error)
  METHOD ExprBuilder : func() github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder
  METHOD ForKeyShare : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForKeyShareNoWait : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForKeyShareSkipLocked : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForNoKeyUpdate : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForNoKeyUpdateNoWait : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForNoKeyUpdateSkipLocked : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForShare : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForShareNoWait : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForShareSkipLocked : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForUpdate : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForUpdateNoWait : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ForUpdateSkipLocked : func(tables ...any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD FullJoin : func(model any, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD FullJoinExpr : func(eBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD FullJoinSubQuery : func(sqBuilder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD FullJoinTable : func(name string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD GetTable : func() *github.com/uptrace/bun/schema.Table
  METHOD GroupBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD GroupByExpr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Having : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD IncludeDeleted : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Intersect : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD IntersectAll : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Join : func(model any, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD JoinExpr : func(eBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD JoinRelations : func(specs ...*github.com/coldsmirk/vef-framework-go/internal/orm.RelationSpec) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD JoinSubQuery : func(sqBuilder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD JoinTable : func(name string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD LeftJoin : func(model any, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD LeftJoinExpr : func(eBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD LeftJoinSubQuery : func(sqBuilder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD LeftJoinTable : func(name string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Limit : func(limit int) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD ModelTable : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Offset : func(offset int) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD OrderByExpr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Paginate : func(pageable github.com/coldsmirk/vef-framework-go/page.Pageable) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Query : func() github.com/uptrace/bun.Query
  METHOD Relation : func(name string, apply ...func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD RightJoin : func(model any, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD RightJoinExpr : func(eBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD RightJoinSubQuery : func(sqBuilder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), cBuilder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD RightJoinTable : func(name string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Rows : func(ctx context.Context) (*database/sql.Rows, error)
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD ScanAndCount : func(ctx context.Context, dest ...any) (int64, error)
  METHOD Select : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD SelectAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD SelectAs : func(column string, alias string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD SelectExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD SelectModelColumns : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD SelectModelPKs : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD String : func() string
  METHOD Table : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD TableExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD TableFrom : func(model any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD TableSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Union : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD UnionAll : func(func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD Where : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD WhereDeleted : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD WherePK : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD With : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD WithOrderedValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD WithRecursive : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD WithValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
CONST StatisticalDefault : github.com/coldsmirk/vef-framework-go/internal/orm.StatisticalMode = 0
TYPE StatisticalMode : github.com/coldsmirk/vef-framework-go/orm.StatisticalMode
  METHOD String : func() string
CONST StatisticalPopulation : github.com/coldsmirk/vef-framework-go/internal/orm.StatisticalMode = 1
CONST StatisticalSample : github.com/coldsmirk/vef-framework-go/internal/orm.StatisticalMode = 2
TYPE StdDevBuilder : github.com/coldsmirk/vef-framework-go/orm.StdDevBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.StdDevBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.StdDevBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.StdDevBuilder
  METHOD Population : func() github.com/coldsmirk/vef-framework-go/internal/orm.StdDevBuilder
  METHOD Sample : func() github.com/coldsmirk/vef-framework-go/internal/orm.StdDevBuilder
TYPE StringAggBuilder : github.com/coldsmirk/vef-framework-go/orm.StringAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD IgnoreNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD RespectNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
  METHOD Separator : func(separator string) github.com/coldsmirk/vef-framework-go/internal/orm.StringAggBuilder
TYPE SumBuilder : github.com/coldsmirk/vef-framework-go/orm.SumBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.SumBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.SumBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.SumBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.SumBuilder
TYPE Table : github.com/coldsmirk/vef-framework-go/orm.Table
  METHOD AppendNamedArg : func(gen github.com/uptrace/bun/schema.QueryGen, b []byte, name string, strct reflect.Value) ([]byte, bool)
  METHOD CheckPKs : func() error
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Field : func(name string) (*github.com/uptrace/bun/schema.Field, error)
  METHOD HasAfterScanRowHook : func() bool
  METHOD HasBeforeAppendModelHook : func() bool
  METHOD HasBeforeScanRowHook : func() bool
  METHOD HasField : func(name string) bool
  METHOD LookupField : func(name string) *github.com/uptrace/bun/schema.Field
  METHOD String : func() string
TYPE TableTarget : github.com/coldsmirk/vef-framework-go/orm.TableTarget[T github.com/coldsmirk/vef-framework-go/internal/orm.Executor]
  METHOD Model : func(model any) T
  METHOD Table : func(tables ...string) T
TYPE TruncateTableQuery : github.com/coldsmirk/vef-framework-go/orm.TruncateTableQuery
  METHOD Cascade : func() github.com/coldsmirk/vef-framework-go/internal/orm.TruncateTableQuery
  METHOD ContinueIdentity : func() github.com/coldsmirk/vef-framework-go/internal/orm.TruncateTableQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.TruncateTableQuery
  METHOD Restrict : func() github.com/coldsmirk/vef-framework-go/internal/orm.TruncateTableQuery
  METHOD String : func() string
  METHOD Table : func(tables ...string) github.com/coldsmirk/vef-framework-go/internal/orm.TruncateTableQuery
TYPE Tx : github.com/coldsmirk/vef-framework-go/orm.Tx
  METHOD BeginTx : func(ctx context.Context, opts *database/sql.TxOptions) (github.com/coldsmirk/vef-framework-go/internal/orm.Tx, error)
  METHOD Commit : func() error
  METHOD Connection : func(ctx context.Context) (*database/sql.Conn, error)
  METHOD ModelPKFields : func(model any) []*github.com/coldsmirk/vef-framework-go/internal/orm.PKField
  METHOD ModelPKs : func(model any) (map[string]any, error)
  METHOD NewAddColumn : func() github.com/coldsmirk/vef-framework-go/internal/orm.AddColumnQuery
  METHOD NewCreateIndex : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateIndexQuery
  METHOD NewCreateTable : func() github.com/coldsmirk/vef-framework-go/internal/orm.CreateTableQuery
  METHOD NewDelete : func() github.com/coldsmirk/vef-framework-go/internal/orm.DeleteQuery
  METHOD NewDropColumn : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropColumnQuery
  METHOD NewDropIndex : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropIndexQuery
  METHOD NewDropTable : func() github.com/coldsmirk/vef-framework-go/internal/orm.DropTableQuery
  METHOD NewInsert : func() github.com/coldsmirk/vef-framework-go/internal/orm.InsertQuery
  METHOD NewMerge : func() github.com/coldsmirk/vef-framework-go/internal/orm.MergeQuery
  METHOD NewRaw : func(query string, args ...any) github.com/coldsmirk/vef-framework-go/internal/orm.RawQuery
  METHOD NewSelect : func() github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD NewTruncateTable : func() github.com/coldsmirk/vef-framework-go/internal/orm.TruncateTableQuery
  METHOD NewUpdate : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD RegisterModel : func(models ...any)
  METHOD ResetModel : func(ctx context.Context, models ...any) error
  METHOD Rollback : func() error
  METHOD RunInReadOnlyTx : func(ctx context.Context, fn func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/internal/orm.DB) error) error
  METHOD RunInTx : func(ctx context.Context, fn func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/internal/orm.DB) error) error
  METHOD ScanRow : func(ctx context.Context, rows *database/sql.Rows, dest ...any) error
  METHOD ScanRows : func(ctx context.Context, rows *database/sql.Rows, dest ...any) error
  METHOD TableOf : func(model any) *github.com/uptrace/bun/schema.Table
  METHOD WithNamedArg : func(name string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.DB
VAR Unique : func() github.com/coldsmirk/vef-framework-go/internal/orm.ColumnConstraint
TYPE UniqueBuilder : github.com/coldsmirk/vef-framework-go/orm.UniqueBuilder
  METHOD Columns : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UniqueBuilder
  METHOD Name : func(name string) github.com/coldsmirk/vef-framework-go/internal/orm.UniqueBuilder
CONST UnitDay : github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit = 2
CONST UnitHour : github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit = 3
CONST UnitMinute : github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit = 4
CONST UnitMonth : github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit = 1
CONST UnitSecond : github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit = 5
CONST UnitYear : github.com/coldsmirk/vef-framework-go/internal/orm.DateTimeUnit = 0
TYPE UpdateQuery : github.com/coldsmirk/vef-framework-go/orm.UpdateQuery
  METHOD Apply : func(fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD ApplyIf : func(condition bool, fns ...github.com/coldsmirk/vef-framework-go/internal/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery]) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD BuildCondition : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) interface{github.com/uptrace/bun/schema.QueryAppender; github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder}
  METHOD BuildSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) *github.com/uptrace/bun.SelectQuery
  METHOD Bulk : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Column : func(name string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD ColumnExpr : func(name string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD CreateSubQuery : func(subQuery *github.com/uptrace/bun.SelectQuery) github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery
  METHOD DB : func() github.com/coldsmirk/vef-framework-go/internal/orm.DB
  METHOD Dialect : func() github.com/uptrace/bun/schema.Dialect
  METHOD Exclude : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD ExcludeAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Exec : func(ctx context.Context, dest ...any) (database/sql.Result, error)
  METHOD ExprBuilder : func() github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder
  METHOD GetTable : func() *github.com/uptrace/bun/schema.Table
  METHOD IncludeDeleted : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Limit : func(limit int) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Model : func(model any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD ModelTable : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD OmitZero : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD OrderByExpr : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Query : func() github.com/uptrace/bun.Query
  METHOD Returning : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD ReturningAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD ReturningNone : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Scan : func(ctx context.Context, dest ...any) error
  METHOD Select : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD SelectAll : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Set : func(name string, value any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD SetExpr : func(name string, builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD String : func() string
  METHOD Table : func(name string, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD TableExpr : func(builder func(github.com/coldsmirk/vef-framework-go/internal/orm.ExprBuilder) any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD TableFrom : func(model any, alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD TableSubQuery : func(builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery), alias ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD Where : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD WhereDeleted : func() github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD WherePK : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD With : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD WithOrderedValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD WithRecursive : func(name string, builder func(query github.com/coldsmirk/vef-framework-go/internal/orm.SelectQuery)) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
  METHOD WithValues : func(name string, model any) github.com/coldsmirk/vef-framework-go/internal/orm.UpdateQuery
TYPE VarianceBuilder : github.com/coldsmirk/vef-framework-go/orm.VarianceBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.VarianceBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.VarianceBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.VarianceBuilder
  METHOD Population : func() github.com/coldsmirk/vef-framework-go/internal/orm.VarianceBuilder
  METHOD Sample : func() github.com/coldsmirk/vef-framework-go/internal/orm.VarianceBuilder
TYPE WindowArrayAggBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowArrayAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD IgnoreNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
  METHOD RespectNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowArrayAggBuilder
TYPE WindowAvgBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowAvgBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowAvgBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowAvgBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowAvgBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowAvgBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowBitAndBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowBitAndBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitAndBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitAndBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitAndBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowBitOrBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowBitOrBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitOrBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitOrBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBitOrBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowBoolAndBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowBoolAndBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolAndBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolAndBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolAndBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowBoolOrBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowBoolOrBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolOrBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolOrBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowBoolOrBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowCountBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowCountBuilder
  METHOD All : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowCountBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowCountBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowCountBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowCountBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowCountBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowJSONArrayAggBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowJSONArrayAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONArrayAggBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowJSONObjectAggBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowJSONObjectAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD KeyColumn : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD KeyExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowJSONObjectAggBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowMaxBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowMaxBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowMaxBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowMaxBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowMaxBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowMinBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowMinBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowMinBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowMinBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowMinBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowStdDevBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowStdDevBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStdDevBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStdDevBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStdDevBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
  METHOD Population : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowStdDevBuilder
  METHOD Sample : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowStdDevBuilder
TYPE WindowStringAggBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowStringAggBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD IgnoreNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD OrderBy : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD OrderByDesc : func(columns ...string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD OrderByExpr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
  METHOD RespectNulls : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
  METHOD Separator : func(separator string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowStringAggBuilder
TYPE WindowSumBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowSumBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowSumBuilder
  METHOD Distinct : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowSumBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowSumBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowSumBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
TYPE WindowVarianceBuilder : github.com/coldsmirk/vef-framework-go/orm.WindowVarianceBuilder
  METHOD Column : func(column string) github.com/coldsmirk/vef-framework-go/internal/orm.WindowVarianceBuilder
  METHOD Expr : func(expr any) github.com/coldsmirk/vef-framework-go/internal/orm.WindowVarianceBuilder
  METHOD Filter : func(func(github.com/coldsmirk/vef-framework-go/internal/orm.ConditionBuilder)) github.com/coldsmirk/vef-framework-go/internal/orm.WindowVarianceBuilder
  METHOD Over : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowFrameablePartitionBuilder
  METHOD Population : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowVarianceBuilder
  METHOD Sample : func() github.com/coldsmirk/vef-framework-go/internal/orm.WindowVarianceBuilder

## github.com/coldsmirk/vef-framework-go/page
CONST DefaultPageNumber : int = 1
CONST DefaultPageSize : int = 15
CONST MaxPageSize : int = 1000
FUNC New : func[T any](pageable github.com/coldsmirk/vef-framework-go/page.Pageable, total int64, items []T) github.com/coldsmirk/vef-framework-go/page.Page[T]
TYPE Page : github.com/coldsmirk/vef-framework-go/page.Page[T any]
  FIELD Page : int [field_order=1 tag="json:\"page\""]
  FIELD Size : int [field_order=2 tag="json:\"size\""]
  FIELD Total : int64 [field_order=3 tag="json:\"total\""]
  FIELD Items : []T [field_order=4 tag="json:\"items\""]
  METHOD HasNext : func() bool
  METHOD HasPrevious : func() bool
  METHOD TotalPages : func() int
TYPE Pageable : github.com/coldsmirk/vef-framework-go/page.Pageable
  FIELD Page : int [field_order=1 tag="json:\"page\""]
  FIELD Size : int [field_order=2 tag="json:\"size\""]
  METHOD Normalize : func(size ...int)
  METHOD Offset : func() int

## github.com/coldsmirk/vef-framework-go/password
TYPE Argon2Option : github.com/coldsmirk/vef-framework-go/password.Argon2Option
TYPE BcryptOption : github.com/coldsmirk/vef-framework-go/password.BcryptOption
TYPE Encoder : github.com/coldsmirk/vef-framework-go/password.Encoder
  METHOD Encode : func(password string) (string, error)
  METHOD Matches : func(password string, encodedPassword string) bool
  METHOD UpgradeEncoding : func(encodedPassword string) bool
CONST EncoderArgon2 : github.com/coldsmirk/vef-framework-go/password.EncoderID = "argon2"
CONST EncoderBcrypt : github.com/coldsmirk/vef-framework-go/password.EncoderID = "bcrypt"
TYPE EncoderID : github.com/coldsmirk/vef-framework-go/password.EncoderID
CONST EncoderMd5 : github.com/coldsmirk/vef-framework-go/password.EncoderID = "md5"
CONST EncoderPbkdf2 : github.com/coldsmirk/vef-framework-go/password.EncoderID = "pbkdf2"
CONST EncoderPlaintext : github.com/coldsmirk/vef-framework-go/password.EncoderID = "plaintext"
CONST EncoderScrypt : github.com/coldsmirk/vef-framework-go/password.EncoderID = "scrypt"
CONST EncoderSha256 : github.com/coldsmirk/vef-framework-go/password.EncoderID = "sha256"
VAR ErrDefaultEncoderNotFound : error
VAR ErrInvalidCost : error
VAR ErrInvalidEncoderID : error
VAR ErrInvalidHashFormat : error
VAR ErrInvalidIterations : error
VAR ErrInvalidMemory : error
VAR ErrInvalidParallelism : error
TYPE Md5Option : github.com/coldsmirk/vef-framework-go/password.Md5Option
FUNC NewArgon2Encoder : func(opts ...github.com/coldsmirk/vef-framework-go/password.Argon2Option) github.com/coldsmirk/vef-framework-go/password.Encoder
FUNC NewBcryptEncoder : func(opts ...github.com/coldsmirk/vef-framework-go/password.BcryptOption) github.com/coldsmirk/vef-framework-go/password.Encoder
FUNC NewCompositeEncoder : func(defaultEncoderID github.com/coldsmirk/vef-framework-go/password.EncoderID, encoders map[github.com/coldsmirk/vef-framework-go/password.EncoderID]github.com/coldsmirk/vef-framework-go/password.Encoder) github.com/coldsmirk/vef-framework-go/password.Encoder
FUNC NewMd5Encoder : func(opts ...github.com/coldsmirk/vef-framework-go/password.Md5Option) github.com/coldsmirk/vef-framework-go/password.Encoder
FUNC NewPbkdf2Encoder : func(opts ...github.com/coldsmirk/vef-framework-go/password.Pbkdf2Option) github.com/coldsmirk/vef-framework-go/password.Encoder
FUNC NewPlaintextEncoder : func() github.com/coldsmirk/vef-framework-go/password.Encoder
FUNC NewScryptEncoder : func(opts ...github.com/coldsmirk/vef-framework-go/password.ScryptOption) github.com/coldsmirk/vef-framework-go/password.Encoder
FUNC NewSha256Encoder : func(opts ...github.com/coldsmirk/vef-framework-go/password.Sha256Option) github.com/coldsmirk/vef-framework-go/password.Encoder
TYPE Pbkdf2Option : github.com/coldsmirk/vef-framework-go/password.Pbkdf2Option
TYPE ScryptOption : github.com/coldsmirk/vef-framework-go/password.ScryptOption
TYPE Sha256Option : github.com/coldsmirk/vef-framework-go/password.Sha256Option
FUNC WithArgon2Iterations : func(iterations uint32) github.com/coldsmirk/vef-framework-go/password.Argon2Option
FUNC WithArgon2Memory : func(memory uint32) github.com/coldsmirk/vef-framework-go/password.Argon2Option
FUNC WithArgon2Parallelism : func(parallelism uint8) github.com/coldsmirk/vef-framework-go/password.Argon2Option
FUNC WithBcryptCost : func(cost int) github.com/coldsmirk/vef-framework-go/password.BcryptOption
FUNC WithMd5Salt : func(salt string) github.com/coldsmirk/vef-framework-go/password.Md5Option
FUNC WithMd5SaltPosition : func(position string) github.com/coldsmirk/vef-framework-go/password.Md5Option
FUNC WithPbkdf2HashFunction : func(hashFunction string) github.com/coldsmirk/vef-framework-go/password.Pbkdf2Option
FUNC WithPbkdf2Iterations : func(iterations int) github.com/coldsmirk/vef-framework-go/password.Pbkdf2Option
FUNC WithScryptN : func(n int) github.com/coldsmirk/vef-framework-go/password.ScryptOption
FUNC WithScryptP : func(p int) github.com/coldsmirk/vef-framework-go/password.ScryptOption
FUNC WithScryptR : func(r int) github.com/coldsmirk/vef-framework-go/password.ScryptOption
FUNC WithSha256Salt : func(salt string) github.com/coldsmirk/vef-framework-go/password.Sha256Option
FUNC WithSha256SaltPosition : func(position string) github.com/coldsmirk/vef-framework-go/password.Sha256Option

## github.com/coldsmirk/vef-framework-go/reflectx
CONST BreadthFirst : github.com/coldsmirk/vef-framework-go/reflectx.TraversalMode = 1
FUNC CollectMethods : func(target reflect.Value) map[string]reflect.Value
FUNC Contains : func(collection any, element any) bool
CONST Continue : github.com/coldsmirk/vef-framework-go/reflectx.VisitAction = 0
FUNC ConvertValue : func(sourceValue reflect.Value, targetType reflect.Type) (reflect.Value, error)
CONST DepthFirst : github.com/coldsmirk/vef-framework-go/reflectx.TraversalMode = 0
FUNC Equal : func(a any, b any) bool
VAR ErrCannotConvertType : error
TYPE FieldTypeVisitor : github.com/coldsmirk/vef-framework-go/reflectx.FieldTypeVisitor
TYPE FieldVisitor : github.com/coldsmirk/vef-framework-go/reflectx.FieldVisitor
FUNC FindMethod : func(target reflect.Value, name string) reflect.Value
FUNC GetStringMapValue : func(v reflect.Value) (map[string]string, bool)
FUNC GetStringSliceValue : func(v reflect.Value) ([]string, bool)
FUNC GetStringValue : func(v reflect.Value) (string, bool)
FUNC Indirect : func(t reflect.Type) reflect.Type
FUNC IsEmpty : func(value any) bool
FUNC IsFloat : func(value any) bool
FUNC IsInteger : func(value any) bool
FUNC IsNotEmpty : func(value any) bool
FUNC IsNumeric : func(value any) bool
FUNC IsPointerToStruct : func(t reflect.Type) bool
FUNC IsSignedInt : func(value any) bool
FUNC IsSimilarType : func(t1 reflect.Type, t2 reflect.Type) bool
FUNC IsStringMapType : func(t reflect.Type) bool
FUNC IsStringSliceType : func(t reflect.Type) bool
FUNC IsStringType : func(t reflect.Type) bool
FUNC IsTypeCompatible : func(sourceType reflect.Type, targetType reflect.Type) bool
FUNC IsUnsignedInt : func(value any) bool
TYPE MethodTypeVisitor : github.com/coldsmirk/vef-framework-go/reflectx.MethodTypeVisitor
TYPE MethodVisitor : github.com/coldsmirk/vef-framework-go/reflectx.MethodVisitor
FUNC SetStringMapValue : func(v reflect.Value, m map[string]string)
FUNC SetStringSliceValue : func(v reflect.Value, s []string)
FUNC SetStringValue : func(v reflect.Value, s string)
CONST SkipChildren : github.com/coldsmirk/vef-framework-go/reflectx.VisitAction = 2
CONST Stop : github.com/coldsmirk/vef-framework-go/reflectx.VisitAction = 1
TYPE StructTypeVisitor : github.com/coldsmirk/vef-framework-go/reflectx.StructTypeVisitor
TYPE StructVisitor : github.com/coldsmirk/vef-framework-go/reflectx.StructVisitor
TYPE TagConfig : github.com/coldsmirk/vef-framework-go/reflectx.TagConfig
  FIELD Name : string [field_order=1 tag=""]
  FIELD Value : string [field_order=2 tag=""]
VAR ToBool : func(i any) bool
VAR ToBoolE : func(i any) (bool, error)
FUNC ToDecimal : func(value any) github.com/coldsmirk/vef-framework-go/decimal.Decimal
FUNC ToDecimalE : func(value any) (github.com/coldsmirk/vef-framework-go/decimal.Decimal, error)
VAR ToFloat32 : func(i any) float32
VAR ToFloat32E : func(i any) (float32, error)
VAR ToFloat64 : func(i any) float64
VAR ToFloat64E : func(i any) (float64, error)
VAR ToInt : func(i any) int
VAR ToInt16 : func(i any) int16
VAR ToInt16E : func(i any) (int16, error)
VAR ToInt32 : func(i any) int32
VAR ToInt32E : func(i any) (int32, error)
VAR ToInt64 : func(i any) int64
VAR ToInt64E : func(i any) (int64, error)
VAR ToInt8 : func(i any) int8
VAR ToInt8E : func(i any) (int8, error)
VAR ToIntE : func(i any) (int, error)
VAR ToString : func(i any) string
VAR ToStringE : func(i any) (string, error)
VAR ToUint : func(i any) uint
VAR ToUint16 : func(i any) uint16
VAR ToUint16E : func(i any) (uint16, error)
VAR ToUint32 : func(i any) uint32
VAR ToUint32E : func(i any) (uint32, error)
VAR ToUint64 : func(i any) uint64
VAR ToUint64E : func(i any) (uint64, error)
VAR ToUint8 : func(i any) uint8
VAR ToUint8E : func(i any) (uint8, error)
VAR ToUintE : func(i any) (uint, error)
TYPE TraversalMode : github.com/coldsmirk/vef-framework-go/reflectx.TraversalMode
TYPE TypeVisitor : github.com/coldsmirk/vef-framework-go/reflectx.TypeVisitor
  FIELD VisitStructType : github.com/coldsmirk/vef-framework-go/reflectx.StructTypeVisitor [field_order=1 tag=""]
  FIELD VisitFieldType : github.com/coldsmirk/vef-framework-go/reflectx.FieldTypeVisitor [field_order=2 tag=""]
  FIELD VisitMethodType : github.com/coldsmirk/vef-framework-go/reflectx.MethodTypeVisitor [field_order=3 tag=""]
FUNC Visit : func(target reflect.Value, visitor github.com/coldsmirk/vef-framework-go/reflectx.Visitor, opts ...github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption)
TYPE VisitAction : github.com/coldsmirk/vef-framework-go/reflectx.VisitAction
FUNC VisitFor : func[T any](visitor github.com/coldsmirk/vef-framework-go/reflectx.TypeVisitor, opts ...github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption)
FUNC VisitOf : func(value any, visitor github.com/coldsmirk/vef-framework-go/reflectx.Visitor, opts ...github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption)
FUNC VisitType : func(targetType reflect.Type, visitor github.com/coldsmirk/vef-framework-go/reflectx.TypeVisitor, opts ...github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption)
TYPE Visitor : github.com/coldsmirk/vef-framework-go/reflectx.Visitor
  FIELD VisitStruct : github.com/coldsmirk/vef-framework-go/reflectx.StructVisitor [field_order=1 tag=""]
  FIELD VisitField : github.com/coldsmirk/vef-framework-go/reflectx.FieldVisitor [field_order=2 tag=""]
  FIELD VisitMethod : github.com/coldsmirk/vef-framework-go/reflectx.MethodVisitor [field_order=3 tag=""]
TYPE VisitorConfig : github.com/coldsmirk/vef-framework-go/reflectx.VisitorConfig
  FIELD TraversalMode : github.com/coldsmirk/vef-framework-go/reflectx.TraversalMode [field_order=1 tag=""]
  FIELD Recursive : bool [field_order=2 tag=""]
  FIELD DiveTag : github.com/coldsmirk/vef-framework-go/reflectx.TagConfig [field_order=3 tag=""]
  FIELD MaxDepth : int [field_order=4 tag=""]
TYPE VisitorOption : github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption
FUNC WithDisableRecursive : func() github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption
FUNC WithDiveTag : func(tagName string, tagValue string) github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption
FUNC WithMaxDepth : func(maxDepth int) github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption
FUNC WithTraversalMode : func(mode github.com/coldsmirk/vef-framework-go/reflectx.TraversalMode) github.com/coldsmirk/vef-framework-go/reflectx.VisitorOption

## github.com/coldsmirk/vef-framework-go/result
FUNC AsErr : func(err error) (github.com/coldsmirk/vef-framework-go/result.Error, bool)
FUNC Err : func(messageOrOptions ...any) github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrAccessDenied : github.com/coldsmirk/vef-framework-go/result.Error
CONST ErrCodeAccessDenied : untyped int = 1100
CONST ErrCodeBadRequest : untyped int = 1400
CONST ErrCodeDangerousSQL : untyped int = 1600
CONST ErrCodeDefault : untyped int = 2000
CONST ErrCodeForeignKeyViolation : untyped int = 2003
CONST ErrCodeNotFound : untyped int = 1200
CONST ErrCodeNotImplemented : untyped int = 1500
CONST ErrCodeRecordAlreadyExists : untyped int = 2002
CONST ErrCodeRecordNotFound : untyped int = 2001
CONST ErrCodeRequestTimeout : untyped int = 1402
CONST ErrCodeTooManyRequests : untyped int = 1401
CONST ErrCodeUnknown : untyped int = 1900
CONST ErrCodeUnsupportedMediaType : untyped int = 1300
VAR ErrDangerousSQL : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrForeignKeyViolation : github.com/coldsmirk/vef-framework-go/result.Error
CONST ErrMessage : untyped string = "error"
CONST ErrMessageAccessDenied : untyped string = "access_denied"
CONST ErrMessageDangerousSQL : untyped string = "dangerous_sql"
CONST ErrMessageForeignKeyViolation : untyped string = "foreign_key_violation"
CONST ErrMessageNotFound : untyped string = "not_found"
CONST ErrMessageRecordAlreadyExists : untyped string = "record_already_exists"
CONST ErrMessageRecordNotFound : untyped string = "record_not_found"
CONST ErrMessageRequestTimeout : untyped string = "request_timeout"
CONST ErrMessageTooManyRequests : untyped string = "too_many_requests"
CONST ErrMessageUnknown : untyped string = "unknown_error"
CONST ErrMessageUnsupportedMediaType : untyped string = "unsupported_media_type"
FUNC ErrNotImplemented : func(message string) github.com/coldsmirk/vef-framework-go/result.Error
TYPE ErrOption : github.com/coldsmirk/vef-framework-go/result.ErrOption
VAR ErrRecordAlreadyExists : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrRecordNotFound : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrRequestTimeout : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTooManyRequests : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUnknown : github.com/coldsmirk/vef-framework-go/result.Error
FUNC Errf : func(format string, args ...any) github.com/coldsmirk/vef-framework-go/result.Error
TYPE Error : github.com/coldsmirk/vef-framework-go/result.Error
  FIELD Code : int [field_order=1 tag=""]
  FIELD Message : string [field_order=2 tag=""]
  FIELD Status : int [field_order=3 tag=""]
  METHOD Error : func() string
  METHOD Is : func(target error) bool
FUNC IsRecordNotFound : func(err error) bool
FUNC Ok : func(dataOrOptions ...any) github.com/coldsmirk/vef-framework-go/result.Result
CONST OkCode : untyped int = 0
CONST OkMessage : untyped string = "ok"
TYPE OkOption : github.com/coldsmirk/vef-framework-go/result.OkOption
TYPE Result : github.com/coldsmirk/vef-framework-go/result.Result
  FIELD Code : int [field_order=1 tag="json:\"code\""]
  FIELD Message : string [field_order=2 tag="json:\"message\""]
  FIELD Data : any [field_order=3 tag="json:\"data\""]
  METHOD IsOk : func() bool
  METHOD Response : func(ctx github.com/gofiber/fiber/v3.Ctx, status ...int) error
FUNC WithCode : func(code int) github.com/coldsmirk/vef-framework-go/result.ErrOption
FUNC WithMessage : func(message string) github.com/coldsmirk/vef-framework-go/result.OkOption
FUNC WithMessagef : func(format string, args ...any) github.com/coldsmirk/vef-framework-go/result.OkOption
FUNC WithStatus : func(status int) github.com/coldsmirk/vef-framework-go/result.ErrOption

## github.com/coldsmirk/vef-framework-go/schema
TYPE Check : github.com/coldsmirk/vef-framework-go/schema.Check
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Expr : string [field_order=2 tag="json:\"expr\""]
TYPE Column : github.com/coldsmirk/vef-framework-go/schema.Column
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Type : string [field_order=2 tag="json:\"type\""]
  FIELD Nullable : bool [field_order=3 tag="json:\"nullable\""]
  FIELD Default : string [field_order=4 tag="json:\"default,omitempty\""]
  FIELD Comment : string [field_order=5 tag="json:\"comment,omitempty\""]
  FIELD IsPrimaryKey : bool [field_order=6 tag="json:\"isPrimaryKey,omitempty\""]
  FIELD IsAutoIncrement : bool [field_order=7 tag="json:\"isAutoIncrement,omitempty\""]
CONST ErrCodeTableNotFound : untyped int = 2300
VAR ErrTableNotFound : github.com/coldsmirk/vef-framework-go/result.Error
TYPE ForeignKey : github.com/coldsmirk/vef-framework-go/schema.ForeignKey
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Columns : []string [field_order=2 tag="json:\"columns\""]
  FIELD RefTable : string [field_order=3 tag="json:\"refTable\""]
  FIELD RefColumns : []string [field_order=4 tag="json:\"refColumns\""]
  FIELD OnUpdate : string [field_order=5 tag="json:\"onUpdate,omitempty\""]
  FIELD OnDelete : string [field_order=6 tag="json:\"onDelete,omitempty\""]
TYPE Index : github.com/coldsmirk/vef-framework-go/schema.Index
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Columns : []string [field_order=2 tag="json:\"columns\""]
TYPE PrimaryKey : github.com/coldsmirk/vef-framework-go/schema.PrimaryKey
  FIELD Name : string [field_order=1 tag="json:\"name,omitempty\""]
  FIELD Columns : []string [field_order=2 tag="json:\"columns\""]
TYPE Service : github.com/coldsmirk/vef-framework-go/schema.Service
  METHOD GetTableSchema : func(ctx context.Context, name string) (*github.com/coldsmirk/vef-framework-go/schema.TableSchema, error)
  METHOD ListTables : func(ctx context.Context) ([]github.com/coldsmirk/vef-framework-go/schema.Table, error)
  METHOD ListViews : func(ctx context.Context) ([]github.com/coldsmirk/vef-framework-go/schema.View, error)
TYPE Table : github.com/coldsmirk/vef-framework-go/schema.Table
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Schema : string [field_order=2 tag="json:\"schema,omitempty\""]
  FIELD Comment : string [field_order=3 tag="json:\"comment,omitempty\""]
TYPE TableSchema : github.com/coldsmirk/vef-framework-go/schema.TableSchema
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Schema : string [field_order=2 tag="json:\"schema,omitempty\""]
  FIELD Comment : string [field_order=3 tag="json:\"comment,omitempty\""]
  FIELD Columns : []github.com/coldsmirk/vef-framework-go/schema.Column [field_order=4 tag="json:\"columns\""]
  FIELD PrimaryKey : *github.com/coldsmirk/vef-framework-go/schema.PrimaryKey [field_order=5 tag="json:\"primaryKey,omitempty\""]
  FIELD Indexes : []github.com/coldsmirk/vef-framework-go/schema.Index [field_order=6 tag="json:\"indexes,omitempty\""]
  FIELD UniqueKeys : []github.com/coldsmirk/vef-framework-go/schema.UniqueKey [field_order=7 tag="json:\"uniqueKeys,omitempty\""]
  FIELD ForeignKeys : []github.com/coldsmirk/vef-framework-go/schema.ForeignKey [field_order=8 tag="json:\"foreignKeys,omitempty\""]
  FIELD Checks : []github.com/coldsmirk/vef-framework-go/schema.Check [field_order=9 tag="json:\"checks,omitempty\""]
TYPE UniqueKey : github.com/coldsmirk/vef-framework-go/schema.UniqueKey
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Columns : []string [field_order=2 tag="json:\"columns\""]
TYPE View : github.com/coldsmirk/vef-framework-go/schema.View
  FIELD Name : string [field_order=1 tag="json:\"name\""]
  FIELD Schema : string [field_order=2 tag="json:\"schema,omitempty\""]
  FIELD Definition : string [field_order=3 tag="json:\"definition\""]
  FIELD Comment : string [field_order=4 tag="json:\"comment,omitempty\""]
  FIELD Columns : []string [field_order=5 tag="json:\"columns,omitempty\""]

## github.com/coldsmirk/vef-framework-go/search
FUNC Applier : func[T any]() func(T) github.com/coldsmirk/vef-framework-go/orm.ApplyFunc[github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder]
CONST AttrAlias : untyped string = "alias"
CONST AttrColumn : untyped string = "column"
CONST AttrDive : untyped string = "dive"
CONST AttrOperator : untyped string = "operator"
CONST AttrParams : untyped string = "params"
CONST Between : github.com/coldsmirk/vef-framework-go/search.Operator = "between"
CONST Contains : github.com/coldsmirk/vef-framework-go/search.Operator = "contains"
CONST ContainsIgnoreCase : github.com/coldsmirk/vef-framework-go/search.Operator = "iContains"
CONST EndsWith : github.com/coldsmirk/vef-framework-go/search.Operator = "endsWith"
CONST EndsWithIgnoreCase : github.com/coldsmirk/vef-framework-go/search.Operator = "iEndsWith"
CONST Equals : github.com/coldsmirk/vef-framework-go/search.Operator = "eq"
CONST GreaterThan : github.com/coldsmirk/vef-framework-go/search.Operator = "gt"
CONST GreaterThanOrEqual : github.com/coldsmirk/vef-framework-go/search.Operator = "gte"
CONST IgnoreField : untyped string = "-"
CONST In : github.com/coldsmirk/vef-framework-go/search.Operator = "in"
CONST IsNotNull : github.com/coldsmirk/vef-framework-go/search.Operator = "isNotNull"
CONST IsNull : github.com/coldsmirk/vef-framework-go/search.Operator = "isNull"
CONST LessThan : github.com/coldsmirk/vef-framework-go/search.Operator = "lt"
CONST LessThanOrEqual : github.com/coldsmirk/vef-framework-go/search.Operator = "lte"
FUNC New : func(typ reflect.Type) github.com/coldsmirk/vef-framework-go/search.Search
FUNC NewFor : func[T any]() github.com/coldsmirk/vef-framework-go/search.Search
CONST NotBetween : github.com/coldsmirk/vef-framework-go/search.Operator = "notBetween"
CONST NotContains : github.com/coldsmirk/vef-framework-go/search.Operator = "notContains"
CONST NotContainsIgnoreCase : github.com/coldsmirk/vef-framework-go/search.Operator = "iNotContains"
CONST NotEndsWith : github.com/coldsmirk/vef-framework-go/search.Operator = "notEndsWith"
CONST NotEndsWithIgnoreCase : github.com/coldsmirk/vef-framework-go/search.Operator = "iNotEndsWith"
CONST NotEquals : github.com/coldsmirk/vef-framework-go/search.Operator = "neq"
CONST NotIn : github.com/coldsmirk/vef-framework-go/search.Operator = "notIn"
CONST NotStartsWith : github.com/coldsmirk/vef-framework-go/search.Operator = "notStartsWith"
CONST NotStartsWithIgnoreCase : github.com/coldsmirk/vef-framework-go/search.Operator = "iNotStartsWith"
TYPE Operator : github.com/coldsmirk/vef-framework-go/search.Operator
CONST ParamDelimiter : untyped string = "delimiter"
CONST ParamType : untyped string = "type"
TYPE Search : github.com/coldsmirk/vef-framework-go/search.Search
  METHOD Apply : func(cb github.com/coldsmirk/vef-framework-go/orm.ConditionBuilder, target any, defaultAlias ...string)
CONST StartsWith : github.com/coldsmirk/vef-framework-go/search.Operator = "startsWith"
CONST StartsWithIgnoreCase : github.com/coldsmirk/vef-framework-go/search.Operator = "iStartsWith"
CONST TagSearch : untyped string = "search"
CONST TypeDate : untyped string = "date"
CONST TypeDateTime : untyped string = "datetime"
CONST TypeDecimal : untyped string = "dec"
CONST TypeInt : untyped string = "int"
CONST TypeTime : untyped string = "time"

## github.com/coldsmirk/vef-framework-go/security
TYPE AllDataScope : github.com/coldsmirk/vef-framework-go/security.AllDataScope
  METHOD Apply : func(*github.com/coldsmirk/vef-framework-go/security.Principal, github.com/coldsmirk/vef-framework-go/orm.SelectQuery) error
  METHOD Key : func() string
  METHOD Priority : func() int
  METHOD Supports : func(*github.com/coldsmirk/vef-framework-go/security.Principal, *github.com/coldsmirk/vef-framework-go/orm.Table) bool
TYPE AuthManager : github.com/coldsmirk/vef-framework-go/security.AuthManager
  METHOD Authenticate : func(ctx context.Context, authentication github.com/coldsmirk/vef-framework-go/security.Authentication) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
CONST AuthSchemeBearer : untyped string = "Bearer"
TYPE AuthTokens : github.com/coldsmirk/vef-framework-go/security.AuthTokens
  FIELD AccessToken : string [field_order=1 tag="json:\"accessToken\""]
  FIELD RefreshToken : string [field_order=2 tag="json:\"refreshToken\""]
TYPE Authentication : github.com/coldsmirk/vef-framework-go/security.Authentication
  FIELD Type : string [field_order=1 tag="json:\"type\""]
  FIELD Principal : string [field_order=2 tag="json:\"principal\""]
  FIELD Credentials : any [field_order=3 tag="json:\"credentials\""]
TYPE Authenticator : github.com/coldsmirk/vef-framework-go/security.Authenticator
  METHOD Authenticate : func(ctx context.Context, authentication github.com/coldsmirk/vef-framework-go/security.Authentication) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
  METHOD Supports : func(authType string) bool
TYPE CachedRolePermissionsLoader : github.com/coldsmirk/vef-framework-go/security.CachedRolePermissionsLoader
  METHOD LoadPermissions : func(ctx context.Context, role string) (map[string]github.com/coldsmirk/vef-framework-go/security.DataScope, error)
TYPE ChallengeProvider : github.com/coldsmirk/vef-framework-go/security.ChallengeProvider
  METHOD Evaluate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.LoginChallenge, error)
  METHOD Order : func() int
  METHOD Resolve : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, response any) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
  METHOD Type : func() string
TYPE ChallengeState : github.com/coldsmirk/vef-framework-go/security.ChallengeState
  FIELD Principal : *github.com/coldsmirk/vef-framework-go/security.Principal [field_order=1 tag=""]
  FIELD Username : string [field_order=2 tag=""]
  FIELD Pending : []string [field_order=3 tag=""]
  FIELD Resolved : []string [field_order=4 tag=""]
CONST ChallengeTokenExpires : time.Duration = 300000000000
TYPE ChallengeTokenStore : github.com/coldsmirk/vef-framework-go/security.ChallengeTokenStore
  METHOD Generate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, username string, pending []string, resolved []string) (string, error)
  METHOD Parse : func(ctx context.Context, token string) (*github.com/coldsmirk/vef-framework-go/security.ChallengeState, error)
CONST ChallengeTypeDepartmentSelection : untyped string = "department_selection"
CONST ChallengeTypeEmail : untyped string = "email_otp"
CONST ChallengeTypePasswordChange : untyped string = "password_change"
CONST ChallengeTypeSMS : untyped string = "sms_otp"
CONST ChallengeTypeTOTP : untyped string = "totp"
CONST ClaimChallengePending : untyped string = "pnd"
CONST ClaimChallengePrincipalName : untyped string = "pnm"
CONST ClaimChallengePrincipalType : untyped string = "ptp"
CONST ClaimChallengeResolved : untyped string = "rsd"
CONST ClaimChallengeUsername : untyped string = "unm"
TYPE DataPermissionApplier : github.com/coldsmirk/vef-framework-go/security.DataPermissionApplier
  METHOD Apply : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery) error
TYPE DataPermissionResolver : github.com/coldsmirk/vef-framework-go/security.DataPermissionResolver
  METHOD ResolveDataScope : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, permission string) (github.com/coldsmirk/vef-framework-go/security.DataScope, error)
TYPE DataScope : github.com/coldsmirk/vef-framework-go/security.DataScope
  METHOD Apply : func(principal *github.com/coldsmirk/vef-framework-go/security.Principal, query github.com/coldsmirk/vef-framework-go/orm.SelectQuery) error
  METHOD Key : func() string
  METHOD Priority : func() int
  METHOD Supports : func(principal *github.com/coldsmirk/vef-framework-go/security.Principal, table *github.com/coldsmirk/vef-framework-go/orm.Table) bool
CONST DefaultJWTAudience : untyped string = "vef-app"
CONST DefaultJWTSecret : untyped string = "af6675678bd81ad7c93c4a51d122ef61e9750fe5d42ceac1c33b293f36bc14c2"
TYPE DeliveredCodeSender : github.com/coldsmirk/vef-framework-go/security.DeliveredCodeSender
  METHOD Send : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) error
TYPE DeliveredCodeVerifier : github.com/coldsmirk/vef-framework-go/security.DeliveredCodeVerifier
  METHOD Verify : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, code string) (bool, error)
TYPE DepartmentLoader : github.com/coldsmirk/vef-framework-go/security.DepartmentLoader
  METHOD LoadDepartments : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) ([]github.com/coldsmirk/vef-framework-go/security.DepartmentOption, error)
TYPE DepartmentOption : github.com/coldsmirk/vef-framework-go/security.DepartmentOption
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD Name : string [field_order=2 tag="json:\"name\""]
TYPE DepartmentSelectionChallengeData : github.com/coldsmirk/vef-framework-go/security.DepartmentSelectionChallengeData
  FIELD Departments : []github.com/coldsmirk/vef-framework-go/security.DepartmentOption [field_order=1 tag="json:\"departments\""]
  FIELD Meta : map[string]any [field_order=2 tag="json:\"meta,omitempty\""]
TYPE DepartmentSelectionChallengeProvider : github.com/coldsmirk/vef-framework-go/security.DepartmentSelectionChallengeProvider
  METHOD Evaluate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.LoginChallenge, error)
  METHOD Order : func() int
  METHOD Resolve : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, response any) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
  METHOD Type : func() string
TYPE DepartmentSelector : github.com/coldsmirk/vef-framework-go/security.DepartmentSelector
  METHOD SelectDepartment : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, departmentID string) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
FUNC ErrAccountLocked : func(retryAfter time.Duration) github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrAppIDRequired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrAuthHeaderInvalid : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrAuthHeaderMissing : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrChallengeResolveFailed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrChallengeTokenInvalid : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrChallengeTypeInvalid : github.com/coldsmirk/vef-framework-go/result.Error
CONST ErrCodeAccountLocked : untyped int = 1023
CONST ErrCodeAppIDRequired : untyped int = 1009
CONST ErrCodeAuthHeaderInvalid : untyped int = 1022
CONST ErrCodeAuthHeaderMissing : untyped int = 1021
CONST ErrCodeChallengeResolveFailed : untyped int = 1034
CONST ErrCodeChallengeTokenInvalid : untyped int = 1031
CONST ErrCodeChallengeTypeInvalid : untyped int = 1033
CONST ErrCodeCredentialsInvalid : untyped int = 1008
CONST ErrCodeDepartmentRequired : untyped int = 1038
CONST ErrCodeExternalAppDisabled : untyped int = 1015
CONST ErrCodeExternalAppNotFound : untyped int = 1014
CONST ErrCodeIPNotAllowed : untyped int = 1016
CONST ErrCodeNewPasswordRequired : untyped int = 1037
CONST ErrCodeNonceAlreadyUsed : untyped int = 1020
CONST ErrCodeNonceInvalid : untyped int = 1019
CONST ErrCodeNonceRequired : untyped int = 1018
CONST ErrCodeOTPCodeInvalid : untyped int = 1036
CONST ErrCodeOTPCodeRequired : untyped int = 1035
CONST ErrCodePasswordPolicyViolation : untyped int = 1050
CONST ErrCodePrincipalInvalid : untyped int = 1007
CONST ErrCodeSignatureExpired : untyped int = 1013
CONST ErrCodeSignatureInvalid : untyped int = 1017
CONST ErrCodeSignatureRequired : untyped int = 1011
CONST ErrCodeTimestampInvalid : untyped int = 1012
CONST ErrCodeTimestampRequired : untyped int = 1010
CONST ErrCodeTokenExpired : untyped int = 1002
CONST ErrCodeTokenInvalid : untyped int = 1003
CONST ErrCodeTokenInvalidAudience : untyped int = 1006
CONST ErrCodeTokenInvalidIssuer : untyped int = 1005
CONST ErrCodeTokenNotValidYet : untyped int = 1004
CONST ErrCodeTooManyConcurrentSessions : untyped int = 1024
CONST ErrCodeUnauthenticated : untyped int = 1000
CONST ErrCodeUnsupportedAuthenticationType : untyped int = 1001
FUNC ErrCredentialsInvalid : func(message string) github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrDecodeJWTSecretFailed : error
VAR ErrDecodeSignatureSecretFailed : error
VAR ErrDepartmentRequired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrExternalAppDetailsNotStruct : error
VAR ErrExternalAppDisabled : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrExternalAppNotFound : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrGenerateJWTSecretFailed : error
VAR ErrIPNotAllowed : github.com/coldsmirk/vef-framework-go/result.Error
CONST ErrMessageAccountLocked : untyped string = "security_account_locked"
CONST ErrMessageChallengeResolveFailed : untyped string = "security_challenge_resolve_failed"
CONST ErrMessageCredentialsFormatInvalid : untyped string = "security_credentials_format_invalid"
CONST ErrMessageExternalAppLoaderNotImplemented : untyped string = "security_external_app_loader_not_implemented"
CONST ErrMessagePasswordTooFewCharClasses : untyped string = "security_password_too_few_char_classes"
CONST ErrMessagePasswordTooLong : untyped string = "security_password_too_long"
CONST ErrMessagePasswordTooShort : untyped string = "security_password_too_short"
CONST ErrMessageUnauthenticated : untyped string = "security_unauthenticated"
CONST ErrMessageUnsupportedAuthenticationType : untyped string = "security_unsupported_authentication_type"
CONST ErrMessageUserInfoLoaderNotImplemented : untyped string = "security_user_info_loader_not_implemented"
CONST ErrMessageUserLoaderNotImplemented : untyped string = "security_user_loader_not_implemented"
VAR ErrNewPasswordRequired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrNonceAlreadyUsed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrNonceInvalid : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrNonceRequired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrOTPCodeInvalid : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrOTPCodeRequired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrPasswordBlocked : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrPasswordContainsIdentity : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrPasswordMissingDigit : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrPasswordMissingLowercase : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrPasswordMissingSymbol : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrPasswordMissingUppercase : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrPasswordReused : github.com/coldsmirk/vef-framework-go/result.Error
FUNC ErrPasswordTooFewCharClasses : func(minClasses int) github.com/coldsmirk/vef-framework-go/result.Error
FUNC ErrPasswordTooLong : func(maxLength int) github.com/coldsmirk/vef-framework-go/result.Error
FUNC ErrPasswordTooShort : func(minLength int) github.com/coldsmirk/vef-framework-go/result.Error
FUNC ErrPrincipalInvalid : func(message string) github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrQueryModelNotSet : error
VAR ErrQueryNotQueryBuilder : error
VAR ErrSignatureExpired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrSignatureInvalid : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrSignatureRequired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrSignatureSecretRequired : error
VAR ErrTimestampInvalid : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTimestampRequired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTokenExpired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTokenInvalid : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTokenInvalidAudience : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTokenInvalidIssuer : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTokenNotValidYet : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTooManyConcurrentSessions : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUnauthenticated : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUserDetailsNotStruct : error
TYPE ExpiryPasswordChangeChecker : github.com/coldsmirk/vef-framework-go/security.ExpiryPasswordChangeChecker
  METHOD Check : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.PasswordChangeChallengeData, error)
TYPE ExternalAppConfig : github.com/coldsmirk/vef-framework-go/security.ExternalAppConfig
  FIELD Enabled : bool [field_order=1 tag="json:\"enabled\""]
  FIELD IPWhitelist : string [field_order=2 tag="json:\"ipWhitelist\""]
TYPE ExternalAppLoader : github.com/coldsmirk/vef-framework-go/security.ExternalAppLoader
  METHOD LoadByID : func(ctx context.Context, id string) (*github.com/coldsmirk/vef-framework-go/security.Principal, string, error)
TYPE Gender : github.com/coldsmirk/vef-framework-go/security.Gender
CONST GenderFemale : github.com/coldsmirk/vef-framework-go/security.Gender = "female"
CONST GenderMale : github.com/coldsmirk/vef-framework-go/security.Gender = "male"
CONST GenderUnknown : github.com/coldsmirk/vef-framework-go/security.Gender = "unknown"
FUNC GenerateOpaqueToken : func() (string, error)
FUNC GenerateSecret : func() (string, error)
FUNC HashOpaqueToken : func(token string) string
TYPE IPWhitelist : github.com/coldsmirk/vef-framework-go/security.IPWhitelist
  FIELD Entries : []string [field_order=1 tag=""]
TYPE IPWhitelistLoader : github.com/coldsmirk/vef-framework-go/security.IPWhitelistLoader
  METHOD LoadByName : func(ctx context.Context, name string) (*github.com/coldsmirk/vef-framework-go/security.IPWhitelist, error)
TYPE IPWhitelistValidator : github.com/coldsmirk/vef-framework-go/security.IPWhitelistValidator
  METHOD IsAllowed : func(ipStr string) bool
  METHOD IsEmpty : func() bool
TYPE JWT : github.com/coldsmirk/vef-framework-go/security.JWT
  METHOD Generate : func(claimsBuilder *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder, expires time.Duration, notBefore time.Duration) (string, error)
  METHOD Parse : func(tokenString string) (*github.com/coldsmirk/vef-framework-go/security.JWTClaimsAccessor, error)
TYPE JWTChallengeTokenStore : github.com/coldsmirk/vef-framework-go/security.JWTChallengeTokenStore
  METHOD Generate : func(_ context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, username string, pending []string, resolved []string) (string, error)
  METHOD Parse : func(_ context.Context, token string) (*github.com/coldsmirk/vef-framework-go/security.ChallengeState, error)
TYPE JWTClaimsAccessor : github.com/coldsmirk/vef-framework-go/security.JWTClaimsAccessor
  METHOD Claim : func(key string) any
  METHOD Details : func() any
  METHOD ID : func() string
  METHOD Roles : func() []string
  METHOD Subject : func() string
  METHOD Type : func() string
TYPE JWTClaimsBuilder : github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
  METHOD Claim : func(key string) (any, bool)
  METHOD Details : func() (any, bool)
  METHOD ID : func() (string, bool)
  METHOD Roles : func() ([]string, bool)
  METHOD Subject : func() (string, bool)
  METHOD Type : func() (string, bool)
  METHOD WithClaim : func(key string, value any) *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
  METHOD WithDetails : func(details any) *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
  METHOD WithID : func(id string) *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
  METHOD WithRoles : func(roles []string) *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
  METHOD WithSubject : func(subject string) *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
  METHOD WithType : func(typ string) *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
TYPE JWTConfig : github.com/coldsmirk/vef-framework-go/security.JWTConfig
  FIELD Secret : string [field_order=1 tag=""]
  FIELD Audience : string [field_order=2 tag=""]
CONST JWTIssuer : untyped string = "vef"
TYPE LockoutKey : github.com/coldsmirk/vef-framework-go/security.LockoutKey
CONST LockoutKeyIP : github.com/coldsmirk/vef-framework-go/security.LockoutKey = "ip"
CONST LockoutKeyUser : github.com/coldsmirk/vef-framework-go/security.LockoutKey = "user"
CONST LockoutKeyUserIP : github.com/coldsmirk/vef-framework-go/security.LockoutKey = "user_ip"
TYPE LockoutPolicy : github.com/coldsmirk/vef-framework-go/security.LockoutPolicy
  FIELD MaxFailures : int [field_order=1 tag=""]
  FIELD Window : time.Duration [field_order=2 tag=""]
  FIELD LockDuration : time.Duration [field_order=3 tag=""]
  FIELD Strategy : github.com/coldsmirk/vef-framework-go/security.LockoutStrategy [field_order=4 tag=""]
  FIELD BackoffBase : time.Duration [field_order=5 tag=""]
  FIELD BackoffMax : time.Duration [field_order=6 tag=""]
  FIELD Key : github.com/coldsmirk/vef-framework-go/security.LockoutKey [field_order=7 tag=""]
TYPE LockoutStrategy : github.com/coldsmirk/vef-framework-go/security.LockoutStrategy
CONST LockoutStrategyBackoff : github.com/coldsmirk/vef-framework-go/security.LockoutStrategy = "backoff"
CONST LockoutStrategyLock : github.com/coldsmirk/vef-framework-go/security.LockoutStrategy = "lock"
TYPE LoginAttempt : github.com/coldsmirk/vef-framework-go/security.LoginAttempt
  FIELD Identity : string [field_order=1 tag=""]
  FIELD ClientIP : string [field_order=2 tag=""]
TYPE LoginChallenge : github.com/coldsmirk/vef-framework-go/security.LoginChallenge
  FIELD Type : string [field_order=1 tag="json:\"type\""]
  FIELD Data : any [field_order=2 tag="json:\"data,omitempty\""]
  FIELD Required : bool [field_order=3 tag="json:\"required\""]
TYPE LoginDecision : github.com/coldsmirk/vef-framework-go/security.LoginDecision
  FIELD Allowed : bool [field_order=1 tag=""]
  FIELD RetryAfter : time.Duration [field_order=2 tag=""]
TYPE LoginEvent : github.com/coldsmirk/vef-framework-go/security.LoginEvent
  FIELD AuthType : string [field_order=1 tag="json:\"authType\""]
  FIELD UserID : *string [field_order=2 tag="json:\"userId\""]
  FIELD Username : string [field_order=3 tag="json:\"username\""]
  FIELD LoginIP : string [field_order=4 tag="json:\"loginIp\""]
  FIELD UserAgent : string [field_order=5 tag="json:\"userAgent\""]
  FIELD TraceID : string [field_order=6 tag="json:\"traceId\""]
  FIELD IsOk : bool [field_order=7 tag="json:\"isOk\""]
  FIELD FailReason : string [field_order=8 tag="json:\"failReason\""]
  FIELD ErrorCode : int [field_order=9 tag="json:\"errorCode\""]
  METHOD EventType : func() string
TYPE LoginEventParams : github.com/coldsmirk/vef-framework-go/security.LoginEventParams
  FIELD AuthType : string [field_order=1 tag=""]
  FIELD UserID : *string [field_order=2 tag=""]
  FIELD Username : string [field_order=3 tag=""]
  FIELD LoginIP : string [field_order=4 tag=""]
  FIELD UserAgent : string [field_order=5 tag=""]
  FIELD TraceID : string [field_order=6 tag=""]
  FIELD IsOk : bool [field_order=7 tag=""]
  FIELD FailReason : string [field_order=8 tag=""]
  FIELD ErrorCode : int [field_order=9 tag=""]
TYPE LoginGuard : github.com/coldsmirk/vef-framework-go/security.LoginGuard
  METHOD Check : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) (github.com/coldsmirk/vef-framework-go/security.LoginDecision, error)
  METHOD RecordFailure : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) (github.com/coldsmirk/vef-framework-go/security.LoginDecision, error)
  METHOD RecordSuccess : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) error
TYPE LoginResult : github.com/coldsmirk/vef-framework-go/security.LoginResult
  FIELD Tokens : *github.com/coldsmirk/vef-framework-go/security.AuthTokens [field_order=1 tag="json:\"tokens,omitempty\""]
  FIELD ChallengeToken : string [field_order=2 tag="json:\"challengeToken,omitempty\""]
  FIELD Challenge : *github.com/coldsmirk/vef-framework-go/security.LoginChallenge [field_order=3 tag="json:\"challenge,omitempty\""]
TYPE MemoryChallengeTokenStore : github.com/coldsmirk/vef-framework-go/security.MemoryChallengeTokenStore
  METHOD Generate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, username string, pending []string, resolved []string) (string, error)
  METHOD Parse : func(ctx context.Context, token string) (*github.com/coldsmirk/vef-framework-go/security.ChallengeState, error)
TYPE MemoryLoginGuard : github.com/coldsmirk/vef-framework-go/security.MemoryLoginGuard
  METHOD Check : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) (github.com/coldsmirk/vef-framework-go/security.LoginDecision, error)
  METHOD RecordFailure : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) (github.com/coldsmirk/vef-framework-go/security.LoginDecision, error)
  METHOD RecordSuccess : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) error
TYPE MemoryNonceStore : github.com/coldsmirk/vef-framework-go/security.MemoryNonceStore
  METHOD StoreIfAbsent : func(ctx context.Context, appID string, nonce string, ttl time.Duration) (bool, error)
TYPE MemorySessionStore : github.com/coldsmirk/vef-framework-go/security.MemorySessionStore
  METHOD Create : func(_ context.Context, tokenHash string, session github.com/coldsmirk/vef-framework-go/security.Session, _ time.Duration) error
  METHOD ListAll : func(context.Context) ([]github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD ListByUser : func(_ context.Context, userID string) ([]github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD Lookup : func(_ context.Context, tokenHash string) (*github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD Renew : func(_ context.Context, tokenHash string, expiresAt time.Time, _ time.Duration) error
  METHOD Revoke : func(_ context.Context, id string) error
  METHOD RevokeUser : func(_ context.Context, userID string) error
FUNC NewAllDataScope : func() github.com/coldsmirk/vef-framework-go/security.DataScope
FUNC NewBlocklistRule : func(entries []string) github.com/coldsmirk/vef-framework-go/security.PasswordRule
FUNC NewCachedRolePermissionsLoader : func(loader github.com/coldsmirk/vef-framework-go/security.RolePermissionsLoader, bus github.com/coldsmirk/vef-framework-go/event.Bus) github.com/coldsmirk/vef-framework-go/security.RolePermissionsLoader
FUNC NewChainValidator : func(validators ...github.com/coldsmirk/vef-framework-go/security.PasswordValidator) github.com/coldsmirk/vef-framework-go/security.PasswordValidator
FUNC NewCharacterClassRule : func(requireUpper bool, requireLower bool, requireDigit bool, requireSymbol bool, minClasses int) github.com/coldsmirk/vef-framework-go/security.PasswordRule
FUNC NewCompositePasswordChangeChecker : func(checkers ...github.com/coldsmirk/vef-framework-go/security.PasswordChangeChecker) github.com/coldsmirk/vef-framework-go/security.PasswordChangeChecker
FUNC NewDeliveredChallengeProvider : func(challengeType string, order int, evaluator github.com/coldsmirk/vef-framework-go/security.OTPEvaluator, store github.com/coldsmirk/vef-framework-go/security.OTPCodeStore, delivery github.com/coldsmirk/vef-framework-go/security.OTPCodeDelivery) *github.com/coldsmirk/vef-framework-go/security.OTPChallengeProvider
FUNC NewDeliveredCodeSender : func(store github.com/coldsmirk/vef-framework-go/security.OTPCodeStore, delivery github.com/coldsmirk/vef-framework-go/security.OTPCodeDelivery) *github.com/coldsmirk/vef-framework-go/security.DeliveredCodeSender
FUNC NewDeliveredCodeVerifier : func(store github.com/coldsmirk/vef-framework-go/security.OTPCodeStore) *github.com/coldsmirk/vef-framework-go/security.DeliveredCodeVerifier
FUNC NewDepartmentSelectionChallengeProvider : func(loader github.com/coldsmirk/vef-framework-go/security.DepartmentLoader, selector github.com/coldsmirk/vef-framework-go/security.DepartmentSelector) *github.com/coldsmirk/vef-framework-go/security.DepartmentSelectionChallengeProvider
FUNC NewDisallowIdentityRule : func() github.com/coldsmirk/vef-framework-go/security.PasswordRule
FUNC NewEmailChallengeProvider : func(evaluator github.com/coldsmirk/vef-framework-go/security.OTPEvaluator, store github.com/coldsmirk/vef-framework-go/security.OTPCodeStore, delivery github.com/coldsmirk/vef-framework-go/security.OTPCodeDelivery) *github.com/coldsmirk/vef-framework-go/security.OTPChallengeProvider
FUNC NewExpiryPasswordChangeChecker : func(loader github.com/coldsmirk/vef-framework-go/security.PasswordMetadataLoader, maxAge time.Duration) *github.com/coldsmirk/vef-framework-go/security.ExpiryPasswordChangeChecker
FUNC NewExternalApp : func(id string, name string, roles ...string) *github.com/coldsmirk/vef-framework-go/security.Principal
FUNC NewHistoryValidator : func(store github.com/coldsmirk/vef-framework-go/security.PasswordHistoryStore, encoder github.com/coldsmirk/vef-framework-go/password.Encoder, depth int) github.com/coldsmirk/vef-framework-go/security.PasswordValidator
FUNC NewIPWhitelistValidator : func(whitelist string) *github.com/coldsmirk/vef-framework-go/security.IPWhitelistValidator
FUNC NewIPWhitelistValidatorFromEntries : func(entries []string) *github.com/coldsmirk/vef-framework-go/security.IPWhitelistValidator
FUNC NewJWT : func(config *github.com/coldsmirk/vef-framework-go/security.JWTConfig) (*github.com/coldsmirk/vef-framework-go/security.JWT, error)
FUNC NewJWTChallengeTokenStore : func(jwt *github.com/coldsmirk/vef-framework-go/security.JWT) github.com/coldsmirk/vef-framework-go/security.ChallengeTokenStore
FUNC NewJWTClaimsAccessor : func(claims github.com/golang-jwt/jwt/v5.MapClaims) *github.com/coldsmirk/vef-framework-go/security.JWTClaimsAccessor
FUNC NewJWTClaimsBuilder : func() *github.com/coldsmirk/vef-framework-go/security.JWTClaimsBuilder
FUNC NewLoginEvent : func(params github.com/coldsmirk/vef-framework-go/security.LoginEventParams) *github.com/coldsmirk/vef-framework-go/security.LoginEvent
FUNC NewMaxLengthRule : func(maxLength int) github.com/coldsmirk/vef-framework-go/security.PasswordRule
FUNC NewMemoryChallengeTokenStore : func() github.com/coldsmirk/vef-framework-go/security.ChallengeTokenStore
FUNC NewMemoryLoginGuard : func(policy github.com/coldsmirk/vef-framework-go/security.LockoutPolicy) github.com/coldsmirk/vef-framework-go/security.LoginGuard
FUNC NewMemoryNonceStore : func() github.com/coldsmirk/vef-framework-go/security.NonceStore
FUNC NewMemorySessionStore : func() github.com/coldsmirk/vef-framework-go/security.SessionStore
FUNC NewMinLengthRule : func(minLength int) github.com/coldsmirk/vef-framework-go/security.PasswordRule
FUNC NewOTPChallengeProvider : func(config github.com/coldsmirk/vef-framework-go/security.OTPChallengeProviderConfig) *github.com/coldsmirk/vef-framework-go/security.OTPChallengeProvider
FUNC NewPasswordChangeChallengeProvider : func(checker github.com/coldsmirk/vef-framework-go/security.PasswordChangeChecker, changer github.com/coldsmirk/vef-framework-go/security.PasswordChanger, validator github.com/coldsmirk/vef-framework-go/security.PasswordValidator) *github.com/coldsmirk/vef-framework-go/security.PasswordChangeChallengeProvider
FUNC NewRedisChallengeTokenStore : func(client *github.com/redis/go-redis/v9.Client) github.com/coldsmirk/vef-framework-go/security.ChallengeTokenStore
FUNC NewRedisLoginGuard : func(client *github.com/redis/go-redis/v9.Client, policy github.com/coldsmirk/vef-framework-go/security.LockoutPolicy) github.com/coldsmirk/vef-framework-go/security.LoginGuard
FUNC NewRedisNonceStore : func(client *github.com/redis/go-redis/v9.Client) github.com/coldsmirk/vef-framework-go/security.NonceStore
FUNC NewRedisSessionStore : func(client *github.com/redis/go-redis/v9.Client) github.com/coldsmirk/vef-framework-go/security.SessionStore
FUNC NewRequestScopedDataPermApplier : func(principal *github.com/coldsmirk/vef-framework-go/security.Principal, dataScope github.com/coldsmirk/vef-framework-go/security.DataScope, logger github.com/coldsmirk/vef-framework-go/logx.Logger) github.com/coldsmirk/vef-framework-go/security.DataPermissionApplier
FUNC NewRuleBasedValidator : func(rules ...github.com/coldsmirk/vef-framework-go/security.PasswordRule) github.com/coldsmirk/vef-framework-go/security.PasswordValidator
FUNC NewSMSChallengeProvider : func(evaluator github.com/coldsmirk/vef-framework-go/security.OTPEvaluator, store github.com/coldsmirk/vef-framework-go/security.OTPCodeStore, delivery github.com/coldsmirk/vef-framework-go/security.OTPCodeDelivery) *github.com/coldsmirk/vef-framework-go/security.OTPChallengeProvider
FUNC NewSelfDataScope : func(createdByColumn string) github.com/coldsmirk/vef-framework-go/security.DataScope
FUNC NewSignature : func(secret string, opts ...github.com/coldsmirk/vef-framework-go/security.SignatureOption) (*github.com/coldsmirk/vef-framework-go/security.Signature, error)
FUNC NewTOTPChallengeProvider : func(loader github.com/coldsmirk/vef-framework-go/security.TOTPSecretLoader, opts ...github.com/coldsmirk/vef-framework-go/security.TOTPOption) *github.com/coldsmirk/vef-framework-go/security.OTPChallengeProvider
FUNC NewTOTPEvaluator : func(loader github.com/coldsmirk/vef-framework-go/security.TOTPSecretLoader, opts ...github.com/coldsmirk/vef-framework-go/security.TOTPOption) *github.com/coldsmirk/vef-framework-go/security.TOTPEvaluator
FUNC NewTOTPVerifier : func(loader github.com/coldsmirk/vef-framework-go/security.TOTPSecretLoader) *github.com/coldsmirk/vef-framework-go/security.TOTPVerifier
FUNC NewUser : func(id string, name string, roles ...string) *github.com/coldsmirk/vef-framework-go/security.Principal
TYPE NonceStore : github.com/coldsmirk/vef-framework-go/security.NonceStore
  METHOD StoreIfAbsent : func(ctx context.Context, appID string, nonce string, ttl time.Duration) (bool, error)
TYPE OTPChallengeData : github.com/coldsmirk/vef-framework-go/security.OTPChallengeData
  FIELD Destination : string [field_order=1 tag="json:\"destination\""]
  FIELD Meta : map[string]any [field_order=2 tag="json:\"meta,omitempty\""]
TYPE OTPChallengeProvider : github.com/coldsmirk/vef-framework-go/security.OTPChallengeProvider
  METHOD Evaluate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.LoginChallenge, error)
  METHOD Order : func() int
  METHOD Resolve : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, response any) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
  METHOD Type : func() string
TYPE OTPChallengeProviderConfig : github.com/coldsmirk/vef-framework-go/security.OTPChallengeProviderConfig
  FIELD ChallengeType : string [field_order=1 tag=""]
  FIELD ChallengeOrder : int [field_order=2 tag=""]
  FIELD Evaluator : github.com/coldsmirk/vef-framework-go/security.OTPEvaluator [field_order=3 tag=""]
  FIELD Sender : github.com/coldsmirk/vef-framework-go/security.OTPCodeSender [field_order=4 tag=""]
  FIELD Verifier : github.com/coldsmirk/vef-framework-go/security.OTPCodeVerifier [field_order=5 tag=""]
TYPE OTPCodeDelivery : github.com/coldsmirk/vef-framework-go/security.OTPCodeDelivery
  METHOD Deliver : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, code string) error
TYPE OTPCodeSender : github.com/coldsmirk/vef-framework-go/security.OTPCodeSender
  METHOD Send : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) error
TYPE OTPCodeStore : github.com/coldsmirk/vef-framework-go/security.OTPCodeStore
  METHOD Generate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (string, error)
  METHOD Verify : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, code string) (bool, error)
TYPE OTPCodeVerifier : github.com/coldsmirk/vef-framework-go/security.OTPCodeVerifier
  METHOD Verify : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, code string) (bool, error)
TYPE OTPEvaluator : github.com/coldsmirk/vef-framework-go/security.OTPEvaluator
  METHOD Evaluate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.OTPChallengeData, error)
TYPE PasswordChangeChallengeData : github.com/coldsmirk/vef-framework-go/security.PasswordChangeChallengeData
  FIELD Reason : string [field_order=1 tag="json:\"reason\""]
  FIELD Meta : map[string]any [field_order=2 tag="json:\"meta,omitempty\""]
TYPE PasswordChangeChallengeProvider : github.com/coldsmirk/vef-framework-go/security.PasswordChangeChallengeProvider
  METHOD Evaluate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.LoginChallenge, error)
  METHOD Order : func() int
  METHOD Resolve : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, response any) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
  METHOD Type : func() string
TYPE PasswordChangeChecker : github.com/coldsmirk/vef-framework-go/security.PasswordChangeChecker
  METHOD Check : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.PasswordChangeChallengeData, error)
CONST PasswordChangeReasonExpired : untyped string = "expired"
CONST PasswordChangeReasonFirstLogin : untyped string = "first_login"
TYPE PasswordChanger : github.com/coldsmirk/vef-framework-go/security.PasswordChanger
  METHOD ChangePassword : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, newPassword string) error
TYPE PasswordDecryptor : github.com/coldsmirk/vef-framework-go/security.PasswordDecryptor
  METHOD Decrypt : func(encryptedPassword string) (string, error)
TYPE PasswordHistoryStore : github.com/coldsmirk/vef-framework-go/security.PasswordHistoryStore
  METHOD Add : func(ctx context.Context, principalID string, encodedPassword string) error
  METHOD Recent : func(ctx context.Context, principalID string, limit int) ([]string, error)
TYPE PasswordMetadataLoader : github.com/coldsmirk/vef-framework-go/security.PasswordMetadataLoader
  METHOD PasswordChangedAt : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (time.Time, error)
TYPE PasswordRule : github.com/coldsmirk/vef-framework-go/security.PasswordRule
  METHOD Check : func(principal *github.com/coldsmirk/vef-framework-go/security.Principal, plaintext string) error
TYPE PasswordValidator : github.com/coldsmirk/vef-framework-go/security.PasswordValidator
  METHOD Validate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, plaintext string) error
TYPE PermissionChecker : github.com/coldsmirk/vef-framework-go/security.PermissionChecker
  METHOD HasPermission : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, permission string) (bool, error)
TYPE Principal : github.com/coldsmirk/vef-framework-go/security.Principal
  FIELD Type : github.com/coldsmirk/vef-framework-go/security.PrincipalType [field_order=1 tag="json:\"type\""]
  FIELD ID : string [field_order=2 tag="json:\"id\""]
  FIELD Name : string [field_order=3 tag="json:\"name\""]
  FIELD Roles : []string [field_order=4 tag="json:\"roles\""]
  FIELD Details : any [field_order=5 tag="json:\"details\""]
  METHOD AttemptUnmarshalDetails : func(details any)
  METHOD UnmarshalJSON : func(data []byte) error
  METHOD WithRoles : func(roles ...string) *github.com/coldsmirk/vef-framework-go/security.Principal
VAR PrincipalAnonymous : *github.com/coldsmirk/vef-framework-go/security.Principal
VAR PrincipalSystem : *github.com/coldsmirk/vef-framework-go/security.Principal
TYPE PrincipalType : github.com/coldsmirk/vef-framework-go/security.PrincipalType
CONST PrincipalTypeExternalApp : github.com/coldsmirk/vef-framework-go/security.PrincipalType = "external_app"
CONST PrincipalTypeSystem : github.com/coldsmirk/vef-framework-go/security.PrincipalType = "system"
CONST PrincipalTypeUser : github.com/coldsmirk/vef-framework-go/security.PrincipalType = "user"
CONST PriorityAll : untyped int = 10000
CONST PriorityCustom : untyped int = 60
CONST PriorityDepartment : untyped int = 20
CONST PriorityDepartmentAndSub : untyped int = 30
CONST PriorityOrganization : untyped int = 40
CONST PriorityOrganizationAndSub : untyped int = 50
CONST PrioritySelf : untyped int = 10
FUNC PublishRolePermissionsChangedEvent : func(ctx context.Context, bus github.com/coldsmirk/vef-framework-go/event.Bus, roles ...string) error
CONST QueryKeyAccessToken : untyped string = "__accessToken"
TYPE RedisChallengeTokenStore : github.com/coldsmirk/vef-framework-go/security.RedisChallengeTokenStore
  METHOD Generate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, username string, pending []string, resolved []string) (string, error)
  METHOD Parse : func(ctx context.Context, token string) (*github.com/coldsmirk/vef-framework-go/security.ChallengeState, error)
TYPE RedisLoginGuard : github.com/coldsmirk/vef-framework-go/security.RedisLoginGuard
  METHOD Check : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) (github.com/coldsmirk/vef-framework-go/security.LoginDecision, error)
  METHOD RecordFailure : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) (github.com/coldsmirk/vef-framework-go/security.LoginDecision, error)
  METHOD RecordSuccess : func(ctx context.Context, attempt github.com/coldsmirk/vef-framework-go/security.LoginAttempt) error
TYPE RedisNonceStore : github.com/coldsmirk/vef-framework-go/security.RedisNonceStore
  METHOD StoreIfAbsent : func(ctx context.Context, appID string, nonce string, ttl time.Duration) (bool, error)
TYPE RedisSessionStore : github.com/coldsmirk/vef-framework-go/security.RedisSessionStore
  METHOD Create : func(ctx context.Context, tokenHash string, session github.com/coldsmirk/vef-framework-go/security.Session, ttl time.Duration) error
  METHOD ListAll : func(ctx context.Context) ([]github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD ListByUser : func(ctx context.Context, userID string) ([]github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD Lookup : func(ctx context.Context, tokenHash string) (*github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD Renew : func(ctx context.Context, tokenHash string, expiresAt time.Time, ttl time.Duration) error
  METHOD Revoke : func(ctx context.Context, id string) error
  METHOD RevokeUser : func(ctx context.Context, userID string) error
TYPE RequestScopedDataPermApplier : github.com/coldsmirk/vef-framework-go/security.RequestScopedDataPermApplier
  METHOD Apply : func(query github.com/coldsmirk/vef-framework-go/orm.SelectQuery) error
TYPE RolePermissionsChangedEvent : github.com/coldsmirk/vef-framework-go/security.RolePermissionsChangedEvent
  FIELD Roles : []string [field_order=1 tag="json:\"roles\""]
  METHOD EventType : func() string
TYPE RolePermissionsLoader : github.com/coldsmirk/vef-framework-go/security.RolePermissionsLoader
  METHOD LoadPermissions : func(ctx context.Context, role string) (map[string]github.com/coldsmirk/vef-framework-go/security.DataScope, error)
TYPE SelfDataScope : github.com/coldsmirk/vef-framework-go/security.SelfDataScope
  METHOD Apply : func(principal *github.com/coldsmirk/vef-framework-go/security.Principal, query github.com/coldsmirk/vef-framework-go/orm.SelectQuery) error
  METHOD Key : func() string
  METHOD Priority : func() int
  METHOD Supports : func(_ *github.com/coldsmirk/vef-framework-go/security.Principal, table *github.com/coldsmirk/vef-framework-go/orm.Table) bool
TYPE Session : github.com/coldsmirk/vef-framework-go/security.Session
  FIELD ID : string [field_order=1 tag=""]
  FIELD UserID : string [field_order=2 tag=""]
  FIELD Principal : *github.com/coldsmirk/vef-framework-go/security.Principal [field_order=3 tag=""]
  FIELD ClientIP : string [field_order=4 tag=""]
  FIELD UserAgent : string [field_order=5 tag=""]
  FIELD CreatedAt : time.Time [field_order=6 tag=""]
  FIELD LastSeenAt : time.Time [field_order=7 tag=""]
  FIELD ExpiresAt : time.Time [field_order=8 tag=""]
CONST SessionExceedEvictOldest : github.com/coldsmirk/vef-framework-go/security.SessionExceedPolicy = "evict_oldest"
TYPE SessionExceedPolicy : github.com/coldsmirk/vef-framework-go/security.SessionExceedPolicy
CONST SessionExceedReject : github.com/coldsmirk/vef-framework-go/security.SessionExceedPolicy = "reject"
TYPE SessionInspector : github.com/coldsmirk/vef-framework-go/security.SessionInspector
  METHOD ListAll : func(ctx context.Context) ([]github.com/coldsmirk/vef-framework-go/security.Session, error)
TYPE SessionMeta : github.com/coldsmirk/vef-framework-go/security.SessionMeta
  FIELD ClientIP : string [field_order=1 tag=""]
  FIELD UserAgent : string [field_order=2 tag=""]
TYPE SessionPolicy : github.com/coldsmirk/vef-framework-go/security.SessionPolicy
  FIELD MaxConcurrent : int [field_order=1 tag=""]
  FIELD OnExceed : github.com/coldsmirk/vef-framework-go/security.SessionExceedPolicy [field_order=2 tag=""]
  FIELD IdleTTL : time.Duration [field_order=3 tag=""]
  FIELD MaxLifetime : time.Duration [field_order=4 tag=""]
  FIELD Sliding : bool [field_order=5 tag=""]
TYPE SessionStore : github.com/coldsmirk/vef-framework-go/security.SessionStore
  METHOD Create : func(ctx context.Context, tokenHash string, session github.com/coldsmirk/vef-framework-go/security.Session, ttl time.Duration) error
  METHOD ListByUser : func(ctx context.Context, userID string) ([]github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD Lookup : func(ctx context.Context, tokenHash string) (*github.com/coldsmirk/vef-framework-go/security.Session, error)
  METHOD Renew : func(ctx context.Context, tokenHash string, expiresAt time.Time, ttl time.Duration) error
  METHOD Revoke : func(ctx context.Context, id string) error
  METHOD RevokeUser : func(ctx context.Context, userID string) error
FUNC SetExternalAppDetailsType : func[T any]()
FUNC SetUserDetailsType : func[T any]()
TYPE Signature : github.com/coldsmirk/vef-framework-go/security.Signature
  METHOD Sign : func(appID string, method string, path string) (*github.com/coldsmirk/vef-framework-go/security.SignatureResult, error)
  METHOD Verify : func(ctx context.Context, appID string, method string, path string, timestamp int64, nonce string, signature string) error
  METHOD VerifyWithSecret : func(ctx context.Context, secret string, appID string, method string, path string, timestamp int64, nonce string, signature string) error
CONST SignatureAlgHmacSHA256 : github.com/coldsmirk/vef-framework-go/security.SignatureAlgorithm = "HMAC-SHA256"
CONST SignatureAlgHmacSHA512 : github.com/coldsmirk/vef-framework-go/security.SignatureAlgorithm = "HMAC-SHA512"
CONST SignatureAlgHmacSM3 : github.com/coldsmirk/vef-framework-go/security.SignatureAlgorithm = "HMAC-SM3"
TYPE SignatureAlgorithm : github.com/coldsmirk/vef-framework-go/security.SignatureAlgorithm
TYPE SignatureCredentials : github.com/coldsmirk/vef-framework-go/security.SignatureCredentials
  FIELD Timestamp : int64 [field_order=1 tag=""]
  FIELD Nonce : string [field_order=2 tag=""]
  FIELD Signature : string [field_order=3 tag=""]
TYPE SignatureOption : github.com/coldsmirk/vef-framework-go/security.SignatureOption
TYPE SignatureResult : github.com/coldsmirk/vef-framework-go/security.SignatureResult
  FIELD AppID : string [field_order=1 tag=""]
  FIELD Timestamp : int64 [field_order=2 tag=""]
  FIELD Nonce : string [field_order=3 tag=""]
  FIELD Signature : string [field_order=4 tag=""]
FUNC SubscribeLoginEvent : func(bus github.com/coldsmirk/vef-framework-go/event.Bus, handler func(context.Context, *github.com/coldsmirk/vef-framework-go/security.LoginEvent) error, opts ...github.com/coldsmirk/vef-framework-go/event.SubscribeOption) (github.com/coldsmirk/vef-framework-go/event.Unsubscribe, error)
CONST TOTPDefaultDestination : untyped string = "Authenticator App"
TYPE TOTPEvaluator : github.com/coldsmirk/vef-framework-go/security.TOTPEvaluator
  METHOD Evaluate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (*github.com/coldsmirk/vef-framework-go/security.OTPChallengeData, error)
TYPE TOTPOption : github.com/coldsmirk/vef-framework-go/security.TOTPOption
TYPE TOTPSecretLoader : github.com/coldsmirk/vef-framework-go/security.TOTPSecretLoader
  METHOD LoadSecret : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal) (string, error)
TYPE TOTPVerifier : github.com/coldsmirk/vef-framework-go/security.TOTPVerifier
  METHOD Verify : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, code string) (bool, error)
TYPE TokenGenerator : github.com/coldsmirk/vef-framework-go/security.TokenGenerator
  METHOD Generate : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, meta github.com/coldsmirk/vef-framework-go/security.SessionMeta) (*github.com/coldsmirk/vef-framework-go/security.AuthTokens, error)
CONST TokenTypeAccess : untyped string = "access"
CONST TokenTypeChallenge : untyped string = "challenge"
CONST TokenTypeRefresh : untyped string = "refresh"
TYPE UserInfo : github.com/coldsmirk/vef-framework-go/security.UserInfo
  FIELD ID : string [field_order=1 tag="json:\"id\""]
  FIELD Name : string [field_order=2 tag="json:\"name\""]
  FIELD Gender : github.com/coldsmirk/vef-framework-go/security.Gender [field_order=3 tag="json:\"gender\""]
  FIELD Avatar : *string [field_order=4 tag="json:\"avatar\""]
  FIELD PermissionTokens : []string [field_order=5 tag="json:\"permissionTokens\""]
  FIELD Menus : []github.com/coldsmirk/vef-framework-go/security.UserMenu [field_order=6 tag="json:\"menus\""]
  FIELD Details : any [field_order=7 tag="json:\"details,omitempty\""]
TYPE UserInfoLoader : github.com/coldsmirk/vef-framework-go/security.UserInfoLoader
  METHOD LoadUserInfo : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, params map[string]any) (*github.com/coldsmirk/vef-framework-go/security.UserInfo, error)
TYPE UserLoader : github.com/coldsmirk/vef-framework-go/security.UserLoader
  METHOD LoadByID : func(ctx context.Context, id string) (*github.com/coldsmirk/vef-framework-go/security.Principal, error)
  METHOD LoadByUsername : func(ctx context.Context, username string) (*github.com/coldsmirk/vef-framework-go/security.Principal, string, error)
TYPE UserMenu : github.com/coldsmirk/vef-framework-go/security.UserMenu
  FIELD Type : github.com/coldsmirk/vef-framework-go/security.UserMenuType [field_order=1 tag="json:\"type\""]
  FIELD Path : string [field_order=2 tag="json:\"path\""]
  FIELD Name : string [field_order=3 tag="json:\"name\""]
  FIELD Icon : *string [field_order=4 tag="json:\"icon\""]
  FIELD Meta : map[string]any [field_order=5 tag="json:\"meta,omitempty\""]
  FIELD Children : []github.com/coldsmirk/vef-framework-go/security.UserMenu [field_order=6 tag="json:\"children,omitempty\""]
TYPE UserMenuType : github.com/coldsmirk/vef-framework-go/security.UserMenuType
CONST UserMenuTypeDashboard : github.com/coldsmirk/vef-framework-go/security.UserMenuType = "dashboard"
CONST UserMenuTypeDirectory : github.com/coldsmirk/vef-framework-go/security.UserMenuType = "directory"
CONST UserMenuTypeMenu : github.com/coldsmirk/vef-framework-go/security.UserMenuType = "menu"
CONST UserMenuTypeReport : github.com/coldsmirk/vef-framework-go/security.UserMenuType = "report"
CONST UserMenuTypeView : github.com/coldsmirk/vef-framework-go/security.UserMenuType = "view"
FUNC WithAlgorithm : func(algorithm github.com/coldsmirk/vef-framework-go/security.SignatureAlgorithm) github.com/coldsmirk/vef-framework-go/security.SignatureOption
FUNC WithNonceStore : func(store github.com/coldsmirk/vef-framework-go/security.NonceStore) github.com/coldsmirk/vef-framework-go/security.SignatureOption
FUNC WithTOTPDestination : func(destination string) github.com/coldsmirk/vef-framework-go/security.TOTPOption
FUNC WithTimestampTolerance : func(tolerance time.Duration) github.com/coldsmirk/vef-framework-go/security.SignatureOption

## github.com/coldsmirk/vef-framework-go/sequence
TYPE DBStore : github.com/coldsmirk/vef-framework-go/sequence.DBStore
  METHOD Init : func(ctx context.Context) error
  METHOD Reserve : func(ctx context.Context, key string, count int, now github.com/coldsmirk/vef-framework-go/timex.DateTime) (*github.com/coldsmirk/vef-framework-go/sequence.Rule, int, error)
CONST DBStoreTableName : untyped string = "sys_sequence_rule"
VAR ErrInvalidCount : error
VAR ErrInvalidStep : error
VAR ErrRuleNotFound : error
VAR ErrSequenceOverflow : error
FUNC FormatDate : func(dt github.com/coldsmirk/vef-framework-go/timex.DateTime, format string) string
TYPE Generator : github.com/coldsmirk/vef-framework-go/sequence.Generator
  METHOD Generate : func(ctx context.Context, key string) (string, error)
  METHOD GenerateN : func(ctx context.Context, key string, count int) ([]string, error)
TYPE MemoryStore : github.com/coldsmirk/vef-framework-go/sequence.MemoryStore
  METHOD Register : func(rules ...*github.com/coldsmirk/vef-framework-go/sequence.Rule)
  METHOD Reserve : func(_ context.Context, key string, count int, now github.com/coldsmirk/vef-framework-go/timex.DateTime) (*github.com/coldsmirk/vef-framework-go/sequence.Rule, int, error)
FUNC NewDBStore : func(db github.com/coldsmirk/vef-framework-go/orm.DB) *github.com/coldsmirk/vef-framework-go/sequence.DBStore
FUNC NewMemoryStore : func() *github.com/coldsmirk/vef-framework-go/sequence.MemoryStore
FUNC NewRedisStore : func(client *github.com/redis/go-redis/v9.Client) *github.com/coldsmirk/vef-framework-go/sequence.RedisStore
CONST OverflowError : github.com/coldsmirk/vef-framework-go/sequence.OverflowStrategy = "error"
CONST OverflowExtend : github.com/coldsmirk/vef-framework-go/sequence.OverflowStrategy = "extend"
CONST OverflowReset : github.com/coldsmirk/vef-framework-go/sequence.OverflowStrategy = "reset"
TYPE OverflowStrategy : github.com/coldsmirk/vef-framework-go/sequence.OverflowStrategy
TYPE RedisStore : github.com/coldsmirk/vef-framework-go/sequence.RedisStore
  METHOD RegisterRule : func(ctx context.Context, rule *github.com/coldsmirk/vef-framework-go/sequence.Rule) error
  METHOD Reserve : func(ctx context.Context, key string, count int, now github.com/coldsmirk/vef-framework-go/timex.DateTime) (*github.com/coldsmirk/vef-framework-go/sequence.Rule, int, error)
TYPE ResetCycle : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle
CONST ResetDaily : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle = "D"
CONST ResetMonthly : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle = "M"
CONST ResetNone : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle = "N"
CONST ResetQuarterly : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle = "Q"
CONST ResetWeekly : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle = "W"
CONST ResetYearly : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle = "Y"
TYPE Rule : github.com/coldsmirk/vef-framework-go/sequence.Rule
  FIELD Key : string [field_order=1 tag=""]
  FIELD Name : string [field_order=2 tag=""]
  FIELD Prefix : string [field_order=3 tag=""]
  FIELD Suffix : string [field_order=4 tag=""]
  FIELD DateFormat : string [field_order=5 tag=""]
  FIELD SeqLength : int [field_order=6 tag=""]
  FIELD SeqStep : int [field_order=7 tag=""]
  FIELD StartValue : int [field_order=8 tag=""]
  FIELD MaxValue : int [field_order=9 tag=""]
  FIELD OverflowStrategy : github.com/coldsmirk/vef-framework-go/sequence.OverflowStrategy [field_order=10 tag=""]
  FIELD ResetCycle : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle [field_order=11 tag=""]
  FIELD CurrentValue : int [field_order=12 tag=""]
  FIELD LastResetAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=13 tag=""]
  FIELD IsActive : bool [field_order=14 tag=""]
  METHOD Clone : func() *github.com/coldsmirk/vef-framework-go/sequence.Rule
TYPE RuleModel : github.com/coldsmirk/vef-framework-go/sequence.RuleModel
  FIELD BaseModel : github.com/coldsmirk/vef-framework-go/orm.BaseModel [field_order=1 tag="bun:\"table:sys_sequence_rule,alias:ssr\""]
  FIELD FullAuditedModel : github.com/coldsmirk/vef-framework-go/orm.FullAuditedModel [field_order=2 tag=""]
  FIELD Key : string [field_order=3 tag="bun:\"key,notnull,unique\""]
  FIELD Name : string [field_order=4 tag="bun:\"name,notnull\""]
  FIELD Prefix : *string [field_order=5 tag="bun:\"prefix,type:varchar(32)\""]
  FIELD Suffix : *string [field_order=6 tag="bun:\"suffix,type:varchar(32)\""]
  FIELD DateFormat : *string [field_order=7 tag="bun:\"date_format,type:varchar(32)\""]
  FIELD SeqLength : int16 [field_order=8 tag="bun:\"seq_length,notnull,default:4\""]
  FIELD SeqStep : int16 [field_order=9 tag="bun:\"seq_step,notnull,default:1\""]
  FIELD StartValue : int [field_order=10 tag="bun:\"start_value,notnull,default:0\""]
  FIELD MaxValue : int [field_order=11 tag="bun:\"max_value,notnull,default:0\""]
  FIELD OverflowStrategy : github.com/coldsmirk/vef-framework-go/sequence.OverflowStrategy [field_order=12 tag="bun:\"overflow_strategy,notnull,default:'error'\""]
  FIELD ResetCycle : github.com/coldsmirk/vef-framework-go/sequence.ResetCycle [field_order=13 tag="bun:\"reset_cycle,notnull,default:'N'\""]
  FIELD CurrentValue : int [field_order=14 tag="bun:\"current_value,notnull,default:0\""]
  FIELD LastResetAt : *github.com/coldsmirk/vef-framework-go/timex.DateTime [field_order=15 tag="bun:\"last_reset_at,type:timestamp\""]
  FIELD IsActive : bool [field_order=16 tag="bun:\"is_active,notnull,default:true\""]
  FIELD Remark : *string [field_order=17 tag="bun:\"remark,type:varchar(256)\""]
TYPE Store : github.com/coldsmirk/vef-framework-go/sequence.Store
  METHOD Reserve : func(ctx context.Context, key string, count int, now github.com/coldsmirk/vef-framework-go/timex.DateTime) (rule *github.com/coldsmirk/vef-framework-go/sequence.Rule, newValue int, err error)

## github.com/coldsmirk/vef-framework-go/sortx
VAR ErrInvalidOrderDirection : error
CONST NullsDefault : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder = 0
CONST NullsFirst : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder = 1
CONST NullsLast : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder = 2
TYPE NullsOrder : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder
  METHOD String : func() string
CONST OrderAsc : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection = 0
CONST OrderDesc : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection = 1
TYPE OrderDirection : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection
  METHOD MarshalJSON : func() ([]byte, error)
  METHOD MarshalText : func() ([]byte, error)
  METHOD String : func() string
  METHOD UnmarshalJSON : func(data []byte) error
  METHOD UnmarshalText : func(text []byte) error
TYPE OrderSpec : github.com/coldsmirk/vef-framework-go/sortx.OrderSpec
  FIELD Column : string [field_order=1 tag=""]
  FIELD Direction : github.com/coldsmirk/vef-framework-go/sortx.OrderDirection [field_order=2 tag=""]
  FIELD NullsOrder : github.com/coldsmirk/vef-framework-go/sortx.NullsOrder [field_order=3 tag=""]
  METHOD IsValid : func() bool

## github.com/coldsmirk/vef-framework-go/storage
TYPE AbortMultipartOptions : github.com/coldsmirk/vef-framework-go/storage.AbortMultipartOptions
  FIELD Key : string [field_order=1 tag=""]
  FIELD UploadID : string [field_order=2 tag=""]
FUNC CanonicalizeMetadataKeys : func(m map[string]string) map[string]string
TYPE ClaimConsumer : github.com/coldsmirk/vef-framework-go/storage.ClaimConsumer
  METHOD Consume : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, principal *github.com/coldsmirk/vef-framework-go/security.Principal, keys []string) error
TYPE CompleteMultipartOptions : github.com/coldsmirk/vef-framework-go/storage.CompleteMultipartOptions
  FIELD Key : string [field_order=1 tag=""]
  FIELD UploadID : string [field_order=2 tag=""]
  FIELD Parts : []github.com/coldsmirk/vef-framework-go/storage.CompletedPart [field_order=3 tag=""]
TYPE CompletedPart : github.com/coldsmirk/vef-framework-go/storage.CompletedPart
  FIELD PartNumber : int [field_order=1 tag=""]
  FIELD ETag : string [field_order=2 tag=""]
TYPE CopyObjectOptions : github.com/coldsmirk/vef-framework-go/storage.CopyObjectOptions
  FIELD SourceKey : string [field_order=1 tag=""]
  FIELD DestKey : string [field_order=2 tag=""]
TYPE DefaultFileACL : github.com/coldsmirk/vef-framework-go/storage.DefaultFileACL
  METHOD CanRead : func(_ context.Context, _ *github.com/coldsmirk/vef-framework-go/security.Principal, key string) (bool, error)
CONST DefaultProxyPrefix : untyped string = "/storage/files/"
TYPE DeleteDeadLetterEvent : github.com/coldsmirk/vef-framework-go/storage.DeleteDeadLetterEvent
  FIELD PendingDeleteID : string [field_order=1 tag="json:\"pendingDeleteId\""]
  FIELD FileKey : string [field_order=2 tag="json:\"fileKey\""]
  FIELD Reason : github.com/coldsmirk/vef-framework-go/storage.DeleteReason [field_order=3 tag="json:\"reason\""]
  FIELD Attempts : int [field_order=4 tag="json:\"attempts\""]
  FIELD LastError : string [field_order=5 tag="json:\"lastError,omitempty\""]
  METHOD EventType : func() string
TYPE DeleteEnqueuer : github.com/coldsmirk/vef-framework-go/storage.DeleteEnqueuer
  METHOD Enqueue : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, keys []string, reason github.com/coldsmirk/vef-framework-go/storage.DeleteReason) error
TYPE DeleteObjectOptions : github.com/coldsmirk/vef-framework-go/storage.DeleteObjectOptions
  FIELD Key : string [field_order=1 tag=""]
TYPE DeleteObjectsOptions : github.com/coldsmirk/vef-framework-go/storage.DeleteObjectsOptions
  FIELD Keys : []string [field_order=1 tag=""]
TYPE DeleteReason : github.com/coldsmirk/vef-framework-go/storage.DeleteReason
CONST DeleteReasonClaimExpired : github.com/coldsmirk/vef-framework-go/storage.DeleteReason = "claim_expired"
CONST DeleteReasonDeleted : github.com/coldsmirk/vef-framework-go/storage.DeleteReason = "deleted"
CONST DeleteReasonReplaced : github.com/coldsmirk/vef-framework-go/storage.DeleteReason = "replaced"
VAR ErrAbortFailed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrAccessDenied : error
VAR ErrBucketNotFound : error
VAR ErrClaimExpired : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrClaimNotFound : error
VAR ErrClaimNotMultipart : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrClaimNotPending : github.com/coldsmirk/vef-framework-go/result.Error
CONST ErrCodeAbortFailed : untyped int = 2219
CONST ErrCodeClaimExpired : untyped int = 2204
CONST ErrCodeClaimNotMultipart : untyped int = 2212
CONST ErrCodeClaimNotPending : untyped int = 2203
CONST ErrCodeFailedToGetFile : untyped int = 2202
CONST ErrCodeFileNotFound : untyped int = 2201
CONST ErrCodeInvalidFileKey : untyped int = 2200
CONST ErrCodeMultipartNotSupported : untyped int = 2206
CONST ErrCodePublicUploadsNotAllowed : untyped int = 2207
CONST ErrCodeTooManyPendingUploads : untyped int = 2209
CONST ErrCodeUploadObjectNotFound : untyped int = 2217
CONST ErrCodeUploadPartNumberOutOfRange : untyped int = 2213
CONST ErrCodeUploadPartTooLarge : untyped int = 2214
CONST ErrCodeUploadPartTooSmall : untyped int = 2215
CONST ErrCodeUploadPartsIncomplete : untyped int = 2216
CONST ErrCodeUploadRequiresFile : untyped int = 2211
CONST ErrCodeUploadRequiresMultipart : untyped int = 2210
CONST ErrCodeUploadSizeExceedsLimit : untyped int = 2205
CONST ErrCodeUploadSizeMismatch : untyped int = 2218
CONST ErrCodeUploadTooManyParts : untyped int = 2208
VAR ErrFailedToGetFile : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrFileNotFound : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrInvalidBucketName : error
VAR ErrInvalidFileKey : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrMultipartNotSupported : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrObjectNotFound : error
VAR ErrPartETagMismatch : error
VAR ErrPartNumberOutOfRange : error
VAR ErrPartTooSmall : error
VAR ErrPublicUploadsNotAllowed : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrTooManyPendingUploads : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadObjectNotFound : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadPartNumberOutOfRange : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadPartTooLarge : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadPartTooSmall : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadPartsIncomplete : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadRequiresFile : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadRequiresMultipart : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadSessionNotFound : error
VAR ErrUploadSizeExceedsLimit : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadSizeMismatch : github.com/coldsmirk/vef-framework-go/result.Error
VAR ErrUploadTooManyParts : github.com/coldsmirk/vef-framework-go/result.Error
CONST EventTypeDeleteDeadLetter : untyped string = "vef.storage.delete.dead_letter"
CONST EventTypeFileClaimed : untyped string = "vef.storage.file.claimed"
CONST EventTypeFileDeleted : untyped string = "vef.storage.file.deleted"
TYPE FileACL : github.com/coldsmirk/vef-framework-go/storage.FileACL
  METHOD CanRead : func(ctx context.Context, principal *github.com/coldsmirk/vef-framework-go/security.Principal, key string) (bool, error)
TYPE FileClaimedEvent : github.com/coldsmirk/vef-framework-go/storage.FileClaimedEvent
  FIELD FileKey : string [field_order=1 tag="json:\"fileKey\""]
  METHOD EventType : func() string
TYPE FileDeletedEvent : github.com/coldsmirk/vef-framework-go/storage.FileDeletedEvent
  FIELD FileKey : string [field_order=1 tag="json:\"fileKey\""]
  FIELD Reason : github.com/coldsmirk/vef-framework-go/storage.DeleteReason [field_order=2 tag="json:\"reason\""]
  METHOD EventType : func() string
TYPE FileRef : github.com/coldsmirk/vef-framework-go/storage.FileRef
  FIELD Key : string [field_order=1 tag=""]
  FIELD MetaType : github.com/coldsmirk/vef-framework-go/storage.MetaType [field_order=2 tag=""]
  FIELD Attrs : map[string]string [field_order=3 tag=""]
TYPE Files : github.com/coldsmirk/vef-framework-go/storage.Files
  METHOD OnCreate : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, principal *github.com/coldsmirk/vef-framework-go/security.Principal, model any) error
  METHOD OnDelete : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, model any) error
  METHOD OnUpdate : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, principal *github.com/coldsmirk/vef-framework-go/security.Principal, oldModel any, newModel any) error
TYPE FilesFor : github.com/coldsmirk/vef-framework-go/storage.FilesFor[T any]
  METHOD OnCreate : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, principal *github.com/coldsmirk/vef-framework-go/security.Principal, model *T) error
  METHOD OnDelete : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, model *T) error
  METHOD OnUpdate : func(ctx context.Context, tx github.com/coldsmirk/vef-framework-go/orm.DB, principal *github.com/coldsmirk/vef-framework-go/security.Principal, oldModel *T, newModel *T) error
TYPE GetObjectOptions : github.com/coldsmirk/vef-framework-go/storage.GetObjectOptions
  FIELD Key : string [field_order=1 tag=""]
TYPE IdentityURLKeyMapper : github.com/coldsmirk/vef-framework-go/storage.IdentityURLKeyMapper
  METHOD KeyToURL : func(key string) string
  METHOD URLToKey : func(rawURL string) (string, bool)
TYPE InitMultipartOptions : github.com/coldsmirk/vef-framework-go/storage.InitMultipartOptions
  FIELD Key : string [field_order=1 tag=""]
  FIELD ContentType : string [field_order=2 tag=""]
  FIELD Metadata : map[string]string [field_order=3 tag=""]
TYPE MetaType : github.com/coldsmirk/vef-framework-go/storage.MetaType
CONST MetaTypeMarkdown : github.com/coldsmirk/vef-framework-go/storage.MetaType = "markdown"
CONST MetaTypeRichText : github.com/coldsmirk/vef-framework-go/storage.MetaType = "rich_text"
CONST MetaTypeUploadedFile : github.com/coldsmirk/vef-framework-go/storage.MetaType = "uploaded_file"
TYPE Multipart : github.com/coldsmirk/vef-framework-go/storage.Multipart
  METHOD AbortMultipart : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.AbortMultipartOptions) error
  METHOD CompleteMultipart : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.CompleteMultipartOptions) (*github.com/coldsmirk/vef-framework-go/storage.ObjectInfo, error)
  METHOD InitMultipart : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.InitMultipartOptions) (*github.com/coldsmirk/vef-framework-go/storage.MultipartSession, error)
  METHOD MaxPartCount : func() int
  METHOD PartSize : func() int64
  METHOD PutPart : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.PutPartOptions) (*github.com/coldsmirk/vef-framework-go/storage.PartInfo, error)
FUNC MultipartFor : func(svc github.com/coldsmirk/vef-framework-go/storage.Service) github.com/coldsmirk/vef-framework-go/storage.Multipart
TYPE MultipartSession : github.com/coldsmirk/vef-framework-go/storage.MultipartSession
  FIELD Key : string [field_order=1 tag=""]
  FIELD UploadID : string [field_order=2 tag=""]
FUNC NewDeleteDeadLetterEvent : func(id string, key string, reason github.com/coldsmirk/vef-framework-go/storage.DeleteReason, attempts int, lastErr string) *github.com/coldsmirk/vef-framework-go/storage.DeleteDeadLetterEvent
FUNC NewFileClaimedEvent : func(key string) *github.com/coldsmirk/vef-framework-go/storage.FileClaimedEvent
FUNC NewFileDeletedEvent : func(key string, reason github.com/coldsmirk/vef-framework-go/storage.DeleteReason) *github.com/coldsmirk/vef-framework-go/storage.FileDeletedEvent
FUNC NewFiles : func(cc github.com/coldsmirk/vef-framework-go/storage.ClaimConsumer, de github.com/coldsmirk/vef-framework-go/storage.DeleteEnqueuer, bus github.com/coldsmirk/vef-framework-go/event.Bus, urlMapper github.com/coldsmirk/vef-framework-go/storage.URLKeyMapper) github.com/coldsmirk/vef-framework-go/storage.Files
FUNC NewFilesFor : func[T any](files github.com/coldsmirk/vef-framework-go/storage.Files) github.com/coldsmirk/vef-framework-go/storage.FilesFor[T]
TYPE ObjectInfo : github.com/coldsmirk/vef-framework-go/storage.ObjectInfo
  FIELD Bucket : string [field_order=1 tag="json:\"bucket\""]
  FIELD Key : string [field_order=2 tag="json:\"key\""]
  FIELD ETag : string [field_order=3 tag="json:\"eTag\""]
  FIELD Size : int64 [field_order=4 tag="json:\"size\""]
  FIELD ContentType : string [field_order=5 tag="json:\"contentType\""]
  FIELD LastModified : time.Time [field_order=6 tag="json:\"lastModified\""]
  FIELD Metadata : map[string]string [field_order=7 tag="json:\"metadata,omitempty\""]
TYPE PartInfo : github.com/coldsmirk/vef-framework-go/storage.PartInfo
  FIELD PartNumber : int [field_order=1 tag=""]
  FIELD ETag : string [field_order=2 tag=""]
  FIELD Size : int64 [field_order=3 tag=""]
CONST PrivatePrefix : untyped string = "priv/"
TYPE ProxyURLKeyMapper : github.com/coldsmirk/vef-framework-go/storage.ProxyURLKeyMapper
  FIELD Prefix : string [field_order=1 tag=""]
  METHOD KeyToURL : func(key string) string
  METHOD URLToKey : func(rawURL string) (string, bool)
CONST PublicPrefix : untyped string = "pub/"
TYPE PutObjectOptions : github.com/coldsmirk/vef-framework-go/storage.PutObjectOptions
  FIELD Key : string [field_order=1 tag=""]
  FIELD Reader : io.Reader [field_order=2 tag=""]
  FIELD Size : int64 [field_order=3 tag=""]
  FIELD ContentType : string [field_order=4 tag=""]
  FIELD Metadata : map[string]string [field_order=5 tag=""]
TYPE PutPartOptions : github.com/coldsmirk/vef-framework-go/storage.PutPartOptions
  FIELD Key : string [field_order=1 tag=""]
  FIELD UploadID : string [field_order=2 tag=""]
  FIELD PartNumber : int [field_order=3 tag=""]
  FIELD Reader : io.Reader [field_order=4 tag=""]
  FIELD Size : int64 [field_order=5 tag=""]
FUNC ReplaceHtmlURLs : func(content string, replacements map[string]string) string
FUNC ReplaceMarkdownURLs : func(content string, replacements map[string]string) string
TYPE Service : github.com/coldsmirk/vef-framework-go/storage.Service
  METHOD CopyObject : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.CopyObjectOptions) (*github.com/coldsmirk/vef-framework-go/storage.ObjectInfo, error)
  METHOD DeleteObject : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.DeleteObjectOptions) error
  METHOD DeleteObjects : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.DeleteObjectsOptions) error
  METHOD GetObject : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.GetObjectOptions) (io.ReadCloser, *github.com/coldsmirk/vef-framework-go/storage.ObjectInfo, error)
  METHOD PutObject : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.PutObjectOptions) (*github.com/coldsmirk/vef-framework-go/storage.ObjectInfo, error)
  METHOD StatObject : func(ctx context.Context, opts github.com/coldsmirk/vef-framework-go/storage.StatObjectOptions) (*github.com/coldsmirk/vef-framework-go/storage.ObjectInfo, error)
TYPE StatObjectOptions : github.com/coldsmirk/vef-framework-go/storage.StatObjectOptions
  FIELD Key : string [field_order=1 tag=""]
TYPE URLKeyMapper : github.com/coldsmirk/vef-framework-go/storage.URLKeyMapper
  METHOD KeyToURL : func(key string) string
  METHOD URLToKey : func(url string) (key string, ok bool)

## github.com/coldsmirk/vef-framework-go/strx
CONST BareAsKey : github.com/coldsmirk/vef-framework-go/strx.BareValueMode = 1
CONST BareAsValue : github.com/coldsmirk/vef-framework-go/strx.BareValueMode = 0
TYPE BareValueMode : github.com/coldsmirk/vef-framework-go/strx.BareValueMode
CONST DefaultKey : untyped string = "__default"
TYPE ParseOption : github.com/coldsmirk/vef-framework-go/strx.ParseOption
FUNC ParseTag : func(input string, opts ...github.com/coldsmirk/vef-framework-go/strx.ParseOption) map[string]string
FUNC WithBareValueMode : func(mode github.com/coldsmirk/vef-framework-go/strx.BareValueMode) github.com/coldsmirk/vef-framework-go/strx.ParseOption
FUNC WithPairDelimiter : func(delimiter rune) github.com/coldsmirk/vef-framework-go/strx.ParseOption
FUNC WithPairDelimiterFunc : func(fn func(rune) bool) github.com/coldsmirk/vef-framework-go/strx.ParseOption
FUNC WithSpacePairDelimiter : func() github.com/coldsmirk/vef-framework-go/strx.ParseOption
FUNC WithValueDelimiter : func(delimiter rune) github.com/coldsmirk/vef-framework-go/strx.ParseOption

## github.com/coldsmirk/vef-framework-go/tabular
CONST AttrDefault : untyped string = "default"
CONST AttrDive : untyped string = "dive"
CONST AttrFormat : untyped string = "format"
CONST AttrFormatter : untyped string = "formatter"
CONST AttrName : untyped string = "name"
CONST AttrOrder : untyped string = "order"
CONST AttrParser : untyped string = "parser"
CONST AttrWidth : untyped string = "width"
FUNC BuildHeaderMapping : func(headerRow []string, schema *github.com/coldsmirk/vef-framework-go/tabular.Schema, opts github.com/coldsmirk/vef-framework-go/tabular.MappingOptions) (map[int]int, error)
TYPE CellValidator : github.com/coldsmirk/vef-framework-go/tabular.CellValidator
TYPE Column : github.com/coldsmirk/vef-framework-go/tabular.Column
  FIELD Key : string [field_order=1 tag=""]
  FIELD Name : string [field_order=2 tag=""]
  FIELD Type : reflect.Type [field_order=3 tag=""]
  FIELD Order : int [field_order=4 tag=""]
  FIELD Width : float64 [field_order=5 tag=""]
  FIELD Default : string [field_order=6 tag=""]
  FIELD Format : string [field_order=7 tag=""]
  FIELD Formatter : string [field_order=8 tag=""]
  FIELD Parser : string [field_order=9 tag=""]
  FIELD FormatterFn : github.com/coldsmirk/vef-framework-go/tabular.Formatter [field_order=10 tag=""]
  FIELD ParserFn : github.com/coldsmirk/vef-framework-go/tabular.ValueParser [field_order=11 tag=""]
  FIELD Required : bool [field_order=12 tag=""]
  FIELD Validators : []github.com/coldsmirk/vef-framework-go/tabular.CellValidator [field_order=13 tag=""]
  FIELD Index : []int [field_order=14 tag=""]
TYPE ColumnMapping : github.com/coldsmirk/vef-framework-go/tabular.ColumnMapping
TYPE ColumnSpec : github.com/coldsmirk/vef-framework-go/tabular.ColumnSpec
  FIELD Key : string [field_order=1 tag=""]
  FIELD Name : string [field_order=2 tag=""]
  FIELD Type : reflect.Type [field_order=3 tag=""]
  FIELD Order : int [field_order=4 tag=""]
  FIELD Width : float64 [field_order=5 tag=""]
  FIELD Default : string [field_order=6 tag=""]
  FIELD Format : string [field_order=7 tag=""]
  FIELD Formatter : string [field_order=8 tag=""]
  FIELD Parser : string [field_order=9 tag=""]
  FIELD FormatterFn : github.com/coldsmirk/vef-framework-go/tabular.Formatter [field_order=10 tag=""]
  FIELD ParserFn : github.com/coldsmirk/vef-framework-go/tabular.ValueParser [field_order=11 tag=""]
  FIELD Required : bool [field_order=12 tag=""]
  FIELD Validators : []github.com/coldsmirk/vef-framework-go/tabular.CellValidator [field_order=13 tag=""]
FUNC DefaultPositionalMapping : func(schema *github.com/coldsmirk/vef-framework-go/tabular.Schema) map[int]int
VAR ErrDataMustBeSlice : error
VAR ErrDuplicateColumnKey : error
VAR ErrDuplicateHeaderName : error
VAR ErrMissingColumnKey : error
VAR ErrMissingColumnType : error
VAR ErrNoDataRowsFound : error
VAR ErrRequiredMissing : error
VAR ErrSchemaMismatch : error
VAR ErrTypedRowMismatch : error
VAR ErrUnknownColumn : error
VAR ErrUnsetField : error
VAR ErrUnsupportedType : error
TYPE ExportError : github.com/coldsmirk/vef-framework-go/tabular.ExportError
  FIELD Row : int [field_order=1 tag=""]
  FIELD Column : string [field_order=2 tag=""]
  FIELD Field : string [field_order=3 tag=""]
  FIELD Err : error [field_order=4 tag=""]
  METHOD Error : func() string
  METHOD Unwrap : func() error
TYPE Exporter : github.com/coldsmirk/vef-framework-go/tabular.Exporter
  METHOD Export : func(data any) (*bytes.Buffer, error)
  METHOD ExportToFile : func(data any, filename string) error
  METHOD RegisterFormatter : func(name string, formatter github.com/coldsmirk/vef-framework-go/tabular.Formatter)
TYPE Formatter : github.com/coldsmirk/vef-framework-go/tabular.Formatter
  METHOD Format : func(value any) (string, error)
TYPE FormatterFunc : github.com/coldsmirk/vef-framework-go/tabular.FormatterFunc
  METHOD Format : func(value any) (string, error)
CONST IgnoreField : untyped string = "-"
TYPE ImportError : github.com/coldsmirk/vef-framework-go/tabular.ImportError
  FIELD Row : int [field_order=1 tag=""]
  FIELD Column : string [field_order=2 tag=""]
  FIELD Field : string [field_order=3 tag=""]
  FIELD Err : error [field_order=4 tag=""]
  METHOD Error : func() string
  METHOD Unwrap : func() error
FUNC ImportRows : func(rows [][]string, adapter github.com/coldsmirk/vef-framework-go/tabular.RowAdapter, parsers map[string]github.com/coldsmirk/vef-framework-go/tabular.ValueParser, opts github.com/coldsmirk/vef-framework-go/tabular.ImportRowsOptions) (any, []github.com/coldsmirk/vef-framework-go/tabular.ImportError, error)
TYPE ImportRowsOptions : github.com/coldsmirk/vef-framework-go/tabular.ImportRowsOptions
  FIELD SkipRows : int [field_order=1 tag=""]
  FIELD HasHeader : bool [field_order=2 tag=""]
  FIELD TrimSpace : bool [field_order=3 tag=""]
TYPE Importer : github.com/coldsmirk/vef-framework-go/tabular.Importer
  METHOD Import : func(reader io.Reader) (any, []github.com/coldsmirk/vef-framework-go/tabular.ImportError, error)
  METHOD ImportFromFile : func(filename string) (any, []github.com/coldsmirk/vef-framework-go/tabular.ImportError, error)
  METHOD RegisterParser : func(name string, parser github.com/coldsmirk/vef-framework-go/tabular.ValueParser)
FUNC IsDefaultFormatter : func(col *github.com/coldsmirk/vef-framework-go/tabular.Column, registry map[string]github.com/coldsmirk/vef-framework-go/tabular.Formatter) bool
FUNC IsEmptyRow : func(cells []string, trimSpace bool) bool
TYPE MapOption : github.com/coldsmirk/vef-framework-go/tabular.MapOption
TYPE MappingOptions : github.com/coldsmirk/vef-framework-go/tabular.MappingOptions
  FIELD TrimSpace : bool [field_order=1 tag=""]
FUNC NewColumnMapping : func(m map[int]int) github.com/coldsmirk/vef-framework-go/tabular.ColumnMapping
FUNC NewDefaultFormatter : func(format string) github.com/coldsmirk/vef-framework-go/tabular.Formatter
FUNC NewDefaultParser : func(format string) github.com/coldsmirk/vef-framework-go/tabular.ValueParser
FUNC NewMapAdapter : func(schema *github.com/coldsmirk/vef-framework-go/tabular.Schema, opts ...github.com/coldsmirk/vef-framework-go/tabular.MapOption) github.com/coldsmirk/vef-framework-go/tabular.RowAdapter
FUNC NewMapAdapterFromSpecs : func(specs []github.com/coldsmirk/vef-framework-go/tabular.ColumnSpec, opts ...github.com/coldsmirk/vef-framework-go/tabular.MapOption) (github.com/coldsmirk/vef-framework-go/tabular.RowAdapter, error)
FUNC NewSchema : func(typ reflect.Type) *github.com/coldsmirk/vef-framework-go/tabular.Schema
FUNC NewSchemaFor : func[T any]() *github.com/coldsmirk/vef-framework-go/tabular.Schema
FUNC NewSchemaFromSpecs : func(specs []github.com/coldsmirk/vef-framework-go/tabular.ColumnSpec) (*github.com/coldsmirk/vef-framework-go/tabular.Schema, error)
FUNC NewStructAdapter : func(typ reflect.Type) github.com/coldsmirk/vef-framework-go/tabular.RowAdapter
FUNC NewStructAdapterFor : func[T any]() github.com/coldsmirk/vef-framework-go/tabular.RowAdapter
FUNC NewTypedExporter : func[T any](inner github.com/coldsmirk/vef-framework-go/tabular.Exporter) github.com/coldsmirk/vef-framework-go/tabular.TypedExporter[T]
FUNC NewTypedImporter : func[T any](inner github.com/coldsmirk/vef-framework-go/tabular.Importer) github.com/coldsmirk/vef-framework-go/tabular.TypedImporter[T]
FUNC ParseRow : func(cells []string, mapping github.com/coldsmirk/vef-framework-go/tabular.ColumnMapping, schema *github.com/coldsmirk/vef-framework-go/tabular.Schema, builder github.com/coldsmirk/vef-framework-go/tabular.RowBuilder, parsers []github.com/coldsmirk/vef-framework-go/tabular.ValueParser, rowNumber int, opts github.com/coldsmirk/vef-framework-go/tabular.ParseRowOptions) []github.com/coldsmirk/vef-framework-go/tabular.ImportError
TYPE ParseRowOptions : github.com/coldsmirk/vef-framework-go/tabular.ParseRowOptions
  FIELD TrimSpace : bool [field_order=1 tag=""]
TYPE ParserFunc : github.com/coldsmirk/vef-framework-go/tabular.ParserFunc
  METHOD Parse : func(cellValue string, targetType reflect.Type) (any, error)
FUNC ResolveFormatter : func(col *github.com/coldsmirk/vef-framework-go/tabular.Column, registry map[string]github.com/coldsmirk/vef-framework-go/tabular.Formatter) github.com/coldsmirk/vef-framework-go/tabular.Formatter
FUNC ResolveFormatters : func(schema *github.com/coldsmirk/vef-framework-go/tabular.Schema, registry map[string]github.com/coldsmirk/vef-framework-go/tabular.Formatter) []github.com/coldsmirk/vef-framework-go/tabular.Formatter
FUNC ResolveParser : func(col *github.com/coldsmirk/vef-framework-go/tabular.Column, registry map[string]github.com/coldsmirk/vef-framework-go/tabular.ValueParser) github.com/coldsmirk/vef-framework-go/tabular.ValueParser
FUNC ResolveParsers : func(schema *github.com/coldsmirk/vef-framework-go/tabular.Schema, registry map[string]github.com/coldsmirk/vef-framework-go/tabular.ValueParser) []github.com/coldsmirk/vef-framework-go/tabular.ValueParser
TYPE RowAdapter : github.com/coldsmirk/vef-framework-go/tabular.RowAdapter
  METHOD Reader : func(data any) (github.com/coldsmirk/vef-framework-go/tabular.RowReader, error)
  METHOD Schema : func() *github.com/coldsmirk/vef-framework-go/tabular.Schema
  METHOD Writer : func(capacity int) github.com/coldsmirk/vef-framework-go/tabular.RowWriter
TYPE RowBuilder : github.com/coldsmirk/vef-framework-go/tabular.RowBuilder
  METHOD Set : func(column *github.com/coldsmirk/vef-framework-go/tabular.Column, value any) error
  METHOD Validate : func() error
  METHOD Value : func() any
TYPE RowReader : github.com/coldsmirk/vef-framework-go/tabular.RowReader
  METHOD All : func() iter.Seq2[int, github.com/coldsmirk/vef-framework-go/tabular.RowView]
TYPE RowValidator : github.com/coldsmirk/vef-framework-go/tabular.RowValidator
TYPE RowView : github.com/coldsmirk/vef-framework-go/tabular.RowView
  METHOD Get : func(column *github.com/coldsmirk/vef-framework-go/tabular.Column) (any, error)
TYPE RowWriter : github.com/coldsmirk/vef-framework-go/tabular.RowWriter
  METHOD Build : func() any
  METHOD Commit : func(row github.com/coldsmirk/vef-framework-go/tabular.RowBuilder) error
  METHOD NewRow : func() github.com/coldsmirk/vef-framework-go/tabular.RowBuilder
TYPE Schema : github.com/coldsmirk/vef-framework-go/tabular.Schema
  METHOD ColumnByKey : func(key string) (*github.com/coldsmirk/vef-framework-go/tabular.Column, bool)
  METHOD ColumnByName : func(name string) (*github.com/coldsmirk/vef-framework-go/tabular.Column, bool)
  METHOD ColumnCount : func() int
  METHOD ColumnNames : func() []string
  METHOD Columns : func() []*github.com/coldsmirk/vef-framework-go/tabular.Column
CONST TagTabular : untyped string = "tabular"
TYPE TypedExporter : github.com/coldsmirk/vef-framework-go/tabular.TypedExporter[T any]
  METHOD Export : func(rows []T) (*bytes.Buffer, error)
  METHOD ExportToFile : func(rows []T, filename string) error
  METHOD Inner : func() github.com/coldsmirk/vef-framework-go/tabular.Exporter
  METHOD RegisterFormatter : func(name string, formatter github.com/coldsmirk/vef-framework-go/tabular.Formatter)
TYPE TypedImporter : github.com/coldsmirk/vef-framework-go/tabular.TypedImporter[T any]
  METHOD Import : func(reader io.Reader) ([]T, []github.com/coldsmirk/vef-framework-go/tabular.ImportError, error)
  METHOD ImportFromFile : func(filename string) ([]T, []github.com/coldsmirk/vef-framework-go/tabular.ImportError, error)
  METHOD Inner : func() github.com/coldsmirk/vef-framework-go/tabular.Importer
  METHOD RegisterParser : func(name string, parser github.com/coldsmirk/vef-framework-go/tabular.ValueParser)
TYPE ValueParser : github.com/coldsmirk/vef-framework-go/tabular.ValueParser
  METHOD Parse : func(cellValue string, targetType reflect.Type) (any, error)
FUNC WithRowValidator : func(validator github.com/coldsmirk/vef-framework-go/tabular.RowValidator) github.com/coldsmirk/vef-framework-go/tabular.MapOption

## github.com/coldsmirk/vef-framework-go/timex
TYPE Date : github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD AddDate : func(years int, months int, days int) github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD AddDays : func(days int) github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD AddMonths : func(months int) github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD AddYears : func(years int) github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD After : func(other github.com/coldsmirk/vef-framework-go/timex.Date) bool
  METHOD Before : func(other github.com/coldsmirk/vef-framework-go/timex.Date) bool
  METHOD BeginOfDay : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD BeginOfMonth : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD BeginOfQuarter : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD BeginOfWeek : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD BeginOfYear : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Between : func(start github.com/coldsmirk/vef-framework-go/timex.Date, end github.com/coldsmirk/vef-framework-go/timex.Date) bool
  METHOD Day : func() int
  METHOD EndOfDay : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD EndOfMonth : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD EndOfQuarter : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD EndOfWeek : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD EndOfYear : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Equal : func(other github.com/coldsmirk/vef-framework-go/timex.Date) bool
  METHOD Format : func(layout string) string
  METHOD Friday : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD IsZero : func() bool
  METHOD Location : func() *time.Location
  METHOD MarshalJSON : func() ([]byte, error)
  METHOD MarshalText : func() ([]byte, error)
  METHOD Monday : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Month : func() time.Month
  METHOD Saturday : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Scan : func(src any) error
  METHOD Since : func() time.Duration
  METHOD String : func() string
  METHOD Sub : func(other github.com/coldsmirk/vef-framework-go/timex.Date) time.Duration
  METHOD Sunday : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Thursday : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Tuesday : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Unix : func() int64
  METHOD UnmarshalJSON : func(bs []byte) error
  METHOD UnmarshalText : func(text []byte) error
  METHOD Until : func() time.Duration
  METHOD Unwrap : func() time.Time
  METHOD Value : func() (database/sql/driver.Value, error)
  METHOD Wednesday : func() github.com/coldsmirk/vef-framework-go/timex.Date
  METHOD Weekday : func() time.Weekday
  METHOD Year : func() int
  METHOD YearDay : func() int
FUNC DateOf : func(t time.Time) github.com/coldsmirk/vef-framework-go/timex.Date
TYPE DateTime : github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Add : func(d time.Duration) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD AddDate : func(years int, months int, days int) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD AddDays : func(days int) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD AddHours : func(hours int) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD AddMinutes : func(minutes int) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD AddMonths : func(months int) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD AddSeconds : func(seconds int) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD AddYears : func(years int) github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD After : func(other github.com/coldsmirk/vef-framework-go/timex.DateTime) bool
  METHOD Before : func(other github.com/coldsmirk/vef-framework-go/timex.DateTime) bool
  METHOD BeginOfDay : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD BeginOfHour : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD BeginOfMinute : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD BeginOfMonth : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD BeginOfQuarter : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD BeginOfWeek : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD BeginOfYear : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Between : func(start github.com/coldsmirk/vef-framework-go/timex.DateTime, end github.com/coldsmirk/vef-framework-go/timex.DateTime) bool
  METHOD Day : func() int
  METHOD EndOfDay : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD EndOfHour : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD EndOfMinute : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD EndOfMonth : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD EndOfQuarter : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD EndOfWeek : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD EndOfYear : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Equal : func(other github.com/coldsmirk/vef-framework-go/timex.DateTime) bool
  METHOD Format : func(layout string) string
  METHOD Friday : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Hour : func() int
  METHOD IsZero : func() bool
  METHOD Location : func() *time.Location
  METHOD MarshalJSON : func() ([]byte, error)
  METHOD MarshalText : func() ([]byte, error)
  METHOD Minute : func() int
  METHOD Monday : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Month : func() time.Month
  METHOD Nanosecond : func() int
  METHOD Saturday : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Scan : func(src any) error
  METHOD Second : func() int
  METHOD Since : func() time.Duration
  METHOD String : func() string
  METHOD Sub : func(other github.com/coldsmirk/vef-framework-go/timex.DateTime) time.Duration
  METHOD Sunday : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Thursday : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Tuesday : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Unix : func() int64
  METHOD UnixMicro : func() int64
  METHOD UnixMilli : func() int64
  METHOD UnixNano : func() int64
  METHOD UnmarshalJSON : func(bs []byte) error
  METHOD UnmarshalText : func(text []byte) error
  METHOD Until : func() time.Duration
  METHOD Unwrap : func() time.Time
  METHOD Value : func() (database/sql/driver.Value, error)
  METHOD Wednesday : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
  METHOD Weekday : func() time.Weekday
  METHOD Year : func() int
  METHOD YearDay : func() int
VAR ErrFailedScan : error
VAR ErrInvalidDateFormat : error
VAR ErrInvalidDateTimeFormat : error
VAR ErrInvalidTimeFormat : error
VAR ErrUnsupportedDestType : error
FUNC FromUnix : func(sec int64, nsec int64) github.com/coldsmirk/vef-framework-go/timex.DateTime
FUNC FromUnixMicro : func(usec int64) github.com/coldsmirk/vef-framework-go/timex.DateTime
FUNC FromUnixMilli : func(msec int64) github.com/coldsmirk/vef-framework-go/timex.DateTime
FUNC Now : func() github.com/coldsmirk/vef-framework-go/timex.DateTime
FUNC NowDate : func() github.com/coldsmirk/vef-framework-go/timex.Date
FUNC NowTime : func() github.com/coldsmirk/vef-framework-go/timex.Time
FUNC Of : func(t time.Time) github.com/coldsmirk/vef-framework-go/timex.DateTime
FUNC Parse : func(value string, pattern ...string) (github.com/coldsmirk/vef-framework-go/timex.DateTime, error)
FUNC ParseDate : func(value string, pattern ...string) (github.com/coldsmirk/vef-framework-go/timex.Date, error)
FUNC ParseTime : func(value string, pattern ...string) (github.com/coldsmirk/vef-framework-go/timex.Time, error)
TYPE Time : github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD Add : func(d time.Duration) github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD AddHours : func(hours int) github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD AddMicroseconds : func(microseconds int64) github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD AddMilliseconds : func(milliseconds int64) github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD AddMinutes : func(minutes int) github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD AddNanoseconds : func(nanoseconds int64) github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD AddSeconds : func(seconds int) github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD After : func(other github.com/coldsmirk/vef-framework-go/timex.Time) bool
  METHOD Before : func(other github.com/coldsmirk/vef-framework-go/timex.Time) bool
  METHOD BeginOfHour : func() github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD BeginOfMinute : func() github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD Between : func(start github.com/coldsmirk/vef-framework-go/timex.Time, end github.com/coldsmirk/vef-framework-go/timex.Time) bool
  METHOD EndOfHour : func() github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD EndOfMinute : func() github.com/coldsmirk/vef-framework-go/timex.Time
  METHOD Equal : func(other github.com/coldsmirk/vef-framework-go/timex.Time) bool
  METHOD Format : func(layout string) string
  METHOD Hour : func() int
  METHOD IsZero : func() bool
  METHOD MarshalJSON : func() ([]byte, error)
  METHOD MarshalText : func() ([]byte, error)
  METHOD Minute : func() int
  METHOD Nanosecond : func() int
  METHOD Scan : func(src any) error
  METHOD Second : func() int
  METHOD String : func() string
  METHOD Sub : func(other github.com/coldsmirk/vef-framework-go/timex.Time) time.Duration
  METHOD ToDuration : func() time.Duration
  METHOD UnmarshalJSON : func(bs []byte) error
  METHOD UnmarshalText : func(text []byte) error
  METHOD Unwrap : func() time.Time
  METHOD Value : func() (database/sql/driver.Value, error)
FUNC TimeOf : func(t time.Time) github.com/coldsmirk/vef-framework-go/timex.Time

## github.com/coldsmirk/vef-framework-go/tree
TYPE Adapter : github.com/coldsmirk/vef-framework-go/tree.Adapter[T any]
  FIELD GetID : func(T) string [field_order=1 tag=""]
  FIELD GetParentID : func(T) *string [field_order=2 tag=""]
  FIELD GetChildren : func(T) []T [field_order=3 tag=""]
  FIELD SetChildren : func(*T, []T) [field_order=4 tag=""]
FUNC Build : func[T any](nodes []T, adapter github.com/coldsmirk/vef-framework-go/tree.Adapter[T]) []T
FUNC FindNode : func[T any](roots []T, targetID string, adapter github.com/coldsmirk/vef-framework-go/tree.Adapter[T]) (T, bool)
FUNC FindNodePath : func[T any](roots []T, targetID string, adapter github.com/coldsmirk/vef-framework-go/tree.Adapter[T]) ([]T, bool)

## github.com/coldsmirk/vef-framework-go/validator
TYPE CustomTypeFunc : github.com/coldsmirk/vef-framework-go/validator.CustomTypeFunc
FUNC RegisterTypeFunc : func(fn github.com/coldsmirk/vef-framework-go/validator.CustomTypeFunc, types ...any)
FUNC RegisterValidationRules : func(rules ...github.com/coldsmirk/vef-framework-go/validator.ValidationRule) error
FUNC Validate : func(value any) error
TYPE ValidationRule : github.com/coldsmirk/vef-framework-go/validator.ValidationRule
  FIELD RuleTag : string [field_order=1 tag=""]
  FIELD ErrMessageTemplate : string [field_order=2 tag=""]
  FIELD ErrMessageI18nKey : string [field_order=3 tag=""]
  FIELD Validate : func(fl github.com/go-playground/validator/v10.FieldLevel) bool [field_order=4 tag=""]
  FIELD ParseParam : func(fe github.com/go-playground/validator/v10.FieldError) []string [field_order=5 tag=""]
  FIELD CallValidationEvenIfNull : bool [field_order=6 tag=""]

## github.com/coldsmirk/vef-framework-go/version
CONST VEFVersion : untyped string = "v0.37.0"
```
