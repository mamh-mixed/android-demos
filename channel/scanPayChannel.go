package channel

import (
	"github.com/CardInfoLink/quickpay/channel/alipay"
	"github.com/CardInfoLink/quickpay/model"
)

// ScanPayChan 扫码支付
type ScanPayChan interface {
	// ProcessBarcodePay 扫条码下单
	ProcessBarcodePay(req *model.ScanPay) *model.ScanPayResponse

	// ProcessQrCodeOfflinePay 扫二维码预下单
	ProcessQrCodeOfflinePay(req *model.ScanPay) *model.ScanPayResponse

	// ProcessRefund 退款
	ProcessRefund(req *model.ScanPay) *model.ScanPayResponse

	// ProcessEnquiry 查询
	ProcessEnquiry(req *model.ScanPay) *model.ScanPayResponse
}

// GetScanPayChan 扫码支付渠道
func GetScanPayChan(chanCode string) ScanPayChan {

	switch chanCode {
	// 支付宝
	case "ALP":
		return &alipay.DefaultClient
	// 微信
	case "WXP":
		return nil
	}
	return nil
}
