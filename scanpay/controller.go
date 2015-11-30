package scanpay

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/logs"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// 专门做监控的商户
var monitorMerId = goconf.Config.App.MonitorMerId

// ScanPayHandle 执行扫码支付逻辑
func ScanPayHandle(reqBytes []byte, isGBK bool) []byte {

	// 解析请求内容
	req := model.NewScanPayRequest()
	// 设置请求方式
	req.IsGBK = isGBK
	// 解析json
	err := json.Unmarshal(reqBytes, req)
	if err != nil {
		log.Errorf("fail to unmarshal json(%s): %s", reqBytes, err)
		return errorResp(req, "DATA_ERROR")
	}

	if req.Mchntid != monitorMerId { // 专门做监控的商户，不打日志
		log.Infof("from merchant message: %s", string(reqBytes))
	}

	// 记录请求时日志
	if req.Mchntid != monitorMerId { // 专门做监控的商户，不记录日志
		logs.SpLogs <- req.GetMerReqLogs()
	}

	// 具体业务
	ret := dispatch(req)

	// 记录返回时日志
	if req.Mchntid != monitorMerId { // 专门做监控的商户，不记录日志
		logs.SpLogs <- req.GetMerRetLogs(ret)
	}

	// 应答
	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("fail to marshal (%+v): %s", ret, err)
		return errorResp(req, "SYSTEM_ERROR")
	}
	if req.Mchntid != monitorMerId { // 专门做监控的商户，不打日志
		log.Infof("to merchant message: %s", retBytes)
	}
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
	case model.Qyzf:
		ret = doScanPay(validateEnterprisePay, core.EnterprisePay, req)
	case model.Jszf:
		ret = doScanPay(validatePublicPay, core.PublicPay, req)
	case model.Veri:
		ret = doScanPay(validatePurchaseCoupons, core.PurchaseCoupons, req)
	case model.Crve:
		ret = doScanPay(validatePurchaseActCoupons, core.PurchaseActCoupons, req)
	case model.Quve:
		ret = doScanPay(validateQueryPurchaseCoupons, core.QueryPurchaseCouponsResult, req)
	case model.Cave:
		ret = doScanPay(validateUndoPurchaseActCoupons, core.UndoPurchaseActCoupons, req)
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

	// 1. 先检查商户代码，如果不存在，直接报错
	mer, err := mongo.MerchantColl.Find(req.Mchntid)
	if err != nil {
		ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("NO_MERCHANT"))
		return
	}
	req.M = *mer

	// 需要验签
	if mer.IsNeedSign {
		signKey = mer.SignKey
	}

	// 2. 开始处理逻辑前，验证字段
	if ret = validateFunc(req); ret != nil {
		return ret
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
	if mer.IsNeedSign {
		content, sig := "", ""
		switch req.Busicd {
		// 公众号支付
		case nonCheckSignBusicd:
			if mer.JsPayVersion == "2.0" {
				content = fmt.Sprintf("backUrl=%s&mchntid=%s&orderNum=%s&txamt=%s", req.NotifyUrl, req.Mchntid, req.OrderNum, req.Txamt)
				sig = security.SHA1WithKey(content, mer.SignKey)
			}
		// 其他接口
		default:
			content = req.SignMsg()
			sig = security.SHA1WithKey(content, mer.SignKey)
		}

		if sig != "" && sig != req.Sign {
			log.Errorf("mer(%s) sign failed: data=%v, expect sign=%s, get sign=%s", req.Mchntid, content, sig, req.Sign)
			ret = model.NewScanPayResponse(*mongo.ScanPayRespCol.Get("SIGN_AUTH_ERROR"))
			return
		}
	}

	// 过滤包含空格字符串
	req.Chcd = strings.TrimSpace(req.Chcd)
	var reqAgentCode = req.AgentCode
	req.AgentCode = mer.AgentCode // 以我们系统的代理代码为准

	// 6. 开始业务处理
	ret = processFunc(req)

	ret.AgentCode = strings.TrimSpace(reqAgentCode) // 返回时送回原代理代码

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

	// 验签, 如果验签失败，只打印日志，不中止逻辑
	// buf, err := util.Query(req)
	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }
	// buf.WriteString("&key=" + signKey)
	// log.Debugf("%s", buf.String())
	//
	// sign := md5.Sum(buf.Bytes())
	// actual := strings.ToUpper(hex.EncodeToString(sign[:]))
	//
	// if actual != req.GetSign() {
	// 	log.Errorf("check sign error: query={%s}, expected=%s, actual=%s", buf.String(), req.GetSign(), actual)
	// }

	return core.ProcessWeixinNotify(req)
}

