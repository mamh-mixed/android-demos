package mongo

import (
	"quickpay/model"

	"gopkg.in/mgo.v2/bson"
)

type respCodeCollection struct {
	name string
}

var RespCodeColl = respCodeCollection{"respCode"}

// Get 根据传入的code类型得到Resp对象
func (c *respCodeCollection) Get(code string) (resp *model.BindingReturn) {
	resp = &model.BindingReturn{}
	database.C(c.name).Find(bson.M{"respCode": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	return resp
}

// GetByCfca 根据传入的cfca的code得到Resp对象
func (c *respCodeCollection) GetByCfca(code string) (resp *model.BindingReturn) {
	resp = &model.BindingReturn{}
	database.C(c.name).Find(bson.M{"cfca.code": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	return resp
}

// GetMsg 根据传入的code类型得到msg
func (c *respCodeCollection) GetMsg(code string) (msg string) {
	resp := &model.BindingReturn{}
	database.C(c.name).Find(bson.M{"respCode": code}).Select(bson.M{"respMsg": 1}).One(resp)
	msg = resp.RespMsg
	return msg
}
