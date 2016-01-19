package cfca

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

// transformResp 转换应答内容
func transformResp(resp *BindingResponse, txCode string) (ret *model.BindingReturn) {

	// 打印渠道返回日志
	log.Infof("CFCA %s: %+v", txCode, resp)

	if resp == nil {
		return mongo.RespCodeColl.Get("000001")
	}
	ret = new(model.BindingReturn)
	// ret.ChanRespCode = resp.Head.Code
	// ret.ChanRespMsg = resp.Head.Message

	// 不成功的受理
	if flag := resp.Head.Code != correctCode; flag {
		ret.RespCode, ret.RespMsg = resp.Head.Code, resp.Head.Message
		return ret
	}

	// 标识是否业务失败
	var isFailed bool

AGAIN:
	// 成功受理的请求
	switch txCode {
	//根据交易类型处理结果
	//建立绑定关系、绑定关系查询
	case BindingCreate, BindingEnquiry:
		// ret.BindingId = resp.Body.TxSNBinding
		//10=绑定处理中 20=绑定失败 30=绑定成功 40=解绑成功
		switch resp.Body.Status {
		case "10":
			ret.RespCode = "000009"
		case "20":
			// 失败
			ret.RespCode = "100040"
			isFailed = true
		case "30":
			ret.RespCode = "000000"
		case "40":
			ret.RespCode = "100050"
		default:
			log.Errorf("渠道返回状态值(%d)错误，无法匹配。", resp.Body.Status)
			ret.RespCode = "000001"
		}
	//解除绑定关系
	case BindingRemove:
		//10=解绑处理中 20=解绑成功 30=解绑失败(等于已绑定)
		switch resp.Body.Status {
		case "10":
			ret.RespCode = "000009"
		case "20":
			ret.RespCode = "000000"
		case "30":
			// 解绑失败
			ret.RespCode = "100060"
			isFailed = true
		default:
			log.Errorf("渠道返回状态值(%d)错误，无法匹配。", resp.Body.Status)
			ret.RespCode = "000001"
		}
	//快捷支付、快捷支付查询
	case MerModePay, MarketModePay, MerModePayEnquiry, MarketModePayEnquiry:
		//10=处理中 20=支付成功 30=支付失败
		switch resp.Body.Status {
		case "10":
			ret.RespCode = "000009"
		case "20":
			ret.RespCode = "000000"
		case "30":
			// 支付失败
			ret.RespCode = "100070"
			isFailed = true
		default:
			log.Errorf("渠道返回状态值(%d)错误，无法匹配。", resp.Body.Status)
			ret.RespCode = "000001"
		}
	case MerModeRefund, MarketModeRefund:
		//都是受理成功
		ret.RespCode = "000000"
	case MerModeRefundEnquiry, MarketModeRefundEnquiry:
		//10=已受理 20=正在退款 30=退款成功 40=退款失败
		switch resp.Body.Status {
		case "10", "20", "30":
			ret.RespCode = "000000"
		case "40":
			ret.RespCode = "100080"
			isFailed = true
		default:
			log.Errorf("渠道返回状态值(%d)错误，无法匹配。", resp.Body.Status)
			ret.RespCode = "000001"
		}
	case MerModePayWithSMS, MarketModePayWithSMS:
		// 验证码状态
		switch resp.Body.VerifyStatus {
		case "40":
			txCode = MerModePay
			// 验证码通过，验证交易状态
			goto AGAIN
		case "20":
			ret.RespCode = "200171"
			isFailed = true
		case "30":
			ret.RespCode = "200172"
			isFailed = true
		}

	case MerModeSendSMS, MarketModeSendSMS, MarketPaySettlement:
		ret.RespCode = "000000"
	case TransChecking:
		ret.RespCode = "000000"
	}

	// 业务失败，在应答码前加44返回
	if isFailed {
		if resp.Body.ResponseCode != "" && resp.Body.ResponseMessage != "" {
			ret.RespCode, ret.RespMsg = "44"+resp.Body.ResponseCode, resp.Body.ResponseMessage
			return ret
		}
	}

	ret.RespMsg = mongo.RespCodeColl.GetMsg(ret.RespCode)

	return ret
}
