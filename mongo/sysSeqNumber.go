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
	sysMaxSN = 1000000000000
)

var SnColl = snCollecion{"sn"}

// GetSysSn 返回一个系统唯一的只包含数字和字母的12位字符串
func (c *snCollecion) GetSysSN() string {
	u1 := bson.M{
		"key":       "sysSN",
		"$isolated": 1,
	}
	u2 := bson.M{
		"$inc": bson.M{
			"value": 1,
		},
	}
	// 先查找，后更新
	var sn = new(model.SN)
	err := database.C(c.name).Find(u1).One(&sn)
	if err != nil {
		log.Errorf("method 'GetSysSn' find SN collection error: %s", err)
	}

	err = database.C(c.name).Update(u1, u2)
	if err != nil {
		log.Errorf("method 'GetSysSn' update SN collection error: %s", err)
	}

	return fmt.Sprintf("%012d", sn.Value%sysMaxSN)
}
