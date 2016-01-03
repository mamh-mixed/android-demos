package settle

import (
	"testing"
)

func TestGenerateTransFlow(t *testing.T) {
	var trans = transFlow{}
	trans.GenerateTransFlow("2015-12-08", "99911888")
}
