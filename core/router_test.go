package core

import (
	"github.com/omigo/g"
	"testing"
)

func TestFindCardBrandByCardNum(t *testing.T) {
	cardNum := "38221234098712"
	cardBrand := doCardBrandTest(cardNum, t)
	if cardBrand != "DNC" {
		t.Errorf("Expected '%s',but get '%s'", "DNC", cardBrand)
	}

	cardNum = "6222020302062061908"
	cardBrand = doCardBrandTest(cardNum, t)
	if cardBrand != "CUP" {
		t.Errorf("Expected '%s',but get '%s'", "CUP", cardBrand)
	}

	cardNum = "400052123456789123"
	cardBrand = doCardBrandTest(cardNum, t)
	if cardBrand != "VIS" {
		t.Errorf("Expected '%s',but get '%s'", "VIS", cardBrand)
	}

	cardNum = "5500521234567891"
	cardBrand = doCardBrandTest(cardNum, t)
	if cardBrand != "MCC" {
		t.Errorf("Expected '%s',but get '%s'", "MCC", cardBrand)
	}

}

func doCardBrandTest(cardNum string, t *testing.T) (cardBrand string) {
	g.Debug("cardNum is: %s;cardNum length: %d", cardNum, len(cardNum))
	cardBrand = FindCardBrandByCardNum(cardNum)
	g.Debug("cardBrand is: %s", cardBrand)
	return cardBrand
}

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
