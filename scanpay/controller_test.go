package scanpay

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/security"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"sync"
	"testing"
	"time"
)

var (
	// 下单支付
	scanPayBarcodePay = &model.ScanPayRequest{
		GoodsInfo: "鞋子,1000.00,2;衣服,1500,3",
		OrderNum:  util.Millisecond(),
		// OrderNum:   "哈哈中文订单号",
		ScanCodeId: "130039663149119746",
		AgentCode:  "99911888",
		Txamt:      "000000000001",
		Chcd:       "WXP",
		Busicd:     "PURC",
		Mchntid:    "999118880000312",
		// Sign:       "ce76927257b57f133f68463c83bbd408e0f25211",
	}
	// 预下单支付
	scanPayQrCodeOfflinePay = &model.ScanPayRequest{
		GoodsInfo: "鞋子,1000,2;衣服,1500,3",
		OrderNum:  util.Millisecond(),
		AgentCode: "CIL00002",
		Txamt:     "100000000000210",
		Busicd:    "PAUT",
		Mchntid:   "100000000000021",
		Chcd:      "WXP",
	}
	// 查询
	scanPayEnquiry = &model.ScanPayRequest{
		Busicd:       "INQY",
		Mchntid:      "100000000000210",
		AgentCode:    "CIL00002",
		OrigOrderNum: "1439884584561",
	}
	// 退款
	scanPayRefund = &model.ScanPayRequest{
		Busicd:       "REFD",
		Mchntid:      "100000000000210",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1440032751947",
		AgentCode:    "CIL00002",
		Txamt:        "000000000001",
		Chcd:         "WXP",
	}
	// 撤销
	scanPayCancel = &model.ScanPayRequest{
		Busicd:       "VOID",
		Mchntid:      "966833200000007",
		OrderNum:     util.Millisecond(),
		OrigOrderNum: "1440040340745",
		AgentCode:    "96683320",
	}
	// 关单
	scanPayClose = &model.ScanPayRequest{
		Busicd:       "CANC",
		Mchntid:      "999118880000312",
		OrderNum:     util.Millisecond() + "1",
		OrigOrderNum: "1439886859870",
		AgentCode:    "99911888",
	}
	// 企业支付
	scanPayEnterprise = &model.ScanPayRequest{
		Busicd:    "QYFK",
		Mchntid:   "888888888888888",
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
	n := util.Millisecond()
	scanPayBarcodePay.OrderNum = n
	scanPayClose.OrigOrderNum = n
	wg.Add(2)
	go func() {
		doOneScanPay(scanPayBarcodePay)
		wg.Done()
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		doOneScanPay(scanPayClose)
		wg.Done()
	}()
	wg.Wait()
}

func TestScanPay(t *testing.T) {
	err := doOneScanPay(scanPayBarcodePay)
	if err != nil {
		t.Error(err)
	}
}

func TestSignMsg(t *testing.T) {
	t.Log(scanPayBarcodePay.SignMsg())
}
