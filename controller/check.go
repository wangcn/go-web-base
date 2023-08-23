package controller

import (
	"go.uber.org/zap"

	"mybase/util"
)

func Ping(ctx *util.Context) *util.Resp {
	util.Log().With(zap.String("aa", "aa1")).Debug("ping check")
	return ctx.Output(333)
}
