/*
   @Time : 2020/5/13 10:14 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : test_auth
   @Software: GoLand
*/

package sso

import (
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/gscf"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"tsf/common/database/orm/structs"
)

// 测试 单点登模块注册接口的处理函数
func TestPostSignUp(t *testing.T) {
	// 初始化上下文
	c := gin.Context{}
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用

	// 初始化测试数据
	initTestData()

	// 测试 未绑定 Body 数据
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "Invalid Json Data")
	// 增加 Body
	c.Request.Body = "{\"MID\":10000,\"Phone\":\"*\"}"

	// 测试 验证手机号码是否有效
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "手机号码不正确")
	// 修正手机号码
	c.Request.Body = "{\"MID\":10000,\"Phone\":\"17866668888\"}"

	// 测试 校验登录模块配置
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "单点登陆模块配置有误")
	// 修正为已存在的登陆模块 ID
	c.Request.Body = "{\"MID\":10001,\"Phone\":\"17866668888\"}"

	// 测试 校验是否发送过验证码
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "请您先获取验证码后再进行注册")
	// 模拟获取验证码
	createTestVerificationCode()
	// 修正为模拟获取过验证码的手机号码
	c.Request.Body = "{\"MID\":10001,\"Phone\":\"17888886666\"}"

	// 测试 校验验证码是正确 ( 此时的上下文中未填写验证码 )
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "验证码有误")
	// 向上请求中添加正确, 但已经失效的验证码
	c.Request.Body = "{\"MID\":10001,\"Phone\":\"17888886666\",\"Code\":9999}"

	// 测试 校验验证码是否有效
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "验证码失效")
	// 更换上下文中的手机号码及验证码未为正确且未失效的内容
	c.Request.Body = "{\"MID\":10001,\"Phone\":\"17866886688\",\"Code\":9999}"

	// 测试 注册成功
	PostSignUp(&c)
	assert.Equal(t, "{\"Code\":0}", c.Response.Body)

	// 判断用户是否存在
	assert.True(t, isSignUp("17866886688"))

	// 判断会话是否存在
	assert.True(t, isSignIn("17866886688", 10001))
}

// 测试 单点登陆模块登陆接口的处理函数
func TestPostSignIn(t *testing.T) {
	// 初始化上下文
	c := gin.Context{}
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用

	// 初始化测试数据
	initTestData()

	// 测试未绑定 Body 数据
	PostSignIn(&c)
	assert.Contains(t, c.Response.Body, "Invalid Json Data")
	// 增加 Body
	c.Request.Body = "{\"MID\":10000,\"Phone\":\"*\"}"

	// 测试 验证手机号码是否有效
	PostSignIn(&c)
	assert.Contains(t, c.Response.Body, "手机号码不正确")
	// 修正手机号码
	c.Request.Body = "{\"MID\":10000,\"Phone\":\"17866886688\"}"

	// 测试 校验登录模块配置 ( 此时的上下文中没有登陆模块的配置 )
	PostSignIn(&c)
	assert.Contains(t, c.Response.Body, "单点登陆模块配置有误")
	// 修正登陆模块 ID
	c.Request.Body = "{\"MID\":10001,\"Phone\":\"17888668866\"}"

	// 测试 校验用户是否已经注册
	PostSignIn(&c)
	assert.Contains(t, c.Response.Body, "请您先进行注册")
	// 修正为已经注册的手机号码
	c.Request.Body = "{\"MID\":10001,\"Phone\":\"17866886688\"}"

	// 测试 登陆成功
	PostSignIn(&c)
	assert.Equal(t, "{\"Code\":0}", c.Response.Body)

	// 判断会话是否存在
	assert.True(t, isSignIn("17866886688", 10001))
}

// 测试 检查用户是否已经注册且未失效
func TestIsSignUp(t *testing.T) {
	userInfo := structs.SingleSignOnUser{}
	userInfo.Phone = "1788710" + time.Now().Format("0405") // 使用当前时间生成手机号尾号用于测试

	assert.False(t, isSignUp(userInfo.Phone))

	orm.PostgreSQL.Create(&userInfo)

	assert.True(t, isSignUp(userInfo.Phone))
}

