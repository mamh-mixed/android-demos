package tools

import (
	"fmt"
	"strings"
	"time"

	u "github.com/nu7hatch/gouuid"
	"github.com/omigo/log"
)

const localDateTimeLayout = "0102150405" // MMDDHHMMSS

// SerialNumber 生成序列号，也就是UUID
func SerialNumber() string {
	u4, err := u.NewV4()
	if err != nil {
		log.Errorf("error: %s", err)
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
		log.Errorf("fail to parese (%s : %s) ", today, err)
	}
	// add a day
	next := to.Add(time.Duration(24) * time.Hour)
	return next.Format(layout)
}

// TimeToGiven 获得当前时间到point的秒数
// 格式hh:mm:ss
func TimeToGiven(point string) (time.Duration, error) {
	layout := "2006-01-02 15:04:05"
	//当前时间
	nowStr := time.Now().Format(layout)
	//TODO current不能直接跟given比较
	now, err := time.Parse(layout, nowStr)
	value := strings.Split(now.Format(layout), " ")[0] + " " + point
	given, err := time.Parse(layout, value)
	if err != nil {
		return 0, err
	}
	//在当前时间之后
	if given.After(now) {
		return time.Duration(given.Unix() - now.Unix()), nil
	}
	//在当前时间之前，应该加一天
	given = given.Add(time.Duration(24) * time.Hour)
	return time.Duration(given.Unix() - now.Unix()), nil
}
