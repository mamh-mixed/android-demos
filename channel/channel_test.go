package channel

import (
	"github.com/omigo/log"
	"testing"
)

func TestGetChan(t *testing.T) {

	c := GetChan("CFCA")
	if c == nil {
		t.Error("fail...")
		t.FailNow()
	}
	log.Debugf("%+v", c)
}
