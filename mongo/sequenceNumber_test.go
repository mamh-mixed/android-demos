package mongo

import (
	"testing"
)

func TestGetDaySN(t *testing.T) {
	// todo 并发测试
	t.Log(SnColl.GetDaySN())
}
