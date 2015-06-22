package weixin

import (
	"encoding/xml"
	"fmt"

	"github.com/CardInfoLink/quickpay/model"
)

/*
 提交被扫支付API
*/

// 请求参数
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

//
type MicroPayResponse struct {
	//当return_code为SUCCESS的时候
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

	//当return_code 和result_code都为SUCCESS的时
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

func (w *MicropayRequest) copyData(scanPayReq *model.ScanPay) {
	*w = MicropayRequest{
		AppId:    appid,
		MchId:    scanPayReq.Mchntid,
		NonceStr: "random string",

		TotalFee:       toInt(scanPayReq.Txamt),
		OutTradeNo:     scanPayReq.OrderNum,
		FeeType:        "CNY",
		SpbillCreateIp: "10.10.10.1",
		Body:           scanPayReq.Subject,
		AuthCode:       scanPayReq.ScanCodeId,
		SubMchId:       sub_mch_id,
		NotifyUrl:      scanPayReq.NotifyUrl,
	}
}

func (sp *MicroPayResponse) convertToScanPayResp() *model.ScanPayResponse {
	ret := new(model.ScanPayResponse)

	if sp.ReturnCode == "SUCCESS" {
		// normal connection
		if sp.ResultCode == "SUCCESS" {

			ret.Busicd = sp.TradeType
			ret.Respcd = sp.ResultCode
			ret.Mchntid = sp.MchId

		} else if sp.ResultCode == "FAIL" {
			ret.Respcd = sp.ResultCode
			ret.ErrorDetail = sp.ReturnMsg
			ret.Mchntid = sp.MchId
			ret.Sign = sp.Sign
		}
	} else {
		// inormal connection
		fmt.Println("connect failure")
	}

	return ret
}

/*
 查询订单
*/

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

func (w *OrderqueryRequest) copyData(scanPayReq *model.ScanPay) {
	*w = OrderqueryRequest{
		AppId:    appid,
		MchId:    scanPayReq.Mchntid,
		NonceStr: "random string",

		OutTradeNo: scanPayReq.OrderNum,
		NotifyUrl:  scanPayReq.NotifyUrl,
		SubMchId:   sub_mch_id,
	}
}

func (sp *OrderqueryResponse) convertToScanPayResp() *model.ScanPayResponse {
	ret := new(model.ScanPayResponse)

	if sp.ReturnCode == "SUCCESS" {
		// normal connection
		if sp.ResultCode == "SUCCESS" {

			ret.Busicd = sp.TradeType
			ret.Respcd = sp.ResultCode
			ret.Mchntid = sp.MchId

		} else if sp.ResultCode == "FAIL" {
			ret.Respcd = sp.ResultCode
			ret.ErrorDetail = sp.ReturnMsg
			ret.Mchntid = sp.MchId
			ret.Sign = sp.Sign
		}
	} else {
		// inormal connection
		fmt.Println("connect failure")
	}

	return ret
}

/*
   Txndir          string `json:"txndir"`                    // 交易方向 M M
   Busicd          string `json:"busicd"`                    // 交易类型 M M
   Respcd          string `json:"respcd"`                    // 交易结果  M
   Inscd           string `json:"inscd,omitempty"`           // 机构号 M M
   Chcd            string `json:"chcd,omitempty"`            // 渠道 C C
   Mchntid         string `json:"mchntid"`                   // 商户号 M M
   Txamt           string `json:"txamt,omitempty"`           // 订单金额 M M
   ChannelOrderNum string `json:"channelOrderNum,omitempty"` // 渠道交易号 C
   ConsumerAccount string `json:"consumerAccount,omitempty"` // 渠道账号  C
   ConsumerId      string `json:"consumerId,omitempty"`      // 渠道账号ID   C
   ErrorDetail     string `json:"errorDetail,omitempty"`     // 错误信息   C
   OrderNum        string `json:"orderNum,omitempty"`        //订单号 M C
   OrigOrderNum    string `json:"origOrderNum,omitempty"`    //源订单号 M C
   Sign            string `json:"sign"`                      //签名 M M
   ChcdDiscount    string `json:"chcdDiscount,omitempty"`    //渠道优惠  C
   MerDiscount     string `json:"merDiscount,omitempty"`     // 商户优惠  C
   QrCode          string `json:"qrcode,omitempty"`          // 二维码 C
   // 辅助字段
   RespCode     string `json:"-"` // 系统应答码
   ChanRespCode string `json:"-"` // 渠道详细应答码
*/
