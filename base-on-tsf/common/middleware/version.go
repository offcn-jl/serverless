/*
   @Time : 2020/5/4 9:50 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : version
   @Software: GoLand
*/

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common"
	"tsf/common/config"
)

// AddVersions 用于向响应头添加版本信息
func AddVersions() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Common-Version", common.Version)
		c.Header("X-"+config.Project+"-Version", config.Version)
		c.Next()
	}
}
