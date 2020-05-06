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
	"github.com/offcn-jl/go-common/database/orm"
	"serverless/common/database/orm/structs"
	"serverless/common/middleware"
)

// 接口版本号
var version = "0.1.0"

// 接口构建时间, 将会在编译时注入
var builtTime = ""

// main 主函数, 作为程序的入口
func main() {
	// 使用默认引擎
	r := chaos.Default()

	// 添加中间件及处理函数
	// 处理函数要作为最后一个参数传入
	r.Use(middleware.AddVersions(version+builtTime), MainHandler)

	// 启动框架, 开始监听请求
	r.Run()

	// 在程序结束时关闭 ORM 的连接
	defer orm.Close()
}

// MainHandler 处理函数
// 本接口用来获取指定 APP 的版本控制信息
// 信息包括版本号、发布时间、更新时间、下载地址
func MainHandler(c *chaos.Context) {
	// 定义版本信息结构
	info := structs.VersionControlInfo{}

	// 获取版本信息
	orm.PostgreSQL.Where("app_id = ?", c.Param("AppID")).Last(&info)

	// 返回数据
	c.JSON(http.StatusOK, info)
}
