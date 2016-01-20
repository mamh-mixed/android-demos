package app

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/qiniu"

	"github.com/CardInfoLink/log"
)

var sha1Key = "eu1dr0c8znpa43blzy1wirzmk8jqdaon"

func registerHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.register(&reqParams{
		UserName:       r.FormValue("username"),
		Password:       r.FormValue("password"),
		Transtime:      r.FormValue("transtime"),
		InvitationCode: r.FormValue("invitationCode"),
		UserFrom:       model.SelfRegister,
		Remark:         "self_register",
		Limit:          "true",
	})

	w.Write(jsonMarshal(result))
}

// loginHandle 登录
func loginHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	req := new(reqParams)
	req.UserName = r.FormValue("username")
	req.Password = r.FormValue("password")
	req.Transtime = r.FormValue("transtime")
	user := new(model.AppUser)
	user.DeviceType = r.FormValue("device_type")
	user.DeviceToken = r.FormValue("device_token")
	req.AppUser = user

	result := User.login(req)

	w.Write(jsonMarshal(result))
}

// reqActivateHandle 请求发送激活邮件
func reqActivateHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.reqActivate(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
	})

	w.Write(jsonMarshal(result))
}

// activateHandle 激活
func activateHandle(w http.ResponseWriter, r *http.Request) {

	result := User.activate(&reqParams{
		UserName: r.FormValue("username"),
		Code:     r.FormValue("code"),
	})

	successPage := "<html><head><title>激活跳转页面</title></head><body>激活成功</body></html>"
	failPage := "<html><head><title>激活跳转页面</title></head><body>激活失败，失败原因:%s</body></html>"

	if result.State == "success" {
		w.Write([]byte(successPage))
	} else {
		w.Write([]byte(fmt.Sprintf(failPage, result.Error)))
	}

}

// improveInfoHandle 补充清算信息
func improveInfoHandle(w http.ResponseWriter, r *http.Request) {

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.improveInfo(&reqParams{
		UserName:   r.FormValue("username"),
		Password:   r.FormValue("password"),
		BankOpen:   r.FormValue("bank_open"),
		Payee:      r.FormValue("payee"),
		PayeeCard:  r.FormValue("payee_card"),
		PhoneNum:   r.FormValue("phone_num"),
		Transtime:  r.FormValue("transtime"),
		Province:   r.FormValue("province"),
		City:       r.FormValue("city"),
		BranchBank: r.FormValue("branch_bank"),
		BankNo:     r.FormValue("bankNo"),
	})

	w.Write(jsonMarshal(result))

}

// getOrderHandle 获得单个订单信息
func getOrderHandle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.getUserTrans(&reqParams{
		UserName:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		OrderNum:     r.FormValue("orderNum"),
		Transtime:    r.FormValue("transtime"),
		BusinessType: "getOrder",
	})

	w.Write(jsonMarshal(result))
}

// billHandle 获取账单信息
func billHandle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	// TODO DELETE:修复客户端bug
	month := r.FormValue("month")
	if month == "201612" {
		month = "201512"
	}

	result := User.getUserBill(&reqParams{
		UserName:    r.FormValue("username"),
		Password:    r.FormValue("password"),
		Month:       month,
		Date:        r.FormValue("day"),
		Status:      r.FormValue("status"),
		Transtime:   r.FormValue("transtime"),
		Index:       r.FormValue("index"),
		OrderDetail: r.FormValue("order_detail"),
		Size:        r.FormValue("size"),
	})

	w.Write(jsonMarshal(result))

}

// getTotalHandle 获取某天总交易金额
func getTotalHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.getTotalTransAmt(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
		Date:      r.FormValue("date"),
	})

	w.Write(jsonMarshal(result))
}

// getRefdHandle 获得某笔交易已退款金额
func getRefdHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.getUserTrans(&reqParams{
		UserName:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		OrderNum:     r.FormValue("orderNum"),
		Transtime:    r.FormValue("transtime"),
		BusinessType: "getRefd",
	})

	w.Write(jsonMarshal(result))
}

// passwordHandle 密码修改
func passwordHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.passwordHandle(&reqParams{
		UserName:    r.FormValue("username"),
		OldPassword: r.FormValue("oldpassword"),
		NewPassword: r.FormValue("newpassword"),
		Transtime:   r.FormValue("transtime"),
	})

	w.Write(jsonMarshal(result))
}

// promoteLimitHandle 提升限额
func promoteLimitHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.promoteLimit(&reqParams{
		UserName: r.FormValue("username"),
		Password: r.FormValue("password"),
		Payee:    r.FormValue("payee"),
		PhoneNum: r.FormValue("phone_num"),
		Email:    r.FormValue("email"),
	})

	w.Write(jsonMarshal(result))
}

// updateSettInfoHandle 修改清算帐号信息
func updateSettInfoHandle(w http.ResponseWriter, r *http.Request) {

	// 暂不支持
	w.Write(jsonMarshal(model.NOT_SUPPORT))
	return

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.updateSettInfo(&reqParams{
		UserName:   r.FormValue("username"),
		Password:   r.FormValue("password"),
		BankOpen:   r.FormValue("bank_open"),
		Payee:      r.FormValue("payee"),
		PayeeCard:  r.FormValue("payee_card"),
		PhoneNum:   r.FormValue("phone_num"),
		Transtime:  r.FormValue("transtime"),
		Province:   r.FormValue("province"),
		City:       r.FormValue("city"),
		BranchBank: r.FormValue("branch_bank"),
		BankNo:     r.FormValue("bankNo"),
	})

	w.Write(jsonMarshal(result))
}

