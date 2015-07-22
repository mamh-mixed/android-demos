package scanpay

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransformX(t *testing.T) {
	resp := &PayResp{
		ReturnCode: "SUCCESS",
		ResultCode: "FAIL",
		ErrCode:    "OUT_TRADE_NO_USED",
	}
	status, msg := transformX("prePay", resp)
	Convey("应该返回非空的应答码", t, func() {
		So(status, ShouldNotEqual, "")
	})

	Convey("应该返回19的应答码", t, func() {
		So(status, ShouldEqual, "19")
	})

	t.Logf("response code is %s; response message is %s", status, msg)

	resp = &PayResp{
		ReturnCode: "SUCCESS",
		ResultCode: "FAIL",
		ErrCode:    "ORDERPAID",
		ErrCodeDes: "商户订单已支付",
	}
	status, msg = transformX("prePay", resp)
	Convey("应该返回非空的应答码", t, func() {
		So(status, ShouldNotEqual, "")
	})

	Convey("应该返回01的应答码", t, func() {
		So(status, ShouldEqual, "01")
	})

	Convey("应该返回非空的应答信息", t, func() {
		So(status, ShouldNotEqual, "")
	})

	t.Logf("response code is %s; response message is %s", status, msg)

}
