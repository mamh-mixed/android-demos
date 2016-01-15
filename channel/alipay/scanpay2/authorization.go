// 授权类接口
package scanpay2

import (
	"encoding/json"
	"github.com/CardInfoLink/quickpay/model"
	"net/url"
)

// GetAuthToken oauth2认证
func GetAuthToken(appID, code string, pem []byte) (*AuthTokenResp, error) {
	req := &AuthTokenReq{
		CommonParams: CommonParams{
			AppID:      appID,
			PrivateKey: LoadPrivateKey(pem),
			Req:        &model.ScanPayRequest{},
		},
		GrantType: "authorization_code",
		Code:      code,
	}

	resp := &AuthTokenResp{}
	return resp, Execute(req, resp)
}

// RefreshAuthToken 刷新交易令牌
func RefreshAuthToken(appID, refreshToken string, pem []byte) (*AuthTokenResp, error) {
	req := &AuthTokenReq{
		CommonParams: CommonParams{
			AppID:      appID,
			PrivateKey: LoadPrivateKey(pem),
			Req:        &model.ScanPayRequest{},
		},
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	resp := &AuthTokenResp{}
	return resp, Execute(req, resp)
}

type AuthTokenReq struct {
	CommonParams
	GrantType    string `json:"grant_type"` // authorization_code/refresh_token
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AuthTokenResp struct {
	CommonBody
	Raw             json.RawMessage `json:"alipay_open_auth_token_app_response"` // 返回消息体
	AppAuthToken    string          `json:"app_auth_token,omitempty"`
	UserID          string          `json:"user_id,omitempty"`
	AuthAppID       string          `json:"auth_app_id,omitempty"`
	ExpiresIn       int             `json:"expires_in,omitempty"`
	ReExpiresIn     int             `json:"re_expires_in,omitempty"`
	AppRefreshToken string          `json:"app_refresh_token,omitempty"`
}

// Values 组装公共参数
func (c *AuthTokenReq) Values() (v url.Values) {
	c.CommonParams.Method = "alipay.open.auth.token.app"
	return c.CommonParams.Values()
}

// SaveLog 重写该方法，不记录日志
func (c *AuthTokenReq) SaveLog() bool {
	return false
}

// GetRaw 报文内容
func (c *AuthTokenResp) GetRaw() []byte {
	return []byte(c.Raw)
}
