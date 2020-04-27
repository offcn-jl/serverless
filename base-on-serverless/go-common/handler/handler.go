/*
   @Time : 2020/4/22 5:04 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : handler
   @Software: GoLand
*/

package handler

import (
	"github.com/offcn-jl/chaos-go-scf"
	"github.com/offcn-jl/chaos-go-scf/fake-http"
	"serverless/go-common/codes"
	"serverless/go-common/configer"
	"serverless/go-common/database"
	"serverless/go-common/database/orm"
	"serverless/go-common/logger"
	"strings"
)

// 向响应头添加版本信息
func AddVersions(apiVersion string) chaos.HandlerFunc {
	return func(c *chaos.Context) {
		c.Header(configer.Conf.Project+"-Version", configer.Conf.Version)
		c.Header(configer.Conf.Project+"-Configer-Version", configer.Version)
		c.Header(configer.Conf.Project+"-Logger-Version", logger.Version)
		c.Header(configer.Conf.Project+"-Database-Version", database.Version)
		c.Header(configer.Conf.Project+"-Database-ORM-Version", orm.Version)
		c.Header(configer.Conf.Project+"-Api-Version", apiVersion)
		c.Next()
	}
}

// 跨域检查与跨域头添加 中间件
func CheckOrigin() chaos.HandlerFunc {
	return func(c *chaos.Context) {
		// 跨域校验
		allowOrigins := configer.GetString("AllowOrigins", "")
		allowOriginsArray := strings.Split(allowOrigins, ",")
		pass := false
		for _, origin := range allowOriginsArray {
			// 遍历配置中的跨域头，寻找匹配项
			if c.GetHeader("origin") == origin {
				c.Header("Access-Control-Allow-Origin", origin)
				pass = true
				// 只要有一个跨域头匹配就跳出循环
				break
			}
		}

		if !pass {
			c.JSON(http.StatusForbidden, chaos.H{"Code": codes.NotCertifiedCORS, "Error": codes.ErrorText(codes.NotCertifiedCORS)})
			c.Abort() // 出错后结束请求
		}

		// 通过跨域校验后，放行所有 OPTIONS 方法，并添加按照客户端的请求添加 Allow Headers
		//if c.Request.Method == "OPTIONS" {
		//	// 请求首部  Access-Control-Request-Headers 出现于 preflight request （预检请求）中，用于通知服务器在真正的请求中会采用哪些请求首部。
		//	c.Header("Access-Control-Allow-Headers", c.GetHeader("Access-Control-Request-Headers")) // 放行预检请求通知的请求首部。
		//	// https://cloud.tencent.com/developer/section/1189896
		//	c.Header("Access-Control-Allow-Methods", c.GetHeader("Access-Control-Request-Method")) // 放行预检请求通知的请求首部。
		//	c.AbortWithStatus(http.StatusNoContent)
		//}
	}
}
