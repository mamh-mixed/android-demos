package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/omigo/log"
)

const (
	host   = "121.41.85.237"
	dbname = "quickpay"
)

var database *mgo.Database

// Connect 程序启动时，或者，单元测试前，先连接到 MongoDB 数据库
func Connect() {
	session, err := mgo.Dial(host)
	if err != nil {
		log.Fatalf("unable connect to mongodb server %s", err)
	}

	session.SetMode(mgo.Monotonic, true)
	database = session.DB(dbname)

	log.Infof("connected to mongodb %s database %s", host, dbname)
}
