package master

import (
	"encoding/json"
	"testing"
)

func TestRespCodeFindOne(t *testing.T) {
	code := "000002"
	result := RespCode.FindOne(code)
	data, _ := json.Marshal(result)
	t.Logf("result is %s", string(data))
}
