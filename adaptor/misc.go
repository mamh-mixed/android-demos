package adaptor

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

// 使用8583应答
// TODO 常用状态码整合到一起
var (
	CloseCode, CloseMsg, _         = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
	FailCode, FailMsg, _           = mongo.ScanPayRespCol.Get8583CodeAndMsg("FAIL")
	InprocessCode, InprocessMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("INPROCESS")
	SuccessCode, SuccessMsg, _     = mongo.ScanPayRespCol.Get8583CodeAndMsg("SUCCESS")
	UnKnownCode, UnKnownMsg, _     = mongo.ScanPayRespCol.Get8583CodeAndMsg("CHAN_UNKNOWN_ERROR")
)

// ReturnWithErrorCode 使用错误码直接返回
func ReturnWithErrorCode(errorCode string) *model.ScanPayResponse {
	spResp := mongo.ScanPayRespCol.Get(errorCode)
	return &model.ScanPayResponse{
		Respcd:      spResp.ISO8583Code,
		ErrorDetail: spResp.ISO8583Msg,
		ErrorCode:   errorCode,
	}
}

// LogicErrorHandler 逻辑错误处理
func LogicErrorHandler(t *model.Trans, errorCode string) *model.ScanPayResponse {
	spResp := mongo.ScanPayRespCol.Get(errorCode)
	// 8583应答
	code, msg := spResp.ISO8583Code, spResp.ISO8583Msg

	// 交易保存
	t.RespCode = code
	t.ErrorDetail = msg
	t.LockFlag = 0
	mongo.SpTransColl.Add(t)

	return &model.ScanPayResponse{
		Respcd:      code,
		ErrorDetail: msg,
		ErrorCode:   errorCode,
	}
}
