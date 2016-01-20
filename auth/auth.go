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

	// 查询存在的大商户
	agentCm, err := mongo.ChanMerColl.Find(channel.ChanCodeAlipay, appID) // TODO 待确定
	if err != nil {
		w.Write([]byte("非法的APPID参数。"))
		return
	}

	resp, err := alipay.GetAppAuthToken(appID, code, []byte(agentCm.PrivateKey))
	if err != nil {
		w.Write([]byte("获取授权信息失败，请重新授权。详情: " + err.Error()))
		return
	}

	log.Infof("alipay auth return: %+v", resp)
	if resp.CommonBody.Code != "10000" {
		log.Infof("alipay auth return error: %+v", resp)
		w.Write([]byte("获取授权信息失败，请重新授权。详情: " + resp.CommonBody.Msg))
		return
	}

	// 如果找到，修改
	cm, err := mongo.ChanMerColl.Find(channel.ChanCodeAlipay, resp.UserID)
	if err != nil {
		// 没找到，保存个新的，直接走2.0接口
		if err = mongo.ChanMerColl.Add(&model.ChanMer{
			ChanMerId:    resp.UserID,
			ChanCode:     channel.ChanCodeAlipay,
			WxpAppId:     resp.AuthAppID,
			AuthToken:    resp.AppAuthToken,
			RefreshToken: resp.AppRefreshToken,
			AuthTime:     time.Now().Format("2006-01-02 15:04:05"),
			IsAgentMode:  true,
			Version:      channel.ALP2_0,
			AgentMer:     agentCm,
		}); err != nil {
			w.Write([]byte("系统错误，请重新授权。"))
			return
		}

		w.Write([]byte("授权成功"))
		return
	}

	// 找到
	cm.WxpAppId = resp.AuthAppID
	cm.RefreshToken = resp.AppRefreshToken
	cm.AuthToken = resp.AppAuthToken
	cm.AuthTime = time.Now().Format("2006-01-02 15:04:05")
	cm.AgentMer = agentCm // 先关联上授权支付宝商户，等确认走2.0时，IsAgentMode=true，并且version=ALP2_0即可切换
	if err = mongo.ChanMerColl.Update(cm); err != nil {
		w.Write([]byte("系统错误，请重新授权。"))
		return
	}

	w.Write([]byte("授权成功"))

}
