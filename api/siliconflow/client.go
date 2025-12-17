// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package siliconflow

import (
	"github.com/zhimaAi/llm_adaptor/api/openai"
	"github.com/zhimaAi/llm_adaptor/common"
)

type Client struct {
	APIKey       string
	EndPoint     string
	ErrResp      common.ErrorResponseInterface
	OpenAIClient *openai.Client
}

func NewClient(apiEndpoint, APIKey, apiVersion string) *Client {
	return &Client{
		EndPoint: apiEndpoint,
		APIKey:   APIKey,
		OpenAIClient: &openai.Client{
			EndPoint: apiEndpoint + "/" + apiVersion,
			APIKey:   APIKey,
			ErrResp:  &openai.ErrorResponse{},
		},
	}
}
