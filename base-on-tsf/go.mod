module tsf

go 1.12

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/jinzhu/gorm v1.9.12
	github.com/offcn-jl/go-common v0.0.0-20200504092729-1134fe9358be
	github.com/stretchr/testify v1.5.1
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.157+incompatible
	github.com/xluohome/phonedata v0.0.0-20200423024337-2be14779ab82
)

replace github.com/offcn-jl/go-common => ../../go-common // 将 go-common 替换为本地版本, 便于框架的的开发和调试, fork 本项目时应当删除本行 replace 配置
