---
sidebar_position: 17
---

# AI Contracts

The `ai` package defines provider-neutral contracts for chat models, tools,
agents, messages, and streams. The framework does not ship a concrete model
provider by default; applications or provider packages register implementations
through the public registry.

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

Roles are exported constants: `RoleSystem`, `RoleUser`, `RoleAssistant`, and
`RoleTool`. Helper constructors cover the common cases:

```go
messages := []*ai.Message{
    ai.NewSystemMessage("Answer briefly."),
    ai.NewUserMessage("Summarize this invoice."),
}
```

Use `NewAssistantMessageWithToolCalls(...)` for assistant tool requests and
`NewToolMessage(callID, content)` for tool results; `NewAssistantMessage(...)`
builds a plain assistant response. `Message.IsSystem`, `IsUser`,
`IsAssistant`, `IsTool`, and `HasToolCalls` are convenience predicates for
routers and tests.

Helper constructors populate exact fields:

| Helper | Fields set |
| --- | --- |
| `NewSystemMessage(content)` | `Role: RoleSystem`, `Content: content` |
| `NewUserMessage(content)` | `Role: RoleUser`, `Content: content` |
| `NewAssistantMessage(content)` | `Role: RoleAssistant`, `Content: content` |
| `NewAssistantMessageWithToolCalls(content, toolCalls)` | `Role: RoleAssistant`, `Content: content`, `ToolCalls: toolCalls` |
| `NewToolMessage(callID, content)` | `Role: RoleTool`, `ToolResult: &ToolResult{CallID: callID, Content: content}`; top-level `Message.Content` is not populated |

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

`WithTools` follows an immutable pattern: it returns a model instance with tools
bound and does not mutate the current one.

Runtime options are functional:

```go
reply, err := model.Generate(ctx, messages,
    ai.WithTemperature(0.2),
    ai.WithMaxTokens(512),
    ai.WithStopSequences("\n\n"),
    ai.WithMeta("trace_id", traceID),
)
```

`NewOptions()` creates the default option accumulator; `Options.Apply(...)`
folds functional options into it. Most application code passes options directly
to `Generate`, `Stream`, `Run`, or `Invoke` instead of touching the accumulator.
`Apply` mutates and returns the same accumulator, applies options in argument
order, and later options can overwrite fields set by earlier options.

Runtime option fields and helpers:

| Field or helper | Contract |
| --- | --- |
| `Options.Temperature` / `WithTemperature(t)` | optional temperature pointer set to `t` |
| `Options.MaxTokens` / `WithMaxTokens(n)` | optional maximum-token pointer set to `n` |
| `Options.StopSequences` / `WithStopSequences(seqs...)` | stop sequence slice replaced with `seqs` |
| `Options.Meta` / `WithMeta(key, value)` | string metadata map; `NewOptions()` initializes it and `WithMeta` creates it if needed |

Model metadata and configuration fields:

| Struct | Public fields |
| --- | --- |
| `ModelConfig` | `Provider`, `Model`, `APIKey`, `BaseURL`, `Temperature`, `MaxTokens`, `Timeout` |
| `ModelInfo` | `Provider`, `Model`, `MaxTokens`, `Temperature` |

## Provider Registry

Model providers implement:

```go
type ModelProvider interface {
    Name() string
    CreateModel(ctx context.Context, cfg *ModelConfig) (ToolableChatModel, error)
}
```

Register them once during application startup or package initialization:

```go
ai.RegisterModelProvider(openaiProvider)

model, err := ai.NewChatModel(ctx, &ai.ModelConfig{
    Provider: "openai",
    Model:    "gpt-4.1-mini",
    APIKey:   apiKey,
})
```

`RegisterModelProvider` panics on duplicate names with
`ai: model provider "{name}" already registered`; `RegisterAgentFactory` uses
`ai: agent factory "{name}" already registered`. `NewChatModel` returns
`ai.ErrProviderNotFound` wrapped with the requested provider name when no
provider is registered for `cfg.Provider`. `NewAgentBuilder` does the same with
`ai.ErrAgentNotFound` and the requested agent type. `cfg`, provider instances,
and factory instances are non-nil preconditions; the registry does not add nil
guards before calling `Name()` or reading `cfg.Provider`.

Agent factories follow the same pattern with `RegisterAgentFactory(...)`,
`NewAgentBuilder(...)`, `ListAgentFactories()`, and `ai.ErrAgentNotFound`.
Use `ListModelProviders()` to inspect registered model provider names. The
registry list helpers return map keys and do not define a stable sort order.

