package core

import (
	// "bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/logs"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/log"
)

// ProcessAlipayNotify 支付宝异步通知处理，接受预下单和下单异步通知
func ProcessAlipayNotify(params url.Values) error {

	// 通知动作类型
	notifyAction := params.Get("notify_action_type")
	// 交易订单号
	orderNum := params.Get("out_trade_no")
	// 系统订单号、异步通知日志关联ID
	sysOrderNum, logReqId := parseAttach(params.Get("extra_common_param"))

	// 系统订单号是全局唯一
	t, err := mongo.SpTransColl.FindByOrderNum(sysOrderNum)
	if err != nil {
		log.Errorf("fail to find trans by sysOrderNum=%s, error: %s", sysOrderNum, err)
		return err
	}

	count, err := mongo.NotifyRecColl.Count(t.MerId, t.OrderNum)
	if err != nil {
		return err
	}

	// 有此记录，表示已经处理过
	if count > 0 {
		return nil
	}

	// 判断是否是原订单
	if t.OrderNum != orderNum {
		log.Errorf("orderNum not match, expect %s, but get %s", t.OrderNum, orderNum)
		return fmt.Errorf("%s", "orderNum error")
	}

	// 锁住交易
	t, err = findAndLockTrans(t.MerId, t.OrderNum)
	if err != nil {
		return err
	}

	// 异步通知数据进入日志
	logs.SpLogs <- &model.SpTransLogs{
		ReqId:        logReqId,
		Direction:    "in",
		MerId:        t.MerId,
		OrderNum:     t.OrderNum,
		OrigOrderNum: t.OrigOrderNum,
		TransType:    t.Busicd,
		MsgType:      3,
		Msg:          params}

	// 解锁
	defer func() {
		if t.LockFlag == 1 {
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
		}
	}()

	ret := &model.ScanPayResponse{}
	bills := params.Get("paytools_pay_amount")
	tradeStatus := params.Get("trade_status")
	tradeNo := params.Get("trade_no")
	account := params.Get("buyer_email")
	payTime := params.Get("gmt_payment")

	// 折扣
	var merDiscount float64
	if bills != "" {
		var arrayBills []map[string]string
		if err := json.Unmarshal([]byte(bills), &arrayBills); err == nil {
			for _, bill := range arrayBills {
				for k, v := range bill {
					if k == "MCOUPON" || k == "MDISCOUNT" {
						f, _ := strconv.ParseFloat(v, 64)
						merDiscount += f
					}
				}
			}
		}
	}
	// 交易状态更新
	switch tradeStatus {
	case "TRADE_SUCCESS":
		ret.ChanRespCode = tradeStatus
		ret.ChannelOrderNum = tradeNo
		ret.ConsumerAccount = account
		ret.Respcd = adaptor.SuccessCode
		ret.ErrorDetail = adaptor.SuccessMsg
		ret.ErrorCode = "SUCCESS"
		ret.MerDiscount = fmt.Sprintf("%0.2f", merDiscount)
		ret.PayTime = payTime
	case "WAIT_BUYER_PAY":
		log.Errorf("alp notify return tradeStatus: WAIT_BUYER_PAY, sysOrderNum=%s", sysOrderNum)
		return fmt.Errorf("%s", "transStatus no change")
	default:
		ret = adaptor.ReturnWithErrorCode("FAIL")
		ret.ChanRespCode = tradeStatus
	}

	if t.TransStatus != model.TransClosed {
		if err = updateTrans(t, ret); err != nil {
			// 如果更新失败，则认为没有处理过
			return err
		}
	} else {
		// 订单已关闭，但是又收到异步通知是成功
		// 表明订单已被退款
		if ret.Respcd == adaptor.SuccessCode {
			t.RefundStatus = model.TransRefunded
			t.Fee, t.NetFee = 0, 0
			t.RespCode = adaptor.SuccessCode
			t.ErrorDetail = adaptor.SuccessMsg
			t.ConsumerAccount = ret.ConsumerAccount
			t.ChanOrderNum = ret.ChannelOrderNum
			t.PayTime = payTime
			if err = mongo.SpTransColl.UpdateAndUnlock(t); err != nil {
				return err
			}
		}
	}

	// 为空时为预下单
	if notifyAction != "" {
		return nil
	}

	reqBytes, _ := json.Marshal(params)
	// 记录
	nr := &model.NotifyRecord{
		MerId:       t.MerId,
		OrderNum:    t.OrderNum,
		FromChanMsg: string(reqBytes),
	}
	// 订单状态正常且有填写通知地址
	if t.NotifyUrl != "" && t.TransStatus == model.TransSuccess {
		sendNotifyToMerchant(t, nr, ret)
	}
	mongo.NotifyRecColl.Add(nr)
	return nil
}

