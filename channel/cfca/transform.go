package cfca

import (
	"quickpay/model"
	"quickpay/mongo"
)

// transformResp 转换应答内容
func transformResp(resp *BindingResponse, txCode string) (ret *model.BindingReturn) {

	// 成功受理的请求
	if flag := resp.Head.Code == correctCode; flag {

		ret = &model.BindingReturn{}

		switch txCode {
		//根据交易类型处理结果
		//建立绑定关系
		case BindingCreateTxCode:
			ret.BindingId = resp.Body.TxSNBinding
			switch resp.Body.Status {
			case 10:
				ret.RespCode = "000009"
			case 30:
				ret.RespCode = "000000"
			default:
				ret.RespCode = "000001"
			}
			ret.RespMsg = resp.Body.ResponseMessage

		//绑定关系查询
		case BindingEnquiryTxCode:
			//10=绑定处理中 20=绑定失败 30=绑定成功 40=解绑成功
			switch resp.Body.Status {
			case 10:
				ret.RespCode = "000009"
			case 20:
				ret.RespCode = "100040"
			case 30:
				ret.RespCode = "000000"
			case 40:
				ret.RespCode = "000040"
			default:
				ret.RespCode = "000001"
			}
			ret.RespMsg = resp.Body.ResponseMessage

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
				ret.RespCode = "000001"
			}
			ret.RespMsg = resp.Body.ResponseMessage

		//快捷支付
		case BindingPaymentTxCode:
			//10=处理中 20=支付成功 30=支付失败
			switch resp.Body.Status {
			case 10:
				ret.RespCode = "000009"
			case 20:
				ret.RespCode = "000000"
			case 30:
				ret.RespCode = "100070"
			default:
				ret.RespCode = "000001"
			}
			ret.RespMsg = resp.Body.ResponseMessage
		}

		return
	}

	// 失败的请求
	// 查找对应关系
	rep := mongo.GetRespCodeByCfca(resp.Head.Code)
	ret = &model.BindingReturn{
		RespCode: rep.RespCode,
		RespMsg:  rep.RespMsg,
	}

	return
}
