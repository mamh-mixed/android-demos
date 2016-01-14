package scanpay2

import (
	"testing"
)

func TestGetAuthToken(t *testing.T) {

	resp, err := GetAuthToken(&AuthTokenReq{
		CommonParams: CommonParams{
			AppID:      "2014122500021754",
			PrivateKey: LoadPrivateKey([]byte(privateKeyPem)),
		},
		GrantType: "authorization_code",
		Code:      "",
	})

}
