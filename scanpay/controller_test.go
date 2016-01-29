package scanpay

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

var (
	// 下单支付
	scanPayBarcodePay = &model.ScanPayRequest{
		GoodsInfo: "鞋子,1000.00,2;衣服,1500,3",
		OrderNum:  util.Millisecond(),
		// OrderNum:   "哈哈中文订单号",
		ScanCodeId: "287528698389366895",
		AgentCode:  "19992900",
		Txamt:      "000000000001",
		Chcd:       "AOS",
		Busicd:     "PURC",
		Currency:   "JPY",
		Mchntid:    "200000000010002",
		// Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}
	// 预下单支付
	scanPayQrCodeOfflinePay = &model.ScanPayRequest{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  util.Millisecond(),
		AgentCode: "10134001",
		Txamt:     "000000000001",
		Busicd:    "PAUT",
		Mchntid:   "100000000010001",
		Chcd:      "WXP",
		// TimeExpire: "201510201050000",
	}
	// 查询
	scanPayEnquiry = &model.ScanPayRequest{
		Busicd:       "INQY",
		Mchntid:      "200000000010001",
		AgentCode:    "19992900",
		OrigOrderNum: "1447168085242",
	}
	// 退款
	scanPayRefund = &model.ScanPayRequest{
		Busicd:       "REFD",
		Mchntid:      "200000000010001",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1447297657383",
		AgentCode:    "19992900",
		Txamt:        "000000000100",
		Currency:     "JPY",
		// Chcd:         "AOS",
	}
	// 撤销
	scanPayCancel = &model.ScanPayRequest{
		Busicd:       "VOID",
		Mchntid:      "200000000010001",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1447297595319",
		AgentCode:    "96683320",
	}
	// 关单
	scanPayClose = &model.ScanPayRequest{
		Busicd:       "CANC",
		Mchntid:      "200000000010001",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1447166282329",
		AgentCode:    "19992900",
	}
	// 企业支付
	scanPayEnterprise = &model.ScanPayRequest{
		Busicd:    "QYZF",
		Mchntid:   "200000000010001",
		OrderNum:  util.Millisecond(),
		AgentCode: "10134001",
		Chcd:      "WXP",
		Txamt:     "000000000100",
		OpenId:    "omYJss7PyKb02j3Y5pnZLm2IL6F4", //omYJss7PyKb02j3Y5pnZLm2IL6F4
		CheckName: "FORCE_CHECK",
		UserName:  "陈芝锐",
		Desc:      "ipad2 mini 64G",
	}
	// 公众号支付
	scanPayPublic = &model.ScanPayRequest{
		GoodsInfo:    "鞋子,1000,2;衣服,1500,3",
		Busicd:       "JSZF",
		Txamt:        "000000000001",
		Mchntid:      "100000000000203",
		OrderNum:     util.Millisecond(),
		AgentCode:    "CIL00002",
		Chcd:         "WXP",
		Code:         "001fbfbe9b2a351311e4212dd30c6f83",
		NeedUserInfo: "YES",
	}

	settQuery = &model.ScanPayRequest{
		Txndir:   "Q",
		Busicd:   "LIST",
		Mchntid:  "100000000000203",
		SettDate: "2015-12-01",
	}

	// 卡券核销
	purchaseCoupons = &model.ScanPayRequest{
		Txndir:    "Q",
		Busicd:    "VERI",
		AgentCode: "10134001",
		// Chcd:       "ULIVE",
		Mchntid:    "100000000010001",
		Terminalid: "30150006",
		OrderNum:   fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		ScanCodeId: "1816086060100100",
		// VeriTime:   "-1",
		// Txamt: "000000021000",
	}
	// 刷卡电子券核销
	purchaseActCoupons = &model.ScanPayRequest{
		Txndir:    "Q",
		Busicd:    "CRVE",
		AgentCode: "10134001",
		// Chcd:       "ULIVE",
		Mchntid:    "100000000010001",
		Terminalid: "30150006",
		OrderNum:   fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		ScanCodeId: "1816086060100100",
		// VeriTime:   "1",
		OrigOrderNum: "14498231841427131847",
		Cardbin:      "622525",
		Txamt:        "000000022000",
		PayType:      "2",
	}
	// 电子券查询
	queryPurchaseCouponsResult = &model.ScanPayRequest{
		Txndir:    "Q",
		Busicd:    "QUVE",
		AgentCode: "10134001",
		// Chcd:       "ULIVE",
		Mchntid:    "100000000010001",
		Terminalid: "30150006",
		OrderNum:   fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		ScanCodeId: "1816086060100100",
		// VeriTime:     "1",
		OrigOrderNum: "14500571571427131847",
	}
	// 刷卡电子券撤销
	undoPurchaseActCoupons = &model.ScanPayRequest{
		Txndir:    "Q",
		Busicd:    "CAVE",
		AgentCode: "10134001",
		// Chcd:       "ULIVE",
		Mchntid:    "100000000010001",
		Terminalid: "30150006",
		OrderNum:   fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		ScanCodeId: "1816086060100100",
		// VeriTime:   "-1",
		OrigOrderNum: "14500571571427131847",
	}
)

func doOneScanPay(scanPay *model.ScanPayRequest) error {
	mer, err := mongo.MerchantColl.Find(scanPay.Mchntid)
	if err != nil {
		return err
	}
	scanPay.Sign = security.SHA1WithKey(scanPay.SignMsg(), mer.SignKey)
	reqBytes, _ := json.Marshal(scanPay)
	respBytes := ScanPayHandle(reqBytes, false)
	resp := new(model.ScanPayResponse)
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return err
	}
	log.Debug(string(respBytes))
	return nil
}