Factory implementations exchange `ModelConfig`, `ModelInfo`, `AgentConfig`,
and `AgentFactory`. These are plain public structs/interfaces for provider
packages; applications usually consume them only when registering their own
model or agent implementation.

## Tools

Tools expose JSON-schema-style metadata plus an invocation method:

```go
type Tool interface {
    Info() *ToolInfo
    Invoke(ctx context.Context, arguments string) (string, error)
}
```

`arguments` is a JSON-encoded string from the model. Tool implementations should
decode and validate it at the tool boundary, then return a string payload for
the model.

Tool metadata uses `ToolInfo`, `ParameterSchema`, and `PropertySchema`. The
schema shape is intentionally small: `type`, `properties`, `required`,
`description`, `enum`, and `items` cover the JSON-schema subset providers need
for tool calling.

Tool-related data shapes:

| Struct | Public fields |
| --- | --- |
| `ToolCall` | `ID`, `Name`, `Arguments` |
| `ToolResult` | `CallID`, `Content` |
| `ToolInfo` | `Name`, `Description`, `Parameters` |
| `ParameterSchema` | `Type`, `Properties`, `Required` |
| `PropertySchema` | `Type`, `Description`, `Enum`, `Items` |

For tools that stream output, implement `StreamableTool`:

```go
type StreamableTool interface {
    Tool
    InvokeStream(ctx context.Context, arguments string) (StringStream, error)
}
```

## Agents

An `Agent` is a higher-level runner that can reason with a model and tools:

```go
type Agent interface {
    Run(ctx context.Context, input string, opts ...Option) (*Message, error)
    Stream(ctx context.Context, input string, opts ...Option) (MessageStream, error)
}
```

`AgentBuilder` lets a registered factory expose fluent setup:

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

`AgentConfig` exposes the same builder inputs as fields: `Model`, `Tools`,
`SystemPrompt`, and `MaxIterations`.

## Streams

The base package defines two stream contracts:

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

`Recv` returns `io.EOF` when the stream is exhausted.

Stream data shapes:

| Struct | Public fields |
| --- | --- |
| `MessageChunk` | `Content`, `ToolCalls`, `Done` |
| `TokenUsage` | `PromptTokens`, `CompletionTokens`, `TotalTokens` |

## UI Message SSE Streams

The `ai/stream` package converts messages into Server-Sent Events compatible
with the AI SDK UI Message Stream protocol. It sets the standard SSE headers,
emits `data: ...` frames, and terminates with `data: [DONE]`.

Common adapters:

| Helper | Source |
| --- | --- |
| `stream.FromChannel(ch)` | a Go channel of `stream.Message` |
| `stream.FromCallback(fn)` | a callback that writes through `CallbackWriter` |
| `stream.FromAiMessageStream(s)` | an `ai.MessageStream` |
| `stream.NewChannelSource(ch)` | raw `MessageSource` adapter for a channel |
| `stream.NewCallbackSource(fn)` | raw callback `MessageSource` |
| `stream.NewAiMessageStreamSource(s)` | raw `ai.MessageStream` adapter |

Built-in adapter behavior:

| Adapter | Behavior |
| --- | --- |
| `NewChannelSource(ch)` | `Recv()` returns channel messages until the channel closes, then returns `io.EOF`; `Close()` marks the source closed and later `Recv()` calls return `io.EOF` without closing the caller-owned channel |
| `NewCallbackSource(fn)` | runs `fn` in a goroutine with a buffered message channel; queued messages are delivered before a callback error is returned; `Close()` waits for the callback goroutine to finish |
| `NewAiMessageStreamSource(s)` | converts each `ai.MessageChunk` into an assistant `stream.Message`, forwarding `Content` and any tool calls; `Close()` delegates to the wrapped `ai.MessageStream` |

Example Fiber handler:

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

The builder emits start/finish chunks by default, plus text, reasoning, tool
input, tool output, and custom `data-{type}` chunks when the source supplies
those message parts. Source and file chunks are exposed through the low-level
constructors for custom writers; `WithSources(bool)` only toggles the option
field and does not add sources to ordinary `MessageSource` output by itself.
`WithReasoning(bool)` gates reasoning output, and reasoning chunks are emitted
only when the incoming `stream.Message.Reasoning` is non-empty.

Default stream options:

