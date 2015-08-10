package master

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	case "/master/trade/report":
		tradeReport(w, r)
		return
	case "/master/merchant/find":
		merId := r.FormValue("merId")
		merStatus := r.FormValue("merStatus")
		ret = Merchant.Find(merId, merStatus)
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
