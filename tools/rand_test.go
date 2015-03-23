package tools

import (
	"testing"

	"github.com/omigo/g"
)

func TestUUID(t *testing.T) {
	uuid := SerialNumber()

	g.Debug("uuid=%s", uuid)

	if uuid == "" {
		t.Error("unable generate uuid")
	}
}

func TestMillisecond(t *testing.T) {
	t.Log("Millisecond: %s", Millisecond())
}

func TestTimeToGiven(t *testing.T) {

	time, err := TimeToGiven("08:35:00")
	if err != nil {
		t.Errorf("fail to get time %s", err)
	}
	g.Debug("time to given (%d)", time)
}
