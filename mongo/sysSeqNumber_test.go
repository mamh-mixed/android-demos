package mongo

import (
	"runtime"
	"testing"
)

func TestGetSysSN(t *testing.T) {
	runtime.GOMAXPROCS(2)

	go loop(t)
	go loop(t)

	for i := 0; i < 2; i++ {
		<-quit
	}
}

// 并行测试
var quit chan int = make(chan int)

func loop(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(SnColl.GetSysSN())
	}

	quit <- 0
}
