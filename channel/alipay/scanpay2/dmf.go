package scanpay2

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"time"
)

var DefaultClient dmf2

var NotifyPath = "/scanpay/upNotify/alipay2"
var Alipay2NotifyUrl = goconf.Config.AlipayScanPay.NotifyUrl + NotifyPath

type dmf2 struct{}

func getCommonParams(m *model.ScanPayRequest) *CommonParams {
	return &CommonParams{
		AppID:      m.AppID,
		PrivateKey: LoadPrivateKey(m.PemKey), // TODO 做个缓存处理
		Req:        m,
		// TODO 预留authToken
	}
}

func (d *dmf2) ProcessBarcodePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	p := &PayReq{}
	p.CommonParams = *getCommonParams(req)
	// p.CommonParams.NotifyUrl = Alipay2NotifyUrl
	p.OutTradeNo = req.OrderNum
	p.Scene = "bar_code"
	p.AuthCode = req.ScanCodeId
	p.Subject = req.Subject
	p.TotalAmount = req.ActTxamt
	p.GoodsDetail = parseGoods(req)
	_, p.TimeExpire = handleExpireTime(req.TimeExpire)
	p.ExtendParams = Params{req.ExtendParams}
	p.Body = ""
	p.StoreID = req.M.Detail.ShopID
	p.OperatorID = ""
	p.TerminalID = ""

	q := &PayResp{}
	err := Execute(p, q)
	if err != nil {
		return nil, err
	}
	ret := &model.ScanPayResponse{}
	ret.Respcd, ret.ErrorDetail = transform("pay", q.Code, q.Msg, q.SubCode, q.SubMsg)
	ret.ChannelOrderNum = q.TradeNo
	ret.PayTime = q.GmtPayment
	ret.ConsumerAccount = q.BuyerLogonID
	// ret.ConsumerId = q.OpenID
	// TODO...金额如何处理

	return ret, err
}

func (d *dmf2) ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	p := &PrecreateReq{}
	p.CommonParams = *getCommonParams(req)
	p.CommonParams.NotifyUrl = Alipay2NotifyUrl
	p.OutTradeNo = req.OrderNum
	p.Subject = req.Subject
	p.TotalAmount = req.ActTxamt
	p.GoodsDetail = parseGoods(req)
	_, p.TimeExpire = handleExpireTime(req.TimeExpire)
	p.ExtendParams = Params{req.ExtendParams}
	p.StoreID = req.M.Detail.ShopID
	q := &PrecreateResp{}
	err := Execute(p, q)
	if err != nil {
		return nil, err
	}
	ret := &model.ScanPayResponse{}
	ret.Respcd, ret.ErrorDetail = transform("precreate", q.Code, q.Msg, q.SubCode, q.SubMsg)
	ret.QrCode = q.QrCode

	return ret, err
}

func (d *dmf2) ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	p := &RefundReq{}
	p.CommonParams = *getCommonParams(req)
	p.RefundAmount = req.ActTxamt
	p.TradeNo = req.OrigChanOrderNum
	p.OutRequestNo = req.OrderNum

	// RefundReason  string `json:"refund_reason,omitempty"`          // 退款原因
	// StoreID       string `json:"store_id,omitempty"`               // 商户的门店编号
	// AlipayStoreID string `json:"alipay_store_id,omitempty"`        // 支付宝店铺编号
	// TerminalID    string `json:"terminal_id,omitempty"`            // 商户的终端编号

	q := &RefundResp{}
	err := Execute(p, q)
	if err != nil {
		return nil, err
	}
	ret := &model.ScanPayResponse{}
	ret.Respcd, ret.ErrorDetail = transform("refund", q.Code, q.Msg, q.SubCode, q.SubMsg)
	ret.PayTime = q.GmtRefundPay
	ret.ConsumerAccount = q.BuyerLogonID
	// ret.ConsumerId = q.OpenID

	// TradeNo              string `json:"trade_no,omitempty"`       // 支付宝交易号
	// OutTradeNo           string `json:"out_trade_no,omitempty"`   // 商户订单号
	// FundChange           string `json:"fund_change"`              // 本次退款请求是否发生资金变动
	// RefundFee            string `json:"fund_change,omitempty"`    // 累计退款金额
	// RefundDetailItemList []struct {
	// 	FundChannel string `json:"fund_channel,omitempty"` // 支付渠道,例 COUPON、DISCOUNT
	// 	Amount      string `json:"amount,omitempty"`       // 支付金额
	// } `json:"refund_detail_item_list"` // 退款资金明细信息集合

	return ret, err
}

