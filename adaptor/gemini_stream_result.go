// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import "github.com/zhimaAi/llm_adaptor/api/gemini"

type GeminiStreamResult struct {
	*gemini.ChatCompletionStream
}

func (c *GeminiStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	responseGemini, err := c.ChatCompletionStream.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	return ZhimaChatCompletionResponse{
		Result:          responseGemini.Candidates[0].Content.Parts[0].Text,
		PromptToken:     responseGemini.UsageMetadata.PromptTokenCount,
		CompletionToken: responseGemini.UsageMetadata.CandidatesTokenCount,
	}, nil
}
