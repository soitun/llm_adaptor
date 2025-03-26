package define

const ApiVersionV1 = `v1`
const ApiVersionV2 = `v2`

type (
	CommonChatCompletionChoiceRes struct {
		Message ChatCompletionResponseMessage `json:"message,omitempty"`
		Delta   ChatCompletionResponseMessage `json:"delta,omitempty"`
	}
	ChatCompletionResponseMessage struct {
		Role             string     `json:"role,omitempty"`
		Content          string     `json:"content"`
		ReasoningContent string     `json:"reasoning_content"`
		ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	}
	ToolCall struct {
		Id       string   `json:"id"`
		Type     string   `json:"type"`
		Function Function `json:"function"`
	}
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	}
)
