// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

import (
	"regexp"
	"strings"

	"github.com/zhimaAi/llm_adaptor/api/baidu"
	"github.com/zhimaAi/llm_adaptor/basics"
)

type BaiduStreamResult struct {
	*baidu.ChatCompletionStream
}

func (r *BaiduStreamResult) Read() (ZhimaChatCompletionResponse, error) {
	res, err := r.Recv()
	if err != nil {
		return ZhimaChatCompletionResponse{}, err
	}
	var functionToolCalls []FunctionToolCall
	if len(res.Choices) > 0 {
		res.Result = res.Choices[0].Delta.Content
		res.ReasoningContent = res.Choices[0].Delta.ReasoningContent
		for _, toolCall := range res.Choices[0].Delta.ToolCalls {
			functionToolCalls = append(functionToolCalls, FunctionToolCall{
				Name:      toolCall.Function.Name,
				Arguments: toolCall.Function.Arguments,
			})
		}
	} else if len(res.FunctionCall.Name) > 0 || len(res.FunctionCall.Arguments) > 0 {
		functionToolCalls = append(functionToolCalls, FunctionToolCall{
			Name:      res.FunctionCall.Name,
			Arguments: res.FunctionCall.Arguments,
		})
		if strings.Contains(res.FunctionCall.Thoughts, `prompt`) {
			arguments := make(map[string]string)
			err := basics.JsonDecode([]byte(res.FunctionCall.Arguments), &arguments)
			if err != nil {
				return ZhimaChatCompletionResponse{}, err
			}
			for k := range arguments {
				arguments[k] = ``
			}
			res.FunctionCall.Arguments, _ = basics.JsonEncodeStr(arguments)

			re := regexp.MustCompile(`"prompt":\s*"([^"]*)"`)
			matches := re.FindStringSubmatch(res.FunctionCall.Thoughts)
			if len(matches) > 1 {
				res.Result = matches[1]
			} else {
				var argumentKeys []string
				for _, argumentKey := range arguments {
					argumentKeys = append(argumentKeys, argumentKey)
				}
				res.Result = `请提供必须参数: ` + strings.Join(argumentKeys, `, `)
			}
		}
	}
	return ZhimaChatCompletionResponse{
		Result:            res.Result,
		ReasoningContent:  res.ReasoningContent,
		FunctionToolCalls: functionToolCalls,
		PromptToken:       res.Usage.PromptTokens,
		CompletionToken:   res.Usage.CompletionTokens,
	}, nil
}
