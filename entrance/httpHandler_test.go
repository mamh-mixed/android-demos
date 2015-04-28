package entrance

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/CardInfoLink/quickpay/model"
// 	"github.com/CardInfoLink/quickpay/mongo"
// 	"github.com/CardInfoLink/quickpay/tools"
// 	"github.com/omigo/log"
// )

// func post(req *http.Request, t *testing.T) {
// 	w := httptest.NewRecorder()
// 	Quickpay(w, req)
// 	log.Infof("%d - %s", w.Code, w.Body.String())
// 	if w.Code != 200 {
// 		t.Errorf("response error with status %d", w.Code)
// 	}

// 	var out model.BindingReturn
// 	err := json.Unmarshal(w.Body.Bytes(), &out)
// 	if err != nil {
// 		t.Errorf("Unmarshal response error (%s)", err)
// 	}

// 	if out.RespCode == "" {
// 		t.Error("测试失败")
// 	}
// }

// func TestBindingCreateWithSignHandle(t *testing.T) {
// 	b := model.BindingCreate{
// 		MerId:         "99001405",            // 商户ID
// 		BindingId:     tools.Millisecond(),   // 银行卡绑定ID
// 		AcctName:      "张三",                  // 账户名称
// 		AcctNum:       "6222020302062061908", // 账户号码
// 		IdentType:     "0",                   // 证件类型
// 		IdentNum:      "350583199009153732",  // 证件号码
// 		PhoneNum:      "18205960039",         // 手机号
// 		AcctType:      "20",                  // 账户类型
// 		ValidDate:     "1903",                // 信用卡有限期
// 		Cvv2:          "232",                 // CVV2
// 		SendSmsId:     "",                    // 发送短信验证码的交易流水
// 		SmsCode:       "",                    // 短信验证码
// 		BankId:        "102",                 // 银行ID
// 		ChanBindingId: "",                    // 渠道绑定ID
// 		ChanMerId:     "",                    // 渠道商户ID
// 	}
// 	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=" + b.MerId

// 	body, _ := json.Marshal(b)

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		t.Error("创建POST请求失败")
// 	}
// 	sign := SignatureUseSha1(body, "0123456789")
// 	req.Header.Set("X-Sign", sign)

// 	post(req, t)
// }

// func TestBindingCreateHandle(t *testing.T) {
// 	// todo 生成一个随机的商户号
// 	rdMerId := tools.Millisecond()
// 	// todo 在路由策略里面插入新商户号码的路由策略
// 	rp := &model.RouterPolicy{
// 		MerId:     rdMerId,
// 		CardBrand: "CUP",
// 		ChanCode:  "CFCA",
// 		ChanMerId: "001405",
// 	}

// 	if err := mongo.RouterPolicyColl.Insert(rp); err != nil {
// 		t.Errorf("Excepted no erro,but get ,error message is %s", err.Error())
// 	}
// 	// 生成一个随机的绑定ID
// 	// rdBindingId := tools.Millisecond()
// 	// merId := rdMerId
// 	url := "https://api.xxxx.com/quickpay/bindingCreate?merId=1000000000002"
// 	// body := `{"bindingId":"` + rdBindingId + `","acctName":"张三","acctNum":"6222022003008481261","identType":"0","identNum":"440583199111031012","phoneNum":"18205960039","acctType":"10","BankId":"102","sendSmsId":"1000000000009","smsCode":"12353"}`
// 	body := `{"bindingId":"vXfD08q1e9e5jvHmv1iDmXXP","acctName":"60202215176842555995459843018306154894ce1da849aa0af4699d5334fa5b","acctNum":"16996319105944721792491686643034c4eb35a9a1d2b3185caa2749f7081c923434dc9287ea9e036d657586fe6a2736","identType":"0","identNum":"311380834235092862316365984917959acde9dec78028d9791b291f0f1cefdaf4d6688567b377f7f4c0cbdee8d78eb4","phoneNum":"12056249536046666949568428698813a60c39514a4f4c69c576a517809f0002","acctType":"10"}`
// 	doPost("POST", url, body, t)
// }

