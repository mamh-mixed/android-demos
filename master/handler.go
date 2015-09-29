package master

import (
	// "bytes"
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

func tradeMsgHandle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	msgType, _ := strconv.Atoi(params.Get("msgType"))
	ret := getTradeMsg(&model.QueryCondition{
		MerId:    params.Get("merId"),
		OrderNum: params.Get("orderNum"),
		Busicd:   params.Get("busicd"),
	}, msgType)
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
		MerId:        merId,
		AgentCode:    params.Get("agentCode"),
		GroupCode:    params.Get("groupCode"),
		TransType:    transType,
		Respcd:       params.Get("respcd"),
		Busicd:       params.Get("busicd"),
		StartTime:    params.Get("startTime"),
		EndTime:      params.Get("endTime"),
		OrderNum:     params.Get("orderNum"),
		OrigOrderNum: params.Get("origOrderNum"),
		Col:          params.Get("pay"),
		BindingId:    params.Get("bindingId"),
		Size:         size,
		Page:         page,
	}

	transStatus := params.Get("transStatus")
	if transStatus != "" {
		cond.TransStatus = []string{transStatus}
	}

	ret := tradeQuery(cond)

	// // 允许跨域
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "*")

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
	params := r.URL.Query()
	filename := params.Get("filename")

	var merId = params.Get("merId")
	cond := &model.QueryCondition{
		MerId:        merId,
		Busicd:       params.Get("busicd"),
		StartTime:    params.Get("startTime"),
		EndTime:      params.Get("endTime"),
		OrderNum:     params.Get("orderNum"),
		OrigOrderNum: params.Get("origOrderNum"),
		Size:         maxReportRec,
		IsForReport:  true,
		Page:         1,
		RefundStatus: model.TransRefunded,
		TransStatus:  []string{model.TransSuccess},
	}

	tradeReport(w, cond, filename)
}

func tradeQueryStatsHandle(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))
	q := &model.QueryCondition{
		MerId:     r.FormValue("merId"),
		AgentCode: r.FormValue("agentCode"),
		Page:      page,
		Size:      size,
		MerName:   r.FormValue("merName"),
		StartTime: r.FormValue("startTime"),
		EndTime:   r.FormValue("endTime"),
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
	tradeQueryStatsReport(w, r)
}

func merchantFindHandle(w http.ResponseWriter, r *http.Request) {

	merId := r.FormValue("merId")
	merStatus := r.FormValue("merStatus")
	merName := r.FormValue("merName")
	groupCode := r.FormValue("groupCode")
	groupName := r.FormValue("groupName")
	agentCode := r.FormValue("agentCode")
	agentName := r.FormValue("agentName")
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	pay := r.FormValue("pay")
	ret := Merchant.Find(merId, merStatus, merName, groupCode, groupName, agentCode, agentName, pay, size, page)

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
	ret := SubAgent.Find(agentCode, agentName, subAgentCode, subAgentName, size, page)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func subAgentDeleteHandle(w http.ResponseWriter, r *http.Request) {

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
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	ret := Group.Find(groupCode, groupName, agentCode, agentName, size, page)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func groupDeleteHandle(w http.ResponseWriter, r *http.Request) {

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
	handleUptoken(w, r)
}

func downURLHandle(w http.ResponseWriter, r *http.Request) {
	handleDownURL(w, r)
}
