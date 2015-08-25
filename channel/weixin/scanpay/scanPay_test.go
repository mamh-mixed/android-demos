package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"

	. "github.com/smartystreets/goconvey/convey"
)

func xTestProcessBarcodePay(t *testing.T) {
	m := &model.ScanPayRequest{
		AppID:      "wx25ac886b6dac7dd2", // 公众账号ID
		ChanMerId:  "1236593202",         // 商户号
		SubMchId:   "1247075201",         // 子商户
		DeviceInfo: "1000",               // 设备号
		Subject:    "被扫支付测试",             // 商品描述
		GoodsInfo:  "",                   // 商品详情
		OrderNum:   util.Millisecond(),   // 商户订单号
		ActTxamt:   "1",                  // 总金额
		ScanCodeId: "130466765198371945", // 授权码
		SignKey:    "12sdffjjguddddd2widousldadi9o0i1",
	}

	ret, err := DefaultWeixinScanPay.ProcessBarcodePay(m)

	Convey("应该不出现错误", t, func() {
		So(err, ShouldBeNil)
	})

	Convey("应该有响应信息", t, func() {
		So(ret, ShouldNotBeNil)
	})

	Convey("应答码应该是14", t, func() {
		So(ret.Respcd, ShouldEqual, "14")
	})

	m.ScanCodeId = "130502284209256489"
	ret, err = DefaultWeixinScanPay.ProcessBarcodePay(m)
	Convey("应答码应该是00", t, func() {
		So(ret.Respcd, ShouldEqual, "00")
	})
	t.Logf("%#v", ret)

}

func TestProcessEnquiry(t *testing.T) {
	m := &model.ScanPayRequest{
		AppID:      "wx25ac886b6dac7dd2", // 公众账号ID
		ChanMerId:  "1236593202",         // 商户号
		SubMchId:   "1247075201",
		DeviceInfo: "1000",               // 设备号
		Subject:    "被扫支付测试",             // 商品描述
		GoodsInfo:  "",                   // 商品详情
		OrderNum:   "1437537877995",      // 商户订单号
		Txamt:      "1",                  // 总金额
		ScanCodeId: "130512005267470788", // 授权码
		SignKey:    "12sdffjjguddddd2widousldadi9o0i1",
	}

	ret, err := DefaultWeixinScanPay.ProcessEnquiry(m)

	t.Logf("%#v", ret)

	if err != nil {
		t.Error(err)
	}
}

func xTestProcessClose(t *testing.T) {
	m := &model.ScanPayRequest{
		AppID:        "wx25ac886b6dac7dd2", // 公众账号ID
		ChanMerId:    "1236593202",         // 商户号
		SubMchId:     "1247075201",
		OrigOrderNum: "1415757673", // 商户订单号
		SignKey:      "12sdffjjguddddd2widousldadi9o0i1",
	}

	ret, err := DefaultWeixinScanPay.ProcessClose(m)

	t.Logf("%#v", ret)

	if err != nil {
		t.Error(err)
	}
}
