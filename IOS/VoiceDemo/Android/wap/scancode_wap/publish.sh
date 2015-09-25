#/bin/bash

set -ex

#/ 生产     workdir="/home/weixin/cloudCashier/payment"
#/ 测试     workdir="/home/weixin/cloudCashier/agent"

workdir="/home/weixin/cloudCashier/agent"

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
