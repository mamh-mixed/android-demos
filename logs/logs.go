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
			default:
				log.Warnf("unknown msgType: %d", l.MsgType)
			}
		}
	}

}
