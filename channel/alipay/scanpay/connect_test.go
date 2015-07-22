package scanpay

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const privateKeyPem = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDXoE7uGddakVk8AygOCtsX/z1EfIESLXRWiqepRy4iAY0QOjv3
rjAoRT0qfTnMlzSizwXbgSPBK6Hw8JRGql8a+guKDxEER2OwxPBojM3j3sHgmVpD
cVmQDVfybQBWSqjr0UCZNNL5jrvKlN7KvYSbTP49jk/KQozbhiSZJXGaowIDAQAB
AoGBAJgt1qbqa/fKby0QmTtX5AsKo4XrTPi0RdAyKWQqDWAIsjMKtnn2YJx7SVDs
cld3O7nP3DVv7fkOP0uZrzw0T8ugjN6tlKcNL3OLyv2jLzzhLbbsSCMYD2C4SsS+
Ka6NnTJO0GArE3PW0tqRqjs5UMzo1bFkKeLMZymneELAlGDhAkEA/oVTW8YG6RgY
afvOEPdVMl0EbQur4hmyyCOSeHFvKORjCqUKJaiX19PolCq9UyFNeKlL++FMmKX0
uEEGm9WmsQJBANjhHZYqx4QP05qNeLnmdsLIznU6Myha+ZJLuGbCq9M8dOq+9wLK
uYFsyYfXbmo/2NTxXcFqL138uCdWQd0P05MCQD6ikOEE2q6CP3/Vd+C0/UJnoa80
MBh0OosGNgVt5O0rRzRXaSfbVYLHo3TTD8RlbatD/m7+AtuN+6tcVUQEUAECQQCZ
yce5HEyuEKr0BS1+ZTYBmXMNHV/5VclzO85ez9wXxd8CNrfheu9gH0woz1K0dOHE
3gKljC5ab0IGYtwBbZ+RAkA9OtNP8bFlW3LpHWgq4hsGVWdGpr5LYewWMLQXAaWN
8E/0d0BdmQr0SR4EIUHD4ociBJHyrttJRTSyQSnpATva
-----END RSA PRIVATE KEY-----
`

func TestPrepareData(t *testing.T) {
	d := &CancelReq{
		CommonParams: CommonParams{
			AppID:      "2015051100069108",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		OutTradeNo: "14141341234",
	}

	v, err := prepareData(d)
	if err != nil {
		t.Errorf("prepare data error: %s", err)
	}

	t.Logf("%s", v.Encode())
}

func TestBarcodePay(t *testing.T) {
	u := "https://openapi.alipay.com/gateway.do?charset=utf-8"

	v := url.Values{}

	v.Set("method", "alipay.trade.pay")
	v.Set("sign", "EG7vAe13LppMu1h7YFden307Hs6TrDX8Uo0g9h29xgqsX/6j9A1lyxAQ5Hj3t4HGHpiOjx3LZ4hLJqER20GjVY1ay+4E/8nCdrD6XvVhd2pPtPTxoVFNxw6dPZ3J9LXWFB05I8qNcQCUhWf2HROBTjCDXTm9NmVeIGBai/4KGLY=")
	v.Set("charset", "utf-8")
	v.Set("sign_type", "RSA")
	v.Set("app_id", "2014072300007148")
	v.Set("timestamp", "2014-07-24 03:07:50")
	v.Set("biz_content", `{
    "out_trade_no": "201503022001",
    "scene": "bar_code",
    "auth_code": "23383838388383883",
    "total_amount": "88.88",
    "discountable_amount":"8.88",
    "undiscountable_amount ": "80",
    "subject": "条码支付",
    "goods_detail": [
    {
    "goods_id": "apple-01",
    "goods_name": "ipad",
    "goods_category": "7788230",
    "price": "88.88",
    "quantity": "1"
    }
    ],
    "operator_id": "op001",
    "store_id": "pudong001",
    "terminal_id": "t_001",
    "time_expire": "2015-01-24 03:07:50"
    }`)

	res, err := http.PostForm(u, v)
	if err != nil {
		t.Errorf("post error: %s", err)
	}
	defer res.Body.Close()

	rbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("read body error: %s", err)
	}
	t.Logf("%s", rbody)
}
