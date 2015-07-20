package mongo

import (
	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

// TODO delete
func init() {
	OffLineCdCol = make(map[string]string)
	OffLineCdCol["00"] = "成功"
	OffLineCdCol["01"] = "交易失败"
	OffLineCdCol["03"] = "商户错误"
	OffLineCdCol["09"] = "处理中"
	OffLineCdCol["12"] = "签名错误"
	OffLineCdCol["13"] = "退款失败"
	OffLineCdCol["14"] = "条码错误或过期"
	OffLineCdCol["15"] = "无此渠道"
	OffLineCdCol["22"] = "撤销失败"
	OffLineCdCol["05"] = "不支持该交易类型"
	OffLineCdCol["19"] = "订单号重复"
	OffLineCdCol["25"] = "订单不存在"
	OffLineCdCol["30"] = "报文错误"
	OffLineCdCol["31"] = "权限不足"
	OffLineCdCol["51"] = "余额不足"
	OffLineCdCol["54"] = "订单已关闭或取消"
	OffLineCdCol["58"] = "未知应答码类型"
	OffLineCdCol["64"] = "退款金额超过原订单金额"
	OffLineCdCol["91"] = "外部系统错误"
	OffLineCdCol["96"] = "内部系统错误"
	OffLineCdCol["98"] = "交易超时"
}

// TODO delete
// OffLineRespCd 扫码支付应答码
func OffLineRespCd(code string) *model.ScanPayResponse {

	errorDetail, respCd := "", ""

	switch code {
	case "SUCCESS":
		respCd = "00"
	case "INPROCESS":
		respCd = "09"
	case "FAIL":
		respCd = "01"
	case "NO_ROUTERPOLICY", "NO_CHANMER", "NO_PERMISSION":
		respCd = "31"
	case "NOT_PAYTRADE", "NOT_SUCESS_TRADE", "TRADE_REFUNDED", "REFUND_TIME_ERROR":
		respCd = "13"
	case "TRADE_AMT_INCONSISTENT", "TRADE_HAS_REFUND":
		respCd = "64"
	case "CANCEL_TIME_ERROR":
		respCd = "22"
	case "SYSTEM_ERROR", "CONNECT_ERROR":
		respCd = "96"
	case "ORDER_DUPLICATE":
		respCd = "19"
	case "SIGN_AUTH_ERROR":
		respCd = "12"
	case "NO_MERCHANT":
		respCd = "03"
	case "TRADE_OVERTIME":
		respCd = "98"
	case "DATA_ERROR":
		respCd = "30"
	case "QRCODE_INVALID":
		respCd = "14"
	case "NO_CHANNEL":
		respCd = "15"
	case "TRADE_NOT_EXIST":
		respCd = "25"
	case "ORDER_CLOSED":
		respCd = "54"
	case "NOT_SUPPORT_TYPE":
		respCd = "05"
	case "INSUFFICIENT_BALANCE":
		respCd = "51"
	case "UNKNOWN_ERROR":
		respCd = "91"
	default:
		respCd = "58"
	}

	errorDetail = OffLineCdCol[respCd]
	return &model.ScanPayResponse{ErrorDetail: errorDetail, Respcd: respCd}
}

var ScanPayRespCol = &scanPayRespCollection{"respCode.sp"}
var OffLineCdCol map[string]string

type scanPayRespCollection struct {
	name string
}

type scanPayResp struct {
	RespCode      string `bson:"respCode"`
	RespMsg       string `bson:"respMsg"`
	Iso8583Code   string `bson:"iso8583Code"`
	Iso8583Msg    string `bson:"iso8583Msg"`
	IsUseChanDesc bool   `bson:"isUseChanDesc"`
	ErrorCode     string `bson:"errorCode"`
}

var spRespCache = cache.New(model.Cache_ScanPayResp)

// Get 根据传入的code类型得到Resp对象
func (c *scanPayRespCollection) Get(code string) (resp *scanPayResp) {

	o, found := spRespCache.Get(code)
	if found {
		resp = o.(*scanPayResp)
		return resp
	}

	resp = &scanPayResp{}
	err := database.C(c.name).Find(bson.M{"errorCode": code}).One(resp)
	if err != nil {
		log.Errorf("can not find scanPayResp for %s: %s", code, err)
		return resp
	}

	// save cache
	spRespCache.Set(code, resp, cache.NoExpiration)

	return resp
}

// GetByAlp 由支付宝应答得到Resp对象
func (c *scanPayRespCollection) GetByAlp(code, busicd string) (resp *scanPayResp) {
	resp = &scanPayResp{}
	database.C(c.name).Find(bson.M{"alp.code": code}).One(resp)
	return resp
}

// GetByWxp 由微信应答得到Resp对象
func (c *scanPayRespCollection) GetByWxp(code, busicd string) (resp *scanPayResp) {
	resp = &scanPayResp{}
	database.C(c.name).Find(bson.M{"wxp.code": code}).One(resp)
	return resp
}
