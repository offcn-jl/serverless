/*
   @Time : 2020/5/14 3:23 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : suffix_test
   @Software: GoLand
*/

package sso

import (
	"github.com/offcn-jl/gscf"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试 获取后缀花名册
func TestGetSuffixList(t *testing.T) {
	// 初始化上下文
	c := gin.Context{}
	c.Response.Headers = make(map[string]string) // 空 map 需要初始化后才可以使用

	GetSuffixList(&c)
	assert.Contains(t, c.Response.Body, "{\"ID\":1,\"Suffix\":\"default\",\"Name\":\"默认后缀(ID=1)\",\"CRMUser\":\"default\",\"CRMUID\":32431,\"CRMChannel\":7,\"CRMOID\":1,\"CRMOFID\":0,\"CRMOCode\":22,\"CRMOName\":\"吉林分校\"},{\"ID\":2,\"Suffix\":\"test\",\"Name\":\"后缀 1\",\"CRMUser\":\"test\",\"CRMUID\":123,\"CRMChannel\":104,\"CRMOID\":2,\"CRMOFID\":1,\"CRMOCode\":2290,\"CRMOName\":\"吉林长春分校\"}")
}
