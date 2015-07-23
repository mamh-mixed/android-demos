package scanpay

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// ScanPayHandle 执行扫码支付逻辑
func ScanPayHandle(reqBytes []byte) []byte {
	log.Infof("from merchant message: %s", string(reqBytes))

	// 解析请求内容
	req := new(model.ScanPayRequest)
	err := json.Unmarshal(reqBytes, req)
	if err != nil {
		log.Errorf("fail to unmarshal json(%s): %s", reqBytes, err)
		return ErrorResponse(req, "DATA_ERROR")
	}

	// 具体业务
	ret := router(req)

	// 应答
	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", ret, err)
		return ErrorResponse(req, "SYSTEM_ERROR")
	}

	log.Infof("to merchant message: %s", retBytes)
	return retBytes
}

// router 分发业务逻辑
func router(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	switch req.Busicd {
	case model.Purc:
		ret = doScanPay(validateBarcodePay, core.BarcodePay, req)
	case model.Paut:
		ret = doScanPay(validateQrCodeOfflinePay, core.QrCodeOfflinePay, req)
	case model.Inqy:
		ret = doScanPay(validateEnquiry, core.Enquiry, req)
	case model.Refd:
		ret = doScanPay(validateRefund, core.Refund, req)
	case model.Void:
		ret = doScanPay(validateCancel, core.Cancel, req)
	case model.Canc:
		ret = doScanPay(validateClose, core.Close, req)
	default:
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("DATA_ERROR"))
	}

	return ret
}

type handleFuc func(req *model.ScanPayRequest) (ret *model.ScanPayResponse)

// doScanPay 执行业务逻辑
func doScanPay(validateFuc, processFuc handleFuc, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 补充原信息返回
	defer fillResponseInfo(req, ret)

	// 验证字段
	if ret = validateFuc(req); ret != nil {
		return ret
	}

	mer, err := mongo.MerchantColl.Find(req.Mchntid)
	if err != nil {
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("NO_MERCHANT"))
		return
	}

	// 验签
	if mer.IsNeedSign {
		sig := security.SHA1WithKey(req.SignMsg(), mer.SignKey)
		if sig != req.Sign {
			log.Errorf("mer(%s) sign failed: data=%v, sign=%s", req.Mchntid, req, sig)
			ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("SIGN_AUTH_ERROR"))
			return
		}
	}

	// 验证接口权限
	if !util.StringInSlice(req.Busicd, mer.Permission) {
		log.Errorf("merchant(%s) request(%s) refused", req.Mchntid, req.Busicd)
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("NO_PERMISSION"))
		return
	}

	// process
	ret = processFuc(req)

	// 签名
	if mer.IsNeedSign {
		ret.Sign = security.SHA1WithKey(ret.SignMsg(), mer.SignKey)
	}

	return ret
}

// fillResponseInfo 如果空白，默认将原信息返回
func fillResponseInfo(req *model.ScanPayRequest, ret *model.ScanPayResponse) {
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
	if ret.Sign == "" {
		ret.Sign = req.Sign
	}
	ret.Txndir = "A"
}

// ErrorResponse 返回错误信息
func ErrorResponse(req *model.ScanPayRequest, errorCode string) []byte {
	spResp := mongo.ScanPayRespCol.Get(errorCode)
	ret := &model.ScanPayResponse{
		Respcd:      spResp.ISO8583Code,
		ErrorDetail: spResp.ISO8583Msg,
	}
	fillResponseInfo(req, ret)

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
	}
	return retBytes
}
