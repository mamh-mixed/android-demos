package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
)

func OffLineRespCd(code string) *model.ScanPayResponse {
	responseCode := ""
	switch code {
	case "ORDER_FAIL":
		responseCode = "01"

	case "ORDER_SUCCESS_PAY_SUCCESS":
		responseCode = "00"

	case "ORDER_SUCCESS_PAY_FAIL":
		responseCode = "01"

	case "ORDER_SUCCESS_PAY_INPROCESS":
		responseCode = "09"

	case "TRADE_SETTLE_ERROR":
		responseCode = "05"

	case "TRADE_BUYER_NOT_MATCH":
		responseCode = "05"

	case "CONTEXT_INCONSISTENT":
		responseCode = "05"

	case "TRADE_HAS_SUCCESS":
		responseCode = "12"

	case "TRADE_HAS_CLOSE":
		responseCode = "12"

	case "REASON_ILLEGAL_STATUS":
		responseCode = "12"

	case "EXIST_FORBIDDEN_WORD":
		responseCode = "58"

	case "PARTNER_ERROR":
		responseCode = "58"

	case "ACCESS_FORBIDDEN":
		responseCode = "58"

	case "SELLER_NOT_EXIST":
		responseCode = "58"

	case "BUYER_NOT_EXIST":
		responseCode = "57"

	case "BUYER_ENABLE_STATUS_FORBID":
		responseCode = "57"

	case "BUYER_SELLER_EQUAL":
		responseCode = "57"

	case "INVALID_PARAMETER":
		responseCode = "58"

	case "UN_SUPPORT_BIZ_TYPE":
		responseCode = "58"

	case "INVALID_RECEIVE_ACCOUNT":
		responseCode = "58"

	case "BUYER_PAYMENT_AMOUNT_DAY_LIMIT_ERROR":
		responseCode = "57"

	case "ERROR_BUYER_CERTIFY_LEVEL_LIMIT":
		responseCode = "57"

	case "ERROR_SELLER_CERTIFY_LEVEL_LIMIT":
		responseCode = "57"

	case "CLIENT_VERSION_NOT_MATCH":
		responseCode = "57"

	case "AUTH_NO_ERROR":
		responseCode = "12"

	case "ILLEGAL_SIGN":
		responseCode = "96"

	case "ILLEGAL_DYN_MD5_KEY":
		responseCode = "96"

	case "ILLEGAL_ENCRYPT":
		responseCode = "96"

	case "ILLEGAL_ARGUMENT":
		responseCode = "96"

	case "ILLEGAL_SERVICE":
		responseCode = "96"

	case "ILLEGAL_USER":
		responseCode = "96"

	case "ILLEGAL_PARTNER":
		responseCode = "96"

	case "ILLEGAL_EXTERFACE":
		responseCode = "96"

	case "ILLEGAL_PARTNER_EXTERFACE":
		responseCode = "96"

	case "ILLEGAL_SECURITY_PROFILE":
		responseCode = "96"

	case "ILLEGAL_AGENT":
		responseCode = "96"

	case "ILLEGAL_SIGN_TYPE":
		responseCode = "96"

	case "ILLEGAL_CHARSET":
		responseCode = "96"

	case "HAS_NO_PRIVILEGE":
		responseCode = "96"

	case "INVALID_CHARACTER_SET":
		responseCode = "96"

	case "SYSTEM_ERROR":
		responseCode = "91"

	case "SESSION_TIMEOUT":
		responseCode = "91"

	case "ILLEGAL_TARGET_SERVICE":
		responseCode = "91"

	case "ILLEGAL_ACCESS_SWITCH_SYSTEM":
		responseCode = "91"

	case "EXTERFACE_IS_CLOSED":
		responseCode = "91"

	case "TRADE_NOT_EXIST":
		responseCode = "25"

	case "WAIT_BUYER_PAY":
		responseCode = "09"

	case "TRADE_CLOSED":
		responseCode = "02"

	case "TRADE_SUCCESS":
		responseCode = "00"

	case "TRADE_PENDING":
		responseCode = "03"

	case "TRADE_FINISHED":
		responseCode = "04"

	case "Time_out":
		responseCode = "98"

	case "失败":
		responseCode = "01"

	case "成功":
		responseCode = "00"

	default:
		responseCode = "58"
	}
	return &model.ScanPayResponse{ErrorDetail: code, Respcd: responseCode}
}
