#/bin/bash

set -e

prog="quickpay"

shortcut=("dev" "test" "app1" "app2")
hosts=("webapp@dev.ipay.so" "webapp@test.ipay.so" \
    "quick@app1.set.shou.money" "quick@app2.set.shou.money" )

input=$1

function main() {
    echo "Host list:"
    for i in ${!shortcut[@]}; do
        printf "  %-4s): %s\n" ${shortcut[$i]} ${hosts[$i]}
    done
    echo
    if [ -z "$input" ]; then
        read -t 30 -p "Select host: ${shortcut[*]}? [dev] " input
        if [ -z "$input" ]; then
            input="dev"
        fi
    fi

    idx=""
    for i in ${!shortcut[@]}; do
        if [ "$input" == "${shortcut[$i]}" ]; then
            idx=$i
            break
        fi
    done
    if [ -z "$idx" ]; then
        echo "Error host selected!"
        exit 1
    fi

    host=${hosts[$idx]}
    echo -e "You selected host:\033[1;31m $input => $host \033[0m"
    if [ "$input" != "dev" ]; then
        read -t 30 -p "Press any key to continue..."
    fi
    echo

    # Golang 跨平台编译
    echo ">>> Compile backend golang code..."
    goBuild $prog
    echo

    # 前端打包
    echo ">>> Use Gulp to package frontend html/js/css..."
    gulpPackage
    echo

    workdir="/opt/$prog"

    if [ "$input" == "app2" ]; then
        sed -i '' 's/port=4160/port=4161/' distrib/config/config_product.ini
    fi

    # 发布
    echo ">>> Rsync executable program..."
    deploy $host $workdir
    echo
    # 重启
    echo ">>> Restart program and tail log..."
    restart $host $workdir
    echo

    rm -rf distrib/
    echo ">>> Publish done."
}

function goBuild() {
    prog=$1
    echo "Building $prog..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o distrib/"$prog" main.go
    cp -r config distrib/
}

function gulpPackage() {
    mkdir -p distrib/static/app
    cd static
    bower install
    cp -r app/fonts/font-roboto bower_components/
    gulp
    cd ..
    cp -r static/dist/ distrib/static/app/
}

function deploy() {
    host=$1
    workdir=$2

    # 上传文件
    echo "Uploading $prog..."
    rsync -rcv --exclude=logs --exclude=.DS_Store \
        --delete --progress distrib/ $host:$workdir/
}

function restart() {
    host=$1
    workdir=$2

    # 远程执行重启命令
    echo "SSH $host"
    ssh $host << EOF

cd $workdir

echo "Killing $prog process..."
ps -ef | grep "$prog"
ps -ef | grep "$prog" | awk '{print \$2}' | xargs kill -9
pwd
echo "Starting $prog process ..."
mkdir -p logs
nohup ./$prog >> logs/$prog.log 2>&1 &
ps -ef | grep $prog

echo "Sleep 5 seconds..."
sleep 5
tail -n 60 logs/$prog.log

exit
EOF
}

main

exit 0
