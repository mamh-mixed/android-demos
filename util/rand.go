package util

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	gouuid "github.com/nu7hatch/gouuid"
	"github.com/CardInfoLink/log"
	"time"
)

const localDateTimeLayout = "0102150405" // MMDDHHMMSS

// SerialNumber 生成序列号，也就是UUID
func SerialNumber() string {
	u4, err := gouuid.NewV4()
	if err != nil {
		log.Errorf("error: %s", err)
		return ""
	}
	return fmt.Sprintf("%x", u4[:])
}

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_-"

// Nonce 生成指定长度位的随机数
func Nonce(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b&63]
	}
	return string(bytes)
}

// Confuse 混淆一个唯一码
func Confuse(uniqueId string) string {

	ub := []byte(uniqueId)
	length := len(ub)
	var sum = make([]byte, length+4)

	rand.Read(sum[0:2])
	for i, b := range ub {
		sum[i+2] = b
	}
	rand.Read(sum[length+2 : length+4])

	mb := md5.Sum(sum)
	return fmt.Sprintf("%x", mb[:])
}

// SignKey 随机生成32的密钥
func SignKey() string {
	var b = make([]byte, 20)
	rand.Read(b)
	mb := md5.Sum(b)
	return fmt.Sprintf("%x", mb[:])
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
// 格式 hh:mm:ss
func TimeToGiven(point string) (time.Duration, error) {
	layout := "2006-01-02 15:04:05"
	now := time.Now()
	given, err := time.ParseInLocation(layout, now.Format("2006-01-02")+" "+point, time.Local)
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
