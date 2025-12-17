// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package xinference

import "time"

type Duration struct {
	time.Duration
}

type EmbeddingRequest struct {
	// Model is the model name.
	Model string `json:"model"`

	// Prompt is the textual prompt to embed.
	Input []string `json:"input"`
}

type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
type EmbeddingResponse struct {
	Data []Embedding `json:"data"`
}
type Embedding struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}
