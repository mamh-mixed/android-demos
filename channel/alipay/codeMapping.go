package alipay

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
		return "12"
	case "TRADE_NOT_EXIST":
		return "25"
	default:
		return "58"
	}
}
func queryCd(service, code string) string {
	switch code {
	// 8.4交易状态
	case "WAIT_BUYER_PAY":
		return "09"
	case "TRADE_CLOSED":
		if service == "void" {
			return "00"
		} else if service == "purc" {
			return "02"
		} else {
			return "05"
		}
	case "TRADE_SUCCESS":
		if service == "void" {
			return "02"
		} else {
			return "00"
		}
	case "TRADE_PENDING":
		return "03"
	case "TRADE_FINISHED":
		if service == "void" || service == "purc" {
			return "04"
		} else {
			return "00"
		}

	// 8.1业务错误码
	case "TRADE_NOT_EXIST":
		return "25"
	default:
		return "58"
	}

}
func cancelCd(code string) string {
	return ""
}
