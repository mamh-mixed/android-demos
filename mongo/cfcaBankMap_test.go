package mongo

import (
	// "encoding/json"
	// "github.com/CardInfoLink/quickpay/model"
	// "io/ioutil"
	"testing"
)

func TestFindCfcaBankMap(t *testing.T) {
	insCode := "803390000"

	cm, err := CfcaBankMapColl.Find(insCode)
	if err != nil {
		t.Errorf("Find cfcaBankMap error: (%s)", err.Error())
	}
	t.Logf("CfcaBankMap is: %+v", cm)
}

/*
func TestImportCfcaBankMap(t *testing.T) {
	var arrays []model.CfcaBankMap

	bytes, err := ioutil.ReadFile("/opt/gowork/src/github.com/CardInfoLink/quickpay/data/cfcaBankMap.json")
	if err != nil {
		t.Error("ERROR:read json file error")
	}

	if err := json.Unmarshal(bytes, &arrays); err != nil {
		t.Errorf("ERROR:unmarshal json to arrays error: %s", err.Error())
	}

	temps := make([]interface{}, len(arrays))
	for idx, value := range arrays {
		temps[idx] = value
	}

	t.Logf("%+v", arrays)

	if err := database.C("cfcaBankMap").Insert(temps...); err != nil {
		t.Errorf("ERROR:insert arrays into mongodb error: %s", err.Error())
	}

}
*/
