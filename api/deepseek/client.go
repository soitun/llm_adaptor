// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package deepseek

import "github.com/zhimaAi/llm_adaptor/api/openai"

type Client struct {
	EndPoint     string
	APIKey       string
	OpenAIClient *openai.Client
}

func NewClient(APIKey string) *Client {
	return &Client{
		EndPoint: "https://api.deepseek.com",
		APIKey:   APIKey,
		OpenAIClient: &openai.Client{
			EndPoint: "https://api.deepseek.com",
			APIKey:   APIKey,
			ErrResp:  &openai.ErrorResponse{},
		},
	}
}
