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
	"strings"
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
					// Partner_transaction_id|Transaction_id|Transaction_amount|Charge_amount|Currency|Payment_time|Transaction_type|Remark
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
		var retry int
		for {
			ts, err = mongo.SpTransSettColl.FindOne(orderNum, chanOrderNum)
			if err != nil {
				retry++
				if retry == 2 {
					// 渠道多清
					mt := &model.TransSett{}
					mt.Trans.OrderNum = orderNum
					mt.Trans.ChanOrderNum = chanOrderNum
					mt.BlendType = CHAN_MORE
					// TODO..
					mongo.SpTransSettColl.Add(mt)
					break
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
		t := ts.Trans

		// 开始勾兑，默认成功
		ts.BlendType = MATCH

		// 原交易要么成功，要么被退款，都认为是成功的
		if t.TransStatus != model.TransSuccess || t.RefundStatus != model.TransRefunded {
			// 渠道多清
			ts.BlendType = CHAN_MORE
		}

		// 不管是支付交易还是逆向交易，成功的交易都是有金额的，所以直接比较金额即可。
		if currency.Str(t.Currency, t.TransAmt) != data[2] {
			// 金额不一致
			ts.BlendType = AMT_ERROR
		}

		// 不管是支付交易还是逆向交易，都是有计算手续费的。
		if currency.Str(t.Currency, t.Fee) != data[3] {
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
