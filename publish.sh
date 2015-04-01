#/bin/bash

set -e

host="webapp@121.41.85.237"
prog="quickpay"

### 这个脚本是通用，下面无需改动 ###

# Golang 跨平台编译
echo "=== Building $prog..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $prog main.go

# 上传文件
echo
echo "=== Uploading $prog..."
rsync -v --progress quickpay $host:~/$prog/

# 远程执行重启命令
echo
echo "=== SSH $host"
ssh $host << EOF

cd ~/$prog

echo
echo "=== Killing $prog process..."
ps -ef | grep $prog
killall $prog

echo
echo "=== Starting $prog process ..."
mkdir -p logs
nohup ./$prog >> logs/$prog.log 2>&1 &
ps -ef | grep $prog

echo
echo "=== Sleep 3 seconds..."
sleep 3
tail -n 10 logs/$prog.log

echo
echo "=== Publish done."
exit

EOF
