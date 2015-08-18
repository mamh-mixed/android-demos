#/bin/bash

set -ex

prog="quickpay"

function main() {
    # Golang 跨平台编译
    echo "=== Building $prog..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o distrib/"$prog" main.go

    workdir="/opt/$prog"

    mkdir -p distrib/static/app
    cp -r config distrib/
    cd static
    bower install
    cp -r app/fonts/font-roboto bower_components/
    gulp
    cd ..
    cp -r static/dist/ distrib/static/app/
    cd admin && bower install && cd ..
    cp -r admin distrib/

    #host="webapp@dev.ipay.so"
    host="webapp@test.ipay.so"

    # host="quick@app1.set.shou.money"

    # sed -i '' 's/port=4160/port=4161/' distrib/config/config_product.ini
    # host="quick@app2.set.shou.money"

    deploy $host $workdir
    rm -rf distrib/
}

function deploy() {
    host=$1
    workdir=$2

    # 上传文件
    echo "=== Uploading $prog..."
    rsync -rcv --exclude=logs --exclude=.DS_Store \
        --delete --progress distrib/ $host:$workdir/

    # 远程执行重启命令
    echo "=== SSH $host"
    ssh $host << EOF

cd $workdir

echo "=== Killing $prog process..."
ps -ef | grep "$prog"
ps -ef | grep "$prog" | awk '{print \$2}' | xargs kill -9
pwd
echo "=== Starting $prog process ..."
mkdir -p logs
nohup ./$prog >> logs/$prog.log 2>&1 &
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
