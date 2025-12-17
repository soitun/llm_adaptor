// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package jina

type ReRankRequest struct {
	Model           string   `json:"model"`
	Query           string   `json:"query"`
	Documents       []string `json:"documents"`
	TopN            int      `json:"top_n,omitempty"`
	ReturnDocuments bool     `json:"return_documents,omitempty"`
}

type ReRankResponse struct {
	Model   string   `json:"model"`
	Results []Result `json:"results"`
	Usage   Usage    `json:"usage"`
}

type Result struct {
	Index          int      `json:"index"`
	Document       Document `json:"document"`
	RelevanceScore float64  `json:"relevance_score"`
}

type Document struct {
	Text string `json:"text"`
}