// func doPost(method, url, body string, t *testing.T) {
// 	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
// 	if err != nil {
// 		t.Error("创建POST请求失败")
// 	}
// 	req.Header.Set("X-Sign", SignatureUseSha1([]byte(body), "0123456789")) // TODO

// 	w := httptest.NewRecorder()
// 	BindingPay(w, req)
// 	log.Infof("%d - %s", w.Code, w.Body.String())
// 	if w.Code != 200 {
// 		t.Errorf("response error with status %d", w.Code)
// 	}

// 	var out model.BindingReturn
// 	err = json.Unmarshal(w.Body.Bytes(), &out)
// 	if err != nil {
// 		t.Errorf("Unmarshal response error (%s)", err)
// 	}

// 	if out.RespCode == "" {
// 		t.Error("测试失败")
// 	}
// }

// func TestBindingRemoveHandle(t *testing.T) {
// 	url := "https://api.xxxx.com/quickpay/bindingRemove?merId=" + removeMerId
// 	body := `{"bindingId": ` + removeBindingId + `}`
// 	doPost("POST", url, body, t)
// }

// func TestBindingEnquiryHandle(t *testing.T) {
// 	url := "https://api.xxxx.com/quickpay/bindingEnquiry?merId=" + merId
// 	body := `{"bindingId": ` + bindingId + `}`
// 	doPost("POST", url, body, t)
// }

// func TestBindingPaymentHandle(t *testing.T) {
// 	url := "https://api.xxxx.com/quickpay/bindingPayment?merId=" + merId
// 	body := `{
// 		"subMerId": "",
// 		"merOrderNum": "` + merOrderNum + `",
// 		"transAmt": 1000,
// 		"bindingId": "` + bindingId + `",
// 		"sendSmsId": "` + sendSmsId + `",
// 		"smsCode": "` + smsCode + `",
// 		"merId":"` + merId + `"
// 	}`
// 	doPost("POST", url, body, t)
// }

// func TestRefundHandle(t *testing.T) {
// 	url := "https://api.xxxx.com/quickpay/refund?merId=" + merId
// 	body := `{
// 		"merOrderNum": "` + merOrderNum + `",
// 		"transAmt": 1000,
// 		"origOrderNum": "` + origOrderNum + `"
// 	}`
// 	doPost("POST", url, body, t)
// }

// func TestNoTrackPaymentHandle(t *testing.T) {
// 	url := "https://api.xxxx.com/quickpay/noTrackPayment?merId=" + merId
// 	body := `{
// 		"subMerId": "",
// 		"merOrderNum": "` + merOrderNum + `",
// 		"transAmt": 1000,
// 		"acctName":"` + acctName + `",
// 		"acctNum":"` + acctNum + `",
// 		"identType":"` + identType + `",
// 		"identNum":"` + identNum + `",
// 		"phoneNum":"` + phoneNum + `",
// 		"acctType":"` + acctType + `",
// 		"validDate":"` + validDate + `",
// 		"cvv2":"` + cvv2 + `",
// 		"sendSmsId": "` + sendSmsId + `",
// 		"smsCode": "` + smsCode + `"
// 	}`
// 	doPost("POST", url, body, t)
// }

// func TestOrderEnquiry(t *testing.T) {
// 	url := "https://api.xxxx.com/quickpay/orderEnquiry?merId=" + merId
// 	body := `{
// 		"origOrderNum":"` + origOrderNum + `",
// 		"merId":"` + merId + `",
// 		"showOrigInfo":"` + showOrigInfo + `"
// 		}`
// 	//"showOrigInfo":1,
// 	doPost("POST", url, body, t)

// 	//1426840770177
// 	//1426840770235
// 	// orderNum 1426841285631
// }
