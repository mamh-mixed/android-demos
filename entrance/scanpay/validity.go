package scanpay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"regexp"
	"strconv"
	"strings"
)

const (
	minTxamt = 0
	maxTxamt = 1e10 - 1
)

// fieldName
const (
	txamt        = "txamt"
	orderNum     = "orderNum"
	origOrderNum = "origOrderNum"
	inscd        = "inscd"
	mchntid      = "mchntid"
	scanCodeId   = "scanCodeId"
	chcd         = "chcd"
)

var (
	emptyError  = mongo.ScanPayRespCol.Get("DATA_EMPTY_ERROR")
	formatError = mongo.ScanPayRespCol.Get("DATA_FORMAT_ERROR")
)

// validateBarcodePay 验证扫码下单的参数
func validateBarcodePay(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Inscd == "":
		return fieldEmptyError(inscd)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.ScanCodeId == "":
		return fieldEmptyError(scanCodeId)
	}

	// 验证格式
	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return fieldFormatError(txamt)
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return fieldFormatError(mchntid)
	}

	if matched, _ := regexp.MatchString(`^\d{14,24}$`, req.ScanCodeId); !matched {
		return fieldFormatError("scanCodeId")
	}

	// 转换金额
	toInt, err := strconv.ParseInt(req.Txamt, 10, 64)
	if err != nil {
		return fieldFormatError(txamt)
	}

	// 金额范围
	if toInt == minTxamt || toInt > maxTxamt {
		return fieldFormatError(txamt)
	}

	req.IntTxamt = toInt

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
	case req.Inscd == "":
		return fieldEmptyError(inscd)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	}

	// 验证格式
	if req.Chcd != "WXP" && req.Chcd != "ALP" {
		return fieldFormatError(chcd)
	}

	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return fieldFormatError(txamt)
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return fieldFormatError(mchntid)
	}

	// 转换金额
	toInt, err := strconv.ParseInt(req.Txamt, 10, 64)
	if err != nil {
		return fieldFormatError(txamt)
	}

	// 金额范围
	if toInt == minTxamt || toInt > maxTxamt {
		return fieldFormatError(txamt)
	}

	req.IntTxamt = toInt

	return
}

// validateEnquiry 验证查询接口的参数
func validateEnquiry(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	switch {
	case req.OrigOrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Inscd == "":
		return fieldEmptyError(inscd)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	}

	// TODO validate format
	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return fieldFormatError(mchntid)
	}

	return
}

// validateRefund 验证退款接口的参数
func validateRefund(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Inscd == "":
		return fieldEmptyError(inscd)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.Txamt == "":
		return fieldEmptyError(txamt)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	}

	// TODO validate format
	if matched, _ := regexp.MatchString(`^\d{12}$`, req.Txamt); !matched {
		return fieldFormatError(txamt)
	}

	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return fieldFormatError(mchntid)
	}

	// 转换金额
	toInt, err := strconv.ParseInt(req.Txamt, 10, 64)
	if err != nil {
		return fieldFormatError(txamt)
	}

	// 金额范围
	if toInt == minTxamt || toInt > maxTxamt {
		return fieldFormatError(txamt)
	}

	req.IntTxamt = toInt

	return
}

// validateCancel 验证撤销接口参数
func validateCancel(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Inscd == "":
		return fieldEmptyError(inscd)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	}

	// 验证格式
	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return fieldFormatError(mchntid)
	}

	return
}

// validateCancel 验证关闭订单接口参数
func validateClose(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {

	// 验证非空
	switch {
	case req.OrderNum == "":
		return fieldEmptyError(orderNum)
	case req.Inscd == "":
		return fieldEmptyError(inscd)
	case req.Mchntid == "":
		return fieldEmptyError(mchntid)
	case req.OrigOrderNum == "":
		return fieldEmptyError(origOrderNum)
	}

	// 验证格式
	if matched, _ := regexp.MatchString(`^\d{15}$`, req.Mchntid); !matched {
		return fieldFormatError(mchntid)
	}

	return
}

// fieldEmptyError 字段为空
func fieldEmptyError(f string) *model.ScanPayResponse {

	errMsg := strings.Replace(emptyError.ISO8583Msg, "fieldName", f, 1)
	return &model.ScanPayResponse{
		Respcd:      emptyError.ISO8583Code,
		ErrorDetail: errMsg,
	}
}

// fieldFormatError 字段格式错误
func fieldFormatError(f string) *model.ScanPayResponse {

	errMsg := strings.Replace(formatError.ISO8583Msg, "fieldName", f, 1)
	return &model.ScanPayResponse{
		Respcd:      formatError.ISO8583Code,
		ErrorDetail: errMsg,
	}
}
