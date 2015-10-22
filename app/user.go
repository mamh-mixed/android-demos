package app

import (
	cr "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/email"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

type user struct{}

var (
	User         user
	timeReplacer = strings.NewReplacer("-", "", ":", "", " ", "")
	dateRegexp   = regexp.MustCompile(`^\d{8}$`)
	monthRegexp  = regexp.MustCompile(`^\d{6}$`)
	b64Encoding  = base64.StdEncoding
	hostAddress  = goconf.Config.App.NotifyURL
	WXPMerId     = goconf.Config.MobileApp.WXPMerId
	ALPMerId     = goconf.Config.MobileApp.ALPMerId
)

// register 注册
func (u *user) register(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	log.Debugf("userName=%s,password=%s,transtime=%s", req.UserName, req.Password, req.Transtime)
	// 参数不能为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 用户是否存在
	num, err := mongo.AppUserCol.FindCountByUserName(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}
	if num != 0 {
		return model.USERNAME_EXIST
	}

	user := &model.AppUser{
		UserName:     req.UserName,
		Password:     req.Password,
		Activate:     "false",
		Limit:        "true",
		Remark:       req.Remark,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		SubAgentCode: req.SubAgentCode,
	}
	user.UpdateTime = user.CreateTime

	// 放进req里
	req.AppUser = user

	// 保存用户信息
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR
	}

	return model.SUCCESS1
}

// login 登录
func (u *user) login(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	log.Debugf("userName=%s,password=%s,transtime=%s", req.UserName, req.Password, req.Transtime)
	// 参数不能为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}
	// 密码是否正确
	if user.Password != req.Password {
		return model.USERNAME_PASSWORD_ERROR
	}
	// 用户是否激活
	if user.Activate == "false" {
		return model.USER_NO_ACTIVATE
	}

	// 查找uniqueId
	if user.MerId != "" {
		merchant, err := mongo.MerchantColl.Find(user.MerId)
		if err != nil {
			log.Errorf("find database err,%s", err)
			return model.SYSTEM_ERROR
		}
		user.SignKey = merchant.SignKey
		user.UniqueId = merchant.UniqueId
		user.AgentCode = merchant.AgentCode
	}

	result = &model.AppResult{
		State: model.SUCCESS,
		Error: "",
		User:  user,
	}

	return result
}

// reqActivate 请求发送激活链接
func (u *user) reqActivate(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	log.Debugf("userName=%s,password=%s,transtime=%s", req.UserName, req.Password, req.Transtime)
	// 参数不能为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}

	// 密码是否正确
	if user.Password != req.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 如果用户已激活，直接返回success
	if user.Activate == "true" {
		return model.SUCCESS1
	}

	// 发送激活链接到注册时提供的邮箱
	code := fmt.Sprintf("%d", rand.Int31())
	// hostAddress := goconf.Config.App.NotifyURL
	activateUrl := fmt.Sprintf("%s/app/activate?username=%s&code=%s", hostAddress, req.UserName, code)

	click := b64Encoding.EncodeToString(randBytes(32))

	email := &email.Email{
		To:    req.UserName,
		Title: activation.Title,
		Body:  fmt.Sprintf(activation.Body, activateUrl, click),
	}

	e := &model.Email{
		UserName:  req.UserName,
		Code:      code,
		Success:   false,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 保存email信息
	err = mongo.EmailCol.Upsert(e)
	if err != nil {
		log.Errorf("save email err: %s", err)
		return model.SYSTEM_ERROR
	}

	// 异步发送邮件
	go func() {
		err := email.Send()
		if err != nil {
			log.Errorf("send email fail: %s", err)
			return
		}

		// update
		e.Success = true
		mongo.EmailCol.Upsert(e)
	}()

	return model.SUCCESS1
}

// activate 激活
func (u *user) activate(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	log.Debugf("userName=%s,code=%s", req.UserName, req.Code)
	// 参数不能为空
	if req.UserName == "" || req.Code == "" {
		return model.PARAMS_EMPTY_CH
	}

	// 判断code是否正确
	e, err := mongo.EmailCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST_CH
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR_CH
	}
	if req.Code != e.Code {
		return model.CODE_ERROR_CH
	}

	// code有效期为2小时
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Timestamp, time.Local)
	if time.Now().Sub(timestamp) > 2*time.Hour {
		return model.CODE_TIME_ERROR_CH
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST_CH
		}
		return model.SYSTEM_ERROR_CH
	}
	// 如果用户已激活，表示码已过期
	if user.Activate == "true" {
		return model.CODE_TIME_ERROR_CH
	}

	// 更新activate为已激活
	user.Activate = "true"
	user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR_CH
	}

	return model.SUCCESS1
}

