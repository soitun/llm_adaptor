// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package voyage

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	Detail         string `json:"detail"`
	HTTPStatusCode int    `json:"-"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("error, status code: %d, detail: %s", e.HTTPStatusCode, e.Detail))
	}
	return errors.New(fmt.Sprintf("error, message: %s", e.Detail))
}
