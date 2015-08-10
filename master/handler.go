package master

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
)

const (
	NormalMerStatus = "Normal"
)

// QuickMaster 后台管理的请求统一入口
func QuickMaster(w http.ResponseWriter, r *http.Request) {
	log.Infof("url = %s", r.URL.String())

	var ret *model.ResultBody

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read all body error: %s", err)
		w.WriteHeader(501)
		return
	}

	switch r.URL.Path {
	case "/quickMaster/merchant/all":
		ret = AllMerchant(data)
	case "/quickMaster/merchant/add":
		ret = AddMerchant(data)
	case "/quickMaster/channelMerchant/all":
		ret = AllChannelMerchant(data)
	case "/quickMaster/channelMerchant/add":
		ret = AddChannelMerchant(data)
	case "/quickMaster/router/save":
		ret = AddRouter(data)
	case "/quickMaster/router/find":
		merId := r.FormValue("merId")
		ret = AllRouterOfOneMerchant(merId)
	case "/quickMaster/agent/all":
		ret = AllAgent(data)
	case "/quickMaster/agent/add":
		ret = AddAgent(data)
	default:
		w.WriteHeader(404)
	}

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

// AllMerchant 处理查找所有商户的请求
func AllMerchant(data []byte) (result *model.ResultBody) {
	cond := new(model.Merchant)
	err := json.Unmarshal(data, cond)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	merchants, err := mongo.MerchantColl.FindAllMerchant(cond)

	if err != nil {
		log.Errorf("查询所有商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    merchants,
	}

	return
}

// AddMerchant 处理新增一个商户的请求
func AddMerchant(data []byte) (result *model.ResultBody) {
	m := new(model.Merchant)
	err := json.Unmarshal(data, m)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if m.MerId == "" {
		log.Error("没有MerId")
		return model.NewResultBody(3, "缺失必要元素merId")
	}

	if m.MerStatus == "" {
		m.MerStatus = NormalMerStatus
	}

	err = mongo.MerchantColl.Insert(m)
	if err != nil {
		log.Errorf("新增商户失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    m,
	}

	return
}

// AllChannelMerchant 处理查找所有商户的请求
func AllChannelMerchant(data []byte) (result *model.ResultBody) {
	cond := new(model.ChanMer)
	err := json.Unmarshal(data, cond)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	merchants, err := mongo.ChanMerColl.FindByCondition(cond)

	if err != nil {
		log.Errorf("查询所有商户出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    merchants,
	}

	return
}

// AddChannelMerchant 处理新增一个渠道商户的请求
func AddChannelMerchant(data []byte) (result *model.ResultBody) {
	m := new(model.ChanMer)
	err := json.Unmarshal(data, m)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if m.ChanCode == "" {
		log.Error("没有ChanCode")
		return model.NewResultBody(3, "缺失必要元素chanCode")
	}

	err = mongo.ChanMerColl.Add(m)
	if err != nil {
		log.Errorf("新增渠道商户失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    m,
	}

	return
}

// AddRouter 处理新增一个路由的请求
func AddRouter(data []byte) (result *model.ResultBody) {
	r := new(model.RouterPolicy)
	err := json.Unmarshal(data, r)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}
	if r.MerId == "" {
		log.Error("MerId")
		return model.NewResultBody(3, "缺失必要元素 merId")
	}

	if r.ChanCode == "" {
		log.Error("没有 ChanCode")
		return model.NewResultBody(3, "缺失必要元素 chanCode")
	}

	if r.ChanMerId == "" {
		log.Error("没有 ChanMerId")
		return model.NewResultBody(3, "缺失必要元素 chanMerId")
	}

	if r.CardBrand == "" {
		log.Error("没有 CardBrand")
		return model.NewResultBody(3, "缺失必要元素 cardBrand")
	}

	err = mongo.RouterPolicyColl.Insert(r)
	if err != nil {
		log.Errorf("保存路由信息失败:%s", err)
		return model.NewResultBody(1, "保存路由信息失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "保存成功",
		Data:    r,
	}

	return
}

// AllRouterOfOneMerchant 处理查找商户的所有路由的请求
func AllRouterOfOneMerchant(merId string) (result *model.ResultBody) {
	routers, err := mongo.RouterPolicyColl.FindAllOfOneMerchant(merId)

	if err != nil {
		log.Errorf("查询商户(%s)的所有路由失败: %s", merId, err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    routers,
	}

	return
}

// AllAgent 处理查找所有代理商的请求
func AllAgent(data []byte) (result *model.ResultBody) {
	cond := new(model.Agent)
	err := json.Unmarshal(data, cond)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	agents, err := mongo.AgentColl.FindByCondition(cond)

	if err != nil {
		log.Errorf("查询代理商出错:%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    agents,
	}

	return
}

// AddAgent 处理新增一个商户的请求
func AddAgent(data []byte) (result *model.ResultBody) {
	a := new(model.Agent)
	err := json.Unmarshal(data, a)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if a.AgentCode == "" {
		log.Error("没有AgentCode")
		return model.NewResultBody(3, "缺失必要元素AgentCode")
	}

	if a.AgentName == "" {
		log.Error("没有AgentName")
		return model.NewResultBody(3, "缺失必要元素AgentName")
	}

	err = mongo.AgentColl.Add(a)
	if err != nil {
		log.Errorf("新增代理商失败:%s", err)
		return model.NewResultBody(1, err.Error())
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "操作成功",
		Data:    a,
	}

	return
}
