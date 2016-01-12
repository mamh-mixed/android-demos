package app

import (
	cr "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/email"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

type user struct{}

const (
	LOCKTIME      = 30000 //锁定三小时
	LOGINDIFFTIME = 10000 //1小时以内
)

var (
	User         user
	timeReplacer = strings.NewReplacer("-", "", ":", "", " ", "")
	dateRegexp   = regexp.MustCompile(`^\d{8}$`)
	monthRegexp  = regexp.MustCompile(`^\d{6}$`)
	b64Encoding  = base64.StdEncoding
	hostAddress  = goconf.Config.App.NotifyURL
	WXPMerId     = goconf.Config.MobileApp.WXPMerId
	ALPMerId     = goconf.Config.MobileApp.ALPMerId
	webAppUrl    = goconf.Config.MobileApp.WebAppUrl
)

// register 注册
func (u *user) register(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
		UserName:       req.UserName,
		Password:       req.Password,
		Activate:       "false",
		Limit:          req.Limit,
		RegisterFrom:   req.UserFrom,
		Remark:         req.Remark,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		SubAgentCode:   req.SubAgentCode,
		BelongsTo:      req.BelongsTo,
		InvitationCode: req.InvitationCode,
	}
	user.UpdateTime = user.CreateTime

	// 保存用户信息
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR
	}

	return model.SUCCESS1
}

// login 登录
func (u *user) login(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
			return model.USERNAME_PASSWORD_ERROR
		}
		return model.SYSTEM_ERROR
	}

	//判断锁定情况
	if user.LockTime != "" {
		localTime := time.Now().Format("2006-01-02 15:04:05")
		stemp := ""
		for _, v := range localTime {
			if (v >= '0') && (v <= '9') {
				stemp += string(v)
			}
		}
		localTime = stemp
		date1 := user.LockTime[:8]
		date2 := localTime[:8]
		time1 := user.LockTime[8:]
		time2 := localTime[8:]
		time1Int, err := strconv.ParseInt(string(time1), 10, 32)
		if err != nil {
			log.Errorf("convert string to int error value:%s, error:", time1, err)
			return model.SYSTEM_ERROR
		}
		time2Int, err := strconv.ParseInt(string(time2), 10, 32)
		if err != nil {
			log.Errorf("convert string to int error value:%s, error:", time2, err)
			return model.SYSTEM_ERROR
		}
		if string(date1) != string(date2) {
			time2Int += 240000
		}

		if (time2Int - time1Int) > LOCKTIME {
			mongo.AppUserCol.UpdateLoginTime(req.UserName, "", "") //解锁
		} else {
			return model.USER_LOCK
		}
	}
	// 密码是否正确
	if user.Password != req.Password {
		localTime := time.Now().Format("2006-01-02 15:04:05")
		stemp := ""
		for _, v := range localTime {
			if (v >= '0') && (v <= '9') {
				stemp += string(v)
			}
		}
		localTime = stemp
		if user.LoginTime != "" {
			timeArray := strings.Split(user.LoginTime, ",")
			var count = 0
			var loginTime = ""
			date2 := localTime[:8]
			time2 := localTime[8:]
			time2Int, err := strconv.ParseInt(string(time2), 10, 32)
			if err != nil {
				log.Errorf("convert string to int is value:%s, error:%s", time2, err)
				return model.SYSTEM_ERROR
			}
			for _, timeElement := range timeArray {
				date1 := timeElement[:8]
				time1 := timeElement[8:]

				time1Int, err := strconv.ParseInt(string(time1), 10, 32)
				if err != nil {
					log.Errorf("convert string to int is value:%s, error:%s", time1, err)
					continue
				}

				if string(date1) != string(date2) { //日期不同，加24小时
					time2Int += 240000
				}

				if (time2Int - time1Int) > LOGINDIFFTIME { //超过一小时，舍弃
					continue
				} else {
					count++
					//fmt.Println("the count is :", count)
					//记录
					if loginTime == "" {
						loginTime = timeElement
					} else {
						loginTime += ","
						loginTime += timeElement
					}
					//fmt.Println("the first loginTime is :", loginTime)
				}
			}
			//判断count是否达到10次
			if count == 9 {
				loginTime += ","
				loginTime += localTime
				mongo.AppUserCol.UpdateLoginTime(req.UserName, loginTime, localTime)
				//fmt.Println("the count is 9")
				//fmt.Println("the Transtime is :", localTime)
				return model.USER_LOCK //锁定
			} else {
				var ret model.AppResult
				if count == 6 {
					ret = model.USER_THREE_TIMES
				} else if count == 7 {
					ret = model.USER_TWO_TIMES
				} else if count == 8 {
					ret = model.USER_ONE_TIMES
				} else {
					ret = model.USERNAME_PASSWORD_ERROR
				}

				if loginTime == "" {
					loginTime = localTime
				} else {
					loginTime += ","
					loginTime += localTime
				}
				//fmt.Println("the second logintime is :", loginTime)
				mongo.AppUserCol.UpdateLoginTime(req.UserName, loginTime, "")

				return ret
			}
		} else {
			mongo.AppUserCol.UpdateLoginTime(req.UserName, localTime, "")
			//fmt.Println("the *** logintime is :", localTime)
		}
		return model.USERNAME_PASSWORD_ERROR
	}

	//密码正确，清空登陆记录
	//mongo.AppUserCol.UpdateLoginTime(req.UserName, "", "")
	var userInfo model.AppUser
	userInfo.LoginTime = ""
	userInfo.LockTime = ""
	userInfo.DeviceType = req.AppUser.DeviceType
	userInfo.DeviceToken = req.AppUser.DeviceToken
	userInfo.UserName = req.UserName
	mongo.AppUserCol.UpdateAppUser(&userInfo, mongo.UPDATE_DEVICE_LOCK_INFO)

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
		user.PayUrl = merchant.Detail.PayUrl
		user.MerName = merchant.Detail.MerName
	}

	result = model.AppResult{
		State: model.SUCCESS,
		Error: "",
		User:  user,
	}

	return result
}

