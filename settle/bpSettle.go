// Package settle used to Settlement
package settle

import (
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/app"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/channel/cfca"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

// yesterday 每天不一样
var yesterday string
var ip = IPv4()

const (
	interval = 24 * time.Hour
	ld       = "2006-01-02"
	lt       = "2006-01-02 15:04:05"
)

// DoSettWork 开启任务
func DoSettWork() {
	go processTransSettle()
}

// ProcessTransSettle 清分
// TODO: 重构到定时任务中
func processTransSettle() {

	// 凌晨10分时将交易数据copy到清分表
	// 距离指定的时间
	// dis, _ := util.TimeToGiven("00:10:00")
	// afterFunc(dis*time.Second, "doTransSett")

	// 扫码支付每天定时copy交易到sett表，随后进行勾兑，出报表
	disOs, _ := util.TimeToGiven("00:30:00")
	afterFunc(disOs*time.Second, "doSpTransSett")

	// 中金渠道
	// disCfca, _ := util.TimeToGiven("08:00:00")
	// afterFunc(disCfca*time.Second, "doCFCATransCheck")

	// 讯联线下渠道
	// disCil, _ := util.TimeToGiven("01:00:00")
	// afterFunc(disCfca*time.Second, "doCILTransCheck")

	// 扫码每天出报表
	// disReport, _ := util.TimeToGiven("08:00:00")
	// afterFunc(disReport*time.Second, "doScanpaySettReport")

	// app用户每天发邮件
	appEmail, _ := util.TimeToGiven("23:00:00")
	afterFunc(appEmail*time.Second, "doAppToolsSendEmail")

	// 主线程阻塞
	select {}
}

func afterFunc(d time.Duration, method string) {
	log.Infof("prepare to process %s method after %s ", method, d)
	time.AfterFunc(d, func() {
		// 到点时先执行一次
		do(method)
		// 24小时后执行
		tick := time.Tick(interval)
		for {
			select {
			case <-tick:
				do(method)
			}
		}
	})
}

func do(method string) {

	// 查找昨天的交易
	now := time.Now()
	d, _ := time.ParseDuration("-24h")
	yesterday = now.Add(d).Format(ld)
	log.Debugf("yesterday : %s", yesterday)

	// 判断是否可执行
	l := &model.TransSettLog{
		Date:       yesterday,
		Method:     method,
		CreateTime: time.Now().Format(lt),
	}
	updated, err := mongo.TransSettLogColl.AtomUpsert(l)
	if err != nil {
		log.Errorf("fail to Upsert transSettlog: %s ", err)
		return
	}
	// 假如updated == 1
	// 说明已经存在该记录
	// 即该任务已被执行
	if updated == 1 {
		return
	}

	defer func() {
		if err != nil {
			log.Errorf("process %s fail: %s", method, err)
			l.Status = 2
		} else {
			l.Status = 1
		}
		l.Addr = ip
		l.ModifyTime = time.Now().Format(lt)
		mongo.TransSettLogColl.AtomUpsert(l)
	}()

	switch method {
	case "doTransSett":
		doTransSett()
	case "doCFCATransCheck":
		doCFCATransCheck()
	case "doCILTransCheck":
		doCilTransCheck()
	case "doAppToolsSendEmail":
		app.NotifySalesman()
	case "doScanpaySettReport":
		err = SpSettReport(yesterday)
	case "doSpTransSett":
		err = DoSpTransSett(yesterday, false)
	default:
		//..
	}
}

func doTransSett() {

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
				ChanMerId:   v.ChanMerId,
				SysOrderNum: v.SysOrderNum,
				PrivateKey:  chanMer.PrivateKey,
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
}

// doCFCATransCheck 中金渠道勾兑
// 勾兑:只需确认系统的交易记录在渠道方是否存在
// 不用勾兑金额
func doCFCATransCheck() {

	chanMers, err := mongo.ChanMerColl.FindByCode("CFCA")
	if err != nil {
		log.Errorf("fail to load all cfca chanMer %s", err)
	}
	// 中金渠道对象
	c := cfca.DefaultClient

	// 遍历渠道商户
	for _, v := range chanMers {
		resp := c.ProcessTransChecking(v.ChanMerId, yesterday, v.PrivateKey)
		if resp != nil && len(resp.Body.Tx) > 0 {
			for _, tx := range resp.Body.Tx {
				// 根据订单号查找
				if transSett, err := mongo.TransSettColl.FindByOrderNum(tx.TxSn); err == nil {
					// 找到记录，修改清分状态
					log.Infof("check success %+v", transSett)
					// TODO:transSett.SettFlag = model.SettSuccess
					if err = mongo.TransSettColl.Update(transSett); err != nil {
						log.Errorf("fail to update transSett record %s,transSett id : %s", err, transSett.Trans.Id)
					}

				} else {
					// 找不到，则是渠道多出的交易
					// 添加该笔交易
					newTrans := &model.Trans{
						Id:          bson.NewObjectId(),
						SysOrderNum: tx.TxSn,
						TransAmt:    tx.TxAmount,
					}
					// 判断交易类型
					switch {
					case tx.TxType == cfca.MerModePay:
						newTrans.TransType = model.PayTrans
					case tx.TxType == cfca.MerModeRefund:
						newTrans.TransType = model.RefundTrans
					}
					addTransSett(newTrans, model.SettChanRemain)
				}
			}
		}
	}
}

// doCilTransCheck 讯联线下渠道勾兑
// 勾兑:默认成功的交易都是勾兑成功
// TODO 从对账文件里分析
func doCilTransCheck() {

	chanMers, err := mongo.ChanMerColl.FindByCode("CIL")
	if err != nil {
		log.Errorf("fail to load all cfca chanMer %s", err)
	}

	for _, v := range chanMers {
		// mongo.TransSettColl
		log.Debugf("%v", v)
	}
}

// addTransSett 保存一条清分数据
// 计算手续费
func addTransSett(t *model.Trans, settFlag int8) {

	// TODO CIL 渠道暂时默认勾兑成功
	if t.ChanCode == "CIL" {
		settFlag = model.SettSuccess
	}

	var rate float64
	// 获得商户详情
	if t.MerId != "" {
		m, err := mongo.MerchantColl.Find(t.MerId)
		if err != nil {
			log.Errorf("fail to find merchant by merId(%s): %s", t.MerId, err)
		}
		log.Debugf("schemecd : %s", m.Detail.BillingScheme)
		if m.Detail.BillingScheme != "" {
			// scheme, err := mongo.SettSchemeCdCol.Find(m.Detail.BillingScheme)
			// if err != nil {
			// 	log.Errorf("fail to find settScheme by cd(%s): %s", m.Detail.BillingScheme, err)
			// }
			// 固定百分比
			schemeCd := m.Detail.BillingScheme
			if strings.HasPrefix(schemeCd, "00") {
				f, err := strconv.ParseFloat(schemeCd, 64)
				if err != nil {
					log.Errorf("fail to conver %s to float64: %s", schemeCd, err)
				}
				if f <= 5000 {
					rate = f / 100000
				}
			}
			// 非固定百分比
			// TODO...
		}

	}

	// 计算费率
	sett := &model.TransSett{}
	sett.Trans = *t
	// TODO:sett.SettFlag = settFlag
	sett.MerFee = int64(math.Floor(float64(t.TransAmt)*rate + 0.5)) // 四舍五入
	sett.MerSettAmt = t.TransAmt - int64(sett.MerFee)

	if err := mongo.TransSettColl.Add(sett); err != nil {
		log.Errorf("add trans sett fail : %s, trans id : %s", err, t.Id)
	}
}

func IPv4() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}
