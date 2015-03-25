package core

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/g"
	"time"
)

// init 开启任务routine
func init() {
	g.Debug("wait to process transSett method")
	go ProcessTransSettle()
}

// ProcessTransSettle 清分
func ProcessTransSettle() {

	// 暂时先每天凌晨将交易信息拷贝到清分表里
	// 距离0点的时间
	dis, err := tools.TimeToGiven("00:30:00")
	if err != nil {
		g.Error("fail to get time second by given %s", err)
		return
	}
	c := make(chan bool)
	time.AfterFunc(time.Duration(dis)*time.Second, func() {
		// time.AfterFunc(10*time.Second, func() {
		//到点时先执行一次
		doTransSett()
		//24小时后执行
		tick := time.Tick(24 * time.Hour)
		// boom := time.After(5 * time.Second)
		time.Since(time.Now())
		for {
			select {
			case <-tick:
				// g.Debug("tick ... %s", "boom")
				doTransSett()
				//for test
				// case <-boom:
				// 	c <- true
				// 	g.Debug("boom break %s", "boom")
				// 	return
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
	// g.Debug("yesterday data : %+v", trans)
	// 计算费率
	for _, v := range trans {
		// 根据交易状态处理
		switch v.TransStatus {
		// 交易成功
		case model.TransSuccess:
			sett := &model.TransSett{
				Tran:        v,
				SettFlag:    1,
				SettDate:    now.Format("2006-01-02 15:04:05"),
				MerSettAmt:  v.TransAmt * 9 / 10,
				MerFee:      v.TransAmt / 10,
				ChanSettAmt: v.TransAmt * 9 / 10,
				ChanFee:     v.TransAmt / 10,
			}
			if err := mongo.TransSettColl.Add(sett); err != nil {
				g.Error("add trans sett fail : %s", err)
			}
		// 处理中
		case model.TransHandling:
			// TODO调用交易查询更新状态
		}

	}

	//TODO 勾兑
	//cfca.ProcessTransChecking(be)
}
