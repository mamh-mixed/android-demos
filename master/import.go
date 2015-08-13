package master

import (
	"archive/zip"
	"bytes"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

const (
	domain = "7xl02q.com1.z0.glb.clouddn.com"
)

var client = kodo.Client{}

func init() {
	conf.ACCESS_KEY = "-OOrgfZJbxz29kiW6HQsJ_OQJcjX6gaPRDf6xOcc"
	conf.SECRET_KEY = "rgBxbGeGJluv8ApEjY1RL2vq9IIfXcQAQqH4ttGo"
}

// importMerchant 接受excel格式文件，导入商户
func importMerchant(w http.ResponseWriter, r *http.Request) {

	// 调用七牛api获取刚上传的图片
	key := r.FormValue("key")
	baseUrl := kodo.MakeBaseUrl(domain, key)
	privateUrl := client.MakePrivateUrl(baseUrl, nil)

	resp, err := http.Get(privateUrl)
	if err != nil {
		log.Error(err)
		return
	}

	ebytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	// 判断内容类型
	contentType := resp.Header.Get("content-type")
	if contentType == "application/json" {
		log.Error(string(ebytes))
		return
	}

	// 包装成zipReader
	reader := bytes.NewReader(ebytes)
	zipReader, err := zip.NewReader(reader, int64(len(ebytes)))
	if err != nil {
		log.Error(err)
		return
	}

	// 转换成excel
	file, err := xlsx.ReadZipReader(zipReader)
	if err != nil {
		log.Error(err)
		return
	}

	// 读取数据
	for _, v := range file.Sheets {
		for _, r := range v.Rows {
			for _, c := range r.Cells {
				log.Debug(c.Value)
			}
		}
	}

	// TODO 校验数据、入库
}
