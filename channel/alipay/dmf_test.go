package alipay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	. "github.com/smartystreets/goconvey/convey"
)

var scanPay = &model.ScanPay{
	GoodsInfo:       "鞋子,1000,2;衣服,1500,3",
	ChannelOrderNum: "awdajwdadn",
	Key:             "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ScanCodeId:      "23131242413",
	Txamt:           "0.01",
}

func TestProcessBarcodePay(t *testing.T) {

	Convey("test the alp BarcodePay", t, func() {
		resp := DefaultClient.ProcessBarcodePay(scanPay)
		Convey("the resp should be fail", func() {
			So(resp.ChannelOrderNum, ShouldEqual, "")
		})
	})

}
