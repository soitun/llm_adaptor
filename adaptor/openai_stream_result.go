// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"github.com/zhimaAi/llm_adaptor/api/openai"
)

type OpenAIStreamResult struct {
	*openai.ChatCompletionStream
}

func (r *OpenAIStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseOpenAI, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}

	var promptTokens int
	var completionTokens int
	var result, reasoningContent = "", ""
	if responseOpenAI.Usage.PromptTokens > 0 {
		promptTokens = responseOpenAI.Usage.PromptTokens
	}
	if responseOpenAI.Usage.CompletionTokens > 0 {
		completionTokens = responseOpenAI.Usage.CompletionTokens
	}
	var functionToolCalls []FunctionToolCall
	if len(responseOpenAI.Choices) > 0 {
		result = responseOpenAI.Choices[0].Delta.Content
		reasoningContent = responseOpenAI.Choices[0].Delta.ReasoningContent
		// Compatible with moonlight
		if responseOpenAI.Choices[0].Usage.PromptTokens > 0 {
			promptTokens = responseOpenAI.Choices[0].Usage.PromptTokens
		}
		if responseOpenAI.Choices[0].Usage.CompletionTokens > 0 {
			completionTokens = responseOpenAI.Choices[0].Usage.CompletionTokens
		}
		for _, toolCall := range responseOpenAI.Choices[0].Delta.ToolCalls {
			functionToolCalls = append(functionToolCalls, FunctionToolCall{
				Name:      toolCall.Function.Name,
				Arguments: toolCall.Function.Arguments,
			})
		}
	}

	return ZhimaChatCompletionResponse{
		Result:            result,
		ReasoningContent:  reasoningContent,
		FunctionToolCalls: functionToolCalls,
		PromptToken:       promptTokens,
		CompletionToken:   completionTokens,
	}, nil
}

type OpenAIImageGenerationStreamResult struct {
	*openai.ImageGenerationStream
	Ext string
}

func (r *OpenAIImageGenerationStreamResult) Read() (ZhimaImageGenerationResp, error) {
	res, err := r.Recv()
	if err != nil {
		return ZhimaImageGenerationResp{}, err
	}
	inputToken := res.Usage.TotalTokens - res.Usage.OutputTokens
	outputToken := res.Usage.OutputTokens
	datas := make([]*ImageGenerationData, 0)
	if res.Type == `image_generation.completed` {
		//
	} else if res.Type == `image_generation.partial_failed` {
		datas = append(datas, &ImageGenerationData{
			Error: DataError{
				Code:    res.Error.Code,
				Message: res.Error.Message,
			},
		})
	} else if res.Type == `image_generation.partial_succeeded` {
		datas = append(datas, &ImageGenerationData{
			Url:     res.Url,
			B64Json: res.B64Json,
			Size:    res.Size,
			Error:   DataError{},
			Ext:     r.Ext,
		})
	}
	return ZhimaImageGenerationResp{
		InputToken:  inputToken,
		OutputToken: outputToken,
		Datas:       datas,
	}, nil
}
