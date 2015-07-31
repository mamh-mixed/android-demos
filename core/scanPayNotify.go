package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

// ProcessAlpNotify 支付宝异步通知处理
// 该接口只接受预下单的异步通知
// 支付宝的其它接口将不接受异步通知
func ProcessAlpNotify(params url.Values) {

	log.Infof("alp async notify: %+v", params)
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
		return
	}

	// 判断是否是原订单
	if t.OrderNum != orderNum {
		log.Errorf("orderNum not match, expect %s, but get %s", t.OrderNum, orderNum)
		return
	}

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
			updateTrans(t, &model.ScanPayResponse{
				ChanRespCode:    tradeStatus,
				ChannelOrderNum: tradeNo,
				ConsumerAccount: account,
				MerDiscount:     fmt.Sprintf("%0.2f", merDiscount),
				Respcd:          adaptor.SuccessCode,
				ErrorDetail:     adaptor.SuccessMsg,
			})
		case "WAIT_BUYER_PAY":
			// do nothing
		default:
			updateTrans(t, returnWithErrorCode(tradeStatus))
		}
	}

	// TODO 可能需要通知接入方
}

// ProcessWeixinNotify 微信异步通知处理(预下单)
func ProcessWeixinNotify(req *weixin.WeixinNotifyReq) {
	log.Infof("weixin paut async request: %#v", req)

	// 上送的订单号
	sysOrderNum := req.Attach
	// 系统订单号是全局唯一
	t, err := mongo.SpTransColl.FindByOrderNum(sysOrderNum)
	if err != nil {
		log.Errorf("fail to find trans by sysOrderNum=%s", sysOrderNum)
		return
	}

	// 判断是否是原订单
	if t.OrderNum != req.OutTradeNo {
		log.Errorf("orderNum not match, expect %s, but get %s", t.OrderNum, req.OutTradeNo)
		return
	}

	// 更新交易信息
	switch req.ResultCode {
	case "SUCCESS":
		updateTrans(t, &model.ScanPayResponse{
			ChanRespCode:    req.ResultCode,
			ChannelOrderNum: req.TransactionId,
			ConsumerAccount: req.OpenID,
			Respcd:          adaptor.SuccessCode,
			ErrorDetail:     adaptor.SuccessMsg,
		})
	default:
		updateTrans(t, returnWithErrorCode(req.ErrCode))
	}

	// TODO 可能需要通知接入方
}
