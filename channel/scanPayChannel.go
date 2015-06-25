package channel

import (
	"github.com/CardInfoLink/quickpay/channel/alipay"
	"github.com/CardInfoLink/quickpay/channel/weixin/scanpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// 扫码支付渠道
const (
	ChanCodeWeixin = "WXP"
	ChanCodeAlipay = "ALP"
)

// ScanPayChan 扫码支付
type ScanPayChan interface {
	// ProcessBarcodePay 扫条码下单
	ProcessBarcodePay(req *model.ScanPay) (*model.ScanPayResponse, error)

	// ProcessQrCodeOfflinePay 扫二维码预下单
	ProcessQrCodeOfflinePay(req *model.ScanPay) (*model.ScanPayResponse, error)

	// ProcessRefund 退款
	ProcessRefund(req *model.ScanPay) (*model.ScanPayResponse, error)

	// ProcessEnquiry 查询
	ProcessEnquiry(req *model.ScanPay) (*model.ScanPayResponse, error)

	// ProcessCancel 撤销
	ProcessCancel(req *model.ScanPay) (*model.ScanPayResponse, error)
}

// GetScanPayChan 扫码支付渠道
func GetScanPayChan(chanCode string) ScanPayChan {
	switch chanCode {
	// 微信
	case ChanCodeWeixin:
		return &scanpay.DefaultWeixinScanPay
	// 支付宝
	case ChanCodeAlipay:
		return &alipay.DefaultClient
	default:
		log.Errorf("unknown scan pay channel `%s`", chanCode)
		return nil
	}
}
