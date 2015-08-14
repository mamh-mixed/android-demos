/* 为避免删除非注释文本，配置文件有如下规定:
 * 1. 单行注释文本前必须加一个空字符，与 http://example.com 区分开
 * 2. 多行注释文本前后都必须加一个空格
 */
{
    "App": {
        "LogLevel": "info", // 日志级别
        // 加密存储密钥 44 位，加密密钥由两部分组成，后 22 位在程序中写死，前 22 位在这里配置
        "EncryptKey": "Ky3ETaC0WjgGTQQtFVfl9",
        "HTTPAddr": ":6800", // HTTP 接口
        "TCPAddr": ":6600", // 扫码 TCP 接口，UTF-8 编码传输，UTF-8 签名
        "TCPGBKAddr": ":6601", // 扫码 TCP 接口，GBK 编码传输，UTF-8 签名
        "DefaultCacheTime": "5m" // 缓存有效时间
    },
    "Mongo": {
        // URL format [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
        "URL": "mongodb://quickpay:dwLDq1OyiZVgY40Mt0SdOrADckJVJyyeNK1BuV1D0Iw@mgo2.set.shou.money:27017,mgo1.set.shou.money:27018,nsq1.set.shou.money:27017/quickpay",
        "DB": "quickpay"
    },
    "CILOnline": {
        "Host": "172.16.30.11",
        "Port": 4160,
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
        "URL": "https://api.mch.weixin.qq.com", // 微信刷卡支付接口地址
        "NotifyURL": "https://api.shou.money" // 异步消息通知地址，路径是固定的，只需要域名和端口
    },
    "AlipayScanPay": {
        "AlipayPubKey": "config/pem/alipay/pubkey.pem", // 支付宝 RSA 公钥
        "OpenAPIURL": "https://openapi.alipay.com/gateway.do", // 支付宝 Open API 地址
        "URL": "https://mapi.alipay.com/gateway.do", // 支付宝扫码支付接口地址
        "NotifyUrl": "https://api.shou.money", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "AgentId": "12010128a1" // 标识讯联交易
    }
}
