// dmf1.0外海接口
package oversea

import (
	"errors"
	"fmt"
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"math"
	"time"
)

var DefaultClient alp
var alipayCurrency map[string]int

// TODO 常用状态码整合到一起
var (
	CloseCode, CloseMsg, _         = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
	InprocessCode, InprocessMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("INPROCESS")
	SuccessCode, SuccessMsg, _     = mongo.ScanPayRespCol.Get8583CodeAndMsg("SUCCESS")
	UnKnownCode, UnKnownMsg, _     = mongo.ScanPayRespCol.Get8583CodeAndMsg("CHAN_UNKNOWN_ERROR")
)

func init() {
	initAvailableCurrency()
}

type alp struct{}

func getCommonParams(m *model.ScanPayRequest) scanpay1.CommonReq {
	return scanpay1.CommonReq{
		InputCharset: "utf-8",
		SignKey:      m.SignKey,
		Partner:      m.ChanMerId,
		SpReq:        m,
	}
}

// ProcessBarcodePay 下单
func (a *alp) ProcessBarcodePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	// pay
	b := NewPayReq()
	b.CommonReq = getCommonParams(req)
	b.AlipaySellerId = req.ChanMerId
	b.Currency = req.Currency
	b.TransName = req.Subject
	b.PartnerTransId = req.OrderNum
	b.BuyerIdentityCode = req.ScanCodeId
	b.ExtendInfo = req.ExtendParams
	b.TransAmount = handleAmt(req.IntTxamt, req.Currency)
	b.TransCreateTime = time.Now().Format("20060102150405") // TODO

	// resp
	p := &PayResp{}
	if err := scanpay1.Execute(b, p); err != nil {
		return nil, err
	}

	resp := &model.ScanPayResponse{}
	resp.ChannelOrderNum = p.Response.Alipay.AlipayTransId
	resp.PayTime = p.Response.Alipay.AlipayPayTime
	resp.ConsumerId = p.Response.Alipay.AlipayBuyerLoginId
	resp.Rate = p.Response.Alipay.ExchangeRate

	// result
	alipayResponseHandle(p, resp, "createandpay")

	return resp, nil
}

func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, errors.New("Not support yet!")
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	b := NewRefundReq()
	b.CommonReq = getCommonParams(req)
	b.PartnerTransId = req.OrigOrderNum
	b.PartnerRefundId = req.OrderNum
	b.RefundAmount = handleAmt(req.IntTxamt, req.Currency)
	b.Currency = req.Currency

	// resp
	p := &RefundResp{}
	if err := scanpay1.Execute(b, p); err != nil {
		return nil, err
	}

	resp := &model.ScanPayResponse{}
	resp.ChannelOrderNum = p.Response.Alipay.AlipayTransId
	resp.Rate = p.Response.Alipay.ExchangeRate

	// result
	alipayResponseHandle(p, resp, "refund")

	return resp, nil
}

// ProcessEnquiry 查询
func (a *alp) ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	// query
	b := NewQueryReq()
	b.CommonReq = getCommonParams(req)
	b.PartnerTransId = req.OrigOrderNum

	// resp
	p := &QueryResp{}
	if err := scanpay1.Execute(b, p); err != nil {
		return nil, err
	}

	resp := &model.ScanPayResponse{}
	resp.ChannelOrderNum = p.Response.Alipay.AlipayTransId
	resp.PayTime = p.Response.Alipay.AlipayPayTime
	resp.ConsumerId = p.Response.Alipay.AlipayBuyerLoginId

	// result
	alipayResponseHandle(p, resp, "query")

	if p.ResultCode() == "SUCCESS" {
		switch p.Response.Alipay.AlipayTransStatus {
		case "WAIT_BUYER_PAY":
			resp.Respcd, resp.ErrorDetail = InprocessCode, InprocessMsg
		case "TRADE_CLOSED":
			resp.Respcd, resp.ErrorDetail = CloseCode, CloseMsg
		default:
			// success
		}
	}

	return resp, nil
}

// ProcessCancel 撤销
func (a *alp) ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return a.ProcessClose(req)
}

// ProcessClose 关闭
func (a *alp) ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	// reverse
	b := NewReverseReq()
	b.CommonReq = getCommonParams(req)
	b.PartnerTransId = req.OrigOrderNum

	// resp
	p := &ReverseResp{}
	if err := scanpay1.Execute(b, p); err != nil {
		return nil, err
	}

	resp := &model.ScanPayResponse{}
	resp.ChannelOrderNum = p.Response.Alipay.AlipayTransId
	resp.PayTime = p.Response.Alipay.AlipayReverseTime

	// result
	alipayResponseHandle(p, resp, "cancel")

	return resp, nil
}

// ProcessRefundQuery 退款查询
func (a *alp) ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, errors.New("Not support yet!")
}

func alipayResponseHandle(p scanpay1.BaseResp, resp *model.ScanPayResponse, service string) {

	var errorCode = p.ErrorCode()
	if errorCode != "" {
		// error
		spCode := mongo.ScanPayRespCol.GetByAlp(errorCode, service)
		resp.Respcd = spCode.ISO8583Code
		resp.ErrorDetail = spCode.ISO8583Msg
		if !spCode.IsUseISO {
			resp.ErrorDetail = errorCode // TODO
		}
		return
	}

	switch p.ResultCode() {
	case "SUCCESS":
		resp.Respcd, resp.ErrorDetail = SuccessCode, SuccessMsg
	case "UNKNOW":
		resp.Respcd, resp.ErrorDetail = InprocessCode, InprocessMsg
	default:
		resp.Respcd, resp.ErrorDetail = UnKnownCode, UnKnownMsg
	}
}

func handleAmt(txamt int64, currency string) string {
	if cur, ok := alipayCurrency[currency]; ok {
		realUnit := math.Pow(10, float64(cur))
		return fmt.Sprintf("%0."+fmt.Sprintf("%d", cur)+"f", float64(txamt)/realUnit)
	}
	return fmt.Sprintf("%d", txamt)
}

func initAvailableCurrency() {
	// 币种、精度
	alipayCurrency = make(map[string]int)
	alipayCurrency["GBP"] = 2
	alipayCurrency["HKD"] = 2
	alipayCurrency["USD"] = 2
	alipayCurrency["CHF"] = 2
	alipayCurrency["SGD"] = 2
	alipayCurrency["SEK"] = 2
	alipayCurrency["DKK"] = 2
	alipayCurrency["NOK"] = 2
	alipayCurrency["JPY"] = 0
	alipayCurrency["CAD"] = 2
	alipayCurrency["AUD"] = 2
	alipayCurrency["EUR"] = 2
	alipayCurrency["KRW"] = 0
}
