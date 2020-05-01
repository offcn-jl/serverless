#! /bin/bash

WorkSpeace=$(pwd)

# 函数 获取文件列表
function getFileList(){ # https://www.jb51.net/article/142325.htm , https://www.jb51.net/article/48832.htm
    for file in `ls $1` # https://blog.csdn.net/paulkg12/article/details/90549999
    do
        if [ -d $1"/"$file ]; then
            getFileList $1"/"$file
        else
            echo $1"/"$file
        fi
    done
}

# 配置交叉编译
export GOOS=linux
export GOARCH=amd64

echo Build Environments :
echo "GOOS -> "$GOOS
echo "GOARCH -> "$GOARCH

# 开始逐个编译
echo "开始编译."

# 遍历文件列表
for file in `getFileList .`
do
    # 找到含有 serverless 配置文件的目录, 并编译目录中的代码
    if [ "${file##*.}"x = "yml"x ] && [[ $file =~ "scf" ]]; then
        echo "前往目录 "`dirname $file`
        cd `dirname $file` # 跳转到文件所在目录
        echo "开始编译 ..."
        go build -ldflags "-X 'main.builtTime= [ `date +"%Y/%m/%d %H:%M:%S"` ]'" -o artifacts/main # 进行编译
        cd $WorkSpeace
    fi
done

# 结束逐个编译
echo "编译完成."
