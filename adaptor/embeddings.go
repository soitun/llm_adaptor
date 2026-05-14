// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"context"
	"errors"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentHunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/zhimaAi/llm_adaptor/api/ali"
	"github.com/zhimaAi/llm_adaptor/api/azure"
	"github.com/zhimaAi/llm_adaptor/api/baai"
	"github.com/zhimaAi/llm_adaptor/api/baichuan"
	"github.com/zhimaAi/llm_adaptor/api/baidu"
	"github.com/zhimaAi/llm_adaptor/api/cohere"
	"github.com/zhimaAi/llm_adaptor/api/gemini"
	"github.com/zhimaAi/llm_adaptor/api/hunyuan"
	"github.com/zhimaAi/llm_adaptor/api/jina"
	"github.com/zhimaAi/llm_adaptor/api/ollama"
	"github.com/zhimaAi/llm_adaptor/api/openai"
	openaiagent "github.com/zhimaAi/llm_adaptor/api/openaiAgent"
	"github.com/zhimaAi/llm_adaptor/api/siliconflow"
	"github.com/zhimaAi/llm_adaptor/api/voyage"
	"github.com/zhimaAi/llm_adaptor/api/xinference"
	"github.com/zhimaAi/llm_adaptor/api/zhipu"
)

type ZhimaEmbeddingRequest struct {
	Input string `json:"input"`
}

type ZhimaEmbeddingResponse struct {
	Result          []float64 `json:"result"`
	PromptToken     int       `json:"prompt_token"`
	CompletionToken int       `json:"completion_token"`
}

