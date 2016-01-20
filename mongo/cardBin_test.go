package mongo

import (
	// "encoding/json"
	// "github.com/CardInfoLink/quickpay/model"
	// "io/ioutil"
	// "strconv"
	"testing"

	"github.com/CardInfoLink/log"
)

/*
func TestCardBinImport(t *testing.T) {
	var arrays []model.CardBin
	bytes, err := ioutil.ReadFile("/opt/gowork/src/github.com/CardInfoLink/quickpay/data/cardBin.json")
	if err != nil {
		t.Error("ERROR:read json file error")
	}

	if err := json.Unmarshal(bytes, &arrays); err != nil {
		t.Errorf("ERROR:unmarshal json to arrays error: %s", err.Error())
	}

	t.Logf("%+v", arrays)
	c := database.C("cardBin")

	// 逐条导入数据，2514条数据用了 0.444秒
	// for idx, cardBin := range arrays {
	// 	if err := c.Insert(cardBin); err != nil {
	// 		t.Errorf("ERROR:insert arrays into mongodb error: %s", err.Error())
	// 	} else {
	// 		t.Logf("INSERT: 插入第%d个数据，内容为%+v", idx+1, cardBin)
	// 	}
	// }

	// 批量导入数据，2514条数据用了0.049秒
	temps := make([]interface{}, len(arrays))
	for idx, value := range arrays {
		bl, _ := strconv.Atoi(value.Bin)
		value.Overflow = strconv.Itoa(bl + 1)
		temps[idx] = value
	}
	if err := c.Insert(temps...); err != nil {
		t.Errorf("ERROR:insert arrays into mongodb error: %s", err.Error())
	}

}
*/

func TestLoadAll(t *testing.T) {
	cbs, err := CardBinColl.LoadAll()
	if err != nil {
		t.Error(err)
	}
	cb, err := CardBinColl.Find("622280193", len(cardNum))
	log.Debugf("%d", len(cbs))
	log.Debugf("%v", cb)
}

func TestFindCardBin(t *testing.T) {

	cardBin, err := CardBinColl.Find("622280193", len(cardNum))
	if err != nil {
		t.Errorf("Find CardBIN error (%s)", err.Error())
	}
	if cardBin != nil && cardBin.Bin != "622280193" {
		t.Errorf("cardNum %s prefix is not %s", cardNum, cardBin.Bin)
	}
}
