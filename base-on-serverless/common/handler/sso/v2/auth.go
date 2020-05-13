/*
   @Time : 2020/5/12 12:47 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : auth
   @Software: GoLand
*/

package sso

import (
	"fmt"
	"github.com/offcn-jl/go-common/codes"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/go-common/logger"
	"github.com/offcn-jl/go-common/verify"
	"github.com/offcn-jl/gscf"
	"github.com/xluohome/phonedata"
	"net/http"
	"net/url"
	"serverless/common/database/orm/structs"
	"time"
)

// PostSignUp 单点登模块注册接口的处理函数
func PostSignUp(c *gin.Context) {
	// 构造会话信息
	sessionInfo := structs.SingleSignOnSession{}
	// 绑定数据
	if err := c.ShouldBindJSON(&sessionInfo); err != nil {
		// 绑定数据错误
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"Code": codes.InvalidJson, "Error": codes.ErrorText(codes.InvalidJson)})
		return
	}
	sessionInfo.SourceIP = c.ClientIP()
	// 从上下文中取出版本信息
	if apiVersion, exist := c.Get("Api-Version"); exist {
		// 存在版本信息, 添加到记录中
		sessionInfo.ApiVersion = apiVersion.(string)
	} else {
		// 不存在版本信息, 将记录中的版本设置为 Unknown
		sessionInfo.ApiVersion = "Unknown"
	}

	// 验证手机号码是否有效
	if !verify.Phone(sessionInfo.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "手机号码不正确!"})
		return
	}

	// 校验登录模块配置
	moduleInfo := structs.SingleSignOnLoginModule{}
	orm.PostgreSQL.Where("id = ?", c.Param("MID")).Find(&moduleInfo)
	if moduleInfo.ID == 0 {
		// 模块不存在
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "单点登陆模块配置有误!"})
		return
	}
	// 保存模块 ID 到会话信息中
	sessionInfo.SSOMID = moduleInfo.ID     // 登陆模块 ID
	sessionInfo.CRMSID = moduleInfo.CRMSID // CRM 活动表单 ID

	// 检查验证码是否正确且未失效
	codeInfo := structs.SingleSignOnVerificationCode{}
	orm.PostgreSQL.Where("phone = ?", sessionInfo.Phone).Find(&codeInfo)
	// 校验是否发送过验证码
	if codeInfo.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "请您先获取验证码后再进行注册!"})
		return
	}
	// 校验验证码是正确
	if c.Param("Code") != fmt.Sprint(codeInfo.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "验证码有误!"})
		return
	}
	// 校验验证码是否有效
	duration, _ := time.ParseDuration("-" + fmt.Sprint(codeInfo.Term) + "m")
	if codeInfo.CreatedAt.Before(time.Now().Add(duration)) {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "验证码失效!"})
		return
	}

	// 校验用户是否已经注册 ( 避免重复注册 )
	if !isSignUp(sessionInfo.Phone) {
		// 保存注册信息
		orm.PostgreSQL.Create(&structs.SingleSignOnUser{Phone: sessionInfo.Phone})
	}

	// 校验用户是否已经参与过当前活动 ( 避免重复创建会话信息, 避免重复推送信息到 CRM )
	if !isSignIn(sessionInfo.Phone, sessionInfo.SSOMID) {
		// 未参与
		// 保存会话
		createSession(&sessionInfo)
	}

	// 注册成功
	c.JSON(http.StatusOK, gin.H{"Code": 0})
}

