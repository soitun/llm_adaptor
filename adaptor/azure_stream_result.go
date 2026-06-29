// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"github.com/zhimaAi/llm_adaptor/api/azure"
	"github.com/zhimaAi/llm_adaptor/basics"
)

type AzureStreamResult struct {
	*azure.ChatCompletionStream
}

func (r *AzureStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseAzure, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	var result string
	var toolCalls basics.ToolCalls
	if len(responseAzure.Choices) > 0 {
		result = responseAzure.Choices[0].Delta.Content
		toolCalls = responseAzure.Choices[0].Delta.ToolCalls
	}
	return ZhimaChatCompletionResponse{
		Result:            result,
		ToolCalls:         toolCalls,
		FunctionToolCalls: toolCalls.FunctionToolCalls(),
		PromptToken:       responseAzure.Usage.PromptTokens,
		CompletionToken:   responseAzure.Usage.CompletionTokens,
	}, nil
}
