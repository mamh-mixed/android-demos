package master

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"io/ioutil"
	"net/http"
	"strconv"

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
	case "/master/trade/report":
		tradeReport(w, r)
		return
	case "/master/trade/query":
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
		ret = tradeQueryStatistics(q)
	case "/master/merchant/find":
		merId := r.FormValue("merId")
		merStatus := r.FormValue("merStatus")
		page, _ := strconv.Atoi(r.FormValue("page"))
		size, _ := strconv.Atoi(r.FormValue("size"))
		ret = Merchant.Find(merId, merStatus, size, page)
	case "/master/merchant/save":
		ret = Merchant.Save(data)
	case "/master/router/save":
		ret = RouterPolicy.Save(data)
	case "/master/router/find":
		merId := r.FormValue("merId")
		ret = RouterPolicy.Find(merId)
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
		ret = ChanMer.Find(chanCode, chanMerId, chanMerName)
	case "/master/channelMerchant/findByMerIdAndCardBrand":
		merId := r.FormValue("merId")
		cardBrand := r.FormValue("cardBrand")
		ret = ChanMer.FindByMerIdAndCardBrand(merId, cardBrand)
	case "/master/channelMerchant/save":
		ret = ChanMer.Save(data)
	case "/master/agent/find":
		agentCode := r.FormValue("agentCode")
		agentName := r.FormValue("agentName")
		ret = Agent.Find(agentCode, agentName)
	case "/master/agent/delete":
		agentCode := r.FormValue("agentCode")
		agentName := r.FormValue("agentName")
		ret = Agent.Delete(agentCode, agentName)
	case "/master/agent/save":
		ret = Agent.Save(data)
	case "/master/group/find":
		groupCode := r.FormValue("groupCode")
		groupName := r.FormValue("groupName")
		ret = Group.Find(groupCode, groupName)
	case "/master/group/delete":
		groupCode := r.FormValue("groupCode")
		groupName := r.FormValue("groupName")
		ret = Group.Delete(groupCode, groupName)
	case "/master/group/save":
		ret = Group.Save(data)

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
