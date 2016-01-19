package goconf

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

const hardHalf = "TEZMUboYmBLVfjnduURAk4="

// Config 系统启动时先读取配置文件，绑定到这个 struct 上
var Config = &configStruct{}

// configStruct 对应于 config_<env>.js 文件
type configStruct struct {
	App struct {
		LogLevel           log.Level
		EncryptKey         string
		HTTPAddr           string
		TCPAddr            string
		TCPGBKAddr         string
		DefaultCacheTime   Duration
		NotifyURL          string
		OrderCloseTime     Duration
		OrderRefreshTime   Duration
		SessionExpiredTime Duration
		MonitorMerId       string
	}

	Qiniu struct {
		Bucket string
		Domain string
	}

	Mongo struct {
		Encrypt    bool
		URL        string
		EncryptURL string
		DB         string
	}

	CILOnline struct {
		Host       string
		Port       int
		ServerCert string
	}

	CFCA struct {
		URL                string
		CheckSignPublicKey string
	}

	WeixinScanPay struct {
		URL                 string
		NotifyURL           string
		DNSCacheRefreshTime Duration
	}

	AlipayScanPay struct {
		AlipayPubKey string
		OpenAPIURL   string
		URL          string
		NotifyUrl    string
		AgentId      string
	}

	MobileApp struct {
		WXPMerId  string
		ALPMerId  string
		WebAppUrl string
	}

	UnionLive struct {
		Encrypt           bool
		URL               string
		EncryptKey        string
		EncryptEncryptKey string
		SignKey           string
		EncryptSignKey    string
		ChannelId         string
	}

	Settle struct {
		OverseasSettPoint string
		DomesticSettPoint string
	}
}

