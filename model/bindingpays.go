package model

// BindingReturn 绑定支付返回
type BindingReturn struct {
	RespCode string `json:"respCode"` // 响应代码
	RespMsg  string `json:"respMsg"`  // 响应信息

	BindingId string `json:"bindingId"` // 银行卡绑定ID

	// 绑定支付响应
	MerOrderNum string `json:"merOrderNum"` // 商户订单号
	OrderNum    string `json:"orderNum"`    // 网关订单号
}

// BindingCreate 建立绑定支付
type BindingCreate struct {
	BindingId string `json:"bindingId"` // 银行卡绑定ID
	AcctName  string `json:"acctName"`  // 账户名称
	AcctNum   string `json:"acctNum"`   // 账户号码
	IdentType string `json:"identType"` // 证件类型
	IdentNum  string `json:"identNum"`  // 证件号码
	PhoneNum  string `json:"phoneNum"`  // 手机号
	AcctType  string `json:"acctType"`  // 账户类型
	ValidDate string `json:"validDate"` // 信用卡有限期
	Cvv2      string `json:"cvv2"`      // CVV2
	SendSmsId string `json:"sendSmsId"` // 发送短信验证码的交易流水
	SmsCode   string `json:"smsCode"`   // 短信验证码
}

// BindingRemove 解除绑定关系请求
type BindingRemove struct {
	BindingId string `json:"bindingId"` // 银行卡绑定ID
}

// BindingEnquiry 绑定关系查询
type BindingEnquiry struct {
	BindingId string `json:"bindingId"` // 银行卡绑定ID
}

// BindingPayment 绑定支付请求
type BindingPayment struct {
	SubMerId    string `json:"subMerId"`    // 子商户号
	MerOrderNum string `json:"merOrderNum"` // 商户订单号
	TransAmt    int    `json:"transAmt"`    // 支付金额
	BindingId   string `json:"bindingId"`   // 银行卡绑定ID
	SendSmsId   string `json:"sendSmsId"`   // 申请短信验证码的交易流水
	SmsCode     string `json:"smsCode"`     // 短信验证码
}

// BindingRefund 退款请求
type BindingRefund struct {
	MerOrderNum  string `json:"merOrderNum"`  // 商户订单号
	OrigOrderNum string `json:"origOrderNum"` // 原支付订单号
	TransAmt     int    `json:"transAmt"`     // 退款金额
}
