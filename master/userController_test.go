package master

import (
	"testing"
)

func TestIsPasswordOk(t *testing.T) {
	p := "un#1016"
	result := isPasswordOk(p)

	t.Log(result)
}
