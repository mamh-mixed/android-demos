package weixin

import (
	"github.com/CardInfoLink/quickpay/mongo"
)

// 应答码
var (
	success     = mongo.ScanPayRespCol.Get("SUCCESS")
	inprocess   = mongo.ScanPayRespCol.Get("INPROCESS")
	systemError = mongo.ScanPayRespCol.Get("SYSTEM_ERROR")
)

// transformX 根据业务类型和错误码查找应答码
// returnCode: 通信标识
// resultCode: 业务结果标识
// errCode: 渠道返回的错误码
func Transform(busicd, returnCode, resultCode, errCode, errCodeDes string) (status, msg, errorCode string) {

	// 成功直接返回
	if returnCode == "SUCCESS" {
		// 如果业务结果标识成功，直接返回给前台成功的应答码
		if resultCode == "SUCCESS" {
			// 预下单时返回09
			if busicd == "prePay" {
				return inprocess.ISO8583Code, inprocess.ISO8583Msg, inprocess.ErrorCode
			}
			return success.ISO8583Code, success.ISO8583Msg, success.ErrorCode
		}
		// 错误码映射
	} else {
		// 查询接口通讯失败，返回处理中
		if busicd == "payQuery" {
			return inprocess.ISO8583Code, inprocess.ISO8583Msg, inprocess.ErrorCode
		}
		// 错误码映射
	}

	// 微信系统错误、银行错误
	if errCode == "SYSTEMERROR" || errCode == "BANKERROR" {
		// 默认为处理中
		return inprocess.ISO8583Code, inprocess.ISO8583Msg, inprocess.ErrorCode
	}

	// 业务结果失败，则根据具体的错误码转换对应的应答码
	respCode := mongo.ScanPayRespCol.GetByWxp(errCode, busicd)
	if respCode.IsUseISO {
		return respCode.ISO8583Code, respCode.ISO8583Msg, respCode.ErrorCode
	}

	errCodeDesRune := []rune(errCodeDes)
	if len(errCodeDesRune) > 64 {
		errCodeDes = string(errCodeDesRune[:64])
	}

	return respCode.ISO8583Code, errCodeDes, respCode.ErrorCode
}
