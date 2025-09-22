// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package ali

import (
	"errors"

	"github.com/zhimaAi/llm_adaptor/api/openai"
	"github.com/zhimaAi/llm_adaptor/common"
)

type Client struct {
	EndPoint     string
	APIKey       string
	OpenAIClient *openai.Client // proxy openai
}

func NewClient(APIKey string) *Client {
	return &Client{
		EndPoint: "https://dashscope.aliyuncs.com",
		APIKey:   APIKey,
		OpenAIClient: &openai.Client{
			EndPoint: "https://dashscope.aliyuncs.com/compatible-mode/v1",
			APIKey:   APIKey,
			ErrResp:  &openai.ErrorResponse{},
		},
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (EmbeddingResponse, error) {

	url := c.EndPoint + "/api/v1/services/embeddings/text-embedding/text-embedding"
	headers := []common.Header{
		{Key: "Authorization", Value: c.APIKey},
	}
	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return EmbeddingResponse{}, err
	}

	var result EmbeddingResponse
	err = common.HttpDecodeResponse(responseRaw, &result)
	return result, err
}

type (
	CreateRerankReq struct {
		Model      string     `json:"model"`
		Input      Input      `json:"input"`
		Parameters Parameters `json:"parameters"`
	}
	Input struct {
		Query     string   `json:"query"`
		Documents []string `json:"documents"`
	}
	Parameters struct {
		ReturnDocuments bool `json:"return_documents"`
		TopN            int  `json:"top_n"`
	}
	CreateRerankRes struct {
		RequestId string `json:"request_id"`
		Code      string `json:"code"`
		Message   string `json:"message"`
		Output    struct {
			Results []Result `json:"results"`
		} `json:"output"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}
	Result struct {
		Index          int     `json:"index"`
		RelevanceScore float64 `json:"relevance_score"`
	}
)

func (r *CreateRerankReq) validate() error {
	return nil
}

func (c *Client) CreateRerank(req *CreateRerankReq) (*CreateRerankRes, error) {
	result := &CreateRerankRes{}
	if err := req.validate(); err != nil {
		return result, err
	}
	url := c.EndPoint + "/api/v1/services/rerank/text-rerank/text-rerank"
	headers := []common.Header{
		{Key: "Authorization", Value: "Bearer " + c.APIKey},
		{Key: "Context-Type", Value: "application/json"},
	}
	responseRaw, err := common.HttpPost(url, headers, nil, req)
	if err != nil {
		return result, err
	}
	err = common.HttpCheckError(responseRaw, &ErrorResponse{})
	if err != nil {
		return result, err
	}
	err = common.HttpDecodeResponse(responseRaw, &result)
	if err != nil {
		return result, err
	}
	if len(result.Code) > 0 && len(result.Output.Results) == 0 {
		return result, errors.New(result.Message)
	}
	return result, err
}
