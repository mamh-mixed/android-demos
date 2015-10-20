package master

import (
	// "bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/omigo/log"

	"github.com/CardInfoLink/quickpay/util"
)

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
		Page:     page,
		Size:     size,
	}

	reqIds := params.Get("reqIds")
	if strings.Contains(reqIds, ",") {
		q.ReqIds = strings.Split(reqIds, ",")
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
		GroupCode:    params.Get("groupCode"),
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

	log.Debugf("GROUP CODE is %s", r.FormValue("groupCode"))
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
	ret := SubAgent.Find(subAgentCode, subAgentName, agentCode, agentName, size, page)
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
	qiniu.HandleUptoken(w, r)
}

func downURLHandle(w http.ResponseWriter, r *http.Request) {
	qiniu.HandleDownURL(w, r)
}

func userCreateHandle(w http.ResponseWriter, r *http.Request) {

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
		UserType:     params.Get("userType"),
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

func loginHandle(w http.ResponseWriter, r *http.Request) {

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
	ret := User.Login(user.UserName, user.Password)

	if ret.Status == 0 {
		log.Debugf("create session begin")

		now := time.Now()
		cValue := util.SerialNumber()
		cExpires := now.Add(expiredTime)

		http.SetCookie(w, &http.Cookie{
			Name:    "QUICKMASTERID",
			Value:   cValue,
			Path:    "/master",
			Expires: cExpires,
		})

		// 创建session
		session := &model.Session{
			SessionID:  cValue,
			User:       ret.Data.(*model.User),
			CreateTime: now,
			UpdateTime: now,
			Expires:    cExpires,
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

// 删除session
func sessionDeleteHandle(w http.ResponseWriter, r *http.Request) {
	sid, err := r.Cookie(SessionKey)
	if err != nil {
		log.Errorf("user not login: %s", err)
		http.Error(w, "用户未登录", http.StatusNotAcceptable)
		return
	}

	ret := Session.Delete(sid.Value)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func userUpdateHandle(w http.ResponseWriter, r *http.Request) {

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
func userUpdatePwdHandle(w http.ResponseWriter, r *http.Request) {

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
