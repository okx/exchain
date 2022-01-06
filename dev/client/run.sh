#!/bin/bash

gomod() {
  export GOPROXY=http://goproxy.io
  if [ "$1" == '1' ]; then
      export GOPROXY=http://mirrors.aliyun.com/goproxy/
  elif [ "$1" == '2' ]; then
      export GOPROXY=https://athens.azurefd.net
  elif [ "$1" == '3' ]; then
      export GOPROXY=https://gocenter.io
  fi
  export GO111MODULE=on
  go mod tidy
  go mod vendor
}

#gomod

export GO111MODULE=on
go build

TYPE=counter
if [ ! -z "$1" ];
then
  TYPE=$1
fi

./client --type ${TYPE}
