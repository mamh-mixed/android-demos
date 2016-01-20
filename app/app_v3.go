package app

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/qiniu"

	"github.com/CardInfoLink/log"
)

// checkSignSha256 APP V3 版本使用SHA256算法签名
func checkSignSha256(r *http.Request) bool {
	sign, content := r.FormValue("sign"), signContent(r.Form)
	log.Debugf("sign content: %s", content)

	valid := fmt.Sprintf("%x", sha256.Sum256([]byte(content+sha1Key)))
	if sign != valid {
		log.Warnf("check sign error, expect %s ,get %s", valid, sign)
		return false
	}

	return true
}

// billV3Handle 获取账单信息
func billV3Handle(w http.ResponseWriter, r *http.Request) {

	// 可跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	result := User.getUserBill(&reqParams{
		UserName:    r.FormValue("username"),
		Password:    r.FormValue("password"),
		Month:       r.FormValue("month"),
		Date:        r.FormValue("day"),
		Status:      r.FormValue("status"),
		Transtime:   r.FormValue("transtime"),
		Index:       r.FormValue("index"),
		OrderDetail: r.FormValue("order_detail"),
		Size:        r.FormValue("size"),
		TransType:   model.PayTrans,
	})

	w.Write(jsonMarshal(result))
}

// qiniuTokenHandler 获取七牛的上传token
func qiniuTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("username is %s; password is %s", r.FormValue("username"), r.FormValue("password"))
	result := User.getQiniuToken(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
	})

	if result.State == "success" {
		result.UploadToken = qiniu.GetUploadtoken()
	}

	w.Write(jsonMarshal(result))
}

// registerHandler 注册处理
func registerHandler(w http.ResponseWriter, r *http.Request) {
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

// loginHandler 登录处理
func loginHandler(w http.ResponseWriter, r *http.Request) {
	req := &reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
	}

	user := new(model.AppUser)
	devictType := r.FormValue("deviceType")
	user.DeviceType = strings.ToUpper(devictType)
	user.DeviceToken = r.FormValue("deviceToken")
	req.AppUser = user

	result := User.login(req)

	w.Write(jsonMarshal(result))
}

// forgetPasswordHandler 忘记密码处理
func forgetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	result := User.forgetPassword(&reqParams{
		UserName:  r.FormValue("username"),
		Transtime: r.FormValue("transtime"),
	})

	w.Write(jsonMarshal(result))
}

// updatePasswordHandler 更新密码处理
func updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	result := User.passwordHandle(&reqParams{
		UserName:    r.FormValue("username"),
		OldPassword: r.FormValue("oldpassword"),
		NewPassword: r.FormValue("newpassword"),
		Transtime:   r.FormValue("transtime"),
	})

	w.Write(jsonMarshal(result))
}

// activateAccountHandler 激活帐户处理。通过用户名和密码激活的
func activateAccountHandler(w http.ResponseWriter, r *http.Request) {
	result := User.reqActivate(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
	})

	w.Write(jsonMarshal(result))
}

// improveAccountHandler 清算银行卡信息完善处理。
func improveAccountHandler(w http.ResponseWriter, r *http.Request) {
	result := User.improveInfo(&reqParams{
		UserName:   r.FormValue("username"),
		Password:   r.FormValue("password"),
		BankOpen:   r.FormValue("bankOpen"),
		Payee:      r.FormValue("payee"),
		PayeeCard:  r.FormValue("payeeCard"),
		PhoneNum:   r.FormValue("phoneNum"),
		Transtime:  r.FormValue("transtime"),
		Province:   r.FormValue("province"),
		City:       r.FormValue("city"),
		BranchBank: r.FormValue("branchBank"),
		BankNo:     r.FormValue("bankNo"),
	})

	w.Write(jsonMarshal(result))
}

// settleInfoHandler 获取银行卡清算信息的处理
func settleInfoHandler(w http.ResponseWriter, r *http.Request) {
	result := User.getSettInfo(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
	})

	w.Write(jsonMarshal(result))
}

// certificateAccountHandler 帐户验证处理，用于提升限额
func certificateAccountHandler(w http.ResponseWriter, r *http.Request) {
	result := User.improveCertInfo(&reqParams{
		UserName:         r.FormValue("username"),
		Password:         r.FormValue("password"),
		Transtime:        r.FormValue("transtime"),
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

// billsHandler 账单处理
func billsHandler(w http.ResponseWriter, r *http.Request) {
	result := UserV3.getUserBills(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
		Month:     r.FormValue("month"),
		Status:    r.FormValue("status"),
		Index:     r.FormValue("index"),
		Size:      r.FormValue("size"),
		TransType: model.PayTrans,
	})

	w.Write(jsonMarshal(result))
}

// totalSummaryHandler 单日汇总处理
func totalSummaryHandler(w http.ResponseWriter, r *http.Request) {
	result := UserV3.getDaySummary(&reqParams{
		UserName:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		Transtime:    r.FormValue("transtime"),
		BusinessType: r.FormValue("reportType"), // 报表类型。1:收款账单；2:卡券账单
		Date:         r.FormValue("day"),
	})

	w.Write(jsonMarshal(result))
}

// ordersHandler 查询订单
func ordersHandler(w http.ResponseWriter, r *http.Request) {
	result := UserV3.findOrderHandle(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
		OrderNum:  r.FormValue("orderNum"),
		PayType:   r.FormValue("payType"),
		RecType:   r.FormValue("recType"),
		Status:    r.FormValue("txnStatus"),
		Index:     r.FormValue("index"),
		Size:      r.FormValue("size"),
	})

	w.Write(jsonMarshal(result))
}

// couponsHandler  卡券列表
func couponsHandler(w http.ResponseWriter, r *http.Request) {
	result := User.couponsHandler(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
		Month:     r.FormValue("month"),
		Index:     r.FormValue("index"),
		Size:      r.FormValue("size"),
	})
	w.Write(jsonMarshal(result))
}

// messagePullHandler 消息接口
func messagePullHandler(w http.ResponseWriter, r *http.Request) {
	result := UserV3.messagePullHandler(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
		Size:      r.FormValue("size"),
		LastTime:  r.FormValue("lastTime"),
		MaxTime:   r.FormValue("maxTime"),
	})

	w.Write(jsonMarshal(result))
}

// messageUpdateHandler 消息修改
func messageUpdateHandler(w http.ResponseWriter, r *http.Request) {
	result := User.updateMessageHandle(&reqParams{
		UserName:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Transtime: r.FormValue("transtime"),
		Message:   r.FormValue("message"),
	})

	w.Write(jsonMarshal(result))
}
