/*
   @Time : 2020/4/22 3:17 下午
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : orm
   @Software: GoLand
*/

package orm

import (
	"database/sql/driver"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"reflect"
	"regexp"
	"serverless/go-common/database"
	"serverless/go-common/logger"
	"strconv"
	"time"
	"unicode"
)

var Version = "0.1.0"

type orm struct {
	PostgreSQL postgreSQL
}

type postgreSQL struct {
	Marketing *gorm.DB
}

func New() *orm {
	orm := &orm{}
	dsn := database.GetDSN()
	var err error
	if orm.PostgreSQL.Marketing, err = gorm.Open("postgres", dsn.PostgreSQL); err != nil {
		logger.Panic(err)
	}
	return orm
}

func (o *orm) Close() {
	o.PostgreSQL.Marketing.Close()
}

// 替换 grom 的 LogFormatter , 移除其中输出日志颜色的逻辑, 优化 gorm 日志的可读性
var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func init() {
	gorm.LogFormatter = func(values ...interface{}) (messages []interface{}) {
		if len(values) > 1 {
			var (
				sql             string
				formattedValues []string
				level           = values[0]
				currentTime     = "\n[" + gorm.NowFunc().Format("2006-01-02 15:04:05") + "]"
				source          = fmt.Sprintf("(%v)", values[1])
			)

			messages = []interface{}{source, currentTime}

			if len(values) == 2 {
				//remove the line break
				currentTime = currentTime[1:]
				//remove the brackets
				source = fmt.Sprintf("%v", values[1])

				messages = []interface{}{currentTime, source}
			}

			if level == "sql" {
				// duration
				messages = append(messages, fmt.Sprintf(" [%.2fms] ", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))
				// sql

				for _, value := range values[4].([]interface{}) {
					indirectValue := reflect.Indirect(reflect.ValueOf(value))
					if indirectValue.IsValid() {
						value = indirectValue.Interface()
						if t, ok := value.(time.Time); ok {
							if t.IsZero() {
								formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
							} else {
								formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
							}
						} else if b, ok := value.([]byte); ok {
							if str := string(b); isPrintable(str) {
								formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
							} else {
								formattedValues = append(formattedValues, "'<binary>'")
							}
						} else if r, ok := value.(driver.Valuer); ok {
							if value, err := r.Value(); err == nil && value != nil {
								formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
							} else {
								formattedValues = append(formattedValues, "NULL")
							}
						} else {
							switch value.(type) {
							case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
								formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
							default:
								formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
							}
						}
					} else {
						formattedValues = append(formattedValues, "NULL")
					}
				}

				// differentiate between $n placeholders or else treat like ?
				if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
					sql = values[3].(string)
					for index, value := range formattedValues {
						placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
						sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
					}
				} else {
					formattedValuesLength := len(formattedValues)
					for index, value := range sqlRegexp.Split(values[3].(string), -1) {
						sql += value
						if index < formattedValuesLength {
							sql += formattedValues[index]
						}
					}
				}

				messages = append(messages, sql)
				messages = append(messages, fmt.Sprintf(" \n[%v] ", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned "))
			} else {
				//messages = append(messages, "\033[31;1m")
				messages = append(messages, values[2:]...)
				//messages = append(messages, "\033[0m")
			}
		}

		return
	}
}
