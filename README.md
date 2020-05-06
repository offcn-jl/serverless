# Serverless APIs
[![Status](https://img.shields.io/badge/Status-Beta-yellow)](#当前版本) [![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/offcn-jl/serverless-apis)](https://goreportcard.com/report/github.com/offcn-jl/serverless-apis) [![Master Build](https://github.com/offcn-jl/serverless-apis/workflows/Master%20Build/badge.svg)](https://github.com/offcn-jl/serverless-apis/actions?query=workflow%3A%22Master+Build%22) [![codecov](https://codecov.io/gh/offcn-jl/serverless-apis/branch/master/graph/badge.svg)](https://codecov.io/gh/offcn-jl/serverless-apis) [![Build](https://github.com/offcn-jl/serverless-apis/workflows/Build/badge.svg)](https://github.com/offcn-jl/serverless-apis/actions?query=workflow%3ABuild) [![codecov](https://codecov.io/gh/offcn-jl/serverless-apis/branch/new-feature/graph/badge.svg)](https://codecov.io/gh/offcn-jl/serverless-apis/branch/new-feature) 

基于无服务器架构的各种 RESTFul 接口

## 食用须知
1. [TSF](https://cloud.tencent.com/document/product/649) 、[SCF](https://cloud.tencent.com/document/product/583) 、 [TKE](https://cloud.tencent.com/document/product/457) 等产品名称、品牌或商标的一切权利归 [腾讯云](https://cloud.tencent.com) 所有。
1. Serverless、Serverless Framework 是 [serverless.com](https://serverless.com) 的产品。
1. 本项目整体基于 MIT 协议开源。
1. 本项目主要包括两个主要分支 : master ( 主分支, 可用于生产环境 )、 new-feature ( 新功能分支, 包含处于测试和验证阶段的新功能 )。
1. 原则上，本项目中的所有模块所使用的各种密码、令牌、口令等敏感信息均配置在环境变量、启动脚本或配置文件中。建议使用者在二次开发的过程中采取同样的方式保存相关设置，并在使用过程中请注意妥善保管相关信息。

Enjoy it. XD

## 当前版本
当前版本 : 测试版

暂定版本发布流程 : Alpha -> Beta -> RC -> GA

> Alpha : 内部测试版, 一般不向外部发布  
> Beta : 也是测试版, 这个阶段的版本会一直加入新的功能  
> RC : 发行候选版本, 基本不再加入新的功能, 主要进行缺陷修复  
> GA : 正式发布的版本, 采用 Release X.Y.Z 作为发布版本号  
> 参考 : [Alpha、Beta、RC、GA版本的区别](http://www.blogjava.net/RomulusW/archive/2008/05/04/197985.html) 、 [软件版本GA、RC、beta等含义](https://blog.csdn.net/gnail_oug/article/details/79998154)

## 目录结构 ( 按文件名排序 )
- marketing-apis  
    - .github ( Github 配置文件 )
    - base-on-tsf  ( 基于 [TSF ( 腾讯微服务平台 )](https://cloud.tencent.com/document/product/649) 的营销接口 )  
    - base-on-serverless ( 基于 [Serverless Framework](https://serverless.com) 及 [腾讯云 Serverless](https://cloud.tencent.com/product/sls) 的营销接口 )  
    - .gitignore ( GIT 配置文件，用于配置需要忽略提交的内容 )  
    - LICENSE ( 版权声明 )  
    - README.md ( 使用说明 )  

## 使用说明
请见各个子目录中的 README.md
