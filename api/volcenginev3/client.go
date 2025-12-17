// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package volcenginev3

import (
	"github.com/zhimaAi/llm_adaptor/api/openai"
	"github.com/zhimaAi/llm_adaptor/common"
	"github.com/zhimaAi/llm_adaptor/define"
)

type Client struct {
	Host       string
	Model      string
	AK         string
	SK         string
	Region     string
	ApiVersion string
}

func NewClient(Host, Model, AK, SK, Region string) *Client {
	if SK == "" {
		return &Client{
			Host:       Host,
			Model:      Model,
			AK:         AK,
			Region:     Region,
			ApiVersion: define.ApiVersionV2,
		}
	}
	return &Client{
		Host:       Host,
		Model:      Model,
		AK:         AK,
		SK:         SK,
		Region:     Region,
		ApiVersion: define.ApiVersionV1,
	}
}

func getAccessToken(Region, Model, AK, SK string) (string, error) {
	tokenManager := common.GetTokenManagerInstance()
	return tokenManager.GetVolcengineAccessToken("https://open.volcengineapi.com", Region, Model, AK, SK)
}

func (c *Client) CreateChatCompletion(req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	var (
		err         error
		accessToken = c.AK
	)
	if c.ApiVersion == define.ApiVersionV1 {
		accessToken, err = getAccessToken(c.Region, c.Model, c.AK, c.SK)
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
	}

	OpenAIClient := openai.NewClient(c.Host, accessToken, &openai.ErrorResponse{})
	return OpenAIClient.CreateChatCompletion(req)
}

func (c *Client) CreateChatCompletionStream(req openai.ChatCompletionRequest) (*openai.ChatCompletionStream, error) {
	var (
		err         error
		accessToken = c.AK
	)
	if c.ApiVersion == define.ApiVersionV1 {
		accessToken, err = getAccessToken(c.Region, c.Model, c.AK, c.SK)
		if err != nil {
			return nil, err
		}
	}
	OpenAIClient := openai.NewClient(c.Host, accessToken, &openai.ErrorResponse{})
	return OpenAIClient.CreateChatCompletionStream(req)
}
