/*
   @Time : 2020/4/24 10:26 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
   @Package: checkout
*/

package main

import (
	"github.com/offcn-jl/chaos-go-scf"
	"github.com/offcn-jl/chaos-go-scf/fake-http"
	"serverless/go-common/database/orm"
	"serverless/go-common/database/orm/structs"
	"serverless/go-common/handler"
)

var version = "0.1.0"
var builtTime = ""

func main() {
	r := chaos.Default()
	r.Use(handler.AddVersions(version+builtTime), MainHandler)
	r.Run()
}

func MainHandler(c *chaos.Context) {
	o := orm.New()

	giftCheckInfo := structs.EventsGift{}
	o.PostgreSQL.Marketing.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	c.JSON(http.StatusOK, giftCheckInfo)

	defer func() {
		o.Close() // 在程序结束时关闭 ORM 的连接
	}()
}
