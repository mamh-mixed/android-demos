package settle

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

// RefreshSpTransSett 重新生成数据
func RefreshSpTransSett(date string) (err error) {

	// 删除交易
	err = mongo.SpTransSettColl.BatchRemove(date)
	if err != nil {
		return err
	}

	// 导数据
	err = DoSpTransSett(date)
	if err != nil {
		return err
	}

	// TODO:勾兑
	return nil
}

// DoSpTransSett 扫码支付清算
func DoSpTransSett(date string) (err error) {

	var ts []*model.Trans
	for i := 0; i < 3; i++ {
		ts, err = mongo.SpTransColl.FindToSett(date)
		if err != nil {
			log.Errorf("find trans error: %s, %d times", err, i+1)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	if len(ts) == 0 {
		return nil
	}

	// 缓存
	var chanMersMap = make(map[string]*model.ChanMer)
	// var routesMap = make(map[string]*model.RouterPolicy)

	var transSetts []model.TransSett
	for _, t := range ts {

		// 得到渠道商户
		var cm *model.ChanMer
		if c, ok := chanMersMap[t.ChanCode+t.ChanMerId]; ok {
			cm = c
		} else {
			cm, err = mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
			if err != nil {
				log.Errorf("chanMer should be found ,but find error: %s", err)
				continue
			}
		}

		// 路由
		// var route *model.RouterPolicy
		// if r, ok := routesMap[t.ChanCode+t.MerId]; ok {
		// 	route = r
		// } else {
		// 	route, err = mongo.RouterPolicyColl.Find(t.MerId, t.ChanCode)
		// 	if err != nil {
		// 		log.Errorf("routerPolicy should be found ,but find error: %s", err)
		// 		continue
		// 	}
		// }

		transSett := model.TransSett{}
		transSett.Trans = *t
		transSett.SettDate = date
		transSett.SettRole = t.SettRole
		transSett.BlendType = CHAN_LESS

		// 计算商户手费率
		if t.TransType == model.PayTrans {
			// 支付交易
			transSett.MerFee = t.Fee
			transSett.MerSettAmt = t.TransAmt - t.Fee
			// TODO...
		} else {
			// 逆向交易
			switch cm.SchemeType {
			case 0, 1:
				// 原路返回，默认就是，不做特殊处理
				transSett.MerFee = t.Fee
			case 2:
				// 渠道不退手续费，也就是逆向交易的时候没有手续费
				transSett.MerFee = 0
			}
		}

		// TODO:计算机构手续费
		// TODO:计算机渠道续费

		transSetts = append(transSetts, transSett)
	}

	// TODO:报表引用数据流整理

	err = mongo.SpTransSettColl.BatchAdd(transSetts)
	if err != nil {
		log.Errorf("batch add transSett error: %s", err)
	}

	return err
}
