package channel

import (
	"github.com/CardInfoLink/quickpay/channel/weixin/enterprisepay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
)

// EnterprisePayChan 企业支付
type EnterprisePayChan interface {
	// ProcessPay 扫条码下单
	ProcessPay(req *model.ScanPayRequest) (*model.ScanPayResponse, error)

	// ProcessEnquiry 查询
	ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error)
}

// GetEnterprisePayChan 企业支付渠道
func GetEnterprisePayChan(chanCode string) EnterprisePayChan {
	switch chanCode {
	// 微信
	case ChanCodeWeixin:
		return &enterprisepay.DefaultClient
	// 支付宝
	// case ChanCodeAlipay:
	// 	return &alipay.DefaultClient
	default:
		log.Errorf("unknown scan pay channel `%s`", chanCode)
		return nil
	}
}
