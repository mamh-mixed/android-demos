package scanpay

import (
	"encoding/json"
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	// 下单支付
	scanPayBarcodePay = &model.ScanPay{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  tools.Millisecond(),
		// ScanCodeId: "285533206604506831", // 支付宝
		ScanCodeId: "130523449261557875", // 微信
		Inscd:      "CIL00002",
		Txamt:      "000000000001",
		Busicd:     "purc",
		Mchntid:    "CIL0001",
	}
	// 预下单支付
	scanPayQrCodeOfflinePay = &model.ScanPay{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  tools.Millisecond(),
		Inscd:     "CIL00002",
		Txamt:     "000000000001",
		Busicd:    "paut",
		Mchntid:   "100000000000203",
		Chcd:      "ALP",
	}
	// 查询
	scanPayEnquiry = &model.ScanPay{
		Busicd:       "inqy",
		Mchntid:      "CIL0001",
		Inscd:        "CIL00002",
		OrigOrderNum: "1435306550752",
	}
	// 退款
	scanPayRefund = &model.ScanPay{
		Busicd:       "refd",
		Mchntid:      "CIL0001",
		OrderNum:     tools.Millisecond(),
		OrigOrderNum: "1435568370974",
		Inscd:        "CIL00002",
		Txamt:        "000000000001",
	}
	// 撤销
	scanPayCancel = &model.ScanPay{
		Busicd:       "void",
		Mchntid:      "CIL0001",
		OrderNum:     tools.Millisecond(),
		OrigOrderNum: "1435568762666",
		Inscd:        "CIL00002",
	}
	// 关单
	scanPayClose = &model.ScanPay{
		Busicd:       "canc",
		Mchntid:      "CIL0001",
		OrderNum:     tools.Millisecond(),
		OrigOrderNum: "1435569257167",
		Inscd:        "CIL00002",
	}
)

func TestScanPay(t *testing.T) {
	log.SetOutputLevel(log.Ldebug)
	scanPay := scanPayQrCodeOfflinePay
	reqBytes, _ := json.Marshal(scanPay)
	respBytes := ScanPayHandle(reqBytes)

	resp := new(model.ScanPayResponse)
	err := json.Unmarshal(respBytes, resp)
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
