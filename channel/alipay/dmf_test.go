package alipay

// 真实测试，如果参数对的话，是会扣钱的！！！！！
// ScanCodeId 从手机获取扫条码
import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	. "github.com/smartystreets/goconvey/convey"
)

var scanPay = &model.ScanPay{
	// GoodsInfo:    "鞋子,1000,2;衣服,1500,3",
	ChanOrderNum: tools.SerialNumber(),
	Key:          "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ScanCodeId:   "289434710505996982",
	Txamt:        "0.01",
	Subject:      "讯联测试",
}

func TestProcessBarcodePay(t *testing.T) {
	Debug = true
	Convey("支付宝下单", t, func() {
		resp := DefaultClient.ProcessBarcodePay(scanPay)
		Convey("期望", func() {
			So(resp.RespCode, ShouldEqual, "000000")
		})
	})

}
