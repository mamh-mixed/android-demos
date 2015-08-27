package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// 微信企业支付
var DefaultClient WeixinEnterprisePay

// WeixinEnterprisePay 微信企业支付
type WeixinEnterprisePay struct{}

func getCommonParams(req *model.ScanPayRequest) *weixin.CommonParams {
	return &weixin.CommonParams{
		NonceStr:     util.Nonce(32), // 随机字符串
		Sign:         "",             // 签名
		WeixinMD5Key: req.SignKey,    // md5key
		ClientCert:   req.WeixinClientCert,
		ClientKey:    req.WeixinClientKey,
	}
}

// ProcessPay 支付
func (w *WeixinEnterprisePay) ProcessPay(req *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {

	q := &EnterprisePayReq{
		CommonParams:   *getCommonParams(req),
		MchID:          req.ChanMerId, // 商户号
		MchAappid:      req.AppID,
		OpenId:         req.OpenId,
		CheckName:      req.CheckName,
		ReUserName:     req.UserName,
		Desc:           req.Desc,
		SpbillCreateIp: util.LocalIP,
		PartnerTradeNo: req.OrderNum,
		Amount:         req.ActTxamt,
	}
	p := &EnterprisePayResp{}

	err = weixin.Execute(q, p)
	if err != nil {
		return nil, err
	}
	log.Debugf("%+v", p)

	status, msg := weixin.Transform("enterprisePay", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	ret = &model.ScanPayResponse{
		Respcd:          status,    // 交易结果  M
		ErrorDetail:     msg,       // 错误信息   C
		ChanRespCode:    p.ErrCode, // 渠道详细应答码
		ChannelOrderNum: p.PaymentNo,
	}

	return ret, nil
}

// ProcessEnquiry 支付
func (w *WeixinEnterprisePay) ProcessEnquiry(req *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {

	query := &EnterprisePayReq{
		CommonParams: *getCommonParams(req),
	}
	resp := &EnterprisePayResp{}

	err = weixin.Execute(query, resp)
	if err != nil {
		return nil, err
	}

	ret = &model.ScanPayResponse{}

	return ret, nil
}
