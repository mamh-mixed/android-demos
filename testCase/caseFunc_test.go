// quickpay 测试用例脚本
package testCase

import (
	"github.com/CardInfoLink/quickpay/channel/cil"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	testMerId      = "012345678901234"
	testSign       = "0123456789"
	testEncryptKey = "AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8="
)

var (
	acctNum         = ""
	amt       int64 = 100000
	bindingId       = ""
	orderNum        = ""
)

func init() {
	// 日志输出级别
	log.SetOutputLevel(log.Lerror)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 连接到 MongoDB
	mongo.Connect()

	// 初始化卡 Bin 树
	core.BuildTree()

	// 连接线下
	cil.Connect()
}

// Test1 1.1.1使用一张有效银联卡发起建立绑定关系请求
func Test1(t *testing.T) {
	acctNum = "6222022003008481261"
	bindingId = tools.Millisecond()
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + testMerId
	Convey("1.1.1使用一张有效银联卡发起建立绑定关系请求", t, func() {
		b, _ := bindingCreate()
		Convey("期望结果：请求处理成功", func() {
			ret, _ := post(url, b)
			So(ret.RespCode, ShouldEqual, "000000")
		})
	})
}

// Test2 1.1.2使用该商户下已存在的绑定ID进行建立绑定关系
func Test2(t *testing.T) {
	bindingId = "1430128629966"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + testMerId
	Convey("1.1.2使用该商户下已存在的绑定ID进行建立绑定关系", t, func() {
		b, _ := bindingCreate()
		Convey("期望结果：绑定ID重复", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "绑定ID重复")
		})
	})
}

// Test3 1.1.3使用一张外卡进行绑定
func Test3(t *testing.T) {
	acctNum = "5311622289073236" // 外卡
	bindingId = tools.Millisecond()
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + testMerId
	Convey("1.1.3使用一张外卡进行绑定", t, func() {
		b, _ := bindingCreate()
		Convey("期望结果：无此交易权限", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "无此交易权限")
		})
	})
}

// Test4 1.1.4使用一个不存在的卡号进行绑定
func Test4(t *testing.T) {
	acctNum = "4931236819413" // 不存在的卡号
	bindingId = tools.Millisecond()
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + testMerId
	Convey("1.1.4使用一个不存在的卡号进行绑定", t, func() {
		b, _ := bindingCreate()
		Convey("期望结果：账户号码有误", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "账户号码有误")
		})
	})
}

// Test5 1.1.5使用不存在商户号发起建立绑定关系请求
func Test5(t *testing.T) {
	acctNum = "6222022003008481261"
	bindingId = tools.Millisecond()
	merId := "111123333"
	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + merId
	Convey("1.1.5使用不存在商户号发起建立绑定关系请求", t, func() {
		b, _ := bindingCreate()
		Convey("期望结果：商户号不存在", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "商户号不存在")
		})
	})
}

// Test6
// func Test6(t *testing.T) {
// 	// acctNum = "6222022003008481261"
// 	bindingId = tools.Millisecond()
// 	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + testMerId
// 	Convey("1.1.7使用错误签名上送建立绑定关系请求", t, func() {
// 		b, err := bindingCreate()
// 		So(err, ShouldBeNil)
// 		Convey("期望结果", func() {
// 			ret, err := post(url, b)
// 			So(err, ShouldBeNil)
// 			So(ret.RespCode, ShouldEqual, "200050")
// 		})
// 	})
// }

// Test7 1.2.1使用1.1.1中建立的绑定关系进行绑定支付
func Test7(t *testing.T) {
	Test1(t)
	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + testMerId
	Convey("1.2.1使用1.1.1中建立的绑定关系进行绑定支付", t, func() {
		amt = 10000
		orderNum = tools.Millisecond()
		b := BindingPayment()
		Convey("期望结果：请求处理成功", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "请求处理成功")
		})
	})
}

