package scanpay

import (
	"encoding/xml"

	"github.com/omigo/log"
	"github.com/omigo/validator"
)

// BaseReq 只是为了注入签名方便
type BaseReq interface {
	GenSign()
}

// CommonParams 微信接口请求公共参数
type CommonParams struct {
	XMLName xml.Name `xml:"xml"`

	// 公共字段
	Appid    string `xml:"appid" validate:"len=18"`       // 公众账号ID
	MchID    string `xml:"mch_id" validate:"nonzero"`     // 商户号
	SubMchId string `xml:"sub_mch_id" validate:"nonzero"` // 子商户号（文档没有该字段）
	NonceStr string `xml:"nonce_str" validate:"nonzero"`  // 随机字符串
	Sign     string `xml:"sign"`                          // 签名

	WeixinMD5Key string `xml:"-" validate:"nonzero"`
}

// BaseResp 只是为了传参方便
type BaseResp interface{}

// CommonBody 微信接口返回公共字段
type CommonBody struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"`          // 返回状态码
	ReturnMsg  string `xml:"return_msg,omitempty"` // 返回信息

	// 当 return_code 为 SUCCESS 的时候，还会包括以下字段：
	Appid      string `xml:"appid"`                  // 公众账号ID
	MchID      string `xml:"mch_id"`                 // 商户号
	SubMchId   string `xml:"sub_mch_id"`             // 子商户号（文档没有该字段）
	SubAppid   string `xml:"sub_appid"`              // 子商户公众账号 ID
	NonceStr   string `xml:"nonce_str"`              // 随机字符串
	Sign       string `xml:"sign"`                   // 签名
	ResultCode string `xml:"result_code"`            // 业务结果
	ErrCode    string `xml:"err_code,omitempty"`     // 错误代码
	ErrCodeDes string `xml:"err_code_des,omitempty"` // 错误代码描述
}

func base(d BaseReq, r BaseResp) (err error) {
	if err := validator.Validate(d); err != nil {
		log.Errorf("validate error, %s", err)
		return err
	}

	err = sendRequest(d, r)
	if err != nil {
		log.Errorf("weixin request error: %s", err)
		return err
	}
	return nil
}
