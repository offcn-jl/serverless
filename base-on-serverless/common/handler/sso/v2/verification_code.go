/*
   @Time : 2020/5/12 10:46 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : verification_code
   @Software: GoLand
*/

package sso

import (
	"encoding/json"
	"fmt"
	"github.com/offcn-jl/go-common/codes"
	"github.com/offcn-jl/go-common/configer"
	"github.com/offcn-jl/go-common/database/orm"
	"github.com/offcn-jl/go-common/logger"
	"github.com/offcn-jl/go-common/verify"
	"github.com/offcn-jl/gscf"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"serverless/common/database/orm/structs"
	"time"
)

// PostSendCode 发送验证码接口的处理函数
func PostSendCode(c *gin.Context) {
	// 验证手机号码是否有效
	if !verify.Phone(c.Param("Phone")) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "手机号码不正确!"})
		return
	}

	// 根据登陆模块 ID, 获取登陆模块的配置
	// 需要使用登陆模块配置中的下发平台、签名、模板 ID
	SSOLoginModuleInfo := structs.SingleSignOnLoginModule{}
	orm.PostgreSQL.Where("id = ?", c.Param("LMID")).Find(&SSOLoginModuleInfo)
	if SSOLoginModuleInfo.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "登陆模块配置有误!"})
	} else {
		fmt.Println(SSOLoginModuleInfo)
		switch SSOLoginModuleInfo.Platform {
		case 1:
			// 使用中公短信下发短信
			sendByOffcn(c, SSOLoginModuleInfo.TemplateID, SSOLoginModuleInfo.Term)
		case 2:
			// 使用腾讯云下发短信
			sendByTSmsV2(c, SSOLoginModuleInfo.Sign, SSOLoginModuleInfo.TemplateID, SSOLoginModuleInfo.Term)
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "登陆模块 SMS 平台配置有误!"})
		}
	}
}

// sendByOffcn 使用中公短信平台发送验证码
// 签名无需设置亦无法变更, 所以忽略签名参数
// 平台没有模板的概念, 但是为了更加通用, 内部模拟一套与腾讯云短信服务相同的模板逻辑, 基于格式化输出实现
func sendByOffcn(c *gin.Context, templateID, term uint) {
	// 验证模板 ID 并配置模板内容
	template := ""
	switch templateID {
	case 391863:
		// 验证码 ( 登陆 )
		template = "%d 为您的登录验证码，请于 %d 分钟内填写。如非本人操作，请忽略本短信。"
	case 392030:
		// 通用验证码
		template = "您的验证码是 %d ，请于 %d 分钟内填写。如非本人操作，请忽略本短信。"
	case 392074:
		// 可复用验证码
		template = "您的验证码是 %d ，%d 分钟内可重复使用。如非本人操作，请忽略本短信。"
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "短信模板配置有误"})
		return
	}

	// 获取配置
	apiURL := configer.GetString("OFFCN_SMS_URL", "")
	if apiURL == "" {
		// 未配置 apiURL ( 接口地址 )
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "未配置 OFFCN_SMS_URL"})
		return
	}
	smsName := configer.GetString("OFFCN_SMS_NAME", "")
	if smsName == "" {
		// 未配置 sname
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "未配置 OFFCN_SMS_NAME"})
		return
	}
	smsPwd := configer.GetString("OFFCN_SMS_PWD", "")
	if smsPwd == "" {
		// 未配置 spwd
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "未配置 OFFCN_SMS_PWD"})
		return
	}
	smsTJCode := configer.GetString("OFFCN_SMS_TJ_Code", "")
	if smsTJCode == "" {
		// 未配置 tjcode
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "未配置 OFFCN_SMS_TJ_Code"})
		return
	}

	// 验证是否具有发送条件
	verificationCodeInfo := structs.SingleSignOnVerificationCode{}
	orm.PostgreSQL.Where("phone = ?", c.Param("Phone")).Find(&verificationCodeInfo)
	if verificationCodeInfo.ID != 0 {
		// 存在发送记录, 继续判断是否失效
		duration, _ := time.ParseDuration("-" + fmt.Sprint(verificationCodeInfo.Term) + "m")
		if verificationCodeInfo.CreatedAt.After(time.Now().Add(duration)) {
			// 上一条验证码未超过有效期
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "请勿重复发送验证码!"})
			return
		}
	}

	// 初始化随机数的资源库, 如果不执行这行, 不管运行多少次都返回同样的值 # https://learnku.com/articles/26011
	rand.Seed(time.Now().UnixNano())
	// 生成随机数作为验证码
	code := uint(rand.Intn(8999) + 1000) // 如果直接用 Intn(9999) 会生成出来不是4位的数字

	// 拼接参数
	data := url.Values{
		"sname":   []string{smsName},
		"spwd":    []string{smsPwd},
		"mobile":  []string{c.Param("Phone")},
		"content": []string{fmt.Sprintf(template, code, term)},
		"tjcode":  []string{smsTJCode},
	}

	// 发送短信
	if resp, err := http.PostForm(apiURL, data); err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": codes.InternalErrorUnknown, "Error": codes.ErrorText(codes.InternalErrorUnknown)})
	} else {
		defer resp.Body.Close()
		// 判断有没有发送成功
		if resp.StatusCode != 200 {
			// 请求出错
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": codes.InternalErrorUnknown, "Error": codes.ErrorText(codes.InternalErrorUnknown)})
		} else {
			// 读取 body
			if respBytes, err := ioutil.ReadAll(resp.Body); err != nil {
				logger.Error(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": codes.InternalErrorUnknown, "Error": codes.ErrorText(codes.InternalErrorUnknown)})
			} else {
				// 解码 body
				var respJsonMap map[string]interface{}
				if err := json.Unmarshal(respBytes, &respJsonMap); err != nil {
					logger.Error(err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": codes.InternalErrorUnknown, "Error": codes.ErrorText(codes.InternalErrorUnknown)})
				} else {
					// 返回请求回来的 Json 的 Map
					if respJsonMap["status"].(float64) != 1 {
						// 发送失败
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "[ " + fmt.Sprint(respJsonMap["status"]) + " ] " + fmt.Sprint(respJsonMap["msg"])})
					} else {
						// 发送成功
						// 构造验证码记录
						ssoVerificationCodeInfo := structs.SingleSignOnVerificationCode{
							Phone:    c.Param("Phone"),
							Term:     term,
							Code:     code,
							SourceIP: c.ClientIP(),
						}
						// 从上下文中取出版本信息
						if apiVersion, exist := c.Get("Api-Version"); exist {
							// 存在版本信息, 添加到记录中
							ssoVerificationCodeInfo.ApiVersion = apiVersion.(string)
						} else {
							// 不存在版本信息, 将记录中的版本设置为 Unknown
							ssoVerificationCodeInfo.ApiVersion = "Unknown"
						}
						// 保存记录
						orm.PostgreSQL.Create(&ssoVerificationCodeInfo)
					}
				}
			}
		}
	}
}

