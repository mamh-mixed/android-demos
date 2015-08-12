package master

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestImportFromCsv(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { importMerchant(w, r) }))

	params := url.Values{}
	params.Add("key", "tple.xlsx")
	resp, err := http.PostForm(ts.URL, params)
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
