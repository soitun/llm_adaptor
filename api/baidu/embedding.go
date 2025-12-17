// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package baidu

type EmbeddingRequest struct {
	Input []string `json:"input"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type EmbeddingData struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Object    string    `json:"object"`
}

type EmbeddingResponse struct {
	Data  []EmbeddingData `json:"data"`
	Usage Usage           `json:"usage"`
}
