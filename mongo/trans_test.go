package mongo

import (
	"github.com/omigo/g"
	"testing"
	"time"
)

func TestTransAdd(t *testing.T) {
	trans := Trans{
		ChanCode: "00010000",
		Time:     time.Now().Unix(),
		Flag:     0,
	}
	err := trans.Add()
	if err != nil {
		t.Errorf("add trans unsunccessful", err)
		t.FailNow()
	}
	g.Debug("add trans success %s", trans)
}

// func TestTransModify(t *testing.T) {

// }
