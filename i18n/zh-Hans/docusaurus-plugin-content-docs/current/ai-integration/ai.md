---
sidebar_position: 1
---

# AI 契约

`ai` 包定义了与 provider 无关的 chat model、tool、agent、message 和 stream
契约。框架默认不内置具体模型 provider；应用或 provider 包通过公开 registry
注册自己的实现。

## Messages

```go
type Message struct {
    Role       Role
    Content    string
    ToolCalls  []ToolCall
    ToolResult *ToolResult
    Usage      *TokenUsage
}
```

Role 公开为常量：`RoleSystem`、`RoleUser`、`RoleAssistant` 和 `RoleTool`。
常见消息可用 helper 构造：

```go
messages := []*ai.Message{
    ai.NewSystemMessage("Answer briefly."),
    ai.NewUserMessage("Summarize this invoice."),
}
```

Assistant tool request 使用 `NewAssistantMessageWithToolCalls(...)`，tool
result 使用 `NewToolMessage(callID, content)`；普通 assistant 回复用
`NewAssistantMessage(...)`。`Message.IsSystem`、`IsUser`、`IsAssistant`、
`IsTool` 和 `HasToolCalls` 是给路由与测试使用的便捷谓词。

helper 构造函数会精确填充这些字段：

| Helper | 填充字段 |
| --- | --- |
| `NewSystemMessage(content)` | `Role: RoleSystem`, `Content: content` |
| `NewUserMessage(content)` | `Role: RoleUser`, `Content: content` |
| `NewAssistantMessage(content)` | `Role: RoleAssistant`, `Content: content` |
| `NewAssistantMessageWithToolCalls(content, toolCalls)` | `Role: RoleAssistant`, `Content: content`, `ToolCalls: toolCalls` |
| `NewToolMessage(callID, content)` | `Role: RoleTool`, `ToolResult: &ToolResult{CallID: callID, Content: content}`；顶层 `Message.Content` 不会被填充 |

## Models

```go
type ChatModel interface {
    Generate(ctx context.Context, messages []*Message, opts ...Option) (*Message, error)
    Stream(ctx context.Context, messages []*Message, opts ...Option) (MessageStream, error)
}

type ToolableChatModel interface {
    ChatModel
    WithTools(tools ...Tool) ToolableChatModel
}
```

`WithTools` 使用不可变模式：它返回一个绑定了 tools 的 model instance，不会
修改当前 instance。

运行时 options 是 functional options：

```go
reply, err := model.Generate(ctx, messages,
    ai.WithTemperature(0.2),
    ai.WithMaxTokens(512),
    ai.WithStopSequences("\n\n"),
    ai.WithMeta("trace_id", traceID),
)
```

`NewOptions()` 会创建默认 option 累加器；`Options.Apply(...)` 把 functional
options 折叠进去。大多数应用代码直接把 options 传给 `Generate`、`Stream`、
`Run` 或 `Invoke`，不需要手工操作累加器。
`Apply` 会修改并返回同一个累加器，按参数顺序应用 options；后传入的 option
可以覆盖之前 option 设置的字段。

运行时 option 字段与 helper：

| 字段或 helper | 契约 |
| --- | --- |
| `Options.Temperature` / `WithTemperature(t)` | 可选 temperature pointer，值为 `t` |
| `Options.MaxTokens` / `WithMaxTokens(n)` | 可选 maximum-token pointer，值为 `n` |
| `Options.StopSequences` / `WithStopSequences(seqs...)` | stop sequence slice 会被替换为 `seqs` |
| `Options.Meta` / `WithMeta(key, value)` | string metadata map；`NewOptions()` 会初始化它，`WithMeta` 在需要时也会创建它 |

Model metadata 和 configuration 字段：

| Struct | 公开字段 |
| --- | --- |
| `ModelConfig` | `Provider`, `Model`, `APIKey`, `BaseURL`, `Temperature`, `MaxTokens`, `Timeout` |
| `ModelInfo` | `Provider`, `Model`, `MaxTokens`, `Temperature` |

## Provider Registry

模型 provider 实现：

```go
type ModelProvider interface {
    Name() string
    CreateModel(ctx context.Context, cfg *ModelConfig) (ToolableChatModel, error)
}
```

在应用启动或 package initialization 阶段注册一次：

```go
ai.RegisterModelProvider(openaiProvider)

model, err := ai.NewChatModel(ctx, &ai.ModelConfig{
    Provider: "openai",
    Model:    "gpt-4.1-mini",
    APIKey:   apiKey,
})
```

