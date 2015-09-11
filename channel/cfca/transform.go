package cfca

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// transformResp 转换应答内容
func transformResp(resp *BindingResponse, txCode string) (ret *model.BindingReturn) {

	// 打印渠道返回日志
	log.Infof("CFCA %s: %+v", txCode, resp)

	// default
	ret = mongo.RespCodeColl.Get("000001")
	if resp == nil {
		return
	}
	ret = new(model.BindingReturn)
	ret.ChanRespCode = resp.Head.Code
	ret.ChanRespMsg = resp.Head.Message
	// 不成功的受理
	if flag := resp.Head.Code != correctCode; flag {
		respObject := mongo.RespCodeColl.GetByCfca(resp.Head.Code)
		ret.RespCode = respObject.RespCode
		ret.RespMsg = respObject.RespMsg
		if ret.RespCode == "" {
			//系统外部错误
			log.Errorf("找不到系统对应的中金应答码:(%s)", resp.Head.Code)
			ret.RespCode = "000002"
			ret.RespMsg = mongo.RespCodeColl.GetMsg(ret.RespCode)
		}
		return
	}

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
			ret.RespCode = "100040"
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
			ret.RespCode = "100060"
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
			ret.RespCode = "100070"
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
			// TODO:验证码超时
			ret.RespCode = "400000"
		case "30":
			// TODO:验证未通过
			ret.RespCode = "400001"
		}

	case MerModeSendSMS, MarketModeSendSMS, MarketPaySettlement:
		ret.RespCode = "000000"
	case TransChecking:
		ret.RespCode = "000000"
	}

	ret.RespMsg = mongo.RespCodeColl.GetMsg(ret.RespCode)
	log.Debugf("resp message %+v", ret)

	return
}
