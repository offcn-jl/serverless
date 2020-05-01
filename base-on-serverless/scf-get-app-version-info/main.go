/*
   @Time : 2020/4/22 3:09 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
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
	r.Use(handler.AddVersions(version+builtTime), MainHandler)
	r.Run()
}

func MainHandler(c *chaos.Context) {
	o := orm.New()
	info := structs.VersionControlInfo{}
	o.PostgreSQL.Marketing.Where("app_id = ?", c.Param("AppID")).Last(&info)

	c.JSON(http.StatusOK, info)

	defer func() {
		o.Close() // 在程序结束时关闭 ORM 的连接
	}()
}
