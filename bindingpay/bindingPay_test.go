package bindingpay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

const (
	// testMerId      = "012345678901234"
	testMerId = "000000001405"
	// testSign       = "0123456789"
	testSign = "0123456789"
	// testEncryptKey = "AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8="
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
	bindingId string = "012345678901234"
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
	req.Header.Set("X-Sign", security.SHA1WithKey(string(j), testSign))

	w := httptest.NewRecorder()
	BindingpayHandle(w, req)
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
	url := "http://quick.ipay.so/bindingpay/bindingCreate?merId=" + testMerId

	bindingId = util.Millisecond()

	b := model.BindingCreate{
		MerId:     testMerId,
		BindingId: bindingId,
		AcctName:  "张三",
		AcctNum:   "6222022003008481261",
		IdentType: "0",
		IdentNum:  "440583199111031012",
		PhoneNum:  "15618103236",
		AcctType:  "10",
		// ValidDate: "0612",
		// Cvv2:      "793",
		SendSmsId: "",
		SmsCode:   "",
		// BankId:    "102",
	}

	var aes = security.NewAESCBCEncrypt(testEncryptKey)
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

func TestPaySettlement(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/bindingPaymentSettlement?merId=" + "000000001406"
	b := model.PaySettlement{
		TerminalId:      "00001000",
		MerOrderNum:     util.Millisecond(),
		SettOrderNum:    "Iris20150917",
		SettAmt:         10,
		SettAccountType: "11",
		SettAccountName: "yYllsyq9k5dCqLJNn/LE1gCKixdiKrrZ49CWVoMg9cY=",
		SettAccountNum:  "goP5u0z9B02CHTCoYt07GKALiW2xKtgTcsuimcA643jNzl3DE7T2omY8cqYffjXF",
		SettBranchName:  "中国工商银行",
		Province:        "上海市",
		City:            "上海市",
	}
	doPost(url, b, t)
}

func TestGetCardInfo(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/getCardInfo?merId=" + testMerId
	b := model.CardInfo{CardNum: "6225768739233847"}
	doPost(url, b, t)
}

func TestBindingEnquiryHandle(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/bindingEnquiry?merId=" + testMerId
	b := model.BindingEnquiry{BindingId: bindingId}
	doPost(url, b, t)
}

func TestBindingPaymentHandle(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/bindingPayment?merId=" + testMerId
	orderNum = util.Millisecond()
	b := model.BindingPayment{
		MerOrderNum: orderNum,
		TransAmt:    1000,
		BindingId:   bindingId,
		MerId:       testMerId,
	}
	doPost(url, b, t)
}

func TestSendBindingPaySMS(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/sendBindingPaySms?merId=" + testMerId
	orderNum = util.Millisecond()
	b := model.BindingPayment{
		MerOrderNum: orderNum,
		TransAmt:    1000,
		BindingId:   bindingId,
		MerId:       testMerId,
	}
	doPost(url, b, t)
}

func TestBindingPayWithSMS(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/bindingPayWithSms?merId=" + testMerId
	orderNum = "1440667764327"
	b := model.BindingPayment{
		MerOrderNum: orderNum,
		MerId:       testMerId,
		SmsCode:     "123456",
	}
	doPost(url, b, t)
}

func TestOrderEnquiry(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/orderEnquiry?merId=" + testMerId
	b := model.OrderEnquiry{
		OrigOrderNum: orderNum,
	}
	doPost(url, b, t)
}

func TestBindingRefundHandle(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/refund?merId=" + testMerId
	b := model.BindingRefund{
		OrigOrderNum: "1440643274621",
		MerOrderNum:  util.Millisecond(),
		TransAmt:     1000,
	}
	doPost(url, b, t)
}

func TestBindingRemoveHandle(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/bindingRemove?merId=" + testMerId
	b := model.BindingEnquiry{BindingId: bindingId}
	doPost(url, b, t)
}

func TestBillingSummaryHandle(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/billingSummary?merId=001405"
	b := model.BillingSummary{SettDate: "2015-05-20"}
	doPost(url, b, t)
}

func TestBillingDetailsHandle(t *testing.T) {
	url := "http://quick.ipay.so/bindingpay/billingDetails?merId=001405"
	b := model.BillingDetails{SettDate: "2015-05-20"}
	doPost(url, b, t)
}
