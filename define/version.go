// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import "github.com/zhimaAi/llm_adaptor/basics"

const ApiVersionV1 = `v1`
const ApiVersionV2 = `v2`

type (
	CommonChatCompletionChoiceRes struct {
		Message ChatCompletionResponseMessage `json:"message,omitempty"`
		Delta   ChatCompletionResponseMessage `json:"delta,omitempty"`
	}
	ChatCompletionResponseMessage struct {
		Role             string           `json:"role,omitempty"`
		Content          string           `json:"content"`
		ReasoningContent string           `json:"reasoning_content"`
		ToolCalls        basics.ToolCalls `json:"tool_calls,omitempty"`
	}
)
