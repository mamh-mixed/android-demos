package mongo

import (
	"github.com/omigo/g"
	"testing"
)

func TestChanMerInit(t *testing.T) {

	chanMer := ChanMer{
		ChanCode:  "000100000",
		ChanMerId: "45672341231",
	}
	err := chanMer.Init()

	if err != nil {
		t.Error("init chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("init chanMer success %s", chanMer)
}

func TestChanMerAdd(t *testing.T) {
	chanMer := ChanMer{
		ChanCode:       "CFCA",
		ChanMerId:      "001405",
		ChanMerName:    "测试渠道商户",
		SettlementFlag: "457",
		SettlementRole: "testRole",
		SignCert:       "cfcaCert",
		CheckSignCert:  "checkcfcaCert",
	}

	err := chanMer.Add()
	if err != nil {
		t.Errorf("add chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("add chanMer success %s", chanMer)
}
