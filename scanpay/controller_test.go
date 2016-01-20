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
		GoodsInfo:      "food,10.00,2;water,1.00,3",
		OrderNum:       util.Millisecond(),
		ScanCodeId:     "283850099094963575",
		AgentCode:      "10134001",
		Txamt:          "000000000001",
		DiscountAmt:    "000000000001",
		PayType:        "5",
		Busicd:         "PURC",
		Currency:       "CNY",
		Mchntid:        "100000000010001",
		CouponOrderNum: "kaquandingdanhao",
		// Chcd:        "AOS",
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

	// 电子券核销
	purchaseCouponsSingle = &model.ScanPayRequest{
		Txndir:    "Q",
		Busicd:    "VERI",
		AgentCode: "10134001",
		// Chcd:       "ULIVE",
		Mchntid:    "999118880000017",
		Terminalid: "30150006",
		OrderNum:   fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		ScanCodeId: "1810037103015010",
		// VeriTime:   "-1",
		// Txamt:   "000000010100",
		// Cardbin: "665523",
		PayType: "4",
	}
	// 电子券验证冲正
	recoverCoupons = &model.ScanPayRequest{
		Txndir:    "Q",
		Busicd:    "CAVE",
		AgentCode: "10134001",
		// Chcd:       "ULIVE",
		Mchntid:      "999118880000017",
		Terminalid:   "30150006",
		OrderNum:     fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31()),
		OrigOrderNum: "14513068011474941318",
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
	t.Logf("Order Number is %s", scanPayBarcodePay.OrderNum)
	err := doOneScanPay(scanPayBarcodePay)
	if err != nil {
		t.Error(err)
	}
}

func TestSignMsg(t *testing.T) {

	//4d045cf4039a420a86824c7132a24d6ff4c559f3
	str := `{"txndir":"Q","busicd":"PURC","inscd":"99911888","chcd":"WXP","mchntid":"991221054110001","txamt":"000000001400","goodsInfo":"6927229221501,福宁鸡蛋肉松面包,3,2.50;6901028001465,双喜（软国际）,1,6.50;","orderNum":"1056000011024917","scanCodeId":"130472120612304529","currency":"CNY","terminalid":"10560001","sign":"99706d3f26df36a33ecd51d928a7181d208f7608"}`

	req := new(model.ScanPayRequest)
	err := json.Unmarshal([]byte(str), req)
	if err != nil {
		t.Error(err)
	}

	err = doOneScanPay(req)
	time.Sleep(5 * time.Second)
	if err != nil {
		t.Error(err)
	}
	// t.Log(security.SHA1WithKey(req.SignMsg(), "8627a2ba43da3ada31b820b788680b99"))
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

func TestPurchaseCouponsSingle(t *testing.T) {
	err := doOneScanPay(purchaseCouponsSingle)
	if err != nil {
		t.Error(err)
	}
}
func TestRecoverCoupons(t *testing.T) {
	err := doOneScanPay(recoverCoupons)
	if err != nil {
		t.Error(err)
	}
}
