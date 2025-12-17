// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package baai

type Client struct {
	EndPoint string
	Model    string
	APIKey   string
}

func NewClient(endPoint, Model, APIKey string) *Client {
	return &Client{
		EndPoint: endPoint,
		Model:    Model,
		APIKey:   APIKey,
	}
}
