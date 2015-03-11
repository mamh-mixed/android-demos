package core

import (
	"quickpay/model"
	"regexp"
)

// 银联卡必须验证证件和手机号
func UnionPayCardValidity(bc *model.BindingCreate) (ret *model.BindingReturn) {
	if bc.IdentType == "" {
		return model.NewBindingReturn("200050", "字段 identType 不能为空")
	}

	if bc.IdentNum == "" {
		return model.NewBindingReturn("200050", "字段 identNum 不能为空")
	}

	if bc.PhoneNum == "" {
		return model.NewBindingReturn("200050", "字段 phoneNum 不能为空")
	}

	if matched, _ := regexp.MatchString(`^[0-9]$|^X$`, bc.IdentType); !matched {
		return model.NewBindingReturn("200111", "证件类型不正确")
	}

	if matched, _ := regexp.MatchString(`^1[3|5|7|8|][0-9]{9}$`, bc.PhoneNum); !matched {
		return model.NewBindingReturn("200114", "手机号码格式错误")
	}

	return nil
}
