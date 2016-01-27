package settle

import (
	"math"
	"strconv"
	"time"

	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/channel/alipay"
	"github.com/CardInfoLink/quickpay/channel/weixin/scanpay"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

func init() {
	needSettles = append(needSettles,
		&scanpayDomestic{
			At: goconf.Config.Settle.DomesticSettPoint,
		})
}

type scanpayDomestic struct {
	At string // "02:00:00" 表示凌晨两点才可以拉取数据
}

// ProcessDuration 可执行duration
func (s scanpayDomestic) ProcessDuration() time.Duration {
	if s.At == "" {
		return time.Duration(0)
	}

	now := time.Now()
	t, err := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" "+s.At, time.Local)
	if err != nil {
		log.Errorf("parse time error: %s", err)
		return time.Duration(-1)
	}

	if now.After(t) {
		return time.Duration(0)
	}

	return t.Sub(now)
}

// Reconciliation 勾兑
func (s scanpayDomestic) Reconciliation(date string) {

	// 本地数据集
	localMMap, alpMers, wxpMers, err := genLocalBlendMap(date)
	if err != nil {
		log.Errorf("find transSett error: %s", err)
		return
	}

	log.Infof("begin blend, localMMap length=%d, alpMers=%d, wxpMers=%d", len(localMMap), len(alpMers), len(wxpMers))
	chanMMap := make(model.ChanBlendMap)
	//微信请求
	for _, k := range wxpMers {
		c, err := mongo.ChanMerColl.Find(channel.ChanCodeWeixin, k)
		if err != nil {
			log.Errorf("find chanMer(%s) error: %s", k, err)
			continue
		}

		req := &model.ScanPayRequest{SettDate: date}
		if c.IsAgentMode && c.AgentMer != nil {
			req.AppID = c.AgentMer.WxpAppId
			req.ChanMerId = c.AgentMer.ChanMerId
			req.SubMchId = c.ChanMerId
			req.SignKey = c.AgentMer.SignKey
		} else {
			req.AppID = c.WxpAppId
			req.ChanMerId = c.ChanMerId
			req.SignKey = c.SignKey
		}
		err = scanpay.DefaultWeixinScanPay.ProcessSettleEnquiry(req, chanMMap)
		if err != nil {
			log.Errorf("the request error , chanMerId=%s, chanCode=%s", req.ChanMerId, "WXP")
		}
	}

	//支付宝
	for _, k := range alpMers {
		c, err := mongo.ChanMerColl.Find("ALP", k)
		if err != nil {
			log.Errorf("find alp mer info error: %s", k)
			continue
		}
		err = alipay.Domestic.ProcessSettleEnquiry(&model.ScanPayRequest{
			SignKey:   c.SignKey,
			ChanMerId: k,
			SettDate:  date,
		}, chanMMap)
		if err != nil {
			log.Errorf("the request error: %s , merId: %s, chanCode: %s", err, c.ChanMerId, "ALP")
		}
	}

	log.Infof("begin blend, chanMMap length=%d", len(chanMMap))

	// 金额有误交易
	amtErrorMap := make(map[string]string)

	// 勾兑过程
	var localSuccess int
	var chanSuccess int
	for chanMerId, localOrderMap := range localMMap {
		chanOrderMap, ok := chanMMap[chanMerId]
		if ok {
			for chanOrderNum, transSetts := range localOrderMap { //查找该商户记录
				blendArray, ok := chanOrderMap[chanOrderNum]
				if ok {
					// 查到该商户的记录
					var blendACT float64 = 0.0
					var orderACT float64 = 0.0
					for _, blendRecord := range blendArray { //计算渠道该订单总金额
						tempAct, err := strconv.ParseFloat(blendRecord.OrderAct, 64)
						if err == nil {
							blendACT += tempAct
						}
					}

					for _, orderRecord := range transSetts { //计算本地该订单总金额
						var fee int64
						switch orderRecord.Trans.ChanCode {
						case channel.ChanCodeAlipay:
							// 支付宝是包含手续费计算的
							fee = orderRecord.MerFee
						case channel.ChanCodeWeixin:
							// TODO:手续费校验
						}
						act := float64(orderRecord.Trans.TransAmt-fee) / 100
						if orderRecord.Trans.TransType == model.PayTrans {
							orderACT += act
						} else {
							orderACT -= act
						}
					}

					// 相等保存
					if math.Abs(blendACT-orderACT) < 0.001 {
						for _, transSett := range transSetts {
							localSuccess++
							transSett.BlendType = MATCH
							transSett.SettTime = time.Now().Format("2006-01-02 15:04:05")
							// log.Infof("blend success, merId=%s, orderNum=%s, chanOrderNum=%s", transSett.Trans.MerId, transSett.Trans.OrderNum, transSett.Trans.ChanOrderNum)
							// mongo.SpTransSettColl.Update(&transSett)
						}
						chanSuccess += len(blendArray)
						delete(localOrderMap, chanOrderNum) //删除本地记录，剩下的进C001
						delete(chanOrderMap, chanOrderNum)  //删除渠道记录，剩下的进C002
					} else {
						log.Errorf("amt error: expect %0.2f, actual %0.2f, chanOrderNum=%s", orderACT, blendACT, chanOrderNum)
						for _, local := range transSetts {
							log.Errorf("merId=%s, orderNum=%s, chanOrderNum=%s", local.Trans.MerId, local.Trans.OrderNum, local.Trans.ChanOrderNum)
						}
						for _, blend := range blendArray {
							log.Errorf("chanMerId=%s, orderNum=%s, chanOrderNum=%s, amt=%s", blend.ChanMerID, blend.LocalID, blend.OrderID, blend.OrderAct)
						}
						// 对上，但金额不一致
						for _, transSett := range transSetts {
							transSett.BlendType = AMT_ERROR
							transSett.SettTime = time.Now().Format("2006-01-02 15:04:05")
							// mongo.SpTransSettColl.Update(&transSett)
						}
						amtErrorMap[chanOrderNum] = chanOrderNum // 只是打个标记
					}
				}
			}

			// 如果内部map为空，那么删除外部key
			if len(chanOrderMap) == 0 {
				delete(chanMMap, chanMerId)
			}
			if len(localOrderMap) == 0 {
				delete(localMMap, chanMerId)
			}
		}
	}
	log.Infof("blend success localSuccess=%d,chanSuccess=%d", localSuccess, chanSuccess)
	/*
			var localTrans int
			for _, v := range localMMap {
				for _, v1 := range v {
					localTrans += len(v1)
					// for _, ts := range v1 {
					// 	localTrans++
					// 	// log.Infof("after blend localMMap: merId=%s,orderNum=%s,chanOrderNum=%s,chanMerId=%s", ts.Trans.MerId, ts.Trans.OrderNum, ts.Trans.ChanOrderNum, ts.Trans.ChanMerId)
					// }
				}
			}
			log.Infof("after blend localMMap, remain=%d", localTrans)
			var chanTrans int
			for _, v := range chanMMap {
				for _, v1 := range v {
					chanTrans += len(v1)
					// for _, b := range v1 {
					// 	log.Infof("after blend chanMMap: orderNum=%s,chanOrderNum=%s,chanMerId=%s", b.LocalID, b.OrderID, b.ChanMerID)
					// }
				}
			}
		log.Infof("after blend chanMMap, remain=%d", chanTrans)
	*/
	log.Infof("after blend localMMap length=%d", len(localMMap))
	log.Infof("after blend chanMMap length=%d", len(chanMMap))
	log.Infof("after blend errAmtMap length=%d", len(amtErrorMap))

	// 处理没有勾兑上的数据
	// 渠道少清
	if len(localMMap) != 0 {
		// 上传并记录
		rs := getRsRecord(ChanLessReport, date)
		if err = upload(rs.ReportName, genC001ReportExcel(localMMap, date)); err == nil {
			if err = mongo.RoleSettCol.Upsert(rs); err != nil {
				log.Errorf("roleSett upsert error: %s", err)
			}
		}
	}

	// 渠道多清
	if len(chanMMap) != 0 {
		for _, v := range localMMap {
			for ik, iv := range v {
				// 没勾兑上的里面包含金额错误的
				if _, ok := amtErrorMap[ik]; ok {
					// 跳过
					continue
				}
				for _, transSett := range iv {
					transSett.BlendType = CHAN_MORE
					// mongo.SpTransSettColl.Update(&transSett)
				}
			}
		}
		rs := getRsRecord(ChanMoreReport, date)
		if err = upload(rs.ReportName, genC002ReportExcel(chanMMap, date)); err == nil {
			if err = mongo.RoleSettCol.Upsert(rs); err != nil {
				log.Errorf("roleSett upsert error: %s", err)
			}
		}
	}
}

