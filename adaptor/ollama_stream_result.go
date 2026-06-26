// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"github.com/zhimaAi/llm_adaptor/api/ollama"
)

type OllamaStreamResult struct {
	*ollama.ChatCompletionStream
}

func (r *OllamaStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseOllama, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	toolCalls := responseOllama.Message.ToolCalls
	return ZhimaChatCompletionResponse{
		Result:            responseOllama.Message.Content,
		ReasoningContent:  responseOllama.Message.ReasoningContent,
		ToolCalls:         toolCalls,
		FunctionToolCalls: toolCalls.FunctionToolCalls(),
		PromptToken:       responseOllama.Metrics.PromptEvalCount,
		CompletionToken:   responseOllama.Metrics.EvalCount,
	}, nil
}
