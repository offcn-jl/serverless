/*
   @Time : 2020/4/24 9:48 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : 2020-gift
   @Software: GoLand
*/

package structs

import "github.com/jinzhu/gorm"

/**
 * Gift 礼品
 */
type EventsGift struct {
	gorm.Model
	Name          string `gorm:"not null"` // 礼品名称
	Detail        string // 礼品详情
	Phone         string // 领取人手机号
	ConsumeDetail string // 消费详情
	// 平台日志
	SourceIP   string
	ApiVersion string
}
