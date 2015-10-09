package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"time"
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
		Req:          req,
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
	var retry int
	for {
		err = weixin.Execute(q, p)
		if err != nil {
			return nil, err
		}

		// 如果是系统错误，重试
		if p.ErrCode == "SYSTEMERROR" {
			log.Warnf("enterprisepay weixin return SYSTEMERROR , retry ..., orderNum=%s,merId=%s", req.OrderNum, req.Mchntid)
			retry++
			if retry == 3 {
				log.Error("enterprisepay retry 3 times, break.")
				p.ReturnCode, p.ResultCode = "SUCCESS", "SUCCESS"
				break
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	status, msg, ec := weixin.Transform("enterprisePay", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	ret = &model.ScanPayResponse{
		Respcd:          status,    // 交易结果  M
		ErrorDetail:     msg,       // 错误信息   C
		ChanRespCode:    p.ErrCode, // 渠道详细应答码
		ChannelOrderNum: p.PaymentNo,
		ErrorCode:       ec,
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
