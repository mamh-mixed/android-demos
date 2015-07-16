/* 注意, 单行注释前必须加一个空格，否则解析失败 */
{
    "App": {
        "LogLevel": "debug", // 日志级别
        // 加密存储密钥 44 位，加密密钥由两部分组成，后 22 位在程序中写死，前 22 位在这里配置
        "EncryptKey": "Ky3ETaC0WjgGTQQtFVfl0",
        "HTTPAddr": ":6800", // HTTP 接口
        "TCPPort": ":6600" // 扫码 TCP 接口
    },
    "Mongo": {
        // URL format [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
        "URL": "mongodb://quickpay:quickpay@quick.ipay.so:27017,quick.ipay.so:27018,quick.ipay.so:27019/quickpay",
        "DB": "quickpay"
    },
    "CILOnline": {
        "Host": "140.207.50.238",
        "Port": 7827,
        "QueueSize": 1000, // 所有交易发送到线下前先排队，队列满后将阻塞
        "InitWindowSize": 100, // 全双工 pipeline 方式，接收数据后，找到对应的请求者
        "KeepaliveTime": "65s", // 每隔一段时间发送一个 keepalive 消息
        "ReconnectTimeout": "5s" // 连接断开后过一会儿再重新连接
    },
    "CFCA": {
        "Cert": "config/pem/cfca/cert_testing.pem", // CFCA  验签公钥
        "CCACert": "config/pem/cfca/evCcaCrt_testing.pem", // HTTPS 证书
        "RootCert": "config/pem/cfca/evRootCrt_testing.pem", // CFCA 根证书
        "URL": "https://test.china-clearing.com/Gateway/InterfaceII" // API 地址
    },
    "WeixinScanPay": {
        "ClientCert": "config/pem/weixin/apiclient_cert.pem", // Weixin 客户端证书
        "ClientKey": "config/pem/weixin/apiclient_key.pem", // Weixin 客户端密钥
        "URL": "https://api.mch.weixin.qq.com", // 微信刷卡支付接口地址
        "NotifyURL": "http://test.ipay.so" // 异步消息通知地址，路径是固定的，只需要域名和端口
    },
    "AlipayScanPay": {
        "URL": "https://mapi.alipay.com/gateway.do", // 支付宝扫码支付接口地址
        "NotifyUrl": "http://test.ipay.so" // 支付宝异步消息通知地址，路径是固定的，只需要域名和端口
    }
}
