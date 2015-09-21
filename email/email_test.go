package email

import (
	"testing"
)

func TestSend(t *testing.T) {
	e := Email{To: "379630413@qq.com", ActUrl: "https://www.baidu.com"}
	err := e.Send()
	if err != nil {
		t.Error(err)
	}
}
