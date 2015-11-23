package coupon

import "github.com/CardInfoLink/quickpay/model"

type QueryPurchaseLogReqHeader struct {
	Version       string `json:"version"`       // 报文版本号	15	M	当前版本1.0
	TransDirect   string `json:"transDirect"`   // 交易方向	1	M	Q：请求
	TransType     string `json:"transType"`     // 交易类型	8	M	本交易固定值W395
	MerchantId    string `json:"merchantId"`    // 商户编号	15	M	由优麦圈后台分配给商户的编号
	SubmitTime    string `json:"submitTime"`    // 交易提交时间	14	M	固定格式：yyyyMMddHHmmss
	ClientTraceNo string `json:"clientTraceNo"` // 客户端交易流水号	40	M	本交易客户端的唯一交易流水号

}
type QueryPurchaseLogReqBody struct {
	TermId    string `json:"termId"`    // 终端编号	8	M	由优麦圈后台分配给该终端的编号
	TermSn    string `json:"termSn"`    // 终端唯一序列号	100	M	商户终端对应的硬件唯一序列号
	ExtMercId string `json:"extMercId"` // 商户自定义商户号 40 C3 商户自定义的商户编号,可选,如果送入则会校验该值
	ExtTermId string `json:"extTermId"` // 商户自定义终端号  40 C3 商户自定义的终端编号,可选
	PageIndex int    `json:"pageIndex"` // 分页索引	5	O	指定查询第几页数据，从0开始。每页最多返回20笔交易，按时间倒序返回
}

// QueryPurchaseLogReq W395-商户券验证流水查询
type QueryPurchaseLogReq struct {
	Header QueryPurchaseLogReqHeader `json:"header"`
	Body   QueryPurchaseLogReqBody   `json:"body"`
	SpReq  *model.ScanPayRequest     `json:"-" url:"-" bson:"-"`
}

func (req *QueryPurchaseLogReq) GetT() string {
	return "QueryPurchaseLog"
}
func (req *QueryPurchaseLogReq) GetSpReq() *model.ScanPayRequest {
	return req.SpReq
}

type QueryPurchaseLogRespHeader struct {
	Version       string `json:"version"`       // 报文版本号	15	M	当前版本1.0
	TransDirect   string `json:"transDirect"`   // 交易方向	1	M	A：应答
	TransType     string `json:"transType"`     // 交易类型	8	M	原样返回。本交易固定值W395
	MerchantId    string `json:"merchantId"`    // 商城编号	15	M	由优麦圈后台分配给商户的编号
	SubmitTime    string `json:"submitTime"`    // 提交时间	14	M	原样返回。固定格式：yyyyMMddHHmmss
	ClientTraceNo string `json:"clientTraceNo"` // 客户端交易流水号	40	M	原样返回。客户端本会话中的唯一交易流水号
	HostTime      string `json:"hostTime"`      // 后台处理时间	14	M	固定格式：yyyyMMddHHmmss
	HostTraceNo   string `json:"hostTraceNo"`   // 后台交易流水号	10	M	后台交易唯一流水号
	ReturnCode    string `json:"returnCode"`    // 返回码	4	M	后台交易返回码。详见附录。
	ReturnMessage string `json:"returnMessage"` // 返回码描述	100	M	后台返回码描述
}
type QueryPurchaseLogRespBody struct {
	TermId    string `json:"termId"`    // 终端编号	8	M	由优麦圈后台分配给该终端的编号
	TermSn    string `json:"termSn"`    // 终端唯一序列号	100	M	商户终端对应的硬件唯一序列号
	ExtMercId string `json:"extMercId"` // 商户自定义商户号 40 C3  原样返回。商户自定义的商户编号
	ExtTermId string `json:"extTermId"` // 商户自定义终端号 40 C3  原样返回。商户自定义的终端编号
	PageNum   int    `json:"pageNum"`   // 分页数量	6	O	正整数，表示总共分了多少页
	PageIndex int    `json:"pageIndex"` // 分页索引	6	O	正整数，表示这是第几页的数据，从0开始
	Count     int    `json:"count"`     // 总验券交易数	6	M	当日成功验券交易总数

	Trans []struct {
		CouponsNo        string `json:"couponsNo"`        // 优麦圈电子券号	50	M	优麦圈电子券号，中间部分以*屏蔽
		Amount           int    `json:"amount"`           // 验证的次数	10	M	验证该券码的次数
		AuthCode         string `json:"authCode"`         // 主机授权码	10	M	后台交易处理成功后的授权码
		ProdName         string `json:"prodName"`         // 券产品名称	32	M	该券的产品名称
		OldReturnCode    string `json:"oldReturnCode"`    // 原交易返回码	4	M	后台交易返回码。详见附录。
		OldReturnMessage string `json:"oldReturnMessage"` // 原交易返回码描述	100	M	后台返回码描述
	} `json:"trans"`
}
type QueryPurchaseLogResp struct {
	Header QueryPurchaseLogRespHeader `json:"header"`
	Body   QueryPurchaseLogRespBody   `json:"body"`
}
