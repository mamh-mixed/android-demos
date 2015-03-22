package core

import (
	"github.com/CardInfoLink/quickpay/channel/cfca"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/g"
	"time"
)

// init 开启任务routine
// func init() {
// 	go ProcessTransSettle()
// }

// ProcessTransSettle 清分
func ProcessTransSettle() {

	//TODO 暂时先每天凌晨将交易信息拷贝到清分表里
	// 距离0点的时间
	dis, err := tools.TimeToGiven("00:00:00")
	if err != nil {
		g.Error("fail to get time second by given %s", err)
		return
	}
	c := make(chan bool)
	time.AfterFunc(time.Duration(dis)*time.Second, func() {
		//for test
		tick := time.Tick(1 * time.Second)
		boom := time.After(5 * time.Second)
		for {
			select {
			case <-tick:
				doTransSett()
			//for test
			case <-boom:
				c <- true
				g.Debug("boom break %s", "boom")
				return
			}
		}

	})
	<-c

}

func doTransSett() {
	layout := "2006-01-02"
	now := time.Now()
	//查找昨天的交易
	d, _ := time.ParseDuration("-24h")
	yestday := now.Add(d)
	g.Debug("yesterday : %s", yestday.Format(layout))
	trans, err := mongo.TransColl.FindByTime(yestday.Format(layout))
	if err != nil {
		g.Error("find trans fail : %s", err)
		return
	}
	for _, v := range trans {
		//暂时是假勾兑
		sett := &model.TransSett{
			Tran:     v,
			SettFlag: 1,
			SettDate: now.Format("2006-01-02 15:04:05"),
			SettAmt:  v.TransAmt / 10,
			MerFee:   v.TransAmt / 10,
		}
		if err := mongo.TransSettColl.Add(sett); err != nil {
			g.Error("add trans sett fail : %s", err)
		}

	}

	//TODO 勾兑
	//cfca.ProcessTransChecking(be)
}
