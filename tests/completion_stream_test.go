// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package tests

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func testChatCompletionStream(Meta adaptor.Meta) {
	client := &adaptor.Adaptor{}
	client.Init(Meta)
	req := adaptor.ZhimaChatCompletionRequest{
		Messages:    []adaptor.ZhimaChatCompletionMessage{{Role: "user", Content: "你好"}},
		Temperature: 0.1,
		MaxToken:    10,
	}
	stream, err := client.CreateChatCompletionStream(req)
	if err != nil {
		panic(err.Error())
	}
	defer stream.Close()
	for {
		response, err := stream.Read()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}
		if err != nil {
			fmt.Printf("\nStream error: %v", err)
			return
		}
		fmt.Print(response.Result)
	}
}

func TestOpenAIChatCompletionStream(t *testing.T) {
	testChatCompletionStream(adaptor.Meta{
		Corp:   "openai",
		Model:  `gpt-3.5-turbo`,
		APIKey: os.Getenv(`OPENAI_KEY`),
	})
}

func TestMinimaxiChatCompletionStream(t *testing.T) {
	testChatCompletionStream(adaptor.Meta{
		Corp:   "minimax",
		Model:  `abab6.5s-chat`,
		APIKey: os.Getenv(`MINIMAX_KEY`),
	})
}

func TestSiliconFlowChatCompletionStream(t *testing.T) {
	testChatCompletionStream(adaptor.Meta{
		EndPoint:   `https://api.siliconflow.cn`,
		Corp:       "siliconflow",
		APIVersion: "v1",
		Model:      `Qwen/Qwen2.5-72B-Instruct`,
		APIKey:     os.Getenv(`SILICONFLOW_KEY`),
	})
}
