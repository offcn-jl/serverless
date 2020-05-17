/*
   @Time : 2020/5/4 11:32 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : structs
   @Software: GoLand
*/

package structs

import (
	"fmt"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/go-common/logger"
	"time"
	"tsf/common/config"
)

// AutoMigrate 提供表结构自动迁移功能
func AutoMigrate() {
	// 判断是否符合表结构自动迁移条件
	if len(config.Version) == 37 {
		builtTime, _ := time.ParseInLocation("2006/01/02 15:04:05", config.Version[16:35], time.Local) // 使用parseInLocation将字符串格式化返回本地时区时间
		// 判断构建时间是否超过一小时
		if time.Since(builtTime).Hours() > 1 {
			// 超过一小时, 不需要进行迁移
			logger.Log("ORM : 构建于 " + fmt.Sprint(time.Since(builtTime)) + " 之前, 跳过自动迁移数据库表结构步骤.")
		} else {
			// 没有超过一小时, 进行表结构自动迁移
			autoMigrate()
		}
	} else {
		// 版本信息中没有包含构建时间, 直接进行表结构自动迁移
		autoMigrate()
	}
}

// autoMigrate 可以进行表结构迁移工作
func autoMigrate() {
	logger.Log("ORM : 开始自动迁移数据库表结构 ...")
	orm.PostgreSQL.AutoMigrate(
		&E20200501{},
		&E20200502{},
		&VersionControlInfo{},
		&SingleSignOnLoginModule{},
		&SingleSignOnVerificationCode{},
		&SingleSignOnUser{},
		&SingleSignOnSession{},
		&SingleSignOnSuffix{},
		&SingleSignOnOrganization{},
		&SingleSignOnCRMRoundLog{},
	)
}