// reqActivate 请求发送激活链接
func (u *user) reqActivate(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
			return model.USERNAME_PASSWORD_ERROR
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
func (u *user) activate(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
	user.LoginTime = "" //登录时间
	user.LockTime = ""  //锁定时间
	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR_CH
	}

	return model.SUCCESS1
}

// improveInfo 信息完善
func (u *user) improveInfo(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
			return model.USERNAME_PASSWORD_ERROR
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

	agentCode, agentName := "99911888", "讯联O2O机构"
	if user.InvitationCode != "" {
		agent, err := mongo.AgentColl.Find(user.InvitationCode)
		if err == nil {
			agentCode, agentName = agent.AgentCode, agent.AgentName
		}
	}

	// 创建商户
	permission := []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf}
	merchant := &model.Merchant{
		AgentCode:  agentCode,
		AgentName:  agentName,
		Permission: permission,
		MerStatus:  model.MerStatusNormal,
		TransCurr:  "156",
		Remark:     "app_register",
		RefundType: model.CurrentDayRefund, // 只能当天退
		IsNeedSign: true,
		SignKey:    fmt.Sprintf("%x", randBytes(16)),
		Detail: model.MerDetail{
			MerName:       "云收银",
			CommodityName: "讯联云收银在线注册商户",
			Province:      req.Province,
			City:          req.City,
			OpenBankName:  req.BranchBank,
			BankName:      req.BankOpen,
			BankId:        req.BankNo,
			AcctName:      req.Payee,
			AcctNum:       req.PayeeCard,
			ContactTel:    req.PhoneNum,
		},
	}

	// 生成商户号，并保存商户
	if err := genMerId(merchant, "999118880"); err != nil {
		return model.SYSTEM_ERROR
	}

	// 创建路由,支付宝，微信
	if err := genRouter(merchant); err != nil {
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
	user.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	err = mongo.AppUserCol.Upsert(user)
	if err != nil {
		log.Errorf("save user err,%s", err)
		return model.SYSTEM_ERROR
	}

	result = model.AppResult{
		State: model.SUCCESS,
		Error: "",
		User:  user,
	}

	return result
}

// getTotalTransAmt 查询某天交易总额
func (u *user) getTotalTransAmt(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
			return model.USERNAME_PASSWORD_ERROR
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

	s, _ := query.TransStatistics(&model.QueryCondition{
		MerId:     user.MerId,
		StartTime: month,
		EndTime:   month,
		Size:      1,
		Page:      1,
	})

	result.Count = s.TotalTransNum
	result.TotalAmt = fmt.Sprintf("%0.2f", float64(s.TotalTransAmt)/100)
	return result
}

