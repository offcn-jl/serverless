/*
   @Time : 2020/5/17 1:13 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : h20200502
   @Software: GoLand
*/

// 可用于礼品 ( 兑换码 ) 发放类活动
// 礼品需要提前添加到数据库中
// 可以配合校验接口和消费接口，对奖品进行校验和消费操作
package h20200502

import (
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common/database/orm"
	"net/http"
	"tsf/common/database/orm/structs"
)

// GetSurplus 获取礼品剩余数量接口的处理函数
func GetSurplus(c *gin.Context) {
	// 定义计数变量
	count := 0

	// 取出数量
	orm.PostgreSQL.Model(structs.E20200502{}).Where("Name = ? AND Phone IS NULL", c.Param("Name")).Count(&count)

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"Surplus": count})
}

// GetGift 获取礼品接口的处理函数
// 如果没有领取过奖品, 将会从奖品表中取出一个奖品然后返回奖品信息
// 如果领取过奖品, 将会直接返回第一次领取到的奖品信息
// 注意 : 本接口使用了事务, 对将要领取奖品所在行进行加锁处理, 在大量未领奖用户同时进行领奖时会因这个加锁逻辑出现数据库连接数被占满而停止响应的情况。所以理论上本接口在大量新用户领奖时, 最大并发数即数据库最大连接数。
func GetGift(c *gin.Context) {
	// 查询是否已经领取过奖品
	giftCheckInfo := structs.E20200502{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	// 已经领取过奖品, 直接返回奖品信息
	if giftCheckInfo.ID != 0 {
		c.JSON(http.StatusOK, gin.H{"Detail": giftCheckInfo.Detail})
		return
	}

	// 开启事务
	tx := orm.PostgreSQL.Begin()

	// 取出一个奖品，并锁定其所在行
	giftInfo := structs.E20200502{}
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
	giftUpdateInfo := structs.E20200502{
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
	tx.Model(structs.E20200502{}).Where("id = ?", giftInfo.ID).Updates(&giftUpdateInfo)

	// 提交事务
	tx.Commit()

	// 返回信息
	c.JSON(http.StatusOK, gin.H{"Detail": giftInfo.Detail})
}

// GetCheckout 查询领取信息接口的处理函数
func GetCheckout(c *gin.Context) {
	giftCheckInfo := structs.E20200502{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	c.JSON(http.StatusOK, giftCheckInfo)
}

// PatchConsume 消费礼品接口的处理函数
// 礼品未被消费时, 将奖品消费信息保存后, 返回 HTTP 状态码 200
// 礼品已经被消费时, 返回奖品消费详情, HTTP 状态码为 403
func PatchConsume(c *gin.Context) {
	// 获取奖品详情
	giftCheckInfo := structs.E20200502{}
	orm.PostgreSQL.Where("Name = ? AND Phone = ?", c.Param("Name"), c.Param("Phone")).First(&giftCheckInfo)

	// 判断奖品是否已经被消费
	if giftCheckInfo.ConsumeDetail != "" {
		// 奖品已经被消费
		c.JSON(http.StatusForbidden, gin.H{"Code": -1, "Error": "礼品已被消费，详情 : " + giftCheckInfo.ConsumeDetail})
	} else {
		// 奖品未被消费, 更新奖品信息
		orm.PostgreSQL.Model(&structs.E20200502{}).Where("name = ? AND phone = ?", c.Param("Name"), c.Param("Phone")).Update("consume_detail", c.Param("ConsumeDetail"))
		// 隐式返回 HTTP 状态码 200
	}
}
