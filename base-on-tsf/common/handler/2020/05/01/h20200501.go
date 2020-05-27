/*
   @Time : 2020/5/4 10:40 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : h20200401
   @Software: GoLand
*/

// 可用于带计数的预约类活动
// 可以灵活的改造成正数、倒数、分批次
// 接口未验证参与者身份、不可设置预约上限, 存在发生超卖的可能性, 所以只可用于接受超卖情况发生的活动
package h20200501

import (
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common/codes"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/serverless-apis/base-on-tsf/common/config"
	"github.com/offcn-jl/serverless-apis/base-on-tsf/common/database/orm/structs"
	"github.com/offcn-jl/serverless-apis/base-on-tsf/common/utils"
	"net/http"
)

// PostSubscribe 处理用户订阅 ( 参与活动 ) 事件
func PostSubscribe(c *gin.Context) {
	// 检查 SQL 注入
	if err := utils.ParameterChecker(c.Params); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Code": codes.InvalidParameterSQLInjection, "Error": codes.ErrorText(codes.InvalidParameterSQLInjection), "Detail": err.Error()})
		return
	}

	// 订阅记录结构
	subscribeInfo := struct {
		ID        uint
		RowNumber uint
	}{}

	// 获取订阅记录
getSubscribeInfo:
	// 先取出参与目标活动的所有手机号, 及按时间升序排列的行号; 然后再取出目标手机号及对应的行号
	orm.PostgreSQL.Raw("SELECT * FROM (SELECT id, ROW_NUMBER() OVER(ORDER BY id ASC), phone FROM e20200501 WHERE \"name\" = ?) as a WHERE a.phone = ?", c.Param("Name"), c.Param("Phone")).Scan(&subscribeInfo)

	// 判断是否存在订阅记录
	if subscribeInfo.RowNumber == 0 {
		// 没有订阅订阅记录, 订阅后重新获取记录
		subscribe := structs.E20200501{}
		// 设置活动名称和手机号码字段的数据
		subscribe.Name = c.Param("Name")
		subscribe.Phone = c.Param("Phone")
		subscribe.ProjectVersion = config.Version
		// 插入数据
		// 将记录插入到数据库后，Gorm 会从数据库加载没有值或值为零值的字段的值
		orm.PostgreSQL.Create(&subscribe)
		// 重新获取订阅记录
		goto getSubscribeInfo
	}

	// 返回订阅记录
	c.JSON(http.StatusOK, gin.H{"Code": 0, "Result": subscribeInfo})
}

// GetCount 获取订阅活动的用户数量
func GetCount(c *gin.Context) {
	count := 0
	orm.PostgreSQL.Model(structs.E20200501{}).Where("name = ?", c.Param("Name")).Count(&count)
	c.JSON(http.StatusOK, gin.H{"Code": 0, "Result": count})
}