`RegisterModelProvider` 遇到重复名称会 panic，消息为
`ai: model provider "{name}" already registered`；`RegisterAgentFactory` 使用
`ai: agent factory "{name}" already registered`。`NewChatModel` 在
`cfg.Provider` 没有对应 provider 时返回带有请求 provider 名称的
`ai.ErrProviderNotFound` wrapped error。`NewAgentBuilder` 也会用
`ai.ErrAgentNotFound` 和请求的 agent type 做同样处理。`cfg`、provider
instance 和 factory instance 都是非 nil 前置条件；registry 在调用 `Name()` 或
读取 `cfg.Provider` 前不会额外做 nil guard。

Agent factory 走同样模式：`RegisterAgentFactory(...)`、
`NewAgentBuilder(...)`、`ListAgentFactories()` 和 `ai.ErrAgentNotFound`。
`ListModelProviders()` 可用于查看已注册的模型 provider 名称。registry list
helpers 返回 map keys，不定义稳定排序。

Factory 实现会交换 `ModelConfig`、`ModelInfo`、`AgentConfig` 和
`AgentFactory`。这些都是给 provider package 使用的普通公开 struct /
interface；应用通常只在注册自己的 model 或 agent 实现时直接接触它们。

## Tools

Tool 暴露 JSON-schema-style metadata 和 invocation 方法：

```go
type Tool interface {
    Info() *ToolInfo
    Invoke(ctx context.Context, arguments string) (string, error)
}
```

`arguments` 是模型传来的 JSON-encoded string。Tool 实现应在工具边界解码并
校验它，再返回给模型使用的 string payload。

Tool metadata 使用 `ToolInfo`、`ParameterSchema` 和 `PropertySchema`。schema
形态刻意保持很小：`type`、`properties`、`required`、`description`、`enum`
和 `items` 覆盖 provider tool calling 需要的 JSON-schema 子集。

Tool 相关数据形状：

| Struct | 公开字段 |
| --- | --- |
| `ToolCall` | `ID`, `Name`, `Arguments` |
| `ToolResult` | `CallID`, `Content` |
| `ToolInfo` | `Name`, `Description`, `Parameters` |
| `ParameterSchema` | `Type`, `Properties`, `Required` |
| `PropertySchema` | `Type`, `Description`, `Enum`, `Items` |

支持流式输出的 tool 实现 `StreamableTool`：

```go
type StreamableTool interface {
    Tool
    InvokeStream(ctx context.Context, arguments string) (StringStream, error)
}
```

## Agents

`Agent` 是更高层的 runner，可以用 model 和 tools 完成推理：

```go
type Agent interface {
    Run(ctx context.Context, input string, opts ...Option) (*Message, error)
    Stream(ctx context.Context, input string, opts ...Option) (MessageStream, error)
}
```

`AgentBuilder` 让已注册 factory 暴露 fluent setup：

```go
builder, err := ai.NewAgentBuilder("tool-loop")
if err != nil {
    return err
}

agent, err := builder.
    WithModel(model).
    WithTools(searchTool, calculatorTool).
    WithSystemPrompt("Use tools when needed.").
    WithMaxIterations(8).
    Build(ctx)
```

`AgentConfig` 暴露同一组 builder 输入字段：`Model`、`Tools`、
`SystemPrompt` 和 `MaxIterations`。

## Streams

base package 定义了两种 stream contract：

```go
type MessageStream interface {
    io.Closer
    Recv() (*MessageChunk, error)
    Collect() (*Message, error)
}

type StringStream interface {
    io.Closer
    Recv() (string, error)
    Collect() (string, error)
}
```

`Recv` 在 stream 耗尽时返回 `io.EOF`。

Stream 数据形状：

| Struct | 公开字段 |
| --- | --- |
| `MessageChunk` | `Content`, `ToolCalls`, `Done` |
| `TokenUsage` | `PromptTokens`, `CompletionTokens`, `TotalTokens` |

## UI Message SSE Streams

`ai/stream` 包会把消息转换成兼容 AI SDK UI Message Stream protocol 的
Server-Sent Events。它设置标准 SSE headers，输出 `data: ...` frame，并以
`data: [DONE]` 结束。

常见 adapter：

