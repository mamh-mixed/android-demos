package alipay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// barcodePayTransform 下单接口返回报文处理
func barcodePayTransform(alpResp *alpResponse) *model.QrCodePayResponse {

	ret := new(model.QrCodePayResponse)
	if alpResp.IsSuccess == "T" {
		alipay := alpResp.Response.Alipay
		ret.ChanRespCode = alipay.ResultCode
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
			log.Errorf("渠道返回状态值(%s)错误，无法匹配。", alipay.ResultCode)
		}
		// TODO get ret.Respcd by ResultCode

	} else {
		ret.ChanRespCode = alpResp.Error
		ret.ErrorDetail = alpResp.Error

		// TODO get ret.Respcd by alpResp.Error
	}

	return ret
}
