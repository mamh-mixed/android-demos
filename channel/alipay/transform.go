package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// transform 支付宝返回报文处理
func transform(service string, alpResp *alpResponse, ret *model.ScanPayResponse) *model.ScanPayResponse {

	if ret == nil {
		ret = new(model.ScanPayResponse)
	}
	if alpResp.IsSuccess == "T" {

		// 成功返回参数
		alipay := alpResp.Response.Alipay
		ret.ErrorDetail = alipay.ResultCode

		switch service {
		// 下单
		case createAndPay:

			switch alipay.ResultCode {
			case "ORDER_SUCCESS_PAY_SUCCESS":
				ret.ChannelOrderNum = alipay.TradeNo
				ret.ConsumerAccount = alipay.BuyerLogonId
				ret.ConsumerId = alipay.BuyerUserId
				// 计算折扣
				ret.MerDiscount, ret.ChcdDiscount = alipay.DisCount()
				// ret.RespCode = "000000"
				ret.ChanRespCode = alipay.ResultCode
			// 下单失败
			case "ORDER_FAIL":
				// ret.RespCode = "100070"
				ret.ChanRespCode = alipay.DetailErrorCode
			case "ORDER_SUCCESS_PAY_INPROCESS", "UNKNOWN":
				// ret.RespCode = "000009"
				ret.ChanRespCode = alipay.DetailErrorCode
				ret.ChannelOrderNum = alipay.TradeNo
			case "ORDER_SUCCESS_PAY_FAIL":
				// ret.RespCode = "100070"
				ret.ChanRespCode = alipay.DetailErrorCode
				ret.ChannelOrderNum = alipay.TradeNo
			default:
				log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", createAndPay, alipay.ResultCode)
			}
			// get ret.Respcd by ResultCode
			// TODO check by ResultCode or by alipay.DetailErrorCode?
			ret.Respcd = createAndPayCd(alipay.ResultCode)

		case preCreate:
			// TODO
		case query:

			switch alipay.ResultCode {
			case "SUCCESS":
				ret.ChannelOrderNum = alipay.TradeNo
				ret.ConsumerAccount = alipay.BuyerLogonId
				ret.ConsumerId = alipay.BuyerUserId
				// 计算折扣
				ret.MerDiscount, ret.ChcdDiscount = alipay.DisCount()
				// ret.RespCode = "000000"
				ret.Respcd = queryCd(ret.Busicd, alipay.TradeStatus)
			case "FAIL", "PROCESS_EXCEPTION":
				ret.ErrorDetail = alipay.DetailErrorCode
				// ret.ChanRespCode = alipay.DetailErrorCode
				ret.Respcd = queryCd(ret.Busicd, ret.ErrorDetail)
			default:
				log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", query, alipay.ResultCode)
				ret.Respcd = queryCd(ret.Busicd, alipay.ResultCode)
			}

		case refund:
			// TODO
		default:
			// TODO
		}
	} else {
		ret.ErrorDetail = alpResp.Error

		// TODO get ret.Respcd by alpResp.Error
	}

	return ret
}
