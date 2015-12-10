package settle

import (
	"bufio"
	// "fmt"
	"fmt"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/currency"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

func init() {
	needSettles = append(needSettles,
		&alipayOverseas{
			At:       goconf.Config.Settle.OverseasSettPoint,
			sftpAddr: "sftp.alipay.com:22",
		})
}

type alipayOverseas struct {
	At       string // "02:00:00" 表示凌晨两点才可以拉取数据
	sftpAddr string
}

// ProcessDuration 返回一个当前时间到可执行时间的duration
func (a *alipayOverseas) ProcessDuration() time.Duration {
	if a.At == "" {
		return time.Duration(0)
	}

	now := time.Now()
	t, err := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" "+a.At, time.Local)
	if err != nil {
		log.Errorf("parse time error: %s", err)
		return time.Duration(-1)
	}

	if now.After(t) {
		return time.Duration(0)
	}

	return t.Sub(now)
}

// download 从sftp下载对账文件到服务器上
func (a *alipayOverseas) download(date string) [][]string {
	var data [][]string
	// 境外商户
	chanMers, err := mongo.ChanMerColl.FindByArea(channel.Oversea)
	if err != nil {
		log.Errorf("fail to find overseas alp merchant: %s", err)
		return data
	}

	path := "/download/%s_transaction_%s.txt"

	// date yyyymmdd
	date = strings.Replace(date, "-", "", -1)
	for _, cm := range chanMers {
		if cm.Sftp != nil {

			var authMethods []ssh.AuthMethod
			// add password
			authMethods = append(authMethods, ssh.Password(cm.Sftp.Password))
			config := &ssh.ClientConfig{
				User: cm.Sftp.Username,
				Auth: authMethods,
			}

			// 建立连接
			conn, err := ssh.Dial("tcp", a.sftpAddr, config)
			if err != nil {
				log.Errorf("fail to connect sftp service, error: %s", err)
				continue
			}

			client, err := sftp.NewClient(conn)
			if err != nil {
				log.Errorf("fail to get sftp client, error: %s", err)
				continue
			}

			f, err := client.Open(fmt.Sprintf(path, cm.ChanMerId, date))
			if err != nil {
				log.Errorf("open file error: %s", err)
				continue
			}

			// read
			s := bufio.NewScanner(f)
			// skip 2
			for i := 0; s.Scan(); i++ {
				if i > 1 {
					// Partner_transaction_id-0|Transaction_id-1|Transaction_amount-2|Charge_amount-3|Currency-4|Payment_time-5|Transaction_type-6|Remark-7
					ts := strings.Split(s.Text(), "|")
					data = append(data, ts)
				}
			}
		} else {
			log.Warnf("overseas alipay merchant no sftp config, check ... pid=%s", cm.ChanMerId)
		}
	}
	return data
}

// Reconciliation 海外接口对账函数
func (a *alipayOverseas) Reconciliation(date string) {

	// 下载渠道数据
	recs := a.download(date)

	for _, data := range recs {
		if len(data) != 8 {
			log.Errorf("invalid reconciliation data length=%d should be 8", len(data))
			continue
		}
		orderNum, chanOrderNum := data[0], data[1]

		// 排除可能一时的系统错误导致查询失败
		var ts *model.TransSett
		var err error

		// Transaction_type
		switch data[6] {
		case "REVERSAL", "CANCEL":
			// 只能当天
			// 暂时没有CANCEL支付类型
			// 该情况下没有原订单，只有撤销和退款的订单
			// 渠道流水中的订单，假如取消的这笔原订单是成功的，那么给出的是原订单号和渠道订单号
			// 假如原订单不是成功的，那么给出的是取消这笔的订单号和渠道订单号
			ts, err = mongo.SpTransSettColl.FindOne(orderNum, chanOrderNum)
			if err != nil {
				// 没找到，说明原订单没有成功
				// 不处理即可
				break
			}
			// 找到，说明是原订单成功下发起的取消
			// 那么比较金额即可，手续费都重置为0
			mongo.SpTransSettColl.RemoveOne(orderNum, chanOrderNum)
		case "PAYMENT", "REFUND":
			ts, err = mongo.SpTransSettColl.FindOne(orderNum, chanOrderNum)
			if err != nil {
				// 渠道多清
				mt := &model.TransSett{}
				mt.Trans.Id = bson.NewObjectId() // 没有ID会报错
				mt.Trans.OrderNum = orderNum
				mt.Trans.ChanOrderNum = chanOrderNum
				mt.Trans.TransAmt = currency.I64(data[4], data[2])
				mt.Trans.Currency = data[4]
				mt.Trans.PayTime = data[5]
				mt.AcqFee = currency.I64(data[4], data[3])
				mt.BlendType = CHAN_MORE // blendType
				if data[6] == "PAYMENT" {
					mt.Trans.Busicd = model.Purc
					mt.Trans.TransType = model.PayTrans
				} else {
					mt.Trans.Busicd = model.Refd
					mt.Trans.TransType = model.RefundTrans
				}
				mongo.SpTransSettColl.Add(mt)
				continue
			}
			t := ts.Trans

			// 开始勾兑，默认成功
			ts.BlendType = MATCH

			// 不管是支付交易还是逆向交易，成功的交易都是有金额的，所以直接比较金额即可。
			if currency.Str(t.Currency, t.TransAmt) != data[2] {
				// 金额不一致
				ts.BlendType = AMT_ERROR
			} else if currency.Str(t.Currency, ts.InsFee) != data[3] {
				// 不管是支付交易还是逆向交易，都是有计算手续费的。
				// 手续费不一致
				ts.BlendType = FEE_ERROR
			}

			// 时间
			ts.SettTime = time.Now().Format("2006-01-02 15:04:05")

			// 更新交易状态
			if err = mongo.SpTransSettColl.Update(ts); err != nil {
				log.Errorf("fail to update transSett error: %s, orderNum=%s", err, orderNum)
			}
		}
	}
}