// improveInfo 信息完善
func (u *user) improveInfo(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		log.Errorf("find database err,%s", err)
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		return model.SYSTEM_ERROR
	}

	// 密码是否正确
	if user.Password != req.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 用户信息是否已更新
	if user.MerId != "" {
		return model.USER_ALREADY_IMPROVED
	}

	// 创建商户
	permission := []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf}
	merchant := &model.Merchant{
		AgentCode:  "99911888",
		AgentName:  "讯联O2O机构",
		Permission: permission,
		MerStatus:  model.MerStatusNormal,
		TransCurr:  "156",
		RefundType: model.CurrentDayRefund, // 只能当天退
		IsNeedSign: true,
		SignKey:    fmt.Sprintf("%x", randBytes(16)),
		Detail: model.MerDetail{
			MerName:       "云收银",
			CommodityName: "讯联云收银在线注册商户",
			Province:      req.Province,
			City:          req.City,
			OpenBankName:  req.BankOpen,
			BankName:      req.BranchBank,
			BankId:        req.BankNo,
			AcctName:      req.Payee,
			AcctNum:       req.PayeeCard,
			ContactTel:    req.PhoneNum,
		},
	}

	// 生成商户号，并保存商户
	if err := genMerId(merchant, "999118880"); err != nil {
		return err
	}

	// 创建路由,支付宝，微信
	if err := genRouter(merchant); err != nil {
		return err
	}

	// 更新用户信息
	user.BankOpen = req.BankOpen
	user.Payee = req.Payee
	user.PayeeCard = req.PayeeCard
	user.PhoneNum = req.PhoneNum
	user.MerId = merchant.MerId
	user.Limit = "true"
	user.SignKey = merchant.SignKey
	user.AgentCode = merchant.AgentCode
	user.UniqueId = merchant.UniqueId
	user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR
	}
	result = &model.AppResult{
		State: model.SUCCESS,
		Error: "",
		User:  user,
	}
	return result
}

// getTotalTransAmt 查询某天交易总额
func (u *user) getTotalTransAmt(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	// 用户名不为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" {
		return model.PARAMS_EMPTY
	}

	if !dateRegexp.MatchString(req.Date) {
		return model.TIME_ERROR
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	result = model.NewAppResult(model.SUCCESS, "")
	month := req.Date
	month = month[:4] + "-" + month[4:6] + "-" + month[6:8]

	ret := query.TransStatistics(&model.QueryCondition{
		MerId:     user.MerId,
		StartTime: month,
		EndTime:   month,
		Size:      1,
		Page:      1,
	})

	if summary, ok := ret.Rec.(model.Summary); ok {
		result.Count = summary.TotalTransNum
		result.TotalAmt = fmt.Sprintf("%0.2f", summary.TotalTransAmt)
	} else {
		return model.SYSTEM_ERROR
	}

	return result
}

// getUserBill 获取用户账单
func (u *user) getUserBill(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	// 用户名不为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" {
		return model.PARAMS_EMPTY
	}

	if !monthRegexp.MatchString(req.Month) {
		return model.TIME_ERROR
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	result = model.NewAppResult(model.SUCCESS, "")
	date := req.Month
	yearNum, _ := strconv.Atoi(date[:4])
	month := date[4:6]
	day := ""
	if month == "01" || month == "03" || month == "05" || month == "07" || month == "08" || month == "10" || month == "12" {
		day = "31"
	} else if month == "02" {
		if (yearNum%4 == 0 && yearNum%100 != 0) || yearNum%400 == 0 {
			day = "29"
		} else {
			day = "28"
		}
	} else {
		day = "30"
	}

	startDate := date[:4] + "-" + date[4:6] + "-" + "01"
	endDate := date[:4] + "-" + date[4:6] + "-" + day

	q := &model.QueryCondition{
		MerId:     user.MerId,
		StartTime: startDate + " 00:00:00",
		EndTime:   endDate + " 23:59:59",
		Size:      15,
		Page:      1,
		Skip:      req.Index,
	}

	switch req.Status {
	case "all":
	case "success":
		q.Respcd = "00"
	case "fail":
		q.RespcdNotIn = "00"
	}

	trans, _, err := mongo.SpTransColl.Find(q)
	if err != nil {
		log.Errorf("find user trans error: %s", err)
		return model.SYSTEM_ERROR
	}

	var txns []*model.AppTxn
	for _, t := range trans {
		txns = append(txns, transToTxn(t))
	}

	var transAmt, refundAmt int64
	var transCount, refundCount int
	typeGroup, err := mongo.SpTransColl.MerBills(&model.QueryCondition{
		MerId:        user.MerId,
		StartTime:    q.StartTime,
		EndTime:      q.EndTime,
		RefundStatus: model.TransRefunded,
		TransStatus:  []string{model.TransSuccess},
	})
	if err != nil {
		return model.SYSTEM_ERROR
	}

	for _, v := range typeGroup {
		switch v.TransType {
		case model.PayTrans:
			transAmt += v.TransAmt
			transCount += v.TransNum
		default:
			refundAmt += v.TransAmt
			transAmt -= v.TransAmt
			refundCount += v.TransNum
		}
	}

	result.Txn = txns
	result.Size = len(trans)
	result.TotalAmt = fmt.Sprintf("%0.2f", float32(transAmt)/100)
	result.RefdCount = refundCount
	result.RefdTotalAmt = fmt.Sprintf("%0.2f", float32(refundAmt)/100)
	result.Count = transCount
	return
}

