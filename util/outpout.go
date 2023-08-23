package util

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Resp struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	ShowErr   bool        `json:"show_err"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func NewResp(code int, msg string, showErr bool) *Resp {
	resp := &Resp{
		Code:      code,
		Message:   msg,
		ShowErr:   showErr,
		Timestamp: time.Now().Unix(),
		Data:      gin.H{},
	}
	return resp
}

// RespSuccess 正常返回数据
func RespSuccess(data interface{}) *Resp {
	resp := NewResp(0, "成功", false)
	resp.Data = data
	return resp
}

// RespError 错误返回数据
func RespError(err error) *Resp {
	// 业务错误
	if e, ok := err.(*Error); ok {
		return NewResp(e.Code, e.Msg, e.ShowErr)
	}

	// 服务器端错误，屏蔽错误把错误记录到日志，并报警
	err = errors.WithStack(err)

	Log().With(zap.Error(err)).Error("系统错误")

	return NewResp(3, "系统错误", false)
}
