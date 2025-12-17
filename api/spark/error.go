// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package spark

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	Header ErrHeader `json:"header"`
}

type ErrHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Status  int    `json:"status"`
}

func (e *ErrorResponse) Error() error {
	return errors.New(fmt.Sprintf("Spark Request error, code: %d, message: %s", e.Header.Code, e.Header.Message))
}
