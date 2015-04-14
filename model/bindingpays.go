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
}

// NewBindingReturn 构造函数
func NewBindingReturn(code, msg string) (ret *BindingReturn) {
	// resp := mongo.GetRespCode(code)
	return &BindingReturn{
		RespCode: code,
		RespMsg:  msg,
	}
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
	BankId        string `json:"bankId" bson:"bankId,omitempty"`               // 银行ID
	ChanBindingId string `json:"chanBindingId" bson:"chanBindingId,omitempty"` // 渠道绑定ID
	ChanMerId     string `json:"chanMerId" bson:"chanMerId,omitempty"`         // 渠道商户ID
	SignCert      string `json:"－"`                                            //签名密钥
}

// BindingRemove 解除绑定关系
type BindingRemove struct {
	MerId         string `json:"merId"`         //商户ID
	BindingId     string `json:"bindingId"`     // 银行卡绑定ID
	TxSNUnBinding string `json:"txSNUnBinding"` //解绑流水号
	ChanBindingId string //渠道绑定ID
	ChanMerId     string //渠道商户ID
	SignCert      string //签名密钥
}

// BindingEnquiry 绑定关系查询
type BindingEnquiry struct {
	MerId         string `json:"merId"`     //商户ID
	BindingId     string `json:"bindingId"` // 银行卡绑定ID
	ChanBindingId string //渠道绑定ID
	ChanMerId     string //渠道商户ID
	SignCert      string //签名密钥
}

// BindingPayment 绑定支付请求
type BindingPayment struct {
	MerId         string `json:"merId"`       //商户ID
	SubMerId      string `json:"subMerId"`    // 子商户号
	MerOrderNum   string `json:"merOrderNum"` // 商户订单号
	TransAmt      int64  `json:"transAmt"`    // 支付金额
	BindingId     string `json:"bindingId"`   // 银行卡绑定ID
	SendSmsId     string `json:"sendSmsId"`   // 申请短信验证码的交易流水
	SmsCode       string `json:"smsCode"`     // 短信验证码
	SettFlag      string `json:"settFlag"`    //清算标识
	Remark        string `json:"remark"`      //备注
	ChanOrderNum  string //渠道订单号
	ChanBindingId string //渠道绑定ID
	ChanMerId     string //渠道商户ID
	SignCert      string //签名密钥
}

// BindingRefund 退款
type BindingRefund struct {
	MerId            string `json:"merId"`        //商户ID
	MerOrderNum      string `json:"merOrderNum"`  // 商户订单号
	OrigOrderNum     string `json:"origOrderNum"` // 原支付订单号
	TransAmt         int64  `json:"transAmt"`     // 退款金额
	Remark           string `json:"remark"`       //备注
	ChanOrderNum     string //渠道订单号
	ChanOrigOrderNum string //渠道原支付订单号
	ChanMerId        string //渠道商户ID
	SignCert         string //签名密钥
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
	ChanOrderNum string `json:"chanOrderNum"` //原网关订单号
	ShowOrigInfo string `json:"showOrigInfo"` //是否需要返回原交易详细信息;0:不需要，1:需要,不送默认为0
	ChanMerId    string //渠道商户Id
	SignCert     string //签名密钥
}

// NoTrackPayment 无卡直接支付
type NoTrackPayment struct {
	MerId       string `json:"merId"`       //商户ID
	SubMerId    string `json:"subMerId"`    // 子商户号
	MerOrderNum string `json:"merOrderNum"` // 商户订单号
	TransAmt    int    `json:"transAmt"`    // 支付金额
	AcctName    string `json:"acctName"`    // 账户名称
	AcctNum     string `json:"acctNum"`     // 账户号码
	IdentType   string `json:"identType"`   // 证件类型
	IdentNum    string `json:"identNum"`    // 证件号码
	PhoneNum    string `json:"phoneNum"`    // 手机号
	AcctType    string `json:"acctType"`    // 账户类型
	ValidDate   string `json:"validDate"`   // 信用卡有限期
	Cvv2        string `json:"cvv2"`        // CVV2
	SendSmsId   string `json:"sendSmsId"`   // 发送短信验证码的交易流水
	SmsCode     string `json:"smsCode"`     // 短信验证码
}

// Appyle pay
type ApplePay struct {
	MerId         string       `json:"merId"`                //商户ID
	TransType     string       `json:"transType"`            //子交易类型
	SubMerId      string       `json:"subMerId,omitempty"`   //子商户号
	TerminalId    string       `json:"terminalId,omitempty"` //终端号
	MerOrderNum   string       `json:"merOrderNum"`          //商户订单号
	TransactionId string       `json:"transactionId"`        //ApplePay标识
	ApplePayData  ApplePayData `json:"applePayData"`         //ApplePay数据
	SignCert      string       //签名密钥
}

// ApplePayData applePay数据
type ApplePayData struct {
	ApplicationPrimaryAccountNumber string      `json:"applicationPrimaryAccountNumber"` // 主账号
	ApplicationExpirationDate       string      `json:"applicationExpirationDate"`       // 有效期截止日
	CurrencyCode                    string      `json:"currencyCode"`                    // 货币代码
	TransactionAmount               int64       `json:"transactionAmount"`               // 交易金额
	DeviceManufacturerIdentifier    string      `json:"deviceManufacturerIdentifier"`    // 设备制造商标识符
	PaymentDataType                 string      `json:"paymentDataType"`                 // 支付数据类型(EMV或3D Secure)
	PaymentData                     PaymentData `json:"paymentData"`                     // 支付数据内容
}

// PaymentData 支付数据
type PaymentData struct {
	OnlinePaymentCryptogram string `json:"merId,omitempty"`        // 3D Secure类型的在线支付密码
	EciIndicator            string `json:"eciIndicator,omitempty"` // 3D Secure类型的Eci指示符
	EmvData                 string `json:"emvData,omitempty"`      // EMV类型的支付数据，到线下网关的时候存到iccdata里面
}
