package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type appUserCollection struct {
	name string
}

var AppUserCol = appUserCollection{"appUser"}

func (col *appUserCollection) Upsert(user *model.AppUser) (err error) {

	bo := bson.M{
		"username": user.UserName,
	}
	_, err = database.C(col.name).Upsert(bo, user)
	if err != nil {
		log.Errorf("upsert user err,%s", err)
		return err
	}
	return nil
}

// func (col *appUserCollection) FindOne(userName string) (err error){
//     bo:=bson.M
// }
