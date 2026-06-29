// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package adaptor

import (
	"errors"
	"io"
	"strings"
)

const (
	thinkStartTag = "<think>"
	thinkEndTag   = "</think>"
)

type thinkTagExtractor struct {
	inThink bool
	pending string
}

func (e *thinkTagExtractor) Process(content string) (result string, reasoningContent string) {
	e.pending += content
	for {
		tag := thinkStartTag
		dst := &result
		if e.inThink {
			tag = thinkEndTag
			dst = &reasoningContent
		}

		idx := strings.Index(e.pending, tag)
		if idx >= 0 {
			*dst += e.pending[:idx]
			e.pending = e.pending[idx+len(tag):]
			e.inThink = !e.inThink
			continue
		}

		emitLen := len(e.pending) - longestSuffixPrefixLen(e.pending, tag)
		*dst += e.pending[:emitLen]
		e.pending = e.pending[emitLen:]
		return result, reasoningContent
	}
}

func (e *thinkTagExtractor) Flush() (result string, reasoningContent string) {
	if e.pending == "" {
		return "", ""
	}
	if e.inThink {
		reasoningContent = e.pending
	} else {
		result = e.pending
	}
	e.pending = ""
	e.inThink = false
	return result, reasoningContent
}

func longestSuffixPrefixLen(s, prefix string) int {
	maxLen := len(prefix) - 1
	if len(s) < maxLen {
		maxLen = len(s)
	}
	for n := maxLen; n > 0; n-- {
		if strings.HasSuffix(s, prefix[:n]) {
			return n
		}
	}
	return 0
}

func normalizeThinkTaggedResponse(resp ZhimaChatCompletionResponse) ZhimaChatCompletionResponse {
	if resp.ReasoningContent != "" || !strings.Contains(resp.Result, thinkStartTag) {
		return resp
	}
	var extractor thinkTagExtractor
	result, reasoningContent := extractor.Process(resp.Result)
	flushResult, flushReasoningContent := extractor.Flush()
	resp.Result = result + flushResult
	resp.ReasoningContent = reasoningContent + flushReasoningContent
	return resp
}

func readThinkTaggedStream(stream ZhimaStreamResult, extractor *thinkTagExtractor, eofFlushed *bool) (ZhimaChatCompletionResponse, error) {
	if *eofFlushed {
		return ZhimaChatCompletionResponse{}, io.EOF
	}
	resp, err := stream.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			result, reasoningContent := extractor.Flush()
			*eofFlushed = true
			if result != "" || reasoningContent != "" {
				return ZhimaChatCompletionResponse{
					Result:           result,
					ReasoningContent: reasoningContent,
				}, nil
			}
		}
		return resp, err
	}
	if resp.ReasoningContent != "" || resp.Result == "" {
		return resp, nil
	}
	resp.Result, resp.ReasoningContent = extractor.Process(resp.Result)
	return resp, nil
}
