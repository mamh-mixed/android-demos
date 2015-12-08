package settle

import (
	"github.com/CardInfoLink/quickpay/channel/alipay"
	"github.com/CardInfoLink/quickpay/channel/weixin/scanpay"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"math"
	"strconv"
	"strings"
	"time"
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

	// 微信大商户
	wxpAgent, err := mongo.ChanMerColl.FindWXPAgent() //SettleDate SignKey
	if err != nil {
		log.Errorf("search the wxp agent error:%s", err)
	}
	chanMMap := make(model.ChanBlendMap)
	blendDateStr := strings.Replace(date, "-", "", -1)
	//微信请求
	wxpreq := new(model.ScanPayRequest)
	wxpreq.SettleDate = blendDateStr
	for _, a := range wxpAgent {
		wxpreq.AppID = a.WxpAppId
		wxpreq.ChanMerId = a.ChanMerId
		wxpreq.SignKey = a.SignKey
		err = scanpay.DefaultWeixinScanPay.ProcessSettleEnquiry(wxpreq, chanMMap)
		if err != nil {
			log.Errorf("the request error , merid:%s, chanCode:%s", wxpreq.ChanMerId, "WXP")
		}
	}

	//支付宝
	alpreq := new(model.ScanPayRequest)
	alpreq.StartTime = date + " 00:00:00"
	alpreq.EndTime = date + " 23:59:59"
	var alpMers []string
	for _, k := range alpMers {
		c, err := mongo.ChanMerColl.Find("ALP", k)
		if err != nil {
			log.Errorf("find alp mer info error:%s", k)
			continue
		}
		alpreq.SignKey = c.SignKey
		alpreq.ChanMerId = k
		err = alipay.Domestic.ProcessSettleEnquiry(alpreq, chanMMap)
		if err != nil {
			log.Errorf("the request error , merid:%s, chanCode:%s", alpreq.ChanMerId, "ALP")
		}
	}
	// 渠道多余
	chanMoreMap := make(map[string]string)

	// 本地数据集
	localMMap := make(model.LocalBlendMap) // TODO 整理数据源

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
						act := float64(orderRecord.Trans.TransAmt / 100)
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
							mongo.SpTransSettColl.Update(&transSett)
						}
						delete(localOrderMap, chanOrderNum) //删除本地记录，剩下的进C001
						delete(chanOrderMap, chanOrderNum)  //删除渠道记录，剩下的进C002
					} else {
						// 对上，但金额不一致
						for _, blendRecord := range blendArray {
							chanMoreMap[blendRecord.OrderID] = blendRecord.ChanMerID
						}
					}
				}
			}
		}
	}

	// 处理没有勾兑上的数据

	// upload
	if len(localMMap) != 0 {
		file401 := "IC401.xlsx"
		upload(file401, genC001ReportExcel(localMMap, date))
	}

	if len(chanMMap) != 0 {
		file402 := "IC402.xlsx"
		upload(file402, genC002ReportExcel(chanMMap, date))
	}
}

// genLocalBlendMap 根据当天交易生成本地勾兑数据集
func genLocalBlendMap(date string) model.LocalBlendMap {

	var lbm model.LocalBlendMap
	_, err := mongo.SpTransSettColl.Find(&model.QueryCondition{
		StartTime: date + " 00:00:00",
		EndTime:   date + " 23:59:59",
	})
	if err != nil {
		return lbm
	}

	return lbm

}
