package scanpay

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
)

const (
	minTxamt = 0
	maxTxamt = 1e10 - 1
)

// fieldName
const (
	curr       = "currency"
	txamt      = "txamt"
	orderNum   = "orderNum 或 origOrderNum"
	agentCode  = "inscd"
	mchntid    = "mchntid"
	scanCodeId = "scanCodeId"
	chcd       = "chcd"
	goodsInfo  = "goodsInfo"
	buiscd     = "busicd"
	terminalid = "terminalid"
	openId     = "openid"
	checkName  = "checkName"
	desc       = "desc"
	userName   = "userName"
	txndir     = "txndir"
	sign       = "sign"
)

var (
	alipayCurrency map[string]int
	success        = mongo.ScanPayRespCol.Get("SUCCESS")
	emptyError     = mongo.ScanPayRespCol.Get("DATA_EMPTY_ERROR")
	formatError    = mongo.ScanPayRespCol.Get("DATA_FORMAT_ERROR")
	contentError   = mongo.ScanPayRespCol.Get("DATA_CONTENT_ERROR")
)

// validateBarcodePay 验证扫码下单的参数
func validateBarcodePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
		// case req.Currency == "":
		// 	return fieldEmptyError(curr)
	}

	// 验证格式
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, _ := regexp.MatchString(`^\d{14,24}$`, req.ScanCodeId); !matched {
		return fieldFormatError(scanCodeId)
	}
	if matched, err := validateTxamt(req); !matched {
		return err
	}
	if matched, err := validateGoodsInfo(req.GoodsInfo); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	if req.Currency != "" {
		if matched, _ := regexp.MatchString(`^[A-Z]{3}$`, req.Currency); !matched {
			return fieldFormatError(curr)
		}
	}

	return
}

// validateQrCodeOfflinePay 验证预下单的参数
func validateQrCodeOfflinePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Chcd == "":
		return fieldEmptyError(chcd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
		// case req.Currency == "":
		// 	return fieldEmptyError(curr)
	}

	// 验证格式
	if req.Chcd != "WXP" && req.Chcd != "ALP" {
		return fieldContentError(chcd)
	}
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateTxamt(req); !matched {
		return err
	}
	if matched, err := validateGoodsInfo(req.GoodsInfo); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}

	if req.Currency != "" {
		if matched, _ := regexp.MatchString(`^[A-Z]{3}$`, req.Currency); !matched {
			return fieldFormatError(curr)
		}
	}

	if req.TimeExpire != "" {
		if matched, err := validateTimeExpire(req.TimeExpire); !matched {
			return err
		}
	}

	return
}

// validateEnquiry 验证查询接口的参数
func validateEnquiry(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.OrigOrderNum == "":
		return fieldEmptyError(orderNum)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	}

	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrigOrderNum); !matched {
		return err
	}

	return
}

// validateRefund 验证退款接口的参数
func validateRefund(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.OrigOrderNum == "":
		return fieldEmptyError(orderNum)
	}

	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateTxamt(req); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrigOrderNum); !matched {
		return err
	}

	return
}

// validateCancel 验证撤销接口参数
func validateCancel(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.OrigOrderNum == "":
		return fieldEmptyError(orderNum)
	}

	// 验证格式
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrigOrderNum); !matched {
		return err
	}

	return
}

// validateCancel 验证关闭订单接口参数
func validateClose(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.OrigOrderNum == "":
		return fieldEmptyError(orderNum)
	}

	// 验证格式
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrigOrderNum); !matched {
		return err
	}

	return
}

// validateEnterprisePay 验证企业付款接口参数
func validateEnterprisePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Chcd == "":
		return fieldEmptyError(chcd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.OpenId == "":
		return fieldEmptyError(openId)
	case req.CheckName == "":
		return fieldEmptyError(checkName)
	case req.Desc == "":
		return fieldEmptyError(desc)
	}

	// 验证格式
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	if matched, err := validateUserName(req); !matched {
		return err
	}
	if matched, err := validateTxamt(req); !matched {
		return err
	}
	if len(req.OpenId) > 64 {
		return fieldFormatError(openId)
	}

	return
}

