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
	"serverless/go-common/configer"
	"serverless/go-common/database"
	"serverless/go-common/database/orm"
	"serverless/go-common/logger"
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
