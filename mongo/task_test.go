package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestAddTask(t *testing.T) {

	err := TaskCol.Add("test", false)
	if err != nil {
		t.Error(err)
	}
}

func TestPopTask(t *testing.T) {

	err := TaskCol.Pop("test")
	if err != nil {
		t.Error(err)
	}
}
