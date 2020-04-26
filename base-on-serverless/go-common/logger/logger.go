/*
   @Time : 2020/4/19 9:36 上午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : logger
   @Software: GoLand
*/

package logger

import (
	"encoding/json"
	"fmt"
	"github.com/offcn-jl/chaos-go-scf"
	"runtime/debug"
	"serverless/go-common/configer"
	"time"
)

var Version = "0.1.0"

func init() {
	Log("Project " + configer.Conf.Project + " Ver." + configer.Conf.Version)
}

/**
 * 日志
 */
func Log(log string) {
	_, _ = fmt.Fprintf(chaos.DefaultWriter, "[%s-Log] %v | %s\n",
		configer.Conf.Project,
		time.Now().Format("2006/01/02 - 15:04:05"),
		log,
	)
}

/**
 * 错误
 */
func Error(err error) {
	_, _ = fmt.Fprintf(chaos.DefaultWriter, "[%s-Error] %v | %s\n",
		configer.Conf.Project,
		time.Now().Format("2006/01/02 - 15:04:05"),
		err,
	)
	_, _ = fmt.Fprintf(chaos.DefaultWriter, "[%s-Error-Stacks] %v\n%s\n",
		configer.Conf.Project,
		time.Now().Format("2006/01/02 - 15:04:05"),
		debug.Stack(), // 输出调用堆栈
	)
	// debug.PrintStack() // 打印调用堆栈
}

/**
 * Panic
 */
func Panic(err error) {
	_, _ = fmt.Fprintf(chaos.DefaultWriter, "[%s-Error] %v | %s\n",
		configer.Conf.Project,
		time.Now().Format("2006/01/02 - 15:04:05"),
		err,
	)
	_, _ = fmt.Fprintf(chaos.DefaultWriter, "[%s-error-Stacks] %v\n%s\n",
		configer.Conf.Project,
		time.Now().Format("2006/01/02 - 15:04:05"),
		debug.Stack(), // 输出调用堆栈
	)
	panic(err)
}

/**
 * 调试输出为 Json 字符串
 */
func DebugToJson(name string, parameters interface{}) {
	if configer.Conf.Debug {
		jsonStrings, _ := json.Marshal(parameters)
		_, _ = fmt.Fprintf(chaos.DefaultWriter, "[%s-Debug-Json] %v | %s --> %s\n",
			configer.Conf.Project,
			time.Now().Format("2006/01/02 - 15:04:05"),
			name,
			jsonStrings,
		)
	}
}

/**
 * 调试输出为字符串
 */
func DebugToString(name string, str interface{}) {
	if configer.Conf.Debug {
		_, _ = fmt.Fprintf(chaos.DefaultWriter, "[%s-Debug-Sting] %v | %s --> %s\n",
			configer.Conf.Project,
			time.Now().Format("2006/01/02 - 15:04:05"),
			name,
			str,
		)
	}
}
