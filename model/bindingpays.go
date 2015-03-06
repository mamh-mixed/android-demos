package model

// 建立绑定支付
type BindingCreateIn struct {
	BindingId string `json:"bindingId"` //银行卡绑定ID
	AcctName  string `json:"acctName"`  //账户名称
	AcctNum   string `json:"acctNum"`   //账户号码
	IdentType string `json:"identType"` //证件类型
	IdentNum  string `json:"identNum"`  //证件号码
	PhoneNum  string `json:"phoneNum"`  //手机号
	AcctType  string `json:"acctType"`  //账户类型
	ValidDate string `json:"validDate"` //信用卡有限期
	Cvv2      string `json:"cvv2"`      //CVV2
	SendSmsId string `json:"sendSmsId"` //发送短信验证码的交易流水
	SmsCode   string `json:"smsCode"`   //短信验证码
}

// 解除绑定支付
type BindingCreateOut struct {
	BindingId string `json:"bindingId"` //银行卡绑定ID
	RespCode  string `json:"respCode"`  //响应代码
	RespMsg   string `json:"respMsg"`   //响应信息
}

// 解除绑定关系请求
type BindingRemoveIn struct {
	BindingId string `json:"bindingId"` //银行卡绑定ID
}

// 解除绑定关系响应
type BindingRemoveOut struct {
	RespCode string `json:"respCode"` //响应代码
	RespMsg  string `json:"respMsg"`  //响应信息
}

// 绑定关系查询请求
type BindingEnquiryIn struct {
	BindingId string `json:"bindingId"` //银行卡绑定ID
}

// 绑定关系查询响应
type BindingEnquiryOut struct {
	RespCode string `json:"respCode"` //响应代码
	RespMsg  string `json:"respMsg"`  //响应信息
}

// 绑定支付请求
type BindingPaymentIn struct {
	SubMerId    string `json:"subMerId"`    //子商户号
	MerOrderNum string `json:"merOrderNum"` //商户订单号
	TransAmt    int    `json:"transAmt"`    //支付金额
	BindingId   string `json:"bindingId"`   //银行卡绑定ID
	SendSmsId   string `json:"sendSmsId"`   //申请短信验证码的交易流水
	SmsCode     string `json:"smsCode"`     //短信验证码
}

// 绑定支付响应
type BindingPaymentOut struct {
	MerOrderNum string `json:"merOrderNum"` //商户订单号
	OrderNum    string `json:"orderNum"`    //网关订单号
	RespCode    string `json:"respCode"`    //响应代码
	RespMsg     string `json:"respMsg"`     //响应信息
}

// 退款请求
type RefundIn struct {
	MerOrderNum  string `json:"merOrderNum"`  //商户订单号
	OrigOrderNum string `json:"origOrderNum"` //原支付订单号
	TransAmt     int    `json:"transAmt"`     //退款金额
}

// 退款响应
type RefundOut struct {
	MerOrderNum string `json:"merOrderNum"` //商户订单号
	OrderNum    string `json:"orderNum"`    //网关订单号
	RespCode    string `json:"respCode"`    //响应代码
	RespMsg     string `json:"respMsg"`     //响应信息
}
