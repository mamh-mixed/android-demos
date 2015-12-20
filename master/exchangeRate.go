package master

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/omigo/log"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
)

// excratQueryHandle 汇率查询
func excratQueryHandle(w http.ResponseWriter, r *http.Request) {
	localCurr := r.FormValue("localCurrency")
	targetCurr := r.FormValue("targetCurrency")
	createTime := r.FormValue("createTime")
	enforcementTime := r.FormValue("enforcementTime")
	enforceUser := r.FormValue("enforceUser")
	rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))

	cond := &model.ExchangeRateManage{
		LocalCurrency:       localCurr,
		TargetCurrency:      targetCurr,
		CreateTime:          createTime,
		PlanEnforcementTime: enforcementTime,
		CheckedUser:         enforceUser,
		Rate:                rate,
	}

	ret := ExcRat.Find(cond, size, page)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// excratActivateHandle 汇率立即生效
func excratActivateHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read all body error: %s", err)
		w.WriteHeader(501)
		return
	}

	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		w.Write([]byte("FIND SESSION ERROR"))
		return
	}

	ret := ExcRat.RateActivate(data, curSession.User.UserName)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// excratCheckHandle 汇率复核
func excratCheckHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read all body error: %s", err)
		w.WriteHeader(501)
		return
	}

	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		w.Write([]byte("FIND SESSION ERROR"))
		return
	}

	ret := ExcRat.RateCheck(data, curSession.User.UserName)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// excratCreateHandle 新增
func excratCreateHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read all body error: %s", err)
		w.WriteHeader(501)
		return
	}

	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		w.Write([]byte("FIND SESSION ERROR"))
		return
	}

	ret := ExcRat.Add(data, curSession.User.UserName)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// excratRangeUpdateHandle 更新汇率环比正常阈值范围的操作
func excratRangeUpdateHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read all body error: %s", err)
		w.WriteHeader(501)
		return
	}

	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		w.Write([]byte("FIND SESSION ERROR"))
		return
	}

	ret := ExcRat.UpdateD2dRatio(data, curSession.User.UserName)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// excratRangeFindHandle 查询汇率环比上下浮动值