// validatePublicPay 验证公众号支付接口参数
func validatePublicPay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	var needUserInfo = "needUserInfo"

	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Chcd == "":
		return fieldEmptyError(chcd)
	case req.AgentCode == "":
		return fieldEmptyError(agentCode)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.NeedUserInfo == "":
		return fieldEmptyError(needUserInfo)
	}

	// 验证格式
	if matched, err := validateMchntid(req.Mchntid); !matched {
		return err
	}
	if matched, err := validateOrderNum(req.OrderNum); !matched {
		return err
	}
	if matched, err := validateTxamt(req); !matched {
		return err
	}
	if req.NeedUserInfo != "YES" && req.NeedUserInfo != "NO" {
		return fieldContentError(needUserInfo)
	}

	if req.TimeExpire != "" {
		if matched, err := validateTimeExpire(req.TimeExpire); !matched {
			return err
		}
	}

	return
}

// validateTimeExpire 验证失效时间
func validateTimeExpire(timeExpire string) (bool, *model.ScanPayResponse) {
	if mactch, _ := regexp.MatchString(`^\d{14}$`, timeExpire); !mactch {
		return false, fieldFormatError("timeExpire")
	}
	return true, nil
}

// validateUserName 验证商户名称
func validateUserName(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {

	switch req.CheckName {
	case "FORCE_CHECK", "OPTION_CHECK":
		if req.UserName == "" {
			return false, fieldEmptyError(userName)
		}
	case "NO_CHECK":
		// do nothing
	default:
		return false, fieldContentError(checkName)
	}
	return true, nil
}

// validateOrderNum 验证订单号
func validateOrderNum(no string) (bool, *model.ScanPayResponse) {

	if len(no) > 64 {
		return false, fieldFormatError(orderNum)
	}
	// 是否包含中文或其他非法字符
	if len([]rune(no)) != len(no) {
		return false, fieldFormatError(orderNum)
	}
	return true, nil
}

// validateGoodsInfo 验证商品格式
func validateGoodsInfo(goods string) (bool, *model.ScanPayResponse) {

	// 不校验，在送到渠道时解析，如果解析错误，那么不送
	// if goods != "" {
	// 	toRunes := []rune(goods)
	// 	if len(toRunes) > 120 {
	// 		return false, fieldFormatError(goodsInfo)
	// 	}
	// 	goodsArray := strings.Split(goods, ";")
	// 	for i, v := range goodsArray {
	// 		// 处理最后多送了;的情况
	// 		if i == len(goodsArray)-1 && v == "" {
	// 			continue
	// 		}
	// 		good := strings.Split(v, ",")
	// 		if len(good) != 3 {
	// 			return false, fieldFormatError(goodsInfo)
	// 		}
	// 		// 金额
	// 		if matched, _ := regexp.MatchString(`^(([1-9]\d*)|0)(\.(\d){1,2})?$`, good[1]); !matched {
	// 			return false, fieldFormatError(goodsInfo)
	// 		}
	// 		// 数量
	// 		if matched, _ := regexp.MatchString(`^\d+$`, good[2]); !matched {
	// 			return false, fieldFormatError(goodsInfo)
	// 		}
	// 	}
	// }

	return true, nil
}

// validateTxamt 验证金额
func validateTxamt(req *model.ScanPayRequest) (bool, *model.ScanPayResponse) {

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
	if toInt == minTxamt {
		return false, fieldFormatError(txamt)
	}

	req.IntTxamt = toInt
	return true, nil
}

// validateMchntid 验证商户号格式
func validateMchntid(mcid string) (bool, *model.ScanPayResponse) {
	if len(mcid) > 15 {
		return false, fieldFormatError(mchntid)
	}
	return true, nil
}

// fieldContentError 字段内容错误
func fieldContentError(f string) *model.ScanPayResponse {

	errMsg := strings.Replace(contentError.ISO8583Msg, "fieldName", f, 1)
	return &model.ScanPayResponse{
		Respcd:      contentError.ISO8583Code,
		ErrorDetail: errMsg,
		ErrorCode:   contentError.ErrorCode,
	}
}

// fieldEmptyError 字段为空
func fieldEmptyError(f string) *model.ScanPayResponse {

	errMsg := strings.Replace(emptyError.ISO8583Msg, "fieldName", f, 1)
	return &model.ScanPayResponse{
		Respcd:      emptyError.ISO8583Code,
		ErrorDetail: errMsg,
		ErrorCode:   emptyError.ErrorCode,
	}
}

// fieldFormatError 字段格式错误
func fieldFormatError(f string) *model.ScanPayResponse {

	errMsg := strings.Replace(formatError.ISO8583Msg, "fieldName", f, 1)
	return &model.ScanPayResponse{
		Respcd:      formatError.ISO8583Code,
		ErrorDetail: errMsg,
		ErrorCode:   formatError.ErrorCode,
	}
}
