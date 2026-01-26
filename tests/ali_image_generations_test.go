// Copyright © 2016- 2026 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package tests

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func TestAliQwenImageGenerations(t *testing.T) {
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		t.Skip("skip: DASHSCOPE_API_KEY not set")
	}

	endpoint := os.Getenv("DASHSCOPE_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://dashscope.aliyuncs.com"
	}
	model := os.Getenv("DASHSCOPE_IMAGE_MODEL")
	if model == "" {
		model = "qwen-image-max"
	}

	client := &adaptor.Adaptor{}
	client.Init(adaptor.Meta{
		Corp:     "ali",
		EndPoint: endpoint,
		Model:    model,
		APIKey:   apiKey,
	})

	resp, err := client.CreateImageGenerate(&adaptor.ZhimaImageGenerationReq{
		Prompt:             "生成一张海边日落的图片，写实风格",
		Size:               getString("2K"),
		Watermark:          getBool(false),
		OptimizePromptMode: getString("standard"),
	})
	if err != nil {
		t.Fatalf("CreateImageGenerate error: %v", err)
	}
	if resp == nil || len(resp.Datas) != 1 {
		t.Fatalf("unexpected resp: %+v", resp)
	}
	if resp.Datas[0].Url == "" {
		t.Fatalf("missing image url: %+v", resp.Datas[0])
	}
	if resp.Datas[0].Ext != "png" {
		t.Fatalf("unexpected ext: %s", resp.Datas[0].Ext)
	}
	if resp.Datas[0].Size != "1664*928" {
		t.Fatalf("unexpected mapped size for 2K: %s", resp.Datas[0].Size)
	}

	resp2, err := client.CreateImageGenerate(&adaptor.ZhimaImageGenerationReq{
		Prompt:             "生成一张海边日落的图片，写实风格",
		Size:               getString("2048x2048"),
		Watermark:          getBool(false),
		OptimizePromptMode: getString("standard"),
	})
	if err != nil {
		t.Fatalf("CreateImageGenerate error: %v", err)
	}
	if resp2 == nil || len(resp2.Datas) != 1 {
		t.Fatalf("unexpected resp: %+v", resp2)
	}
	if resp2.Datas[0].Size != "1328*1328" {
		t.Fatalf("unexpected mapped size for 2048x2048: %s", resp2.Datas[0].Size)
	}
}

func TestAliQwenImageGenerationsStream(t *testing.T) {
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		t.Skip("skip: DASHSCOPE_API_KEY not set")
	}

	endpoint := os.Getenv("DASHSCOPE_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://dashscope.aliyuncs.com"
	}
	model := os.Getenv("DASHSCOPE_IMAGE_MODEL")
	if model == "" {
		model = "qwen-image-max"
	}

	client := &adaptor.Adaptor{}
	client.Init(adaptor.Meta{
		Corp:     "ali",
		EndPoint: endpoint,
		Model:    model,
		APIKey:   apiKey,
	})

	stream, err := client.CreateImageGenerateStream(&adaptor.ZhimaImageGenerationReq{
		Prompt:             "生成一张海边日落的图片，写实风格",
		Size:               getString("1664*928"),
		Watermark:          getBool(false),
		OptimizePromptMode: getString("standard"),
		Stream:             true,
	})
	if err != nil {
		t.Fatalf("CreateImageGenerateStream error: %v", err)
	}
	defer stream.Close()

	_, err = stream.Read()
	if err != nil {
		t.Fatalf("stream first read error: %v", err)
	}
	_, err = stream.Read()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got: %v", err)
	}
}
