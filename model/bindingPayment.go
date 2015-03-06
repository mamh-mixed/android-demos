package model

type BindingPaymentIn struct {
	SubMerId    string `json:"subMerId"`    //子商户号
	MerOrderNum string `json:"merOrderNum"` //商户订单号
	TransAmt    int    `json:"transAmt"`    //支付金额
	BindingId   string `json:"bindingId"`   //银行卡绑定ID
	SendSmsId   string `json:"sendSmsId"`   //申请短信验证码的交易流水
	SmsCode     string `json:"smsCode"`     //短信验证码
}

type BindingPaymentOut struct {
	MerOrderNum string `json:"merOrderNum"` //商户订单号
	OrderNum    string `json:"orderNum"`    //网关订单号
	RespCode    string `json:"respCode"`    //响应代码
	RespMsg     string `json:"respMsg"`     //响应信息
}
