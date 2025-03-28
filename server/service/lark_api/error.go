package lark_api

import (
	"fmt"
)

func NewErrResponseNotSuccess(code int, msg string) error {
	return fmt.Errorf("response is not success, code: %d, msg: %s", code, msg)
}
