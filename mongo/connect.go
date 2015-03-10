package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/omigo/g"
)

const (
	host   = "121.41.85.237"
	dbname = "quickpay"
)

var db *mgo.Database

func init() {
	session, err := mgo.Dial(host)
	if err != nil {
		g.Fatal("unable connect to mongodb server ", err)
	}

	session.SetMode(mgo.Monotonic, true)
	db = session.DB(dbname)

	g.Info("connected to mongodb %s database %s", host, dbname)
}
