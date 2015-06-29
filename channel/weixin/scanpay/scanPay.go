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
func (sp *WeixinScanPay) ProcessBarcodePay(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	d := &PayReq{
		// 必填
		Appid:          m.AppID,         // 公众账号ID
		MchID:          m.ChanMerId,     // 商户号
		SubMchId:       m.SubMchId,      // 子商户号
		NonceStr:       tools.Nonce(16), // 随机字符串
		Body:           m.Subject,       // 商品描述
		OutTradeNo:     m.OrderNum,      // 商户订单号
		Sign:           "",              // 签名
		AuthCode:       m.ScanCodeId,    // 授权码
		TotalFee:       m.ActTxamt,      // 总金额
		WeixinMD5Key:   m.SignCert,      // md5key
		SpbillCreateIP: tools.LocalIP,   // 终端IP

		// 非必填
		DeviceInfo: m.DeviceInfo,     // 设备号
		Detail:     m.MarshalGoods(), // 商品详情
		Attach:     m.Attach,         // 附加数据
		FeeType:    m.CurrType,       // 货币类型
		GoodsGag:   m.GoodsGag,       // 商品标记
	}

	p := &PayResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform(p.ReturnCode, p.ReturnMsg, p.ResultCode, p.ErrCode)
	ret = &model.ScanPayResponse{
		Txndir:          "A",             // 交易方向 M M
		Busicd:          m.Busicd,        // 交易类型 M M
		Respcd:          status,          // 交易结果  M
		Inscd:           m.Inscd,         // 机构号 M M
		Chcd:            m.Chcd,          // 渠道 C C
		Mchntid:         p.MchID,         // 商户号 M M
		Txamt:           p.TotalFee,      // 订单金额 M M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: p.OpenID,        // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		OrderNum:        m.OrderNum,      // 订单号 M C
		OrigOrderNum:    "",              // 源订单号 M C
		QrCode:          m.ScanCodeId,    // 二维码 C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
	}

	return ret, err
}

// ProcessEnquiry 查询
func (sp *WeixinScanPay) ProcessEnquiry(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	d := &PayQueryReq{
		Appid:         m.AppID,         // 公众账号ID
		MchID:         m.ChanMerId,     // 商户号
		SubMchId:      m.SubMchId,      // 子商户号
		TransactionId: "",              // 微信支付订单号
		OutTradeNo:    m.OrderNum,      // 商户订单号
		NonceStr:      tools.Nonce(32), // 商品详情
		Sign:          "",
		WeixinMD5Key:  m.SignCert,
	}

	p := &PayQueryResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform(p.ReturnCode, p.ReturnMsg, p.ResultCode, p.ErrCode)
	ret = &model.ScanPayResponse{
		Txndir:          "A",             // 交易方向 M M
		Busicd:          m.Busicd,        // 交易类型 M M
		Respcd:          status,          // 交易结果  M
		Inscd:           m.Inscd,         // 机构号 M M
		Chcd:            m.Chcd,          // 渠道 C C
		Mchntid:         p.MchID,         // 商户号 M M
		Txamt:           p.TotalFee,      // 订单金额 M M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: p.OpenID,        // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		OrderNum:        m.OrderNum,      // 订单号 M C
		OrigOrderNum:    m.OrigOrderNum,  // 源订单号 M C
		QrCode:          m.ScanCodeId,    // 二维码 C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
	}

	return ret, err
}

// ProcessQrCodeOfflinePay 扫二维码预下单
func (sp *WeixinScanPay) ProcessQrCodeOfflinePay(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	d := &PrePayReq{

		// 公共字段
		Appid:    m.AppID,         // 公众账号ID
		MchID:    m.ChanMerId,     // 商户号
		SubMchId: m.SubMchId,      // 子商户号
		NonceStr: tools.Nonce(32), // 随机字符串
		Sign:     "",              // 签名

		WeixinMD5Key: m.SignCert, // md5key

		DeviceInfo:     m.DeviceInfo,     // 设备号
		Body:           m.Subject,        // 商品描述
		Detail:         m.MarshalGoods(), // 商品详情
		Attach:         m.Attach,         // 附加数据
		OutTradeNo:     m.OrderNum,       // 商户订单号
		TotalFee:       m.ActTxamt,       // 总金额
		FeeType:        m.CurrType,       // 货币类型
		SpbillCreateIP: tools.LocalIP,    // 终端IP
		TimeStart:      "",               // 交易起始时间
		TimeExpire:     "",               // 交易结束时间
		GoodsGag:       m.GoodsGag,       // 商品标记
		NotifyURL:      weixinNotifyURL,  // 通知地址
		TradeType:      "NATIVE",         // 交易类型
		ProductID:      "",               // 商品ID
		Openid:         "",               // 用户标识
	}

	p := &PrePayResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform(p.ReturnCode, p.ReturnMsg, p.ResultCode, p.ErrCode)
	ret = &model.ScanPayResponse{
		Txndir:          "A",        // 交易方向 M M
		Busicd:          m.Busicd,   // 交易类型 M M
		Respcd:          status,     // 交易结果  M
		Inscd:           m.Inscd,    // 机构号 M M
		Chcd:            m.Chcd,     // 渠道 C C
		Mchntid:         p.MchID,    // 商户号 M M
		Txamt:           m.ActTxamt, // 订单金额 M M
		ChannelOrderNum: "",         // 渠道交易号 C
		ConsumerAccount: "",         // 渠道账号  C
		ConsumerId:      "",         // 渠道账号ID   C
		ErrorDetail:     msg,        // 错误信息   C
		OrderNum:        m.OrderNum, // 订单号 M C
		OrigOrderNum:    "",         // 源订单号 M C
		QrCode:          "",         // 二维码 C
		ChanRespCode:    p.ErrCode,  // 渠道详细应答码
	}

	return ret, err
}

