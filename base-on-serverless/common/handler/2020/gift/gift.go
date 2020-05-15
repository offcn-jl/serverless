/*
   @Time : 2020/5/12 11:08 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : gift
   @Software: GoLand
*/

package gift

import (
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/gscf"
	"github.com/offcn-jl/gscf/fake-http"
	"serverless/common/database/orm/structs"
)

// GetCheckout 是查询获奖信息接口的处理函数
func GetCheckout(c *gin.Context) {
	giftCheckInfo := structs.EventsGift{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	c.JSON(http.StatusOK, giftCheckInfo)
}

// PatchConsume 是消费奖品接口的处理函数
// 奖品未被消费时, 将奖品消费信息保存后, 返回 HTTP 状态码 200
// 奖品已经被消费时, 返回奖品消费详情, HTTP 状态码为 403
func PatchConsume(c *gin.Context) {
	// 获取奖品详情
	giftCheckInfo := structs.EventsGift{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	// 判断奖品是否已经被消费
	if giftCheckInfo.ConsumeDetail != "" {
		// 奖品已经被消费
		c.JSON(http.StatusForbidden, gin.H{"Code": -1, "Error": "礼品已被消费，详情 : " + giftCheckInfo.ConsumeDetail})
	} else {
		// 奖品未被消费, 更新奖品信息
		orm.PostgreSQL.Model(&structs.EventsGift{}).Where("name = ? AND phone = ?", c.Param("Name"), c.Param("Phone")).Update("consume_detail", c.Param("ConsumeDetail"))
		// 隐式返回 HTTP 状态码 200
	}
}

// GetGift 是获取奖品接口的处理函数
// 如果没有领取过奖品, 将会从奖品表中取出一个奖品然后返回奖品信息
// 如果领取过奖品, 将会直接返回第一次领取到的奖品信息
// 注意 : 本接口使用了事务, 对将要领取奖品所在行进行加锁处理, 在大量未领奖用户同时进行领奖时会因这个加锁逻辑出现数据库连接数被占满而停止响应的情况。所以理论上本接口在大量新用户领奖时, 最大并发数即数据库最大连接数。
func GetGift(c *gin.Context) {
	// 查询是否已经领取过奖品
	giftCheckInfo := structs.EventsGift{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	// 已经领取过奖品, 直接返回奖品信息
	if giftCheckInfo.ID != 0 {
		c.JSON(http.StatusOK, gin.H{"Detail": giftCheckInfo.Detail})
		return
	}

	// 开启事务
	tx := orm.PostgreSQL.Begin()

	// 取出一个奖品，并锁定其所在行
	giftInfo := structs.EventsGift{}
	tx.Set("gorm:query_option", "FOR UPDATE").Where("Name = ? AND Phone IS NULL", c.Param("Name")).First(&giftInfo)

	// 判断是否还有剩余奖品
	if giftInfo.ID == 0 {
		// 没有剩余奖品
		c.JSON(http.StatusForbidden, gin.H{"Code": -1, "Error": "剩余数量为 0"})
		// 回滚事务
		tx.Rollback()
		return
	}

	// 构造更新中奖信息的内容
	giftUpdateInfo := structs.EventsGift{
		Phone:    c.Param("Phone"),
		SourceIP: c.ClientIP(),
	}
	// 从上下文中取出版本信息
	if apiVersion, exist := c.Get("Api-Version"); exist {
		// 存在版本信息, 添加到记录中
		giftUpdateInfo.ApiVersion = apiVersion.(string)
	} else {
		// 不存在版本信息, 将记录中的版本设置为 Unknown
		giftUpdateInfo.ApiVersion = "Unknown"
	}
	// 更新中奖信息
	tx.Model(structs.EventsGift{}).Where("id = ?", giftInfo.ID).Updates(&giftUpdateInfo)

	// 提交事务
	tx.Commit()

	// 返回信息
	c.JSON(http.StatusOK, gin.H{"Detail": giftInfo.Detail})
}

// GetSurplus 是获取礼品剩余数量接口的处理函数
func GetSurplus(c *gin.Context) {
	// 定义计数变量
	count := 0

	// 取出数量
	orm.PostgreSQL.Model(structs.EventsGift{}).Where("Name = ? AND Phone IS NULL", c.Param("Name")).Count(&count)

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"Surplus": count})
}
