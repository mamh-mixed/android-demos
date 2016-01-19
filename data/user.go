package data

import (
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

func AddUserFromCSV(path string) error {
	users, err := readUserCSV(path)
	if err != nil {
		return err
	}

	log.Infof("add users len=%d", len(users))
	log.Debugf("%+v", users)

	for _, u := range users {
		err := mongo.UserColl.Add(&u)
		if err != nil {
			log.Errorf("新增用户帐号 重复 userName=%s", u.UserName)
			continue
		}
	}
	return nil
}
