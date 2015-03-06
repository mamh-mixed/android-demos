package model

type BindingRemoveIn struct {
	BindingId string `json:"bindingId"`
}

type BindingRemoveOut struct {
	RespCode string `json:"respCode"`
	RespMsg  string `json:"respMsg"`
}
