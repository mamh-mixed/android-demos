package settle

import (
	"bufio"
	"fmt"
	"github.com/CardInfoLink/quickpay/currency"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"os"
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

// copyTrans 假设每天23:59:59时拷贝到清算表
func copyTrans(date string) {}

func readSftpTxt(fn string) ([][]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(f)
	var data [][]string
	for s.Scan() {
		ts := strings.Split(s.Text(), "|")
		data = append(data, ts)
	}

	return data, nil
}

type alipayOverseas struct {
	chanData [][]string
}

func (a *alipayOverseas) dataHandle() {

}

func (a *alipayOverseas) Reconciliation() error {
	for _, data := range a.chanData {
		if len(data) != 8 {
			return fmt.Errorf("invalid reconciliation data length=%d should be 8", len(data))
		}
		orderNum, chanOrderNum := data[0], data[1]

		// 排除可能一时的系统错误导致查询失败
		var ts *model.TransSett
		var err error
		var retry int
		for {
			ts, err = mongo.TransSettColl.FindOne(orderNum, chanOrderNum)
			if err != nil {
				retry++
				if retry == 2 {
					// 渠道多清
					mt := &model.TransSett{}
					mt.Trans.OrderNum = orderNum
					mt.Trans.ChanOrderNum = chanOrderNum
					mt.BlendType = CHAN_MORE
					// TODO..
					mongo.TransSettColl.Add(mt)
					break
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}

		if currency.Str(ts.Trans.Currency, ts.Trans.TransAmt) != data[2] {
			// 金额不一致
		}

		if currency.Str(ts.Trans.Currency, ts.Trans.NetFee) != data[3] {
			// 手续费不一致
		}

	}
	return nil
}
