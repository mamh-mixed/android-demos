package domain

type BindingCreateRequest struct {
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

type BindingCreateResponse struct {
	BindingId string `json:"bindingId"` //银行卡绑定ID
	RespCode  string `json:"respCode"`  //响应代码
	RespMsg   string `json:"respMsg"`   //响应信息
}