| Helper | Source |
| --- | --- |
| `stream.FromChannel(ch)` | `stream.Message` 的 Go channel |
| `stream.FromCallback(fn)` | 通过 `CallbackWriter` 写入的 callback |
| `stream.FromAiMessageStream(s)` | 一个 `ai.MessageStream` |
| `stream.NewChannelSource(ch)` | channel 的底层 `MessageSource` adapter |
| `stream.NewCallbackSource(fn)` | callback 的底层 `MessageSource` |
| `stream.NewAiMessageStreamSource(s)` | `ai.MessageStream` 的底层 adapter |

内置 adapter 行为：

| Adapter | 行为 |
| --- | --- |
| `NewChannelSource(ch)` | `Recv()` 会一直返回 channel message，直到 channel 关闭后返回 `io.EOF`；`Close()` 只把 source 标记为已关闭，后续 `Recv()` 返回 `io.EOF`，不会关闭调用方持有的 channel |
| `NewCallbackSource(fn)` | 在 goroutine 中运行 `fn`，并使用带缓冲的 message channel；已排队的消息会先交付，然后才返回 callback error；`Close()` 会等待 callback goroutine 结束 |
| `NewAiMessageStreamSource(s)` | 把每个 `ai.MessageChunk` 转成 assistant `stream.Message`，转发 `Content` 和 tool calls；`Close()` 委托给被包装的 `ai.MessageStream` |

Fiber handler 示例：

```go
import aistream "github.com/coldsmirk/vef-framework-go/ai/stream"

func Chat(ctx fiber.Ctx) error {
    return aistream.FromCallback(func(w aistream.CallbackWriter) error {
        w.WriteText("Hello")
        w.WriteText(" from VEF")
        w.WriteData("usage", map[string]int{"totalTokens": 12})

        return nil
    }).Stream(ctx)
}
```

Builder 默认输出 start/finish chunks；当 source 提供对应 message parts 时，
还会输出 text、reasoning、tool input、tool output 和自定义 `data-{type}`
chunks。source 与 file chunks 通过低层构造函数开放给自定义 writer；
`WithSources(bool)` 目前只是切换公开 option 字段，不会让普通
`MessageSource` 自动输出 sources。
`WithReasoning(bool)` 会控制 reasoning 输出，而且只有传入的
`stream.Message.Reasoning` 非空时才会输出 reasoning chunks。

默认 stream options：

| Option | 默认值 |
| --- | --- |
| `SendReasoning` | `true` |
| `SendSources` | `true` |
| `SendStart` | `true` |
| `SendFinish` | `true` |
| `OnError` | 返回 `err.Error()` |
| `OnFinish` | `nil` |
| `GenerateID` | `DefaultOptions()` 中为 `nil`；builder fallback 为 `prefix + "_" + id.GenerateUUID()` |

Builder 控制项：

| 方法 | 作用 |
| --- | --- |
| `WithSource(source)` | 设置自定义 `MessageSource` |
| `WithMessageID(id)` / `WithIDGenerator(fn)` | 控制生成的 chunk ID |
| `WithStart(bool)` / `WithFinish(bool)` | 开关 start/finish chunk |
| `WithReasoning(bool)` / `WithSources(bool)` | 开关 reasoning 输出与公开的 source-output option |
| `WithHeader(key, value)` | 增加 SSE response header |
| `OnError(fn)` / `OnFinish(fn)` | 自定义错误文本或完成回调 |
| `Stream(ctx)` / `StreamToWriter(w)` | 写入 Fiber 或 buffered writer |

执行路径契约：

| 路径 | 行为 |
| --- | --- |
| `Stream(ctx)` 未配置 `WithSource(...)` | 写 headers 前返回 `stream.ErrSourceRequired` |
| `Stream(ctx)` 已配置 source | 先写 `SseHeaders`，再写 `WithHeader` 的覆盖/新增 header，然后通过 `fiber.Ctx.SendStreamWriter` 输出 |
| `StreamToWriter(w)` | 直接写入 buffered writer，必须先配置 source；不同于 `Stream(ctx)`，它没有缺失 source 的 guard |
| source lifecycle | stream loop 退出时会关闭已配置的 source |
| source 正常返回 `io.EOF` | 关闭已打开的 text/reasoning chunks；当 `SendFinish` 为 true 时写 `finish-step` 和 `finish`；写 `data: [DONE]`；配置了 `OnFinish` 时调用 `OnFinish(fullContent)` |
| 禁用 start/finish | `WithStart(false)` 会同时省略 `start` 和 `start-step`；`WithFinish(false)` 会同时省略 `finish-step` 和 `finish`，但仍然写 `data: [DONE]` |
| source 返回非 EOF 错误 | 用 `OnError(err)` 写 `error` chunk，写 `data: [DONE]`，不会调用 `OnFinish`；已打开的 text/reasoning chunks 不会先被显式结束 |
| `OnFinish(fullContent)` | 只收到拼接后的 `Content` text chunks，不包含 reasoning、tool payload 或自定义 data |
| tool call arguments | JSON-unmarshal 后作为 `input` payload；非法 JSON 会按原始 string 输出；空 tool call ID 会被替换为生成的 `call_*` ID |
| tool result content | JSON-unmarshal 后作为 `output` payload；非法 JSON 会按原始 string 输出 |
| 自定义 data messages | `Data` map 的每个 entry 会输出一个 `data-{type}` chunk；同一条 message 包含多个 data key 时，map iteration order 不稳定 |
| writer errors | stream loop 内部的 chunk/write error 会被忽略，因为 `StreamToWriter` 和 Fiber stream callback 不会从 loop 返回 error |

