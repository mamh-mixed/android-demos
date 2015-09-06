package model

// BindingReturn 绑定支付返回
type BindingReturn struct {
	RespCode string `json:"respCode" bson:"respCode"` // 响应代码
	RespMsg  string `json:"respMsg" bson:"respMsg"`   // 响应信息

	BindingId string `json:"bindingId,omitempty"` // 银行卡绑定ID

	// 绑定支付响应
	MerOrderNum string `json:"merOrderNum,omitempty"` // 商户订单号
	OrderNum    string `json:"orderNum,omitempty"`    // 网关订单号

	// 交易对账汇总
	SettDate string            `json:"settDate,omitempty"` // 对账日期
	Data     []SummarySettData `json:"data,omitempty"`     // 对账数据集

	// 交易对账明细
	Count        int             `json:"count,omitempty"`        // 拉取的记录条数
	NextOrderNum string          `json:"nextOrderNum,omitempty"` // 拉取列表的后一个记录的订单号
	Rec          []TransSettInfo `json:"rec,omitempty"`          // 交易记录

	// 查询订单状态
	OrigRespCode    string     `json:"origRespCode,omitempty"`    //原交易响应代码
	OrigTransDetail *TransInfo `json:"origTransDetail,omitempty"` //原交易明细信息

	// 渠道返回信息
	ChanRespCode string `json:"-"`
	ChanRespMsg  string `json:"-"`

	// 交易状态
	TransStatus string `json:"transStatus,omitempty"`

	// 绑定状态查询响应
	BindingStatus string `json:"bindingStatus,omitempty"` // 绑定状态 10=绑定处理中；20=绑定失败；30=绑定成功；40=已解绑（绑定成功过，后续解绑也成功）

	// 卡片信息
	AcctType          string `json:"acctType,omitempty"`
	CardBrand         string `json:"cardBrand,omitempty"`
	CardNum           string `json:"cardNum,omitempty"`
	IssBankName       string `json:"issBankName,omitempty"`
	IssBankNum        string `json:"issBankNum,omitempty"`
	BankCode          string `json:"bankCode,omitempty"`
	BindingPaySupport string `json:"bindingPaySupport,omitempty"`
}

// NewBindingReturn 构造函数
func NewBindingReturn(code, msg string) (ret *BindingReturn) {
	// resp := mongo.GetRespCode(code)
	return &BindingReturn{
		RespCode: code,
		RespMsg:  msg,
	}
}

// CardInfo 获取卡片信息
type CardInfo struct {
	MerId   string `json:"merId"`
	CardNum string `json:"cardNum"`
}

// BindingCreate 建立绑定关系
type BindingCreate struct {
	MerId         string `json:"merId" bson:"merId,omitempty"`                 // 商户ID
	BindingId     string `json:"bindingId" bson:"bindingId,omitempty"`         // 银行卡绑定ID
	AcctName      string `json:"acctName" bson:"acctName,omitempty"`           // 账户名称
	AcctNum       string `json:"acctNum" bson:"acctNum,omitempty"`             // 账户号码
	IdentType     string `json:"identType" bson:"identType,omitempty"`         // 证件类型
	IdentNum      string `json:"identNum" bson:"identNum,omitempty"`           // 证件号码
	PhoneNum      string `json:"phoneNum" bson:"phoneNum,omitempty"`           // 手机号
	AcctType      string `json:"acctType" bson:"acctType,omitempty"`           // 账户类型
	ValidDate     string `json:"validDate" bson:"validDate,omitempty"`         // 信用卡有限期
	Cvv2          string `json:"cvv2" bson:"cvv2,omitempty"`                   // CVV2
	SendSmsId     string `json:"sendSmsId" bson:"sendSmsId,omitempty"`         // 发送短信验证码的交易流水
	SmsCode       string `json:"smsCode" bson:"smsCode,omitempty"`             // 短信验证码
	BankCode      string `json:"bankCode" bson:"bankCode,omitempty"`           // 银行ID
	ChanBindingId string `json:"chanBindingId" bson:"chanBindingId,omitempty"` // 渠道绑定ID
	ChanMerId     string `json:"chanMerId" bson:"chanMerId,omitempty"`         // 渠道商户ID
	PrivateKey    string `json:"-"`                                            // 签名密钥
	// 存储解密字段，辅助
	AcctNumDecrypt   string `json:"-"`
	AcctNameDecrypt  string `json:"-"`
	IdentNumDecrypt  string `json:"-"`
	PhoneNumDecrypt  string `json:"-"`
	ValidDateDecrypt string `json:"-"`
	Cvv2Decrypt      string `json:"-"`
}

// BindingRemove 解除绑定关系
type BindingRemove struct {
	MerId         string `json:"merId"`         //商户ID
	BindingId     string `json:"bindingId"`     // 银行卡绑定ID
	TxSNUnBinding string `json:"txSNUnBinding"` //解绑流水号
	ChanBindingId string //渠道绑定ID
	ChanMerId     string //渠道商户ID
	PrivateKey    string //签名密钥
}

