package scanpay

import (
	"encoding/json"
	"net/url"
)

// 参考文档 https://app.alipay.com/market/document.htm?name=saomazhifu#page-14
// 收银员通过本接口将订单信息上送至支付宝后，将支付宝返回的二维码信息展示给用户，由用户扫描二维码完成订单支付。

// PrecreateReq 预下单
type PrecreateReq struct {
	CommonParams

	OutTradeNo           string `json:"out_trade_no" validate:"nonzero"` // 商户订单号，64个字符以内、只能包含字母、数字、下划线;需保证在商户端不重复
	SellerID             string `json:"seller_id,omitempty"`             // 卖家支付宝用户 ID，如果该值为空，则默认为商户签约账号对应的支付宝用户 ID
	TotalAmount          string `json:"total_amount" validate:"nonzero"` // 订单总金额
	DiscountableAmount   string `json:"discountable_amount,omitempty"`   // 可打折金额
	UndiscountableAmount string `json:"undiscountable_amount,omitempty"` // 不可打折金额
	Subject              string `json:"subject" validate:"nonzero"`      // 订单标题
	Body                 string `json:"body,omitempty"`                  // 订单描述
	GoodsDetail          string `json:"goods_detail,omitempty"`          // 商品明细列表，Json 格式，其它说明详见:“商品明细说明” [{"goods_id": "apple-01","goods_name":"ipad","goods_category":"7788230","price":" 2000.00","quantity":"1"}]
	OperatorID           string `json:"operator_id,omitempty"`           // 商户操作员编号
	StoreID              string `json:"store_id,omitempty"`              // 商户门店编号
	TerminalID           string `json:"terminal_id,omitempty"`           // 机具终端编号
	ExtendParams         string `json:"extend_params,omitempty"`         // 扩展参数
	TimeExpire           string `json:"time_expire,omitempty"`           // 支付超时时间，该笔订单允许的最晚付款时间，逾期将关闭交易。格式为: yyyy-MM-dd HH:mm:ss	2015-01-01 11:01:01

}

// Values 组装公共参数
func (c *PrecreateReq) Values() (v url.Values) {
	c.CommonParams.Method = "alipay.trade.precreate"
	return c.CommonParams.Values()
}

// PrecreateResp 预下单
type PrecreateResp struct {
	CommonBody

	Raw json.RawMessage `json:"alipay_trade_precreate_response"` // 返回消息体

	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	QrCode     string `json:"qr_code"`      // 二维码码串
}

// GetRaw 报文内容
func (c *PrecreateResp) GetRaw() []byte {
	return []byte(c.Raw)
}
