package scanpay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

// WeixinScanPay 微信扫码支付
type WeixinScanPay struct{}

// DefaultWeixinScanPay 微信扫码支付默认实现
var DefaultWeixinScanPay WeixinScanPay

// ProcessBarcodePay 扫条码下单
func (p *WeixinScanPay) ProcessBarcodePay(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	// TODO validate params...

	d := &ScanPayReqData{
		// 必填
		Appid:          m.AppID,         // 公众账号ID
		MchID:          m.ChanMerId,     // 商户号
		NonceStr:       tools.Nonce(16), // 随机字符串
		Body:           m.Subject,       // 商品描述
		OutTradeNo:     m.OrderNum,      // 商户订单号
		Sign:           "",              // 签名
		AuthCode:       m.ScanCodeId,    // 授权码
		TotalFee:       m.ActTxamt,      // 总金额
		WeixinMD5Key:   m.SignCert,      // md5key
		SpbillCreateIP: tools.LocalIP,   // 终端IP

		// 非必填
		SubMchId:   m.SubMchId,       // 子商户号
		DeviceInfo: m.DeviceInfo,     // 设备号
		Detail:     m.MarshalGoods(), // 商品详情
		Attach:     m.Attach,         // 附加数据
		FeeType:    m.CurrType,       // 货币类型
		GoodsGag:   m.GoodsGag,       // 商品标记

	}

	var respData ScanPayRespData
	err = sendRequest(ScanPayURI, d, &respData)
	if err != nil {
		log.Errorf("weixin device scan phone request error: %s", err)
		return nil, err
	}

	ret = &model.ScanPayResponse{
		Txndir:          "A",                 // 交易方向 M M
		Busicd:          m.Busicd,            // 交易类型 M M
		Respcd:          respData.ReturnCode, // 交易结果  M
		Inscd:           m.Inscd,             // 机构号 M M
		Chcd:            m.Chcd,              // 渠道 C C
		Mchntid:         respData.MchID,      // 商户号 M M
		Txamt:           respData.TotalFee,   // 订单金额 M M
		ChannelOrderNum: respData.OutTradeNo, // 渠道交易号 C
		ConsumerAccount: m.ChanMerId,         // 渠道账号  C
		ConsumerId:      respData.OpenID,     // 渠道账号ID   C
		ErrorDetail:     respData.ErrCodeDes, // 错误信息   C
		OrderNum:        m.OrderNum,          // 订单号 M C
		OrigOrderNum:    m.OrigOrderNum,      // 源订单号 M C
		QrCode:          m.ScanCodeId,        // 二维码 C
		ChanRespCode:    respData.ReturnCode, // 渠道详细应答码
	}

	// TODO 应答码转换
	return ret, err
}

// ProcessEnquiry 查询
func (p *WeixinScanPay) ProcessEnquiry(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	// TODO validate params...

	d := &ScanPayQueryReqData{
		Appid:         m.AppID,         // 公众账号ID
		MchID:         m.ChanMerId,     // 商户号
		SubMchId:      m.SubMchId,      // 子商户号
		TransactionId: "",              // 微信支付订单号
		OutTradeNo:    m.OrderNum,      // 商户订单号
		NonceStr:      tools.Nonce(32), // 商品详情
		Sign:          "",
		WeixinMD5Key:  m.SignCert,
	}

	var respData ScanPayQueryRespData
	err = sendRequest(ScanPayQueryURI, d, &respData)
	if err != nil {
		log.Errorf("weixin device scan phone request error: %s", err)
		return nil, err
	}

	ret = &model.ScanPayResponse{
		Txndir:          "A",                 // 交易方向 M M
		Busicd:          m.Busicd,            // 交易类型 M M
		Respcd:          respData.ReturnCode, // 交易结果  M
		Inscd:           m.Inscd,             // 机构号 M M
		Chcd:            m.Chcd,              // 渠道 C C
		Mchntid:         respData.MchID,      // 商户号 M M
		Txamt:           respData.TotalFee,   // 订单金额 M M
		ChannelOrderNum: respData.OutTradeNo, // 渠道交易号 C
		ConsumerAccount: m.ChanMerId,         // 渠道账号  C
		ConsumerId:      respData.OpenID,     // 渠道账号ID   C
		ErrorDetail:     respData.ErrCodeDes, // 错误信息   C
		OrderNum:        m.OrderNum,          // 订单号 M C
		OrigOrderNum:    m.OrigOrderNum,      // 源订单号 M C
		QrCode:          m.ScanCodeId,        // 二维码 C
		ChanRespCode:    respData.ReturnCode, // 渠道详细应答码
	}

	// TODO 应答码转换
	return ret, err
}

// ProcessQrCodeOfflinePay 扫二维码预下单
func (p *WeixinScanPay) ProcessQrCodeOfflinePay(req *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	return
}

// ProcessRefund 退款
func (p *WeixinScanPay) ProcessRefund(req *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	return
}

// ProcessCancel 撤销
func (p *WeixinScanPay) ProcessCancel(req *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	return
}
