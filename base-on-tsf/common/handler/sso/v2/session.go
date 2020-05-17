/*
   @Time : 2020/5/14 1:37 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : session
   @Software: GoLand
*/

package sso

import (
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/go-common/verify"
	"net/http"
	"tsf/common/database/orm/structs"
)

// GetSessionInfo 获取会话信息
func GetSessionInfo(c *gin.Context) {
	response := struct {
		Sign           string // 发信签名
		CRMEID         string // CRM 活动 ID
		CRMSID         string // CRM 活动表单 ID
		CRMChannel     uint   // CRM 所属渠道
		CRMOCode       uint   // CRM 组织代码
		CRMOName       string // CRM 组织名称
		CRMUID         uint   // CRM 用户 ID
		CRMUser        string // CRM 用户名
		Suffix         string // 后缀 ( 19课堂后缀 )
		IsLogin        bool   // 是否登陆
		NeedToRegister bool   // 是否需要注册
	}{}

	// 验证手机号是否有效
	if c.Param("Phone") != "0" && !verify.Phone(c.Param("Phone")) {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "手机号码不正确!"})
		return
	}

	// 校验登陆模块配置
	moduleInfo := structs.SingleSignOnLoginModule{}
	orm.PostgreSQL.Where("id = ?", c.Param("MID")).Find(&moduleInfo)
	if moduleInfo.ID == 0 {
		// 模块不存在
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "单点登陆模块配置有误!"})
		return
	}
	// 保存模块信息到会话信息中
	response.Sign = moduleInfo.Sign     // 发信签名
	response.CRMEID = moduleInfo.CRMEID // CRM 活动 ID
	response.CRMSID = moduleInfo.CRMSID // CRM 活动表单 ID

	// 校验后缀
	suffixInfo := structs.SingleSignOnSuffix{}
	orm.PostgreSQL.Where("suffix = ?", c.Param("Suffix")).Find(&suffixInfo)
	if suffixInfo.ID == 0 {
		// 后缀不存在
		// 获取默认后缀 ( ID = 1, 第一条 )
		defaultSuffixInfo := structs.SingleSignOnSuffix{}
		orm.PostgreSQL.First(&defaultSuffixInfo)
		response.Suffix = defaultSuffixInfo.Suffix         // 后缀 ( 19课堂后缀 )
		response.CRMChannel = defaultSuffixInfo.CRMChannel // CRM 所属渠道
		response.CRMUID = defaultSuffixInfo.CRMUID         // CRM 用户 ID
		response.CRMUser = defaultSuffixInfo.CRMUser       // CRM 用户名
		// 获取默认后缀对应的 CRM 组织信息
		organizationInfo := structs.SingleSignOnOrganization{}
		if defaultSuffixInfo.CRMOID == 0 {
			// 当前后缀未配置归属组织, 获取省级分部的信息
			orm.PostgreSQL.Where("f_id = 0").Find(&organizationInfo)
		} else {
			// 获取当前后缀配置的归属分部信息
			orm.PostgreSQL.Where("id = ?", defaultSuffixInfo.CRMOID).Find(&organizationInfo)
			if organizationInfo.ID == 0 {
				// 获取当前后缀配置的归属分部信息失败, 获取省级分部信息
				orm.PostgreSQL.Where("f_id = 0").Find(&organizationInfo)
			}
		}
		response.CRMOCode = organizationInfo.Code
		response.CRMOName = organizationInfo.Name
	} else {
		// 后缀存在
		response.Suffix = suffixInfo.Suffix         // 后缀 ( 19课堂后缀 )
		response.CRMChannel = suffixInfo.CRMChannel // CRM 所属渠道
		response.CRMUID = suffixInfo.CRMUID         // CRM 用户 ID
		response.CRMUser = suffixInfo.CRMUser       // CRM 用户名
		// 获取 CRM 组织信息
		organizationInfo := structs.SingleSignOnOrganization{}
		if suffixInfo.CRMOID == 0 {
			// 当前后缀未配置归属组织, 获取省级分部的信息
			orm.PostgreSQL.Where("f_id = 0").Find(&organizationInfo)
		} else {
			// 获取当前后缀配置的归属分部信息
			orm.PostgreSQL.Where("id = ?", suffixInfo.CRMOID).Find(&organizationInfo)
			if organizationInfo.ID == 0 {
				// 获取当前后缀配置的归属分部信息失败, 获取省级分部信息
				orm.PostgreSQL.Where("f_id = 0").Find(&organizationInfo)
			}
		}
		response.CRMOCode = organizationInfo.Code
		response.CRMOName = organizationInfo.Name
	}

	// 校验是否需要注册
	if !isSignUp(c.Param("Phone")) {
		// 未进行注册, 需要注册
		response.NeedToRegister = true
	}

	// 校验是否需要登陆
	if isSignIn(c.Param("Phone"), moduleInfo.ID) {
		// 未进行登陆, 需要登陆
		response.IsLogin = true
	}

	// 返回会话信息
	c.JSON(http.StatusOK, response)
}
