package core

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/channel/cfca"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// init 开启任务routine
func init() {
	log.Debug("wait to process transSett method")
	go ProcessTransSettle()
}

// ProcessTransSettle 清分
func ProcessTransSettle() {

	// 暂时先每天早上8点将交易信息拷贝到清分表里
	// 距离指定的时间
	dis, err := tools.TimeToGiven("08:00:00")
	if err != nil {
		log.Errorf("fail to get time second by given %s", err)
		return
	}
	c := make(chan bool)
	time.AfterFunc(time.Duration(dis)*time.Second, func() {
		// time.AfterFunc(10*time.Second, func() {
		// 到点时先执行一次
		doTransSett()
		// 24小时后执行
		tick := time.Tick(24 * time.Hour)
		for {
			select {
			case <-tick:
				// g.Debug("tick ... %s", "boom")
				doTransSett()
			}
		}

	})

	<-c

}

func doTransSett() {

	layout := "2006-01-02"
	// 查找昨天的交易
	now := time.Now()
	d, _ := time.ParseDuration("-24h")
	yesterday := now.Add(d).Format(layout)
	log.Debugf("yesterday : %s", yesterday)
	trans, err := mongo.TransColl.FindByTime(yesterday)
	if err != nil {
		log.Errorf("fail to load trans by time : %s", err)
		return
	}

	// 交易数据
	for _, v := range trans {
		// 根据交易状态处理
		switch v.TransStatus {
		// 交易成功
		case model.TransSuccess:
			addTransSett(v, model.SettSysRemain)
		// 处理中
		case model.TransHandling:
			// TODO根据渠道代码得到渠道实例，暂时默认cfca

			// 得到渠道商户，获取签名密钥
			chanMer, err := mongo.ChanMerColl.Find(v.ChanCode, v.ChanMerId)
			if err != nil {
				log.Errorf("fail to find chanMer(%s,%s) : %s", v.ChanCode, v.ChanMerId, err)
				continue
			}
			// 封装参数
			be := &model.OrderEnquiry{
				ChanMerId:    v.ChanMerId,
				ChanOrderNum: v.ChanOrderNum,
				SignCert:     chanMer.SignCert,
			}

			// 根据交易类型处理
			ret := new(model.BindingReturn)
			c := channel.GetChan(chanMer.ChanCode)
			switch v.TransType {
			// 支付交易
			case model.PayTrans:
				ret = c.ProcessPaymentEnquiry(be)
			// 退款交易
			case model.RefundTrans:
				ret = c.ProcessRefundEnquiry(be)
			}

			// 处理结果
			if ret.RespCode == "000000" {
				// 支付成功、退款成功
				v.RespCode = ret.RespCode
				v.TransStatus = model.TransSuccess
				// 更新交易状态
				mongo.TransColl.Update(v)
				// 添加到清分表
				addTransSett(v, model.SettSysRemain)
			} else if ret.RespCode == "100070" || ret.RespCode == "100080" {
				// 支付失败、退款失败
				v.RespCode = ret.RespCode
				v.TransStatus = model.TransFail
				// 更新交易状态
				mongo.TransColl.Update(v)
			} else {
				// 不处理
			}
		}

	}

	// 勾兑:只需确认系统的交易记录在渠道方是否存在
	// 不用勾兑金额
	doTransCheck(yesterday)
}

// addTransSett 保存一条清分数据
// 计算手续费
func addTransSett(t *model.Trans, settFlag int8) {
	sett := &model.TransSett{
		Tran:     *t,
		SettFlag: settFlag,
		// TODO
		MerSettAmt:  t.TransAmt * 9 / 10,
		MerFee:      t.TransAmt / 10,
		ChanSettAmt: t.TransAmt * 9 / 10,
		ChanFee:     t.TransAmt / 10,
	}
	if err := mongo.TransSettColl.Add(sett); err != nil {
		log.Errorf("add trans sett fail : %s, trans id : %s", err, t.Id)
	}
}

// doTransCheck 勾兑
func doTransCheck(settDate string) {
	chanMers, err := mongo.ChanMerColl.FindAll()
	if err != nil {
		log.Errorf("fail to load all chanMer %s", err)
	}
	// 遍历渠道商户
	for _, v := range chanMers {

		// TODO 应该根据chanCode获得渠道实例
		// 暂时先默认cfca
		// c := channel.GetChan(v.ChanCode)
		if v.ChanCode == "CFCA" {
			c := cfca.Obj
			resp := c.ProcessTransChecking(v.ChanMerId, settDate, v.SignCert)
			if resp != nil && len(resp.Body.Tx) > 0 {
				for _, tx := range resp.Body.Tx {
					// 根据订单号查找
					if transSett, err := mongo.TransSettColl.FindByOrderNum(tx.TxSn); err == nil {
						// 找到记录，修改清分状态
						log.Infof("check success %+v", transSett)
						transSett.SettFlag = model.SettSuccess
						if err = mongo.TransSettColl.Update(transSett); err != nil {
							log.Errorf("fail to update transSett record %s,transSett id : %s", err, transSett.Tran.Id)
						}

					} else {
						// 找不到，则是渠道多出的交易
						// 添加该笔交易
						newTrans := &model.Trans{
							Id:           bson.NewObjectId(),
							ChanOrderNum: tx.TxSn,
							TransAmt:     tx.TxAmount,
						}
						// 判断交易类型
						switch {
						case tx.TxType == cfca.BindingPaymentTxCode:
							newTrans.TransType = model.PayTrans
						case tx.TxType == cfca.BindingRefundTxCode:
							newTrans.TransType = model.RefundTrans
						}
						addTransSett(newTrans, model.SettChanRemain)
					}
				}
			}
		}
	}
}
