#/bin/bash

# set -ex

prog="quickpay"

function main() {
    # Golang 跨平台编译
    echo "=== Building $prog..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $prog main.go

    host="webapp@121.40.86.222"
    args="-all -port=6800"
    # args="-master -port 6700"
    # args="-pay -port 6800"
    # args="-settle -port 6900"

    workdir="/opt/$prog"

    deploy $host "$args" $workdir

    host="webapp@121.40.86.222"
    args="-all -port=6801"
    workdir="/opt/${prog}2"

    # deploy $host "$args" $workdir
    deploy2 $host "$args" $workdir

    rm -f $prog
}

function deploy() {
    host=$1
    args=$2
    workdir=$3

    # 上传文件
    echo "=== Uploading $prog..."
    rsync -rcv --progress $prog $host:$workdir/
    rsync -rcv --progress config/ --exclude=*.go $host:$workdir/config/
    rsync -rcv --progress static/ $host:$workdir/static/

    # 远程执行重启命令
    echo "=== SSH $host"
    ssh $host << EOF
export QUICKPAY_ENV=testing

cd $workdir

echo "=== Killing $prog process..."
ps -ef | grep "$prog $args"
ps -ef | grep "$prog $args" | awk '{print \$2}' | xargs kill -9

echo "=== Starting $prog process ..."
mkdir -p logs
nohup ./$prog $args >> logs/$prog.log 2>&1 &
ps -ef | grep $prog

echo "=== Sleep 3 seconds..."
sleep 2
tail -n 30 logs/$prog.log

echo "=== Publish done."
exit

EOF
}

function deploy2() {
    host=$1
    args=$2
    workdir=$3

    # 远程执行重启命令
    echo "=== SSH $host"
    ssh $host << EOF
export QUICKPAY_ENV=testing

cd $workdir

echo "=== Killing $prog process..."
ps -ef | grep "$prog $args"
ps -ef | grep "$prog $args" | awk '{print \$2}' | xargs kill -9

echo "=== Copying file from another directory..."
cp -rf ../$prog/$prog ../$prog/static ./

echo "=== Starting $prog process ..."
mkdir -p logs
nohup ./$prog $args >> logs/$prog.log 2>&1 &
ps -ef | grep $prog

echo "=== Sleep 3 seconds..."
sleep 2
tail -n 30 logs/$prog.log

echo "=== Publish done."
exit

EOF
}

main

exit 0
