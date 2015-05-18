// Package check 定时检查数据库 checkAndNotify 文档所有记录，如果发生变化，
// 通知响应业务模块，重新加载缓存，或者重新构建业务逻辑
package check

import (
	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

func DoCheck() {
	// do something

	go checking()
}

func checking() {
	// tick ...
	log.Debug("wait to check conf...")
	tick := time.Tick(30 * time.Second)

	for {
		<-tick
		// 取 checkAndNotify 文档所有记录
		cans, err := mongo.NotifyColl.GetAll()
		if err != nil {
			log.Errorf("fail to load all CheckAndNotify info : %s", err)
			continue
		}

		changes := make([]*model.CheckAndNotify, 0, len(cans))
		// 遍历
		for _, v := range cans {
			// TODO read from config
			if v.App1Tag != v.CurTag {
				// 放到slice里
				changes = append(changes, v)
			}
		}

		// 处理
		go notifying(changes)
	}
}

func notifying(changes []*model.CheckAndNotify) {

	for _, v := range changes {
		// 获得缓存
		c := cache.Client.Get(v.BizType)

		// 清空
		if c != nil {
			log.Debugf("clear %s cache...", v.BizType)
			c.Clear()
		}

		// cardBin需要重建树
		if v.BizType == "cardBin" {
			// TODO
			log.Infof("cardBin had been updated (%s -> %s), begin to rebuild the cardBin tree ", v.CurTag, v.App1Tag)
			err := core.ReBuildTree()
			if err != nil {
				log.Error(err)
				continue
			}
		}

		// 成功，更新当前版本
		// TODO
		v.PrevTag = v.CurTag
		v.CurTag = v.App1Tag
		err := mongo.NotifyColl.Update(v)
		if err != nil {
			log.Errorf("fail to update NotifyColl(%+v) : %s", v, err)
		}
	}

}
