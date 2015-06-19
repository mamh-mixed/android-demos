package weixin

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/mahonia"
	//"github.com/omigo/log"
)

type WeixinPay struct{}

var DefaultClient WeixinPay

const (
	MicroPay = iota
	OrderQuery
	Reverse
	Refund
	RefundQuery
)

// ProcessBarcodePay 扫条码下单
func (c *WeixinPay) ProcessBarcodePay(scanPayReq *model.ScanPay) *model.ScanPayResponse {

	microPayReq := PerpareRequestData(scanPayReq, MicroPay)

	microPayResp := RequestWeixin(microPayReq, scanPayReq.NotifyUrl)

	//log.Debugf("micropay response: %+v", buf)
	return transformToScanPayResp(microPayResp)
}

//
func (c *WeixinPay) ProcessEnquiry(scanPayReq *model.ScanPay) *model.ScanPayResponse {

	orderqueryReq := PerpareRequestData(scanPayReq, OrderQuery)

	orderqueryResp := RequestWeixin(orderqueryReq, scanPayReq.NotifyUrl)

	return transformToScanPayResp(orderqueryResp)
}

func PerpareRequestData(scanPayReq *model.ScanPay, businessType int) WeixinRequest {

	var weixinRequest WeixinRequest

	switch businessType {
	case MicroPay:
		weixinRequest = &MicropayRequest{
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

	case OrderQuery:
		weixinRequest = &OrderqueryRequest{
			AppId:    appid,
			MchId:    scanPayReq.Mchntid,
			NonceStr: "random string",

			OutTradeNo: scanPayReq.OrderNum,
			NotifyUrl:  scanPayReq.NotifyUrl,
			SubMchId:   sub_mch_id,
		}
	default:
		log.Fatal(errors.New("should not be here"))
	}
	sign := calculateSign(weixinRequest, md5Key)

	weixinRequest.setSign(sign)

	return weixinRequest
}

func RequestWeixin(m WeixinRequest, url string) WeixinResponse {
	var weixinResp WeixinResponse

	buf, err := xml.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))

	bytebuf := bytes.NewBuffer(buf)

	// the http response from weixin
	rep, _ := http.Post(url, "text/xml", bytebuf)

	fmt.Println("rep", rep)

	body := rep.Body
	fmt.Println("Body:", body)
	defer rep.Body.Close()
	d := xml.NewDecoder(body)

	d.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
		dec := mahonia.NewDecoder(s)
		if dec == nil {
			return nil, fmt.Errorf("not support %s", s)
		}
		return dec.NewReader(r), nil
	}

	switch m.(type) {

	case *MicropayRequest:
		weixinResp = new(MicroPayResponse)
	case *OrderqueryRequest:
		weixinResp = new(OrderqueryResponse)
	default:
		log.Fatal("should not be here")
	}
	err = d.Decode(weixinResp)

	if err != nil {
		// log.Errorf("unmarsal body fail : %s", err)
		log.Fatalf("unmarsal body fail : %s", err)
	}
	return weixinResp
}

func transformToScanPayResp(sp WeixinResponse) *model.ScanPayResponse {
	fmt.Println("weixinResponse:", sp)
	switch sp.(type) {
	case *MicroPayResponse:
		fmt.Println("it is MicroPayResponse")
	case *OrderqueryResponse:
		fmt.Println("it is OrderqueryResponse")
	default:
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
	ret := new(model.ScanPayResponse)
	//
	// if sp.ReturnCode == "SUCCESS" {
	// 	// normal connection
	// 	if sp.ResultCode == "SUCCESS" {
	// 		fmt.Println("request success")
	//
	// 		ret.Busicd = sp.TradeType
	// 		ret.Respcd = sp.ResultCode
	// 		ret.Mchntid = sp.MchId
	//
	// 	} else if sp.ResultCode == "FAIL" {
	// 		fmt.Println("request fail")
	// 		ret.Respcd = sp.ResultCode
	// 		ret.ErrorDetail = sp.ReturnMsg
	// 		ret.Mchntid = sp.MchId
	// 		ret.Sign = sp.Sign
	// 	}
	// } else {
	// 	// inormal connection
	// 	fmt.Println("connect failure")
	// }
	return ret
}
