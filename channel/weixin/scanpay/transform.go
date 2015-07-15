package scanpay

import (
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type respCodeMap struct {
	respCode       string // 微信应答码
	respMsg        string // 微信应答码描述
	classification string // 分类
	iso8583Code    string // 8583 应答码
	iso8583Name    string // 8583 含义
}

var weixinRespCodeMap map[string]*respCodeMap

// 应答码
var (
	success      = mongo.OffLineRespCd("SUCCESS")
	cilError     = mongo.OffLineRespCd("SYSTEM_ERROR")
	unknownError = mongo.OffLineRespCd("UNKNOWN_ERROR")
)

func init() {
	weixinRespCodeMap = make(map[string]*respCodeMap)

	weixinRespCodeMap["SYSTEMERROR"] = &respCodeMap{"SYSTEMERROR", "接口返回错误", "微信问题", "91", "发卡方不能操作"}
	weixinRespCodeMap["INVALID_TRANSACTIONID"] = &respCodeMap{"INVALID_TRANSACTIONID", "无效transaction_id", "讯联问题", "96", "交换中心异常、失效"}
	weixinRespCodeMap["PARAM_ERROR"] = &respCodeMap{"PARAM_ERROR", "参数错误", "讯联问题", "96", "交换中心异常、失效"}
	weixinRespCodeMap["ORDERPAID"] = &respCodeMap{"ORDERPAID", "订单已支付", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["OUT_TRADE_NO_USED"] = &respCodeMap{"OUT_TRADE_NO_USED", "商户订单号重复", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["NOAUTH"] = &respCodeMap{"NOAUTH", "商户无权限", "交易权限问题", "58", "不允许终端进行的交易"}
	weixinRespCodeMap["AUTHCODEEXPIRE"] = &respCodeMap{"AUTHCODEEXPIRE", "条码已过期,请刷新再试", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["NOTENOUGH"] = &respCodeMap{"NOTENOUGH", "余额不足", "交易权限问题", "57", "不允许终端进行的交易"}
	weixinRespCodeMap["NOTSUPORTCARD"] = &respCodeMap{"NOTSUPORTCARD", "不支持卡类型", "交易权限问题", "58", "不允许终端进行的交易"}
	weixinRespCodeMap["ORDERCLOSED"] = &respCodeMap{"ORDERCLOSED", "订单已关闭", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["ORDERREVERSED"] = &respCodeMap{"ORDERREVERSED", "订单已撤销", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["BANKERROR"] = &respCodeMap{"BANKERROR", "银行系统异常", "微信问题", "91", "发卡方不能操作"}
	weixinRespCodeMap["USERPAYING"] = &respCodeMap{"USERPAYING", "用户支付中,需要输入密码", "没问题", "09（自定义）", "交易处理中（重试）"}
	weixinRespCodeMap["AUTH_CODE_ERROR"] = &respCodeMap{"AUTH_CODE_ERROR", "授权码参数错误", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["AUTH_CODE_INVALID"] = &respCodeMap{"AUTH_CODE_INVALID", "授权码检验错误", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["TRADE_STATE_ERROR"] = &respCodeMap{"TRADE_STATE_ERROR", "请重新发起（只在撤销交易时出现）", "微信问题", "01", "查发卡方"}
	weixinRespCodeMap["REFUND_FEE_INVALID"] = &respCodeMap{"REFUND_FEE_INVALID", "退款金额大于支付金额", "交易状态问题", "12", "无效交易"}
	weixinRespCodeMap["USERPAYING"] = &respCodeMap{"USERPAYING", "用户支付中", "没问题", "09", "交易处理中（重试）"}
	weixinRespCodeMap["CLOSED"] = &respCodeMap{"CLOSED", "已关闭", "没问题", "12", "交易关闭"}
	weixinRespCodeMap["SUCCESS"] = &respCodeMap{"SUCCESS", "支付成功", "没问题", "00", "交易成功"}
	weixinRespCodeMap["PAYERROR"] = &respCodeMap{"PAYERROR", "支付失败（其他原因，如银行返回失败）", "没问题", "01", "支付失败（其他原因，如银行返回失败）"}
	weixinRespCodeMap["NOTPAY"] = &respCodeMap{"NOTPAY", "未支付", "没问题", "12", "未支付"}
	weixinRespCodeMap["NOPAY"] = &respCodeMap{"NOPAY", "未支付（确认支付超时）", "没问题", "12", "未支付（确认支付超时）"}
	weixinRespCodeMap["REVOKED"] = &respCodeMap{"REVOKED", "已撤销", "没问题", "12", "已撤销"}
	weixinRespCodeMap["REFUND"] = &respCodeMap{"REFUND", "转入退款", "没问题", "12", "转入退款"}
}

func transform(returnCode, returnMsg, resultCode, errCode, errCodeDes string) (status, msg string) {

	// 描述长度限制
	returnMsgRune := []rune(returnMsg)
	if len(returnMsgRune) > 64 {
		returnMsg = string(returnMsgRune[:64])
	}
	errCodeDesRune := []rune(errCodeDes)
	if len(errCodeDesRune) > 64 {
		errCodeDes = string(errCodeDesRune[:64])
	}

	// 如果通信失败，返回错误代码 96，并返回错误原因
	if returnCode != "SUCCESS" {
		return cilError.Respcd, returnMsg
	}

	// 如果结果正确，返回代码 00，消息为交易成功
	if resultCode == "SUCCESS" {
		return success.Respcd, success.ErrorDetail
	}

	if m, ok := weixinRespCodeMap[errCode]; ok {
		// 使用微信业务错误描述
		if errCodeDes != "" {
			return m.iso8583Code, errCodeDes
		}
		return m.iso8583Code, mongo.OffLineCdCol[m.iso8583Code]
	}

	log.Errorf("unknown weixin error code `%s`", errCode)
	return unknownError.Respcd, unknownError.ErrorDetail
}
