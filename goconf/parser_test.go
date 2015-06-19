package goconf

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	// init()

	jsonBytes, err := json.MarshalIndent(Config, "", "\t")
	if err != nil {
		t.Errorf("marshal json error: %s", err)
		return
	}

	fmt.Printf("%s\n", string(jsonBytes))
}