// getSettInfoHandle 获取清算帐号信息
func getSettInfoHandle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.getSettInfo(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
	})

	w.Write(jsonMarshal(result))
}

// ticketHandle 处理小票接口
func ticketHandle(w http.ResponseWriter, r *http.Request) {

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.ticketHandle(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		OrderNum:  r.FormValue("ordernum"),
		TicketNum: r.FormValue("receiptnum"),
	})

	w.Write(jsonMarshal(result))
}

// findOrderHandle 订单模糊搜索
func findOrderHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.findOrderHandle(&reqParams{
		UserName: r.FormValue("username"),
		Password: r.FormValue("password"),
		OrderNum: r.FormValue("orderNum"),
		PayType:  r.FormValue("payType"),
		RecType:  r.FormValue("recType"),
		Status:   r.FormValue("txnStatus"),
		Index:    r.FormValue("index"),
		Size:     r.FormValue("size"),
	})

	w.Write(jsonMarshal(result))
}

// updateMessageHandle 更新消息状态
func updateMessageHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.updateMessageHandle(&reqParams{
		UserName: r.FormValue("username"),
		Password: r.FormValue("password"),
		Message:  r.FormValue("message"),
	})

	w.Write(jsonMarshal(result))
}

// pullInfoHandle 推送消息接口
func pullInfoHandle(w http.ResponseWriter, r *http.Request) {

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.findPushMessage(&reqParams{
		UserName: r.FormValue("username"),
		Password: r.FormValue("password"),
		Size:     r.FormValue("size"),
		LastTime: r.FormValue("lasttime"),
		MaxTime:  r.FormValue("maxtime"),
	})

	w.Write(jsonMarshal(result))
}

// 重置密码
func forgetPasswordHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.forgetPassword(&reqParams{
		UserName: r.FormValue("username"),
	})

	w.Write(jsonMarshal(result))
}

//获取七牛token
func getQiniuTokenHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.getQiniuToken(&reqParams{
		UserName: r.FormValue("username"),
		Password: r.FormValue("password"),
	})

	result.UploadToken = qiniu.GetUploadtoken()

	w.Write(jsonMarshal(result))
}

//修改证书信息
func improveCertInfoHandle(w http.ResponseWriter, r *http.Request) {
	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.improveCertInfo(&reqParams{
		UserName:         r.FormValue("username"),
		Password:         r.FormValue("password"),
		CertName:         r.FormValue("certName"),
		CertAddr:         r.FormValue("certAddr"),
		LegalCertPos:     r.FormValue("legalCertPos"),
		LegalCertOpp:     r.FormValue("legalCertOpp"),
		BusinessLicense:  r.FormValue("businessLicense"),
		TaxRegistCert:    r.FormValue("taxRegistCert"),
		OrganizeCodeCert: r.FormValue("organizeCodeCert"),
	})

	w.Write(jsonMarshal(result))
}

func checkSign(r *http.Request) bool {

	sign := r.FormValue("sign")
	content := signContent(r.Form)
	log.Debugf("sign content: %s", content)
	valid := fmt.Sprintf("%x", sha1.Sum([]byte(content+sha1Key)))
	if sign != valid {
		log.Warnf("check sign error, expect %s ,get %s", valid, sign)
		return false
	}
	return true
}

func signContent(values url.Values) string {
	var keys []string
	for k, _ := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		// sign不参与签名
		if k == "sign" {
			continue
		}
		// 支持多个同名参数
		values := values[k]
		for _, v := range values {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k + "=" + v)
		}
	}
	return buf.String()
}

func jsonMarshal(result model.AppResult) []byte {
	data, err := json.Marshal(result)
	if err != nil {
		log.Error("json marshal error: %s", err)
		return []byte(model.JSON_ERROR)
	}
	log.Debugf("response message: %s", string(data))
	return data
}

type reqParams struct {
	OrderDetail      string
	UserName         string
	InvitationCode   string
	Password         string
	Transtime        string
	Sign             string
	Code             string
	BankOpen         string
	Payee            string
	PayeeCard        string
	PhoneNum         string
	Email            string
	OldPassword      string
	NewPassword      string
	OrderNum         string
	BusinessType     string
	Status           string
	Index            string
	Date             string
	Size             string
	Month            string
	Province         string
	City             string
	BranchBank       string
	BankNo           string
	Remark           string
	SubAgentCode     string
	MerName          string
	Images           []string
	UserFrom         int
	BelongsTo        string
	Limit            string
	TicketNum        string
	CertName         string
	CertAddr         string
	LegalCertPos     string
	LegalCertOpp     string
	BusinessLicense  string
	TaxRegistCert    string
	OrganizeCodeCert string
	PayType          string
	RecType          string
	LastTime         string
	MaxTime          string
	Message          string
	TransType        int
	ClientId         string
	AppUser          *model.AppUser
	m                *model.Merchant
}
