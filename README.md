[![Build Status](https://magnum.travis-ci.com/CardInfoLink/quickpay.svg?token=zWvvzH6Ca6HFV3cUQVQD)](https://magnum.travis-ci.com/CardInfoLink/quickpay)


快捷支付
========

快捷支付平台分为 3 个部分：快捷支付、快捷管理、快捷清算。

* __快捷支付__  核心部分，主要对商户提供统一 API 接口，屏蔽多渠道的接口差异，记录绑定和交易信息
* __快捷管理__ 管理、服务、支持快捷支付业务流程正常运转
* __快捷清算__ 对快捷支付的交易定时或准实时清算，并生成报表


安装依赖
-------

安装后端依赖

```
go get github.com/omigo/log
go get github.com/nu7hatch/gouuid
go get gopkg.in/mgo.v2
go get gopkg.in/mgo.v2/bson
go get github.com/axgle/mahonia
```

安装前端依赖

```
# 通过 bower 安装前端依赖
cd static
bower install
cd ..
```


编译安装
-------

```
go install github.com/CardInfoLink/quickpay
```


启动
-------

```
# 查看帮助
$ quickpay
Usage of quickpay:
  -master=false: Startup QuickMaster
  -pay=false: Startup Quickpay
  -port=3800: server listen port, default QuickMaster 3700, Quickpay 3800, QuickSettle 3900
  -settle=false: Startup QuickSettle

# 启动 QuickMaster
$ quickpay -master -port=3700

# 启动 Quickpay
$ quickpay -pay -port=3800

# 启动 QuickSettle
$ quickpay -settle -port=3900
```

免翻墙安装golang.org/x/tools包
-------
```
git clone https://github.com/golang/tools $GOPATH/src/golang.org/x/tools 
cd $GOPATH
go install golang.org/x/tools/cmd/cover
```
