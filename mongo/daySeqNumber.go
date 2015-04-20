package mongo

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type daySeqNumCollecion struct {
	name string
}

const (
	dayMaxSN = 1000000
)

var DaySNColl = snCollecion{"counter"}

// GetDaySN 返回一个当天唯一的六位数字
func (c *snCollecion) GetDaySN(merId, termId string) string {
	q := bson.M{
		"merId":     merId,
		"termId":    termId,
		"$isolated": 1,
	}
	u := bson.M{
		"$inc": bson.M{
			"sn": 1,
		},
	}
	change := mgo.Change{
		Update:    u,
		ReturnNew: true,
	}

	var sn = new(model.DaySN)
	_, err := database.C(c.name).Find(q).Apply(change, &sn)

	if err == nil {
		return fmt.Sprintf("%06d", sn.Sn%dayMaxSN)
	}

	// 如果没找到，添加文档，并返回初始值
	if err.Error() == "not found" {
		ds := model.DaySN{
			MerId:  merId,
			TermId: termId,
			Sn:     0,
		}
		err = database.C(c.name).Insert(ds)
		if err != nil {
			log.Errorf("method 'GetDaySN' insert daySN collection error: %s\n", err)
			return ""
		}

		return fmt.Sprintf("%06d", 0)
	}

	log.Errorf("GetDaySN error: %s\n", err)
	return ""
}
