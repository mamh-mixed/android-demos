package weixin

import (
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// 应答码
var (
	success      = mongo.ScanPayRespCol.Get("SUCCESS")
	inprocess    = mongo.ScanPayRespCol.Get("INPROCESS")
	unknownError = mongo.ScanPayRespCol.Get("CHAN_UNKNOWN_ERROR")
)

// transformX 根据业务类型和错误码查找应答码
// returnCode: 通信标识
// resultCode: 业务结果标识
// errCode: 渠道返回的错误码
func Transform(busicd, returnCode, resultCode, errCode, errCodeDes string) (status, msg string) {
	// 如果通信标识为失败，一般‘签名失败’，‘参数格式校验失败’都会返回失败的通信标识
	if returnCode == "FAIL" {
		log.Error("weixin request fail, return code is FAIL")
		return unknownError.ISO8583Code, unknownError.ISO8583Msg
	}

	// 如果业务结果标识成功，直接返回给前台成功的应答码
	if resultCode == "SUCCESS" {
		// 预下单时返回09
		if busicd == "prePay" {
			return inprocess.ISO8583Code, inprocess.ISO8583Msg
		}
		return success.ISO8583Code, success.ISO8583Msg
	}

	// 业务结果失败，则根据具体的错误码转换对应的应答码
	respCode := mongo.ScanPayRespCol.GetByWxp(errCode, busicd)
	log.Debugf("response code is %#v", respCode)

	if respCode.IsUseISO {
		return respCode.ISO8583Code, respCode.ISO8583Msg
	}

	errCodeDesRune := []rune(errCodeDes)
	if len(errCodeDesRune) > 64 {
		errCodeDes = string(errCodeDesRune[:64])
	}

	return respCode.ISO8583Code, errCodeDes
}
