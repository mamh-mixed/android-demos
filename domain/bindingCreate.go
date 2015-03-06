package domain

type BindingCreateRequest struct {
	BindingId string `json:"bindingId"`
	AcctName  string `json:"acctName"`
	AcctNum   string `json:"acctNum"`
	IdentType string `json:"identType"`
	IdentNum  string `json:"identNum"`
	PhoneNum  string `json:"phoneNum"`
	AcctType  string `json:"acctType"`
	ValidDate string `json:"validDate"`
	Cvv2      string `json:"cvv2"`
	SendSmsId string `json:"sendSmsId"`
	SmsCode   string `json:"smsCode"`
}

type BindingCreateResponse struct {
	BindingId string `json:"bindingId"`
	RespCode  string `json:"respCode"`
	RespMsg   string `json:"respMsg"`
}
