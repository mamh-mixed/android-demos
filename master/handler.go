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
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/omigo/log"

	"github.com/CardInfoLink/quickpay/util"
)

// tradeSettleReportHandle 清算报表查询的
func tradeSettleReportHandle(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")
	date := r.FormValue("date")
	size, _ := strconv.Atoi(r.FormValue("size"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	ret := tradeSettleReportQuery(role, date, size, page)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("marshal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
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

	log.Debugf("col is %s", params.Get("pay"))

	cond := &model.QueryCondition{
		MerId:          merId,
		AgentCode:      params.Get("agentCode"),
		SubAgentCode:   params.Get("subAgentCode"),
		GroupCode:      params.Get("groupCode"),
		TransType:      transType,
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
		Size:           size,
		Page:           page,
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
	params := r.URL.Query()
	filename := params.Get("filename")

	var merId = params.Get("merId")
	cond := &model.QueryCondition{
		MerId:        merId,
		Busicd:       params.Get("busicd"),
		AgentCode:    params.Get("agentCode"),
		SubAgentCode: params.Get("subAgentCode"),
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

	ret := Merchant.Find(merchant, pay, size, page)

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
	InsertMasterLog(r, user, data)
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

	http.SetCookie(w, &http.Cookie{
		Name:   "QUICKMASTERID",
		Value:  "",
		Path:   "/master",
		MaxAge: -1,
	})

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)

	HandleMasterLog(w, r, session.User)
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

// userUpdatePwdHandle 修改密码
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

// userDeleteHandle 删除用户
func userDeleteHandle(w http.ResponseWriter, r *http.Request) {
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
	params := r.URL.Query()
	userName := params.Get("userName")
	ret := User.ResetPwd(userName)
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
