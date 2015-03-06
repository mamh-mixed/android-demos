package validity

import (
	"errors"
	"quickpay/domain"
	"quickpay/model"
	"regexp"
)

//建立绑定关系的时候验证请求报文
func BindingCreateRequestValidity(request domain.BindingCreateRequest) (string, error) {
	cardNum := request.AcctNum
	if request.BindingId == "" || request.AcctName == "" || request.AcctNum == "" || request.AcctType == "" {
		return "200050", errors.New("报文要素缺失")
	}
	if isUnionPayCard(cardNum) {
		//银联卡
		if request.IdentType == "" || request.IdentNum == "" || request.PhoneNum == "" || request.SendSmsId == "" || request.SmsCode == "" {
			return "200050", errors.New("报文要素缺失")
		}
		if matched, _ := regexp.MatchString(`^\d{1}$|^X$`, request.IdentType); !matched {
			return "200120", errors.New("证件类型有误")
		}
		if matched, _ := regexp.MatchString(`^(13[0-9]|14[57]|15[0-9]|18[0-9])\d{8}$`, request.PhoneNum); !matched {
			return "200130", errors.New("手机号有误")
		}
	} else {
		//外卡

	}
	if request.AcctType == "20" {
		// 贷记卡
		if request.ValidDate == "" || request.Cvv2 == "" {
			return "200050", errors.New("报文要素缺失")
		}
		if matched, _ := regexp.MatchString(`\d{2}(0[1-9]|1[1-2])`, request.ValidDate); !matched {
			return "200140", errors.New("卡片有效期有误")
		}
		if matched, _ := regexp.MatchString(`^\d{3}$`, request.Cvv2); !matched {
			return "200150", errors.New("CVV2有误")
		}
	}

	return "00", nil
}

// 移除绑定关系的时候验证请求报文
func BindingRemoveRequestValidity(in model.BindingRemoveIn) (string, error) {
	if in.BindingId == "" {
		return "200050", errors.New("报文要素缺失")
	} else {
		return "00", nil
	}
}

// 查询绑定关系的时候验证请求报文
func BindingEnquiryRequestValidity(in model.BindingEnquiryIn) (string, error) {
	if in.BindingId == "" {
		return "200050", errors.New("报文要素缺失")
	} else {
		return "00", nil
	}
}

// 绑定支付的请求报文验证
func BindingPaymentRequestValidity(in model.BindingPaymentIn) (string, error) {
	if in.MerOrderNum == "" || in.TransAmt == 0 || in.BindingId == "" {
		return "200050", errors.New("报文要素缺失")
	}

	if in.TransAmt < 0 {
		return "200180", errors.New("金额有误")
	}
	// 短信验证码
	if in.SendSmsId != "" && in.SmsCode == "" {
		return "200050", errors.New("报文要素缺失")
	}

	return "00", nil
}

// todo 根据卡BIN验证是否是银联卡
func isUnionPayCard(cardNum string) bool {
	return true
}