// postProcess 后续处理
func (c *configStruct) postProcess() {
	// 拼接出完整的加密密钥
	whole := Config.App.EncryptKey + hardHalf
	encryptKey2, err := base64.StdEncoding.DecodeString(whole)
	if err != nil {
		log.Error("系统密钥配置错误", err)
		os.Exit(1)
	}
	Config.App.EncryptKey = string(encryptKey2)

	// 相对路径变成绝对路径
	Config.CILOnline.ServerCert = util.WorkDir + "/" + Config.CILOnline.ServerCert
	Config.CFCA.CheckSignPublicKey = util.WorkDir + "/" + Config.CFCA.CheckSignPublicKey
	Config.AlipayScanPay.AlipayPubKey = util.WorkDir + "/" + Config.AlipayScanPay.AlipayPubKey

	// 把加密字段解开，赋值给不加密字段
	if Config.Mongo.Encrypt {
		// fmt.Printf("use encrypt url: %s\n", encryptURL)
		url2, err := security.RSADecryptBase64(Config.Mongo.EncryptURL, privateKey)
		if err != nil {
			fmt.Printf("Mongo.EncryptURL decrypt error: %s\n", err)
			os.Exit(1)
		}
		Config.Mongo.URL = string(url2)
	}
	if Config.UnionLive.Encrypt {
		encryptKey2, err := security.RSADecryptBase64(Config.UnionLive.EncryptEncryptKey, privateKey)
		if err != nil {
			fmt.Printf("UnionLive.EncryptEncryptKey decrypt error: %s\n", err)
			os.Exit(1)
		}
		Config.UnionLive.EncryptKey = string(encryptKey2)

		signKey2, err := security.RSADecryptBase64(Config.UnionLive.EncryptSignKey, privateKey)
		if err != nil {
			fmt.Printf("UnionLive.EncryptSignKey decrypt error: %s\n", err)
			os.Exit(1)
		}
		Config.UnionLive.SignKey = string(signKey2)
	}
}

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAywFiuwYKjQZOqyqMayIj/SnMkEsaE1G0JQ15jO4vPdSdpMDy
Hn1fXOw8Z9T0fmjk5D2NPi9zLUbDGo1+22RYpbGs8hL5PPqLrFaGY8kGjMGa3e9Y
noo0lqoK53iUm4Zqd6FMFl2jBa7u15T8VWHwZSBzXn+slhlKFkJHMTruxih0U5BN
yGPbf4y37+6kLNilFBYVQ9837VdwT2eW4sP5ah2E07DCPUqwT8UlWe11rfPPqCZh
1wLEqYeQkvqo0J5juE9S/c/PNiOVl+mwsWGyV71OrVmqvKEFlxKOK906u24lrUHK
ek7Kq4oQui5Wt89m4yr5TwAsHeIcyV4kYYIe1QIDAQABAoIBAC1AnNKV8Soomsa7
EFwdWypm8+vCYgimcOLFky/gNHWy/IUqYY58YhKjsn9u0CWRmlxqgB65kxInsPwt
SHb9cmlVJvk7U4XNT+9VxlVeDXC5A52vafDFXB2twAqDLZVRrFAIi558twdgTGuQ
EYOy9lSEnFMXYNCAyKwXkCkgOvO+0PJzewVlo1fjDSLpR5B/kONWFrHM5M+up0nZ
VTObtB0TmB3lYTxVI9+sBFYPoOIeGnA0RLo/0UgFUxNLDsjODsq0n2kIclkje9dz
r6VDf6iA8NWa5jk4ICCyXO5APQXKYJqDt1UWxJV8e3tRSCuZznVFIfvmsk4Vcjr3
II9MLQ0CgYEA9NoexZpROUOvM5oVxi/wwOkWGSi7I+k7/GrbRR3N0GusZozacsO3
ToeFy3RKJoWh4cm5wQWcedMDpoNJPHqDc2HVGgvN2oPYF65vv6xAl/Qe+Ns1dHAQ
Js0IKt2tuItTVu8i1F0KELf1HdJharx4cQAydIWkPxA1D1iqGI1BPncCgYEA1D+F
YLOo9FJtt4ocvp/1jYQo3c9XezOWbsjeIq9W2//vvUIh1Y1D9kv0oYyMooVeqLLL
zJmYauifuOosS9UwIbkXlLi+V/3XZj0tNGGLoST2MR4BTvPWrXa+OFPBTH7dl17X
Iq4Ln4oJdnmZhFmmuecLG0t57CngJTZcit3cZBMCgYBgJ9yXy2+EZpFCWYudhiwt
BhxYiwdbJfgZu7kanoa6B97vcvdCxJuTKmOfr66DDE2zhu384IA/01+Gn9498vr2
cAApN2ODIe3V9voJstK3Gfaj0ipe7LdbFX/UnbPgWk7DQCxUa6lNQYDwUjNRoGxI
LESkP+ttnKbJvQ8njymFJwKBgG6e8LgWyy8TqwVm4VZk7kkkoVwBzblziKsS29u+
AQpGmT/NsO6pYsuCiOyN4VpvIofQMDHht7O4rE5nFlEruptI6cZkhyg7L8GkjuPn
FywUpI+y8MGiirf71GZtGKjy0jErh/sWNQ6glg/+jomRZDkt9vbx3oi8xor+izsB
KPWrAoGBAIdJz5K1EHhJWowgkaAhUiMvK8NaB2UJ9GQoK7b/SPVirsXKbDKQS/AG
ca+TKm+By6vG6qfqeDnpLxFxM/jF1u01ASvEr9usAjVpxnGamlFqJwyYmeV/yyKj
V/C3a5OiMf8AOU5rGfbb0SWXehiyiX4nH9q1fgO8pkqIWQ/C0UNr
-----END RSA PRIVATE KEY-----
`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAywFiuwYKjQZOqyqMayIj
/SnMkEsaE1G0JQ15jO4vPdSdpMDyHn1fXOw8Z9T0fmjk5D2NPi9zLUbDGo1+22RY
pbGs8hL5PPqLrFaGY8kGjMGa3e9Ynoo0lqoK53iUm4Zqd6FMFl2jBa7u15T8VWHw
ZSBzXn+slhlKFkJHMTruxih0U5BNyGPbf4y37+6kLNilFBYVQ9837VdwT2eW4sP5
ah2E07DCPUqwT8UlWe11rfPPqCZh1wLEqYeQkvqo0J5juE9S/c/PNiOVl+mwsWGy
V71OrVmqvKEFlxKOK906u24lrUHKek7Kq4oQui5Wt89m4yr5TwAsHeIcyV4kYYIe
1QIDAQAB
-----END PUBLIC KEY-----
`)
