package channel

import (
	"github.com/omigo/log"
	"testing"
)

func TestGetChan(t *testing.T) {

	c := GetChan("cfca")
	if c == nil {
		t.Error("fail...")
		t.FailNow()
	}
	log.Debugf("%+v", c)
}
