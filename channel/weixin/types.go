package weixin

import (
	"encoding/xml"
	"fmt"
)

type WeixinRequest interface {
	setSign(sign string)
	//	display()
}
type WeixinResponse interface {
	show()
}

// func (w *MicropayRequest) display() {
// 	fmt.Println("MicropayRequest")
// }
// func (w *OrderqueryRequest) display() {
// 	fmt.Println("OrderqueryRequest")
// }

func (w *MicropayRequest) setSign(sign string) {
	w.Sign = sign
}
func (w *OrderqueryRequest) setSign(sign string) {
	w.Sign = sign
}

// micropayRequest 请求参数
type MicropayRequest struct {
	XMLName xml.Name `xml:"xml"`
	// 可选
	DeviceInfo string `xml:"device_info,omitempty"` //设备号
	Detail     string `xml:"detail,omitempty"`      //商品详情
	Attach     string `xml:"attach,omitempty"`      //附加数据
	FeeType    string `xml:"fee_type,omitempty"`    //货币类型
	GoodsTag   string `xml:"goods_tag,omitempty"`   //商品标记

	// 必选
	AppId          string `xml:"appid"`                      //公众账号ID
	MchId          string `xml:"mch_id,omitempty"`           //商户号
	NonceStr       string `xml:"nonce_str,omitempty"`        //随机字符串
	Body           string `xml:"body,omitempty"`             //商品描述
	OutTradeNo     string `xml:"out_trade_no,omitempty"`     //商户订单号
	TotalFee       int    `xml:"total_fee,omitempty"`        //总金额
	SpbillCreateIp string `xml:"spbill_create_ip,omitempty"` //终端IP
	AuthCode       string `xml:"auth_code,omitempty"`        //授权码
	Sign           string `xml:"sign"`                       //签名
	SubMchId       string `xml:"sub_mch_id"`

	//
	NotifyUrl string `xml:"-"`
}

func (m *MicroPayResponse) show() {
	fmt.Println("MicroPayResponse")
}
func (m *OrderqueryResponse) show() {
	fmt.Println("OrderqueryResponse")
}

type MicroPayResponse struct {
	//
	ReturnCode string `xml:"return_code,omitempty"`           //返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"`            //返回信息
	AppId      string `xml:"appid,omitempty"`                 //公众账号ID
	MchId      string `xml:"mch_id,omitempty,omitempty"`      //商户号
	DeviceInfo string `xml:"device_info,omitempty,omitempty"` //设备号
	NonceStr   string `xml:"nonce_str,omitempty,omitempty"`   //随机字符串
	Sign       string `xml:"sign,omitempty,omitempty"`        //签名
	ResultCode string `xml:"result_code,omitempty"`           //业务结果
	ErrCode    string `xml:"err_code,omitempty"`              //错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"`          //错误代码描述

	//
	OpenId        string `xml:"openid,omitempty"`         //用户标识
	IsSubscribe   string `xml:"is_subscribe,omitempty"`   //是否关注公众账号
	TradeType     string `xml:"trade_type,omitempty"`     //交易类型
	BankType      string `xml:"bank_type,omitempty"`      //付款银行
	FeeType       string `xml:"fee_type,omitempty"`       //货币类型
	TotalFee      int    `xml:"total_fee,omitempty"`      //总金额
	CashFeeType   string `xml:"cash_fee_type,omitempty"`  //现金支付货币类型
	CashFee       int    `xml:"cash_fee,omitempty"`       //现金支付金额
	CouponFee     int    `xml:"coupon_fee,omitempty"`     //代金券或立减优惠金额
	TransactionId string `xml:"transaction_id,omitempty"` //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no,omitempty"`   //商户订单号
	Attach        string `xml:"attach,omitempty"`         //商家数据包
	TimeEnd       string `xml:"time_end,omitempty"`       //支付完成时间
}

type OrderqueryRequest struct {
	XMLName xml.Name `xml:"xml"`
	// optional
	TransactionId string `xml:"transaction_id,omitempty"` //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no,omitempty"`   //商户订单号

	// needed
	AppId string `xml:"appid"`            //公众账号ID
	MchId string `xml:"mch_id,omitempty"` //商户号

	Sign      string `xml:"sign"`                //签名
	NonceStr  string `xml:"nonce_str,omitempty"` //随机字符串
	NotifyUrl string `xml:"-"`
	SubMchId  string `xml:"sub_mch_id,omitempty"`
}
type OrderqueryResponse struct {
	//
	ReturnCode string `xml:"return_code,omitempty"`           //返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"`            //返回信息
	AppId      string `xml:"appid,omitempty"`                 //公众账号ID
	MchId      string `xml:"mch_id,omitempty,omitempty"`      //商户号
	DeviceInfo string `xml:"device_info,omitempty,omitempty"` //设备号
	NonceStr   string `xml:"nonce_str,omitempty,omitempty"`   //随机字符串
	Sign       string `xml:"sign,omitempty,omitempty"`        //签名
	ResultCode string `xml:"result_code,omitempty"`           //业务结果
	ErrCode    string `xml:"err_code,omitempty"`              //错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"`          //错误代码描述

	//
	OpenId        string `xml:"openid,omitempty"`         //用户标识
	IsSubscribe   string `xml:"is_subscribe,omitempty"`   //是否关注公众账号
	TradeType     string `xml:"trade_type,omitempty"`     //交易类型
	BankType      string `xml:"bank_type,omitempty"`      //付款银行
	FeeType       string `xml:"fee_type,omitempty"`       //货币类型
	TotalFee      int    `xml:"total_fee,omitempty"`      //总金额
	CashFeeType   string `xml:"cash_fee_type,omitempty"`  //现金支付货币类型
	CashFee       int    `xml:"cash_fee,omitempty"`       //现金支付金额
	CouponFee     int    `xml:"coupon_fee,omitempty"`     //代金券或立减优惠金额
	TransactionId string `xml:"transaction_id,omitempty"` //微信支付订单号
	OutTradeNo    string `xml:"out_trade_no,omitempty"`   //商户订单号
	Attach        string `xml:"attach,omitempty"`         //商家数据包
	TimeEnd       string `xml:"time_end,omitempty"`       //支付完成时间
}
