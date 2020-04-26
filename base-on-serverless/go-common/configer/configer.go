/*
   @Time : 2020/4/22 4:38 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : configer
   @Software: GoLand
*/

package configer

import (
	"os"
	"strings"
)

var Version = "0.1.0"

type config struct {
	Project string // 项目名
	Version string // 项目版本
	Debug   bool   // 调试模式
}

// 项目配置
var Conf = config{
	Project: "Chaos-Serverless",
	Version: "0.1.0",
	Debug:   GetBool("Debug", false),
}

/**
 * 获取字符型系统常量
 */
func GetString(parameter string, def string) string {
	if os.Getenv(parameter) != "" {
		return os.Getenv(parameter)
	} else {
		return def
	}
}

/**
 * 获取 Bool 型系统常量
 */
func GetBool(parameter string, def bool) bool {
	if os.Getenv(parameter) != "" {
		if strings.ToLower(os.Getenv(parameter)) == "true" {
			return true
		} else {
			return false
		}
	} else {
		return def
	}
}
