// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"errors"

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
	default:
		return &ZhimaImageGenerationStreamRes{}, errors.New("corp not support")
	}
}
