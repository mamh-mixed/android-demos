package config

import "testing"

func TestGetValue(t *testing.T) {
	mongoHost := GetValue("mongo", "host")
	expected := "121.41.85.237"
	if mongoHost != expected {
		t.Errorf("read config value error: expect `%s`, but get `%s`", expected, mongoHost)
	}
}
