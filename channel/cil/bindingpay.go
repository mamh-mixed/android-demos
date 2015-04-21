package cil

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"time"
)

// 讯联交易类型
const (
	version               = "15.1"
	consumeBusicd         = "000000" // 消费
	orderConsumeBusicd    = "020000" // 订购消费
	consumeUndoBusicd     = "201000" // 消费撤销
	consumeReversalBusicd = "040000" // 消费冲正

	reversalFlag          = "TIME_OUT"          // 冲正标识
	reversalTime          = 10                  // 超时时间
	reversalTimeDuration1 = reversalTime * 1    // 超时间隔1
	reversalTimeDuration2 = reversalTime * 8    // 超时间隔2
	reversalTimeDuration3 = reversalTime * 50   // 超时间隔3
	reversalTimeDuration4 = reversalTime * 1150 // 超时间隔3
)

// reversalHandle 冲正处理方法
func reversalHandle(om *CilMsg) {
	log.Debug("源交易请求超时，发送冲正报文")
	// 创建冲正报文
	rm := &CilMsg{
		Busicd:       consumeReversalBusicd,
		Txndir:       "Q",
		Posentrymode: om.Posentrymode,
		Chcd:         om.Chcd,
		Clisn:        mongo.DaySNColl.GetDaySN(om.Mchntid, om.Terminalid),
		Mchntid:      om.Mchntid,
		Terminalid:   om.Terminalid,
		Txamt:        om.Txamt,
		Txcurrcd:     om.Txcurrcd,
		Cardcd:       om.Cardcd,
		Syssn:        om.Syssn,
		Origclisn:    om.Clisn,
		Localdt:      tools.LocalDt(),
	}

	// TODO 报文入库

	// 发送冲正消息的时间点信道
	dc := make(chan time.Duration)
	// 结束信道
	qc := make(chan int)

	// 冲正时间节点和结束标志入信道
	go func() {
		dc <- reversalTime
		dc <- reversalTimeDuration1
		dc <- reversalTimeDuration2
		dc <- reversalTimeDuration3
		dc <- reversalTimeDuration4
		qc <- 0
	}()

	for isOK := false; !isOK; {
		select {
		case drt := <-dc:
			// 先休眠，再发送
			log.Debugf("reversal request time out ,sleep %d second", drt)
			time.Sleep(drt * time.Second)

			back := send(rm)
			log.Debugf("冲正请求响应结果: %+v", back)

			// 发送成功，跳出循环
			if back.Respcd != reversalFlag {
				isOK = true
				// TODO 更新已存储的报文
			}

		case <-qc:
			log.Errorf("冲正请求发送失败，报文信息如下: %+v", rm)
			isOK = true
		}
	}
}

// 无卡直接消费（订购消费）
func Consume(p *model.NoTrackPayment) (ret *model.BindingReturn) {
	// 构建消费报文
	m := CilMsg{
		Busicd:       orderConsumeBusicd,
		Txndir:       "Q",
		Posentrymode: "012",
		Chcd:         p.Chcd,
		Clisn:        p.CliSN,
		Mchntid:      p.Mchntid,
		Terminalid:   p.Terminalid,
		Txamt:        fmt.Sprintf("%012d", p.TransAmt),
		Txcurrcd:     p.CurrCode,
		Cardcd:       p.AcctNum,
		Syssn:        p.SysSN,
		Localdt:      tools.LocalDt(),
		Expiredate:   p.ValidDate,
		Cvv2:         p.Cvv2,
	}
	log.Debugf("无卡直接支付开始向线下网关发送报文: %+v", m)
	// TODO 报文入库

	resp := send(&m)
	log.Debugf("无卡直接支付返回结果:%+v", resp)

	// 不超时，转换应答码后返回
	if resp.Respcd != reversalFlag {
		ret = transformResp(resp.Respcd)
		return
	}

	// 超时需要冲正
	reversalHandle(&m)

	// 返回‘外部系统错误’的应答码
	ret = mongo.RespCodeColl.Get("000002")
	return
}

// ConsumeByApplePay ApplePay消费
func ConsumeByApplePay(ap *model.ApplePay) (ret *model.BindingReturn) {
	m := CilMsg{
		Busicd:        consumeBusicd,
		Txndir:        "Q",
		Posentrymode:  "992", // todo 如果是3dsecure的，992；EMV的规范还没出
		Chcd:          ap.Chcd,
		Clisn:         ap.CliSN,
		Mchntid:       ap.Mchntid,
		Terminalid:    ap.TerminalId,
		Txamt:         fmt.Sprintf("%012d", ap.ApplePayData.TransactionAmount),
		Txcurrcd:      ap.ApplePayData.CurrencyCode,
		Cardcd:        ap.ApplePayData.ApplicationPrimaryAccountNumber,
		Expiredate:    ap.ApplePayData.ApplicationExpirationDate[0:4],
		Syssn:         ap.SysSN,
		Localdt:       tools.LocalDt(),
		Transactionid: ap.TransactionId,
	}

	if ap.ApplePayData.PaymentDataType == "EMV" {
		// EMV 支付数据类型
		// m.Posentrymode = ""
		m.Iccdata = ap.ApplePayData.PaymentData.EmvData
	} else {
		// 3DSecure 支付数据类型
		// 3D交易发卡行验证结果转换:'5,6,7' ==> '05,06,07'
		m.EciIndicator = "0" + ap.ApplePayData.PaymentData.EciIndicator
		m.Onlinesecuredata = ap.ApplePayData.PaymentData.OnlinePaymentCryptogram
	}

	log.Infof("～～～～～～Apple Pay请求信息: %+v", m)

	resp := send(&m)

	if resp.Respcd != reversalFlag {
		// 应答码转换
		ret = transformResp(resp.Respcd)
		return
	}
	// 冲正处理
	reversalHandle(&m)

	// 返回‘外部系统错误’的应答码
	ret = mongo.RespCodeColl.Get("000002")
	return
}

// 消费撤销
func ConsumeUndo() {

}

// 消费冲正
func ConsumeReversal() {

}
