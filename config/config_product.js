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
        "OrderCloseTime": "24h"
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
        "ServerCert": "config/pem/cil/server_product.cert" // SSL 证书
    },
    "CFCA": {
        "URL": "https://www.china-clearing.com/Gateway/InterfaceII", // API 地址
        "CheckSignPublicKey": "config/pem/cfca/checkSignPublicKey_product.pem" // CFCA 验签公钥
    },
    "WeixinScanPay": {
        "URL": "https://api.mch.weixin.qq.com", // 微信刷卡支付接口地址
        "NotifyURL": "https://api.shou.money", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "DNSCacheRefreshTime": "10m" // 微信域名解析慢，程序内部做了缓存，这里配置缓存刷新时间
    },
    "AlipayScanPay": {
        "AlipayPubKey": "config/pem/alipay/pubkey.pem", // 支付宝 RSA 公钥
        "OpenAPIURL": "https://openapi.alipay.com/gateway.do", // 支付宝 Open API 地址
        "URL": "https://mapi.alipay.com/gateway.do", // 支付宝扫码支付接口地址
        "NotifyUrl": "https://api.shou.money", // 异步消息通知地址，路径是固定的，只需要域名和端口
        "AgentId": "12010128a1" // 标识讯联交易
    }
}
