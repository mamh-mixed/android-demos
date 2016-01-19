package mongo

import (
	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

type cfcaBankMapCollection struct {
	name string
}

var CfcaBankMapColl = cfcaBankMapCollection{"cfcaBankMap"}

var cfcaCache = cache.New(model.Cache_CfcaBankMap)

// Find 根据卡BIN中的发卡行号查找中金支持的银行映射
func (c *cfcaBankMapCollection) Find(insCode string) (cb *model.CfcaBankMap, err error) {

	// get from cache
	o, found := cfcaCache.Get(insCode)
	if found {
		cb = o.(*model.CfcaBankMap)
		return cb, nil
	}

	cb = new(model.CfcaBankMap)
	q := bson.M{"insCode": insCode}

	err = database.C(c.name).Find(q).One(cb)

	if err != nil {
		log.Errorf("'Find CfcaBankMap ERROR!' Error message is: %s\n; Condition is: %+v", err.Error(), q)
		return nil, err
	}

	// save
	cfcaCache.Set(insCode, cb, cache.NoExpiration)

	return
}
