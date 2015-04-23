package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

// CardBinColl 卡Bin Collection
var CardBinColl = cardBinCollection{"cardBin"}

type cardBinCollection struct {
	name string
}

// Find 根据卡长度查找卡BIN列表
// TODO 查找的时候应该匹配卡bin的优先级字段
func (c *cardBinCollection) Find(cardBin string, length int) (cb *model.CardBin, err error) {
	cb = new(model.CardBin)
	q := bson.M{
		"bin":     cardBin,
		"cardLen": length,
	}
	err = database.C(c.name).Find(q).One(cb)

	return
}

// LoadAll 加载所有卡bin
func (c *cardBinCollection) LoadAll() ([]*model.CardBin, error) {
	var cardBins []*model.CardBin
	err := database.C(c.name).Find(nil).All(&cardBins)
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
