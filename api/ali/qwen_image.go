// Copyright © 2016- 2026 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package ali

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zhimaAi/llm_adaptor/common"
)

type QwenImageGenerationRequest struct {
	Model      string                 `json:"model"`
	Input      QwenImageInput         `json:"input"`
	Parameters *QwenImageGenParameter `json:"parameters,omitempty"`
}

type QwenImageInput struct {
	Messages []QwenImageMessage `json:"messages"`
}

type QwenImageMessage struct {
	Role    string             `json:"role"`
	Content []QwenImageContent `json:"content"`
}

type QwenImageContent struct {
	Text string `json:"text"`
}

type QwenImageGenParameter struct {
	Size         *string `json:"size,omitempty"`
	PromptExtend *bool   `json:"prompt_extend,omitempty"`
	Watermark    *bool   `json:"watermark,omitempty"`
	Seed         *int    `json:"seed,omitempty"`
}

type QwenImageGenerationResponse struct {
	RequestID string `json:"request_id"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`

	Output struct {
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				Role    string `json:"role"`
				Content []struct {
					Image string `json:"image"`
				} `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	} `json:"output"`

	Usage struct {
		Width      int `json:"width"`
		Height     int `json:"height"`
		ImageCount int `json:"image_count"`
	} `json:"usage"`
}

type DashScopeErrorResponse struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	RequestID      string `json:"request_id,omitempty"`
	HTTPStatusCode int    `json:"-"`
}

func (e *DashScopeErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *DashScopeErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("ALI request error, status code: %d, business code: %s, request_id: %s, message: %s", e.HTTPStatusCode, e.Code, e.RequestID, e.Message))
	}
	return errors.New(fmt.Sprintf("ALI request error, message: %s", e.Message))
}

func (c *Client) CreateQwenImageGeneration(req *QwenImageGenerationRequest) (*QwenImageGenerationResponse, error) {
	base := strings.TrimRight(c.EndPoint, "/")
	url := base + "/api/v1/services/aigc/multimodal-generation/generation"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
	}

	resp, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return nil, err
	}
	if err := common.HttpCheckError(resp, &DashScopeErrorResponse{}); err != nil {
		return nil, err
	}

	var result QwenImageGenerationResponse
	if err := common.HttpDecodeResponse(resp, &result); err != nil {
		return nil, err
	}
	if result.Code != "" {
		return nil, errors.New(fmt.Sprintf("ALI request error, business code: %s, request_id: %s, message: %s", result.Code, result.RequestID, result.Message))
	}
	return &result, nil
}
