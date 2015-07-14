package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
)

var OffLineCdCol map[string]string

func init() {
	OffLineCdCol = make(map[string]string)
	OffLineCdCol["00"] = "成功"
	OffLineCdCol["01"] = "交易失败"
	OffLineCdCol["03"] = "商户错误"
	OffLineCdCol["09"] = "处理中"
	OffLineCdCol["12"] = "签名错误"
	OffLineCdCol["13"] = "退款失败"
	OffLineCdCol["14"] = "条码错误或过期"
	OffLineCdCol["15"] = "无此渠道"
	OffLineCdCol["16"] = "撤销失败"
	OffLineCdCol["17"] = "关闭失败"
	OffLineCdCol["19"] = "订单号重复"
	OffLineCdCol["25"] = "订单不存在"
	OffLineCdCol["30"] = "报文错误"
	OffLineCdCol["31"] = "权限不足"
	OffLineCdCol["51"] = "余额不足"
	OffLineCdCol["54"] = "订单已关闭或取消"
	OffLineCdCol["58"] = "未知应答码类型"
	OffLineCdCol["91"] = "外部系统错误"
	OffLineCdCol["96"] = "内部系统错误"
	OffLineCdCol["98"] = "交易超时"

}

// OffLineRespCd 扫码支付应答码
func OffLineRespCd(code string) *model.ScanPayResponse {

	errorDetail, respCd := "", ""

	switch code {
	case "SUCCESS":
		respCd = "00"
	case "INPROCESS":
		respCd = "09"
	case "FAIL":
		respCd = "01"
	case "NO_ROUTERPOLICY", "NO_CHANMER", "NO_PERMISSION":
		respCd = "31"
	case "NOT_PAYTRADE", "NOT_SUCESS_TRADE", "TRADE_REFUNDED", "REFUND_TIME_ERROR":
		respCd = "13"
	case "TRADE_AMT_INCONSISTENT":
		respCd = "13"
	case "CANCEL_TIME_ERROR", "TRADE_HAS_REFUND":
		respCd = "16"
	case "SYSTEM_ERROR", "CONNECT_ERROR":
		respCd = "96"
	case "ORDER_DUPLICATE":
		respCd = "19"
	case "SIGN_AUTH_ERROR":
		respCd = "12"
	case "NO_MERCHANT":
		respCd = "03"
	case "TRADE_OVERTIME":
		respCd = "98"
	case "DATA_ERROR":
		respCd = "30"
	case "QRCODE_INVALID":
		respCd = "14"
	case "NO_CHANNEL":
		respCd = "15"
	case "TRADE_NOT_EXIST":
		respCd = "25"
	case "ORDER_CLOSED":
		respCd = "54"
	case "NOT_SUPPORT_TYPE":
		respCd = "17"
	case "INSUFFICIENT_BALANCE":
		respCd = "51"
	case "UNKNOWN_ERROR":
		respCd = "91"
	default:
		respCd = "58"
	}

	errorDetail = OffLineCdCol[respCd]
	return &model.ScanPayResponse{ErrorDetail: errorDetail, Respcd: respCd}
}
