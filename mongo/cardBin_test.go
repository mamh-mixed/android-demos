package mongo

import
// "encoding/json"
// "io/ioutil"
"testing"

/*
func TestCardBinImport(t *testing.T) {
	var arrays []CardBin
	bytes, err := ioutil.ReadFile("/opt/gowork/src/quickpay/card_bin.json")
	if err != nil {
		t.Error("ERROR:read json file error")
	}

	if err := json.Unmarshal(bytes, &arrays); err != nil {
		t.Errorf("ERROR:unmarshal json to arrays error: %s", err.Error())
	}

	t.Logf("%+v", arrays)
	c := db.C("cardBin")

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
		temps[idx] = value
	}
	if err := c.Insert(temps...); err != nil {
		t.Errorf("ERROR:insert arrays into mongodb error: %s", err.Error())
	}

}
*/

func TestFindCardBin(t *testing.T) {

	var cardNum = "6222801932062061908"
	cardBin := CardBinColl.Find(cardNum)

	if cardBin.Bin != "622280193" {
		t.Errorf("cardNum %s prefix is not %s", cardNum, cardBin.Bin)
	}
}
