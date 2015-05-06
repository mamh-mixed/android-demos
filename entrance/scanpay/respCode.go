package scanpay

func responseCode(detail, chcd string) string {

	if chcd == "Alipay" {
		return alipayRespCd(detail)
	} else if chcd == "Weixin" {
		return weixinRespCd(detail)
	} else {
		return "58"
	}
}

func alipayRespCd(detail string) string {
	respCd := ""
	switch detail {
	case "ORDER_FAIL":
		respCd = "01"

	case "ORDER_SUCCESS_PAY_SUCCESS":
		respCd = "00"

	case "ORDER_SUCCESS_PAY_FAIL":
		respCd = "01"

	case "ORDER_SUCCESS_PAY_INPROCESS":
		respCd = "09"

	case "TRADE_SETTLE_ERROR":
		respCd = "05"

	case "TRADE_BUYER_NOT_MATCH":
		respCd = "05"

	case "CONTEXT_INCONSISTENT":
		respCd = "05"

	case "TRADE_HAS_SUCCESS":
		respCd = "12"

	case "TRADE_HAS_CLOSE":
		respCd = "12"

	case "REASON_ILLEGAL_STATUS":
		respCd = "12"

	case "EXIST_FORBIDDEN_WORD":
		respCd = "58"

	case "PARTNER_ERROR":
		respCd = "58"

	case "ACCESS_FORBIDDEN":
		respCd = "58"

	case "SELLER_NOT_EXIST":
		respCd = "58"

	case "BUYER_NOT_EXIST":
		respCd = "57"

	case "BUYER_ENABLE_STATUS_FORBID":
		respCd = "57"

	case "BUYER_SELLER_EQUAL":
		respCd = "57"

	case "INVALID_PARAMETER":
		respCd = "58"

	case "UN_SUPPORT_BIZ_TYPE":
		respCd = "58"

	case "INVALID_RECEIVE_ACCOUNT":
		respCd = "58"

	case "BUYER_PAYMENT_AMOUNT_DAY_LIMIT_ERROR":
		respCd = "57"

	case "ERROR_BUYER_CERTIFY_LEVEL_LIMIT":
		respCd = "57"

	case "ERROR_SELLER_CERTIFY_LEVEL_LIMIT":
		respCd = "57"

	case "CLIENT_VERSION_NOT_MATCH":
		respCd = "57"

	case "AUTH_NO_ERROR":
		respCd = "12"

	case "ILLEGAL_SIGN":
		respCd = "96"

	case "ILLEGAL_DYN_MD5_KEY":
		respCd = "96"

	case "ILLEGAL_ENCRYPT":
		respCd = "96"

	case "ILLEGAL_ARGUMENT":
		respCd = "96"

	case "ILLEGAL_SERVICE":
		respCd = "96"

	case "ILLEGAL_USER":
		respCd = "96"

	case "ILLEGAL_PARTNER":
		respCd = "96"

	case "ILLEGAL_EXTERFACE":
		respCd = "96"

	case "ILLEGAL_PARTNER_EXTERFACE":
		respCd = "96"

	case "ILLEGAL_SECURITY_PROFILE":
		respCd = "96"

	case "ILLEGAL_AGENT":
		respCd = "96"

	case "ILLEGAL_SIGN_TYPE":
		respCd = "96"

	case "ILLEGAL_CHARSET":
		respCd = "96"

	case "HAS_NO_PRIVILEGE":
		respCd = "96"

	case "INVALID_CHARACTER_SET":
		respCd = "96"

	case "SYSTEM_ERROR":
		respCd = "91"

	case "SESSION_TIMEOUT":
		respCd = "91"

	case "ILLEGAL_TARGET_SERVICE":
		respCd = "91"

	case "ILLEGAL_ACCESS_SWITCH_SYSTEM":
		respCd = "91"

	case "EXTERFACE_IS_CLOSED":
		respCd = "91"

	case "TRADE_NOT_EXIST":
		respCd = "25"

	case "WAIT_BUYER_PAY":
		respCd = "09"

	case "TRADE_CLOSED":
		respCd = "02"

	case "TRADE_SUCCESS":
		respCd = "00"

	case "TRADE_PENDING":
		respCd = "03"

	case "TRADE_FINISHED":
		respCd = "04"

	case "Time_out":
		respCd = "98"

	case "失败":
		respCd = "01"

	case "成功":
		respCd = "00"

	default:
		respCd = "58"

	}
	return respCd
}

func weixinRespCd(detail string) string {

	respCd := ""
	switch detail {
	case "Time_out":
		respCd = "98"
	case "SYSTEMERROR":
		respCd = "91"
	case "INVALID_TRANSACTIONID":
		respCd = "96"
	case "PARAM_ERROR ":
		respCd = "96"
	case "ORDERPAID":
		respCd = "12"
	case "OUT_TRADE_NO_USED":
		respCd = "12"
	case "NOAUTH":
		respCd = "58"
	case "AUTHCODEEXPIRE":
		respCd = "12"
	case "NOTENOUGH":
		respCd = "57"
	case "NOTSUPORTCARD":
		respCd = "58"
	case "ORDERCLOSED":
		respCd = "12"
	case "ORDERREVERSED":
		respCd = "12"
	case "BANKERROR":
		respCd = "91"
	case "USERPAYING":
		respCd = "09"
	case "AUTH_CODE_ERROR":
		respCd = "12"
	case "AUTH_CODE_INVALID":
		respCd = "12"
	case "TRADE_STATE_ERROR":
		respCd = "01"
	case "REFUND_FEE_INVALID":
		respCd = "12"
	case "SUCCESS":
		respCd = "00"
	case "ILLEGAL_SIGN":
		respCd = "96"
	case "CLOSED":
		respCd = "12"
	case "PAYERROR":
		respCd = "01"
	case "NOTPAY":
		respCd = "12"
	case "NOPAY":
		respCd = "12"
	case "REVOKED":
		respCd = "12"
	case "REFUND":
		respCd = "12"
	default:
		respCd = "58"
	}
	return respCd
}
