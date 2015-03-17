package mongo

import (
	"quickpay/mongo"
	"testing"
)

func TestGetRespCode(t *testing.T) {
	respCode := "200125"
	ret := mongo.RespCodeColl.Get(respCode)
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
	code := "270032"
	ret := mongo.RespCodeColl.GetByCfca(code)
	// t.Logf("%#v", ret)

	if ret == nil {
		t.Error("cfca code not exist")
	}

	expected := "金额过大"
	if ret.RespMsg != expected {
		t.Errorf("respCode %s message must be `%s`, but get %s", code, expected, ret.RespMsg)
	}
}
