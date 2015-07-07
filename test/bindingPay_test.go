package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CardInfoLink/quickpay/core"
	. "github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

const (
	testMerId      = "012345678901234"
	testSign       = "0123456789"
	testEncryptKey = "AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8="

	// 银联卡测试数据
	testCUPCard      = "6225220100740059"
	testCUPCVV2      = "111"
	testCUPValidDate = "1605"
	testCUPPhone     = "13611111111"
	testCUPIdentNum  = "130412"

	testMSCCard       = "5457210001000019"
	testMSCCVV2       = "300"
	testMSCValidDate  = "1412"
	testMSCTrackdata2 = "5457210001000019=1412101080080748"
	// 无卡直接支付相关
	testMerID = "APPTEST"
)

var (
	bindingId string = "1431654779479"
	orderNum  string
)

func init() {
	// 日志输出级别
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 连接到 MongoDB
	// mongo.Connect()

	// 初始化卡 Bin 树
	core.BuildTree()

	// 连接线下
	// cil.Connect()
}

func doPost(url string, m interface{}, t *testing.T) {
	j, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Sign", SignatureUseSha1(j, testSign))

	w := httptest.NewRecorder()
	Quickpay(w, req)
	log.Infof("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		t.Errorf("response error with status %d", w.Code)
	}

	var out model.BindingReturn
	err = json.Unmarshal(w.Body.Bytes(), &out)
	if err != nil {
		t.Errorf("Unmarshal response error (%s)", err)
	}

	if out.RespCode != "000000" {
		t.Error(out)
	}
}

func TestBindingCreate(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/bindingCreate?merId=" + testMerId

	bindingId = util.Millisecond()

	b := model.BindingCreate{
		MerId:     testMerId,
		BindingId: bindingId,
		AcctName:  "测试账号",
		AcctNum:   "6222022003008481261",
		IdentType: "0",
		IdentNum:  "440583199111031012",
		PhoneNum:  "18205960039",
		AcctType:  "20",
		ValidDate: "0612",
		Cvv2:      "793",
		SendSmsId: "",
		SmsCode:   "",
		BankId:    "102",
	}

	var aes = util.NewAESCBCEncrypt(testEncryptKey)
	b.AcctName = aes.Encrypt(b.AcctName)
	b.AcctNum = aes.Encrypt(b.AcctNum)
	b.IdentNum = aes.Encrypt(b.IdentNum)
	b.PhoneNum = aes.Encrypt(b.PhoneNum)
	b.ValidDate = aes.Encrypt(b.ValidDate)
	b.Cvv2 = aes.Encrypt(b.Cvv2)
	if aes.Err != nil {
		panic(aes.Err)
	}
	doPost(url, b, t)
}

func TestBindingEnquiryHandle(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/bindingEnquiry?merId=" + testMerId
	b := model.BindingEnquiry{BindingId: bindingId}
	doPost(url, b, t)
}

func TestBindingPaymentHandle(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/bindingPayment?merId=" + testMerId
	orderNum = util.Millisecond()
	b := model.BindingPayment{
		MerOrderNum: orderNum,
		TransAmt:    1000,
		BindingId:   bindingId,
		MerId:       testMerId,
	}
	doPost(url, b, t)
}

func TestOrderEnquiry(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/orderEnquiry?merId=" + testMerId
	b := model.OrderEnquiry{
		OrigOrderNum: orderNum,
	}
	doPost(url, b, t)
}

func TestBindingRefundHandle(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/refund?merId=" + testMerId
	b := model.BindingRefund{
		OrigOrderNum: orderNum,
		MerOrderNum:  util.Millisecond(),
		TransAmt:     100,
	}
	doPost(url, b, t)
}

func TestBindingRemoveHandle(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/bindingRemove?merId=" + testMerId
	b := model.BindingEnquiry{BindingId: bindingId}
	doPost(url, b, t)
}

func TestBillingSummaryHandle(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/billingSummary?merId=001405"
	b := model.BillingSummary{SettDate: "2015-05-20"}
	doPost(url, b, t)
}

func TestBillingDetailsHandle(t *testing.T) {
	url := "http://quick.ipay.so/quickpay/billingDetails?merId=001405"
	b := model.BillingDetails{SettDate: "2015-05-20"}
	doPost(url, b, t)
}
