package auth

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/channel/alipay"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"net/http"
	"time"
)

// AuthHandle TODO 授权后的H5需要设计
func AuthHandle(w http.ResponseWriter, r *http.Request) {

	appID := r.FormValue("app_id")
	code := r.FormValue("app_auth_code")

	log.Infof("user authorization, appID=%s, code=%s", appID, code)

	cm, err := mongo.ChanMerColl.Find(channel.ChanCodeAlipay, appID) // TODO 待确定
	if err != nil {
		w.Write([]byte("invalid appID"))
		return
	}

	resp, err := alipay.GetAppAuthToken(appID, code, []byte(cm.PrivateKey))
	if err != nil {
		w.Write([]byte("get app_auth_token error: " + err.Error()))
		return
	}

	log.Infof("alipay auth return: %+v", resp)
	if resp.CommonBody.Code != "10000" {
		log.Infof("alipay auth return: %+v", resp)
		w.Write([]byte("get app_auth_token error: " + resp.CommonBody.Msg))
		return
	}

	// 如果找到，修改
	if err = mongo.ChanMerColl.Upsert(&model.ChanMer{
		ChanMerId:    resp.AuthAppID,
		ChanCode:     channel.ChanCodeAlipay,
		WxpAppId:     resp.AuthAppID,
		AuthToken:    resp.AppAuthToken,
		RefreshToken: resp.AppRefreshToken,
		AuthTime:     time.Now().Format("2006-01-02 15:04:05"),
		IsAgentMode:  true,
		Version:      channel.ALP2_0,
		AgentMer:     cm,
	}); err != nil {
		log.Errorf("alipay auth, upsert chanMer err: %s", err)
		w.Write([]byte("system error, please retry."))
		return
	}

	w.Write([]byte("success"))

}
