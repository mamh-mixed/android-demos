package mongo

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
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
		t.Errorf("update chanMer unsuccessful ", err)
		t.FailNow()
	}
	log.Debugf("update chanMer success %s", chanMer)
}