// 测试 检查用户是否已经登陆
func TestIsSignIn(t *testing.T) {
	session := structs.SingleSignOnSession{}
	session.MID = 1
	session.Phone = "1788710" + time.Now().Format("0405") // 使用当前时间生成手机号尾号用于测试

	assert.False(t, isSignIn(session.Phone, session.MID))

	createSession(&session)

	assert.True(t, isSignIn(session.Phone, session.MID))
}

// 测试 创建会话
func TestCreateSession(t *testing.T) {
	session := structs.SingleSignOnSession{}
	session.Phone = "1788710" + time.Now().Format("0405") // 使用当前时间生成手机号尾号用于测试
	session.CRMSID = "6edbf791cfbaaa68442dd75bfd10ae5b"   // 测试 CRM 活动表单

	session.CustomerName = "测试姓名"
	session.CustomerIdentityID = 1 // 在校生-大一
	session.CustomerColleage = "测试学校"
	session.CustomerMayor = "测试专业"
	session.Remark = "测试备注"

	// 校验校正前的数据
	assert.Equal(t, "", session.ActualSuffix)    // 后缀
	assert.Equal(t, "", session.CurrentSuffix)   // 校正后的后缀
	assert.Equal(t, uint(0), session.CRMChannel) // CRM 所属渠道
	assert.Equal(t, uint(0), session.CRMUID)     // CRM 用户 ID
	assert.Equal(t, uint(0), session.CRMOCode)   // CRM 组织代码

	createSession(&session)

	// 校验校正后的数据
	assert.Equal(t, "", session.ActualSuffix)         // 后缀
	assert.Equal(t, "default", session.CurrentSuffix) // 校正后的后缀
	assert.Equal(t, uint(7), session.CRMChannel)      // CRM 所属渠道
	assert.Equal(t, uint(32431), session.CRMUID)      // CRM 用户 ID
	assert.Equal(t, uint(2290), session.CRMOCode)     // CRM 组织 ID

	// 验证测试后缀
	session.ID = 0 // 将 ID 恢复为 0, 令 ORM 认为这条 Session 是新记录
	session.ActualSuffix = "test"
	session.Phone = "1868648" + time.Now().Format("0405") // 使用当前时间生成手机号尾号用于测试
	createSession(&session)

	// 校验测试后缀信息
	assert.Equal(t, "test", session.ActualSuffix)  // 后缀
	assert.Equal(t, "test", session.CurrentSuffix) // 校正后的后缀
	assert.Equal(t, uint(22), session.CRMChannel)  // CRM 所属渠道
	assert.Equal(t, uint(123), session.CRMUID)     // CRM 用户 ID
	assert.Equal(t, uint(2290), session.CRMOCode)  // CRM 组织 ID
}

// 测试 获取默认后缀配置
func TestGetDefaultSuffix(t *testing.T) {
	session := structs.SingleSignOnSession{}

	session.Phone = "17887106666"

	assert.Equal(t, session.CRMChannel, uint(0))
	assert.Equal(t, session.CRMOCode, uint(0))

	getDefaultSuffix(&session)

	assert.Equal(t, session.CRMOCode, uint(2290))
	assert.Equal(t, session.CRMChannel, uint(7))
}

// 测试 按照手机号码归属地进行归属分部分配
// 号段数据来自 http://www.bixinshui.com
func TestDistributionByPhoneNumber(t *testing.T) {
	session := structs.SingleSignOnSession{}

	// 长春
	session.Phone = "17887106666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2290))

	// 吉林
	session.Phone = "13009156666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2305))

	// 延边
	session.Phone = "18943306666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2277))

	// 通化
	session.Phone = "13009196666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2271))

	// 白山
	session.Phone = "13009076666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2310))

	// 四平
	session.Phone = "13009026666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2263))

	// 松原
	session.Phone = "13009056666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2284))

	// 白城
	session.Phone = "13009066666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2315))

	// 辽源
	session.Phone = "13009046666"
	session.CRMOCode = 0
	assert.Equal(t, session.CRMOCode, uint(0))
	distributionByPhoneNumber(&session)
	assert.Equal(t, session.CRMOCode, uint(2268))
}

