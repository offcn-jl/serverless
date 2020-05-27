/*
   @Time : 2020/5/22 4:02 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : push
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
	"net/http"
	"net/url"
	"tsf/common/database/orm/structs"
)

// PostPush 推送 CRM 数据并保存记录
func PostPush(c *gin.Context) {
	pushInfo := structs.SingleSignOnPushLog{}
	// 绑定数据
	if err := c.ShouldBindJSON(&pushInfo); err != nil {
		// 绑定数据错误
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"Code": codes.InvalidJson, "Error": codes.ErrorText(codes.InvalidJson)})
		return
	}

	// 验证手机号码是否有效
	if !verify.Phone(pushInfo.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "手机号码不正确!"})
		return
	}

	// 检查是否已经进行过推送
	pushInfo4Check := structs.SingleSignOnPushLog{}
	orm.PostgreSQL.Where("crm_sid = ? AND phone = ?", pushInfo.CRMSID, pushInfo.Phone).Find(&pushInfo4Check)
	if pushInfo4Check.ID != 0 {
		// 已经进行过推送, 跳过后续步骤
		c.JSON(http.StatusOK, gin.H{"Code": 0})
		return
	}

	// 校验后缀
	if pushInfo.ActualSuffix == "" {
		// 后缀未填写, 使用默认后缀配置
		tempSession := structs.SingleSignOnSession{Phone: pushInfo.Phone}
		getDefaultSuffix(&tempSession)
		pushInfo.CRMChannel = tempSession.CRMChannel
		pushInfo.CurrentSuffix = tempSession.CurrentSuffix
		pushInfo.CRMUID = tempSession.CRMUID
		pushInfo.CRMOCode = tempSession.CRMOCode
	} else {
		suffixInfo := structs.SingleSignOnSuffix{}
		orm.PostgreSQL.Where("suffix = ?", pushInfo.ActualSuffix).Find(&suffixInfo)
		if suffixInfo.ID == 0 {
			// 后缀无效, 使用默认后缀配置
			tempSession := structs.SingleSignOnSession{Phone: pushInfo.Phone}
			getDefaultSuffix(&tempSession)
			pushInfo.CRMChannel = tempSession.CRMChannel
			pushInfo.CurrentSuffix = tempSession.CurrentSuffix
			pushInfo.CRMUID = tempSession.CRMUID
			pushInfo.CRMOCode = tempSession.CRMOCode
		} else {
			pushInfo.CRMChannel = suffixInfo.CRMChannel
			pushInfo.CRMUID = suffixInfo.CRMUID
			pushInfo.CurrentSuffix = suffixInfo.Suffix
			if suffixInfo.CRMOID > 1 {
				// 配置了 CRMOID 并且不是省级
				organizationInfo := structs.SingleSignOnOrganization{}
				orm.PostgreSQL.Where("id = ?", suffixInfo.CRMOID).Find(&organizationInfo)
				pushInfo.CRMOCode = organizationInfo.Code
			} else {
				// 未配置 CRMOID 或者是省级 ( 等于 1 ), 按手机号码归属地分配 CRM 信息
				tempSession := structs.SingleSignOnSession{Phone: pushInfo.Phone}
				distributionByPhoneNumber(&tempSession)
				pushInfo.CRMOCode = tempSession.CRMOCode
			}
		}
	}

	// 推送信息到 CRM
	if urlObject, err := url.Parse("https://dc.offcn.com:8443/a.gif"); err != nil {
		// 请求失败
		logger.Error(err)
		// 推送失败, 保存推送失败记录
		orm.PostgreSQL.Create(&structs.SingleSignOnErrorLog{
			Phone:      pushInfo.Phone,
			MID:        0, // 0 代表本推送接口
			CRMChannel: pushInfo.CRMChannel,
			CRMUID:     pushInfo.CRMUID,
			CRMOCode:   pushInfo.CRMOCode,
			Error:      "推送接口 > Parse GET URL Fail : " + err.Error(),
		})

	} else {
		// 构建参数 queryObject
		queryObject := urlObject.Query()
		queryObject.Set("sid", pushInfo.CRMSID)
		queryObject.Set("mobile", pushInfo.Phone)
		queryObject.Set("channel", fmt.Sprint(pushInfo.CRMChannel))
		queryObject.Set("orgn", fmt.Sprint(pushInfo.CRMOCode))
		if pushInfo.CRMUID != 0 {
			queryObject.Set("owner", fmt.Sprint(pushInfo.CRMUID))
		}
		if pushInfo.CustomerName != "" {
			queryObject.Set("name", pushInfo.CustomerName)
		}
		if pushInfo.CustomerIdentityID != 0 {
			queryObject.Set("khsf", fmt.Sprint(pushInfo.CustomerIdentityID))
		}
		if pushInfo.CustomerColleage != "" {
			queryObject.Set("colleage", pushInfo.CustomerColleage)
		}
		if pushInfo.CustomerMayor != "" {
			queryObject.Set("mayor", pushInfo.CustomerMayor)
		}
		if pushInfo.Remark != "" {
			queryObject.Set("remark", pushInfo.Remark)
		}
		// 发送 GET 请求
		urlObject.RawQuery = queryObject.Encode()
		if getResponse, err := http.Get(urlObject.String()); err != nil {
			// 发送 GET 请求出错
			logger.Error(err)
			// 推送失败, 保存推送失败记录
			orm.PostgreSQL.Create(&structs.SingleSignOnErrorLog{
				Phone:      pushInfo.Phone,
				MID:        0, // 0 代表本推送接口
				CRMChannel: pushInfo.CRMChannel,
				CRMUID:     pushInfo.CRMUID,
				CRMOCode:   pushInfo.CRMOCode,
				Error:      "推送接口 > GET Fail : " + err.Error(),
			})
		} else {
			if getResponse.StatusCode != 200 {
				// 请求出错
				logger.Error(err)
				// 推送失败, 保存推送失败记录
				orm.PostgreSQL.Create(&structs.SingleSignOnErrorLog{
					Phone:      pushInfo.Phone,
					MID:        0, // 0 代表本推送接口
					CRMChannel: pushInfo.CRMChannel,
					CRMUID:     pushInfo.CRMUID,
					CRMOCode:   pushInfo.CRMOCode,
					Error:      "推送接口 > Json POST Response Fail : " + err.Error(),
				})
			}
		}
	}

	// 保存会话信息
	orm.PostgreSQL.Create(&pushInfo)

	// 返回状态码 0 代表成功
	c.JSON(http.StatusOK, gin.H{"Code": 0})
}
