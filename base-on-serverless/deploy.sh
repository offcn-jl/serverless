#!/bin/bash -e
set -e # 确保脚本在出错时退出

# 编译 Golang 二进制文件
sh build.sh

# 部署
if [ "$1" == "prod" ]; then
    sls deploy --all --stage prod
else
    sls deploy --all --stage dev
fi
