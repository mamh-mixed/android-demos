package app

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"regexp"

	"github.com/CardInfoLink/quickpay/email"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/omigo/log"
)

type user struct{}

var User user
var timeReplacer = strings.NewReplacer("-", "", ":", "", " ", "")
var dateRegexp = regexp.MustCompile(`^\d{8}$`)
var monthRegexp = regexp.MustCompile(`^\d{6}$`)

// register 注册
func (u *user) register(req *reqParams) (result *model.AppResult) {
	log.Debugf("userName=%s,password=%s,transtime=%s", req.UserName, req.Password, req.Transtime)
	// 参数不能为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	user := &model.AppUser{
		UserName: req.UserName,
		Password: req.Password,
		Activate: "false",
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
	hostAddress := goconf.Config.App.NotifyURL
	activateUrl := fmt.Sprintf("%s/app/activate?username=%s&code=%s", hostAddress, req.UserName, code)

	email := &email.Email{
		To:    req.UserName,
		Title: activation.Title,
		Body:  fmt.Sprintf(activation.Body, activateUrl, "点我激活"),
	}
	err = email.Send()

	e := &model.Email{
		UserName:  req.UserName,
		Code:      code,
		Success:   true,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	if err != nil {
		e.Success = false
	}

	// 保存email信息
	err = mongo.EmailCol.Upsert(e)
	if err != nil {
		log.Errorf("save email err")
		return model.SYSTEM_ERROR
	}

	if e.Success {
		return model.SUCCESS1
	} else {
		return model.SYSTEM_ERROR
	}
}

// activate 激活
func (u *user) activate(req *reqParams) (result *model.AppResult) {
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
	// 如果用户已激活，则返回成功
	if user.Activate == "true" {
		return model.SUCCESS1
	}

	// 更新activate为已激活
	user.Activate = "true"
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR_CH
	}

	return model.SUCCESS1
}

// improveInfo 信息完善
func (u *user) improveInfo(req *reqParams) (result *model.AppResult) {
	if req.UserName == "" || req.Password == "" || req.BankOpen == "" || req.Payee == "" || req.PayeeCard == "" ||
		req.PhoneNum == "" || req.Transtime == "" {
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
	uniqueId := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int31())
	randStr := fmt.Sprintf("%d", rand.Int31())
	permission := []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf}
	merchant := &model.Merchant{
		AgentCode:  "99911888",
		AgentName:  "讯联O2O机构",
		Permission: permission,
		MerStatus:  model.MerStatusNormal,
		TransCurr:  "156",
		UniqueId:   uniqueId,
		IsNeedSign: true,
		SignKey:    fmt.Sprintf("%x", base64.StdEncoding.EncodeToString([]byte(randStr))),
		Detail: model.MerDetail{
			MerName:       "云收银",
			CommodityName: "讯联云收银在线注册商户",
		},
	}
	for {
		// 设置merId
		maxMerId, err := mongo.MerchantColl.FindMaxMerId()
		if err != nil {
			if err.Error() == "not found" {
				log.Infof(" set max merId is 999118880000001")
				merchant.MerId = "999118880000001"
			} else {
				log.Errorf("find database  err,%s", err)
				return model.SYSTEM_ERROR
			}

		} else {
			maxMerIdNum, err := strconv.Atoi(maxMerId)
			if err != nil {
				log.Errorf("format maxMerId(%s) err", maxMerId)
				return model.SYSTEM_ERROR
			}
			merchant.MerId = fmt.Sprintf("%d", maxMerIdNum+1)
		}

		err = mongo.MerchantColl.Insert2(merchant)
		if err != nil {
			isDuplicateMerId := strings.Contains(err.Error(), "E11000 duplicate key error index")
			if !isDuplicateMerId {
				log.Errorf("create merchant err,%s", err)
				return model.SYSTEM_ERROR
			}

		} else {
			break
		}
	}

	// 创建路由,支付宝，微信
	alpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "ALP",
		ChanCode:  "ALP",
		ChanMerId: "2088811767473826",
	}
	err = mongo.RouterPolicyColl.Insert(alpRoute)
	if err != nil {
		log.Errorf("create routePolicy err,%s", err)
		return model.SYSTEM_ERROR
	}

	wxpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "WXP",
		ChanCode:  "WXP",
		ChanMerId: "1239305502",
	}
	err = mongo.RouterPolicyColl.Insert(wxpRoute)
	if err != nil {
		log.Errorf("create routePolicy err,%s", err)
		return model.SYSTEM_ERROR
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

	// 用户名不为空
	if req.UserName == "" {
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

	// 用户名不为空
	if req.UserName == "" {
		return model.PARAMS_EMPTY
	}

	if !monthRegexp.MatchString(req.Date) {
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
	date := req.Date
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
		q.RefundStatus = model.TransRefunded
		q.TransStatus = []string{model.TransSuccess}
	case "fail":
		q.TransStatus = []string{model.TransFail, model.TransHandling}
	}

	trans, total, err := mongo.SpTransColl.Find(q)
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
	if total != 0 {
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
	// 用户名不为空
	if req.UserName == "" || req.NewPassword == "" {
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
		return model.USERNAME_PASSWORD_ERROR
	}

	user.Password = req.NewPassword
	if err = mongo.AppUserCol.Upsert(user); err != nil {
		return model.SYSTEM_ERROR
	}

	return model.SUCCESS1
}

// promoteLimit 提升限额
func (u *user) promoteLimit(req *reqParams) (result *model.AppResult) {
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

	log.Debugf("%+v", user)
	// 返回
	result = model.NewAppResult(model.SUCCESS, "")
	settInfo := &model.SettInfo{
		Payee:     user.Payee,
		BankOpen:  user.BankOpen,
		PayeeCard: user.PayeeCard,
		PhoneNum:  user.PhoneNum,
	}

	result.SettInfo = settInfo
	return
}

// updateSettInfo 更新清算信息
func (u *user) updateSettInfo(req *reqParams) (result *model.AppResult) {
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

	// 修改
	user.Payee = req.Payee
	user.BankOpen = req.BankOpen
	user.PayeeCard = req.PayeeCard
	user.PhoneNum = req.PhoneNum

	if err = mongo.AppUserCol.Upsert(user); err != nil {
		return model.SYSTEM_ERROR
	}

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
