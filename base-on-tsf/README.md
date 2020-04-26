# 基于 TSF 的营销接口

## 食用须知
1. [TSF](https://cloud.tencent.com/document/product/649) 、[SCF](https://cloud.tencent.com/document/product/583)、[TKE](https://cloud.tencent.com/document/product/457) 等产品名称、品牌或商标的一切权利归[腾讯云](https://cloud.tencent.com)所有。
1. Serverless、Serverless Framework 是 [serverless.com](https://serverless.com) 的产品。
1. 本项目整体基于 MIT 协议开源。
1. 本项目主要包括两个主要分支 : master ( 主分支, 可用于生产环境 )、 new-feature ( 新功能分支, 包含处于测试和验证阶段的新功能 )。
1. 原则上，本项目中的所有模块所使用的各种密码、令牌、口令等敏感信息均配置在环境变量、启动脚本或配置文件中。建议使用者在二次开发的过程中采取同样的方式保存相关设置，并在使用过程中请注意妥善保管相关信息。

Enjoy it. XD

## 目录结构 ( 按文件名排序 )

|--  base-on-tsf  // 基于 [TSF ( 腾讯微服务平台 )](https://cloud.tencent.com/document/product/649) 的营销接口  
|	|-- artifacts // 制品 ( 用于最终部署 )  
|		|-- cmdline // 用于检查应用进程是否存在，没有.sh后缀  
|		|-- main // 主程序，执行构建步骤后生成  
|		|-- result.zip // 程序包，用于最终部署  
|		|-- start.example.sh // 示例启动脚本，需要自行修改为启动脚本  
|		|-- start.sh // 启动脚本  
|		|-- stop.sh // 停止脚本  
|	|-- .gitignore // GIT 配置文件，用于配置需要忽略提交的内容  
|	|-- build.sh // 构建脚本  
|	|-- go.mod // go mod 配置文件  
|	|-- go.sum // go mod checksums  
|	|-- main.go // 主程序源码  
|	|-- README.md // 使用说明  

## 功能

1. 照片处理 ( 证件照生成 ) [ [腾讯云人像分割](https://cloud.tencent.com/document/product/1208/42970) & [腾讯云人脸美颜](https://cloud.tencent.com/document/product/1172/40715) ]

## 部署
1. 准备工作
	1. [创建 API 密钥](https://console.cloud.tencent.com/capi) （ 注意妥善保管 ）
	1. [开通腾讯微服务平台 ( TSF )](https://cloud.tencent.com/document/product/649)
	1. 初始化设置
		1. 创建集群
			1. 进入 [TSF 集群管理](https://console.cloud.tencent.com/tsf/cluster)
			1. 新建 Serverless 类型的集群
		1. [创建 TSF 应用](https://console.cloud.tencent.com/tsf/app)
1. 安装 Golang
1. 按需修改代码
1. 修改配置
	1. 进入程序包目录 ( artifacts )
	1. 复制一份启动脚本模板 ( start.example.sh  )
	1. 更名为可用的启动脚本 ( start.sh )
	1. 打开启动脚本，将其中的 TencentCloudAPISecretID 字段后的 your_id 替换为 API 密钥 ID
	1. 打开启动脚本，将其中的 TencentCloudAPISecretKey 字段后的 your_key 替换为 API 密钥 Key
1. 构建
	1. MacOS / Linux
		1. 直接在本目录执行 sh build.sh
	2. Windows
		1. 交叉编译 Linux 64 位二进制可执行程序，在本目录执行 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o artifacts/main
		1. 压缩程序包，压缩 artifacts 目录中的 cmdline main start.sh stop.sh 四个文件
1. 上传代码包并部署应用
	1. 进入 [TSF 应用管理](https://console.cloud.tencent.com/tsf/app)
	1. 点击第一步中创建的应用名称，进入应用管理
	1. 点击程序包管理选项卡，上传程序包
	1. 点击部署组选项卡，创建部署组
	1. 点击上一步创建的部署组名称，进入部署组管理
	1. 点击外网访问选项卡，打开外网访问开关
	1. 外网访问路径即接口根域名
1. 完成部署
