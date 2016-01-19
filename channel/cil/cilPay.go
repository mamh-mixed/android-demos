package cil

import (
	"fmt"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

// 讯联交易类型
const (
	version               = "15.1"
	consumeBusicd         = "000000" // 消费
	orderConsumeBusicd    = "020000" // 订购消费
	consumeUndoBusicd     = "201000" // 消费撤销
	consumeReversalBusicd = "040000" // 消费冲正

	reversalFlag = "TIME_OUT" // 冲正标识
)

var (
	DefaultCILPayClient CILPay             // 线下网关交易的入口
	transTimeout        = 50 * time.Second // 超时时间
	reversalTimeouts    = [...]time.Duration{transTimeout, transTimeout * 1, transTimeout * 8, transTimeout * 50, transTimeout * 1140}
)

// CILPay 表示线下网关的支付对象
type CILPay struct{}

// Consume 直接消费（订购消费）
func (c *CILPay) Consume(p *model.NoTrackPayment) (ret *model.BindingReturn) {
	// 构建消费报文
	m := &model.CilMsg{
		Busicd:       orderConsumeBusicd,
		Txndir:       "Q",
		Posentrymode: "012",
		Chcd:         p.Chcd,
		Clisn:        p.CliSN,
		Mchntid:      p.Mchntid,
		Terminalid:   p.TerminalId,
		Txamt:        fmt.Sprintf("%012d", p.TransAmt),
		Txcurrcd:     p.CurrCode,
		Cardcd:       p.AcctNumDecrypt,
		Syssn:        p.SysSN,
		Localdt:      time.Now().Format("0102150405"),
		Expiredate:   p.ValidDateDecrypt,
		Cvv2:         p.Cvv2Decrypt,
	}
	if len(p.MerOrderNum) > 20 {
		m.Transactionid = p.MerOrderNum[len(p.MerOrderNum)-20:]
	} else {
		m.Transactionid = p.MerOrderNum
	}

	// 报文入库
	// m.UUID = util.SerialNumber()
	// log.Debugf("直接消费（订购消费）向线下网关发送报文内容: %+v", m)
	// mongo.CilMsgColl.Upsert(m)

	resp := send(m, transTimeout)
	log.Debugf("直接消费（订购消费）的线下网关返回结果: %+v", resp)
	if resp == nil {
		return mongo.RespCodeColl.Get("000002")
	}
	// 如果超时，请冲正
	if resp.Respcd == reversalFlag {
		log.Warn("请求超时!!!")
		// 另起线程，冲正处理
		go reversalHandle(m)
		// 返回‘外部系统错误’的应答码
		ret = mongo.RespCodeColl.Get("000002")
		return
	}

	// 应答码转换
	ret = transformResp(resp.Respcd)

	// 更新已存储的报文
	m.Respcd = resp.Respcd
	// mongo.CilMsgColl.Upsert(m)

	return
}

// ConsumeByApplePay ApplePay 消费
func (c *CILPay) ConsumeByApplePay(ap *model.ApplePay) (ret *model.BindingReturn) {
	m := &model.CilMsg{
		Busicd:        consumeBusicd,
		Txndir:        "Q",
		Posentrymode:  "992", // todo 如果是 3DSecure 的，992；EMV的规范还没出
		Chcd:          ap.Chcd,
		Clisn:         ap.CliSN,
		Mchntid:       ap.Mchntid,
		Terminalid:    ap.TerminalId,
		Txamt:         fmt.Sprintf("%012d", ap.ApplePayData.TransactionAmount),
		Txcurrcd:      ap.ApplePayData.CurrencyCode,
		Cardcd:        ap.ApplePayData.ApplicationPrimaryAccountNumber,
		Expiredate:    ap.ApplePayData.ApplicationExpirationDate[0:4],
		Syssn:         ap.SysSN,
		Localdt:       time.Now().Format("0102150405"),
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

	resp := send(m, transTimeout)
	log.Debugf("ApplePay 消费的线下网关返回结果: %+v", resp)

	if resp.Respcd == reversalFlag {
		log.Warn("请求超时!!!")
		// 另起线程，冲正处理
		go reversalHandle(m)
		// 返回‘外部系统错误’的应答码
		ret = mongo.RespCodeColl.Get("000002")
		return
	}

	// 应答码转换
	ret = transformResp(resp.Respcd)

	// 更新已存储的报文
	m.Respcd = resp.Respcd
	// mongo.CilMsgColl.Upsert(m)

	return
}

// ConsumeUndo 消费撤销
func ConsumeUndo() {

}

// reversalHandle 冲正处理方法
func reversalHandle(om *model.CilMsg) {
	log.Warn("源交易请求超时，发送冲正报文")
	// 创建冲正报文
	rm := &model.CilMsg{
		Busicd:       consumeReversalBusicd,
		Txndir:       "Q",
		Posentrymode: om.Posentrymode,
		Chcd:         om.Chcd,
		Clisn:        mongo.SnColl.GetDaySN(om.Mchntid, om.Terminalid),
		Mchntid:      om.Mchntid,
		Terminalid:   om.Terminalid,
		Txamt:        om.Txamt,
		Txcurrcd:     om.Txcurrcd,
		Cardcd:       om.Cardcd,
		Syssn:        om.Syssn,
		Origclisn:    om.Clisn,
		Localdt:      time.Now().Format("0102150405"),
	}

	for _, i := range reversalTimeouts {
		log.Debugf("Send reversal request, overtime is %s", i)

		back := send(rm, i)
		if back.Respcd != reversalFlag {
			log.Info("reversal operation success")

			return
		}
	}

	log.Errorf("冲正失败,报文数据是：%+v", rm)
}
