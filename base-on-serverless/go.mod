module serverless

go 1.12

require (
	github.com/jinzhu/gorm v1.9.12
	github.com/offcn-jl/chaos-go-scf v0.0.0-20200422072856-dfa76f029278
)

replace github.com/offcn-jl/chaos-go-scf => ../../chaos-go-scf // 将 chaos-go-scf 框架替换为本地版本, 便于框架的的开发和调试, fork 本项目时应当删除本行 replace 配置
