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
        - database ( 数据库 )  
            - orm ( ORM 库 )  
                - auto-migrate-tool ( 自动迁移工具 )  
                    - main.go ( 自动迁移工具主程序 )  
                - structs ( 结构体库 )  
                    - 2020-gift.go ( 礼品功能的结构体 )  
                    - apps.go ( 应用管理功能的结构体 )  
        - middleware ( 处理函数库 )  
            - middleware.go ( 处理函数库 )  
    - db-postgre-sql ( Serverless 组件, 数据库, Serverless PG 数据库 )  
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-2020-gift-checkout ( Serverless 组件, 云函数, 查询礼品信息接口 )  
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-2020-gift-consume ( Serverless 组件, 云函数, 消费礼品接口 )  
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-2020-gift-get ( Serverless 组件, 云函数, 获取礼品接口 )  
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-2020-gift-surplus ( Serverless 组件, 云函数, 查询礼品剩余数量接口 )  
        - main.go ( 主程序 )
        - serverless.yml ( Serverless 组件配置文件 )
    - scf-get-app-version-info ( Serverless 组件, 云函数, 获取应用版本控制信息接口 )  
        - main.go ( 主程序 )
        - main_test.go ( 主程序单元测试 )
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
	1. [创建 API 密钥](https://console.cloud.tencent.com/capi) （ 注意妥善保管 ）
	1. 获取 API 网关的 Service ID, 替换到每个组件的配置文件 ( serverless.yml ) 中
1. 安装 Golang
1. 按需修改代码 ( 不需要的组件可直接删除 )
1. 修改配置
	1. 复制环境变量配置文件模板 ( .env.dev.example 及 .env.prod.example )
	1. 更名为可用的环境变量配置文件 ( .env.dev 及 .env.prod )
	1. 打开生产环境的环境变量配置文件 ( .env.prod )，将其中各个字段后的 your_xxx 替换为对应的数据
	1. 打开测试环境的环境变量配置文件 ( .env.dev )，将其中各个字段后的 your_xxx 替换为对应的数据
1. 部署
	1. MacOS / Linux
		1. 直接在本目录执行 sh deploy.sh
	1. Windows  
	    1. 手动进行构建及部署操作
1. 完成部署
    1. Serverless Framework 会在每个组件部署成功后返回访问方式等信息
