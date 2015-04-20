package cil

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

// 讯联交易类型
const (
	version               = "15.1"
	reversalTime          = 50       // 超时时间
	consumeBusicd         = "000000" // 消费
	orderConsumeBusicd    = "020000" // 订购消费
	consumeUndoBusicd     = "201000" // 消费撤销
	consumeReversalBusicd = "040000" // 消费冲正
)

// 无卡直接消费（订购消费）
func Consume(p *model.NoTrackPayment) (ret *model.BindingReturn) {
	log.Info("无卡直接支付开始向线下网关发送报文")
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

	resp := send(&m)
	log.Debugf("无卡直接支付返回结果:%+v", resp)
	return nil
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
		Expiredate:    "",
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
		m.Eclindicator = "0" + ap.ApplePayData.PaymentData.EciIndicator
		m.Onlinesecuredata = ap.ApplePayData.PaymentData.OnlinePaymentCryptogram
	}

	log.Info(m)

	resp := send(&m)

	if resp == nil {
		return nil
	}

	// 应答码转换
	ret = transformResp(resp.Respcd)
	return
}

// 消费撤销
func ConsumeUndo() {

}

// 消费冲正
func ConsumeReversal() {

}
