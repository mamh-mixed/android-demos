package channel

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1/domestic"
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1/oversea"
	"github.com/CardInfoLink/quickpay/channel/weixin/scanpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// 扫码支付渠道
const (
	ChanCodeWeixin     = "WXP"
	ChanCodeAlipay     = "ALP"
	ChanCodeAliOversea = "AOS" // TODO:未定
)

// ScanPayChan 扫码支付
type ScanPayChan interface {
	// ProcessBarcodePay 扫条码下单
	ProcessBarcodePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error)

	// ProcessQrCodeOfflinePay 扫二维码预下单
	ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error)

	// ProcessRefund 退款
	ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error)

	// ProcessEnquiry 查询
	ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error)

	// ProcessCancel 撤销
	ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error)

	// ProcessClose 关闭
	ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error)

	// ProcessRefundQuery 退款查询
	ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error)
}

// GetScanPayChan 扫码支付渠道
func GetScanPayChan(chanCode string) ScanPayChan {
	switch chanCode {
	// 微信
	case ChanCodeWeixin:
		return &scanpay.DefaultWeixinScanPay
	// 支付宝
	case ChanCodeAlipay:
		return &domestic.DefaultClient
	case ChanCodeAliOversea:
		return &oversea.DefaultClient
	default:
		log.Errorf("unknown scan pay channel `%s`", chanCode)
		return nil
	}
}
