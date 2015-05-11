#/bin/bash

set -e

host="webapp@121.40.86.222"
prog="quickpay"
args="-all -port 6800"
# args="-master -port 6700"
# args="-pay -port 6800"
# args="-settle -port 6900"


# Golang 跨平台编译
echo "=== Building $prog..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $prog main.go

# 上传文件
echo
echo "=== Uploading $prog..."
rsync -rcv --progress $prog $host:~/$prog/
rm -f $prog
rsync -rcv --progress config/*.ini $host:~/$prog/config/
rsync -rcv --progress static/ $host:~/$prog/static/


# 远程执行重启命令
echo
echo "=== SSH $host"
ssh $host << EOF
export QUICKPAY_ENV=testing

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
