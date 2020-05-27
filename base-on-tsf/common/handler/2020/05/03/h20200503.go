/*
   @Time : 2020/5/20 11:27 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : h20200503
   @Software: GoLand
*/

// 可用于需要进行参与次数计数的活动
package h20200503

import (
	"github.com/gin-gonic/gin"
	"github.com/offcn-jl/go-common/codes"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/go-common/logger"
	"github.com/offcn-jl/go-common/verify"
	"github.com/offcn-jl/serverless-apis/base-on-tsf/common/database/orm/structs"
	"net/http"
)

// PostAdd 参与活动
// 完成参与后会返回参与活动的状态
// 即, 返回 活动参与总人次, 当前用户参与次数
func PostAdd(c *gin.Context) {
	// 绑定数据
	eventInfo := structs.E20200503{}
	if err := c.ShouldBindJSON(&eventInfo); err != nil {
		// 绑定数据错误
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"Code": codes.InvalidJson, "Error": codes.ErrorText(codes.InvalidJson)})
		return
	}

	// 检查手机号码是否合法
	if !verify.Phone(eventInfo.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "手机号码不正确!"})
		return
	}

	// 添加平台信息
	eventInfo.SourceIP = c.ClientIP()

	// 定义用于 count 的变量 活动参与总人次, 当前用户参与次数
	total, count := 0, 0

	// 获取 活动参与总人次
	orm.PostgreSQL.Model(structs.E20200503{}).Where("event = ?", eventInfo.Event).Count(&total)

	// 获取 当前用户参与次数
	orm.PostgreSQL.Model(structs.E20200503{}).Where("event = ? AND phone = ?", eventInfo.Event, eventInfo.Phone).Count(&count)

	// 保存当前参与记录
	orm.PostgreSQL.Create(&eventInfo)

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"Total": total, "Count": count})
}
