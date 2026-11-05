// Package errno 定义业务错误码。
//
// 通过 CodeError 接口统一携带错误码和错误信息。
package errno

import (
	"strconv"

	"github.com/pkg/errors"
)

// CodeError 可携带错误码的错误。
type CodeError interface {
	error
	ErrCode() int64
}

type errnoError struct {
	code int64
	msg  string
}

// NewCodeError 创建错误码错误。
func NewCodeError(errCode int64, errMsg string) *errnoError {
	if errCode < 0 {
		panic("error code cannot less zero")
	}
	return &errnoError{code: errCode, msg: errMsg}
}

func (e *errnoError) ErrCode() int64 {
	return e.code
}

func (e *errnoError) Error() string {
	return e.msg
}

// Cause 提取 err 链上的 CodeError，非业务错误统一转为 ServerErr。
func Cause(err error) CodeError {
	if err == nil {
		return OK
	}
	if ce, ok := errors.Cause(err).(CodeError); ok {
		return ce
	}
	return String(err.Error())
}

// String 将字符串转为 CodeError；无法解析为数字时返回 ServerErr。
func String(e string) CodeError {
	if e == "" {
		return OK
	}
	i, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return ServerErr
	}
	return NewCodeError(i, "")
}
