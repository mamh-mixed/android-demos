package core

import (
	"regexp"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

const (
	IDCardFlag = "0" // 0 表示身份证
)

// 银联卡必须验证证件和手机号
func UnionPayCardValidity(bc *model.BindingCreate) (ret *model.BindingReturn) {
	if bc.IdentType == "" {
		return model.NewBindingReturn("200050", "银联卡 identType 不能为空")
	}

	if bc.IdentNumDecrypt == "" {
		return model.NewBindingReturn("200050", "银联卡 identNum 不能为空")
	}

	if bc.PhoneNumDecrypt == "" {
		return model.NewBindingReturn("200050", "银联卡 phoneNum 不能为空")
	}

	if matched, _ := regexp.MatchString(`^[0-9]$|^X$`, bc.IdentType); !matched {
		return mongo.RespCodeColl.Get("200120")
	}
	// 判断格式，需要使用解密后的参数
	if matched, _ := regexp.MatchString(`^1[3|5|7|8|][0-9]{9}$`, bc.PhoneNumDecrypt); !matched {
		return mongo.RespCodeColl.Get("200130")
	}

	if matched, _ := regexp.MatchString(`^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}(x|X|[0-9])$`, bc.IdentNumDecrypt); !matched {
		return mongo.RespCodeColl.Get("200240")
	}
	return nil
}

// 银联卡必须验证证件和手机号
func UnionPayCardCommonValidity(identType, identNum, phoneNum string) (ret *model.BindingReturn) {
	if identType == "" {
		return model.NewBindingReturn("200050", "银联卡 identType 不能为空")
	}

	if identNum == "" {
		return model.NewBindingReturn("200050", "银联卡 identNum 不能为空")
	}

	if phoneNum == "" {
		return model.NewBindingReturn("200050", "银联卡 phoneNum 不能为空")
	}

	if matched, _ := regexp.MatchString(`^[0-9]$|^X$`, identType); !matched {
		return mongo.RespCodeColl.Get("200120")
	}
	// 判断格式，需要使用解密后的参数
	if matched, _ := regexp.MatchString(`^1[3|5|7|8|][0-9]{9}$`, phoneNum); !matched {
		return mongo.RespCodeColl.Get("200130")
	}

	// 如果是身份证
	if identType == IDCardFlag {
		matched, _ := regexp.MatchString(`^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}(x|X|[0-9])$`, identNum)
		if !matched {
			return mongo.RespCodeColl.Get("200240")
		}
	}

	return nil
}
