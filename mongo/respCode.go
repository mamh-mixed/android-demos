package mongo

import (
	"quickpay/model"

	"gopkg.in/mgo.v2/bson"
)

// GetRespCode 根据传入的code类型得到Resp对象
func GetRespCode(code string) (resp *model.BindingReturn) {
	resp = &model.BindingReturn{}
	db.respCode.Find(bson.M{"respCode": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	return resp
}

// GetRespCodeByCfca 根据传入的cfca的code得到Resp对象
func GetRespCodeByCfca(code string) (resp *model.BindingReturn) {
	resp = &model.BindingReturn{}
	db.respCode.Find(bson.M{"cfca.code": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	return resp
}
