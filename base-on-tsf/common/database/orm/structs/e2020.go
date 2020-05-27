/*
   @Time : 2020/5/4 11:08 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : e20200401
   @Software: GoLand
*/

package structs

import "github.com/jinzhu/gorm"

// 可用于带计数的预约类活动
type E20200501 struct {
	gorm.Model
	Name           string `gorm:"not null"` // 活动名称
	Phone          string `gorm:"not null"` // 参与者手机号码, 用于识别参与者
	ProjectVersion string `gorm:"not null"` // 项目版本
}

// 可用于礼品 ( 兑换码 ) 发放类活动
type E20200502 struct {
	gorm.Model
	Name          string `gorm:"not null"` // 礼品名称
	Detail        string // 礼品详情
	Phone         string // 领取人手机号
	ConsumeDetail string // 消费详情
	// 平台日志
	SourceIP string
}

// 可用于需要进行参与次数计数的活动
type E20200503 struct {
	gorm.Model
	Event string `gorm:"not null" json:"Event" binding:"required"` // 活动标识
	Phone string `gorm:"not null" json:"Phone" binding:"required"` // 参与者手机号码 ( 参与者标识 )
	// 平台日志
	SourceIP string // 用户 IP
}
