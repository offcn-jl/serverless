/*
   @Time : 2020/5/22 5:08 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : push_test
   @Software: GoLand
*/

package sso

import (
	"fmt"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/gscf"
	"github.com/stretchr/testify/assert"
	"testing"
	"tsf/common/database/orm/structs"
)

func TestPostPush(t *testing.T) {
	round := 0
	orm.PostgreSQL.Model(&structs.SingleSignOnPushLog{}).Count(&round)

	// 初始化上下文
	c := gin.Context{}
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用

	// 初始化测试数据
	initTestData()

	// 测试 未绑定 Body 数据
	PostPush(&c)
	assert.Contains(t, c.Response.Body, "Invalid Json Data")
	// 增加 Body
	c.Request.Body = "{\"CRMSID\":\"6edbf791cfbaaa68442dd75bfd10ae5b\",\"Phone\":\"*\"}"

	// 测试 验证手机号码是否有效
	PostPush(&c)
	assert.Contains(t, c.Response.Body, "手机号码不正确")
	// 生成伪随机手机号码
	tempPhone1 := "17887101" + fmt.Sprintf("%03d", round) // 使用当前时间生成手机号尾号用于测试
	// 修正手机号码
	c.Request.Body = "{\"CRMSID\":\"6edbf791cfbaaa68442dd75bfd10ae5b\",\"Phone\":\"" + tempPhone1 + "\",\"CustomerName\":\"客户姓名\",\"CustomerIdentityID\":1,\"CustomerColleage\":\"客户毕业院校\",\"CustomerMayor\":\"客户专业\",\"Remark\":\"备注\"}"

	// 测试 后缀未填写, 使用默认后缀配置, 并且配置了自定义字段 fixme
	PostPush(&c)
	pushLog1 := structs.SingleSignOnPushLog{}
	orm.PostgreSQL.Where("crm_sid = '6edbf791cfbaaa68442dd75bfd10ae5b' AND phone = ?", tempPhone1).Find(&pushLog1)
	assert.Equal(t, uint(7), pushLog1.CRMChannel) // CRM 所属渠道
	// 获取默认后缀
	suffixInfo := structs.SingleSignOnSuffix{}
	orm.PostgreSQL.First(&suffixInfo)
	assert.Equal(t, "", pushLog1.ActualSuffix)                 // 校正前的后缀 空
	assert.Equal(t, suffixInfo.Suffix, pushLog1.CurrentSuffix) // 校正后的后缀 默认后缀
	assert.Equal(t, suffixInfo.CRMUID, pushLog1.CRMUID)        // CRM 所属用户 默认用户
	assert.Equal(t, uint(2290), pushLog1.CRMOCode)             // CRM 所属组织 按归属地分配 长春
	// 校验自定义字段
	assert.Equal(t, "客户姓名", pushLog1.CustomerName)        // CRM 所属组织 按归属地分配 长春
	assert.Equal(t, uint(1), pushLog1.CustomerIdentityID) // CRM 所属组织 按归属地分配 长春
	assert.Equal(t, "客户毕业院校", pushLog1.CustomerColleage)  // CRM 所属组织 按归属地分配 长春
	assert.Equal(t, "客户专业", pushLog1.CustomerMayor)       // CRM 所属组织 按归属地分配 长春
	assert.Equal(t, "备注", pushLog1.Remark)                // CRM 所属组织 按归属地分配 长春
	// 生成伪随机手机号码
	tempPhone2 := "17887102" + fmt.Sprintf("%03d", round) // 使用当前时间生成手机号尾号用于测试
	// 增加无效的后缀
	c.Request.Body = "{\"CRMSID\":\"6edbf791cfbaaa68442dd75bfd10ae5b\",\"Phone\":\"" + tempPhone2 + "\",\"Suffix\":\"wrong\"}"

	// 测试 后缀无效, 使用默认后缀配置
	PostPush(&c)
	pushLog2 := structs.SingleSignOnPushLog{}
	orm.PostgreSQL.Where("crm_sid = '6edbf791cfbaaa68442dd75bfd10ae5b' AND phone = ?", tempPhone2).Find(&pushLog2)
	assert.Equal(t, uint(7), pushLog2.CRMChannel)              // CRM 所属渠道
	assert.Equal(t, "wrong", pushLog2.ActualSuffix)            // 校正前的后缀 空
	assert.Equal(t, suffixInfo.Suffix, pushLog2.CurrentSuffix) // 校正后的后缀 默认后缀
	assert.Equal(t, suffixInfo.CRMUID, pushLog2.CRMUID)        // CRM 所属用户 默认用户
	assert.Equal(t, uint(2290), pushLog2.CRMOCode)             // CRM 所属组织 按归属地分配 长春
	// 生成伪随机手机号码
	tempPhone3 := "17887103" + fmt.Sprintf("%03d", round) // 使用当前时间生成手机号尾号用于测试
	// 增加不是省级的后缀
	c.Request.Body = "{\"CRMSID\":\"6edbf791cfbaaa68442dd75bfd10ae5b\",\"Phone\":\"" + tempPhone3 + "\",\"Suffix\":\"test\"}"

	// 测试 配置了 CRMOID 并且不是省级
	PostPush(&c)
	pushLog3 := structs.SingleSignOnPushLog{}
	orm.PostgreSQL.Where("crm_sid = '6edbf791cfbaaa68442dd75bfd10ae5b' AND phone = ?", tempPhone3).Find(&pushLog3)
	assert.Equal(t, uint(22), pushLog3.CRMChannel)  // CRM 所属渠道
	assert.Equal(t, "test", pushLog3.ActualSuffix)  // 校正前的后缀 test
	assert.Equal(t, "test", pushLog3.CurrentSuffix) // 校正后的后缀 test
	assert.Equal(t, uint(123), pushLog3.CRMUID)     // CRM 所属用户 高**
	assert.Equal(t, uint(2290), pushLog3.CRMOCode)  // CRM 所属组织 长春

	// 测试 检查是否已经进行过推送 ( 重复推送同一个 CRMSID + Phone )
	PostPush(&c)
	count := 0
	orm.PostgreSQL.Model(&structs.SingleSignOnPushLog{}).Where("crm_sid = '6edbf791cfbaaa68442dd75bfd10ae5b' AND phone = ?", tempPhone3).Count(&count)
	assert.Equal(t, 1, count)
}
