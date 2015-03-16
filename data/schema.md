
# 名字解释

|名词       | 缩写 | 解释 |
|----------|------|------|
|merchant  | mer  | 商户 |
|channel   | chan | 渠道 |
|CFCA      | cfca | 中金 |
|settlement|sett  | 清分 |


# 数据模式

1. 商户基本信息	merchant
商户号	merId
商户状态（Normal，Deleted）	merStatus
商户交易币种	transCurr
商户签名密钥	signKey
商户加密密钥	encryptKey

2. 商户详细信息 merDetail	merDetail
商户号	merId
商户名称	merName
商户简称	shortName
商户城市	city
商户国家	nation
商户类型	type
商户计费方案代码	billingScheme
商户清算币种	SettCurr
商户账户名称	acctName
商户账户	acctNum
法人代表	corp
商户负责人	master
商户联系人	contact
商户联系电话	contactTel
商户传真	fax
商户邮箱	email
商户地址	addr
商户邮编	postcode
商户密码	password

3. 渠道商户信息：PK（渠道代码，商户号）	chanMer
渠道代码	chanCode
渠道商户号	chanMerId
渠道商户名称	chanMerName
清算角色	settRole
结算标识（中金渠道商户特有）	settFlag
签名证书	signCert
验签证书	checkSignCert
支付宝 MD5 Key	alpMd5Key
微信支付App Id	wxpAppId
微信支付Partner Key	wxpPartnerKey
微信支付加密证书	wxpEncryptCert

4. 应答码	respCode
应答码	respCode
应答信息	respMsg
CFCA	cfca
CFCA 应答码	cfca.code
CFCA 应答信息	cfca.msg

5. 卡Bin	cardBin
卡Bin	bin
卡Bin长度	binLen
卡号长度	cardLen
卡品牌	cardBrand
发卡行代码 issBankCode
发卡行名称 issBankName

6. 路由策略	routerPolicy
源商户号	merId
卡品牌	cardBrand
卡类型	cardType
交易类型	transType
卡Bin组	binGroup
输入方式	inputWay
起始金额	minAmount
最大金额（与起始金额配套使用，该金额范围）	maxAmount
渠道代码	chanCode
渠道商户号	chanMerId

7. 商家绑定信息（可能要调用加密机）	bindingInfo
商户号	merId
绑定ID	bindingId
卡品牌	cardBrand
卡类型	cardType
账户名称	acctName
账户号码	acctName
开户行	depositBank
证件类型	credType
证件号	credNum
手机号	phoneNum
有效期	expired
CVV2	cvv2

8. 绑定关系映射	bindingMap
商户号	merId
源绑定ID	bindingId
目标渠道代码	chanCode
渠道商户号	chanMerId
目标渠道绑定ID	chanBindingId

9. 交易记录	trans
商户订单号	orderNum
网关订单号	chanOrderNum
绑定ID	chanBindingId
交易账户	acctNum
网关应答码	respCode
商户号	merId
交易金额	transAmount
交易币种	transCurr
交易状态	transStatus
转换前交易类型(支付、退货)	beforeType
转换后交易类型（支付、退货、预授权）	afterType
渠道商户号	chanMerId
渠道代码	chanCode
渠道应答码	chanRespCode
交易创建时间	createTime
交易更新时间	updateTime