func excratRangeFindHandle(w http.ResponseWriter, r *http.Request) {
	ret := ExcRat.d2dRatioRange()

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

type exchangeRate struct{}

var ExcRat exchangeRate

// d2dRatioRange 查询汇率环比值有效浮动范围
func (e *exchangeRate) d2dRatioRange() (result *model.ResultBody) {
	sysConst, err := mongo.SysConstColl.Find()
	if err != nil {
		log.Errorf("FIND D2D RATIO RANGE ERROR: %s", err)
		return model.NewResultBody(2, "EXCHANGE_RATE.TOAST.FIND_RATE_FLOATING_RANGE_ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "SAVE_SUCCESS",
		Data: model.SystemConstant{
			RateFloatingUpper: sysConst.RateFloatingUpper,
			RateFloatingLower: sysConst.RateFloatingLower,
		},
	}
	return result
}

// UpdateD2dRatio 更新汇率环比值合理范围
func (e *exchangeRate) UpdateD2dRatio(data []byte, username string) (result *model.ResultBody) {
	sysConst := new(model.SystemConstant)
	err := json.Unmarshal(data, sysConst)
	if err != nil {
		log.Errorf("JSON(%s) UNMARSHAL ERROR: %s", string(data), err)
		return model.NewResultBody(2, "JSON_RESOLVE_FAIL")
	}

	sc, err := mongo.SysConstColl.Find()
	if err != nil {
		log.Errorf("FIND SYSTEM CONSTANT ERROR: %s", err)
		return model.NewResultBody(2, "SAVE_FAIL")
	}

	sc.RateFloatingUpper = sysConst.RateFloatingUpper
	sc.RateFloatingLower = sysConst.RateFloatingLower
	sc.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	sc.UpdateUser = username

	err = mongo.SysConstColl.Upsert(sc)
	if err != nil {
		log.Errorf("SAVE SYSTEM CONSTANT ERROR: %s", err)
		return model.NewResultBody(2, "SAVE_FAIL")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "SAVE_SUCCESS",
	}
	return result
}

// Add 新增一个
func (e *exchangeRate) Add(data []byte, username string) (result *model.ResultBody) {
	t := new(model.ExchangeRateManage)
	err := json.Unmarshal(data, t)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_RESOLVE_FAIL")
	}

	if t.LocalCurrency == "" {
		log.Errorf("MISSING_REQUIRED_ITEM: %s", "localCurrency")
		result = model.NewResultBody(400, "MISSING_REQUIRED_ITEM")
		result.Data = "localCurrency"
		return result
	}

	if t.TargetCurrency == "" {
		log.Errorf("MISSING_REQUIRED_ITEM: %s", "targetCurrency")
		result = model.NewResultBody(400, "MISSING_REQUIRED_ITEM")
		result.Data = "targetCurrency"
		return result
	}

	if t.Rate == 0.0 {
		log.Errorf("MISSING_REQUIRED_ITEM: %s", "rate")
		result = model.NewResultBody(400, "MISSING_REQUIRED_ITEM")
		result.Data = "rate"
		return result
	}

	if t.PlanEnforcementTime == "" {
		log.Errorf("MISSING_REQUIRED_ITEM: %s", "planEnforcementTime")
		result = model.NewResultBody(400, "MISSING_REQUIRED_ITEM")
		result.Data = "planEnforcementTime"
		return result
	}

	// 校验日期时间格式
	planTime, err := time.ParseInLocation("2006-01-02 15:04:05", t.PlanEnforcementTime, time.Local)
	if err != nil {
		log.Errorf("Unsupport time format: %s.Time format must be 'YYYY-MM-DD hh:mm:ss'", t.PlanEnforcementTime)
		return model.NewResultBody(2, "UNSUPPORT_TIME_FORMAT")
	}

	// 生效时点必须在未来
	if planTime.Before(time.Now()) {
		log.Errorf("Exchange rate enforcement time must be in the future.But now(%s) is %s", time.Now(), planTime)
		return model.NewResultBody(2, "ENFORCEMENT_TIME_OUTDATE")
	}

	// 查询当前生效的汇率
	enforceRate, err := mongo.ExchangeRateColl.FindOne(t.LocalCurrency, t.TargetCurrency)
	// 新增汇率是否已经生效
	rateExist := true
	if err != nil {
		if err.Error() != "not found" {
			log.Errorf("FIND OLD RATES(%s<=>%s) ERROR: %s", t.LocalCurrency, t.TargetCurrency, err)
			return model.NewResultBody(2, "SAVE_RATE_FAIL")
		}

		rateExist = false
	}

	// 新录入的汇率关联的币种已经存在已激活汇率中，校验阈值
	if rateExist {
		// 校验是否超过阈值
		sysConst, err := mongo.SysConstColl.Find()
		if err != nil {
			log.Errorf("FIND SYSTEM CONSTANT ERROR: %s", err)
			return model.NewResultBody(2, "SAVE_RATE_FAIL")
		}

		// 环比值
		d2dRatio := t.Rate / enforceRate.Rate
		log.Debugf("**************%f=%f/%f********", d2dRatio, t.Rate, enforceRate.Rate)
		if d2dRatio < sysConst.RateFloatingLower || d2dRatio > sysConst.RateFloatingUpper {
			log.Errorf("D2D_RATIO(%f=%f/%f) OUT OF ALLOW RANGE[%f, %f]", d2dRatio, t.Rate, enforceRate.Rate, sysConst.RateFloatingLower, sysConst.RateFloatingUpper)
			return model.NewResultBody(3, "OUT_OF_RANGE")
		}
	}

	t.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	t.CreateUser = username
	t.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	t.UpdateUser = username
	t.Status = model.ER_UNCHECKED
	t.EId = util.SerialNumber()

	err = mongo.ExchangeRateManageColl.Add(t)
	if err != nil {
		log.Errorf("ADD_RATE_ERROR: %s", err)
		return model.NewResultBody(2, "SAVE_RATE_FAIL")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "SAVE_RATE_SUCCESS",
		Data:    t,
	}
	return result
}

// RateActivate 处理汇率立即生效的方法
func (e *exchangeRate) RateActivate(data []byte, username string) (result *model.ResultBody) {
	t := new(model.ExchangeRateManage)
	err := json.Unmarshal(data, t)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "TOAST.JSON_RESOLVE_FAIL")
	}

	if t.EId == "" {
		log.Error("缺失必要参数 eID")
		result = &model.ResultBody{
			Status:  400,
			Message: "TOAST.MISSING_REQUIRED_ITEM",
			Data:    "eId",
		}
		return result
	}

	// 查找要立即生效的汇率记录
	rate, err := mongo.ExchangeRateManageColl.FindOne(t.EId)
	if err != nil {
		log.Errorf("Find exchange Rate（%s）FAIL： %s", t.EId, err)
		return model.NewResultBody(401, "EXCHANGE_RATE.TOAST.ACTIVATE_FAIL")
	}

	// 待存入的有效汇率表中的数据
	acRt := &model.ExchangeRate{
		CurrencyPair:    rate.LocalCurrency + "<=>" + rate.TargetCurrency,
		Rate:            rate.Rate,
		EnforcementTime: time.Now().Format("2006-01-02 15:04:05"),
		EnforceUser:     username,
	}
	err = mongo.ExchangeRateColl.Upsert(acRt)
	if err != nil {
		log.Errorf("Activate exchange rate when upsert into database error: %s", err)
		return model.NewResultBody(2, "EXCHANGE_RATE.TOAST.ACTIVATE_FAIL")
	}

	// 更新到汇率管理表中
	rate.Status = model.ER_ACTIVATED
	rate.UpdateUser = acRt.EnforceUser
	rate.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	err = mongo.ExchangeRateManageColl.Update(rate)
	if err != nil {
		log.Errorf("Update rate into exchangeRateManage error: %s", err)
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "EXCHANGE_RATE.TOAST.ACTIVATE_SUCCESS",
		Data:    rate,
	}

	return result
}

