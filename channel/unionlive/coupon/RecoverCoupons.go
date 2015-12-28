package coupon

import "github.com/CardInfoLink/quickpay/model"

type RecoverCouponsReqHeader struct {
	Version       string `json:"version"`       // 报文版本号  15 M 当前版本 1.0
	TransDirect   string `json:"transDirect"`   // 交易方向  1 M Q:请求
	TransType     string `json:"transType"`     // 交易类型 8 M 本交易固定值W493
	MerchantId    string `json:"merchantId"`    // 商户编号 15 M 由优麦圈后台分配给商户的编号
	SubmitTime    string `json:"submitTime"`    // 交易提交时间 14 M 固定格式:yyyyMMddHHmmss
	ClientTraceNo string `json:"clientTraceNo"` // 客户端交易流水号 40 M 客户端的唯一交易流水号
}
type RecoverCouponsReqBody struct {
	CouponsNo        string `json:"couponsNo"`             // 优麦圈电子券号 50 M 原验证交易的电子券号
	TermId           string `json:"termId"`                // 终端编号 8 C1 原验证交易值
	TermSn           string `json:"termSn"`                // 终端唯一序列号 100 C2 原验证交易值。商户终端对应的硬件唯一序列号
	ExtMercId        string `json:"extMercId"`             // 商户自定义商户号 40 O 原验证交易值。商户自定义的商户编号，可选，如果送入则会校验该值
	ExtTermId        string `json:"extTermId"`             // 商户自定义终端号  40 O 原验证交易值。商户自定义的终端编号，可选
	Amount           int    `json:"amount"`                // 要验证的次数  10 0 原验证交易值。要验证该券码的次数，次数必须大于0，默认为1
	Cardbin          string `json:"cardbin"`               // 银行卡cardbin或者用户标识等 30 O 原验证交易值。如果券配置了限制cardbin，则传入该值会做相应的限制判断
	TransAmount      int64  `json:"transAmount,omitempty"` // 交易原始金额 10 O 原验证交易值。交易的原始金额，即抵扣前的原价
	PayType          int    `json:"payType,omitempty"`     // 支付方式 2 O 原验证交易值。2：银行卡支付；4：微信支付；5：支付宝支付
	OldSubmitTime    string `json:"oldSubmitTime"`         // 原验证交易提交时间 14 M  原验证交易的提交时间
	OldClientTraceNo string `json:"oldClientTraceNo"`      // 原验证交易客户端流水号 40 M 原验证交易送入的客户端流水号
}

// RecoverCouponsReq W493-电子券验证冲正
type RecoverCouponsReq struct {
	Header RecoverCouponsReqHeader `json:"header"`
	Body   RecoverCouponsReqBody   `json:"body"`
	SpReq  *model.ScanPayRequest   `json:"-" url:"-" bson:"-"`
}

func (req *RecoverCouponsReq) GetT() string {
	return "RecoverCoupons"
}
func (req *RecoverCouponsReq) GetSpReq() *model.ScanPayRequest {
	return req.SpReq
}

type RecoverCouponsRespHeader struct {
	Version       string `json:"version"`       // 报文版本号	15	M	当前版本1.0
	TransDirect   string `json:"transDirect"`   // 交易方向	1	M	A：应答
	TransType     string `json:"transType"`     // 交易类型	8	M	原样返回。本交易固定值W492
	MerchantId    string `json:"merchantId"`    // 商户编号	15	M	由优麦圈后台分配给商户的编号
	SubmitTime    string `json:"submitTime"`    // 提交时间	14	M	原样返回。固定格式：yyyyMMddHHmmss
	ClientTraceNo string `json:"clientTraceNo"` // 客户端交易流水号	40	M	原样返回。客户端本会话中的唯一交易流水号
	HostTime      string `json:"hostTime"`      // 后台处理时间	14	M	固定格式：yyyyMMddHHmmss
	HostTraceNo   string `json:"hostTraceNo"`   // 后台交易流水号	10	M	后台交易唯一流水号
	ReturnCode    string `json:"returnCode"`    // 返回码	4	M	后台交易返回码。详见附录。
	ReturnMessage string `json:"returnMessage"` // 返回码描述	100	M	后台返回码描述
}
type RecoverCouponsRespBody struct {
	CouponsNo        string `json:"couponsNo"`        // 优麦圈电子券号 50 M 原验证交易的电子券号
	TermId           string `json:"termId"`           // 终端编号 8 C1 原验证交易值
	TermSn           string `json:"termSn"`           // 终端唯一序列号 100 C2 原验证交易值。商户终端对应的硬件唯一序列号
	ExtMercId        string `json:"extMercId"`        // 商户自定义商户号 40 O 原验证交易值。商户自定义的商户编号，可选，如果送入则会校验该值
	ExtTermId        string `json:"extTermId"`        // 商户自定义终端号  40 O 原验证交易值。商户自定义的终端编号，可选
	Amount           int    `json:"amount"`           // 要验证的次数  10 0 原验证交易值。要验证该券码的次数，次数必须大于0，默认为1
	Cardbin          string `json:"cardbin"`          // 银行卡cardbin或者用户标识等 30 O 原验证交易值。如果券配置了限制cardbin，则传入该值会做相应的限制判断
	TransAmount      int64  `json:"transAmount"`      // 交易原始金额 10 O 原验证交易值。交易的原始金额，即抵扣前的原价
	PayType          int    `json:"payType"`          // 支付方式 2 O 原验证交易值。2：银行卡支付；4：微信支付；5：支付宝支付
	OldSubmitTime    string `json:"oldSubmitTime"`    // 原验证交易提交时间 14 M  原验证交易的提交时间
	OldClientTraceNo string `json:"oldClientTraceNo"` // 原验证交易客户端流水号 40 M 原验证交易送入的客户端流水号
	AuthCode         int    `json:"authCode"`         // 冲正交易主机授权码	10	C1	后台冲正交易处理成功后的授权码
}
type RecoverCouponsResp struct {
	Header RecoverCouponsRespHeader `json:"header"`
	Body   RecoverCouponsRespBody   `json:"body"`
}
