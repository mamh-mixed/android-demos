package scanpay

import (
	"encoding/json"
	"fmt"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// TcpScanPayHandle tcp处理
func TcpScanPayHandle(reqBytes []byte) []byte {

	// gbk解编码
	dgbk, _ := util.GBKTranscoder.Decode(string(reqBytes))

	// 处理
	retBytes := ScanPayHandle([]byte(dgbk))

	// gbk编码
	retStr := string(retBytes)

	egbk, _ := util.GBKTranscoder.Encode(retStr)

	// 长度位
	retLen := fmt.Sprintf("%04d", len(egbk))

	return []byte(retLen + egbk)
}

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
		ret = fieldFormatError("busicd")
		fillResponseInfo(req, ret)
	}

	return ret
}

type handleFunc func(req *model.ScanPayRequest) (ret *model.ScanPayResponse)

// doScanPay 执行业务逻辑
func doScanPay(validateFuc, processFunc handleFunc, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 验证字段
	if ret = validateFuc(req); ret != nil {
		fillResponseInfo(req, ret)
		return ret
	}

	mer, err := mongo.MerchantColl.Find(req.Mchntid)
	if err != nil {
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("NO_MERCHANT"))
		fillResponseInfo(req, ret)
		return
	}

	// 验签
	sign := req.Sign
	if mer.IsNeedSign {
		req.Sign = "" // 置空
		s := security.SHA1WithKey(req.SignMsg(), mer.SignKey)
		if s != sign {
			log.Errorf("sign should be %s, but get %s", s, sign)
			ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("SIGN_AUTH_ERROR"))
			fillResponseInfo(req, ret)
			return
		}
	}

	// 验证接口权限
	if !util.StringInSlice(req.Busicd, mer.Permission) {
		log.Errorf("merchant %s request %s interface without permission!", req.Mchntid, req.Busicd)
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("NO_PERMISSION"))
		fillResponseInfo(req, ret)
		return
	}

	// process
	ret = processFunc(req)

	// 补充原信息返回
	fillResponseInfo(req, ret)

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