// getUserBill 获取用户账单
func (u *user) getUserBill(req *reqParams) (result model.AppResult) {

	// 用户名不为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" || req.Status == "" || req.Index == "" {
		return model.PARAMS_EMPTY
	}

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}

	// 不同时为空
	if req.Month == "" && req.Date == "" {
		return model.TIME_ERROR
	}

	if req.Month != "" && !monthRegexp.MatchString(req.Month) {
		return model.TIME_ERROR
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_PASSWORD_ERROR
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	result = model.NewAppResult(model.SUCCESS, "")

	dsDate, deDate, ssDate, seDate := "", "", "", ""
	// 按month来统计
	ym := req.Month
	yearNum, _ := strconv.Atoi(ym[:4])
	month := ym[4:6]
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
	ssDate = ym[:4] + "-" + ym[4:6] + "-" + "01" + " 00:00:00"
	seDate = ym[:4] + "-" + ym[4:6] + "-" + day + " 23:59:59"

	// 按送的日期天数拉取数据
	if req.Date != "" {
		day, _ := strconv.Atoi(req.Date)
		now := time.Now()
		dsDate = now.Add(-time.Hour*24*time.Duration(day)).Format("2006-01-02") + " 00:00:00"
		deDate = now.Format("2006-01-02") + " 23:59:59"
	} else {
		// 没送天数的话，默认当月
		dsDate = ssDate
		deDate = seDate
	}

	index, size := pagingParams(req)
	q := &model.QueryCondition{
		MerId:     user.MerId,
		StartTime: dsDate,
		EndTime:   deDate,
		Size:      size,
		Page:      1,
		Skip:      index,
	}

	// 只包含支付交易
	if req.TransType != 0 {
		q.TransType = req.TransType
	}

	// TODO�������日本项目字段 只包含支付交易
	if req.OrderDetail == "pay" {
		q.TransType = model.PayTrans
	}

	switch req.Status {
	case "all":
	case "success":
		q.Respcd = "00"
	case "fail":
		q.RespcdNotIn = "00"
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
	typeGroup, err := mongo.SpTransColl.MerBills(&model.QueryCondition{
		MerId:        user.MerId,
		StartTime:    ssDate,
		EndTime:      seDate,
		RefundStatus: []int{model.TransRefunded},
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

	if len(txns) == 0 {
		txns = make([]*model.AppTxn, 0)
	}

	result.Txn = txns
	result.Size = len(trans)
	result.RefdCount = refundCount
	result.Count = transCount
	result.TotalRecord = total

	// TODO:先用该字段做判断是日币还是元
	if req.OrderDetail == "pay" {
		result.TotalAmt = fmt.Sprintf("%d", transAmt)
		result.RefdTotalAmt = fmt.Sprintf("%d", refundAmt)
	} else {
		result.TotalAmt = fmt.Sprintf("%0.2f", float32(transAmt)/100)
		result.RefdTotalAmt = fmt.Sprintf("%0.2f", float32(refundAmt)/100)
	}

	return
}

// getUserTrans 获取用户某笔交易信息
func (u *user) getUserTrans(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
			return model.USERNAME_PASSWORD_ERROR
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
		// TODO 兼容所有币种
		if t.Currency == "JPY" {
			result.RefdTotalAmt = fmt.Sprintf("%d", t.RefundAmt)
		} else {
			result.RefdTotalAmt = fmt.Sprintf("%0.2f", float32(t.RefundAmt)/100)
		}
	case "getOrder":
		result.Txn = transToTxn(t)
	}

	return result
}

// passwordHandle 修改密码
func (u *user) passwordHandle(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
			return model.USERNAME_PASSWORD_ERROR
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
func (u *user) promoteLimit(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
			return model.USERNAME_PASSWORD_ERROR
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
func (u *user) getSettInfo(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}

	// 用户名不为空
	if req.UserName == "" || req.Password == "" || req.Transtime == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_PASSWORD_ERROR
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
		BankOpen:   mer.Detail.BankName,
		PayeeCard:  mer.Detail.AcctNum,
		PhoneNum:   mer.Detail.ContactTel,
		Province:   mer.Detail.Province,
		City:       mer.Detail.City,
		BranchBank: mer.Detail.OpenBankName,
		BankNo:     mer.Detail.BankId,
	}

	result.SettInfo = settInfo
	return
}

// updateSettInfo 更新清算信息
func (u *user) updateSettInfo(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
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
				return model.USERNAME_PASSWORD_ERROR
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
	if req.BranchBank != "" {
		m.Detail.OpenBankName = req.BranchBank // 支行对应开户行
	}
	if req.BankOpen != "" {
		m.Detail.BankName = req.BankOpen // 银行对应银行
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

// ticketHandle 处理小票接口
func (u *user) ticketHandle(req *reqParams) (result model.AppResult) {
	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}

	// 必填参数不为空
	if req.UserName == "" || req.Password == "" || req.TicketNum == "" || req.OrderNum == "" {
		return model.PARAMS_EMPTY
	}

	// 根���用���名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_PASSWORD_ERROR
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	// 更新交易
	if user.MerId != "" {
		err = mongo.SpTransColl.UpdateFields(user.MerId, req.OrderNum, "ticketNum", req.TicketNum)
		if err != nil {
			if err.Error() == "not found" {
				return model.NO_TRANS
			} else {
				log.Errorf("update fields fail: %s", err)
				return model.SYSTEM_ERROR
			}
		}
	} else {
		return model.NO_PAY_MER
	}

	return model.SUCCESS1

}

// findOrderHandle 条件组合查找
func (u *user) findOrderHandle(req *reqParams) (result model.AppResult) {

	user, errResult := checkPWD(req)
	if errResult != nil {
		return *errResult
	}

	index, size := pagingParams(req)
	q := &model.QueryCondition{
		OrderNum:  req.OrderNum,
		TransType: model.PayTrans, // 只是支付交易
		MerId:     user.MerId,
		Skip:      index,
		Size:      size,
		Page:      1,
	}

	// 初始化查询参数
	findOrderParams(req, q)
	trans, total, err := mongo.SpTransColl.Find(q)

	if err != nil {
		return model.SYSTEM_ERROR
	}

	var txns []*model.AppTxn
	for _, t := range trans {
		txns = append(txns, transToTxn(t))
	}

	if len(txns) == 0 {
		txns = make([]*model.AppTxn, 0)
	}
	result = model.SUCCESS1
	result.Txn = txns
	result.Count = total
	result.Size = len(txns)
	return
}

// updateMessageHandle 更新推送信息状态
func (u *user) updateMessageHandle(req *reqParams) (result model.AppResult) {

	_, errResult := checkPWD(req)
	if errResult != nil {
		return *errResult
	}

	type msg struct {
		MsgId  string `json:"msgId"`
		Status int    `json:"status"`
	}

	var ms []msg
	err := json.Unmarshal([]byte(req.Message), &ms)
	if err != nil {
		log.Errorf("unmarshal message=%s error: %s", req.Message, err)
		return model.PARAMS_FORMAT_ERROR
	}

	// update
	for _, m := range ms {
		if err = mongo.PushMessageColl.UpdateStatusByID(m.MsgId, m.Status); err != nil {
			log.Errorf("update push message fail, MsgId=%s, error: %s", m.MsgId, err)
		}
	}

	return model.SUCCESS1
}

// 重置密码
func (u *user) forgetPassword(req *reqParams) (result model.AppResult) {
	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}

	log.Debugf("userName=%s", req.UserName)
	// 参数不能为空
	if req.UserName == "" {
		return model.PARAMS_EMPTY
	}

	// 根据用户名查找用户
	_, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		return model.USERNAME_NO_EXIST
	}

	code := fmt.Sprintf("%d", rand.Int31())
	resetPasswordUrl := fmt.Sprintf("%s/master/resetPassword?code=%s", hostAddress, code)
	click := b64Encoding.EncodeToString(randBytes(32))

	email := &email.Email{
		To:    req.UserName,
		Title: resetPassword.Title,
		Body:  fmt.Sprintf(resetPassword.Body, resetPasswordUrl, click),
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

//获取七牛token
func (u *user) getQiniuToken(req *reqParams) (result model.AppResult) {
	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}

	if _, err := checkPWD(req); err != nil {
		return *err
	}

	return model.SUCCESS1
}

//完善证书信息
func (u *user) improveCertInfo(req *reqParams) (result model.AppResult) {

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}

	var user *model.AppUser
	var err error

	// 云收银app
	if req.AppUser == nil {
		appUser, err := checkPWD(req)
		if err != nil {
			return *err
		}
		user = appUser
	} else {
		// 从销售工具接口过来的
		user = req.AppUser
	}

	m, err := mongo.MerchantColl.Find(user.MerId)
	if err != nil {
		return model.MERID_NO_EXIST
	}

	if req.CertName != "" {
		m.Detail.CertName = req.CertName
	}
	if req.CertAddr != "" {
		m.Detail.CertAddr = req.CertAddr
	}
	if req.LegalCertPos != "" {
		m.Detail.LegalCertPos = req.LegalCertPos
	}
	if req.LegalCertOpp != "" {
		m.Detail.LegalCertOpp = req.LegalCertOpp
	}
	if req.BusinessLicense != "" {
		m.Detail.BusinessLicense = req.BusinessLicense
	}
	if req.TaxRegistCert != "" {
		m.Detail.TaxRegistCert = req.TaxRegistCert
	}
	if req.OrganizeCodeCert != "" {
		m.Detail.OrganizeCodeCert = req.OrganizeCodeCert
	}

	if err = mongo.MerchantColl.Update(m); err != nil {
		return model.SYSTEM_ERROR
	}

	req.m = m

	return model.SUCCESS1
}

// findPushMessage 查询某个用户下的推送消息
func (u *user) findPushMessage(req *reqParams) (result model.AppResult) {

	if _, errResult := checkPWD(req); errResult != nil {
		return *errResult
	}

	size, _ := strconv.Atoi(req.Size)
	messages, err := mongo.PushMessageColl.Find(&model.PushMessage{
		UserName: req.UserName,
		LastTime: req.LastTime,
		MaxTime:  req.MaxTime,
		MsgType:  "MSG_TYPE_C", // 只返回MSG_TYPE_C类型消息
		Size:     size,
	})
	if err != nil {
		return model.SYSTEM_ERROR
	}

	if len(messages) == 0 {
		messages = make([]model.PushMessage, 0)
	}

	result = model.SUCCESS1
	result.Count = len(messages)
	result.Message = messages
	return result
}

// 卡券列表
func (u *user) couponsHandler(req *reqParams) (result model.AppResult) {
	//判断必填项不能为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" || req.ClientId == "" || req.Index == "" {
		return model.PARAMS_EMPTY
	}
	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}
	if req.Month != "" && !monthRegexp.MatchString(req.Month) {
		return model.TIME_ERROR
	}

	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		if err.Error() == "not found" {
			return model.USERNAME_PASSWORD_ERROR
		}
		log.Errorf("find database err,%s", err)
		return model.SYSTEM_ERROR
	}

	// 密码不对
	if req.Password != user.Password {
		return model.USERNAME_PASSWORD_ERROR
	}

	if req.Size == "" {
		req.Size = "15"
	}

	result = model.NewAppResult(model.SUCCESS, "")

	dsDate, deDate, ssDate, seDate := "", "", "", ""
	if req.Month != "" {
		// 按month来统计
		ym := req.Month
		yearNum, _ := strconv.Atoi(ym[:4])
		month := ym[4:6]
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
		ssDate = ym[:4] + "-" + ym[4:6] + "-" + "01" + " 00:00:00"
		seDate = ym[:4] + "-" + ym[4:6] + "-" + day + " 23:59:59"
		//month不为空，则返回该月账单情况
		dsDate = ssDate
		deDate = seDate
	} else {
		// month不填默认返回当月账单
		now := time.Now()
		dsDate = now.Format("2006-01")[0:7] + "-" + "01" + " 00:00:00"
		deDate = now.Format("2006-01-02") + " 23:59:59"
	}

	index, size := pagingParams(req)

	q := &model.QueryCondition{
		MerId:       user.MerId,
		StartTime:   dsDate,
		EndTime:     deDate,
		Size:        size,
		Page:        1,
		Skip:        index,
		TransStatus: []string{model.TransSuccess},
		Respcd:      "00",
		Busicd:      "VERI",
	}

	trans, total, err := mongo.CouTransColl.Find(q)
	if err != nil {
		log.Errorf("find user coupon trans error: %s", err)
		return model.SYSTEM_ERROR
	}
	var coupons []*model.Coupon
	for _, t := range trans {
		coupons = append(coupons, transToCoupon(t))
	}
	result.Coupons = coupons
	result.Size = len(coupons)
	result.TotalRecord = total

	if req.Month != "" {
		result.Count = total
	}

	return result
}

