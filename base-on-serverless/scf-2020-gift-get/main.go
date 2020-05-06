/*
   @Time : 2020/4/24 10:23 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
   @Package: get
*/

package main

import (
	"github.com/offcn-jl/cscf"
	"github.com/offcn-jl/cscf/fake-http"
	"github.com/offcn-jl/go-common/database/orm"
	"serverless/common/config"
	"serverless/common/database/orm/structs"
	"serverless/common/middleware"
)

// 接口版本号
var version = "0.1.0"

// 接口构建时间, 将会在编译时注入
var builtTime = ""

func main() {
	// 使用默认引擎
	r := chaos.Default()

	// 添加中间件及处理函数
	// 处理函数要作为最后一个参数传入
	r.Use(middleware.AddVersions(version+builtTime), middleware.CheckOrigin(), MainHandler)

	// 启动框架, 开始监听请求
	r.Run()

	// 在程序结束时关闭 ORM 的连接
	defer orm.Close()
}

// MainHandler 处理函数
// 本接口用来获取奖品
// 如果没有领取过奖品, 将会从奖品表中取出一个奖品然后返回奖品信息
// 如果领取过奖品, 将会直接返回第一次领取到的奖品信息
// 注意 : 本接口使用了事务, 对将要领取奖品所在行进行加锁处理, 在大量未领奖用户同时进行领奖时会因这个加锁逻辑出现数据库连接数被占满而停止响应的情况。所以理论上本接口在大量新用户领奖时, 最大并发数即数据库最大连接数。
func MainHandler(c *chaos.Context) {
	// 查询是否已经领取过奖品
	giftCheckInfo := structs.EventsGift{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	// 已经领取过奖品, 直接返回奖品信息
	if giftCheckInfo.ID != 0 {
		c.JSON(http.StatusOK, chaos.H{"Detail": giftCheckInfo.Detail})
		return
	}

	// 开启事务
	tx := orm.PostgreSQL.Begin()

	// 取出一个奖品，并锁定其所在行
	giftInfo := structs.EventsGift{}
	tx.Set("gorm:query_option", "FOR UPDATE").Where("Name = ? AND Phone IS NULL", c.Param("Name")).First(&giftInfo)

	// 判断是否还有剩余奖品
	if giftInfo.ID == 0 {
		// 没有剩余奖品
		c.JSON(http.StatusForbidden, chaos.H{"Code": -1, "Error": "剩余数量为 0"})
		// 回滚事务
		tx.Rollback()
		return
	}

	// 更新中奖信息
	giftUpdateInfo := structs.EventsGift{
		Phone:      c.Param("Phone"),
		SourceIP:   c.ClientIP(),
		ApiVersion: config.Version + " ( " + version + builtTime + " )",
	}
	tx.Model(structs.EventsGift{}).Where("id = ?", giftInfo.ID).Updates(&giftUpdateInfo)

	// 提交事务
	tx.Commit()

	// 返回信息
	c.JSON(http.StatusOK, chaos.H{"Detail": giftInfo.Detail})
}
