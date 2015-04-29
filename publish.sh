#/bin/bash

set -e

host="webapp@121.40.86.222"
prog="quickpay"
# args="-master -port 3700"
args="-pay -port 3800"
# args="-settle -port 3900"


# Golang 跨平台编译
echo "=== Building $prog..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $prog main.go

# 上传文件
echo
echo "=== Uploading $prog..."
rsync -rcv --progress $prog $host:~/$prog/
rm -f $prog
rsync -rcv --progress static/ $host:~/$prog/static/


# 远程执行重启命令
echo
echo "=== SSH $host"
ssh $host << EOF
export CIL_HOST=140.207.50.238
export CIL_PORT=7826
export MONGO_PORT_27017_TCP_ADDR=121.40.86.222
export MONGO_PORT_27017_TCP_PORT=27017

cd ~/$prog

echo
echo "=== Killing $prog process..."
ps -ef | grep "$prog $args"
ps -ef | grep "$prog $args" | awk '{print \$2}' | xargs kill -9

echo
echo "=== Starting $prog process ..."
mkdir -p logs
nohup ./$prog $args >> logs/$prog.log 2>&1 &
ps -ef | grep $prog

echo
echo "=== Sleep 3 seconds..."
sleep 2
tail -n 10 logs/$prog.log

echo
echo "=== Publish done."
exit

EOF
