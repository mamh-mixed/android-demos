package scanpay

import (
	"bufio"
	"fmt"
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	"strconv"
	"strings"
	"time"
)

// WeixinScanPay 微信扫码支付
type WeixinScanPay struct{}

// DefaultWeixinScanPay 微信扫码支付默认实现
var DefaultWeixinScanPay WeixinScanPay

func getCommonParams(m *model.ScanPayRequest) *weixin.CommonParams {
	return &weixin.CommonParams{
		Appid:        m.AppID,        // 公众账号ID
		MchID:        m.ChanMerId,    // 商户号
		SubAppid:     m.SubAppID,     // 子公众账号ID
		SubMchId:     m.SubMchId,     // 子商户号
		NonceStr:     util.Nonce(32), // 随机字符串
		Sign:         "",             // 签名
		WeixinMD5Key: m.SignKey,      // md5key
		ClientCert:   m.PemCert,
		ClientKey:    m.PemKey,
		Req:          m,
	}
}

// ProcessBarcodePay 扫条码下单
func (sp *WeixinScanPay) ProcessBarcodePay(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {

	// 设置失效时间
	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour)

	d := &PayReq{
		CommonParams: *getCommonParams(m),

		Body:           parseGoods(m), // 商品描述
		OutTradeNo:     m.OrderNum,    // 商户订单号
		AuthCode:       m.ScanCodeId,  // 授权码
		TotalFee:       m.ActTxamt,    // 总金额
		SpbillCreateIP: util.LocalIP,  // 终端IP

		TimeStart:  startTime.Format("20060102150405"), // 交易起始时间
		TimeExpire: endTime.Format("20060102150405"),   // 交易结束时间

		// 非必填
		DeviceInfo: m.DeviceInfo, // 设备号
		GoodsGag:   m.GoodsTag,   // 商品标记
		// Detail:     m.WxpMarshalGoods(), // 商品详情
		// Attach:     m.Attach,         // 附加数据
		// FeeType:    m.CurrType,       // 货币类型
	}

	p := &PayResp{}
	if err = weixin.Execute(d, p); err != nil {
		return nil, err
	}

	status, msg, ec := weixin.Transform("pay", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	ret = &model.ScanPayResponse{
		Respcd:          status,          // 交易结果  M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: p.SubOpenid,     // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
		ErrorCode:       ec,
		PayTime:         p.TimeEnd,
	}
	// 如果非大商户模式，用自己的 openid
	if d.SubMchId == "" || p.SubOpenid == "" {
		ret.ConsumerAccount = p.Openid
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

		TransactionId: "",             // 微信支付订单号
		OutTradeNo:    m.OrigOrderNum, // 商户订单号
	}

	p := &PayQueryResp{}
	if err = weixin.Execute(d, p); err != nil {
		return nil, err
	}

	status, msg, ec := weixin.Transform("payQuery", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	// 如果返回的是成功的，要对 trade_state 做判断
	// SUCCESS—支付成功;REFUND—转入退款;NOTPAY—未支付;CLOSED—已关闭;REVOKED—已撤销;
	// USERPAYING-用户支付中;PAYERROR-支付失败(其他原因，如银行返回失败)
	if status == "00" {
		respCode := mongo.ScanPayRespCol.GetByWxp(p.TradeState, "payQuery")
		status, msg, ec = respCode.ISO8583Code, respCode.ISO8583Msg, respCode.ErrorCode
	}

	ret = &model.ScanPayResponse{
		Respcd:          status,          // 交易结果  M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: p.SubOpenid,     // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
		ErrorCode:       ec,
		PayTime:         p.TimeEnd,
	}
	// 如果非大商户模式，用自己的 openid
	if d.SubMchId == "" || p.SubOpenid == "" {
		ret.ConsumerAccount = p.Openid
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

	// 设置失效时间
	startTime, endTime := handleExpireTime(m.TimeExpire)

	// 判断tradeType
	tradeType := ""
	if m.OpenId == "" {
		tradeType = "NATIVE"
	} else {
		tradeType = "JSAPI"
	}

	d := &PrePayReq{
		CommonParams: *getCommonParams(m),

		DeviceInfo:     m.DeviceInfo,                  // 设备号
		Body:           parseGoods(m),                 // 商品描述
		Attach:         m.SysOrderNum + "," + m.ReqId, // 格式：系统订单号,日志Id
		OutTradeNo:     m.OrderNum,                    // 商户订单号
		TotalFee:       m.ActTxamt,                    // 总金额
		SpbillCreateIP: util.LocalIP,                  // 终端IP
		TimeStart:      startTime,                     // 交易起始时间
		TimeExpire:     endTime,                       // 交易结束时间
		NotifyURL:      weixinNotifyURL,               // 通知地址
		TradeType:      tradeType,                     // 交易类型
		ProductID:      "",                            // 商品ID
		Openid:         m.OpenId,                      // 用户标识
		GoodsGag:       m.GoodsTag,                    // 商品标记
		// FeeType:        m.CurrType,                         // 货币类型
		// Detail:         m.WxpMarshalGoods(),                // 商品详情
	}

	p := &PrePayResp{}
	if err = weixin.Execute(d, p); err != nil {
		return nil, err
	}

	status, msg, ec := weixin.Transform("prePay", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)

	ret = &model.ScanPayResponse{
		Respcd:       status,     // 交易结果  M
		ErrorDetail:  msg,        // 错误信息   C
		QrCode:       p.CodeURL,  // 二维码 C
		ChanRespCode: p.ErrCode,  // 渠道详细应答码
		PrePayId:     p.PrepayID, // 预支付标识
		ErrorCode:    ec,
	}

	return ret, err
}

// ProcessRefund 退款
func (sp *WeixinScanPay) ProcessRefund(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	// log.Debugf("%#c", m)
	d := &RefundReq{
		CommonParams: *getCommonParams(m),

		DeviceInfo:    m.DeviceInfo,   // 设备号
		TransactionId: "",             // 微信订单号
		OutTradeNo:    m.OrigOrderNum, // 商户订单号
		OutRefundNo:   m.OrderNum,     // 商户退款单号
		TotalFee:      m.TotalTxamt,   // 总金额
		RefundFee:     m.ActTxamt,     // 退款金额
		OpUserId:      m.ChanMerId,    // 操作员
		// RefundFeeType: m.CurrType,     // 货币种类
	}

	p := &RefundResp{}
	if err = weixin.Execute(d, p); err != nil {
		return nil, err
	}

	status, msg, ec := weixin.Transform("refund", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
	ret = &model.ScanPayResponse{
		Respcd:          status,          // 交易结果  M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ConsumerAccount: "",              // 渠道账号  C
		ConsumerId:      "",              // 渠道账号ID   C
		ErrorDetail:     msg,             // 错误信息   C
		QrCode:          m.ScanCodeId,    // 二维码 C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
		ErrorCode:       ec,
	}

	return ret, err
}

// ProcessRefundQuery 退款查询
func (sp *WeixinScanPay) ProcessRefundQuery(m *model.ScanPayRequest) (ret *model.ScanPayResponse, err error) {
	d := &RefundQueryReq{
		CommonParams: *getCommonParams(m),

		TransactionId: "",             // 微信订单号
		OutTradeNo:    m.OrigOrderNum, // 商户订单号
		OutRefundNo:   m.OrderNum,     // 商户退款单号
		RefundId:      "",             // 操作员
	}

	p := &RefundQueryResp{}
	if err = weixin.Execute(d, p); err != nil {
		return nil, err
	}

	status, msg, ec := weixin.Transform("refundQuery", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
	ret = &model.ScanPayResponse{
		Txndir:          "A",             // 交易方向 M M
		Respcd:          status,          // 交易结果  M
		ChannelOrderNum: p.TransactionId, // 渠道交易号 C
		ErrorDetail:     msg,             // 错误信息   C
		ChanRespCode:    p.ErrCode,       // 渠道详细应答码
		ErrorCode:       ec,
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
	if err = weixin.Execute(d, p); err != nil {
		return nil, err
	}

	status, msg, ec := weixin.Transform("reverse", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
	ret = &model.ScanPayResponse{
		Respcd:          status,       // 交易结果  M
		ChannelOrderNum: "",           // 渠道交易号 C
		ConsumerAccount: "",           // 渠道账号  C
		ConsumerId:      "",           // 渠道账号ID   C
		ErrorDetail:     msg,          // 错误信息   C
		QrCode:          m.ScanCodeId, // 二维码 C
		ChanRespCode:    p.ErrCode,    // 渠道详细应答码
		ErrorCode:       ec,
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
	if err = weixin.Execute(d, p); err != nil {
		return nil, err
	}

	status, msg, ec := weixin.Transform("close", p.ReturnCode, p.ResultCode, p.ErrCode, p.ErrCodeDes)
	ret = &model.ScanPayResponse{
		Respcd:       status,    // 交易结果  M
		ErrorDetail:  msg,       // 错误信息   C
		ChanRespCode: p.ErrCode, // 渠道详细应答码
		ErrorCode:    ec,
	}

	return ret, err
}

//微信对账接口
func (sp *WeixinScanPay) ProcessSettleEnquiry(m *model.ScanPayRequest, cbd model.ChanBlendMap) error {

	if cbd == nil {
		return fmt.Errorf("%s", "nil map found")
	}

	d := &SettleQueryReq{
		CommonParams: *getCommonParams(m),
		SettleDate:   strings.Replace(m.SettDate, "-", "", -1),
		SettleType:   "ALL",
	}

	d.CommonParams.Req = m

	p := &SettleQueryResp{}
	dataStr, err := weixin.SettleExecute(d, p)
	if err != nil {
		return err
	}

	analysisSettleData(dataStr, cbd)

	return nil
}

//分析数据 如：交易时间,公众账号ID,商户号,子商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,总金额,代金券或立减优惠金额,微信退款单号,商户退款单号,退款金额, 代金券或立减优惠退款金额，退款类型，退款状态,商品名称,商户数据包,手续费,费率
func analysisSettleData(dataStr string, cbd model.ChanBlendMap) { //外部map key为商户号，内部map key为订单号

	if dataStr == "" {
		return
	}

	dataStr = strings.Replace(dataStr, "`", "", -1)
	log.Debugf("weixin settle csv: %s", dataStr)
	//modelMMap := make(map[string]map[string][]model.BlendElement)
	strStream := strings.NewReader(dataStr)
	rd := bufio.NewReader(strStream)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil {
			break
		}

		if strings.Contains(line, "交易时间") {
			continue
		}

		dataArray := strings.Split(line, ",")
		if len(dataArray) != 24 {
			break
		}

		var elementModel model.BlendElement
		elementModel.Chcd = "WXP"
		elementModel.ChcdName = "微信"
		elementModel.ChanMerID = dataArray[3] // 子商户号
		elementModel.OrderTime = dataArray[0] // 时间
		elementModel.OrderID = dataArray[5]   // 微信交易号
		elementModel.LocalID = dataArray[6]   // 商户订单号
		elementModel.IsBlend = false
		// init
		recsMap, ok := cbd[elementModel.ChanMerID]
		if !ok {
			recsMap = make(map[string][]model.BlendElement)
		}

		switch dataArray[9] {
		case "SUCCESS":
			elementModel.OrderAct = dataArray[12] //金额
		case "REFUND", "REVOKED":
			elementModel.OrderAct = "-" + dataArray[16]
			elementModel.RefundOrderID = dataArray[15]
		}

		elementModel.OrderType = dataArray[9]
		elementArray, ok := recsMap[elementModel.OrderID]
		if !ok {
			elementArray = make([]model.BlendElement, 0)
		}
		elementArray = append(elementArray, elementModel)
		recsMap[elementModel.OrderID] = elementArray

		// back
		cbd[elementModel.ChanMerID] = recsMap
	}
}

func handleExpireTime(expirtTime string) (string, string) {

	layout := "20060102150405"
	startTime := time.Now()
	defaultEntTime := startTime.Add(24 * time.Hour)

	var stStr, etStr = startTime.Format(layout), defaultEntTime.Format(layout)

	if expirtTime == "" {
		return stStr, etStr
	}

	et, err := time.ParseInLocation(layout, expirtTime, time.Local)
	if err != nil {
		return stStr, etStr
	}

	d := et.Sub(startTime)
	if d < 5*time.Minute {
		return stStr, startTime.Add(5 * time.Minute).Format(layout)
	}

	return stStr, expirtTime
}

// parseGoods 按照微信格式输出商品详细
func parseGoods(req *model.ScanPayRequest) string {

	goods, err := req.MarshalGoods()
	if err != nil {
		// 格式不对，送配置的商品名称，防止商户送的内容过长
		return req.Subject
	}

	var goodsName []string
	if len(goods) > 0 {
		for _, v := range goods {
			goodsName = append(goodsName, v.GoodsName)
		}

		body := strings.Join(goodsName, ",")
		bodySizes := []rune(body)
		if len(bodySizes) > 20 {
			body = string(bodySizes[:20]) + "..."
		}
		return body
	}

	// 假如商品详细为空，送配置的商品名称
	return req.Subject
}
