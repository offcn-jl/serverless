#!/bin/bash
set -e # 确保脚本在出错时退出

# 交叉 Linux 二进制可执行文件
echo "交叉编译 Linux 二进制可执行文件..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'tsf/common/config.builtTime= [ `date +"%Y/%m/%d %H:%M:%S"` ]'" -o artifacts/main

# 配置启动脚本
if [ "$1" == "release" ]; then
    cp artifacts/start.release.sh artifacts/start.sh
else
    cp artifacts/start.test.sh artifacts/start.sh
fi

# 打包
echo "打包..."
cd ./artifacts
zip result.zip cmdline main start.sh stop.sh