默认 SSE headers：

| Header | Value |
| --- | --- |
| `Content-Type` | `text/event-stream` |
| `Cache-Control` | `no-cache` |
| `Connection` | `keep-alive` |
| `Transfer-Encoding` | `chunked` |
| `X-Vercel-AI-UI-Message-Stream` | `v1` |
| `X-Accel-Buffering` | `no` |

低层 chunk 构造函数也是公开 API，适合测试和自定义 writer 使用：

| Chunk 家族 | 构造函数 |
| --- | --- |
| 生命周期 | `NewStartChunk`, `NewFinishChunk`, `NewStartStepChunk`, `NewFinishStepChunk`, `NewErrorChunk` |
| text | `NewTextStartChunk`, `NewTextDeltaChunk`, `NewTextEndChunk` |
| reasoning | `NewReasoningStartChunk`, `NewReasoningDeltaChunk`, `NewReasoningEndChunk` |
| tools | `NewToolInputStartChunk`, `NewToolInputDeltaChunk`, `NewToolInputAvailableChunk`, `NewToolOutputAvailableChunk` |
| sources/files/data | `NewSourceURLChunk`, `NewSourceDocumentChunk`, `NewFileChunk`, `NewDataChunk` |

精确 chunk wire shapes：

| 构造函数 | Wire `type` | 字段 |
| --- | --- | --- |
| `NewStartChunk(messageID)` | `start` | `type`, `messageID` |
| `NewFinishChunk()` | `finish` | `type` |
| `NewStartStepChunk()` | `start-step` | `type` |
| `NewFinishStepChunk()` | `finish-step` | `type` |
| `NewErrorChunk(errorText)` | `error` | `type`, `errorText` |
| `NewTextStartChunk(id)` | `text-start` | `type`, `id` |
| `NewTextDeltaChunk(id, delta)` | `text-delta` | `type`, `id`, `delta` |
| `NewTextEndChunk(id)` | `text-end` | `type`, `id` |
| `NewReasoningStartChunk(id)` | `reasoning-start` | `type`, `id` |
| `NewReasoningDeltaChunk(id, delta)` | `reasoning-delta` | `type`, `id`, `delta` |
| `NewReasoningEndChunk(id)` | `reasoning-end` | `type`, `id` |
| `NewToolInputStartChunk(toolCallID, toolName)` | `tool-input-start` | `type`, `toolCallID`, `toolName` |
| `NewToolInputDeltaChunk(toolCallID, delta)` | `tool-input-delta` | `type`, `toolCallID`, `inputTextDelta` |
| `NewToolInputAvailableChunk(toolCallID, toolName, input)` | `tool-input-available` | `type`, `toolCallID`, `toolName`, `input` |
| `NewToolOutputAvailableChunk(toolCallID, output)` | `tool-output-available` | `type`, `toolCallID`, `output` |
| `NewSourceURLChunk(sourceID, url, title)` | `source-url` | `type`, `sourceID`, `url`；`title` 仅在非空时出现 |
| `NewSourceDocumentChunk(sourceID, mediaType, title)` | `source-document` | `type`, `sourceID`, `mediaType`；`title` 仅在非空时出现 |
| `NewFileChunk(fileID, mediaType, url)` | `file` | `type`, `fileID`, `mediaType`, `url` |
| `NewDataChunk(dataType, data)` | `data-{dataType}` | `type`, `data`；`data-{dataType}` 是动态值，不是 `ChunkType` 常量 |

`CallbackWriter` 写入这些 source messages：