// sendByTSmsV2 使用腾讯云短信平台发送验证码
func sendByTSmsV2(c *gin.Context, sign string, templateID, term uint) {
	// 获取配置
	tencentCloudAPISecretID4SMS := configer.GetString("TENCENT_SECRET_ID_SMS", "")
	if tencentCloudAPISecretID4SMS == "" {
		// 未配置 apiURL ( 接口地址 )
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "未配置 TENCENT_SECRET_ID_SMS"})
		return
	}
	tencentCloudSecretKey4SMS := configer.GetString("TENCENT_SECRET_KEY_SMS", "")
	if tencentCloudSecretKey4SMS == "" {
		// 未配置 sname
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "未配置 TENCENT_SECRET_KEY_SMS"})
		return
	}
	tencentCloudSmsSdkAppid := configer.GetString("TENCENT_SMS_APPID", "")
	if tencentCloudSmsSdkAppid == "" {
		// 未配置 spwd
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "未配置 TENCENT_SMS_APPID"})
		return
	}

	// 初始化随机数的资源库, 如果不执行这行, 不管运行多少次都返回同样的值 # https://learnku.com/articles/26011
	rand.Seed(time.Now().UnixNano())
	// 生成随机数作为验证码
	code := uint(rand.Intn(8999) + 1000) // 如果直接用 Intn(9999) 会生成出来不是4位的数字

	// 验证是否具有发送条件
	verificationCodeInfo := structs.SingleSignOnVerificationCode{}
	orm.PostgreSQL.Where("phone = ?", c.Param("Phone")).Find(&verificationCodeInfo)
	if verificationCodeInfo.ID != 0 {
		// 存在发送记录, 继续判断是否失效
		duration, _ := time.ParseDuration("-" + fmt.Sprint(verificationCodeInfo.Term) + "m")
		if verificationCodeInfo.CreatedAt.After(time.Now().Add(duration)) {
			// 上一条验证码未超过有效期
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": "请勿重复发送验证码!"})
			return
		}
	}

	// # https://cloud.tencent.com/document/product/382/43199
	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对 secretId 和 secretKey
	 * 本示例采用从环境变量读取的方式，需要预先在环境变量中设置这两个值
	 * 您也可以直接在代码中写入密钥对，但需谨防泄露，不要将代码复制、上传或者分享给他人
	 * CAM 密匙查询: https://console.cloud.tencent.com/cam/capi*/
	credential := common.NewCredential(tencentCloudAPISecretID4SMS, tencentCloudSecretKey4SMS)

	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK 默认使用 POST 方法
	 * 如需使用 GET 方法，可以在此处设置，但 GET 方法无法处理较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK 有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	//cpf.HttpProfile.ReqTimeout = 5

	/* SDK 会自动指定域名，通常无需指定域名，但访问金融区的服务时必须手动指定域名
	 * 例如 SMS 的上海金融区域名为 sms.ap-shanghai-fsi.tencentcloudapi.com */
	//cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.HttpProfile.Endpoint = "sms.internal.tencentcloudapi.com" // 使用内网接口地址

	/* SDK 默认用 TC3-HMAC-SHA256 进行签名，非必要请不要修改该字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化 SMS 的 client 对象
	 * 第二个参数是地域信息，可以直接填写字符串 ap-guangzhou，或者引用预设的常量 */
	//client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	client, _ := sms.NewClient(credential, regions.Beijing, cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	   * 您可以直接查询 SDK 源码确定接口有哪些属性可以设置
	    * 属性可能是基本类型，也可能引用了另一个数据结构
	    * 推荐使用 IDE 进行开发，可以方便地跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK 采用的是指针风格指定参数，即使对于基本类型也需要用指针来对参数赋值。
	 * SDK 提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台：https://console.cloud.tencent.com/smsv2
	 * sms helper：https://cloud.tencent.com/document/product/382/3773 */

	/* 短信应用 ID: 在 [短信控制台] 添加应用后生成的实际 SDKAppID，例如1400006666 */
	request.SmsSdkAppid = common.StringPtr(tencentCloudSmsSdkAppid)
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，可登录 [短信控制台] 查看签名信息 */
	request.Sign = common.StringPtr(sign)
	/* 国际/港澳台短信 senderid: 国内短信填空，默认未开通，如需开通请联系 [sms helper] */
	//request.SenderId = common.StringPtr("xxx")
	/* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	//request.SessionContext = common.StringPtr("xxx")
	/* 短信码号扩展号: 默认未开通，如需开通请联系 [sms helper] */
	//request.ExtendCode = common.StringPtr("0")
	/* 模板参数: 若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs([]string{fmt.Sprint(code), fmt.Sprint(term)})
	/* 模板 ID: 必须填写已审核通过的模板 ID，可登录 [短信控制台] 查看模板 ID */
	request.TemplateID = common.StringPtr(fmt.Sprint(templateID))
	/* 下发手机号码，采用 e.164 标准，+[国家或地区码][手机号]
	 * 例如+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + c.Param("Phone")})

	// 通过 client 对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": -1, "Error": "[ " + err.(*errors.TencentCloudSDKError).GetCode() + " ] " + err.(*errors.TencentCloudSDKError).GetMessage()})
		return
	}
	// 非 SDK 异常，直接失败。实际代码中可以加入其他的处理
	if err != nil {
		logger.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Code": codes.InternalErrorUnknown, "Error": codes.ErrorText(codes.InternalErrorUnknown)})
		return
	}
	if *response.Response.SendStatusSet[0].Code != "Ok" {
		logger.DebugToJson("Response", response.Response)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Code": -1, "Error": *response.Response.SendStatusSet[0].Message})
		return
	} else {
		// 发送成功
		// 构造验证码记录
		ssoVerificationCodeInfo := structs.SingleSignOnVerificationCode{
			Phone:    c.Param("Phone"),
			Term:     term,
			Code:     code,
			SourceIP: c.ClientIP(),
		}
		// 从上下文中取出版本信息
		if apiVersion, exist := c.Get("Api-Version"); exist {
			// 存在版本信息, 添加到记录中
			ssoVerificationCodeInfo.ApiVersion = apiVersion.(string)
		} else {
			// 不存在版本信息, 将记录中的版本设置为 Unknown
			ssoVerificationCodeInfo.ApiVersion = "Unknown"
		}
		// 保存记录
		orm.PostgreSQL.Create(&ssoVerificationCodeInfo)
	}
}
