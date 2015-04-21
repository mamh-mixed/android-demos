// 启动一个 Goruntine，检测数据库数据变化，重新初始化
package conf

import (
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

var defaultVn string

func CheckConf() {
	// cardBin
	version, err := mongo.VersionColl.Find("cardBin")
	if err != nil {
		log.Panicf("fail to load cardBin version : %s ", err)
	}
	defaultVn = version.Vn
	go doCheckConf()
}

func doCheckConf() {

	tick := time.Tick(30 * time.Second)

	for {
		select {
		case <-tick:
			o, _ := mongo.VersionColl.Find("cardBin")
			if o.Vn != defaultVn {
				log.Infof("cardBin had been updated (%s -> %s), begin to rebuild the cardBin tree ", defaultVn, o.Vn)
				// ... rebuild the tree
				err := core.ReBuildTree()
				if err != nil {
					log.Error(err)
					continue
				}
				defaultVn = o.Vn
			}
		}
	}
}
