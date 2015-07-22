package mongo

import (
	"testing"
)

func TestGetRespCode(t *testing.T) {
	respCode = "200021"

	ret := RespCodeColl.Get(respCode)

	if ret == nil {
		t.Error("respCode not exist")
	}

	expected := "金额过大"
	if ret.RespMsg != expected {
		t.Errorf("respCode %s message must be `%s`, but get %s", respCode, expected, ret.RespMsg)
	}
}

func TestGetRespCodeByCfca(t *testing.T) {

	// ret := RespCodeColl.GetByCfca(cfcacode)
	ret := ScanPayRespCol.GetByWxp("SYSTEMERROR", "prePay")

	if ret == nil {
		t.Error("cfca code not exist")
	}

	expected := "外部系统错误"
	if ret.ISO8583Msg != expected {
		t.Errorf("respCode %s message must be `%s`, but get %s", cfcacode, expected, ret.RespMsg)
	}
}

func TestGetRespCodeByCIL(t *testing.T) {

	ret := RespCodeColl.GetByCIL("00")

	if ret == nil {
		t.Error("cfca code not exist")
	}

	expected := "000000"
	if ret.RespCode != expected {
		t.Errorf("respCode %s message must be `%s`, but get %s", cfcacode, expected, ret.RespMsg)
	}
	t.Logf("result is %+v", ret)
}