// ProcessWeixinNotify 微信异步通知处理(预下单)
func ProcessWeixinNotify(req *weixin.WeixinNotifyReq) error {
	// 上送的订单号、日志关联ID
	sysOrderNum, logReqId := parseAttach(req.Attach)
	// 系统订单号是全局唯一
	t, err := mongo.SpTransColl.FindByOrderNum(sysOrderNum)
	if err != nil {
		log.Errorf("fail to find trans by sysOrderNum=%s, error: %s", sysOrderNum, err)
		return err
	}

	count, err := mongo.NotifyRecColl.Count(t.MerId, t.OrderNum)
	if err != nil {
		return err
	}

	// 有此记录，表示已经处理过
	if count > 0 {
		return nil
	}

	// 判断是否是原订单
	if t.OrderNum != req.OutTradeNo {
		log.Errorf("orderNum not match, expect %s, but get %s", t.OrderNum, req.OutTradeNo)
		return fmt.Errorf("%s", "orderNum error")
	}

	// 锁住交易
	t, err = findAndLockTrans(t.MerId, t.OrderNum)
	if err != nil {
		return err
	}

	// 异步通知数据进入日志
	logs.SpLogs <- &model.SpTransLogs{
		ReqId:        logReqId,
		Direction:    "in",
		MerId:        t.MerId,
		OrderNum:     t.OrderNum,
		OrigOrderNum: t.OrigOrderNum,
		TransType:    t.Busicd,
		MsgType:      3,
		Msg:          req,
	}

	// 解锁
	defer func() {
		if t.LockFlag == 1 {
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
		}
	}()

	ret := &model.ScanPayResponse{}
	// 更新交易信息
	switch req.ResultCode {
	case "SUCCESS":
		ret.ChanRespCode = req.ResultCode
		ret.ChannelOrderNum = req.TransactionId
		if req.SubOpenid != "" {
			ret.ConsumerAccount = req.SubOpenid
		} else {
			ret.ConsumerAccount = req.OpenID
		}
		ret.Respcd = adaptor.SuccessCode
		ret.ErrorDetail = adaptor.SuccessMsg
		ret.PayTime = req.TimeEnd
		ret.ErrorCode = "SUCCESS"
	default:
		ret = adaptor.ReturnWithErrorCode("FAIL")
		ret.ChanRespCode = req.ResultCode
	}

	// 如果交易状态不是关闭，才去更新，否则认为该笔交易已被正确关闭
	// 该异步通知只送往商户，不做更新。
	if t.TransStatus != model.TransClosed {
		if err = updateTrans(t, ret); err != nil {
			// 如果更新失败，则认为没有处理过
			return err
		}
	}

	reqBytes, _ := json.Marshal(req)
	// 记录
	nr := &model.NotifyRecord{
		MerId:       t.MerId,
		OrderNum:    t.OrderNum,
		FromChanMsg: string(reqBytes),
	}
	// 订单状态正常且有填写通知地址
	if t.NotifyUrl != "" && t.TransStatus == model.TransSuccess {
		sendNotifyToMerchant(t, nr, ret)
	}
	mongo.NotifyRecColl.Add(nr)

	return nil
}

