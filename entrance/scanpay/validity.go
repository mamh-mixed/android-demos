package scanpay

import (
	"github.com/CardInfoLink/quickpay/model"
	// "regexp"
)

// validateBarcodePay 验证扫码下单的参数
func validateBarcodePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	if req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" || req.ScanCodeId == "" {
		return &model.ScanPayResponse{ErrorDetail: "INVALID_PARAMETER"}
	}
	return
}
