package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

var success = mongo.OffLineRespCd("SUCCESS")

// transform 支付宝返回报文处理
func transform(service string, alpResp *alpResponse) (*model.ScanPayResponse, error) {

	if alpResp.IsSuccess != "T" {
		ret := requestError(alpResp.Error)
		ret.ChanRespCode = alpResp.Error
		return ret, nil
	}

	ret := new(model.ScanPayResponse)
	// 成功返回参数
	alipay := alpResp.Response.Alipay

	// 中文长度限制
	r := []rune(alipay.DetailErrorDes)
	if len(r) > 64 {
		alipay.DetailErrorDes = string(r[:64])
	}

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
		// TODO
	}

	// 响应成功返回
	return ret, nil
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
		alipay.DetailErrorDes = success.ErrorDetail
	// 下单失败
	case "ORDER_FAIL":
		// ret.ChanRespCode = alipay.DetailErrorCode
	case "ORDER_SUCCESS_PAY_INPROCESS", "UNKNOWN", "ORDER_SUCCESS_PAY_FAIL":
		// ret.ChanRespCode = alipay.DetailErrorCode
		ret.ChannelOrderNum = alipay.TradeNo
	default:
		ret.ChanRespCode = alipay.ResultCode
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", createAndPay, alipay.ResultCode)
	}

	ret.ChanRespCode = alipay.ResultCode
	ret.ErrorDetail = alipay.DetailErrorDes
	ret.Respcd = createAndPayCd(ret.ChanRespCode)
}

// preCreateHandle 预下单处理
func preCreateHandle(ret *model.ScanPayResponse, alipay alpDetail) {

	switch alipay.ResultCode {
	case "SUCCESS":
		ret.QrCode = alipay.QrCode
		ret.ChanRespCode = alipay.ResultCode
		alipay.DetailErrorDes = success.ErrorDetail
	case "FAIL", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
	default:
		ret.ChanRespCode = alipay.ResultCode
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", preCreate, alipay.ResultCode)
	}

	ret.ErrorDetail = alipay.DetailErrorDes
	ret.Respcd = preCreateCd(ret.ChanRespCode)
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

	ret.Respcd = queryCd(ret.ChanRespCode)
	ret.ErrorDetail = mongo.OffLineCdCol[ret.Respcd]
}

// refundHandle 退款处理
func refundHandle(ret *model.ScanPayResponse, alipay alpDetail) {

	switch alipay.ResultCode {
	case "SUCCESS":
		ret.ChanRespCode = alipay.ResultCode
		alipay.DetailErrorDes = success.ErrorDetail
	case "FAIL", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", refund, alipay.ResultCode)
		ret.ChanRespCode = alipay.ResultCode
	}

	ret.ErrorDetail = alipay.DetailErrorDes
	ret.Respcd = refundCd(ret.ChanRespCode)
}

// cancelHandle 撤销处理
func cancelHandle(ret *model.ScanPayResponse, alipay alpDetail) {
	switch alipay.ResultCode {
	case "SUCCESS":
		ret.ChanRespCode = alipay.ResultCode
		ret.ChannelOrderNum = alipay.TradeNo
		alipay.DetailErrorDes = success.ErrorDetail
	case "FAIL", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", refund, alipay.ResultCode)
		ret.ChanRespCode = alipay.ResultCode
	}
	ret.ErrorDetail = alipay.DetailErrorDes
	ret.Respcd = cancelCd(ret.ChanRespCode)
}
