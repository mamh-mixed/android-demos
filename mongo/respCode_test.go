package mongo

import (
	"testing"
)

func TestGetRespCode(t *testing.T) {
	ret := RespCodeColl.Get(respCode)
	// t.Logf("%#v", ret)

	if ret == nil {
		t.Error("respCode not exist")
	}

	expected := "金额过大"
	if ret.RespMsg != expected {
		t.Errorf("respCode %s message must be `%s`, but get %s", respCode, expected, ret.RespMsg)
	}
}

func TestGetRespCodeByCfca(t *testing.T) {

	ret := RespCodeColl.GetByCfca(cfcacode)
	// t.Logf("%#v", ret)

	if ret == nil {
		t.Error("cfca code not exist")
	}

	expected := "金额过大"
	if ret.RespMsg != expected {
		t.Errorf("respCode %s message must be `%s`, but get %s", cfcacode, expected, ret.RespMsg)
	}
}
