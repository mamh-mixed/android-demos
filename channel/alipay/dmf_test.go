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
	SysOrderNum: tools.SerialNumber(),
	Key:         "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
	ScanCodeId:  "289253580324978839",
	Txamt:       "10",
	Subject:     "讯联测试",
}

func TestProcessBarcodePay(t *testing.T) {

	// 默认开启调试
	// Debug = false
	Convey("支付宝下单", t, func() {
		resp := DefaultClient.ProcessBarcodePay(scanPay)
		Convey("期望", func() {
			So(resp.RespCode, ShouldEqual, "000000")
		})
	})

}
