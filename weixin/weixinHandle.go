package weixin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/omigo/log"
)

const (
	authAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	authUserInfoURL    = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	refreshAATokenURL  = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	// 智慧微商会支付
	appID     = "wxdb1a8e71a25753a0"
	appSECRET = "fc17ad6aaa0da0434d75f51a2c96c3f9"
)

// AuthAccessTokenResp 获取网页授权用户信息时用的access_token
type AuthAccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
	Errcode      int32  `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

// AuthUserInfoResp  网页授权用户信息响应
type AuthUserInfoResp struct {
	OpenId     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int8     `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
	Errcode    int32    `json:"errcode"`
	Errmsg     string   `json:"errmsg"`
}

// RefreshAATokenResp 刷新AuthAccessToken
type RefreshAATokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Errcode      int32  `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

// GetAuthAccessToken 获取AuthAccessToken
func GetAuthAccessToken(code string) (authAccessTokenResp *AuthAccessTokenResp, err error) {
	authAccessTokenURLT := fmt.Sprintf(authAccessTokenURL, appID, appSECRET, code)
	resp, err := http.Get(authAccessTokenURLT)
	if err != nil {
		log.Errorf("http.Get authAccessToken err,%s", err)
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	log.Debugf("authAccessTokenResp:%s", string(bs))
	if err != nil {
		log.Errorf("read body of AuthAccessToken err,%s", err)
		return nil, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal(bs, &authAccessTokenResp)
	if err != nil {
		log.Errorf("json unmarshal AuthAccessToken err,%s", err)
		return nil, err
	}
	return authAccessTokenResp, nil
}

// GetAuthUserInfo 获取网页授权用户信息
func GetAuthUserInfo(authAccessToken, openId string) (authUserInfoResp *AuthUserInfoResp, err error) {
	authUserInfoURLT := fmt.Sprintf(authUserInfoURL, authAccessToken, openId)
	resp, err := http.Get(authUserInfoURLT)
	if err != nil {
		log.Errorf("http.Get authUserInfo err,%s", err)
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	log.Debugf("authUserInfoResp:", string(bs))
	if err != nil {
		log.Errorf("read body of resp authUserInfo err,%s", err)
		return nil, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal(bs, &authUserInfoResp)
	if err != nil {
		log.Errorf("json unmarshal authUserInfo err,%s", err)
		return nil, err
	}
	return authUserInfoResp, err
}

// RefreshAuthAccessToken 刷新authAccessToken  refreshToken:通过access_token获取到的refresh_token参数
func RefreshAuthAccessToken(appid, refreshToken string) (refreshAATokenResp *RefreshAATokenResp, err error) {
	refreshAATokenURLT := fmt.Sprintf(refreshAATokenURL, appid, refreshToken)
	resp, err := http.Get(refreshAATokenURLT)
	if err != nil {
		log.Errorf("http.Get refreshAAToken err,%s", err)
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	log.Debugf("refreshAATokenResp:%s", string(bs))
	if err != nil {
		log.Errorf("read body of refreshAATokenResp err,%s", err)
		return nil, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal(bs, &refreshAATokenResp)
	if err != nil {
		log.Errorf("json unmarshal refreshAATokenResp err,%s", err)
		return nil, err
	}
	return refreshAATokenResp, err
}
