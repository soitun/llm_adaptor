// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"errors"

	"github.com/zhimaAi/llm_adaptor/api/openai"
	"github.com/zhimaAi/llm_adaptor/api/volcenginev3"
)

type ZhimaImageGenerationStreamResult interface {
	Read() (ZhimaImageGenerationResp, error)
	Close() error
}

type ZhimaImageGenerationStreamResponse struct {
	ZhimaImageGenerationStreamResult
}

func (a *Adaptor) CreateImageGenerateStream(params *ZhimaImageGenerationReq) (ZhimaImageGenerationStreamResponse, error) {
	switch a.meta.Corp {
	case "doubao":
		client := volcenginev3.NewClient("https://ark.cn-beijing.volces.com/api/v3/images/generations", a.meta.Model, a.meta.APIKey, a.meta.SecretKey, a.meta.Region)
		req := openai.ImageGenerationRequest{
			Model:                            a.meta.Model,
			Prompt:                           params.Prompt,
			Image:                            params.Image,
			Size:                             params.Size,
			Seed:                             -1,
			SequentialImageGeneration:        params.SequentialImageGeneration,
			SequentialImageGenerationOptions: openai.SequentialImageGenerationOptions{MaxImages: params.MaxImages},
			Stream:                           true,
			ResponseFormat:                   params.ResponseFormat,
			Watermark:                        params.Watermark,
			OptimizePromptOptions:            openai.OptimizePromptOptions{Mode: params.OptimizePromptMode},
		}
		stream, err := client.CreateImageGenerateStream(req)
		if err != nil {
			return ZhimaImageGenerationStreamResponse{}, err
		}
		return ZhimaImageGenerationStreamResponse{
			&OpenAIImageGenerationStreamResult{stream},
		}, nil
	default:
		return ZhimaImageGenerationStreamResponse{}, errors.New("corp not support")
	}
}
