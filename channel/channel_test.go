package channel

import (
	"testing"

	"github.com/omigo/log"
)

func TestGetChan(t *testing.T) {

	c := GetChan("CFCA")
	if c == nil {
		t.Error("fail...")
		t.FailNow()
	}
	log.Debugf("%+v", c)
}
