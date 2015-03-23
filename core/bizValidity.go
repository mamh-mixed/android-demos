package core

import (
	"regexp"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

// 银联卡必须验证证件和手机号
func UnionPayCardValidity(bc *model.BindingCreate) (ret *model.BindingReturn) {
	if bc.IdentType == "" {
		return model.NewBindingReturn("200050", "银联卡 identType 不能为空")
	}

	if bc.IdentNum == "" {
		return model.NewBindingReturn("200050", "银联卡 identNum 不能为空")
	}

	if bc.PhoneNum == "" {
		return model.NewBindingReturn("200050", "银联卡 phoneNum 不能为空")
	}

	if matched, _ := regexp.MatchString(`^[0-9]$|^X$`, bc.IdentType); !matched {
		return mongo.RespCodeColl.Get("200120")
	}

	if matched, _ := regexp.MatchString(`^1[3|5|7|8|][0-9]{9}$`, bc.PhoneNum); !matched {
		return mongo.RespCodeColl.Get("200130")
	}

	return nil
}
