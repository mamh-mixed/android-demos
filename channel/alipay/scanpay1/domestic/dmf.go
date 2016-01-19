package domestic

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"strings"
	"time"
)

var DefaultClient alp
var NotifyPath = "/scanpay/upNotify/alipay"
var alipayNotifyUrl = goconf.Config.AlipayScanPay.NotifyUrl + NotifyPath

// alp 当面付，扫码支付
type alp struct{}

// service
const (
	createAndPay  = "alipay.acquire.createandpay"
	preCreate     = "alipay.acquire.precreate"
	refund        = "alipay.acquire.refund"
	query         = "alipay.acquire.query"
	cancel        = "alipay.acquire.cancel"
	settleService = "export_trade_account_report"
)

// params
const (
	// partner  = "2088811767473826" // only for test
	charSet  = "utf-8"
	currency = "156"
)

// ProcessBarcodePay 条码支付/下单
func (a *alp) ProcessBarcodePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:        req.ChanMerId,
		Service:        createAndPay,
		NotifyUrl:      alipayNotifyUrl,
		OutTradeNo:     req.OrderNum,                      // 送的是原订单号，不转换
		PassbackParams: req.SysOrderNum + "," + req.ReqId, // 格式：系统订单号,日志Id
		Subject:        req.Subject,
		GoodsDetail:    parseGoods(req),
		ProductCode:    "BARCODE_PAY_OFFLINE",
		TotalFee:       req.ActTxamt,
		ExtendParams:   req.ExtendParams, //...
		ItBPay:         "1d",             // 超时时间
		DynamicIdType:  "bar_code",
		DynamicId:      req.ScanCodeId,
		SpReq:          req,
	}

	alpResp, err := sendRequest(alpReq)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, createAndPay)
		return nil, err
	}

	// 处理结果返回
	return transform(alpReq.Service, alpResp)
}

// ProcessQrCodeOfflinePay 扫码支付/预下单
func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:        req.ChanMerId,
		Service:        preCreate,
		NotifyUrl:      alipayNotifyUrl,
		OutTradeNo:     req.OrderNum, // 送的是原订单号，不转换,
		Subject:        req.Subject,
		GoodsDetail:    parseGoods(req),
		PassbackParams: req.SysOrderNum + "," + req.ReqId, // 格式：系统订单号,日志Id
		ProductCode:    "QR_CODE_OFFLINE",
		TotalFee:       req.ActTxamt,
		ExtendParams:   req.ExtendParams,
		ItBPay:         handleItBpay(req.TimeExpire), // 超时时间
		SpReq:          req,
	}

	alpResp, err := sendRequest(alpReq)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, preCreate)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:      req.ChanMerId,
		Service:      refund,
		NotifyUrl:    req.NotifyUrl,
		OutTradeNo:   req.OrigOrderNum,
		RefundAmount: req.ActTxamt,
		OutRequestNo: req.OrderNum, //该字段上送才能部分退款，如果不送则只能全额退款
		SpReq:        req,
	}

	alpResp, err := sendRequest(alpReq)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, refund)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessEnquiry 查询，包含支付、退款
func (a *alp) ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:    req.ChanMerId,
		Service:    query,
		OutTradeNo: req.OrigOrderNum, // 送的是原订单号，不转换
		SpReq:      req,
	}

	alpResp, err := sendRequest(alpReq)
	if err != nil {
		log.Errorf("sendRequest fail, origOrderNum=%s, service=%s, channel=alp", req.OrigOrderNum, query)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessCancel 撤销
func (a *alp) ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	alpReq := &alpRequest{
		Partner:    req.ChanMerId,
		Service:    cancel,
		NotifyUrl:  req.NotifyUrl,
		OutTradeNo: req.OrigOrderNum,
		SpReq:      req,
	}

	alpResp, err := sendRequest(alpReq)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=%s, channel=alp", req.OrderNum, cancel)
		return nil, err
	}

	return transform(alpReq.Service, alpResp)
}

// ProcessClose 关闭接口即撤销接口
func (a *alp) ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return a.ProcessCancel(req)
}

// ProcessRefundQuery 退款查询
func (a *alp) ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, errors.New("not support now")
}

func toMap(req *alpRequest) map[string]string {

	dict := make(map[string]string)
	// 固定参数
	dict["_input_charset"] = charSet
	dict["partner"] = req.Partner
	dict["currency"] = currency
	dict["seller_id"] = req.Partner
	// 参数转换
	dict["service"] = req.Service
	dict["notify_url"] = req.NotifyUrl
	dict["product_code"] = req.ProductCode
	dict["out_trade_no"] = req.OutTradeNo
	dict["total_fee"] = req.TotalFee
	dict["it_b_pay"] = req.ItBPay
	dict["dynamic_id_type"] = req.DynamicIdType
	dict["dynamic_id"] = req.DynamicId
	dict["refund_amount"] = req.RefundAmount
	dict["out_request_no"] = req.OutRequestNo
	dict["passback_parameters"] = req.PassbackParams

	// utf-8 -> gbk
	// e := mahonia.NewEncoder("gbk")
	if req.Subject != "" {
		dict["subject"] = req.Subject
		// dict["subject"] = e.ConvertString(req.Subject)
	}
	if req.GoodsDetail != "" {
		dict["goods_detail"] = req.GoodsDetail
		// dict["goods_detail"] = e.ConvertString(req.GoodsDetail)
	}
	if req.ExtendParams != "" {
		dict["extend_params"] = req.ExtendParams
		// dict["extend_params"] = e.ConvertString(req.ExtendParams)
	}

	return dict
}

