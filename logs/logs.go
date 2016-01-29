package logs

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

func init() {
	go persistSpLogs()
}

var SpLogs = make(chan *model.SpTransLogs, 1e6)

// persistSpLogs 持久化扫码交易日志
func persistSpLogs() {

	for {
		select {
		case l := <-SpLogs:
			switch l.MsgType {
			case 1:
				mongo.SpMerLogsCol.Add(l)
			case 2, 3:
				// 2 渠道同步返回报文
				// 3 渠道异步返回报文
				mongo.SpChanLogsCol.Add(l)
			default:
				log.Warnf("unknown msgType: %d", l.MsgType)
			}
		}
	}

}
