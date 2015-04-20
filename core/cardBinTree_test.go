package core

import (
	"github.com/omigo/log"
	"testing"
)

func TestCardBinMatch(t *testing.T) {

	s := tree.match("6222801932062061908")
	if s != "622280193" {
		t.Errorf("expect cardBin : 622280193,but get : %s", s)
		t.FailNow()
	}
	log.Debugf("%+s , %+v\n", s, tree.root)
}
