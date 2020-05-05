/*
   @Time : 2020/4/24 5:30 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
*/

package main

import (
	"github.com/offcn-jl/go-common/database/orm"
	"serverless/common/database/orm/structs"
)

// 数据库结构自动迁移工具
func main() {
	// 自动迁移数据库结构
	orm.PostgreSQL.AutoMigrate(
		&structs.VersionControlInfo{},
		&structs.EventsGift{},
	)
	// 退出前关闭数据库连接
	defer orm.Close()
}
