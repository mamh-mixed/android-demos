package master

import (
	"encoding/csv"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"net/http"
	"strings"
)

const (
	totalFieldNum = 22
)

// importMerchant 接受csv格式文件，导入商户
func importMerchant(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("merchant")
	if err != nil {
		w.Write([]byte("文件上传错误，请重新上传:" + err.Error()))
		return
	}
	defer file.Close()

	// 判断是否csv文件
	filename := handler.Filename
	log.Debugf("file name : %s", filename)
	if !strings.HasSuffix(filename, "csv") {
		w.Write([]byte("文件类型错误，请上传csv文件"))
		return
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comment = '\n'
	rawCSVdata, err := reader.ReadAll()

	// 顺序 [商户编号 商户名称 机构号 权限（空即默认全部开放） 是否开启验签 签名密钥 商户商品名称 支付宝商户号 支付宝商户名称 讯联跟支付宝费率 商户跟讯联费率(支付宝) 支付宝密钥 微信商户号 微信子商户号 微信商户名称 是否代理商模式 微信AppId 讯联跟微信费率 商户跟讯联费率(微信) 微信签名密钥 微信http证书 微信http密钥]
	for i, mer := range rawCSVdata {
		log.Debug(len(mer))
		// 跳过头行
		if i == 0 {
			continue
		}
		if len(mer) != totalFieldNum {
			m := &model.Merchant{}
			m.MerId = mer[0]
		}
	}

}
