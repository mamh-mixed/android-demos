package scanpay

import (
	"regexp"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

// validateBarcodePay 验证扫码下单的参数
func validateBarcodePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" || req.ScanCodeId == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	return
}

// validateQrCodeOfflinePay 验证预下单的参数
func validateQrCodeOfflinePay(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrderNum == "" || req.Chcd == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}
	// TODO ..
	if req.Chcd != "WXP" && req.Chcd != "ALP" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}
	return
}

// validateEnquiry 验证查询接口的参数
func validateEnquiry(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.Inscd == "" || req.Mchntid == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	// TODO validate format
	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	return
}

// validateRefund 验证退款接口的参数
func validateRefund(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	// TODO validate format
	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	return
}

// validateCancel 验证撤销接口参数
func validateCancel(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	return
}

// validateCancel 验证关闭订单接口参数
func validateClose(req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	return
}
