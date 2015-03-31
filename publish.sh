#/bin/bash

set -ex

host="webapp@121.41.85.237"

rsync -av --exclude=.git --exclude=.idea --exclude=static/node_modules --exclude=logs --delete \
 . $host:~/gowork/src/github.com/CardInfoLink/quickpay/

 ssh $host 'cd  ~/gowork/src/github.com/CardInfoLink/quickpay/ && ./run.sh'
