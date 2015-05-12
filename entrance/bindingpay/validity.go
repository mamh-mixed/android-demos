package bindingpay

import (
	"regexp"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

// validateBindingCreate 建立绑定关系的时候验证请求报文
func validateBindingCreate(request *model.BindingCreate) (ret *model.BindingReturn) {
	if request.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}

	if request.AcctNameDecrypt == "" {
		return model.NewBindingReturn("200050", "字段 acctName 不能为空")
	}

	if request.AcctNumDecrypt == "" {
		return model.NewBindingReturn("200050", "字段 acctNum 不能为空")
	}

	if request.AcctType == "" {
		return model.NewBindingReturn("200050", "字段 acctType 不能为空")
	}

	if request.AcctType == "20" {
		// 贷记卡
		if request.ValidDateDecrypt == "" {
			return model.NewBindingReturn("200050", "字段 validDate 不能为空")
		}

		if request.Cvv2Decrypt == "" {
			return model.NewBindingReturn("200050", "字段 cvv2 不能为空")
		}

		// 判断格式，需要使用解密后的参数
		if matched, _ := regexp.MatchString(`^\d{2}(0[1-9]|1[1-2])$`, request.ValidDateDecrypt); !matched {
			return mongo.RespCodeColl.Get("200140")
		}

		if matched, _ := regexp.MatchString(`^\d{3}$`, request.Cvv2Decrypt); !matched {
			return mongo.RespCodeColl.Get("200150")
		}
	}

	return nil
}

// validateBindingRemove 移除绑定关系的时候验证请求报文
func validateBindingRemove(in *model.BindingRemove) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}
	return nil
}

// validateBindingEnquiry 查询绑定关系的时候验证请求报文
func validateBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	if be.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}
	return nil
}

// validateBindingPayment 绑定支付的请求报文验证
func validateBindingPayment(in *model.BindingPayment) (ret *model.BindingReturn) {
	if in.BindingId == "" {
		return model.NewBindingReturn("200050", "字段 bindingId 不能为空")
	}

	if in.TransAmt == 0 {
		return model.NewBindingReturn("200050", "字段 transAmt 不能为空")
	}

	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}

	if in.TransAmt < 0 {
		return mongo.RespCodeColl.Get("200180")
	}
	// 验证短信验证码是否填写
	if in.SendSmsId != "" && in.SmsCode == "" {
		return model.NewBindingReturn("200050", "字段 smsCode 不能为空")
	}

	return nil
}

// validateBindingRefund 退款请求报文验证
func validateBindingRefund(in *model.BindingRefund) (ret *model.BindingReturn) {
	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}

	if in.OrigOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 origOrderNum 不能为空")
	}

	if in.TransAmt == 0 {
		return model.NewBindingReturn("200050", "字段 transAmt 不能为空")
	}

	if in.TransAmt < 0 {
		return mongo.RespCodeColl.Get("200180")
	}

	return nil
}

// validateOrderEnquiry 订单查询报文验证
func validateOrderEnquiry(in *model.OrderEnquiry) (ret *model.BindingReturn) {
	if in.OrigOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 origOrderNum 不能为空")
	}
	if matched, _ := regexp.MatchString(`^[1|0]?$`, in.ShowOrigInfo); !matched {
		//TODO check respCode
		return model.NewBindingReturn("200050", "字段 showOrigInfo 取值错误")
	}
	return
}

// validateBillingSummary 交易对账汇总验证
func validateBillingSummary(in *model.BillingSummary) (ret *model.BindingReturn) {
	if matched, _ := regexp.MatchString(`^[1-2][0-9][0-9][0-9]-(0[1-9]|1[0-2])-[0-3]{0,1}[0-9]$`, in.SettDate); !matched {
		return model.NewBindingReturn("200200", "日期 SettDate 格式错误")
	}
	return
}

// validateBillingSummary 交易对账汇总验证
func validateBillingDetails(in *model.BillingDetails) (ret *model.BindingReturn) {
	if matched, _ := regexp.MatchString(`^[1-2][0-9][0-9][0-9]-(0[1-9]|1[0-2])-[0-3]{0,1}[0-9]$`, in.SettDate); !matched {
		return model.NewBindingReturn("200200", "日期 SettDate 格式错误")
	}
	if len(in.NextOrderNum) > 32 {
		return model.NewBindingReturn("200080", "订单号 NextOrderNum 不正确")
	}
	return
}

// validateNoTrackPayment 无卡直接支付请求报文验证
func validateNoTrackPayment(in *model.NoTrackPayment) (ret *model.BindingReturn) {
	if in.TransType == "" {
		return model.NewBindingReturn("200050", "字段 transType 不能为空")
	}

	if in.MerOrderNum == "" {
		return model.NewBindingReturn("200050", "字段 merOrderNum 不能为空")
	}

	if in.TransAmt == 0 {
		return model.NewBindingReturn("200050", "字段 transAmt 不能为空")
	}

	if in.AcctName == "" {
		return model.NewBindingReturn("200050", "字段 acctName 不能为空")
	}

	if in.AcctNum == "" {
		return model.NewBindingReturn("200050", "字段 acctNum 不能为空")
	}

	if in.AcctType == "" {
		return model.NewBindingReturn("200050", "字段 acctType 不能为空")
	}

	if in.TransAmt < 0 {
		return mongo.RespCodeColl.Get("200180")
	}
	if matched, _ := regexp.MatchString(`^10$|^20$`, in.AcctType); !matched {
		return mongo.RespCodeColl.Get("200230")
	}

	if in.AcctType == "20" {
		// 贷记卡
		if in.ValidDate == "" {
			return model.NewBindingReturn("200050", "字段 validDate 不能为空")
		}

		if in.Cvv2 == "" {
			return model.NewBindingReturn("200050", "字段 cvv2 不能为空")
		}

		// 判断格式，需要使用解密后的参数
		if matched, _ := regexp.MatchString(`^\d{2}(0[1-9]|1[1-2])$`, in.ValidDateDecrypt); !matched {
			return mongo.RespCodeColl.Get("200140")
		}

		if matched, _ := regexp.MatchString(`^\d{3}$`, in.Cvv2Decrypt); !matched {
			return mongo.RespCodeColl.Get("200150")
		}
	}

	if in.CurrCode != "" {
		// 判断交易币种格式
		if matched, _ := regexp.MatchString(`^\d{3}$`, in.CurrCode); !matched {
			return mongo.RespCodeColl.Get("200251")
		}
	}

	if matched, _ := regexp.MatchString(`^SALE|AUTH$`, in.TransType); !matched {
		return mongo.RespCodeColl.Get("100030")
	}

	return nil
}
