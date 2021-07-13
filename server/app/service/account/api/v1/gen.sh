#!/bin/sh

path=$(cd `dirname $0`; pwd)
function read_dir(){
  cd $1
  mkdir .proto.bak
  touch _.micro.go _.pb.go
  mv -f *.micro.go ./.proto.bak
  mv -f *.pb.go ./.proto.bak
  if [ `find ./ -maxdepth 1 -type f -name *.proto | wc -l` -gt 0 ]; then
    protoc -I=. \
      -I=/usr/local/include \
      -I=$GOPATH/src \
      -I=$GOPATH/src/github.com/gogo/protobuf/protobuf \
      -I=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --swagger_out=logtostderr=true:. \
      --gofast_out=plugins=grpc:. \
      *.proto

    if [ $? -ne 0 ]; then
      rm -rf *.micro.go *.pb.go
      mv -f .proto.bak/* ./
      rm -rf .proto.bak _.micro.go _.pb.go
      echo ""
      echo "\033[31m编译失败\033[0m"
      echo ""
      exit 1
    fi
  fi
  rm -rf .proto.bak _.micro.go _.pb.go
  for file in `ls`; do
    if [ -d $1"/"$file ]; then
      read_dir $1"/"$file
    fi
  done
}
read_dir ${path}

echo "编译成功"
exit 0

