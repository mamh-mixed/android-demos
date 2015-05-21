package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
)

var (
	debug            = false
	orderNum         = "123456"
	hexId            = "550ee5e36a3dd74f93000001"
	transType        = 1
	transStatus      = "00"
	respCode         = "200125"
	merId            = "012345678901234"
	merStatus        = "Normal"
	transCurr        = "156"
	signKey          = ""
	encryptKey       = ""
	cardNum          = "6222801932062061908"
	createTime       = "2015-03-26"
	chanMerId        = "alp0002"
	chanCode         = "ALP"
	settFlag         = "457"
	cardBrand        = "ALP"
	settRole         = "testRole"
	checkSignCert    = "check"
	chanBingingId    = "cf00fd61d5ef4d924485db88b584897e"
	chanOrderNum     = "aaaaaaaaaaaaaaaaaabb"
	chanOrigOrderNum = "aaaaaaaaaaaaaaaaaabb"
	chanMerName      = "讯联测试商户"
	acctName         = "张三"
	acctNum          = "6222020302062061908"
	identType        = "0"
	identNum         = "350583199009153732"
	phoneNum         = "18205960039"
	acctType         = "10"
	validDate        = ""
	cvv2             = ""
	sendSmsId        = "1000000000009"
	smsCode          = "12353"
	bankId           = "102"
	cfcacode         = "270032"
	priKeyPem        = "eu1dr0c8znpa43blzy1wirzmk8jqdaon"

//      priKeyPem        = `-----BEGIN RSA PRIVATE KEY-----
// MIICXQIBAAKBgQCvJC9MMGRKmxRBI0KMjDtz2KooIc6XOljHPWhTfAamhV3A5v5y
// PiZr4haMDpulU08Y0JxsegwDwfbscQrhG7nvilIqIa+HiI1xkfFxjtNUrMN5hpvO
// 8HUUfwqzb5EdllQcv/C0xxBkeCECIb86JJry7ty4mNBkN2idbGxldMi90QIDAQAB
// AoGATvTIIdfbDss06Vyk/smlb8dohmkfQov6Q/AKHUDXmrCbIIDCiuw70/z73y4i
// uviAuxYovrqSugryb4tStUMTogmft4methz1/O/083XHwBNKBPnS2fobYDfBxqkX
// tH26woCjrEr/O/wngo6iFp7b5yJlyXapN0x+iOF3CShIhAECQQD2gZ6LLYdxSP8i
// aRYAPOh10mF5IHt2dl89eOjNiqVGMlkV5aXNT80jAQr/kWGZfIjscb/xkawSKQKs
// ovcn99GRAkEAteL02mBrCLfn2idBwXTdil+yeigReAZmRpqQuAfTRZN4RM+5Dw3q
// X0IiCkR3oyiwx89n1eGmz1JTZRxoY1AIQQJAWVbQ5xAxLlWOYiJD3wI0Hb+JpCSp
// ml18VwMjHJtLGw3US6NXW/m4Fx+hpM5D2STRWyA+uIZbHpnOZlMJ0Gp4gQJBAK38
// 66JV5y1Q1r2tHc6UHzQ1tMH7wDIjVQSm6FbSTXxZxAt29Rx8gD0dQvi1ZAg0bV7F
// fRtwnqPlqZaoJQcTUMECQQD1Dh+Mu3OMb5AHnrtbk9l1qjM3U81QBKdyF0RY+djo
// b3cR9I7+hurpqhJmQ7yuvAWe2xWc+YNTQ48FDJTogXlB
// -----END RSA PRIVATE KEY-----`
)

func TestChanMerFind(t *testing.T) {
	chanMer, err := ChanMerColl.Find(chanCode, chanMerId)
	if err != nil {
		t.Error("find chanMer unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("find chanMer success %s", chanMer)
}

func TestChanMerAdd(t *testing.T) {
	chanMer := &model.ChanMer{
		ChanCode:      "TEST1",
		ChanMerId:     "TEST",
		ChanMerName:   "TEST",
		SettFlag:      "TEST",
		SettRole:      "TEST",
		SignCert:      "TEST",
		CheckSignCert: "TEST",
	}

	err := ChanMerColl.Add(chanMer)
	// test:update
	// err := ChanMerColl.Add(chanMer)
	if err != nil {
		t.Errorf("add chanMer unsuccessful %s", err)
		t.FailNow()
	}
	log.Debugf("add chanMer success %s", chanMer)
}

func TestChanMerModify(t *testing.T) {
	chanMer := &model.ChanMer{
		ChanCode:      chanCode,
		ChanMerId:     chanMerId,
		ChanMerName:   chanMerName,
		SettFlag:      settFlag,
		SettRole:      settRole,
		SignCert:      priKeyPem,
		CheckSignCert: checkSignCert,
	}

	err := ChanMerColl.Update(chanMer)
	if err != nil {
		t.Errorf("update chanMer unsuccessful %s", err)
		t.FailNow()
	}
	log.Debugf("update chanMer success %s", chanMer)
}

func TestChanMerFindAll(t *testing.T) {

	cs, err := ChanMerColl.FindByCode(chanCode)
	if err != nil {
		t.Errorf("findAll chanMer unsuccessful %s", err)
		t.FailNow()
	}
	log.Debugf("%+v", cs)
}
