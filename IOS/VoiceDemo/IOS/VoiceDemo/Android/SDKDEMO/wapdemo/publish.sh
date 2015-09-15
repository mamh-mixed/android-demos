#/bin/bash

set -ex

workdir="/home/weixin/cloudCashier"
host="weixin@139.129.116.65"
function deploy() {
    host=$1
    workdir=$2

    # 上传文件
    echo "=== Uploading $prog..."
    rsync -rcv --progress * --exclude=.DS_Store $host:$workdir/
}

deploy $host $workdir

exit 0
