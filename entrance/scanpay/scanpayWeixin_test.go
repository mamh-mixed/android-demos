package scanpay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"

	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 混淆的交易类型：报文为下单，类型为预下单
func TestWXPayConfusionBusicd(t *testing.T) {
	orderNum := util.Millisecond()
	req := &model.ScanPayRequest{
		Chcd:       "WXP",
		GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
		OrderNum:   orderNum,
		ScanCodeId: "131008245544067366",
		Inscd:      "CIL00002",
		Txamt:      "000000000001",
		Busicd:     "PAUT",
		Mchntid:    "100000000000203",
		Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}

	resp := doSendWXPayRequest(req)

	Convey("混淆的交易类型", t, func() {
		So(resp.Respcd, ShouldEqual, "30")
	})
}

// 不存在的交易类型
func TestWXPayWrongBusicd(t *testing.T) {
	orderNum := util.Millisecond()
	req := &model.ScanPayRequest{
		Chcd:       "WXP",
		GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
		OrderNum:   orderNum,
		ScanCodeId: "",
		Inscd:      "CIL00002",
		Txamt:      "000000000001",
		Busicd:     "ANSD",
		Mchntid:    "100000000000203",
		Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}

	resp := doSendWXPayRequest(req)

	Convey("不存在的交易类型", t, func() {
		So(resp.Respcd, ShouldEqual, "30")
	})
}

// 下单取消 成功
func TestWXPayCloseSuccess(t *testing.T) {
	orderNum := util.Millisecond()
	req := &model.ScanPayRequest{
		Chcd:       "WXP",
		GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
		OrderNum:   orderNum,
		ScanCodeId: "130939087138632762",
		Inscd:      "CIL00002",
		Txamt:      "000000000001",
		Busicd:     "PURC",
		Mchntid:    "100000000000203",
		Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}

	resp := doSendWXPayRequest(req)

	Convey("成功下单", t, func() {
		So(resp.Respcd, ShouldEqual, "00")
	})

	req = &model.ScanPayRequest{
		Busicd:       "CANC",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: orderNum,
		Inscd:        "CIL00002",
	}

	resp = doSendWXPayRequest(req)
	Convey("下单取消", t, func() {
		So(resp.Respcd, ShouldEqual, "00")
	})
}

// 预下单取消 成功
func TestWXPrePayCloseSuccess(t *testing.T) {
	req := &model.ScanPayRequest{
		Busicd:       "CANC",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1437718003292",
		Inscd:        "CIL00002",
	}

	resp := doSendWXPayRequest(req)
	Convey("预下单取消", t, func() {
		So(resp.Respcd, ShouldEqual, "00")
	})
}

// 已撤销订单取消 成功
func TestWXReverseCloseSuccess(t *testing.T) {
	req := &model.ScanPayRequest{
		Busicd:       "CANC",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1437638534585",
		Inscd:        "CIL00002",
	}

	resp := doSendWXPayRequest(req)
	Convey("预下单取消", t, func() {
		So(resp.Respcd, ShouldEqual, "00")
	})
}

// 已退款订单取消 成功
func TestWXRefundCloseSuccess(t *testing.T) {
	req := &model.ScanPayRequest{
		Busicd:       "CANC",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1437641588899",
		Inscd:        "CIL00002",
	}

	resp := doSendWXPayRequest(req)
	Convey("预下单取消", t, func() {
		So(resp.Respcd, ShouldEqual, "00")
	})
}

// 支付宝有效支付条码送微信
func TestWXPayUseAlPQrcode(t *testing.T) {
	orderNum := util.Millisecond()
	req := &model.ScanPayRequest{
		Chcd:       "WXP",
		GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
		OrderNum:   orderNum,
		ScanCodeId: "289948455182572311",
		Inscd:      "CIL00002",
		Txamt:      "000000000001",
		Busicd:     "PURC",
		Mchntid:    "100000000000203",
		Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}

	resp := doSendWXPayRequest(req)

	Convey("支付宝有效支付条码送微信", t, func() {
		So(resp.Respcd, ShouldEqual, "01")
	})
}

func TestWXScanPay(t *testing.T) {
	// 下单支付
	orderNum := util.Millisecond()
	req := &model.ScanPayRequest{
		Chcd:       "WXP",
		GoodsInfo:  "鞋子,1000,2;衣服,1500,3",
		OrderNum:   orderNum,
		ScanCodeId: "130322701451895248",
		Inscd:      "CIL00002",
		Txamt:      "000000000001",
		Busicd:     "PURC",
		Mchntid:    "100000000000203",
		Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}

	reqBytes, _ := json.Marshal(req)

	respBytes := ScanPayHandle(reqBytes)

	var resp = model.ScanPayResponse{}

	_ = json.Unmarshal(respBytes, &resp)

	t.Logf("response is %#v", resp)

	Convey("微信下单支付", t, func() {
		So(resp.Respcd, ShouldEqual, "00")
	})

	time.Sleep(2 * time.Second)

	// 条码已使用
	req.OrderNum = util.Millisecond()
	reqBytes, _ = json.Marshal(req)
	respBytes = ScanPayHandle(reqBytes)
	_ = json.Unmarshal(respBytes, &resp)
	Convey("微信下单支付条码已使用", t, func() {
		So(resp.Respcd, ShouldEqual, "12")
	})

	time.Sleep(5 * time.Second)
	// 撤销交易
	req = &model.ScanPayRequest{
		Busicd:       "VOID",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: orderNum,
		Inscd:        "CIL00002",
		Chcd:         "WXP",
	}
	t.Logf("撤销：%#v", req)

	reqBytes, _ = json.Marshal(req)
	respBytes = ScanPayHandle(reqBytes)

	_ = json.Unmarshal(respBytes, &resp)
	Convey("微信下单支付撤销", t, func() {
		So(resp.Respcd, ShouldEqual, "00")
	})
}

func doSendWXPayRequest(req *model.ScanPayRequest) (resp *model.ScanPayResponse) {
	reqBytes, _ := json.Marshal(req)

	respBytes := ScanPayHandle(reqBytes)

	resp = &model.ScanPayResponse{}

	_ = json.Unmarshal(respBytes, resp)

	return resp
}