func (d *dmf2) ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	p := &QueryReq{}
	p.CommonParams = *getCommonParams(req)
	p.OutTradeNo = req.OrigOrderNum

	q := &QueryResp{}
	err := Execute(p, q)
	if err != nil {
		return nil, err
	}

	ret := &model.ScanPayResponse{}
	// 查询操作如果不成功，返回09
	if q.Code != "10000" {
		ret.Respcd, ret.ErrorDetail = inprocessCode, inprocessMsg
		return ret, nil
	}

	switch q.TradeStatus {
	case "TRADE_SUCCESS":
		ret.Respcd, ret.ErrorDetail = successCode, successMsg
		ret.ChannelOrderNum = q.TradeNo
		// ret.PayTime = q.GmtPayment
		ret.ConsumerAccount = q.BuyerLogonID
		// ret.ConsumerId = q.OpenID
		ret.PayTime = q.SendPayDate
	case "WAIT_BUYER_PAY":
		ret.Respcd, ret.ErrorDetail = inprocessCode, inprocessMsg
	case "TRADE_FINISHED", "TRADE_CLOSED":
		ret.Respcd, ret.ErrorDetail = closeCode, closeMsg
	}

	return ret, err
}

func (d *dmf2) ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	p := &CancelReq{}
	p.CommonParams = *getCommonParams(req)
	p.OutTradeNo = req.OrigOrderNum

	q := &CancelResp{}
	err := Execute(p, q)
	if err != nil {
		return nil, err
	}
	ret := &model.ScanPayResponse{}
	ret.Respcd, ret.ErrorDetail = transform("cancel", q.Code, q.Msg, q.SubCode, q.SubMsg)
	ret.ChannelOrderNum = q.TradeNo

	return ret, err
}

// ProcessClose 取消走撤销接口
func (d *dmf2) ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return d.ProcessCancel(req)
}

func (d *dmf2) ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, fmt.Errorf("%s", "not support yet!")
}

func handleExpireTime(expirtTime string) (string, string) {

	layout := "2006-01-02 15:04:05"
	startTime := time.Now()
	defaultEntTime := startTime.Add(24 * time.Hour)

	var stStr, etStr = startTime.Format(layout), defaultEntTime.Format(layout)

	if expirtTime == "" {
		return stStr, etStr
	}

	et, err := time.ParseInLocation("20060102150405", expirtTime, time.Local)
	if err != nil {
		return stStr, etStr
	}

	d := et.Sub(startTime)
	if d < 5*time.Minute {
		return stStr, startTime.Add(5 * time.Minute).Format(layout)
	}

	return stStr, expirtTime
}

// parseGoods 输出2.0要求的商品格式
func parseGoods(req *model.ScanPayRequest) []GoodsDetail {
	details, err := req.MarshalGoods()
	if err != nil {
		return nil
	}
	if len(details) > 0 {
		var gs []GoodsDetail
		for _, g := range details {
			gs = append(gs, GoodsDetail{
				GoodsId:   fmt.Sprintf("%d", g.GoodsId),
				GoodsName: g.GoodsName,
				Price:     g.Price,
				Quantity:  g.Quantity,
			})
		}
		return gs
	}
	return nil
}
