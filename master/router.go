package master

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/CardInfoLink/quickpay/model"

	"github.com/omigo/log"
)

// MasterRoute 后台管理的请求统一入口
func MasterRoute(w http.ResponseWriter, r *http.Request) {
	log.Infof("url = %s", r.URL.String())

	var ret *model.ResultBody

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read all body error: %s", err)
		w.WriteHeader(501)
		return
	}

	switch r.URL.Path {
	case "/master/trade/query":
		tradeQuery(w, data)
		return
	case "/master/trade/report":
		tradeReport(w, r)
		return
	case "/master/trade/stat":
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
		ret = tradeQueryStats(q)
	case "/master/trade/stat/report":
		tradeQueryStatsReport(w, r)
		return
	case "/master/merchant/find":
		merId := r.FormValue("merId")
		merStatus := r.FormValue("merStatus")
		merName := r.FormValue("merName")
		groupCode := r.FormValue("groupCode")
		groupName := r.FormValue("groupName")
		agentCode := r.FormValue("agentCode")
		agentName := r.FormValue("agentName")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = Merchant.Find(merId, merStatus, merName, groupCode, groupName, agentCode, agentName, size, page)
	case "/master/merchant/one":
		merId := r.FormValue("merId")
		ret = Merchant.FindOne(merId)
	case "/master/merchant/save":
		ret = Merchant.Save(data)
	case "/master/router/save":
		ret = RouterPolicy.Save(data)
	case "/master/router/find":
		merId := r.FormValue("merId")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = RouterPolicy.Find(merId, size, page)
	case "/master/router/one":
		merId := r.FormValue("merId")
		cardBrand := r.FormValue("cardBrand")
		ret = RouterPolicy.FindOne(merId, cardBrand)
	case "/master/router/delete":
		merId := r.FormValue("merId")
		chanCode := r.FormValue("chanCode")
		cardBrand := r.FormValue("cardBrand")
		ret = RouterPolicy.Delete(merId, chanCode, cardBrand)
	case "/master/channelMerchant/find":
		chanCode := r.FormValue("chanCode")
		chanMerId := r.FormValue("chanMerId")
		chanMerName := r.FormValue("chanMerName")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = ChanMer.Find(chanCode, chanMerId, chanMerName, size, page)
	case "/master/channelMerchant/match":
		chanCode := r.FormValue("chanCode")
		chanMerId := r.FormValue("chanMerId")
		chanMerName := r.FormValue("chanMerName")
		maxSize, _ := strconv.Atoi(r.FormValue("maxSize"))
		ret = ChanMer.Match(chanCode, chanMerId, chanMerName, maxSize)
	case "/master/channelMerchant/findByMerIdAndCardBrand":
		merId := r.FormValue("merId")
		cardBrand := r.FormValue("cardBrand")
		ret = ChanMer.FindByMerIdAndCardBrand(merId, cardBrand)
	case "/master/channelMerchant/save":
		ret = ChanMer.Save(data)
	case "/master/agent/find":
		agentCode := r.FormValue("agentCode")
		agentName := r.FormValue("agentName")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = Agent.Find(agentCode, agentName, size, page)
	case "/master/agent/delete":
		agentCode := r.FormValue("agentCode")
		ret = Agent.Delete(agentCode)
	case "/master/agent/save":
		ret = Agent.Save(data)
	case "/master/group/find":
		groupCode := r.FormValue("groupCode")
		groupName := r.FormValue("groupName")
		agentCode := r.FormValue("agentCode")
		agentName := r.FormValue("agentName")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = Group.Find(groupCode, groupName, agentCode, agentName, size, page)
	case "/master/group/delete":
		groupCode := r.FormValue("groupCode")
		ret = Group.Delete(groupCode)
	case "/master/group/save":
		ret = Group.Save(data)
	case "/master/user/login":
		userName := r.FormValue("userName")
		password := r.FormValue("password")
		ret = User.Login(userName, password)
	case "/master/user/create":
		ret = User.CreateUser(data)
	case "/master/user/update":
		ret = User.UpdateUser(data)
	case "/master/user/find":
		userName := r.FormValue("userName")
		nickName := r.FormValue("nickName")
		roleName := r.FormValue("roleName")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = User.Find(userName, nickName, roleName, size, page)
	case "/master/user/remove":
		userName := r.FormValue("userName")
		ret = User.RemoveUser(userName)
	case "/master/menu/create":
		ret = Menu.CreateMenu(data)
	case "/master/menu/update":
		ret = Menu.UpdateMenu(data)
	case "/master/menu/find":
		nameCN := r.FormValue("nameCN")
		route := r.FormValue("route")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = Menu.Find(nameCN, route, size, page)
	case "/master/menu/remove":
		route := r.FormValue("route")
		ret = Menu.RemoveMenu(route)
	case "/master/role/create":
		ret = Role.CreateRole(data)
	case "/master/role/update":
		ret = Role.UpdateRole(data)
	case "/master/role/find":
		roleID := r.FormValue("roleID")
		name := r.FormValue("name")
		size, _ := strconv.Atoi(r.FormValue("size"))
		page, _ := strconv.Atoi(r.FormValue("page"))
		ret = Role.Find(roleID, name, size, page)
	case "/master/role/remove":
		roleID := r.FormValue("roleID")
		ret = Role.RemoveRole(roleID)
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
