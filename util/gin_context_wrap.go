package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Context 与 gin.Context 组合，可以扩展方法，其他用法不变
type Context struct {
	*gin.Context
}

// HandlerFunc 自定义的 Handler 结构，传递自定义的 Context，自定义的返回
type HandlerFunc func(ctx *Context) *Resp

// Handle 对原 gin.handler 的封装，内部调用 HandlerFunc，并处理返回
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := &Context{
			ginCtx,
		}
		resp := h(ctx)
		if resp != nil {
			ctx.JSON(200, resp)
		}
	}
}

// Output 对自定义 handler 返回的封装，不再需要 return 单独写一行
// 1. 正常返回 JSON
// 2. 自定义 Error 返回 JSON，内部有自定义的 code 和 msg
// 3. 系统 error 处理，直接返回【系统错误】JSON，并记录 err 日志和堆栈信息，收集 err 并报警
// 4. 当返回 nil 时，不做任何处理，一般用于输出其他自定义格式，例如图片等
func (ctx *Context) Output(data interface{}) *Resp {
	if err, ok := data.(error); ok {
		return RespError(err)
	}
	if data == nil {
		return nil
	}
	return RespSuccess(data)
}

func (ctx *Context) SetBinaryHeader(filename string) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
}

// Input 直接获取 GET 或 POST 里的参数
func (ctx *Context) Input(name string, def ...string) string {
	return GetParam(ctx.Context, name, def...)
}

// Inputs 直接获取 GET 或 POST 里的所有参数
func (ctx *Context) Inputs() map[string]string {
	return GetParams(ctx.Context)
}