func (a *Adaptor) CreateEmbeddings(req ZhimaEmbeddingRequest) (ZhimaEmbeddingResponse, error) {
	if req.Input == "" {
		return ZhimaEmbeddingResponse{}, errors.New("input empty")
	}
	a.meta.EndPoint = strings.TrimRight(strings.TrimSpace(a.meta.EndPoint), `/`)
	switch a.meta.Corp {
	case "openai":
		client := openai.NewClient(GenerateOpenAiApiUrl(a), a.meta.APIKey, &openai.ErrorResponse{})
		r := openai.EmbeddingRequest{
			Model: a.meta.Model,
			Input: []string{req.Input},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.TotalTokens - res.Usage.PromptTokens,
		}, nil
	case "baichuan", "zhipu", "openaiAgent":
		var client *openai.Client
		if a.meta.Corp == "baichuan" {
			c := baichuan.NewClient(a.meta.APIKey)
			if len(a.meta.EndPoint) > 0 {
				c.EndPoint, c.OpenAIClient.EndPoint = GenerateClientEndPoint(a)
			}
			client = c.OpenAIClient
		} else if a.meta.Corp == "zhipu" {
			c := zhipu.NewClient(a.meta.APIKey)
			if len(a.meta.EndPoint) > 0 {
				c.EndPoint, c.OpenAIClient.EndPoint = GenerateClientEndPoint(a)
			}
			client = c.OpenAIClient
		} else if a.meta.Corp == "openaiAgent" {
			client = openaiagent.NewClient(a.meta.EndPoint, a.meta.APIKey, a.meta.APIVersion).OpenAIClient
		}
		r := openai.EmbeddingRequest{
			Model: a.meta.Model,
			Input: []string{req.Input},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.TotalTokens - res.Usage.PromptTokens,
		}, nil
	case "azure":
		client := azure.NewClient(
			a.meta.EndPoint,
			a.meta.APIVersion,
			a.meta.APIKey,
			a.meta.Model,
		)
		r := azure.EmbeddingRequest{
			Input: []string{req.Input},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.TotalTokens - res.Usage.PromptTokens,
		}, nil
	case "baidu":
		client := baidu.NewClient(
			a.meta.APIKey,
			a.meta.SecretKey,
			a.meta.Model,
		)
		r := baidu.EmbeddingRequest{
			Input: []string{req.Input},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.CompletionTokens,
		}, nil
	case "ali":
		client := ali.NewClient(a.meta.APIKey)
		if len(a.meta.EndPoint) > 0 {
			client.EndPoint, _ = GenerateClientEndPoint(a)
		}
		dimension := a.meta.Dimension
		if dimension == 0 {
			if a.meta.Model == `text-embedding-v3` {
				dimension = 1024 // should be in [64, 128, 256, 512, 768, 1024]
			} else {
				dimension = 1536
			}
		}
		textType := a.meta.TextType
		if textType == "" {
			textType = "document"
		}
		r := ali.EmbeddingRequest{
			Input:      ali.Texts{Texts: []string{req.Input}},
			Model:      a.meta.Model,
			Parameters: ali.QwenParameters{TextType: textType, Dimension: dimension},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Output.Embeddings[0].Embedding,
			PromptToken:     res.Usage.TotalTokens,
			CompletionToken: 0,
		}, nil
	case "voyage":
		client := voyage.NewClient(
			a.meta.APIKey,
		)
		r := voyage.EmbeddingRequest{
			Input: []string{req.Input},
			Model: a.meta.Model,
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.TotalTokens - res.Usage.PromptTokens,
		}, nil
	case "gemini":
		client := gemini.NewClient(
			a.meta.APIKey,
			a.meta.Model,
		)
		if len(a.meta.EndPoint) > 0 {
			client.EndPoint, _ = GenerateClientEndPoint(a)
		}
		r := gemini.EmbeddingRequest{
			Content: gemini.Content{Parts: []gemini.Part{{Text: req.Input}}},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result: res.Embedding.Values,
		}, nil

	case "baai":
		client := baai.NewClient(a.meta.EndPoint, a.meta.Model, a.meta.APIKey)
		r := baai.EmbeddingRequest{
			Model: a.meta.Model,
			Input: []string{req.Input},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.TotalTokens - res.Usage.PromptTokens,
		}, nil
	case "doubao":
		baseUrl := `https://ark.cn-beijing.volces.com/api/v3`
		if len(a.meta.EndPoint) > 0 {
			baseUrl, _ = GenerateClientEndPoint(a)
		}
		var client *arkruntime.Client
		if len(a.meta.SecretKey) == 0 {
			client = arkruntime.NewClientWithApiKey(a.meta.APIKey,
				arkruntime.WithBaseUrl(baseUrl),
				arkruntime.WithRegion(a.meta.Region))
		} else {
			client = arkruntime.NewClientWithAkSk(a.meta.APIKey, a.meta.SecretKey,
				arkruntime.WithBaseUrl(baseUrl),
				arkruntime.WithRegion(a.meta.Region))
		}
		res, err := client.CreateEmbeddings(context.Background(),
			model.EmbeddingRequestStrings{
				Input: []string{req.Input},
				Model: a.meta.Model,
			})
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		result := make([]float64, len(res.Data[0].Embedding))
		for idx := range res.Data[0].Embedding {
			result[idx] = float64(res.Data[0].Embedding[idx])
		}
		return ZhimaEmbeddingResponse{
			Result:          result,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.CompletionTokens,
		}, nil
	case "cohere":
		client := cohere.NewClient(a.meta.APIKey)
		if len(a.meta.EndPoint) > 0 {
			client.EndPoint, _ = GenerateClientEndPoint(a)
		}
		r := cohere.EmbeddingRequest{
			Texts:     []string{req.Input},
			Model:     a.meta.Model,
			InputType: "classification",
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Embeddings[0],
			PromptToken:     res.Meta.Tokens.InputTokens,
			CompletionToken: res.Meta.Tokens.OutputTokens,
		}, nil
	case "hunyuan":
		client := hunyuan.NewClient(a.meta.APIKey, a.meta.SecretKey, a.meta.Region)
		r := tencentHunyuan.NewGetEmbeddingRequest()
		r.Input = common.StringPtr(req.Input)
		res, err := client.CreateEmbeddings(*r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		var result []float64
		for _, v := range res.Data[0].Embedding {
			result = append(result, *v)
		}
		return ZhimaEmbeddingResponse{
			Result:          result,
			PromptToken:     int(*res.Usage.PromptTokens),
			CompletionToken: int(*res.Usage.TotalTokens) - int(*res.Usage.PromptTokens),
		}, nil
	case "jina":
		client := jina.NewClient(a.meta.APIKey)
		if len(a.meta.EndPoint) > 0 {
			client.EndPoint = a.meta.EndPoint
		}
		r := jina.EmbeddingRequest{
			Input:        []string{req.Input},
			Model:        a.meta.Model,
			EncodingType: "float",
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.TotalTokens - res.Usage.PromptTokens,
		}, nil
	case "ollama":
		client := ollama.NewClient(a.meta.EndPoint, a.meta.Model)
		r := ollama.EmbeddingRequest{
			Prompt: req.Input,
			Model:  a.meta.Model,
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result: res.Embedding,
		}, nil
	case "xinference":
		client := xinference.NewClient(a.meta.EndPoint, a.meta.APIVersion, a.meta.Model)
		r := xinference.EmbeddingRequest{
			Input: []string{req.Input},
			Model: a.meta.Model,
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result: res.Data[0].Embedding,
		}, nil
	case "siliconflow":
		client := siliconflow.NewClient(a.meta.EndPoint, a.meta.APIKey, a.meta.APIVersion)
		r := siliconflow.EmbeddingRequest{
			Model: a.meta.Model,
			Input: []string{req.Input},
		}
		res, err := client.CreateEmbeddings(r)
		if err != nil {
			return ZhimaEmbeddingResponse{}, err
		}
		return ZhimaEmbeddingResponse{
			Result:          res.Data[0].Embedding,
			PromptToken:     res.Usage.PromptTokens,
			CompletionToken: res.Usage.TotalTokens - res.Usage.PromptTokens,
		}, nil
	}
	return ZhimaEmbeddingResponse{}, nil
}