// 消息接口
func (u *user) messagePullHandler(req *reqParams) (result model.AppResult) {
	if req.Size == "" {
		req.Size = "15"
	}

	result = u.findPushMessage(req)

	return result
}

func transToTxn(t *model.Trans) *model.AppTxn {
	txn := &model.AppTxn{
		Response:          t.RespCode,
		SystemDate:        timeReplacer.Replace(t.CreateTime),
		ConsumerAccount:   t.ConsumerAccount,
		TransStatus:       t.TransStatus,
		RefundAmt:         t.RefundAmt,
		TicketNum:         t.TicketNum,
		NickName:          t.NickName,
		AvatarUrl:         t.HeadImgUrl,
		CheckCode:         t.VeriCode,
		CouponName:        t.Prodname,
		CouponChannel:     t.ChanCode,
		CouponOrderNo:     t.CouponOrderNum,
		CouponDiscountAmt: t.DiscountAmt,
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
	txn.ReqData.TotalFee = t.TransAmt
	txn.ReqData.ChanCode = t.ChanCode
	txn.ReqData.Currency = t.Currency
	if t.Currency == "" {
		txn.ReqData.Currency = "CNY"
	}

	return txn
}

func randBytes(length int) []byte {
	var randBytes = make([]byte, length)
	if _, err := io.ReadFull(cr.Reader, randBytes[:]); err != nil {
		log.Errorf("io.ReadFull error: %s", err)
	}
	return randBytes
}

func genRouter(merchant *model.Merchant) error {

	// 创建路由,支付宝，微信
	alpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "ALP",
		ChanCode:  "ALP",
		ChanMerId: ALPMerId,
		SettFlag:  "CIL",
		SettRole:  "CIL",
		MerFee:    0.006,
		AcqFee:    0.0,
	}
	err := mongo.RouterPolicyColl.Insert(alpRoute)
	if err != nil {
		log.Errorf("create routePolicy err: %s", err)
		return err
	}

	wxpRoute := &model.RouterPolicy{
		MerId:     merchant.MerId,
		CardBrand: "WXP",
		ChanCode:  "WXP",
		ChanMerId: WXPMerId,
		SettFlag:  "CIL",
		SettRole:  "CIL",
		MerFee:    0.006,
		AcqFee:    0.003,
	}
	err = mongo.RouterPolicyColl.Insert(wxpRoute)
	if err != nil {
		log.Errorf("create routePolicy err: %s", err)
		return err
	}

	return nil
}

