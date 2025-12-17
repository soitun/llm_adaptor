// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package basics

import (
	"encoding/json"
	"fmt"
	"io"
)

func JsonDecode(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf(`json解析错误,原文:%s,错误:%s`, data, err)
	}
	return nil
}

func JsonEncode(v any) ([]byte, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return buf, fmt.Errorf(`json编码错误,原文:%#v,错误:%s`, v, err)
	}
	return buf, nil
}

func JsonDecoder(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf(`reader读取错误:%s`, err)
	}
	return JsonDecode(data, v)
}

func JsonEncodeStr(v any) (string, error) {
	buf, err := JsonEncode(v)
	return string(buf), err
}
