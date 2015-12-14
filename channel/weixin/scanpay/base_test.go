package scanpay

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/util"
)

var testCommonParams = weixin.CommonParams{
	Appid: "wx25ac886b6dac7dd2", // 公众账号ID
	MchID: "1236593202",         // 商户号
	// SubMchId: "1247075201",         // 子商户号
	SubMchId: "1272504101",   // 子商户号
	NonceStr: util.Nonce(32), // 随机字符串
	Sign:     "",             // 签名

	WeixinMD5Key: "12sdffjjguddddd2widousldadi9o0i1",

	// ClientCert: readPEMBlock(goconf.Config.WeixinScanPay.ClientCert),
	// ClientKey:  readPEMBlock(goconf.Config.WeixinScanPay.ClientKey),

	ClientCert: []byte(`-----BEGIN CERTIFICATE-----\nMIIEaDCCA9GgAwIBAgIDAfqVMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD\nTjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE\nChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w\nHQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MDQwMzAyMzAwNFoX\nDTI1MDMzMTAyMzAwNFowgZgxCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv\nbmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL\nEwVNTVBheTEtMCsGA1UEAxQk5LiK5rW36K6v6IGU5pWw5o2u5pyN5Yqh5pyJ6ZmQ\n5YWs5Y+4MREwDwYDVQQEEwgxMDE2NzQxMTCCASIwDQYJKoZIhvcNAQEBBQADggEP\nADCCAQoCggEBAMxA1RYF2AMSCh0RePIV9FVqtwgTMGYLx9AxMfCcq/0NcWx3RlBG\nCfixd7KrI0GT7WGBj6PN2U/yMegkyTnSdHUlfyvDpIztTzYIxAPUg7cIqB+ixaF1\n5yOFiDLsrPxVKb8viU1vUuSb0N9i/beEia8Bfq2Jk+mrZi1I7ohoSkrCAKoxFbWw\n084bZz4T1U7hUQ6abXyEgzBtL8KGriXrr5+XV+JF0BQM0w2JuE3UQxgbPglOoTWI\ntp7cdRrzL5bN2iPn02Q1EDCkNv8m2KRaDOeloaACe4jN4SW6hzQZU6z9WzHjVjak\nV0NcJXBI5mk8fN1EhqwAt2Sop55OZ1ClNAsCAwEAAaOCAUYwggFCMAkGA1UdEwQC\nMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2VydGlmaWNhdGUi\nMB0GA1UdDgQWBBR5JD3KNtUZo81jYIlCwBBvRsQy7TCBvwYDVR0jBIG3MIG0gBQ+\nBSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMCQ04xEjAQBgNV\nBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1RlbmNl\nbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEfMB0GCSqGSIb3\nDQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1UdDwEB/wQEAwIG\nwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOBgQBoedyK\nGXJ1pklDV1vgIYT+lrog8dE2U/TBxhwL65mSVT7Litgmxy2Mylm726+FoGiy4Mkx\nBoSj6A0Dfb3rtd3q4fmG8cK5eq1Uz0KTjMhlQs+WR3AYy18vQgKb3YmhGVnoPFLU\n5LF6DTYvAnRgGA2I5UrTRXYsuqXm7qKXK9E+7Q==\n-----END CERTIFICATE-----\n`),
	ClientKey:  []byte(`-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAzEDVFgXYAxIKHRF48hX0VWq3CBMwZgvH0DEx8Jyr/Q1xbHdG\nUEYJ+LF3sqsjQZPtYYGPo83ZT/Ix6CTJOdJ0dSV/K8OkjO1PNgjEA9SDtwioH6LF\noXXnI4WIMuys/FUpvy+JTW9S5JvQ32L9t4SJrwF+rYmT6atmLUjuiGhKSsIAqjEV\ntbDTzhtnPhPVTuFRDpptfISDMG0vwoauJeuvn5dX4kXQFAzTDYm4TdRDGBs+CU6h\nNYi2ntx1GvMvls3aI+fTZDUQMKQ2/ybYpFoM56WhoAJ7iM3hJbqHNBlTrP1bMeNW\nNqRXQ1wlcEjmaTx83USGrAC3ZKinnk5nUKU0CwIDAQABAoIBAQCvl8jYnvt+YELL\njJrSW+dqi0yAp6aDA/uqUrChLr9409a/raaIGj42S7MgqZmspdR8b9qhsrTw0sDu\n1rkbeX7euvaiFBZhhR4E0PJabJczgkCuuct3LBoiYoidZvSsFTbHgsFiDaNQn1eo\nw7xkyY9oITvbSpwbVVuI8NsH78h2jNqdxBmEo9JAWQiJwVG28Gvcf0+KxZydAj9q\nGlfiCVDoYeqejCVVSU1yjzvJwhNC4fK4Jy0JzpmFYBMj3BMRT/lnzvJdXABbxb/7\n2FWojkq4IRaUL1g/r0W0pErJavx5PSHgm3sQawXdHcmZ8ZLCDZZmG0mZOqAohtH+\nRNdJsTcZAoGBAOsT7Dc44hJhMvZKkjvHhixWWl6L3E6yXZ8LxhaQ0qw5vYwMQ6BU\nJjY+PzjMjskJs9hYz5MWB9bcWPHkBLHajiUqp53W4S6/5046VrcQyIGxwlea0s8w\nufNrcQUXVT9EnYTvBvfW1DheUJiyn3rnr7hKyXVxrX5PY+DT64REHCLFAoGBAN5u\nmB0xTpXuiR5BdtVa8de7ChCOqTwXYIKtNiX6Q30UxKZJ7OwZoyASDKRsTXSeGW/r\n6Tb3jUEa/wcdngonacrid4BdGmKaMwHaKaZowxHMn5j3g6901mUu/AtV0q3nbJq5\nuCuFXeOnDA4vSpaml/UDvWvRqJbiPzLvqyhNcyiPAoGBAI7bfZqVi/VlckXwTWvc\ntfItzB9W2VxN0s07p3bBLfYR5Nm9/j7pxIsESwFmdoM/zTZ1yjd1lPAC2l6tlhjL\nW8TEZjZqhlAVuSh2FYqMvXzrnNIGOYRF9UsziOxyIJEhTqShadelizRyRIJ3Uqmr\nMMNLV6Byo991uZnAz4iCp6KNAoGBAJyTM0bRb6VBHYqLwI/djgIzKpmPIvgm6Iv0\nS/qd2aYR2X/I6BsmzNqFehrAFiHyLKvJYAiOaAOdckpbAeXZ6rGji0VzxGAGdcNn\nBAydEDvWU75E9ZCr6UOeuFNuXXiHQL8F3uvb3MSk0WqmxZWYvbz+nfdoxYk4yA4e\nAdjD9D1nAoGAVDUAqgbUMYp7toDTgamj0Qlu+Uz1QK8QXN7O0QN+iDi9Gr1BSUAh\ngkztyH0g0nY1tD7WPAwM/kNupd0SI9z9q5K36wGrNlETF7wdOHgWl3sQlTDepNzr\nlCxbgNyTgFhC6N+1YhY1QKOvjjPS+tHAhMR0FnH9gUW1RtBQLUxpfFM=\n-----END RSA PRIVATE KEY-----\n`),
}

func readPEMBlock(file string) (certPEMBlock []byte) {
	var err error
	certPEMBlock, err = ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("read cert file err: %s\n", err)
	}
	return certPEMBlock
}

func TestReadPEMBlock(t *testing.T) {
	t.Errorf("%s", testCommonParams.ClientCert)
	t.Errorf("%s", testCommonParams.ClientKey)
}