func genMerId(merchant *model.Merchant, prefix string) error {
	if prefix == "" {
		return fmt.Errorf("%s", "prefix should not be empty")
	}

	var length = fmt.Sprintf("%d", 15-len(prefix))
	for {
		// 设置merId
		maxMerId, err := mongo.MerchantColl.FindMaxMerId(prefix)
		if err != nil {
			if err.Error() == "not found" {
				merchant.MerId = prefix + fmt.Sprintf("%0"+length+"d", 1)
			} else {
				log.Errorf("find merchant err,%s", err)
				return err
			}

		} else {
			// 找到的话  len(maxMerId)==15
			order := maxMerId[len(prefix):15]
			maxMerIdNum, err := strconv.Atoi(order)
			if err != nil {
				log.Errorf("format maxMerId(%s) err", maxMerId)
				return err
			}
			merchant.MerId = fmt.Sprintf("%s%0"+fmt.Sprintf("%d", len(order))+"d", prefix, maxMerIdNum+1)
		}

		merchant.UniqueId = util.Confuse(merchant.MerId)
		if merchant.Detail.TitleOne != "" || merchant.Detail.TitleTwo != "" {
			merchant.Detail.BillUrl = fmt.Sprintf("%s/trade.html?merchantCode=%s", webAppUrl, merchant.UniqueId)
			merchant.Detail.PayUrl = fmt.Sprintf("%s/index.html?merchantCode=%s", webAppUrl, b64Encoding.EncodeToString([]byte(merchant.MerId)))
		}
		err = mongo.MerchantColl.Insert(merchant)
		if err != nil {
			isDuplicateMerId := strings.Contains(err.Error(), "E11000 duplicate key error index")
			if !isDuplicateMerId {
				log.Errorf("add merchant err: %s, merId=%s", err, merchant.MerId)
				return err
			}
		}
		break
	}
	return nil
}

