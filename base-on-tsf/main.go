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
	"tsf/common/handler/2020/05/01"
	"tsf/common/handler/2020/05/02"
	h20200503 "tsf/common/handler/2020/05/03"
	"tsf/common/handler/app"
	"tsf/common/handler/photo"
	"tsf/common/handler/sso/v2"
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
	r.POST("/photo/:Beauty", photo.PostHandler)

	// 应用程序路由组
	appGroup := r.Group("/app")
	{
		// 版本控制路由组
		versionControlGroup := appGroup.Group("/version-control")
		{
			// GetVersion 获取应用程序版本控制信息
			versionControlGroup.GET("/get/:AppID", app.GetVersion)
		}
	}

	// 2020 年路由组 ( Group Year 2020 )
	GY2020 := r.Group("/2020", middleware.CheckOrigin())
	{
		// 5 月路由组 ( Group Month 05 )
		GM05 := GY2020.Group("/05")
		{
			// 活动 01 路由组 ( Group Event 01 )
			// 可用于带计数的预约类活动
			GE01 := GM05.Group("/01")
			{
				// PostSubscribe 处理用户订阅 ( 参与活动 ) 事件
				GE01.POST("/:Name/subscribe/:Phone", h20200501.PostSubscribe)

				// GetCount 获取订阅活动的用户数量
				GE01.GET("/:Name/count", h20200501.GetCount)
			}

			// 活动 02 路由组
			// 礼品 ( 兑换码 ) 发放类
			GE02 := GM05.Group("/02")
			{
				// GetSurplus 获取礼品剩余数量
				GE02.GET("/:Name/surplus", h20200502.GetSurplus)

				// GetGift 获取礼品
				GE02.GET("/:Name/get/:Phone", h20200502.GetGift)

				// GetCheckout 查询领取信息
				GE02.GET("/:Name/checkout/:Phone", h20200502.GetCheckout)

				// PatchConsume 消费礼品
				GE02.PATCH("/:Name/consume/:Phone/detail/:ConsumeDetail", h20200502.PatchConsume)
			}

			// 活动 03 路由
			// 需要进行参与次数计数的活动
			GM05.PATCH("/03/add", h20200503.PatchAdd)
		}
	}

	// SSO 路由组
	ssoV2Group := r.Group("/sso/v2")
	{
		// GetSessionInfo 获取会话信息
		ssoV2Group.GET("/sessions/info/:MID/:Suffix/:Phone", sso.GetSessionInfo)

		// GetSuffixList 获取后缀花名册
		ssoV2Group.GET("/suffix/list", sso.GetSuffixList)
	}

	// 启动服务
	err := r.Run()
	if err != nil {
		logger.Error(err)
	}

	// 退出前关闭数据库连接
	defer orm.Close()
}
