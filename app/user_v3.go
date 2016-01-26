package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/push"

	"github.com/CardInfoLink/log"
)

type userV3 struct{}

// UserV3 app v3版本的相关逻辑代码
var UserV3 userV3

func PushMessageToClient(title, content string) (int, error) {
	apps, err := mongo.AppUserCol.Find(&model.AppUserContiditon{})
	if err != nil {
		return 0, err
	}

	var ss int
	for _, u := range apps {
		if u.DeviceToken != "" && u.MerId != "" {
			ss++
			var to = strings.ToLower(u.DeviceType)
			push.Do(&model.PushMessageReq{
				MerID:       u.MerId,
				UserName:    u.UserName,
				Title:       content,
				Message:     title,
				DeviceToken: u.DeviceToken,
				MsgType:     MsgType_C,
				To:          to,
				// OrderNum:    "",
			})
		}
	}

	log.Infof("send %d message success ...", ss)
	return ss, nil
}

// 拉取消息的处理器
func (u *userV3) messagePullHandler(req *reqParams) (result model.AppResult) {
	// 非空校验
	if req.UserName == "" || req.Transtime == "" || req.Password == "" || req.Size == "" {
		return model.PARAMS_EMPTY
	}

	// 校验size是不是非数字的
	if !regexDigit.MatchString(req.Size) {
		return model.InvalidSizeParams
	}

	// lastTime 和 maxTime 不能同时出现
	if req.LastTime != "" {
		req.MaxTime = ""
	}

	result = u.findPushMessage(req)

	return result
}

// findPushMessage 查询某个用户下的推送消息
func (u *userV3) findPushMessage(req *reqParams) (result model.AppResult) {

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
	result.Size = len(messages)
	result.Message = messages
	return result
}

// getUserBills 获取账单
func (u *userV3) getUserBills(req *reqParams) (result model.AppResult) {
	// 用户名不为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" || req.Status == "" || req.Index == "" {
		return model.PARAMS_EMPTY
	}

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}

	hasMonth := true
	// 如果为空，默认当月
	if req.Month == "" {
		hasMonth = false
		req.Month = time.Now().Format("200601")
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

	// 统计的开始时间和结束时间
	startTime, err := time.ParseInLocation("200601", req.Month, time.Local)
	if err != nil {
		log.Errorf("Invalid date format is 'month': %s", err)
		return model.TIME_ERROR
	}
	endTime := startTime.AddDate(0, 1, 0).Add(-time.Second)

	// 构建查询条件
	startTimeStr, endTimeStr := startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05")
	index, size := pagingParams(req)
	q := &model.QueryCondition{
		MerId:     user.MerId,
		StartTime: startTimeStr,
		EndTime:   endTimeStr,
		Size:      size,
		Page:      1,
		Skip:      index,
		TransType: model.PayTrans, // APP v3版本只返回支付订单
	}

	// 交易状态
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

	// 承载返回结果的数组
	result = model.NewAppResult(model.SUCCESS, "")
	// 如果刚好刷新到最后一条
	if index+len(trans) == total {
		// 查找接下来哪个月有数据
		lt, err := mongo.SpTransColl.FindLastRecord(q)
		if err == nil {
			if len(lt.CreateTime) > 7 {
				result.NextMonth = lt.CreateTime[0:4] + lt.CreateTime[5:7] // 200601
			}
		}
	}

	var txns []*model.AppTxn
	for _, t := range trans {
		// 遍历查询结果
		txns = append(txns, transToTxn(t))
	}

	if len(txns) == 0 {
		txns = make([]*model.AppTxn, 0)
	}
	result.Txn = txns
	result.Size = len(txns)
	result.TotalRecord = total

	// 如果APP传递了月份，则需要返回total，count，fefdtotal，refdcount
	if hasMonth {
		result.TotalFee = 0
		result.Count = 0
		result.RefdTotalFee = 0
		result.RefdCount = 0

		typeGroup, err := mongo.SpTransColl.MerBills(&model.QueryCondition{
			MerId:        user.MerId,
			StartTime:    startTimeStr,
			EndTime:      endTimeStr,
			RefundStatus: []int{model.TransRefunded},   // 1: 已退款
			TransStatus:  []string{model.TransSuccess}, // '30': 交易成功
		})
		if err != nil {
			return model.SYSTEM_ERROR
		}

		for _, v := range typeGroup {
			switch v.TransType {
			case model.PayTrans:
				result.TotalFee += v.TransAmt
				result.Count += v.TransNum
			default:
				result.RefdTotalFee += v.TransAmt
				result.TotalFee -= v.TransAmt
				result.RefdCount += v.TransNum
			}
		}
	}

	return result
}