// Test8 1.2.2使用已解除的绑定ID进行支付
func Test8(t *testing.T) {
	Test1(t)
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + testMerId
	Convey("1.2.2使用已解除的绑定ID进行支付", t, func() {

		Convey("解除绑定", func() {
			b := BindingRemove()
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "请求处理成功")
			Convey("期望结果：绑定ID已失效", func() {
				orderNum = tools.Millisecond()
				c := BindingPayment()
				url = "https://api.xxxx.com/quickpay/bindingPayment?merId=" + testMerId
				ret, _ := post(url, c)
				So(ret.RespCode, ShouldEqual, "200074")
			})
		})
	})
}

// Test9 1.2.4使用已使用的订单号进行绑定支付
func Test9(t *testing.T) {
	Test7(t)
	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + testMerId
	Convey("1.2.4使用已使用的订单号进行绑定支付", t, func() {

		b := BindingPayment()
		Convey("期望结果：订单号重复", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "订单号重复")
		})
	})
}

// Test10
// func Test10(t *testing.T) {
// 	Test1(t)
// 	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + testMerId
// 	Convey("1.2.8进行一笔大额绑定支付（可输入最大额）", t, func() {

// 		orderNum = tools.Millisecond()
// 		amt = 100000000000000
// 		b := BindingPayment()
// 		Convey("预期结果", func() {
// 			ret, _ := post(url, b)
// 			So(ret.RespMsg, ShouldEqual, "金额有误")
// 		})
// 	})
// }

// Test11 1.2.10进行一笔0，00元的绑定支付
func Test11(t *testing.T) {
	Test1(t)
	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + testMerId
	Convey("1.2.10进行一笔0，00元的绑定支付", t, func() {

		orderNum = tools.Millisecond()
		amt = 0
		b := BindingPayment()
		Convey("期望结果：字段xx不能位空", func() {
			ret, _ := post(url, b)
			So(ret.RespCode, ShouldEqual, "200050")
		})
	})
}

// Test12 1.3.1对1.2.1已支付交易进行退款
func Test12(t *testing.T) {
	Test7(t)
	url := "https://api.xxxx.com/quickpay/refund?merId=" + testMerId
	Convey("1.3.1对1.2.1已支付交易进行退款", t, func() {
		b := BindingRefund()
		// 后续的退款查询
		Convey("期望结果：请求处理成功", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "请求处理成功")
		})
	})
}

// Test13 退一笔原交易不存在的订单
func Test13(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/refund?merId=" + testMerId
	Convey("1.3.2退一笔原交易不存在的订单", t, func() {
		orderNum = tools.Millisecond()
		b := BindingRefund()
		Convey("期望结果：原交易不成功，不能退款", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "原交易不成功，不能退款")
		})
	})
}

// Test14 退一笔原交易已退款的订单
func Test14(t *testing.T) {
	Test12(t)
	url := "https://api.xxxx.com/quickpay/refund?merId=" + testMerId
	Convey("1.3.2退一笔原交易已退款的订单", t, func() {
		b := BindingRefund()
		Convey("期望结果：该笔订单已经存在退款交易，不能再次退款", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "该笔订单已经存在退款交易，不能再次退款")
		})
	})
}

// Test15 做一笔退款金额大于原订单金额的退款
func Test15(t *testing.T) {
	Test7(t)
	url := "https://api.xxxx.com/quickpay/refund?merId=" + testMerId
	Convey("1.3.4做一笔退款金额大于原订单金额的退款", t, func() {
		amt = amt * 2
		b := BindingRefund()
		Convey("期望结果：退款金额（累计）大于可退金额", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "退款金额（累计）大于可退金额")
		})
	})
}

// Test16 做一笔退款金额小于原订单订单金额的退款
func Test16(t *testing.T) {
	Test7(t)
	url := "https://api.xxxx.com/quickpay/refund?merId=" + testMerId
	Convey("1.3.5做一笔退款金额小于原订单订单金额的退款", t, func() {
		amt = amt / 2
		b := BindingRefund()
		Convey("期望结果：退款金额有误", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "退款金额有误")
		})
	})
}

