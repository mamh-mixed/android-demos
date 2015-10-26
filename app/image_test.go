package app

import (
	// "bytes"
	"image/jpeg"
	// "io/ioutil"
	"os"
	"testing"
)

func TestGenImageWithQrCode(t *testing.T) {

	text := ""
	shopName := "测试中文"
	newImage := genImageWithQrCode(payImageTemplate, text, shopName)
	f, _ := os.Create("/Users/zhiruichen/Desktop/bill1.jpg")
	jpeg.Encode(f, newImage, nil)
}
