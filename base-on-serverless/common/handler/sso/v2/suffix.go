/*
   @Time : 2020/5/14 1:48 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : suffix
   @Software: GoLand
*/

package sso

import (
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/go-common/logger"
	"github.com/offcn-jl/gscf"
	"github.com/offcn-jl/gscf/fake-http"
)

// GetSuffixList 获取后缀花名册
func GetSuffixList(c *gin.Context) {
	if rows, err := orm.PostgreSQL.Raw("SELECT A.\"id\",A.suffix,A.\"name\",A.crm_user,A.crm_uid,A.crm_channel,b.\"id\",b.f_id,b.code,b.\"name\" FROM single_sign_on_suffixes AS A,single_sign_on_organizations AS b WHERE A.crm_oid=b.\"id\";").Rows(); err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "执行 SQL 查询出错!"})
	} else {
		type Result struct {
			ID         uint   // 后缀 ID
			Suffix     string // 后缀
			Name       string // 后缀名称
			CRMUser    string // CRM 用户名
			CRMUID     uint   // CRM 用户 ID
			CRMChannel uint   // CRM 所属渠道
			CRMOID     uint   // CRM 组织 ID
			CRMOFID    uint   // CRM 上级组织 ID
			CRMOCode   uint   // CRM 组织代码
			CRMOName   string // CRM 组织名称
		}
		results := make([]Result, 0)
		for rows.Next() {
			tempResult := Result{}
			if err := rows.Scan(
				&tempResult.ID,
				&tempResult.Suffix,
				&tempResult.Name,
				&tempResult.CRMUser,
				&tempResult.CRMUID,
				&tempResult.CRMChannel,
				&tempResult.CRMOID,
				&tempResult.CRMOFID,
				&tempResult.CRMOCode,
				&tempResult.CRMOName,
			); err != nil {
				logger.Error(err)
			}
			results = append(results, tempResult)
		}
		if err := rows.Close(); err != nil {
			logger.Error(err)
		}
		c.JSON(http.StatusOK, results)
	}
}