// findOrderHandle 查找订单处理器
func (u *userV3) findOrderHandle(req *reqParams) (result model.AppResult) {
	// 用户名不为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" || req.Index == "" {
		return model.PARAMS_EMPTY
	}

	// 校验index是不是非数字的
	if !regexDigit.MatchString(req.Index) {
		return model.InvalidIndexParams
	}

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
	for i, t := range trans {
		log.Debugf("%d: trans is %+v", i, t)
		txns = append(txns, transToTxn(t))
	}
	log.Debugf("txns's length is %d", len(txns))

	if len(txns) == 0 {
		txns = make([]*model.AppTxn, 0)
	}
	result = model.SUCCESS1
	result.Txn = txns
	result.TotalRecord = total
	result.Size = len(txns)
	return
}

// getDaySummary 获取单日汇总的处理
func (u *userV3) getDaySummary(req *reqParams) (result model.AppResult) {
	// req.BusinessType 表示 报表类型。"1":收款账单；"2":卡券账单
	// 必填字段不为空
	if req.UserName == "" || req.Transtime == "" || req.Password == "" || req.BusinessType == "" || req.Date == "" {
		return model.PARAMS_EMPTY
	}

	// 字段长度验证
	if result, ok := requestDataValidate(req); !ok {
		return result
	}
	// 报表类型验证
	if req.BusinessType != "1" && req.BusinessType != "2" {
		return model.INVALID_REPORT_TYPE
	}
	// 验证日期格式
	startDate, err := time.ParseInLocation("20060102", req.Date, time.Local)
	if err != nil {
		return model.TIME_ERROR
	}
	formatDate := startDate.Format("2006-01-02")
	dsDate := formatDate + " 00:00:00"
	deDate := formatDate + " 23:59:59"

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
	q := &model.QueryCondition{
		MerId:        user.MerId,
		StartTime:    dsDate,
		EndTime:      deDate,
		RefundStatus: []int{model.TransRefunded},   // 1: 已退款
		TransStatus:  []string{model.TransSuccess}, // '30': 交易成功
	}
	if req.BusinessType == "1" {
		q.IsRelatedCoupon = false
		typeGroup, err := mongo.SpTransColl.MerBills(q)
		if err != nil {
			return model.SYSTEM_ERROR
		}
		for _, v := range typeGroup {
			switch v.TransType {
			case model.PayTrans:
				result.TotalFee += v.TransAmt
				result.Count += v.TransNum
			default:
				result.TotalFee -= v.TransAmt
			}
		}

	} else if req.BusinessType == "2" {
		//totalFee:返回的是原始总金额
		q.IsRelatedCoupon = true
		typeGroup, err := mongo.SpTransColl.MerBills(q)
		if err != nil {
			return model.SYSTEM_ERROR
		}
		for _, v := range typeGroup {
			switch v.TransType {
			case model.PayTrans:
				origTransAmt := v.TransAmt + v.DiscountAmt
				result.TotalFee += origTransAmt
			}
		}
		q := &model.QueryCondition{
			MerId:       user.MerId,
			StartTime:   dsDate,
			EndTime:     deDate,
			Size:        15,
			Page:        1,
			Skip:        0,
			TransStatus: []string{model.TransSuccess},
			Respcd:      "00",
			Busicd:      "VERI",
		}

		_, total, err := mongo.CouTransColl.Find(q)
		if err != nil {
			log.Errorf("find user coupon trans error: %s", err)
			return model.SYSTEM_ERROR
		}
		result.Count = total

	}
	return result
}
