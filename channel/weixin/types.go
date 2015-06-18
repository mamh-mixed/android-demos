package weixinpay

var wxr = weiXinRequest{}

// weiXinRequest 请求参数
type weiXinRequest struct {
	// 可选
	DeviceInfo string `xml:"device_info"` //设备号
	Detail     string `xml:"detail"`      //商品详情
	Attach     string `xml:"attach"`      //附加数据
	FeeType    string `xml:"fee_type"`    //货币类型
	GoodsTag   string `xml:"goods_tag"`   //商品标记

	// 必选
	AppId          string `xml:"appid"`            //公众账号ID
	MchId          string `xml:"mch_id"`           //商户号
	NonceStr       string `xml:"nonce_str"`        //随机字符串
	Body           string `xml:"body"`             //商品描述
	OutTradeNo     string `xml"out_trade_no"`      //商户订单号
	TotalFee       int    `xml:"total_fee"`        //总金额
	SpbillCreateIp string `xml:"spbill_create_ip"` //终端IP
	AuthCode       string `xml:"auto_code"`        //授权码
	Sign           string `xml:"sign"`             //签名
}

// weiXinResponse 微信接口返回信息
type weiXinResponse struct {
	ReturnCode string
	ReturnMsg  string
}
