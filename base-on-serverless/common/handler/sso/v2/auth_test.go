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
	"serverless/common/database/orm/structs"
	"testing"
	"time"
)

// 测试 单点登陆模块登陆接口的处理函数
func TestPostSignIn(t *testing.T) {
	c := gin.Context{}
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用

	// 测试未绑定 Body 数据
	PostSignIn(&c)
	assert.Contains(t, c.Response.Body, "Invalid Json Data")
	// 增加 Body
	c.Request.Body = "{\"Phone\":\"*\"}"

	// 测试 验证手机号码是否有效
	PostSignIn(&c)
	assert.Contains(t, c.Response.Body, "手机号码不正确")
	// 修正手机号码
	c.Request.Body = "{\"Phone\":\"17866886688\"}"

	// 测试 校验登录模块配置 ( 此时的上下文中没有登陆模块的配置 )
	PostSignIn(&c)
	assert.Contains(t, c.Response.Body, "单点登陆模块配置有误")
	// 向上下文中添加登陆模块配置, 这一步操作隐式的初始化了 PathParameters
	c.Request.PathParameters = map[string]string{"MID": "1"}

	// 测试 登陆成功
	PostSignIn(&c)
	assert.Equal(t, "{\"Code\":0}", c.Response.Body)

	// 判断会话是否存在
	assert.True(t, isSignIn("17866886688", 1))
}

// 测试 单点登模块注册接口的处理函数
func TestPostSignUp(t *testing.T) {
	c := gin.Context{}
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用

	// 测试未绑定 Body 数据
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "Invalid Json Data")
	// 增加 Body
	c.Request.Body = "{\"Phone\":\"*\"}"

	// 测试 验证手机号码是否有效
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "手机号码不正确")
	// 修正手机号码
	c.Request.Body = "{\"Phone\":\"17866668888\"}"

	// 测试 校验登录模块配置 ( 此时的上下文中没有登陆模块的配置 )
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "单点登陆模块配置有误")
	// 向上下文中添加登陆模块配置, 这一步操作隐式的初始化了 PathParameters
	c.Request.PathParameters = map[string]string{"MID": "1"}

	// 测试 校验是否发送过验证码
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "请您先获取验证码后再进行注册")
	// 修正为获取过验证码的手机号码
	c.Request.Body = "{\"Phone\":\"17888886666\"}"

	// 测试 校验验证码是正确 ( 此时的上下文中未填写验证码 )
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "验证码有误")
	// 向上下文中添加正确, 但已经失效的验证码
	c.Request.PathParameters["Code"] = "0"

	// 测试 校验验证码是否有效
	PostSignUp(&c)
	assert.Contains(t, c.Response.Body, "验证码失效")
	// 更换上下文中但手机号码及验证码未为正确且未失效的内容
	c.Request.Body = "{\"Phone\":\"17866886688\"}"
	c.Request.PathParameters["Code"] = "1"

	// 测试 注册成功
	PostSignUp(&c)
	assert.Equal(t, "{\"Code\":0}", c.Response.Body)

	// 判断用户是否存在
	assert.True(t, isSignUp("17866886688"))

	// 判断会话是否存在
	assert.True(t, isSignIn("17866886688", 1))
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
	session.SSOMID = 1
	session.Phone = "1788710" + time.Now().Format("0405") // 使用当前时间生成手机号尾号用于测试

	assert.False(t, isSignIn(session.Phone, session.SSOMID))

	createSession(&session)

	assert.True(t, isSignIn(session.Phone, session.SSOMID))
}

// 测试 创建会话
func TestCreateSession(t *testing.T) {
	session := structs.SingleSignOnSession{}
	session.ActualSuffix = "0"                            // 默认后缀
	session.Phone = "1788710" + time.Now().Format("0405") // 使用当前时间生成手机号尾号用于测试
	session.CRMSID = "6edbf791cfbaaa68442dd75bfd10ae5b"   // 测试 CRM 活动表单

	session.CustomerName = "测试姓名"
	session.CustomerIdentityID = 1 // 在校生-大一
	session.CustomerColleage = "测试学校"
	session.CustomerMayor = "测试专业"
	session.Remark = "测试备注"

	// 校验校正前的数据
	assert.Equal(t, session.ActualSuffix, "0")   // 后缀
	assert.Equal(t, session.CurrentSuffix, "")   // 校正后的后缀
	assert.Equal(t, session.CRMChannel, uint(0)) // CRM 所属渠道
	assert.Equal(t, session.CRMUID, uint(0))     // CRM 用户 ID
	assert.Equal(t, session.CRMOCode, uint(0))   // CRM 组织代码

	createSession(&session)

	// 校验校正后的数据
	assert.Equal(t, session.ActualSuffix, "0")     // 后缀
	assert.Equal(t, session.CurrentSuffix, "test") // 校正后的后缀
	assert.Equal(t, session.CRMChannel, uint(7))   // CRM 所属渠道
	assert.Equal(t, session.CRMUID, uint(123))     // CRM 用户 ID
	assert.Equal(t, session.CRMOCode, uint(2290))  // CRM 组织 ID
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
