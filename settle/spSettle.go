package settle

import (
	"github.com/CardInfoLink/log"
	"github.com/CardInfoLink/quickpay/channel/alipay"
	"github.com/CardInfoLink/quickpay/channel/weixin/scanpay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"math"
	"strconv"
	"time"
)

// func init() {
// 	needSettles = append(needSettles,
// 		&scanpayDomestic{
// 			At: goconf.Config.Settle.DomesticSettPoint,
// 		})
// }

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
	localMMap, alpMers, err := genLocalBlendMap(date)
	if err != nil {
		log.Errorf("find transSett error: %s", err)
		return
	}

	log.Infof("local MMap length=%s", len(localMMap))
	log.Infof("found alp mers: %d", len(alpMers))

	// 微信大商户
	wxpAgent, err := mongo.ChanMerColl.FindWXPAgent() //SettleDate SignKey
	if err != nil {
		log.Errorf("search the wxp agent error:%s", err)
	}
	chanMMap := make(model.ChanBlendMap)
	//微信请求
	for _, a := range wxpAgent {
		err = scanpay.DefaultWeixinScanPay.ProcessSettleEnquiry(&model.ScanPayRequest{
			AppID:     a.WxpAppId,
			ChanMerId: a.ChanMerId,
			SignKey:   a.SignKey,
			SettDate:  date,
		}, chanMMap)
		if err != nil {
			log.Errorf("the request error , merid:%s, chanCode:%s", a.ChanMerId, "WXP")
		}
	}

	//支付宝
	for _, k := range alpMers {
		c, err := mongo.ChanMerColl.Find("ALP", k)
		if err != nil {
			log.Errorf("find alp mer info error:%s", k)
			continue
		}
		err = alipay.Domestic.ProcessSettleEnquiry(&model.ScanPayRequest{
			SignKey:   c.SignKey,
			ChanMerId: k,
			SettDate:  date,
		}, chanMMap)
		if err != nil {
			log.Errorf("the request error: %s , merid:%s, chanCode:%s", err, c.ChanMerId, "ALP")
		}
	}

	log.Infof("chan MMap length=%d", len(chanMMap))

	// 金额有误交易
	amtErrorMap := make(map[string]string)

	// 勾兑过程
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
						act := float64(orderRecord.Trans.TransAmt) / 100
						if orderRecord.Trans.TransType == model.PayTrans {
							orderACT += act
						} else {
							orderACT -= act
						}
					}

					// 相等保存
					if math.Abs(blendACT-orderACT) < 0.001 {
						for _, transSett := range transSetts {
							transSett.BlendType = MATCH
							transSett.SettTime = time.Now().Format("2006-01-02 15:04:05")
							mongo.SpTransSettColl.Update(&transSett)
						}
						delete(localOrderMap, chanOrderNum) //删除本地记录，剩下的进C001
						delete(chanOrderMap, chanOrderNum)  //删除渠道记录，剩下的进C002
					} else {
						log.Errorf("amt error: expect %0.2f, actual %0.2f, chanOrderNum=%s", orderACT, blendACT, chanOrderNum)
						// 对上，但金额不一致
						for _, transSett := range transSetts {
							transSett.BlendType = AMT_ERROR
							transSett.SettTime = time.Now().Format("2006-01-02 15:04:05")
							mongo.SpTransSettColl.Update(&transSett)
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
		if err = upload(rs.ReportName, genC002ReportExcel(chanMMap, date)); err != nil {
			if err = mongo.RoleSettCol.Upsert(rs); err != nil {
				log.Errorf("roleSett upsert error: %s", err)
			}
		}
	}
}

// genLocalBlendMap 根据当天交易生成本地勾兑数据集
func genLocalBlendMap(date string) (lbm model.LocalBlendMap, alpChanMer map[string]string, err error) {

	var transSetts []model.TransSett
	transSetts, err = mongo.SpTransSettColl.Find(&model.QueryCondition{Date: date})
	if err != nil {
		return
	}

	if len(transSetts) == 0 {
		return
	}

	lbm = make(model.LocalBlendMap)
	alpChanMer = make(map[string]string)
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

		// 存储支付宝发生交易渠道号
		if ts.Trans.ChanCode == "ALP" {
			if _, ok := alpChanMer[ts.Trans.ChanMerId]; !ok {
				alpChanMer[ts.Trans.ChanMerId] = ts.Trans.ChanMerId
			}
		}

	}

	return

}
