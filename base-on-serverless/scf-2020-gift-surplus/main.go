/*
   @Time : 2020/4/24 9:09 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
   @Package: surplus
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
	r.Use(handler.AddVersions(version+builtTime), handler.CheckOrigin(), MainHandler)
	r.Run()
}

func MainHandler(c *chaos.Context) {
	o := orm.New()

	count := 0
	o.PostgreSQL.Marketing.Model(structs.EventsGift{}).Where("Name = ? AND Phone IS NULL", c.Param("Name")).Count(&count)

	// 返回数据
	c.JSON(http.StatusOK, chaos.H{"Surplus": count})

	defer func() {
		o.Close() // 在程序结束时关闭 ORM 的连接
	}()
}
