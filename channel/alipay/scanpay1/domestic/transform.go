package domestic

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

var (
	buscMap   map[string]string
	inprocess = mongo.ScanPayRespCol.Get("INPROCESS")
	success   = mongo.ScanPayRespCol.Get("SUCCESS")
)

func init() {
	buscMap = make(map[string]string)
	buscMap[createAndPay] = "createandpay"
	buscMap[preCreate] = "precreate"
	buscMap[refund] = "refund"
	buscMap[query] = "query"
	buscMap[cancel] = "cancel"
}

// transform 支付宝返回报文处理
func transform(service string, alpResp *alpResponse) (*model.ScanPayResponse, error) {

	ret := new(model.ScanPayResponse)
	if alpResp.IsSuccess != "T" {
		ret.ChanRespCode = alpResp.Error
		errorCodeMapping(ret, "", buscMap[service])
		return ret, nil
	}

	// 成功返回参数
	alipay := alpResp.Response.Alipay
	switch service {
	// 下单
	case createAndPay:
		createAndPayHandle(ret, alipay)
	case preCreate:
		preCreateHandle(ret, alipay)
	case query:
		queryHandle(ret, alipay)
	case refund:
		refundHandle(ret, alipay)
	case cancel:
		cancelHandle(ret, alipay)
	default:
		log.Errorf("unknown alp service: %s", service)
	}

	// 如果不成功
	if ret.Respcd == "" {
		errorCodeMapping(ret, alipay.DetailErrorDes, buscMap[service])
	}

	// 响应成功返回
	return ret, nil
}

// errorCodeMapping 错误码映射
func errorCodeMapping(ret *model.ScanPayResponse, errorDetail, service string) {
	spCode := mongo.ScanPayRespCol.GetByAlp(ret.ChanRespCode, service)
	ret.Respcd = spCode.ISO8583Code
	ret.ErrorCode = spCode.ErrorCode

	if spCode.IsUseISO || errorDetail == "" {
		ret.ErrorDetail = spCode.ISO8583Msg
		return
	}

	// 中文长度限制
	r := []rune(errorDetail)
	if len(r) > 64 {
		errorDetail = string(r[:64])
	}

	// 使用渠道应答
	log.Infof("use alipay errorDetail info: %s", errorDetail)
	ret.ErrorDetail = errorDetail
}

// createAndPayHandle 下单处理
func createAndPayHandle(ret *model.ScanPayResponse, alipay alpDetail) {
	switch alipay.ResultCode {
	case "ORDER_SUCCESS_PAY_SUCCESS":
		ret.ChannelOrderNum = alipay.TradeNo
		ret.ConsumerAccount = alipay.BuyerLogonId
		ret.ConsumerId = alipay.BuyerUserId
		// 计算折扣
		ret.MerDiscount, ret.ChcdDiscount = alipay.DisCount()
		ret.ChanRespCode = alipay.ResultCode
		ret.Respcd = success.ISO8583Code
		ret.ErrorDetail = success.ISO8583Msg
		ret.ErrorCode = success.ErrorCode
	case "ORDER_SUCCESS_PAY_INPROCESS":
		ret.ChanRespCode = alipay.ResultCode
		ret.Respcd = inprocess.ISO8583Code
		ret.ErrorDetail = inprocess.ISO8583Msg
		ret.ErrorCode = inprocess.ErrorCode
	// 下单失败
	case "ORDER_FAIL", "UNKNOWN", "ORDER_SUCCESS_PAY_FAIL":
		ret.ChanRespCode = alipay.DetailErrorCode
		ret.ChannelOrderNum = alipay.TradeNo
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", createAndPay, alipay.ResultCode)
	}
}

// preCreateHandle 预下单处理
func preCreateHandle(ret *model.ScanPayResponse, alipay alpDetail) {

	switch alipay.ResultCode {
	case "SUCCESS":
		ret.QrCode = alipay.QrCode
		ret.ChanRespCode = alipay.ResultCode
		// 预下单为支付中
		ret.Respcd = inprocess.ISO8583Code
		ret.ErrorDetail = inprocess.ISO8583Msg
		ret.ErrorCode = inprocess.ErrorCode
	case "FAIL", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
	default:
		ret.ChanRespCode = alipay.ResultCode
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", preCreate, alipay.ResultCode)
	}
}

// queryHandle 查询处理
func queryHandle(ret *model.ScanPayResponse, alipay alpDetail) {

	switch alipay.ResultCode {
	case "SUCCESS":
		ret.ChannelOrderNum = alipay.TradeNo
		ret.ConsumerAccount = alipay.BuyerLogonId
		ret.ConsumerId = alipay.BuyerUserId
		// 计算折扣
		ret.MerDiscount, ret.ChcdDiscount = alipay.DisCount()
		ret.ChanRespCode = alipay.TradeStatus
	case "FAIL", "PROCESS_EXCEPTION":
		ret.ChanRespCode = alipay.DetailErrorCode
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", query, alipay.ResultCode)
		ret.ChanRespCode = alipay.ResultCode
	}
}

// refundHandle 退款处理
func refundHandle(ret *model.ScanPayResponse, alipay alpDetail) {

	switch alipay.ResultCode {
	case "SUCCESS":
		ret.ChanRespCode = alipay.ResultCode
		ret.ChannelOrderNum = alipay.TradeNo
		ret.Respcd = success.ISO8583Code
		ret.ErrorDetail = success.ISO8583Msg
		ret.ErrorCode = success.ErrorCode
	case "FAIL", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", refund, alipay.ResultCode)
		ret.ChanRespCode = alipay.ResultCode
	}
}

// cancelHandle 撤销处理
func cancelHandle(ret *model.ScanPayResponse, alipay alpDetail) {
	switch alipay.ResultCode {
	case "SUCCESS":
		ret.ChanRespCode = alipay.ResultCode
		ret.ChannelOrderNum = alipay.TradeNo
		ret.Respcd = success.ISO8583Code
		ret.ErrorDetail = success.ISO8583Msg
		ret.ErrorCode = success.ErrorCode
	case "FAIL", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", refund, alipay.ResultCode)
		ret.ChanRespCode = alipay.ResultCode
	}
}
