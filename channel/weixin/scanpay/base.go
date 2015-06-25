package scanpay

// BaseReq 只是为了注入签名方便
type BaseReq interface {
	GenSign()
}

// BaseResp 只是为了传参方便
type BaseResp interface{}
