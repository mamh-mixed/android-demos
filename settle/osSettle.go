package settle

import (
	"bufio"
	// "fmt"
	"fmt"
	"github.com/CardInfoLink/quickpay/currency"
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

type alipayOverseas struct {
	Date     string
	sftpAddr string
	Data     [][]string
}

// download 从sftp下载对账文件到服务器上
func (a *alipayOverseas) download() {
	var chanMers []model.ChanMer

	path := "/download/%s_transaction_%s.txt"
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

			f, err := client.Open(fmt.Sprintf(path, cm.ChanMerId, a.Date))
			if err != nil {
				log.Errorf("open file error: %s", err)
				continue
			}

			// read
			s := bufio.NewScanner(f)
			// skip 2
			for i := 0; s.Scan(); i++ {
				if i > 1 {
					// orderNum|chanOrderNum|amt|fee|time|type|
					ts := strings.Split(s.Text(), "|")
					a.Data = append(a.Data, ts)
				}
			}
		} else {
			log.Warnf("overseas alipay merchant no sftp config, check ... pid=%s", cm.ChanMerId)
		}
	}
}

// Reconciliation 海外接口对账函数
func (a *alipayOverseas) Reconciliation() error {
	for _, data := range a.Data {
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

		// 更新交易状态
		mongo.SpTransSettColl.Update(ts)

	}
	return nil
}