// BindingEnquiry 绑定关系查询
type BindingEnquiry struct {
	MerId         string `json:"merId"`     //商户ID
	BindingId     string `json:"bindingId"` // 银行卡绑定ID
	ChanBindingId string //渠道绑定ID
	ChanMerId     string //渠道商户ID
	PrivateKey    string //签名密钥
}

// BindingPayment 绑定支付请求
type BindingPayment struct {
	MerId       string `json:"merId"`       //商户ID
	SubMerId    string `json:"subMerId"`    // 子商户号
	MerOrderNum string `json:"merOrderNum"` // 商户订单号
	TransAmt    int64  `json:"transAmt"`    // 支付金额
	BindingId   string `json:"bindingId"`   // 银行卡绑定ID
	SendSmsId   string `json:"sendSmsId"`   // 申请短信验证码的交易流水
	SmsCode     string `json:"smsCode"`     // 短信验证码
	SettFlag    string `json:"settFlag"`    // 清算标识
	Remark      string `json:"remark"`      // 备注
	TerminalId  string `json:"terminalId"`  // 终端Id
	// 辅助参数
	SysOrderNum   string //系统订单号
	ChanBindingId string //渠道绑定ID
	ChanMerId     string //渠道商户ID
	PrivateKey    string //签名密钥
}

// BindingRefund 退款
type BindingRefund struct {
	MerId           string `json:"merId"`        //商户ID
	MerOrderNum     string `json:"merOrderNum"`  // 商户订单号
	OrigOrderNum    string `json:"origOrderNum"` // 原支付订单号
	TransAmt        int64  `json:"transAmt"`     // 退款金额
	Remark          string `json:"remark"`       //备注
	SysOrderNum     string //系统订单号
	SysOrigOrderNum string //系统原支付订单号
	ChanMerId       string //渠道商户ID
	PrivateKey      string //签名密钥
}

// BillingSummary 交易对账汇总
type BillingSummary struct {
	MerId    string `json:"merId"`    //商户ID
	SettDate string `json:"settDate"` // 对账日期，格式为‘YYYYMMDD’
}

// BillingDetails 交易对账明细
type BillingDetails struct {
	MerId        string `json:"merId"`        //商户ID
	SettDate     string `json:"settDate"`     // 对账日期，格式为‘YYYYMMDD’
	NextOrderNum string `json:"nextOrderNum"` // 拉取的第一条记录的商户订单号,不填默认从头开始拉取，使用上一次调用返回的nextOrderNum可连续拉取
}

// OrderEnquiry 查询订单状态
type OrderEnquiry struct {
	MerId        string `json:"merId"`        //商户ID
	OrigOrderNum string `json:"origOrderNum"` //原交易订单号
	SysOrderNum  string `json:"sysOrderNum"`  //原网关订单号
	ShowOrigInfo string `json:"showOrigInfo"` //是否需要返回原交易详细信息;0:不需要，1:需要,不送默认为0
	ChanMerId    string //渠道商户Id
	PrivateKey   string //签名密钥
}

// NoTrackPayment 无卡直接支付
type NoTrackPayment struct {
	MerId       string `json:"merId"`                // 商户ID
	TransType   string `json:"transType"`            // 交易子类型 SALE:消费（直接扣款）AUTH:预授权
	SubMerId    string `json:"subMerId"`             // 子商户号
	MerOrderNum string `json:"merOrderNum"`          // 商户订单号
	TransAmt    int64  `json:"transAmt"`             // 支付金额
	CurrCode    string `json:"currCode"`             // 交易币种
	AcctName    string `json:"acctName"`             // 账户名称
	AcctNum     string `json:"acctNum"`              // 账户号码
	IdentType   string `json:"identType"`            // 证件类型
	IdentNum    string `json:"identNum"`             // 证件号码
	PhoneNum    string `json:"phoneNum"`             // 手机号
	AcctType    string `json:"acctType"`             // 账户类型
	ValidDate   string `json:"validDate"`            // 信用卡有限期
	Cvv2        string `json:"cvv2"`                 // CVV2
	SendSmsId   string `json:"sendSmsId"`            // 发送短信验证码的交易流水
	SmsCode     string `json:"smsCode"`              // 短信验证码
	Chcd        string `json:"chcd,omitempty"`       //下游商户配置的渠道机构号
	Mchntid     string `json:"mchntid,omitempty"`    //下游商户配置的渠道商户号
	TerminalId  string `json:"terminalId,omitempty"` //下游商户配置的渠道商户的终端号
	CliSN       string `json:"cliSN,omitempty"`      //商户的终端在当天对应的一个序列号
	SysSN       string `json:"sysSN,omitempty"`      //系统序列号
	// 存储解密字段，辅助
	AcctNumDecrypt   string `json:"-"`
	AcctNameDecrypt  string `json:"-"`
	IdentNumDecrypt  string `json:"-"`
	PhoneNumDecrypt  string `json:"-"`
	ValidDateDecrypt string `json:"-"`
	Cvv2Decrypt      string `json:"-"`
}
