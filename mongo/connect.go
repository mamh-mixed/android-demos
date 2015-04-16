package mongo

import (
	"gopkg.in/mgo.v2"

	"github.com/omigo/log"
)

const (
	host = "tdd.ipay.so"
	// host   = "127.0.0.1"
	dbname = "quickpay"
)

var database *mgo.Database

func init() {
	session, err := mgo.Dial(host)
	if err != nil {
		log.Fatalf("unable connect to mongodb server ", err)
	}

	session.SetMode(mgo.Monotonic, true)
	database = session.DB(dbname)

	log.Infof("connected to mongodb %s database %s", host, dbname)

	buildTree()
}
