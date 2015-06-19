package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
)

var database *mgo.Database

// Connect 程序启动时，或者，单元测试前，先连接到 MongoDB 数据库
func init() {
	url := goconf.GetValue("mongo", "url")
	dbname := goconf.GetValue("mongo", "db")

	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("unable connect to mongodb server %s", err)
	}

	session.SetMode(mgo.Eventual, true) //需要指定为Eventual
	session.SetSafe(&mgo.Safe{})

	database = session.DB(dbname)

	log.Infof("connected to mongodb host `%s` and database `%s`", url, dbname)
}
