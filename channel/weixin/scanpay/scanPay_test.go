package scanpay

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"testing"

	"github.com/omigo/log"
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

func TestProcessSettleEnquiry(t *testing.T) {
	spReq := &model.ScanPayRequest{
		AppID:     "wxba40dfa5d37b4ad4",
		ChanMerId: "1294619701", // 商户号
		SubMchId:  "11015876",   // 子商户号
		SignKey:   "7319f3ae133ec5b065388c2eb88bd969",
		SignType:  "HMAC-SHA256",
		WeixinClientCert: []byte(`-----BEGIN CERTIFICATE-----
MIIESjCCA7OgAwIBAgIDC74OMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MTIwNDA5MDAyNloX
DTI1MTIwMTA5MDAyNlowezELMAkGA1UEBhMCQ04xEjAQBgNVBAgTCUd1YW5nZG9u
ZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1RlbmNlbnQxDjAMBgNVBAsT
BU1NUGF5MRAwDgYDVQQDEwd0ZW5jZW50MREwDwYDVQQEEwgxMDk5NDE0NzCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALsnFopfktZG/LM4aRAyfFsSaJMS
Vj+0BE2yBO2xboQR1dx74vpZCTfprasrCpuj6igedUSqbuUIsw2wr11NlaFRxFVs
A9aOvYBjHYKMyOPwJGBuUzlucmpdFdpEbw35ipvwYsuUN1IOlNdLuIJzREzh3A1X
MA9cSkzLATbLJ3vCa7uE3WqmoJ/vuQyf3RZrsgh259CsABSV5dGndR8E148K1KGw
hD0jjX6WccekOeQnynjXWevQTBaE+06BxcmfTUaues3Vd5RldOdqFCH1rfknr6X9
P5/KPu5keA8Vov0+9TCxHMkEACbm9sykZoggcINpZAw1UJq0xwe7b7jQNZ8CAwEA
AaOCAUYwggFCMAkGA1UdEwQCMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2Vu
ZXJhdGUgQ2VydGlmaWNhdGUiMB0GA1UdDgQWBBTkfI5Od5vnilA0ne+YsxgzwYLD
VDCBvwYDVR0jBIG3MIG0gBQ+BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijEL
MAkGA1UEBhMCQ04xEjAQBgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpo
ZW4xEDAOBgNVBAoTB1RlbmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1w
YXltY2hDQTEfMB0GCSqGSIb3DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8
AOhXMA4GA1UdDwEB/wQEAwIGwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkq
hkiG9w0BAQUFAAOBgQCM3w4NC+GutK1UdbAd1hgLuzE+kIaadxmvVeOLym48kQjR
VlMUFWZLYlVhBGwGpGCxV15mTtC2QnGv5GQ82M9zxLBpUMl+I58Ado0QcoqmkOJ8
WRe20KNX7okjvLrDqwJ+j6HHbdtr4qbsKBWzC1d9mp/vbuft6xzTm4X7D0vIAA==
-----END CERTIFICATE-----
`),
		WeixinClientKey: []byte(`-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC7JxaKX5LWRvyz
OGkQMnxbEmiTElY/tARNsgTtsW6EEdXce+L6WQk36a2rKwqbo+ooHnVEqm7lCLMN
sK9dTZWhUcRVbAPWjr2AYx2CjMjj8CRgblM5bnJqXRXaRG8N+Yqb8GLLlDdSDpTX
S7iCc0RM4dwNVzAPXEpMywE2yyd7wmu7hN1qpqCf77kMn90Wa7IIdufQrAAUleXR
p3UfBNePCtShsIQ9I41+lnHHpDnkJ8p411nr0EwWhPtOgcXJn01GrnrN1XeUZXTn
ahQh9a35J6+l/T+fyj7uZHgPFaL9PvUwsRzJBAAm5vbMpGaIIHCDaWQMNVCatMcH
u2+40DWfAgMBAAECggEAWBShayZn+SkXrVGTQOhB1qrnRLCQnxKeI+Lwpt2m4clz
GX3E6YYV7JayAakUKQQQJCmRj2uXHXvmqT7KieMF1RKikIFxnP04+r+rF9IiigRv
WmMSECmG98AvlLY4fh8/uPx9wspS7u+l5V3hzKNrNbPm2PB70f7hiRrRy3P5dSjR
Qx35A2Pc21Us0LPDmOqeFKLZ8hZZgXnZC3DSOXT+DPfhiWSrhSprXcNn9UDFC5G4
xBji1tiu5FxiOGS1qX2F2jhcpEVC91TX6NyLKx5qT7jIRfG8OPCCDUA+fjTxm6n6
F02Sx4ySzvz68bQANQzTn5YE8kxx+Iu7o2QXAlALIQKBgQDlKAUogxbeIgM2cjKl
HGQCsfN5BBv7zkXMtIOaJcCvdahmUFflVyfzhPLwkKcWJDqEazpcE7opR8ngrcRY
w8hK5N7Yo+DfHHnRegDIec7gENF/QWZISORV2RplMPjYCx9qjZb8Efeo5v76c28P
RQxsa935jG6ts1O7LxXgXbUODQKBgQDRE3SVdBZIBrivTqC1VOlQ294LqHjyULWq
yCJhenwJAiEB6VB6siR/pYWT7woze77kMx1h17mpUpJH1sOd+TscWJiysSCXTTHL
fB1m0RK9dMcAGTXoFClQQIdtnT7j1xo3EFSkE1CDPgx9LDtmhmXYw2/jrIzeBgqG
/WaPmf5TWwKBgHo1oKKdXqmyi/ISbw39GR0S2BzW4zGkLTdhobmonp00a029VVLa
61SEt6cyDdnSEl8ibGpEnwh635/yK/+G0h+W9X2l0DxMjG752McRpxF6BrAOKcUt
EpDyCpCi0GMvc165Cid+UY0DPEJHI5iKI3kDAcufaDs8os/16X3Rm8hVAoGAUVWH
OgjSDi0HNbOZ01D3/uu5osFkY6fANHLs5Qz4Kaz5WDwCgejBSCMFJvqL9mFCbpXv
7Ts1z8f+fhHvQYpOU4WnyYQckJ+IKofbDD9zUd6W99PW47TMKZsg1CSi7ZGlchxZ
QQb+gD/oLd3CMe1Au6Wz5lce6dRhzgfQGTcn+BcCgYBcPZBRdWMzAPrPNDeid8A1
BZjzM1Gdd8RaA0WTvPBOKYs36AG65/65/KUTRgn2+/g2Ya420qRiVdokRpw1mQh6
THBu6M11DaONH0iLLHf5LfMiBE9GIMd5De3wlZHKeZU9pPOA+E+0C2vpBMQBUF4r
kRmEuLIiZIyI3CuocnM2Jg==
-----END PRIVATE KEY-----
`),
	}
	spReq.SettDate = "20151231"
	cbd := make(model.ChanBlendMap)
	err := DefaultWeixinScanPay.ProcessSettleEnquiry(spReq, cbd)
	if err != nil {
		log.Debug("test")
		fmt.Printf("error execut \n")
		t.Error(err)
		t.FailNow()
	}
}

