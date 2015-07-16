package scanpay

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"regexp"
)

// validateBarcodePay 验证扫码下单的参数
func validateBarcodePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" || req.ScanCodeId == "" {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{14,24}$`, req.ScanCodeId); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	return
}

// validateQrCodeOfflinePay 验证预下单的参数
func validateQrCodeOfflinePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	fmt.Println(req.OrderNum == "")
	// 验证非空
	if req.OrderNum == "" || req.Chcd == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" {
		return mongo.OffLineRespCd("DATA_ERROR")
	}
	// TODO ..
	if req.Chcd != "WXP" && req.Chcd != "ALP" {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}
	return
}

// validateEnquiry 验证查询接口的参数
func validateEnquiry(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.Inscd == "" || req.Mchntid == "" {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	// TODO validate format
	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	return
}

// validateRefund 验证退款接口的参数
func validateRefund(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" || req.Txamt == "" {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	// TODO validate format
	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	return
}

// validateCancel 验证撤销接口参数
func validateCancel(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	return
}

// validateCancel 验证关闭订单接口参数
func validateClose(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	if req.OrigOrderNum == "" || req.OrderNum == "" || req.Inscd == "" || req.Mchntid == "" {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return mongo.OffLineRespCd("DATA_ERROR")
	}

	return
}
