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
	Version   = "0.1.0"
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
	// 检查是否可以获取到配置信息
	if configer.GetString("TencentCloudAPISecretID", "") == "" {
		logger.Panic(errors.New("未配置 TencentCloudAPISecretID"))
	}
	// 获取配置信息
	TencentCloud.APISecretID = configer.GetString("TencentCloudAPISecretID", "")
	if configer.GetString("TencentCloudAPISecretKey", "") == "" {
		logger.Panic(errors.New("未配置 TencentCloudAPISecretKey"))
	}
	TencentCloud.SecretKey = configer.GetString("TencentCloudAPISecretKey", "")
}
