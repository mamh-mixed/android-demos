package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
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

	// 请求微信
	err = weixin.Execute(q, p)
	if err != nil {
		return nil, err
	}

	// 如果是系统错误，重试
	if p.ErrCode == "SYSTEMERROR" {

		query := &EnterpriseQueryReq{
			CommonParams:   *getCommonParams(req),
			AppId:          req.AppID,
			MchId:          req.ChanMerId,
			PartnerTradeNo: req.OrderNum,
		}
		resp := &EnterpriseQueryResp{}
		var queryDuration = []time.Duration{3 * time.Second, 6 * time.Second, 9 * time.Second, 12 * time.Second}
	Tag:
		for i, d := range queryDuration {
			time.Sleep(d)
			// query
			weixin.Execute(query, resp)
			switch resp.Status {
			case "PROCESSING":
				log.Infof("enterprise query %d times:", i+1)
			case "SUCCESS":
				p.ReturnCode, p.ResultCode = "SUCCESS", "SUCCESS"
				break Tag
			case "FAILED":
				p.ReturnCode, p.ResultCode = "FAIL", "FAIL"
				break Tag
			}
		}
	}

	// p.ResultCode = "FAIL"
	// p.ErrCode = "DEBUG_ERROR"
	// p.ErrCodeDes = "企业付款调试错误"

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

	q := &EnterpriseQueryReq{
		CommonParams:   *getCommonParams(req),
		AppId:          req.AppID,
		MchId:          req.ChanMerId,
		PartnerTradeNo: req.OrigOrderNum,
	}
	p := &EnterpriseQueryResp{}

	err = weixin.Execute(q, p)
	if err != nil {
		return nil, err
	}

	log.Debugf("%+v", p)

	status, msg, ec := weixin.Transform("enterpriseQuery", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	ret = &model.ScanPayResponse{
		Respcd:          status,    // 交易结果  M
		ErrorDetail:     msg,       // 错误信息   C
		ChanRespCode:    p.ErrCode, // 渠道详细应答码
		ChannelOrderNum: p.DetailId,
		ErrorCode:       ec,
	}

	return ret, nil
}
