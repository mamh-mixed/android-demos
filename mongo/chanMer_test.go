package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"

	"github.com/omigo/g"
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

func TestChanMerFind(t *testing.T) {

	// chanMer := &model.ChanMer{
	// 	ChanCode:  "000100000",
	// 	ChanMerId: "45672341231",
	// }
	chanMer, err := ChanMerColl.Find("000100000", "45672341231")
	if err != nil {
		t.Error("find chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("find chanMer success %s", chanMer)
}

func TestChanMerAdd(t *testing.T) {
	chanMer := &model.ChanMer{
		ChanCode:      "CFCA",
		ChanMerId:     "001405",
		ChanMerName:   "测试渠道商户",
		SettFlag:      "457",
		SettRole:      "testRole",
		SignCert:      "cfcaCert",
		CheckSignCert: "checkcfcaCert",
	}

	err := ChanMerColl.Add(chanMer)
	if err != nil {
		t.Errorf("add chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("add chanMer success %s", chanMer)
}

func TestChanMerModify(t *testing.T) {
	chanMer := &model.ChanMer{
		ChanCode:      "CFCA",
		ChanMerId:     "001405",
		ChanMerName:   "测试渠道商户",
		SettFlag:      "457",
		SettRole:      "testRole",
		SignCert:      priKeyPem,
		CheckSignCert: "checkcfcaCert",
	}

	err := ChanMerColl.Update(chanMer)
	if err != nil {
		t.Errorf("update chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("update chanMer success %s", chanMer)
}
