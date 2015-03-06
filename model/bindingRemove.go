package model

type BindingRemoveIn struct {
	BindingId string `json:"bindingId"` //银行卡绑定ID
}

type BindingRemoveOut struct {
	RespCode string `json:"respCode"` //响应代码
	RespMsg  string `json:"respMsg"`  //响应信息
}
