package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type respCodeCollection struct {
	name string
}

func init() {
	Connect()
}

// RespCodeColl 应答码 Collection
var RespCodeColl = respCodeCollection{"respCode"}

// Get 根据传入的code类型得到Resp对象
func (c *respCodeCollection) Get(code string) (resp *model.BindingReturn) {
	resp = &model.BindingReturn{}
	err := database.C(c.name).Find(bson.M{"respCode": code}).Select(bson.M{"respCode": 1, "respMsg": 1}).One(resp)
	if err != nil {
		log.Errorf("can not find respCode for %s: %s", code, err)
	}

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

/* only use for import respCode */

func (c *respCodeCollection) Add(r *model.QuickpayCsv) error {
	err := database.C(c.name).Insert(r)
	return err
}

func (c *respCodeCollection) FindOne(code string) (*model.QuickpayCsv, error) {
	q := new(model.QuickpayCsv)
	err := database.C(c.name).Find(bson.M{"respCode": code}).One(q)
	return q, err
}

func (c *respCodeCollection) Update(r *model.QuickpayCsv) error {
	err := database.C(c.name).Update(bson.M{"respCode": r.RespCode}, r)
	return err
}
