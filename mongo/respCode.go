package mongo

import (
	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type respCodeCollection struct {
	name string
}

// RespCodeColl 应答码 Collection
var RespCodeColl = respCodeCollection{"respCode"}

var respCodeCache = cache.New(model.Cache_RespCode)
var unKnown = &model.BindingReturn{RespCode: "000004", RespMsg: "未知应答，请联系系统管理员", IsRetChanRespMsg: true}

// Get 根据传入的code类型得到Resp对象
func (c *respCodeCollection) Get(code string) (resp *model.BindingReturn) {

	// o, found := respCodeCache.Get(code)
	// if found {
	// 	resp = o.(*model.BindingReturn)
	// 	return resp
	// }

	resp = &model.BindingReturn{}
	err := database.C(c.name).Find(bson.M{"respCode": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	if err != nil {
		log.Errorf("can not find respCode for %s: %s", code, err)
		return resp
	}

	// save cache
	// respCodeCache.Set(code, resp, cache.NoExpiration)

	return resp
}

// GetByCfca 根据传入的cfca的code得到Resp对象
func (c *respCodeCollection) GetByCfca(code string) (resp *model.BindingReturn) {
	resp = &model.BindingReturn{}
	err := database.C(c.name).Find(bson.M{"cfca.code": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	if err != nil {
		log.Errorf("find cfca code(%s) error: %s", code, err)
		return unKnown
	}
	return resp
}

// GetByCIL 由线下应答码获得系统应答码
func (c *respCodeCollection) GetByCIL(code string) (resp *model.BindingReturn) {
	resp = &model.BindingReturn{}
	database.C(c.name).Find(bson.M{"cil.code": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	return resp
}

// GetMsg 根据传入的code类型得到msg
func (c *respCodeCollection) GetMsg(code string) (msg string) {
	resp := &model.BindingReturn{}
	database.C(c.name).Find(bson.M{"respCode": code}).Select(bson.M{"respMsg": 1}).One(resp)
	msg = resp.RespMsg
	return msg
}

/* only use for import respCode */

func (c *respCodeCollection) Add(r *model.QuickpayCSV) error {
	err := database.C(c.name).Insert(r)
	return err
}

func (c *respCodeCollection) FindOne(code string) (*model.QuickpayCSV, error) {
	q := new(model.QuickpayCSV)
	err := database.C(c.name).Find(bson.M{"respCode": code}).One(q)
	return q, err
}

func (c *respCodeCollection) Update(r *model.QuickpayCSV) error {
	err := database.C(c.name).Update(bson.M{"respCode": r.RespCode}, r)
	return err
}
