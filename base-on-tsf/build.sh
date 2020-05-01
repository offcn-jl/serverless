# 交叉 Linux 二进制可执行文件
echo "交叉 Linux 二进制可执行文件..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o artifacts/main

# 打包
echo "打包..."
cd ./artifacts
zip result.zip cmdline main start.sh stop.sh 
