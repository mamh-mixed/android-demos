package mongo

import (
	"runtime"
	"testing"
)

// 并行测试sysSN
func TestGetSysSN(t *testing.T) {
	runtime.GOMAXPROCS(2)

	go loopSysSN(t)
	go loopSysSN(t)

	for i := 0; i < 2; i++ {
		<-quit
	}
}

// 并行测试
var quit chan int = make(chan int)

func loopSysSN(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(SnColl.GetSysSN())
	}

	quit <- 0
}

// GetDaySN 返回一个当天唯一的六位数字
func TestGetDaySN(t *testing.T) {
	s := SnColl.GetDaySN("M123", "T123")
	if s == "" {
		t.Error("TestGetDaySN error")
	}
	t.Log(s)

}

// 并行测试daySN
func TestGetDaySNConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(4)

	go loopDaySN(t)
	go loopDaySN(t)
	go loopDaySN(t)
	go loopDaySN(t)

	for i := 0; i < 4; i++ {
		<-quit1
	}
}

var quit1 chan int = make(chan int)

func loopDaySN(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(SnColl.GetDaySN("M123", "T123"))
	}

	quit1 <- 0
}
