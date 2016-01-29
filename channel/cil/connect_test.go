package cil

import (
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
)

func TestShortConnect(t *testing.T) {
	msg := newTestCilMsg()
	// t.Logf("msg  = %#v", msg)
	msg.Clisn = time.Now().Format("150405")
	back := send(msg, transTimeout) // 线下网关发送报文

	_ = back
	log.Debug("--------------------------------------------")
}

func TestLongConnect(t *testing.T) {

	Connect()

	msg := newTestCilMsg()
	// t.Logf("msg  = %#v", msg)
	msg.Clisn = time.Now().Format("150405")

	back := send0(msg, transTimeout) // 线下网关发送报文

	_ = back
	log.Debug("--------------------------------------------")

}

func newTestCilMsg() (m *model.CilMsg) {
	m = &model.CilMsg{
		Busicd:          "500000",
		Txndir:          "Q",
		Posentrymode:    "022",
		Chcd:            "00022900",
		Txamt:           "000000001000",
		Txdt:            "0926115934",
		Localdt:         "0926115934",
		Cardcd:          "9559970030000000215",
		Trackdata2:      "9559970030000000215=00002101815546",
		Trackdata3:      "",
		Cardpin:         "",
		Syssn:           "101213113013",
		Clisn:           "115934",
		Inscd:           "30512900",
		Mchntid:         "0002220F0002804",
		Terminalid:      "60000005",
		Mcc:             "4816",
		Txcurrcd:        "156",
		Billingcurr:     "156",
		Regioncd:        "0156",
		Mchntnm:         "shanghai test                           ",
		Nminfo:          "PKE",
		Cardseqnum:      "001",
		Iccdata:         "",
		Termreadability: "5",
		Icccondcode:     "0",
		Outgoingacct:    "9559970030000000215",
		Incomingacct:    "4682030210337444",
		Custmrtp:        "01",
		Custmracnt:      "130412",
		Paymd:           "01",
		Goodscd:         "19100059",
		Billyymm:        "201201",
		Chname:          "",
		Inchname:        "我们都是好孩子",
		Phonenum:        "13611111111",
		Cvv2:            "111",
		Paymethod:       "3",
		Billinscd:       "888880000502900",
		Barcd:           "539100060832536001034816",
		Psamcd:          "1234567890123456",
		Txnmode:         "1",
		Termserialcd:    "1234567890123",
		Expiredate:      "1605",
		Usagetags:       "12",
	}
	return m
}
