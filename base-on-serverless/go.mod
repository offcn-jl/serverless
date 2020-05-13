module serverless

go 1.12

require (
	github.com/jinzhu/gorm v1.9.12
	github.com/offcn-jl/cscf v0.0.0-20200511103748-5fcad49f70a7
	github.com/offcn-jl/go-common v0.0.0-20200504092729-1134fe9358be
	github.com/offcn-jl/gscf v0.0.0
	github.com/stretchr/testify v1.5.1
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.168+incompatible
	github.com/tencentyun/scf-go-lib v0.0.0-20200116145541-9a6ea1bf75b8
	github.com/xluohome/phonedata v0.0.0-20200423024337-2be14779ab82
	gopkg.in/go-playground/assert.v1 v1.2.1
)

replace github.com/offcn-jl/gscf => ../../gscf // 将 gscf 框架替换为本地版本, 便于框架的的开发和调试, fork 本项目时应当删除本行 replace 配置

replace github.com/offcn-jl/go-common => ../../go-common // 将 go-common 替换为本地版本, 便于框架的的开发和调试, fork 本项目时应当删除本行 replace 配置
