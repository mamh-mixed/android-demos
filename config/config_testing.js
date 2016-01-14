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
        "DefaultCacheTime": "1s", // 缓存有效时间
        "NotifyURL": "http://test.overseas.ipay.so", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "OrderCloseTime": "24h", // 未支付订单关闭时间
        "OrderRefreshTime":"10m",
        // 平台用户登录后，一段时间内无操作失效时间，比如 30m， 表示 30m 无操作后会话失效，
        // 同时，当会话接近失效时间（1/2）时，如果用户有操作，那么延长失效时间，
        // 比如 用户在第 16 分钟有操作，那么失效时间再延长 30m
        "SessionExpiredTime": "30m",
        "MonitorMerId" : "100001000010110" // 专门做监控的商户，不记日志
    },
    "Qiniu": {
        "Bucket": "test",
        "Domain": "dn-yun-test.qbox.me"
    },
    "Mongo": {
        "Encrypt": false, // URL（密码）是否需要加密，如果需要加密，用 EncryptXXXX 这个字段，避免直接暴露密码，否则用 XXXX 字段
        // URL format [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
        "URL": "mongodb://quickpay:quickpay@test.overseas.ipay.so:27017,test.overseas.ipay.so:27018,test.overseas.ipay.so:27019/quickpay",
        "EncryptURL": "GYeB1+/dUEg41y5AOjDWSCjJM7WyxQOazhHTikDZzQ61SzASrOZYyZ1FDLqc6rD4ByM1uY4bAFJGJglJqiGSPDEl73iBH3v9NbEHMDm61DBE+qp5AzD6tqIPDtCbMdjtGuEzMp5TUis1XHWKDuzhfrHGSKYWKBu13XftXFg7tutqTmdHLGbO1iM8C9Dk1b8zGnZCHLqBTJe9U/hoSQslCf8GZjUC5wJoTfD31KP0gL0kJQ8iLLzM4rLtonnDUwSIilHeHrh6Mx9q6Ia7IJveOIHYfQWu668vJyUNVvMoN9XOEVe3u0WUlxs/NvkkhLgNMKkuK00kZvyIkgoXsNbPZw==",
        "DB": "quickpay"
    },
    "CILOnline": {
        "Host": "140.207.50.238",
        "Port": 7823,
        "ServerCert": "config/pem/cil/server_testing.cert" // TLS 证书
    },
    "CFCA": {
        "URL": "https://test.cpcn.com.cn/Gateway/InterfaceII", // API 地址
        "CheckSignPublicKey": "config/pem/cfca/checkSignPublicKey_testing.pem" // CFCA 验签公钥
    },
    "WeixinScanPay": {
        "URL": "https://api.mch.weixin.qq.com", // 微信刷卡支付接口地址
        "NotifyURL": "http://test.overseas.ipay.so", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "DNSCacheRefreshTime": "5m" // 微信域名解析慢，程序内部做了缓存，这里配置缓存刷新时间
    },
    "AlipayScanPay": {
        "AlipayPubKey": "config/pem/alipay/pubkey.pem", // 支付宝 RSA 公钥
        "OpenAPIURL": "https://openapi.alipay.com/gateway.do", // 支付宝 Open API 地址
        "URL": "https://mapi.alipay.com/gateway.do?_input_charset=utf-8", // 支付宝扫码支付接口地址
        "NotifyUrl": "http://test.overseas.ipay.so", // 支付宝异步消息通知地址，路径是固定的，只需要域名和端口
        "AgentId": "12010128a1" // 标识讯联交易
    },
    "MobileApp":{
        "WXPMerId":"1247075201",
        "ALPMerId":"2088811767473826",
        "WebAppUrl":"http://qrcode.cardinfolink.net/agent"
    },
    "UnionLive": {
        "Encrypt": false, // 密码是否需要加密，如果需要加密，用 EncryptXXXX 这个字段，避免直接暴露密码，否则用 XXXX 字段
        "URL": "http://d.umq.me/PosService/CouponsPurchaseService.ashx", // 优方卡券接口
        "EncryptKey": "tommy123", // 加密密钥
        "EncryptEncryptKey": "CpZQpCXovGXTHHfNILT1rBljRYpgBK4sKNxUOZjCQ5ah3bAYIjQZWk6Vhz7g4uKEtEL+rl2NXMPpF2DKsSOQ4/TLfJQNEcV7+Swnl+ZLzaDPrxDAGxFKdNEQYxiLwxhK78D/EP8qxx0WVUut2MXC0ok6DQiaJNXmQQpkiRGYb7AnGJycifwGFG/4dE++Iqn8nUPgdiGI5xy1VSEkcuOXwZKVRQGIe0wvoi1tzYECICoidL1u+6gqQ2u96821g7P2AYTs2yIMkuwBpdOzIuKvNbKlvEz5iBOPFriu7sl+RirWBtQSn5WSp5VxiYJ6gQzEnh6hkywJhyoSr2Yj10dDMA==", // 加密密钥
        "SignKey": "tommy123", // 签名密钥
        "EncryptSignKey": "IKOxKmCdK86pSgIUUJZFC1wDJBrFxNW7U2x9RuI1tzVi0ubRFqGJ84pAXr5RmoKrPsZeQ+CUOrPlulWrV8BFmyzp9aqu3VuyRfRbONKsgAyz+QIlnRXAQkpNK7raKzH5O+7Arva4TyVY2EaQG6igcRl4Eh79ezpLEaWxtX/EvpXg1jPrIWVBCrYD1qH5kOJHqWHbSiKBO+BKBRy27u0cC70yhWe63sUCoPvsYQg7jl/P3zP8UKn+xz4DFE0JYz81v8ON55F2DKmnOcFj1M2TV4xq5bYy56bQbVrsZxbTLKuw48BOLyM/bUPx2j9zSulhXmKc9eWhi3zFxbdlAdpQPg==", // 签名密钥
        "ChannelId": "182000000000001" // 渠道 ID
    },
    "Settle":{
        "OverseasSettPoint":"", // 海外接口清算时间点
        "DomesticSettPoint":""
    }
}
