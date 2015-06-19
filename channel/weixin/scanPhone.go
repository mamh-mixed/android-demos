package weixin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/mahonia"
	//"github.com/omigo/log"
)

type WeixinPay struct{}

var DefaultClient WeixinPay

// ProcessBarcodePay 扫条码下单
func (c *WeixinPay) ProcessBarcodePay(scanPayReq *model.ScanPay) {

	microPayReq := PerpareRequestToWeiXin(scanPayReq)

	microPayResp := RequestWeixin(microPayReq, scanPayReq.NotifyUrl)

	//log.Debugf("micropay response: %+v", buf)
	transformToScanPayResp(microPayResp, scanPayReq.Response)
}
func (c *WeixinPay) ProcessEnquiry(scanPayReq *model.ScanPay) *model.ScanPayResponse {
	return nil
}

func PerpareRequestToWeiXin(req *model.ScanPay) *MicropayRequest {

	microPayReq := &MicropayRequest{
		AppId:    appid,
		MchId:    req.Mchntid,
		NonceStr: "random string",

		TotalFee:       toInt(req.Txamt),
		OutTradeNo:     req.OrderNum,
		FeeType:        "CNY",
		SpbillCreateIp: "10.10.10.1",
		Body:           req.Subject,
		AuthCode:       req.ScanCodeId,
		SubMchId:       sub_mch_id,
	}

	microPayReq.setSign(md5Key)

	return microPayReq
}

func (microPay *MicropayRequest) setSign(md5Key string) {
	dict := toMapWithValueNotNil(microPay)

	var keys []string
	for k, _ := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	for _, v := range keys {
		buffer.WriteString(v + "=" + dict[v] + "&")
	}
	buffer.WriteString("key=" + md5Key)

	seq := buffer.String()
	signSlice := md5.Sum([]byte(seq))

	microPay.Sign = strings.ToUpper(hex.EncodeToString(signSlice[:]))
	fmt.Println("sign:", microPay.Sign)
}

func RequestWeixin(m *MicropayRequest, url string) *MicroPayResponse {
	buf, err := xml.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))

	body := bytes.NewBuffer(buf)

	r, _ := http.Post(url, "text/xml", body)

	return transformToMicroPayResponse(r)
}

func transformToMicroPayResponse(rep *http.Response) *MicroPayResponse {
	fmt.Println("rep:", rep)
	ret := new(MicroPayResponse)
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
	err := d.Decode(ret)

	if err != nil {
		// log.Errorf("unmarsal body fail : %s", err)
		log.Fatalf("unmarsal body fail : %s", err)
	}
	return ret
}

func transformToScanPayResp(sp *MicroPayResponse, ret *model.ScanPayResponse) {
	fmt.Println("microPayResponse:", sp)
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

	if sp.ReturnCode == "SUCCESS" {
		// normal connection
		if sp.ResultCode == "SUCCESS" {
			fmt.Println("request success")

			ret.Busicd = sp.TradeType
			ret.Respcd = sp.ResultCode
			ret.Mchntid = sp.MchId

		} else if sp.ResultCode == "FAIL" {
			fmt.Println("request fail")
			ret.Respcd = sp.ResultCode
			ret.ErrorDetail = sp.ReturnMsg
			ret.Mchntid = sp.MchId
			ret.Sign = sp.Sign
		}
	} else {
		// inormal connection
		fmt.Println("connect failure")
	}
}
