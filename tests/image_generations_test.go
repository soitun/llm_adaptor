// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package tests

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func TestDoubaoImageGenerations(t *testing.T) {
	imageGenerations(adaptor.Meta{
		Corp:   "doubao",
		Model:  `doubao-seedream-4-5-251128`,
		APIKey: os.Getenv(`doubao_apikey`),
	})
}

func imageGenerations(meta adaptor.Meta) {
	client := &adaptor.Adaptor{}
	client.Init(meta)
	ret, err := client.CreateImageGenerate(&adaptor.ZhimaImageGenerationReq{
		Prompt:                    "帮我生成一个大海的图片",
		Image:                     nil,
		Size:                      getString("2048x2048"),
		SequentialImageGeneration: getString("disabled"),
		MaxImages:                 1,
		Stream:                    false,
		ResponseFormat:            getString("b64_json"),
		Watermark:                 getBool(false),
		OptimizePromptMode:        getString("standard"),
	})
	if err != nil {
		fmt.Println(fmt.Sprintf(`生成失败：%s`, err.Error()))
	} else {
		for _, val := range ret.Datas {
			fmt.Println(fmt.Sprintf(`图片生成成功：%s`, val.B64Json))
		}
	}
}

func TestDoubaoImageGenerationsStream(t *testing.T) {
	imageGenerationsStream(adaptor.Meta{
		Corp:   "doubao",
		Model:  `doubao-seedream-4-5-251128`,
		APIKey: os.Getenv(`doubao_apikey`),
	})
}

func imageGenerationsStream(Meta adaptor.Meta) {
	client := &adaptor.Adaptor{}
	client.Init(Meta)
	req := &adaptor.ZhimaImageGenerationReq{
		Prompt:                    "帮我生成一个大海的图片，我要2张",
		Image:                     nil,
		Size:                      getString("2048x2048"),
		SequentialImageGeneration: getString("auto"),
		MaxImages:                 2,
		Stream:                    true,
		ResponseFormat:            getString("url"),
		Watermark:                 getBool(false),
		OptimizePromptMode:        getString("standard"),
	}
	stream, err := client.CreateImageGenerateStream(req)

	if err != nil {
		fmt.Println(fmt.Sprintf(`run error: %s`, err.Error()))
		return
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
		fmt.Println(fmt.Sprintf(`stream result %s`, tool.JsonEncodeNoError(response)))
	}
}

func getString(msg string) *string {
	return &msg
}

func getBool(msg bool) *bool {
	return &msg
}
