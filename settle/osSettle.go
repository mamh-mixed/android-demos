package settle

import (
	"bufio"
	// "fmt"
	"github.com/CardInfoLink/quickpay/currency"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
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
// 拷贝所有交易成功的交易，包括交易关闭中被退款的
// 这时勾兑状态默认是渠道少清的
func copyTrans(date string) {}

// readSftpTxt 返回的数据格式如下：
// orderNum|chanOrderNum|amt|fee|time|type|
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

func (a *alipayOverseas) Reconciliation() error {
	for _, data := range a.chanData {
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
