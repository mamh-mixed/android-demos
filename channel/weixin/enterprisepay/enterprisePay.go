package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// 微信企业支付
var DefaultClient WeixinEnterprisePay

// WeixinEnterprisePay 微信企业支付
type WeixinEnterprisePay struct{}

// ProcessPay 支付
func (w *WeixinEnterprisePay) ProcessPay(req *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {

	pay := &EnterprisePayReq{
		MchAappid:      req.AppID,
		MchID:          req.ChanMerId,
		NonceStr:       util.Nonce(32),
		OpenId:         req.OpenId,
		CheckName:      req.CheckName,
		ReUserName:     req.UserName,
		Desc:           req.Desc,
		SpbillCreateIp: util.LocalIP,
		PartnerTradeNo: req.OrderNum,
		WeixinMD5Key:   req.SignCert,
		Amount:         req.ActTxamt,
	}
	resp := &EnterprisePayResp{}

	err = request(pay, resp)
	if err != nil {
		return nil, err
	}
	log.Debugf("%+v", resp)

	ret = &model.ScanPayResponse{}

	return ret, nil
}

// ProcessEnquiry 支付
func (w *WeixinEnterprisePay) ProcessEnquiry(req *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {

	pay := &EnterprisePayReq{}
	resp := &EnterprisePayResp{}

	err = request(pay, resp)
	if err != nil {
		return nil, err
	}

	ret = &model.ScanPayResponse{}

	return ret, nil
}
