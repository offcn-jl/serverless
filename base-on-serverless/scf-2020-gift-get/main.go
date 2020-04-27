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
	"github.com/offcn-jl/chaos-go-scf"
	"github.com/offcn-jl/chaos-go-scf/fake-http"
	"serverless/go-common/configer"
	"serverless/go-common/database/orm"
	"serverless/go-common/database/orm/structs"
	"serverless/go-common/handler"
)

var version = "0.1.0"
var builtTime = ""

func main() {
	r := chaos.Default()
	r.Use(handler.AddVersions(version+builtTime), handler.CheckOrigin(), MainHandler)
	r.Run()
}

func MainHandler(c *chaos.Context) {
	o := orm.New()

	// 查询是否已经领取过奖品
	giftCheckInfo := structs.EventsGift{}
	o.PostgreSQL.Marketing.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	if giftCheckInfo.ID != 0 {
		c.JSON(http.StatusOK, chaos.H{"Detail": giftCheckInfo.Detail})
		return
	}

	// 开启事务
	tx := o.PostgreSQL.Marketing.Begin()

	// 取出一个奖品，并锁定其所在行
	giftInfo := structs.EventsGift{}
	tx.Set("gorm:query_option", "FOR UPDATE").Where("Name = ? AND Phone IS NULL", c.Param("Name")).First(&giftInfo)

	if giftInfo.ID == 0 {
		c.JSON(http.StatusForbidden, chaos.H{"Code": -1, "Error": "剩余数量为 0"})
		// 回滚
		tx.Rollback()
		return
	}

	// 更新中奖信息
	giftUpdateInfo := structs.EventsGift{
		Phone:      c.Param("Phone"),
		SourceIP:   c.ClientIP(),
		ApiVersion: configer.Conf.Version,
	}
	tx.Model(structs.EventsGift{}).Where("id = ?", giftInfo.ID).Updates(&giftUpdateInfo)

	// 提交
	tx.Commit()

	// 返回信息
	c.JSON(http.StatusOK, chaos.H{"Detail": giftInfo.Detail})

	defer func() {
		o.Close() // 在程序结束时关闭 ORM 的连接
	}()
}
