package goconf

import "testing"

func TestGetValue(t *testing.T) {
	mongoHost := GetValue("mongo", "db")
	expected := "angrycard"
	if mongoHost != expected {
		t.Errorf("read config value error: expect `%s`, but get `%s`", expected, mongoHost)
	}
}
