package weixin

import "github.com/CardInfoLink/quickpay/model"

//"github.com/omigo/log"

var DefaultClient WeixinPay

const (
	MicroPay = iota
	OrderQuery
	RefundQuery
)

// ProcessBarcodePay 扫条码下单
func (c *WeixinPay) ProcessBarcodePay(scanPayReq *model.ScanPay) *model.ScanPayResponse {
	micropayReq := c.createRequestData(MicroPay)
	micropayReq.copyData(scanPayReq)
	setSign(micropayReq, calculateSign(micropayReq, md5Key))

	microPayResp := c.requestWeixin(micropayReq, scanPayReq.NotifyUrl)

	//log.Debugf("micropay response: %+v", buf)
	return microPayResp.convertToScanPayResp()
}

// ProcessEnquiry
func (c *WeixinPay) ProcessEnquiry(scanPayReq *model.ScanPay) *model.ScanPayResponse {
	orderqueryReq := c.createRequestData(OrderQuery)
	orderqueryReq.copyData(scanPayReq)
	setSign(orderqueryReq, calculateSign(orderqueryReq, md5Key))

	orderqueryResp := c.requestWeixin(orderqueryReq, scanPayReq.NotifyUrl)

	//log.Debugf("micropay response: %+v", buf)
	return orderqueryResp.convertToScanPayResp()
}
