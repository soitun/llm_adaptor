// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"errors"

	"github.com/zhimaAi/llm_adaptor/api/volcenginev3"
)

type SequentialImageGenerationOptions struct {
	MaxImages int `json:"max_images"`
}

type OptimizePromptOptions struct {
	Mode string `json:"mode"`
}

type ZhimaImageGenerationReq struct {
	Prompt                    string    `json:"prompt"`
	Image                     *[]string `json:"image"`
	Size                      *string   `json:"size"`
	SequentialImageGeneration *string   `json:"sequential_image_generation"`
	MaxImages                 int       `json:"max_images"`
	Stream                    bool      `json:"stream"`
	ResponseFormat            *string   `json:"response_format"`
	Watermark                 *bool     `json:"watermark"`
	OptimizePromptMode        *string   `json:"optimize_prompt_mode"`
	Seed                      *int      `json:"seed"`
	GuidanceScale             *float32  `json:"guidance_scale"`
}

type DataError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type ImageGenerationData struct {
	Url     string    `json:"url"`
	B64Json string    `json:"b64_json"`
	Size    string    `json:"size"`
	Error   DataError `json:"error"`
	Ext     string    `json:"ext"`
}
type ZhimaImageGenerationResp struct {
	Datas       []*ImageGenerationData `json:"data"`
	InputToken  int                    `json:"input_token"`
	OutputToken int                    `json:"output_token"`
}

func (a *Adaptor) CreateImageGenerate(params *ZhimaImageGenerationReq) (*ZhimaImageGenerationResp, error) {
	switch a.meta.Corp {
	case "doubao":
		client := volcenginev3.NewClient("https://ark.cn-beijing.volces.com/api/v3/images/generations", a.meta.Model, a.meta.APIKey, a.meta.SecretKey, a.meta.Region)
		req := map[string]any{
			`model`:  a.meta.Model,
			`prompt`: params.Prompt,
			`stream`: false,
		}
		formatDoubaoParams(params, req)
		res, err := client.CreateImageGenerate(req)
		if err != nil {
			return &ZhimaImageGenerationResp{}, err
		}
		datas := make([]*ImageGenerationData, len(res.Data))
		for i, item := range res.Data {
			datas[i] = &ImageGenerationData{
				Url:     item.Url,
				B64Json: item.B64Json,
				Size:    item.Size,
				Error:   DataError{Code: item.Error.Code, Message: item.Error.Message},
				Ext:     `jpeg`,
			}
		}
		return &ZhimaImageGenerationResp{
			InputToken:  res.Usage.TotalTokens - res.Usage.OutputTokens,
			OutputToken: res.Usage.OutputTokens,
			Datas:       datas,
		}, nil
	default:
		return &ZhimaImageGenerationResp{}, errors.New("corp not support")
	}
}

func formatDoubaoParams(params *ZhimaImageGenerationReq, req map[string]any) {
	if params.Seed != nil {
		req[`seed`] = *params.Seed
	}
	if params.Image != nil {
		req[`image`] = *params.Image
	}
	if params.Size != nil {
		req[`size`] = *params.Size
	}
	if params.SequentialImageGeneration != nil {
		req[`sequential_image_generation`] = *params.SequentialImageGeneration
		if *params.SequentialImageGeneration == `auto` {
			req[`sequential_image_generation_options `] = map[string]any{
				`max_images`: min(max(params.MaxImages, 1), 15),
			}
		}
	}
	if params.GuidanceScale != nil {
		req[`guidance_scale`] = min(max(1, *params.GuidanceScale), 10)
	}
	if params.ResponseFormat != nil {
		req[`response_format`] = *params.ResponseFormat
	}
	if params.Watermark != nil {
		req[`watermark`] = *params.Watermark
	}
	if params.OptimizePromptMode != nil {
		req[`optimize_prompt_options`] = map[string]any{
			`mode `: *params.OptimizePromptMode,
		}
	}
}