func TestProcessCancel(t *testing.T) {
	spReq := &model.ScanPayRequest{
		AppID:        "wxba40dfa5d37b4ad4",
		ChanMerId:    "1294619701", // 商户号
		SubMchId:     "11015876",   // 子商户号
		SignKey:      "7319f3ae133ec5b065388c2eb88bd969",
		OrigOrderNum: "1450165113707",
		Currency:     "",

		WeixinClientCert: []byte(`-----BEGIN CERTIFICATE-----
MIIESjCCA7OgAwIBAgIDC74OMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MTIwNDA5MDAyNloX
DTI1MTIwMTA5MDAyNlowezELMAkGA1UEBhMCQ04xEjAQBgNVBAgTCUd1YW5nZG9u
ZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1RlbmNlbnQxDjAMBgNVBAsT
BU1NUGF5MRAwDgYDVQQDEwd0ZW5jZW50MREwDwYDVQQEEwgxMDk5NDE0NzCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALsnFopfktZG/LM4aRAyfFsSaJMS
Vj+0BE2yBO2xboQR1dx74vpZCTfprasrCpuj6igedUSqbuUIsw2wr11NlaFRxFVs
A9aOvYBjHYKMyOPwJGBuUzlucmpdFdpEbw35ipvwYsuUN1IOlNdLuIJzREzh3A1X
MA9cSkzLATbLJ3vCa7uE3WqmoJ/vuQyf3RZrsgh259CsABSV5dGndR8E148K1KGw
hD0jjX6WccekOeQnynjXWevQTBaE+06BxcmfTUaues3Vd5RldOdqFCH1rfknr6X9
P5/KPu5keA8Vov0+9TCxHMkEACbm9sykZoggcINpZAw1UJq0xwe7b7jQNZ8CAwEA
AaOCAUYwggFCMAkGA1UdEwQCMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2Vu
ZXJhdGUgQ2VydGlmaWNhdGUiMB0GA1UdDgQWBBTkfI5Od5vnilA0ne+YsxgzwYLD
VDCBvwYDVR0jBIG3MIG0gBQ+BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijEL
MAkGA1UEBhMCQ04xEjAQBgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpo
ZW4xEDAOBgNVBAoTB1RlbmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1w
YXltY2hDQTEfMB0GCSqGSIb3DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8
AOhXMA4GA1UdDwEB/wQEAwIGwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkq
hkiG9w0BAQUFAAOBgQCM3w4NC+GutK1UdbAd1hgLuzE+kIaadxmvVeOLym48kQjR
VlMUFWZLYlVhBGwGpGCxV15mTtC2QnGv5GQ82M9zxLBpUMl+I58Ado0QcoqmkOJ8
WRe20KNX7okjvLrDqwJ+j6HHbdtr4qbsKBWzC1d9mp/vbuft6xzTm4X7D0vIAA==
-----END CERTIFICATE-----
`),
		WeixinClientKey: []byte(`-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC7JxaKX5LWRvyz
OGkQMnxbEmiTElY/tARNsgTtsW6EEdXce+L6WQk36a2rKwqbo+ooHnVEqm7lCLMN
sK9dTZWhUcRVbAPWjr2AYx2CjMjj8CRgblM5bnJqXRXaRG8N+Yqb8GLLlDdSDpTX
S7iCc0RM4dwNVzAPXEpMywE2yyd7wmu7hN1qpqCf77kMn90Wa7IIdufQrAAUleXR
p3UfBNePCtShsIQ9I41+lnHHpDnkJ8p411nr0EwWhPtOgcXJn01GrnrN1XeUZXTn
ahQh9a35J6+l/T+fyj7uZHgPFaL9PvUwsRzJBAAm5vbMpGaIIHCDaWQMNVCatMcH
u2+40DWfAgMBAAECggEAWBShayZn+SkXrVGTQOhB1qrnRLCQnxKeI+Lwpt2m4clz
GX3E6YYV7JayAakUKQQQJCmRj2uXHXvmqT7KieMF1RKikIFxnP04+r+rF9IiigRv
WmMSECmG98AvlLY4fh8/uPx9wspS7u+l5V3hzKNrNbPm2PB70f7hiRrRy3P5dSjR
Qx35A2Pc21Us0LPDmOqeFKLZ8hZZgXnZC3DSOXT+DPfhiWSrhSprXcNn9UDFC5G4
xBji1tiu5FxiOGS1qX2F2jhcpEVC91TX6NyLKx5qT7jIRfG8OPCCDUA+fjTxm6n6
F02Sx4ySzvz68bQANQzTn5YE8kxx+Iu7o2QXAlALIQKBgQDlKAUogxbeIgM2cjKl
HGQCsfN5BBv7zkXMtIOaJcCvdahmUFflVyfzhPLwkKcWJDqEazpcE7opR8ngrcRY
w8hK5N7Yo+DfHHnRegDIec7gENF/QWZISORV2RplMPjYCx9qjZb8Efeo5v76c28P
RQxsa935jG6ts1O7LxXgXbUODQKBgQDRE3SVdBZIBrivTqC1VOlQ294LqHjyULWq
yCJhenwJAiEB6VB6siR/pYWT7woze77kMx1h17mpUpJH1sOd+TscWJiysSCXTTHL
fB1m0RK9dMcAGTXoFClQQIdtnT7j1xo3EFSkE1CDPgx9LDtmhmXYw2/jrIzeBgqG
/WaPmf5TWwKBgHo1oKKdXqmyi/ISbw39GR0S2BzW4zGkLTdhobmonp00a029VVLa
61SEt6cyDdnSEl8ibGpEnwh635/yK/+G0h+W9X2l0DxMjG752McRpxF6BrAOKcUt
EpDyCpCi0GMvc165Cid+UY0DPEJHI5iKI3kDAcufaDs8os/16X3Rm8hVAoGAUVWH
OgjSDi0HNbOZ01D3/uu5osFkY6fANHLs5Qz4Kaz5WDwCgejBSCMFJvqL9mFCbpXv
7Ts1z8f+fhHvQYpOU4WnyYQckJ+IKofbDD9zUd6W99PW47TMKZsg1CSi7ZGlchxZ
QQb+gD/oLd3CMe1Au6Wz5lce6dRhzgfQGTcn+BcCgYBcPZBRdWMzAPrPNDeid8A1
BZjzM1Gdd8RaA0WTvPBOKYs36AG65/65/KUTRgn2+/g2Ya420qRiVdokRpw1mQh6
THBu6M11DaONH0iLLHf5LfMiBE9GIMd5De3wlZHKeZU9pPOA+E+0C2vpBMQBUF4r
kRmEuLIiZIyI3CuocnM2Jg==
-----END PRIVATE KEY-----
`),
	}

	resp, err := DefaultWeixinScanPay.ProcessCancel(spReq)
	if err != nil {
		t.Logf("error is %s", err)
	}
	t.Logf("response is %#v", resp)
}
