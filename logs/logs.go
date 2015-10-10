package logs

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

func init() {
	go persistSpLogs()
}

var SpLogs = make(chan *model.SpTransLogs, 1e5)

// persistSpLogs 持久化扫码交易日志
func persistSpLogs() {

	for {
		select {
		case l := <-SpLogs:
			switch l.MsgType {
			case 1:
				mongo.SpMerLogsCol.Add(l)
			case 2:
				mongo.SpChanLogsCol.Add(l)
			case 3:
				// 异步消息通知，先关联reqId
				pl, err := mongo.SpChanLogsCol.FindOne(&model.QueryCondition{MerId: l.MerId, Busicd: l.TransType, OrderNum: l.OrderNum})
				if err != nil {
					log.Errorf("find paut logs error: %s, find condition: %v", err, l)
					continue
				}
				// 关联reqId
				l.ReqId = pl.ReqId
				mongo.SpChanLogsCol.Add(l)

			default:
				log.Warnf("unknown msgType: %d", l.MsgType)
			}
		}
	}

}