// PostSignIn 单点登陆模块登陆接口的处理函数
func PostSignIn(c *gin.Context) {
	// 构造会话信息
	sessionInfo := structs.SingleSignOnSession{}
	// 绑定数据
	if err := c.ShouldBindJSON(&sessionInfo); err != nil {
		// 绑定数据错误
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"Code": codes.InvalidJson, "Error": codes.ErrorText(codes.InvalidJson)})
		return
	}
	sessionInfo.SourceIP = c.ClientIP()
	// 从上下文中取出版本信息
	if apiVersion, exist := c.Get("Api-Version"); exist {
		// 存在版本信息, 添加到记录中
		sessionInfo.ApiVersion = apiVersion.(string)
	} else {
		// 不存在版本信息, 将记录中的版本设置为 Unknown
		sessionInfo.ApiVersion = "Unknown"
	}

	// 验证手机号码是否有效
	if !verify.Phone(sessionInfo.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "手机号码不正确!"})
		return
	}

	// 校验登录模块配置
	moduleInfo := structs.SingleSignOnLoginModule{}
	orm.PostgreSQL.Where("id = ?", c.Param("MID")).Find(&moduleInfo)
	if moduleInfo.ID == 0 {
		// 模块不存在
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "单点登陆模块配置有误!"})
		return
	}
	// 保存模块 ID 到会话信息中
	sessionInfo.SSOMID = moduleInfo.ID     // 登陆模块 ID
	sessionInfo.CRMSID = moduleInfo.CRMSID // CRM 活动表单 ID

	// 校验用户是否已经注册 ( 避免重复注册 )
	if !isSignUp(sessionInfo.Phone) {
		// 保存注册信息
		c.JSON(http.StatusForbidden, gin.H{"Code": -1, "Error": "请您先进行注册!"})
		return
	}

	// 校验用户是否已经参与过当前活动 ( 避免重复创建会话信息, 避免重复推送信息到 CRM )
	if !isSignIn(sessionInfo.Phone, sessionInfo.SSOMID) {
		// 未参与
		// 保存会话
		createSession(&sessionInfo)
	}

	// 登陆成功
	c.JSON(http.StatusOK, gin.H{"Code": 0})
}

// isSignUp 检查用户是否已经注册且未失效
func isSignUp(phone string) bool {
	userInfo := structs.SingleSignOnUser{}
	orm.PostgreSQL.Where("phone = ? and created_at > ?", phone, time.Now().AddDate(0, 0, -30)).Find(&userInfo)
	if userInfo.ID != 0 {
		return true
	}
	return false
}

// isSignIn 检查用户是否已经登陆
func isSignIn(phone string, mID uint) bool {
	sessionInfo := structs.SingleSignOnSession{}
	orm.PostgreSQL.Where("phone = ? and sso_mid = ?", phone, mID).Find(&sessionInfo)
	if sessionInfo.ID != 0 {
		return true
	}
	return false
}

// createSession 创建会话
// 创建会话前会进行推送信息到 CRM 的操作
func createSession(session *structs.SingleSignOnSession) {
	// 校验后缀
	if session.ActualSuffix == "0" {
		// 后缀为 0 ( 等同于未填写 ), 使用默认后缀配置
		getDefaultSuffix(session)
	} else {
		suffixInfo := structs.SingleSignOnSuffix{}
		orm.PostgreSQL.Where("suffix = ?", session.ActualSuffix).Find(&suffixInfo)
		if suffixInfo.ID == 0 {
			// 后缀无效, 使用默认后缀配置
			getDefaultSuffix(session)
		} else {
			session.CRMChannel = suffixInfo.Channel
			session.CRMUID = suffixInfo.CRMUID
			session.CurrentSuffix = suffixInfo.Suffix
			if suffixInfo.CRMOID > 1 {
				// 配置了 CRMOID 并且不是省级
				organizationInfo := structs.SingleSignOnOrganization{}
				orm.PostgreSQL.Where("id = ?", suffixInfo.CRMOID)
				session.CRMOCode = organizationInfo.Code
			} else {
				// 未配置 CRMOID 或者是省级 ( 等于 1 ), 按手机号码归属地分配 CRM 信息
				distributionByPhoneNumber(session)
			}
		}
	}

	// 推送信息到 CRM
	if urlObject, err := url.Parse("https://dc.offcn.com:8443/a.gif"); err != nil {
		// 请求失败
		logger.Error(err)
		// 推送失败, 保存推送失败记录
		orm.PostgreSQL.Create(&structs.SingleSignOnErrorLog{
			Phone:      session.Phone,
			SSOMID:     session.SSOMID,
			CRMChannel: session.CRMChannel,
			CRMUID:     session.CRMUID,
			CRMOCode:   session.CRMOCode,
			Error:      "Parse GET URL Fail : " + err.Error(),
		})
	} else {
		// 构建参数 queryObject
		queryObject := urlObject.Query()
		queryObject.Set("sid", session.CRMSID)
		queryObject.Set("mobile", session.Phone)
		queryObject.Set("channel", fmt.Sprint(session.CRMChannel))
		queryObject.Set("orgn", fmt.Sprint(session.CRMOCode))
		if session.CRMUID != 0 {
			queryObject.Set("owner", fmt.Sprint(session.CRMUID))
		}
		if session.CustomerName != "" {
			queryObject.Set("name", session.CustomerName)
		}
		if session.CustomerIdentityID != 0 {
			queryObject.Set("khsf", fmt.Sprint(session.CustomerIdentityID))
		}
		if session.CustomerColleage != "" {
			queryObject.Set("colleage", session.CustomerColleage)
		}
		if session.CustomerMayor != "" {
			queryObject.Set("mayor", session.CustomerMayor)
		}
		if session.Remark != "" {
			queryObject.Set("remark", session.Remark)
		}
		// 发送 GET 请求
		urlObject.RawQuery = queryObject.Encode()
		if getResponse, err := http.Get(urlObject.String()); err != nil {
			// 发送 GET 请求出错
			logger.Error(err)
			// 推送失败, 保存推送失败记录
			orm.PostgreSQL.Create(&structs.SingleSignOnErrorLog{
				Phone:      session.Phone,
				SSOMID:     session.SSOMID,
				CRMChannel: session.CRMChannel,
				CRMUID:     session.CRMUID,
				CRMOCode:   session.CRMOCode,
				Error:      "GET Fail : " + err.Error(),
			})
		} else {
			if getResponse.StatusCode != 200 {
				// 请求出错
				logger.Error(err)
				// 推送失败, 保存推送失败记录
				orm.PostgreSQL.Create(&structs.SingleSignOnErrorLog{
					Phone:      session.Phone,
					SSOMID:     session.SSOMID,
					CRMChannel: session.CRMChannel,
					CRMUID:     session.CRMUID,
					CRMOCode:   session.CRMOCode,
					Error:      "Json POST Response Fail : " + err.Error(),
				})
			}
		}
	}

	// 保存会话信息
	orm.PostgreSQL.Create(&session)
}

