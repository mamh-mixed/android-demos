package alipay

// 真实测试，如果参数对的话，是会扣钱的！！！！！
// ScanCodeId 从手机获取扫条码
import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	. "github.com/smartystreets/goconvey/convey"
)

var pay = &model.ScanPayRequest{
	GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
	OrderNum:   util.SerialNumber(),
	SignCert:   "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ScanCodeId: "281095002166152011",
	ActTxamt:   "0.01",
	Subject:    "讯联测试",
	ChanMerId:  "2088811767473826",
}

var prePay = &model.ScanPayRequest{
	GoodsInfo: "鞋子,1000,2;衣服,1500,3",
	OrderNum:  util.Millisecond(),
	SignCert:  "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ActTxamt:  "0.01",
	Subject:   "讯联测试",
	ChanMerId: "2088811767473826",
}

var cancelPay = &model.ScanPayRequest{
	SignCert:     "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	Subject:      "讯联测试",
	ChanMerId:    "2088811767473826",
	OrigOrderNum: "1435564308178",
}

var enquiry = &model.ScanPayRequest{
	OrderNum:  "e148a25a84f14024511c5f3cde5d4594",
	SignCert:  "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ChanMerId: "2088811767473826",
}

var refundPay = &model.ScanPayRequest{
	OrderNum:     util.Millisecond(),
	OrigOrderNum: "e148a25a84f14024511c5f3cde5d4594",
	SignCert:     "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ChanMerId:    "2088811767473826",
	ActTxamt:     "0.01",
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
		resp, _ := DefaultClient.ProcessQrCodeOfflinePay(prePay)
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

func TestDmf2(t *testing.T) {
	ProcessPreOrder()
}
