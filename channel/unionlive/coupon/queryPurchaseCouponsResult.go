package coupon

import "github.com/CardInfoLink/quickpay/model"

type QueryPurchaseCouponsResultReqHeader struct {
	Version       string `json:"version"`       // 报文版本号	15	M	当前版本1.0
	Transdirect   string `json:"transDirect"`   // 交易方向	1	M	Q：请求
	Transtype     string `json:"transType"`     // 交易类型	8	M	本交易固定值W394
	Merchantid    string `json:"merchantId"`    // 商户编号	15	M	由优麦圈后台分配给商户的编号
	Submittime    string `json:"submitTime"`    // 交易提交时间	14	M	固定格式：yyyyMMddHHmmss
	Clienttraceno string `json:"clientTraceNo"` // 客户端交易流水号	40	M	本交易客户端的唯一交易流水号
}
type QueryPurchaseCouponsResultReqBody struct {
	Couponsno        string `json:"couponsNo"`        // 优麦圈电子券号	50	M	优麦圈电子券号
	Termid           string `json:"termId"`           // 终端编号	8	M	由优麦圈后台分配给该终端的编号
	Termsn           string `json:"termSn"`           // 终端唯一序列号	100	M	商户终端对应的硬件唯一序列号
	Amount           int    `json:"amount"`           // 要验证的次数	10	M	要验证该券码的次数，次数必须大于0
	Oldclienttraceno string `json:"oldClientTraceNo"` // 原验证交易客户端交易流水号	40	M	原验证交易客户端的唯一交易流水号
	Oldsubmittime    string `json:"oldSubmitTime"`    // 原交易提交时间	14	M	固定格式：yyyyMMddHHmmss
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
	Transdirect   string `json:"transDirect"`   // 交易方向	1	M	A：应答
	Transtype     string `json:"transType"`     // 交易类型	8	M	原样返回。本交易固定值W394
	Merchantid    string `json:"merchantId"`    // 商城编号	15	M	由优麦圈后台分配给商户的编号
	Submittime    string `json:"submitTime"`    // 提交时间	14	M	原样返回。固定格式：yyyyMMddHHmmss
	Clienttraceno string `json:"clientTraceNo"` // 客户端交易流水号	40	M	原样返回。客户端本会话中的唯一交易流水号
	Hosttime      string `json:"hostTime"`      // 后台处理时间	14	M	固定格式：yyyyMMddHHmmss
	Hosttraceno   string `json:"hostTraceNo"`   // 后台交易流水号	10	M	后台交易唯一流水号
	Returncode    string `json:"returnCode"`    // 返回码	4	M	后台交易返回码。详见附录。
	Returnmessage string `json:"returnMessage"` // 返回码描述	100	M	后台返回码描述
}
type QueryPurchaseCouponsResultRespBody struct {
	Couponsno        string `json:"couponsNo"`        // 优麦圈电子券号	50	M	优麦圈电子券号，中间部分以*屏蔽
	Termid           string `json:"termId"`           // 终端编号	8	M	由优麦圈后台分配给该终端的编号
	Termsn           string `json:"termSn"`           // 终端唯一序列号	100	M	商户终端对应的硬件唯一序列号
	Amount           int    `json:"amount"`           // 要验证的次数	10	M	要验证该券码的次数
	Oldreturncode    string `json:"oldReturnCode"`    // 原交易返回码	4	C1	原验证交易的后台交易返回码。详见附录。
	Oldreturnmessage string `json:"oldReturnMessage"` // 原交易返回码描述	100	C1	原验证交易的后台返回码描述
	Authcode         int    `json:"authCode"`         // 主机授权码	10	C2	后台交易处理成功后的授权码
	Prodname         string `json:"prodName"`         // 券产品名称	32	C2	该券的产品名称
	Proddesc         string `json:"prodDesc"`         // 券描述	100	C2	该券的产品描述
	Availcount       int    `json:"availCount"`       // 券剩余可用次数	10	C2	该券的剩余可用次数
	Expdate          string `json:"expDate"`          // 券有效期	10	C2	券的最后可用日期，格式：yyyyMMdd
}
type QueryPurchaseCouponsResultResp struct {
	Header QueryPurchaseCouponsResultRespHeader `json:"header"`
	Body   QueryPurchaseCouponsResultRespBody   `json:"body"`
}
