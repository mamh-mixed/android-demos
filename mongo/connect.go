package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

const (
	// 配置多个连接地址，取第一个可用地址
	// 变量 MONGO_PORT_27017_TCP_ADDR 是为了 Docker 环境下自动取得 MongoDB 地址
	host   = "$MONGO_PORT_27017_TCP_ADDR | 121.41.85.237 | 121.40.86.222 "
	port   = "$MONGO_PORT_27017_TCP_PORT | 27017"
	dbname = "quickpay"
)

var database *mgo.Database

// Connect 程序启动时，或者，单元测试前，先连接到 MongoDB 数据库
func Connect() {
	favHost := tools.FirstExistValue(host)
	favPort := tools.FirstExistValue(port)

	addr := favHost + ":" + favPort
	session, err := mgo.Dial(addr)
	if err != nil {
		log.Fatalf("unable connect to mongodb server %s", err)
	}

	session.SetMode(mgo.Monotonic, true)
	database = session.DB(dbname)

	log.Infof("connected to mongodb host `%s` and database `%s`", addr, dbname)
}
