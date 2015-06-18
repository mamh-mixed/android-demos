package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// transform 支付宝返回报文处理
func transform(service string, alpResp *alpResponse, err error) *model.ScanPayResponse {

	// TODO check return error
	if err != nil {
		return mongo.OffLineRespCd("Time_out")
	}

	if alpResp.IsSuccess != "T" {
		return mongo.OffLineRespCd(alpResp.Error)
	}

	ret := new(model.ScanPayResponse)
	// 成功返回参数
	alipay := alpResp.Response.Alipay
	ret.ErrorDetail = alipay.ResultCode

	switch service {
	// 下单
	case createAndPay:
		createAndPayHandle(ret, alipay)
	case preCreate:
		preCreateHandle(ret, alipay)
	case query:
		queryHandle(ret, alipay)
	case refund:
		// TODO
	case cancel:
		// TODO
	default:
		// TODO
	}

	// 响应成功返回
	return ret
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
	// 下单失败
	case "ORDER_FAIL":
		ret.ChanRespCode = alipay.DetailErrorCode
	case "ORDER_SUCCESS_PAY_INPROCESS", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
		ret.ChannelOrderNum = alipay.TradeNo
	case "ORDER_SUCCESS_PAY_FAIL":
		ret.ChanRespCode = alipay.DetailErrorCode
		ret.ChannelOrderNum = alipay.TradeNo
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", createAndPay, alipay.ResultCode)
	}
	// get ret.Respcd by ResultCode
	// TODO check by ResultCode or by alipay.DetailErrorCode?
	ret.Respcd = createAndPayCd(alipay.ResultCode)
}

// preCreateHandle 预下单处理
func preCreateHandle(ret *model.ScanPayResponse, alipay alpDetail) {

	switch alipay.ResultCode {
	case "SUCCESS":
		ret.QrCode = alipay.QrCode
	case "FAIL", "UNKNOWN":
		ret.ChanRespCode = alipay.DetailErrorCode
		ret.ErrorDetail = alipay.DetailErrorCode
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", createAndPay, alipay.ResultCode)
	}
	ret.Respcd = preCreateCd(ret.ErrorDetail)
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
		ret.Respcd = queryCd(ret.Busicd, alipay.TradeStatus)
	case "FAIL", "PROCESS_EXCEPTION":
		ret.ErrorDetail = alipay.DetailErrorCode
		ret.Respcd = queryCd(ret.Busicd, ret.ErrorDetail)
	default:
		log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", query, alipay.ResultCode)
		ret.Respcd = queryCd(ret.Busicd, alipay.ResultCode)
	}
}
