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

var database *mgo.Database

func init() {
	session, err := mgo.Dial(host)
	if err != nil {
		g.Fatal("unable connect to mongodb server ", err)
	}

	session.SetMode(mgo.Monotonic, true)
	database = session.DB(dbname)

	g.Info("connected to mongodb %s database %s", host, dbname)
}
