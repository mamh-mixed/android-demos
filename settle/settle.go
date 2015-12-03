package settle

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

// Settle 清算接口
type Settle interface {
	ProcessDuration() time.Duration // 何时可以执行
	Reconciliation(date string)     // 勾兑过程
}

// 需要清算
var needSettles []Settle

// DoSettle 勾兑
func DoSettle(date string) {
	for _, ns := range needSettles {
		d := ns.ProcessDuration()
		switch d {
		case 0:
			go ns.Reconciliation(date)
		case -1:
			// ignore
		default:
			time.AfterFunc(d, func() {
				ns.Reconciliation(date)
			})
		}
	}
}

// RefreshSpTransSett 重新生成数据
func RefreshSpTransSett(date string) (err error) {

	// 删除交易
	err = mongo.SpTransSettColl.BatchRemove(date)
	if err != nil {
		return err
	}

	// 重新执行
	err = DoSpTransSett(date)
	if err != nil {
		return err
	}

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
		transSett.BlendType = CHAN_LESS // 默认是渠道少清的

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
		return err
	}

	// 进行勾兑
	DoSettle(date)

	return err
}
