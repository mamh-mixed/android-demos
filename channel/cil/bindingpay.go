package cil

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"strconv"
)

// 讯联交易类型
const (
	version               = "15.1"
	consumeBusicd         = "000000" // 消费
	consumeUndoBusicd     = "201000" // 消费撤销
	consumeReversalBusicd = "040000" // 消费冲正
)

// 消费
func Consume(p *model.NoTrackPayment) (ret *model.BindingReturn) {
	// m := CilMsg{
	// 	Busicd:        consumeBusicd,
	// 	Txndir:        "Q",
	// 	Posentrymode:  "992", // todo 如果是3dsecure的，992；EMV的规范还没出
	// 	Chcd:          "",
	// 	Clisn:         "",
	// 	Mchntid:       "",
	// 	Terminalid:    "",
	// 	Txamt:         "",
	// 	Txcurrcd:      "",
	// 	Cardcd:        "",
	// 	Syssn:         "",
	// 	Localdt:       tools.LocalDt(),
	// 	Transactionid: "",
	// }
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
		Mchntid:       ap.MerId,
		Terminalid:    ap.TerminalId,
		Txamt:         strconv.FormatInt(ap.ApplePayData.TransactionAmount, 10),
		Txcurrcd:      ap.ApplePayData.CurrencyCode,
		Cardcd:        ap.ApplePayData.ApplicationPrimaryAccountNumber,
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
