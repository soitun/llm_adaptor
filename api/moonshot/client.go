// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package moonshot

import "github.com/zhimaAi/llm_adaptor/api/openai"

type Client struct {
	APIKey       string
	EndPoint     string
	OpenAIClient *openai.Client
}

func NewClient(APIKey string) *Client {
	return &Client{
		APIKey:   APIKey,
		EndPoint: "https://api.moonshot.cn/v1",
		OpenAIClient: &openai.Client{
			EndPoint: "https://api.moonshot.cn/v1",
			APIKey:   APIKey,
			ErrResp:  &ErrorResponse{},
		},
	}
}
