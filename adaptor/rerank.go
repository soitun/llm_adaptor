// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"sort"
	"strings"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/llm_adaptor/api/ali"
	"github.com/zhimaAi/llm_adaptor/api/baai"
	"github.com/zhimaAi/llm_adaptor/api/cohere"
	"github.com/zhimaAi/llm_adaptor/api/jina"
	"github.com/zhimaAi/llm_adaptor/api/siliconflow"
	"github.com/zhimaAi/llm_adaptor/api/xinference"
)

type ZhimaRerankReq struct {
	Enable   bool
	Query    string        `json:"query" toml:"query"`
	Passages []string      `json:"passages" toml:"passages"`
	Data     []msql.Params `json:"data"`
	TopK     int           `json:"top_k" toml:"top_k"`
}
type RerankData struct {
	Index          int     `json:"index"`
	Text           string  `json:"text"`
	RelevanceScore float64 `json:"relevance_score"`
}
type ZhimaRerankResp struct {
	Data        []*RerankData `json:"data"`
	InputToken  int           `json:"input_token"`
	OutputToken int           `json:"output_token"`
}

func (a *Adaptor) CreateRerank(params *ZhimaRerankReq) (ZhimaRerankResp, error) {
	a.meta.EndPoint = strings.TrimRight(strings.TrimSpace(a.meta.EndPoint), `/`)
	logs.Debug(`CreateRerank endpoint %s`, a.meta.EndPoint)
	zhimaRes := ZhimaRerankResp{}
	switch a.meta.Corp {
	case "baai":
		client := baai.NewClient(a.meta.EndPoint, a.meta.Model, a.meta.APIKey)
		req := &baai.CreateRerankReq{
			Model:    a.meta.Model,
			Query:    params.Query,
			Passages: params.Passages,
			TopK:     params.TopK,
		}
		res, err := client.CreateRerank(req)
		if err != nil || len(res.Results) <= 0 {
			return ZhimaRerankResp{}, err
		}
		for _, item := range res.Results {
			zhimaRes.Data = append(zhimaRes.Data, &RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "cohere":
		client := cohere.NewClient(a.meta.APIKey)
		if len(a.meta.EndPoint) > 0 {
			client.EndPoint, _ = GenerateClientEndPoint(a)
		}
		req := cohere.ReRankRequest{
			Model:     a.meta.Model,
			Query:     params.Query,
			Documents: params.Passages,
			TopN:      params.TopK,
		}
		res, err := client.ReRank(req)
		if err != nil || len(res.Results) <= 0 {
			return ZhimaRerankResp{}, err
		}
		for _, item := range res.Results {
			zhimaRes.Data = append(zhimaRes.Data, &RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "jina":
		client := jina.NewClient(a.meta.APIKey)
		if len(a.meta.EndPoint) > 0 {
			client.EndPoint, _ = GenerateClientEndPoint(a)
		}
		req := jina.ReRankRequest{
			Model:     a.meta.Model,
			Query:     params.Query,
			Documents: params.Passages,
			TopN:      params.TopK,
		}
		res, err := client.ReRank(req)
		if err != nil {
			return ZhimaRerankResp{}, err
		}
		for _, item := range res.Results {
			zhimaRes.Data = append(zhimaRes.Data, &RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "xinference":
		client := xinference.NewClient(a.meta.EndPoint, a.meta.APIVersion, a.meta.Model)
		req := &xinference.CreateRerankReq{
			Model:     a.meta.Model,
			Query:     params.Query,
			Documents: params.Passages,
			TopN:      params.TopK,
		}
		res, err := client.CreateRerank(req)
		if err != nil {
			return ZhimaRerankResp{}, err
		}
		for _, item := range res.Results {
			zhimaRes.Data = append(zhimaRes.Data, &RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "siliconflow":
		client := siliconflow.NewClient(a.meta.EndPoint, a.meta.APIKey, a.meta.Model)
		req := &siliconflow.CreateRerankReq{
			Model:     a.meta.Model,
			Query:     params.Query,
			Documents: params.Passages,
			TopK:      params.TopK,
		}
		res, err := client.CreateRerank(req)
		if err != nil || len(res.Results) <= 0 {
			return ZhimaRerankResp{}, err
		}
		zhimaRes.InputToken = res.Meta.Tokens.InputTokens
		zhimaRes.OutputToken = res.Meta.Tokens.OutputTokens
		for _, item := range res.Results {
			zhimaRes.Data = append(zhimaRes.Data, &RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	case "ali":
		client := ali.NewClient(a.meta.APIKey)
		if len(a.meta.EndPoint) > 0 {
			client.EndPoint, client.OpenAIClient.EndPoint = GenerateClientEndPoint(a)
		}
		req := &ali.CreateRerankReq{
			Model: a.meta.Model,
			Input: ali.Input{
				Query:     params.Query,
				Documents: params.Passages,
			},
			Parameters: ali.Parameters{
				ReturnDocuments: false,
				TopN:            params.TopK,
			},
		}
		res, err := client.CreateRerank(req)
		if err != nil {
			return ZhimaRerankResp{}, err
		}
		zhimaRes.InputToken = res.Usage.TotalTokens / 2 //没办法,输入输出个半
		zhimaRes.OutputToken = res.Usage.TotalTokens - zhimaRes.InputToken
		for _, item := range res.Output.Results {
			zhimaRes.Data = append(zhimaRes.Data, &RerankData{
				Index:          item.Index,
				RelevanceScore: item.RelevanceScore,
			})
		}
	}
	rerankData(params, zhimaRes.Data)
	return zhimaRes, nil
}

func rerankData(req *ZhimaRerankReq, rerankData []*RerankData) []*RerankData {
	newData := make([]*RerankData, 0)
	if len(rerankData) <= 0 {
		return newData
	}
	sort.Slice(rerankData, func(i, j int) bool {
		return rerankData[i].RelevanceScore > rerankData[j].RelevanceScore
	})
	for key, item := range rerankData {
		// topN filter
		if req.TopK > 0 && len(newData) >= req.TopK {
			continue
		}
		if req.Passages[item.Index] != "" {
			rerankData[key].Text = req.Passages[item.Index]
		}
	}
	return newData
}
