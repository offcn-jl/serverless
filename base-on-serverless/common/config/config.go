/*
   @Time : 2020/5/4 5:46 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : config
   @Software: GoLand
*/

package config

import (
	"github.com/offcn-jl/go-common/logger"
)

var (
	Project = "SCF-APIs"
	Version = "0.2.0"
)

// 初始化
func init() {
	// 打印版本信息
	logger.Log("Project " + Project + " Ver." + Version)
}
