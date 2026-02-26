// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/zhimaAi/llm_adaptor/api/ali"
	"github.com/zhimaAi/llm_adaptor/api/openai"
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
	OutputFormat              *string   `json:"output_format"`
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
	case "302ai":
		apiUrl := "https://api.302ai.cn/302/images/generations"
		client := openai.NewClient(apiUrl, a.meta.APIKey, &openai.ErrorResponse{})
		req := map[string]any{
			`model`:  a.meta.Model,
			`prompt`: params.Prompt,
			`stream`: false,
		}
		formatOpenaiParams(params, req)
		res, err := client.CreateImageGenerate(req)
		if err != nil {
			return &ZhimaImageGenerationResp{}, err
		}
		datas := make([]*ImageGenerationData, len(res.Data))
		for i, item := range res.Data {
			datas[i] = &ImageGenerationData{
				Url:     item.Url,
				B64Json: item.B64Json,
				Ext:     *params.OutputFormat,
			}
		}
		return &ZhimaImageGenerationResp{
			InputToken:  res.Usage.TotalTokens - res.Usage.OutputTokens,
			OutputToken: res.Usage.OutputTokens,
			Datas:       datas,
		}, nil
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
	case "ali":
		client := ali.NewClient(a.meta.APIKey)
		if strings.TrimSpace(a.meta.EndPoint) != "" {
			client.EndPoint = a.meta.EndPoint
		}

		req := &ali.QwenImageGenerationRequest{
			Model: a.meta.Model,
			Input: ali.QwenImageInput{
				Messages: []ali.QwenImageMessage{
					{
						Role: "user",
						Content: []ali.QwenImageContent{
							{Text: params.Prompt},
						},
					},
				},
			},
		}
		mappedSize := formatAliQwenImageParams(params, req)

		res, err := client.CreateQwenImageGeneration(req)
		if err != nil {
			return &ZhimaImageGenerationResp{}, err
		}

		var imageURL string
		if len(res.Output.Choices) > 0 && len(res.Output.Choices[0].Message.Content) > 0 {
			imageURL = res.Output.Choices[0].Message.Content[0].Image
		}
		if imageURL == "" {
			return &ZhimaImageGenerationResp{}, errors.New("ali qwen-image response missing image url")
		}

		b64 := ""
		if params.ResponseFormat != nil && *params.ResponseFormat == "b64_json" {
			b64, err = downloadImageAsBase64(imageURL)
			if err != nil {
				return &ZhimaImageGenerationResp{}, err
			}
			imageURL = ""
		}
		return &ZhimaImageGenerationResp{
			InputToken:  0,
			OutputToken: 0,
			Datas: []*ImageGenerationData{
				{
					Url:     imageURL,
					B64Json: b64,
					Size:    mappedSize,
					Error:   DataError{},
					Ext:     "png",
				},
			},
		}, nil
	default:
		return &ZhimaImageGenerationResp{}, errors.New("corp not support")
	}
}

func formatAliQwenImageParams(params *ZhimaImageGenerationReq, req *ali.QwenImageGenerationRequest) string {
	mappedSize := ""
	var p ali.QwenImageGenParameter
	var hasParam bool
	if params.Size != nil && strings.TrimSpace(*params.Size) != "" {
		mappedSize = mapAliImageSize(*params.Size)
		size := mappedSize
		p.Size = &size
		hasParam = true
	}
	if params.Watermark != nil {
		p.Watermark = params.Watermark
		hasParam = true
	}
	if params.Seed != nil {
		p.Seed = params.Seed
		hasParam = true
	}
	if params.OptimizePromptMode != nil {
		extend := *params.OptimizePromptMode != "off"
		p.PromptExtend = &extend
		hasParam = true
	}
	if hasParam {
		req.Parameters = &p
	}
	return mappedSize
}

func mapAliImageSize(size string) string {
	s := strings.TrimSpace(size)
	if s == "" {
		return "1664*928"
	}
	switch strings.ToUpper(s) {
	case "2K":
		return "1664*928"
	case "4K":
		return "1328*1328"
	}
	s = strings.NewReplacer("x", "*", "X", "*", "×", "*").Replace(s)
	parts := strings.Split(s, "*")
	if len(parts) != 2 {
		return "1664*928"
	}
	w, okW := atoi(parts[0])
	h, okH := atoi(parts[1])
	if !okW || !okH || w <= 0 || h <= 0 {
		return "1664*928"
	}
	r := float64(w) / float64(h)

	type cand struct {
		size  string
		ratio float64
		pix   int
	}
	cands := []cand{
		{size: "1664*928", ratio: 1664.0 / 928.0, pix: 1664 * 928},
		{size: "1472*1104", ratio: 1472.0 / 1104.0, pix: 1472 * 1104},
		{size: "1328*1328", ratio: 1, pix: 1328 * 1328},
		{size: "1104*1472", ratio: 1104.0 / 1472.0, pix: 1104 * 1472},
		{size: "928*1664", ratio: 928.0 / 1664.0, pix: 928 * 1664},
	}

	best := cands[0]
	bestDiff := absFloat(r - best.ratio)
	for i := 1; i < len(cands); i++ {
		d := absFloat(r - cands[i].ratio)
		if d < bestDiff || (d == bestDiff && cands[i].pix > best.pix) {
			best = cands[i]
			bestDiff = d
		}
	}
	return best.size
}

func downloadImageAsBase64(url string) (string, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("download image failed")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(body), nil
}

func atoi(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	n := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func formatDoubaoParams(params *ZhimaImageGenerationReq, req map[string]any) {
	if params.Seed != nil {
		req[`seed`] = *params.Seed
	}
	if params.Image != nil {
		req[`image`] = *params.Image
	}
	if params.Size != nil && *params.Size != "auto" {
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

func formatOpenaiParams(params *ZhimaImageGenerationReq, req map[string]any) {
	if params.Size != nil && *params.Size != "auto" {
		req[`size`] = *params.Size
	}
	if params.SequentialImageGeneration != nil {
		if *params.SequentialImageGeneration == `auto` {
			req[`n`] = params.MaxImages
		}
	}
	if params.ResponseFormat != nil {
		req[`response_format`] = *params.ResponseFormat
	}
	if params.OutputFormat != nil {
		req[`output_format`] = params.OutputFormat
	} else {
		params.OutputFormat = tea.String(`jpeg`)
		req[`output_format`] = params.OutputFormat
	}
}
