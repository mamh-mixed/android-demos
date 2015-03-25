#!/bin/bash

set -x

. /etc/profile.d/go.sh
export GOPATH="/home/webapp/gowork"

go build main.go

mkdir -p logs

killall main
pwd
nohup ./main > logs/quickpay.log 2>&1 &

sleep 3

tail -f logs/quickpay.log