| Option | Default |
| --- | --- |
| `SendReasoning` | `true` |
| `SendSources` | `true` |
| `SendStart` | `true` |
| `SendFinish` | `true` |
| `OnError` | returns `err.Error()` |
| `OnFinish` | `nil` |
| `GenerateID` | `nil` in `DefaultOptions()`; the builder falls back to `prefix + "_" + id.GenerateUUID()` |

Builder controls:

| Method | Purpose |
| --- | --- |
| `WithSource(source)` | set a custom `MessageSource` |
| `WithMessageID(id)` / `WithIDGenerator(fn)` | control generated chunk IDs |
| `WithStart(bool)` / `WithFinish(bool)` | toggle start/finish chunks |
| `WithReasoning(bool)` / `WithSources(bool)` | toggle reasoning output and the public source-output option |
| `WithHeader(key, value)` | add an SSE response header |
| `OnError(fn)` / `OnFinish(fn)` | customize error text or completion callback |
| `Stream(ctx)` / `StreamToWriter(w)` | write to Fiber or a buffered writer |

Execution-path contracts:

| Path | Behavior |
| --- | --- |
| `Stream(ctx)` without `WithSource(...)` | returns `stream.ErrSourceRequired` before writing headers |
| `Stream(ctx)` with a source | writes `SseHeaders`, then any `WithHeader` overrides/additions, then streams through `fiber.Ctx.SendStreamWriter` |
| `StreamToWriter(w)` | writes directly to the buffered writer and must be used only after configuring a source; unlike `Stream(ctx)`, it has no missing-source guard |
| source lifecycle | the configured source is closed when the stream loop exits |
| normal `io.EOF` from the source | closes open text/reasoning chunks, writes `finish-step` and `finish` when `SendFinish` is true, writes `data: [DONE]`, then calls `OnFinish(fullContent)` if configured |
| disabled start/finish | `WithStart(false)` omits both `start` and `start-step`; `WithFinish(false)` omits both `finish-step` and `finish`, but still writes `data: [DONE]` |
| non-EOF source error | writes an `error` chunk using `OnError(err)`, writes `data: [DONE]`, and does not call `OnFinish`; any open text/reasoning chunks are not explicitly ended first |
| `OnFinish(fullContent)` | receives the concatenated `Content` text chunks only, not reasoning, tool payloads, or custom data |
| tool call arguments | JSON-unmarshaled into the `input` payload; invalid JSON is sent as the original string; an empty tool call ID is replaced with a generated `call_*` ID |
| tool result content | JSON-unmarshaled into the `output` payload; invalid JSON is sent as the original string |
| custom data messages | each `Data` map entry becomes a `data-{type}` chunk; map iteration order is not stable when one message contains multiple data keys |
| writer errors | chunk/write errors inside the stream loop are ignored because `StreamToWriter` and the Fiber stream callback do not return an error from the loop |

Default SSE headers:

| Header | Value |
| --- | --- |
| `Content-Type` | `text/event-stream` |
| `Cache-Control` | `no-cache` |
| `Connection` | `keep-alive` |
| `Transfer-Encoding` | `chunked` |
| `X-Vercel-AI-UI-Message-Stream` | `v1` |
| `X-Accel-Buffering` | `no` |

Low-level chunk constructors are public for tests and custom writers:

| Chunk family | Constructors |
| --- | --- |
| lifecycle | `NewStartChunk`, `NewFinishChunk`, `NewStartStepChunk`, `NewFinishStepChunk`, `NewErrorChunk` |
| text | `NewTextStartChunk`, `NewTextDeltaChunk`, `NewTextEndChunk` |
| reasoning | `NewReasoningStartChunk`, `NewReasoningDeltaChunk`, `NewReasoningEndChunk` |
| tools | `NewToolInputStartChunk`, `NewToolInputDeltaChunk`, `NewToolInputAvailableChunk`, `NewToolOutputAvailableChunk` |
| sources/files/data | `NewSourceURLChunk`, `NewSourceDocumentChunk`, `NewFileChunk`, `NewDataChunk` |

Exact chunk wire shapes:

| Constructor | Wire `type` | Fields |
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
| `NewSourceURLChunk(sourceID, url, title)` | `source-url` | `type`, `sourceID`, `url`; `title` is included only when non-empty |
| `NewSourceDocumentChunk(sourceID, mediaType, title)` | `source-document` | `type`, `sourceID`, `mediaType`; `title` is included only when non-empty |
| `NewFileChunk(fileID, mediaType, url)` | `file` | `type`, `fileID`, `mediaType`, `url` |
| `NewDataChunk(dataType, data)` | `data-{dataType}` | `type`, `data`; `data-{dataType}` is dynamic and is not a `ChunkType` constant |

