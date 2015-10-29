package unionlive

import "github.com/CardInfoLink/quickpay/model"

type BaseReq interface {
	GetT() string
	GetSpReq() *model.ScanPayRequest
}

type BaseResp interface {
}

const (
	Version      = "1.0"
	TransDirectQ = "Q"
)
