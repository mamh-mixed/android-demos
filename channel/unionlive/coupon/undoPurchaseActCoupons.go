package coupon

import "github.com/CardInfoLink/quickpay/model"

type UndoPurchaseActCouponsReqHeader struct {
	Version       string `json:"version"`       // 报文版本号  15 M 当前版本 1.0
	TransDirect   string `json:"transDirect"`   // 交易方向  1 M Q:请求
	TransType     string `json:"transType"`     // 交易类型 8 M 本交易固定值W492
	MerchantId    string `json:"merchantId"`    // 商户编号 15 M 由优麦圈后台分配给商户的编号
	SubmitTime    string `json:"submitTime"`    // 交易提交时间 14 M 固定格式:yyyyMMddHHmmss
	ClientTraceNo string `json:"clientTraceNo"` // 客户端交易流水号 40 M 客户端的唯一交易流水号
}
type UndoPurchaseActCouponsReqBody struct {
	CouponsNo        string `json:"couponsNo"`        // 优麦圈电子券号 50 M 优麦圈电子券号
	TermId           string `json:"termId"`           // 终端编号 8 C1 由优麦圈后台分配给该终端的编号
	TermSn           string `json:"termSn"`           // 终端唯一序列号 100 C2 商户终端对应的硬件唯一序列号
	ExtMercId        string `json:"extMercId"`        // 商户自定义商户号 40 C3 商户自定义的商户编号,可选,如果送入则会校验该值
	ExtTermId        string `json:"extTermId"`        // 商户自定义终端号  40 C3 商户自定义的终端编号,可选
	OldTransAmount   int    `json:"oldTransAmount"`   // 原验证交易验证的次数 10 M 原验证交易验证该券码的次数 (amount字段),次数必须大于0,默 认为1
	OldSubmitTime    string `json:"oldSubmitTime"`    // 原验证交易提交时间 14 M  原验证交易的提交时间
	OldClientTraceNo string `json:"oldClientTraceNo"` // 原验证交易客户端流水号 40 M 原验证交易送入的客户端流水号
	OldHostTraceNo   string `json:"oldHostTraceNo"`   // 原交易后台交易流水号 10 M 原W452交易返回里的hostTraceNo值
}

// UndoPurchaseActCouponsReq W492-刷卡活动券验证撤销
type UndoPurchaseActCouponsReq struct {
	Header UndoPurchaseActCouponsReqHeader `json:"header"`
	Body   UndoPurchaseActCouponsReqBody   `json:"body"`
	SpReq  *model.ScanPayRequest           `json:"-" url:"-" bson:"-"`
}

func (req *UndoPurchaseActCouponsReq) GetT() string {
	return "UndoPurchaseActCoupons"
}
func (req *UndoPurchaseActCouponsReq) GetSpReq() *model.ScanPayRequest {
	return req.SpReq
}

type UndoPurchaseActCouponsRespHeader struct {
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
type UndoPurchaseActCouponsRespBody struct {
	CouponsNo        string `json:"couponsNo"`        // 优麦圈电子券号	50	M	优麦圈电子券号，中间部分以*屏蔽
	TermId           string `json:"termId"`           // 终端编号	8 C	由优麦圈后台分配给该终端的编号
	TermSn           string `json:"termSn"`           // 终端唯一序列号	100	C	商户终端对应的硬件唯一序列号
	ExtMercId        string `json:"extMercId"`        // 商户自定义商户号 40 C  原样返回。商户自定义的商户编号
	ExtTermId        string `json:"extTermId"`        // 商户自定义终端号 40 C  原样返回。商户自定义的终端编号
	AuthCode         int    `json:"authCode"`         // 主机授权码	10	C1	后台交易处理成功后的授权码
	OldTransAmount   int    `json:"oldTransAmount"`   // 原验证交易验证的次数 10 M 原验证交易验证该券码的次数 (amount字段),次数必须大于0,默 认为1
	OldSubmitTime    string `json:"oldSubmitTime"`    // 原验证交易提交时间 14 M  原验证交易的提交时间
	OldClientTraceNo string `json:"oldClientTraceNo"` // 原验证交易客户端流水号 40 M 原验证交易送入的客户端流水号
	OldHostTraceNo   string `json:"oldHostTraceNo"`   // 原交易后台交易流水号 10 M 原W452交易返回里的hostTraceNo值
}
type UndoPurchaseActCouponsResp struct {
	Header UndoPurchaseActCouponsRespHeader `json:"header"`
	Body   UndoPurchaseActCouponsRespBody   `json:"body"`
}
