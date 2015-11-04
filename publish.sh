#/bin/bash

set -e

prog="quickpay"

shortcut=("dev" "test" "app1" "app2")
envs=("develop" "testing" "product" "product")
hosts=("webapp@dev.ipay.so" "webapp@test.ipay.so" \
    "quick@app1.set.shou.money" "quick@app2.set.shou.money")

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
    goBuild $prog ${envs[$idx]}
    echo

    version=$(git rev-parse --verify HEAD --short)
    if [ ${envs[$idx]} == "product" ]; then
        version=$(git describe --abbrev=0 --tags)
    fi

    # 前端打包
    echo ">>> Use Gulp to package frontend html/js/css..."
    gulpPackage $version
    echo

    workdir="/opt/$prog"

    # 发布
    echo ">>> Rsync executable program..."
    deploy $host $workdir $version
    echo
    # 重启
    echo ">>> Restart program and tail log..."
    restart $host $workdir $version
    echo

    rm -rf distrib/
    echo ">>> Publish done."
}

function goBuild() {
    prog=$1
    env=$2
    echo "Running go generate..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate
    echo "Building $prog..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o distrib/"$prog" main.go
    cp -r config distrib/
    mkdir -p distrib/app/material
    cp -r app/material/ distrib/app/material
    if [ "$env" == "develop" ]; then
        rm distrib/config/*testing*
        rm distrib/config/*product*
    elif [ "$env" == "testing" ]; then
        rm distrib/config/*develop*
        rm distrib/config/*product*
    elif [ "$env" == "product" ]; then
        rm distrib/config/*testing*
        rm distrib/config/*develop*
    fi
}

function gulpPackage() {
    version=$1

    echo "frontend version $version"

    mkdir -p distrib/static/$version
    cd static
    bower install # 安装前端依赖
    gulp # 压缩文件
    cd ..
    cp static/index.html distrib/static/index.html # 跳转页面
    sed -i '' 's/url=v0.0.1"/url='$version'"/g' distrib/static/index.html # 替换版本号
    cp -r static/dist/ distrib/static/$version/ # 新版本静态文件路径
}

function deploy() {
    host=$1
    workdir=$2
    version=$3

    # 远程执行复制，以使用 rsync 加速传输
    echo "SSH $host"
    ssh $host << EOF
cd $workdir
prev=\$(ls -t static | grep -v 'index.html' | head -n  1)
if [ "\$prev" != "$version" ]; then
    cp -r static/\$prev static/$version
fi
exit
EOF

    # 上传文件
    echo "Uploading $prog..."
    rsync -rcv --exclude=logs --exclude=.DS_Store \
        --progress distrib/ $host:$workdir/
}

function restart() {
    host=$1
    workdir=$2
    version=$3

    # 远程执行重启命令
    echo "SSH $host"
    ssh $host << EOF
cd $workdir
# 备份
cp $prog ${prog}_$version
cp -r config config_$version
cp -r app app_$version

echo "Killing $prog process..."
ps -ef | grep $prog
killall $prog

echo "Starting $prog process ..."
mkdir -p logs
nohup ./$prog >> logs/$prog.log 2>&1 &
ps -ef | grep $prog

echo "Sleep 5 seconds..."
sleep 5
tail -n 30 logs/$prog.log

EOF
}

main

exit 0