// Test17 查询1.2.1绑定支付结果
func Test17(t *testing.T) {
	Test7(t)
	url := "https://api.xxxx.com/quickpay/orderEnquiry?merId=" + testMerId
	Convey("1.4.1查询1.2.1绑定支付结果", t, func() {
		b := OrderEnquiry()
		Convey("期望结果：交易成功", func() {
			ret, _ := post(url, b)
			So(ret.TransStatus, ShouldEqual, "30")
		})
	})
}

// Test18 查询1.3.1退款结果
func Test18(t *testing.T) {
	Test12(t)
	url := "https://api.xxxx.com/quickpay/orderEnquiry?merId=" + testMerId
	Convey("1.4.2查询1.3.1退款结果", t, func() {
		b := OrderEnquiry()
		Convey("期望结果：请求处理成功", func() {
			ret, _ := post(url, b)
			So(ret.RespCode, ShouldEqual, "000000")
		})
	})
}

// Test20 查询一笔不存在的退款交易订单
func Test20(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/orderEnquiry?merId=" + testMerId
	Convey("1.4.6查询一笔不存在的退款交易订单", t, func() {
		orderNum = tools.Millisecond()
		b := OrderEnquiry()
		Convey("期望结果：订单号有误", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "订单号有误")
		})
	})
}

// Test21 查询一笔返回已退款成功的退款订单
func Test21(t *testing.T) {
	Test18(t)
}

// Test22 查询1.1.1绑定关系
func Test22(t *testing.T) {
	Test1(t)
	url := "https://api.xxxx.com/quickpay/bindingEnquiry?merId=" + testMerId
	Convey("1.5.1查询1.1.1绑定关系", t, func() {
		b := BindingEnquiry()
		Convey("期望结果：绑定成功", func() {
			ret, _ := post(url, b)
			So(ret.BindingStatus, ShouldEqual, "30")
		})
	})
}

// Test23 查询任意一个该商户下不存在的绑定关系
func Test23(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/bindingEnquiry?merId=" + testMerId
	Convey("1.5.2查询任意一个该商户下不存在的绑定关系", t, func() {
		bindingId = tools.Millisecond()
		b := BindingEnquiry()
		Convey("期望结果：绑定ID有误", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "绑定ID有误")
		})
	})
}

// Test24 查询已解除绑定的绑定ID
func Test24(t *testing.T) {
	Test1(t)
	Convey("1.5.5查询已解除绑定的绑定ID", t, func() {
		Convey("解除绑定", func() {
			url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + testMerId
			b := BindingRemove()
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "请求处理成功")
			Convey("期望结果：已解绑", func() {
				url = "https://api.xxxx.com/quickpay/bindingEnquiry?merId=" + testMerId
				b := BindingEnquiry()
				ret, _ := post(url, b)
				So(ret.BindingStatus, ShouldEqual, "40")
			})
		})
	})
}

// Test25 上送1.1.1中绑定ID的解除绑定关系请求
func Test25(t *testing.T) {
	Test1(t)
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + testMerId
	Convey("1.6.1上送1.1.1中绑定ID的解除绑定关系请求", t, func() {
		b := BindingRemove()
		Convey("期望结果：请求处理成功", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "请求处理成功")
		})
	})
}

// Test26 对已解除绑定关系的绑定ID再次解绑
func Test26(t *testing.T) {
	Test25(t)
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + testMerId
	Convey("1.6.2对已解除绑定关系的绑定ID再次解绑", t, func() {
		b := BindingRemove()
		Convey("期望结果：该绑定ID的已经解绑过", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "该绑定ID的已经解绑过")
		})
	})
}

// Test27 对一个不存在的绑定ID进行解绑
func Test27(t *testing.T) {
	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + testMerId
	Convey("1.6.3对一个不存在的绑定ID进行解绑", t, func() {
		bindingId = tools.Millisecond()
		b := BindingRemove()
		Convey("期望结果：绑定ID有误", func() {
			ret, _ := post(url, b)
			So(ret.RespMsg, ShouldEqual, "绑定ID有误")
		})
	})
}
