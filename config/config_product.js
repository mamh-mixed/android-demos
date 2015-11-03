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
        "DefaultCacheTime": "5m", // 缓存有效时间
        "NotifyURL": "https://api.shou.money", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "OrderCloseTime": "24h", // 未支付订单关闭时间
        "OrderRefreshTime": "10m",
        // 平台用户登录后，一段时间内无操作失效时间，比如 30m， 表示 30m 无操作后会话失效，
        // 同时，当会话接近失效时间（1/2）时，如果用户有操作，那么延长失效时间，
        // 比如 用户在第 16 分钟有操作，那么失效时间再延长 30m
        "SessionExpiredTime": "30m",
        "MonitorMerId" : "100001000010110" // 专门做监控的商户，不记日志
    },
    "Qiniu": {
        "Bucket": "yun-master",
        "Domain": "dn-yun-master.qbox.me"
    },
    "Mongo": {
        "Encrypt": true, // URL（密码）是否需要加密，如果需要加密，用 EncryptXXXX 这个字段，避免直接暴露密码，否则用 XXXX 字段
        // URL format [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
        "URL": "mongodb://quickpay:******@mgo2.set.shou.money:27017,mgo1.set.shou.money:27018,nsq1.set.shou.money:27017/quickpay",
        "EncryptURL": "gVHAPh0W/zta5rtNsgiqFSUA4kexR0s2wfyu7bIocTGGZntW7oLlqILV2OFJAX3YQAWywM6JBZNEucWXERINsMQm1OoyXcukRwTxyl+i7B/aEWlA6mumm2iz1zdI+R/hYO5K/mJ8T4fp9qnorVkQixsUvXQuOs2+8S0He/2V8QmsnIjouZF8h7XvFM21yuelmPwKw1zoGRyvUNt2YjA9jAQ5J+YPzxJWkfyLBjoCaKWxVM4DgIOtBGMZQVYSedAWtCsAlrRw6GAEa4TD74vkpvKSX9LmGxtt9SeTj04cNpViqZDmPS4c9Jrl3J/G35cICfI3meSvq78XYjmYjnSohA==",
        "DB": "quickpay"
    },
    "CILOnline": {
        "Host": "211.147.72.70",
        "Port": 10010,
        "ServerCert": "config/pem/cil/server_product.cert" // TLS 证书
    },
    "CFCA": {
        "URL": "https://www.china-clearing.com/Gateway/InterfaceII", // API 地址
        "CheckSignPublicKey": "config/pem/cfca/checkSignPublicKey_product.pem" // CFCA 验签公钥
    },
    "WeixinScanPay": {
        "URL": "https://api.mch.weixin.qq.com", // 微信刷卡支付接口地址
        "NotifyURL": "https://api.shou.money", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "DNSCacheRefreshTime": "5m" // 微信域名解析慢，程序内部做了缓存，这里配置缓存刷新时间
    },
    "AlipayScanPay": {
        "AlipayPubKey": "config/pem/alipay/pubkey.pem", // 支付宝 RSA 公钥
        "OpenAPIURL": "https://openapi.alipay.com/gateway.do", // 支付宝 Open API 地址
        "URL": "https://mapi.alipay.com/gateway.do", // 支付宝扫码支付接口地址
        "NotifyUrl": "https://api.shou.money", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "AgentId": "12010128a1" // 标识讯联交易
    },
    "MobileApp": {
        "WXPMerId": "1239305502",
        "ALPMerId": "2088811767473826",
        "WebAppUrl": "http://qrcode.cardinfolink.net/payment"
    },
    "UnionLive": {
        "Encrypt": true, // 密码是否需要加密，如果需要加密，用 EncryptXXXX 这个字段，避免直接暴露密码，否则用 XXXX 字段
        "URL": "http://umq.cc/posservice/CouponsPurchaseService.ashx", // 优方卡券接口
        "EncryptKey": "BA823FCE", // 加密密钥
        "EncryptEncryptKey": "RkRcF8ofgwV8G5stnhtxZa7jebuAL8lOgmjxeNosgqvzVatwPxMkh+yPGD4uPrmuLXg/LeJ3kAb65AeluOltbbGhdafIA8rUftiKlq1ULmE16lcv8NrXZ/TtONcHBDwud5Ggzxo0aBND8qy3kfGi2DQSa8eFZPG31BjTDcgh0B1euV/X+pEzIW1z8+x3Izmf2TJzdUGqiwUFog437k47+zdVhtgyhWMX3m9GGpQx0CSEtJQtzxZfutNNONzJu+eVstr3en7RWMplTzxtWRZO+uy2W3LNEa1r1pmcTPbQEhr6YsjBgpP+TJ7RNaSq9fojTVAkiy9aacwq/+JvZJAKgg==", // 加密密钥
        "SignKey": "291BFD0289B24317BF9A520CA6DFDAD2", // 签名密钥
        "EncryptSignKey": "bNw1qCugpH6J7xILeOLRiHrBJRrEceLmJa0wZ9xA7gcvc7pdiCnH/1OaO2kfTKrZ0L9PtGjj3TilkBvBz7w2x/PkhyhwhGxnxpNzZSBpUNo84gskSRaJSd4Xv/d6VuyMzezUVnmsK2DtqCdcE/90Li22aOyvo5r8yJR8YRcLsEOpR4mh1su7MNkqG6zw3o4R+SEI9htqGCEugJ9b0VxE0es+o3zn0Kb5yr6zDN8fIgv5n7HRTj2ogLGF5SSMYt3bLlAaqNeoseZLursMYynd+XJo+uNq0EYDJfVzhoWtgq/rBUMHvGoIqihNvsrkM/ZJ95ZkgCEBWPEq6FTyDJN1Cg==", // 签名密钥
        "ChannelId": "180210000270000" // 渠道 ID
    }
}
