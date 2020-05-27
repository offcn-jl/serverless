/*
   @Time : 2020/5/15 11:45 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
   @Package: main
*/

package main

import (
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/gscf"
	"github.com/offcn-jl/serverless-apis/base-on-serverless/common/handler/sso/v2"
	"github.com/offcn-jl/serverless-apis/base-on-serverless/common/middleware"
)

// 接口版本号
var version = "0.1.1"

// 接口构建时间, 将会在编译时注入
var builtTime = ""

// main 主函数, 作为程序的入口
func main() {
	// 使用默认引擎
	r := gin.Default()

	// 添加中间件及处理函数
	// 处理函数要作为最后一个参数传入
	r.Use(middleware.AddVersions(version+builtTime), middleware.CheckOrigin(), sso.PostSignUp)

	// 启动框架, 开始监听请求
	r.Run()

	// 在程序结束时关闭 ORM 的连接
	defer orm.Close()
}
