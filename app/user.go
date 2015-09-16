package app

import (
	"crypto/sha1"
	"fmt"
	"sort"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

type user struct{}

var User user

func (u *user) register(userName, password, transtime, sign string) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s,sign=%s", userName, password, transtime, sign)
	//验签
	strs := sort.StringSlice{userName, password, transtime, model.KEY}
	sort.Strings(strs)
	var str string
	for _, s := range strs {
		str += s
	}
	value := sha1.Sum([]byte(str))
	valueStr := fmt.Sprintf("%x", value)
	log.Debugf("sign(%s)=%s", str, valueStr)
	if sign != valueStr {
		log.Errorf("check signature err")
		return model.NewAppResult(model.FAIL, model.SIGN_FAIL)
	}

	user := &model.AppUser{
		UserName: userName,
		Password: password,
	}

	err := mongo.AppUserCol.Upsert(user)
	if err != nil {

	}

	return result
}
