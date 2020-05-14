/*
   @Time : 2020/5/14 2:21 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : session_test
   @Software: GoLand
*/

package sso

import (
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/gscf"
	"github.com/stretchr/testify/assert"
	"serverless/common/database/orm/structs"
	"testing"
	"time"
)

func TestGetSessionInfo(t *testing.T) {
	// 初始化上下文
	c := gin.Context{}
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用

	// 测试 验证手机号是否有效 ( 需要是 0 或正确的手机号 )
	// 此时没有填写手机号
	GetSessionInfo(&c)
	assert.Contains(t, c.Response.Body, "手机号码不正确")
	// 此时填写了错误的手机号
	c.Request.PathParameters = map[string]string{"Phone": "1788710666"} // 10 位数字, 此时隐式的初始化了 c.Request.PathParameters 这个 map
	GetSessionInfo(&c)
	assert.Contains(t, c.Response.Body, "手机号码不正确")
	// 修正为正确的手机号码
	c.Request.PathParameters["Phone"] = "17887106666"

	// 测试 校验登陆模块配置
	// 此时没有填写登陆模块配置
	GetSessionInfo(&c)
	assert.Contains(t, c.Response.Body, "单点登陆模块配置有误")
	// 此时填写了错误的登陆模块 ID
	c.Request.PathParameters["MID"] = "1"
	GetSessionInfo(&c)
	assert.Contains(t, c.Response.Body, "单点登陆模块配置有误")
	// 修正为正确的登陆模块 ID
	c.Request.PathParameters["MID"] = "10001"

	// 测试 后缀不存在或错误
	defaultInfo := "{\"CRMEID\":\"HD202003061144\",\"CRMSID\":\"6edbf791cfbaaa68442dd75bfd10ae5b\",\"CRMChannel\":7,\"CRMOCode\":22,\"CRMOName\":\"吉林分校\",\"CRMUID\":32431,\"CRMUser\":\"default\",\"Suffix\":\"default\",\"IsLogin\":false,\"NeedToRegister\":true}"
	// 此时未配置后缀, 即后缀不存在, 可以认为等同为后缀错误
	GetSessionInfo(&c)
	assert.Equal(t, defaultInfo, c.Response.Body)
	// 设置设置默认后缀 CRM 组织 ID 为 0
	orm.PostgreSQL.Model(structs.SingleSignOnSuffix{}).Where("suffix = 'default'").Update("crm_oid", "0")
	GetSessionInfo(&c)
	assert.Equal(t, defaultInfo, c.Response.Body)
	// 设置默认后缀 CRM 组织 ID 为不存在的 ID
	orm.PostgreSQL.Model(structs.SingleSignOnSuffix{}).Where("suffix = 'default'").Update("crm_oid", "10000")
	GetSessionInfo(&c)
	assert.Equal(t, defaultInfo, c.Response.Body)
	// 还原默认后缀 CRM 组织 ID
	orm.PostgreSQL.Model(structs.SingleSignOnSuffix{}).Where("suffix = 'default'").Update("crm_oid", "1")

	// 测试 配置了后缀
	testInfo := "{\"CRMEID\":\"HD202003061144\",\"CRMSID\":\"6edbf791cfbaaa68442dd75bfd10ae5b\",\"CRMChannel\":104,\"CRMOCode\":2290,\"CRMOName\":\"吉林长春分校\",\"CRMUID\":123,\"CRMUser\":\"test\",\"Suffix\":\"test\",\"IsLogin\":false,\"NeedToRegister\":true}"
	testInfoWithDefauleOrgnation := "{\"CRMEID\":\"HD202003061144\",\"CRMSID\":\"6edbf791cfbaaa68442dd75bfd10ae5b\",\"CRMChannel\":104,\"CRMOCode\":22,\"CRMOName\":\"吉林分校\",\"CRMUID\":123,\"CRMUser\":\"test\",\"Suffix\":\"test\",\"IsLogin\":false,\"NeedToRegister\":true}"
	// 配置后缀为测试后缀
	c.Request.PathParameters["Suffix"] = "test"
	GetSessionInfo(&c)
	assert.Equal(t, testInfo, c.Response.Body)
	// 设置设置默认后缀 CRM 组织 ID 为 0
	orm.PostgreSQL.Model(structs.SingleSignOnSuffix{}).Where("suffix = 'test'").Update("crm_oid", "0")
	GetSessionInfo(&c)
	assert.Equal(t, testInfoWithDefauleOrgnation, c.Response.Body)
	// 设置默认后缀 CRM 组织 ID 为不存在的 ID
	orm.PostgreSQL.Model(structs.SingleSignOnSuffix{}).Where("suffix = 'test'").Update("crm_oid", "10000")
	GetSessionInfo(&c)
	assert.Equal(t, testInfoWithDefauleOrgnation, c.Response.Body)
	// 还原默认后缀 CRM 组织 ID
	orm.PostgreSQL.Model(structs.SingleSignOnSuffix{}).Where("suffix = 'test'").Update("crm_oid", "2")

	// 测试 校验是否需要注册
	// 模拟注册手机号
	userInfo := structs.SingleSignOnUser{}
	userInfo.Phone = "1788710" + time.Now().Format("0405") // 使用当前时间生成手机号尾号用于测试
	orm.PostgreSQL.Create(&userInfo)
	// 更换手机号为已经注册的手机号
	c.Request.PathParameters["Phone"] = userInfo.Phone
	GetSessionInfo(&c)
	assert.Contains(t, c.Response.Body, "\"NeedToRegister\":false")

	// 测试 校验是否需要登陆
	// 模拟创建会话
	session := structs.SingleSignOnSession{}
	session.MID = 10001
	session.Phone = userInfo.Phone
	createSession(&session)
	GetSessionInfo(&c)
	assert.Contains(t, c.Response.Body, "\"IsLogin\":true")
}
