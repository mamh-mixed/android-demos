package app

import (
	// "bytes"
	"github.com/golang/freetype"
	qr "github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	// "io/ioutil"
	"os"
	"testing"
)

func TestGenImageWithQrCode(t *testing.T) {

	text := "112323131312313131"
	shopName := "测试中文"
	newImage := genImageWithQrCode(payImageTemplate, text, shopName)
	f, _ := os.Create("/Users/zhiruichen/Desktop/bill1.jpg")
	jpeg.Encode(f, newImage, nil)
}
