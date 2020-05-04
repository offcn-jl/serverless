module tsf

go 1.12

require (
	github.com/gin-gonic/gin v1.6.2
	github.com/jinzhu/gorm v1.9.12
	github.com/offcn-jl/go-common v0.0.0-20200429095944-c2b1e7c076b7
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.157+incompatible
)

replace github.com/offcn-jl/go-common => ../../go-common // 将 go-common 替换为本地版本, 便于框架的的开发和调试, fork 本项目时应当删除本行 replace 配置
