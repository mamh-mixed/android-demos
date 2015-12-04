package scanpay

import "github.com/CardInfoLink/quickpay/channel/weixin"

// 参考文档https://pay.weixin.qq.com/wiki/doc/api/native_sl.php?chapter=9_6
// 商户可以通过该接口下载历史交易清单。比如掉单、系统错误等导致商户侧和微信侧数据不一致，通过对账单核对后可校正支付状态。
//注意：
//1、微信侧未成功下单的交易不会出现在对账单中。支付成功后撤销的交易会出现在对账单中，跟原支付单订单号一致，bill_type为REVOKED；
//2、微信在次日9点启动生成前一天的对账单，建议商户10点后再获取；
//3、对账单中涉及金额的字段单位为“元”。
// 接口链接 https://api.mch.weixin.qq.com/pay/orderquery
// 是否需要证书 :不需要

// SettleQueryReq 请求账单提交的数据
type SettleQueryReq struct {
	weixin.CommonParams

	SettleDate string `xml:"bill_date,omitempty" url:"bill_date,omitempty"` // 对账日期
	SettleType string `xml:"bill_type,omitempty" url:"bill_type,omitempty"` // 对账日期
}

// GetURI 取接口地址
func (p *SettleQueryReq) GetURI() string {
	return "/pay/downloadbill"
}

// SettleQueryResp 账单请求数据被Post到API之后，API会返回退货的数据，这个类用来装这些数据
type SettleQueryResp struct {
	weixin.CommonBody
}
