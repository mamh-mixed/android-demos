
# 名字解释

| 名词        | 缩写   | 解释         |
|------------|--------|-------------|
| merchant   | mer    | 商户         |
| channel    | chan   | 渠道         |
| CFCA       | cfca   | 中金         |
| settlement | sett   | 清分         |
| cil-online | cil    | 线下联机系统  |
| alipay     | alipay | 支付宝       |
| weixin     | weixin | 微信支付     |


# 数据模式


* 商户基本信息	merchant

| 名称                 | 字段           | 备注                             |
| ------------------- | -------------  | -------------                    |
| 商户号               | merId          |                                  |
| 商户状态             | merStatus      | （Normal，Deleted）               |
| 商户交易币种          | transCurr      |                                  |
| 商户签名密钥          | signKey        |                                  |
| 商户加密密钥          | encryptKey     |                                  |
| 备注信息          | remark     |                                  |


* 商户详细信息 merDetail

| 名称              | 字段             | 备注                           |
| --------------- | -------------  | -------------                    |
| 商户号             | merId          |                                 |
| 商户名称            | merName        |                                  |
| 商户简称            | shortName      |                                  |
| 商户城市            | city           |                                  |
| 商户国家            | nation         |                                  |
| 商户类型            | merType        |                                  |
| 商户计费方案代码        | billingScheme  |                                  |
| 商户清算币种          | SettCurr       |                                  |
| 商户账户名称          | acctName       |                                  |
| 商户账户            | acctNum        |                                  |
| 法人代表            | corp           |                                  |
| 商户负责人           | master         |                                  |
| 商户联系人           | contact        |                                  |
| 商户联系电话          | contactTel     |                                  |
| 商户传真            | fax            |                                  |
| 商户邮箱            | email          |                                  |
| 商户地址            | addr           |                                  |
| 商户邮编            | postcode       |                                  |
| 商户密码            | password       |                                  |


* 渠道商户信息：PK（渠道代码，商户号）	chanMer

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 渠道代码            | chanCode       |                                  |
| 渠道商户号           | chanMerId      |                                  |
| 渠道商户名称          | chanMerName    |                                  |
| 清算角色            | settRole       |                                  |
| 结算标识（中金渠道商户特有）  | settFlag       |                                  |
| 签名证书            | signCert       |                                  |
| 验签证书            | checkSignCert  |                                  |
| 支付宝 MD5Key      | alpMd5Key      |                                  |
| 微信支付AppId       | wxpAppId       |                                  |
| 微信支付PartnerKey  | wxpPartnerKey  |                                  |
| 微信支付加密证书        | wxpEncryptCert |                                  |
| 机构号             | insCode        | (Apple Pay用到，对应到线下网关的chcd)       |
| 终端号             | terminalId     | (Apple Pay用到，对应到线下网关的terminalid) |


* 应答码	respCode

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 应答码             | respCode       |                                  |
| 应答信息            | respMsg        |                                  |
| CFCA            | cfca           |                                  |
| CFCA 应答码        | cfca.code      |                                |
| CFCA 应答信息       | cfca.msg       |                                |


* 卡Bin	cardBin

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 卡Bin            | bin            |                                  |
| 卡Bin长度          | binLen         |                                  |
| 卡号长度            | cardLen        |                                  |
| 卡品牌             | cardBrand      |                                  |
| 发卡行代码           | issBankCode    |                                  |
| 发卡行名称           | issBankName    |                                  |


* 路由策略	routerPolicy

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 源商户号            | merId          |                                  |
| 卡品牌             | cardBrand      |                                  |
| 卡类型             | cardType       |                                  |
| 交易类型            | transType      |                                  |
| 卡Bin组           | binGroup       |                                  |
| 输入方式            | inputWay       |                                  |
| 起始金额            | minAmount      |                                  |
| 最大金额            | maxAmount      | （与起始金额配套使用，该金额范围）                |
| 渠道代码            | chanCode       |                                  |
| 渠道商户号           | chanMerId      |                                  |


