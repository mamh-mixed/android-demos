package core

import "testing"

func TestIsUnionPayCard(t *testing.T) {
	// 银联卡
	cardNum := "6222020302062061908"
	if result := IsUnionPayCard(cardNum, "CUP"); !result {
		t.Errorf("Expected 'true',but get '%t'", result)
	}
	// MCC卡
	cardNum = "5500521234567891"
	if result := IsUnionPayCard(cardNum, ""); result {
		t.Errorf("Expected 'false',but get '%t'", result)
	}
}
