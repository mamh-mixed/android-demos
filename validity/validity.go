package validity

import (
	"errors"
	"quickpay/domain"
	"regexp"
)

func BindingCreateRequestValidity(request domain.BindingCreateRequest) (string, error) {
	cardNum := request.AcctNum
	if request.MerBindingId == "" || request.AcctName == "" || request.AcctNum == "" || request.AcctType == "" {
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

func isUnionPayCard(cardNum string) bool {
	return true
}
