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
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/gscf"
	"serverless/common/handler/2020/gift"
	"serverless/common/middleware"
)

// 接口版本号
var version = "0.2.0"

// 接口构建时间, 将会在编译时注入
var builtTime = ""

func main() {
	// 使用默认引擎
	r := gin.Default()

	// 添加中间件及处理函数
	// 处理函数要作为最后一个参数传入
	r.Use(middleware.AddVersions(version+builtTime), middleware.CheckOrigin(), gift.GetSurplus)

	// 启动框架, 开始监听请求
	r.Run()

	// 在程序结束时关闭 ORM 的连接
	defer orm.Close()
}
