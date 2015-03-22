package tools

import (
	"fmt"
	u "github.com/nu7hatch/gouuid"
	"github.com/omigo/g"
	"strings"
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

// TimeToGiven 获得当前时间到point的秒数
// 格式hh:mm:ss
func TimeToGiven(point string) (int64, error) {
	layout := "2006-01-02 15:04:05"
	//当前时间
	current := time.Now()
	//TODO current不能直接跟given比较
	now, err := time.Parse(layout, current.Format(layout))
	value := strings.Split(now.Format(layout), " ")[0] + " " + point
	given, err := time.Parse(layout, value)
	if err != nil {
		return 0, err
	}
	//在当前时间之后
	if given.After(now) {
		return given.Unix() - now.Unix(), nil
	}
	//在当前时间之前，应该加一天
	given = given.Add(time.Duration(24) * time.Hour)
	return given.Unix() - now.Unix(), nil
}
