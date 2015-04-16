package mongo

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type snCollecion struct {
	name string
}

const (
	dayMaxSN = 999999
)

var SnColl = snCollecion{"sn"}

// GetDaySN 返回一个当天唯一的六位数字
func (c *snCollecion) GetDaySN() string {
	u1 := bson.M{
		"key":       "daySN",
		"$isolated": 1,
	}
	u2 := bson.M{
		"$inc": bson.M{
			"value": 1,
		},
	}
	// 先查找，后更新
	var sn = new(model.SN)
	err := database.C(c.name).Find(bson.M{"key": "daySN"}).One(&sn)
	if err != nil {
		log.Errorf("method 'GetDaySN' find SN collection error: %s", err)
	}

	if sn.Value > dayMaxSN {
		err = database.C(c.name).Update(u1, bson.M{"key": "daySN", "value": 2})
		if err != nil {
			log.Errorf("method 'GetDaySN' update SN collection error: %s", err)
		}
		return fmt.Sprintf("%06d", 1)
	}

	err = database.C(c.name).Update(u1, u2)
	if err != nil {
		log.Errorf("method 'GetDaySN' update SN collection error: %s", err)
	}

	return fmt.Sprintf("%06d", sn.Value)
}

// GetSysSn 返回一个系统唯一的只包含数字和字母的12位字符串
func (c *snCollecion) GetSysSn() string {
	return ""
}
