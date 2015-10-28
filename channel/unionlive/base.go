package unionlive

type BaseReq interface {
	GetT() string
}

type BaseResp interface {
}

const (
	Version      = "1.0"
	TransDirectQ = "Q"
)
