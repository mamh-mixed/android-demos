package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"testing"
)

func TestVersionAdd(t *testing.T) {
	if false {
		v := &model.Version{
			Vn:     "20150421120000",
			LastVn: "",
			VnType: "cardBin",
		}

		err := VersionColl.Add(v)
		if err != nil {
			t.Errorf("fail to add one version(%+v): %s", v, err)
			t.FailNow()
		}
	}
}

func TestVersionFind(t *testing.T) {
	v, err := VersionColl.FindOne("cardBin")
	if err != nil {
		t.Errorf("fail to find cardBin version : %s", err)
		t.FailNow()
	}
	log.Debugf("%+v", v)
}
