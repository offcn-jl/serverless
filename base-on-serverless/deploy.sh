#! /bin/bash

# 编译 Golang 二进制文件
sh build.sh

# 部署
if [ "$1" == "release" ]; then
    sls deploy --all --stage release
else
    sls deploy --all --stage test
fi