// getDefaultSuffix 获取默认后缀配置
func getDefaultSuffix(session *structs.SingleSignOnSession) {
	session.CRMChannel = 7 // 默认所属渠道 19课堂
	// 获取默认后缀
	suffixInfo := structs.SingleSignOnSuffix{}
	orm.PostgreSQL.First(&suffixInfo)
	// 配置默认后缀
	session.CurrentSuffix = suffixInfo.Suffix
	session.CRMUID = suffixInfo.CRMUID
	// 所属组织按照手机号码归属地进行分配
	distributionByPhoneNumber(session)
}

// distributionByPhoneNumber 按照手机号码归属地进行归属分部分配
func distributionByPhoneNumber(session *structs.SingleSignOnSession) {
	if record, err := phonedata.Find(session.Phone); err != nil {
		logger.Error(err)
		// 解析出错，循环分配给九个地市
		roundCrmList(session)
	} else {
		switch record.City {
		case "长春":
			session.CRMOCode = 2290
		case "吉林":
			session.CRMOCode = 2305
		case "延边":
			session.CRMOCode = 2277
		case "通化":
			session.CRMOCode = 2271
		case "白山":
			session.CRMOCode = 2310
		case "四平":
			session.CRMOCode = 2263
		case "松原":
			session.CRMOCode = 2284
		case "白城":
			session.CRMOCode = 2315
		case "辽源":
			session.CRMOCode = 2268
		default:
			// 循环分配给九个地市
			roundCrmList(session)
		}
	}
}

// roundCrmList 循环分配手机号给九个地市分部
// 高并发时存在数据库读写延迟，可能无法确保幂等性
// 即, 可能出现同时某分部重复分配, 或跳过某分部进行分配的情况
func roundCrmList(session *structs.SingleSignOnSession) {
	// 取出地市分校列表
	crmOrganizations := make([]structs.SingleSignOnOrganization, 0)
	orm.PostgreSQL.Where("f_id = 1").Find(&crmOrganizations)
	logger.DebugToJson("crmOrganizations", crmOrganizations)

	// 获取当前分配计数
	count := 0
	orm.PostgreSQL.Model(structs.SingleSignOnCRMRoundLog{}).Count(&count)
	logger.DebugToString("count", count)

	// 分配
	session.CRMOCode = crmOrganizations[count%len(crmOrganizations)].Code

	// 保存分配记录
	orm.PostgreSQL.Create(&structs.SingleSignOnCRMRoundLog{
		SSOMID: session.SSOMID,
		Phone:  session.Phone,
		CRMOID: crmOrganizations[count%len(crmOrganizations)].ID,
	})
}
