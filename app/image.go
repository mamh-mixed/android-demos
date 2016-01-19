package app

import (
	"github.com/CardInfoLink/quickpay/util"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/CardInfoLink/log"

	qr "github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
)

var huawenFont *truetype.Font
var (
	payImageTemplate  image.Image
	billImageTemplate image.Image
)

func init() {
	ff := util.WorkDir + "/app/material/huawen.ttf"
	pp := util.WorkDir + "/app/material/pay.jpg"
	bp := util.WorkDir + "/app/material/bill.jpg"

	huawenFont = getFont(ff)
	payImageTemplate = getImage(pp)
	billImageTemplate = getImage(bp)
}

func getFont(path string) *truetype.Font {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("fail to load font: %s", err)
		os.Exit(2)
	}
	fb, err := ioutil.ReadAll(f)
	if err != nil {
		log.Errorf("fail to read font: %s", err)
		os.Exit(2)
	}
	font, err := freetype.ParseFont(fb)
	if err != nil {
		log.Errorf("fail to parse font: %s, use default", err)
		os.Exit(2)
	}
	return font
}

func getImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("fail to load image: %s, path=%s", err, path)
		os.Exit(2)
	}
	jpgImage, err := jpeg.Decode(f)
	if err != nil {
		log.Errorf("fail to decode jpg image: %s", err)
		os.Exit(2)
	}
	return jpgImage
}

// genImageWithQrCode 根据内容、店名生成图片
func genImageWithQrCode(template image.Image, text, shopName string) image.Image {
	// 准备空的图片
	dst := image.NewNRGBA(template.Bounds())

	// 截取logo图片
	logo := image.NewNRGBA(image.Rect(0, 0, 245, 245))
	draw.Draw(logo, logo.Bounds(), template, image.Pt(437+350, 420+350), draw.Src)

	// 将image对象包装成draw.Image
	draw.Draw(dst, dst.Bounds(), template, image.ZP, draw.Src)

	// 画最下方的门店名
	c := freetype.NewContext()
	c.SetFont(huawenFont)
	c.SetFontSize(20)
	c.SetDPI(144)
	c.SetClip(image.Rect(0, 2350, 1825, 2500))
	c.SetDst(dst)
	c.SetSrc(image.White)

	// 计算居中，一个字节大概13个像素点
	length := len(shopName) * 13
	sp := (1825 - length) / 2

	_, err := c.DrawString(shopName, freetype.Pt(sp, 2400))
	if err != nil {
		log.Errorf("fail to draw text on image: %s", err)
	}

	// 生成二维码
	qrCode, err := qr.New(text, qr.Highest) // 高容错等级
	if err != nil {
		log.Errorf("fail to generate qrCode image: %s", err)
		return dst
	}
	codeImage := qrCode.Image(1000) // 1000像素

	draw.Draw(dst, image.Rect(437-36, 420-36, 1387-36, 1370-36), codeImage, image.ZP, draw.Src)
	draw.Draw(dst, image.Rect(437+355, 420+355, 437+355+240, 420+355+240), logo, image.ZP, draw.Src)

	return dst
}
