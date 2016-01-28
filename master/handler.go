package master

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/quickpay/settle"
	"github.com/CardInfoLink/quickpay/util"

	"github.com/omigo/log"
)

var maxReportRec = 10000

// appLocaleHandle 网关展示语言
func appLocaleHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	locale := r.FormValue("locale")
	log.Debugf("LOCALE is %s", locale)
	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		w.Write([]byte("FIND SESSION ERROR"))
		return
	}

	// 查看是否支持该语言
	if !IsLocaleExist(locale) {
		log.Warnf("Unsupport language: %s", locale)
		w.Write([]byte("UNSUPPORT LANGUAGE"))
		return
	}

	curSession.Locale = locale
	err = mongo.SessionColl.Update(curSession)
	if err != nil {
		log.Errorf("update session(id=%s) fail:%s", curSession.SessionID, err)
	}

	w.Write([]byte("SUCCESS"))
}

// tradeSettleQueryHandle 交易划款报表查询的
func tradeTransferQueryHandle(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")
	date := r.FormValue("date")
	reportType, _ := strconv.Atoi(r.FormValue("reportType"))
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	ret := tradeSettleReportQuery(role, date, reportType, size, page)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// tradeTransferReportHandle 交易划款报表明细
func tradeTransferReportHandle(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")
	date := r.FormValue("date")
	utcOffset, _ := strconv.Atoi(r.FormValue("utcOffset"))
	fn := strings.Replace(date, "-", "", -1) + "_" + role + ".xlsx"

	tradeTransferReport(w, &model.QueryCondition{
		Date:      date,
		SettRole:  role,
		UtcOffset: utcOffset * 60, // second
	}, fn)
}

// tradeSettleJournalHandle 交易流水，勾兑后的交易
func tradeSettleJournalHandle(w http.ResponseWriter, r *http.Request) {

	// 页面参数
	date := r.FormValue("date")
	utcOffset, _ := strconv.Atoi(r.FormValue("utcOffset"))

	// session参数
	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		return
	}

	// 设置交易查询权限
	q := &model.QueryCondition{
		IsForReport:  true,
		Date:         date, // 北京时间
		AgentCode:    curSession.User.AgentCode,
		MerId:        curSession.User.MerId,
		SubAgentCode: curSession.User.SubAgentCode,
		GroupCode:    curSession.User.GroupCode,
		UtcOffset:    utcOffset * 60, // second
		Locale:       curSession.Locale,
		UserType:     curSession.UserType,
	}

	// 下载
	tradeSettJournalReport(w, q)

}

// tradeSettleReportHandle 交易流水汇总，勾兑后的交易
func tradeSettleReportHandle(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date") // 北京时间

	// session参数
	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		return
	}

	// condition
	q := &model.QueryCondition{
		IsForReport:  true,
		Date:         date, // 北京时间
		AgentCode:    curSession.User.AgentCode,
		MerId:        curSession.User.MerId,
		SubAgentCode: curSession.User.SubAgentCode,
		GroupCode:    curSession.User.GroupCode,
		Locale:       curSession.Locale,
		UserType:     curSession.UserType,
	}

	// 导出
	tradeSettReport(w, q)
}

// tradeSettleRefreshHandle 重新勾兑交易数据
func tradeSettleRefreshHandle(w http.ResponseWriter, r *http.Request) {
	// TODO
	date := r.FormValue("date")
	key := r.FormValue("key")
	if key != "cilxl12345$" {
		return
	}
	log.Infof("process refresh transSettle")
	go settle.RefreshSpTransSett(date)
}

