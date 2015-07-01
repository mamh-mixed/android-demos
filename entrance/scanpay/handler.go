package scanpay

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"strings"
)

type HandleFuc func(req *model.ScanPay) (ret *model.ScanPayResponse)

// ScanPayHandle 执行扫码支付逻辑
func ScanPayHandle(reqBytes []byte) []byte {

	log.Debugf("request body: %s", string(reqBytes))
	// 解析请求内容
	req := new(model.ScanPay)
	err := json.Unmarshal(reqBytes, req)
	if err != nil {
		log.Errorf("fail to unmarshal jsonStr(%s): %s", reqBytes, err)
		return ErrorResponse(req, "INVALID_PARAMETER")
	}

	// 具体业务
	ret := router(req)

	// 应答
	retBytes, err := json.Marshal(ret)

	log.Debugf("handled body: %s", retBytes)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", ret, err)
		return ErrorResponse(req, "SYSTEM_ERROR")
	}

	retStr := string(retBytes)
	retLen := fmt.Sprintf("%0.4d", len(retStr))

	return []byte(retLen + retStr)
}

// router 分发业务逻辑
func router(req *model.ScanPay) (ret *model.ScanPayResponse) {

	switch req.Busicd {
	case model.Purc:
		// ret = barcodePay(req)
		ret = doScanPay(validateBarcodePay, core.BarcodePay, req)
	case model.Paut:
		// ret = qrCodeOfflinePay(req)
		ret = doScanPay(validateQrCodeOfflinePay, core.QrCodeOfflinePay, req)
	case model.Inqy:
		// ret = enquiry(req)
		ret = doScanPay(validateEnquiry, core.Enquiry, req)
	case model.Refd:
		// ret = refund(req)
		ret = doScanPay(validateRefund, core.Refund, req)
	case model.Void:
		// ret = cancel(req)
		ret = doScanPay(validateCancel, core.Cancel, req)
	case model.Canc:

		ret = doScanPay(validateClose, core.Close, req)
	default:
		ret = mongo.OffLineRespCd("INVALID_PARAMETER")
	}

	// 补充原信息返回
	fillResponseInfo(req, ret)

	return ret
}

// doScanPay 执行业务逻辑
func doScanPay(validateFuc, processFuc HandleFuc, req *model.ScanPay) (ret *model.ScanPayResponse) {

	// 验证字段
	if validateFuc(req); ret != nil {
		return ret
	}

	mer, err := mongo.MerchantColl.Find(req.Mchntid)
	if err != nil {
		return mongo.OffLineRespCd("NO_MERCHANT_MATCH") // todo check error code
	}

	// 验签
	sign := req.Sign
	if mer.IsNeedSign {
		req.Sign = "" // 置空
		signBytes := sha1.Sum([]byte(req.DictSortMsg() + mer.SignKey))
		s := fmt.Sprintf("%x", signBytes[:])
		if s != sign {
			log.Errorf("sign should be %s, but get %s", s, sign)
			return mongo.OffLineRespCd("AUTH_NO_ERROR")
		}
	}

	// 验证接口权限
	if !strings.Contains(strings.Join(mer.Permission, ","), req.Busicd) {
		log.Errorf("merchant %s request %s interface without permission!", req.Mchntid, req.Busicd)
		return mongo.OffLineRespCd("NO_PERMISSION") // todo check error code
	}

	// process
	ret = processFuc(req)

	// TODO sign
	req.Sign = sign

	return ret

}

func fillResponseInfo(req *model.ScanPay, ret *model.ScanPayResponse) {

	// 如果空白，默认将原信息返回
	if ret.Busicd == "" {
		ret.Busicd = req.Busicd
	}
	if ret.Inscd == "" {
		ret.Inscd = req.Inscd
	}
	if ret.Chcd == "" {
		ret.Chcd = req.Chcd
	}
	if ret.Mchntid == "" {
		ret.Mchntid = req.Mchntid
	}
	if ret.Txamt == "" {
		ret.Txamt = req.Txamt
	}
	if ret.OrigOrderNum == "" {
		ret.OrigOrderNum = req.OrigOrderNum
	}
	if ret.OrderNum == "" {
		ret.OrderNum = req.OrderNum
	}
	// TODO
	if ret.Sign == "" {
		ret.Sign = req.Sign
	}
	ret.Txndir = "A"
}

// ErrorResponse 返回错误信息
func ErrorResponse(req *model.ScanPay, errorCode string) []byte {

	ret := mongo.OffLineRespCd(errorCode)
	fillResponseInfo(req, ret)

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
	}
	retStr := string(retBytes)
	retLen := fmt.Sprintf("%0.4d", len(retStr))

	return []byte(retLen + retStr)
}
