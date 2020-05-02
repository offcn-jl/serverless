module serverless

go 1.12

require (
	github.com/jinzhu/gorm v1.9.12
	github.com/offcn-jl/cscf v0.0.0-20200502025411-fd0546ea314c
	github.com/offcn-jl/go-common v0.0.0-20200429024900-a9522f28c623
)

replace github.com/offcn-jl/cscf => ../../cscf // 将 cscf 框架替换为本地版本, 便于框架的的开发和调试, fork 本项目时应当删除本行 replace 配置

replace github.com/offcn-jl/go-common => ../../go-common // 将 go-common 替换为本地版本, 便于框架的的开发和调试, fork 本项目时应当删除本行 replace 配置
