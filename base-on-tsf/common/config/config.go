/*
   @Time : 2020/5/4 9:16 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : config
   @Software: GoLand
*/

package config

import (
	"errors"
	"github.com/offcn-jl/go-common/configer"
	"github.com/offcn-jl/go-common/logger"
)

var (
	Project   = "TSF-APIs"
	Version   = "0.4.1"
	builtTime = ""
)

// 腾讯云配置
var TencentCloud struct {
	// API 密钥 ID
	APISecretID string
	// API 密钥 Key
	SecretKey string
}

// 初始化
func init() {
	// 拼接版本信息
	Version = Version + builtTime
	// 打印版本信息
	logger.Log("Project " + Project + " Ver." + Version)
	// 获取配置信息, 并检查是否正确获取
	if TencentCloud.APISecretID = configer.GetString("TENCENT_SECRET_ID", ""); TencentCloud.APISecretID == "" {
		logger.Panic(errors.New("未配置 TENCENT_SECRET_ID"))
	}
	// 获取配置信息, 并检查是否正确获取
	if TencentCloud.SecretKey = configer.GetString("TENCENT_SECRET_KEY", ""); TencentCloud.SecretKey == "" {
		logger.Panic(errors.New("未配置 TENCENT_SECRET_KEY"))
	}
}