// handleItBpay 处理过期时间，默认为一天
func handleItBpay(timeExpired string) string {

	var defaultInterval = "1d"

	if timeExpired == "" {
		return defaultInterval
	}
	et, err := time.ParseInLocation("20060102150405", timeExpired, time.Local)
	if err != nil {
		log.Warnf("timeExpired(%s) format error:%s", timeExpired, err)
		return defaultInterval
	}

	now := time.Now()
	d := et.Sub(now)

	// 如果小于5分钟，按5分钟的来
	if d < 5*time.Minute {
		return "5m"
	}

	return fmt.Sprintf("%0.fm", d.Minutes()+1)
}

func toSettleMap(req *alpRequest) map[string]string {
	dict := make(map[string]string)

	dict["service"] = req.Service
	dict["partner"] = req.Partner
	dict["_input_charset"] = charSet
	dict["gmt_create_end"] = req.Gmt_create_end
	dict["gmt_create_start"] = req.Gmt_create_start

	return dict
}

//账单明细查询
func (a *alp) ProcessSettleEnquiry(req *model.ScanPayRequest, cbd model.ChanBlendMap) error {

	if cbd == nil {
		return fmt.Errorf("%s", "nil map found")
	}

	alpReq := &alpRequest{
		Partner:          req.ChanMerId,
		Service:          settleService,
		Gmt_create_start: req.SettDate + " 00:00:00",
		Gmt_create_end:   req.SettDate + " 23:59:59",
	}

	alpReq.SpReq = req

	alpResp, err := sendSettleRequest(alpReq)
	if err != nil {
		log.Errorf("sendRequest fail, error:%s", err)
		return err
	}

	if alpResp.Error != "" {
		return fmt.Errorf("%s", alpResp.Error)
	}

	analysisSettleData(alpResp.Response.CsvResult, alpReq.Partner, cbd)

	return nil
}

func analysisSettleData(csvData csvDetail, chanMerId string, cbd model.ChanBlendMap) {

	csv := csvData.CsvStr
	// 截取内容
	if strings.Contains(csv, "<![CDATA[") {
		csv = csv[9 : len(csv)-3]
		// log.Debugf("csv content: %s", csv)
	}

	// 检查要取关键位置是否变化 如：
	// 外部订单号,账户余额（元）,时间,流水号,支付宝交易号,交易对方Email,交易对方,用户编号,收入（元）,支出（元）,交易场所,商品名称,类型,说明,
	s := bufio.NewScanner(strings.NewReader(csv))
	// skip 1
	var data [][]string
	for i := 0; s.Scan(); i++ {
		if i > 0 {
			ts := strings.Split(s.Text(), ",")
			if len(ts) != 15 {
				log.Errorf("length invalid , skip, data=%s", ts)
				continue
			}
			data = append(data, ts)
		}
	}

	recsMap, ok := cbd[chanMerId]
	if !ok {
		recsMap = make(map[string][]model.BlendElement)
	}

	for _, d := range data {
		var elementModel model.BlendElement
		elementModel.Chcd = "ALP"
		elementModel.ChcdName = "支付宝"
		elementModel.LocalID = d[0]
		elementModel.OrderTime = d[2] // 时间
		elementModel.OrderID = d[4]   // 支付宝交易号
		elementModel.Account = d[5]
		elementModel.OrderAct = d[8] // 收入
		if elementModel.OrderAct == "" {
			elementModel.OrderAct = d[9] // 支出
		}
		elementModel.OrderType = d[12]

		if strings.Contains(elementModel.LocalID, "-r-") {
			elementModel.LocalID = strings.Split(elementModel.LocalID, "-r-")[0]
			// log.Debugf("退费交易: %+v", elementModel)
		}

		var key string
		switch elementModel.OrderType {
		case "收费", "退费":
			// 取的是外部订单号
			key = elementModel.LocalID
		default:
			key = elementModel.OrderID
		}

		if elementArray, ok := recsMap[key]; ok {
			elementArray = append(elementArray, elementModel)
			recsMap[key] = elementArray
		} else {
			recsMap[key] = []model.BlendElement{elementModel}
		}

	}
	log.Debugf("debug data: %+v", recsMap["2015122221001004770044413116"])
	cbd[chanMerId] = recsMap
}

// parseGoods 按照dmf1.0格式输出商品详情
func parseGoods(req *model.ScanPayRequest) string {
	goods, err := req.MarshalGoods()
	if err != nil {
		return ""
	}

	if len(goods) > 0 {
		gbs, err := json.Marshal(goods)
		if err != nil {
			log.Errorf("goodsInfo marshal error:%s", err)
			return ""
		}
		return string(gbs)
	}

	return ""
}
