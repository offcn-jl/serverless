# 基于 Serverless Framework 的接口

## 食用须知
1. [TSF](https://cloud.tencent.com/document/product/649) 、 [SCF](https://cloud.tencent.com/document/product/583) 、 [TKE](https://cloud.tencent.com/document/product/457) 等产品名称、品牌或商标的一切权利归 [腾讯云](https://cloud.tencent.com) 所有。
1. Serverless、Serverless Framework 是 [serverless.com](https://serverless.com) 的产品。
1. 本项目整体基于 MIT 协议开源。
1. 本项目主要包括两个主要分支 : master ( 主分支, 可用于生产环境 )、 new-feature ( 新功能分支, 包含处于测试和验证阶段的新功能 )。
1. 原则上，本项目中的所有模块所使用的各种密码、令牌、口令等敏感信息均配置在环境变量、启动脚本或配置文件中。建议使用者在二次开发的过程中采取同样的方式保存相关设置，并在使用过程中请注意妥善保管相关信息。

Enjoy it. XD

## 功能

1. 礼品功能 ( 剩余数量查询、领取、查询、消费 )  
1. 应用管理功能 ( 版本控制 )  

## 目录结构 ( 按文件名排序 )

- base-on-serverless ( 基于 [Serverless Framework](https://serverless.com) 及 [腾讯云 Serverless](https://cloud.tencent.com/product/sls) 的接口 )  
    - common ( 组件库 )  
        - config ( 配置库 )  
            - config.go ( 配置库 )  
        - handler ( 处理函数包 )  
            - sso ( 单点登录系统接口的处理函数包 )  
                - v2 ( 第二代单点登录系统接口的处理函数包 )  
                    - auth.go ( 鉴权相关接口的处理函数包 ) 
                    - auth_test.go ( 鉴权相关接口的处理函数包的单元测试 ) 
                    - push.go ( 推送相关接口的处理函数包 ) 
                    - push_test.go ( 推送相关接口的处理函数包的单元测试 ) 
                    - verification_code.go ( 验证码相关接口的处理函数包 ) 
        - middleware ( 处理函数库 )  
            - middleware.go ( 处理函数库 )  
    - db-postgre-sql ( Serverless 组件, 数据库, Serverless PG 数据库 )  
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-sso-v2-auth-sign-in ( Serverless 组件, 云函数, 第二代单点登陆, 鉴权, 登陆 )  
        - artifacts ( 制品目录, 执行构建步骤后自动生成 )
            - main ( 主程序二进制文件, 执行构建步骤后自动生成 )
            - phone.dat ( 手机号码数据库, 需要自行从 github.com/xluohome/phonedata 包中复制到此处 )
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-sso-v2-auth-sign-up ( Serverless 组件, 云函数, 第二代单点登陆, 鉴权, 注册 )  
        - artifacts ( 制品目录, 执行构建步骤后自动生成 )
            - main ( 主程序二进制文件, 执行构建步骤后自动生成 )
            - phone.dat ( 手机号码数据库, 需要自行从 github.com/xluohome/phonedata 包中复制到此处 )
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-sso-v2-crm-push ( Serverless 组件, 云函数, 第二代单点登陆, CRM, 推送 )  
        - artifacts ( 制品目录, 执行构建步骤后自动生成 )
            - main ( 主程序二进制文件, 执行构建步骤后自动生成 )
            - phone.dat ( 手机号码数据库, 需要自行从 github.com/xluohome/phonedata 包中复制到此处 )
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-sso-v2-send-verification-code ( Serverless 组件, 云函数, 第二代单点登陆, 验证码, 发送 )  
        - artifacts ( 制品目录, 执行构建步骤后自动生成 )
            - main ( 主程序二进制文件, 执行构建步骤后自动生成 )
            - phone.dat ( 手机号码数据库, 需要自行从 github.com/xluohome/phonedata 包中复制到此处 )
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - vpc ( Serverless 组件, VPC, 仅提供演示配置 )
        -   serverless.yml.example ( 演示配置 )
    - .env.dev.example ( 测试环境的环境变量, 演示配置 )
    - .env.prod.example ( 生产环境的环境变量, 演示配置 )
    - .gitignore ( GIT 配置文件，用于配置需要忽略提交的内容 )  
    - build.sh ( 构建脚本 )  
    - deploy.sh ( 部署脚本 )  
    - go.mod ( go mod 配置文件 )  
    - go.sum ( go mod checksums )  
    - README.md ( 使用说明 )  

## 部署
1. 准备工作
	1. 创建 [API 密钥](https://console.cloud.tencent.com/capi) ( 注意妥善保管 )
	1. 创建 [API 网关](https://console.cloud.tencent.com/apigateway/index)
	1. [ 可选, 建议 ] 创建 [私有网络 ( VPC )](https://console.cloud.tencent.com/vpc/vpc) ( 使用 VPC Component 自动编排时存在会为不同的 stage 创建不同的 vpc 的问题; 还存在修改 name 后, 会创建并使用新的 vpc 的问题; 后续使用会存在很多不确定性, 并且在与 TKE 或 TSF 等业务进行打通时会增加很多的工作量; 如果对 stage 的隔离性有更高要求的话, 可以选择使用 Component 自动编排 )
1. 安装 Golang
1. 按需修改代码 ( 不需要的组件可直接删除 )
1. 修改配置
	1. 复制环境变量配置文件模板 ( .env.dev.example 及 .env.prod.example )
	1. 更名为可用的环境变量配置文件 ( .env.dev 及 .env.prod )
	1. 打开生产环境的环境变量配置文件 ( .env.prod )，将其中各个字段后的 your_xxx 替换为对应的数据
	1. 打开测试环境的环境变量配置文件 ( .env.dev )，将其中各个字段后的 your_xxx 替换为对应的数据
1. 部署
	1. MacOS / Linux
		1. 直接在本目录执行 sh deploy.sh ( 测试环境 )
		1. 直接在本目录执行 sh deploy.sh prod ( 生产环境 )
	1. Windows  
	    1. 手动进行构建及部署操作
1. 完成部署
    1. Serverless Framework 会在每个组件部署成功后返回访问方式等信息 ( output )
