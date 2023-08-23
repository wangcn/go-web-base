package util

import (
	"fmt"
	"strings"
)

// 用户自定义错误
var (
	ErrSystem   = NewErr(3, "系统错误", false)
	ErrNotFound = NewErr(4, "资源未找到", true)
	ErrParams   = NewErr(7, "参数错误", true)
	ErrSign     = NewErr(12, "签名错误", true)
)

func Err2Code(err error) int64 {
	customErr := err.(*Error)
	return int64(customErr.Code)
}

// 是否为表不存在的错误
func IsTableNoExistErr(err error) bool {
	if strings.Contains(err.Error(), "Error 1146") {
		return true
	}
	return false
}

// Error 对系统 error 的封装，code 与错误消息一次设置，并且可以使用比较
type Error struct {
	Code    int
	Msg     string
	ShowErr bool
	error
}

// NewErr 新建一个自定义的 Error
func NewErr(code int, msg string, showErr bool) error {
	err := &Error{
		Code:    code,
		Msg:     msg,
		ShowErr: showErr,
	}
	return err
}

// Error 打印时显示的内容
func (e *Error) Error() string {
	return fmt.Sprintf("code=%d, msg=%s, show_err=%t", e.Code, e.Msg, e.ShowErr)
}