// getUserTrans 获取用户某笔交易信息
func (u *user) getUserTrans(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	// 用户名不为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// TODO:密码不对，兼容IOS客户端bug，暂时不验证密码
	// if req.Password != user.Password {
	// 	return model.USERNAME_PASSWORD_ERROR
	// }

	// 没有机构用户
	if user.MerId == "" {
		return model.NO_PAY_MER
	}

	// 查找交易
	t, err := mongo.SpTransColl.FindOne(user.MerId, req.OrderNum)
	if err != nil {
		return model.NO_TRANS
	}

	result = model.NewAppResult(model.SUCCESS, "")
	switch req.BusinessType {
	case "getRefd":
		result.RefdTotalAmt = fmt.Sprintf("%0.2f", float32(t.RefundAmt)/100)
	case "getOrder":
		result.Txn = transToTxn(t)
	}

	return result
}

// passwordHandle 修改密码
func (u *user) passwordHandle(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	// 用户名不为空
	if req.UserName == "" || req.NewPassword == "" || req.Transtime == "" || req.OldPassword == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.OldPassword != user.Password {
		return model.OLD_PASSWORD_ERROR
	}

	user.Password = req.NewPassword
	user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	if err = mongo.AppUserCol.Upsert(user); err != nil {
		return model.SYSTEM_ERROR
	}

	return model.SUCCESS1
}

// promoteLimit 提升限额
func (u *user) promoteLimit(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	// 用户名不为空
	if req.UserName == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	if user.Limit == "false" {
		return model.SUCCESS1
	}

	// 发送邮件通知Andy.Li
	email := &email.Email{
		To:    andyLi,
		Title: promote.Title,
		Body:  fmt.Sprintf(promote.Body, req.Payee, req.Email, req.PhoneNum, user.MerId),
	}

	go func() {
		if err = email.Send(); err != nil {
			log.Errorf("send email error: %s", err)
		}
	}()

	return model.SUCCESS1
}

