// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"errors"
	"io"

	"github.com/zhimaAi/llm_adaptor/api/openai"
	"github.com/zhimaAi/llm_adaptor/api/volcenginev3"
)

type ZhimaImageGenerationStream interface {
	Read() (ZhimaImageGenerationResp, error)
	Close() error
}

type ZhimaImageGenerationStreamRes struct {
	ZhimaImageGenerationStream
}

func (a *Adaptor) CreateImageGenerateStream(params *ZhimaImageGenerationReq) (*ZhimaImageGenerationStreamRes, error) {
	switch a.meta.Corp {
	case "302ai", "openrouter":
		apiUrl := "https://api.302ai.cn/302/images/generations"
		if a.meta.Corp == "openrouter" {
			apiUrl = "https://openrouter.ai/api/v1/chat/completions"
		}
		client := openai.NewClient(apiUrl, a.meta.APIKey, &openai.ErrorResponse{})
		req := map[string]any{
			`model`:  a.meta.Model,
			`prompt`: params.Prompt,
			`stream`: false,
		}
		formatOpenaiParams(params, req)
		stream, err := client.CreateImageGenerateStream(req)
		if err != nil {
			return &ZhimaImageGenerationStreamRes{}, err
		}
		return &ZhimaImageGenerationStreamRes{
			&OpenAIImageGenerationStreamResult{stream, *params.OutputFormat},
		}, nil
	case "doubao":
		client := volcenginev3.NewClient("https://ark.cn-beijing.volces.com/api/v3/images/generations", a.meta.Model, a.meta.APIKey, a.meta.SecretKey, a.meta.Region)
		req := map[string]any{
			`model`:  a.meta.Model,
			`prompt`: params.Prompt,
			`stream`: true,
		}
		formatDoubaoParams(params, req)
		stream, err := client.CreateImageGenerateStream(req)
		if err != nil {
			return &ZhimaImageGenerationStreamRes{}, err
		}
		return &ZhimaImageGenerationStreamRes{
			&OpenAIImageGenerationStreamResult{stream, `jpeg`},
		}, nil
	case "ali":
		resp, err := a.CreateImageGenerate(params)
		if err != nil {
			return &ZhimaImageGenerationStreamRes{}, err
		}
		return &ZhimaImageGenerationStreamRes{
			&SingleImageGenerationStream{resp: *resp},
		}, nil
	default:
		return &ZhimaImageGenerationStreamRes{}, errors.New("corp not support")
	}
}

type SingleImageGenerationStream struct {
	resp ZhimaImageGenerationResp
	done bool
}

func (s *SingleImageGenerationStream) Read() (ZhimaImageGenerationResp, error) {
	if s.done {
		return ZhimaImageGenerationResp{}, io.EOF
	}
	s.done = true
	return s.resp, nil
}

func (s *SingleImageGenerationStream) Close() error {
	s.done = true
	return nil
}
