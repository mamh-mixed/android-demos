package bindingpay

import (
	"quickpay/model"
	"regexp"
)

//建立绑定关系的时候验证请求报文
func bindingCreateRequestValidity(request model.BindingCreate) (ret *model.BindingReturn) {
	cardNum := request.AcctNum
	if request.BindingId == "" || request.AcctName == "" || request.AcctNum == "" || request.AcctType == "" {
		return model.NewBindingReturn("200050", "报文要素缺失")
	}
	if isUnionPayCard(cardNum) {
		//银联卡
		if request.IdentType == "" || request.IdentNum == "" || request.PhoneNum == "" || request.SendSmsId == "" || request.SmsCode == "" {
			return model.NewBindingReturn("200050", "报文要素缺失")
		}
		if matched, _ := regexp.MatchString(`^\d{1}$|^X$`, request.IdentType); !matched {
			return model.NewBindingReturn("200111", "证件类型不正确")
		}
		if matched, _ := regexp.MatchString(`^(13[0-9]|14[57]|15[0-9]|18[0-9])\d{8}$`, request.PhoneNum); !matched {
			return model.NewBindingReturn("200113", "手机号不正确")
		}
	} else {
		//外卡

	}
	if request.AcctType == "20" {
		// 贷记卡
		if request.ValidDate == "" || request.Cvv2 == "" {
			return model.NewBindingReturn("200050", "报文要素缺失")
		}
		if matched, _ := regexp.MatchString(`\d{2}(0[1-9]|1[1-2])`, request.ValidDate); !matched {
			return model.NewBindingReturn("200116", "信用卡有效期不正确")
		}
		if matched, _ := regexp.MatchString(`^\d{3}$`, request.Cvv2); !matched {
			return model.NewBindingReturn("200118", "CVV2不正确")
		}
	}

	return nil
}

// 移除绑定关系的时候验证请求报文
func bindingRemoveRequestValidity(in model.BindingRemove) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		return model.NewBindingReturn("200050", "报文要素缺失")
	}
	return nil
}

// 查询绑定关系的时候验证请求报文
func bindingEnquiryRequestValidity(be model.BindingEnquiry) (ret *model.BindingReturn) {
	if be.BindingId == "" {
		return model.NewBindingReturn("200050", "报文要素缺失")
	}
	return nil
}

// 绑定支付的请求报文验证
func bindingPaymentRequestValidity(in model.BindingPayment) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" || in.TransAmt == 0 || in.BindingId == "" {
		return model.NewBindingReturn("200050", "报文要素缺失")
	}

	if in.TransAmt < 0 {
		return model.NewBindingReturn("200124", "金额错误")
	}
	// 验证短信验证码是否填写
	if in.SendSmsId != "" && in.SmsCode == "" {
		return model.NewBindingReturn("200050", "报文要素缺失")
	}

	return nil
}

// 退款请求报文验证
func bindingRefundRequestValidity(in *model.BindingRefund) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" || in.OrigOrderNum == "" || in.TransAmt == 0 {
		return model.NewBindingReturn("200050", "报文要素缺失")
	}

	if in.TransAmt < 0 {
		return model.NewBindingReturn("200124", "金额错误")
	}

	return nil
}

// 无卡直接支付请求报文验证
func noTrackPaymentRequestValidity(in *model.NoTrackPayment) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "报文要素缺失---merOrderNum")
	}

	if in.TransAmt == 0 {
		return model.NewBindingReturn("200050", "报文要素缺失---transAmt")
	}

	if in.AcctName == "" {
		return model.NewBindingReturn("200050", "报文要素缺失---acctName")
	}

	if in.AcctNum == "" {
		return model.NewBindingReturn("200050", "报文要素缺失---acctNum")
	}

	if in.AcctType == "" {
		return model.NewBindingReturn("200050", "报文要素缺失---acctType")
	}

	if in.TransAmt < 0 {
		return model.NewBindingReturn("200180", "金额错误")
	}
	if matched, _ := regexp.MatchString(`^10$|^20$`, in.AcctType); !matched {
		return model.NewBindingReturn("200115", "账户类型不正确")
	}

	return nil
}

// todo 根据卡BIN验证是否是银联卡
func isUnionPayCard(cardNum string) bool {
	return true
}
