/*
   @Time : 2020/4/22 3:38 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : database
   @Software: GoLand
*/

package database

import (
	"errors"
	"github.com/offcn-jl/go-common/configer"
	"github.com/offcn-jl/go-common/logger"
)

/**
 * 数据库配置
 */
type DSN struct {
	PostgreSQL string
}

func GetDSN() (d DSN) {
	if configer.GetString("PostgreSQLHost", "") == "" {
		logger.Panic(errors.New("未配置 PostgreSQLHost"))
	}
	if configer.GetString("PostgreSQLPort", "") == "" {
		logger.Panic(errors.New("未配置 PostgreSQLPort"))
	}
	if configer.GetString("PostgreSQLUser", "") == "" {
		logger.Panic(errors.New("未配置 PostgreSQLUser"))
	}
	if configer.GetString("PostgreSQLDBName", "") == "" {
		logger.Panic(errors.New("未配置 PostgreSQLDBName"))
	}
	if configer.GetString("PostgreSQLPassword", "") == "" {
		logger.Panic(errors.New("未配置 PostgreSQLPassword"))
	}
	if configer.GetString("PostgreSQLSSLMode", "") == "" {
		logger.Panic(errors.New("未配置 PostgreSQLSSLMode"))
	}
	d.PostgreSQL = "host=" + configer.GetString("PostgreSQLHost", "") + " port=" + configer.GetString("PostgreSQLPort", "") + " user=" + configer.GetString("PostgreSQLUser", "") + " dbname=" + configer.GetString("PostgreSQLDBName", "") + " password=" + configer.GetString("PostgreSQLPassword", "") + " sslmode=" + configer.GetString("PostgreSQLSSLMode", "")
	return
}
