package model

// BindingReturn 绑定支付返回
type BindingReturn struct {
	RespCode string `json:"respCode"` // 响应代码
	RespMsg  string `json:"respMsg"`  // 响应信息

	BindingId string `json:"bindingId,omitempty"` // 银行卡绑定ID

	// 绑定支付响应
	MerOrderNum string `json:"merOrderNum,omitempty"` // 商户订单号
	OrderNum    string `json:"orderNum,omitempty"`    // 网关订单号

	// 交易对账汇总
	SettDate string   `json:"settDate,omitempty"` // 对账日期
	Data     []string `json:"data,omitempty"`     // 对账数据集
	// 交易对账明细
	Count        int      `json:"count,omitempty"`        // 拉取的记录条数
	NextOrderNum string   `json:"nextOrderNum,omitempty"` // 拉取列表的后一个记录的订单号
	Rec          []string `json:"rec,omitempty"`          // 交易记录

	// 查询订单状态
	OrigRespCode string `json:"origRespCode,omitempty"` //原交易响应代码
	// OrigTransDetail object `json:"origTransDetail,omitempty"` //原交易明细信息

}

//bindingReturn的构造函数
func NewBindingReturn(code, msg string) (ret *BindingReturn) {
	return &BindingReturn{
		RespCode: code,
		RespMsg:  msg,
	}
}

// BindingCreate 建立绑定支付
type BindingCreate struct {
	BindingId string `json:"bindingId" bson:"bindingId,omitempty"`     // 银行卡绑定ID
	MerId     string `json:"merId" bson:"merId,omitempty"`             // 商户ID
	AcctName  string `json:"acctName" bson:"acctName,omitempty"`       // 账户名称
	AcctNum   string `json:"acctNum" bson:"acctNum,omitempty"`         // 账户号码
	IdentType string `json:"identType" bson:"identType,omitempty"`     // 证件类型
	IdentNum  string `json:"identNum" bson:"identNum,omitempty"`       // 证件号码
	PhoneNum  string `json:"phoneNum" bson:"phoneNum,omitempty"`       // 手机号
	AcctType  string `json:"acctType" bson:"acctType,omitempty"`       // 账户类型
	ValidDate string `json:"validDate" bson:"validDate,omitempty"`     // 信用卡有限期
	Cvv2      string `json:"cvv2" bson:"cvv2,omitempty"`               // CVV2
	SendSmsId string `json:"sendSmsId" bson:"sendSmsId,omitempty"`     // 发送短信验证码的交易流水
	SmsCode   string `json:"smsCode" bson:"smsCode,omitempty"`         // 短信验证码
	BankId    string `json:"bankId,omitempty" bson:"bankId,omitempty"` //银行ID
}

// BindingRemove 解除绑定关系请求
type BindingRemove struct {
	BindingId     string `json:"bindingId"`     // 银行卡绑定ID
	MerId         string `json:"merId"`         //商户ID
	TxSNUnBinding string `json:"txSNUnBinding"` //解绑流水号
}

// BindingEnquiry 绑定关系查询
type BindingEnquiry struct {
	BindingId string `json:"bindingId"` // 银行卡绑定ID
	MerId     string `json:"merId"`     //商户ID

}

// BindingPayment 绑定支付请求
type BindingPayment struct {
	SubMerId       string `json:"subMerId"`       // 子商户号
	MerOrderNum    string `json:"merOrderNum"`    // 商户订单号
	TransAmt       int64  `json:"transAmt"`       // 支付金额
	BindingId      string `json:"bindingId"`      // 银行卡绑定ID
	SendSmsId      string `json:"sendSmsId"`      // 申请短信验证码的交易流水
	SmsCode        string `json:"smsCode"`        // 短信验证码
	SettlementFlag string `json:"settlementFlag"` //清算标识
	MerId          string `json:"merId"`          //商户ID
	Remark         string `json:"remark"`         //备注
}

// BindingRefund 退款请求
type BindingRefund struct {
	MerOrderNum  string `json:"merOrderNum"`  // 商户订单号
	OrigOrderNum string `json:"origOrderNum"` // 原支付订单号
	TransAmt     int64  `json:"transAmt"`     // 退款金额
	MerId        string `json:"merId"`        //商户ID
	Remark       string `json:"remark"`       //备注
}

// 交易对账汇总请求
type BillingSummary struct {
	SettDate string `json:"settDate"` // 对账日期，格式为‘YYYYMMDD’
}

// 交易对账明细
type BillingDetails struct {
	SettDate     string `json:"settDate"`     // 对账日期，格式为‘YYYYMMDD’
	NextOrderNum string `json:"nextOrderNum"` // 拉取的第一条记录的商户订单号,不填默认从头开始拉取，使用上一次调用返回的nextOrderNum可连续拉取
}

// 查询订单状态
type OrderEnquiry struct {
	OrigOrderNum string `json:"origOrderNum"` //原交易订单号
	OrderNum     string `json:"orderNum"`     //原网关订单号
	ShowOrigInfo string `json:"showOrigInfo"` //是否需要返回原交易详细信息;0:不需要，1:需要,不送默认为0
}

// 无卡直接支付
type NoTrackPayment struct {
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
