// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package baidu

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode      int    `json:"error_code"`
	ErrorMessage   string `json:"error_msg"`
	HTTPStatusCode int    `json:"-"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.ErrorCode > 0 || e.HTTPStatusCode < http.StatusOK || e.HTTPStatusCode >= http.StatusBadRequest {
		return errors.New(fmt.Sprintf("Baidu request error, http status code: %d, err code: %d, message: %s", e.HTTPStatusCode, e.ErrorCode, e.ErrorMessage))
	}
	return nil
}
