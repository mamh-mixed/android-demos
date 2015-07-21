package scanpay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
)

func TestProcessBarcodePay(t *testing.T) {
	m := &model.ScanPayRequest{
		AppID:      "wx25ac886b6dac7dd2", // 公众账号ID
		ChanMerId:  "1236593202",         // 商户号
		SubMchId:   "1247075201",         // 子商户
		DeviceInfo: "1000",               // 设备号
		Subject:    "被扫支付测试",             // 商品描述
		GoodsInfo:  "",                   // 商品详情
		Attach:     "订单额外描述",             // 附加数据
		OrderNum:   util.Millisecond(),   // 商户订单号
		ActTxamt:   "1",                  // 总金额
		CurrType:   "CNY",                // 货币类型
		GoodsGag:   "",                   // 商品标记
		ScanCodeId: "130466765198371945", // 授权码
		SignCert:   "12sdffjjguddddd2widousldadi9o0i1",
	}

	ret, err := DefaultWeixinScanPay.ProcessBarcodePay(m)

	t.Logf("%#v", ret)

	if err != nil {
		t.Error(err)
	}
}

func TestProcessEnquiry(t *testing.T) {
	m := &model.ScanPayRequest{
		AppID:      "wx25ac886b6dac7dd2", // 公众账号ID
		ChanMerId:  "1236593202",         // 商户号
		SubMchId:   "1247075201",
		DeviceInfo: "1000",               // 设备号
		Subject:    "被扫支付测试",             // 商品描述
		GoodsInfo:  "",                   // 商品详情
		Attach:     "订单额外描述",             // 附加数据
		OrderNum:   "1415757673",         // 商户订单号
		Txamt:      "1",                  // 总金额
		CurrType:   "CNY",                // 货币类型
		GoodsGag:   "",                   // 商品标记
		ScanCodeId: "130512005267470788", // 授权码
		SignCert:   "12sdffjjguddddd2widousldadi9o0i1",
	}

	ret, err := DefaultWeixinScanPay.ProcessEnquiry(m)

	t.Logf("%#v", ret)

	if err != nil {
		t.Error(err)
	}
}

func TestProcessClose(t *testing.T) {
	m := &model.ScanPayRequest{
		AppID:        "wx25ac886b6dac7dd2", // 公众账号ID
		ChanMerId:    "1236593202",         // 商户号
		SubMchId:     "1247075201",
		OrigOrderNum: "1415757673", // 商户订单号
		CurrType:     "CNY",        // 货币类型
		SignCert:     "12sdffjjguddddd2widousldadi9o0i1",
	}

	ret, err := DefaultWeixinScanPay.ProcessClose(m)

	t.Logf("%#v", ret)

	if err != nil {
		t.Error(err)
	}
}
