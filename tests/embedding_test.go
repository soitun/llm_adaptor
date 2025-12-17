// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func testEmbedding(Meta adaptor.Meta) {
	client := &adaptor.Adaptor{}
	client.Init(Meta)
	req := adaptor.ZhimaEmbeddingRequest{
		Input: "数学",
	}
	res, err := client.CreateEmbeddings(req)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
}
func TestSiliconFlowEmbedding(t *testing.T) {
	testEmbedding(adaptor.Meta{
		EndPoint:   `https://api.siliconflow.cn`,
		Corp:       "siliconflow",
		APIVersion: "v1",
		Model:      `BAAI/bge-m3`,
		APIKey:     os.Getenv(`SILICONFLOW_KEY`),
	})
}
