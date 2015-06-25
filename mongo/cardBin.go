package mongo

import (
	"fmt"

	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

// CardBinColl 卡Bin Collection
var CardBinColl = cardBinCollection{"cardBin"}

var cardBinCache = cache.New(model.Cache_CardBin)

type cardBinCollection struct {
	name string
}

// Find 根据卡长度查找卡BIN列表
// TODO 查找的时候应该匹配卡bin的优先级字段
func (c *cardBinCollection) Find(cardBin string, length int) (cb *model.CardBin, err error) {

	k := cardBin + fmt.Sprintf("%s", length)
	o, found := cardBinCache.Get(k)
	if found {
		cb = o.(*model.CardBin)
		return cb, nil
	}

	cb = new(model.CardBin)
	q := bson.M{
		"bin":     cardBin,
		"cardLen": length,
	}
	err = database.C(c.name).Find(q).One(cb)

	if err != nil {
		return nil, err
	}

	// save cache
	cardBinCache.Set(k, cb, cache.NoExpiration)

	return
}

// LoadAll 加载所有卡bin，并刷新缓存
func (c *cardBinCollection) LoadAll() ([]*model.CardBin, error) {

	var cardBins []*model.CardBin
	err := database.C(c.name).Find(nil).All(&cardBins)

	if err != nil {
		return nil, err
	}

	// 清空缓存
	cardBinCache.Clear()
	// 重新初始化
	for _, v := range cardBins {
		k := v.Bin + fmt.Sprintf("%s", v.CardLen)
		cardBinCache.Set(k, v, cache.NoExpiration)
	}

	return cardBins, err
}

/*  only use for import/update cardBin data from csv   */

func (c *cardBinCollection) Upsert(cb *model.CardBin) error {
	s := bson.M{
		"bin":     cb.Bin,
		"cardLen": cb.CardLen,
	}
	_, err := database.C(c.name).Upsert(s, cb)
	return err
}

func (c *cardBinCollection) Drop() error {
	return database.C(c.name).DropCollection()
}
