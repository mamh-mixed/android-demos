package domestic

import (
	"errors"
	"fmt"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	// "github.com/omigo/mahonia"
	"strconv"
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
		OutTradeNo:     req.OrderNum,    // 送的是原订单号，不转换
		PassbackParams: req.SysOrderNum, // 传系统订单号，异步通知时可用
		Subject:        req.Subject,
		GoodsDetail:    req.AlpMarshalGoods(),
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
		GoodsDetail:    req.AlpMarshalGoods(),
		PassbackParams: req.SysOrderNum, // 传系统订单号，异步通知时可用
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
func (a *alp) ProcessSettleEnquiry(req *model.ScanPayRequest, modelMMap map[string]map[string][]model.BlendElement) error {

	if modelMMap == nil {
		return errors.New("the map is nil")
	}

	alpReq := &alpRequest{
		Partner:          req.ChanMerId,
		Service:          settleService,
		Gmt_create_start: req.StartTime,
		Gmt_create_end:   req.EndTime,
	}

	alpReq.SpReq = req

	alpRsp, err := sendSettleRequest(alpReq)

	if err != nil {
		log.Errorf("sendRequest fail, func name is ProcessSettleEnquiry, the settle time start:%s, end:%s, error:%s", req.StartTime, req.EndTime, err)
		return err
	}

	err1 := analysisSettleData(alpRsp.Response.Csv_result, alpReq.Partner, modelMMap)
	/*
		for _, array := range modelMap {
			for _, element := range array {
				fmt.Println("the value is:", element)
			}
		}
	*/
	if err1 != nil {
		log.Errorf("settle data change to map error:%s", err)
		return err1
	}

	return err
}

func analysisSettleData(csvData csv_detail, chanMer string, modelMMap map[string]map[string][]model.BlendElement) error { //key订单号
	dataStr := csvData.Csv_data
	count, err := strconv.Atoi(csvData.Count)
	if err != nil {
		log.Errorf("change data count errDetail:%s", err)
		return err
	}
	istart := strings.LastIndex(dataStr, "[")
	iend := strings.Index(dataStr, "]")
	dataArray := []byte(dataStr)
	stemp := string(dataArray[istart+1 : iend])
	element := strings.Split(stemp, ",")
	//检查要取关键位置是否变化 如：外部订单号,账户余额（元）,时间,流水号,支付宝交易号,交易对方Email,交易对方,用户编号,收入（元）,支出（元）,交易场所,商品名称,类型,说明,
	if element[0] != "外部订单号" {
		log.Errorf("the first position is different")
		err = errors.New("the first position is different")
		return err
	}
	if element[2] != "时间" {
		log.Errorf("the third position is different")
		err = errors.New("the third position is different")
		return err
	}
	if element[4] != "支付宝交易号" {
		log.Errorf("the fifth position is different")
		err = errors.New("the fifth position is different")
		return err
	}
	if element[8] != "收入（元）" {
		log.Errorf("the eighth position is different")
		err = errors.New("the eighth position is different")
		return err
	}
	if element[9] != "支出（元）" {
		log.Errorf("the ninth position is different")
		err = errors.New("the ninth position is different")
		return err
	}
	if element[12] != "类型" {
		log.Errorf("the twelve position is different")
		err = errors.New("the twelve position is different")
		return err
	}

	//elementArray := make([]model.BlendElement)
	//mmap := make(map[string]map[string][]model.BlendElement)

	modelMap, ret := modelMMap[chanMer]
	if !ret {
		modelMap = make(map[string][]model.BlendElement)
	}

	i := 0
	for i < count {
		var elementModel model.BlendElement
		//element[14*(i+1)]  //订单号
		elementModel.Chcd = "ALP"
		elementModel.ChcdName = "支付宝"
		elementModel.ChanMerID = chanMer
		elementModel.OrderTime = element[14*(i+1)+2] //时间
		elementModel.OrderID = element[14*(i+1)+4]   //支付宝交易号
		elementModel.IsBlend = false
		elementModel.OrderType = element[14*(i+1)+12]
		if elementModel.OrderType == "在线支付" {
			elementModel.OrderAct = element[14*(i+1)+8] //收入
		} else if elementModel.OrderType == "交易退款" {
			elementModel.OrderAct = element[14*(i+1)+9] //支出
		}

		//append(elementArray, elementModel)
		elementArray, ret := modelMap[elementModel.OrderID]
		if !ret {
			elementArray = make([]model.BlendElement, 0)
		}
		elementArray = append(elementArray, elementModel)
		modelMap[elementModel.OrderID] = elementArray

		i += 1
	}

	//modelMMap := make(map[string]map[string][]model.BlendElement)
	modelMMap[chanMer] = modelMap

	return nil
}
