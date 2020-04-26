# Serverless APIs
基于无服务器架构的各种 RESTFul 接口

## 食用须知
1. [TSF](https://cloud.tencent.com/document/product/649) 、[SCF](https://cloud.tencent.com/document/product/583) 、 [TKE](https://cloud.tencent.com/document/product/457) 等产品名称、品牌或商标的一切权利归[腾讯云](https://cloud.tencent.com)所有。
1. Serverless、Serverless Framework 是 [serverless.com](https://serverless.com) 的产品。
1. 本项目整体基于 MIT 协议开源。
1. 本项目主要包括两个主要分支 : master ( 主分支, 可用于生产环境 )、 new-feature ( 新功能分支, 包含处于测试和验证阶段的新功能 )。
1. 原则上，本项目中的所有模块所使用的各种密码、令牌、口令等敏感信息均配置在环境变量、启动脚本或配置文件中。建议使用者在二次开发的过程中采取同样的方式保存相关设置，并在使用过程中请注意妥善保管相关信息。

Enjoy it. XD

## 目录结构 ( 按文件名排序 )
|-- marketing-apis  
|	|-- base-on-tsf  // 基于 [TSF ( 腾讯微服务平台 )](https://cloud.tencent.com/document/product/649) 的营销接口  
|	|-- base-on-serverless // 基于 [Serverless Framework](https://serverless.com) 及 [腾讯云 Serverless](https://cloud.tencent.com/product/sls) 的营销接口  
|	|-- .gitignore // GIT 配置文件，用于配置需要忽略提交的内容  
|	|-- LICENSE // 开源协议  
|	|-- README.md // 使用说明  

## 使用说明
请见各个子目录中的 README.md
