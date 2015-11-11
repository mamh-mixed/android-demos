// dmf1.0外海接口
package oversea

import (
	"errors"
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	// "github.com/omigo/log"
	"time"
)

var DefaultClient alp

// TODO 常用状态码整合到一起
var (
	CloseCode, _, CloseMsg         = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
	InprocessCode, _, InprocessMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("INPROCESS")
	SuccessCode, _, SuccessMsg     = mongo.ScanPayRespCol.Get8583CodeAndMsg("SUCCESS")
	UnKnownCode, _, UnKnownMsg     = mongo.ScanPayRespCol.Get8583CodeAndMsg("CHAN_UNKNOWN_ERROR")
)

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
	b.TransAmount = req.ActTxamt
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
	b.RefundAmount = req.ActTxamt
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
		resp.ErrorDetail = spCode.ErrorCode
		if !spCode.IsUseISO {
			resp.ErrorDetail = errorCode // TODO
		}
		return
	}

	switch p.ResultCode() {
	case "SUCCESS":
		resp.Respcd, resp.ErrorCode = SuccessCode, SuccessMsg
	case "UNKNOW":
		resp.Respcd, resp.ErrorCode = InprocessCode, InprocessMsg
	default:
		resp.Respcd, resp.ErrorCode = UnKnownCode, UnKnownMsg
	}

}
