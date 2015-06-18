package weixin

// micropayRequest 请求参数
type MicropayRequest struct {
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
	OutTradeNo     string `xml"out_trade_no,omitempty"`      //商户订单号
	TotalFee       int    `xml:"total_fee,omitempty"`        //总金额
	SpbillCreateIp string `xml:"spbill_create_ip,omitempty"` //终端IP
	AuthCode       string `xml:"auto_code,omitempty"`        //授权码
	Sign           string `xml:"sign,omitempty"`             //签名
}

type OrderqueryRequest struct {
}
