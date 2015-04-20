package mongo

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type snCollecion struct {
	name string
}

const (
	sysMaxSN = 1000000000000
)

var SnColl = snCollecion{"counter"}

// GetSysSn 返回一个系统唯一的只包含数字和字母的12位字符串
func (c *snCollecion) GetSysSN() string {
	q := bson.M{
		"key":       "sysSN",
		"$isolated": 1,
	}
	change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"value": 1,
			},
		},
		ReturnNew: true,
	}

	var sn = new(model.SN)
	_, err := database.C(c.name).Find(q).Apply(change, &sn)
	if err != nil {
		log.Errorf("Find and modify error: %s\n", err)
		return ""
	}

	return fmt.Sprintf("%012d", sn.Value%sysMaxSN)
}
