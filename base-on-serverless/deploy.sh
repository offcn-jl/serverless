#! /bin/bash

# 编译 Golang 二进制文件
sh build.sh

# 部署
if [ "$1" == "prod" ]; then
    sls deploy --all --stage prod
else
    sls deploy --all --stage dev
fi