func checkPWD(req *reqParams) (user *model.AppUser, errResult *model.AppResult) {
	// 参数不能为空
	if req.UserName == "" || req.Password == "" {
		return nil, &model.PARAMS_EMPTY
	}
	// 根据用户名查找用户
	user, err := mongo.AppUserCol.FindOne(req.UserName)
	if err != nil {
		return nil, &model.USERNAME_PASSWORD_ERROR
	}
	// 密码不对
	if req.Password != user.Password {
		return nil, &model.USERNAME_PASSWORD_ERROR
	}
	return user, nil
}

// 解析分页参数
func pagingParams(req *reqParams) (index, size int) {
	var err error
	index, _ = strconv.Atoi(req.Index) // 默认0
	size, err = strconv.Atoi(req.Size) // 默认15
	if err != nil {
		size = 15
	}
	return
}

func findOrderParams(req *reqParams, q *model.QueryCondition) {
	recType, _ := strconv.Atoi(req.RecType)
	payType, _ := strconv.Atoi(req.PayType)
	transStatus, _ := strconv.Atoi(req.Status)

	switch recType {
	case 1:
		q.TradeFrom = []string{model.IOS, model.Android}
	case 2, 8:
		// 暂时没有
		q.TradeFrom = []string{model.Pc}
	case 4:
		q.TradeFrom = []string{model.Wap}
	case 3:
		q.TradeFrom = []string{model.IOS, model.Android}
	case 5:
		q.TradeFrom = []string{model.IOS, model.Android, model.Wap}
	case 15:
		// ignore
	}

	switch payType {
	case 1:
		q.ChanCode = channel.ChanCodeAlipay
	case 2:
		q.ChanCode = channel.ChanCodeWeixin
	case 3:
		// ignore
	}

	switch transStatus {
	case 1:
		// 交易成功
		q.TransStatus = []string{model.TransSuccess}
		q.RefundStatus = []int{model.TransNoRefunded}
	case 2:
		// 部分退
		q.RefundStatus = []int{model.TransPartRefunded}
	case 4:
		// 全额退
		q.RefundStatus = []int{model.TransRefunded}
	case 5:
		// 交易成功、全额退款
		q.RefundStatus = []int{model.TransNoRefunded, model.TransRefunded}
	case 6:
		// 部分、全额退款
		q.RefundStatus = []int{model.TransRefunded, model.TransPartRefunded}
	case 7:
		// 交易成功、部分、全额退
		q.RefundStatus = []int{model.TransNoRefunded, model.TransPartRefunded, model.TransRefunded}
	}
}

