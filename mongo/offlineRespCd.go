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
var defaultResp = &model.ScanPayRespCode{"", "", "58", "未知应答,请联系管理员", true, "UNKNOWN"}

type scanPayRespCollection struct {
	name string
}

var spRespCache = cache.New(model.Cache_ScanPayResp)

// Get 根据传入的errorCode类型得到Resp对象
// 屏蔽8583与6位应答码的差别
func (c *scanPayRespCollection) Get(errorCode string) (resp *model.ScanPayRespCode) {

	o, found := spRespCache.Get(errorCode)
	if found {
		resp = o.(*model.ScanPayRespCode)
		return resp
	}

	resp = &model.ScanPayRespCode{}
	err := database.C(c.name).Find(bson.M{"errorCode": errorCode}).One(resp)
	if err != nil {
		log.Errorf("can not find scanPayResp for %s: %s", errorCode, err)
		// 没找到对应应答码，返回默认应答
		return defaultResp
	}

	// save cache
	spRespCache.Set(errorCode, resp, cache.NoExpiration)

	return resp
}

// GetByAlp 由支付宝应答得到Resp对象
func (c *scanPayRespCollection) GetByAlp(code string) (resp *model.ScanPayRespCode) {
	resp = &model.ScanPayRespCode{}
	err := database.C(c.name).Find(bson.M{"alp.code": code}).One(resp)
	if err != nil {
		log.Errorf("can not find scanPayResp for %s: %s", code, err)
		return defaultResp
	}

	return resp
}

// GetByWxp 由微信应答得到Resp对象
func (c *scanPayRespCollection) GetByWxp(code, busicd string) (resp *model.ScanPayRespCode) {
	resp = &model.ScanPayRespCode{}

	q := bson.M{
		"wxp": bson.M{
			"$elemMatch": bson.M{
				"code":   code,
				"busicd": busicd,
			},
		},
	}
	err := database.C(c.name).Find(q).One(resp)
	if err != nil {
		log.Errorf("can not find scanPayResp for %s: %s", code, err)
		return defaultResp
	}
	return resp
}

/* only use for import respCode */

func (c *scanPayRespCollection) Add(r *model.ScanPayCSV) error {
	err := database.C(c.name).Insert(r)
	return err
}

func (c *scanPayRespCollection) FindOne(code string) (*model.ScanPayCSV, error) {
	q := new(model.ScanPayCSV)
	err := database.C(c.name).Find(bson.M{"ISO8583Code": code}).One(q)
	return q, err
}

func (c *scanPayRespCollection) Update(r *model.ScanPayCSV) error {
	err := database.C(c.name).Update(bson.M{"ISO8583Code": r.ISO8583Code}, r)
	return err
}
