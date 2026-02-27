// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"errors"
	"io"
	"strings"

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
		baseUrl := "https://ark.cn-beijing.volces.com/api/v3"
		if strings.TrimSpace(a.meta.EndPoint) != "" {
			baseUrl = strings.TrimSpace(a.meta.EndPoint)
		}
		client := volcenginev3.NewClient(baseUrl+"/images/generations", a.meta.Model, a.meta.APIKey, a.meta.SecretKey, a.meta.Region)
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
