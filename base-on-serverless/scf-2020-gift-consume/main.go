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
// 本接口用来消费奖品
// 奖品未被消费时, 将奖品消费信息保存后, 返回 HTTP 状态码 200
// 奖品已经被消费时, 返回奖品消费详情, HTTP 状态码为 403
func MainHandler(c *chaos.Context) {
	// 获取奖品详情
	giftCheckInfo := structs.EventsGift{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	// 判断奖品是否已经被消费
	if giftCheckInfo.ConsumeDetail != "" {
		// 奖品已经被消费
		c.JSON(http.StatusForbidden, chaos.H{"Code": -1, "Error": "礼品已被消费，详情 : " + giftCheckInfo.ConsumeDetail})
	} else {
		// 奖品未被消费, 更新奖品信息
		orm.PostgreSQL.Model(&structs.EventsGift{}).Where("name = ? AND phone = ?", c.Param("Name"), c.Param("Phone")).Update("consume_detail", c.Param("ConsumeDetail"))
		// 隐式返回 HTTP 状态码 200
	}
}