// getSettInfo 获得清算信息
func (u *user) getSettInfo(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	// 用户名不为空
	if req.UserName == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}
	mer, err := mongo.MerchantColl.Find(user.MerId)
	if err != nil {
		if err.Error() == "not found" {
			return model.MERID_NO_EXIST
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// log.Debugf("%+v", user)
	// 返回
	result = model.NewAppResult(model.SUCCESS, "")
	settInfo := &model.SettInfo{
		Payee:      mer.Detail.AcctName,
		BankOpen:   mer.Detail.OpenBankName,
		PayeeCard:  mer.Detail.AcctNum,
		PhoneNum:   mer.Detail.ContactTel,
		Province:   mer.Detail.Province,
		City:       mer.Detail.City,
		BranchBank: mer.Detail.BankName,
		BankNo:     mer.Detail.BankId,
	}

	result.SettInfo = settInfo
	return
}

// updateSettInfo 更新清算信息
func (u *user) updateSettInfo(req *reqParams) (result *model.AppResult) {

	// 字段长度验证
	if result = requestDataValidate(req); result != nil {
		return result
	}

	var user *model.AppUser
	var err error

	// 云收银app
	if req.AppUser == nil {
		// 用户名不为空
		if req.UserName == "" || req.Transtime == "" || req.Password == "" {
			return model.PARAMS_EMPTY
		}
		// 根据用户名查找用户
		user, err = mongo.AppUserCol.FindOne(req.UserName)
		if err != nil {
			if err.Error() == "not found" {
				return model.USERNAME_NO_EXIST
			}
			log.Errorf("find database err,%s", err)
			return model.SYSTEM_ERROR
		}

		// 密码不对
		if req.Password != user.Password {
			return model.USERNAME_PASSWORD_ERROR
		}
	} else {
		// 从销售工具接口过来的
		user = req.AppUser
	}

	m, err := mongo.MerchantColl.Find(user.MerId)
	if err != nil {
		return model.MERID_NO_EXIST
	}

	if req.MerName != "" {
		m.Detail.MerName = req.MerName
	}
	if req.Province != "" {
		m.Detail.Province = req.Province
	}
	if req.City != "" {
		m.Detail.City = req.City
	}
	if req.BankOpen != "" {
		m.Detail.OpenBankName = req.BankOpen
	}
	if req.BranchBank != "" {
		m.Detail.BankName = req.BranchBank
	}
	if req.BankNo != "" {
		m.Detail.BankId = req.BankNo
	}
	if req.Payee != "" {
		m.Detail.AcctName = req.Payee
	}
	if req.PayeeCard != "" {
		m.Detail.AcctNum = req.PayeeCard
	}
	if req.PhoneNum != "" {
		m.Detail.ContactTel = req.PhoneNum
	}
	if len(req.Images) > 0 {
		m.Detail.Images = req.Images
	}

	if err = mongo.MerchantColl.Update(m); err != nil {
		return model.SYSTEM_ERROR
	}

	req.m = m

	return model.SUCCESS1
}

func transToTxn(t *model.Trans) *model.AppTxn {
	txn := &model.AppTxn{
		Response:        t.RespCode,
		SystemDate:      timeReplacer.Replace(t.CreateTime),
		ConsumerAccount: t.ConsumerAccount,
	}
	txn.ReqData.Busicd = t.Busicd
	txn.ReqData.AgentCode = t.AgentCode
	txn.ReqData.Txndir = "Q"
	txn.ReqData.Terminalid = t.Terminalid
	txn.ReqData.OrigOrderNum = t.OrigOrderNum
	txn.ReqData.OrderNum = t.OrderNum
	txn.ReqData.MerId = t.MerId
	txn.ReqData.TradeFrom = t.TradeFrom
	txn.ReqData.Txamt = fmt.Sprintf("%012d", t.TransAmt)
	txn.ReqData.ChanCode = t.ChanCode
	txn.ReqData.Currency = t.TransCurr
	return txn
}

func randBytes(length int) []byte {
	var randBytes = make([]byte, length)
	if _, err := io.ReadFull(cr.Reader, randBytes[:]); err != nil {
		log.Errorf("io.ReadFull error: %s", err)
	}
	return randBytes
}

func genRouter(merchant *model.Merchant) *model.AppResult {

	// 创建路由,支付宝，微信
	alpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "ALP",
		ChanCode:  "ALP",
		ChanMerId: ALPMerId,
	}
	err := mongo.RouterPolicyColl.Insert(alpRoute)
	if err != nil {
		log.Errorf("create routePolicy err: %s", err)
		return model.SYSTEM_ERROR
	}

	wxpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "WXP",
		ChanCode:  "WXP",
		ChanMerId: WXPMerId,
	}
	err = mongo.RouterPolicyColl.Insert(wxpRoute)
	if err != nil {
		log.Errorf("create routePolicy err: %s", err)
		return model.SYSTEM_ERROR
	}

	return nil
}

func genMerId(merchant *model.Merchant, prefix string) *model.AppResult {
	for {
		// 设置merId
		maxMerId, err := mongo.MerchantColl.FindMaxMerId(prefix)
		if err != nil {
			if err.Error() == "not found" {
				// 从第一个开始编
				merchant.MerId = prefix + "000001"
			} else {
				log.Errorf("find merchant err,%s", err)
				return model.SYSTEM_ERROR
			}

		} else {
			log.Debugf("maxMerId: %s", maxMerId)
			var maxMerIdNum int
			if len(maxMerId) == 15 {
				// TODO:
				order := maxMerId[len(prefix):15]
				maxMerIdNum, err = strconv.Atoi(order)
				if err != nil {
					log.Errorf("format maxMerId(%s) err", maxMerId)
					return model.SYSTEM_ERROR
				}
				merchant.MerId = fmt.Sprintf("%s%0"+fmt.Sprintf("%d", len(order))+"d", prefix, maxMerIdNum+1)
			} else if len(maxMerId) < 15 {
				// TODO:
				l := fmt.Sprintf("%d", 14-len(prefix))
				rp := fmt.Sprintf("%0"+l+"d", 0)
				merchant.MerId = fmt.Sprintf("%s%s%d", prefix, rp, 1)
			} else {
				// TODO:
				log.Errorf("format maxMerId(%s) err", maxMerId)
				return model.SYSTEM_ERROR
			}
		}

		merchant.UniqueId = util.Confuse(merchant.MerId)
		err = mongo.MerchantColl.Insert2(merchant)
		if err != nil {
			isDuplicateMerId := strings.Contains(err.Error(), "E11000 duplicate key error index")
			if !isDuplicateMerId {
				log.Errorf("add merchant err: %s, merId=%s", err, merchant.MerId)
				return model.SYSTEM_ERROR
			}
		}

		break
	}
	return nil
}
