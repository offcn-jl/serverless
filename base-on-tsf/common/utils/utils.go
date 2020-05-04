/*
   @Time : 2020/5/4 1:47 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : utils
   @Software: GoLand
*/

package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common/logger"
	"regexp"
	"strings"
)

// 参数检查器
func ParameterChecker(parameters interface{}) (err error) {
	switch param := parameters.(type) {
	case gin.Params:
		for _, value := range param {
			if pass, err := CheckSqlInject(value.Value); err != nil {
				logger.Error(err)
				return errors.New("[ FAILURE ] 参数 " + value.Key + " 检查 SQL 注入失败")
			} else {
				if !pass {
					return errors.New("[ FAILURE ] 参数 " + value.Key + " 存在 SQL 注入")
				}
			}
		}
	default:
		err = errors.New("[ FAILURE ] 类型错误, SQL 注入检查失败")
	}
	return
}

// SQL注入过滤 # https://www.cnblogs.com/mafeng/p/6207988.html
func CheckSqlInject(parameter string) (pass bool, err error) {
	// 关键字过滤
	// 正则条件不能用 "" 因为 "" 内的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return !re.MatchString(strings.ToLower(parameter)), nil
}
