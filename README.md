[![Build Status](https://magnum.travis-ci.com/CardInfoLink/quickpay.svg?token=zWvvzH6Ca6HFV3cUQVQD)](https://magnum.travis-ci.com/CardInfoLink/quickpay)


快捷支付
========

快捷支付平台分为 3 个部分：快捷支付、快捷管理、快捷清算。

* __快捷支付__ 核心部分，主要对商户提供统一 API 接口，屏蔽多渠道的接口差异，记录绑定和交易信息
* __快捷管理__ 管理、服务、支持快捷支付业务流程正常运转
* __快捷清算__ 对快捷支付的交易定时或准实时清算，并生成报表


安装依赖
-------

* 安装后端依赖

```bash
go get ./...
```

* 安装前端依赖

```bash
# 通过 bower 安装前端依赖
# polymer 0.5 绑定支付
cd static && bower install && cd ..
# polymer 1.0 扫码支付
cd admin && bower install && cd ..
```


编译安装
-------

```bash
go install github.com/CardInfoLink/quickpay
```


启动
----

1. 启动前，需要在系统中配置一个环境变量，表明是开发环境、测试环境、还是生产环境。把如下配置加入到
`~/.bashrc` 或 `~/.profile` 或 `~/.bash_profile` 中：

```bash
# 配置快捷支付环境变量，QUICKPAY_ENV 的值只能是 develop、testing 或 product 中的一个
export QUICKPAY_ENV=develop
```

2. 修改对应环境的配置文件，配置文件应该放在程序的启动目录下，结构如下：

```
config/
├── config_develop.js
├── config_product.js
├── config_testing.js
└── pem
    ├── alipay
    │   └── pubkey.pem
    ├── cfca
    │   ├── cert_testing.pem
    │   ├── evCcaCrt_testing.pem
    │   └── evRootCrt_testing.pem
    └── weixin
        ├── apiclient_cert.pem
        └── apiclient_key.pem
```

3. 在程序的启动目录下创建 logs 目录，用于存放日志文件

4. 启动命令 `nohup quickpay >> logs/quickpay.log 2>&1 &`

5. 查看日志 `tail -f logs/quickpay.log`
