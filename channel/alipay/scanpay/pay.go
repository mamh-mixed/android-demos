package scanpay

import (
	"encoding/json"
	"net/url"
)

//参考文档 https://app.alipay.com/market/document.htm?name=tiaomazhifu#page-13
//收银员使用扫码设备读取用户支付宝钱包“付款码”后，将二维码或条码信息通过本接口上送至支付宝发起支付。

//PayReq 条码支付请求
type PayReq struct {
	CommonParams

	OutTradeNo           string `json:"out_trade_no" validate:"nonzero"` //商户订单号
	Scene                string `json:"scene" validate:"nonzero"`        //支付场景
	AuthCode             string `json:"auth_code" validate:"nonzero"`    //支付授权码
	SellerID             string `json:"seller_id,omitempty"`             //卖家支付宝用户 ID
	TotalAmount          string `json:"total_amount" validate:"nonzero"` //订单总金额
	DiscountableAmount   string `json:"discountable_amount,omitempty"`   //可打折金额
	UndiscountableAmount string `json:"undiscountable_amount,omitempty"` //不可打折金额
	Subject              string `json:"subject" validate:"nonzero"`      //订单标题
	Body                 string `json:"body,omitempty"`                  //订单描述
	GoodsDetail          string `json:"goods_detail,omitempty"`          //商品明细列表信息,json格式
	OperatorID           string `json:"operator_id,omitempty"`           //商户操作员编号
	StoreID              string `json:"store_id,omitempty"`              //商户门店编号
	TerminalID           string `json:"terminal_id,omitempty"`           //机具终端编号
	ExtendParams         string `json:"extend_params,omitempty"`         //业务扩展参数
	TimeExpire           string `json:"time_expire,omitempty"`           //该笔订单允许的最晚付款时间，逾期将关闭交易。格式为: yyyy-MM-dd HH:mm:ss
}

// Values 组装公共参数
func (c *PayReq) Values() (v url.Values) {
	c.CommonParams.Method = "alipay.trade.pay"
	return c.CommonParams.Values()
}

//PayResp 条码支付请求
type PayResp struct {
	CommonBody

	Raw json.RawMessage `json:"alipay_trade_pay_response"` // 返回消息体

	TradeNo        string `json:"trade_no"`         //支付宝交易号
	OutTradeNo     string `json:"out_trade_no"`     //商户订单号
	OpenID         string `json:"open_id"`          //买家支付宝用户号
	BuyerLogonID   string `json:"buyer_logon_id"`   //买家支付宝账号
	TotalAmount    string `json:"total_amount"`     //交易金额
	ReceiptAmount  string `json:"receipt_amount"`   //实收金额
	InvoiceAmount  string `json:"invoice_amount"`   //开票金额
	BuyerPayAmount string `json:"buyer_pay_amount"` //付款金额
	PointAmount    string `json:"point_amount"`     //积分宝金额
	GmtPayment     string `json:"gmt_payment"`      //买家付款时间。 格式为 yyyy-MM-dd HH:mm:ss
	FundBillList   []struct {
		FundChannel string `json:"fund_channel,omitempty"` //支付渠道
		Amount      string `json:"amount,omitempty"`       //支付金额
	} `json:"fund_bill_list"` //交易资金明细信息集合

}

// GetRaw 报文内容
func (c *PayResp) GetRaw() []byte {
	return []byte(c.Raw)
}
