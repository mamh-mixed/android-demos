package core

import (
	"testing"

	"github.com/CardInfoLink/log"
)

func TestCardBinMatch(t *testing.T) {

	s := tree.match("6222801932062061908")
	if s != "622280193" {
		t.Errorf("expect cardBin : 622280193,but get : %s", s)
		t.FailNow()
	}
	log.Debugf("%+s , %+v\n", s, tree.root)
}

func TestFindCardBin(t *testing.T) {

	cardBin, err := findCardBin("6222801932062061908")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	log.Debugf("%+v", cardBin)
}
