package core

import (
	// "bytes"
	"encoding/json"
	"fmt"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/omigo/log"
	// "io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ProcessAlipayNotify 支付宝异步通知处理
// 该接口只接受预下单的异步通知
// 支付宝的其它接口将不接受异步通知
func ProcessAlipayNotify(params url.Values) error {
	// 通知动作类型
	notifyAction := params.Get("notify_action_type")
	// 交易订单号
	orderNum := params.Get("out_trade_no")
	// 系统订单号
	sysOrderNum := params.Get("extra_common_param")

	// 系统订单号是全局唯一
	t, err := mongo.SpTransColl.FindByOrderNum(sysOrderNum)
	if err != nil {
		log.Errorf("fail to find trans by sysOrderNum=%s", sysOrderNum)
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

	// 解锁
	defer func() {
		if t.LockFlag == 1 {
			mongo.SpTransColl.Unlock(t.MerId, t.OrderNum)
		}
	}()

	ret := &model.ScanPayResponse{}
	switch notifyAction {
	// 退款
	case "refundFPAction":
		// 将优惠信息更新为0.00，貌似为了打单用
		// mongo.SpTransColl.UpdateFields(&model.Trans{
		// 	Id:           t.Id,
		// 	MerDiscount:  "0.00",
		// 	ChanDiscount: "0.00",
		// })

	// 预下单时支付异步通知
	// TODO 是否需要校验
	default:
		bills := params.Get("paytools_pay_amount")
		tradeStatus := params.Get("trade_status")
		tradeNo := params.Get("trade_no")
		account := params.Get("buyer_email")

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

		case "WAIT_BUYER_PAY":
			log.Errorf("alp notify return tradeStatus: WAIT_BUYER_PAY, sysOrderNum=%s", sysOrderNum)
			return fmt.Errorf("%s", "transStatus no change")
		default:
			ret = adaptor.ReturnWithErrorCode("FAIL")
			ret.ChanRespCode = tradeStatus
		}
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
			if err = mongo.SpTransColl.UpdateAndUnlock(t); err != nil {
				return err
			}
		}
	}

	reqBytes, _ := json.Marshal(params)
	// 记录
	nr := &model.NotifyRecord{
		MerId:       t.MerId,
		OrderNum:    t.OrderNum,
		FromChanMsg: string(reqBytes),
	}
	// 可能需要通知接入方
	if t.NotifyUrl != "" {
		sendNotifyToMerchant(t, nr, ret)
	}
	mongo.NotifyRecColl.Add(nr)
	return nil
}

// ProcessWeixinNotify 微信异步通知处理(预下单)
func ProcessWeixinNotify(req *weixin.WeixinNotifyReq) error {
	// 上送的订单号
	sysOrderNum := req.Attach
	// 系统订单号是全局唯一
	t, err := mongo.SpTransColl.FindByOrderNum(sysOrderNum)
	if err != nil {
		log.Errorf("fail to find trans by sysOrderNum=%s", sysOrderNum)
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
		ret.ConsumerAccount = req.OpenID
		ret.Respcd = adaptor.SuccessCode
		ret.ErrorDetail = adaptor.SuccessMsg
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
	// 可能需要通知接入方
	if t.NotifyUrl != "" {
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

	// bs, err := json.Marshal(ret)
	// if err != nil {
	// 	log.Errorf("json marshal error: %s", err)
	// 	return
	// }
	nr.ToMerMsg = parms
	notifyUrl := t.NotifyUrl
	if strings.Contains(notifyUrl, "?") {
		notifyUrl += "&" + parms
	} else {
		notifyUrl += "?" + parms
	}
	log.Infof("send notify: %s", notifyUrl)

	go func() {
		var interval = []time.Duration{15, 15, 30, 180, 1800, 1800, 1800, 1800, 3600, 0}
		// var interval = []time.Duration{1, 1, 1, 1, 1, 1, 1, 1, 1, 0} // for test
		for i, d := range interval {
			log.Infof("merId=%s,orderNum=%s,url=%s, send notify %d times", ret.Mchntid, ret.OrderNum, t.NotifyUrl, i+1)
			// resp, err := http.Post(t.NotifyUrl, "application/json", strings.NewReader(parms))
			resp, err := http.Get(notifyUrl)
			if err != nil {
				time.Sleep(time.Second * d)
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				time.Sleep(time.Second * d)
				continue
			}
			// rs, err := ioutil.ReadAll(resp.Body)
			// if err != nil {
			// 	time.Sleep(time.Second * d)
			// 	continue
			// }
			// clientResp := &model.ScanPayResponse{}
			// err = json.Unmarshal(rs, clientResp)
			// if err != nil {
			// 	time.Sleep(time.Second * d)
			// 	continue
			// }
			// if clientResp.Respcd != adaptor.SuccessCode {
			// 	time.Sleep(time.Second * d)
			// 	continue
			// }
			// 异步通知成功，返回
			return
		}

		// 异步通知失败
		nr.IsToMerFail = true
		mongo.NotifyRecColl.Update(nr)
	}()
}

func copyNotifyProperties(ret *model.ScanPayResponse, t *model.Trans) {
	// ret.Busicd = "NOTI"
	ret.Txndir = "A"
	ret.Busicd = model.Paut
	ret.Mchntid = t.MerId
	ret.AgentCode = t.AgentCode
	ret.Terminalid = t.Terminalid
	ret.OrderNum = t.OrderNum
	ret.Chcd = t.ChanCode
	ret.Txamt = fmt.Sprintf("%012d", t.TransAmt)
	// ret.MerDiscount = t.MerDiscount
	// ret.ChcdDiscount = t.ChanDiscount
	// ret.QrCode = t.QrCode
}