* 商家绑定信息（可能要调用加密机）	bindingInfo

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 商户号             | merId          |                                  |
| 绑定ID            | bindingId      |                                  |
| 卡品牌             | cardBrand      |                                  |
| 卡类型             | acctType       |                                  |
| 账户名称            | acctName       |                                  |
| 账户号码            | acctNum        |                                  |
| 开户行             | bankId         |                                  |
| 证件类型            | identType      |                                  |
| 证件号             | identNum       |                                  |
| 手机号             | phoneNum       |                                  |
| 有效期             | validDate      |                                  |
| CVV2            | cvv2           |                                  |


* 绑定关系映射	bindingMap

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 商户号             | merId          |                                  |
| 源绑定ID           | bindingId      |                                  |
| 目标渠道代码          | chanCode       |                                  |
| 渠道商户号           | chanMerId      |                                  |
| 目标渠道绑定ID        | chanBindingId  |                                  |
| 绑定状态            | bindingStatus  | （成功，失败，或者处理中）                    |

* 交易记录	trans

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 商户订单号           | orderNum       |                                  |
| 网关订单号           | sysOrderNum   |                                  |
| 渠道订单号           | chanOrderNum   |                                  |
| 退款订单号           | refundOrderNum |                                  |
| 绑定ID            | chanBindingId  |                                  |
| 交易账户            | acctNum        |                                  |
| 网关应答码           | respCode       |                                  |
| 商户号             | merId          |                                  |
| 交易金额            | transAmount    |                                  |
| 交易币种            | transCurr      |                                  |
| 交易状态            | transStatus    |                                  |
| 交易类型            | transType      |                                  |
| 渠道商户号           | chanMerId      |                                  |
| 渠道代码            | chanCode       |                                  |
| 渠道应答码           | chanRespCode   |                                  |
| 交易创建时间          | createTime     |                                  |
| 交易更新时间          | updateTime     |                                  |
| 短信流水号           | sendSmsId      |                                  |
| 短信验证码           | smsCode        |                                  |
| 备注              | remark           |                                  |
| 子商户id           | subMerId         |                                  |
| 渠道折扣			| chanDiscount     |	支付宝、微信						|
| 商户折扣			| merDiscount      |	支付宝、微信						|
| 消费帐号			| consumerAccount  |	支付宝、微信						|
| 消费类型			| consumerId       |	支付宝、微信						|
| 业务类型			| busicd           |	支付宝、微信						|
| 机构号			    | inscd            |	支付宝、微信						|

* 清算信息 TransSett

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 交易记录            | tran           |                                  |
| 清算标志            | settFlag       |                                  |
| 清算时间            | settDate       |                                  |
| 商户清算金额          | merSettAmt     |                                  |
| 商户手续费           | merFee         |                                  |
| 渠道清算金额          | chanSettAmt    |                                  |
| 渠道手续费           | chanFee        |                                  |

* 清算日志 TransSettLog

| 名称              | 字段             | 备注                        |
| --------------- | -------------  | -------------                 |
| 清算状态            | status         | 1-成功 2-失败                |
| 清算机器IP          | addr           |                 				|
| 清算日期            | date           | 格式 yyyy-mm-dd              |
| 清算开始时间         | createTime     | yyyy-mm-dd hh:mm:ss         |
| 清算结束时间         | modifyTime     | yyyy-mm-dd hh:mm:ss         |
| 清算任务名称         | method         |                             |


* 计数器 counter

| 名称              | 字段             | 备注                               |
| --------------- | -------------  | -------------                    |
| 序列号的键           | key            | 属于系统唯一序列号，对应结构体SN                |
| 序列号的值           | value          | 属于系统唯一序列号，对应结构体SN                |
| 商户号             | merId          | 属于单日序列号信息，对应结构体DaySN             |
| 终端号             | termId         | 属于单日序列号信息，对应结构体DaySN             |
| 序列号             | sn             | 属于单日序列号信息，对应结构体DaySN             |


* 检查并通知 checkAndNotify

| 名称             | 字段             | 备注                         |
| ---------------- | --------------- | ---------------------------  |
| 业务类型          | bizType         |                              |
| 最新版本          | curTag          |                              |
| 前一版本          | prevTag         |                              |
| 应用1版本         | app1Tag         | 集群部署下，每个应用对应一个字段 |
| 应用2版本         | app2Tag         |                              |