func TestConcurrentScanPay(t *testing.T) {
	log.SetOutputLevel(log.Ldebug)
	var wg sync.WaitGroup
	n := "14417647179553"
	scanPayBarcodePay.OrderNum = n
	scanPayClose.OrigOrderNum = n
	wg.Add(2)
	go func() {
		doOneScanPay(scanPayBarcodePay)
		wg.Done()
	}()

	// scanPayClose.OrderNum += "2"
	go func() {
		time.Sleep(500 * time.Millisecond)
		doOneScanPay(scanPayClose)
		wg.Done()
	}()
	wg.Wait()
}

func TestScanPay(t *testing.T) {
	// scanPayEnterprise.OrderNum = "1444639800979"
	// scanPayClose.OrigOrderNum = "14417647179551"
	err := doOneScanPay(settQuery)
	if err != nil {
		t.Error(err)
	}
}

func TestSignMsg(t *testing.T) {

	str := `{"sign":"ed1838760bbde16ca708a49a4b5f5d3279374519","txndir":"Q","scanCodeId":"281223029725731233","mchntid":"991663048160001","orderNum":"2015092217294332704","busicd":"PURC","inscd":"99911888","txamt":"000000000001","terminalid":"00000379"}`

	req := new(model.ScanPayRequest)
	err := json.Unmarshal([]byte(str), req)
	if err != nil {
		t.Error(err)
	}
	t.Log(req.SignMsg())
}

// 测试卡券核销
func TestPurchaseCoupons(t *testing.T) {
	err := doOneScanPay(purchaseCoupons)
	if err != nil {
		t.Error(err)
	}
}
func TestPurchaseActCoupons(t *testing.T) {
	err := doOneScanPay(purchaseActCoupons)
	if err != nil {
		t.Error(err)
	}
}
func TestQueryPurchaseCouponsResult(t *testing.T) {
	err := doOneScanPay(queryPurchaseCouponsResult)
	if err != nil {
		t.Error(err)
	}
}
func TestUndoPurchaseActCoupons(t *testing.T) {
	err := doOneScanPay(undoPurchaseActCoupons)
	if err != nil {
		t.Error(err)
	}
}
