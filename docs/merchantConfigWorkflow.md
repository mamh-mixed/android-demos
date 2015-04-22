商户配置流程
-------

假定有一个商户，商户号为`APPTEST`

###配置密钥信息
在`merchant`集合中配置商户的密钥信息：

```json
{
    "merId" : "APPTEST",
    "merStatus" : "Normal",
    "transCurr" : "156",
    "signKey" : "0123456789",
    "encryptKey" : "AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8="
}
```

###配置路由
在`routerPolicy`集合中配置路由信息
```json
{
    "merId" : "APPTEST",
    "cardBrand" : "MCC",
    "chanCode" : "APT",
    "chanMerId" : "050310058120002"
}
```

说明：`chanCode`是渠道商户信息`chanMer`的主键

###渠道商户
```json
{
    "chanCode" : "APT",
    "chanMerId" : "050310058120002",
    "chanMerName" : "Apple Pay测试渠道商户",
    "terminalId" : "00000001",
    "insCode" : "00000050"
}
```
