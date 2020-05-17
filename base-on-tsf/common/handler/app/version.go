/*
   @Time : 2020/5/17 1:51 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : version
   @Software: GoLand
*/

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common/database/orm"
	"net/http"
	"tsf/common/database/orm/structs"
)

// GetVersion 获取应用程序版本控制信息接口的处理函数
// 版本控制信息包括版本号、发布时间、更新时间、下载地址
func GetVersion(c *gin.Context) {
	// 定义版本信息结构
	info := structs.VersionControlInfo{}

	// 获取版本信息
	orm.PostgreSQL.Where("app_id = ?", c.Param("AppID")).Last(&info)

	// 返回数据
	c.JSON(http.StatusOK, info)
}