// ProcessRefund 退款
func (sp *WeixinScanPay) ProcessRefund(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	log.Debugf("%#c", m)
	d := &RefundReq{
		// 公共字段
		Appid:        m.AppID,         // 公众账号ID
		MchID:        m.ChanMerId,     // 商户号
		SubMchId:     m.SubMchId,      // 子商户号
		NonceStr:     tools.Nonce(16), // 随机字符串
		Sign:         "",              // 签名
		WeixinMD5Key: m.SignCert,

		DeviceInfo:    m.DeviceInfo,   // 设备号
		TransactionId: "",             // 微信订单号
		OutTradeNo:    m.OrigOrderNum, // 商户订单号
		OutRefundNo:   m.OrderNum,     // 商户退款单号
		TotalFee:      m.TotalTxamt,   // 总金额
		RefundFee:     m.ActTxamt,     // 退款金额
		RefundFeeType: m.CurrType,     // 货币种类
		OpUserId:      m.ChanMerId,    // 操作员
	}

	p := &RefundResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform(p.ReturnCode, p.ReturnMsg, p.ResultCode, p.ErrCode)
	ret = &model.ScanPayResponse{
		Txndir:          "A",             // 交易方向 M M
		Busicd:          m.Busicd,        // 交易类型 M M
		Respcd:          status,          // 交易结果  M
		Inscd:           m.Inscd,         // 机构号 M M
		Chcd:            m.Chcd,          // 渠道 C C
		Mchntid:         p.MchID,         // 商户号 M M
		Txamt:           p.RefundFee,     // 订单金额 M M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: "",              // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		OrderNum:        m.OrderNum,      // 订单号 M C
		OrigOrderNum:    "",              // 源订单号 M C
		QrCode:          m.ScanCodeId,    // 二维码 C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
	}

	return ret, err
}

// ProcessRefundQuery 退款查询
func (sp *WeixinScanPay) ProcessRefundQuery(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	d := &RefundQueryReq{
		// 公共字段
		Appid:        m.AppID,         // 公众账号ID
		MchID:        m.ChanMerId,     // 商户号
		SubMchId:     m.SubMchId,      // 子商户号
		NonceStr:     tools.Nonce(16), // 随机字符串
		Sign:         "",              // 签名
		WeixinMD5Key: m.SignCert,

		DeviceInfo:    m.DeviceInfo,   // 设备号
		TransactionId: "",             // 微信订单号
		OutTradeNo:    m.OrderNum,     // 商户订单号
		OutRefundNo:   m.OrigOrderNum, // 商户退款单号
		RefundId:      "",             // 操作员
	}

	p := &RefundQueryResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform(p.ReturnCode, p.ReturnMsg, p.ResultCode, p.ErrCode)
	ret = &model.ScanPayResponse{
		Txndir:          "A",             // 交易方向 M M
		Busicd:          m.Busicd,        // 交易类型 M M
		Respcd:          status,          // 交易结果  M
		Inscd:           m.Inscd,         // 机构号 M M
		Chcd:            m.Chcd,          // 渠道 C C
		Mchntid:         p.MchID,         // 商户号 M M
		Txamt:           m.Txamt,         // 订单金额 M M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: "",              // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		OrderNum:        m.OrderNum,      // 订单号 M C
		OrigOrderNum:    "",              // 源订单号 M C
		QrCode:          m.ScanCodeId,    // 二维码 C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
	}

	return ret, err
}

// ProcessCancel 撤销
func (sp *WeixinScanPay) ProcessCancel(m *model.ScanPay) (ret *model.ScanPayResponse, err error) {
	d := &ReverseReq{
		// 公共字段
		Appid:        m.AppID,         // 公众账号ID
		MchID:        m.ChanMerId,     // 商户号
		SubMchId:     m.SubMchId,      // 子商户号
		NonceStr:     tools.Nonce(16), // 随机字符串
		Sign:         "",              // 签名
		WeixinMD5Key: m.SignCert,

		TransactionId: "",         // 微信订单号
		OutTradeNo:    m.OrderNum, // 商户订单号
	}

	p := &ReverseResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform(p.ReturnCode, p.ReturnMsg, p.ResultCode, p.ErrCode)
	ret = &model.ScanPayResponse{
		Txndir:          "A",          // 交易方向 M M
		Busicd:          m.Busicd,     // 交易类型 M M
		Respcd:          status,       // 交易结果  M
		Inscd:           m.Inscd,      // 机构号 M M
		Chcd:            m.Chcd,       // 渠道 C C
		Mchntid:         p.MchID,      // 商户号 M M
		Txamt:           m.Txamt,      // 订单金额 M M
		ChannelOrderNum: "",           // 渠道交易号 C
		ConsumerAccount: "",           // 渠道账号  C
		ConsumerId:      "",           // 渠道账号ID   C
		ErrorDetail:     msg,          // 错误信息   C
		OrderNum:        m.OrderNum,   // 订单号 M C
		OrigOrderNum:    "",           // 源订单号 M C
		QrCode:          m.ScanCodeId, // 二维码 C
		ChanRespCode:    p.ErrCode,    // 渠道详细应答码
	}

	return ret, err
}
