// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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
	var toolCalls basics.ToolCalls
	if len(res.Choices) > 0 {
		res.Result = res.Choices[0].Delta.Content
		res.ReasoningContent = res.Choices[0].Delta.ReasoningContent
		toolCalls = res.Choices[0].Delta.ToolCalls
	} else if len(res.FunctionCall.Name) > 0 || len(res.FunctionCall.Arguments) > 0 {
		toolCalls = append(toolCalls, basics.NewFunctionToolCall("", res.FunctionCall.Name, res.FunctionCall.Arguments))
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
		ToolCalls:         toolCalls,
		FunctionToolCalls: toolCalls.FunctionToolCalls(),
		PromptToken:       res.Usage.PromptTokens,
		CompletionToken:   res.Usage.CompletionTokens,
	}, nil
}
