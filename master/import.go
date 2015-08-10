package master

import (
	"encoding/csv"
	"github.com/omigo/log"
	"net/http"
	"strings"
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

	for _, mer := range rawCSVdata {
		log.Debug(mer)
	}

}
