package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

// requestError 请求错误
// 分为接入错误和支付宝系统错误
func requestError(code string) *model.ScanPayResponse {
	switch code {
	// 支付宝系统错误
	case "ILLEGAL_TARGET_SERVICE",
		"ILLEGAL_ACCESS_SWITCH_SYSTEM",
		"SYSTEM_ERROR",
		"SESSION_TIMEOUT",
		"EXTERFACE_IS_CLOSED":
		return mongo.OffLineRespCd("UNKNOWN_ERROR")
	// 讯联系统错误
	case "ILLEGAL_SIGN", "ILLEGAL_DYN_MD5_KEY",
		"ILLEGAL_ENCRYPT",
		"ILLEGAL_ARGUMENT",
		"ILLEGAL_SERVICE",
		"ILLEGAL_USER",
		"ILLEGAL_PARTNER",
		"ILLEGAL_EXTERFACE",
		"ILLEGAL_PARTNER_EXTERFACE",
		"ILLEGAL_SECURITY_PROFILE",
		"ILLEGAL_AGENT",
		"ILLEGAL_SIGN_TYPE",
		"ILLEGAL_CHARSET",
		"HAS_NO_PRIVILEGE",
		"INVALID_CHARACTER_SET":
		return mongo.OffLineRespCd("SYSTEM_ERROR")
	default:
		return mongo.OffLineRespCd("UNKNOWN_ERROR")
	}
}

func preCreateCd(code string) string {

	switch code {
	case "SUCCESS":
		return "09"
	case "TRADE_SETTLE_ERROR", "CONTEXT_INCONSISTENT", "TRADE_BUYER_NOT_MATCH":
		return "05"
	case "TRADE_HAS_SUCCESS", "TRADE_HAS_CLOSE":
		return "12"
	default:
		return "58"
	}
}
func createAndPayCd(code string) string {
	switch code {
	case "ORDER_SUCCESS_PAY_SUCCESS":
		return "00"
	case "ORDER_FAIL", "ORDER_SUCCESS_PAY_FAIL", "UNKNOWN":
		return "01"
	case "ORDER_SUCCESS_PAY_INPROCESS":
		return "09"
	default:
		return "58"
	}
}
func refundCd(code string) string {

	switch code {
	case "SUCCESS":
		return "00"
	case "INVALID_PARAMETER":
		return "58"
	case "TRADE_HAS_CLOSE":
		return "54"
	case "TRADE_NOT_EXIST":
		return "25"
	default:
		return "58"
	}
}
func queryCd(code string) string {
	switch code {
	// 8.4交易状态
	case "WAIT_BUYER_PAY":
		return "09"
	case "TRADE_CLOSED":
		return "54"
	case "TRADE_SUCCESS":
		return "00"
	case "TRADE_PENDING":
		return "01"
	case "TRADE_FINISHED":
		// TODO check
		return "00"
	// 8.1业务错误码
	case "TRADE_NOT_EXIST":
		return "25"
	default:
		return "58"
	}

}
func cancelCd(code string) string {
	switch code {
	case "SUCCESS":
		return "00"
	case "TRADE_NOT_EXIST":
		return "25"
	default:
		return "58"
	}
}
