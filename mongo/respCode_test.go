package mongo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/omigo/g"
	"gopkg.in/mgo.v2/bson"
)

//
// func TestGetRespCode(t *testing.T) {
//
// 	//channel
// 	channel := GetRespCode("240021", "cfca")
//
// 	if channel == nil {
// 		t.Error("codeType does not exist")
// 	}
// 	//sys
// 	sys := GetRespCode(channel.RespCode, "sys")
// 	if sys == nil {
// 		t.Error("codeType does not exist")
// 	}
//
// 	g.Debug("sucess,resp obj (%s),(%s)", channel, sys)
//
// }

func TestInitRespCode(t *testing.T) {

	rc, _ := os.Open("../respCode.json")
	g.Debug("file path ", rc.Name())
	bytes, _ := ioutil.ReadAll(rc)
	var resps []Resp
	err := json.Unmarshal(bytes, &resps)
	if err != nil {
		t.Error(err)
	}
	g.Debug("read respCode from .json", resps)

	vals := make([]interface{}, len(resps))
	for i, v := range resps {
		vals[i] = v
	}
	//delete
	db.respCode.DropCollection()

	db.respCode.Insert(vals...)

}

func TestInitCFCARespCode(t *testing.T) {

	rc, _ := os.Open("../channel/cfca/respCodeMap.json")
	bytes, _ := ioutil.ReadAll(rc)
	var cfcas []CfcaReader
	err := json.Unmarshal(bytes, &cfcas)
	if err != nil {
		t.Error(err)
	}
	g.Debug("read respCode from .json", len(cfcas))

	for _, c := range cfcas {
		o := Cfca{c.CfcaCode, c.CfcaMsg}
		var resp RespC
		db.respCode.Find(bson.M{"respcode": c.RespCode}).One(&resp)
		ccs := resp.Cfca
		ccs = append(ccs, o)
		resp.Cfca = ccs
		// g.Debug("updated resp ", resp)
		db.respCode.Update(bson.M{"respcode": c.RespCode}, resp)
	}
}

type CfcaReader struct {
	RespCode string
	RespMsg  string
	CfcaCode string
	CfcaMsg  string
}

type Cfca struct {
	Code string
	Msg  string
}

type RespC struct {
	RespCode string
	RespMsg  string
	Cfca     []Cfca `bson:",omitempty"`
}
