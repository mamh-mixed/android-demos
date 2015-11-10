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
	// log.Debugf("payResp: %+v", p)

	resp := &model.ScanPayResponse{}
	resp.ChannelOrderNum = p.Response.Alipay.AlipayTransId
	resp.PayTime = p.Response.Alipay.AlipayPayTime
	resp.ConsumerId = p.Response.Alipay.AlipayBuyerLoginId
	// TODO...

	// result
	alipayResponseHandle(p, resp, "createandpay")

	return resp, nil
}

func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, errors.New("Not support yet!")
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessEnquiry 查询
func (a *alp) ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessCancel 撤销
func (a *alp) ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

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

// ProcessClose 关闭
func (a *alp) ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessRefundQuery 退款查询
func (a *alp) ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

func alipayResponseHandle(p scanpay1.BaseResp, resp *model.ScanPayResponse, service string) {

	var errorCode string = p.ErrorCode()
	if p.ReqFlag() {
		switch p.ResultCode() {
		case "SUCCESS":
			resp.Respcd = "00"
			resp.ErrorCode = "成功"
			return
		case "UNKNOW":
			// TODO: unknow 如何处理
			if errorCode == "" {
				resp.Respcd = "09"
				resp.ErrorCode = "处理中"
				return
			}
		default:
		}
	}
	// error
	spCode := mongo.ScanPayRespCol.GetByAlp(errorCode, service)
	resp.Respcd = spCode.ISO8583Code
	resp.ErrorDetail = spCode.ErrorCode
	if !spCode.IsUseISO {
		resp.ErrorDetail = errorCode // TODO
	}
}
