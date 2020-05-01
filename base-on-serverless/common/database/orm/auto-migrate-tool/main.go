/*
   @Time : 2020/4/24 5:30 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
*/

package main

import (
	"serverless/common/database/orm"
	"serverless/common/database/orm/structs"
)

func main() {
	o := orm.New()
	o.PostgreSQL.Marketing.AutoMigrate(
		&structs.VersionControlInfo{},
		&structs.EventsGift{},
	)
}