| 方法 | Message 效果 |
| --- | --- |
| `WriteText(content)` | 带 `Content` 的 assistant message |
| `WriteToolCall(id, name, arguments)` | 带一个 `ToolCall` 的 assistant message |
| `WriteToolResult(toolCallID, content)` | 带 `ToolCallID` 和 `Content` 的 tool message |
| `WriteReasoning(reasoning)` | 带 `Reasoning` 的 assistant message |
| `WriteData(dataType, data)` | 带 `Data[dataType]` 的 assistant message |
| `WriteMessage(msg)` | 原样转发传入的 `stream.Message` |

stream package 数据形状：

| Struct 或 interface | 公开字段或方法 |
| --- | --- |
| `stream.Message` | `Role`, `Content`, `ToolCalls`, `ToolCallID`, `Reasoning`, `Data` |
| `stream.ToolCall` | `ID`, `Name`, `Arguments` |
| `stream.Source` | `Type`, `ID`, `URL`, `Title`, `MediaType` |
| `stream.Options` | `SendReasoning`, `SendSources`, `SendStart`, `SendFinish`, `OnError`, `OnFinish`, `GenerateID` |
| `MessageSource` | `Recv()`, `Close()`；`Recv()` 完成时返回 `io.EOF` |
| `StreamWriter` | `WriteChunk(chunk)`, `Flush()` |
| `ResponseWriter` | 与 Fiber stream writers 兼容的 `io.Writer` |

`stream.SseHeaders` 包含默认 AI SDK UI Message Stream headers。
builder 缺少 source 时返回 `stream.ErrSourceRequired`。`stream.ErrSourceClosed`
是留给自定义 `MessageSource` 实现使用的公开 sentinel；内置 channel 和
callback adapter 正常结束时返回 `io.EOF`。

其他 stream 公开 API：

| API 组 | 公开 surface |
| --- | --- |
| role 常量 | `stream.RoleUser`, `stream.RoleAssistant`, `stream.RoleTool`, `stream.RoleSystem` |
| chunk type 常量 | `ChunkType`, `ChunkTypeStart`, `ChunkTypeFinish`, `ChunkTypeStartStep`, `ChunkTypeFinishStep`, `ChunkTypeError`, `ChunkTypeTextStart`, `ChunkTypeTextDelta`, `ChunkTypeTextEnd`, `ChunkTypeReasoningStart`, `ChunkTypeReasoningDelta`, `ChunkTypeReasoningEnd`, `ChunkTypeToolInputStart`, `ChunkTypeToolInputDelta`, `ChunkTypeToolInputAvailable`, `ChunkTypeToolOutputAvailable`, `ChunkTypeSourceURL`, `ChunkTypeSourceDocument`, `ChunkTypeFile` |
| stream 数据结构 | `stream.Message`, `stream.ToolCall`, `stream.Source`, `stream.Chunk`, `stream.Options` |
| writer/source interface | `CallbackWriter`, `MessageSource`, `ResponseWriter`, `StreamWriter` |
| 默认值与构造器 | `stream.New()`, `DefaultOptions()`, `SseHeaders` |

## Error Sentinels

| Error | 含义 |
| --- | --- |
| `ai.ErrProviderNotFound` | 没有注册与 `ModelConfig.Provider` 匹配的 model provider |
| `ai.ErrModelNotSupported` | provider 不支持请求的 model |
| `ai.ErrAgentNotFound` | 没有注册与请求类型匹配的 agent factory |
| `ai.ErrToolNotFound` | agent/model 请求了不可用的 tool |
| `ai.ErrInvalidArguments` | tool arguments 校验失败 |
| `ai.ErrMaxIterationsReached` | agent 超过最大迭代次数 |
| `ai.ErrNoContent` | model 没有返回可用内容 |
| `ai.ErrStreamClosed` | stream 关闭后继续读取 |
| `stream.ErrSourceRequired` | `ai/stream` builder 没有 message source |
| `stream.ErrSourceClosed` | 给自定义 `MessageSource` 实现使用的公开 sentinel |

`NewToolError(...)` 与 `NewModelError(...)` 用于创建结构化的 `ToolError` /
`ModelError`。`ToolError.Unwrap()` 会暴露底层 tool error。

结构化 error 形状：

| Type | 公开字段与行为 |
| --- | --- |
| `ToolError` | `ToolName`, `Err`；`Error()` 返回 `ai: tool {ToolName}: {Err}`，`Unwrap()` 返回 `Err` |
| `ModelError` | `Provider`, `StatusCode`, `Message`；`Error()` 返回 `ai: {Provider}: {Message}` |
