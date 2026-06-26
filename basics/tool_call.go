// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package basics

import "encoding/json"

type ToolCall struct {
	Index    *int           `json:"index,omitempty"`
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Function FunctionCall   `json:"function"`
	Extra    map[string]any `json:"extra,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

func (f *FunctionCall) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name      string          `json:"name,omitempty"`
		Arguments json.RawMessage `json:"arguments,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	f.Name = raw.Name
	if len(raw.Arguments) == 0 || string(raw.Arguments) == "null" {
		f.Arguments = ""
		return nil
	}
	if err := json.Unmarshal(raw.Arguments, &f.Arguments); err == nil {
		return nil
	}
	f.Arguments = string(raw.Arguments)
	return nil
}

// NewFunctionToolCall creates a function tool call with the shared basics shape.
func NewFunctionToolCall(id string, name string, arguments string) ToolCall {
	return ToolCall{
		ID:   id,
		Type: "function",
		Function: FunctionCall{
			Name:      name,
			Arguments: arguments,
		},
	}
}

// ToolCalls keeps raw model tool calls and exposes legacy projections.
type ToolCalls []ToolCall

type FunctionToolCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// FunctionToolCalls returns the legacy FunctionToolCall view.
func (toolCalls ToolCalls) FunctionToolCalls() []FunctionToolCall {
	functionToolCalls := make([]FunctionToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		if toolCall.Type != "" && toolCall.Type != "function" {
			continue
		}
		functionToolCalls = append(functionToolCalls, FunctionToolCall{
			Name:      toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		})
	}
	if len(functionToolCalls) == 0 {
		return nil
	}
	return functionToolCalls
}
