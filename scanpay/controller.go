package scanpay

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// ScanPayHandle 执行扫码支付逻辑
func ScanPayHandle(reqBytes []byte, isGBK bool) []byte {
	log.Infof("from merchant message: %s", string(reqBytes))

	// 解析请求内容
	req := new(model.ScanPayRequest)
	// 设置请求方式
	req.IsGBK = isGBK
	// 解析json
	err := json.Unmarshal(reqBytes, req)
	if err != nil {
		log.Errorf("fail to unmarshal json(%s): %s", reqBytes, err)
		return errorResp(req, "DATA_ERROR")
	}

	// 具体业务
	ret := dispatch(req)

	// 应答
	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", ret, err)
		return errorResp(req, "SYSTEM_ERROR")
	}

	log.Infof("to merchant message: %s", retBytes)
	return retBytes
}

// dispatch 分发业务逻辑
func dispatch(req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
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
	case model.Qyfk:
		ret = doScanPay(validateEnterprisePay, core.EnterprisePay, req)
	case model.Jszf:
		ret = doScanPay(validatePublicPay, core.PublicPay, req)
	default:
		ret = fieldContentError(buiscd)
		ret.FillWithRequest(req)
	}

	return ret
}

var nonCheckSignBusicd = model.Jszf

type handleFunc func(req *model.ScanPayRequest) (ret *model.ScanPayResponse)

// doScanPay 执行业务逻辑
func doScanPay(validateFunc, processFunc handleFunc, req *model.ScanPayRequest) (ret *model.ScanPayResponse) {
	// 修复请求失败时，应答签名也失败的 bug
	var signKey string
	defer func() {
		// 7. 补充信息
		ret.FillWithRequest(req)

		// 8. 如果是 gbk 进来的，兼容老插件和商户，不返回中文，不返回 errorCode
		if req.IsGBK {
			ret.ErrorDetail = ret.ErrorCode
			ret.ErrorCode = ""
		}

		// 9. 对返回报文签名
		if signKey != "" {
			log.Debug("sign content to return : " + ret.SignMsg())
			ret.Sign = security.SHA1WithKey(ret.SignMsg(), signKey)
		}
	}()

	// 1. 开始处理逻辑前，验证字段
	if ret = validateFunc(req); ret != nil {
		return ret
	}

	// 2. 先检查商户代码，如果不存在，直接报错
	mer, err := mongo.MerchantColl.Find(req.Mchntid)
	if err != nil {
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("NO_MERCHANT"))
		return
	}

	if mer.IsNeedSign {
		signKey = mer.SignKey
	}

	// 3. 检查机构号
	// ac := strings.Trim(req.AgentCode, " ")
	// if mer.AgentCode != ac {
	// 	ret = fieldContentError(agentCode)
	// 	return
	// }

	// 4. 商户、机构号都通过后，验证接口权限
	if !util.StringInSlice(req.Busicd, mer.Permission) {
		log.Errorf("merchant(%s) request(%s) refused", req.Mchntid, req.Busicd)
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("NO_PERMISSION"))
		return
	}

	// 5. 商户存在，则验签
	if mer.IsNeedSign && req.Busicd != nonCheckSignBusicd {
		log.Debug("sign msg : " + req.SignMsg())
		sig := security.SHA1WithKey(req.SignMsg(), mer.SignKey)
		if sig != req.Sign {
			log.Errorf("mer(%s) sign failed: data=%v, sign=%s", req.Mchntid, req, sig)
			ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("SIGN_AUTH_ERROR"))
			return
		}
	}

	// 过滤包含空格字符串
	req.Chcd = strings.TrimSpace(req.Chcd)
	req.AgentCode = strings.TrimSpace(req.AgentCode)

	// 6. 开始业务处理
	ret = processFunc(req)

	return ret
}

// errorResp 返回错误信息
func errorResp(req *model.ScanPayRequest, errorCode string) []byte {
	spResp := mongo.ScanPayRespCol.Get(errorCode)
	ret := &model.ScanPayResponse{
		Respcd:      spResp.ISO8583Code,
		ErrorDetail: spResp.ISO8583Msg,
	}

	// 如果是gbk端口，采用英文描述应答
	if req.IsGBK {
		ret.ErrorDetail = errorCode
	}

	ret.FillWithRequest(req)

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
	}
	return retBytes
}

// weixinNotifyCtrl 微信异步通知处理(预下单)
func weixinNotifyCtrl(req *weixin.WeixinNotifyReq) error {
	return core.ProcessWeixinNotify(req)
}

// alipayNotifyCtrl 支付宝异步通知处理(预下单)
func alipayNotifyCtrl(v url.Values) error {
	return core.ProcessAlipayNotify(v)
}
