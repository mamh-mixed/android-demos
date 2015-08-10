package master

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestImportFromCsv(t *testing.T) {

	buf := &bytes.Buffer{}
	fw := multipart.NewWriter(buf)
	w, err := fw.CreateFormFile("merchant", "merchant.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	file, err := os.Open("/Users/zhiruichen/Desktop/respCode_scanpay.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = io.Copy(w, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contentType := fw.FormDataContentType()
	t.Log("contentType: " + contentType)
	fw.Close()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { importMerchant(w, r) }))

	resp, err := http.Post(ts.URL, contentType, buf)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(string(data))
}
