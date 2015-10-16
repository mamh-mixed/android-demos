package mongo

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/security"
)

var database *mgo.Database

// Connect 程序启动时，或者，单元测试前，先连接到 MongoDB 数据库
func init() {
	encrypt := goconf.Config.Mongo.Encrypt
	url := goconf.Config.Mongo.URL
	encryptURL := goconf.Config.Mongo.EncryptURL
	dbname := goconf.Config.Mongo.DB

	if encrypt {
		// fmt.Printf("use encrypt url: %s\n", encryptURL)
		url2, err := security.RSADecryptBase64(encryptURL, privateKey)
		if err != nil {
			fmt.Printf("unable connect to mongodb server %s\n", err)
			os.Exit(1)
		}
		url = string(url2)
	}

	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("unable connect to mongodb server %s\n", err)
		os.Exit(1)
	}

	session.SetMode(mgo.Eventual, true) //需要指定为Eventual
	session.SetSafe(&mgo.Safe{})

	database = session.DB(dbname)

	// 不能在日志中出现数据库密码
	// fmt.Println("connected to mongodb host `%s` and database `%s`", url, dbname)
	fmt.Printf("connected to mongodb host `%s` and database `%s`\n", "***", dbname)
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
