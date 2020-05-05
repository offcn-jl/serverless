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
	"github.com/offcn-jl/cscf"
	"github.com/offcn-jl/cscf/fake-http"
	"github.com/offcn-jl/go-common/database/orm"
	"serverless/common/database/orm/structs"
	"serverless/common/handler"
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
	r.Use(handler.AddVersions(version+builtTime), handler.CheckOrigin(), MainHandler)

	// 启动框架, 开始监听请求
	r.Run()

	// 在程序结束时关闭 ORM 的连接
	defer orm.Close()
}

// MainHandler 处理函数
// 本接口用来获取礼品的剩余数量
func MainHandler(c *chaos.Context) {
	// 定义计数变量
	count := 0

	// 取出数量
	orm.PostgreSQL.Model(structs.EventsGift{}).Where("Name = ? AND Phone IS NULL", c.Param("Name")).Count(&count)

	// 返回数据
	c.JSON(http.StatusOK, chaos.H{"Surplus": count})
}
