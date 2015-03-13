package tools

import (
	"testing"

	"github.com/omigo/g"
)

func TestUUID(t *testing.T) {
	uuid := serialNumber()

	g.Debug("uuid=%s", uuid)

	if uuid == "" {
		t.Error("unable generate uuid")
	}
}