func transToCoupon(t *model.Trans) *model.Coupon {
	coupon := &model.Coupon{
		Name:       t.Prodname,
		Channel:    t.ChanCode,
		TradeFrom:  t.TradeFrom,
		Response:   t.RespCode,
		SystemDate: timeReplacer.Replace(t.CreateTime),
		Terminalid: t.Terminalid,
		OrderNum:   t.OrderNum,
	}
	if t.ScanPayCoupon != nil {
		coupon.ReqData.Busicd = t.ScanPayCoupon.Busicd
		coupon.ReqData.AgentCode = t.ScanPayCoupon.AgentCode
		coupon.ReqData.Txndir = "Q"
		coupon.ReqData.Terminalid = t.ScanPayCoupon.Terminalid
		coupon.ReqData.OrderNum = t.ScanPayCoupon.OrderNum
		coupon.ReqData.MerId = t.ScanPayCoupon.MerId
		coupon.ReqData.TradeFrom = t.ScanPayCoupon.TradeFrom
		coupon.ReqData.Txamt = fmt.Sprintf("%012d", t.ScanPayCoupon.TransAmt)
		coupon.ReqData.TotalFee = t.ScanPayCoupon.TransAmt
		coupon.ReqData.ChanCode = t.ScanPayCoupon.ChanCode
		coupon.ReqData.Currency = t.ScanPayCoupon.Currency
		if t.ScanPayCoupon.Currency == "" {
			coupon.ReqData.Currency = "CNY"
		}
		coupon.ReqData.CouponDiscountAmt = t.ScanPayCoupon.DiscountAmt
	}

	// coupon.ReqData.OrigTransAmt = t.ScanPayCoupon.TransAmt + t.ScanPayCoupon.DiscountAmt
	couponType := ""
	if len(t.VoucherType) == 2 {
		couponType = t.VoucherType[1:2]
		if couponType == "2" {
			couponType = "1"
		}
	} else {
		couponType = t.VoucherType
	}
	coupon.Type = couponType
	return coupon
}
