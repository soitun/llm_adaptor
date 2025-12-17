// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package voyage

type EmbeddingRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

type EmbeddingResponse struct {
	Data  []EmbeddingData `json:"data"`
	Usage Usage           `json:"usage"`
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type EmbeddingData struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Object    string    `json:"object"`
}