// respCodeMatchHandle 查找应答码处理器
func respCodeMatchHandle(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	ret := RespCode.FindOne(code)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func tradeMsgHandle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	msgType, _ := strconv.Atoi(params.Get("msgType"))
	size, _ := strconv.Atoi(params.Get("size"))
	page, _ := strconv.Atoi(params.Get("page"))

	q := &model.QueryCondition{
		MerId:    params.Get("merId"),
		OrderNum: params.Get("orderNum"),
		ReqIds:   params["reqIds"],
		Page:     page,
		Size:     size,
	}

	ret := getTradeMsg(q, msgType)
	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
		http.Error(w, "system error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(retBytes)
}

func tradeQueryHandle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var merId = params.Get("merId")
	size, err := strconv.Atoi(params.Get("size"))
	if err != nil {
		http.Error(w, "参数 `size` 必须为整数", http.StatusBadRequest)
	}
	if size > 100 || size <= 0 {
		size = 20
	}
	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		http.Error(w, "参数 `page` 必须为整数", http.StatusBadRequest)
	}
	if page <= 0 {
		page = 1
	}

	transType, _ := strconv.Atoi(params.Get("transType"))
	cond := &model.QueryCondition{
		MerName:         params.Get("merName"),
		Terminalid:      params.Get("terminalId"),
		MerId:           merId,
		AgentCode:       params.Get("agentCode"),
		SubAgentCode:    params.Get("subAgentCode"),
		GroupCode:       params.Get("groupCode"),
		TransType:       transType,
		Respcd:          params.Get("respcd"),
		Busicd:          params.Get("busicd"),
		StartTime:       params.Get("startTime"),
		EndTime:         params.Get("endTime"),
		OrderNum:        params.Get("orderNum"),
		OrigOrderNum:    params.Get("origOrderNum"),
		Col:             params.Get("pay"),
		BindingId:       params.Get("bindingId"),
		CouponsNo:       params.Get("couponsNo"),
		WriteoffStatus:  params.Get("writeoffStatus"),
		VoucherType:     params.Get("voucherType"),
		CouponPayStatus: params.Get("couponPayStatus"),
		Prodname:        params.Get("prodname"),
		Size:            size,
		Page:            page,
	}

	transStatus := params.Get("transStatus")
	if transStatus != "" {
		cond.TransStatus = []string{transStatus}
	}

	ret := tradeQuery(cond)

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
		http.Error(w, "system error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(retBytes)
}
func tradeFindOneHandle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	cond := &model.QueryCondition{
		Busicd:       params.Get("busicd"),
		OrigOrderNum: params.Get("origOrderNum"),
		AgentCode:    params.Get("agentCode"),
		SubAgentCode: params.Get("subAgentCode"),
		GroupCode:    params.Get("groupCode"),
		MerId:        params.Get("merId"),
	}
	ret := tradeFindOne(cond)

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
		http.Error(w, "system error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(retBytes)
}

func tradeReportHandle(w http.ResponseWriter, r *http.Request) {
	// get session
	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		return
	}

	// 查询参数
	params := r.URL.Query()

	// 时区偏移量，前端传过来是分
	utcOffset, err := strconv.Atoi(params.Get("utcOffset"))
	if err != nil {
		return
	}

	cond := &model.QueryCondition{
		MerName:        params.Get("merName"),
		Terminalid:     params.Get("terminalId"),
		MerId:          params.Get("merId"),
		AgentCode:      params.Get("agentCode"),
		SubAgentCode:   params.Get("subAgentCode"),
		GroupCode:      params.Get("groupCode"),
		Respcd:         params.Get("respcd"),
		Busicd:         params.Get("busicd"),
		StartTime:      params.Get("startTime"),
		EndTime:        params.Get("endTime"),
		OrderNum:       params.Get("orderNum"),
		OrigOrderNum:   params.Get("origOrderNum"),
		Col:            params.Get("pay"),
		BindingId:      params.Get("bindingId"),
		CouponsNo:      params.Get("couponsNo"),
		WriteoffStatus: params.Get("writeoffStatus"),
		Page:           1,
		Locale:         curSession.Locale,
		UtcOffset:      utcOffset * 60,
		IsForReport:    true,
	}

	transStatus := params.Get("transStatus")
	if transStatus != "" {
		cond.TransStatus = []string{transStatus}
	}

	// 报表需求来自汇总
	if params.Get("from") == "summary" {
		cond.TransStatus = []string{model.TransSuccess}
		cond.RefundStatus = []int{model.TransRefunded}
	}

	// 如果前台传过来‘按商户号分组’的条件，解析成bool成功的话就赋值，不���功的话���不处理，默认为false
	isAggreByGroup, err := strconv.ParseBool(r.FormValue("isAggregateByGroup"))
	if err == nil {
		cond.IsAggregateByGroup = isAggreByGroup
	}

	var merId = params.Get("merId")
	if !isAggreByGroup {
		cond.MerId = merId
	} else {
		// 按照商户号分组的，如果 ‘merId’ 以 ‘GC-’ 开头的，只取 ‘GC-’ 之后的部分，放置到 ‘merId’ 之中
		// 否则把 ‘merId’ 赋值给 groupCode
		if matched, _ := regexp.MatchString(`^GC-`, merId); !matched {
			cond.GroupCode = merId
		} else {
			cond.MerId = merId[3:]
		}
	}

	tradeReport(w, cond, "trade_detail.xlsx", params.Get("from"))
}

