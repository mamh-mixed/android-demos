package mongo

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/CardInfoLink/quickpay/goconf"
)

var masterDB *mgo.Database
var database *mgo.Database

// Connect 程序启动时，或者，单元测试前，先连接到 MongoDB 数据库
func init() {
	url := goconf.Config.Mongo.URL
	dbname := goconf.Config.Mongo.DB

	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("unable connect to mongo %s: %s\n", url, err)
		os.Exit(1)
	}

	// 连接master的session
	// Strong session 强制要求与master连接
	// 当数据库master移位的时候，Strong session 可能会发生EOF
	// 需要重新Refresh
	masterSession := session.Copy()
	masterSession.SetMode(mgo.Strong, true) // 从master读写，处理对于实时性要求高的操作
	masterSession.SetSafe(&mgo.Safe{})

	session.SetMode(mgo.Eventual, true) // 最终一致性即可，读写分离
	session.SetSafe(&mgo.Safe{})

	database = session.DB(dbname)
	masterDB = masterSession.DB(dbname)

	// 不能在日志中出现数据库密码
	// fmt.Println("connected to mongodb host `%s` and database `%s`", url, dbname)
	fmt.Printf("connected to mongodb host `%s` and database `%s`\n", "***", dbname)
}
