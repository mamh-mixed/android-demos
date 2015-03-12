package mongo

import (
	"github.com/omigo/g"
	"quickpay/model"
	"testing"
)

func TestChanMerFind(t *testing.T) {

	chanMer := &model.ChanMer{
		ChanCode:  "000100000",
		ChanMerId: "45672341231",
	}
	err := FindChanMer(chanMer)

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

	err := AddChanMer(chanMer)
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
		SignCert:      "cfcaCert",
		CheckSignCert: "checkcfcaCert",
	}

	err := ModifyChanMer(chanMer)
	if err != nil {
		t.Errorf("update chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("update chanMer success %s", chanMer)
}
