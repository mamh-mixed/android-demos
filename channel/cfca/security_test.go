package cfca

import (
	"encoding/base64"
	"testing"
)

const (
	priKeyPem = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCvJC9MMGRKmxRBI0KMjDtz2KooIc6XOljHPWhTfAamhV3A5v5y
PiZr4haMDpulU08Y0JxsegwDwfbscQrhG7nvilIqIa+HiI1xkfFxjtNUrMN5hpvO
8HUUfwqzb5EdllQcv/C0xxBkeCECIb86JJry7ty4mNBkN2idbGxldMi90QIDAQAB
AoGATvTIIdfbDss06Vyk/smlb8dohmkfQov6Q/AKHUDXmrCbIIDCiuw70/z73y4i
uviAuxYovrqSugryb4tStUMTogmft4methz1/O/083XHwBNKBPnS2fobYDfBxqkX
tH26woCjrEr/O/wngo6iFp7b5yJlyXapN0x+iOF3CShIhAECQQD2gZ6LLYdxSP8i
aRYAPOh10mF5IHt2dl89eOjNiqVGMlkV5aXNT80jAQr/kWGZfIjscb/xkawSKQKs
ovcn99GRAkEAteL02mBrCLfn2idBwXTdil+yeigReAZmRpqQuAfTRZN4RM+5Dw3q
X0IiCkR3oyiwx89n1eGmz1JTZRxoY1AIQQJAWVbQ5xAxLlWOYiJD3wI0Hb+JpCSp
ml18VwMjHJtLGw3US6NXW/m4Fx+hpM5D2STRWyA+uIZbHpnOZlMJ0Gp4gQJBAK38
66JV5y1Q1r2tHc6UHzQ1tMH7wDIjVQSm6FbSTXxZxAt29Rx8gD0dQvi1ZAg0bV7F
fRtwnqPlqZaoJQcTUMECQQD1Dh+Mu3OMb5AHnrtbk9l1qjM3U81QBKdyF0RY+djo
b3cR9I7+hurpqhJmQ7yuvAWe2xWc+YNTQ48FDJTogXlB
-----END RSA PRIVATE KEY-----`
)

func TestSignatureUseSha1WithRsa(t *testing.T) {
	data := `<Request version="2.0"><Head><TxCode>2501</TxCode><InstitutionID>001405</InstitutionID></Head><Body><TxSNBinding>15030622072014626553</TxSNBinding><BankID>700</BankID><AccountName>张三</AccountName><AccountNumber>1503063124684673</AccountNumber><IdentificationType>0</IdentificationType><IdentificationNumber>1503063937742309</IdentificationNumber><PhoneNumber>13333333333</PhoneNumber><CardType>10</CardType></Body></Request>`

	hexSign := signatureUseSha1WithRsa([]byte(data), priKeyPem)

	expected := "0c958e3fa28e5b4b4c112276510386cb53f1cb080c70d3905fafc764f1daea59e7e1ecb093f50ff85f26b6ee9364c5a278cec8420cd1d480ce8d6a57cfb01fefa2be61f4dcc7e20295bacc95cbbf7847d7089bff651efa19299f324eb0f143751e907af0606ab9e2be79702ebe33043ff0d7d668202a98f0ef577f1fed51cb6d"
	if hexSign != expected {
		t.Error("签名不正确")
	}
}

func TestCheckSignatureUseSha1WithRsa(t *testing.T) {
	b64Data := "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjxSZXNwb25zZSB2ZXJzaW9uPSIyLjAiPgo8SGVhZD4KPENvZGU+MzAyNTAxMzE8L0NvZGU+CjxNZXNzYWdlPui6q+S7veivgeWPt+S4jeWQiOazle+8jOivt+ajgOafpe+8gTwvTWVzc2FnZT4KPC9IZWFkPgo8L1Jlc3BvbnNlPg=="
	hexSign := "79C0AC23DCD547D5640F85B5216B21BB19FFA420CE5AD581520717BCC25ABCD7BABAC4B2590FEE2869FF5FE8931CB8FBFA4D066944B3CC5FDBC9C19BB8F9E39E933FD3FF40D2F38F4D7714925621C5353E8F59098FEA00159859FC3FF93C064DD83ADB161B0AFD636D5C043C6FD11B4A72D083B55BF4B1E2ACEAECC2B2F933F5"

	data, _ := base64.StdEncoding.DecodeString(b64Data)
	t.Log(string(data))

	err := checkSignatureUseSha1WithRsa(data, hexSign)
	if err != nil {
		t.Error("验签失败")
	}
}

func TestCheckSignatureUseSha1WithRsa2(t *testing.T) {
	b64Data := "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjxSZXNwb25zZSB2ZXJzaW9uPSIyLjAiPgo8SGVhZD4KPENvZGU+MjQwMDAyPC9Db2RlPgo8TWVzc2FnZT7ns7vnu5/kuK3kuI3lrZjlnKjmjIflrprnmoTmnLrmnoTvvIzmn6XnnIvlj4LmlbBJbnN0aXR1dGlvbklEPC9NZXNzYWdlPgo8L0hlYWQ+CjwvUmVzcG9uc2U+"
	hexSign := "40E88F745B9F050D116E627CE852F4CD8CE85AE25066DDAC01A758F55B9393F4DF4F6AB7099A8D212EB4EF3F6BD29958D00BEA553DA8367BA257AD0634147F3DD35AD17D68C3C6FC3966AE97814124F36B28A930409C841E1FD0B81B2FEEA16389C56A9314195557394583ED00B6FE5F2BE2C5F7F8C6D7889B724A1656EAB5B3"

	data, _ := base64.StdEncoding.DecodeString(b64Data)

	t.Log(string(data))

	err := checkSignatureUseSha1WithRsa(data, hexSign)
	if err != nil {
		t.Error("验签失败")
	}
}
