package core

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"net/url"
	"strconv"
)

// ProcessAlpNotify 支付宝异步通知处理
func ProcessAlpNotify(params url.Values) {

	// 通知动作类型
	notifyAction := params.Get("notify_action_type")
	// 交易订单号
	orderNum := params.Get("out_trade_no")
	// 系统订单号
	sysOrderNum := params.Get("schema")

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
		mongo.SpTransColl.UpdateFields(&model.Trans{
			Id:           t.Id,
			MerDiscount:  "0.00",
			ChanDiscount: "0.00",
		})
	// 其他
	default:
		// TODO 是否需要校验
		bills := params.Get("paytools_pay_amount")
		if bills != "" {
			var merDiscount float64
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
			// 更新指定字段，注意，这里不能全部更新
			// 否则可能会覆盖同步返回的结果
			mongo.SpTransColl.UpdateFields(&model.Trans{
				Id:          t.Id,
				MerDiscount: fmt.Sprintf("%0.2f", merDiscount),
			})
		}
	}

	// TODO 可能需要通知接入方
}

// ProcessWeixinNotify 微信异步通知处理(预下单)
func ProcessWeixinNotify(req *model.WeixinNotifyReq) {
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
	// TODO 这里会有并发问题
	switch req.ResultCode {
	case "SUCCESS":
		updateTrans(t, &model.ScanPayResponse{
			ChanRespCode:    req.ResultCode,
			ChannelOrderNum: req.TransactionId,
			ConsumerAccount: req.OpenID,
			Respcd:          "00",
		})
	default:
		updateTrans(t, &model.ScanPayResponse{
			ChanRespCode: req.ErrCode,
			Respcd:       "01",
		})
	}

	// TODO 可能需要通知接入方
}
