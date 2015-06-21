package alipay

// 真实测试，如果参数对的话，是会扣钱的！！！！！
// ScanCodeId 从手机获取扫条码
import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var pay = &model.ScanPay{
	// GoodsInfo:    "鞋子,1000,2;衣服,1500,3",
	SysOrderNum: tools.SerialNumber(),
	Key:         "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ScanCodeId:  "289253580324978839",
	Txamt:       "0.01",
	Subject:     "讯联测试",
}

var prePay = &model.ScanPay{
	// GoodsInfo:    "鞋子,1000,2;衣服,1500,3",
	SysOrderNum: tools.SerialNumber(),
	Key:         "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	Txamt:       "0.01",
	Subject:     "讯联测试",
}

var enquiry = &model.ScanPay{
	// SysOrderNum: "fc718816621f4bc47fc09ccba1c66304",
	SysOrderNum: "30723e05292b42df78c0e8fd8acb474e",
	Key:         "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
}

func TestProcessBarcodePay(t *testing.T) {

	// 默认开启调试
	// Debug = false
	Convey("支付宝下单", t, func() {
		resp := DefaultClient.ProcessBarcodePay(pay)
		Convey("期望", func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	})

}

func TestProcessQrCodeOfflinePay(t *testing.T) {

	// 默认开启调试
	log.SetOutputLevel(log.Linfo)
	log.Infof("%+v", prePay)
	Convey("支付宝预下单", t, func() {
		resp := DefaultClient.ProcessQrCodeOfflinePay(prePay)
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
		resp := DefaultClient.ProcessEnquiry(enquiry)
		Convey("期望", func() {
			So(resp.Respcd, ShouldEqual, "00")
		})
	})

}
