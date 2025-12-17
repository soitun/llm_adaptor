package openai

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/zhimaAi/llm_adaptor/basics"
	"github.com/zhimaAi/llm_adaptor/common"
)

type SequentialImageGenerationOptions struct {
	MaxImages int `json:"max_images"`
}

type OptimizePromptOptions struct {
	Mode string `json:"mode"`
}
type ImageGenerationRequest struct {
	Model                            string                           `json:"model"`
	Prompt                           string                           `json:"prompt"`
	Image                            []string                         `json:"image"`
	Size                             string                           `json:"size"`
	Seed                             int                              `json:"seed"`
	SequentialImageGeneration        string                           `json:"sequential_image_generation"`
	SequentialImageGenerationOptions SequentialImageGenerationOptions `json:"sequential_image_generation_options"`
	Stream                           bool                             `json:"stream"`
	ResponseFormat                   string                           `json:"response_format"`
	Watermark                        bool                             `json:"watermark"`
	OptimizePromptOptions            OptimizePromptOptions            `json:"optimize_prompt_options"`
}
type DataItem struct {
	Url     string    `json:"url"`
	B64Json string    `json:"b64_json"`
	Size    string    `json:"size"`
	Error   DataError `json:"error"`
}
type DataError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Usage struct {
	GeneratedImages int `json:"generated_images"`
	OutputTokens    int `json:"output_tokens"`
	TotalTokens     int `json:"total_tokens"`
}

type ImageGenerationResponse struct {
	Model   string     `json:"model"`
	Created int        `json:"created"`
	Data    []DataItem `json:"data"`
	Usage   Usage      `json:"usage"`
	Error   Error      `json:"error"`
}

type ImageGenerationStreamResponse struct {
	Type       string    `json:"type"`
	Model      string    `json:"model"`
	Created    int       `json:"created"`
	ImageIndex int       `json:"image_index"`
	Url        string    `json:"url"`
	B64Json    string    `json:"b64_json"`
	Size       string    `json:"size"`
	Error      DataError `json:"error"`
	Usage      Usage     `json:"usage"`
}

type ImageGenerationStream struct {
	*common.StreamReader[ImageGenerationStreamResponse]
}

func (c *ImageGenerationStream) Recv() (ImageGenerationStreamResponse, error) {
	if c.StreamReader.IsFinished {
		return ImageGenerationStreamResponse{}, io.EOF
	}

	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)

	for {
		rawLine, readErr := c.StreamReader.Reader.ReadBytes('\n')
		if readErr != nil || hasErrorPrefix {
			if readErr != io.EOF {
				c.StreamReader.UnmarshalError()
				if c.StreamReader.ErrorResponse != nil {
					return *new(ImageGenerationStreamResponse), fmt.Errorf("unmarshal error, %w", c.StreamReader.ErrorResponse.Error())
				}
				return *new(ImageGenerationStreamResponse), readErr
			} else {
				c.StreamReader.IsFinished = true
				return *new(ImageGenerationStreamResponse), io.EOF
			}
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			writeErr := c.StreamReader.ErrAccumulator.Write(noSpaceLine)
			if writeErr != nil {
				return *new(ImageGenerationStreamResponse), writeErr
			}
			emptyMessagesCount++
			if emptyMessagesCount > c.StreamReader.EmptyMessagesLimit {
				return *new(ImageGenerationStreamResponse), errors.New("stream has sent too many empty messages")
			}

			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if strings.TrimSpace(string(noPrefixLine)) == "[DONE]" {
			c.StreamReader.IsFinished = true
			return *new(ImageGenerationStreamResponse), io.EOF
		}

		var response ImageGenerationStreamResponse
		unmarshalErr := basics.JsonDecode(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(ImageGenerationStreamResponse), unmarshalErr
		}

		return response, nil
	}
}
