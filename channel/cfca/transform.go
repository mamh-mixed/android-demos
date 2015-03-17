package cfca

import (
	"quickpay/model"
	"quickpay/mongo"

	"github.com/omigo/g"
)

// transformResp 转换应答内容
func transformResp(resp *BindingResponse, txCode string) (ret *model.BindingReturn) {

	// default
	ret = model.NewBindingReturn("000001", "系统内部错误")
	if resp == nil {
		return
	}
	ret = &model.BindingReturn{}
	ret.ChanRespCode = resp.Head.Code
	ret.ChanRespMsg = resp.Head.Message
	// 不成功的受理
	if flag := resp.Head.Code != correctCode; flag {
		respObject := mongo.RespCodeColl.GetByCfca(resp.Head.Code)
		ret.RespCode = respObject.RespCode
		ret.RespMsg = respObject.RespMsg
		return
	}
	// 成功受理的请求
	switch txCode {
	//根据交易类型处理结果
	//建立绑定关系、绑定关系查询
	case BindingCreateTxCode, BindingEnquiryTxCode:
		// ret.BindingId = resp.Body.TxSNBinding
		//10=绑定处理中 20=绑定失败 30=绑定成功 40=解绑成功
		switch resp.Body.Status {
		case 10:
			ret.RespCode = "000009"
		case 20:
			ret.RespCode = "100040"
		case 30:
			ret.RespCode = "000000"
		case 40:
			ret.RespCode = "100050"
		default:
			g.Error("渠道返回状态值(%d)错误，无法匹配。", resp.Body.Status)
			ret.RespCode = "000001"
		}
	//解除绑定关系
	case BindingRemoveTxCode:
		//10=解绑处理中 20=解绑成功 30=解绑失败(等于已绑定)
		switch resp.Body.Status {
		case 10:
			ret.RespCode = "000009"
		case 20:
			ret.RespCode = "000000"
		case 30:
			ret.RespCode = "100060"
		default:
			g.Error("渠道返回状态值(%d)错误，无法匹配。", resp.Body.Status)
			ret.RespCode = "000001"
		}
	//快捷支付、快捷支付查询
	case BindingPaymentTxCode, PaymentEnquiryTxCode:
		//10=处理中 20=支付成功 30=支付失败
		switch resp.Body.Status {
		case 10:
			ret.RespCode = "000009"
		case 20:
			ret.RespCode = "000000"
		case 30:
			ret.RespCode = "100070"
		default:
			g.Error("渠道返回状态值(%d)错误，无法匹配。", resp.Body.Status)
			ret.RespCode = "000001"
		}
	}
	ret.RespMsg = mongo.RespCodeColl.GetMsg(ret.RespCode)
	g.Debug("resp message %+v", ret)

	return
}