// sendNotifyToMerchant 向接入方发送异步消息通知
func sendNotifyToMerchant(t *model.Trans, nr *model.NotifyRecord, ret *model.ScanPayResponse) {

	// 从交易中复制要发送给商户的数据
	copyNotifyProperties(ret, t)

	mer, err := mongo.MerchantColl.Find(ret.Mchntid)
	if err != nil {
		log.Errorf("send notify error: %s", err)
		return
	}
	// 发送异步消息采用英文描述
	ret.ErrorDetail = ret.ErrorCode
	ret.ErrorCode = ""

	// 签名
	signContent, sign := ret.SignMsg(), ""
	if mer.IsNeedSign {
		log.Debug("send notify sign content to return : " + signContent)
		sign = security.SHA1WithKey(signContent, mer.SignKey)
	}
	parms := signContent + "&sign=" + sign

	nr.ToMerMsg = parms
	notifyUrl := t.NotifyUrl
	if strings.Contains(notifyUrl, "?") {
		notifyUrl += "&" + parms
	} else {
		notifyUrl += "?" + parms
	}

	go func() {
		var interval = []time.Duration{15, 15, 30, 180, 1800, 1800, 1800, 1800, 3600, 0}
		// var interval = []time.Duration{1, 1, 1, 1, 1, 1, 1, 1, 1, 0} // for test
		for i, d := range interval {
			// 跳过第一次
			if i != 0 {
				nowT, err := mongo.SpTransColl.FindOne(t.MerId, t.OrderNum)
				if err != nil {
					log.Errorf("find trans error: %s", err)
					time.Sleep(time.Second * d)
					continue
				}

				// 假如在发送时，交易的状态改变了
				if nowT.TransStatus != t.TransStatus {
					log.Warnf("merId=%s, orderNum=%s, trans status change from %s to %s", t.MerId, t.OrderNum, t.TransStatus, nowT.TransStatus)
					nr.Remark = "transStatus_change"
					break
				}
			}

			// resp, err := http.Post(t.NotifyUrl, "application/json", strings.NewReader(parms))
			resp, err := http.Get(notifyUrl)
			if err != nil {
				log.Warnf("send notify %d times fail, merId=%s, orderNum=%s, channelOrderNum=%s, notifyUrl=%s: %s",
					i+1, ret.Mchntid, ret.OrderNum, ret.ChannelOrderNum, notifyUrl, err)
				time.Sleep(time.Second * d)
				continue
			}
			log.Infof("send notify %d times successfully, merId=%s, orderNum=%s, channelOrderNum=%s, notifyUrl=%s",
				i+1, ret.Mchntid, ret.OrderNum, ret.ChannelOrderNum, notifyUrl)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Infof("remote service return http status: %d", resp.StatusCode)
				time.Sleep(time.Second * d)
				continue
			}

			// 异步通知成功，返回
			return
		}

		// 异步通知失败
		nr.IsToMerFail = true
		mongo.NotifyRecColl.Update(nr)
	}()
}

func copyNotifyProperties(ret *model.ScanPayResponse, t *model.Trans) {
	ret.Txndir = "A"
	ret.Busicd = model.Paut
	ret.Mchntid = t.MerId
	ret.AgentCode = t.AgentCode
	ret.Terminalid = t.Terminalid
	ret.OrderNum = t.OrderNum
	ret.Chcd = t.ChanCode
	ret.Txamt = fmt.Sprintf("%012d", t.TransAmt)
	ret.Attach = t.Attach
	ret.GoodsInfo = t.GoodsInfo
}

func parseAttach(attach string) (sysOrderNum string, logReqId string) {
	data := strings.Split(attach, ",")
	switch len(data) {
	case 0:
		// ignore
	case 1:
		sysOrderNum = data[0]
	case 2:
		sysOrderNum = data[0]
		logReqId = data[1]
	}
	return
}
