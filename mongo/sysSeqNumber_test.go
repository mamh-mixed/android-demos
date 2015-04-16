package mongo

import (
	"testing"
)

func TestGetSysSN(t *testing.T) {
	t.Log(SnColl.GetSysSN())
}
