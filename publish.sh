#/bin/bash

set -ex

prog="quickpay"

function main() {
    # Golang 跨平台编译
    echo "=== Building $prog..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $prog main.go

    host="webapp@121.40.86.222"
    args="-all -port 6800"
    # args="-master -port 6700"
    # args="-pay -port 6800"
    # args="-settle -port 6900"

    workdir="/opt/$prog"

    deploy $host "$args" $workdir

    # host="webapp@121.40.86.222"
    # args="-all -port 68001"
    # workdir="/opt/${prog}2"
    #
    # deploy $host "$args" $workdir

    rm -f $prog
}

function deploy() {
    host=$1
    args=$2
    workdir=$3

    # 上传文件
    echo
    echo "=== Uploading $prog..."
    rsync -rcv --progress $prog $host:$workdir/
    rsync -rcv --progress config/*.ini $host:$workdir/config/
    rsync -rcv --progress static/ $host:$workdir/static/

    # 远程执行重启命令
    echo
    echo "=== SSH $host"
    ssh $host << EOF
export QUICKPAY_ENV=testing

cd $workdir

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
tail -n 30 logs/$prog.log

echo
echo "=== Publish done."
echo
exit

EOF
}

main

exit 0
