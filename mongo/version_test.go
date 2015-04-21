package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestVersionAdd(t *testing.T) {
	// connect to mongodb
	Connect()

	v := &model.Version{
		Vn:     "20150421120000",
		LastVn: "",
		VnType: "cardBin",
	}

	err := versionColl.Add(v)
	if err != nil {
		t.Errorf("fail to add one version(%+v): %s", v, err)
		t.FailNow()
	}

}
