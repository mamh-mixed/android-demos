package scanpay

import (
	"fmt"
	"strconv"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// WeixinScanPay 微信扫码支付
type WeixinScanPay struct{}

// DefaultWeixinScanPay 微信扫码支付默认实现
var DefaultWeixinScanPay WeixinScanPay

func getCommonParams(m *model.ScanPayRequest) *CommonParams {
	return &CommonParams{
		Appid:    m.AppID,        // 公众账号ID
		MchID:    m.ChanMerId,    // 商户号
		SubMchId: m.SubMchId,     // 子商户号
		NonceStr: util.Nonce(32), // 随机字符串
		Sign:     "",             // 签名

		WeixinMD5Key: m.SignCert, // md5key
	}
}

// ProcessBarcodePay 扫条码下单
func (sp *WeixinScanPay) ProcessBarcodePay(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	d := &PayReq{
		CommonParams: *getCommonParams(m),

		Body:           m.Subject,    // 商品描述
		OutTradeNo:     m.OrderNum,   // 商户订单号
		AuthCode:       m.ScanCodeId, // 授权码
		TotalFee:       m.ActTxamt,   // 总金额
		SpbillCreateIP: util.LocalIP, // 终端IP

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

	status, msg := transform("pay", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	ret = &model.ScanPayResponse{
		Respcd:          status,          // 交易结果  M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: p.SubOpenid,     // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
	}

	ret.MerDiscount, ret.ChcdDiscount = "0.00", "0.00"
	if p.CouponFee != "" {
		f, _ := strconv.ParseFloat(p.CouponFee, 64)
		ret.MerDiscount = fmt.Sprintf("%0.2f", f/100)
	}

	return ret, err
}

// ProcessEnquiry 查询
func (sp *WeixinScanPay) ProcessEnquiry(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	d := &PayQueryReq{
		CommonParams: *getCommonParams(m),

		TransactionId: "",         // 微信支付订单号
		OutTradeNo:    m.OrderNum, // 商户订单号
	}

	p := &PayQueryResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	log.Debugf("ProcessEnquiry response data is %#v", p)

	status, msg := transform("payQuery", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	// 如果返回的是成功的，要对 trade_state 做判断
	// SUCCESS—支付成功;REFUND—转入退款;NOTPAY—未支付;CLOSED—已关闭;REVOKED—已撤销;
	// USERPAYING-用户支付中;PAYERROR-支付失败(其他原因，如银行返回失败)
	if status == "00" {
		respCode := mongo.ScanPayRespCol.GetByWxp(p.TradeState, "payQuery")
		status, msg = respCode.ISO8583Code, respCode.ISO8583Msg
	}

	ret = &model.ScanPayResponse{
		Respcd:          status,          // 交易结果  M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: p.SubOpenid,     // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
	}

	ret.MerDiscount, ret.ChcdDiscount = "0.00", "0.00"
	if p.CouponFee != "" {
		f, _ := strconv.ParseFloat(p.CouponFee, 64)
		ret.MerDiscount = fmt.Sprintf("%0.2f", f/100)
	}

	return ret, err
}

// ProcessQrCodeOfflinePay 扫二维码预下单
func (sp *WeixinScanPay) ProcessQrCodeOfflinePay(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	d := &PrePayReq{
		CommonParams: *getCommonParams(m),

		DeviceInfo:     m.DeviceInfo,     // 设备号
		Body:           m.Subject,        // 商品描述
		Detail:         m.MarshalGoods(), // 商品详情
		Attach:         m.SysOrderNum,    // 附加数据 这里送系统订单号
		OutTradeNo:     m.OrderNum,       // 商户订单号
		TotalFee:       m.ActTxamt,       // 总金额
		FeeType:        m.CurrType,       // 货币类型
		SpbillCreateIP: util.LocalIP,     // 终端IP
		TimeStart:      "",               // 交易起始时间
		TimeExpire:     "",               // 交易结束时间
		GoodsGag:       m.GoodsGag,       // 商品标记
		NotifyURL:      m.NotifyUrl,      // 通知地址
		TradeType:      "NATIVE",         // 交易类型
		ProductID:      "",               // 商品ID
		Openid:         "",               // 用户标识
	}

	p := &PrePayResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform("prePay", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	ret = &model.ScanPayResponse{
		Respcd:          status,    // 交易结果  M
		ChannelOrderNum: "",        // 渠道交易号 C
		ConsumerAccount: "",        // 渠道账号  C
		ConsumerId:      "",        // 渠道账号ID   C
		ErrorDetail:     msg,       // 错误信息   C
		QrCode:          p.CodeURL, // 二维码 C
		ChanRespCode:    p.ErrCode, // 渠道详细应答码
	}

	return ret, err
}

// ProcessRefund 退款
func (sp *WeixinScanPay) ProcessRefund(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	log.Debugf("%#c", m)
	d := &RefundReq{
		CommonParams: *getCommonParams(m),

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

	status, msg := transform("refund", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
	ret = &model.ScanPayResponse{
		Respcd:          status,          // 交易结果  M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: "",              // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		QrCode:          m.ScanCodeId,    // 二维码 C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
	}

	return ret, err
}

// ProcessRefundQuery 退款查询
func (sp *WeixinScanPay) ProcessRefundQuery(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	d := &RefundQueryReq{
		CommonParams: *getCommonParams(m),

		DeviceInfo:    m.DeviceInfo,   // 设备号
		TransactionId: "",             // 微信订单号
		OutTradeNo:    m.OrigOrderNum, // 商户订单号
		OutRefundNo:   m.OrderNum,     // 商户退款单号
		RefundId:      "",             // 操作员
	}

	p := &RefundQueryResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform("refundQuery", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
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
func (sp *WeixinScanPay) ProcessCancel(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	d := &ReverseReq{
		CommonParams: *getCommonParams(m),

		TransactionId: "",             // 微信订单号
		OutTradeNo:    m.OrigOrderNum, // 商户订单号
	}

	p := &ReverseResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform("reverse", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
	ret = &model.ScanPayResponse{
		Respcd:          status,       // 交易结果  M
		ChannelOrderNum: "",           // 渠道交易号 C
		ConsumerAccount: "",           // 渠道账号  C
		ConsumerId:      "",           // 渠道账号ID   C
		ErrorDetail:     msg,          // 错误信息   C
		QrCode:          m.ScanCodeId, // 二维码 C
		ChanRespCode:    p.ErrCode,    // 渠道详细应答码
	}

	return ret, err
}

// ProcessClose 关闭接口
func (sp *WeixinScanPay) ProcessClose(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	d := &CloseReq{
		CommonParams: *getCommonParams(m),

		TransactionId: "",             // 微信订单号
		OutTradeNo:    m.OrigOrderNum, // 商户订单号
	}

	p := &CloseResp{}
	if err = base(d, p); err != nil {
		return nil, err
	}

	status, msg := transform("close", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
	ret = &model.ScanPayResponse{
		Respcd:       status,    // 交易结果  M
		ErrorDetail:  msg,       // 错误信息   C
		ChanRespCode: p.ErrCode, // 渠道详细应答码
	}

	return ret, err
}
