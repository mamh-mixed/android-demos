package scanpay2

import (
	"github.com/CardInfoLink/quickpay/mongo"
)

// 应答码
var (
	successCode, successMsg, _     = mongo.ScanPayRespCol.Get8583CodeAndMsg("SUCCESS")
	inprocessCode, inprocessMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("INPROCESS")
	closeCode, closeMsg, _         = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
)

func transform(transType, code, msg, subCode, subMsg string) (respCode, respMsg string) {

	switch code {
	case "10000":
		if transType == "precreate" {
			return inprocessCode, inprocessMsg
		}
		// 业务处理成功
		respCode, respMsg = successCode, successMsg
	case "40004":
		// 业务处理失败
		resp := mongo.ScanPayRespCol.GetByAlp2(subCode, transType)
		respCode, respMsg = resp.ISO8583Code, resp.ISO8583Msg
		if !resp.IsUseISO {
			respMsg = subMsg // 使用渠道应答
		}
	case "20000", "10003":
		// 20000 业务出现未知错误或者系统异常，如果支付接口返回，需要调用查询接口确认订单状态或者发起撤销
		// 10003 业务处理中
		respCode, respMsg = inprocessCode, inprocessMsg
	}

	return
}
