package mongo

import (
	"fmt"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type snCollecion struct {
	name string
}

const (
	sysMaxSN = 1000000000000
	dayMaxSN = 1000000
)

var SnColl = snCollecion{"counter"}

// GetSysSn 返回一个系统唯一的只包含数字和字母的12位字符串
func (c *snCollecion) GetSysSN() string {
	q := bson.M{
		"type":      "sys",
		"$isolated": 1,
	}
	change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"sn": 1,
			},
		},
		ReturnNew: true,
		Upsert:    true,
	}

	var sn = new(model.SN)
	_, err := database.C(c.name).Find(q).Apply(change, &sn)
	if err != nil {
		log.Errorf("Find and modify error: %s\n", err)
		return ""
	}

	return fmt.Sprintf("%012d", sn.Sn%sysMaxSN)
}

// GetDaySN 返回一个指定商户和终端当天唯一的六位数字
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

	var sn = new(model.SN)
	_, err := database.C(c.name).Find(q).Apply(change, &sn)

	if err == nil {
		return fmt.Sprintf("%06d", sn.Sn%dayMaxSN)
	}

	// 如果没找到，添加文档，并返回初始值
	if err.Error() == "not found" {
		ds := model.SN{
			Type:   "day",
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
