package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/CardInfoLink/quickpay/config"
	"github.com/omigo/log"
)

var database *mgo.Database

// Connect 程序启动时，或者，单元测试前，先连接到 MongoDB 数据库
func init() {
	host := config.GetValue("mongo", "host")
	port := config.GetValue("mongo", "port")
	dbname := config.GetValue("mongo", "dbname")

	addr := host + ":" + port
	session, err := mgo.Dial(addr)
	if err != nil {
		log.Fatalf("unable connect to mongodb server %s", err)
	}

	session.SetMode(mgo.Monotonic, true)
	database = session.DB(dbname)

	log.Infof("connected to mongodb host `%s` and database `%s`", addr, dbname)
}
