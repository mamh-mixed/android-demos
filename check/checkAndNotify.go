// Package check 定时检查数据库 checkAndNotify 文档所有记录，如果发生变化，
// 通知响应业务模块，重新加载缓存，或者重新构建业务逻辑
package check

import (
	"time"

	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

var appId = util.Hostname

// DoCheck 定时检查配置
func DoCheck() {
	go check()
}

func check() {
	log.Debug("wait to check conf...")
	tick := time.Tick(30 * time.Second)

	for {
		<-tick

		// 取 checkAndNotify 文档所有记录
		notifies, err := mongo.NotifyColl.GetAll()
		if err != nil {
			log.Errorf("fail to load all CheckAndNotify info : %s", err)
			continue
		}

		changes := make([]*model.CheckAndNotify, 0, len(notifies))
		// 遍历
		for _, v := range notifies {
			// TODO 利用反射，得到对应字段
			switch appId {
			case "app1":
				if v.App1Tag != v.CurTag {
					changes = append(changes, v)
				}
			case "app2":
				if v.App2Tag != v.CurTag {
					changes = append(changes, v)
				}
			case "app3":
				if v.App1Tag != v.CurTag {
					changes = append(changes, v)
				}
			case "app4":
				if v.App1Tag != v.CurTag {
					changes = append(changes, v)
				}
			}
		}

		// 处理
		go notify(changes)
	}
}

func notify(changes []*model.CheckAndNotify) {

	for _, v := range changes {
		switch v.BizType {
		case model.Cache_CardBin:
			// TODO
			log.Infof("cardBin had been updated (%s -> %s), begin to rebuild the cardBin tree ", v.App1Tag, v.CurTag)
			err := core.ReBuildTree()
			if err != nil {
				log.Error(err)
				continue
			}
		case model.Cache_Merchant, model.Cache_ChanMer, model.Cache_CfcaBankMap, model.Cache_CfcaMerRSAPrivKey,
			model.Cache_RespCode:
			// 获得缓存
			c, ok := cache.Client.Get(v.BizType)
			// 清空
			if ok {
				c.Clear()
				log.Infof("clear %s cache...", v.BizType)
			}
		default:
			log.Errorf("unimplement business type %s", v.BizType)
		}

		// 成功，更新当前应用的版本，不要更新其他值
		// TODO 利用反射，得到对应字段
		switch appId {
		case "app1":
			v.App1Tag = v.CurTag
		case "app2":
			v.App2Tag = v.CurTag
		case "app3":
			v.App3Tag = v.CurTag
		case "app4":
			v.App4Tag = v.CurTag
		}

		err := mongo.NotifyColl.Update(v)
		if err != nil {
			log.Errorf("fail to update NotifyColl(%+v) : %s", v, err)
		}
	}

}
