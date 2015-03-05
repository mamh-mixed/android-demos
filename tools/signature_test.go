package tools

import (
	"fmt"
	"testing"
)

func TestSignatureSha1WithRsa(t *testing.T) {

	t.Error("unimplement")
}

// TODO 每个方法都要有单元测试

func TestCheckChinaPaySignature(t *testing.T) {

	data := "{json:中文}"
	message, signature := ChinaPaySignature(data)

	pass := CheckChinaPaySignature(message, signature)

	fmt.Println(pass)
}
