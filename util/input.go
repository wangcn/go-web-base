package util

import "github.com/gin-gonic/gin"

// GetParam
func GetParam(ctx *gin.Context, name string, def ...string) string {
	if value, exists := ctx.GetQuery(name); exists {
		return value
	}

	if value, exists := ctx.GetPostForm(name); exists {
		return value
	}

	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func GetParams(ctx *gin.Context) map[string]string {

	input := make(map[string]string)

	if queries := ctx.Request.URL.Query(); queries != nil {
		for key, item := range queries {
			if len(item) == 0 {
				input[key] = ""
			} else {
				input[key] = item[0]
			}
		}
	}

	_ = ctx.Request.ParseForm()
	if ctx.Request.Form != nil {
		for key, item := range ctx.Request.Form {
			if len(item) == 0 {
				input[key] = ""
			} else {
				input[key] = item[0]
			}
		}
	}

	_ = ctx.Request.ParseMultipartForm(2 * 1024 * 1024)
	if ctx.Request.PostForm != nil {
		for key, item := range ctx.Request.PostForm {
			if len(item) == 0 {
				input[key] = ""
			} else {
				input[key] = item[0]
			}
		}
	}

	return input
}
