/*
   @Time : 2020/5/4 11:08 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : e20200401
   @Software: GoLand
*/

package structs

import "github.com/jinzhu/gorm"

// 20200501
// 可用于带计数的预约类活动
type E20200501 struct {
	gorm.Model
	Name           string `gorm:"not null"` // 活动名称
	Phone          string `gorm:"not null"` // 参与者手机号码, 用于识别参与者
	ProjectVersion string `gorm:"not null"` // 项目版本
}
