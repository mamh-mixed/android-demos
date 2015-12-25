package scanpay

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/util"
)

var testCommonParams = weixin.CommonParams{
	Appid: "wx8854422b20240ed2", // 公众账号ID
	MchID: "1276970101",         // 商户号
	// SubMchId: "1247075201",         // 子商户号
	SubMchId: "1295117001",   // 子商户号
	NonceStr: util.Nonce(32), // 随机字符串
	Sign:     "",             // 签名

	WeixinMD5Key: "tgsdhaysysdglzjjhdgbyyyhdgsbxhsh",

	// ClientCert: readPEMBlock(goconf.Config.WeixinScanPay.ClientCert),
	// ClientKey:  readPEMBlock(goconf.Config.WeixinScanPay.ClientKey),

	ClientCert: []byte(`-----BEGIN CERTIFICATE-----\nMIIEazCCA9SgAwIBAgIDB2/iMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD\nTjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE\nChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w\nHQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MTAxMjExMTEwMloX\nDTI1MTAwOTExMTEwMlowgZsxCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv\nbmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL\nEwVNTVBheTEwMC4GA1UEAxQn5LiK5rW35ou86LCx5ZWm572R57uc56eR5oqA5pyJ\n6ZmQ5YWs5Y+4MREwDwYDVQQEEwgxMDY4MjE1NzCCASIwDQYJKoZIhvcNAQEBBQAD\nggEPADCCAQoCggEBAOHjs3QHldv6ToS43rE3AYgyVOmb5GP850E9Z5q9RqlLYqgR\ndEDwtpmPhDkbtPZnJ8cPJ2ZDdnLx5CP3vIvQYg1bqPqFOCyIud72gTBUdGPuY+vY\ngavVlt8us4k5ioNa2f8UgtmMNQtvN3ohep7yLyXx1MkfEpl+IBYb4Nrj/4PbhgVt\npeTNAQMl/7EXKG8JSYAQR4TmSOVuB4Xv6ZOJgHPwxo9MWqGBV0tLc8Sdl6DBAeGP\n1AdbYq3UBjgxEdOb83ro4Z56+O37V+lNC9CKMqTbf8GmE4uI7RT1usdRX3gu5u2G\n/jVvBJJkJrjAdP3I+PZsJN/PiQOhjxJNReTGSx0CAwEAAaOCAUYwggFCMAkGA1Ud\nEwQCMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2VydGlmaWNh\ndGUiMB0GA1UdDgQWBBRwNiH4y7qrA/wkrIz5S1kDQxNqcDCBvwYDVR0jBIG3MIG0\ngBQ+BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMCQ04xEjAQ\nBgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1Rl\nbmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEfMB0GCSqG\nSIb3DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1UdDwEB/wQE\nAwIGwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOBgQBl\nFvqlyYRhRqWk2aiPZ47WR/QGv0VGDT46Wk8LrjpwjbkkbJNiLJ5ckK3tf7/Nx1F0\nTGnxZtqLyygSg/2ElhNzN9ejXC4D6c+QxhfdRY/j+X2b1VH75PQQKPh42SIInbO1\nRKkyCQG8sw0R9jciTjKl3f5NBzBqfeKt4366Iin8Tw==\n-----END CERTIFICATE-----\n`),
	ClientKey:  []byte(`-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDh47N0B5Xb+k6E\nuN6xNwGIMlTpm+Rj/OdBPWeavUapS2KoEXRA8LaZj4Q5G7T2ZyfHDydmQ3Zy8eQj\n97yL0GINW6j6hTgsiLne9oEwVHRj7mPr2IGr1ZbfLrOJOYqDWtn/FILZjDULbzd6\nIXqe8i8l8dTJHxKZfiAWG+Da4/+D24YFbaXkzQEDJf+xFyhvCUmAEEeE5kjlbgeF\n7+mTiYBz8MaPTFqhgVdLS3PEnZegwQHhj9QHW2Kt1AY4MRHTm/N66OGeevjt+1fp\nTQvQijKk23/BphOLiO0U9brHUV94Lubthv41bwSSZCa4wHT9yPj2bCTfz4kDoY8S\nTUXkxksdAgMBAAECggEBAMJ/GawFP/6ZxnO27mAuWY5YsA45YWzKfKAK7CMraCUq\nuLa32J515PPRw+qcNbOX3IMkRCtkWR/dsS9bByhnc5XG33ddr6GA1HHrVA82GMVW\npQiUcgpvrSlb/9BfECnL1zowAf6pH59J0r3BB+DF7NzCHhJSQ+SP2bbPqEsw13hC\ntUO1rsepyzHzG0jio4zDk1TweDF7T8QgS+HKMZHiGmRLPrnVxfqbp13Dz2mrWsSD\nYUT6dm403w+kmwdCwJmxsMiqIa/kT8BtU0YLKygJrLMQHrp/f6I0Bfa3N2oubU8y\nJU9kJlu8yASjevNgCssbKFH+MQ4XxcJ8fCUUeVCEBi0CgYEA8YbwXyjonBl7Eo5g\nHH4G5treSDvmIWiTakTupUTUqGu9Ocl8IQFi5SiglgyURXnsj4P0DotZMArTu+F7\n8fr7I22ThFWmPSxm8BdtqNMvXdCFprTdAadUe4FyMs4tKNBo1tbLFzizqOJztwNb\ncy743o4S8dI6hGC5yjJmVZLnlEcCgYEA72zg062DSY2hfAN/nsA9T8+sT1SZPRkz\neoZ5hfrhFZMjf8GyRAfmzJJ2wvbFRfnQ+ifgmG/Ntjh7RaTTVIh6OLmkKjJn/a8u\n4pzALGdJc9rL3/qrYZx49v8YfgSXvgI98luH07KjZYQuQ5cOzzpRlkRoFYNzQcDV\nKnK84vNFC3sCgYAznfkE/UMpCTEKOC9GJ5DmCWRz34lBHo5SqcSuwVUJYW0hSnQi\nwZ8XBmW7a5jMeFAcI8Em0pUO9WFmx7urbU36tlJOd9d9P14IdZlT+T4oOIY3qHOL\nBO3DL1jujq2MCW4+a80fe1i6ARtlw2vp4+H//jECSUGERP+vvLGuHCUtxQKBgQDh\n0gqX/H7owAoAgvg6zkzF2zVFOaCy7PMN7InwIXlstPP1isbNvbolVztmlgPpBT/i\nwfvnKwSWit1SCa09fN/yYr4BArvsnO+W04u6Fc1E1agXYEGG9mNta5s5OLG6iDjP\nPx90P3g3xp0wKOjR8cqD9Y9KQ0pRSUSFHeUkFZkYwwKBgA3wSqRlHSkZTTL88JmS\nh09iZ6+QY9L5gVzEtf4fw3jXEvfrPS5LerSnbuJxZdfPzvmzIZTymVzqn0lwu8PH\ncarZPxfvaH9CYxLXs4yrO0BNZ890gYttyOB9puC8zJintpXwRSeFTfenaWIvRU/e\nMNsF4YM0sGPwRiAffxOunr74\n-----END PRIVATE KEY-----\n`),
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
