package scanpay2

import (
	"github.com/CardInfoLink/quickpay/model"
	"time"
)

var DefaultClient dmf2

type dmf2 struct{}

func getCommonParams(m *model.ScanPayRequest) *CommonParams {
	return &CommonParams{
		AppID: "2014122500021754",
		// PrivateKey: LoadPrivateKey([]byte(privateKeyPem)), // TODO 做个缓存处理
		Req: m,
	}
}

func (d *dmf2) ProcessBarcodePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	p := &PayReq{}
	p.CommonParams = *getCommonParams(req)
	p.OutTradeNo = req.OrderNum
	p.Scene = "bar_code"
	p.AuthCode = req.ScanCodeId
	p.Subject = req.Subject
	p.TotalAmount = req.ActTxamt
	p.GoodsDetail = req.AlpMarshalGoods()
	_, p.TimeExpire = handleExpireTime(req.TimeExpire)

	p.Body = ""
	p.StoreID = ""
	p.OperatorID = ""
	p.TerminalID = ""
	p.ExtendParams = ""

	q := &PayResp{}
	err := Execute(p, q)
	if err != nil {
		return nil, err
	}
	ret := &model.ScanPayResponse{}
	ret.Respcd, ret.ErrorDetail = transform("pay", q.Code, q.Msg, q.SubCode, q.SubMsg)
	ret.ChannelOrderNum = q.TradeNo
	ret.PayTime = q.GmtPayment
	// TODO...

	return ret, err
}

func (d *dmf2) ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	ret := new(model.ScanPayResponse)
	var err error

	return ret, err
}

func (d *dmf2) ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	ret := new(model.ScanPayResponse)
	var err error

	return ret, err
}

func (d *dmf2) ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	ret := new(model.ScanPayResponse)
	var err error

	return ret, err
}

func (d *dmf2) ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	ret := new(model.ScanPayResponse)
	var err error

	return ret, err
}

func (d *dmf2) ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	ret := new(model.ScanPayResponse)
	var err error

	return ret, err
}

func (d *dmf2) ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	ret := new(model.ScanPayResponse)
	var err error

	return ret, err
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