// 测试 循环分配手机号给九个地市分部
func TestRoundCrmList(t *testing.T) {
	session := structs.SingleSignOnSession{}

	assert.Equal(t, session.CRMOCode, uint(0))
	roundCrmList(&session)
	assert.NotEqual(t, session.CRMOCode, uint(0))
}

// createTestVerificationCode 模拟获取验证码
func createTestVerificationCode() {
	// 验证码正确, 但是已经失效
	orm.PostgreSQL.Create(&structs.SingleSignOnVerificationCode{
		Phone: "17888886666",
		Term:  0,
		Code:  9999,
	})
	// 验证码正确, 并且有效
	orm.PostgreSQL.Create(&structs.SingleSignOnVerificationCode{
		Phone: "17866886688",
		Term:  5,
		Code:  9999,
	})
}

// initTestData 初始化测试数据
func initTestData() {
	// 创建 测试用组织信息
	// 省级分校
	rootOriginationInfo := structs.SingleSignOnOrganization{}
	rootOriginationInfo.ID = 1
	rootOriginationInfo.FID = 0
	rootOriginationInfo.Code = 22
	rootOriginationInfo.Name = "吉林分校"
	orm.PostgreSQL.Create(&rootOriginationInfo)
	// 地市分校 1
	originationInfo1 := structs.SingleSignOnOrganization{}
	originationInfo1.ID = 2
	originationInfo1.FID = rootOriginationInfo.ID
	originationInfo1.Code = 2290
	originationInfo1.Name = "吉林长春分校"
	orm.PostgreSQL.Create(&originationInfo1)
	// 地市分校 2
	originationInfo2 := structs.SingleSignOnOrganization{}
	originationInfo2.ID = 3
	originationInfo2.FID = rootOriginationInfo.ID
	originationInfo2.Code = 2305
	originationInfo2.Name = "吉林市分校"
	orm.PostgreSQL.Create(&originationInfo2)

	// 创建 测试用后缀信息
	// 默认后缀 ( ID = 1 )
	defaultSuffixInfo := structs.SingleSignOnSuffix{}
	defaultSuffixInfo.ID = 1
	defaultSuffixInfo.Suffix = "default"
	defaultSuffixInfo.Name = "默认后缀(ID=1)"
	defaultSuffixInfo.CRMUser = "default"
	defaultSuffixInfo.CRMUID = 32431 // 齐*
	defaultSuffixInfo.CRMOID = 1     // 吉林分校
	defaultSuffixInfo.CRMChannel = 7 // 19 课堂 ( 网推 )
	orm.PostgreSQL.Create(&defaultSuffixInfo)
	// 后缀 1
	suffixInfo := structs.SingleSignOnSuffix{}
	suffixInfo.ID = 2
	tempTime := time.Now().Add(8760 * time.Hour) // 一年后
	suffixInfo.DeletedAt = &tempTime
	suffixInfo.Suffix = "test"
	suffixInfo.Name = "后缀 1"
	suffixInfo.CRMUser = "test"
	suffixInfo.CRMUID = 123    // 高**
	suffixInfo.CRMOID = 2      // 吉林长春分校
	suffixInfo.CRMChannel = 22 // 户外推广 ( 市场 )
	orm.PostgreSQL.Create(&suffixInfo)

	// 创建 测试用登陆模块信息
	testModelInfo := structs.SingleSignOnLoginModule{}
	testModelInfo.ID = 10001
	testModelInfo.CreatedAt = time.Now()
	testModelInfo.UpdatedAt = testModelInfo.CreatedAt
	testModelInfo.Name = "测试活动"
	testModelInfo.CRMEID = "HD202003061144"
	testModelInfo.CRMSID = "6edbf791cfbaaa68442dd75bfd10ae5b"
	testModelInfo.Term = 9999
	testModelInfo.Platform = 1
	testModelInfo.Sign = "中公教育"
	testModelInfo.TemplateID = 392074
	orm.PostgreSQL.Create(&testModelInfo)
}
