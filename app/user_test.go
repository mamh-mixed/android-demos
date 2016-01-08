package app

import (
	"testing"
	"time"
)

func TestFindRefundOrdersOfOrder(t *testing.T) {
	req := &reqParams{
		UserName:  "200000000010002",
		Password:  "670b14728ad9902aecba32e22fa4f6bd",
		OrderNum:  "15120318035173616",
		Transtime: time.Now().Format("20060102150405"),
	}

	result := User.FindRefundOrdersOfOrder(req)

	t.Logf("reesult is %+v", result)
}
