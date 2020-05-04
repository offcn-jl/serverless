/*
   @Time : 2020-04-13 15:28
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
*/
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/go-common/logger"
	"tsf/common/database/orm/structs"
	"tsf/common/handler"
	"tsf/common/handler/2020/05/01"
	"tsf/common/middleware"
)

// 主函数
func main() {
	// 自动迁移数据库结构
	structs.AutoMigrate()

	// 初始化 Gin 引擎，采用默认配置
	r := gin.Default()

	// 配置中间件
	r.Use(middleware.AddVersions())

	// 照片处理路由
	r.POST("/photo/:Beauty", handler.Photo)

	// 2020 年路由组 ( Group Year 2020 )
	GY2020 := r.Group("/2020", middleware.CheckOrigin())
	{
		// 4 月路由组 ( Group Month 05 )
		GM04 := GY2020.Group("/05")
		{
			// 活动 01 路由组 ( Group Event 01 )
			GE01 := GM04.Group("/01")
			{
				// PostSubscribe 处理用户订阅 ( 参与活动 ) 事件
				GE01.POST("/:Name/subscribe/:Phone", h20200501.PostSubscribe)

				// GetCount 获取订阅活动的用户数量
				GE01.GET("/:Name/count", h20200501.GetCount)
			}
		}
	}

	// 启动服务
	err := r.Run()
	if err != nil {
		logger.Error(err)
	}

	// 退出前关闭数据库连接
	defer orm.Close()
}
