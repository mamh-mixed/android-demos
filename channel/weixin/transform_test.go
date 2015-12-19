package weixin

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func xTestTransform(t *testing.T) {
	returnCode, resultCode, errCode, errCodeDes := "SUCCESS", "FAIL", "OUT_TRADE_NO_USED", ""
	status, msg, _ := Transform("prePay", returnCode, resultCode, errCode, errCodeDes)
	Convey("应该返回非空的应答码", t, func() {
		So(status, ShouldNotEqual, "")
	})

	Convey("应该返回19的应答码", t, func() {
		So(status, ShouldEqual, "19")
	})

	t.Logf("response code is %s; response message is %s", status, msg)

	returnCode, resultCode, errCode, errCodeDes = "SUCCESS", "FAIL", "ORDERPAID", "商户订单已支付"
	status, msg, _ = Transform("prePay", returnCode, resultCode, errCode, errCodeDes)
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
