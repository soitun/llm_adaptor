package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func testRerank(Meta adaptor.Meta) {
	client := &adaptor.Adaptor{}
	client.Init(Meta)
	req := &adaptor.ZhimaRerankReq{
		Query:    "数学",
		Passages: []string{`语文`, `英语`, `应用数学`, `高数`, `汉言文学`, `小学算术`},
		TopK:     10,
	}
	res, err := client.CreateRerank(req)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
}
func TestSiliconFlowRerank(t *testing.T) {
	testRerank(adaptor.Meta{
		EndPoint:   `https://api.siliconflow.cn`,
		Corp:       "siliconflow",
		APIVersion: "v1",
		Model:      `BAAI/bge-reranker-v2-m3`,
		APIKey:     os.Getenv(`SILICONFLOW_KEY`),
	})
}
