package scanpay

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	// 下单支付
	scanPayBarcodePay = &model.ScanPayRequest{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  util.Millisecond(),
		// OrderNum:   "289526980697905109",
		ScanCodeId: "286235706060278591",
		Inscd:      "CIL00002",
		Txamt:      "000000000002",
		Busicd:     "PURC",
		Mchntid:    "100000000000203",
		Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}
	// 预下单支付
	scanPayQrCodeOfflinePay = &model.ScanPayRequest{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  util.Millisecond(),
		Inscd:     "CIL00002",
		Txamt:     "000000000001",
		Busicd:    "PAUT",
		Mchntid:   "100000000000203",
		Chcd:      "ALP",
	}
	// 查询
	scanPayEnquiry = &model.ScanPayRequest{
		Busicd:       "INQY",
		Mchntid:      "100000000000203",
		Inscd:        "CIL00002",
		OrigOrderNum: "1437632220399",
	}
	// 退款
	scanPayRefund = &model.ScanPayRequest{
		Busicd:       "REFD",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1437634914781",
		Inscd:        "CIL00002",
		Txamt:        "000000000002",
		Chcd:         "WXP",
	}
	// 撤销
	scanPayCancel = &model.ScanPayRequest{
		Busicd:       "VOID",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "14376326273610",
		Inscd:        "CIL00002",
	}
	// 关单
	scanPayClose = &model.ScanPayRequest{
		Busicd:       "CANC",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1437632003684",
		Inscd:        "CIL00002",
	}

	scanPay = scanPayRefund
)

func TestScanPay(t *testing.T) {

	log.SetOutputLevel(log.Ldebug)
	reqBytes, _ := json.Marshal(scanPay)
	e := mahonia.NewEncoder("gbk")
	gbk := e.ConvertString(string(reqBytes))

	respBytes := TcpScanPayHandle([]byte(gbk))
	respStr := string(respBytes)

	d := mahonia.NewDecoder("gbk")
	utf8 := d.ConvertString(respStr)

	resp := new(model.ScanPayResponse)
	err := json.Unmarshal([]byte(utf8[4:]), resp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 预期结果
	switch scanPay.Busicd {
	case "purc":
		Convey("下单", t, func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	case "paut":
		Convey("预下单", t, func() {
			So(resp.Respcd, ShouldEqual, "09")
		})
	case "inqy":
		Convey("查询", t, func() {
			So(resp.Respcd, ShouldNotEqual, "")
		})
	case "refd":
		Convey("退款", t, func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	case "void":
		Convey("撤销", t, func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	case "canc":
		Convey("关单", t, func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	}
	t.Logf("%+v", resp)
}

func TestSignMsg(t *testing.T) {
	t.Log(scanPayBarcodePay.SignMsg())
}
