package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type cfcaBankMapCollection struct {
	name string
}

var CfcaBankMapColl = cfcaBankMapCollection{"cfcaBankMap"}

// Find 根据卡BIN中的发卡行号查找中金支持的银行映射
func (c *cfcaBankMapCollection) Find(insCode string) (cb *model.CfcaBankMap, err error) {
	cb = new(model.CfcaBankMap)
	q := bson.M{"insCode": insCode}

	err = database.C(c.name).Find(q).One(cb)

	if err != nil {
		log.Errorf("'Find CfcaBankMap ERROR!' Error message is: %s\n; Condition is: %+v", err.Error(), q)
		return nil, err
	}

	return
}
