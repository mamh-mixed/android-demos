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

	cond := &model.ExchangeRate{
		LocalCurrency:         localCurr,
		TargetCurrency:        targetCurr,
		CreateTime:            createTime,
		ActualEnforcementTime: enforcementTime,
		EnforceUser:           enforceUser,
		Rate:                  rate,
	}

	ret := ExcRat.Find(cond, size, page)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// excratActivateHandle 激活
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

	ret := ExcRat.Activate(data, curSession.User.UserName)

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
	t := new(model.ExchangeRate)
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

	// 校验阈值
	checkCond := &model.ExchangeRate{
		LocalCurrency:  t.LocalCurrency,
		TargetCurrency: t.TargetCurrency,
		IsEnforced:     true,
	}

	log.Debugf("************ checkCond is %#v", checkCond)

	olds, _, err := mongo.ExchangeRateColl.PaginationFind(checkCond, 10, 1)
	if err != nil {
		log.Errorf("FIND OLD RATES(%s<=>%s) ERROR: %s", checkCond.LocalCurrency, checkCond.TargetCurrency, err)
		return model.NewResultBody(2, "SAVE_RATE_FAIL")
	}

	if len(olds) > 0 {
		// 校验是否超过阈值
		lastRate := olds[0].Rate
		sysConst, err := mongo.SysConstColl.Find()
		if err != nil {
			log.Errorf("FIND SYSTEM CONSTANT ERROR: %s", err)
			return model.NewResultBody(2, "SAVE_RATE_FAIL")
		}

		d2dRatio := t.Rate / lastRate
		log.Debugf("**************%f=%f/%f********", d2dRatio, t.Rate, lastRate)
		if d2dRatio < sysConst.RateFloatingLower || d2dRatio > sysConst.RateFloatingUpper {
			log.Errorf("D2D_RATIO(%f=%f/%f) OUT OF ALLOW RANGE[%f, %f]", d2dRatio, t.Rate, lastRate, sysConst.RateFloatingLower, sysConst.RateFloatingUpper)
			return model.NewResultBody(3, "OUT_OF_RANGE")
		}
	}

	t.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	t.CreateUser = username
	t.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	t.UpdateUser = username

	t.EId = util.SerialNumber()

	err = mongo.ExchangeRateColl.Add(t)
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

// FindOne 查找一个
func (e *exchangeRate) Activate(data []byte, username string) (result *model.ResultBody) {
	t := new(model.ExchangeRate)
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

	rate, err := mongo.ExchangeRateColl.FindOne(t.EId)
	if err != nil {
		log.Errorf("查找汇率（%s）失败： %s", t.EId, err)
		return model.NewResultBody(401, "ACTIVATE_FAIL")
	}

	rate.IsEnforced = true
	rate.EnforceUser = username
	rate.UpdateUser = username
	rate.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	rate.ActualEnforcementTime = time.Now().Format("2006-01-02 15:04:05")

	err = mongo.ExchangeRateColl.Update(rate)
	if err != nil {
		log.Errorf("更新汇率（%s）失败： %s", t.EId, err)
		return model.NewResultBody(401, "ACTIVATE_FAIL")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "ACTIVATE_SUCCESS",
		Data:    rate,
	}

	return result
}

// Find 分页查询
func (e *exchangeRate) Find(cond *model.ExchangeRate, size, page int) (result *model.ResultBody) {
	if page < 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 10
	}

	results, total, err := mongo.ExchangeRateColl.PaginationFind(cond, size, page)
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
