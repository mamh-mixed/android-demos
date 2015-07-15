#/bin/bash

set -ex

prog="quickpay"

function main() {
    # Golang 跨平台编译
    echo "=== Building $prog..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o distrib/"$prog" main.go

    workdir="/opt/$prog"

    cp -r config distrib/
    cd static && bower install && cd ..
    cp -r admin distrib/
    cp -r static distrib/

    # host="webapp@dev.ipay.so"
    host="webapp@test.ipay.so"

    # host="quick@app1.set.shou.money"

    # sed -i '' 's/port=4160/port=4161/' distrib/config/config_product.ini
    # host="quick@app2.set.shou.money"

    deploy $host $workdir
}

function deploy() {
    host=$1
    workdir=$2

    # 上传文件
    echo "=== Uploading $prog..."
    rsync -rcv --progress distrib/ --exclude=.DS_Store $host:$workdir/

    # 远程执行重启命令
    echo "=== SSH $host"
    ssh $host << EOF

cd $workdir

echo "=== Killing $workdir/$prog process..."
ps -ef | grep "$workdir/$prog"
ps -ef | grep "$workdir/$prog" | awk '{print \$2}' | xargs kill -9
pwd
echo "=== Starting $prog process ..."
mkdir -p logs
nohup $workdir/$prog >> $workdir/logs/$prog.log 2>&1 &
ps -ef | grep $workdir/$prog

echo "=== Sleep 3 seconds..."
sleep 2
tail -n 30 logs/$prog.log

echo "=== Publish done."
exit

EOF
}

main

exit 0
