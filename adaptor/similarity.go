// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/llm_adaptor/api/baai"
)

type ZhimaSimilarityRequest struct {
	Model string      `json:"model"`
	Query []float64   `json:"query"`
	Input [][]float64 `json:"input"`
}

type ZhimaSimilarityResponse struct {
	Result []float64 `json:"result"`
}

func (a *Adaptor) CreateSimilarity(req ZhimaSimilarityRequest) (ZhimaSimilarityResponse, error) {
	logs.Debug(`CreateSimilarity endpoint %s`, a.meta.EndPoint)
	switch a.meta.Corp {
	case "baai":
		client := baai.NewClient(a.meta.EndPoint, a.meta.Model, a.meta.APIKey)
		req := baai.SimilarityRequest{
			Model: req.Model,
			Query: req.Query,
			Input: req.Input,
		}
		res, err := client.ComputeSimilarity(req)
		if err != nil {
			return ZhimaSimilarityResponse{}, err
		}
		return ZhimaSimilarityResponse{
			Result: res.Data,
		}, nil
	}
	return ZhimaSimilarityResponse{}, nil
}
