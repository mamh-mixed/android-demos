package mongo

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type daySeqNumCollecion struct {
	name string
}

const (
	dayMaxSN = 1000000
)

var DaySNColl = snCollecion{"daySN"}

// GetDaySN 返回一个当天唯一的六位数字
func (c *snCollecion) GetDaySN(merId, termId string) string {
	u1 := bson.M{
		"merId":     merId,
		"termId":    termId,
		"$isolated": 1,
	}
	u2 := bson.M{
		"$inc": bson.M{
			"sn": 1,
		},
	}
	// 先查找，后更新
	var sn = new(model.DaySN)
	err := database.C(c.name).Find(u1).One(&sn)

	// 如果成功找到，加1，并返回
	if err == nil {
		err = database.C(c.name).Update(u1, u2)
		if err != nil {
			log.Errorf("method 'GetDaySN' update daySN collection error: %s", err)
		}
		return fmt.Sprintf("%06d", sn.Sn%dayMaxSN)
	}

	// 如果没找到，添加文档，并返回初始值
	if err.Error() == "not found" {
		ds := model.DaySN{
			MerId:  merId,
			TermId: termId,
			Sn:     1,
		}
		err = database.C(c.name).Insert(ds)
		if err != nil {
			log.Errorf("method 'GetDaySN' insert daySN collection error: %s", err)
			return ""
		}

		return fmt.Sprintf("%06d", 0)
	}

	// 其他错误
	log.Errorf("method 'GetDaySN' find daySN collection error: %s", err)
	return ""
}
