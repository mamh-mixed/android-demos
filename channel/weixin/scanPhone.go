package weixin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/CardInfoLink/quickpay/model"
)

const (
	md5Key     = "12sdffjjguddddd2widousldadi9o0i1"
	mch_id     = "1236593202"
	appid      = "wx25ac886b6dac7dd2"
	acqfee     = "0.02"
	merfee     = "0.03"
	fee        = "0.01"
	sub_mch_id = "1247075201"
	url        = "https://api.mch.weixin.qq.com/pay/micropay"
)

// ProcessBarcodePay 扫条码下单
func ProcessBarcodePay(scanPayReq *model.ScanPay) *model.ScanPayResponse {

	microPayReq := PerpareRequestToWeiXin(scanPayReq)

	buf := RequestWeixin(microPayReq)

	log.Debugf("micropay response: %+v", microPayResp)

	scanPayRep := transformToScanPayResp(buf, scanPayReq.Response)

	return scanPayRep
}

func PerpareRequestToWeiXin(req *model.ScanPay) *MicropayRequest {

	// initial request struct to weixin
	microPay := &MicropayRequest{
		/*
			// data used to test
				AppId:          "sdfsd",
				MchId:          "werwer1231",
				NonceStr:       "xxxxoooo",
				TotalFee:       12,
				OutTradeNo:     "fdsfsfdsf",
				FeeType:        "CNY",
				SpbillCreateIp: "10.10.10.1",
				Body:           "sdfdfsdds",
				AuthCode:       "sfdsfafd",
		*/
		AppId:    appid,
		MchId:    mch_id,
		NonceStr: getRandomStr(),
		//TotalFee:       req.Txamt,
		TotalFee:       67,
		OutTradeNo:     req.OrderNum,
		FeeType:        "CNY",
		SpbillCreateIp: "10.10.10.1",
		Body:           "sdfdfsdds",
		AuthCode:       req.ScanCodeId,
		/*
			// optional data
				DeviceInfo :req.
				GoodsTag       :req.
				Detail     :req.
				Attach     :req.
				Sign           :req.
		*/
	}

	setSign(microPay, md5Key)

	return microPay
}

func setSign(microPay *MicropayRequest, md5Key string) {
	dict := toMapWithKeySortedAndValueNotNil(microPay)
	var keys []string
	for k, _ := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	for _, v := range keys {
		buffer.WriteString(v)
		buffer.WriteString("=")
		buffer.WriteString(dict[v])
		buffer.WriteString("&")
	}
	seq := buffer.String() + md5Key
	signSlice := md5.Sum([]byte(seq))
	microPay.Sign = hex.EncodeToString(signSlice[:])
}

func RequestWeixin(m *MicropayRequest) []byte {
	buf, err := xml.MarshalIndent(microPay, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))

	body := bytes.NewBuffer(buf)
	r, _ := http.Post(url, "text/xml", body)

	//
	//d := xml.NewDecoder(r.Body)

	buf, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(buf))

	return buf
}

func transformToScanPayResp(buf []byte, response *ScanPayResponse) {

}