// RateCheck 汇率审核通过，建立定时任务让汇率生效
func (e *exchangeRate) RateCheck(data []byte, username string) (result *model.ResultBody) {
	t := new(model.ExchangeRateManage)
	err := json.Unmarshal(data, t)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "JSON_RESOLVE_FAIL")
	}

	if t.EId == "" {
		log.Error("缺失必要参数 eID")
		result = &model.ResultBody{
			Status:  400,
			Message: "MISSING_REQUIRED_ITEM",
			Data:    "eId",
		}
		return result
	}

	// 查找要审核的汇率记录
	rate, err := mongo.ExchangeRateManageColl.FindOne(t.EId)
	if err != nil {
		log.Errorf("Find exchange Rate（%s）FAIL： %s", t.EId, err)
		return model.NewResultBody(401, "ACTIVATE_FAIL")
	}

	// 计划生效时点必须在将来
	planTime, err := time.ParseInLocation("2006-01-02 15:04:05", rate.PlanEnforcementTime, time.Local)
	if err != nil {
		log.Errorf("Unsupport time format: %s format must be 'YYYY-MM-DD hh:mm:ss'", rate.PlanEnforcementTime)
		return model.NewResultBody(2, "UNSUPPORT_TIME_FORMAT")
	}

	// 生效时点必须在未来
	if planTime.Before(time.Now()) {
		log.Errorf("Exchange rate enforcement time must be in the future.But now(%s) is %s", time.Now(), planTime)
		return model.NewResultBody(2, "ENFORCEMENT_TIME_OUTDATE")
	}

	// 更新汇率管理表中的记录
	rate.Status = model.ER_CHECKED
	rate.CheckedUser = username
	rate.CheckedTime = time.Now().Format("2006-01-02 15:04:05")
	rate.UpdateUser = username
	rate.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	err = mongo.ExchangeRateManageColl.Update(rate)
	if err != nil {
		log.Errorf("Update exchange rate（%s）fail： %s", t.EId, err)
		return model.NewResultBody(401, "ACTIVATE_FAIL")
	}

	// 定时任务，到时候激活。既把当前记录存储到汇率表中
	d := planTime.Sub(time.Now())
	go func() {
		log.Infof("The rate of %s<=>%s will be activate after %s", rate.LocalCurrency, rate.TargetCurrency, d)
		time.AfterFunc(d, func() {
			log.Infof("******Exchange rate check pass: %#v*****", rate)
			// 待存储数据库的已激活的数据
			acRt := &model.ExchangeRate{
				CurrencyPair:    rate.LocalCurrency + "<=>" + rate.TargetCurrency,
				Rate:            rate.Rate,
				EnforcementTime: time.Now().Format("2006-01-02 15:04:05"),
				EnforceUser:     rate.CheckedUser,
			}

			err = mongo.ExchangeRateColl.Upsert(acRt)
			if err != nil {
				log.Errorf("Activate exchange rate when upsert into database error: %s", err)
				// TODO 汇率激活失败应该有个处理机制
			} else {
				rate.Status = model.ER_ACTIVATED
				rate.UpdateUser = acRt.EnforceUser
				rate.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
				err = mongo.ExchangeRateManageColl.Update(rate)
				if err != nil {
					// TODO 汇率更新失败应该有个处理机制
				}
			}
		})
	}()
	result = &model.ResultBody{
		Status:  0,
		Message: "ACTIVATE_SUCCESS",
		Data:    rate,
	}

	return result
}

// Find 分页查询
func (e *exchangeRate) Find(cond *model.ExchangeRateManage, size, page int) (result *model.ResultBody) {
	if page < 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 10
	}

	results, total, err := mongo.ExchangeRateManageColl.PaginationFind(cond, size, page)
	if err != nil {
		log.Errorf("查询汇率列表出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(results),
		Data:  results,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}
