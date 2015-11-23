package coupon

import "github.com/CardInfoLink/quickpay/model"

type QueryPurchaseCouponsResultReqHeader struct {
	Version       string `json:"version"`       // 报文版本号	15	M	当前版本1.0
	TransDirect   string `json:"transDirect"`   // 交易方向	1	M	Q：请求
	TransType     string `json:"transType"`     // 交易类型	8	M	本交易固定值W394
	MerchantId    string `json:"merchantId"`    // 商户编号	15	M	由优麦圈后台分配给商户的编号
	SubmitTime    string `json:"submitTime"`    // 交易提交时间	14	M	固定格式：yyyyMMddHHmmss
	ClientTraceNo string `json:"clientTraceNo"` // 客户端交易流水号	40	M	本交易客户端的唯一交易流水号
}
type QueryPurchaseCouponsResultReqBody struct {
	CouponsNo        string `json:"couponsNo"`        // 优麦圈电子券号	50	M	优麦圈电子券号
	TermId           string `json:"termId"`           // 终端编号	8	M	由优麦圈后台分配给该终端的编号
	TermSn           string `json:"termSn"`           // 终端唯一序列号	100	M	商户终端对应的硬件唯一序列号
	ExtMercId        string `json:"extMercId"`        // 商户自定义商户号 40 C3 商户自定义的商户编号,可选,如果送入则会校验该值
	ExtTermId        string `json:"extTermId"`        // 商户自定义终端号  40 C3 商户自定义的终端编号,可选
	Amount           int    `json:"amount"`           // 要验证的次数	10	M	要验证该券码的次数，次数必须大于0
	OldClientTraceNo string `json:"oldClientTraceNo"` // 原验证交易客户端交易流水号	40	M	原验证交易客户端的唯一交易流水号
	OldSubmitTime    string `json:"oldSubmitTime"`    // 原交易提交时间	14	M	固定格式：yyyyMMddHHmmss
}

// QueryPurchaseCouponsResultReq W394-电子券验证结果查询
type QueryPurchaseCouponsResultReq struct {
	Header QueryPurchaseCouponsResultReqHeader `json:"header"`
	Body   QueryPurchaseCouponsResultReqBody   `json:"body"`
	SpReq  *model.ScanPayRequest               `json:"-" url:"-" bson:"-"`
}

func (req *QueryPurchaseCouponsResultReq) GetT() string {
	return "QueryPurchaseCouponsResult"
}
func (req *QueryPurchaseCouponsResultReq) GetSpReq() *model.ScanPayRequest {
	return req.SpReq
}

type QueryPurchaseCouponsResultRespHeader struct {
	Version       string `json:"version"`       // 报文版本号	15	M	当前版本1.0
	TransDirect   string `json:"transDirect"`   // 交易方向	1	M	A：应答
	TransType     string `json:"transType"`     // 交易类型	8	M	原样返回。本交易固定值W394
	MerchantId    string `json:"merchantId"`    // 商城编号	15	M	由优麦圈后台分配给商户的编号
	SubmitTime    string `json:"submitTime"`    // 提交时间	14	M	原样返回。固定格式：yyyyMMddHHmmss
	ClientTraceNo string `json:"clientTraceNo"` // 客户端交易流水号	40	M	原样返回。客户端本会话中的唯一交易流水号
	HostTime      string `json:"hostTime"`      // 后台处理时间	14	M	固定格式：yyyyMMddHHmmss
	HostTraceNo   string `json:"hostTraceNo"`   // 后台交易流水号	10	M	后台交易唯一流水号
	ReturnCode    string `json:"returnCode"`    // 返回码	4	M	后台交易返回码。详见附录。
	ReturnMessage string `json:"returnMessage"` // 返回码描述	100	M	后台返回码描述
}
type QueryPurchaseCouponsResultRespBody struct {
	CouponsNo        string `json:"couponsNo"`        // 优麦圈电子券号	50	M	优麦圈电子券号，中间部分以*屏蔽
	TermId           string `json:"termId"`           // 终端编号	8	C	由优麦圈后台分配给该终端的编号
	TermSn           string `json:"termSn"`           // 终端唯一序列号	100	C	商户终端对应的硬件唯一序列号
	ExtMercId        string `json:"extMercId"`        // 商户自定义商户号 40 C  原样返回。商户自定义的商户编号
	ExtTermId        string `json:"extTermId"`        // 商户自定义终端号 40 C  原样返回。商户自定义的终端编号
	Amount           int    `json:"amount"`           // 要验证的次数	10	M	要验证该券码的次数
	OldReturnCode    string `json:"oldReturnCode"`    // 原交易返回码	4	C1	原验证交易的后台交易返回码。详见附录。
	OldReturnMessage string `json:"oldReturnMessage"` // 原交易返回码描述	100	C1	原验证交易的后台返回码描述
	AuthCode         int    `json:"authCode"`         // 主机授权码	10	C2	后台交易处理成功后的授权码
	ProdName         string `json:"prodName"`         // 券产品名称	32	C2	该券的产品名称
	ProdDesc         string `json:"prodDesc"`         // 券描述	100	C2	该券的产品描述
	AvailCount       int    `json:"availCount"`       // 券剩余可用次数	10	C2	该券的剩余可用次数
	VoucherType      int    `json:"voucherType"`      // 券类型 2 M  券的类型。21:刷卡活动满减券;22: 刷卡活动固定金额券;23:刷卡活动 满折券;31:礼包券;其他:普通电子券
	SaleMinAmount    int    `json:"saleMinAmount"`    // 满足优惠条件的最小金额 12 M 满足优惠条件的最小金额,满折、满 减等优惠中需要满足的金额。单位: 分。
	SaleDiscount     int    `json:"saleDiscount"`     // 抵扣值 12 M 如果是满减券,则是减免的金额,以 分为单位;如果是满折券,则是折扣 率,如 9.5 折则值为 95;如果是固定 金额刷卡券,则是固定要扣款的金额, 以分为单位,如值为 100 则表示固定 付款 1 元钱
	Expdate          string `json:"expDate"`          // 券有效期	10	C2	券的最后可用日期，格式：yyyyMMdd
}
type QueryPurchaseCouponsResultResp struct {
	Header QueryPurchaseCouponsResultRespHeader `json:"header"`
	Body   QueryPurchaseCouponsResultRespBody   `json:"body"`
}
