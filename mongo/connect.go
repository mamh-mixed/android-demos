package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/omigo/g"
)

const (
	host = "121.41.85.237"
	// host   = "127.0.0.1"
	dbname = "quickpay"
)

type mgodb struct {
	database *mgo.Database
	// 应答码表
	respCode *mgo.Collection
	// for test
	people *mgo.Collection
	// 卡bin
	cardBin *mgo.Collection
	// 绑定关系集合
	bindingRelation *mgo.Collection
	// 路由策略集合
	routerPolicy *mgo.Collection
}

var db mgodb

func init() {
	session, err := mgo.Dial(host)
	if err != nil {
		g.Fatal("unable connect to mongodb server ", err)
	}

	session.SetMode(mgo.Monotonic, true)
	database := session.DB(dbname)

	g.Info("connected to mongodb %s database %s", host, dbname)

	//init
	db = mgodb{
		database:        database,
		respCode:        database.C("respCode"),
		people:          database.C("people"),
		cardBin:         database.C("cardBin"),
		bindingRelation: database.C("bindingRelation"),
		routerPolicy:    database.C("routerPolicy"),
	}
}
