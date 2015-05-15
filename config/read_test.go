package config

import "testing"

func TestGetValue(t *testing.T) {
	mongoHost := GetValue("mongo", "db")
	expected := "quickpay"
	if mongoHost != expected {
		t.Errorf("read config value error: expect `%s`, but get `%s`", expected, mongoHost)
	}
}
