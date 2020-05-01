/*
   @Time : 2020-04-13 15:28
   @Author : Rebeta
   @Email : master@rebeta.cn
   @File : main
   @Software: GoLand
*/
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	bda "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bda/v20200324"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	fmu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fmu/v20191213"
	"io/ioutil"
	"net/http"
	"os"
)

// 配置信息
var config struct {
	// 腾讯云配置
	tencentCloud struct {
		// API 密钥 ID
		APISecretID string
		// API 密钥 Key
		SecretKey string
	}
}

// 初始化
// 获取配置信息
func init() {
	config.tencentCloud.APISecretID = os.Getenv("TencentCloudAPISecretID")
	config.tencentCloud.SecretKey = os.Getenv("TencentCloudAPISecretKey")
}

// 主函数
func main() {
	r := gin.Default()
	r.POST("/photo/:Beauty", handlePhoto)
	err := r.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func handlePhoto(c *gin.Context) {
	// 从请求 Body 中读取 POST 提交的图片二进制 Buffer
	body, _ := ioutil.ReadAll(c.Request.Body)
	// 将图片 Buffer 编码为 Base64
	base64BodyString := base64.StdEncoding.EncodeToString(body)

	// 生成 腾讯云 API SDK 配置信息
	credential := common.NewCredential(config.tencentCloud.APISecretID, config.tencentCloud.SecretKey)
	cpf := profile.NewClientProfile()

	// 人像分割
	fmt.Println("开始进行 : 人像分割")
	//cpf.HttpProfile.Endpoint = "bda.tencentcloudapi.com"
	// 使用内网地址调用接口
	cpf.HttpProfile.Endpoint = "bda.internal.tencentcloudapi.com"
	client, _ := bda.NewClient(credential, "ap-beijing", cpf)
	request := bda.NewSegmentPortraitPicRequest()
	params := "{\"Image\":\"" + base64BodyString + "\"}"
	err := request.FromJsonString(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	response, err := client.SegmentPortraitPic(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	//fmt.Printf("%s", response.ToJsonString())

	base64Image := *response.Response.ResultImage

	if c.Param("Beauty") == "true" {
		// 人脸美颜
		fmt.Println("开始进行 : 人脸美颜")
		//cpf.HttpProfile.Endpoint = "fmu.tencentcloudapi.com"
		// 使用内网地址调用接口
		cpf.HttpProfile.Endpoint = "fmu.internal.tencentcloudapi.com"
		client, _ := fmu.NewClient(credential, "ap-beijing", cpf)

		request := fmu.NewBeautifyPicRequest()

		params := "{\"Image\":\"" + base64Image + "\",\"Whitening\":50,\"Smoothing\":50,\"FaceLifting\":50,\"EyeEnlarging\":50}"
		err := request.FromJsonString(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		response, err := client.BeautifyPic(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		//fmt.Printf("%s", response.ToJsonString())
		base64Image = *response.Response.ResultImage
	}

	// 将 Base64 编码的图片解码为 Buffer
	imageBuffer, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	// 返回 Buffer
	c.Data(http.StatusOK, "arraybuffer", imageBuffer)
}