// alipayNotifyCtrl 支付宝异步通知处理(预下单)
func alipayNotifyCtrl(v url.Values) error {
	return core.ProcessAlipayNotify(v)
}

var noMerCode, noMerMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("NO_MERCHANT")
var dataErrCode, dataErrMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("DATA_ERROR")
var signErrCode, signErrMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("SIGN_AUTH_ERROR")
var sysErrCode, sysErrMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("SYSTEM_ERROR")

// getBillsCtrl 获取商户对账单
func getBillsCtrl(reqBytes []byte) []byte {

	var errorResp = `{"respcd":"%s","errorDetail":"%s"}`

	var q = &struct {
		MerId        string `json:"mchntid" url:"mchntid,omitempty"`
		Busicd       string `json:"busicd" url:"busicd,omitempty"`
		SettDate     string `json:"settDate" url:"settDate,omitempty"`
		NextOrderNum string `json:"nextOrderNum" url:"nextOrderNum,omitempty"`
		Sign         string `json:"sign" url:"-"`
	}{}

	err := json.Unmarshal(reqBytes, q)
	if err != nil {
		return []byte(fmt.Sprintf(errorResp, dataErrCode, dataErrMsg))

	}

	// 查找商户
	m, err := mongo.MerchantColl.Find(q.MerId)
	if err != nil {
		return []byte(fmt.Sprintf(errorResp, noMerCode, noMerMsg))
	}

	// 验签
	if m.IsNeedSign {
		rbuf, _ := util.Query(q)
		rsign := security.SHA1WithKey(rbuf.String(), m.SignKey)
		if rsign != q.Sign {
			log.Warnf("check sign error, expect=%s,get=%s", rsign, q.Sign)
			return []byte(fmt.Sprintf(errorResp, signErrCode, signErrMsg))
		}
	}

	// 获取对账单
	p := query.GetBills(&model.QueryCondition{
		MerId:        q.MerId,
		Busicd:       q.Busicd,
		StartTime:    q.SettDate + " 00:00:00",
		EndTime:      q.SettDate + " 23:59:59",
		NextOrderNum: q.NextOrderNum,
	})

	// 默认返回
	var result = &struct {
		Respcd       string      `json:"respcd,omitempty" url:"respcd,omitempty"`
		ErrorDetail  string      `json:"errorDetail,omitempty" url:"errorDetail,omitempty"`
		Count        int         `json:"count" url:"count"`
		Rec          interface{} `json:"rec,omitempty" url:"-"`
		RecStr       string      `json:"-" url:"rec,omitempty"`
		NextOrderNum string      `json:"nextOrderNum,omitempty" url:"nextOrderNum,omitempty"`
		Sign         string      `json:"sign" url:"-"`
	}{}

	recBytes, _ := json.Marshal(p.Rec)
	result.RecStr = string(recBytes)
	result.Rec = p.Rec
	result.Respcd = p.RespCode
	result.ErrorDetail = p.RespMsg
	result.Count = p.Count
	result.NextOrderNum = p.NextOrderNum

	// 签名
	bbuf, _ := util.Query(result)
	result.Sign = security.SHA1WithKey(bbuf.String(), m.SignKey)

	retBytes, _ := json.Marshal(result)

	return retBytes
}
