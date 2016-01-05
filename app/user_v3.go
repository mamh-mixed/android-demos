package app

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"

	"github.com/omigo/log"
)

type userV3 struct{}

var UserV3 userV3

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
	var txns []*model.AppTxn
	for _, t := range trans {
		// 遍历查询结果
		txns = append(txns, transToTxn(t))
	}
	result.Txn = txns
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
				result.Count += 1
			case model.RefundTrans:
				result.RefdTotalFee += v.TransAmt
				result.TotalFee -= v.TransAmt
				result.RefdCount += 1
			}
		}
	}

	return result
}

// getDaySummary 获取单日汇总的处理
func (u *userV3) getDaySummary(req *reqParams) (result model.AppResult) {
	// // req.BusinessType 表示 报表类型。"1":收款账单；"2":卡券账单
	// // 必填字段不为空
	// if req.UserName == "" || req.Transtime == "" || req.Password == "" || req.BusinessType == "" || req.Date == "" {
	// 	return model.PARAMS_EMPTY
	// }
	//
	// // 字段长度验证
	// if result, ok := requestDataValidate(req); !ok {
	// 	return result
	// }
	//
	// // 报表类型
	// if req.BusinessType != "1" && req.BusinessType != "2" {
	// 	return model.INVALID_REPORT_TYPE
	// }
	//
	// // 验证日期格式
	// startTime, err := time.Parse("20060102", req.Date)
	// if err != nil {
	// 	return model.TIME_ERROR
	// }
	// endTime := startTime.AddDate(0, 0, 1).Add(-time.Second)
	// startTimeStr, endTimeStr := startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05")

	return result
}
