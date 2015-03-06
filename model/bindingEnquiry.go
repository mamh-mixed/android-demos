package model

type BindingEnquiryIn struct {
	BindingId string `json:"bindingId"`
}
type BindingEnquiryOut struct {
	RespCode string `json:"respCode"`
	RespMsg  string `json:"respMsg"`
}
