package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/g"
)

func TestChanMerFind(t *testing.T) {
	chanMer, err := ChanMerColl.Find(chanCode, chanMerId)
	if err != nil {
		t.Error("find chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("find chanMer success %s", chanMer)
}

func TestChanMerAdd(t *testing.T) {
	chanMer := &model.ChanMer{
		ChanCode:      chanCode,
		ChanMerId:     chanMerId,
		ChanMerName:   chanMerName,
		SettFlag:      settFlag,
		SettRole:      settRole,
		SignCert:      priKeyPem,
		CheckSignCert: checkSignCert,
	}

	// err := ChanMerColl.Add(chanMer)
	// test:update
	err := ChanMerColl.Update(chanMer)
	if err != nil {
		t.Errorf("add chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("add chanMer success %s", chanMer)
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
		t.Errorf("update chanMer unsuccessful ", err)
		t.FailNow()
	}
	g.Debug("update chanMer success %s", chanMer)
}
