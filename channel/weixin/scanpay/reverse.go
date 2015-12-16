package scanpay

import (
	"net/http"

	"github.com/CardInfoLink/quickpay/channel/weixin"
)

// 参考文档 https://pay.weixin.qq.com/wiki/doc/api/micropay_sl.php?chapter=9_11&index=3
// 应用场景
// 支付交易返回失败或支付系统超时，调用该接口撤销交易。如果此订单用户支付失败，微信支付系统会将此订单关闭；如果用户支付成功，微信支付系统会将此订单资金退还给用户。
// 注意：7天以内的交易单可调用撤销，其他正常支付的单如需实现相同功能请调用申请退款API。提交支付交易后调用【查询订单API】，没有明确的支付结果再调用【撤销订单API】。
// 接口链接
// https://api.mch.weixin.qq.com/secapi/pay/reverse
// 是否需要证书
// 请求需要双向证书。

// ReverseReq 撤销订单
type ReverseReq struct {
	weixin.CommonParams

	TransactionId string `xml:"transaction_id,omitempty" url:"transaction_id,omitempty"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no" url:"out_trade_no" validate:"nonzero"`      // 商户订单号
	FeeType       string `xml:"fee_type,omitempty" url:"fee_type,omitempty"`             // 标价币种
}

// GetURI 取接口地址
func (r *ReverseReq) GetURI() string {
	return "/secapi/pay/reverse"
}

// GetHTTPClient 使用双向 HTTPS 认证
func (r *ReverseReq) GetHTTPClient() *http.Client {
	return r.GetHTTPSClient()
}

// ReverseResp 撤销订单
type ReverseResp struct {
	weixin.CommonBody

	Recall string `xml:"recall" url:"recall,omitempty"` // 是否重调
}
