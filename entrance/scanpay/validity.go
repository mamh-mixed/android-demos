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
	txamt      = "txamt"
	orderNum   = "orderNum或origOrderNum"
	inscd      = "inscd"
	mchntid    = "mchntid"
	scanCodeId = "scanCodeId"
	chcd       = "chcd"
	goodsInfo  = "goodsInfo"
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
	case req.Inscd == "":
		return fieldEmptyError(inscd)
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
	case req.Inscd == "":
		return fieldEmptyError(inscd)
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
	case req.Inscd == "":
		return fieldEmptyError(inscd)
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

	if goods != "" {
		toRunes := []rune(goods)
		if len(toRunes) > 120 {
			return false, fieldFormatError(goodsInfo)
		}
		// todo 验证格式
	}

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
	if toInt == minTxamt || toInt > maxTxamt {
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
