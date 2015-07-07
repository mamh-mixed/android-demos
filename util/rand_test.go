package util

import (
	"fmt"
	"testing"

	"github.com/omigo/log"
)

func TestUUID(t *testing.T) {
	uuid := SerialNumber()

	log.Debugf("uuid=%s", uuid)

	if uuid == "" {
		t.Error("unable generate uuid")
	}
}

func TestMillisecond(t *testing.T) {
	t.Logf("Millisecond: %s", Millisecond())
}

func TestTimeToGiven(t *testing.T) {

	time, err := TimeToGiven("08:35:00")
	if err != nil {
		t.Errorf("fail to get time %s", err)
	}
	log.Debugf("time to given (%d)", time)
}

func TestNonce(t *testing.T) {
	fmt.Println(Nonce(16))
	fmt.Println(Nonce(32))
	fmt.Println(Nonce(64))
	fmt.Println(Nonce(128))
	fmt.Println(Nonce(256))
}
