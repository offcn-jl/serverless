/*
   @Time : 2020/4/24 10:27 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
   @Package: consume
*/

package main

import (
	"github.com/offcn-jl/cscf"
	"github.com/offcn-jl/cscf/fake-http"
	"serverless/common/database/orm"
	"serverless/common/database/orm/structs"
	"serverless/common/handler"
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

	giftCheckInfo := structs.EventsGift{}
	o.PostgreSQL.Marketing.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)
	if giftCheckInfo.ConsumeDetail != "" {
		c.JSON(http.StatusForbidden, chaos.H{"Code": -1, "Error": "礼品已被消费，详情 : " + giftCheckInfo.ConsumeDetail})
	} else {
		o.PostgreSQL.Marketing.Model(&structs.EventsGift{}).Where("name = ? AND phone = ?", c.Param("Name"), c.Param("Phone")).Update("consume_detail", c.Param("ConsumeDetail"))
	}

	defer func() {
		o.Close() // 在程序结束时关闭 ORM 的连接
	}()
}
