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
	spReq := &model.ScanPayRequest{}

	spReq.AppID = "wx25ac886b6dac7dd2" // 公众账号ID
	spReq.ChanMerId = "1236593202"     // 商户号
	// SubMchId: "1247075201",        // 子商户号
	spReq.SubMchId = "1244891002" // 子商户号
	spReq.SignKey = "12sdffjjguddddd2widousldadi9o0i1"

	// ClientCert: readPEMBlock(goconf.Config.WeixinScanPay.ClientCert),
	// ClientKey:  readPEMBlock(goconf.Config.WeixinScanPay.ClientKey),
	/*
		spReq.WeixinClientCert = []byte(`-----BEGIN CERTIFICATE-----
			MIIEaDCCA9GgAwIBAgIDAfqVMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
			TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
			ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
			HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MDQwMzAyMzAwNFoX
			DTI1MDMzMTAyMzAwNFowgZgxCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv
			bmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL
			EwVNTVBheTEtMCsGA1UEAxQk5LiK5rW36K6v6IGU5pWw5o2u5pyN5Yqh5pyJ6ZmQ
			5YWs5Y+4MREwDwYDVQQEEwgxMDE2NzQxMTCCASIwDQYJKoZIhvcNAQEBBQADggEP
			ADCCAQoCggEBAMxA1RYF2AMSCh0RePIV9FVqtwgTMGYLx9AxMfCcq/0NcWx3RlBG
			Cfixd7KrI0GT7WGBj6PN2U/yMegkyTnSdHUlfyvDpIztTzYIxAPUg7cIqB+ixaF1
			5yOFiDLsrPxVKb8viU1vUuSb0N9i/beEia8Bfq2Jk+mrZi1I7ohoSkrCAKoxFbWw
			084bZz4T1U7hUQ6abXyEgzBtL8KGriXrr5+XV+JF0BQM0w2JuE3UQxgbPglOoTWI
			tp7cdRrzL5bN2iPn02Q1EDCkNv8m2KRaDOeloaACe4jN4SW6hzQZU6z9WzHjVjak
			V0NcJXBI5mk8fN1EhqwAt2Sop55OZ1ClNAsCAwEAAaOCAUYwggFCMAkGA1UdEwQC
			MAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2VydGlmaWNhdGUi
			MB0GA1UdDgQWBBR5JD3KNtUZo81jYIlCwBBvRsQy7TCBvwYDVR0jBIG3MIG0gBQ+
			BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMCQ04xEjAQBgNV
			BAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1RlbmNl
			bnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEfMB0GCSqGSIb3
			DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1UdDwEB/wQEAwIG
			wDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOBgQBoedyK
			GXJ1pklDV1vgIYT+lrog8dE2U/TBxhwL65mSVT7Litgmxy2Mylm726+FoGiy4Mkx
			BoSj6A0Dfb3rtd3q4fmG8cK5eq1Uz0KTjMhlQs+WR3AYy18vQgKb3YmhGVnoPFLU
			5LF6DTYvAnRgGA2I5UrTRXYsuqXm7qKXK9E+7Q==
			-----END CERTIFICATE-----`)
		spReq.WeixinClientKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
			MIIEpQIBAAKCAQEAzEDVFgXYAxIKHRF48hX0VWq3CBMwZgvH0DEx8Jyr/Q1xbHdG
			UEYJ+LF3sqsjQZPtYYGPo83ZT/Ix6CTJOdJ0dSV/K8OkjO1PNgjEA9SDtwioH6LF
			oXXnI4WIMuys/FUpvy+JTW9S5JvQ32L9t4SJrwF+rYmT6atmLUjuiGhKSsIAqjEV
			tbDTzhtnPhPVTuFRDpptfISDMG0vwoauJeuvn5dX4kXQFAzTDYm4TdRDGBs+CU6h
			NYi2ntx1GvMvls3aI+fTZDUQMKQ2/ybYpFoM56WhoAJ7iM3hJbqHNBlTrP1bMeNW
			NqRXQ1wlcEjmaTx83USGrAC3ZKinnk5nUKU0CwIDAQABAoIBAQCvl8jYnvt+YELL
			jJrSW+dqi0yAp6aDA/uqUrChLr9409a/raaIGj42S7MgqZmspdR8b9qhsrTw0sDu
			1rkbeX7euvaiFBZhhR4E0PJabJczgkCuuct3LBoiYoidZvSsFTbHgsFiDaNQn1eo
			w7xkyY9oITvbSpwbVVuI8NsH78h2jNqdxBmEo9JAWQiJwVG28Gvcf0+KxZydAj9q
			GlfiCVDoYeqejCVVSU1yjzvJwhNC4fK4Jy0JzpmFYBMj3BMRT/lnzvJdXABbxb/7
			2FWojkq4IRaUL1g/r0W0pErJavx5PSHgm3sQawXdHcmZ8ZLCDZZmG0mZOqAohtH+
			RNdJsTcZAoGBAOsT7Dc44hJhMvZKkjvHhixWWl6L3E6yXZ8LxhaQ0qw5vYwMQ6BU
			JjY+PzjMjskJs9hYz5MWB9bcWPHkBLHajiUqp53W4S6/5046VrcQyIGxwlea0s8w
			ufNrcQUXVT9EnYTvBvfW1DheUJiyn3rnr7hKyXVxrX5PY+DT64REHCLFAoGBAN5u
			mB0xTpXuiR5BdtVa8de7ChCOqTwXYIKtNiX6Q30UxKZJ7OwZoyASDKRsTXSeGW/r
			6Tb3jUEa/wcdngonacrid4BdGmKaMwHaKaZowxHMn5j3g6901mUu/AtV0q3nbJq5
			uCuFXeOnDA4vSpaml/UDvWvRqJbiPzLvqyhNcyiPAoGBAI7bfZqVi/VlckXwTWvc
			tfItzB9W2VxN0s07p3bBLfYR5Nm9/j7pxIsESwFmdoM/zTZ1yjd1lPAC2l6tlhjL
			W8TEZjZqhlAVuSh2FYqMvXzrnNIGOYRF9UsziOxyIJEhTqShadelizRyRIJ3Uqmr
			MMNLV6Byo991uZnAz4iCp6KNAoGBAJyTM0bRb6VBHYqLwI/djgIzKpmPIvgm6Iv0
			S/qd2aYR2X/I6BsmzNqFehrAFiHyLKvJYAiOaAOdckpbAeXZ6rGji0VzxGAGdcNn
			BAydEDvWU75E9ZCr6UOeuFNuXXiHQL8F3uvb3MSk0WqmxZWYvbz+nfdoxYk4yA4e
			AdjD9D1nAoGAVDUAqgbUMYp7toDTgamj0Qlu+Uz1QK8QXN7O0QN+iDi9Gr1BSUAh
			gkztyH0g0nY1tD7WPAwM/kNupd0SI9z9q5K36wGrNlETF7wdOHgWl3sQlTDepNzr
			lCxbgNyTgFhC6N+1YhY1QKOvjjPS+tHAhMR0FnH9gUW1RtBQLUxpfFM=
			-----END RSA PRIVATE KEY-----`)*/
	spReq.SettDate = "20151015"
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