// genLocalBlendMap 根据当天交易生成本地勾兑数据集
func genLocalBlendMap(date string) (lbm model.LocalBlendMap, alpChanMer map[string]string, wxpChanMer map[string]string, err error) {

	var transSetts []model.TransSett
	transSetts, err = mongo.SpTransSettColl.Find(&model.QueryCondition{
		Date:        date,
		IsForReport: true,
		ChanCode:    channel.ChanCodeWeixin,
	})
	if err != nil {
		return
	}

	log.Infof("the trans len is %d", len(transSetts))

	if len(transSetts) == 0 {
		return
	}

	lbm = make(model.LocalBlendMap)
	alpChanMer = make(map[string]string)
	wxpChanMer = make(map[string]string)
	for _, ts := range transSetts {
		if chanOrderMap, ok := lbm[ts.Trans.ChanMerId]; ok {
			if tss, found := chanOrderMap[ts.Trans.ChanOrderNum]; found {
				tss = append(tss, ts)
				chanOrderMap[ts.Trans.ChanOrderNum] = tss
			} else {
				chanOrderMap[ts.Trans.ChanOrderNum] = []model.TransSett{ts}
			}
		} else {
			chanOrderMap := make(map[string][]model.TransSett)
			chanOrderMap[ts.Trans.ChanOrderNum] = []model.TransSett{ts}
			lbm[ts.Trans.ChanMerId] = chanOrderMap
		}

		switch ts.Trans.ChanCode {
		case channel.ChanCodeAlipay:
			if _, ok := alpChanMer[ts.Trans.ChanMerId]; !ok {
				alpChanMer[ts.Trans.ChanMerId] = ts.Trans.ChanMerId
			}
		case channel.ChanCodeWeixin:
			if _, ok := wxpChanMer[ts.Trans.ChanMerId]; !ok {
				wxpChanMer[ts.Trans.ChanMerId] = ts.Trans.ChanMerId
			}
		}
	}

	return

}
