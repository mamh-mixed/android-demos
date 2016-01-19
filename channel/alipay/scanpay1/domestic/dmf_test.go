package domestic

// 真实测试，如果参数对的话，是会扣钱的！！！！！
// ScanCodeId 从手机获取扫条码
import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var pay = &model.ScanPayRequest{
	GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
	OrderNum:   util.SerialNumber(),
	SignKey:    "86l3l20oagn2afs0r0ztkizut1il66ec",
	ScanCodeId: "280499860770919934",
	ActTxamt:   "0.01",
	Subject:    "讯联测试",
	ChanMerId:  "2088811767473826",
}

var prePay = &model.ScanPayRequest{
	GoodsInfo: "鞋子,1000,2;衣服,1500,3",
	OrderNum:  util.Millisecond(),
	SignKey:   "86l3l20oagn2afs0r0ztkizut1il66ec",
	ActTxamt:  "0.01",
	Subject:   "讯联测试",
	ChanMerId: "2088811767473826",
}

var cancelPay = &model.ScanPayRequest{
	SignKey:      "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	Subject:      "讯联测试",
	ChanMerId:    "2088811767473826",
	OrigOrderNum: "1435564308178",
}

var enquiry = &model.ScanPayRequest{
	OrderNum:  "e148a25a84f14024511c5f3cde5d4594",
	SignKey:   "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ChanMerId: "2088811767473826",
}

var refundPay = &model.ScanPayRequest{
	OrderNum:     util.Millisecond(),
	OrigOrderNum: "00a4371518554214622c801a9a158128",
	SignKey:      "86l3l20oagn2afs0r0ztkizut1il66ec",
	ChanMerId:    "2088811767473826",
	ActTxamt:     "0.01",
}

var settle = &model.ScanPayRequest{
	SettDate:  "2015-12-10",
	SignKey:   "wg5txarw1shatk0boc61di3971lgl8xe",
	ChanMerId: "2088121476326615",
}

func TestProcessBarcodePay(t *testing.T) {

	// 默认开启调试
	// Debug = false
	log.SetOutputLevel(log.Ldebug)
	Convey("支付宝下单", t, func() {
		resp, _ := DefaultClient.ProcessBarcodePay(pay)
		Convey("期望", func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	})

}

func TestProcessQrCodeOfflinePay(t *testing.T) {

	// 默认开启调试
	log.SetOutputLevel(log.Ldebug)
	log.Infof("%+v", prePay)
	Convey("支付宝预下单", t, func() {
		resp, err := DefaultClient.ProcessQrCodeOfflinePay(prePay)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		log.Infof("%+v", resp)
		Convey("期望", func() {
			So(resp.Respcd, ShouldEqual, "09")
		})
	})

}

func TestProcessEnquiry(t *testing.T) {

	// 默认开启调试
	// log.SetOutputLevel(log.Linfo)
	Convey("支付宝订单查询", t, func() {
		resp, _ := DefaultClient.ProcessEnquiry(enquiry)
		Convey("期望", func() {
			So(resp.Respcd, ShouldNotEqual, "")
		})
	})

}

func TestProcessRefund(t *testing.T) {

	// 默认开启调试
	log.SetOutputLevel(log.Ldebug)
	Convey("支付宝退款订单", t, func() {
		resp, _ := DefaultClient.ProcessRefund(refundPay)
		Convey("期望", func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	})
}

func TestProcessCancel(t *testing.T) {

	// 默认开启调试
	log.SetOutputLevel(log.Ldebug)
	Convey("支付宝撤销订单", t, func() {
		resp, _ := DefaultClient.ProcessCancel(cancelPay)
		Convey("期望", func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	})
}

func TestProcessSettleEnquiry(t *testing.T) {
	cbd := make(model.ChanBlendMap)
	err := DefaultClient.ProcessSettleEnquiry(settle, cbd)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("%+v", cbd)

}
