package model

import (
	"encoding/json"
	"strings"
)

// ScanPay 扫码支付
type ScanPay struct {
	Txndir       string //交易方向
	Busicd       string //交易类型
	Inscd        string //机构号
	Chcd         string //渠道机构
	Mchntid      string //商户号
	Txamt        string //订单金额
	GoodsInfo    string //商品详情
	OrderNum     string //订单号
	OrigOrderNum string //原订单号
	ScanCodeId   string //扫码号
	Sign         string //签名
	NotifyUrl    string //异步通知地址
	// 辅助字段
	Key             string // md5key
	Subject         string // 商品名称
	SysOrderNum     string //渠道交易号
	OrigSysOrderNum string // 原系统订单号
}

// ScanPayResponse 下单支付返回体
// M:返回时必须带上
// C:可选
type ScanPayResponse struct {
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
	ChanRespCode string `json:"-"` // 渠道详细应答码
}

// MarshalGoods 将商品详情解析成字符json字符串
// 格式: 商品名称,价格,数量;商品名称,价格,数量;...
func (s *ScanPay) MarshalGoods() string {

	if s.GoodsInfo == "" {
		return ""
	}

	goods := strings.Split(s.GoodsInfo, ";")
	gs := make([]interface{}, 0, len(goods))

	for i, v := range goods {
		good := strings.Split(v, ",")
		if len(good) != 3 {
			return ""
		}
		g := &struct {
			GoodsId   int    `json:"goodsId"`
			GoodsName string `json:"goodsName"`
			Price     string `json:"price"`
			Quantity  string `json:"quantity"`
		}{
			i, good[0], good[1], good[2],
		}
		gs = append(gs, g)
	}
	formated, err := json.Marshal(gs)
	if err != nil {
		return ""
	}
	return string(formated)
}