func tradeQueryStatsHandle(w http.ResponseWriter, r *http.Request) {

	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))
	q := &model.QueryCondition{
		AgentCode:    r.FormValue("agentCode"),
		SubAgentCode: r.FormValue("subAgentCode"),
		GroupCode:    r.FormValue("groupCode"),
		MerId:        r.FormValue("merId"),
		Page:         page,
		Size:         size,
		MerName:      r.FormValue("merName"),
		StartTime:    r.FormValue("startTime"),
		EndTime:      r.FormValue("endTime"),
	}

	// 如果前台传过来‘按商户号分组’的条件，解析成bool成功的话就赋值，不成功的话就不处理，默认为空
	isAggreByGroup, err := strconv.ParseBool(r.FormValue("isAggregateByGroup"))
	if err == nil {
		q.IsAggregateByGroup = isAggreByGroup
	}

	ret := tradeQueryStats(q)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func tradeQueryStatsReportHandle(w http.ResponseWriter, r *http.Request) {

	// 语言环境
	curSession, err := Session.Get(r)
	if err != nil {
		log.Error("fail to find session")
		return
	}
	params := r.URL.Query()

	// 时区偏移量，前端传过来是分
	utcOffset, err := strconv.Atoi(params.Get("utcOffset"))
	if err != nil {
		return
	}

	// 查询条件
	q := &model.QueryCondition{
		MerId:        params.Get("merId"),
		AgentCode:    params.Get("agentCode"),
		SubAgentCode: params.Get("subAgentCode"),
		MerName:      params.Get("merName"),
		GroupCode:    params.Get("groupCode"),
		StartTime:    params.Get("startTime"),
		EndTime:      params.Get("endTime"),
		Page:         1,
		Size:         maxReportRec,
		Locale:       curSession.Locale,
		UtcOffset:    utcOffset * 60,
	}

	// 如果前台传过来‘按商户号分组’的条件，解析成bool成功的话就赋值，不成功的话就不处理，默认为空
	isAggreByGroup, err := strconv.ParseBool(r.FormValue("isAggregateByGroup"))
	if err == nil {
		q.IsAggregateByGroup = isAggreByGroup
	}

	// 设置content-type
	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`, fmt.Sprintf(`attachment; filename="%s"`, "trade_summary.xlsx"))

	// 导出
	statTradeReport(w, q)
}

func merchantFindHandle(w http.ResponseWriter, r *http.Request) {
	createTime := r.FormValue("createTime")
	createStartTime := ""
	createEndTime := ""
	if createTime != "" {
		createStartTime = createTime + " 00:00:00"
		createEndTime = createTime + " 23:59:59"
	}

	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	pay := r.FormValue("pay")
	isNeedSignStr := r.FormValue("isNeedSign")
	isNeedSign := false
	if isNeedSignStr == "on" {
		isNeedSign = true
	}
	merchant := model.Merchant{
		MerId:        r.FormValue("merId"),
		AgentCode:    r.FormValue("agentCode"),
		AgentName:    r.FormValue("agentName"),
		SubAgentCode: r.FormValue("subAgentCode"),
		SubAgentName: r.FormValue("subAgentName"),
		GroupCode:    r.FormValue("groupCode"),
		GroupName:    r.FormValue("groupName"),
		IsNeedSign:   isNeedSign,
		MerStatus:    r.FormValue("merStatus"),
		Detail: model.MerDetail{
			MerName:       r.FormValue("merName"),
			AcctNum:       r.FormValue("acctNum"),
			GoodsTag:      r.FormValue("goodsTag"),
			CommodityName: r.FormValue("commodityName"),
		},
	}

	ret := Merchant.Find(merchant, pay, createStartTime, createEndTime, size, page)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func merchantFindOneHandle(w http.ResponseWriter, r *http.Request) {
	merId := r.FormValue("merId")
	ret := Merchant.FindOne(merId)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func merchantSaveHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := Merchant.Save(data)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func merchantUpdateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := Merchant.Update(data)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func merchantDeleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	merId := r.FormValue("merId")
	ret := Merchant.Delete(merId)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func routerSaveHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := RouterPolicy.Save(data)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func routerFindHandle(w http.ResponseWriter, r *http.Request) {
	merId := r.FormValue("merId")
	cardBrand := r.FormValue("cardBrand")
	chanCode := r.FormValue("chanCode")
	chanMerId := r.FormValue("chanMerId")
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	pay := r.FormValue("pay")
	ret := RouterPolicy.Find(merId, cardBrand, chanCode, chanMerId, pay, size, page)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func routerFindOneHandle(w http.ResponseWriter, r *http.Request) {
	merId := r.FormValue("merId")
	cardBrand := r.FormValue("cardBrand")
	ret := RouterPolicy.FindOne(merId, cardBrand)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func routerDeleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	merId := r.FormValue("merId")
	chanCode := r.FormValue("chanCode")
	cardBrand := r.FormValue("cardBrand")
	ret := RouterPolicy.Delete(merId, chanCode, cardBrand)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func channelMerchantFindHandle(w http.ResponseWriter, r *http.Request) {
	chanCode := r.FormValue("chanCode")
	chanMerId := r.FormValue("chanMerId")
	chanMerName := r.FormValue("chanMerName")
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	pay := r.FormValue("pay")
	ret := ChanMer.Find(chanCode, chanMerId, chanMerName, pay, size, page)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func channelMerchantMatchHandle(w http.ResponseWriter, r *http.Request) {
	chanCode := r.FormValue("chanCode")
	chanMerId := r.FormValue("chanMerId")
	chanMerName := r.FormValue("chanMerName")
	maxSize, _ := strconv.Atoi(r.FormValue("maxSize"))
	ret := ChanMer.Match(chanCode, chanMerId, chanMerName, maxSize)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func channelFindByMerIdAndCardBrandHandle(w http.ResponseWriter, r *http.Request) {

	merId := r.FormValue("merId")
	cardBrand := r.FormValue("cardBrand")
	ret := ChanMer.FindByMerIdAndCardBrand(merId, cardBrand)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func channelMerchantSaveHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := ChanMer.Save(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func channelMerchantDeleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	v := r.URL.Query()
	chanCode := v.Get("chanCode")
	chanMerId := v.Get("chanMerId")
	ret := ChanMer.Delete(chanCode, chanMerId)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func agentFindHandle(w http.ResponseWriter, r *http.Request) {
	agentCode := r.FormValue("agentCode")
	agentName := r.FormValue("agentName")
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	ret := Agent.Find(agentCode, agentName, size, page)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func agentDeleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	agentCode := r.FormValue("agentCode")
	ret := Agent.Delete(agentCode)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func agentSaveHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := Agent.Save(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func subAgentFindHandle(w http.ResponseWriter, r *http.Request) {
	agentCode := r.FormValue("agentCode")
	agentName := r.FormValue("agentName")
	subAgentCode := r.FormValue("subAgentCode")
	subAgentName := r.FormValue("subAgentName")
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	ret := SubAgent.Find(subAgentCode, subAgentName, agentCode, agentName, size, page)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func subAgentDeleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	subAgentCode := r.FormValue("subAgentCode")
	ret := SubAgent.Delete(subAgentCode)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func subAgentSaveHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := SubAgent.Save(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func groupFindHandle(w http.ResponseWriter, r *http.Request) {
	groupCode := r.FormValue("groupCode")
	groupName := r.FormValue("groupName")
	agentCode := r.FormValue("agentCode")
	agentName := r.FormValue("agentName")
	subAgentCode := r.FormValue("subAgentCode")
	subAgentName := r.FormValue("subAgentName")
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	ret := Group.Find(groupCode, groupName, agentCode, agentName, subAgentCode, subAgentName, size, page)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func groupDeleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	groupCode := r.FormValue("groupCode")
	ret := Group.Delete(groupCode)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func groupSaveHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := Group.Save(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func uptokenHandle(w http.ResponseWriter, r *http.Request) {
	qiniu.HandleUptoken(w, r)
}

func downURLHandle(w http.ResponseWriter, r *http.Request) {
	qiniu.HandleDownURL(w, r)
}

// qiniuDownloadHandle 用key换取七牛的私密下载链接
func qiniuDownloadHandle(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	url := qiniu.MakePrivateUrl(key)
	log.Debugf("redirect url is %s", url)
	var result *model.ResultBody
	if url == "" {
		result = &model.ResultBody{
			Status:  1,
			Message: "未找到下载链接，请确认您输入的key是否有误",
		}
	} else {
		result = &model.ResultBody{
			Status:  0,
			Message: url,
		}
	}

	rdata, err := json.Marshal(result)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func userCreateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := User.CreateUser(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}
func userFindHandle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	size, err := strconv.Atoi(params.Get("size"))
	if err != nil {
		http.Error(w, "参数 `size` 必须为整数", http.StatusBadRequest)
	}
	if size > 100 || size <= 0 {
		size = 20
	}
	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		http.Error(w, "参数 `page` 必须为整数", http.StatusBadRequest)
	}
	if page <= 0 {
		page = 1
	}

	cond := &model.User{
		UserName:     params.Get("userName"),
		NickName:     params.Get("nickName"),
		Mail:         params.Get("mail"),
		PhoneNum:     params.Get("phoneNum"),
		UserType:     params.Get("userRole"),
		AgentCode:    params.Get("agentCode"),
		SubAgentCode: params.Get("subAgentCode"),
		// AgentName: params.Get("agentName"),
		GroupCode: params.Get("groupCode"),
		// GroupName: params.Get("groupName"),
		MerId: params.Get("merId"),
		// MerName:   params.Get("merName"),
	}

	ret := User.Find(cond, size, page)

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}
	w.Write(retBytes)
}

// 登录操作���只允许get请���
func loginHandle(w http.ResponseWriter, r *http.Request) {
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

	user := &model.User{}
	err = json.Unmarshal(data, user)
	if err != nil {
		log.Errorf("json unmarshal error: %s", err)
		w.WriteHeader(501)
		return
	}
	log.Infof("user login,username=%s", user.UserName)

	// 密码解密
	pwd, err := rsaDecryptFromBrowser(user.Password)
	if err != nil {
		log.Errorf("escrypt password error %s", err)
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	log.Debugf("decrypted password is %s", pwd)

	ret := User.Login(user.UserName, pwd)

	if ret.Status == 0 {
		log.Debugf("create session begin")

		now := time.Now()
		cValue := util.SerialNumber()
		cExpires := now.Add(expiredTime)

		http.SetCookie(w, &http.Cookie{
			Name:     SessionKey,
			Value:    cValue,
			HttpOnly: true,
			Path:     "/master",
			Expires:  cExpires,
		})

		// 创建session
		session := &model.Session{
			SessionID:  cValue,
			User:       ret.Data.(*model.User),
			CreateTime: now,
			UpdateTime: now,
			Expires:    cExpires,
			Locale:     DefaultLocale,
		}

		ret = Session.Save(session)
		log.Debugf("create session end")
	}
	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}
	w.Write(retBytes)
	InsertMasterLog(r, user, []byte(""))
}

// 查找
func findSessionHandle(w http.ResponseWriter, r *http.Request) {
	sid, err := r.Cookie(SessionKey)
	if err != nil {
		log.Errorf("user not login: %s", err)
		http.Error(w, "用户未登录", http.StatusNotAcceptable)
		return
	}

	ret := Session.FindOne(sid.Value)

	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}
	w.Write(retBytes)
}

// 删除session。 登出操作，无论后台出什么错都返回成功。
func sessionDeleteHandle(w http.ResponseWriter, r *http.Request) {
	// 清除cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "QUICKMASTERID",
		Value:    "",
		HttpOnly: true,
		Path:     "/master",
		MaxAge:   -1,
	})

	sid, err := r.Cookie(SessionKey)
	if err != nil {
		log.Errorf("user not login when doing logout operation: %s", err)
		return
	}

	// 用于记录日志
	session, err := mongo.SessionColl.Find(sid.Value)
	if err != nil {
		log.Errorf("find session(%s) err: %s", sid.Value, err)
		return
	}

	ret := Session.Delete(sid.Value)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)

	HandleMasterLog(w, r, session.User)
}

func userUpdateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := User.UpdateUser(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// userUpdatePwdHandle 修改密码
func userUpdatePwdHandle(w http.ResponseWriter, r *http.Request) {
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
	ret := User.UpdatePwd(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// userDeleteHandle 删除用户
func userDeleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	userName := params.Get("userName")
	ret := User.RemoveUser(userName)
	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// userResetPwdHandle 重置密码
func userResetPwdHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// params := r.URL.Query()
	// userName := params.Get("userName")
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
	ret := User.ResetPwd(data, curSession.User)
	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// 商户导出
func merchantExportHandle(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	createTime := r.FormValue("createTime")
	createStartTime := ""
	createEndTime := ""
	if createTime != "" {
		createStartTime = createTime + " 00:00:00"
		createEndTime = createTime + " 23:59:59"
	}
	pay := r.FormValue("pay")
	isNeedSignStr := r.FormValue("isNeedSign")
	isNeedSign := false
	if isNeedSignStr == "true" {
		isNeedSign = true
	}
	merchant := model.Merchant{
		MerId:        r.FormValue("merId"),
		AgentCode:    r.FormValue("agentCode"),
		AgentName:    r.FormValue("agentName"),
		SubAgentCode: r.FormValue("subAgentCode"),
		SubAgentName: r.FormValue("subAgentName"),
		GroupCode:    r.FormValue("groupCode"),
		GroupName:    r.FormValue("groupName"),
		IsNeedSign:   isNeedSign,
		MerStatus:    r.FormValue("merStatus"),
		Detail: model.MerDetail{
			MerName:       r.FormValue("merName"),
			AcctNum:       r.FormValue("acctNum"),
			GoodsTag:      r.FormValue("goodsTag"),
			CommodityName: r.FormValue("commodityName"),
		},
	}

	Merchant.Export(w, merchant, pay, filename, createStartTime, createEndTime)
}

func agentUpdateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := Agent.Update(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func subAgentUpdateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := SubAgent.Update(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func groupUpdateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := Group.Update(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func channelMerchantUpdateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := ChanMer.Update(data)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func routerUpdateHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := RouterPolicy.Update(data)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// appResetPwdHandle 重置app用户密码
func appResetPwdHandle(w http.ResponseWriter, r *http.Request) {
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
	ret := AppUser.ResetPwd(data, curSession.User)
	rdata, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("mashal data error: %s", err)
		w.WriteHeader(501)
		w.Write([]byte("mashal data error"))
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

//邮件重置密码
func passwordResetHandle(w http.ResponseWriter, r *http.Request) {
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

	ret := User.PasswordReset(data)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}
