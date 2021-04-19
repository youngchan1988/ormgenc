#!/bin/sh
#获取脚本所在的路径
work_path=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
#开始构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags '-w -s' -o $work_path/bin/linux/ormgenc_$1 $work_path/main.go
#upx 压缩
$work_path/tools/upx $work_path/bin/linux/ormgenc_$1
