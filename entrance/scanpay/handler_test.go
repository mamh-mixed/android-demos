package scanpay

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"github.com/omigo/mahonia"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	// 下单支付
	scanPayBarcodePay = &model.ScanPay{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  tools.Millisecond(),
		// ScanCodeId: "282259453320278456", // 支付宝
		ScanCodeId: "130381911129127781", // 微信
		Inscd:      "CIL00002",
		Txamt:      "000000000001",
		Busicd:     "PURC",
		Mchntid:    "100000000000203",
		Subject:    "test",
	}
	// 预下单支付
	scanPayQrCodeOfflinePay = &model.ScanPay{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  tools.Millisecond(),
		Inscd:     "CIL00002",
		Txamt:     "000000000001",
		Busicd:    "PAUT",
		Mchntid:   "100000000000203",
		Chcd:      "ALP",
	}
	// 查询
	scanPayEnquiry = &model.ScanPay{
		Busicd:       "INQY",
		Mchntid:      "100000000000203",
		Inscd:        "CIL00002",
		OrigOrderNum: "1435658612934",
	}
	// 退款
	scanPayRefund = &model.ScanPay{
		Busicd:       "REFD",
		Mchntid:      "100000000000203",
		OrderNum:     tools.Millisecond(),
		OrigOrderNum: "1435658612934",
		Inscd:        "CIL00002",
		Txamt:        "000000000001",
	}
	// 撤销
	scanPayCancel = &model.ScanPay{
		Busicd:       "VOID",
		Mchntid:      "100000000000203",
		OrderNum:     tools.Millisecond(),
		OrigOrderNum: "1435658612934",
		Inscd:        "CIL00002",
	}
	// 关单
	scanPayClose = &model.ScanPay{
		Busicd:       "CANC",
		Mchntid:      "100000000000203",
		OrderNum:     tools.Millisecond(),
		OrigOrderNum: "1435726232974",
		Inscd:        "CIL00002",
	}

	scanPay = scanPayClose
)

func TestScanPay(t *testing.T) {

	log.SetOutputLevel(log.Ldebug)
	reqBytes, _ := json.Marshal(scanPay)
	e := mahonia.NewEncoder("gbk")
	gbk := e.ConvertString(string(reqBytes))

	respBytes := ScanPayHandle([]byte(gbk))
	respStr := string(respBytes)
	resp := new(model.ScanPayResponse)
	err := json.Unmarshal([]byte(respStr[4:]), resp)
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
