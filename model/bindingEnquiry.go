package model

type BindingEnquiryIn struct {
	BindingId string `json:"bindingId"` //银行卡绑定ID
}
type BindingEnquiryOut struct {
	RespCode string `json:"respCode"` //响应代码
	RespMsg  string `json:"respMsg"`  //响应信息
}
