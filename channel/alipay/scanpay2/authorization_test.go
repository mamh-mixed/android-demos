package scanpay2

import (
	"testing"
)

func TestGetAuthToken(t *testing.T) {

	resp, err := GetAuthToken("2016010501065650", "wadawdawdad", []byte(privateKeyPem))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", resp)

}
