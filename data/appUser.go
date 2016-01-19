package data

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"time"
	// "gopkg.in/mgo.v2/bson"
	// "math"
	// "net/http"
	// "strconv"
)

type appUser struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Activate string `bson:"activate"`
	Clientid string `bson:"clientid"`
	Limit    string `bson:"limit"`
}

// AsyncAppUser 同步旧系统app用户信息
func AsyncAppUser() error {
	appUsers, err := readAppUserFromOldDB()
	if err != nil {
		return err
	}

	for _, app := range appUsers {
		user := &model.AppUser{}
		user.UserName = app.Username
		user.Password = app.Password
		user.Activate = app.Activate
		user.MerId = app.Clientid
		user.Limit = app.Limit
		user.Remark = "old_system_appUsers"
		user.RegisterFrom = model.SelfRegister
		user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		user.UpdateTime = user.CreateTime
		err = mongo.AppUserCol.Upsert(user)
		if err != nil {
			return err
		}

		if app.Clientid != "" {
			m, err := mongo.MerchantColl.Find(app.Clientid)
			if err != nil {
				log.Errorf("未找到商户号 %s", app.Clientid)
				continue
			}
			if m.RefundType != model.CurrentDayRefund {
				// 将商户设置为只能当天退款
				m.RefundType = model.CurrentDayRefund
				mongo.MerchantColl.Update(m)
			}
		}
	}

	log.Infof("成功导入 %d 条用户信息", len(appUsers))
	return nil
}

type settInfo struct {
	Clientid  string `bson:"clientid"`
	BankOpen  string `bson:"bank_open"`
	Payee     string `bson:"payee"`
	PayeeCard string `bson:"payee_card"`
	PhoneNum  string `bson:"phone_num"`
}

// AsyncSettInfo 同步清算信息
func AsyncSettInfo() error {
	ss, err := readSettInfoFromOldDB()
	if err != nil {
		return err
	}

	for _, s := range ss {
		m, err := mongo.MerchantColl.Find(s.Clientid)
		if err != nil {
			log.Warnf("找不到商户(%s)，请检查。", s.Clientid)
			continue
		}

		// 开户行
		if m.Detail.OpenBankName == "" {
			m.Detail.OpenBankName = s.BankOpen
		} else {
			// check
			if m.Detail.OpenBankName != s.BankOpen {
				log.Warnf("新旧系统数据不一致，商户(%s)，开户行名称，新=%s，旧=%s", s.Clientid, m.Detail.OpenBankName, s.BankOpen)
			}
		}

		// 手机
		if m.Detail.ContactTel == "" {
			m.Detail.ContactTel = s.PhoneNum
		} else {
			if m.Detail.ContactTel != s.PhoneNum {
				log.Warnf("新旧系统数据不一致，商户(%s)，手机，新=%s，旧=%s", s.Clientid, m.Detail.ContactTel, s.PhoneNum)
			}
		}

		// 账户名称
		if m.Detail.AcctName == "" {
			m.Detail.AcctName = s.Payee
		} else {
			if m.Detail.AcctName != s.Payee {
				log.Warnf("新旧系统数据不一致，商户(%s)，账户名称，新=%s，旧=%s", s.Clientid, m.Detail.AcctName, s.Payee)
			}
		}

		// 帐号
		if m.Detail.AcctNum == "" {
			m.Detail.AcctNum = s.PayeeCard
		} else {
			if m.Detail.AcctNum != s.PayeeCard {
				log.Warnf("新旧系统数据不一致，商户(%s)，银行帐号，新=%s，旧=%s", s.Clientid, m.Detail.AcctNum, s.PayeeCard)
			}
		}
		mongo.MerchantColl.Update(m)
	}
	return nil
}

func readSettInfoFromOldDB() ([]settInfo, error) {
	var settInfos []settInfo
	err := saomaDB.C("cloudCashierAccount").Find(nil).All(&settInfos)
	return settInfos, err
}

func readAppUserFromOldDB() ([]appUser, error) {
	var appUsers []appUser
	err := saomaDB.C("user").Find(nil).All(&appUsers)
	return appUsers, err
}
