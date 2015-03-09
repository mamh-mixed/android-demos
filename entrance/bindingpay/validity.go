package bindingpay

import (
	"quickpay/model"
	"regexp"
)

//建立绑定关系的时候验证请求报文
func bindingCreateRequestValidity(request model.BindingCreate) (ret *model.BindingReturn) {
	cardNum := request.AcctNum
	if request.BindingId == "" || request.AcctName == "" || request.AcctNum == "" || request.AcctType == "" {
		return newBindingReturn("200050", "报文要素缺失")
	}
	if isUnionPayCard(cardNum) {
		//银联卡
		if request.IdentType == "" || request.IdentNum == "" || request.PhoneNum == "" || request.SendSmsId == "" || request.SmsCode == "" {
			return newBindingReturn("200050", "报文要素缺失")
		}
		if matched, _ := regexp.MatchString(`^\d{1}$|^X$`, request.IdentType); !matched {
			return newBindingReturn("200111", "证件类型不正确")
		}
		if matched, _ := regexp.MatchString(`^(13[0-9]|14[57]|15[0-9]|18[0-9])\d{8}$`, request.PhoneNum); !matched {
			return newBindingReturn("200113", "手机号不正确")
		}
	} else {
		//外卡

	}
	if request.AcctType == "20" {
		// 贷记卡
		if request.ValidDate == "" || request.Cvv2 == "" {
			return newBindingReturn("200050", "报文要素缺失")
		}
		if matched, _ := regexp.MatchString(`\d{2}(0[1-9]|1[1-2])`, request.ValidDate); !matched {
			return newBindingReturn("200116", "信用卡有效期不正确")
		}
		if matched, _ := regexp.MatchString(`^\d{3}$`, request.Cvv2); !matched {
			return newBindingReturn("200118", "CVV2不正确")
		}
	}

	return nil
}

// 移除绑定关系的时候验证请求报文
func bindingRemoveRequestValidity(in model.BindingRemove) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		return newBindingReturn("200050", "报文要素缺失")
	}
	return nil
}

// 查询绑定关系的时候验证请求报文
func bindingEnquiryRequestValidity(be model.BindingEnquiry) (ret *model.BindingReturn) {
	if be.BindingId == "" {
		return newBindingReturn("200050", "报文要素缺失")
	}
	return nil
}

// 绑定支付的请求报文验证
func bindingPaymentRequestValidity(in model.BindingPayment) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" || in.TransAmt == 0 || in.BindingId == "" {
		return newBindingReturn("200050", "报文要素缺失")
	}

	if in.TransAmt < 0 {
		return newBindingReturn("200124", "金额错误")
	}
	// 验证短信验证码是否填写
	if in.SendSmsId != "" && in.SmsCode == "" {
		return newBindingReturn("200050", "报文要素缺失")
	}

	return nil
}

// 退款请求报文验证
func bindingRefundRequestValidity(in *model.BindingRefund) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" || in.OrigOrderNum == "" || in.TransAmt == 0 {
		return newBindingReturn("200050", "报文要素缺失")
	}

	if in.TransAmt < 0 {
		return newBindingReturn("200124", "金额错误")
	}

	return nil
}

// 无卡直接支付请求报文验证
func noTrackPaymentRequestValidity(in *model.NoTrackPayment) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" {
		return newBindingReturn("200050", "报文要素缺失---merOrderNum")
	}

	if in.TransAmt == 0 {
		return newBindingReturn("200050", "报文要素缺失---transAmt")
	}

	if in.AcctName == "" {
		return newBindingReturn("200050", "报文要素缺失---acctName")
	}

	if in.AcctNum == "" {
		return newBindingReturn("200050", "报文要素缺失---acctNum")
	}

	if in.AcctType == "" {
		return newBindingReturn("200050", "报文要素缺失---acctType")
	}

	if in.TransAmt < 0 {
		return newBindingReturn("200180", "金额错误")
	}
	if matched, _ := regexp.MatchString(`^10$|^20$`, in.AcctType); !matched {
		return newBindingReturn("200115", "账户类型不正确")
	}

	return nil
}

// 生成BindingReturn的方法
func newBindingReturn(code, msg string) (ret *model.BindingReturn) {
	return &model.BindingReturn{
		RespCode: code,
		RespMsg:  msg,
	}
}

// todo 根据卡BIN验证是否是银联卡
func isUnionPayCard(cardNum string) bool {
	return true
}
