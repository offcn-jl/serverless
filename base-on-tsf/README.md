# 基于 TSF 的接口
使用 [TSF](https://cloud.tencent.com/document/product/649) 作为计算资源的接口, 也可以直接将编译后的二进制主程序在任何平台上运行后提供服务。  
还可以使用 [CSCF 框架](https://github.com/offcn-jl/cscf) 将 handler 包中的处理函数迁移到 [SCF](https://cloud.tencent.com/document/product/583) 平台运行。

## 食用须知
1. [TSF](https://cloud.tencent.com/document/product/649) 、 [SCF](https://cloud.tencent.com/document/product/583) 、 [TKE](https://cloud.tencent.com/document/product/457) 等产品名称、品牌或商标的一切权利归 [腾讯云](https://cloud.tencent.com) 所有。
1. Serverless、Serverless Framework 是 [serverless.com](https://serverless.com) 的产品。
1. 本项目整体基于 MIT 协议开源。
1. 本项目主要包括两个主要分支 : master ( 主分支, 可用于生产环境 )、 new-feature ( 新功能分支, 包含处于测试和验证阶段的新功能 )。
1. 原则上，本项目中的所有模块所使用的各种密码、令牌、口令等敏感信息均配置在环境变量、启动脚本或配置文件中。建议使用者在二次开发的过程中采取同样的方式保存相关设置，并在使用过程中请注意妥善保管相关信息。

Enjoy it. XD

## 接口列表

1. 照片处理 ( 证件照生成 ) [ [腾讯云人像分割](https://cloud.tencent.com/document/product/1208/42970) & [腾讯云人脸美颜](https://cloud.tencent.com/document/product/1172/40715) ]
1. 带计数的预约类活动接口

## 目录结构 ( 按文件名排序 )

- base-on-tsf ( 基于 [TSF ( 腾讯微服务平台 )](https://cloud.tencent.com/document/product/649) 的接口 )  
    - artifacts ( 用于最终部署的制品 )  
        - cmdline ( 用于检查应用进程是否存在，没有 .sh 后缀 )  
        - main ( 主程序, 执行构建步骤后生成 )  
        - result.zip ( 程序包，用于最终部署, 执行构建步骤后生成 )  
        - start.release.sh.example ( 生产环境示例启动脚本，需要自行修改为启动脚本 )  
        - start.test.sh.example ( 测试环境示例启动脚本，需要自行修改为启动脚本 )  
        - start.sh ( 启动脚本 )  
        - stop.sh ( 停止脚本 )  
    - common ( 组件库 )  
        - config  
            - config.go ( 配置包 )  
        - database ( 数据库包 )
            - orm ( ORM 包 )  
                - structs ( 结构体包 )  
                    - app.go ( 应用程序相关数据结构 )  
                    - e2020.go ( 2020 年的各类活动使用的数据结构 )  
                    - sso_v2.go ( 第二代单点登录系统相关数据结构 )  
                    - structs.go ( 基础数据结构及工具 )  
                    - structs_test.go ( 基础数据结构及工具单元测试 )  
        - handler ( 处理函数 )  
            - 2020 ( 2020 年各类活动接口的处理函数包 )  
                - 05 ( 5 月各类活动接口的处理函数包 )  
                    - 01 ( 本月第一个活动的处理函数包 )  
                        - h20200501.go ( 带计数的预约类活动接口的处理函数包 )  
                    - 02 ( 本月第二个活动的处理函数包 )  
                        - h20200502.go ( 礼品 ( 兑换码 ) 发放类活动接口的处理函数包 )  
                    - 03 ( 本月第三个活动的处理函数包 )  
                        - h20200503.go ( 需要进行参与次数计数的活动接口的处理函数包 )  
                        - h20200503_test.go ( 需要进行参与次数计数的活动接口的处理函数包的单元测试 )  
            - app ( 应用程序接口的处理函数包 )  
                - version.go ( 版本相关接口的处理函数包 )  
            - photo ( 照片处理接口的处理函数包 )  
                - photo.go ( 照片处理接口的处理函数包 )  
            - sso ( 单点登录系统接口的处理函数包 )  
                - v2 ( 第二代单点登录系统接口的处理函数包 )  
                    - session.go ( 会话相关接口的处理函数包 )  
                    - session_test.go ( 会话相关接口的处理函数包的单元测试 )  
                    - suffix.go ( 后缀相关接口的处理函数包 )  
                    - suffix_test.go ( 后缀相关接口的处理函数包的单元测试 )  
        - middleware ( 中间件包 )  
            - cors.go ( 处理跨域资源共享的中间件包 )  
            - version.go ( 添加版本信息的中间件包 )  
        - utils ( 工具包 )  
            - utils.go ( 工具包 )  
    - front-end-example ( 前端演示页面 )  
        - 2020 ( 2020 年各类活动的演示页面 )  
            - 05 ( 5 月各类活动的演示页面 )  
                - 01 ( 带计数的预约类活动接口的演示页面 **!!可能发生非预期的超卖!!** )  
                    - 01 ( 显示预约人数或剩余数量, 可显示预约序号及预约 ID, 可限制总数或主动超卖 )  
                    - 02 ( 配合登陆模块，按分部显示预约人数或剩余数量, 可显示预约序号及预约 ID, 可限制总数或主动超卖 )  
                    - public.js ( 公共 JS )  
- .gitignore ( GIT 配置文件，用于配置需要忽略提交的内容 )  
- build.sh ( 构建脚本 )  
- go.mod ( go mod 配置文件 )  
- go.sum ( go mod checksums )  
- main.go ( 主程序源码 )  
- README.md ( 使用说明 )  

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
	1. 复制启动脚本模板 ( start.release.sh.example 及  start.test.sh.example )
	1. 更名为可用的启动脚本 ( start.release.sh 及 start.test.sh )
	1. 打开测试环境启动脚本 ( start.test.sh ) ，将其中各个字段后的 your_xxx 替换为对应的数据
	1. 打开生产环境启动脚本 ( start.release.sh )，将其中各个字段后的 your_xxx 替换为对应的数据
1. 构建
	1. MacOS / Linux
		1. 直接在本目录执行 sh build.sh
	1. Windows
		1. 交叉编译 Linux 64 位二进制可执行程序，在本目录执行 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o artifacts/main
		1. 根据要部署的环境, 将对应环境的启动脚本复制并重命名为 start.sh
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
