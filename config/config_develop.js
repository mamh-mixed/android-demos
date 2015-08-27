/* 为避免删除非注释文本，配置文件有如下规定:
 * 1. 单行注释文本前必须加一个空字符，与 http://example.com 区分开
 * 2. 多行注释文本前后都必须加一个空格
 */
{
    "App": {
        "LogLevel": "debug", // 日志级别
        // 加密存储密钥 44 位，加密密钥由两部分组成，后 22 位在程序中写死，前 22 位在这里配置
        "EncryptKey": "Ky3ETaC0WjgGTQQtFVfl0",
        "HTTPAddr": ":6800", // HTTP 接口
        "TCPAddr": ":6600", // 扫码 TCP 接口，UTF-8 编码传输，UTF-8 签名
        "TCPGBKAddr": ":6601", // 扫码 TCP 接口，GBK 编码传输，UTF-8 签名
        "DefaultCacheTime": "1s" // 缓存有效时间，0 表示永不过期（慎用），比如 "1h2m3s"
    },
    "Mongo": {
        // URL format [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
        "URL": "mongodb://quickpay:quickpay@dev.ipay.so:27017/quickpay",
        "DB": "quickpay"
    },
    "CILOnline": {
        "Host": "192.168.1.102",
        "Port": 7823,
        "QueueSize": 1000, // 所有交易发送到线下前先排队，队列满后将阻塞
        "InitWindowSize": 100, // 全双工 pipeline 方式，接收数据后，找到对应的请求者
        "KeepaliveTime": "65s", // 每隔一段时间发送一个 keepalive 消息
        "ReconnectTimeout": "5s" // 连接断开后过一会儿再重新连接
    },
    "CFCA": {
        "CFCAPublicKey": "config/pem/cfca/sign_cert_testing.pem", // CFCA  验签公钥
        "CPCNCert": "config/pem/cfca/test.cpcn.com.cn.pem", // HTTPS 证书
        "RootCert": "config/pem/cfca/test.cpcn.com.cn_root.pem", // CFCA 根证书
        "URL": "https://test.cpcn.com.cn/Gateway/InterfaceII" // API 地址
    },
    "WeixinScanPay": {
        "URL": "https://api.mch.weixin.qq.com", // 微信刷卡支付接口地址
        "NotifyURL": "http://dev.quick.ipay.so" // 异步消息通知地址，路径是固定的，只需要域名和端口
    },
    "AlipayScanPay": {
        "AlipayPubKey": "config/pem/alipay/pubkey.pem", // 支付宝 RSA 公钥
        "OpenAPIURL": "https://openapi.alipay.com/gateway.do", // 支付宝 Open API 地址
        "URL": "https://mapi.alipay.com/gateway.do", // 支付宝扫码支付接口地址
        "NotifyUrl": "http://dev.quick.ipay.so", // 支付宝异步消息通知地址，路径是固定的，只需要域名和端口
        "AgentId": "12010128a1" // 标识讯联交易
    }
}
