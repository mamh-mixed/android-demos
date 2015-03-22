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
	return fmt.Sprintf("%d", int64(time.Now().UnixNano()/1000000))
}

// NextDay 获得第二天时间
func NextDay(today string) string {

	layout := "2006-01-02"
	to, err := time.Parse(layout, today)
	if err != nil {
		g.Error("fail to parese (%s : %s) ", today, err)
	}
	// add a day
	next := to.Add(time.Duration(24) * time.Hour)
	return next.Format(layout)
}
