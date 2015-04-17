package mongo

import (
	"runtime"
	"testing"
)

// GetDaySN 返回一个当天唯一的六位数字
func TestGetDaySN(t *testing.T) {
	s := DaySNColl.GetDaySNA("M129", "T129")
	if s == "" {
		t.Error("TestGetDaySN error")
	}
	t.Log(s)

}

func TestGetDaySNConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(4)

	go loop1(t)
	go loop1(t)
	go loop1(t)
	go loop1(t)

	for i := 0; i < 4; i++ {
		<-quit1
	}
}

var quit1 chan int = make(chan int)

func loop1(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(DaySNColl.GetDaySNA("M129", "T129"))
	}

	quit1 <- 0
}
