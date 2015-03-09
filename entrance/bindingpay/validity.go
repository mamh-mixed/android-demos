package bindingpay

import (
	"quickpay/model"
	"regexp"
)

//建立绑定关系的时候验证请求报文
func bindingCreateRequestValidity(request model.BindingCreate) (ret *model.BindingReturn) {
	cardNum := request.AcctNum
	if request.BindingId == "" || request.AcctName == "" || request.AcctNum == "" || request.AcctType == "" {
		ret = &model.BindingReturn{
			RespCode: "200050",
			RespMsg:  "报文要素缺失",
		}
		return ret
	}
	if isUnionPayCard(cardNum) {
		//银联卡
		if request.IdentType == "" || request.IdentNum == "" || request.PhoneNum == "" || request.SendSmsId == "" || request.SmsCode == "" {
			ret = &model.BindingReturn{
				RespCode: "200050",
				RespMsg:  "报文要素缺失",
			}
			return ret
		}
		if matched, _ := regexp.MatchString(`^\d{1}$|^X$`, request.IdentType); !matched {
			ret = &model.BindingReturn{
				RespCode: "200120",
				RespMsg:  "证件类型有误",
			}
			return ret
		}
		if matched, _ := regexp.MatchString(`^(13[0-9]|14[57]|15[0-9]|18[0-9])\d{8}$`, request.PhoneNum); !matched {
			ret = &model.BindingReturn{
				RespCode: "200130",
				RespMsg:  "手机号有误",
			}
			return ret
		}
	} else {
		//外卡

	}
	if request.AcctType == "20" {
		// 贷记卡
		if request.ValidDate == "" || request.Cvv2 == "" {
			ret = &model.BindingReturn{
				RespCode: "200050",
				RespMsg:  "报文要素缺失",
			}
			return ret
		}
		if matched, _ := regexp.MatchString(`\d{2}(0[1-9]|1[1-2])`, request.ValidDate); !matched {
			ret = &model.BindingReturn{
				RespCode: "200140",
				RespMsg:  "卡片有效期有误",
			}
			return ret
		}
		if matched, _ := regexp.MatchString(`^\d{3}$`, request.Cvv2); !matched {
			ret = &model.BindingReturn{
				RespCode: "200150",
				RespMsg:  "CVV2有误",
			}
			return ret
		}
	}

	return nil
}

// 移除绑定关系的时候验证请求报文
func bindingRemoveRequestValidity(in model.BindingRemove) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		ret = &model.BindingReturn{
			RespCode: "200050",
			RespMsg:  "报文要素缺失",
		}
		return ret
	}
	return nil
}

// 查询绑定关系的时候验证请求报文
func bindingEnquiryRequestValidity(be model.BindingEnquiry) (ret *model.BindingReturn) {
	if be.BindingId == "" {
		ret = &model.BindingReturn{
			RespCode: "200050",
			RespMsg:  "报文要素缺失",
		}
		return ret
	}
	return nil
}

// 绑定支付的请求报文验证
func bindingPaymentRequestValidity(in model.BindingPayment) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" || in.TransAmt == 0 || in.BindingId == "" {
		ret = &model.BindingReturn{
			RespCode: "200050",
			RespMsg:  "报文要素缺失",
		}
		return ret
	}

	if in.TransAmt < 0 {
		ret = &model.BindingReturn{
			RespCode: "200180",
			RespMsg:  "金额有误",
		}
		return ret
	}
	// 短信验证码
	if in.SendSmsId != "" && in.SmsCode == "" {
		ret = &model.BindingReturn{
			RespCode: "200050",
			RespMsg:  "报文要素缺失",
		}
		return ret
	}

	return nil
}

// 退款请求报文验证
func bindingRefundRequestValidity(in *model.BindingRefund) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" || in.OrigOrderNum == "" || in.TransAmt == 0 {
		ret = &model.BindingReturn{
			RespCode: "200050",
			RespMsg:  "报文要素缺失",
		}
		return ret
	}

	if in.TransAmt < 0 {
		ret = &model.BindingReturn{
			RespCode: "200190",
			RespMsg:  "退款金额有误",
		}
		return ret
	}

	return nil
}

// todo 根据卡BIN验证是否是银联卡
func isUnionPayCard(cardNum string) bool {
	return true
}
