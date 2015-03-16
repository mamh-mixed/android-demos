package tools

import (
	"fmt"
	u "github.com/nu7hatch/gouuid"
	"github.com/omigo/g"
	"time"
)

// SerialNumber 生成序列号，也就是UUID
func SerialNumber() string {
	u4, err := u.NewV4()
	if err != nil {
		g.Error("error: ", err)
		return ""
	}
	return fmt.Sprintf("%x", u4[:])
}

// Millisecond 获取新世纪以来到目前为止的毫秒数
func Millisecond() string {
	t0 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	d := time.Since(t0)
	return fmt.Sprintf("%d", int64(d.Nanoseconds()/1000000))
}
