package settle

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"math"
	"time"
)

// 勾兑状态
const (
	MATCH     = 0
	CHAN_MORE = 1
	CHAN_LESS = 2
	AMT_ERROR = 3
	FEE_ERROR = 4
)

// Settle 清算接口
type Settle interface {
	ProcessDuration() time.Duration // 何时可以执行
	Reconciliation(date string)     // 勾兑过程
}

// 需要清算
var needSettles []Settle

// DoSettle 勾兑
func DoSettle(date string, immediately bool) {
	for _, ns := range needSettles {
		var d time.Duration
		// 马上执行，用于网页发起
		if immediately {
			d = 0
		} else {
			d = ns.ProcessDuration()
		}

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
	err = DoSpTransSett(date, true)
	if err != nil {
		return err
	}

	return nil
}

// DoSpTransSett 扫码支付清算
func DoSpTransSett(date string, immediately bool) (err error) {

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

	log.Infof("%d trans prepare to handle", len(ts))

	// 缓存
	var chanMersMap = make(map[string]*model.ChanMer)
	var agentsMap = make(map[string]*model.Agent)
	// var routesMap = make(map[string]*model.RouterPolicy)

	var transSetts []model.TransSett
	for _, t := range ts {

		// 得到渠道商户
		var cm *model.ChanMer
		if c, ok := chanMersMap[t.ChanCode+t.ChanMerId]; ok {
			cm = c
			chanMersMap[t.ChanCode+t.ChanMerId] = cm
		} else {
			cm, err = mongo.ChanMerColl.Find(t.ChanCode, t.ChanMerId)
			if err != nil {
				log.Errorf("chanMer should be found ,but find error: %s", err)
				continue
			}
		}

		// 机构
		var agent *model.Agent
		if a, ok := agentsMap[t.AgentCode]; ok {
			agent = a
		} else {
			a, err = mongo.AgentColl.Find(t.AgentCode)
			if err != nil {
				log.Errorf("trans skip: agent not found. merId=%s, orderNum=%s", t.MerId, t.OrderNum)
				continue
			}
			agent = a
			agentsMap[t.AgentCode] = a
		}

		// TODO DELETE:
		// 兼容逆向交易没有存渠道订单号
		// 勾兑时用到
		if t.TransType != model.PayTrans {
			if t.ChanOrderNum == "" {
				// 查找原订单
				ot, err := mongo.SpTransColl.FindOne(t.MerId, t.OrigOrderNum)
				if err == nil {
					t.ChanOrderNum = ot.ChanOrderNum
				}
			}
		}

		transSett := model.TransSett{}
		transSett.Trans = *t
		transSett.SettDate = date
		transSett.SettRole = t.SettRole
		transSett.BlendType = CHAN_LESS // 默认是渠道少清的

		// 商户手续费
		transSett.MerFee = t.Fee
		transSett.MerSettAmt = t.TransAmt - t.Fee

		// 计算讯联、机构手续费
		switch t.ChanCode {
		case channel.ChanCodeAlipay:
			transSett.AcqFee = 0 // 支付宝默认0
			transSett.InsFee = int64(math.Floor(float64(t.TransAmt)*agent.AlpCost + 0.5))
		case channel.ChanCodeWeixin:
			transSett.AcqFee = int64(math.Floor(float64(t.TransAmt)*0.003 + 0.5)) // 微信默认0.003
			transSett.InsFee = int64(math.Floor(float64(t.TransAmt)*agent.WxpCost + 0.5))
		}

		// 逆向交易
		if t.TransType != model.PayTrans {
			switch cm.SchemeType {
			case 0, 1:
				// 原路返回，默认就是，不做特殊处理
			case 2:
				// 渠道不退手续费，也就是逆向交易的时候没有手续费
				transSett.MerFee = 0
				transSett.InsFee = 0
			}
		}

		// 讯联、机构清算金额
		transSett.InsSettAmt = t.TransAmt - transSett.InsFee
		transSett.AcqSettAmt = t.TransAmt - transSett.AcqFee
		transSetts = append(transSetts, transSett)
	}

	err = mongo.SpTransSettColl.BatchAdd(transSetts)
	if err != nil {
		log.Errorf("batch add transSett error: %s", err)
		return err
	}

	// 报表引用数据流整理
	// go SpSettReport(date)
	// go SpReconciliatReport(date, transSetts...)

	// 进行勾兑
	DoSettle(date, immediately)
	log.Info("Do SpTransSett success... gen report... do settle...")
	return err
}
