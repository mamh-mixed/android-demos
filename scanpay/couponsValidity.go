package scanpay

import (
	"regexp"
	"strconv"

	"github.com/CardInfoLink/quickpay/model"
)

const (
	veriTime     = "veriTime"
	origOrderNum = "origOrderNum"
	payType      = "payType"
)

// validatePurchaseCoupons 验证卡券核销的参数
func validatePurchaseCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)

	}
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	// 验证格式
	if matched, err := validateChcd(req); !matched {
		return err
	}
	if matched, err := validateVeriTime(req); !matched {
		return err
	}
	return
}

// validatePurchaseActCoupons 验证‘刷卡活动券验证’的参数
func validatePurchaseActCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.PayType == "":
		return fieldEmptyError(payType)

	}

	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}

	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}

	if matched, err := validateChcd(req); !matched {
		return err
	}
	if matched, err := validateVeriTime(req); !matched {
		return err
	}
	if matched, err := validateCouponTxamt(req); !matched {
		return err
	}
	if matched, err := validatePayType(req); !matched {
		return err
	}

	return
}

// validateQueryPurchaseCoupons 验证‘电子券验证结果查询’的参数
func validateQueryPurchaseCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	}

	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}

	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}

	if matched, err := validateChcd(req); !matched {
		return err
	}
	if matched, err := validateVeriTime(req); !matched {
		return err
	}
	// if req.Txamt != "" {
	// 	if matched, err := validateTxamt(req); !matched {
	// 		return err
	// 	}
	// }
	// if req.PayType != "" {
	// 	if matched, err := validatePayType(req); !matched {
	// 		return err
	// 	}
	// }
	return
}

// validateUndoPurchaseActCoupons 验证‘卡券撤销’的参数
func validateUndoPurchaseActCoupons(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.Txndir == "":
		return fieldEmptyError(txndir)
	case req.Busicd == "":
		return fieldEmptyError(buiscd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Terminalid == "":
		return fieldEmptyError(terminalid)
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	case req.Sign == "":
		return fieldEmptyError(sign)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	}

	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}

	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}

	if matched, err := validateChcd(req); !matched {
		return err
	}
	return
}

// validatePayType 验证支付方式
func validatePayType(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {
	intPayType, err := strconv.Atoi(req.PayType)
	if err != nil {
		return false, fieldFormatError(payType)
	}
	req.IntPayType = intPayType
	return true, nil
}

// validateVeriTime 验证核销次数
func validateVeriTime(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {
	if req.VeriTime != "" {
		intVeriTime, err := strconv.Atoi(req.VeriTime)
		if err != nil {
			return false, fieldFormatError(veriTime)
		}
		if intVeriTime <= 0 {
			return false, fieldFormatError(veriTime)
		}
		req.IntVeriTime = intVeriTime
	} else {
		req.IntVeriTime = 1
	}
	return true, nil
}

// validateChcd 验证渠道代码
func validateChcd(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {
	if req.Chcd != "" && req.Chcd != "ULIVE" {
		return false, fieldContentError(chcd)
	}
	return true, nil
}

// validateTxamt 验证金额
func validateCouponTxamt(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {

	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return false, fieldFormatError(txamt)
	}

	// 转换金额
	toInt, err := strconv.ParseInt(req.Txamt, 10, 64)
	if err != nil {
		return false, fieldFormatError(txamt)
	}

	// 金额范围
	// || toInt > maxTxamt 不限制金额，就按12位最大值来
	// if toInt == minTxamt {
	// 	return false, fieldFormatError(txamt)
	// }

	req.IntTxamt = toInt
	return true, nil
}