`CallbackWriter` writes these source messages:

| Method | Message effect |
| --- | --- |
| `WriteText(content)` | assistant message with `Content` |
| `WriteToolCall(id, name, arguments)` | assistant message with one `ToolCall` |
| `WriteToolResult(toolCallID, content)` | tool message with `ToolCallID` and `Content` |
| `WriteReasoning(reasoning)` | assistant message with `Reasoning` |
| `WriteData(dataType, data)` | assistant message with `Data[dataType]` |
| `WriteMessage(msg)` | forwards the provided `stream.Message` unchanged |

Stream-package data shapes:

| Struct or interface | Public fields or methods |
| --- | --- |
| `stream.Message` | `Role`, `Content`, `ToolCalls`, `ToolCallID`, `Reasoning`, `Data` |
| `stream.ToolCall` | `ID`, `Name`, `Arguments` |
| `stream.Source` | `Type`, `ID`, `URL`, `Title`, `MediaType` |
| `stream.Options` | `SendReasoning`, `SendSources`, `SendStart`, `SendFinish`, `OnError`, `OnFinish`, `GenerateID` |
| `MessageSource` | `Recv()`, `Close()`; `Recv()` returns `io.EOF` when complete |
| `StreamWriter` | `WriteChunk(chunk)`, `Flush()` |
| `ResponseWriter` | `io.Writer` compatibility for Fiber stream writers |

`stream.SseHeaders` contains the default AI SDK UI Message Stream headers.
`stream.ErrSourceRequired` is returned when a builder has no source. The
`stream.ErrSourceClosed` sentinel is public for custom `MessageSource`
implementations; the built-in channel and callback adapters report normal
completion as `io.EOF`.

Supporting stream APIs:

| API group | Public surface |
| --- | --- |
| roles | `stream.RoleUser`, `stream.RoleAssistant`, `stream.RoleTool`, `stream.RoleSystem` |
| chunk type constants | `ChunkType`, `ChunkTypeStart`, `ChunkTypeFinish`, `ChunkTypeStartStep`, `ChunkTypeFinishStep`, `ChunkTypeError`, `ChunkTypeTextStart`, `ChunkTypeTextDelta`, `ChunkTypeTextEnd`, `ChunkTypeReasoningStart`, `ChunkTypeReasoningDelta`, `ChunkTypeReasoningEnd`, `ChunkTypeToolInputStart`, `ChunkTypeToolInputDelta`, `ChunkTypeToolInputAvailable`, `ChunkTypeToolOutputAvailable`, `ChunkTypeSourceURL`, `ChunkTypeSourceDocument`, `ChunkTypeFile` |
| stream data structs | `stream.Message`, `stream.ToolCall`, `stream.Source`, `stream.Chunk`, `stream.Options` |
| writer/source interfaces | `CallbackWriter`, `MessageSource`, `ResponseWriter`, `StreamWriter` |
| defaults | `stream.New()`, `DefaultOptions()`, `SseHeaders` |

## Error Sentinels

| Error | Meaning |
| --- | --- |
| `ai.ErrProviderNotFound` | no registered model provider matched `ModelConfig.Provider` |
| `ai.ErrModelNotSupported` | provider does not support the requested model |
| `ai.ErrAgentNotFound` | no registered agent factory matched the requested type |
| `ai.ErrToolNotFound` | an agent/model requested an unavailable tool |
| `ai.ErrInvalidArguments` | tool arguments failed validation |
| `ai.ErrMaxIterationsReached` | agent exceeded its iteration limit |
| `ai.ErrNoContent` | model returned no usable content |
| `ai.ErrStreamClosed` | stream read after close |
| `stream.ErrSourceRequired` | `ai/stream` builder has no message source |
| `stream.ErrSourceClosed` | public sentinel for custom `MessageSource` implementations |

`NewToolError(...)` and `NewModelError(...)` create structured `ToolError` and
`ModelError` values for tool and provider failures. `ToolError.Unwrap()`
exposes the underlying tool error.

Structured error shapes:

| Type | Public fields and behavior |
| --- | --- |
| `ToolError` | `ToolName`, `Err`; `Error()` returns `ai: tool {ToolName}: {Err}` and `Unwrap()` returns `Err` |
| `ModelError` | `Provider`, `StatusCode`, `Message`; `Error()` returns `ai: {Provider}: {Message}` |
