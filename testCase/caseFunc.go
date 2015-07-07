package testCase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// bindingCreate 创建绑定关系
func bindingCreate() (b *model.BindingCreate, err error) {

	b = &model.BindingCreate{
		MerId:     testMerId,
		BindingId: bindingId,
		AcctName:  "测试账号",
		AcctNum:   acctNum,
		IdentType: "0",
		IdentNum:  "440583199111031012",
		PhoneNum:  "18205960039",
		AcctType:  "20",
		ValidDate: "0612",
		Cvv2:      "793",
		SendSmsId: "",
		SmsCode:   "",
		BankId:    "102",
	}

	var aes = util.NewAESCBCEncrypt(testEncryptKey)
	b.AcctName = aes.Encrypt(b.AcctName)
	b.AcctNum = aes.Encrypt(b.AcctNum)
	b.IdentNum = aes.Encrypt(b.IdentNum)
	b.PhoneNum = aes.Encrypt(b.PhoneNum)
	b.ValidDate = aes.Encrypt(b.ValidDate)
	b.Cvv2 = aes.Encrypt(b.Cvv2)
	err = aes.Err
	return
}

// BindingPayment 绑定支付
func BindingPayment() (b *model.BindingPayment) {
	b = &model.BindingPayment{
		MerOrderNum: orderNum,
		TransAmt:    amt,
		BindingId:   bindingId,
		MerId:       testMerId,
	}
	return
}

// BindingRemove 解除绑定
func BindingRemove() (b *model.BindingRemove) {
	b = &model.BindingRemove{
		BindingId:     bindingId,
		TxSNUnBinding: util.Millisecond(),
	}
	return
}

// BindingRefund 退款
func BindingRefund() (b *model.BindingRefund) {

	b = &model.BindingRefund{
		OrigOrderNum: orderNum,
		MerOrderNum:  util.Millisecond(),
		TransAmt:     amt,
	}
	return
}

// OrderEnquiry 订单查询
func OrderEnquiry() (b *model.OrderEnquiry) {

	b = &model.OrderEnquiry{
		OrigOrderNum: orderNum,
	}

	return
}

// BindingEnquiry 绑定查询
func BindingEnquiry() (b *model.BindingEnquiry) {

	b = &model.BindingEnquiry{BindingId: bindingId}
	return
}

func post(url string, m interface{}) (*model.BindingReturn, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Sign", SignatureUseSha1(j, testSign))

	w := httptest.NewRecorder()
	Quickpay(w, req)
	log.Infof("%d - %s", w.Code, w.Body.String())
	if w.Code != 200 {
		return nil, fmt.Errorf("response error with status %d", w.Code)
	}

	var out model.BindingReturn
	err = json.Unmarshal(w.Body.Bytes(), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
