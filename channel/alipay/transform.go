package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// transform 支付宝返回报文处理
func transform(service string, alpResp *alpResponse) *model.ScanPayResponse {

	ret := new(model.ScanPayResponse)
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
				ret.RespCode = "000000"
			// 下单失败
			case "ORDER_FAIL":
				ret.RespCode = "100070"
			case "ORDER_SUCCESS_PAY_INPROCESS", "UNKNOWN":
				ret.RespCode = "000009"
				ret.ChannelOrderNum = alipay.TradeNo
			case "ORDER_SUCCESS_PAY_FAIL":
				ret.RespCode = "100070"
				ret.ChannelOrderNum = alipay.TradeNo
			default:
				log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", createAndPay, alipay.ResultCode)
			}
			// TODO get ret.Respcd by ResultCode
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
				ret.RespCode = "000000"
			case "FAIL", "PROCESS_EXCEPTION":
				ret.ErrorDetail = alipay.DetailErrorCode
			default:
				log.Errorf("支付宝服务(%s),返回状态值(%s)错误，无法匹配。", query, alipay.ResultCode)
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
