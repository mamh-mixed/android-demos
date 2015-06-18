package weixin

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
)

var pay = &model.ScanPay{
	// GoodsInfo:    "鞋子,1000,2;衣服,1500,3",
	OrderNum: tools.SerialNumber(),
	Key:      "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	Txamt:    "10",
}

// var quiry = &model.ScanPay{
// 	// SysOrderNum: "fc718816621f4bc47fc09ccba1c66304",
// 	Key: "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
// }

func TestProcessBarcodePay(t *testing.T) {

	// 默认开启调试
	// Debug = false
	Convey("微信下单", t, func() {
		resp := DefaultClient.ProcessBarcodePay(pay)
		Convey("期望", func() {
			So(resp.RespCode, ShouldEqual, "000000")
		})
	})

}
