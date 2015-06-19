package scanpay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	// "regexp"
)

// validateBarcodePay 验证扫码下单的参数
func validateBarcodePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" || req.ScanCodeId == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	// TODO validate format

	return
}

// validateQrCodeOfflinePay 验证预下单的参数
func validateQrCodeOfflinePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrderNum == "" || req.Chcd == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	// TODO ..
	return
}

// validateEnquiry 验证查询接口的参数
func validateEnquiry(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.Inscd == "" || req.Mchntid == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	// TODO validate format

	return
}
