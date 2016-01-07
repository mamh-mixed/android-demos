package settle

import (
	"testing"
)

func TestGenerateTransFlow(t *testing.T) {
	var trans = transFlow{}
	trans.GenerateTransFlow("2016-01-06", "97491888")
}
