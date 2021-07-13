#!/bin/bash

cd `dirname ${0}`
cd ../
root=`pwd`
projectName="${root##*/}"
cd cmd


export GOROOT=$HOME/gosdk/go1.16.4
alias go=$GOROOT/bin/go
shopt -s expand_aliases
basepath=`pwd`
os=$1
arch=$2

if [ "${os}_" = "_" ]; then
  echo "选择你要编译的目标平台?"
  select os in linux darwin windows;do
    break
  done
fi

if [ "${os}_" = "_" ]; then
  os=linux
fi

if [ "${arch}_" = "_" ]; then
  echo "选择内核版本?"
  select arch in amd64 386 arm;do
    break
  done
fi

if [ "${arch}_" = "_" ]; then
  arch=amd64
fi

echo ""
echo -e "当前位置：${basepath}"
echo -e `go version`
echo -e "GOPATH：${GOPATH}"
echo -e "GOROOT：${GOROOT}"
echo -e "你即将编译发布 \033[41m ${projectName}.${os} \033[0m 版"
echo ""
echo -e "\033[31m注意：若配置文件有变动，请先修改服务器上的配置文件！\033[0m"
echo ""
echo "请按回车继续"
read


if [ ! -f version ]; then
  echo "v_1.0.0" > version
fi
version=$(cat version)
increment_version ()
{
  declare -a part=( ${1//\./ } )
  declare    new
  declare -i carry=1
  CNTR=${#part[@]}-1
  new=$((part[CNTR]+carry))
  part[CNTR]=${new}
  new="${part[*]}"
  version="${new// /.}"
}

increment_version ${version}
buildtime=`date +%Y-%m-%d_%H:%M:%S`

CGO_ENABLED=0 GOARCH=${arch} GOOS=${os} go build -a -v -ldflags "-s -w -X main.buildTime=${buildtime} -X main.version=${version}" -o cmd_${os} main.go
if [ $? -ne 0 ]; then
  echo "编译失败"
  exit 1
fi

echo ${version} > version

rm -rf cmd
upx -o cmd cmd_${os}

echo "成功编译 ${os} 版（${version}）"

exit 0
