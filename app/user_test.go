package app

import (
	"testing"
)

func TestGetQiniuToken(t *testing.T) {
	result := User.getQiniuToken(&reqParams{
		UserName: "fnghwsj@qq.com",
		Password: "83d90a0f21db74e4cb78d6f2cbccb387",
	})

	t.Logf("result is %+v", result)
}
