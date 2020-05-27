/*
   @Time : 2020/5/20 2:22 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : h20200503_test
   @Software: GoLand
*/

package h20200503

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPatchAdd(t *testing.T) {
	// 创建上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 生成随机活动名
	evnetName := "TEST-" + time.Now().Format("20060102-150405")

	// 测试 未绑定 Body 数据
	PatchAdd(c)
	assert.Contains(t, w.Body.String(), "Invalid Json Data")
	// 增加 Body
	c.Request, _ = http.NewRequest("PATCH", "/", bytes.NewBufferString("{\"Event\":\""+evnetName+"\",\"Phone\":\"*\"}"))

	// 测试手机号码正确
	w.Body.Reset() // 再次测试前重置 body
	PatchAdd(c)
	assert.Contains(t, w.Body.String(), "手机号码不正确")
	// 修正手机号码
	c.Request, _ = http.NewRequest("PATCH", "/", bytes.NewBufferString("{\"Event\":\""+evnetName+"\",\"Phone\":\"17866668888\"}"))

	// 测试 第一个用户第一次参加活动
	w.Body.Reset() // 再次测试前重置 body
	PatchAdd(c)
	assert.Equal(t, "{\"Count\":0,\"Total\":0}", w.Body.String())
	c.Request, _ = http.NewRequest("PATCH", "/", bytes.NewBufferString("{\"Event\":\""+evnetName+"\",\"Phone\":\"17866668888\"}"))

	// 测试 第一个用户第二次参加活动
	w.Body.Reset() // 再次测试前重置 body
	PatchAdd(c)
	assert.Equal(t, "{\"Count\":1,\"Total\":1}", w.Body.String())
	// 修改手机号码为第二个用户的号码
	c.Request, _ = http.NewRequest("PATCH", "/", bytes.NewBufferString("{\"Event\":\""+evnetName+"\",\"Phone\":\"17866886688\"}"))

	// 测试 第二个用户第一次参加活动
	w.Body.Reset() // 再次测试前重置 body
	PatchAdd(c)
	assert.Equal(t, "{\"Count\":0,\"Total\":2}", w.Body.String())
	c.Request, _ = http.NewRequest("PATCH", "/", bytes.NewBufferString("{\"Event\":\""+evnetName+"\",\"Phone\":\"17866886688\"}"))

	// 测试 第二个用户第二次参加活动
	w.Body.Reset() // 再次测试前重置 body
	PatchAdd(c)
	assert.Equal(t, "{\"Count\":1,\"Total\":3}", w.Body.String())
}
