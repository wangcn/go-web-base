package middleware

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mybase/util"
)

func AccessLog() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path
		// 跳过健康检查
		if path == "/ping" || path == "/hc" || path == "/metrics" {
			return
		}

		start := time.Now()

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			body = make([]byte, 0)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		ctx.Next()

		stop := time.Since(start)
		// latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		latency := stop.Milliseconds()
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}

		if len(body) > 1024 {
			body = body[:1024]
		}

		util.Log().With(
			zap.Int64("latency", latency),
			zap.String("hoetname", hostName),
			zap.String("method", ctx.Request.Method),
			zap.Int("status_code", ctx.Writer.Status()),
			zap.String("client_ip", ctx.ClientIP()),
			zap.String("ua", ctx.Request.UserAgent()),
			zap.String("referer", ctx.Request.Referer()),
			zap.Int64("response_size", int64(ctx.Writer.Size())),
			zap.String("path", path),
			zap.String("content_type", ctx.GetHeader("Content-Type")),
			zap.Any("params", util.GetParams(ctx)),
			zap.ByteString("body", body),
		).Info("access_log")
	}
}
