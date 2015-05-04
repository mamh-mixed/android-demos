// 启动一个 Goruntine，检测数据库数据变化，重新初始化
package pay

import (
	"time"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

var defaultVn string

// CheckConf 定时检查数据库配置，如果 卡Bin 变化，重建搜索树
func CheckConf() {
	// cardBin
	version, err := mongo.VersionColl.FindOne("cardBin")
	if err != nil {
		log.Errorf("fail to load cardBin version : %s ", err)
		return
	}
	defaultVn = version.Vn
	go doCheckConf()
}

func doCheckConf() {

	tick := time.Tick(30 * time.Second)

	for {
		<-tick

		o, err := mongo.VersionColl.FindOne("cardBin")
		if err != nil {
			log.Errorf("fail to find cardBin version: %s", err)
			continue
		}
		if o.Vn != defaultVn {
			log.Infof("cardBin had been updated (%s -> %s), begin to rebuild the cardBin tree ", defaultVn, o.Vn)
			// ... rebuild the tree
			err = core.ReBuildTree()
			if err != nil {
				log.Error(err)
				continue
			}
			defaultVn = o.Vn
		}
	}
}
