package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58/base58"

	"mybase/util"
)

func RequestID() gin.HandlerFunc {
	node := util.NewNode()
	return func(c *gin.Context) {
		c.Set("request_id", base58.Encode(node.Generate()))
		c.Next()
	}
}
