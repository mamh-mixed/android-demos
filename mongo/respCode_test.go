package mongo

import "testing"

func TestGetRespCode(t *testing.T) {
	respCode := "200125"
	ret := GetRespCode(respCode)
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
	ret := GetRespCodeByCfca(code)
	// t.Logf("%#v", ret)

	if ret == nil {
		t.Error("cfca code not exist")
	}

	expected := "金额过大"
	if ret.RespMsg != expected {
		t.Errorf("respCode %s message must be `%s`, but get %s", code, expected, ret.RespMsg)
	}
}
