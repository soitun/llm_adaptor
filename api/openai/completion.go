// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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

type ChatCompletionRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatCompletionResponseImageUrl struct {
	Url string `json:"url"`
}
type ChatCompletionResponseImage struct {
	Type     string                         `json:"type"`
	ImageUrl ChatCompletionResponseImageUrl `json:"image_url"`
}
type ChatCompletionResponseMessage struct {
	Role             string                        `json:"role"`
	Content          string                        `json:"content"`
	ReasoningContent string                        `json:"reasoning_content"`
	ToolCalls        basics.ToolCalls              `json:"tool_calls"`
	Images           []ChatCompletionResponseImage `json:"images"`
}

type ThinkingType string

const (
	ThinkingTypeEnabled  ThinkingType = "enabled"
	ThinkingTypeDisabled ThinkingType = "disabled"
)

type Thinking struct {
	Type ThinkingType `json:"type"`
}

type ChatCompletionRequest struct {
	Model            string         `json:"model"`
	Messages         any            `json:"messages"`
	Stream           bool           `json:"stream,omitempty"`
	StreamOptions    *StreamOptions `json:"stream_options,omitempty"`
	FrequencyPenalty int            `json:"frequency_penalty,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	N                int            `json:"n,omitempty"`
	PresencePenalty  int            `json:"presence_penalty,omitempty"`
	ResponseFormat   string         `json:"response_format,omitempty"`
	Seed             int            `json:"seed,omitempty"`
	Temperature      float64        `json:"temperature,omitempty"`
	TopP             int            `json:"top_p,omitempty"`
	User             string         `json:"user,omitempty"`
	Tools            []interface{}  `json:"tools,omitempty"`
	Thinking         *Thinking      `json:"thinking,omitempty"`
	EnableThinking   *bool          `json:"enable_thinking,omitempty"`
	Modalities       []string       `json:"modalities,omitempty"`   //openrouter
	ImageConfig      any            `json:"image_config,omitempty"` //openrouter
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

type ChatCompletionChoice struct {
	Message ChatCompletionResponseMessage `json:"message"`
}
type ChatCompletionStreamChoice struct {
	Index        int                           `json:"index"`
	Delta        ChatCompletionResponseMessage `json:"delta"`
	FinishReason string                        `json:"finish_reason"`
	Usage        ChatCompletionUsage           `json:"usage"`
}

type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type ChatCompletionResponse struct {
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   ChatCompletionUsage    `json:"usage"`
}
type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Choices []ChatCompletionStreamChoice `json:"choices,omitempty"`
	Usage   ChatCompletionUsage          `json:"usage"`
}

type ChatCompletionStream struct {
	*common.StreamReader[ChatCompletionStreamResponse]
}

var (
	headerData  = []byte("data:")
	errorPrefix = []byte(`data: {"error":`)
)

func (c *ChatCompletionStream) Recv() (ChatCompletionStreamResponse, error) {
	if c.StreamReader.IsFinished {
		return ChatCompletionStreamResponse{}, io.EOF
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
					return *new(ChatCompletionStreamResponse), fmt.Errorf("unmarshal error, %w", c.StreamReader.ErrorResponse.Error())
				}
				return *new(ChatCompletionStreamResponse), readErr
			} else {
				c.StreamReader.IsFinished = true
				return *new(ChatCompletionStreamResponse), io.EOF
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
				return *new(ChatCompletionStreamResponse), writeErr
			}
			emptyMessagesCount++
			if emptyMessagesCount > c.StreamReader.EmptyMessagesLimit {
				return *new(ChatCompletionStreamResponse), errors.New("stream has sent too many empty messages")
			}

			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if strings.TrimSpace(string(noPrefixLine)) == "[DONE]" {
			c.StreamReader.IsFinished = true
			return *new(ChatCompletionStreamResponse), io.EOF
		}

		var response ChatCompletionStreamResponse
		unmarshalErr := basics.JsonDecode(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(ChatCompletionStreamResponse), unmarshalErr
		}

		return response, nil
	}
}
