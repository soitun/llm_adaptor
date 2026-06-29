// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"errors"
	"io"

	"github.com/zhimaAi/llm_adaptor/api/spark"
	"github.com/zhimaAi/llm_adaptor/basics"
)

type SparkStreamResult struct {
	*spark.ChatCompletionStream
}

func (r *SparkStreamResult) Read() (resp ZhimaChatCompletionResponse, err error) {
	responseSpark, err := r.Recv()
	var toolCalls basics.ToolCalls
	if err != nil {
		if errors.Is(err, io.EOF) {
			if len(responseSpark.Payload.Choices.Text[0].FunctionCall.Name) > 0 {
				toolCalls = append(toolCalls, basics.NewFunctionToolCall("", responseSpark.Payload.Choices.Text[0].FunctionCall.Name, responseSpark.Payload.Choices.Text[0].FunctionCall.Arguments))
			}
			resp = ZhimaChatCompletionResponse{
				Result:            responseSpark.Payload.Choices.Text[0].Content,
				ToolCalls:         toolCalls,
				FunctionToolCalls: toolCalls.FunctionToolCalls(),
				PromptToken:       responseSpark.Payload.Usage.Text.PromptTokens,
				CompletionToken:   responseSpark.Payload.Usage.Text.CompletionTokens,
			}
		}
	} else {
		if len(responseSpark.Payload.Choices.Text[0].FunctionCall.Name) > 0 {
			toolCalls = append(toolCalls, basics.NewFunctionToolCall("", responseSpark.Payload.Choices.Text[0].FunctionCall.Name, responseSpark.Payload.Choices.Text[0].FunctionCall.Arguments))
		}
		resp = ZhimaChatCompletionResponse{
			Result:            responseSpark.Payload.Choices.Text[0].Content,
			ToolCalls:         toolCalls,
			FunctionToolCalls: toolCalls.FunctionToolCalls(),
			PromptToken:       responseSpark.Payload.Usage.Text.PromptTokens,
			CompletionToken:   responseSpark.Payload.Usage.Text.CompletionTokens,
		}
	}

	return
}
