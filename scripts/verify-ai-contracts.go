package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type corpus struct {
	label   string
	content string
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	checks := []struct {
		sourcePath      string
		sourceTerms     []string
		docTerms        []string
		englishDocTerms []string
		chineseDocTerms []string
	}{
		{
			sourcePath: "ai/agent.go",
			sourceTerms: []string{
				"Agent interface", "Run(ctx context.Context, input string, opts ...Option)",
				"Stream(ctx context.Context, input string, opts ...Option)", "AgentConfig",
				"Model ToolableChatModel", "Tools []Tool", "SystemPrompt string",
				"MaxIterations int", "AgentBuilder interface", "WithModel(model ToolableChatModel)",
				"WithTools(tools ...Tool)", "WithSystemPrompt(prompt string)",
				"WithMaxIterations(n int)", "Build(ctx context.Context)",
			},
			docTerms: []string{
				"type Agent interface", "Run(ctx context.Context, input string, opts ...Option)",
				"Stream(ctx context.Context, input string, opts ...Option)", "AgentBuilder",
				"WithModel(model)", "WithTools(searchTool, calculatorTool)",
				"WithSystemPrompt(\"Use tools when needed.\")", "WithMaxIterations(8)",
				"Build(ctx)", "AgentConfig", "Model", "Tools", "SystemPrompt",
				"MaxIterations",
			},
		},
		{
			sourcePath: "ai/message.go",
			sourceTerms: []string{
				"RoleSystem Role = \"system\"", "RoleUser Role = \"user\"",
				"RoleAssistant Role = \"assistant\"", "RoleTool Role = \"tool\"",
				"ToolCall struct", "ID string", "Name string", "Arguments string",
				"ToolResult struct", "CallID string", "TokenUsage struct",
				"PromptTokens int", "CompletionTokens int", "TotalTokens int",
				"Message struct", "ToolCalls []ToolCall", "ToolResult *ToolResult",
				"Usage *TokenUsage", "NewSystemMessage", "NewUserMessage",
				"NewAssistantMessage", "NewAssistantMessageWithToolCalls",
				"NewToolMessage", "IsSystem", "IsUser", "IsAssistant", "IsTool",
				"HasToolCalls", "Role:    RoleSystem", "Role:    RoleUser",
				"Role:    RoleAssistant", "ToolResult: &ToolResult",
				"CallID:  callID", "Content: content",
			},
			docTerms: []string{
				"type Message struct", "RoleSystem", "RoleUser", "RoleAssistant",
				"RoleTool", "NewSystemMessage", "NewUserMessage",
				"NewAssistantMessageWithToolCalls", "NewToolMessage(callID, content)",
				"NewAssistantMessage(...)", "Message.IsSystem", "IsUser",
				"IsAssistant", "IsTool", "HasToolCalls", "ToolCall",
				"ID", "Name", "Arguments", "ToolResult", "CallID", "Content",
				"TokenUsage", "PromptTokens", "CompletionTokens", "TotalTokens",
				"NewSystemMessage(content)", "Role: RoleSystem",
				"NewUserMessage(content)", "Role: RoleUser",
				"NewAssistantMessage(content)", "Role: RoleAssistant",
				"NewAssistantMessageWithToolCalls(content, toolCalls)",
				"ToolCalls: toolCalls",
				"ToolResult: &ToolResult{CallID: callID, Content: content}",
				"Message.Content",
			},
		},
		{
			sourcePath: "ai/model.go",
			sourceTerms: []string{
				"ChatModel interface", "Generate(ctx context.Context, messages []*Message, opts ...Option)",
				"Stream(ctx context.Context, messages []*Message, opts ...Option)",
				"ToolableChatModel interface", "WithTools(tools ...Tool) ToolableChatModel",
				"ModelInfo struct", "Provider string", "Model string", "MaxTokens int",
				"Temperature float64",
			},
			docTerms: []string{
				"type ChatModel interface", "Generate(ctx context.Context, messages []*Message, opts ...Option)",
				"Stream(ctx context.Context, messages []*Message, opts ...Option)",
				"type ToolableChatModel interface", "WithTools(tools ...Tool) ToolableChatModel",
				"ModelInfo", "Provider", "Model", "MaxTokens", "Temperature",
			},
			englishDocTerms: []string{"immutable pattern"},
			chineseDocTerms: []string{"不可变模式"},
		},
		{
			sourcePath: "ai/option.go",
			sourceTerms: []string{
				"Option func(*Options)", "Options struct", "Temperature *float64",
				"MaxTokens *int", "StopSequences []string", "Meta map[string]string",
				"NewOptions", "Meta: make(map[string]string)", "Apply(opts ...Option)",
				"WithTemperature", "WithMaxTokens", "WithStopSequences", "WithMeta",
				"if o.Meta == nil", "for _, opt := range opts", "opt(o)",
				"return o",
			},
			docTerms: []string{
				"NewOptions()", "Options.Apply(...)", "WithTemperature(t)",
				"WithMaxTokens(n)", "WithStopSequences(seqs...)", "WithMeta(key, value)",
				"Options.Temperature", "Options.MaxTokens", "Options.StopSequences",
				"Options.Meta",
			},
			englishDocTerms: []string{
				"NewOptions()` initializes it", "WithMeta` creates it if needed",
				"applies options in argument", "later options can overwrite",
			},
			chineseDocTerms: []string{
				"NewOptions()` 会初始化它", "WithMeta` 在需要时也会创建它",
				"按参数顺序应用 options", "可以覆盖之前 option 设置的字段",
			},
		},
		{
			sourcePath: "ai/provider.go",
			sourceTerms: []string{
				"ModelProvider interface", "Name() string", "CreateModel(ctx context.Context, cfg *ModelConfig)",
				"AgentFactory interface", "CreateBuilder() AgentBuilder", "ModelConfig struct",
				"Provider string", "Model string", "APIKey string", "BaseURL string",
				"Temperature float64", "MaxTokens int", "Timeout time.Duration",
				"RegisterModelProvider", "RegisterAgentFactory", "panic(fmt.Sprintf",
				"NewChatModel", "ErrProviderNotFound", "NewAgentBuilder",
				"ErrAgentNotFound", "ListModelProviders", "ListAgentFactories",
				"maps.Keys", "ai: model provider %q already registered",
				"ai: agent factory %q already registered", "modelProviders[cfg.Provider]",
			},
			docTerms: []string{
				"type ModelProvider interface", "Name() string",
				"CreateModel(ctx context.Context, cfg *ModelConfig)", "RegisterModelProvider",
				"NewChatModel", "ErrProviderNotFound", "RegisterAgentFactory",
				"NewAgentBuilder", "ListAgentFactories()", "ErrAgentNotFound",
				"ListModelProviders()", "ModelConfig", "Provider", "Model", "APIKey", "BaseURL",
				"Temperature", "MaxTokens", "Timeout", "AgentFactory",
				"ai: model provider", "already registered", "ai: agent factory",
				"cfg.Provider", "Name()",
			},
			englishDocTerms: []string{
				"do not define a stable sort order", "non-nil preconditions",
				"does not add nil",
			},
			chineseDocTerms: []string{
				"不定义稳定排序", "非 nil 前置条件", "不会额外做 nil guard",
			},
		},
		{
			sourcePath: "ai/tool.go",
			sourceTerms: []string{
				"ToolInfo struct", "Name string", "Description string",
				"Parameters *ParameterSchema", "ParameterSchema struct",
				"Properties map[string]*PropertySchema", "Required []string",
				"PropertySchema struct", "Enum []string", "Items *PropertySchema",
				"Tool interface", "Info() *ToolInfo", "Invoke(ctx context.Context, arguments string)",
				"StreamableTool interface", "InvokeStream(ctx context.Context, arguments string)",
			},
			docTerms: []string{
				"type Tool interface", "Info() *ToolInfo",
				"Invoke(ctx context.Context, arguments string)", "JSON-encoded string",
				"ToolInfo", "ParameterSchema", "PropertySchema", "type",
				"properties", "required", "description", "enum", "items",
				"StreamableTool", "InvokeStream(ctx context.Context, arguments string)",
			},
		},
		{
			sourcePath: "ai/stream.go",
			sourceTerms: []string{
				"MessageStream interface", "io.Closer", "Recv() (*MessageChunk, error)",
				"Collect() (*Message, error)", "MessageChunk struct", "Content string",
				"ToolCalls []ToolCall", "Done bool", "StringStream interface",
				"Recv() (string, error)", "Collect() (string, error)",
			},
			docTerms: []string{
				"type MessageStream interface", "io.Closer",
				"Recv() (*MessageChunk, error)", "Collect() (*Message, error)",
				"type StringStream interface", "Recv() (string, error)",
				"Collect() (string, error)", "io.EOF", "MessageChunk",
				"Content", "ToolCalls", "Done",
			},
		},
		{
			sourcePath: "ai/stream/chunk.go",
			sourceTerms: []string{
				"NewStartChunk", "messageID", "NewFinishChunk", "NewStartStepChunk",
				"NewFinishStepChunk", "NewErrorChunk", "errorText", "NewTextStartChunk",
				"NewTextDeltaChunk", "NewTextEndChunk", "NewReasoningStartChunk",
				"NewReasoningDeltaChunk", "NewReasoningEndChunk", "NewToolInputStartChunk",
				"toolCallID", "toolName", "NewToolInputDeltaChunk", "inputTextDelta",
				"NewToolInputAvailableChunk", "input", "NewToolOutputAvailableChunk",
				"output", "NewSourceURLChunk", "sourceID", "NewSourceDocumentChunk",
				"mediaType", "NewFileChunk", "fileID", "NewDataChunk", "data-",
			},
			docTerms: []string{
				"NewStartChunk(messageID)", "messageID", "NewFinishChunk()",
				"NewStartStepChunk()", "NewFinishStepChunk()", "NewErrorChunk(errorText)",
				"errorText", "NewTextStartChunk(id)", "NewTextDeltaChunk(id, delta)",
				"NewTextEndChunk(id)", "NewReasoningStartChunk(id)",
				"NewReasoningDeltaChunk(id, delta)", "NewReasoningEndChunk(id)",
				"NewToolInputStartChunk(toolCallID, toolName)", "toolCallID",
				"toolName", "NewToolInputDeltaChunk(toolCallID, delta)",
				"inputTextDelta", "NewToolInputAvailableChunk(toolCallID, toolName, input)",
				"input", "NewToolOutputAvailableChunk(toolCallID, output)", "output",
				"NewSourceURLChunk(sourceID, url, title)", "sourceID",
				"NewSourceDocumentChunk(sourceID, mediaType, title)", "mediaType",
				"NewFileChunk(fileID, mediaType, url)", "fileID",
				"NewDataChunk(dataType, data)", "data-{dataType}",
				"ChunkType",
			},
		},
		{
			sourcePath: "ai/stream/writer.go",
			sourceTerms: []string{
				"data: ", "\\n\\n", "data: [DONE]", "fiber.HeaderContentType",
				"text/event-stream", "fiber.HeaderCacheControl", "no-cache",
				"fiber.HeaderConnection", "keep-alive", "fiber.HeaderTransferEncoding", "chunked",
				"X-Vercel-AI-UI-Message-Stream", "v1", "X-Accel-Buffering", "no",
				"prefix + \"_\" + id.GenerateUUID()",
			},
			docTerms: []string{
				"data: ...", "data: [DONE]", "Content-Type", "text/event-stream",
				"Cache-Control", "no-cache", "Connection", "keep-alive",
				"Transfer-Encoding", "chunked", "X-Vercel-AI-UI-Message-Stream",
				"v1", "X-Accel-Buffering", "prefix + \"_\" + id.GenerateUUID()",
			},
		},
		{
			sourcePath: "ai/stream/options.go",
			sourceTerms: []string{
				"SendReasoning: true", "SendSources:   true", "SendStart:     true",
				"SendFinish:    true", "return err.Error()",
			},
			docTerms: []string{
				"SendReasoning", "SendSources", "SendStart", "SendFinish", "err.Error()",
			},
		},
		{
			sourcePath: "ai/stream/builder.go",
			sourceTerms: []string{
				"ErrSourceRequired", "WithHeader", "StreamToWriter",
				"OnFinish(fullContent)", "NewErrorChunk(errorText)", "json.Unmarshal",
				"input = tc.Arguments", "output = msg.Content",
				"defer func() { _ = b.source.Close() }()", "if toolCallID == \"\"",
				"toolCallID = generateID(\"call\")", "if b.opts.SendFinish",
				"writer.writeDone", "fullContent += msg.Content",
				"_ = writer.WriteChunk", "_ = writer.writeDone()",
			},
			docTerms: []string{
				"stream.ErrSourceRequired", "WithHeader", "StreamToWriter(w)",
				"OnFinish(fullContent)", "NewErrorChunk(errorText)", "JSON",
				"WithStart(false)", "WithFinish(false)", "data: [DONE]",
				"call_*", "source lifecycle",
			},
			englishDocTerms: []string{
				"original string", "has no missing-source guard",
				"configured source is closed", "does not call `OnFinish`",
				"writer errors",
			},
			chineseDocTerms: []string{
				"原始 string", "没有缺失 source 的 guard",
				"关闭已配置的 source", "不会调用 `OnFinish`",
				"writer errors",
			},
		},
		{
			sourcePath: "ai/stream/adapters.go",
			sourceTerms: []string{
				"RoleUser      Role = \"user\"", "RoleAssistant Role = \"assistant\"",
				"RoleTool      Role = \"tool\"", "RoleSystem    Role = \"system\"",
				"Message struct", "ToolCallID string", "Reasoning",
				"Data       map[string]any", "ToolCall struct", "MessageSource interface",
				"Recv() (Message, error)", "Close() error", "NewChannelSource",
				"FromChannel", "NewCallbackSource", "FromCallback",
				"WriteText(content string)", "WriteToolCall(id, name, arguments string)",
				"WriteToolResult(toolCallID, content string)", "WriteReasoning(reasoning string)",
				"WriteData(dataType string, data any)", "WriteMessage(msg Message)",
				"NewAiMessageStreamSource", "FromAiMessageStream",
				"messages: make(chan Message, 16)", "defer close(s.messages)",
				"if c.err != nil", "<-c.done", "return a.stream.Close()",
				"Role:    RoleAssistant", "Content: chunk.Content",
			},
			docTerms: []string{
				"stream.FromChannel(ch)", "stream.FromCallback(fn)",
				"stream.FromAiMessageStream(s)", "stream.NewChannelSource(ch)",
				"stream.NewCallbackSource(fn)", "stream.NewAiMessageStreamSource(s)",
				"CallbackWriter", "WriteText(content)", "WriteToolCall(id, name, arguments)",
				"WriteToolResult(toolCallID, content)", "WriteReasoning(reasoning)",
				"WriteData(dataType, data)", "WriteMessage(msg)", "stream.Message",
				"Role", "Content", "ToolCalls", "ToolCallID", "Reasoning", "Data",
				"stream.ToolCall", "ID", "Name", "Arguments", "MessageSource",
				"Recv()", "Close()", "stream.RoleUser", "stream.RoleAssistant",
				"stream.RoleTool", "stream.RoleSystem", "io.EOF",
				"ai.MessageChunk", "Close()",
			},
			englishDocTerms: []string{
				"queued messages are delivered before a callback error",
				"Close()` waits for the callback goroutine to finish",
				"converts each `ai.MessageChunk` into an assistant `stream.Message`",
				"caller-owned channel",
			},
			chineseDocTerms: []string{
				"已排队的消息会先交付，然后才返回 callback error",
				"Close()` 会等待 callback goroutine 结束",
				"把每个 `ai.MessageChunk` 转成 assistant `stream.Message`",
				"调用方持有的 channel",
			},
		},
		{
			sourcePath: "ai/stream/errors.go",
			sourceTerms: []string{
				"ErrSourceRequired", "message source is required",
				"ErrSourceClosed", "message source is closed",
			},
			docTerms: []string{
				"stream.ErrSourceRequired", "stream.ErrSourceClosed",
				"message source",
			},
			englishDocTerms: []string{"custom `MessageSource` implementations"},
			chineseDocTerms: []string{"自定义 `MessageSource` 实现"},
		},
		{
			sourcePath: "ai/errors.go",
			sourceTerms: []string{
				"ErrStreamClosed", "ErrNoContent", "ErrMaxIterationsReached",
				"ErrToolNotFound", "ErrInvalidArguments", "ErrProviderNotFound",
				"ErrModelNotSupported", "ErrAgentNotFound", "ToolError",
				"ToolName", "Unwrap", "ModelError", "StatusCode",
			},
			docTerms: []string{
				"ai.ErrStreamClosed", "ai.ErrNoContent", "ai.ErrMaxIterationsReached",
				"ai.ErrToolNotFound", "ai.ErrInvalidArguments", "ai.ErrProviderNotFound",
				"ai.ErrModelNotSupported", "ai.ErrAgentNotFound", "ToolError",
				"ToolName", "Unwrap()", "ModelError", "StatusCode",
			},
		},
	}

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)
	englishDocs := readCorpus("English AI docs", filepath.Join(docsRoot, "docs/features/ai.md"))
	chineseDocs := readCorpus("Chinese AI docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/ai.md"))
	docs := []corpus{englishDocs, chineseDocs}

	var failures []string
	for _, check := range checks {
		source := readCorpus(check.sourcePath, filepath.Join(sourceRoot, check.sourcePath))
		failures = append(failures, missingTerms(source, check.sourceTerms)...)
		for _, doc := range docs {
			failures = append(failures, missingTerms(doc, check.docTerms)...)
		}
		failures = append(failures, missingTerms(englishDocs, check.englishDocTerms)...)
		failures = append(failures, missingTerms(chineseDocs, check.chineseDocTerms)...)
	}

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("AI contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("AI contract docs verified: %d source files, %d doc mirrors\n", len(checks), len(docs))
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func missingTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !strings.Contains(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term: %s", c.label, term))
		}
	}

	return failures
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
