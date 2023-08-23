package person

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mybase/service/user"
	"mybase/util"
)

type profileInfoReq struct {
	Id int64 `form:"id" binding:"required"`
}

// 接口
func ProfileInfo(ctx *util.Context) *util.Resp {
	// 获取参数，并验证
	req := new(profileInfoReq)
	if err := ctx.ShouldBind(req); err != nil {
		util.Log().With(zap.Error(err)).Debug("check params")
		return ctx.Output(util.ErrParams)
	}

	profileSvs := &user.Profile{}

	info, err := profileSvs.GetInfoById(ctx, req.Id)

	if err != nil {
		return ctx.Output(err)
	}

	result := gin.H{
		"id":   info.Id,
		"name": info.Name,
	}

	return ctx.Output(result)
}
